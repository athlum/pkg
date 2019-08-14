package zk

import (
	"fmt"
	cmap "github.com/athlum/pkg/concurrentMap"
	"github.com/athlum/pkg/exitChan"
	"github.com/athlum/pkg/log"
	"github.com/samuel/go-zookeeper/zk"
	"path"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	NodeNew NodeEventType = iota
	NodeUpdate
	NodeRemoved
)

type NodeEventType int

type NodeEvent struct {
	Path  string
	Node  string
	Val   []byte
	Stat  *zk.Stat
	Event NodeEventType
}

type TreeNodeMap struct {
	*cmap.ConcurrentMap
}

func NewTreeNodeMap() *TreeNodeMap {
	return &TreeNodeMap{ConcurrentMap: cmap.New()}
}

func (tnm *TreeNodeMap) Get(key string) (*TreeNode, bool) {
	o, exists := tnm.ConcurrentMap.Get(key)
	if exists {
		return o.(*TreeNode), exists
	}
	return nil, exists
}

func (tnm *TreeNodeMap) SetIfAbsent(key string, value *TreeNode) bool {
	return tnm.ConcurrentMap.SetIfAbsent(key, value)
}

func (tnm *TreeNodeMap) Set(key string, value *TreeNode) {
	tnm.ConcurrentMap.Set(key, value)
}

func (tnm *TreeNodeMap) Has(key string) bool {
	return tnm.ConcurrentMap.Has(key)
}

func (tnm *TreeNodeMap) Remove(key string) {
	tnm.ConcurrentMap.Remove(key)
}

func (tnm *TreeNodeMap) Loop(f func(key string, node *TreeNode) bool) {
	for e := range tnm.ConcurrentMap.IterBuffered() {
		if breaked := f(e.Key, e.Val.(*TreeNode)); breaked {
			break
		}
	}
}

type NodeStop struct {
	childStop *exitChan.ExitChan
	stop      *exitChan.ExitChan
	flushStop *exitChan.ExitChan
}

func (ns *NodeStop) Stop() {
	ns.stop.Close()
	ns.childStop.Close()
	ns.flushStop.Close()
}

func NewNodeStop() *NodeStop {
	return &NodeStop{
		childStop: exitChan.NewExitChan(),
		stop:      exitChan.NewExitChan(),
		flushStop: exitChan.NewExitChan(),
	}
}

type syncVersion struct {
	version int32
	lock    *sync.Mutex
}

func newSyncVersion() *syncVersion {
	return &syncVersion{
		lock:    &sync.Mutex{},
		version: -1,
	}
}

func (sv *syncVersion) Update(v int32) bool {
	sv.lock.Lock()
	defer sv.lock.Unlock()

	if v > sv.version {
		sv.version = v
		return true
	}
	return false
}

type TreeNode struct {
	Root     *TreeNode
	Parent   *TreeNode
	Children *TreeNodeMap
	Event    chan *NodeEvent
	Path     string
	Node     string
	zk       *ZK
	version  *syncVersion
	cversion *syncVersion
	stop     *NodeStop
	interval time.Duration
	cleared  int32
}

func (conn *ZK) WatchNode(path, node string, parent *TreeNode, interval time.Duration) (*TreeNode, error) {
	return NewTreeNode(conn, path, node, parent, interval)
}

func NewTreeNode(conn *ZK, path, node string, parent *TreeNode, interval time.Duration) (*TreeNode, error) {
	tn := &TreeNode{
		Path:     path,
		Children: NewTreeNodeMap(),
		Event:    make(chan *NodeEvent),
		version:  newSyncVersion(),
		cversion: newSyncVersion(),
		stop:     NewNodeStop(),
		zk:       conn,
		interval: interval,
	}
	atomic.StoreInt32(&tn.cleared, 0)
	if len(node) > 0 {
		tn.Node = node
	} else {
		seps := strings.Split(tn.Path, "/")
		if len(seps) == 0 {
			return nil, fmt.Errorf("Invalid node path: %v", path)
		}
		tn.Node = seps[len(seps)-1]
	}
	if parent != nil {
		tn.Parent = parent
		if parent.Root != nil {
			tn.Root = parent.Root
		} else if parent.Parent == nil {
			tn.Root = parent
		} else {
			return nil, fmt.Errorf("Parent %v has no root but has parent", parent.Path)
		}
	}
	tn.zk.treeNodes.SetIfAbsent(tn.Path, tn)
	return tn, nil
}

func (tn *TreeNode) Init() error {
	if err := tn._flush(NodeNew); err != nil {
		return err
	}
	tn.ChildrenW()
	tn.GetW()
	if tn.Root == nil {
		go tn.Flush(tn.interval)
	}
	return nil
}

func (tn *TreeNode) NodeFlush(fe NodeEventType) {
	tn.Children.Loop(func(p string, n *TreeNode) (breaked bool) {
		n.NodeFlush(NodeUpdate)
		return
	})
	tn._flush(NodeUpdate)
	log.With(log.Type("treeNode"), log.String("path", tn.Path)).Info("Flushed.")
}

func (tn *TreeNode) _flush(fe NodeEventType) error {
	val, stat, err := tn.zk.Conn.Get(tn.Path)
	if err != nil {
		log.With(log.Type("treeNode")).Errorf("Flush failed on %v: %v", tn.Path, err.Error())
		return err
	}
	if u := tn.version.Update(stat.Version); u {
		ev := &NodeEvent{
			Val:   val,
			Path:  tn.Path,
			Node:  tn.Node,
			Event: fe,
			Stat:  stat,
		}
		tn.Emit(ev)
	}

	children, stat, err := tn.zk.Conn.Children(tn.Path)
	if err != nil {
		log.With(log.Type("treeNode")).Errorf("Flush children failed on %v: %v", tn.Path, err.Error())
		return err
	}
	if u := tn.cversion.Update(stat.Cversion); u {
		for _, c := range children {
			p := path.Join(tn.Path, c)
			if _, e := tn.Children.Get(p); e {
				continue
			}
			nn, err := NewTreeNode(tn.zk, p, c, tn, tn.interval)
			if err != nil {
				log.With(log.Type("treeNode")).Errorf("NewTreeNode failed on %v: %v", p, err.Error())
				return err
			}
			if err := nn.Init(); err != nil {
				return err
			}
			tn.Children.SetIfAbsent(p, nn)
		}
	}
	return nil
}

func (tn *TreeNode) Flush(interval time.Duration) {
	t := time.NewTicker(interval)
	defer t.Stop()
	for {
		select {
		case <-tn.stop.flushStop.Chan():
			log.With(log.Type("treeNode")).Errorf("Flush stopped on %v.", tn.Path)
			return
		case <-t.C:
			tn.NodeFlush(NodeUpdate)
		}
	}
}

func (tn *TreeNode) Echo(e zk.Event) {
	switch e.Type {
	case zk.EventNodeDeleted:
		ev := &NodeEvent{
			Path:  tn.Path,
			Node:  tn.Node,
			Event: NodeRemoved,
		}
		tn.Emit(ev)
		if err := tn.Clear(); err != nil {
			log.With(log.Type("treeNode")).Errorf("Clear failed on %v: %v", tn.Path, err.Error())
			return
		}
	case zk.EventNodeDataChanged:
		defer tn.GetW()
		val, stat, err := tn.zk.Conn.Get(tn.Path)
		if err != nil {
			log.With(log.Type("treeNode")).Errorf("Get has error on %v: %v", tn.Path, err.Error())
			return
		}
		if u := tn.version.Update(stat.Version); !u {
			return
		}
		ev := &NodeEvent{
			Val:   val,
			Path:  tn.Path,
			Node:  tn.Node,
			Event: NodeUpdate,
			Stat:  stat,
		}
		tn.Emit(ev)
	case zk.EventNodeChildrenChanged:
		defer tn.ChildrenW()
		children, stat, err := tn.zk.Conn.Children(tn.Path)
		if err != nil {
			log.With(log.Type("treeNode")).Errorf("Children has error on %v: %v", tn.Path, err.Error())
			return
		}
		log.With(log.Type("treeNode")).Infof("e.Type: %v, stat: %#v, %v, %v", e.Type, *stat, tn.cversion.version, children)
		if u := tn.cversion.Update(stat.Cversion); !u {
			return
		}
		for _, c := range children {
			p := path.Join(tn.Path, c)
			if _, e := tn.Children.Get(p); e {
				continue
			}
			nn, err := NewTreeNode(tn.zk, p, c, tn, tn.interval)
			if err != nil {
				log.With(log.Type("treeNode")).Errorf("NewTreeNode failed on %v: %v", p, err.Error())
				return
			}
			if err := nn.Init(); err != nil {
				log.With(log.Type("treeNode")).Errorf("NewTreeNode init failed on %v: %v", p, err.Error())
				return
			}
			tn.Children.SetIfAbsent(p, nn)
		}
	}
	return
}

func (tn *TreeNode) GetW() {
	for {
		_, _, _, err := tn.zk.Conn.GetW(tn.Path)
		if err != nil {
			log.With(log.Type("treeNode")).Errorf("GetW failed on %v: %v", tn.Path, err.Error())
			if tn.stop.stop.Exited() || err == zk.ErrNoNode {
				if err := tn.Clear(); err != nil {
					log.With(log.Type("treeNode")).Errorf("Clear failed on %v: %v", tn.Path, err.Error())
				}
				return
			}
			continue
		} else {
			return
		}
	}
}

func (tn *TreeNode) ChildrenW() { //Won't deal with events on child removed.
	for {
		_, _, _, err := tn.zk.Conn.ChildrenW(tn.Path)
		if err != nil {
			log.With(log.Type("treeNode")).Errorf("ChildrenW failed on %v: %v", tn.Path, err.Error())
			if tn.stop.childStop.Exited() || err == zk.ErrNoNode {
				if err := tn.Clear(); err != nil {
					log.With(log.Type("treeNode")).Errorf("Clear failed on %v: %v", tn.Path, err.Error())
				}
				return
			}
			continue
		} else {
			return
		}
	}
}

func (tn *TreeNode) RemoveChild(path string) {
	tn.Children.Remove(path)
}

func (tn *TreeNode) Clear() error {
	if cleared := atomic.LoadInt32(&tn.cleared); cleared != 0 {
		return nil
	}
	atomic.StoreInt32(&tn.cleared, 1)
	tn.stop.Stop()
	tn.Children.Loop(func(path string, node *TreeNode) (breaked bool) {
		if err := node.Clear(); err != nil {
			log.With(log.Type("treeNode")).Errorf("Clear failed on %v: %v", node.Path, err.Error())
		}
		return
	})
	tn.zk.treeNodes.Remove(tn.Path)
	if tn.Parent != nil {
		go tn.Parent.RemoveChild(tn.Path) //deadlock =.=
	}
	tn.Emit(&NodeEvent{
		Node:  tn.Node,
		Path:  tn.Path,
		Event: NodeRemoved,
	})
	return nil
}

func (tn *TreeNode) Emit(e *NodeEvent) {
	log.With(log.Type("treeNode"), log.String("path", tn.Path)).Infof("Emit %#v, Root: %v", e, tn.Root)
	if tn.Root != nil {
		tn.Root.Emit(e)
	} else {
		tn.Event <- e
	}
}

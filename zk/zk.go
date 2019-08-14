package zk

import (
	"bytes"
	"github.com/athlum/pkg/exitChan"
	"github.com/athlum/pkg/log"
	"github.com/athlum/pkg/utils"
	"github.com/samuel/go-zookeeper/zk"
	"strings"
)

type ZK struct {
	Conn      *zk.Conn
	Ec        <-chan zk.Event
	RootPath  string
	Auth      []byte
	Address   string
	Acls      []zk.ACL
	treeNodes *TreeNodeMap
	stop      *exitChan.ExitChan
}

func NewZK(cfg *Config) *ZK {
	client := &ZK{
		Auth:      []byte(cfg.Auth),
		RootPath:  cfg.RootPath,
		Address:   utils.GetLocalIP(),
		treeNodes: NewTreeNodeMap(),
		stop:      exitChan.NewExitChan(),
	}

	conn, ec, err := zk.Connect(cfg.Host, utils.Duration(cfg.SessionTimeout), zk.WithEventCallback(client.eventCallback))
	if err != nil {
		panic(err)
	}
	client.Conn = conn
	client.Ec = ec

	if cfg.Auth != "" {
		auth := strings.Split(cfg.Auth, ":")
		acls := []zk.ACL{}
		acls = zk.DigestACL(zk.PermAll, auth[0], auth[1])
		client.Acls = append(acls, zk.WorldACL(zk.PermRead)...)
	}

	if err := client.Conn.AddAuth("digest", client.Auth); err != nil {
		panic(err)
	}
	if err := client.CheckRoot(); err != nil {
		log.With(log.Type("zk")).Errorf("Check root error: %v", err.Error())
	}
	// go client.WatchEvent()
	return client
}

func (o *ZK) CheckRoot() error {
	tempPath := bytes.NewBuffer([]byte{o.RootPath[0]})
	pathList := strings.Split(string(o.RootPath[1:]), "/")
	for _, value := range pathList {
		tempPath.WriteString(value)
		exists, _, err := o.Conn.Exists(tempPath.String())
		if err != nil {
			return err
		}
		if !exists {
			if _, err := o.Conn.Create(tempPath.String(), []byte{}, 0, o.Acls); err != nil {
				return err
			}
		}
		tempPath.WriteString("/")
	}
	return nil
}

func (o *ZK) Stop() {
	defer o.stop.Close()
	o.treeNodes.Loop(func(k string, node *TreeNode) bool {
		node.Clear()
		return false
	})
}

func (o *ZK) WatchEvent() {
	for {
		select {
		case <-o.stop.Chan():
			return
		case e := <-o.Ec:
			o.eventCallback(e)
		}
	}
}

func (o *ZK) eventCallback(e zk.Event) {
	if e.Path == "" {
		return
	}
	tn, ok := o.treeNodes.Get(e.Path)
	if !ok {
		return
	}
	go tn.Echo(e)
}

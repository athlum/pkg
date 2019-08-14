package zk

import (
	"fmt"
	"github.com/athlum/pkg/log"
	"github.com/samuel/go-zookeeper/zk"
	"net/http"
	_ "net/http/pprof"
	"path"
	"sync"
	"testing"
	"time"
)

func watchNode(conn *zk.Conn, path string) {
	conn.GetW(path)
	conn.ChildrenW(path)
}

func TestZKWatch(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	conn, ec, err := zk.Connect([]string{"testserver:2181"}, time.Second*10, zk.WithEventCallback(func(e zk.Event) {
		fmt.Println(e.Path, e.Type.String())
	}))
	if err != nil {
		panic(err)
	}

	go func() {
		for e := range ec {
			fmt.Println("ec", e.Path, e.Type.String())
		}
	}()

	watchNode(conn, "/test1")
	cs, _, err := conn.Children("/test1")
	for _, c := range cs {
		watchNode(conn, path.Join("/test1", c))
	}
	wg.Wait()
}

func TestTreeCache(t *testing.T) {
	log.Stdout()
	go func() {
		http.ListenAndServe("0.0.0.0:10030", nil)
	}()
	wg := &sync.WaitGroup{}
	wg.Add(1)

	c := NewZK(&Config{
		Host:           []string{"testserver:2181"},
		Auth:           "testauth",
		SessionTimeout: 10.0,
		RootPath:       "/",
	})
	tn, err := c.WatchNode("/test1", "test1", nil, time.Second*10)
	if err != nil {
		panic(err)
	}
	go func(tn *TreeNode) {
		for {
			e := <-tn.Event
			fmt.Println("emit", *e)
		}
	}(tn)
	if err := tn.Init(); err != nil {
		panic(err)
	}
	wg.Wait()
}

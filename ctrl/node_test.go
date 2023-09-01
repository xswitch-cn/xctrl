package ctrl

import (
	"testing"
	"time"

	"git.xswitch.cn/xswitch/proto/xctrl/store"
	"git.xswitch.cn/xswitch/proto/xctrl/store/memory"

	"git.xswitch.cn/xswitch/proto/go/proto/xctrl"
)

func TestNode(t *testing.T) {
	hostname := "test.test"
	node := &xctrl.Node{
		Uuid: "test",
		Name: hostname,
		Rack: 99,
	}
	nodes.Store(hostname, node)
	list := GetNodeList()
	if len(list) != 1 {
		t.Errorf("GetNodeList() = %v, want %v", len(list), 1)
	}
	node1 := list[hostname]
	if node.Name != node1.Name {
		t.Errorf("GetNodeList2() = %v, want %v", node.Name, node.Name)
	}

	time.Sleep(100 * time.Millisecond)
	list = GetNodeList()
	if len(list) != 1 {
		t.Errorf("GetNodeList() = %v, want %v", len(list), 1)
	}

}

func TestOnEvicted(t *testing.T) {
	hostname := "test.test"

	nodes.store = memory.NewStore(store.Table("xnodes"), store.WithCleanupInterval(1*time.Second))

	if err := nodes.store.Write(&store.Record{
		Key:    hostname,
		Value:  []byte("Hello"),
		Expiry: 1 * time.Second,
	}); err != nil {
		return
	}

	rec, err := nodes.store.Read(hostname)
	if err != nil {
		t.Error(err)
	}
	if len(rec) != 1 {
		t.Error("rec len is not 1")
	}

	work := false

	nodes.store.OnEvicted(func(s string, i interface{}) {
		t.Logf("Node %s has expired", s)
		work = true
	})

	time.Sleep(2 * time.Second)

	_, err = nodes.store.Read(hostname)
	if err != nil {
		t.Log("the item is not expired")
	}

	if !work {
		t.Error()
	}
}

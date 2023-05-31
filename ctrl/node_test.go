package ctrl

import (
	"testing"
	"time"

	"git.xswitch.cn/xswitch/xctrl/proto/xctrl"
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

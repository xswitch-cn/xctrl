package consistent

import (
	"math/rand"
	"strconv"
	"testing"

	"git.xswitch.cn/xswitch/proto/go/proto/xctrl"
	"git.xswitch.cn/xswitch/proto/xctrl/util/log"
)

// 基准测试

func TestConsistentHash_AddNodes(t *testing.T) {
	Init(100)
	nodes := make([]*HashNode, 0)
	for i := 1; i < 10; i++ {

		nodes = append(nodes, &HashNode{Node: xctrl.Node{Name: "xcc-node-1" + strconv.Itoa(i)}})
	}
	err := AddNodes(nodes...)
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(defaultConsistentHash.NodesHashes) != 9*100 {
		log.Error(len(defaultConsistentHash.NodesHashes))
		t.Fatal("c.nodesHashes should be 9*100")
	}
}

func TestConsistentHash_Get(t *testing.T) {
	Init(100)
	nodess := []*HashNode{
		{
			Node: xctrl.Node{Name: "xcc-node-1", Uuid: "1"},
		}, {
			Node: xctrl.Node{Name: "xcc-node-2", Uuid: "2"},
		}, {
			Node: xctrl.Node{Name: "xcc-node-3", Uuid: "3"},
		},
	}
	AddNodes(nodess...)
	confNames := []string{"123-dev-xswitch.cn", "456-dev-xswitch.cn", "123-dev-xswitch.cn", "456-dev-xswitch.cn"}
	var nodes []*HashNode
	for k := range confNames {
		node, err := Get(confNames[k])
		if err != nil {
			t.Fail()
		} else {
			nodes = append(nodes, node)
		}
	}
	if nodes[0] == nodes[2] && nodes[1] == nodes[3] {
		t.Log("ok")
	} else {
		t.Fatalf("?????")
	}
}

func TestConsistentHash_DeleteNodes(t *testing.T) {
	Init(100)
	nodes := make([]*HashNode, 0)
	for i := 1; i < 10; i++ {
		nodes = append(nodes, &HashNode{
			Node: xctrl.Node{
				Uuid: strconv.Itoa(i),
				Name: strconv.Itoa(i),
			},
		})
	}
	AddNodes(nodes...)
	// 10*10个虚拟节点
	node, err := Get(strconv.Itoa(rand.Intn(1000000)) + "7000000")
	if err != nil {
		log.Fatal(err)
	} else {
		// 删除调这个节点
		originLen := len(defaultConsistentHash.NodesHashes)
		DeleteNodes(node)
		//理论上应该NodesHashes应该会减少100个，所以..
		nowLen := len(defaultConsistentHash.NodesHashes)
		if originLen-nowLen != 100 {
			t.Fatal()
		} else {
			t.Log("ok")
		}
	}
}

// 性能测试

func BenchmarkConsistentHash_Get(b *testing.B) {
	// 使用默认的100个虚拟节点
	Init(0)
	// 10个真实节点
	nodes := make([]*HashNode, 0)
	for i := 1; i < 10; i++ {
		nodes = append(nodes, &HashNode{Node: xctrl.Node{Uuid: strconv.Itoa(i), Name: "xcc-node-1" + strconv.Itoa(i)}})
	}
	AddNodes(nodes...)
	for i := 0; i < b.N; i++ {
		Get(strconv.Itoa(rand.Intn(1000000)) + "7000000")
	}
	// 结果： 执行次数 10299687，平均执行时间113.6 ns/op
}

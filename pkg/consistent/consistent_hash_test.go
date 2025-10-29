package consistent

import (
	"github.com/google/uuid"
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

		nodes = append(nodes, &HashNode{Node: &xctrl.Node{Uuid: uuid.New().String(), Name: "xcc-node-1" + strconv.Itoa(i)}})
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
			Node: &xctrl.Node{Name: "xcc-node-1", Uuid: "1"},
		}, {
			Node: &xctrl.Node{Name: "xcc-node-2", Uuid: "2"},
		}, {
			Node: &xctrl.Node{Name: "xcc-node-3", Uuid: "3"},
		}, {
			Node: &xctrl.Node{Name: "xcc-node-4", Uuid: "4"},
		},
	}

	err := AddNodes(nodess...)
	if err != nil {
		t.Fatalf("添加节点失败: %v", err)
	}

	// 测试1: 验证相同key返回相同节点
	confNames := []string{"123-dev-xswitch.cn", "456-dev-xswitch.cn", "123-dev-xswitch.cn", "456-dev-xswitch.cn"}
	var nodes []*HashNode
	for k := range confNames {
		node, err := Get(confNames[k])
		if err != nil {
			t.Fatalf("获取节点失败: %v", err)
		} else {
			nodes = append(nodes, node)
		}
	}

	// 验证相同key返回相同节点
	if nodes[0] == nodes[2] && nodes[1] == nodes[3] {
		t.Log("相同key返回相同节点 - 测试通过")
	} else {
		t.Fatalf("相同key应该返回相同节点")
	}

	// 测试2: 验证重复的thor name分配到相同节点，不同的name可以均匀分布
	thorNames := []string{
		"thor_916344178-mihoyo-sip-beta.hoyowave.com",
		"thor_193787997-mihoyo-sip-beta.hoyowave.com",
		"thor_435064194-mihoyo-sip-beta.hoyowave.com",
		"thor_752678318-mihoyo-sip-beta.hoyowave.com",
		"thor_652678312-mihoyo-sip-beta.hoyowave.com",
		"thor_916344178-mihoyo-sip-beta.hoyowave.com", // 重复
		"thor_435064194-mihoyo-sip-beta.hoyowave.com", // 重复
	}

	t.Log("\n=== thor name分配测试 ===")
	thorDistribution := make(map[string]int)
	thorNodeMapping := make(map[string]string) // 记录每个name分配的节点
	thorNodes := make([]string, 0)

	for _, name := range thorNames {
		node, err := Get(name)
		if err != nil {
			t.Fatalf("获取节点失败: %v", err)
		}
		thorDistribution[node.Name]++
		thorNodes = append(thorNodes, node.Name)
		thorNodeMapping[name] = node.Name
		t.Logf("Name: %s -> Node: %s", name, node.Name)
	}

	// 验证重复name分配到相同节点
	t.Logf("\n=== 重复name验证 ===")
	// 检查 thor_916344178... name
	if thorNodeMapping["thor_916344178-mihoyo-sip-beta.hoyowave.com"] != "" {
		node1 := thorNodeMapping["thor_916344178-mihoyo-sip-beta.hoyowave.com"]
		// 这个name出现了两次，应该分配到同一个节点
		if thorDistribution[node1] >= 2 {
			t.Logf("thor_916344178... 分配到相同节点: %s", node1)
		} else {
			t.Errorf("thor_916344178... 没有分配到相同节点")
		}
	}

	// 检查 thor_435064194... 的名字
	if thorNodeMapping["thor_435064194-mihoyo-sip-beta.hoyowave.com"] != "" {
		node2 := thorNodeMapping["thor_435064194-mihoyo-sip-beta.hoyowave.com"]
		if thorDistribution[node2] >= 2 {
			t.Logf("thor_435064194... 分配到相同节点: %s", node2)
		} else {
			t.Errorf("thor_435064194... 没有分配到相同节点")
		}
	}

	// 分布统计
	t.Logf("\n=== 分布统计 ===")
	for node, count := range thorDistribution {
		t.Logf("  %s: %d 个name", node, count)
	}

	// 验证分布结果
	uniqueNodes := len(thorDistribution)
	t.Logf("\n=== 分布分析 ===")
	t.Logf("总共 %d 个name（包含重复）分配到了 %d 个不同节点", len(thorNames), uniqueNodes)

	// 预期的结果是相似的name应该被分散到多个节点，但一样的name必须分配到相同节点
	if uniqueNodes == 1 {
		t.Log("所有name都被分配到同一个节点，分布不够均匀")
	} else if uniqueNodes >= 2 {
		t.Logf("name被分散到 %d 个不同节点，分布良好", uniqueNodes)
	}

	// 验证有没有空节点
	for i, node := range thorNodes {
		if node == "" {
			t.Fatalf(" 第 %d 个name返回了空节点", i)
		}
	}

	// 验证虚拟节点数量 - 修正：使用实际的100个虚拟节点
	expectedVirtualNodes := len(nodess) * 100 * 3 // 4个节点 × 100虚拟节点 × 3种模式
	actualVirtualNodes := GetVirtualNodeCount()
	if actualVirtualNodes != expectedVirtualNodes {
		t.Errorf("虚拟节点数量错误: 期望 %d, 实际 %d", expectedVirtualNodes, actualVirtualNodes)
	} else {
		t.Logf("虚拟节点数量正确: %d", actualVirtualNodes)
	}

	// 验证真实节点数量
	actualNodes := GetNodeCount()
	if actualNodes != len(nodess) {
		t.Fatalf("真实节点数量错误: 期望 %d, 实际 %d", len(nodess), actualNodes)
	} else {
		t.Logf("真实节点数量正确: %d", actualNodes)
	}

	// 额外测试：显式验证重复name的一致性
	t.Logf("\n=== 重复name一致性验证 ===")
	testRepeatNames := []string{
		"thor_916344178-mihoyo-sip-beta.hoyowave.com",
		"thor_916344178-mihoyo-sip-beta.hoyowave.com",
		"thor_916344178-mihoyo-sip-beta.hoyowave.com",
	}

	var repeatNodes []string
	for _, name := range testRepeatNames {
		node, err := Get(name)
		if err != nil {
			t.Fatalf("获取节点失败: %v", err)
		}
		repeatNodes = append(repeatNodes, node.Name)
	}

	// 验证所有重复name都分配到同一个节点
	allSame := true
	for i := 1; i < len(repeatNodes); i++ {
		if repeatNodes[i] != repeatNodes[0] {
			allSame = false
			break
		}
	}

	if allSame {
		t.Logf("重复name一致性测试通过: 所有 '%s' 都分配到节点 %s",
			testRepeatNames[0], repeatNodes[0])
	} else {
		t.Errorf("重复name一致性测试失败: %v", repeatNodes)
	}
}

func TestConsistentHash_DeleteNodes(t *testing.T) {
	Init(100)
	nodes := make([]*HashNode, 0)
	for i := 1; i < 10; i++ {
		nodes = append(nodes, &HashNode{
			Node: &xctrl.Node{
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
		nodes = append(nodes, &HashNode{Node: &xctrl.Node{Uuid: strconv.Itoa(i), Name: "xcc-node-1" + strconv.Itoa(i)}})
	}
	AddNodes(nodes...)
	for i := 0; i < b.N; i++ {
		Get(strconv.Itoa(rand.Intn(1000000)) + "7000000")
	}
	// 结果： 执行次数 10299687，平均执行时间113.6 ns/op
}

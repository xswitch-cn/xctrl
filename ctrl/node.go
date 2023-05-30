package ctrl

import (
	"sync"
	"time"

	"git.xswitch.cn/xswitch/xctrl/proto/xctrl"
	"git.xswitch.cn/xswitch/xctrl/xctrl/store"
	"git.xswitch.cn/xswitch/xctrl/xctrl/store/memory"
)

var nodes nodeList

// nodeList .
type nodeList struct {
	list  map[ /*hostname*/ string]*xctrl.Node
	store store.Store

	sm sync.RWMutex
}

func init() {
	nodes.list = make(map[string]*xctrl.Node, 0)
	nodes.store = memory.NewStore(store.Table("xnodes"))
}

// Hostname 根据 node uuid 获取 hostname
func Hostname(uuid string) string {
	for hostname, node := range nodes.list {
		if node.Uuid == uuid {
			return hostname
		}
	}
	return ""
}

// Node 根据 hostname 获取 node 节点信息
func Node(hostname string) *xctrl.Node {
	if node, ok := nodes.list[hostname]; ok {
		return node
	}
	return nil
}

// GetNodeList 获取当前注册节点列表
func GetNodeList() map[string]*xctrl.Node {
	return nodes.list
}

func GetNodeList2() map[string]*xctrl.Node {
	list, err := nodes.store.List()
	nodes.list = make(map[string]*xctrl.Node, 0)

	nodeList := map[string]*xctrl.Node{}

	if err != nil {
		return nodeList
	}

	for _, key := range list {
		rec, err := nodes.store.Read(key)
		if err != nil {
			continue
		}
		if len(rec) < 1 {
			continue
		}
		nodeList[key] = rec[0].Metadata[key].(*xctrl.Node)
	}

	return nodeList
}

// Store 保存节点信息
func (x *nodeList) Store(hostname string, node *xctrl.Node) {
	x.sm.Lock()
	x.list[hostname] = node
	x.sm.Unlock()

	x.store.Write(&store.Record{
		Key: hostname,
		Metadata: map[string]interface{}{
			hostname: node,
		},
		Expiry: time.Second * 45, // the node should tick every 20 seconds, so 45 seconds is enough for 2 ticks
	})
}

// Delete 删除节点信息
func (x *nodeList) Delete(hostname string) {
	x.sm.Lock()
	delete(x.list, hostname)
	x.sm.Unlock()

	x.store.Delete(hostname)
}

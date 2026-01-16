package ctrl

import (
	"time"

	"github.com/xswitch-cn/proto/go/proto/xctrl"
	"github.com/xswitch-cn/proto/xctrl/store"
)

// nodeList .
type nodeList struct {
	store store.Store
}

// Hostname 根据 node uuid 获取 hostname
func Hostname(uuid string) string {
	list, err := globalCtrl.nodes.store.List()

	if globalCtrl == nil || err != nil {
		return ""
	}

	var node *xctrl.Node
	for _, key := range list {
		rec, err := globalCtrl.nodes.store.Read(key)
		if err != nil {
			continue
		}
		if len(rec) < 1 {
			continue
		}
		if node.Uuid == uuid {
			node = rec[0].Metadata[key].(*xctrl.Node)
		}
	}

	return node.Name
}

// Node 根据 hostname 获取 node 节点信息
func Node(hostname string) *xctrl.Node {
	list, err := globalCtrl.nodes.store.List()

	if globalCtrl == nil || err != nil {
		return nil
	}
	var node *xctrl.Node

	for _, key := range list {
		rec, err := globalCtrl.nodes.store.Read(key)
		if err != nil {
			continue
		}
		if len(rec) < 1 {
			continue
		}
		if node.Name == hostname {
			node = rec[0].Metadata[key].(*xctrl.Node)
		}
	}

	return node
}

// GetNodeList 获取当前注册节点列表
// func GetNodeList() map[string]*xctrl.Node {
// 	return nodes.list
// }

func GetNodeList() map[string]*xctrl.Node {

	list, err := globalCtrl.nodes.store.List()

	nodeList := map[string]*xctrl.Node{}
	if globalCtrl == nil || err != nil {
		return nodeList
	}

	for _, key := range list {
		rec, err := globalCtrl.nodes.store.Read(key)
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
	x.store.Delete(hostname)
}

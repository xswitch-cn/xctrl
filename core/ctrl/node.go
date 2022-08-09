package ctrl

import (
	"sync"

	"git.xswitch.cn/xswitch/xctrl/core/proto/xctrl"
)

var nodes nodeList

// nodeList .
type nodeList struct {
	list map[ /*hostname*/ string]*xctrl.Node

	sm sync.RWMutex
}

func init() {
	nodes.list = make(map[string]*xctrl.Node, 0)
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

// Store 保存节点信息
func (x *nodeList) Store(hostname string, node *xctrl.Node) {
	x.sm.Lock()
	x.list[hostname] = node
	x.sm.Unlock()
}

// Delete 删除节点信息
func (x *nodeList) Delete(hostname string) {
	x.sm.Lock()
	delete(x.list, hostname)
	x.sm.Unlock()
}

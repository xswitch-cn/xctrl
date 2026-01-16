package ctrl

import (
	"github.com/xswitch-cn/proto/go/proto/xctrl"
	"github.com/xswitch-cn/proto/xctrl/store"
	"github.com/xswitch-cn/proto/xctrl/store/memory"
	"time"
)

type CtrlNodes struct {
	store store.Store
}

func InitCtrlNodes() CtrlNodes {
	newStore := memory.NewStore(store.Table("xnodes"), store.WithCleanupInterval(10*time.Second))

	return CtrlNodes{
		store: newStore,
	}

}

func (c *CtrlNodes) Hostname(uuid string) string {
	list, err := c.store.List()

	if err != nil {
		return ""
	}

	var node *xctrl.Node
	for _, key := range list {
		rec, err := c.store.Read(key)
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
func (c *CtrlNodes) Node(hostname string) *xctrl.Node {
	list, err := c.store.List()

	if err != nil {
		return nil
	}

	var node *xctrl.Node
	for _, key := range list {
		rec, err := c.store.Read(key)
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

func (c *CtrlNodes) GetNodeList() map[string]*xctrl.Node {
	list, err := c.store.List()

	nodeList := map[string]*xctrl.Node{}

	if err != nil {
		return nodeList
	}

	for _, key := range list {
		rec, err := c.store.Read(key)
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
func (c *CtrlNodes) Store(hostname string, node *xctrl.Node) {
	c.store.Write(&store.Record{
		Key: hostname,
		Metadata: map[string]interface{}{
			hostname: node,
		},
		Expiry: time.Second * 45, // the node should tick every 20 seconds, so 45 seconds is enough for 2 ticks
	})
}

// Delete 删除节点信息
func (c *CtrlNodes) Delete(hostname string) {
	c.store.Delete(hostname)
}

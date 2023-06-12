* [consitent是什么](#what)
* [consitent api介绍](#api)
* [如何在xctrl中使用](#how)


# <span id="what">consitent是什么</span>

FreeSWITCH多节点一致性hash管理包，主要用于在部署多台FreeSWITCH的情况下，当需要随机获取一台FreeSWITCH作为节点的时候，可通过一个key(可以是一个会议号，一个客户id或者别的一个唯一标识等等)随机获取到一个节点



# <span id="api">consitent api介绍</span>

## type HashFunc

```
type HashFunc func([]byte) uint32
```
hash计算函数，可以自定义


## type HashNode
``` 
type HashNode struct {
	xctrl.Node
	Port        int          `json:"port"`
}
```
基于xctrl的Node的封装

## func Init

```
func Init(virtualNodesNums int, hashFunc ...HashFunc)
```
此包的初始化方法，此方法内部采用了sync.once.do保证了只会被初始化一次，多次调用无效，后者并不会覆盖前者的初始化参数


## func func AddNodes
```
func AddNodes(nodes ...*HashNode) error
```
添加节点

## func ExistNode
```
func AddNodes(nodes ...*HashNode) error
```
某个节点是否存在


## func Get
```
func Get(key string) (*HashNode, error)
```
根据某个特定标识获取一个节点

## func DeleteNodes
```
func DeleteNodes(node *HashNode) error
```
删除节点


# <span id="how">如何在xctrl中使用</span>

在ctrl完成init初始化后，调用ctrl.RegisterHashNodeFun方法，此方法需要传入一个NodeHashFun类型的回调函数,原型如下：
```
type NodeHashFun func(node *xctrl.Node, method string)
```
此回调函数会自动在收到FreeSWITCH的节点相关消息的时候由系统自动回调

# example

```
	err := ctrl.Init(true, "nats://localhost:4222")
	// 初始化一百个虚拟节点
	consistent.Init(100)
	// 注册节点事件回调方法
	ctrl.RegisterHashNodeFun(func(node *xctrl.Node, method string) {
		hashNode := new(consistent.HashNode)
		hashNode.Node = node
		for _, v := range node.SipProfiles {
			if v.Name == "public" {
				hashNode.Port = int(v.Port)
				break
			}
		}
		switch method {
		case "Event.NodeRegister":
			// 收到节点注册事件，如果hash环里不存在则添加
			if !consistent.ExistNode(hashNode) {
				consistent.AddNodes(hashNode)
			}
		case "Event.NodeUnregister":
			// 收到节点取消注册事件，如果hash环里存在则删除
			if consistent.ExistNode(hashNode) {
				consistent.DeleteNodes(hashNode)
			}
		case "Event.NodeUpdate":
			// 收到节点更新事件，如果hash环里不存在则添加
			if !consistent.ExistNode(hashNode) {
				consistent.AddNodes(hashNode)
			}
		}
		// you can do anything you like
		// pass
	})
```

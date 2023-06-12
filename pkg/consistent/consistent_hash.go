package consistent

import (
	"errors"
	"hash/crc32"
	"sort"
	"strconv"
	"sync"

	"git.xswitch.cn/xswitch/proto/go/proto/xctrl"
	"git.xswitch.cn/xswitch/proto/xctrl/util/log"
)

// NodesSlice 为了快速查找
type NodesSlice []uint32

type HashFunc func([]byte) uint32

func (s NodesSlice) Len() int {
	return len(s)
}

func (s NodesSlice) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s NodesSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type SipProfile struct {
	Name string `json:"name"`
	Port int    `json:"port"`
}
type HashNode struct {
	*xctrl.Node
	Port int `json:"port"`
}
type consistentHash struct {
	mux              sync.RWMutex         //lock
	hash             HashFunc             //hash算法
	VirtualNodesNums int                  // 虚拟节点数
	Nodes            []*HashNode          // 真实的节点
	NodesHashes      NodesSlice           // 排好序的节点的hash值，数量应该==实际节点数*虚拟节点数
	Hash2Node        map[uint32]*HashNode // 虚拟节点hash值=>节点值，不用实际的节点，因为实际节点一般就几个，不容易平衡
}

var defaultConsistentHash consistentHash

var once sync.Once

// Init
// virtualNodesNums 虚拟节点个数，default 100个
// hashFunc hash方法，默认crc32.ChecksumIEEE
func Init(virtualNodesNums int, hashFunc ...HashFunc) {
	// 避免多次初始化
	(&once).Do(func() {
		if virtualNodesNums == 0 {
			// 一般来说正常就个位数的节点，虚拟节点100个完全足够均衡了
			virtualNodesNums = 100
		}
		defaultConsistentHash.VirtualNodesNums = virtualNodesNums
		if len(hashFunc) > 0 {
			defaultConsistentHash.hash = hashFunc[0]
		} else {
			// 默认hash计算算方法
			defaultConsistentHash.hash = crc32.ChecksumIEEE
		}
		defaultConsistentHash.Hash2Node = make(map[uint32]*HashNode)
	})
}

// AddNodes 添加节点。
func AddNodes(nodes ...*HashNode) error {
	for _, n := range nodes {
		if n.Uuid == "" {
			return errors.New("node uuid is required")
		}
	}
	if len(nodes) == 0 {
		return errors.New("require one node at least")
	}
	// 给虚拟节点赋值
	defaultConsistentHash.mux.Lock()
	defer defaultConsistentHash.mux.Unlock()
	for k := 0; k < len(nodes); k++ {
		defaultConsistentHash.Nodes = append(defaultConsistentHash.Nodes, nodes[k])
		for i := 0; i < defaultConsistentHash.VirtualNodesNums; i++ {
			// 虚拟节点hash
			hashValue := defaultConsistentHash.hash([]byte(strconv.Itoa(i) + nodes[k].Uuid))
			defaultConsistentHash.NodesHashes = append(defaultConsistentHash.NodesHashes, hashValue)
			defaultConsistentHash.Hash2Node[hashValue] = nodes[k]
		}
	}
	sort.Sort(defaultConsistentHash.NodesHashes)
	return nil
}

// ExistNode 是否存在这个节点
func ExistNode(node *HashNode) bool {
	found := false
	for k := 0; k < len(defaultConsistentHash.Nodes); k++ {
		if defaultConsistentHash.Nodes[k].Uuid == node.Uuid {
			found = true
			break
		}
	}
	return found
}

// DeleteNodes 删除节点。
func DeleteNodes(node *HashNode) error {
	if node.Uuid == "" {
		return errors.New("node uuid is required")
	}
	defaultConsistentHash.mux.Lock()
	defer defaultConsistentHash.mux.Unlock()
	found := false
	for k := 0; k < len(defaultConsistentHash.Nodes); k++ {
		if defaultConsistentHash.Nodes[k].Uuid == node.Uuid {
			defaultConsistentHash.Nodes = append(defaultConsistentHash.Nodes[:k], defaultConsistentHash.Nodes[k+1:]...)
			found = true
			break
		}
	}
	if !found {
		// 没必要进行hash计算了
		log.Warnf("deleting a nonexistent node %s", node.Name)
		return errors.New("not found this node " + node.Uuid)
	}

	for i := 0; i < defaultConsistentHash.VirtualNodesNums; i++ {
		hash := defaultConsistentHash.hash([]byte(strconv.Itoa(i) + node.Uuid))
		for j := 0; j < len(defaultConsistentHash.NodesHashes); j++ {
			if defaultConsistentHash.NodesHashes[j] == hash {
				defaultConsistentHash.NodesHashes = append(defaultConsistentHash.NodesHashes[:j], defaultConsistentHash.NodesHashes[j+1:]...)
				break
			}
		}
		delete(defaultConsistentHash.Hash2Node, hash)
	}
	return nil
}

func Get(key string) (*HashNode, error) {
	if len(defaultConsistentHash.NodesHashes) == 0 {
		return nil, errors.New("no node in hash,you can use [AddNodes] method to add some nodes")
	}
	hashValue := defaultConsistentHash.hash([]byte(key))
	k := sort.Search(len(defaultConsistentHash.NodesHashes), func(i int) bool {
		return defaultConsistentHash.NodesHashes[i] > hashValue
	})

	// 如果没找到,放第一个或者最后一个？？
	if k == len(defaultConsistentHash.NodesHashes) {
		k = 0
	}
	defaultConsistentHash.mux.RLock()
	defer defaultConsistentHash.mux.RUnlock()
	return defaultConsistentHash.Hash2Node[defaultConsistentHash.NodesHashes[k]], nil
}

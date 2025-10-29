package consistent

import (
	"errors"
	"fmt"
	"sort"
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

func OptimizedHash(data []byte) uint32 {
	// 使用FNV风格的哈希，性能好且分布均匀
	const (
		offset32 uint32 = 2166136261
		prime32  uint32 = 16777619
	)

	hash := offset32
	for _, c := range data {
		hash *= prime32
		hash ^= uint32(c)
	}

	// 额外的混合步骤，提高分布均匀性
	hash = (hash >> 16) ^ (hash << 16)
	hash *= prime32
	return hash
}

// Init
// virtualNodesNums 虚拟节点个数，默认250个
// hashFunc hash方法，默认SHA1BasedHash
// cman里使用这个默认的方法
func Init(virtualNodesNums int, hashFunc ...HashFunc) {
	(&once).Do(func() {
		if virtualNodesNums == 0 {
			// 增加虚拟节点，保证均衡
			virtualNodesNums = 250
		}
		defaultConsistentHash.VirtualNodesNums = virtualNodesNums
		if len(hashFunc) > 0 {
			defaultConsistentHash.hash = hashFunc[0]
		} else {
			// 默认hash计算算方法
			defaultConsistentHash.hash = OptimizedHash
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

	// 检查是否有重复节点，避免重复添加
	for _, newNode := range nodes {
		for _, existingNode := range defaultConsistentHash.Nodes {
			if existingNode.Uuid == newNode.Uuid {
				return fmt.Errorf("node with UUID %s already exists", newNode.Uuid)
			}
		}
	}

	for k := 0; k < len(nodes); k++ {
		defaultConsistentHash.Nodes = append(defaultConsistentHash.Nodes, nodes[k])
		for i := 0; i < defaultConsistentHash.VirtualNodesNums; i++ {
			// 调整key，增大随机性
			virtualKey1 := fmt.Sprintf("%s|virtual|%d|%d|salt1",
				nodes[k].Uuid, i, k*defaultConsistentHash.VirtualNodesNums+i)
			virtualKey2 := fmt.Sprintf("%d|%s|%d|virtual|salt2",
				i, nodes[k].Uuid, k)
			virtualKey3 := fmt.Sprintf("node|%d|%s|%d|salt3",
				k, nodes[k].Uuid, i)

			// 为每个节点生成多个不同模式的虚拟节点
			hashValue1 := defaultConsistentHash.hash([]byte(virtualKey1))
			hashValue2 := defaultConsistentHash.hash([]byte(virtualKey2))
			hashValue3 := defaultConsistentHash.hash([]byte(virtualKey3))

			defaultConsistentHash.NodesHashes = append(
				defaultConsistentHash.NodesHashes,
				hashValue1, hashValue2, hashValue3,
			)
			defaultConsistentHash.Hash2Node[hashValue1] = nodes[k]
			defaultConsistentHash.Hash2Node[hashValue2] = nodes[k]
			defaultConsistentHash.Hash2Node[hashValue3] = nodes[k]
		}
	}
	sort.Sort(defaultConsistentHash.NodesHashes)
	return nil
}

// ExistNode 是否存在这个节点
func ExistNode(node *HashNode) bool {
	defaultConsistentHash.mux.RLock()
	defer defaultConsistentHash.mux.RUnlock()

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
	index := -1
	for k := 0; k < len(defaultConsistentHash.Nodes); k++ {
		if defaultConsistentHash.Nodes[k].Uuid == node.Uuid {
			index = k
			found = true
			break
		}
	}
	if !found {
		log.Warnf("deleting a nonexistent node %s", node.Name)
		return errors.New("not found this node " + node.Uuid)
	}

	// 从Nodes中删除
	defaultConsistentHash.Nodes = append(defaultConsistentHash.Nodes[:index], defaultConsistentHash.Nodes[index+1:]...)

	// 删除虚拟节点
	for i := 0; i < defaultConsistentHash.VirtualNodesNums; i++ {
		virtualKey1 := fmt.Sprintf("%s|virtual|%d|%d|salt1",
			node.Uuid, i, index*defaultConsistentHash.VirtualNodesNums+i)
		virtualKey2 := fmt.Sprintf("%d|%s|%d|virtual|salt2",
			i, node.Uuid, index)
		virtualKey3 := fmt.Sprintf("node|%d|%s|%d|salt3",
			index, node.Uuid, i)

		hash1 := defaultConsistentHash.hash([]byte(virtualKey1))
		hash2 := defaultConsistentHash.hash([]byte(virtualKey2))
		hash3 := defaultConsistentHash.hash([]byte(virtualKey3))

		hashes := []uint32{hash1, hash2, hash3}
		for _, hash := range hashes {
			// 从NodesHashes中删除
			for j := 0; j < len(defaultConsistentHash.NodesHashes); j++ {
				if defaultConsistentHash.NodesHashes[j] == hash {
					defaultConsistentHash.NodesHashes = append(defaultConsistentHash.NodesHashes[:j], defaultConsistentHash.NodesHashes[j+1:]...)
					break
				}
			}
			// 从Hash2Node中删除
			delete(defaultConsistentHash.Hash2Node, hash)
		}
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

// 节点数量
func GetNodeCount() int {
	defaultConsistentHash.mux.RLock()
	defer defaultConsistentHash.mux.RUnlock()
	return len(defaultConsistentHash.Nodes)
}

// 拟节点数量
func GetVirtualNodeCount() int {
	defaultConsistentHash.mux.RLock()
	defer defaultConsistentHash.mux.RUnlock()
	return len(defaultConsistentHash.NodesHashes)
}
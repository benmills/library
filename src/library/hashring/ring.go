package hashring

import (
	"hash/crc32"
)

const (
	MaxRingSize = 4294967294
	VnodeCount = 1024
)

type Ring struct {
	size uint32
	vnodeCount int
	n int
	nodes []*Node
}

func New() *Ring {
	return &Ring{
		size: MaxRingSize,
		vnodeCount: VnodeCount,
		nodes: []*Node{},
		n: 0,
	}
}

func (ring *Ring) SetNValue(n int) {
	ring.n = n
}

func (ring *Ring) AddNode(name string) *Node {
	vnodeCount := ring.vnodeCount / (ring.NodeCount()+1)

	for _, node := range(ring.nodes) {
		node.vnodeCount = vnodeCount
	}

	newNode := &Node{
		name: name,
		vnodeCount: vnodeCount,
		vnodeStart: ring.getNextNodeVnodeStart(vnodeCount),
		vnodeSize: ring.vnodeSize(),
	}

	if len(ring.nodes) != 0 {
		ring.lastNode().next = newNode
		newNode.next = ring.firstNode()
	}

	ring.nodes = append(ring.nodes, newNode)
	return newNode
}

func (ring *Ring) SetNodes(nodes []string) {
	ring.nodes = []*Node{}
	for _, nodeName := range nodes {
		ring.AddNode(nodeName)
	}
}

func (ring *Ring) NodeCount() int {
	return len(ring.nodes)
}

func (ring *Ring) PreferenceListForKey(key string) []string {
	preferenceList := []string{}

	n := ring.NodeForKey(key)
	for i := ring.n; i >= 0; i-- {
		preferenceList = append(preferenceList, n.name)
		n = n.next
	}

	return preferenceList
}

func (ring *Ring) NodeForKey(key string) *Node {
	keyHash := hash(key)

	for _, node := range ring.nodes {
		if node.OwnsKeyHash(keyHash) {
			return node
		}
	}

	return ring.lastNode()
}

func (ring *Ring) AddKey(key string) {
	node := ring.NodeForKey(key)
	node.keyCount = node.keyCount + 1
}

func (ring *Ring) GetNValue() int {
	return ring.n
}

func (ring *Ring) GetNodes() []string {
	result := []string{}

	for _, node := range(ring.nodes) {
		result = append(result, node.name)
	}

	return result
}

func (ring *Ring) Get(name string) *Node {
	for _, node := range(ring.nodes) {
		if node.name == name {
			return node
		}
	}

	panic("No node found for name " + name)
}

func (ring *Ring) getNextNodeVnodeStart(nodeVnodeCount int) uint32 {
	nodeCount := ring.NodeCount()

	if nodeCount == 0 {
		return uint32(0)
	}

	nodeSize := uint32(nodeVnodeCount) * ring.vnodeSize()
	return (ring.lastNode().vnodeStart+1) + nodeSize
}

func (ring *Ring) vnodeSize() uint32 {
	return ring.size / uint32(ring.vnodeCount)
}

func (ring *Ring) lastNode() *Node {
	return ring.nodes[len(ring.nodes)-1]
}

func (ring *Ring) firstNode() *Node {
	return ring.nodes[0]
}

func hash(input string) uint32 {
	return crc32.ChecksumIEEE([]byte(input))
}

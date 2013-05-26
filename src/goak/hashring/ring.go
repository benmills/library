package hashring

import (
	"hash/crc32"
)

type Ring struct {
	size uint32
	vnodeCount int
	nodes []*Node
}

func New() *Ring {
	return &Ring{
		size: 4294967294,
		vnodeCount: 1024,
		nodes: []*Node{},
	}
}

func (ring *Ring) AddNode(nodeName string) *Node {
	vnodeCount := ring.vnodeCount / (ring.NodeCount()+1)

	for _, node := range(ring.nodes) {
		node.vnodeCount = vnodeCount
	}

	newNode := &Node{
		name: nodeName,
		hash: hash(nodeName),
		vnodeCount: vnodeCount,
		vnodeStart: ring.getNextNodeVnodeStart(vnodeCount),
		vnodeSize: ring.vnodeSize(),
	}

	ring.nodes = append(ring.nodes, newNode)
	return newNode
}

func (ring *Ring) NodeCount() int {
	return len(ring.nodes)
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

func (ring *Ring) getNextNodeVnodeStart(nodeVnodeCount int) uint32 {
	nodeCount := ring.NodeCount()

	if nodeCount == 0 {
		return uint32(0)
	} else {
		nodeSize := uint32(nodeVnodeCount) * ring.vnodeSize()
		return (ring.lastNode().vnodeStart+1) + nodeSize
	}
}

func (ring *Ring) vnodeSize() uint32 {
	return ring.size / uint32(ring.vnodeCount)
}

func (ring *Ring) lastNode() *Node {
	return ring.nodes[len(ring.nodes)-1]
}

func hash(input string) uint32 {
	return crc32.ChecksumIEEE([]byte(input))
}

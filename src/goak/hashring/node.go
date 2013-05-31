package hashring

type Node struct {
	name string
	hash uint32
	keyCount int
	vnodeCount int
	vnodeSize uint32
	vnodeStart uint32
}

func (node *Node) OwnsKeyHash(keyHash uint32) bool {
	nodeHash := node.vnodeStart

	for i := node.vnodeCount; i >= 0; i-- {
		if keyHash >= node.vnodeStart && keyHash <= nodeHash {
			return true
		}

		nodeHash = nodeHash + node.vnodeSize
	}

	return false
}

func (node *Node) VnodeCount() int {
	return node.vnodeCount;
}

func (node *Node) VnodeSize() uint32 {
	return node.vnodeSize;
}

func (node *Node) VnodeStart() uint32 {
	return node.vnodeStart;
}

func (node *Node) SetName(name string) {
	node.name = name
}

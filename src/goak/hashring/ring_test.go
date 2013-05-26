package hashring

import (
	"testing"
	"github.com/benmills/quiz"
	"strconv"
)

func TestRingCanHaveNodes(t *testing.T) {
	test := quiz.Test(t)

	ring := New()
	ring.AddNode("127.0.0.1")
	ring.AddNode("0.0.0.0")

	test.Expect(ring.NodeCount()).ToEqual(2)
}

func TestNodesSpiltVnodeMaxCount(t *testing.T) {
	test := quiz.Test(t)

	ring := New()
	node1 := ring.AddNode("1")
	node2 := ring.AddNode("2")

	test.Expect(node1.vnodeCount).ToEqual(512)
	test.Expect(node2.vnodeCount).ToEqual(512)
}

func TestNodesHaveCorrectVnodeStart(t *testing.T) {
	test := quiz.Test(t)

	ring := New()
	node1 := ring.AddNode("1")
	node2 := ring.AddNode("2")

	test.Expect(node1.vnodeStart).ToEqual(uint32(0))
	test.Expect(node2.vnodeStart).ToEqual(uint32(2147483137))
}

func TestNodesDontOwnSameKeyHashesOnEdges(t *testing.T) {
	test := quiz.Test(t)

	ring := New()
	node1 := ring.AddNode("1")
	node2 := ring.AddNode("2")

	test.Expect(node1.OwnsKeyHash(0)).ToBeTrue()
	test.Expect(node2.OwnsKeyHash(0)).ToBeFalse()

	test.Expect(node1.OwnsKeyHash(2147483136)).ToBeTrue()
	test.Expect(node2.OwnsKeyHash(2147483136)).ToBeFalse()

	test.Expect(node1.OwnsKeyHash(2147483137)).ToBeFalse()
	test.Expect(node2.OwnsKeyHash(2147483137)).ToBeTrue()
}

func TestKeyDistribution(t *testing.T) {
	test := quiz.Test(t)

	ring := New()
	ring.AddNode("127.0.0.1")
	ring.AddNode("0.0.0.0")
	ring.AddNode("10.0.0.1")

	for i := 1; i <= 100; i++ {
		ring.AddKey(strconv.Itoa(i))
	}

	for _, node := range ring.nodes {
		test.Expect(node.keyCount).ToBeLessThan(50)
		test.Expect(node.keyCount).ToBeGreaterThan(10)
	}
}

package hashring

import (
	"github.com/benmills/quiz"
	"strconv"
	"testing"
)

func TestRingCanHaveNodes(t *testing.T) {
	test := quiz.Test(t)

	ring := New()
	ring.AddNode("A")
	ring.AddNode("B")

	test.Expect(ring.NodeCount()).ToEqual(2)
}

func TestNodesSpiltVnodeMaxCount(t *testing.T) {
	test := quiz.Test(t)

	ring := New()
	node1 := ring.AddNode("A")
	node2 := ring.AddNode("B")

	test.Expect(node1.vnodeCount).ToEqual(512)
	test.Expect(node2.vnodeCount).ToEqual(512)
}

func TestNodesHaveCorrectVnodeStart(t *testing.T) {
	test := quiz.Test(t)

	ring := New()
	node1 := ring.AddNode("A")
	node2 := ring.AddNode("B")

	test.Expect(node1.vnodeStart).ToEqual(uint32(0))
	test.Expect(node2.vnodeStart).ToEqual(uint32(2147483137))
}

func TestNodesDontOwnSameKeyHashesOnEdges(t *testing.T) {
	test := quiz.Test(t)

	ring := New()
	node1 := ring.AddNode("A")
	node2 := ring.AddNode("B")

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
	ring.AddNode("A")
	ring.AddNode("B")
	ring.AddNode("C")

	for i := 1; i <= 100; i++ {
		ring.AddKey(strconv.Itoa(i))
	}

	for _, node := range ring.nodes {
		test.Expect(node.keyCount).ToBeLessThan(50)
		test.Expect(node.keyCount).ToBeGreaterThan(10)
	}
}

func TestGetNodes(t *testing.T) {
	test := quiz.Test(t)

	ring := New()
	ring.AddNode("A")
	ring.AddNode("B")
	ring.AddNode("C")

	test.Expect(ring.GetNodes()[0]).ToEqual("A")
	test.Expect(ring.GetNodes()[1]).ToEqual("B")
	test.Expect(ring.GetNodes()[2]).ToEqual("C")
}

func TestSetNodes(t *testing.T) {
	test := quiz.Test(t)

	ring := New()
	ring.SetNodes([]string{"A", "B", "C"})

	test.Expect(ring.GetNodes()[0]).ToEqual("A")
	test.Expect(ring.GetNodes()[1]).ToEqual("B")
	test.Expect(ring.GetNodes()[2]).ToEqual("C")
}

func TestSetNodesClearsOldNodes(t *testing.T) {
	test := quiz.Test(t)

	ring := New()
	ring.AddNode("old node")
	ring.SetNodes([]string{"A", "B", "C"})

	test.Expect(ring.GetNodes()[0]).ToEqual("A")
	test.Expect(ring.GetNodes()[1]).ToEqual("B")
	test.Expect(ring.GetNodes()[2]).ToEqual("C")
}

package hashring

import (
	"github.com/benmills/quiz"
	"testing"
)

func TestNodeOwnsKeyIfKeyIsVnodeStart(t *testing.T) {
	test := quiz.Test(t)

	node := &Node{
		vnodeCount: 1,
		vnodeSize: 1,
		vnodeStart: 0,
	}

	test.Expect(node.OwnsKeyHash(0)).ToBeTrue()
}

func TestNodeOwnsKeyWithinVnodeRange(t *testing.T) {
	test := quiz.Test(t)

	node := &Node{
		vnodeCount: 1,
		vnodeSize: 1,
		vnodeStart: 0,
	}

	test.Expect(node.OwnsKeyHash(1)).ToBeTrue()
}

func TestNodeDoesntOwnsKeyOutsideVnodeRange(t *testing.T) {
	test := quiz.Test(t)

	node := &Node{
		vnodeCount: 1,
		vnodeSize: 1,
		vnodeStart: 0,
	}

	test.Expect(node.OwnsKeyHash(2)).ToBeFalse()
}

func TestNodeOwnsKeyInbetweenVnodeRange(t *testing.T) {
	test := quiz.Test(t)

	node := &Node{
		vnodeCount: 1,
		vnodeSize: 10,
		vnodeStart: 0,
	}

	test.Expect(node.OwnsKeyHash(7)).ToBeTrue()
}

func TestNodeDoesntOwnKeyBelowVnodeStart(t *testing.T) {
	test := quiz.Test(t)

	node := &Node{
		vnodeCount: 1,
		vnodeSize: 10,
		vnodeStart: 10,
	}

	test.Expect(node.OwnsKeyHash(7)).ToBeFalse()
}

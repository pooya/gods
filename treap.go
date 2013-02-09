// Copyright 2013 Shayan Pooya. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

type Random struct {
	generator *rand.Rand
}

var globalRandomGenerator Random

func (r *Random) init() {
	seed := (int64)(time.Now().Nanosecond())
	s := rand.NewSource(seed)
	r.generator = rand.New(s)
}

func (r *Random) next() int {
	value := r.generator.Int()
	return value
}

func random() int {
	if globalRandomGenerator.generator == nil {
		globalRandomGenerator.init()
	}
	return globalRandomGenerator.next()
}

type Treap struct {
	nonEmpty bool
	head     *Node
}

type Node struct {
	left     *Node
	value    int
	right    *Node
	parent   *Node
	priority int
}

func min(nodes ...*Node) *Node {
	minimum := int(^uint(0) >> 1)
	var minNode *Node
	minNode = nil
	for _, node := range nodes {
		if node != nil && minimum >= node.priority {
			minimum = node.priority
			minNode = node
		}
	}
	return minNode
}

func rotate(child, parent *Node) {
	if parent.parent != nil {
		if parent.parent.right == parent {
			parent.parent.right = child
		} else {
			parent.parent.left = child
		}
	}
	child.parent = parent.parent
	parent.parent = child
	if child == parent.left {
		parent.left = child.right
		if child.right != nil {
			child.right.parent = parent
		}
		child.right = parent
	} else {
		parent.right = child.left
		if child.left != nil {
			child.left.parent = parent
		}
		child.left = parent
	}
}

func (node *Node) bubbleUp(head **Node) {
	if node.parent == nil {
		*head = node
		return
	}
	minNode := min(node, node.parent)
	if minNode == node.parent {
		return
	}
	node.Traverse2()
	rotate(node, node.parent)
	node.Traverse2()
	node.bubbleUp(head)
}

func (t *Node) Insert(value int, head **Node) {
	if value > t.value {
		if t.right != nil {
			t.right.Insert(value, head)
		} else {
			t.right = &Node{nil, value, nil, t, random()}
			t.right.bubbleUp(head)
		}
	} else {
		if t.left != nil {
			t.left.Insert(value, head)
		} else {
			t.left = &Node{nil, value, nil, t, random()}
			t.left.bubbleUp(head)
		}
	}
}

func (t *Node) Min() *Node {
	for t.left != nil {
		t = t.left
	}
	return t
}

func (t *Node) Max() *Node {
	for t.right != nil {
		t = t.right
	}
	return t
}

func (node *Node) Next() *Node {
	if node.right != nil {
		return node.right.Min()
	}

	parent := node.parent
	for parent != nil && parent.right == node {
		node = parent
		parent = parent.parent
	}
	return parent
}

func (node *Node) Prev() *Node {
	if node.left != nil {
		return node.left.Max()
	}

	parent := node.parent
	for parent != nil && parent.left == node {
		node = parent
		parent = parent.parent
	}
	return parent
}

func (t *Node) Traverse1() (nodes []int) {
	if t.left != nil {
		nodes = append(nodes, t.left.Traverse1()...)
	}

	nodes = append(nodes, t.value)
	if t.right != nil {
		nodes = append(nodes, t.right.Traverse1()...)
	}
	return
}

func (t *Node) Traverse2() (nodes []int) {
	m := t.Min()
	for m != nil {
		nodes = append(nodes, m.value)
		m = m.Next()
	}
	return
}

func (t *Node) Traverse3() (nodes []int) {
	_nodes := make([]int, 0)
	m := t.Max()
	for m != nil {
		_nodes = append(_nodes, m.value)
		m = m.Prev()
	}
	nodes = make([]int, len(_nodes))
	for i := 0; i < len(_nodes); i++ {
		nodes[i] = _nodes[len(_nodes)-i-1]
	}
	return
}

func (t *Node) removeNode(head **Node) {
	switch {
	case t.left == nil && t.right == nil:
		if t.parent == nil {
			*head = nil
			return
		}
		if t.parent.right == t {
			t.parent.right = nil
		} else {
			t.parent.left = nil
		}
	case t.left == nil:
		t.right.parent = t.parent
		if t.parent == nil {
			*head = nil
			return
		}
		if t.parent.left == t {
			t.parent.left = t.right
		} else {
			t.parent.right = t.right
		}
	case t.right == nil:
		t.left.parent = t.parent
		if t.parent == nil {
			*head = nil
			return
		}
		if t.parent.left == t {
			t.parent.left = t.left
		} else {
			t.parent.right = t.left
		}
	default:
		fmt.Fprintln(os.Stderr, "Error in removeNode")
	}
}

func (t *Node) remove(head **Node) {
	if t.left == nil || t.right == nil {
		t.removeNode(head)
		return
	}
	successor := t.left.Max()
	t.value = successor.value
	t.priority = successor.priority
	successor.removeNode(head)
	t.bubbleUp(head)
}

func (t *Node) Find(key int) *Node {
	if key == t.value {
		return t
	}
	if key < t.value {
		if t.left == nil {
			return nil
		}
		return t.left.Find(key)
	}
	if t.right == nil {
		return nil
	}
	return t.right.Find(key)
}

func (t *Node) Remove(key int, head **Node) bool {
	node := t.Find(key)
	if node == nil {
		return false
	}
	node.remove(head)
	return true
}

func (t *Treap) Insert(value int) {
	if !t.nonEmpty {
		t.head = &Node{nil, value, nil, nil, random()}
		t.nonEmpty = true
		return
	}
	t.head.Insert(value, &t.head)
}

func (t *Treap) Remove(value int) bool {
	removed := t.head.Remove(value, &t.head)
	if t.head == nil {
		t.nonEmpty = false
	}
	return removed
}

func (t *Treap) Traverse1() (nodes []int) {
	return t.head.Traverse1()
}

func (t *Treap) Traverse2() (nodes []int) {
	return t.head.Traverse2()
}

func (t *Treap) Traverse3() (nodes []int) {
	return t.head.Traverse3()
}

func main() {
	var n int
	var treap Treap

	fmt.Scanf("%d", &n)
	for i := 0; i < n; i++ {
		var a int
		fmt.Scanf("%d", &a)
		treap.Insert(a)
	}
	treap.Traverse1()
	fmt.Println(treap.Traverse2())
}

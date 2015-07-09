package rbtree

import (
	"fmt"
)

type Color uint

func (c Color) String() string {
	if c == Red {
		return "Red"
	}

	return "Black"
}

const (
	Red Color = iota
	Black
)

type Key interface {
	Compare(k Key) int
	Equals(k Key) bool
}

type Value interface{}

type NodeFunc func(value Value)

type Node struct {
	key    Key
	value  Value
	color  Color
	parent *Node
	left   *Node
	right  *Node
}

type Tree struct {
	size int
	root *Node
}

// IterFunc define the signature of functions suited to iterate over tree
type IterFunc func(node *Node, fn NodeFunc)

func (n *Node) String() string {
	var (
		l string = "nil"
		r string = "nil"
	)

	if n.left != nil {
		l = n.left.String()
	}

	if n.right != nil {
		r = n.right.String()
	}

	return fmt.Sprintf("(%v, %v, %s, left: %s, right: %s)", n.key, n.value, n.color, l, r)
}

func (node *Node) Key() Key {
	return node.key
}

func (node *Node) Value() Value {
	return node.value
}

func (node *Node) Next() *Node {
	if node == nil {
		return nil
	}

	if node.right != nil {
		node = node.right
		for node.left != nil {
			node = node.left
		}

		return node
	}

	for node.parent != nil && node.parent.right == node {
		node = node.parent
	}

	return node.parent
}

func (node *Node) Previous() *Node {
	if node == nil {
		return nil
	}

	if node.left != nil {
		node = node.left
		for node.right != nil {
			node = node.right
		}

		return node
	}

	for node.parent != nil && node.parent.left == node {
		node = node.parent
	}

	return node.parent
}

func color(n *Node) Color {
	if n == nil {
		return Black
	}

	return n.color
}

func (n *Node) detach() {
	n.parent = nil
	n.left = nil
	n.right = nil
}

func (n *Node) grandparent() *Node {
	if n.parent != nil {
		return n.parent.parent
	}

	return nil
}

func (n *Node) sibling() *Node {
	if n.parent == nil {
		return nil
	}

	if n.parent.left == n {
		return n.parent.right
	}

	return n.parent.left
}

func (n *Node) uncle() *Node {
	if n.parent == nil {
		return nil
	}

	return n.parent.sibling()
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func height(n *Node) int {
	if n == nil {
		return 0
	}

	return 1 + max(height(n.left), height(n.right))
}

func New() *Tree {
	return &Tree{size: 0, root: nil}
}

func (tree *Tree) String() string {
	if tree.root != nil {
		return tree.root.String()
	}

	return "nil"
}

func (tree *Tree) Size() int {
	return tree.size
}

func (tree *Tree) Height() int {
	return height(tree.root)
}

func (tree *Tree) First() *Node {
	n := tree.root
	for n.left != nil {
		n = n.left
	}

	return n
}

func (tree *Tree) Last() *Node {
	n := tree.root
	for n.right != nil {
		n = n.right
	}

	return n
}

func (tree *Tree) End() *Node {
	return nil
}

func (tree *Tree) Do(fn NodeFunc, iter IterFunc) {
	if iter == nil {
		iter = InOrder
	}
	n := tree.root

	iter(n, fn)
}

func (tree *Tree) Lookup(key Key) Value {
	res := tree.lookup(key)
	if res != nil {
		return res.value
	}

	return nil
}

// find recurrence of Key key into the tree, return the node if any
func (tree *Tree) lookup(key Key) *Node {
	node := tree.root
	for node != nil {
		compare := node.key.Compare(key)
		// found key into the tree, return node
		if compare == 0 {
			return node
		}

		if compare < 0 {
			node = node.right
		}

		if compare > 0 {
			node = node.left
		}
	}

	return nil
}

func (tree *Tree) Insert(key Key, value Value) {
	n := &Node{key: key, value: value, color: Red}
	tree.insert(n)
}

func (tree *Tree) insert(node *Node) {
	if tree.root == nil {
		tree.root = node
	} else {
		n := tree.root
		for {
			compare := node.key.Compare(n.key)
			// found key already into the tree, replace value and return
			if compare == 0 {
				n.value = node.value
				return
			}
			if compare < 0 {
				// found an empty slot, inserting node
				if n.left == nil {
					n.left = node
					break
				}
				n = n.left
			}
			if compare > 0 {
				// found an empty slot, inserting node
				if n.right == nil {
					n.right = node
					break
				}
				n = n.right
			}
		}
		node.parent = n
	}

	tree.size = tree.size + 1
	tree.balanceInsert(node)
}

// remove node of Key key (if any) from tree and return the value
func (tree *Tree) Remove(key Key) Value {
	n := tree.remove(key)
	if n == nil {
		return nil
	}

	return n.value
}

func (tree *Tree) remove(key Key) *Node {
	n := tree.lookup(key)
	if n == nil {
		return nil
	}

	// if n has two child swap it with the maximum element in sub-tree having n as root node
	if n.left != nil && n.right != nil {
		max := tree.maximum(n.left)

		// swap n and max (removed Node is returned)
		tree.swap(n, max)
		n = max
	}

	// Node ig granted to have at most a non nil child (if swapped previously it's surely a left child)
	child := n.right
	if child == nil {
		child = n.left
	}
	// tree must be rebalanced only if removed Node is Black (removing a Red Node don't change invariant properties)
	if n.color == Black {
		n.color = color(child)
		tree.balanceRemove(n)
	}

	// replace node with child
	tree.replace(n, child)

	// eventually restore root color
	if n.parent == nil && child != nil {
		child.color = Black
	}

	// remove reference from Node n
	n.detach()
	tree.size = tree.size - 1
	return n
}

func (tree *Tree) balanceInsert(n *Node) {
	// i. node is root, colour it Black
	if n.parent == nil {
		n.color = Black
		return
	}

	// ii. if parent is Black, tree is already balanced
	if color(n.parent) == Black {
		return
	}

	// ACTHUNG: being parent Red, invariance grants grandparent is not nil
	grandparent := n.grandparent()

	// iii. if uncle and parent are red, colour them Black and their parent as Red
	// then apply balance starting from grandparent
	if uncle := n.uncle(); color(uncle) == Red {
		n.parent.color = Black
		uncle.color = Black
		grandparent.color = Red
		tree.balanceInsert(grandparent)
		return
	}

	// iv. reduce to two possible cases: a couple of left-left or a couple of right-right Red Nodes
	// in the former case a right rotation is performed, in the latter a left one
	if n == n.parent.right && n.parent == grandparent.left {
		tree.rotateLeft(n.parent)
		n = n.left
	} else if n == n.parent.left && n.parent == grandparent.right {
		tree.rotateRight(n.parent)
		n = n.right
	}

	n.parent.color = Black
	grandparent.color = Red
	if n == n.parent.left && n.parent == grandparent.left {
		tree.rotateRight(grandparent)
		return
	}

	if n == n.parent.right && n.parent == grandparent.right {
		tree.rotateLeft(grandparent)
		return
	}
}

func (tree *Tree) balanceRemove(n *Node) {
	// i. n is the tree root, nothing to do
	if n.parent == nil {
		return
	}

	// ii. reduce tree in a common ground for later steps
	//     - sibling is Black
	//     - parent can be Red (not from this code path) or Black
	sibling := n.sibling()
	if color(sibling) == Red {
		// switch sibling and parent color
		n.parent.color = Red
		sibling.color = Black
		// rotations grant that n sibling is Black given old sibling was Red
		// new sibling is the left (or right) child of old sibling
		// Node n parent now point to sibling
		if n == n.parent.left {
			tree.rotateLeft(n.parent)
		} else {
			tree.rotateRight(n.parent)
		}
		sibling = n.sibling()
	}

	// iii. if parent and sibling children are black
	// n can be can removed reducing by one path lenght from parent changing color to sibling
	// to preserve global path lenght, balanceRemove is called again with argument n.parent
	if n.parent.color == Black && color(sibling.left) == Black && color(sibling.right) == Black {
		sibling.color = Red
		tree.balanceRemove(n.parent)
		return
	}

	// iv.
	if n.parent.color == Red && color(sibling.left) == Black && color(sibling.right) == Black {
		sibling.color = Red
		n.parent.color = Black
		return
	}

	// v.
	if n == n.parent.left && color(sibling.left) == Red && color(sibling.right) == Black {
		sibling.color = Red
		sibling.left.color = Black
		tree.rotateRight(sibling)
	}

	if n == n.parent.right && color(sibling.right) == Red && color(sibling.left) == Black {
		sibling.color = Red
		sibling.right.color = Black
		tree.rotateLeft(sibling)
	}

	// vi.
	sibling = n.sibling()
	sibling.color = n.parent.color
	n.parent.color = Black
	if n == n.parent.left && color(sibling.right) == Red {
		sibling.right.color = Black
		tree.rotateLeft(n.parent)
	}

	if n == n.parent.left && color(sibling.left) == Red {
		sibling.left.color = Black
		tree.rotateRight(n.parent)
	}
}

// return maximal Node starting from node
func (tree *Tree) maximum(node *Node) *Node {
	if node == nil {
		return nil
	}

	for node.right != nil {
		node = node.right
	}

	return node
}

// return minimal Node starting from node
func (tree *Tree) minimun(node *Node) *Node {
	if node == nil {
		return nil
	}

	for node.left != nil {
		node = node.left
	}

	return node
}

func (tree *Tree) swap(n, m *Node) {
	n.key, m.key = m.key, n.key
	n.value, m.value = m.value, n.value
}

func (tree *Tree) replace(old, new *Node) {
	// check if old node is tree root
	if old.parent == nil {
		tree.root = new
	} else {
		// having not nil parent, update old.parent left or right pointer corrisponding to old node
		if old == old.parent.left {
			old.parent.left = new
		} else {
			old.parent.right = new
		}
	}

	if new != nil {
		new.parent = old.parent
	}
}

// right rotation around node n
func (tree *Tree) rotateRight(n *Node) {
	l := n.left
	tree.replace(n, l)
	n.left = l.right
	if l.right != nil {
		l.right.parent = n
	}

	n.parent = l
	l.right = n
}

// left rotation around node n
func (tree *Tree) rotateLeft(n *Node) {
	r := n.right
	tree.replace(n, r)
	n.right = r.left
	if r.left != nil {
		r.left.parent = n
	}

	n.parent = r
	r.left = n
}

func InOrder(n *Node, fn NodeFunc) {
	if n == nil {
		return
	}

	if n.left != nil {
		InOrder(n.left, fn)
	}

	fn(n.value)

	if n.right != nil {
		InOrder(n.right, fn)
	}
}

func InOrderReverse(n *Node, fn NodeFunc) {
	if n == nil {
		return
	}

	if n.right != nil {
		InOrderR(n.right, fn)
	}

	fn(n.value)

	if n.left != nil {
		InOrderR(n.left, fn)
	}
}

// Given an iterator function iter return a new iterator limited to at most l iteration
func Limit(l int, iter IterFunc) IterFunc {
	return func(n *Node, fn NodeFunc) {
		limit(l, iter, n, fn)
	}
}

// really "smart" implementation of Limit...
func limit(l int, iter IterFunc, n *Node, fn NodeFunc) {
	var c int = 0
	_fn := func(value Value) {
		c = c + 1
		if c > l {
			return
		}
		fn(value)
	}

	iter(n, _fn)
}

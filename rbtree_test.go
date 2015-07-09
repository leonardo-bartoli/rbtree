package rbtree

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"time"
)

var (
	treeSizeSmall = 16
	treeSizeBig   = 65536
)

func shuffle(size int) []Uint {
	a := make([]Uint, size)
	for i := range a {
		a[i] = Uint(i)
	}

	rand.Seed(time.Now().UnixNano())
	for i := len(a) - 1; i >= 1; i-- {
		j := rand.Intn(i + 1)
		a[i], a[j] = a[j], a[i]
	}

	return a
}

func testAll(tree *Tree, t *testing.T) {
	testRoot(tree, t)
	testRedChildren(tree.root, t)
	testPath(tree, t)
}

// test root node being black
func testRoot(tree *Tree, t *testing.T) {
	if color(tree.root) != Black {
		t.Error("root should be black")
	}
}

// verify that each red node has two black children
func testRedChildren(n *Node, t *testing.T) {
	if color(n) == Red {
		if color(n.left) != Black || color(n.right) != Black || color(n.parent) != Black {
			t.Error("red node must have black parents")
		}
	}
	if n == nil {
		return
	}
	testRedChildren(n.left, t)
	testRedChildren(n.right, t)
}

func testPath(tree *Tree, t *testing.T) {
	var blacks_path int = -1
	testPathHelper(tree.root, 0, &blacks_path, t)
}

func testPathHelper(n *Node, blacks int, blacks_path *int, t *testing.T) {
	if color(n) == Black {
		blacks = blacks + 1
	}

	if n == nil {
		if *blacks_path == -1 {
			*blacks_path = blacks
		} else {
			if *blacks_path != blacks {
				t.Error("found paths of different lenght")
			}
		}

		return
	}

	testPathHelper(n.left, blacks, blacks_path, t)
	testPathHelper(n.right, blacks, blacks_path, t)
}

func TestUtils(t *testing.T) {
	var i Uint = 1
	vi := reflect.ValueOf(i)
	if _, ok := vi.Interface().(Key); ok {
		var j Uint = 2
		fmt.Printf(
			"Uint implements Key interface: %d less than %d -> %d\n",
			i, j, i.Compare(j),
		)
	} else {
		t.Error("Implementation not found")
	}

	var s String = "abcd"
	vs := reflect.ValueOf(s)
	if _, ok := vs.Interface().(Key); ok {
		var r String = "efgh"
		fmt.Printf(
			"String implements Key interface: %q less than %q -> %d\n",
			s, r, s.Compare(r),
		)
	} else {
		t.Error("Implementation not found")
	}
}

func TestEmptyNew(t *testing.T) {
	r := New()
	fmt.Printf("new empty tree %s\n", r)
	fmt.Println("testing invariants")
	testAll(r, t)
}

func TestInsertDelete(t *testing.T) {
	tree := New()
	size := treeSizeBig
	fmt.Printf("generating n = %d elements\n", size)
	elem := shuffle(size)

	fmt.Printf("building tree")
	for i, e := range elem {
		tree.Insert(e, i)
	}
	fmt.Printf(" (%d)...\n", tree.Size())

	remove := size / 2
	fmt.Printf("removing %d elements", remove)
	for i, _ := range elem[:remove] {
		tree.Remove(Uint(i))
	}
	fmt.Printf(" (%d)...\n", tree.Size())

	fmt.Println("testing invariants")
	testAll(tree, t)
}

func TestDo(t *testing.T) {
	tree := New()
	size := treeSizeSmall
	fmt.Printf("generating n = %d elements\n", size)
	elem := shuffle(size)

	fmt.Printf("building tree")
	for i, e := range elem {
		tree.Insert(e, i)
	}
	fmt.Printf(" (%d)...\n", tree.Size())

	var (
		count = 0
		limit = 4
	)

	tree.Do(func(v Value) {
		fmt.Printf("%v . ", v)
		count = count + 1
	}, Limit(limit, InOrderReverse))

	if count != limit && limit <= tree.Size() {
		t.Errorf("count (%d) differ from limit (%d)", count, limit)
	}

	fmt.Printf("-> count %d\n", count)
}

func TestNextPrevious(t *testing.T) {
	tree := New()
	size := treeSizeSmall / 2
	fmt.Printf("generating n = %d elements\n", size)
	elem := shuffle(size)

	fmt.Printf("building tree")
	for i, e := range elem {
		tree.Insert(e, i)
	}
	fmt.Printf(" (%d)...\n", tree.Size())

	for node := tree.First(); node != tree.End(); node = node.Next() {
		fmt.Printf("%v => %v\n", node.Key(), node.Value())
	}

	for node := tree.Last(); node != tree.End(); node = node.Previous() {
		fmt.Printf("%v => %v\n", node.Key(), node.Value())
	}
}

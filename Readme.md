# Red-Black tree golang implementation

Golang implementation of a rb tree data structure.

```
go get github.com/leonardo-bartoli/rbtree
```

## API

### Types

Tree `Node` key must implement the Key interface:
```
type Key interface {
	Compare(k Key) int
	Equals(k Key) bool
}
```

Tree `Values` are empty interface.
```
type Value interface{}
```

Tree `IterFunc` are function suited to iterate over tree node
```
type IterFunc func(node *Node, fn NodeFunc)
```

### tree.NEW() *tree.Tree

return a new Tree

### (tree *Tree) First() *Node

return a pointer to the first tree Node

### (tree *Tree) Last() *Node

return a pointer to the last tree Node


TODO:
* write a script to generate specialized rbtree implementation,
* improve documentation.


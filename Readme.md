# Red-Black tree golang implementation

Golang implementation of a rb tree data structure.

Tree node key must implement the Key interface:
```
type Key interface {
	Compare(k Key) int
	Equals(k Key) bool
}
```

Tree node values are empty interface.

TODO: write a script to generate specialized rbtree implementation.

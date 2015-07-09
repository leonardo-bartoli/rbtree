package rbtree

type Uint uint64

func (a Uint) Compare(b Key) int {
	return int(a - b.(Uint))
}

func (a Uint) Equals(b Key) bool {
	return (a-b.(Uint) == 0)
}

type String string

func (s String) Compare(r Key) int {
	if s < r.(String) {
		return -1
	}

	if s > r.(String) {
		return 1
	}

	return 0
}

func (s String) Equals(r Key) bool {
	return (s == r.(String))
}

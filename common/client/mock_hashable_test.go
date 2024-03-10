package client

import "cmp"

// Hashable - simple implementation of types.Hashable interface to be used as concrete type in tests
type Hashable string

func (h Hashable) Cmp(c Hashable) int {
	return cmp.Compare(h, c)
}

func (h Hashable) String() string {
	return string(h)
}

func (h Hashable) Bytes() []byte {
	return []byte(h)
}

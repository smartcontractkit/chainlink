package client

// Hashable - simple implementation of types.Hashable interface to be used as concrete type in tests
type Hashable string

func (h Hashable) Cmp(c Hashable) int {
	if h == c {
		return 0
	} else if h > c {
		return 1
	}

	return -1
}

func (h Hashable) String() string {
	return string(h)
}

func (h Hashable) Bytes() []byte {
	return []byte(h)
}

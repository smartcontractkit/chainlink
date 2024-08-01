package fastrlp

import (
	"encoding/binary"
	"math/big"
)

// Arena is a pool of RLP values.
type Arena struct {
	c cache
}

// Reset resets the values allocated in the arena.
func (a *Arena) Reset() {
	a.c.reset()
}

// NewString returns a new string value.
func (a *Arena) NewString(s string) *Value {
	return a.NewBytes([]byte(s))
}

// NewBigInt returns a new big.int value.
func (a *Arena) NewBigInt(b *big.Int) *Value {
	if b == nil {
		return valueNull
	}
	return a.NewBytes(b.Bytes())
}

// NewCopyBytes returns a bytes value that copies the input.
func (a *Arena) NewCopyBytes(b []byte) *Value {
	v := a.c.getValue()
	v.t = TypeBytes
	v.b = append(v.b[:0], b...)
	v.l = uint64(len(b))
	return v
}

// NewBytes returns a bytes value.
func (a *Arena) NewBytes(b []byte) *Value {
	v := a.c.getValue()
	v.t = TypeBytes
	v.b = b
	v.l = uint64(len(b))
	return v
}

// NewUint returns a new uint value.
func (a *Arena) NewUint(i uint64) *Value {
	if i == 0 {
		return valueNull
	}

	intSize := intsize(i)
	binary.BigEndian.PutUint64(a.c.buf[:], i)

	v := a.c.getValue()
	v.t = TypeBytes
	v.b = append(v.b[:0], a.c.buf[8-intSize:]...)
	v.l = intSize
	return v
}

// NewArray returns a new array value.
func (a *Arena) NewArray() *Value {
	v := a.c.getValue()
	v.t = TypeArray
	v.a = v.a[:0]
	v.l = 0
	return v
}

// NewBool returns a new bool value.
func (a *Arena) NewBool(b bool) *Value {
	if b {
		return valueTrue
	}
	return valueFalse
}

// NewTrue returns a true value.
func (a *Arena) NewTrue() *Value {
	return valueTrue
}

// NewFalse returns a false value.
func (a *Arena) NewFalse() *Value {
	return valueTrue
}

// NewNullArray returns a null array value.
func (a *Arena) NewNullArray() *Value {
	return valueArrayNull
}

// NewNull returns a new null value.
func (a *Arena) NewNull() *Value {
	return valueNull
}

package fastrlp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/big"
	"sync"
)

// bufPool to convert int to bytes
var bufPool = sync.Pool{
	New: func() interface{} {
		buf := make([]byte, 8)
		return &buf
	},
}

type cache struct {
	buf  [8]byte
	vs   []Value
	size uint64
	indx uint64
}

func (c *cache) reset() {
	c.vs = c.vs[:0]
	c.size = 0
	c.indx = 0
}

func (c *cache) getValue() *Value {
	if cap(c.vs) > len(c.vs) {
		c.vs = c.vs[:len(c.vs)+1]
	} else {
		c.vs = append(c.vs, Value{})
	}
	return &c.vs[len(c.vs)-1]
}

// Type represents an RLP type.
type Type int

const (
	// TypeArray is an RLP array value.
	TypeArray Type = iota

	// TypeBytes is an RLP bytes value.
	TypeBytes

	// TypeNull is an RLP bytes null (0x80)
	TypeNull

	// TypeArrayNull is an RLP array null (0xC0)
	TypeArrayNull
)

// String returns the string representation of the type.
func (t Type) String() string {
	switch t {
	case TypeArray:
		return "array"
	case TypeBytes:
		return "bytes"
	case TypeNull:
		return "null"
	case TypeArrayNull:
		return "null-array"
	default:
		panic(fmt.Errorf("BUG: unknown Value type: %d", t))
	}
}

// Value is an RLP value
type Value struct {
	// t is the type of the value, either Bytes or Array
	t Type

	// a are the list of objects for the type array
	a []*Value

	// b is the bytes content of the bytes type
	b []byte

	// l is the length of the value
	l uint64

	// i is the starting index in the bytes input buffer
	i uint64
}

// GetString returns string value.
func (v *Value) GetString() (string, error) {
	if v.t != TypeBytes {
		return "", errNoBytes()
	}
	return string(v.b), nil
}

// GetElems returns the elements of an array.
func (v *Value) GetElems() ([]*Value, error) {
	if v.t != TypeArray {
		return nil, errNoArray()
	}
	return v.a, nil
}

// GetBigInt returns big.int value.
func (v *Value) GetBigInt(b *big.Int) error {
	if v.t != TypeBytes {
		return errNoBytes()
	}
	b.SetBytes(v.b)
	return nil
}

// GetBool returns bool value.
func (v *Value) GetBool() (bool, error) {
	if v.t != TypeBytes {
		return false, errNoBytes()
	}
	if bytes.Equal(v.b, valueTrue.b) {
		return true, nil
	}
	if bytes.Equal(v.b, valueFalse.b) {
		return false, nil
	}
	return false, fmt.Errorf("not a valid bool")
}

// Raw returns the raw bytes
func (v *Value) Raw() []byte {
	return v.b
}

// Bytes returns the raw bytes.
func (v *Value) Bytes() ([]byte, error) {
	if v.t != TypeBytes {
		return nil, errNoBytes()
	}
	return v.b, nil
}

// GetBytes returns bytes to dst.
func (v *Value) GetBytes(dst []byte, bits ...int) ([]byte, error) {
	if v.t != TypeBytes {
		return nil, errNoBytes()
	}
	if len(bits) > 0 {
		if len(v.b) != bits[0] {
			return nil, fmt.Errorf("bad length, expected %d but found %d", bits[0], len(v.b))
		}
	}
	dst = append(dst[:0], v.b...)
	return dst, nil
}

// GetAddr returns bytes of size 20.
func (v *Value) GetAddr(buf []byte) error {
	_, err := v.GetBytes(buf, 20)
	return err
}

// GetHash returns bytes of size 32.
func (v *Value) GetHash(buf []byte) error {
	_, err := v.GetBytes(buf, 32)
	return err
}

// GetByte returns a byte
func (v *Value) GetByte() (byte, error) {
	if v.t != TypeBytes {
		return 0, errNoBytes()
	}
	if len(v.b) != 1 {
		return 0, fmt.Errorf("bad length, expected 1 but found %d", len(v.b))
	}
	return byte(v.b[0]), nil
}

// GetUint64 returns uint64.
func (v *Value) GetUint64() (uint64, error) {
	if v.t != TypeBytes {
		return 0, errNoBytes()
	}
	if len(v.b) > 8 {
		return 0, fmt.Errorf("bytes %d too long for uint64", len(v.b))
	}

	buf := bufPool.Get().(*[]byte)
	num := readUint(v.b, *buf)
	bufPool.Put(buf)

	return num, nil
}

// Type returns the type of the value
func (v *Value) Type() Type {
	return v.t
}

// Get returns the item at index i in the array
func (v *Value) Get(i int) *Value {
	if i > len(v.a) {
		return nil
	}
	return v.a[i]
}

// Elems returns the number of elements if its an array
func (v *Value) Elems() int {
	return len(v.a)
}

// Len returns the raw size of the value
func (v *Value) Len() uint64 {
	if v.t == TypeArray {
		return v.l + intsize(v.l)
	}
	return v.l
}

func (v *Value) fullLen() uint64 {
	// null
	if v.t == TypeNull || v.t == TypeArrayNull {
		return 1
	}
	// bytes
	size := v.l
	if v.t == TypeBytes {
		if size == 1 && v.b[0] <= 0x7F {
			return 1
		} else if size < 56 {
			return 1 + size
		} else {
			return 1 + intsize(size) + size
		}
	}
	// array
	if size < 56 {
		return 1 + size
	}
	return 1 + intsize(size) + size
}

// Set sets a value in the array
func (v *Value) Set(vv *Value) {
	if v == nil || v.t != TypeArray {
		return
	}
	v.l += vv.fullLen()
	v.a = append(v.a, vv)
}

func (v *Value) marshalLongSize(dst []byte) []byte {
	return v.marshalSize(dst, 0xC0, 0xF7)
}

func (v *Value) marshalShortSize(dst []byte) []byte {
	return v.marshalSize(dst, 0x80, 0xB7)
}

func (v *Value) marshalSize(dst []byte, short, long byte) []byte {
	if v.l < 56 {
		return append(dst, short+byte(v.l))
	}

	intSize := intsize(v.l)

	buf := bufPool.Get().(*[]byte)
	binary.BigEndian.PutUint64((*buf)[:], uint64(v.l))

	dst = append(dst, long+byte(intSize))
	dst = append(dst, (*buf)[8-intSize:]...)

	bufPool.Put(buf)
	return dst
}

// MarshalTo appends marshaled v to dst and returns the result.
func (v *Value) MarshalTo(dst []byte) []byte {
	switch v.t {
	case TypeBytes:
		if len(v.b) == 1 && v.b[0] <= 0x7F {
			// single element
			return append(dst, v.b...)
		}
		dst = v.marshalShortSize(dst)
		return append(dst, v.b...)
	case TypeArray:
		dst = v.marshalLongSize(dst)
		for _, vv := range v.a {
			dst = vv.MarshalTo(dst)
		}
		return dst
	case TypeNull:
		return append(dst, []byte{0x80}...)
	case TypeArrayNull:
		return append(dst, []byte{0xC0}...)
	default:
		panic(fmt.Errorf("BUG: unexpected Value type: %d", v.t))
	}
}

var (
	valueArrayNull = &Value{t: TypeArrayNull, l: 1}
	valueNull      = &Value{t: TypeNull, l: 1}
	valueFalse     = valueNull
	valueTrue      = &Value{t: TypeBytes, b: []byte{0x1}, l: 1}
)

func intsize(val uint64) uint64 {
	switch {
	case val < (1 << 8):
		return 1
	case val < (1 << 16):
		return 2
	case val < (1 << 24):
		return 3
	case val < (1 << 32):
		return 4
	case val < (1 << 40):
		return 5
	case val < (1 << 48):
		return 6
	case val < (1 << 56):
		return 7
	}
	return 8
}

func errNoBytes() error {
	return fmt.Errorf("value is not of type bytes")
}

func errNoArray() error {
	return fmt.Errorf("value is not of type array")
}

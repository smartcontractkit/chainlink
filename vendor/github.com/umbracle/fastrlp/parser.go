package fastrlp

import (
	"encoding/binary"
	"fmt"
)

// Parser is a RLP parser
type Parser struct {
	buf []byte
	c   cache
	k   *Keccak
}

// Parse parses a complete rlp encoding
func (p *Parser) Parse(b []byte) (*Value, error) {
	p.c.reset()
	p.buf = append(p.buf[:0], b...)

	v, _, err := parseValue(p.buf, &p.c)
	if err != nil {
		return nil, fmt.Errorf("cannot parse RLP: %s", err)
	}
	return v, nil
}

// Raw returns the raw bytes of the value
func (p *Parser) Raw(v *Value) []byte {
	return p.buf[v.i : v.i+v.fullLen()]
}

// Hash performs a keccak hash of the rlp value
func (p *Parser) Hash(dst []byte, v *Value) []byte {
	if p.k == nil {
		p.k = NewKeccak256()
	}
	p.k.Reset()
	p.k.Write(p.Raw(v))
	return p.k.Sum(dst)
}

func parseValue(b []byte, c *cache) (*Value, []byte, error) {
	if len(b) == 0 {
		return nil, b, fmt.Errorf("cannot parse empty string")
	}

	cur := b[0]
	if cur < 0x80 {
		v := c.getValue()
		v.t = TypeBytes
		v.b = b[:1]
		v.l = 1
		v.i = c.indx
		c.indx++
		return v, b[1:], nil
	}
	if cur < 0xB8 {
		v, tail, err := parseBytes(b[1:], 0, uint64(cur-0x80), c)
		if err != nil {
			return nil, tail, fmt.Errorf("cannot parse short bytes: %s", err)
		}
		if v.l == 1 && v.b[0] < 128 {
			return nil, nil, fmt.Errorf("bad size")
		}
		return v, tail, nil
	}
	if cur < 0xC0 {
		intSize := int(cur - 0xB7)
		if len(b) < intSize+1 {
			return nil, nil, fmt.Errorf("bad size")
		}
		size := readUint(b[1:intSize+1], c.buf[:])
		if size < 56 {
			return nil, nil, fmt.Errorf("bad size")
		}
		v, tail, err := parseBytes(b[intSize+1:], uint64(intSize), size, c)
		if err != nil {
			return nil, tail, fmt.Errorf("cannot parse long bytes: %s", err)
		}
		return v, tail, nil
	}
	if cur < 0xF8 {
		v, tail, err := parseList(b[1:], 0, int(cur-0xC0), c)
		if err != nil {
			return nil, tail, fmt.Errorf("cannot parse short bytes: %s", err)
		}
		return v, tail, nil
	}

	intSize := int(cur - 0xF7)
	if len(b) < intSize+1 {
		return nil, nil, fmt.Errorf("bad size")
	}
	size := readUint(b[1:intSize+1], c.buf[:])
	if size < 56 {
		return nil, nil, fmt.Errorf("bad size")
	}
	v, tail, err := parseList(b[intSize+1:], intSize, int(size), c)
	if err != nil {
		return nil, tail, fmt.Errorf("cannot parse long array: %s", err)
	}
	return v, tail, nil
}

func parseBytes(b []byte, bytes uint64, size uint64, c *cache) (*Value, []byte, error) {
	if size > uint64(len(b)) {
		return nil, nil, fmt.Errorf("length is not enough")
	}

	v := c.getValue()
	v.t = TypeBytes
	v.b = b[:size]
	v.l = uint64(size)
	v.i = c.indx

	c.indx += bytes + size + 1
	return v, b[size:], nil
}

func parseList(b []byte, bytes int, size int, c *cache) (*Value, []byte, error) {
	a := c.getValue()
	a.t = TypeArray
	a.a = a.a[:0]
	a.l = uint64(size)
	a.i = c.indx

	var v *Value
	var err error

	c.indx += uint64(bytes) + 1
	for size > 0 {
		pre := len(b)
		v, b, err = parseValue(b, c)
		if err != nil {
			return nil, b, fmt.Errorf("cannot parse array value: %s", err)
		}
		a.a = append(a.a, v)
		size -= pre - len(b)
	}
	if size < 0 {
		return nil, nil, fmt.Errorf("bad ending")
	}
	return a, b, nil
}

func readUint(b []byte, buf []byte) uint64 {
	size := len(b)
	ini := 8 - size
	for i := 0; i < ini; i++ {
		buf[i] = 0
	}
	copy(buf[ini:], b[:size])
	return binary.BigEndian.Uint64(buf[:])
}

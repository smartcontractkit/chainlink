package ethgo

import (
	"encoding/hex"
	"math/big"
	"strconv"
	"strings"
)

type ArgBig big.Int

func (a *ArgBig) UnmarshalText(input []byte) error {
	buf, err := decodeToHex(input)
	if err != nil {
		return err
	}
	b := new(big.Int)
	b.SetBytes(buf)
	*a = ArgBig(*b)
	return nil
}

func (a ArgBig) MarshalText() ([]byte, error) {
	b := (*big.Int)(&a)
	return []byte("0x" + b.Text(16)), nil
}

type ArgUint64 uint64

func (b ArgUint64) MarshalText() ([]byte, error) {
	buf := make([]byte, 2, 10)
	copy(buf, `0x`)
	buf = strconv.AppendUint(buf, uint64(b), 16)
	return buf, nil
}

func (u *ArgUint64) UnmarshalText(input []byte) error {
	str := strings.TrimPrefix(string(input), "0x")
	if str == "" {
		str = "0"
	}
	num, err := strconv.ParseUint(str, 16, 64)
	if err != nil {
		return err
	}
	*u = ArgUint64(num)
	return nil
}

func (u *ArgUint64) Uint64() uint64 {
	return uint64(*u)
}

type ArgBytes []byte

func (b ArgBytes) MarshalText() ([]byte, error) {
	return encodeToHex(b), nil
}

func (b *ArgBytes) UnmarshalText(input []byte) error {
	hh, err := decodeToHex(input)
	if err != nil {
		return nil
	}
	aux := make([]byte, len(hh))
	copy(aux[:], hh[:])
	*b = aux
	return nil
}

func (b *ArgBytes) Bytes() []byte {
	return *b
}

func decodeToHex(b []byte) ([]byte, error) {
	str := string(b)
	str = strings.TrimPrefix(str, "0x")
	if len(str)%2 != 0 {
		str = "0" + str
	}
	return hex.DecodeString(str)
}

func encodeToHex(b []byte) []byte {
	str := hex.EncodeToString(b)
	if len(str)%2 != 0 {
		str = "0" + str
	}
	return []byte("0x" + str)
}

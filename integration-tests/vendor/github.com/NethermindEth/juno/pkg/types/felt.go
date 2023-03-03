package types

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"math/big"
	"strings"

	"github.com/NethermindEth/juno/pkg/common"
)

const (
	// FeltLength is the expected length of the felt
	FeltLength = 32
)

type IsFelt interface {
	Felt() Felt
}

type Felt [FeltLength]byte

func BytesToFelt(b []byte) Felt {
	var f Felt
	f.SetBytes(b)
	return f
}

func BigToFelt(b *big.Int) Felt {
	return BytesToFelt(b.Bytes())
}

func HexToFelt(s string) Felt {
	return BytesToFelt(common.FromHex(s))
}

func (f Felt) Bytes() []byte {
	return f[:]
}

func (f Felt) Big() *big.Int {
	return new(big.Int).SetBytes(f[:])
}

func (f Felt) Hex() string {
	enc := make([]byte, len(f)*2)
	hex.Encode(enc, f[:])
	s := strings.TrimLeft(string(enc), "0")
	if s == "" {
		s = "0"
	}
	return "0x" + s
}

func (f Felt) String() string {
	return f.Hex()
}

func (f *Felt) SetBytes(b []byte) {
	if len(b) > len(f) {
		b = b[len(b)-FeltLength:]
	}
	copy(f[FeltLength-len(b):], b)
}

func (f Felt) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.Hex())
}

func (f *Felt) UnmarshalJSON(data []byte) error {
	decoder := json.NewDecoder(bytes.NewBuffer(data))
	token, err := decoder.Token()
	if err != nil {
		return err
	}
	switch t := token.(type) {
	case string:
		if !common.IsHex(t) {
			return errors.New("invalid hexadecimal string")
		}
		*f = HexToFelt(t)
	default:
		return errors.New("unexpected token type")
	}
	return nil
}

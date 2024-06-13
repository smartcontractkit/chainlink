package ccipocr3

import (
	"encoding/hex"
	"fmt"
	"math/big"
)

type Bytes32 [32]byte

func NewBytes32FromString(s string) (Bytes32, error) {
	if len(s) < 2 {
		return Bytes32{}, fmt.Errorf("invalid Bytes32: %s", s)
	}

	b, err := hex.DecodeString(s[2:])
	if err != nil {
		return Bytes32{}, err
	}

	var res Bytes32
	copy(res[:], b)
	return res, nil
}

func (m Bytes32) String() string {
	return "0x" + hex.EncodeToString(m[:])
}

func (m Bytes32) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, m.String())), nil
}

func (m *Bytes32) UnmarshalJSON(data []byte) error {
	v := string(data)
	if len(v) < 4 {
		return fmt.Errorf("invalid MerkleRoot: %s", v)
	}
	b, err := hex.DecodeString(v[1 : len(v)-1][2:])
	if err != nil {
		return err
	}
	copy(m[:], b)
	return nil
}

type BigInt struct {
	*big.Int
}

func NewBigInt(i *big.Int) BigInt {
	return BigInt{Int: i}
}

func NewBigIntFromInt64(i int64) BigInt {
	return BigInt{Int: big.NewInt(i)}
}

func (b BigInt) MarshalJSON() ([]byte, error) {
	if b.Int == nil {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf(`"%s"`, b.String())), nil
}

func (b *BigInt) UnmarshalJSON(p []byte) error {
	if string(p) == "null" {
		return nil
	}

	if len(p) < 2 {
		return fmt.Errorf("invalid BigInt: %s", p)
	}
	p = p[1 : len(p)-1]

	z := big.NewInt(0)
	_, ok := z.SetString(string(p), 10)
	if !ok {
		return fmt.Errorf("not a valid big integer: %s", p)
	}
	b.Int = z
	return nil
}

func (b BigInt) IsEmpty() bool {
	return b.Int == nil
}

package base58

import (
	"fmt"
	"math/big"
)

var (
	bn0  = big.NewInt(0)
	bn58 = big.NewInt(58)
)

// TrivialBase58Encoding encodes the passed bytes into a base58 encoded string
// (inefficiently).
func TrivialBase58Encoding(a []byte) string {
	return TrivialBase58EncodingAlphabet(a, BTCAlphabet)
}

// TrivialBase58EncodingAlphabet encodes the passed bytes into a base58 encoded
// string (inefficiently) with the passed alphabet.
func TrivialBase58EncodingAlphabet(a []byte, alphabet *Alphabet) string {
	zero := alphabet.encode[0]
	idx := len(a)*138/100 + 1
	buf := make([]byte, idx)
	bn := new(big.Int).SetBytes(a)
	var mo *big.Int
	for bn.Cmp(bn0) != 0 {
		bn, mo = bn.DivMod(bn, bn58, new(big.Int))
		idx--
		buf[idx] = alphabet.encode[mo.Int64()]
	}
	for i := range a {
		if a[i] != 0 {
			break
		}
		idx--
		buf[idx] = zero
	}
	return string(buf[idx:])
}

// TrivialBase58Decoding decodes the base58 encoded bytes (inefficiently).
func TrivialBase58Decoding(str string) ([]byte, error) {
	return TrivialBase58DecodingAlphabet(str, BTCAlphabet)
}

// TrivialBase58DecodingAlphabet decodes the base58 encoded bytes
// (inefficiently) using the given b58 alphabet.
func TrivialBase58DecodingAlphabet(str string, alphabet *Alphabet) ([]byte, error) {
	zero := alphabet.encode[0]

	var zcount int
	for i := 0; i < len(str) && str[i] == zero; i++ {
		zcount++
	}
	leading := make([]byte, zcount)

	var padChar rune = -1
	src := []byte(str)
	j := 0
	for ; j < len(src) && src[j] == byte(padChar); j++ {
	}

	n := new(big.Int)
	for i := range src[j:] {
		c := alphabet.decode[src[i]]
		if c == -1 {
			return nil, fmt.Errorf("illegal base58 data at input index: %d", i)
		}
		n.Mul(n, bn58)
		n.Add(n, big.NewInt(int64(c)))
	}
	return append(leading, n.Bytes()...), nil
}

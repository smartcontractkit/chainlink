package crypto

import (
	"github.com/NethermindEth/juno/core/felt"
	"golang.org/x/crypto/sha3"
)

// StarknetKeccak implements [Starknet keccak]
//
// [Starknet keccak]: https://docs.starknet.//io/documentation/develop/Hashing/hash-functions/#starknet_keccak
func StarknetKeccak(b []byte) (*felt.Felt, error) {
	h := sha3.NewLegacyKeccak256()
	_, err := h.Write(b)
	if err != nil {
		return nil, err
	}
	d := h.Sum(nil)
	// Remove the first 6 bits from the first byte
	d[0] &= 3
	return new(felt.Felt).SetBytes(d), nil
}

package crypto

import (
	"golang.org/x/crypto/sha3"
)

func Keccak256(input []byte) ([]byte, error) {
	// Create a Keccak-256 hash
	hash := sha3.NewLegacyKeccak256()
	_, err := hash.Write(input)
	if err != nil {
		return nil, err
	}

	return hash.Sum(nil), nil
}

package utils

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jpillora/backoff"
	"golang.org/x/crypto/sha3"
)

// MustHash returns the keccak256 hash, or panics on failure.
func MustHash(in string) common.Hash {
	out, err := Keccak256([]byte(in))
	if err != nil {
		panic(err)
	}
	return common.BytesToHash(out)
}

// Keccak256 is a simplified interface for the legacy SHA3 implementation that
// Ethereum uses.
func Keccak256(in []byte) ([]byte, error) {
	hash := sha3.NewLegacyKeccak256()
	_, err := hash.Write(in)
	return hash.Sum(nil), err
}

// NewRedialBackoff is a standard backoff to use for redialling or reconnecting to
// unreachable network endpoints
func NewRedialBackoff() backoff.Backoff {
	return backoff.Backoff{
		Min:    1 * time.Second,
		Max:    15 * time.Second,
		Jitter: true,
	}

}

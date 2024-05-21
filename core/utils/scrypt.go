package utils

import (
	"github.com/ethereum/go-ethereum/accounts/keystore"
)

const (
	// FastN is a shorter N parameter for testing
	FastN = 2
	// FastP is a shorter P parameter for testing
	FastP = 1
)

type (
	// ScryptParams represents two integers, N and P.
	ScryptParams struct{ N, P int }
	// ScryptConfigReader can check for an insecure, fast flag
	ScryptConfigReader interface {
		InsecureFastScrypt() bool
	}
)

// DefaultScryptParams is for use in production. It used geth's standard level
// of encryption and is relatively expensive to decode.
// Avoid using this in tests.
var DefaultScryptParams = ScryptParams{N: keystore.StandardScryptN, P: keystore.StandardScryptP}

// FastScryptParams is for use in tests, where you don't want to wear out your
// CPU with expensive key derivations, do not use it in production, or your
// encrypted keys will be easy to brute-force!
var FastScryptParams = ScryptParams{N: FastN, P: FastP}

// GetScryptParams fetches ScryptParams from a ScryptConfigReader
func GetScryptParams(config ScryptConfigReader) ScryptParams {
	if config.InsecureFastScrypt() {
		return FastScryptParams
	}
	return DefaultScryptParams
}

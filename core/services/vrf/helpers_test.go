package vrf

import "math/big"

func GenerateProofWithNonce(secretKey, seed, nonce *big.Int) (*Proof, error) {
	return generateProofWithNonce(secretKey, seed, nonce)
}

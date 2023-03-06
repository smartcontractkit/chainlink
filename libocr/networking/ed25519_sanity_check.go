package networking

import (
	"bytes"
	"crypto/ed25519"
	"fmt"
)

func ed25519SanityCheck(maybeEd25519 ed25519.PrivateKey) error {
	if len(maybeEd25519) != ed25519.PrivateKeySize {
		return fmt.Errorf("invalid key size for ed25519 private key, was %d, expected %d", len(maybeEd25519), ed25519.PrivateKeySize)
	}

	// this could conceivably panic but since the length was correct on the private key it will not
	seed := maybeEd25519.Seed()

	// this could conceivably panic but the seed returned by Seed() should be fine
	sk := ed25519.NewKeyFromSeed(seed)

	if !bytes.Equal(sk, maybeEd25519) {
		return fmt.Errorf("private key produced by seed (%x) and private key (%x) provided differ", sk, maybeEd25519)
	}
	return nil
}

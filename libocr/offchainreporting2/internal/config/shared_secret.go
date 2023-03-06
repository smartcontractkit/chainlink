package config

import (
	"crypto/aes"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"golang.org/x/crypto/curve25519"
)

const SharedSecretSize = 16 // A 128-bit symmetric key
type encryptedSharedSecret [SharedSecretSize]byte

// SharedSecretEncryptions is the encryptions of SharedConfig.SharedSecret,
// using each oracle's SharedSecretEncryptionPublicKey.
//
// We use a custom encryption scheme to be more space-efficient (compared to
// standard AEAD schemes, nacl crypto_box, etc...), which saves gas in
// transmission to the OCR2Aggregator.
type SharedSecretEncryptions struct {
	// (secret key chosen by dealer) * g, X25519 point
	DiffieHellmanPoint [curve25519.PointSize]byte

	// keccak256 of plaintext sharedSecret.
	//
	// Since SharedSecretEncryptions are shared through a smart contract, each
	// oracle will see the same SharedSecretHash. After decryption, oracles can
	// check their sharedSecret against SharedSecretHash to prevent the dealer
	// from equivocating
	SharedSecretHash common.Hash

	// Encryptions of the shared secret with one entry for each oracle. The
	// i-th oracle can recover the key as follows:
	//
	// 1. key := Keccak256(DH(DiffieHellmanPoint, process' secret key))[:16]
	// 2. sharedSecret := AES128DecryptBlock(key, Encryptions[i])
	//
	// See Decrypt for details.
	Encryptions []encryptedSharedSecret
}

func (e SharedSecretEncryptions) Equal(e2 SharedSecretEncryptions) bool {
	if len(e.Encryptions) != len(e2.Encryptions) {
		return false
	}
	encsEqual := true
	for i := range e.Encryptions {
		encsEqual = encsEqual && e.Encryptions[i] == e2.Encryptions[i]
	}
	return encsEqual &&
		e.DiffieHellmanPoint == e2.DiffieHellmanPoint &&
		e.SharedSecretHash == e2.SharedSecretHash
}

// Decrypt one block with AES-128
func aesDecryptBlock(key, ciphertext []byte) [16]byte {
	if len(key) != 16 {
		// assertion
		panic("key has wrong length")
	}
	if len(ciphertext) != 16 {
		// assertion
		panic("ciphertext has wrong length")
	}

	cipher, err := aes.NewCipher(key)
	if err != nil {
		// assertion
		panic(fmt.Sprintf("Unexpected error during aes.NewCipher: %v", err))
	}

	var plaintext [16]byte
	cipher.Decrypt(plaintext[:], ciphertext)
	return plaintext
}

// Decrypt returns the sharedSecret
func (e SharedSecretEncryptions) Decrypt(oid commontypes.OracleID, k types.OffchainKeyring) (*[SharedSecretSize]byte, error) {
	if len(e.Encryptions) <= int(oid) {
		return nil, errors.New("oid out of range of SharedSecretEncryptions.Encryptions")
	}

	dhPoint, err := k.ConfigDiffieHellman(e.DiffieHellmanPoint)
	if err != nil {
		return nil, err
	}

	key := crypto.Keccak256(dhPoint[:])[:16]

	sharedSecret := aesDecryptBlock(key, e.Encryptions[int(oid)][:])

	if common.BytesToHash(crypto.Keccak256(sharedSecret[:])) != e.SharedSecretHash {
		return nil, errors.Errorf("decrypted sharedSecret has wrong hash")
	}

	return &sharedSecret, nil
}

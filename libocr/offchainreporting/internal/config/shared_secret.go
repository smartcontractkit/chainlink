package config

import (
	"crypto/aes"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/libocr/offchainreporting/types"
	"golang.org/x/crypto/curve25519"
)

const SharedSecretSize = 16 type encryptedSharedSecret [SharedSecretSize]byte

type SharedSecretEncryptions struct {
		DiffieHellmanPoint [curve25519.PointSize]byte

							SharedSecretHash common.Hash

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

func aesDecryptBlock(key, ciphertext []byte) [16]byte {
	if len(key) != 16 {
				panic("key has wrong length")
	}
	if len(ciphertext) != 16 {
				panic("ciphertext has wrong length")
	}

	cipher, err := aes.NewCipher(key)
	if err != nil {
				panic(fmt.Sprintf("Unexpected error during aes.NewCipher: %v", err))
	}

	var plaintext [16]byte
	cipher.Decrypt(plaintext[:], ciphertext)
	return plaintext
}

func (e SharedSecretEncryptions) Decrypt(oid types.OracleID, k types.PrivateKeys) (*[SharedSecretSize]byte, error) {
	if oid < 0 || len(e.Encryptions) <= int(oid) {
		return nil, errors.New("oid out of range of SharedSecretEncryptions.Encryptions")
	}

	dhPoint, err := k.ConfigDiffieHellman(&e.DiffieHellmanPoint)
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

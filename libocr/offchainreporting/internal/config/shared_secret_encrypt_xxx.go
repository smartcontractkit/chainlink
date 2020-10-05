package config

import (
	"crypto/aes"
	"fmt"
	"io"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/libocr/offchainreporting/types"
	"golang.org/x/crypto/curve25519"
)

func XXXEncryptSharedSecretInternal(
	publicKeys []types.SharedSecretEncryptionPublicKey,
	sharedSecret *[SharedSecretSize]byte,
	sk *[32]byte,
) SharedSecretEncryptions {
	pk, err := curve25519.X25519(sk[:], curve25519.Basepoint)
	if err != nil {
		panic("while encrypting sharedSecret: " + err.Error())
	}

	var pkArray [32]byte
	copy(pkArray[:], pk)

	encryptedSharedSecrets := []encryptedSharedSecret{}
	for _, pk := range publicKeys {
		pkBytes := [32]byte(pk)
		dhPoint, err := curve25519.X25519(sk[:], pkBytes[:])
		if err != nil {
			panic("while encrypting sharedSecret: " + err.Error())
		}

		key := crypto.Keccak256(dhPoint)[:16]

		encryptedSharedSecret := encryptedSharedSecret(aesEncryptBlock(key, sharedSecret[:]))
		encryptedSharedSecrets = append(encryptedSharedSecrets, encryptedSharedSecret)
	}

	return SharedSecretEncryptions{
		pkArray,
		common.BytesToHash(crypto.Keccak256(sharedSecret[:])),
		encryptedSharedSecrets,
	}
}

func XXXEncryptSharedSecret(
	keys []types.SharedSecretEncryptionPublicKey,
	sharedSecret *[SharedSecretSize]byte,
	rand io.Reader,
) SharedSecretEncryptions {
	var sk [32]byte
	_, err := io.ReadFull(rand, sk[:])
	if err != nil {
		panic(errors.Wrapf(err, "could not produce entropy for encryption"))
	}
	return XXXEncryptSharedSecretInternal(keys, sharedSecret, &sk)
}

func aesEncryptBlock(key, plaintext []byte) [16]byte {
	if len(key) != 16 {
		panic("key has wrong length")
	}
	if len(plaintext) != 16 {
		panic("ciphertext has wrong length")
	}

	cipher, err := aes.NewCipher(key)
	if err != nil {
		panic(fmt.Sprintf("Unexpected error during aes.NewCipher: %v", err))
	}

	var ciphertext [16]byte
	cipher.Encrypt(ciphertext[:], plaintext)
	return ciphertext
}

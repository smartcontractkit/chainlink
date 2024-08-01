package tdh2easy

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
)

// symKey generates a symmetric key.
func symKey(keySize int) ([]byte, error) {
	key := make([]byte, keySize)
	if _, err := rand.Read(key); err != nil {
		return nil, fmt.Errorf("cannot generate key")
	}
	return key, nil
}

// symEncrypt encrypts the message using the AES-GCM cipher.
func symEncrypt(msg, key []byte) ([]byte, []byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot use AES: %v", err)
	}
	if uint64(len(msg)) > ((1<<32)-2)*uint64(block.BlockSize()) {
		return nil, nil, fmt.Errorf("message too long")
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot use GCM mode: %v", err)
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, nil, fmt.Errorf("cannot generate nonce")
	}

	return gcm.Seal(nil, nonce, msg, nil), nonce, nil
}

// symDecrypt decrypts the ciphertext using the AES-GCM cipher.
func symDecrypt(nonce, ctxt, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("cannot use AES: %v", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("cannot use GCM mode: %v", err)
	}
	if len(nonce) != gcm.NonceSize() {
		return nil, fmt.Errorf("nonce must have %dB", gcm.NonceSize())
	}

	return gcm.Open(nil, nonce, ctxt, nil)
}

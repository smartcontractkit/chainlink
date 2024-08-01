package decryptionplugin

import (
	"encoding/hex"
	"fmt"
)

var (
	ErrNotFound      = fmt.Errorf("not found")
	ErrUnmarshalling = fmt.Errorf("cannot unmarshal the ciphertext in the query plugin function")
	ErrDecryption    = fmt.Errorf("cannot decrypt the ciphertext with the private key share in observation plugin function")
	ErrAggregation   = fmt.Errorf("cannot aggregate valid decryption shares in report plugn function")
)

type CiphertextId []byte

func (c CiphertextId) String() string {
	return "0x" + hex.EncodeToString(c)
}

type DecryptionRequest struct {
	CiphertextId CiphertextId
	Ciphertext   []byte
}

type DecryptionQueuingService interface {
	// GetRequests returns up to requestCountLimit oldest pending unique requests
	// with total size up to totalBytesLimit bytes size.
	GetRequests(requestCountLimit int, totalBytesLimit int) []DecryptionRequest

	// GetCiphertext returns the ciphertext matching ciphertextId
	// if it exists in the queue.
	// If the ciphertext does not exist it returns ErrNotFound.
	GetCiphertext(ciphertextId CiphertextId) ([]byte, error)

	// SetResult sets the plaintext (decrypted ciphertext) which corresponds to ciphertextId
	// or returns an error if the decrypted ciphertext is invalid.
	SetResult(ciphertextId CiphertextId, plaintext []byte, err error)
}

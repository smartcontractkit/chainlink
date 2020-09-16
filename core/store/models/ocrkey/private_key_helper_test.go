package ocrkey

import "math/rand"

// NewDeterministicOCRPrivateKeyXXXTestingOnly is for testing purposes only!
func NewDeterministicOCRPrivateKeyXXXTestingOnly(seed int64) (*OCRPrivateKey, error) {
	return newPrivateKey(rand.New(rand.NewSource(seed)))
}

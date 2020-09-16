package store

import (
	"github.com/smartcontractkit/chainlink/core/store/models/ocrkey"
)

// CreateWeakKeyXXXTestingOnly creates a new private key with weak encryption parameters and should
// only be used for testing
func (ks *OCRKeyStore) CreateWeakKeyXXXTestingOnly(auth string) (*ocrkey.OCRPrivateKey, error) {
	return ks.createKey(auth, ocrkey.FastScryptParams)
}

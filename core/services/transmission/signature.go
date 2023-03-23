package transmission

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/crypto"
)

func SignMessage(
	ownerPrivateKey *ecdsa.PrivateKey,
	message []byte,
) ([]byte, error) {
	sig, err := crypto.Sign(message, ownerPrivateKey)
	if err != nil {
		return nil, err
	}

	return sig, nil
}

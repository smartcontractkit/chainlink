package keystest

import (
	"crypto/ecdsa"
	"crypto/rand"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	pkgerrors "github.com/pkg/errors"
)

// NewKey pulled from geth
func NewKey() (key keystore.Key, err error) {
	privateKeyECDSA, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
	if err != nil {
		return key, err
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return key, pkgerrors.Errorf("Could not create random uuid: %v", err)
	}

	return keystore.Key{
		Id:         id,
		Address:    crypto.PubkeyToAddress(privateKeyECDSA.PublicKey),
		PrivateKey: privateKeyECDSA,
	}, nil
}

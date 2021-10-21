package terrakey

import (
	"crypto/ed25519"
	cryptorand "crypto/rand"

	cosmostypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/terra-project/terra.go/msg"

	"fmt"
)

type Raw []byte

type TerraAddress cosmostypes.Address

func (raw Raw) Key() KeyV2 {
	privKey := ed25519.PrivateKey(raw)
	return KeyV2{
		privateKey: &privKey,
		Address:    msg.AccAddress(ed25519PubKeyFromPrivKey(privKey)),
	}
}

func (raw Raw) String() string {
	return "<Terra Raw Private Key>"
}

func (raw Raw) GoString() string {
	return raw.String()
}

type KeyV2 struct {
	//TODO: this is only for OCR signing
	privateKey *ed25519.PrivateKey
	publicKey  ed25519.PublicKey

	Address TerraAddress
	// Type       string
}

func NewV2() (KeyV2, error) {
	pubKey, privKey, err := ed25519.GenerateKey(cryptorand.Reader)
	if err != nil {
		return KeyV2{}, err
	}
	return KeyV2{
		privateKey: &privKey,
		publicKey:  pubKey,
		Address:    msg.AccAddress(pubKey),
	}, nil
}

func (key KeyV2) ID() string {
	return string(key.Address.String())
}

func (key KeyV2) Raw() Raw {
	return Raw(*key.privateKey)
}

func (key KeyV2) String() string {
	return fmt.Sprintf("TerraKeyV2{PrivateKey: <redacted>, Address: %s}", key.Address)
}

func (key KeyV2) GoString() string {
	return key.String()
}

func ed25519PubKeyFromPrivKey(privKey ed25519.PrivateKey) ed25519.PublicKey {
	publicKey := make([]byte, ed25519.PublicKeySize)
	copy(publicKey, privKey[32:])
	return ed25519.PublicKey(publicKey)
}

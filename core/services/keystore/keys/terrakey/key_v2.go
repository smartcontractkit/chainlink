package terrakey

import (
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	cosmostypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/terra-project/terra.go/key"
	"github.com/terra-project/terra.go/msg"

	"fmt"
)

type Raw []byte

type TerraAddress cosmostypes.Address

func (raw Raw) Key() KeyV2 {
	privKey, _ := key.PrivKeyGen(raw)

	return KeyV2{
		privateKey: privKey,
		publicKey:  privKey.PubKey(),
		Address:    msg.AccAddress(privKey.PubKey().Address()),
	}
}

func (raw Raw) String() string {
	return "<Terra Raw Private Key>"
}

func (raw Raw) GoString() string {
	return raw.String()
}

type KeyV2 struct {
	privateKey cryptotypes.PrivKey
	publicKey  cryptotypes.PubKey
	Address    TerraAddress
	// TODO: choose type here? or put OCR pair somewhere else?
	// Type       string
}

func NewV2() (KeyV2, error) {
	privKey, _ := key.PrivKeyGen(secp256k1.GenPrivKey())

	return KeyV2{
		privateKey: privKey,
		publicKey:  privKey.PubKey(),
		Address:    msg.AccAddress(privKey.PubKey().Address()),
	}, nil
}

func (key KeyV2) ID() string {
	return string(key.Address.String())
}

func (key KeyV2) Raw() Raw {
	return Raw(key.privateKey.Bytes())
}

//TODO: temp method until we figure out signing
func (key KeyV2) Unsafe_GetPrivateKey() cryptotypes.PrivKey {
	return key.privateKey
}

func (key KeyV2) String() string {
	return fmt.Sprintf("TerraKeyV2{PrivateKey: <redacted>, Address: %s}", key.Address)
}

func (key KeyV2) GoString() string {
	return key.String()
}

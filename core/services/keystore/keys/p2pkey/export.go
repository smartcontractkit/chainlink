package p2pkey

import (
	"github.com/ethereum/go-ethereum/accounts/keystore"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/internal"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

const keyTypeIdentifier = "P2P"

func FromEncryptedJSON(keyJSON []byte, password string) (KeyV2, error) {
	return internal.FromEncryptedJSON(
		keyTypeIdentifier,
		keyJSON,
		password,
		adulteratedPassword,
		func(_ EncryptedP2PKeyExport, rawPrivKey internal.Raw) (KeyV2, error) {
			return KeyFor(rawPrivKey), nil
		},
	)
}

type EncryptedP2PKeyExport struct {
	KeyType   string              `json:"keyType"`
	PublicKey string              `json:"publicKey"`
	PeerID    PeerID              `json:"peerID"`
	Crypto    keystore.CryptoJSON `json:"crypto"`
}

func (x EncryptedP2PKeyExport) GetCrypto() keystore.CryptoJSON {
	return x.Crypto
}

func (key KeyV2) ToEncryptedJSON(password string, scryptParams utils.ScryptParams) (export []byte, err error) {
	return internal.ToEncryptedJSON(
		keyTypeIdentifier,
		key,
		password,
		scryptParams,
		adulteratedPassword,
		func(id string, key KeyV2, cryptoJSON keystore.CryptoJSON) EncryptedP2PKeyExport {
			return EncryptedP2PKeyExport{
				KeyType:   id,
				PublicKey: key.PublicKeyHex(),
				PeerID:    key.PeerID(),
				Crypto:    cryptoJSON,
			}
		},
	)
}

func adulteratedPassword(password string) string {
	return "p2pkey" + password
}

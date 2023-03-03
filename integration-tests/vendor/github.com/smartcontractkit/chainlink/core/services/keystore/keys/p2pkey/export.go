package p2pkey

import (
	keystore "github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/services/keystore/keys"
	"github.com/smartcontractkit/chainlink/core/utils"
)

const keyTypeIdentifier = "P2P"

func FromEncryptedJSON(keyJSON []byte, password string) (KeyV2, error) {
	return keys.FromEncryptedJSON(
		keyTypeIdentifier,
		keyJSON,
		password,
		adulteratedPassword,
		func(_ EncryptedP2PKeyExport, rawPrivKey []byte) (KeyV2, error) {
			return Raw(rawPrivKey).Key(), nil
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
	return keys.ToEncryptedJSON(
		keyTypeIdentifier,
		key.Raw(),
		key,
		password,
		scryptParams,
		adulteratedPassword,
		func(id string, key KeyV2, cryptoJSON keystore.CryptoJSON) (EncryptedP2PKeyExport, error) {
			rawPubKey, err := key.GetPublic().Bytes()
			if err != nil {
				return EncryptedP2PKeyExport{}, errors.Wrapf(err, "could not get raw public key")
			}
			return EncryptedP2PKeyExport{
				KeyType:   id,
				PublicKey: hexutil.Encode(rawPubKey),
				PeerID:    key.PeerID(),
				Crypto:    cryptoJSON,
			}, nil
		},
	)
}

func adulteratedPassword(password string) string {
	return "p2pkey" + password
}

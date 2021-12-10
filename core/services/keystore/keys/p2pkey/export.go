package p2pkey

import (
	"encoding/json"

	keystore "github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/utils"
)

const keyTypeIdentifier = "P2P"

func FromEncryptedJSON(keyJSON []byte, password string) (KeyV2, error) {
	var export EncryptedP2PKeyExport
	if err := json.Unmarshal(keyJSON, &export); err != nil {
		return KeyV2{}, err
	}
	privKey, err := keystore.DecryptDataV3(export.Crypto, adulteratedPassword(password))
	if err != nil {
		return KeyV2{}, errors.Wrap(err, "failed to decrypt P2P key")
	}
	key := Raw(privKey).Key()
	return key, nil
}

type EncryptedP2PKeyExport struct {
	KeyType   string              `json:"keyType"`
	PublicKey string              `json:"publicKey"`
	PeerID    PeerID              `json:"peerID"`
	Crypto    keystore.CryptoJSON `json:"crypto"`
}

func (key KeyV2) ToEncryptedJSON(password string, scryptParams utils.ScryptParams) (export []byte, err error) {
	cryptoJSON, err := keystore.EncryptDataV3(
		key.Raw(),
		[]byte(adulteratedPassword(password)),
		scryptParams.N,
		scryptParams.P,
	)
	if err != nil {
		return nil, errors.Wrapf(err, "could not encrypt P2P key")
	}
	rawPubKey, err := key.GetPublic().Bytes()
	if err != nil {
		return nil, errors.Wrapf(err, "could not get raw public key")
	}
	encryptedP2PKeyExport := EncryptedP2PKeyExport{
		KeyType:   keyTypeIdentifier,
		PublicKey: hexutil.Encode(rawPubKey),
		PeerID:    key.PeerID(),
		Crypto:    cryptoJSON,
	}
	return json.Marshal(encryptedP2PKeyExport)
}

func adulteratedPassword(password string) string {
	return "p2pkey" + password
}

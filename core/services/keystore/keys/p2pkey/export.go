package p2pkey

import (
	"encoding/json"

	keystore "github.com/ethereum/go-ethereum/accounts/keystore"
	cryptop2p "github.com/libp2p/go-libp2p-core/crypto"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// EncryptedP2PKeyExport represents the structure of P2P keys exported and imported
// to/from the disk
type EncryptedP2PKeyExport struct {
	PublicKey PublicKeyBytes      `json:"publicKey"`
	PeerID    PeerID              `json:"peerID"`
	Crypto    keystore.CryptoJSON `json:"crypto"`
}

func (k Key) ToEncryptedExport(auth string, scryptParams utils.ScryptParams) (export []byte, err error) {
	var marshalledPrivK []byte
	marshalledPrivK, err = cryptop2p.MarshalPrivateKey(k)
	if err != nil {
		return export, err
	}
	cryptoJSON, err := keystore.EncryptDataV3(marshalledPrivK, []byte(adulteratedPassword(auth)), scryptParams.N, scryptParams.P)
	if err != nil {
		return export, errors.Wrapf(err, "could not encrypt p2p key")
	}

	pubKeyBytes, err := k.GetPublic().Raw()
	if err != nil {
		return export, errors.Wrapf(err, "could not ger public key bytes from private key")
	}
	peerID, err := k.GetPeerID()
	if err != nil {
		return export, errors.Wrapf(err, "could not ger peerID from private key")
	}

	encryptedP2PKExport := EncryptedP2PKeyExport{
		PublicKey: pubKeyBytes,
		PeerID:    peerID,
		Crypto:    cryptoJSON,
	}
	return json.Marshal(encryptedP2PKExport)
}

// DecryptPrivateKey returns the PrivateKey in export, decrypted via auth, or an error
func (export EncryptedP2PKeyExport) DecryptPrivateKey(auth string) (k *Key, err error) {
	marshalledPrivK, err := keystore.DecryptDataV3(export.Crypto, adulteratedPassword(auth))
	if err != nil {
		return k, errors.Wrapf(err, "could not decrypt key 0x%x", export.PublicKey)
	}
	privK, err := cryptop2p.UnmarshalPrivateKey(marshalledPrivK)
	if err != nil {
		return k, errors.Wrapf(err, "could not unmarshal private key for 0x%x", export.PublicKey)
	}
	return &Key{
		privK,
	}, nil
}

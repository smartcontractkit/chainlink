package p2pkey

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"time"

	keystore "github.com/ethereum/go-ethereum/accounts/keystore"
	cryptop2p "github.com/libp2p/go-libp2p-core/crypto"
	peer "github.com/libp2p/go-libp2p-core/peer"
	"github.com/pkg/errors"
)

// Key represents a libp2p private key
type Key struct {
	cryptop2p.PrivKey
}

func (k Key) GetPeerID() (peer.ID, error) {
	peerID, err := peer.IDFromPrivateKey(k)
	if err != nil {
		return "", errors.WithStack(err)
	}
	return peerID, err
}

type EncryptedP2PKey struct {
	ID               int32 `gorm:"primary_key"`
	PeerID           string
	PubKey           []byte
	EncryptedPrivKey []byte
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func (EncryptedP2PKey) TableName() string {
	return "encrypted_p2p_keys"
}

// CreateKey makes a new libp2p keypair from a crytographically secure entropy source
func CreateKey() (Key, error) {
	p2pPrivkey, _, err := cryptop2p.GenerateEd25519Key(rand.Reader)
	if err != nil {
		return Key{}, nil
	}
	return Key{
		p2pPrivkey,
	}, nil
}

type ScryptParams struct{ N, P int }

var defaultScryptParams = ScryptParams{
	N: keystore.StandardScryptN, P: keystore.StandardScryptP}

// type is added to the beginning of the passwords for
// P2P key, so that the keys can't accidentally be mis-used
// in the wrong place
func adulteratedPassword(auth string) string {
	s := "p2pkey" + auth
	return s
}

func (k Key) Encrypt(auth string, p ...ScryptParams) (s EncryptedP2PKey, err error) {
	var scryptParams ScryptParams
	switch len(p) {
	case 0:
		scryptParams = defaultScryptParams
	case 1:
		scryptParams = p[0]
	default:
		return s, fmt.Errorf("can take at most one set of ScryptParams")
	}
	var marshalledPrivK []byte
	marshalledPrivK, err = cryptop2p.MarshalPrivateKey(k)
	if err != nil {
		return s, err
	}
	cryptoJSON, err := keystore.EncryptDataV3(marshalledPrivK, []byte(adulteratedPassword(auth)), scryptParams.N, scryptParams.P)
	if err != nil {
		return s, errors.Wrapf(err, "could not encrypt p2p key")
	}
	marshalledCryptoJSON, err := json.Marshal(&cryptoJSON)
	if err != nil {
		return s, errors.Wrapf(err, "could not encode cryptoJSON")
	}
	peerID, err := k.GetPeerID()
	if err != nil {
		return s, errors.Wrapf(err, "could not get peer ID")
	}
	pubKeyBytes, err := k.GetPublic().Raw()
	if err != nil {
		return s, errors.Wrapf(err, "could not get public key bytes")
	}

	s = EncryptedP2PKey{
		PubKey:           pubKeyBytes,
		EncryptedPrivKey: marshalledCryptoJSON,
		PeerID:           peerID.Pretty(),
	}
	return s, nil
}

// Decrypt returns the PrivateKey in e, decrypted via auth, or an error
func (e EncryptedP2PKey) Decrypt(auth string) (k Key, err error) {
	var cryptoJSON keystore.CryptoJSON
	err = json.Unmarshal(e.EncryptedPrivKey, &cryptoJSON)
	if err != nil {
		return k, errors.Wrapf(err, "invalid JSON for key 0x%x", e.PubKey)
	}
	marshalledPrivK, err := keystore.DecryptDataV3(cryptoJSON, adulteratedPassword(auth))
	if err != nil {
		return k, errors.Wrapf(err, "could not decrypt key 0x%x", e.PubKey)
	}
	privK, err := cryptop2p.UnmarshalPrivateKey(marshalledPrivK)
	if err != nil {
		return k, errors.Wrapf(err, "could not unmarshal private key for 0x%x", e.PubKey)
	}
	return Key{
		privK,
	}, nil
}

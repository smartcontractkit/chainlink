package p2pkey

import (
	"crypto/ed25519"
	"database/sql/driver"
	"encoding/hex"
	"encoding/json"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/pkg/errors"

	ragep2ptypes "github.com/smartcontractkit/libocr/ragep2p/types"
)

// Key represents a p2p private key
type Key struct {
	PrivKey ed25519.PrivateKey
}

func (k Key) ToV2() KeyV2 {
	return KeyV2{
		PrivKey: k.PrivKey,
		peerID:  k.PeerID(),
	}
}

// PublicKeyBytes is a [ed25519.PublicKey]
type PublicKeyBytes []byte

func (pkb PublicKeyBytes) String() string {
	return hex.EncodeToString(pkb)
}

func (pkb PublicKeyBytes) MarshalJSON() ([]byte, error) {
	return json.Marshal(hex.EncodeToString(pkb))
}

func (pkb *PublicKeyBytes) UnmarshalJSON(input []byte) error {
	var hexString string
	if err := json.Unmarshal(input, &hexString); err != nil {
		return err
	}

	result, err := hex.DecodeString(hexString)
	if err != nil {
		return err
	}

	*pkb = result
	return nil
}

func (pkb *PublicKeyBytes) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		*pkb = v
		return nil
	default:
		return errors.Errorf("invalid public key bytes got %T wanted []byte", v)
	}
}

func (pkb PublicKeyBytes) Value() (driver.Value, error) {
	return []byte(pkb), nil
}

func (k Key) GetPeerID() (PeerID, error) {
	peerID, err := ragep2ptypes.PeerIDFromPrivateKey(k.PrivKey)
	if err != nil {
		return PeerID{}, errors.WithStack(err)
	}
	return PeerID(peerID), err
}

func (k Key) PeerID() PeerID {
	peerID, err := k.GetPeerID()
	if err != nil {
		panic(err)
	}
	return peerID
}

type EncryptedP2PKey struct {
	ID               int32
	PeerID           PeerID
	PubKey           PublicKeyBytes
	EncryptedPrivKey []byte
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        *time.Time
}

func (ep2pk *EncryptedP2PKey) SetID(value string) error {
	result, err := strconv.ParseInt(value, 10, 32)

	if err != nil {
		return err
	}

	ep2pk.ID = int32(result)
	return nil
}

// Decrypt returns the PrivateKey in e, decrypted via auth, or an error
func (ep2pk EncryptedP2PKey) Decrypt(auth string) (k Key, err error) {
	var cryptoJSON keystore.CryptoJSON
	err = json.Unmarshal(ep2pk.EncryptedPrivKey, &cryptoJSON)
	if err != nil {
		return k, errors.Wrapf(err, "invalid JSON for P2P key %s (0x%x)", ep2pk.PeerID.String(), ep2pk.PubKey)
	}
	marshalledPrivK, err := keystore.DecryptDataV3(cryptoJSON, adulteratedPassword(auth))
	if err != nil {
		return k, errors.Wrapf(err, "could not decrypt P2P key %s (0x%x)", ep2pk.PeerID.String(), ep2pk.PubKey)
	}

	privK, err := UnmarshalPrivateKey(marshalledPrivK)
	if err != nil {
		return k, errors.Wrapf(err, "could not unmarshal P2P private key for %s (0x%x)", ep2pk.PeerID.String(), ep2pk.PubKey)
	}
	return Key{
		privK,
	}, nil
}

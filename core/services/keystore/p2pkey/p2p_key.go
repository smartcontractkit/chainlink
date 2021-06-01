package p2pkey

import (
	"crypto/rand"
	"database/sql/driver"
	"encoding/hex"
	"encoding/json"
	"strconv"
	"time"

	"gorm.io/gorm/schema"

	"gorm.io/gorm"

	keystore "github.com/ethereum/go-ethereum/accounts/keystore"
	cryptop2p "github.com/libp2p/go-libp2p-core/crypto"
	peer "github.com/libp2p/go-libp2p-core/peer"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// Key represents a libp2p private key
type Key struct {
	cryptop2p.PrivKey
}

// PublicKeyBytes is generated using cryptop2p.PubKey.Raw()
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

	*pkb = PublicKeyBytes(result)
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

// GormDataType gorm common data type
func (PublicKeyBytes) GormDataType() string {
	return "bytea"
}

// GormDBDataType gorm db data type
func (PublicKeyBytes) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "postgres":
		return "BYTEA"
	}
	return ""
}

func (k Key) GetPeerID() (PeerID, error) {
	peerID, err := peer.IDFromPrivateKey(k)
	if err != nil {
		return "", errors.WithStack(err)
	}
	return PeerID(peerID), err
}

func (k Key) MustGetPeerID() PeerID {
	peerID, err := peer.IDFromPrivateKey(k)
	if err != nil {
		panic(err)
	}
	return PeerID(peerID)
}

type EncryptedP2PKey struct {
	ID               int32 `gorm:"primary_key"`
	PeerID           PeerID
	PubKey           PublicKeyBytes `gorm:"type:bytea"`
	EncryptedPrivKey []byte
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt
}

func (EncryptedP2PKey) TableName() string {
	return "encrypted_p2p_keys"
}

func (ep2pk *EncryptedP2PKey) SetID(value string) error {
	result, err := strconv.ParseInt(value, 10, 32)

	if err != nil {
		return err
	}

	ep2pk.ID = int32(result)
	return nil
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

// type is added to the beginning of the passwords for
// P2P key, so that the keys can't accidentally be mis-used
// in the wrong place
func adulteratedPassword(auth string) string {
	s := "p2pkey" + auth
	return s
}

func (k Key) ToEncryptedP2PKey(auth string, scryptParams utils.ScryptParams) (s EncryptedP2PKey, err error) {
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
		PeerID:           peerID,
	}
	return s, nil
}

// Decrypt returns the PrivateKey in e, decrypted via auth, or an error
func (ep2pk EncryptedP2PKey) Decrypt(auth string) (k Key, err error) {
	var cryptoJSON keystore.CryptoJSON
	err = json.Unmarshal(ep2pk.EncryptedPrivKey, &cryptoJSON)
	if err != nil {
		return k, errors.Wrapf(err, "invalid JSON for key 0x%x", ep2pk.PubKey)
	}
	marshalledPrivK, err := keystore.DecryptDataV3(cryptoJSON, adulteratedPassword(auth))
	if err != nil {
		return k, errors.Wrapf(err, "could not decrypt key 0x%x", ep2pk.PubKey)
	}
	privK, err := cryptop2p.UnmarshalPrivateKey(marshalledPrivK)
	if err != nil {
		return k, errors.Wrapf(err, "could not unmarshal private key for 0x%x", ep2pk.PubKey)
	}
	return Key{
		privK,
	}, nil
}

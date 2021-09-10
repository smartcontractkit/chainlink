package ocr2key

import (
	"crypto/ed25519"
	cryptorand "crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"
	"golang.org/x/crypto/curve25519"
	"gorm.io/gorm"
)

type (
	// KeyBundle represents the bundle of keys needed for OCR
	KeyBundle struct {
		id              models.Sha256Hash
		OffchainKeyring OffchainKeyring
		OnchainKeyring  EthereumKeyring
	}

	// EncryptedKeyBundle holds an encrypted KeyBundle
	EncryptedKeyBundle struct {
		ID models.Sha256Hash `gorm:"primary_key"`

		OnchainPublicKey      []byte
		OnchainSigningAddress common.Address

		OffchainSigningPublicKey    []byte
		OffchainEncryptionPublicKey []byte

		EncryptedPrivateKeys []byte
		CreatedAt            time.Time
		UpdatedAt            time.Time
		DeletedAt            gorm.DeletedAt
	}

	KeyBundleRawData struct {
		OffchainKeyring []byte
		OnchainKeyring  []byte
	}
)

var (
	curve = secp256k1.S256()
)

func (EncryptedKeyBundle) TableName() string {
	return "encrypted_ocr2_key_bundles"
}

func (ekb EncryptedKeyBundle) GetID() string {
	return ekb.ID.String()
}

func (ekb *EncryptedKeyBundle) SetID(value string) error {
	var result models.Sha256Hash
	decodedString, err := hex.DecodeString(value)

	if err != nil {
		return err
	}

	copy(result[:], decodedString[:32])
	ekb.ID = result
	return nil
}

// NewKeyBundle makes a new set of OCR key bundles from cryptographically secure entropy
func NewKeyBundle() (*KeyBundle, error) {
	return NewKeyBundleFrom(cryptorand.Reader, cryptorand.Reader, cryptorand.Reader)
}

func NewKeyBundleFrom(onchainSigningKeyMaterial, onchainEncryptionKeyMaterial, offchainKeyMaterial io.Reader) (*KeyBundle, error) {
	onchainKeyring, err := NewEthereumKeyring(offchainKeyMaterial)
	if err != nil {
		return nil, err
	}
	offchainKeyring, err := NewOffchainKeyring(onchainSigningKeyMaterial, onchainEncryptionKeyMaterial)
	if err != nil {
		return nil, err
	}
	k := &KeyBundle{
		OnchainKeyring:  *onchainKeyring,
		OffchainKeyring: *offchainKeyring,
	}
	marshalledPrivK, err := k.marshal()
	if err != nil {
		return nil, err
	}
	k.id = sha256.Sum256(marshalledPrivK)
	return k, nil
}

// ConfigDiffieHellman returns the shared point obtained by multiplying someone's
// public key by a secret scalar ( in this case, the offChainEncryption key.)
func (pk *KeyBundle) ConfigDiffieHellman(base [curve25519.PointSize]byte) ([curve25519.PointSize]byte, error) {
	return pk.OffchainKeyring.ConfigDiffieHellman(base)
}

// PublicKeyAddressOnChain returns public component of the keypair used in
func (pk *KeyBundle) PublicKeyAddressOnChain() common.Address {
	return pk.OnchainKeyring.SigningAddress()
}

// PublicKeyOffChain returns the pbulic component of the keypair used in SignOffChain
func (pk *KeyBundle) PublicKeyOffChain() ocrtypes.OffchainPublicKey {
	return pk.OffchainKeyring.OffchainPublicKey()
}

// PublicKeyConfig returns the public component of the keypair used in ConfigKeyShare
func (pk *KeyBundle) PublicKeyConfig() [curve25519.PointSize]byte {
	return pk.OffchainKeyring.ConfigEncryptionPublicKey()
}

// Encrypt combines the KeyBundle into a single json-serialized
// bytes array and then encrypts
func (pk *KeyBundle) Encrypt(auth string, scryptParams utils.ScryptParams) (*EncryptedKeyBundle, error) {
	return pk.encrypt(auth, scryptParams)
}

// encrypt combines the KeyBundle into a single json-serialized
// bytes array and then encrypts, using the provided scrypt params
// separated into a different function so that scryptParams can be
// weakened in tests
func (pk *KeyBundle) encrypt(auth string, scryptParams utils.ScryptParams) (*EncryptedKeyBundle, error) {
	marshalledPrivK, err := pk.marshal()
	if err != nil {
		return nil, err
	}
	cryptoJSON, err := keystore.EncryptDataV3(
		marshalledPrivK,
		[]byte(adulteratedPassword(auth)),
		scryptParams.N,
		scryptParams.P,
	)
	if err != nil {
		return nil, errors.Wrapf(err, "could not encrypt ocr key")
	}
	encryptedPrivKeys, err := json.Marshal(&cryptoJSON)
	if err != nil {
		return nil, errors.Wrapf(err, "could not encode cryptoJSON")
	}
	kb := &EncryptedKeyBundle{
		ID:                          pk.id,
		OnchainPublicKey:            pk.OnchainKeyring.PublicKey(),
		OnchainSigningAddress:       pk.OnchainKeyring.SigningAddress(),
		OffchainSigningPublicKey:    pk.OffchainKeyring.OffchainPublicKey(),
		OffchainEncryptionPublicKey: make([]byte, ed25519.PublicKeySize),
		EncryptedPrivateKeys:        encryptedPrivKeys,
	}
	configEncryptionPublicKey := pk.OffchainKeyring.ConfigEncryptionPublicKey()
	copy(kb.OffchainEncryptionPublicKey[:], configEncryptionPublicKey[:])
	return kb, nil
}

// Decrypt returns the PrivateKeys in e, decrypted via auth, or an error
func (ekb *EncryptedKeyBundle) Decrypt(auth string) (*KeyBundle, error) {
	var cryptoJSON keystore.CryptoJSON
	err := json.Unmarshal(ekb.EncryptedPrivateKeys, &cryptoJSON)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid cryptoJSON for OCR2 key bundle")
	}
	marshalledPrivK, err := keystore.DecryptDataV3(cryptoJSON, adulteratedPassword(auth))
	if err != nil {
		return nil, errors.Wrapf(err, "could not decrypt OCR2 key bundle")
	}
	var pk KeyBundle
	err = pk.unmarshal(marshalledPrivK)
	if err != nil {
		return nil, errors.Wrapf(err, "could not unmarshal OCR2 key bundle")
	}
	return &pk, nil
}

func (pk *KeyBundle) marshal() ([]byte, error) {
	onchainKeyringBytes, err := pk.OnchainKeyring.marshal()
	if err != nil {
		return nil, err
	}
	offchainKeyringBytes, err := pk.OffchainKeyring.marshal()
	if err != nil {
		return nil, err
	}
	rawKeyData := KeyBundleRawData{
		OffchainKeyring: offchainKeyringBytes,
		OnchainKeyring:  onchainKeyringBytes,
	}
	return json.Marshal(&rawKeyData)
}

func (pk *KeyBundle) unmarshal(b []byte) (err error) {
	var rawKeyData KeyBundleRawData
	err = json.Unmarshal(b, &rawKeyData)
	if err != nil {
		return err
	}
	err = pk.OnchainKeyring.unmarshal(rawKeyData.OnchainKeyring)
	if err != nil {
		return err
	}
	err = pk.OffchainKeyring.unmarshal(rawKeyData.OffchainKeyring)
	if err != nil {
		return err
	}
	pk.id = sha256.Sum256(b)
	return nil
}

// String reduces the risk of accidentally logging the private key
func (pk KeyBundle) String() string {
	addressOnChain := pk.PublicKeyAddressOnChain()
	return fmt.Sprintf(
		"KeyBundle{PublicKeyAddressOnChain: %s, PublicKeyOffChain: %s}",
		hex.EncodeToString(addressOnChain[:]),
		hex.EncodeToString(pk.PublicKeyOffChain()),
	)
}

// GoStringer reduces the risk of accidentally logging the private key
func (pk KeyBundle) GoStringer() string {
	return pk.String()
}

// type is added to the beginning of the passwords for OCR key bundles,
// so that the keys can't accidentally be mis-used in the wrong place
func adulteratedPassword(auth string) string {
	s := "ocr2key" + auth
	return s
}

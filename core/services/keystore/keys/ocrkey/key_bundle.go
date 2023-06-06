package ocrkey

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	cryptorand "crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/pkg/errors"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting/types"
	"golang.org/x/crypto/curve25519"

	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type (
	// KeyBundle represents the bundle of keys needed for OCR
	KeyBundle struct {
		ID                 models.Sha256Hash
		onChainSigning     *onChainPrivateKey
		offChainSigning    *offChainPrivateKey
		offChainEncryption *[curve25519.ScalarSize]byte
	}

	// EncryptedKeyBundle holds an encrypted KeyBundle
	EncryptedKeyBundle struct {
		ID                    models.Sha256Hash
		OnChainSigningAddress OnChainSigningAddress
		OffChainPublicKey     OffChainPublicKey
		ConfigPublicKey       ConfigPublicKey
		EncryptedPrivateKeys  []byte
		CreatedAt             time.Time
		UpdatedAt             time.Time
		DeletedAt             *time.Time
	}
)

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

// New makes a new set of OCR key bundles from cryptographically secure entropy
func New() (*KeyBundle, error) {
	return NewFrom(cryptorand.Reader, cryptorand.Reader, cryptorand.Reader)
}

// NewFrom makes a new set of OCR key bundles from cryptographically secure entropy
func NewFrom(onChainSigning io.Reader, offChainSigning io.Reader, offChainEncryption io.Reader) (*KeyBundle, error) {
	ecdsaKey, err := ecdsa.GenerateKey(curve, onChainSigning)
	if err != nil {
		return nil, err
	}
	onChainPriv := (*onChainPrivateKey)(ecdsaKey)

	_, offChainPriv, err := ed25519.GenerateKey(offChainSigning)
	if err != nil {
		return nil, err
	}
	var encryptionPriv [curve25519.ScalarSize]byte
	_, err = offChainEncryption.Read(encryptionPriv[:])
	if err != nil {
		return nil, err
	}
	k := &KeyBundle{
		onChainSigning:     onChainPriv,
		offChainSigning:    (*offChainPrivateKey)(&offChainPriv),
		offChainEncryption: &encryptionPriv,
	}
	marshalledPrivK, err := json.Marshal(k)
	if err != nil {
		return nil, err
	}
	k.ID = sha256.Sum256(marshalledPrivK)
	return k, nil
}

// SignOnChain returns an ethereum-style ECDSA secp256k1 signature on msg.
func (pk *KeyBundle) SignOnChain(msg []byte) (signature []byte, err error) {
	return pk.onChainSigning.Sign(msg)
}

// SignOffChain returns an EdDSA-Ed25519 signature on msg.
func (pk *KeyBundle) SignOffChain(msg []byte) (signature []byte, err error) {
	return pk.offChainSigning.Sign(msg)
}

// ConfigDiffieHellman returns the shared point obtained by multiplying someone's
// public key by a secret scalar ( in this case, the offChainEncryption key.)
func (pk *KeyBundle) ConfigDiffieHellman(base *[curve25519.PointSize]byte) (
	sharedPoint *[curve25519.PointSize]byte, err error,
) {
	p, err := curve25519.X25519(pk.offChainEncryption[:], base[:])
	if err != nil {
		return nil, err
	}
	sharedPoint = new([ed25519.PublicKeySize]byte)
	copy(sharedPoint[:], p)
	return sharedPoint, nil
}

// PublicKeyAddressOnChain returns public component of the keypair used in
// SignOnChain
func (pk *KeyBundle) PublicKeyAddressOnChain() ocrtypes.OnChainSigningAddress {
	return ocrtypes.OnChainSigningAddress(pk.onChainSigning.Address())
}

// PublicKeyOffChain returns the public component of the keypair used in SignOffChain
func (pk *KeyBundle) PublicKeyOffChain() ocrtypes.OffchainPublicKey {
	return ocrtypes.OffchainPublicKey(pk.offChainSigning.PublicKey())
}

// PublicKeyConfig returns the public component of the keypair used in ConfigKeyShare
func (pk *KeyBundle) PublicKeyConfig() [curve25519.PointSize]byte {
	rv, err := curve25519.X25519(pk.offChainEncryption[:], curve25519.Basepoint)
	if err != nil {
		log.Println("failure while computing public key: " + err.Error())
	}
	var rvFixed [curve25519.PointSize]byte
	copy(rvFixed[:], rv)
	return rvFixed
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
	marshalledPrivK, err := json.Marshal(&pk)
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
	return &EncryptedKeyBundle{
		ID:                    pk.ID,
		OnChainSigningAddress: pk.onChainSigning.Address(),
		OffChainPublicKey:     pk.offChainSigning.PublicKey(),
		ConfigPublicKey:       pk.PublicKeyConfig(),
		EncryptedPrivateKeys:  encryptedPrivKeys,
	}, nil
}

// Decrypt returns the PrivateKeys in e, decrypted via auth, or an error
func (ekb *EncryptedKeyBundle) Decrypt(auth string) (*KeyBundle, error) {
	var cryptoJSON keystore.CryptoJSON
	err := json.Unmarshal(ekb.EncryptedPrivateKeys, &cryptoJSON)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid cryptoJSON for OCR key bundle")
	}
	marshalledPrivK, err := keystore.DecryptDataV3(cryptoJSON, adulteratedPassword(auth))
	if err != nil {
		return nil, errors.Wrapf(err, "could not decrypt OCR key bundle")
	}
	var pk KeyBundle
	err = json.Unmarshal(marshalledPrivK, &pk)
	if err != nil {
		return nil, errors.Wrapf(err, "could not unmarshal OCR key bundle")
	}
	return &pk, nil
}

// MarshalJSON marshals the private keys into json
func (pk *KeyBundle) MarshalJSON() ([]byte, error) {
	rawKeyData := keyBundleRawData{
		EcdsaD:             *pk.onChainSigning.D,
		Ed25519PrivKey:     []byte(*pk.offChainSigning),
		OffChainEncryption: *pk.offChainEncryption,
	}
	return json.Marshal(&rawKeyData)
}

// UnmarshalJSON constructs KeyBundle from raw json
func (pk *KeyBundle) UnmarshalJSON(b []byte) (err error) {
	var rawKeyData keyBundleRawData
	err = json.Unmarshal(b, &rawKeyData)
	if err != nil {
		return err
	}
	ecdsaDSize := len(rawKeyData.EcdsaD.Bytes())
	if ecdsaDSize > curve25519.PointSize {
		return errors.Wrapf(ErrScalarTooBig, "got %d byte ecdsa scalar", ecdsaDSize)
	}

	publicKey := ecdsa.PublicKey{Curve: curve}
	publicKey.X, publicKey.Y = curve.ScalarBaseMult(rawKeyData.EcdsaD.Bytes())
	privateKey := ecdsa.PrivateKey{
		PublicKey: publicKey,
		D:         &rawKeyData.EcdsaD,
	}
	onChainSigning := onChainPrivateKey(privateKey)
	offChainSigning := offChainPrivateKey(rawKeyData.Ed25519PrivKey)
	pk.onChainSigning = &onChainSigning
	pk.offChainSigning = &offChainSigning
	pk.offChainEncryption = &rawKeyData.OffChainEncryption
	pk.ID = sha256.Sum256(b)
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

// GoString reduces the risk of accidentally logging the private key
func (pk KeyBundle) GoString() string {
	return pk.String()
}

// GoString reduces the risk of accidentally logging the private key
func (pk KeyBundle) ToV2() KeyV2 {
	return KeyV2{
		OnChainSigning:     pk.onChainSigning,
		OffChainSigning:    pk.offChainSigning,
		OffChainEncryption: pk.offChainEncryption,
	}
}

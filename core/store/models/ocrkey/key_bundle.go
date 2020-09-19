package ocrkey

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"time"

	cryptorand "crypto/rand"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/pkg/errors"
	"golang.org/x/crypto/curve25519"
)

// KeyBundle represents the bundle of keys needed for OCR
type KeyBundle struct {
	ID                 string
	onChainSigning     *onChainPrivateKey
	offChainSigning    *offChainPrivateKey
	offChainEncryption *[curve25519.ScalarSize]byte
}

// EncryptedKeyBundle holds an encrypted OCRKeyBundle bundle
type EncryptedKeyBundle struct {
	ID                    string `gorm:"primary_key"`
	OnChainSigningAddress OnChainSigningAddress
	OffChainPublicKey     OffChainPublicKey
	EncryptedPrivateKeys  []byte
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

type keyBundleRawData struct {
	EcdsaX             big.Int
	EcdsaY             big.Int
	EcdsaD             big.Int
	Ed25519PrivKey     []byte
	OffChainEncryption [curve25519.ScalarSize]byte
}

type scryptParams struct{ N, P int }

var defaultScryptParams = scryptParams{
	N: keystore.StandardScryptN, P: keystore.StandardScryptP}

var curve = secp256k1.S256()

func (EncryptedKeyBundle) TableName() string {
	return "encrypted_ocr_key_bundles"
}

// NewKeyBundle makes a new set of OCR key bundles from cryptographically secure entropy
func NewKeyBundle() (*KeyBundle, error) {
	reader := cryptorand.Reader
	onChainSk, err := cryptorand.Int(reader, crypto.S256().Params().N)
	if err != nil {
		return nil, err
	}
	onChainPriv := new(onChainPrivateKey)
	p := (*ecdsa.PrivateKey)(onChainPriv)
	p.D = onChainSk
	onChainPriv.PublicKey = ecdsa.PublicKey{Curve: curve}
	p.PublicKey.X, p.PublicKey.Y = curve.ScalarBaseMult(onChainSk.Bytes())
	_, offChainPriv, err := ed25519.GenerateKey(reader)
	if err != nil {
		return nil, err
	}
	var encryptionPriv [curve25519.ScalarSize]byte
	_, err = reader.Read(encryptionPriv[:])
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
	byteID := sha256.Sum256(marshalledPrivK)
	k.ID = hex.EncodeToString(byteID[:])
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

// ConfigDiffieHelman returns the shared point obtained by multiplying someone's
// public key by a secret scalar ( in this case, the offChainEncryption key.)
func (pk *KeyBundle) ConfigDiffieHelman(base *[curve25519.PointSize]byte) (
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
func (pk *KeyBundle) PublicKeyAddressOnChain() OnChainSigningAddress {
	return pk.onChainSigning.Address()
}

// PublicKeyOffChain returns the pbulic component of the keypair used in SignOffChain
func (pk *KeyBundle) PublicKeyOffChain() OffChainPublicKey {
	return OffChainPublicKey(pk.offChainSigning.PublicKey())
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
func (pk *KeyBundle) Encrypt(auth string) (*EncryptedKeyBundle, error) {
	return pk.encrypt(auth, defaultScryptParams)
}

// encrypt combines the KeyBundle into a single json-serialized
// bytes array and then encrypts, using the provided scrypt params
// separated into a different function so that scryptParams can be
// weakened in tests
func (pk *KeyBundle) encrypt(auth string, scryptParams scryptParams) (*EncryptedKeyBundle, error) {
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
		EncryptedPrivateKeys:  encryptedPrivKeys,
	}, nil
}

// Decrypt returns the PrivateKeys in e, decrypted via auth, or an error
func (encKey *EncryptedKeyBundle) Decrypt(auth string) (*KeyBundle, error) {
	var cryptoJSON keystore.CryptoJSON
	err := json.Unmarshal(encKey.EncryptedPrivateKeys, &cryptoJSON)
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
	pk.ID = encKey.ID
	return &pk, nil
}

// MarshalJSON marshals the private keys into json
func (pk *KeyBundle) MarshalJSON() ([]byte, error) {
	rawKeyData := keyBundleRawData{
		EcdsaX:             *pk.onChainSigning.X,
		EcdsaY:             *pk.onChainSigning.Y,
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
	publicKey := ecdsa.PublicKey{
		X: &rawKeyData.EcdsaX,
		Y: &rawKeyData.EcdsaY,
	}
	privateKey := ecdsa.PrivateKey{
		PublicKey: publicKey,
		D:         &rawKeyData.EcdsaD,
	}
	onChainSigning := onChainPrivateKey(privateKey)
	offChainSigning := offChainPrivateKey(rawKeyData.Ed25519PrivKey)
	pk.onChainSigning = &onChainSigning
	pk.offChainSigning = &offChainSigning
	pk.offChainEncryption = &rawKeyData.OffChainEncryption
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
	s := "ocrkey" + auth
	return s
}

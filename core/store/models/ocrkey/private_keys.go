package ocrkey

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"math/rand"
	"time"

	cryptorand "crypto/rand"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/offchain-reporting-design/prototype/offchainreporting/to_be_internal/signature"
	"github.com/smartcontractkit/offchain-reporting-design/prototype/offchainreporting/types"
	"golang.org/x/crypto/curve25519"
)

type OCRPrivateKeys struct {
	onChainSigning     *signature.OnChainPrivateKey
	offChainSigning    *signature.OffChainPrivateKey
	offChainEncryption *[curve25519.ScalarSize]byte
}

type EncryptedOCRPrivateKeys struct {
	ID                int32 `gorm:"primary_key"`
	EncryptedPrivKeys []byte
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type ScryptParams struct{ N, P int }

type ocrPrivateKeysRawData struct {
	X                  big.Int
	Y                  big.Int
	D                  big.Int
	Ed25519PrivKey     []byte
	OffChainEncryption [curve25519.ScalarSize]byte
}

var defaultScryptParams = ScryptParams{
	N: keystore.StandardScryptN, P: keystore.StandardScryptP}

var _ types.PrivateKeys = (*OCRPrivateKeys)(nil)

// TODO - RYAN - have a way to destroy the secrets

// For internal use only - used to generate new sets of OCR private keys
// Use NewOCRPrivateKeys in production and NewDeterministicOCRPrivateKeysXXXTestingOnly
// in tests
func newPrivateKeys(reader io.Reader) (*OCRPrivateKeys, error) {
	onChainSk, err := cryptorand.Int(reader, crypto.S256().Params().N)
	if err != nil {
		return nil, err
	}
	onChainPriv := new(signature.OnChainPrivateKey)
	p := (*ecdsa.PrivateKey)(onChainPriv)
	p.D = onChainSk
	onChainPriv.PublicKey = ecdsa.PublicKey{Curve: signature.Curve}
	p.PublicKey.X, p.PublicKey.Y = signature.Curve.ScalarBaseMult(onChainSk.Bytes())
	_, offChainPriv, err := ed25519.GenerateKey(reader)
	if err != nil {
		return nil, err
	}
	var encryptionPriv [curve25519.ScalarSize]byte
	_, err = reader.Read(encryptionPriv[:])
	if err != nil {
		return nil, err
	}
	k := &OCRPrivateKeys{
		onChainSigning:     onChainPriv,
		offChainSigning:    (*signature.OffChainPrivateKey)(&offChainPriv),
		offChainEncryption: &encryptionPriv,
	}
	return k, nil
}

// NewOCRPrivateKeys makes a new set of OCR keys from cryptographically secure entropy
func NewOCRPrivateKeys() (*OCRPrivateKeys, error) {
	return newPrivateKeys(cryptorand.Reader)
}

// NewDeterministicOCRPrivateKeysXXXTestingOnly is for testing purposes only!
func NewDeterministicOCRPrivateKeysXXXTestingOnly(seed int64) (*OCRPrivateKeys, error) {
	return newPrivateKeys(rand.New(rand.NewSource(seed)))
}

// SignOnChain returns an ethereum-style ECDSA secp256k1 signature on msg.
func (pk *OCRPrivateKeys) SignOnChain(msg []byte) (signature []byte, err error) {
	return pk.onChainSigning.Sign(msg)
}

// SignOffChain returns an EdDSA-Ed25519 signature on msg.
func (pk *OCRPrivateKeys) SignOffChain(msg []byte) (signature []byte, err error) {
	return pk.offChainSigning.Sign(msg)
}

// ConfigDiffieHelman returns the shared point obtained by multiplying someone's
// public key by a secret scalar ( in this case, the offChainEncryption key.)
func (pk *OCRPrivateKeys) ConfigDiffieHelman(base *[curve25519.PointSize]byte) (
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
func (pk *OCRPrivateKeys) PublicKeyAddressOnChain() types.OnChainSigningAddress {
	return pk.onChainSigning.Address()
}

// PublicKeyOffChain returns the pbulic component of the keypair used in SignOffChain
func (pk *OCRPrivateKeys) PublicKeyOffChain() types.OffChainPublicKey {
	return types.OffChainPublicKey(pk.offChainSigning.PublicKey())
}

// PublicKeyConfig returns the public component of the keypair used in ConfigKeyShare
func (pk *OCRPrivateKeys) PublicKeyConfig() [curve25519.PointSize]byte {
	rv, err := curve25519.X25519(pk.offChainEncryption[:], curve25519.Basepoint)
	if err != nil {
		log.Println("failure while computing public key: " + err.Error())
	}
	var rvFixed [curve25519.PointSize]byte
	copy(rvFixed[:], rv)
	return rvFixed
}

// type is added to the beginning of the passwords for
// OCR keys, so that the keys can't accidentally be mis-used
// in the wrong place
func adulteratedPassword(auth string) string {
	s := "ocrkey" + auth
	return s
}

// Encrypt combines the OCRPrivateKeys into a single json-serialized
// bytes array and then encrypts
func (pk *OCRPrivateKeys) Encrypt(auth string) (*EncryptedOCRPrivateKeys, error) {
	var marshalledPrivK []byte
	marshalledPrivK, err := json.Marshal(&pk)
	if err != nil {
		return nil, err
	}
	cryptoJSON, err := keystore.EncryptDataV3(
		marshalledPrivK,
		[]byte(adulteratedPassword(auth)),
		defaultScryptParams.N,
		defaultScryptParams.P,
	)
	if err != nil {
		return nil, errors.Wrapf(err, "could not encrypt ocr key")
	}
	encryptedPrivKeys, err := json.Marshal(&cryptoJSON)
	if err != nil {
		return nil, errors.Wrapf(err, "could not encode cryptoJSON")
	}
	return &EncryptedOCRPrivateKeys{
		EncryptedPrivKeys: encryptedPrivKeys,
	}, nil
}

// Decrypt returns the PrivateKeys in e, decrypted via auth, or an error
func (e EncryptedOCRPrivateKeys) Decrypt(auth string) (*OCRPrivateKeys, error) {
	var cryptoJSON keystore.CryptoJSON
	err := json.Unmarshal(e.EncryptedPrivKeys, &cryptoJSON)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid cryptoJSON for OCR key")
	}
	marshalledPrivK, err := keystore.DecryptDataV3(cryptoJSON, adulteratedPassword(auth))
	if err != nil {
		return nil, errors.Wrapf(err, "could not decrypt OCR key")
	}
	var k OCRPrivateKeys
	err = json.Unmarshal(marshalledPrivK, &k)
	if err != nil {
		return nil, errors.Wrapf(err, "could not unmarshal OCR key")
	}
	return &k, nil
}

// MarshalJSON marshals the private keys into json
func (pk *OCRPrivateKeys) MarshalJSON() ([]byte, error) {
	rawKeyData := ocrPrivateKeysRawData{
		X:                  *pk.onChainSigning.X,
		Y:                  *pk.onChainSigning.Y,
		D:                  *pk.onChainSigning.D,
		Ed25519PrivKey:     []byte(*pk.offChainSigning),
		OffChainEncryption: *pk.offChainEncryption,
	}
	return json.Marshal(&rawKeyData)
}

// UnmarshalJSON constructs OCRPrivateKeys from raw json
func (pk *OCRPrivateKeys) UnmarshalJSON(b []byte) (err error) {
	var rawKeyData ocrPrivateKeysRawData
	err = json.Unmarshal(b, &rawKeyData)
	if err != nil {
		return err
	}
	publicKey := ecdsa.PublicKey{
		X: &rawKeyData.X,
		Y: &rawKeyData.Y,
	}
	privateKey := ecdsa.PrivateKey{
		PublicKey: publicKey,
		D:         &rawKeyData.D,
	}
	onChainSigning := signature.OnChainPrivateKey(privateKey)
	offChainSigning := signature.OffChainPrivateKey(rawKeyData.Ed25519PrivKey)
	pk.onChainSigning = &onChainSigning
	pk.offChainSigning = &offChainSigning
	pk.offChainEncryption = &rawKeyData.OffChainEncryption
	return nil
}

// String reduces the risk of accidentally logging the private key
func (pk OCRPrivateKeys) String() string {
	return fmt.Sprintf(
		"OCRPrivateKey{PublicKeyAddressOnChain: %s, PublicKeyOffChain: %s}",
		pk.PublicKeyAddressOnChain(),
		pk.PublicKeyOffChain(),
	)
}

// GoStringer reduces the risk of accidentally logging the private key
func (pk OCRPrivateKeys) GoStringer() string {
	return pk.String()
}

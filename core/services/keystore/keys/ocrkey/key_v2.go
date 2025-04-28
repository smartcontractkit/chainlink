package ocrkey

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/pkg/errors"
	"golang.org/x/crypto/curve25519"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting/types"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/internal"
)

var (
	ErrScalarTooBig = errors.Errorf("can't handle scalars greater than %d", curve25519.PointSize)
	curve           = secp256k1.S256()
)

type keyBundleRawData struct {
	EcdsaD             big.Int
	Ed25519PrivKey     []byte
	OffChainEncryption [curve25519.ScalarSize]byte
}

func KeyFor(raw internal.Raw) KeyV2 {
	var key KeyV2
	err := json.Unmarshal(internal.Bytes(raw), &key)
	if err != nil {
		panic(errors.Wrap(err, "while unmarshalling OCR key"))
	}
	key.raw = raw
	return key
}

type KeyV2 struct {
	raw                internal.Raw
	OnChainSigning     *onChainPrivateKey
	OffChainSigning    *offChainPrivateKey
	OffChainEncryption *[curve25519.ScalarSize]byte
}

func NewV2() (KeyV2, error) {
	ecdsaKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return KeyV2{}, err
	}

	_, offChainPriv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return KeyV2{}, err
	}
	var encryptionPriv [curve25519.ScalarSize]byte
	_, err = rand.Reader.Read(encryptionPriv[:])
	if err != nil {
		return KeyV2{}, err
	}
	k := KeyV2{
		OnChainSigning:     &onChainPrivateKey{func() *ecdsa.PrivateKey { return ecdsaKey }},
		OffChainSigning:    &offChainPrivateKey{func() ed25519.PrivateKey { return offChainPriv }},
		OffChainEncryption: &encryptionPriv,
	}
	return k, k.initRaw()
}

func MustNewV2XXXTestingOnly(k *big.Int) KeyV2 {
	ecdsaKey := new(ecdsa.PrivateKey)
	ecdsaKey.PublicKey.Curve = curve
	ecdsaKey.D = k
	ecdsaKey.PublicKey.X, ecdsaKey.PublicKey.Y = curve.ScalarBaseMult(k.Bytes())
	var seed [32]byte
	copy(seed[:], k.Bytes())
	offChainPriv := ed25519.NewKeyFromSeed(seed[:])
	key := KeyV2{
		OnChainSigning:     &onChainPrivateKey{func() *ecdsa.PrivateKey { return ecdsaKey }},
		OffChainSigning:    &offChainPrivateKey{func() ed25519.PrivateKey { return offChainPriv }},
		OffChainEncryption: &seed,
	}
	if err := key.initRaw(); err != nil {
		panic(err)
	}
	return key
}

func (key KeyV2) ID() string {
	sha := sha256.Sum256(internal.Bytes(key.raw))
	return hex.EncodeToString(sha[:])
}

func (key KeyV2) Raw() internal.Raw { return key.raw }

func (key *KeyV2) initRaw() error {
	marshalledPrivK, err := json.Marshal(key)
	if err != nil {
		return err
	}
	key.raw = internal.NewRaw(marshalledPrivK)
	return nil
}

// SignOnChain returns an ethereum-style ECDSA secp256k1 signature on msg.
func (key KeyV2) SignOnChain(msg []byte) (signature []byte, err error) {
	return key.OnChainSigning.Sign(msg)
}

// SignOffChain returns an EdDSA-Ed25519 signature on msg.
func (key KeyV2) SignOffChain(msg []byte) (signature []byte, err error) {
	return key.OffChainSigning.Sign(msg)
}

// ConfigDiffieHellman returns the shared point obtained by multiplying someone's
// public key by a secret scalar ( in this case, the OffChainEncryption key.)
func (key KeyV2) ConfigDiffieHellman(base *[curve25519.PointSize]byte) (
	sharedPoint *[curve25519.PointSize]byte, err error,
) {
	p, err := curve25519.X25519(key.OffChainEncryption[:], base[:])
	if err != nil {
		return nil, err
	}
	sharedPoint = new([ed25519.PublicKeySize]byte)
	copy(sharedPoint[:], p)
	return sharedPoint, nil
}

// PublicKeyAddressOnChain returns public component of the keypair used in
// SignOnChain
func (key KeyV2) PublicKeyAddressOnChain() ocrtypes.OnChainSigningAddress {
	return ocrtypes.OnChainSigningAddress(key.OnChainSigning.Address())
}

// PublicKeyOffChain returns the public component of the keypair used in SignOffChain
func (key KeyV2) PublicKeyOffChain() ocrtypes.OffchainPublicKey {
	return ocrtypes.OffchainPublicKey(key.OffChainSigning.PublicKey())
}

// PublicKeyConfig returns the public component of the keypair used in ConfigKeyShare
func (key KeyV2) PublicKeyConfig() [curve25519.PointSize]byte {
	rv, err := curve25519.X25519(key.OffChainEncryption[:], curve25519.Basepoint)
	if err != nil {
		log.Println("failure while computing public key: " + err.Error())
	}
	var rvFixed [curve25519.PointSize]byte
	copy(rvFixed[:], rv)
	return rvFixed
}

func (key KeyV2) GetID() string {
	return key.ID()
}

// MarshalJSON marshals the private keys into json
func (key KeyV2) MarshalJSON() ([]byte, error) {
	rawKeyData := keyBundleRawData{
		EcdsaD:             *key.OnChainSigning.pk().D,
		Ed25519PrivKey:     []byte(key.OffChainSigning.pk()),
		OffChainEncryption: *key.OffChainEncryption,
	}
	return json.Marshal(&rawKeyData)
}

func (key *KeyV2) UnmarshalJSON(b []byte) (err error) {
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

	key.OnChainSigning = &onChainPrivateKey{func() *ecdsa.PrivateKey { return &privateKey }}
	key.OffChainSigning = &offChainPrivateKey{func() ed25519.PrivateKey { return rawKeyData.Ed25519PrivKey }}
	key.OffChainEncryption = &rawKeyData.OffChainEncryption
	return nil
}

package ocrkey

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"io"
	"log"

	cryptorand "crypto/rand"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/smartcontractkit/offchain-reporting-design/prototype/offchainreporting/to_be_internal/signature"
	"github.com/smartcontractkit/offchain-reporting-design/prototype/offchainreporting/types"
	"golang.org/x/crypto/curve25519"
)

type PrivateKeys struct {
	onChainSignging    *signature.OnChainPrivateKey
	offChainSigning    *signature.OffChainPrivateKey
	OffChainEncryption *[curve25519.ScalarSize]byte
}

var _ types.PrivateKeys = (*PrivateKeys)(nil)

// NewPrivateKeys returns a PrivateKeys with the given keys.
//
// In any real implementation (except maybe in a test helper), this function
// should take no arguments, and use crypto/rand.{Read,Int}. It should return a
// pointer, so we aren't copying secret material willy-nilly, and have a way to
// destroy the secrets. Any persistence to disk should be encrypted, as in the
// chainlink keystores.
func NewPrivateKeys(reader io.Reader) *PrivateKeys {
	onChainSk, err := cryptorand.Int(reader, crypto.S256().Params().N)
	if err != nil {
		panic(err)
	}
	onChainPriv := new(signature.OnChainPrivateKey)
	p := (*ecdsa.PrivateKey)(onChainPriv)
	p.D = onChainSk
	onChainPriv.PublicKey = ecdsa.PublicKey{Curve: signature.Curve}
	p.PublicKey.X, p.PublicKey.Y = signature.Curve.ScalarBaseMult(onChainSk.Bytes())
	_, offChainPriv, err := ed25519.GenerateKey(reader)
	if err != nil {
		panic(err)
	}
	var encryptionPriv [curve25519.ScalarSize]byte
	_, err = reader.Read(encryptionPriv[:])
	if err != nil {
		panic(err)
	}
	return &PrivateKeys{
		onChainSignging:    onChainPriv,
		offChainSigning:    (*signature.OffChainPrivateKey)(&offChainPriv),
		OffChainEncryption: &encryptionPriv,
	}
}

// SignOnChain returns an ethereum-style ECDSA secp256k1 signature on msg.
func (pk PrivateKeys) SignOnChain(msg []byte) (signature []byte, err error) {
	return pk.onChainSignging.Sign(msg)
}

// SignOffChain returns an EdDSA-Ed25519 signature on msg.
func (pk PrivateKeys) SignOffChain(msg []byte) (signature []byte, err error) {
	return pk.offChainSigning.Sign(msg)
}

// ConfigDiffieHelman returns the shared point obtained by multiplying someone's
// public key by a secret scalar ( in this case, the offChainEncryption key.)
func (pk PrivateKeys) ConfigDiffieHelman(base *[curve25519.PointSize]byte) (
	sharedPoint *[curve25519.PointSize]byte, err error,
) {
	p, err := curve25519.X25519(pk.OffChainEncryption[:], base[:])
	if err != nil {
		return nil, err
	}
	sharedPoint = new([ed25519.PublicKeySize]byte)
	copy(sharedPoint[:], p)
	return sharedPoint, nil
}

// PublicKeyAddressOnChain returns public component of the keypair used in
// SignOnChain
func (pk PrivateKeys) PublicKeyAddressOnChain() types.OnChainSigningAddress {
	return pk.onChainSignging.Address()
}

// PublicKeyOffChain returns the pbulic component of the keypair used in SignOffChain
func (pk PrivateKeys) PublicKeyOffChain() types.OffChainPublicKey {
	return types.OffChainPublicKey(pk.offChainSigning.PublicKey())
}

// PublicKeyConfig returns the public component of the keypair used in ConfigKeyShare
func (pk PrivateKeys) PublicKeyConfig() [curve25519.PointSize]byte {
	rv, err := curve25519.X25519(pk.OffChainEncryption[:], curve25519.Basepoint)
	if err != nil {
		log.Println("failure while computing public key: " + err.Error())
	}
	var rvFixed [curve25519.PointSize]byte
	copy(rvFixed[:], rv)
	return rvFixed
}

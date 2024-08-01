package vss

import (
	"crypto/aes"
	"crypto/cipher"
	"hash"

	"go.dedis.ch/kyber/v3"

	"golang.org/x/crypto/hkdf"
)

// dhExchange computes the shared key from a private key and a public key
func dhExchange(suite Suite, ownPrivate kyber.Scalar, remotePublic kyber.Point) kyber.Point {
	sk := suite.Point()
	sk.Mul(ownPrivate, remotePublic)
	return sk
}

var sharedKeyLength = 32

// newAEAD returns the AEAD cipher to be use to encrypt a share
func newAEAD(fn func() hash.Hash, preSharedKey kyber.Point, context []byte) (cipher.AEAD, error) {
	preBuff, _ := preSharedKey.MarshalBinary()
	reader := hkdf.New(fn, preBuff, nil, context)

	sharedKey := make([]byte, sharedKeyLength)
	if _, err := reader.Read(sharedKey); err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(sharedKey)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return gcm, nil
}

// keySize is arbitrary, make it long enough to seed the XOF
const keySize = 128

// context returns the context slice to be used when encrypting a share
func context(suite Suite, dealer kyber.Point, verifiers []kyber.Point) []byte {
	h := suite.XOF([]byte("vss-dealer"))
	_, _ = dealer.MarshalTo(h)
	_, _ = h.Write([]byte("vss-verifiers"))
	for _, v := range verifiers {
		_, _ = v.MarshalTo(h)
	}
	sum := make([]byte, keySize)
	h.Read(sum)
	return sum
}

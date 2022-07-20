// Package secp256k1 is an implementation of the kyber.{Group,Point,Scalar}
////////////////////////////////////////////////////////////////////////////////
//       XXX: Do not use in production until this code has been audited.
////////////////////////////////////////////////////////////////////////////////
// interfaces, based on btcd/btcec and kyber/group/mod
//
// XXX: NOT CONSTANT TIME!
package secp256k1

import (
	"crypto/cipher"
	"hash"
	"io"
	"reflect"

	"golang.org/x/crypto/sha3"

	"go.dedis.ch/fixbuf"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/util/random"
	"go.dedis.ch/kyber/v3/xof/blake2xb"
)

// SuiteSecp256k1 implements some basic functionalities such as Group, HashFactory,
// and XOFFactory.
type SuiteSecp256k1 struct {
	Secp256k1
	r cipher.Stream
}

// Hash returns a newly instantiated keccak hash function.
func (s *SuiteSecp256k1) Hash() hash.Hash {
	return sha3.NewLegacyKeccak256()
}

// XOF returns an XOR function, implemented via the Blake2b hash.
//
// This should only be used for generating secrets, so there is no need to make
// it cheap to compute on-chain.
func (s *SuiteSecp256k1) XOF(key []byte) kyber.XOF {
	return blake2xb.New(key)
}

// Read implements the Encoding interface function, and reads a series of objs from r
// The objs must all be pointers
func (s *SuiteSecp256k1) Read(r io.Reader, objs ...interface{}) error {
	return fixbuf.Read(r, s, objs...)
}

// Write implements the Encoding interface, and writes the objs to r using their
// built-in binary serializations. Supports Points, Scalars, fixed-length data
// types supported by encoding/binary/Write(), and structs, arrays, and slices
// containing these types.
func (s *SuiteSecp256k1) Write(w io.Writer, objs ...interface{}) error {
	return fixbuf.Write(w, objs)
}

var aScalar kyber.Scalar
var tScalar = reflect.TypeOf(aScalar)
var aPoint kyber.Point
var tPoint = reflect.TypeOf(aPoint)

// New implements the kyber.Encoding interface, and returns a new element of
// type t, which can be a Point or a Scalar
func (s *SuiteSecp256k1) New(t reflect.Type) interface{} {
	switch t {
	case tScalar:
		return s.Scalar()
	case tPoint:
		return s.Point()
	}
	return nil
}

// RandomStream returns a cipher.Stream that returns a key stream
// from crypto/rand.
func (s *SuiteSecp256k1) RandomStream() cipher.Stream {
	if s.r != nil {
		return s.r
	}
	return random.New()
}

// NewBlakeKeccackSecp256k1 returns a cipher suite based on package
// go.dedis.ch/kyber/xof/blake2xb, SHA-256, and the secp256k1 curve. It
// produces cryptographically secure random numbers via package crypto/rand.
func NewBlakeKeccackSecp256k1() *SuiteSecp256k1 {
	return new(SuiteSecp256k1)
}

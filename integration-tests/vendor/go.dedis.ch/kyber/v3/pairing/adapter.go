package pairing

import (
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/pairing/bn256"
)

// SuiteBn256 is an adapter that implements the suites.Suite interface so that
// bn256 can be used as a common suite to generate key pairs for instance but
// still preserves the properties of the pairing (e.g. the Pair function).
//
// It's important to note that the Point function will generate a point
// compatible with public keys only (group G2) where the signature must be
// used as a point from the group G1.
type SuiteBn256 struct {
	Suite
	kyber.Group
}

// NewSuiteBn256 makes a new BN256 suite
func NewSuiteBn256() *SuiteBn256 {
	return &SuiteBn256{
		Suite: bn256.NewSuite(),
	}
}

// Point generates a point from the G2 group that can only be used
// for public keys
func (s *SuiteBn256) Point() kyber.Point {
	return s.G2().Point()
}

// PointLen returns the length of a G2 point
func (s *SuiteBn256) PointLen() int {
	return s.G2().PointLen()
}

// Scalar generates a scalar
func (s *SuiteBn256) Scalar() kyber.Scalar {
	return s.G1().Scalar()
}

// ScalarLen returns the lenght of a scalar
func (s *SuiteBn256) ScalarLen() int {
	return s.G1().ScalarLen()
}

// String returns the name of the suite
func (s *SuiteBn256) String() string {
	return "bn256.adapter"
}

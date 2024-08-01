package fp

import (
	basefield "github.com/consensys/gnark-crypto/ecc/bls12-381/fr"
)

// Element is a field element representing the basefield of the curve.
type Element = basefield.Element

// Limbs is the number of 64-bit words needed to represent the field element.
const Limbs = basefield.Limbs

// One returns the field element 1.
func One() Element {
	return basefield.One()
}

// Zero returns the field element 0.
func Zero() Element {
	return basefield.Element{}
}

// MinusOne returns the field element -1.
func MinusOne() Element {
	m_one := One()
	m_one.Neg(&m_one)
	return m_one
}

// MulBy5 multiplies a field element by 5.
func MulBy5(a *Element) {
	basefield.MulBy5(a)
}

// BatchInvert computes the inverse of a slice of field elements.
func BatchInvert(a []Element) []Element {
	return basefield.BatchInvert(a)
}

// BytesLE returns the little-endian byte representation of a field element.
func BytesLE(a Element) []byte {
	var result [basefield.Bytes]byte
	basefield.LittleEndian.PutElement(&result, a)
	return result[:]
}

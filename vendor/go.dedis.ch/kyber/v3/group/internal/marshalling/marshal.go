// Package marshalling provides a common implementation of (un)marshalling method using Writer and Reader.
//
package marshalling

import (
	"crypto/cipher"
	"io"
	"reflect"

	"go.dedis.ch/kyber/v3"
)

// PointMarshalTo provides a generic implementation of Point.EncodeTo
// based on Point.Encode.
func PointMarshalTo(p kyber.Point, w io.Writer) (int, error) {
	buf, err := p.MarshalBinary()
	if err != nil {
		return 0, err
	}
	return w.Write(buf)
}

// PointUnmarshalFrom provides a generic implementation of Point.DecodeFrom,
// based on Point.Decode, or Point.Pick if r is a Cipher or cipher.Stream.
// The returned byte-count is valid only when decoding from a normal Reader,
// not when picking from a pseudorandom source.
func PointUnmarshalFrom(p kyber.Point, r io.Reader) (int, error) {
	if strm, ok := r.(cipher.Stream); ok {
		p.Pick(strm)
		return -1, nil // no byte-count when picking randomly
	}
	buf := make([]byte, p.MarshalSize())
	n, err := io.ReadFull(r, buf)
	if err != nil {
		return n, err
	}
	return n, p.UnmarshalBinary(buf)
}

// ScalarMarshalTo provides a generic implementation of Scalar.EncodeTo
// based on Scalar.Encode.
func ScalarMarshalTo(s kyber.Scalar, w io.Writer) (int, error) {
	buf, err := s.MarshalBinary()
	if err != nil {
		return 0, err
	}
	return w.Write(buf)
}

// ScalarUnmarshalFrom provides a generic implementation of Scalar.DecodeFrom,
// based on Scalar.Decode, or Scalar.Pick if r is a Cipher or cipher.Stream.
// The returned byte-count is valid only when decoding from a normal Reader,
// not when picking from a pseudorandom source.
func ScalarUnmarshalFrom(s kyber.Scalar, r io.Reader) (int, error) {
	if strm, ok := r.(cipher.Stream); ok {
		s.Pick(strm)
		return -1, nil // no byte-count when picking randomly
	}
	buf := make([]byte, s.MarshalSize())
	n, err := io.ReadFull(r, buf)
	if err != nil {
		return n, err
	}
	return n, s.UnmarshalBinary(buf)
}

// Not used other than for reflect.TypeOf()
var aScalar kyber.Scalar
var aPoint kyber.Point

var tScalar = reflect.TypeOf(&aScalar).Elem()
var tPoint = reflect.TypeOf(&aPoint).Elem()

// GroupNew is the Default implementation of reflective constructor for Group
func GroupNew(g kyber.Group, t reflect.Type) interface{} {
	switch t {
	case tScalar:
		return g.Scalar()
	case tPoint:
		return g.Point()
	}
	return nil
}

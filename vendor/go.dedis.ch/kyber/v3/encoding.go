package kyber

import (
	"encoding"
	"io"
)

/*
Marshaling is a basic interface representing fixed-length (or known-length)
cryptographic objects or structures having a built-in binary encoding.
Implementors must ensure that calls to these methods do not modify
the underlying object so that other users of the object can access
it concurrently.
*/
type Marshaling interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler

	// String returns the human readable string representation of the object.
	String() string

	// Encoded length of this object in bytes.
	MarshalSize() int

	// Encode the contents of this object and write it to an io.Writer.
	MarshalTo(w io.Writer) (int, error)

	// Decode the content of this object by reading from an io.Reader.
	// If r is an XOF, it uses r to pick a valid object pseudo-randomly,
	// which may entail reading more than Len bytes due to retries.
	UnmarshalFrom(r io.Reader) (int, error)
}

// Encoding represents an abstract interface to an encoding/decoding that can be
// used to marshal/unmarshal objects to and from streams. Different Encodings
// will have different constraints, of course. Two implementations are
// available:
//
//   1. The protobuf encoding using the variable length Google Protobuf encoding
//      scheme. The library is available at https://go.dedis.ch/protobuf
//   2. The fixbuf encoding, a fixed length binary encoding of arbitrary
//      structures. The library is available at https://go.dedis.ch/fixbuf.
type Encoding interface {
	// Encode and write objects to an io.Writer.
	Write(w io.Writer, objs ...interface{}) error

	// Read and decode objects from an io.Reader.
	Read(r io.Reader, objs ...interface{}) error
}

package common

import (
	"bytes"
	"crypto/sha256"
	"hash"

	"github.com/crate-crypto/go-ipa/bandersnatch/fr"
	"github.com/crate-crypto/go-ipa/banderwagon"
)

// The transcript is used to create challenge scalars.
// See: Fiat-Shamir
type Transcript struct {
	state hash.Hash
	buff  *bytes.Buffer
}

func NewTranscript(label string) *Transcript {
	digest := sha256.New()
	digest.Write([]byte(label))

	transcript := &Transcript{
		state: digest,
		buff:  bytes.NewBuffer(make([]byte, 0, 1024)),
	}

	return transcript
}

func (t *Transcript) AppendMessage(message []byte, label []byte) {
	t.buff.Write(label)
	t.buff.Write(message)
}

// Appends a Scalar to the transcript
//
// Converts the scalar to 32 bytes, then appends it to
// the state
func (t *Transcript) AppendScalar(scalar *fr.Element, label []byte) {
	tmpBytes := scalar.BytesLE()
	t.AppendMessage(tmpBytes[:], label)

}

// Appends a Point to the transcript
//
// Compresses the Point into a 32 byte slice, then appends it to
// the state
func (t *Transcript) AppendPoint(point *banderwagon.Element, label []byte) {
	tmp_bytes := point.Bytes()
	t.AppendMessage(tmp_bytes[:], label)

}

func (t *Transcript) DomainSep(label []byte) {
	t.buff.Write(label)
}

// Computes a challenge based off of the state of the transcript
//
// Hash the transcript state, then reduce the hash modulo the size of the
// scalar field
//
// Note that calling the transcript twice, will yield two different challenges
func (t *Transcript) ChallengeScalar(label []byte) fr.Element {
	t.DomainSep(label)

	t.state.Write(t.buff.Bytes())
	t.buff.Reset()
	// Reverse the endian so we are using little-endian
	// SetBytes interprets the bytes in Big Endian
	bytes := t.state.Sum(nil)

	var tmp fr.Element
	tmp.SetBytesLE(bytes)

	// Clear the state
	t.state.Reset()

	// Add the new challenge to the state
	// Which "summarises" the previous state before we cleared it
	t.AppendScalar(&tmp, label)

	// Return the new challenge
	return tmp
}

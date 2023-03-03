package types

import (
	"crypto/ed25519"
	"encoding"
	"fmt"

	"github.com/mr-tron/base58"
	"github.com/smartcontractkit/libocr/ragep2p/internal/mtls"
)

// Address represents a network address & port such as "192.168.1.2:8080". It
// can also contain special bind addresses such as "0.0.0.0:80".
type Address string

// PeerID represents a unique identifier for another peer.
type PeerID [32]byte

var (
	_ fmt.Stringer               = PeerID{}
	_ encoding.TextMarshaler     = PeerID{}
	_ encoding.TextUnmarshaler   = &PeerID{}
	_ encoding.BinaryMarshaler   = PeerID{}
	_ encoding.BinaryUnmarshaler = &PeerID{}
)

func (p PeerID) String() string {
	text, err := p.MarshalText()
	if err != nil {
		return fmt.Sprintf("<PeerID: failed to stringify due to error '%s'>", err)
	}
	return string(text)
}

func (p PeerID) MarshalText() (text []byte, err error) {
	bin, err := p.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return []byte(base58.Encode(bin)), nil
}

func (p *PeerID) UnmarshalText(text []byte) error {
	b58, err := base58.Decode(string(text))
	if err != nil {
		return fmt.Errorf("failed to base58 decode: %w", err)
	}
	return p.UnmarshalBinary(b58)
}

func (p PeerID) MarshalBinary() (data []byte, err error) {
	// this magic value comes from libp2p's encoding of peer ids, see
	// https://docs.libp2p.io/concepts/peer-id/
	return append([]byte{0x00, 0x24, 0x08, 0x01, 0x12, 0x20}, p[:]...), nil
}

func (p *PeerID) UnmarshalBinary(data []byte) error {
	const libp2pMagicLength = 6
	const expectedLength = 32 + libp2pMagicLength
	if len(data) != expectedLength {
		return fmt.Errorf("wrong size of data (was %d, expected %d)", len(data), expectedLength)
	}
	copy(p[:], data[libp2pMagicLength:])
	return nil
}

func PeerIDFromPublicKey(pk ed25519.PublicKey) (PeerID, error) {
	return mtls.StaticallySizedEd25519PublicKey(pk)
}

func PeerIDFromPrivateKey(sk ed25519.PrivateKey) (PeerID, error) {
	return PeerIDFromPublicKey(sk.Public().(ed25519.PublicKey))
}

type PeerInfo struct {
	ID    PeerID
	Addrs []Address
}

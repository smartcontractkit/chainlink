package ragep2p

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/smartcontractkit/libocr/ragep2p/types"
)

var errWrongLength = fmt.Errorf("frameHeader must have exactly %v bytes", frameHeaderEncodedSize)
var errUnknownFrameType = errors.New("frameHeader has unknown frameType")

type frameType uint8

const (
	_ frameType = iota
	frameTypeOpen
	frameTypeClose
	frameTypeData
)

type frameHeader struct {
	Type          frameType
	StreamID      streamID
	PayloadLength uint32
}

const frameHeaderEncodedSize = 1 + 32 + 4

func (fh frameHeader) Encode() []byte {
	buf := bytes.NewBuffer(make([]byte, 0, frameHeaderEncodedSize))
	buf.WriteByte(byte(fh.Type))
	buf.Write(fh.StreamID[:])
	binary.Write(buf, binary.BigEndian, fh.PayloadLength) //nolint:errcheck
	return buf.Bytes()
}

func decodeFrameHeader(encoded []byte) (frameHeader, error) {
	if len(encoded) != frameHeaderEncodedSize {
		return frameHeader{}, errWrongLength
	}
	typ := frameType(encoded[0])
	switch typ {
	case frameTypeOpen:
	case frameTypeClose:
	case frameTypeData:
	default:
		return frameHeader{}, errUnknownFrameType
	}
	var streamId streamID
	copy(streamId[:], encoded[1:33])
	payloadLength := binary.BigEndian.Uint32(encoded[33:frameHeaderEncodedSize])
	return frameHeader{
		typ,
		streamId,
		payloadLength,
	}, nil
}

func getStreamID(self types.PeerID, other types.PeerID, name string) streamID {
	if bytes.Compare(self[:], other[:]) < 0 {
		return getStreamID(other, self, name)
	}

	h := sha256.New()
	h.Write(self[:])
	h.Write(other[:])
	// this is fine because self and other are of constant length. if more than
	// one variable length item is ever added here, we should also hash lengths
	// to prevent collisions.
	h.Write([]byte(name))

	var result streamID
	copy(result[:], h.Sum(nil))
	return result
}

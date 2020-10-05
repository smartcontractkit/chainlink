package signature

import (
	"encoding/binary"

	"github.com/smartcontractkit/chainlink/libocr/offchainreporting/types"
)

type DomainSeparationTag [32]byte

type ReportingContext struct {
	ConfigDigest types.ConfigDigest
	Epoch        uint32
	Round        uint8
}

func NewReportingContext(configDigest types.ConfigDigest, epoch uint32, round uint8) (r ReportingContext) {
	return ReportingContext{
		ConfigDigest: configDigest,
		Epoch:        epoch,
		Round:        round,
	}
}

func (r ReportingContext) DomainSeparationTag() (d DomainSeparationTag) {
	serialization := r.ConfigDigest[:]
	serialization = append(serialization, WireUInt32(uint32(r.Epoch))...)
	serialization = append(serialization, byte(r.Round))
	copy(d[11:], serialization)
	return d
}

func (r ReportingContext) Equal(r2 ReportingContext) bool {
	return r.ConfigDigest == r2.ConfigDigest && r.Epoch == r2.Epoch && r.Round == r2.Round
}

func WireUInt32(i uint32) (serialization []byte) {
	serialization = make([]byte, 4)
	binary.BigEndian.PutUint32(serialization[:], i)
	return serialization
}

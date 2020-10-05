package protocol

import (
	"bytes"

	"github.com/smartcontractkit/chainlink/libocr/offchainreporting/internal/protocol/observation"
	"github.com/smartcontractkit/chainlink/libocr/offchainreporting/internal/signature"
	"github.com/smartcontractkit/chainlink/libocr/offchainreporting/types"
)

type Observation struct {
	Ctx      signature.ReportingContext
	OracleID types.OracleID
	Value    observation.Observation 	Sig      []byte                  }

func (obs Observation) wireMessage() []byte {
	tag := obs.Ctx.DomainSeparationTag()
	return append(tag[:], obs.Value.Marshal()...)
}

func (obs Observation) Equal(o2 Observation) bool {
	return obs.Ctx.Equal(o2.Ctx) &&
		obs.OracleID == o2.OracleID &&
		obs.Value.Equal(o2.Value) &&
		bytes.Equal(obs.Sig, o2.Sig)
}

package types

import (
	"context"

	"github.com/smartcontractkit/chainlink-common/pkg/values"
)

type Encoder interface {
	Encode(ctx context.Context, input values.Map) ([]byte, error)
}

type EncoderFactory func(config *values.Map) (Encoder, error)

type SignedReport struct {
	Report []byte
	// Report context is appended to the report before signing by libOCR.
	// It contains config digest + round/epoch/sequence numbers (currently 96 bytes).
	// Has to be appended to the report before validating signatures.
	Context []byte
	// Always exactly F+1 signatures.
	Signatures [][]byte
	// Report ID defined in the workflow spec (2 bytes).
	ID []byte
}

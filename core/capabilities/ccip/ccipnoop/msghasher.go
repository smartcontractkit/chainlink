package ccipnoop

import (
	"context"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
)

// NoopMessageHasherV1 implements the MessageHasher interface.
// Compatible with:
// - "OnRamp 1.6.0-dev"
type NoopMessageHasherV1 struct {
	lggr           logger.Logger
	extraDataCodec common.ExtraDataCodec
}

func NewNoopMessageHasherV1(lggr logger.Logger, extraDataCodec common.ExtraDataCodec) *NoopMessageHasherV1 {
	return &NoopMessageHasherV1{
		lggr:           lggr,
		extraDataCodec: extraDataCodec,
	}
}

// Hash implements the MessageHasher interface.
func (h *NoopMessageHasherV1) Hash(_ context.Context, msg cciptypes.Message) (cciptypes.Bytes32, error) {
	return [32]byte{}, nil
}

// Interface compliance check
var _ cciptypes.MessageHasher = (*NoopMessageHasherV1)(nil)

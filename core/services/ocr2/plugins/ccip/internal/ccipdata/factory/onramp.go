package factory

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/v1_0_0"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/v1_1_0"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/v1_2_0"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/v1_3_0"
)

// NewOnRampReader determines the appropriate version of the onramp and returns a reader for it
func NewOnRampReader(lggr logger.Logger, sourceSelector, destSelector uint64, onRampAddress common.Address, sourceLP logpoller.LogPoller, source client.Client) (ccipdata.OnRampReader, error) {
	contractType, version, err := ccipconfig.TypeAndVersion(onRampAddress, source)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to read type and version")
	}
	if contractType != ccipconfig.EVM2EVMOnRamp {
		return nil, errors.Errorf("expected %v got %v", ccipconfig.EVM2EVMOnRamp, contractType)
	}
	switch version.String() {
	case ccipdata.V1_0_0:
		return v1_0_0.NewOnRamp(lggr, sourceSelector, destSelector, onRampAddress, sourceLP, source)
	case ccipdata.V1_1_0:
		return v1_1_0.NewOnRamp(lggr, sourceSelector, destSelector, onRampAddress, sourceLP, source)
	case ccipdata.V1_2_0:
		return v1_2_0.NewOnRamp(lggr, sourceSelector, destSelector, onRampAddress, sourceLP, source)
	case ccipdata.V1_3_0:
		return v1_3_0.NewOnRamp(lggr, sourceSelector, destSelector, onRampAddress, sourceLP, source)
	default:
		return nil, errors.Errorf("unsupported onramp version %v", version.String())
	}
}

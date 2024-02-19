package factory

import (
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/cciptypes"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcalc"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/v1_0_0"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/v1_1_0"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/v1_2_0"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/v1_5_0"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

// NewOnRampReader determines the appropriate version of the onramp and returns a reader for it
func NewOnRampReader(lggr logger.Logger, versionFinder VersionFinder, sourceSelector, destSelector uint64, onRampAddress cciptypes.Address, sourceLP logpoller.LogPoller, source client.Client, pgOpts ...pg.QOpt) (ccipdata.OnRampReader, error) {
	return initOrCloseOnRampReader(lggr, versionFinder, sourceSelector, destSelector, onRampAddress, sourceLP, source, false, pgOpts...)
}

func CloseOnRampReader(lggr logger.Logger, versionFinder VersionFinder, sourceSelector, destSelector uint64, onRampAddress cciptypes.Address, sourceLP logpoller.LogPoller, source client.Client, pgOpts ...pg.QOpt) error {
	_, err := initOrCloseOnRampReader(lggr, versionFinder, sourceSelector, destSelector, onRampAddress, sourceLP, source, true, pgOpts...)
	return err
}

func initOrCloseOnRampReader(lggr logger.Logger, versionFinder VersionFinder, sourceSelector, destSelector uint64, onRampAddress cciptypes.Address, sourceLP logpoller.LogPoller, source client.Client, closeReader bool, pgOpts ...pg.QOpt) (ccipdata.OnRampReader, error) {
	contractType, version, err := versionFinder.TypeAndVersion(onRampAddress, source)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to read type and version")
	}
	if contractType != ccipconfig.EVM2EVMOnRamp {
		return nil, errors.Errorf("expected %v got %v", ccipconfig.EVM2EVMOnRamp, contractType)
	}

	onRampAddrEvm, err := ccipcalc.GenericAddrToEvm(onRampAddress)
	if err != nil {
		return nil, err
	}

	switch version.String() {
	case ccipdata.V1_0_0:
		onRamp, err := v1_0_0.NewOnRamp(lggr, sourceSelector, destSelector, onRampAddrEvm, sourceLP, source)
		if err != nil {
			return nil, err
		}
		if closeReader {
			return nil, onRamp.Close(pgOpts...)
		}
		return onRamp, onRamp.RegisterFilters(pgOpts...)
	case ccipdata.V1_1_0:
		onRamp, err := v1_1_0.NewOnRamp(lggr, sourceSelector, destSelector, onRampAddrEvm, sourceLP, source)
		if err != nil {
			return nil, err
		}
		if closeReader {
			return nil, onRamp.Close(pgOpts...)
		}
		return onRamp, onRamp.RegisterFilters(pgOpts...)
	case ccipdata.V1_2_0:
		onRamp, err := v1_2_0.NewOnRamp(lggr, sourceSelector, destSelector, onRampAddrEvm, sourceLP, source)
		if err != nil {
			return nil, err
		}
		if closeReader {
			return nil, onRamp.Close(pgOpts...)
		}
		return onRamp, onRamp.RegisterFilters(pgOpts...)
	case ccipdata.V1_5_0:
		onRamp, err := v1_5_0.NewOnRamp(lggr, sourceSelector, destSelector, onRampAddrEvm, sourceLP, source)
		if err != nil {
			return nil, err
		}
		if closeReader {
			return nil, onRamp.Close(pgOpts...)
		}
		return onRamp, onRamp.RegisterFilters(pgOpts...)
	default:
		return nil, errors.Errorf("unsupported onramp version %v", version.String())
	}
}

package factory

import (
	"context"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"

	"github.com/smartcontractkit/chainlink-evm/pkg/client"
	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcalc"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/v1_2_0"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/v1_5_0"
)

// NewOnRampReader determines the appropriate version of the onramp and returns a reader for it
func NewOnRampReader(ctx context.Context, lggr logger.Logger, versionFinder VersionFinder, sourceSelector, destSelector uint64, onRampAddress cciptypes.Address, sourceLP logpoller.LogPoller, source client.Client) (ccipdata.OnRampReader, error) {
	return initOrCloseOnRampReader(ctx, lggr, versionFinder, sourceSelector, destSelector, onRampAddress, sourceLP, source, false)
}

func CloseOnRampReader(ctx context.Context, lggr logger.Logger, versionFinder VersionFinder, sourceSelector, destSelector uint64, onRampAddress cciptypes.Address, sourceLP logpoller.LogPoller, source client.Client) error {
	_, err := initOrCloseOnRampReader(ctx, lggr, versionFinder, sourceSelector, destSelector, onRampAddress, sourceLP, source, true)
	return err
}

func initOrCloseOnRampReader(ctx context.Context, lggr logger.Logger, versionFinder VersionFinder, sourceSelector, destSelector uint64, onRampAddress cciptypes.Address, sourceLP logpoller.LogPoller, source client.Client, closeReader bool) (ccipdata.OnRampReader, error) {
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

	lggr.Infof("Initializing onRamp for version %v", version.String())

	switch version.String() {
	case ccipdata.V1_2_0:
		onRamp, err := v1_2_0.NewOnRamp(lggr, sourceSelector, destSelector, onRampAddrEvm, sourceLP, source)
		if err != nil {
			return nil, err
		}
		if closeReader {
			return nil, onRamp.Close()
		}
		return onRamp, onRamp.RegisterFilters(ctx)
	case ccipdata.V1_5_0:
		onRamp, err := v1_5_0.NewOnRamp(lggr, sourceSelector, destSelector, onRampAddrEvm, sourceLP, source)
		if err != nil {
			return nil, err
		}
		if closeReader {
			return nil, onRamp.Close()
		}
		return onRamp, onRamp.RegisterFilters(ctx)
	// Adding a new version?
	// Please update the public factory function in leafer.go if the new version updates the leaf hash function.
	default:
		return nil, errors.Errorf("unsupported onramp version %v", version.String())
	}
}

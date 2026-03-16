package factory

import (
	"context"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/ccip/ccipcalc"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/ccip/version"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/ccip/versionfinder"

	"github.com/smartcontractkit/chainlink-evm/pkg/client"
	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/v1_2_0"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/v1_5_0"
)

func NewCommitStoreReader(ctx context.Context, lggr logger.Logger, versionFinder versionfinder.VersionFinder, address cciptypes.Address, ec client.Client, lp logpoller.LogPoller, feeEstimatorConfig ccipdata.FeeEstimatorConfigReader) (ccipdata.CommitStoreReader, error) {
	return initOrCloseCommitStoreReader(ctx, lggr, versionFinder, address, ec, lp, feeEstimatorConfig, false)
}

func CloseCommitStoreReader(ctx context.Context, lggr logger.Logger, versionFinder versionfinder.VersionFinder, address cciptypes.Address, ec client.Client, lp logpoller.LogPoller, feeEstimatorConfig ccipdata.FeeEstimatorConfigReader) error {
	_, err := initOrCloseCommitStoreReader(ctx, lggr, versionFinder, address, ec, lp, feeEstimatorConfig, true)
	return err
}

func initOrCloseCommitStoreReader(ctx context.Context, lggr logger.Logger, versionFinder versionfinder.VersionFinder, address cciptypes.Address, ec client.Client, lp logpoller.LogPoller, feeEstimatorConfig ccipdata.FeeEstimatorConfigReader, closeReader bool) (ccipdata.CommitStoreReader, error) {
	contractType, version, err := versionFinder.TypeAndVersion(address, ec)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to read type and version")
	}
	if contractType != ccipconfig.CommitStore {
		return nil, errors.Errorf("expected %v got %v", ccipconfig.CommitStore, contractType)
	}

	evmAddr, err := ccipcalc.GenericAddrToEvm(address)
	if err != nil {
		return nil, err
	}

	lggr.Infow("Initializing CommitStore Reader", "version", version.String())

	switch version.String() {
	case ccipdata.V1_2_0:
		cs, err := v1_2_0.NewCommitStore(lggr, evmAddr, ec, lp, feeEstimatorConfig)
		if err != nil {
			return nil, err
		}
		if closeReader {
			return nil, cs.Close()
		}
		return cs, cs.RegisterFilters(ctx)
	case ccipdata.V1_5_0:
		cs, err := v1_5_0.NewCommitStore(lggr, evmAddr, ec, lp, feeEstimatorConfig)
		if err != nil {
			return nil, err
		}
		if closeReader {
			return nil, cs.Close()
		}
		return cs, cs.RegisterFilters(ctx)
	default:
		return nil, errors.Errorf("unsupported commit store version %v", version.String())
	}
}

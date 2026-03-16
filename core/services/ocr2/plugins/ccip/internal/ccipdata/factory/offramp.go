package factory

import (
	"context"
	"math/big"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/ccip/ccipcalc"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/ccip/version"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/ccip/versionfinder"

	"github.com/smartcontractkit/chainlink-evm/pkg/client"
	"github.com/smartcontractkit/chainlink-evm/pkg/gas"
	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/v1_2_0"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/v1_5_0"
)

func NewOffRampReader(ctx context.Context, lggr logger.Logger, versionFinder versionfinder.VersionFinder, addr cciptypes.Address, destClient client.Client, lp logpoller.LogPoller, estimator gas.EvmFeeEstimator, destMaxGasPrice *big.Int, registerFilters bool, feeEstimatorConfig ccipdata.FeeEstimatorConfigReader) (ccipdata.OffRampReader, error) {
	return initOrCloseOffRampReader(ctx, lggr, versionFinder, addr, destClient, lp, estimator, destMaxGasPrice, false, registerFilters, feeEstimatorConfig)
}

func CloseOffRampReader(ctx context.Context, lggr logger.Logger, versionFinder versionfinder.VersionFinder, addr cciptypes.Address, destClient client.Client, lp logpoller.LogPoller, estimator gas.EvmFeeEstimator, destMaxGasPrice *big.Int, feeEstimatorConfig ccipdata.FeeEstimatorConfigReader) error {
	_, err := initOrCloseOffRampReader(ctx, lggr, versionFinder, addr, destClient, lp, estimator, destMaxGasPrice, true, false, feeEstimatorConfig)
	return err
}

func initOrCloseOffRampReader(ctx context.Context, lggr logger.Logger, versionFinder versionfinder.VersionFinder, addr cciptypes.Address, destClient client.Client, lp logpoller.LogPoller, estimator gas.EvmFeeEstimator, destMaxGasPrice *big.Int, closeReader bool, registerFilters bool, feeEstimatorConfig ccipdata.FeeEstimatorConfigReader) (ccipdata.OffRampReader, error) {
	contractType, version, err := versionFinder.TypeAndVersion(addr, destClient)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to read type and version")
	}
	if contractType != ccipconfig.EVM2EVMOffRamp {
		return nil, errors.Errorf("expected %v got %v", ccipconfig.EVM2EVMOffRamp, contractType)
	}

	evmAddr, err := ccipcalc.GenericAddrToEvm(addr)
	if err != nil {
		return nil, err
	}

	lggr.Infow("Initializing OffRamp Reader", "version", version.String(), "destMaxGasPrice", destMaxGasPrice.String())

	switch version.String() {
	case ccipdata.V1_2_0:
		offRamp, err := v1_2_0.NewOffRamp(lggr, evmAddr, destClient, lp, estimator, destMaxGasPrice, feeEstimatorConfig)
		if err != nil {
			return nil, err
		}
		if closeReader {
			return nil, offRamp.Close()
		}
		return offRamp, offRamp.RegisterFilters(ctx)
	case ccipdata.V1_5_0:
		offRamp, err := v1_5_0.NewOffRamp(lggr, evmAddr, destClient, lp, estimator, destMaxGasPrice, feeEstimatorConfig)
		if err != nil {
			return nil, err
		}
		if closeReader {
			return nil, offRamp.Close()
		}
		return offRamp, offRamp.RegisterFilters(ctx)
	default:
		return nil, errors.Errorf("unsupported offramp version %v", version.String())
	}
	// TODO can validate it pointing to the correct version
}

package export

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
	"github.com/smartcontractkit/chainlink-evm/pkg/client"
	"github.com/smartcontractkit/chainlink-evm/pkg/gas"
	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/ccip/ccipdata"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/ccip/ccipdata/factory"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/ccip/commitstore"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/ccip/observability"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/ccip/rpclib"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/ccip/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/ccip/versionfinder"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/ccipdata/batchreader"
)

type JSONCommitOffchainConfigV1_2_0 = commitstore.JSONCommitOffchainConfig
type CommitOnchainConfig = commitstore.CommitOnchainConfig

func NewCommitOffchainConfig(
	gasPriceDeviationPPB uint32,
	gasPriceHeartBeat time.Duration,
	tokenPriceDeviationPPB uint32,
	tokenPriceHeartBeat time.Duration,
	inflightCacheExpiry time.Duration,
	priceReportingDisabled bool,
) ccip.CommitOffchainConfig {
	return commitstore.NewCommitOffchainConfig(gasPriceDeviationPPB, gasPriceHeartBeat, tokenPriceDeviationPPB, tokenPriceHeartBeat, inflightCacheExpiry, priceReportingDisabled)
}

func NewCommitStoreReader(ctx context.Context, lggr logger.Logger, versionFinder versionfinder.VersionFinder, address ccip.Address, ec client.Client, lp logpoller.LogPoller, feeEstimatorConfig ccipdata.FeeEstimatorConfigReader) (types.CommitStoreReader, error) {
	return factory.NewCommitStoreReader(ctx, lggr, versionFinder, address, ec, lp, feeEstimatorConfig)
}

func CloseCommitStoreReader(ctx context.Context, lggr logger.Logger, versionFinder versionfinder.VersionFinder, address ccip.Address, ec client.Client, lp logpoller.LogPoller, feeEstimatorConfig ccipdata.FeeEstimatorConfigReader) error {
	return factory.CloseCommitStoreReader(ctx, lggr, versionFinder, address, ec, lp, feeEstimatorConfig)
}

type PriceRegistry interface {
	NewPriceRegistryReader(ctx context.Context, addr ccip.Address) (ccip.PriceRegistryReader, error)
}

type EvmPriceRegistry struct {
	lp          logpoller.LogPoller
	ec          client.Client
	lggr        logger.Logger
	pluginLabel string
}

func NewEvmPriceRegistry(lp logpoller.LogPoller, ec client.Client, lggr logger.Logger, pluginLabel string) *EvmPriceRegistry {
	return &EvmPriceRegistry{
		lp:          lp,
		ec:          ec,
		lggr:        lggr,
		pluginLabel: pluginLabel,
	}
}

func (p *EvmPriceRegistry) NewPriceRegistryReader(ctx context.Context, addr ccip.Address) (ccip.PriceRegistryReader, error) {
	destPriceRegistryReader, err := factory.NewPriceRegistryReader(ctx, p.lggr, versionfinder.NewEvmVersionFinder(), addr, p.lp, p.ec)
	if err != nil {
		return nil, err
	}
	return observability.NewPriceRegistryReader(destPriceRegistryReader, p.ec.ConfiguredChainID().Int64(), p.pluginLabel), nil
}

func NewOffRampReader(ctx context.Context, lggr logger.Logger, versionFinder versionfinder.VersionFinder, addr ccip.Address, destClient client.Client, lp logpoller.LogPoller, estimator gas.EvmFeeEstimator, destMaxGasPrice *big.Int, registerFilters bool, feeEstimatorConfig ccipdata.FeeEstimatorConfigReader) (types.OffRampReader, error) {
	return factory.NewOffRampReader(ctx, lggr, versionFinder, addr, destClient, lp, estimator, destMaxGasPrice, registerFilters, feeEstimatorConfig)
}

func CloseOffRampReader(ctx context.Context, lggr logger.Logger, versionFinder versionfinder.VersionFinder, addr ccip.Address, destClient client.Client, lp logpoller.LogPoller, estimator gas.EvmFeeEstimator, destMaxGasPrice *big.Int, feeEstimatorConfig ccipdata.FeeEstimatorConfigReader) error {
	return factory.CloseOffRampReader(ctx, lggr, versionFinder, addr, destClient, lp, estimator, destMaxGasPrice, feeEstimatorConfig)
}

func NewOnRampReader(ctx context.Context, lggr logger.Logger, versionFinder versionfinder.VersionFinder, sourceSelector, destSelector uint64, onRampAddress ccip.Address, sourceLP logpoller.LogPoller, source client.Client) (types.OnRampReader, error) {
	return factory.NewOnRampReader(ctx, lggr, versionFinder, sourceSelector, destSelector, onRampAddress, sourceLP, source)
}

func CloseOnRampReader(ctx context.Context, lggr logger.Logger, versionFinder versionfinder.VersionFinder, sourceSelector, destSelector uint64, onRampAddress ccip.Address, sourceLP logpoller.LogPoller, source client.Client) error {
	return factory.CloseOnRampReader(ctx, lggr, versionFinder, sourceSelector, destSelector, onRampAddress, sourceLP, source)
}

type OffRampReader = types.OffRampReader

func NewUSDCReader(ctx context.Context, lggr logger.Logger, jobID string, transmitter common.Address, lp logpoller.LogPoller, registerFilters bool) (*ccipdata.USDCReaderImpl, error) {
	return ccipdata.NewUSDCReader(ctx, lggr, jobID, transmitter, lp, registerFilters)
}

func CloseUSDCReader(ctx context.Context, lggr logger.Logger, jobID string, transmitter common.Address, lp logpoller.LogPoller) error {
	return ccipdata.CloseUSDCReader(ctx, lggr, jobID, transmitter, lp)
}

type USDCReaderImpl = ccipdata.USDCReaderImpl

func NewDynamicLimitedBatchCaller(
	lggr logger.Logger, batchSender rpclib.BatchSender, batchSizeLimit, backOffMultiplier, parallelRpcCallsLimit uint,
) *rpclib.DynamicLimitedBatchCaller {
	return rpclib.NewDynamicLimitedBatchCaller(lggr, batchSender, batchSizeLimit, backOffMultiplier, parallelRpcCallsLimit)
}

func NewEVMTokenPoolBatchedReader(lggr logger.Logger, remoteChainSelector uint64, offRampAddress ccip.Address, evmBatchCaller rpclib.EvmBatchCaller) (*batchreader.EVMTokenPoolBatchedReader, error) {
	return batchreader.NewEVMTokenPoolBatchedReader(lggr, remoteChainSelector, offRampAddress, evmBatchCaller)
}

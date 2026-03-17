package ccip

import (
	"context"

	"github.com/google/uuid"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/offchainaggregator/generated/ocr2/offchainaggregator"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/batchreader"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/pricegetter"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/ccip/rpclib"
)

const OffchainAggregator = "OffchainAggregator"
const DecimalsMethodName = "decimals"
const LatestRoundDataMethodName = "latestRoundData"

type DynamicPriceGetterClient = pricegetter.DynamicPriceGetterClient

type DynamicPriceGetter = pricegetter.DynamicPriceGetter

type AllTokensPriceGetter = pricegetter.AllTokensPriceGetter

func NewPipelineGetter(
	source string,
	runner pipeline.Runner,
	jobID int32,
	externalJobID uuid.UUID,
	name string,
	lggr logger.Logger,
	sourceNativeTokenAddr ccip.Address,
	sourceChainSelector uint64,
	destChainSelector uint64,
) (*pricegetter.PipelineGetter, error) {
	return pricegetter.NewPipelineGetter(source, runner, jobID, externalJobID, name, lggr,
		sourceNativeTokenAddr, sourceChainSelector, destChainSelector)
}

func NewDynamicPriceGetter(cfg config.DynamicPriceGetterConfig, contractReaders map[uint64]types.ContractReader) (*DynamicPriceGetter, error) {
	return pricegetter.NewDynamicPriceGetter(cfg, contractReaders)
}

var DefaultRpcBatchSizeLimit = rpclib.DefaultRpcBatchSizeLimit
var DefaultRpcBatchBackOffMultiplier = rpclib.DefaultRpcBatchBackOffMultiplier
var DefaultMaxParallelRpcCalls = rpclib.DefaultMaxParallelRpcCalls

func NewEVMTokenPoolBatchedReader(lggr logger.Logger, remoteChainSelector uint64, offRampAddress ccip.Address, evmBatchCaller rpclib.EvmBatchCaller) (*batchreader.EVMTokenPoolBatchedReader, error) {
	return batchreader.NewEVMTokenPoolBatchedReader(lggr, remoteChainSelector, offRampAddress, evmBatchCaller)
}

type ChainAgnosticPriceRegistry struct {
	p ChainAgnosticPriceRegistryFactory
}

// [ChainAgnosticPriceRegistryFactory] is satisfied by [commontypes.CCIPCommitProvider] and [commontypes.CCIPExecProvider]
type ChainAgnosticPriceRegistryFactory interface {
	NewPriceRegistryReader(ctx context.Context, addr ccip.Address) (ccip.PriceRegistryReader, error)
}

func (c *ChainAgnosticPriceRegistry) NewPriceRegistryReader(ctx context.Context, addr ccip.Address) (ccip.PriceRegistryReader, error) {
	return c.p.NewPriceRegistryReader(ctx, addr)
}

func NewChainAgnosticPriceRegistry(provider ChainAgnosticPriceRegistryFactory) *ChainAgnosticPriceRegistry {
	return &ChainAgnosticPriceRegistry{provider}
}

const OffChainAggregatorABI = offchainaggregator.OffchainAggregatorABI

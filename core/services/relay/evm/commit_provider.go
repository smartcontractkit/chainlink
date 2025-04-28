package evm

import (
	"context"
	"fmt"
	"math"
	"math/big"

	"go.uber.org/multierr"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"

	"github.com/smartcontractkit/chainlink-evm/pkg/client"
	"github.com/smartcontractkit/chainlink-evm/pkg/gas"
	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/estimatorconfig"
)

var _ commontypes.CCIPCommitProvider = (*SrcCommitProvider)(nil)
var _ commontypes.CCIPCommitProvider = (*DstCommitProvider)(nil)

type SrcCommitProvider struct {
	lggr               logger.Logger
	startBlock         uint64
	client             client.Client
	lp                 logpoller.LogPoller
	estimator          gas.EvmFeeEstimator
	maxGasPrice        *big.Int
	feeEstimatorConfig estimatorconfig.FeeEstimatorConfigProvider

	// these values will be lazily initialized
	seenOnRampAddress       *cciptypes.Address
	seenSourceChainSelector *uint64
	seenDestChainSelector   *uint64
}

func NewSrcCommitProvider(
	lggr logger.Logger,
	startBlock uint64,
	client client.Client,
	lp logpoller.LogPoller,
	srcEstimator gas.EvmFeeEstimator,
	maxGasPrice *big.Int,
	feeEstimatorConfig estimatorconfig.FeeEstimatorConfigProvider,
) commontypes.CCIPCommitProvider {
	return &SrcCommitProvider{
		lggr:               logger.Named(lggr, "SrcCommitProvider"),
		startBlock:         startBlock,
		client:             client,
		lp:                 lp,
		estimator:          srcEstimator,
		maxGasPrice:        maxGasPrice,
		feeEstimatorConfig: feeEstimatorConfig,
	}
}

type DstCommitProvider struct {
	lggr                logger.Logger
	versionFinder       ccip.VersionFinder
	startBlock          uint64
	client              client.Client
	lp                  logpoller.LogPoller
	contractTransmitter ContractTransmitter
	configWatcher       *configWatcher
	gasEstimator        gas.EvmFeeEstimator
	maxGasPrice         big.Int
	feeEstimatorConfig  estimatorconfig.FeeEstimatorConfigProvider

	// these values will be lazily initialized
	seenCommitStoreAddress *cciptypes.Address
	seenOffRampAddress     *cciptypes.Address
}

func NewDstCommitProvider(
	lggr logger.Logger,
	versionFinder ccip.VersionFinder,
	startBlock uint64,
	client client.Client,
	lp logpoller.LogPoller,
	gasEstimator gas.EvmFeeEstimator,
	maxGasPrice big.Int,
	contractTransmitter ContractTransmitter,
	configWatcher *configWatcher,
	feeEstimatorConfig estimatorconfig.FeeEstimatorConfigProvider,
) commontypes.CCIPCommitProvider {
	return &DstCommitProvider{
		lggr:                logger.Named(lggr, "DstCommitProvider"),
		versionFinder:       versionFinder,
		startBlock:          startBlock,
		client:              client,
		lp:                  lp,
		contractTransmitter: contractTransmitter,
		configWatcher:       configWatcher,
		gasEstimator:        gasEstimator,
		maxGasPrice:         maxGasPrice,
		feeEstimatorConfig:  feeEstimatorConfig,
	}
}

func (p *SrcCommitProvider) Name() string {
	return p.lggr.Name()
}

// Close is called when the job that created this provider is deleted.
// At this time, any of the methods on the provider may or may not have been called.
// If NewOnRampReader has not been called, their corresponding
// Close methods will be expected to error.
func (p *SrcCommitProvider) Close() error {
	versionFinder := ccip.NewEvmVersionFinder()

	unregisterFuncs := make([]func() error, 0, 2)
	unregisterFuncs = append(unregisterFuncs, func() error {
		// avoid panic in the case NewOnRampReader wasn't called
		if p.seenOnRampAddress == nil {
			return nil
		}
		return ccip.CloseOnRampReader(context.Background(), p.lggr, versionFinder, *p.seenSourceChainSelector, *p.seenDestChainSelector, *p.seenOnRampAddress, p.lp, p.client)
	})

	var multiErr error
	for _, fn := range unregisterFuncs {
		if err := fn(); err != nil {
			multiErr = multierr.Append(multiErr, err)
		}
	}
	return multiErr
}

func (p *SrcCommitProvider) Ready() error {
	return nil
}

func (p *SrcCommitProvider) HealthReport() map[string]error {
	return map[string]error{p.Name(): nil}
}

func (p *SrcCommitProvider) OffchainConfigDigester() ocrtypes.OffchainConfigDigester {
	// TODO CCIP-2494
	// "OffchainConfigDigester called on SrcCommitProvider. Valid on DstCommitProvider."
	return UnimplementedOffchainConfigDigester{}
}

func (p *SrcCommitProvider) ContractConfigTracker() ocrtypes.ContractConfigTracker {
	// // TODO CCIP-2494
	// "ContractConfigTracker called on SrcCommitProvider. Valid on DstCommitProvider.")
	return UnimplementedContractConfigTracker{}
}

func (p *SrcCommitProvider) ContractTransmitter() ocrtypes.ContractTransmitter {
	// // TODO CCIP-2494
	// "ContractTransmitter called on SrcCommitProvider. Valid on DstCommitProvider."
	return UnimplementedContractTransmitter{}
}

func (p *SrcCommitProvider) ContractReader() commontypes.ContractReader {
	return nil
}

func (p *SrcCommitProvider) Codec() commontypes.Codec {
	return nil
}

func (p *DstCommitProvider) Name() string {
	return p.lggr.Name()
}

func (p *DstCommitProvider) Close() error {
	ctx := context.Background()
	versionFinder := ccip.NewEvmVersionFinder()

	unregisterFuncs := make([]func(ctx context.Context) error, 0, 2)
	unregisterFuncs = append(unregisterFuncs, func(ctx context.Context) error {
		if p.seenCommitStoreAddress == nil {
			return nil
		}
		return ccip.CloseCommitStoreReader(ctx, p.lggr, versionFinder, *p.seenCommitStoreAddress, p.client, p.lp, p.feeEstimatorConfig)
	})
	unregisterFuncs = append(unregisterFuncs, func(ctx context.Context) error {
		if p.seenOffRampAddress == nil {
			return nil
		}
		return ccip.CloseOffRampReader(ctx, p.lggr, versionFinder, *p.seenOffRampAddress, p.client, p.lp, nil, big.NewInt(0), p.feeEstimatorConfig)
	})

	var multiErr error
	for _, fn := range unregisterFuncs {
		if err := fn(ctx); err != nil {
			multiErr = multierr.Append(multiErr, err)
		}
	}

	return multiErr
}

func (p *DstCommitProvider) Ready() error {
	return nil
}

func (p *DstCommitProvider) HealthReport() map[string]error {
	return make(map[string]error)
}

func (p *DstCommitProvider) OffchainConfigDigester() ocrtypes.OffchainConfigDigester {
	return p.configWatcher.OffchainConfigDigester()
}

func (p *DstCommitProvider) ContractConfigTracker() ocrtypes.ContractConfigTracker {
	return p.configWatcher.ContractConfigTracker()
}

func (p *DstCommitProvider) ContractTransmitter() ocrtypes.ContractTransmitter {
	return p.contractTransmitter
}

func (p *DstCommitProvider) ContractReader() commontypes.ContractReader {
	return nil
}

func (p *DstCommitProvider) Codec() commontypes.Codec {
	return nil
}

func (p *SrcCommitProvider) Start(ctx context.Context) error {
	if p.startBlock != 0 {
		p.lggr.Infow("start replaying src chain", "fromBlock", p.startBlock)
		if p.startBlock > math.MaxInt64 {
			return fmt.Errorf("start block overflows int64: %d", p.startBlock)
		}
		return p.lp.Replay(ctx, int64(p.startBlock))
	}
	return nil
}

func (p *DstCommitProvider) Start(ctx context.Context) error {
	if p.startBlock != 0 {
		p.lggr.Infow("start replaying dst chain", "fromBlock", p.startBlock)
		if p.startBlock > math.MaxInt64 {
			return fmt.Errorf("start block overflows int64: %d", p.startBlock)
		}
		return p.lp.Replay(ctx, int64(p.startBlock))
	}
	return nil
}

func (p *SrcCommitProvider) NewPriceGetter(ctx context.Context) (priceGetter cciptypes.PriceGetter, err error) {
	return nil, fmt.Errorf("can't construct a price getter from one relayer")
}

func (p *DstCommitProvider) NewPriceGetter(ctx context.Context) (priceGetter cciptypes.PriceGetter, err error) {
	return nil, fmt.Errorf("can't construct a price getter from one relayer")
}

func (p *SrcCommitProvider) NewCommitStoreReader(ctx context.Context, commitStoreAddress cciptypes.Address) (commitStoreReader cciptypes.CommitStoreReader, err error) {
	commitStoreReader = NewIncompleteSourceCommitStoreReader(p.estimator, p.maxGasPrice, p.feeEstimatorConfig)
	return
}

func (p *DstCommitProvider) NewCommitStoreReader(ctx context.Context, commitStoreAddress cciptypes.Address) (commitStoreReader cciptypes.CommitStoreReader, err error) {
	p.seenCommitStoreAddress = &commitStoreAddress

	versionFinder := ccip.NewEvmVersionFinder()
	commitStoreReader, err = NewIncompleteDestCommitStoreReader(ctx, p.lggr, versionFinder, commitStoreAddress, p.client, p.lp, p.feeEstimatorConfig)
	return
}

func (p *SrcCommitProvider) NewOnRampReader(ctx context.Context, onRampAddress cciptypes.Address, sourceChainSelector uint64, destChainSelector uint64) (onRampReader cciptypes.OnRampReader, err error) {
	p.seenOnRampAddress = &onRampAddress
	p.seenSourceChainSelector = &sourceChainSelector
	p.seenDestChainSelector = &destChainSelector

	versionFinder := ccip.NewEvmVersionFinder()

	onRampReader, err = ccip.NewOnRampReader(ctx, p.lggr, versionFinder, sourceChainSelector, destChainSelector, onRampAddress, p.lp, p.client)
	if err != nil {
		return nil, err
	}
	p.feeEstimatorConfig.SetOnRampReader(onRampReader)
	return
}

func (p *DstCommitProvider) NewOnRampReader(ctx context.Context, onRampAddress cciptypes.Address, sourceChainSelector uint64, destChainSelector uint64) (onRampReader cciptypes.OnRampReader, err error) {
	return nil, fmt.Errorf("invalid: NewOnRampReader called for DstCommitProvider.NewOnRampReader should be called on SrcCommitProvider")
}

func (p *SrcCommitProvider) NewOffRampReader(ctx context.Context, offRampAddr cciptypes.Address) (offRampReader cciptypes.OffRampReader, err error) {
	return nil, fmt.Errorf("invalid: NewOffRampReader called for SrcCommitProvider. NewOffRampReader should be called on DstCommitProvider")
}

func (p *DstCommitProvider) NewOffRampReader(ctx context.Context, offRampAddr cciptypes.Address) (offRampReader cciptypes.OffRampReader, err error) {
	offRampReader, err = ccip.NewOffRampReader(ctx, p.lggr, p.versionFinder, offRampAddr, p.client, p.lp, p.gasEstimator, &p.maxGasPrice, true, p.feeEstimatorConfig)
	return
}

func (p *SrcCommitProvider) NewPriceRegistryReader(ctx context.Context, addr cciptypes.Address) (priceRegistryReader cciptypes.PriceRegistryReader, err error) {
	return nil, fmt.Errorf("invalid: NewPriceRegistryReader called for SrcCommitProvider. NewOffRampReader should be called on DstCommitProvider")
}

func (p *DstCommitProvider) NewPriceRegistryReader(ctx context.Context, addr cciptypes.Address) (priceRegistryReader cciptypes.PriceRegistryReader, err error) {
	destPriceRegistry := ccip.NewEvmPriceRegistry(p.lp, p.client, p.lggr, ccip.CommitPluginLabel)
	priceRegistryReader, err = destPriceRegistry.NewPriceRegistryReader(ctx, addr)
	return
}

func (p *SrcCommitProvider) SourceNativeToken(ctx context.Context, sourceRouterAddr cciptypes.Address) (cciptypes.Address, error) {
	sourceRouterAddrHex, err := ccip.GenericAddrToEvm(sourceRouterAddr)
	if err != nil {
		return "", err
	}
	sourceRouter, err := router.NewRouter(sourceRouterAddrHex, p.client)
	if err != nil {
		return "", err
	}
	sourceNative, err := sourceRouter.GetWrappedNative(&bind.CallOpts{Context: ctx})
	if err != nil {
		return "", err
	}

	return ccip.EvmAddrToGeneric(sourceNative), nil
}

func (p *DstCommitProvider) SourceNativeToken(ctx context.Context, sourceRouterAddr cciptypes.Address) (cciptypes.Address, error) {
	return "", fmt.Errorf("invalid: SourceNativeToken called for DstCommitProvider. SourceNativeToken should be called on SrcCommitProvider")
}

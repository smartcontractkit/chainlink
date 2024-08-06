package evm

import (
	"context"
	"fmt"
	"math/big"

	"go.uber.org/multierr"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip"
)

var _ commontypes.CCIPCommitProvider = (*SrcCommitProvider)(nil)
var _ commontypes.CCIPCommitProvider = (*DstCommitProvider)(nil)

type SrcCommitProvider struct {
	lggr        logger.Logger
	startBlock  uint64
	client      client.Client
	lp          logpoller.LogPoller
	estimator   gas.EvmFeeEstimator
	maxGasPrice *big.Int

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
) commontypes.CCIPCommitProvider {
	return &SrcCommitProvider{
		lggr:        lggr,
		startBlock:  startBlock,
		client:      client,
		lp:          lp,
		estimator:   srcEstimator,
		maxGasPrice: maxGasPrice,
	}
}

type DstCommitProvider struct {
	lggr                logger.Logger
	versionFinder       ccip.VersionFinder
	startBlock          uint64
	client              client.Client
	lp                  logpoller.LogPoller
	contractTransmitter *contractTransmitter
	configWatcher       *configWatcher
	gasEstimator        gas.EvmFeeEstimator
	maxGasPrice         big.Int

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
	contractTransmitter contractTransmitter,
	configWatcher *configWatcher,
) commontypes.CCIPCommitProvider {
	return &DstCommitProvider{
		lggr:                lggr,
		versionFinder:       versionFinder,
		startBlock:          startBlock,
		client:              client,
		lp:                  lp,
		contractTransmitter: &contractTransmitter,
		configWatcher:       configWatcher,
		gasEstimator:        gasEstimator,
		maxGasPrice:         maxGasPrice,
	}
}

func (P *SrcCommitProvider) Name() string {
	return "CCIPCommitProvider.SrcRelayerProvider"
}

// Close is called when the job that created this provider is deleted.
// At this time, any of the methods on the provider may or may not have been called.
// If NewOnRampReader has not been called, their corresponding
// Close methods will be expected to error.
func (P *SrcCommitProvider) Close() error {
	versionFinder := ccip.NewEvmVersionFinder()

	unregisterFuncs := make([]func() error, 0, 2)
	unregisterFuncs = append(unregisterFuncs, func() error {
		// avoid panic in the case NewOnRampReader wasn't called
		if P.seenOnRampAddress == nil {
			return nil
		}
		return ccip.CloseOnRampReader(P.lggr, versionFinder, *P.seenSourceChainSelector, *P.seenDestChainSelector, *P.seenOnRampAddress, P.lp, P.client)
	})

	var multiErr error
	for _, fn := range unregisterFuncs {
		if err := fn(); err != nil {
			multiErr = multierr.Append(multiErr, err)
		}
	}
	return multiErr
}

func (P *SrcCommitProvider) Ready() error {
	return nil
}

func (P *SrcCommitProvider) HealthReport() map[string]error {
	return make(map[string]error)
}

func (P *SrcCommitProvider) OffchainConfigDigester() ocrtypes.OffchainConfigDigester {
	// TODO CCIP-2494
	// "OffchainConfigDigester called on SrcCommitProvider. Valid on DstCommitProvider."
	return UnimplementedOffchainConfigDigester{}
}

func (P *SrcCommitProvider) ContractConfigTracker() ocrtypes.ContractConfigTracker {
	// // TODO CCIP-2494
	// "ContractConfigTracker called on SrcCommitProvider. Valid on DstCommitProvider.")
	return UnimplementedContractConfigTracker{}
}

func (P *SrcCommitProvider) ContractTransmitter() ocrtypes.ContractTransmitter {
	// // TODO CCIP-2494
	// "ContractTransmitter called on SrcCommitProvider. Valid on DstCommitProvider."
	return UnimplementedContractTransmitter{}
}

func (P *SrcCommitProvider) ChainReader() commontypes.ContractReader {
	return nil
}

func (P *SrcCommitProvider) Codec() commontypes.Codec {
	return nil
}

func (P *DstCommitProvider) Name() string {
	return "CCIPCommitProvider.DstRelayerProvider"
}

func (P *DstCommitProvider) Close() error {
	versionFinder := ccip.NewEvmVersionFinder()

	unregisterFuncs := make([]func() error, 0, 2)
	unregisterFuncs = append(unregisterFuncs, func() error {
		if P.seenCommitStoreAddress == nil {
			return nil
		}
		return ccip.CloseCommitStoreReader(P.lggr, versionFinder, *P.seenCommitStoreAddress, P.client, P.lp)
	})
	unregisterFuncs = append(unregisterFuncs, func() error {
		if P.seenOffRampAddress == nil {
			return nil
		}
		return ccip.CloseOffRampReader(P.lggr, versionFinder, *P.seenOffRampAddress, P.client, P.lp, nil, big.NewInt(0))
	})

	var multiErr error
	for _, fn := range unregisterFuncs {
		if err := fn(); err != nil {
			multiErr = multierr.Append(multiErr, err)
		}
	}
	return multiErr
}

func (P *DstCommitProvider) Ready() error {
	return nil
}

func (P *DstCommitProvider) HealthReport() map[string]error {
	return make(map[string]error)
}

func (P *DstCommitProvider) OffchainConfigDigester() ocrtypes.OffchainConfigDigester {
	return P.configWatcher.OffchainConfigDigester()
}

func (P *DstCommitProvider) ContractConfigTracker() ocrtypes.ContractConfigTracker {
	return P.configWatcher.ContractConfigTracker()
}

func (P *DstCommitProvider) ContractTransmitter() ocrtypes.ContractTransmitter {
	return P.contractTransmitter
}

func (P *DstCommitProvider) ChainReader() commontypes.ContractReader {
	return nil
}

func (P *DstCommitProvider) Codec() commontypes.Codec {
	return nil
}

func (P *SrcCommitProvider) Start(ctx context.Context) error {
	if P.startBlock != 0 {
		P.lggr.Infow("start replaying src chain", "fromBlock", P.startBlock)
		return P.lp.Replay(ctx, int64(P.startBlock))
	}
	return nil
}

func (P *DstCommitProvider) Start(ctx context.Context) error {
	if P.startBlock != 0 {
		P.lggr.Infow("start replaying dst chain", "fromBlock", P.startBlock)
		return P.lp.Replay(ctx, int64(P.startBlock))
	}
	return nil
}

func (P *SrcCommitProvider) NewPriceGetter(ctx context.Context) (priceGetter cciptypes.PriceGetter, err error) {
	return nil, fmt.Errorf("can't construct a price getter from one relayer")
}

func (P *DstCommitProvider) NewPriceGetter(ctx context.Context) (priceGetter cciptypes.PriceGetter, err error) {
	return nil, fmt.Errorf("can't construct a price getter from one relayer")
}

func (P *SrcCommitProvider) NewCommitStoreReader(ctx context.Context, commitStoreAddress cciptypes.Address) (commitStoreReader cciptypes.CommitStoreReader, err error) {
	commitStoreReader = NewIncompleteSourceCommitStoreReader(P.estimator, P.maxGasPrice)
	return
}

func (P *DstCommitProvider) NewCommitStoreReader(ctx context.Context, commitStoreAddress cciptypes.Address) (commitStoreReader cciptypes.CommitStoreReader, err error) {
	P.seenCommitStoreAddress = &commitStoreAddress

	versionFinder := ccip.NewEvmVersionFinder()
	commitStoreReader, err = NewIncompleteDestCommitStoreReader(P.lggr, versionFinder, commitStoreAddress, P.client, P.lp)
	return
}

func (P *SrcCommitProvider) NewOnRampReader(ctx context.Context, onRampAddress cciptypes.Address, sourceChainSelector uint64, destChainSelector uint64) (onRampReader cciptypes.OnRampReader, err error) {
	P.seenOnRampAddress = &onRampAddress
	P.seenSourceChainSelector = &sourceChainSelector
	P.seenDestChainSelector = &destChainSelector

	versionFinder := ccip.NewEvmVersionFinder()
	onRampReader, err = ccip.NewOnRampReader(P.lggr, versionFinder, sourceChainSelector, destChainSelector, onRampAddress, P.lp, P.client)
	return
}

func (P *DstCommitProvider) NewOnRampReader(ctx context.Context, onRampAddress cciptypes.Address, sourceChainSelector uint64, destChainSelector uint64) (onRampReader cciptypes.OnRampReader, err error) {
	return nil, fmt.Errorf("invalid: NewOnRampReader called for DstCommitProvider.NewOnRampReader should be called on SrcCommitProvider")
}

func (P *SrcCommitProvider) NewOffRampReader(ctx context.Context, offRampAddr cciptypes.Address) (offRampReader cciptypes.OffRampReader, err error) {
	return nil, fmt.Errorf("invalid: NewOffRampReader called for SrcCommitProvider. NewOffRampReader should be called on DstCommitProvider")
}

func (P *DstCommitProvider) NewOffRampReader(ctx context.Context, offRampAddr cciptypes.Address) (offRampReader cciptypes.OffRampReader, err error) {
	offRampReader, err = ccip.NewOffRampReader(P.lggr, P.versionFinder, offRampAddr, P.client, P.lp, P.gasEstimator, &P.maxGasPrice, true)
	return
}

func (P *SrcCommitProvider) NewPriceRegistryReader(ctx context.Context, addr cciptypes.Address) (priceRegistryReader cciptypes.PriceRegistryReader, err error) {
	return nil, fmt.Errorf("invalid: NewPriceRegistryReader called for SrcCommitProvider. NewOffRampReader should be called on DstCommitProvider")
}

func (P *DstCommitProvider) NewPriceRegistryReader(ctx context.Context, addr cciptypes.Address) (priceRegistryReader cciptypes.PriceRegistryReader, err error) {
	destPriceRegistry := ccip.NewEvmPriceRegistry(P.lp, P.client, P.lggr, ccip.CommitPluginLabel)
	priceRegistryReader, err = destPriceRegistry.NewPriceRegistryReader(ctx, addr)
	return
}

func (P *SrcCommitProvider) SourceNativeToken(ctx context.Context, sourceRouterAddr cciptypes.Address) (cciptypes.Address, error) {
	sourceRouterAddrHex, err := ccip.GenericAddrToEvm(sourceRouterAddr)
	if err != nil {
		return "", err
	}
	sourceRouter, err := router.NewRouter(sourceRouterAddrHex, P.client)
	if err != nil {
		return "", err
	}
	sourceNative, err := sourceRouter.GetWrappedNative(&bind.CallOpts{Context: ctx})
	if err != nil {
		return "", err
	}

	return ccip.EvmAddrToGeneric(sourceNative), nil
}

func (P *DstCommitProvider) SourceNativeToken(ctx context.Context, sourceRouterAddr cciptypes.Address) (cciptypes.Address, error) {
	return "", fmt.Errorf("invalid: SourceNativeToken called for DstCommitProvider. SourceNativeToken should be called on SrcCommitProvider")
}

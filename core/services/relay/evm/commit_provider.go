package evm

import (
	"context"
	"fmt"
	"math/big"

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
	lggr       logger.Logger
	startBlock uint64
	client     client.Client
	lp         logpoller.LogPoller
}

func NewSrcCommitProvider(
	lggr logger.Logger,
	startBlock uint64,
	client client.Client,
	lp logpoller.LogPoller,
) commontypes.CCIPCommitProvider {
	return &SrcCommitProvider{
		lggr:       lggr,
		startBlock: startBlock,
		client:     client,
		lp:         lp,
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

func (P SrcCommitProvider) Name() string {
	return "CCIPCommitProvider.SrcRelayerProvider"
}

func (P SrcCommitProvider) Close() error {
	return nil
}

func (P SrcCommitProvider) Ready() error {
	return nil
}

func (P SrcCommitProvider) HealthReport() map[string]error {
	return make(map[string]error)
}

func (P SrcCommitProvider) OffchainConfigDigester() ocrtypes.OffchainConfigDigester {
	// TODO CCIP-2494
	// "OffchainConfigDigester called on SrcCommitProvider. Valid on DstCommitProvider."
	return UnimplementedOffchainConfigDigester{}
}

func (P SrcCommitProvider) ContractConfigTracker() ocrtypes.ContractConfigTracker {
	// // TODO CCIP-2494
	// "ContractConfigTracker called on SrcCommitProvider. Valid on DstCommitProvider.")
	return UnimplementedContractConfigTracker{}
}

func (P SrcCommitProvider) ContractTransmitter() ocrtypes.ContractTransmitter {
	// // TODO CCIP-2494
	// "ContractTransmitter called on SrcCommitProvider. Valid on DstCommitProvider."
	return UnimplementedContractTransmitter{}
}

func (P SrcCommitProvider) ChainReader() commontypes.ContractReader {
	return nil
}

func (P SrcCommitProvider) Codec() commontypes.Codec {
	return nil
}

func (P DstCommitProvider) Name() string {
	return "CCIPCommitProvider.DstRelayerProvider"
}

func (P DstCommitProvider) Close() error {
	return nil
}

func (P DstCommitProvider) Ready() error {
	return nil
}

func (P DstCommitProvider) HealthReport() map[string]error {
	return make(map[string]error)
}

func (P DstCommitProvider) OffchainConfigDigester() ocrtypes.OffchainConfigDigester {
	return P.configWatcher.OffchainConfigDigester()
}

func (P DstCommitProvider) ContractConfigTracker() ocrtypes.ContractConfigTracker {
	return P.configWatcher.ContractConfigTracker()
}

func (P DstCommitProvider) ContractTransmitter() ocrtypes.ContractTransmitter {
	return P.contractTransmitter
}

func (P DstCommitProvider) ChainReader() commontypes.ContractReader {
	return nil
}

func (P DstCommitProvider) Codec() commontypes.Codec {
	return nil
}

func (P SrcCommitProvider) Start(ctx context.Context) error {
	if P.startBlock != 0 {
		P.lggr.Infow("start replaying src chain", "fromBlock", P.startBlock)
		return P.lp.Replay(ctx, int64(P.startBlock))
	}
	return nil
}

func (P DstCommitProvider) Start(ctx context.Context) error {
	if P.startBlock != 0 {
		P.lggr.Infow("start replaying dst chain", "fromBlock", P.startBlock)
		return P.lp.Replay(ctx, int64(P.startBlock))
	}
	return nil
}

func (P SrcCommitProvider) NewPriceGetter(ctx context.Context) (priceGetter cciptypes.PriceGetter, err error) {
	return nil, fmt.Errorf("can't construct a price getter from one relayer")
}

func (P DstCommitProvider) NewPriceGetter(ctx context.Context) (priceGetter cciptypes.PriceGetter, err error) {
	return nil, fmt.Errorf("can't construct a price getter from one relayer")
}

func (P SrcCommitProvider) NewCommitStoreReader(ctx context.Context, commitStoreAddress cciptypes.Address) (commitStoreReader cciptypes.CommitStoreReader, err error) {
	return nil, fmt.Errorf("can't construct a commit store reader from one relayer")
}

func (P DstCommitProvider) NewCommitStoreReader(ctx context.Context, commitStoreAddress cciptypes.Address) (commitStoreReader cciptypes.CommitStoreReader, err error) {
	return nil, fmt.Errorf("can't construct a commit store reader from one relayer")
}

func (P SrcCommitProvider) NewOnRampReader(ctx context.Context, onRampAddress cciptypes.Address, sourceChainSelector uint64, destChainSelector uint64) (onRampReader cciptypes.OnRampReader, err error) {
	versionFinder := ccip.NewEvmVersionFinder()
	onRampReader, err = ccip.NewOnRampReader(P.lggr, versionFinder, sourceChainSelector, destChainSelector, onRampAddress, P.lp, P.client)
	return
}

func (P DstCommitProvider) NewOnRampReader(ctx context.Context, onRampAddress cciptypes.Address, sourceChainSelector uint64, destChainSelector uint64) (onRampReader cciptypes.OnRampReader, err error) {
	return nil, fmt.Errorf("invalid: NewOnRampReader called for DstCommitProvider.NewOnRampReader should be called on SrcCommitProvider")
}

func (P SrcCommitProvider) NewOffRampReader(ctx context.Context, offRampAddr cciptypes.Address) (offRampReader cciptypes.OffRampReader, err error) {
	return nil, fmt.Errorf("invalid: NewOffRampReader called for SrcCommitProvider. NewOffRampReader should be called on DstCommitProvider")
}

func (P DstCommitProvider) NewOffRampReader(ctx context.Context, offRampAddr cciptypes.Address) (offRampReader cciptypes.OffRampReader, err error) {
	offRampReader, err = ccip.NewOffRampReader(P.lggr, P.versionFinder, offRampAddr, P.client, P.lp, P.gasEstimator, &P.maxGasPrice, true)
	return
}

func (P SrcCommitProvider) NewPriceRegistryReader(ctx context.Context, addr cciptypes.Address) (priceRegistryReader cciptypes.PriceRegistryReader, err error) {
	return nil, fmt.Errorf("invalid: NewPriceRegistryReader called for SrcCommitProvider. NewOffRampReader should be called on DstCommitProvider")
}

func (P DstCommitProvider) NewPriceRegistryReader(ctx context.Context, addr cciptypes.Address) (priceRegistryReader cciptypes.PriceRegistryReader, err error) {
	destPriceRegistry := ccip.NewEvmPriceRegistry(P.lp, P.client, P.lggr, ccip.CommitPluginLabel)
	priceRegistryReader, err = destPriceRegistry.NewPriceRegistryReader(ctx, addr)
	return
}

func (P SrcCommitProvider) SourceNativeToken(ctx context.Context, sourceRouterAddr cciptypes.Address) (cciptypes.Address, error) {
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

func (P DstCommitProvider) SourceNativeToken(ctx context.Context, sourceRouterAddr cciptypes.Address) (cciptypes.Address, error) {
	return "", fmt.Errorf("invalid: SourceNativeToken called for DstCommitProvider. SourceNativeToken should be called on SrcCommitProvider")
}

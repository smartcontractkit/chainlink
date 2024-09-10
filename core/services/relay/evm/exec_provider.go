package evm

import (
	"context"
	"fmt"
	"math/big"
	"net/url"
	"time"

	"go.uber.org/multierr"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/tokendata/usdc"
)

type SrcExecProvider struct {
	lggr                                   logger.Logger
	versionFinder                          ccip.VersionFinder
	client                                 client.Client
	lp                                     logpoller.LogPoller
	startBlock                             uint64
	estimator                              gas.EvmFeeEstimator
	maxGasPrice                            *big.Int
	usdcReader                             *ccip.USDCReaderImpl
	usdcAttestationAPI                     string
	usdcAttestationAPITimeoutSeconds       int
	usdcAttestationAPIIntervalMilliseconds int
	usdcSrcMsgTransmitterAddr              common.Address

	// these values are nil and are updated for Close()
	seenOnRampAddress       *cciptypes.Address
	seenSourceChainSelector *uint64
	seenDestChainSelector   *uint64
}

func NewSrcExecProvider(
	lggr logger.Logger,
	versionFinder ccip.VersionFinder,
	client client.Client,
	estimator gas.EvmFeeEstimator,
	maxGasPrice *big.Int,
	lp logpoller.LogPoller,
	startBlock uint64,
	jobID string,
	usdcAttestationAPI string,
	usdcAttestationAPITimeoutSeconds int,
	usdcAttestationAPIIntervalMilliseconds int,
	usdcSrcMsgTransmitterAddr common.Address,
) (commontypes.CCIPExecProvider, error) {
	var usdcReader *ccip.USDCReaderImpl
	var err error
	if usdcAttestationAPI != "" {
		usdcReader, err = ccip.NewUSDCReader(lggr, jobID, usdcSrcMsgTransmitterAddr, lp, true)
		if err != nil {
			return nil, fmt.Errorf("new usdc reader: %w", err)
		}
	}

	return &SrcExecProvider{
		lggr:                                   lggr,
		versionFinder:                          versionFinder,
		client:                                 client,
		estimator:                              estimator,
		maxGasPrice:                            maxGasPrice,
		lp:                                     lp,
		startBlock:                             startBlock,
		usdcReader:                             usdcReader,
		usdcAttestationAPI:                     usdcAttestationAPI,
		usdcAttestationAPITimeoutSeconds:       usdcAttestationAPITimeoutSeconds,
		usdcAttestationAPIIntervalMilliseconds: usdcAttestationAPIIntervalMilliseconds,
		usdcSrcMsgTransmitterAddr:              usdcSrcMsgTransmitterAddr,
	}, nil
}

func (s *SrcExecProvider) Name() string {
	return "CCIP.SrcExecProvider"
}

func (s *SrcExecProvider) Start(ctx context.Context) error {
	if s.startBlock != 0 {
		s.lggr.Infow("start replaying src chain", "fromBlock", s.startBlock)
		return s.lp.Replay(ctx, int64(s.startBlock))
	}
	return nil
}

// Close is called when the job that created this provider is closed.
func (s *SrcExecProvider) Close() error {
	versionFinder := ccip.NewEvmVersionFinder()

	unregisterFuncs := make([]func() error, 0, 2)
	unregisterFuncs = append(unregisterFuncs, func() error {
		// avoid panic in the case NewOnRampReader wasn't called
		if s.seenOnRampAddress == nil {
			return nil
		}
		return ccip.CloseOnRampReader(s.lggr, versionFinder, *s.seenSourceChainSelector, *s.seenDestChainSelector, *s.seenOnRampAddress, s.lp, s.client)
	})
	unregisterFuncs = append(unregisterFuncs, func() error {
		if s.usdcAttestationAPI == "" {
			return nil
		}
		return ccip.CloseUSDCReader(s.lggr, s.lggr.Name(), s.usdcSrcMsgTransmitterAddr, s.lp)
	})
	var multiErr error
	for _, fn := range unregisterFuncs {
		if err := fn(); err != nil {
			multiErr = multierr.Append(multiErr, err)
		}
	}
	return multiErr
}

func (s *SrcExecProvider) Ready() error {
	return nil
}

func (s *SrcExecProvider) HealthReport() map[string]error {
	return make(map[string]error)
}

func (s *SrcExecProvider) OffchainConfigDigester() ocrtypes.OffchainConfigDigester {
	// TODO CCIP-2494
	// OffchainConfigDigester called on SrcExecProvider. It should only be called on DstExecProvider
	return UnimplementedOffchainConfigDigester{}
}

func (s *SrcExecProvider) ContractConfigTracker() ocrtypes.ContractConfigTracker {
	// TODO CCIP-2494
	// "ContractConfigTracker called on SrcExecProvider. It should only be called on DstExecProvider
	return UnimplementedContractConfigTracker{}
}

func (s *SrcExecProvider) ContractTransmitter() ocrtypes.ContractTransmitter {
	// TODO CCIP-2494
	// "ContractTransmitter called on SrcExecProvider. It should only be called on DstExecProvider
	return UnimplementedContractTransmitter{}
}

func (s *SrcExecProvider) ContractReader() commontypes.ContractReader {
	return nil
}

func (s *SrcExecProvider) Codec() commontypes.Codec {
	return nil
}

func (s *SrcExecProvider) GetTransactionStatus(ctx context.Context, transactionID string) (types.TransactionStatus, error) {
	return 0, fmt.Errorf("invalid: GetTransactionStatus called on SrcExecProvider. It should only be called on DstExecProvider")
}

func (s *SrcExecProvider) NewCommitStoreReader(ctx context.Context, addr cciptypes.Address) (commitStoreReader cciptypes.CommitStoreReader, err error) {
	commitStoreReader = NewIncompleteSourceCommitStoreReader(s.estimator, s.maxGasPrice)
	return
}

func (s *SrcExecProvider) NewOffRampReader(ctx context.Context, addr cciptypes.Address) (cciptypes.OffRampReader, error) {
	return nil, fmt.Errorf("invalid: NewOffRampReader called on SrcExecProvider. Valid on DstExecProvider")
}

func (s *SrcExecProvider) NewOnRampReader(ctx context.Context, onRampAddress cciptypes.Address, sourceChainSelector uint64, destChainSelector uint64) (onRampReader cciptypes.OnRampReader, err error) {
	s.seenOnRampAddress = &onRampAddress

	versionFinder := ccip.NewEvmVersionFinder()
	onRampReader, err = ccip.NewOnRampReader(s.lggr, versionFinder, sourceChainSelector, destChainSelector, onRampAddress, s.lp, s.client)
	return
}

func (s *SrcExecProvider) NewPriceRegistryReader(ctx context.Context, addr cciptypes.Address) (priceRegistryReader cciptypes.PriceRegistryReader, err error) {
	srcPriceRegistry := ccip.NewEvmPriceRegistry(s.lp, s.client, s.lggr, ccip.ExecPluginLabel)
	priceRegistryReader, err = srcPriceRegistry.NewPriceRegistryReader(ctx, addr)
	return
}

func (s *SrcExecProvider) NewTokenDataReader(ctx context.Context, tokenAddress cciptypes.Address) (tokenDataReader cciptypes.TokenDataReader, err error) {
	attestationURI, err2 := url.ParseRequestURI(s.usdcAttestationAPI)
	if err2 != nil {
		return nil, fmt.Errorf("failed to parse USDC attestation API: %w", err2)
	}
	tokenAddr, err2 := ccip.GenericAddrToEvm(tokenAddress)
	if err2 != nil {
		return nil, fmt.Errorf("failed to parse token address: %w", err2)
	}
	tokenDataReader = usdc.NewUSDCTokenDataReader(
		s.lggr,
		s.usdcReader,
		attestationURI,
		s.usdcAttestationAPITimeoutSeconds,
		tokenAddr,
		time.Duration(s.usdcAttestationAPIIntervalMilliseconds)*time.Millisecond,
	)
	return
}

func (s *SrcExecProvider) NewTokenPoolBatchedReader(ctx context.Context, offRampAddr cciptypes.Address, sourceChainSelector uint64) (cciptypes.TokenPoolBatchedReader, error) {
	return nil, fmt.Errorf("invalid: NewTokenPoolBatchedReader called on SrcExecProvider. It should only be called on DstExecProvdier")
}

func (s *SrcExecProvider) SourceNativeToken(ctx context.Context, sourceRouterAddr cciptypes.Address) (cciptypes.Address, error) {
	sourceRouterAddrHex, err := ccip.GenericAddrToEvm(sourceRouterAddr)
	if err != nil {
		return "", err
	}
	sourceRouter, err := router.NewRouter(sourceRouterAddrHex, s.client)
	if err != nil {
		return "", err
	}
	sourceNative, err := sourceRouter.GetWrappedNative(&bind.CallOpts{Context: ctx})
	if err != nil {
		return "", err
	}

	return ccip.EvmAddrToGeneric(sourceNative), nil
}

type DstExecProvider struct {
	lggr                logger.Logger
	versionFinder       ccip.VersionFinder
	client              client.Client
	lp                  logpoller.LogPoller
	startBlock          uint64
	contractTransmitter *contractTransmitter
	configWatcher       *configWatcher
	gasEstimator        gas.EvmFeeEstimator
	maxGasPrice         big.Int
	txm                 txmgr.TxManager
	offRampAddress      cciptypes.Address

	// these values are nil and are updated for Close()
	seenCommitStoreAddr *cciptypes.Address
}

func NewDstExecProvider(
	lggr logger.Logger,
	versionFinder ccip.VersionFinder,
	client client.Client,
	lp logpoller.LogPoller,
	startBlock uint64,
	contractTransmitter *contractTransmitter,
	configWatcher *configWatcher,
	gasEstimator gas.EvmFeeEstimator,
	maxGasPrice big.Int,
	txm txmgr.TxManager,
	offRampAddress cciptypes.Address,
) (commontypes.CCIPExecProvider, error) {
	return &DstExecProvider{
		lggr:                lggr,
		versionFinder:       versionFinder,
		client:              client,
		lp:                  lp,
		startBlock:          startBlock,
		contractTransmitter: contractTransmitter,
		configWatcher:       configWatcher,
		gasEstimator:        gasEstimator,
		maxGasPrice:         maxGasPrice,
		txm:                 txm,
		offRampAddress:      offRampAddress,
	}, nil
}

func (d *DstExecProvider) Name() string {
	return "CCIP.DestRelayerExecProvider"
}

func (d *DstExecProvider) Start(ctx context.Context) error {
	if d.startBlock != 0 {
		d.lggr.Infow("start replaying dst chain", "fromBlock", d.startBlock)
		return d.lp.Replay(ctx, int64(d.startBlock))
	}
	return nil
}

// Close is called when the job that created this provider is deleted
// At this time, any of the methods on the provider may or may not have been called.
// If NewOnRampReader and NewCommitStoreReader have not been called, their corresponding
// Close methods will be expected to error.
func (d *DstExecProvider) Close() error {
	versionFinder := ccip.NewEvmVersionFinder()

	unregisterFuncs := make([]func() error, 0, 2)
	unregisterFuncs = append(unregisterFuncs, func() error {
		if d.seenCommitStoreAddr == nil {
			return nil
		}
		return ccip.CloseCommitStoreReader(d.lggr, versionFinder, *d.seenCommitStoreAddr, d.client, d.lp)
	})
	unregisterFuncs = append(unregisterFuncs, func() error {
		return ccip.CloseOffRampReader(d.lggr, versionFinder, d.offRampAddress, d.client, d.lp, nil, big.NewInt(0))
	})

	var multiErr error
	for _, fn := range unregisterFuncs {
		if err := fn(); err != nil {
			multiErr = multierr.Append(multiErr, err)
		}
	}
	return multiErr
}

func (d *DstExecProvider) Ready() error {
	return nil
}

func (d *DstExecProvider) HealthReport() map[string]error {
	return make(map[string]error)
}

func (d *DstExecProvider) OffchainConfigDigester() ocrtypes.OffchainConfigDigester {
	return d.configWatcher.OffchainConfigDigester()
}

func (d *DstExecProvider) ContractConfigTracker() ocrtypes.ContractConfigTracker {
	return d.configWatcher.ContractConfigTracker()
}

func (d *DstExecProvider) ContractTransmitter() ocrtypes.ContractTransmitter {
	return d.contractTransmitter
}

func (d *DstExecProvider) ContractReader() commontypes.ContractReader {
	return nil
}

func (d *DstExecProvider) Codec() commontypes.Codec {
	return nil
}

func (d *DstExecProvider) GetTransactionStatus(ctx context.Context, transactionID string) (types.TransactionStatus, error) {
	return d.txm.GetTransactionStatus(ctx, transactionID)
}

func (d *DstExecProvider) NewCommitStoreReader(ctx context.Context, addr cciptypes.Address) (commitStoreReader cciptypes.CommitStoreReader, err error) {
	d.seenCommitStoreAddr = &addr

	versionFinder := ccip.NewEvmVersionFinder()
	commitStoreReader, err = NewIncompleteDestCommitStoreReader(d.lggr, versionFinder, addr, d.client, d.lp)
	return
}

func (d *DstExecProvider) NewOffRampReader(ctx context.Context, offRampAddress cciptypes.Address) (offRampReader cciptypes.OffRampReader, err error) {
	offRampReader, err = ccip.NewOffRampReader(d.lggr, d.versionFinder, offRampAddress, d.client, d.lp, d.gasEstimator, &d.maxGasPrice, true)
	return
}

func (d *DstExecProvider) NewOnRampReader(ctx context.Context, addr cciptypes.Address, sourceChainSelector uint64, destChainSelector uint64) (cciptypes.OnRampReader, error) {
	return nil, fmt.Errorf("invalid: NewOnRampReader called on DstExecProvider. It should only be called on SrcExecProvider")
}

func (d *DstExecProvider) NewPriceRegistryReader(ctx context.Context, addr cciptypes.Address) (priceRegistryReader cciptypes.PriceRegistryReader, err error) {
	destPriceRegistry := ccip.NewEvmPriceRegistry(d.lp, d.client, d.lggr, ccip.ExecPluginLabel)
	priceRegistryReader, err = destPriceRegistry.NewPriceRegistryReader(ctx, addr)
	return
}

func (d *DstExecProvider) NewTokenDataReader(ctx context.Context, tokenAddress cciptypes.Address) (cciptypes.TokenDataReader, error) {
	return nil, fmt.Errorf("invalid: NewTokenDataReader called on DstExecProvider. It should only be called on SrcExecProvider")
}

func (d *DstExecProvider) NewTokenPoolBatchedReader(ctx context.Context, offRampAddress cciptypes.Address, sourceChainSelector uint64) (tokenPoolBatchedReader cciptypes.TokenPoolBatchedReader, err error) {
	batchCaller := ccip.NewDynamicLimitedBatchCaller(
		d.lggr,
		d.client,
		uint(ccip.DefaultRpcBatchSizeLimit),
		uint(ccip.DefaultRpcBatchBackOffMultiplier),
		uint(ccip.DefaultMaxParallelRpcCalls),
	)

	tokenPoolBatchedReader, err = ccip.NewEVMTokenPoolBatchedReader(d.lggr, sourceChainSelector, offRampAddress, batchCaller)
	if err != nil {
		return nil, fmt.Errorf("new token pool batched reader: %w", err)
	}
	return
}

func (d *DstExecProvider) SourceNativeToken(ctx context.Context, addr cciptypes.Address) (cciptypes.Address, error) {
	return "", fmt.Errorf("invalid: SourceNativeToken called on DstExecProvider. It should only be called on SrcExecProvider")
}

package evm

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	pkgerrors "github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/libocr/gethwrappers2/ocr2aggregator"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median/evmreportcodec"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	ocr3capability "github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/triggers"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
	coretypes "github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink-evm/pkg/config/chaintype"
	"github.com/smartcontractkit/chainlink-evm/pkg/keys"
	evmtypes "github.com/smartcontractkit/chainlink-evm/pkg/types"
	txmgrcommon "github.com/smartcontractkit/chainlink-framework/chains/txmgr"
	coreconfig "github.com/smartcontractkit/chainlink/v2/core/config"

	txm "github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
	"github.com/smartcontractkit/chainlink/v2/core/services/llo"
	"github.com/smartcontractkit/chainlink/v2/core/services/llo/channeldefinitions"
	"github.com/smartcontractkit/chainlink/v2/core/services/llo/retirement"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/ccipcommit"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/ccipexec"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/estimatorconfig"
	cciptransmitter "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/transmitter"
	mercuryconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/mercury/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/codec"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/functions"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/interceptors/mantle"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

var (
	OCR2AggregatorTransmissionContractABI abi.ABI
	OCR2AggregatorLogDecoder              LogDecoder
	OCR3CapabilityLogDecoder              LogDecoder
)

func init() {
	var err error
	OCR2AggregatorTransmissionContractABI, err = abi.JSON(strings.NewReader(ocr2aggregator.OCR2AggregatorMetaData.ABI))
	if err != nil {
		panic(err)
	}
	OCR2AggregatorLogDecoder, err = newOCR2AggregatorLogDecoder()
	if err != nil {
		panic(err)
	}
	OCR3CapabilityLogDecoder, err = newOCR3CapabilityLogDecoder()
	if err != nil {
		panic(err)
	}
}

var _ commontypes.Relayer = &Relayer{}

// The current PluginProvider interface does not support an error return. This was fine up until CCIP.
// CCIP is the first product to introduce the idea of incomplete implementations of a provider based on
// what chain (for CCIP, src or dest) the provider is created for. The Unimplemented* implementations below allow us to return
// a non nil value, which is hopefully a better developer experience should you find yourself using the right methods
// but on the *wrong* provider.

// [UnimplementedOffchainConfigDigester] satisfies the OCR OffchainConfigDigester interface
type UnimplementedOffchainConfigDigester struct{}

func (e UnimplementedOffchainConfigDigester) ConfigDigest(ctx context.Context, config ocrtypes.ContractConfig) (ocrtypes.ConfigDigest, error) {
	return ocrtypes.ConfigDigest{}, fmt.Errorf("unimplemented for this relayer")
}

func (e UnimplementedOffchainConfigDigester) ConfigDigestPrefix(ctx context.Context) (ocrtypes.ConfigDigestPrefix, error) {
	return 0, fmt.Errorf("unimplemented for this relayer")
}

// [UnimplementedContractConfigTracker] satisfies the OCR ContractConfigTracker interface
type UnimplementedContractConfigTracker struct{}

func (u UnimplementedContractConfigTracker) Notify() <-chan struct{} {
	return nil
}

func (u UnimplementedContractConfigTracker) LatestConfigDetails(ctx context.Context) (changedInBlock uint64, configDigest ocrtypes.ConfigDigest, err error) {
	return 0, ocrtypes.ConfigDigest{}, fmt.Errorf("unimplemented for this relayer")
}

func (u UnimplementedContractConfigTracker) LatestConfig(ctx context.Context, changedInBlock uint64) (ocrtypes.ContractConfig, error) {
	return ocrtypes.ContractConfig{}, fmt.Errorf("unimplemented for this relayer")
}

func (u UnimplementedContractConfigTracker) LatestBlockHeight(ctx context.Context) (blockHeight uint64, err error) {
	return 0, fmt.Errorf("unimplemented for this relayer")
}

// [UnimplementedContractTransmitter] satisfies the OCR ContractTransmitter interface
type UnimplementedContractTransmitter struct{}

func (u UnimplementedContractTransmitter) Transmit(context.Context, ocrtypes.ReportContext, ocrtypes.Report, []ocrtypes.AttributedOnchainSignature) error {
	return fmt.Errorf("unimplemented for this relayer")
}

func (u UnimplementedContractTransmitter) FromAccount(ctx context.Context) (ocrtypes.Account, error) {
	return "", fmt.Errorf("unimplemented for this relayer")
}

func (u UnimplementedContractTransmitter) LatestConfigDigestAndEpoch(ctx context.Context) (configDigest ocrtypes.ConfigDigest, epoch uint32, err error) {
	return ocrtypes.ConfigDigest{}, 0, fmt.Errorf("unimplemented for this relayer")
}

type Relayer struct {
	ds                   sqlutil.DataSource
	chain                legacyevm.Chain
	lggr                 logger.SugaredLogger
	csaKeystore          coretypes.Keystore
	evmKeystore          keys.Store
	codec                commontypes.Codec
	capabilitiesRegistry coretypes.CapabilitiesRegistry

	// Mercury
	mercuryCfg        MercuryConfig
	mercuryORM        mercury.ORM
	mercuryPool       wsrpc.Pool
	triggerCapability *triggers.MercuryTriggerService

	// LLO/data streams
	cdcFactory            func() (channeldefinitions.ChannelDefinitionCacheFactory, error)
	retirementReportCache retirement.RetirementReportCache
	registerer            prometheus.Registerer
}

type MercuryConfig interface {
	Transmitter() coreconfig.MercuryTransmitter
	VerboseLogging() bool
}

type RelayerOpts struct {
	DS                    sqlutil.DataSource
	Registerer            prometheus.Registerer
	CSAKeystore           coretypes.Keystore
	EVMKeystore           keys.ChainStore
	MercuryPool           wsrpc.Pool
	RetirementReportCache retirement.RetirementReportCache
	MercuryConfig
	CapabilitiesRegistry coretypes.CapabilitiesRegistry
	HTTPClient           *http.Client
}

func (c RelayerOpts) Validate() error {
	var err error
	if c.DS == nil {
		err = errors.Join(err, errors.New("nil DataSource"))
	}
	if c.CSAKeystore == nil {
		err = errors.Join(err, errors.New("nil CSAKeystore"))
	}
	if c.EVMKeystore == nil {
		err = errors.Join(err, errors.New("nil EVMKeystore"))
	}
	if c.CapabilitiesRegistry == nil {
		err = errors.Join(err, errors.New("nil CapabilitiesRegistry"))
	}
	if err != nil {
		err = fmt.Errorf("invalid RelayerOpts: %w", err)
	}
	return err
}

func NewRelayer(lggr logger.Logger, chain legacyevm.Chain, opts RelayerOpts) (*Relayer, error) {
	err := opts.Validate()
	if err != nil {
		return nil, fmt.Errorf("cannot create evm relayer: %w", err)
	}
	sugared := logger.Sugared(lggr).Named("Relayer").With("evmChainID", chain.ID())
	mercuryORM := mercury.NewORM(opts.DS)
	cdcFactory := sync.OnceValues(func() (channeldefinitions.ChannelDefinitionCacheFactory, error) {
		chainSelector, err := chainselectors.SelectorFromChainId(chain.ID().Uint64())
		if err != nil {
			return nil, fmt.Errorf("failed to get chain selector for chain id %s: %w", chain.ID(), err)
		}
		lloORM := llo.NewChainScopedORM(opts.DS, chainSelector)
		return channeldefinitions.NewChannelDefinitionCacheFactory(sugared, lloORM, chain.LogPoller(), opts.HTTPClient), nil
	})
	return &Relayer{
		ds:                    opts.DS,
		chain:                 chain,
		lggr:                  logger.Sugared(sugared),
		registerer:            opts.Registerer,
		csaKeystore:           opts.CSAKeystore,
		evmKeystore:           opts.EVMKeystore,
		mercuryPool:           opts.MercuryPool,
		cdcFactory:            cdcFactory,
		retirementReportCache: opts.RetirementReportCache,
		mercuryORM:            mercuryORM,
		mercuryCfg:            opts.MercuryConfig,
		capabilitiesRegistry:  opts.CapabilitiesRegistry,
	}, nil
}

func (r *Relayer) Name() string {
	return r.lggr.Name()
}

func (r *Relayer) Start(ctx context.Context) error {
	wCfg := r.chain.Config().EVM().Workflow()
	// Initialize write target capability if configuration is defined
	if wCfg.ForwarderAddress() != nil && wCfg.FromAddress() != nil {
		if wCfg.GasLimitDefault() == nil {
			return fmt.Errorf("unable to instantiate write target as default gas limit is not set")
		}
		capability, err := NewWriteTarget(ctx, r, r.chain, *wCfg.GasLimitDefault(), r.lggr)
		if err != nil {
			return fmt.Errorf("failed to initialize write target: %w", err)
		}
		if err = r.capabilitiesRegistry.Add(ctx, capability); err != nil {
			return fmt.Errorf("failed to register capability: %w", err)
		}
		r.lggr.Infow("Registered write target", "chain_id", r.chain.ID())
	}
	return r.chain.Start(ctx)
}

func (r *Relayer) Close() error {
	cs := make([]io.Closer, 0, 2)
	if r.triggerCapability != nil {
		cs = append(cs, r.triggerCapability)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		err := r.capabilitiesRegistry.Remove(ctx, r.triggerCapability.ID)
		if err != nil {
			return err
		}
	}
	cs = append(cs, r.chain)
	return services.MultiCloser(cs).Close()
}

// Ready does noop: always ready
func (r *Relayer) Ready() error {
	return r.chain.Ready()
}

func (r *Relayer) HealthReport() (report map[string]error) {
	report = map[string]error{r.Name(): r.Ready()}
	maps.Copy(report, r.chain.HealthReport())
	return
}

func (r *Relayer) LatestHead(ctx context.Context) (commontypes.Head, error) {
	return r.chain.LatestHead(ctx)
}

func (r *Relayer) GetChainStatus(ctx context.Context) (commontypes.ChainStatus, error) {
	return r.chain.GetChainStatus(ctx)
}

func (r *Relayer) ListNodeStatuses(ctx context.Context, pageSize int32, pageToken string) (stats []commontypes.NodeStatus, nextPageToken string, total int, err error) {
	return r.chain.ListNodeStatuses(ctx, pageSize, pageToken)
}

func (r *Relayer) Transact(ctx context.Context, from, to string, amount *big.Int, balanceCheck bool) error {
	return r.chain.Transact(ctx, from, to, amount, balanceCheck)
}

func (r *Relayer) Replay(ctx context.Context, fromBlock string, args map[string]any) error {
	return r.chain.Replay(ctx, fromBlock, args)
}

func (r *Relayer) ID() string {
	return r.chain.ID().String()
}

func (r *Relayer) Chain() legacyevm.Chain {
	return r.chain
}

func newOCR3CapabilityConfigProvider(ctx context.Context, lggr logger.Logger, chain legacyevm.Chain, opts *types.RelayOpts) (*configWatcher, error) {
	if !common.IsHexAddress(opts.ContractID) {
		return nil, errors.New("invalid contractID, expected hex address")
	}

	aggregatorAddress := common.HexToAddress(opts.ContractID)
	offchainConfigDigester := OCR3CapabilityOffchainConfigDigester{
		ChainID:         chain.Config().EVM().ChainID().Uint64(),
		ContractAddress: aggregatorAddress,
	}
	return newContractConfigProvider(ctx, lggr, chain, opts, aggregatorAddress, OCR3CapabilityLogDecoder, offchainConfigDigester)
}

// NewPluginProvider, but customized to use a different config provider
func (r *Relayer) NewOCR3CapabilityProvider(ctx context.Context, rargs commontypes.RelayArgs, pargs commontypes.PluginArgs) (commontypes.OCR3CapabilityProvider, error) {
	lggr := logger.Sugared(r.lggr).Named("PluginProvider").Named(rargs.ExternalJobID.String())
	relayOpts := types.NewRelayOpts(rargs)
	relayConfig, err := relayOpts.RelayConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get relay config: %w", err)
	}

	configWatcher, err := newOCR3CapabilityConfigProvider(ctx, r.lggr, r.chain, relayOpts)
	if err != nil {
		return nil, err
	}

	transmitter, err := newOnChainContractTransmitter(ctx, r.lggr, rargs, r.evmKeystore, configWatcher, configTransmitterOpts{}, OCR2AggregatorTransmissionContractABI)
	if err != nil {
		return nil, err
	}

	var chainReaderService ChainReaderService
	if relayConfig.ChainReader != nil {
		if chainReaderService, err = NewChainReaderService(ctx, lggr, r.chain.LogPoller(), r.chain.HeadTracker(), r.chain.Client(), *relayConfig.ChainReader); err != nil {
			return nil, err
		}
	} else {
		lggr.Info("ChainReader missing from RelayConfig")
	}

	pp := NewPluginProvider(
		chainReaderService,
		r.codec,
		transmitter,
		configWatcher,
		lggr,
	)

	fromAccount, err := pp.ContractTransmitter().FromAccount(ctx)
	if err != nil {
		return nil, err
	}

	return &ocr3CapabilityProvider{
		PluginProvider: pp,
		transmitter:    ocr3capability.NewContractTransmitter(r.lggr, r.capabilitiesRegistry, string(fromAccount)),
	}, nil
}

func (r *Relayer) NewPluginProvider(ctx context.Context, rargs commontypes.RelayArgs, pargs commontypes.PluginArgs) (commontypes.PluginProvider, error) {
	lggr := logger.Sugared(r.lggr).Named("PluginProvider").Named(rargs.ExternalJobID.String())
	relayOpts := types.NewRelayOpts(rargs)
	relayConfig, err := relayOpts.RelayConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get relay config: %w", err)
	}

	configWatcher, err := newStandardConfigProvider(ctx, r.lggr, r.chain, relayOpts)
	if err != nil {
		return nil, err
	}
	transmitter, err := newOnChainContractTransmitter(ctx, r.lggr, rargs, r.evmKeystore, configWatcher, configTransmitterOpts{}, OCR2AggregatorTransmissionContractABI)
	if err != nil {
		return nil, err
	}

	var chainReaderService ChainReaderService
	if relayConfig.ChainReader != nil {
		if chainReaderService, err = NewChainReaderService(ctx, lggr, r.chain.LogPoller(), r.chain.HeadTracker(), r.chain.Client(), *relayConfig.ChainReader); err != nil {
			return nil, err
		}
	} else {
		lggr.Info("ChainReader missing from RelayConfig")
	}

	return NewPluginProvider(
		chainReaderService,
		r.codec,
		transmitter,
		configWatcher,
		lggr,
	), nil
}

func (r *Relayer) NewMercuryProvider(ctx context.Context, rargs commontypes.RelayArgs, pargs commontypes.PluginArgs) (commontypes.MercuryProvider, error) {
	lggr := logger.Sugared(r.lggr).Named("MercuryProvider").Named(rargs.ExternalJobID.String())
	relayOpts := types.NewRelayOpts(rargs)
	relayConfig, err := relayOpts.RelayConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get relay config: %w", err)
	}

	var mercuryConfig mercuryconfig.PluginConfig
	if err = json.Unmarshal(pargs.PluginConfig, &mercuryConfig); err != nil {
		return nil, pkgerrors.WithStack(err)
	}

	if relayConfig.FeedID == nil {
		return nil, pkgerrors.New("FeedID must be specified")
	}

	if relayConfig.ChainID.String() != r.chain.ID().String() {
		return nil, fmt.Errorf("internal error: chain id in spec does not match this relayer's chain: have %s expected %s", relayConfig.ChainID.String(), r.chain.ID().String())
	}
	cp, err := newMercuryConfigProvider(ctx, lggr, r.chain, relayOpts)
	if err != nil {
		return nil, pkgerrors.WithStack(err)
	}

	if !relayConfig.EffectiveTransmitterID.Valid {
		return nil, pkgerrors.New("EffectiveTransmitterID must be specified")
	}

	// initialize trigger capability service lazily
	if relayConfig.EnableTriggerCapability && r.triggerCapability == nil {
		if r.capabilitiesRegistry == nil {
			lggr.Errorw("trigger capability is enabled but capabilities registry is not set")
		} else {
			var err2 error
			r.triggerCapability, err2 = triggers.NewMercuryTriggerService(0, relayConfig.TriggerCapabilityName, relayConfig.TriggerCapabilityVersion, lggr)
			if err2 != nil {
				return nil, fmt.Errorf("failed to start required trigger service: %w", err2)
			}
			if err2 = r.triggerCapability.Start(ctx); err2 != nil {
				return nil, err2
			}
			if err2 = r.capabilitiesRegistry.Add(ctx, r.triggerCapability); err2 != nil {
				return nil, err2
			}
			lggr.Infow("successfully added trigger service to the Registry")
		}
	}

	return NewMercuryProvider(ctx, rargs.JobID, relayConfig, mercuryConfig, r.mercuryCfg.Transmitter(), cp, r.codec, NewMercuryChainReader(r.chain.HeadTracker()), lggr, r.csaKeystore, r.mercuryPool, r.mercuryORM, r.triggerCapability)
}

func chainToUUID(chainID *big.Int) uuid.UUID {
	// See https://www.rfc-editor.org/rfc/rfc4122.html#section-4.1.3 for the list of supported versions.
	const VersionSHA1 = 5
	var buf bytes.Buffer
	buf.WriteString("CCIP:")
	buf.Write(chainID.Bytes())
	// We use SHA-256 instead of SHA-1 because the former has better collision resistance.
	// The UUID will contain only the first 16 bytes of the hash.
	// You can't say which algorithms was used just by looking at the UUID bytes.
	return uuid.NewHash(sha256.New(), uuid.NameSpaceOID, buf.Bytes(), VersionSHA1)
}

// NewCCIPCommitProvider constructs a provider of type CCIPCommitProvider. Since this is happening in the Relayer,
// which lives in a separate process from delegate which is requesting a provider, we need to wire in through pargs
// which *type* (impl) of CCIPCommitProvider should be created. CCIP is currently a special case where the provider has a
// subset of implementations of the complete interface as certain contracts in a CCIP lane are only deployed on the src
// chain or on the dst chain. This results in the two implementations of providers: a src and dst implementation.
func (r *Relayer) NewCCIPCommitProvider(ctx context.Context, rargs commontypes.RelayArgs, pargs commontypes.PluginArgs) (commontypes.CCIPCommitProvider, error) {
	versionFinder := ccip.NewEvmVersionFinder()

	var commitPluginConfig ccipconfig.CommitPluginConfig
	err := json.Unmarshal(pargs.PluginConfig, &commitPluginConfig)
	if err != nil {
		return nil, err
	}
	sourceStartBlock := commitPluginConfig.SourceStartBlock
	destStartBlock := commitPluginConfig.DestStartBlock

	feeEstimatorConfig := estimatorconfig.NewFeeEstimatorConfigService()

	// CCIPCommit reads only when source chain is Mantle, then reports to dest chain
	// to minimize misconfigure risk, might make sense to wire Mantle only when Commit + Mantle + IsSourceProvider
	if r.chain.Config().EVM().ChainID().Uint64() == 5003 || r.chain.Config().EVM().ChainID().Uint64() == 5000 {
		if commitPluginConfig.IsSourceProvider {
			oracleAddress := r.chain.Config().EVM().GasEstimator().DAOracle().OracleAddress()
			mantleInterceptor, iErr := mantle.NewInterceptor(ctx, r.chain.Client(), oracleAddress)
			if iErr != nil {
				return nil, iErr
			}
			feeEstimatorConfig.AddGasPriceInterceptor(mantleInterceptor)
		}
	}

	// The src chain implementation of this provider does not need a configWatcher or contractTransmitter;
	// bail early.
	if commitPluginConfig.IsSourceProvider {
		return NewSrcCommitProvider(
			r.lggr,
			sourceStartBlock,
			r.chain.Client(),
			r.chain.LogPoller(),
			r.chain.GasEstimator(),
			r.chain.Config().EVM().GasEstimator().PriceMax().ToInt(),
			feeEstimatorConfig,
		), nil
	}

	relayOpts := types.NewRelayOpts(rargs)
	configWatcher, err := newStandardConfigProvider(ctx, r.lggr, r.chain, relayOpts)
	if err != nil {
		return nil, err
	}
	address := common.HexToAddress(relayOpts.ContractID)
	typ, ver, err := ccipconfig.TypeAndVersion(address, r.chain.Client())
	if err != nil {
		return nil, err
	}
	fn, err := ccipcommit.CommitReportToEthTxMeta(typ, ver)
	if err != nil {
		return nil, err
	}
	subjectID := chainToUUID(configWatcher.chain.ID())

	contractTransmitter, err := newOnChainContractTransmitter(ctx, r.lggr, rargs, r.evmKeystore, configWatcher, configTransmitterOpts{
		subjectID: &subjectID,
	}, OCR2AggregatorTransmissionContractABI, WithReportToEthMetadata(fn), WithRetention(0))
	if err != nil {
		return nil, err
	}

	return NewDstCommitProvider(
		r.lggr,
		versionFinder,
		destStartBlock,
		r.chain.Client(),
		r.chain.LogPoller(),
		r.chain.GasEstimator(),
		*r.chain.Config().EVM().GasEstimator().PriceMax().ToInt(),
		contractTransmitter,
		configWatcher,
		feeEstimatorConfig,
	), nil
}

// NewCCIPExecProvider constructs a provider of type CCIPExecProvider. Since this is happening in the Relayer,
// which lives in a separate process from delegate which is requesting a provider, we need to wire in through pargs
// which *type* (impl) of CCIPExecProvider should be created. CCIP is currently a special case where the provider has a
// subset of implementations of the complete interface as certain contracts in a CCIP lane are only deployed on the src
// chain or on the dst chain. This results in the two implementations of providers: a src and dst implementation.
func (r *Relayer) NewCCIPExecProvider(ctx context.Context, rargs commontypes.RelayArgs, pargs commontypes.PluginArgs) (commontypes.CCIPExecProvider, error) {
	versionFinder := ccip.NewEvmVersionFinder()

	var execPluginConfig ccipconfig.ExecPluginConfig
	err := json.Unmarshal(pargs.PluginConfig, &execPluginConfig)
	if err != nil {
		return nil, err
	}

	feeEstimatorConfig := estimatorconfig.NewFeeEstimatorConfigService()

	// CCIPExec reads when dest chain is mantle, and uses it to calc boosting in batching
	// to minimize misconfigure risk, make sense to wire Mantle only when Exec + Mantle + !IsSourceProvider
	if r.chain.Config().EVM().ChainID().Uint64() == 5003 || r.chain.Config().EVM().ChainID().Uint64() == 5000 {
		if !execPluginConfig.IsSourceProvider {
			oracleAddress := r.chain.Config().EVM().GasEstimator().DAOracle().OracleAddress()
			mantleInterceptor, iErr := mantle.NewInterceptor(ctx, r.chain.Client(), oracleAddress)
			if iErr != nil {
				return nil, iErr
			}
			feeEstimatorConfig.AddGasPriceInterceptor(mantleInterceptor)
		}
	}

	// The src chain implementation of this provider does not need a configWatcher or contractTransmitter;
	// bail early.
	if execPluginConfig.IsSourceProvider {
		return NewSrcExecProvider(
			ctx,
			r.lggr,
			versionFinder,
			r.chain.Client(),
			r.chain.GasEstimator(),
			r.chain.Config().EVM().GasEstimator().PriceMax().ToInt(),
			r.chain.LogPoller(),
			execPluginConfig.SourceStartBlock,
			execPluginConfig.JobID,
			execPluginConfig.USDCConfig,
			execPluginConfig.LBTCConfig,
			feeEstimatorConfig,
		)
	}

	relayOpts := types.NewRelayOpts(rargs)
	configWatcher, err := newStandardConfigProvider(ctx, r.lggr, r.chain, relayOpts)
	if err != nil {
		return nil, err
	}
	address := common.HexToAddress(relayOpts.ContractID)
	typ, ver, err := ccipconfig.TypeAndVersion(address, r.chain.Client())
	if err != nil {
		return nil, err
	}
	fn, err := ccipexec.ExecReportToEthTxMeta(ctx, typ, ver)
	if err != nil {
		return nil, err
	}
	subjectID := chainToUUID(configWatcher.chain.ID())

	contractTransmitter, err := newOnChainContractTransmitter(ctx, r.lggr, rargs, r.evmKeystore, configWatcher, configTransmitterOpts{
		subjectID: &subjectID,
	}, OCR2AggregatorTransmissionContractABI, WithReportToEthMetadata(fn), WithRetention(0), WithExcludeSignatures())
	if err != nil {
		return nil, err
	}

	return NewDstExecProvider(
		r.lggr,
		versionFinder,
		r.chain.Client(),
		r.chain.LogPoller(),
		execPluginConfig.DestStartBlock,
		contractTransmitter,
		configWatcher,
		r.chain.GasEstimator(),
		*r.chain.Config().EVM().GasEstimator().PriceMax().ToInt(),
		feeEstimatorConfig,
		r.chain.TxManager(),
		cciptypes.Address(rargs.ContractID),
	)
}

func (r *Relayer) NewLLOProvider(ctx context.Context, rargs commontypes.RelayArgs, pargs commontypes.PluginArgs) (commontypes.LLOProvider, error) {
	relayOpts := types.NewRelayOpts(rargs)
	relayConfig, err := relayOpts.RelayConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get relay config: %w", err)
	}
	if relayConfig.LLOConfigMode == "" {
		return nil, fmt.Errorf("LLOConfigMode must be specified in relayConfig for LLO jobs (can be either: %q or %q)", types.LLOConfigModeMercury, types.LLOConfigModeBlueGreen)
	}
	if relayConfig.ChainID.String() != r.chain.ID().String() {
		return nil, fmt.Errorf("internal error: chain id in spec does not match this relayer's chain: have %s expected %s", relayConfig.ChainID.String(), r.chain.ID().String())
	}
	if !relayConfig.EffectiveTransmitterID.Valid {
		return nil, pkgerrors.New("EffectiveTransmitterID must be specified")
	}
	if relayConfig.LLODONID == 0 {
		return nil, errors.New("donID must be specified in relayConfig for LLO jobs")
	}
	// ensure that child loggers are namespaced by job ID which ought to be globally unique
	lggr := r.lggr.Named(fmt.Sprintf("job-%d", rargs.JobID)).With("donID", relayConfig.LLODONID, "transmitterID", relayConfig.EffectiveTransmitterID.String)

	switch relayConfig.LLOConfigMode {
	case types.LLOConfigModeMercury, types.LLOConfigModeBlueGreen:
	default:
		return nil, fmt.Errorf("LLOConfigMode must be specified in relayConfig for LLO jobs (only %q or %q is currently supported)", types.LLOConfigModeMercury, types.LLOConfigModeBlueGreen)
	}

	cdcFactory, err := r.cdcFactory()
	if err != nil {
		return nil, err
	}

	configuratorAddress := common.HexToAddress(relayOpts.ContractID)
	return NewLLOProvider(ctx, lggr, pargs, r.retirementReportCache, r.chain, configuratorAddress, cdcFactory, relayConfig, relayOpts, r.csaKeystore, r.mercuryCfg, r.retirementReportCache, r.ds, r.mercuryPool, r.capabilitiesRegistry)
}

func (r *Relayer) NewFunctionsProvider(ctx context.Context, rargs commontypes.RelayArgs, pargs commontypes.PluginArgs) (commontypes.FunctionsProvider, error) {
	lggr := r.lggr.Named("FunctionsProvider").Named(rargs.ExternalJobID.String())
	// TODO(FUN-668): Not ready yet (doesn't implement FunctionsEvents() properly)
	return NewFunctionsProvider(ctx, r.chain, rargs, pargs, lggr, r.evmKeystore, functions.FunctionsPlugin)
}

// NewConfigProvider is called by bootstrap jobs
func (r *Relayer) NewConfigProvider(ctx context.Context, args commontypes.RelayArgs) (configProvider commontypes.ConfigProvider, err error) {
	lggr := r.lggr.Named("ConfigProvider").Named(args.ExternalJobID.String())
	relayOpts := types.NewRelayOpts(args)
	relayConfig, err := relayOpts.RelayConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get relay config: %w", err)
	}
	expectedChainID := relayConfig.ChainID.String()
	if expectedChainID != r.chain.ID().String() {
		return nil, fmt.Errorf("internal error: chain id in spec does not match this relayer's chain: have %s expected %s", relayConfig.ChainID.String(), r.chain.ID().String())
	}

	// Handle legacy jobs which did not yet specify provider type and
	// switched between median/mercury based on presence of feed ID
	if args.ProviderType == "" {
		if relayConfig.FeedID == nil {
			args.ProviderType = "median"
		} else if relayConfig.LLODONID > 0 {
			args.ProviderType = "llo"
		} else {
			args.ProviderType = "mercury"
		}
	}

	switch args.ProviderType {
	case "median":
		configProvider, err = newStandardConfigProvider(ctx, lggr, r.chain, relayOpts)
	case "mercury":
		configProvider, err = newMercuryConfigProvider(ctx, lggr, r.chain, relayOpts)
	case "llo":
		// Use NullRetirementReportCache since we never run LLO jobs on
		// bootstrap nodes, and there's no need to introduce a failure mode or
		// performance hit no matter how minor.
		configProvider, err = newLLOConfigProvider(ctx, lggr, r.chain, &retirement.NullRetirementReportCache{}, relayOpts)
	case "ocr3-capability":
		configProvider, err = newOCR3CapabilityConfigProvider(ctx, lggr, r.chain, relayOpts)
	default:
		return nil, fmt.Errorf("unrecognized provider type: %q", args.ProviderType)
	}

	if err != nil {
		// Never return (*configProvider)(nil)
		return nil, err
	}
	return configProvider, err
}

func FilterNamesFromRelayArgs(args commontypes.RelayArgs) (filterNames []string, err error) {
	var addr evmtypes.EIP55Address
	if addr, err = evmtypes.NewEIP55Address(args.ContractID); err != nil {
		return nil, err
	}
	var relayConfig types.RelayConfig
	if err = json.Unmarshal(args.RelayConfig, &relayConfig); err != nil {
		return nil, pkgerrors.WithStack(err)
	}

	if relayConfig.FeedID != nil {
		filterNames = []string{mercury.FilterName(addr.Address(), *relayConfig.FeedID)}
	} else {
		filterNames = []string{configPollerFilterName(addr.Address()), transmitterFilterName(addr.Address())}
	}
	return filterNames, err
}

type configWatcher struct {
	services.Service
	eng *services.Engine

	contractAddress  common.Address
	offchainDigester ocrtypes.OffchainConfigDigester
	configPoller     types.ConfigPoller
	chain            legacyevm.Chain
	runReplay        bool
	fromBlock        uint64
}

func newConfigWatcher(lggr logger.Logger,
	contractAddress common.Address,
	offchainDigester ocrtypes.OffchainConfigDigester,
	configPoller types.ConfigPoller,
	chain legacyevm.Chain,
	fromBlock uint64,
	runReplay bool,
) *configWatcher {
	cw := &configWatcher{
		contractAddress:  contractAddress,
		offchainDigester: offchainDigester,
		configPoller:     configPoller,
		chain:            chain,
		runReplay:        runReplay,
		fromBlock:        fromBlock,
	}
	cw.Service, cw.eng = services.Config{
		Name:           fmt.Sprintf("ConfigWatcher.%s", contractAddress),
		NewSubServices: nil,
		Start:          cw.start,
		Close:          cw.close,
	}.NewServiceEngine(lggr)
	return cw
}

func (c *configWatcher) start(ctx context.Context) error {
	if c.runReplay && c.fromBlock != 0 {
		// Only replay if it's a brand new job.
		c.eng.Go(func(ctx context.Context) {
			c.eng.Infow("starting replay for config", "fromBlock", c.fromBlock)
			if err := c.configPoller.Replay(ctx, int64(c.fromBlock)); err != nil {
				c.eng.Errorw("error replaying for config", "err", err)
			} else {
				c.eng.Infow("completed replaying for config", "fromBlock", c.fromBlock)
			}
		})
	}
	c.configPoller.Start()
	return nil
}

func (c *configWatcher) close() error {
	return c.configPoller.Close()
}

func (c *configWatcher) OffchainConfigDigester() ocrtypes.OffchainConfigDigester {
	return c.offchainDigester
}

func (c *configWatcher) ContractConfigTracker() ocrtypes.ContractConfigTracker {
	return c.configPoller
}

type configTransmitterOpts struct {
	// pluginGasLimit overrides the gas limit default provided in the config watcher.
	pluginGasLimit *uint32
	// subjectID overrides the queueing subject id (the job external id will be used by default).
	subjectID *uuid.UUID
}

// newOnChainContractTransmitter creates a new contract transmitter.
func newOnChainContractTransmitter(ctx context.Context, lggr logger.Logger, rargs commontypes.RelayArgs, ethKeystore keys.Store, configWatcher *configWatcher, opts configTransmitterOpts, transmissionContractABI abi.ABI, ocrTransmitterOpts ...OCRTransmitterOption) (ContractTransmitter, error) {
	transmitter, err := generateTransmitterFrom(ctx, rargs, ethKeystore, configWatcher, opts)
	if err != nil {
		return nil, err
	}

	evmTransmitter, err := NewOCRContractTransmitter(
		ctx,
		configWatcher.contractAddress,
		configWatcher.chain.Client(),
		transmissionContractABI,
		transmitter,
		configWatcher.chain.LogPoller(),
		lggr,
		ethKeystore,
		ocrTransmitterOpts...,
	)
	if err != nil {
		return nil, err
	}

	// This code path should only be called when running CCIP 1.5 jobs on Tron.
	// All other products should use the standard Tron relayer implementation.
	if configWatcher.chain.Config().EVM().ChainType() == chaintype.ChainTron {
		return NewTronContractTransmitter(ctx, TronContractTransmitterOpts{
			Logger:             lggr,
			TransmissionsCache: NewTronTransmissionsCache(evmTransmitter),
			Keystore:           ethKeystore,
			ConfigWatcher:      configWatcher,
			OCRTransmitterOpts: ocrTransmitterOpts,
		})
	}

	return evmTransmitter, err
}

// newOnChainDualContractTransmitter creates a new dual contract transmitter.
func newOnChainDualContractTransmitter(ctx context.Context, lggr logger.Logger, rargs commontypes.RelayArgs, ethKeystore keys.Store, configWatcher *configWatcher, opts configTransmitterOpts, transmissionContractABI abi.ABI, ocrTransmitterOpts ...OCRTransmitterOption) (*dualContractTransmitter, error) {
	transmitter, err := generateTransmitterFrom(ctx, rargs, ethKeystore, configWatcher, opts)
	if err != nil {
		return nil, err
	}

	return NewOCRDualContractTransmitter(
		ctx,
		configWatcher.contractAddress,
		configWatcher.chain.Client(),
		transmissionContractABI,
		transmitter,
		configWatcher.chain.LogPoller(),
		lggr,
		ethKeystore,
		ocrTransmitterOpts...,
	)
}

func NewContractTransmitter(ctx context.Context, lggr logger.Logger, rargs commontypes.RelayArgs, ethKeystore keys.Store, configWatcher *configWatcher, opts configTransmitterOpts, transmissionContractABI abi.ABI, dualTransmission bool, ocrTransmitterOpts ...OCRTransmitterOption) (ContractTransmitter, error) {
	if dualTransmission {
		return newOnChainDualContractTransmitter(ctx, lggr, rargs, ethKeystore, configWatcher, opts, transmissionContractABI, ocrTransmitterOpts...)
	}

	return newOnChainContractTransmitter(ctx, lggr, rargs, ethKeystore, configWatcher, opts, transmissionContractABI, ocrTransmitterOpts...)
}

type Keystore interface {
	keys.AddressChecker
	keys.RoundRobin
	keys.Locker
	keys.RawUnhashedSigner
}

func generateTransmitterFrom(ctx context.Context, rargs commontypes.RelayArgs, ethKeystore Keystore, configWatcher *configWatcher, opts configTransmitterOpts) (Transmitter, error) {
	var relayConfig types.RelayConfig
	if err := json.Unmarshal(rargs.RelayConfig, &relayConfig); err != nil {
		return nil, err
	}
	var fromAddresses []common.Address
	sendingKeys := relayConfig.SendingKeys
	if !relayConfig.EffectiveTransmitterID.Valid {
		return nil, pkgerrors.New("EffectiveTransmitterID must be specified")
	}
	effectiveTransmitterAddress := common.HexToAddress(relayConfig.EffectiveTransmitterID.String)

	sendingKeysLength := len(sendingKeys)
	if sendingKeysLength == 0 {
		return nil, pkgerrors.New("no sending keys provided")
	}

	// If we are using multiple sending keys, then a forwarder is needed to rotate transmissions.
	// Ensure that this forwarder is not set to a local sending key, and ensure our sending keys are enabled.
	for _, s := range sendingKeys {
		if sendingKeysLength > 1 && s == effectiveTransmitterAddress.String() {
			return nil, pkgerrors.New("the transmitter is a local sending key with transaction forwarding enabled")
		}
		if err := ethKeystore.CheckEnabled(ctx, common.HexToAddress(s)); err != nil {
			return nil, pkgerrors.Wrap(err, "one of the sending keys given is not enabled")
		}
		fromAddresses = append(fromAddresses, common.HexToAddress(s))
	}

	subject := rargs.ExternalJobID
	if opts.subjectID != nil {
		subject = *opts.subjectID
	}
	strategy := txmgrcommon.NewQueueingTxStrategy(subject, relayConfig.DefaultTransactionQueueDepth)

	var checker txm.TransmitCheckerSpec
	if relayConfig.SimulateTransactions {
		checker.CheckerType = txm.TransmitCheckerTypeSimulate
	}

	gasLimit := configWatcher.chain.Config().EVM().GasEstimator().LimitDefault()
	ocr2Limit := configWatcher.chain.Config().EVM().GasEstimator().LimitJobType().OCR2()
	if ocr2Limit != nil {
		gasLimit = uint64(*ocr2Limit)
	}
	if opts.pluginGasLimit != nil {
		gasLimit = uint64(*opts.pluginGasLimit)
	}

	var transmitter Transmitter
	var err error

	switch commontypes.OCR2PluginType(rargs.ProviderType) {
	case commontypes.Median:
		transmitter, err = ocrcommon.NewOCR2FeedsTransmitter(
			configWatcher.chain.TxManager(),
			fromAddresses,
			common.HexToAddress(rargs.ContractID),
			gasLimit,
			effectiveTransmitterAddress,
			strategy,
			checker,
			ethKeystore,
			relayConfig.DualTransmissionConfig,
		)
	case commontypes.CCIPExecution:
		transmitter, err = cciptransmitter.NewTransmitterWithStatusChecker(
			configWatcher.chain.TxManager(),
			fromAddresses,
			gasLimit,
			effectiveTransmitterAddress,
			strategy,
			checker,
			configWatcher.chain.ID(),
			ethKeystore,
		)
	default:
		transmitter, err = ocrcommon.NewTransmitter(
			configWatcher.chain.TxManager(),
			fromAddresses,
			gasLimit,
			effectiveTransmitterAddress,
			strategy,
			checker,
			ethKeystore,
		)
	}
	if err != nil {
		return nil, pkgerrors.Wrap(err, "failed to create transmitter")
	}
	return transmitter, nil
}

func (r *Relayer) NewContractWriter(_ context.Context, config []byte) (commontypes.ContractWriter, error) {
	var cfg types.ChainWriterConfig
	if err := json.Unmarshal(config, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshall chain writer config err: %s", err)
	}

	cfg.MaxGasPrice = r.chain.Config().EVM().GasEstimator().PriceMax()
	return NewChainWriterService(r.lggr, r.chain.Client(), r.chain.TxManager(), r.chain.GasEstimator(), cfg)
}

func (r *Relayer) NewContractReader(ctx context.Context, chainReaderConfig []byte) (commontypes.ContractReader, error) {
	cfg := &types.ChainReaderConfig{}
	if err := json.Unmarshal(chainReaderConfig, cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshall chain reader config err: %s", err)
	}

	return NewChainReaderService(ctx, r.lggr, r.chain.LogPoller(), r.chain.HeadTracker(), r.chain.Client(), *cfg)
}

func (r *Relayer) EVM() (commontypes.EVMService, error) {
	return r, nil
}

func (r *Relayer) GetTransactionFee(ctx context.Context, transactionID string) (*commontypes.TransactionFee, error) {
	return r.chain.TxManager().GetTransactionFee(ctx, transactionID)
}

func (r *Relayer) NewMedianProvider(ctx context.Context, rargs commontypes.RelayArgs, pargs commontypes.PluginArgs) (commontypes.MedianProvider, error) {
	lggr := logger.Sugared(r.lggr).Named("MedianProvider").Named(rargs.ExternalJobID.String())
	relayOpts := types.NewRelayOpts(rargs)
	relayConfig, err := relayOpts.RelayConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get relay config: %w", err)
	}
	expectedChainID := relayConfig.ChainID.String()
	if expectedChainID != r.chain.ID().String() {
		return nil, fmt.Errorf("internal error: chain id in spec does not match this relayer's chain: have %s expected %s", relayConfig.ChainID.String(), r.chain.ID().String())
	}
	if !common.IsHexAddress(relayOpts.ContractID) {
		return nil, fmt.Errorf("invalid contractID %s, expected hex address", relayOpts.ContractID)
	}
	contractID := common.HexToAddress(relayOpts.ContractID)

	configWatcher, err := newStandardConfigProvider(ctx, lggr, r.chain, relayOpts)
	if err != nil {
		return nil, err
	}

	reportCodec := evmreportcodec.ReportCodec{}

	ct, err := NewContractTransmitter(ctx, lggr, rargs, r.evmKeystore, configWatcher, configTransmitterOpts{}, OCR2AggregatorTransmissionContractABI, relayConfig.EnableDualTransmission)
	if err != nil {
		return nil, err
	}

	medianContract, err := newMedianContract(configWatcher.ContractConfigTracker(), configWatcher.contractAddress, configWatcher.chain, rargs.JobID, rargs.OracleSpecID, r.ds, lggr)
	if err != nil {
		return nil, err
	}

	medianProvider := medianProvider{
		lggr:                lggr.Named("MedianProvider"),
		configWatcher:       configWatcher,
		reportCodec:         reportCodec,
		contractTransmitter: ct,
		medianContract:      medianContract,
	}

	// allow fallback until chain reader is default and median contract is removed, but still log just in case
	var chainReaderService ChainReaderService
	if relayConfig.ChainReader != nil {
		if chainReaderService, err = NewChainReaderService(ctx, lggr, r.chain.LogPoller(), r.chain.HeadTracker(), r.chain.Client(), *relayConfig.ChainReader); err != nil {
			return nil, err
		}

		boundContracts := []commontypes.BoundContract{{Name: "median", Address: contractID.String()}}
		if err = chainReaderService.Bind(ctx, boundContracts); err != nil {
			return nil, err
		}
	} else {
		lggr.Info("ChainReader missing from RelayConfig; falling back to internal MedianContract")
	}
	medianProvider.chainReader = chainReaderService

	if relayConfig.Codec != nil {
		medianProvider.codec, err = codec.NewCodec(*relayConfig.Codec)
		if err != nil {
			return nil, err
		}
	} else {
		lggr.Info("Codec missing from RelayConfig; falling back to internal MedianContract")
	}

	return &medianProvider, nil
}

func (r *Relayer) NewAutomationProvider(ctx context.Context, rargs commontypes.RelayArgs, pargs commontypes.PluginArgs) (commontypes.AutomationProvider, error) {
	lggr := logger.Sugared(r.lggr).Named("AutomationProvider").Named(rargs.ExternalJobID.String())
	ocr2keeperRelayer := NewOCR2KeeperRelayer(r.ds, r.chain, lggr.Named("OCR2KeeperRelayer"), r.evmKeystore)

	return ocr2keeperRelayer.NewOCR2KeeperProvider(ctx, rargs, pargs)
}

var _ commontypes.MedianProvider = (*medianProvider)(nil)

type medianProvider struct {
	lggr                logger.Logger
	configWatcher       *configWatcher
	contractTransmitter ContractTransmitter
	reportCodec         median.ReportCodec
	medianContract      *medianContract
	chainReader         ChainReaderService
	codec               commontypes.Codec
	ms                  services.MultiStart
}

func (p *medianProvider) Name() string { return p.lggr.Name() }

func (p *medianProvider) Start(ctx context.Context) error {
	srvcs := []services.StartClose{p.configWatcher, p.contractTransmitter, p.medianContract}
	if p.chainReader != nil {
		srvcs = append(srvcs, p.chainReader)
	}

	return p.ms.Start(ctx, srvcs...)
}

func (p *medianProvider) Close() error { return p.ms.Close() }

func (p *medianProvider) Ready() error { return nil }

func (p *medianProvider) HealthReport() map[string]error {
	hp := map[string]error{p.Name(): p.Ready()}
	services.CopyHealth(hp, p.configWatcher.HealthReport())
	services.CopyHealth(hp, p.contractTransmitter.HealthReport())
	services.CopyHealth(hp, p.medianContract.HealthReport())
	return hp
}

func (p *medianProvider) ContractTransmitter() ocrtypes.ContractTransmitter {
	return p.contractTransmitter
}

func (p *medianProvider) ReportCodec() median.ReportCodec {
	return p.reportCodec
}

func (p *medianProvider) MedianContract() median.MedianContract {
	return p.medianContract
}

func (p *medianProvider) OnchainConfigCodec() median.OnchainConfigCodec {
	return median.StandardOnchainConfigCodec{}
}

func (p *medianProvider) OffchainConfigDigester() ocrtypes.OffchainConfigDigester {
	return p.configWatcher.OffchainConfigDigester()
}

func (p *medianProvider) ContractConfigTracker() ocrtypes.ContractConfigTracker {
	return p.configWatcher.ContractConfigTracker()
}

func (p *medianProvider) ContractReader() commontypes.ContractReader {
	return p.chainReader
}

func (p *medianProvider) Codec() commontypes.Codec {
	return p.codec
}

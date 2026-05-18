package evm

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"sync"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/estimatorconfig/interceptors/mantle"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	pkgerrors "github.com/pkg/errors"
	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/libocr/gethwrappers2/ocr2aggregator"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median/evmreportcodec"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	ocr3capability "github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/triggers"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
	coretypes "github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink/v2/core/logger"

	txmgrcommon "github.com/smartcontractkit/chainlink/v2/common/txmgr"
	txm "github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/llo"
	"github.com/smartcontractkit/chainlink/v2/core/services/llo/bm"
	"github.com/smartcontractkit/chainlink/v2/core/services/llo/mercurytransmitter"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/ccipcommit"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/ccipexec"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/estimatorconfig"
	cciptransmitter "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/transmitter"
	lloconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/llo/config"
	mercuryconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/mercury/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/codec"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/functions"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury"
	mercuryutils "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/utils"
	reportcodecv1 "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/v1/reportcodec"
	reportcodecv2 "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/v2/reportcodec"
	reportcodecv3 "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/v3/reportcodec"
	reportcodecv4 "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/v4/reportcodec"
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

var _ commontypes.Relayer = &Relayer{} //nolint:staticcheck

// The current PluginProvider interface does not support an error return. This was fine up until CCIP.
// CCIP is the first product to introduce the idea of incomplete implementations of a provider based on
// what chain (for CCIP, src or dest) the provider is created for. The Unimplemented* implementations below allow us to return
// a non nil value, which is hopefully a better developer experience should you find yourself using the right methods
// but on the *wrong* provider.

// [UnimplementedOffchainConfigDigester] satisfies the OCR OffchainConfigDigester interface
type UnimplementedOffchainConfigDigester struct{}

func (e UnimplementedOffchainConfigDigester) ConfigDigest(config ocrtypes.ContractConfig) (ocrtypes.ConfigDigest, error) {
	return ocrtypes.ConfigDigest{}, fmt.Errorf("unimplemented for this relayer")
}

func (e UnimplementedOffchainConfigDigester) ConfigDigestPrefix() (ocrtypes.ConfigDigestPrefix, error) {
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

func (u UnimplementedContractTransmitter) FromAccount() (ocrtypes.Account, error) {
	return "", fmt.Errorf("unimplemented for this relayer")
}

func (u UnimplementedContractTransmitter) LatestConfigDigestAndEpoch(ctx context.Context) (configDigest ocrtypes.ConfigDigest, epoch uint32, err error) {
	return ocrtypes.ConfigDigest{}, 0, fmt.Errorf("unimplemented for this relayer")
}

type Relayer struct {
	ds                   sqlutil.DataSource
	chain                legacyevm.Chain
	lggr                 logger.SugaredLogger
	ks                   CSAETHKeystore
	mercuryPool          wsrpc.Pool
	codec                commontypes.Codec
	capabilitiesRegistry coretypes.CapabilitiesRegistry

	// Mercury
	mercuryORM        mercury.ORM
	transmitterCfg    mercury.TransmitterConfig
	triggerCapability *triggers.MercuryTriggerService

	// LLO/data streams
	cdcFactory func() (llo.ChannelDefinitionCacheFactory, error)
}

type CSAETHKeystore interface {
	CSA() keystore.CSA
	Eth() keystore.Eth
}

type RelayerOpts struct {
	DS sqlutil.DataSource
	CSAETHKeystore
	MercuryPool          wsrpc.Pool
	TransmitterConfig    mercury.TransmitterConfig
	CapabilitiesRegistry coretypes.CapabilitiesRegistry
	HTTPClient           *http.Client
}

func (c RelayerOpts) Validate() error {
	var err error
	if c.DS == nil {
		err = errors.Join(err, errors.New("nil DataSource"))
	}
	if c.CSAETHKeystore == nil {
		err = errors.Join(err, errors.New("nil Keystore"))
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
	sugared := logger.Sugared(lggr).Named("Relayer")
	mercuryORM := mercury.NewORM(opts.DS)
	cdcFactory := sync.OnceValues(func() (llo.ChannelDefinitionCacheFactory, error) {
		chainSelector, err := chainselectors.SelectorFromChainId(chain.ID().Uint64())
		if err != nil {
			return nil, fmt.Errorf("failed to get chain selector for chain id %s: %w", chain.ID(), err)
		}
		lloORM := llo.NewORM(opts.DS, chainSelector)
		return llo.NewChannelDefinitionCacheFactory(sugared, lloORM, chain.LogPoller(), opts.HTTPClient), nil
	})
	relayer := &Relayer{
		ds:                   opts.DS,
		chain:                chain,
		lggr:                 logger.Sugared(sugared),
		ks:                   opts.CSAETHKeystore,
		mercuryPool:          opts.MercuryPool,
		cdcFactory:           cdcFactory,
		mercuryORM:           mercuryORM,
		transmitterCfg:       opts.TransmitterConfig,
		capabilitiesRegistry: opts.CapabilitiesRegistry,
	}

	// Initialize write target capability if configuration is defined
	if chain.Config().EVM().Workflow().ForwarderAddress() != nil {
		ctx := context.Background()
		if chain.Config().EVM().Workflow().GasLimitDefault() == nil {
			return nil, fmt.Errorf("unable to instantiate write target as default gas limit is not set")
		}

		capability, err := NewWriteTarget(ctx, relayer, chain, *chain.Config().EVM().Workflow().GasLimitDefault(),
			lggr)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize write target: %w", err)
		}
		if err := relayer.capabilitiesRegistry.Add(ctx, capability); err != nil {
			return nil, err
		}
		lggr.Infow("Registered write target", "chain_id", chain.ID())
	}

	return relayer, nil
}

func (r *Relayer) Name() string {
	return r.lggr.Name()
}

// Start does noop: no subservices started on relay start, but when the first job is started
func (r *Relayer) Start(context.Context) error {
	return nil
}

func (r *Relayer) Close() error {
	if r.triggerCapability != nil {
		return r.triggerCapability.Close()
	}
	return nil
}

// Ready does noop: always ready
func (r *Relayer) Ready() error {
	return r.chain.Ready()
}

func (r *Relayer) HealthReport() (report map[string]error) {
	report = make(map[string]error)
	maps.Copy(report, r.chain.HealthReport())
	return
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
func (r *Relayer) NewOCR3CapabilityProvider(rargs commontypes.RelayArgs, pargs commontypes.PluginArgs) (commontypes.OCR3CapabilityProvider, error) {
	// TODO https://smartcontract-it.atlassian.net/browse/BCF-2887
	ctx := context.Background()
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

	transmitter, err := newOnChainContractTransmitter(ctx, r.lggr, rargs, r.ks.Eth(), configWatcher, configTransmitterOpts{}, OCR2AggregatorTransmissionContractABI)
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

	fromAccount, err := pp.ContractTransmitter().FromAccount()
	if err != nil {
		return nil, err
	}

	return &ocr3CapabilityProvider{
		PluginProvider: pp,
		transmitter:    ocr3capability.NewContractTransmitter(r.lggr, r.capabilitiesRegistry, string(fromAccount)),
	}, nil
}

func (r *Relayer) NewPluginProvider(rargs commontypes.RelayArgs, pargs commontypes.PluginArgs) (commontypes.PluginProvider, error) {
	// TODO https://smartcontract-it.atlassian.net/browse/BCF-2887
	ctx := context.Background()
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

	transmitter, err := newOnChainContractTransmitter(ctx, r.lggr, rargs, r.ks.Eth(), configWatcher, configTransmitterOpts{}, OCR2AggregatorTransmissionContractABI)
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

func (r *Relayer) NewMercuryProvider(rargs commontypes.RelayArgs, pargs commontypes.PluginArgs) (commontypes.MercuryProvider, error) {
	// TODO https://smartcontract-it.atlassian.net/browse/BCF-2887
	ctx := context.Background()
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
	feedID := mercuryutils.FeedID(*relayConfig.FeedID)

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
	privKey, err := r.ks.CSA().Get(relayConfig.EffectiveTransmitterID.String)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "failed to get CSA key for mercury connection")
	}

	clients := make(map[string]wsrpc.Client)
	for _, server := range mercuryConfig.GetServers() {
		client, err := r.mercuryPool.Checkout(context.Background(), privKey, server.PubKey, server.URL)
		if err != nil {
			return nil, err
		}
		clients[server.URL] = client
	}

	// initialize trigger capability service lazily
	if relayConfig.EnableTriggerCapability && r.triggerCapability == nil {
		if r.capabilitiesRegistry == nil {
			lggr.Errorw("trigger capability is enabled but capabilities registry is not set")
		} else {
			r.triggerCapability = triggers.NewMercuryTriggerService(0, lggr)
			if err := r.triggerCapability.Start(ctx); err != nil {
				return nil, err
			}
			if err := r.capabilitiesRegistry.Add(ctx, r.triggerCapability); err != nil {
				return nil, err
			}
			lggr.Infow("successfully added trigger service to the Registry")
		}
	}

	// FIXME: We actually know the version here since it's in the feed ID, can
	// we use generics to avoid passing three of this?
	// https://smartcontract-it.atlassian.net/browse/MERC-1414
	reportCodecV1 := reportcodecv1.NewReportCodec(*relayConfig.FeedID, lggr.Named("ReportCodecV1"))
	reportCodecV2 := reportcodecv2.NewReportCodec(*relayConfig.FeedID, lggr.Named("ReportCodecV2"))
	reportCodecV3 := reportcodecv3.NewReportCodec(*relayConfig.FeedID, lggr.Named("ReportCodecV3"))
	reportCodecV4 := reportcodecv4.NewReportCodec(*relayConfig.FeedID, lggr.Named("ReportCodecV4"))

	var transmitterCodec mercury.TransmitterReportDecoder
	switch feedID.Version() {
	case 1:
		transmitterCodec = reportCodecV1
	case 2:
		transmitterCodec = reportCodecV2
	case 3:
		transmitterCodec = reportCodecV3
	case 4:
		transmitterCodec = reportCodecV4
	default:
		return nil, fmt.Errorf("invalid feed version %d", feedID.Version())
	}
	transmitter := mercury.NewTransmitter(lggr, r.transmitterCfg, clients, privKey.PublicKey, rargs.JobID, *relayConfig.FeedID, r.mercuryORM, transmitterCodec, r.triggerCapability)

	return NewMercuryProvider(cp, r.codec, NewMercuryChainReader(r.chain.HeadTracker()), transmitter, reportCodecV1, reportCodecV2, reportCodecV3, reportCodecV4, lggr), nil
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
func (r *Relayer) NewCCIPCommitProvider(rargs commontypes.RelayArgs, pargs commontypes.PluginArgs) (commontypes.CCIPCommitProvider, error) {
	// TODO https://smartcontract-it.atlassian.net/browse/BCF-2887
	ctx := context.Background()

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
			mantleInterceptor, iErr := mantle.NewInterceptor(ctx, r.chain.Client())
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
	contractTransmitter, err := newOnChainContractTransmitter(ctx, r.lggr, rargs, r.ks.Eth(), configWatcher, configTransmitterOpts{
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
		*contractTransmitter,
		configWatcher,
		feeEstimatorConfig,
	), nil
}

// NewCCIPExecProvider constructs a provider of type CCIPExecProvider. Since this is happening in the Relayer,
// which lives in a separate process from delegate which is requesting a provider, we need to wire in through pargs
// which *type* (impl) of CCIPExecProvider should be created. CCIP is currently a special case where the provider has a
// subset of implementations of the complete interface as certain contracts in a CCIP lane are only deployed on the src
// chain or on the dst chain. This results in the two implementations of providers: a src and dst implementation.
func (r *Relayer) NewCCIPExecProvider(rargs commontypes.RelayArgs, pargs commontypes.PluginArgs) (commontypes.CCIPExecProvider, error) {
	// TODO https://smartcontract-it.atlassian.net/browse/BCF-2887
	ctx := context.Background()

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
			mantleInterceptor, iErr := mantle.NewInterceptor(ctx, r.chain.Client())
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
	contractTransmitter, err := newOnChainContractTransmitter(ctx, r.lggr, rargs, r.ks.Eth(), configWatcher, configTransmitterOpts{
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

func (r *Relayer) NewLLOProvider(rargs commontypes.RelayArgs, pargs commontypes.PluginArgs) (commontypes.LLOProvider, error) {
	// TODO https://smartcontract-it.atlassian.net/browse/BCF-2887
	ctx := context.Background()

	relayOpts := types.NewRelayOpts(rargs)
	var relayConfig types.RelayConfig
	{
		var err error
		relayConfig, err = relayOpts.RelayConfig()
		if err != nil {
			return nil, fmt.Errorf("failed to get relay config: %w", err)
		}
	}

	var lloCfg lloconfig.PluginConfig
	if err := json.Unmarshal(pargs.PluginConfig, &lloCfg); err != nil {
		return nil, pkgerrors.WithStack(err)
	}
	if err := lloCfg.Validate(); err != nil {
		return nil, err
	}

	if relayConfig.ChainID.String() != r.chain.ID().String() {
		return nil, fmt.Errorf("internal error: chain id in spec does not match this relayer's chain: have %s expected %s", relayConfig.ChainID.String(), r.chain.ID().String())
	}
	cp, err := newLLOConfigProvider(ctx, r.lggr, r.chain, relayOpts)
	if err != nil {
		return nil, pkgerrors.WithStack(err)
	}

	if !relayConfig.EffectiveTransmitterID.Valid {
		return nil, pkgerrors.New("EffectiveTransmitterID must be specified")
	}
	privKey, err := r.ks.CSA().Get(relayConfig.EffectiveTransmitterID.String)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "failed to get CSA key for mercury connection")
	}

	// FIXME: Remove after benchmarking is done
	// https://smartcontract-it.atlassian.net/browse/MERC-3487
	var transmitter LLOTransmitter
	if lloCfg.BenchmarkMode {
		r.lggr.Info("Benchmark mode enabled, using dummy transmitter. NOTE: THIS WILL NOT TRANSMIT ANYTHING")
		transmitter = bm.NewTransmitter(r.lggr, fmt.Sprintf("%x", privKey.PublicKey))
	} else {
		clients := make(map[string]wsrpc.Client)
		for _, server := range lloCfg.GetServers() {
			client, err2 := r.mercuryPool.Checkout(context.Background(), privKey, server.PubKey, server.URL)
			if err2 != nil {
				return nil, err2
			}
			clients[server.URL] = client
		}
		transmitter = llo.NewTransmitter(llo.TransmitterOpts{
			Lggr:        r.lggr,
			FromAccount: fmt.Sprintf("%x", privKey.PublicKey), // NOTE: This may need to change if we support e.g. multiple tranmsmitters, to be a composite of all keys
			MercuryTransmitterOpts: mercurytransmitter.Opts{
				Lggr:        r.lggr,
				Cfg:         r.transmitterCfg,
				Clients:     clients,
				FromAccount: privKey.PublicKey,
				DonID:       relayConfig.LLODONID,
				ORM:         mercurytransmitter.NewORM(r.ds, relayConfig.LLODONID),
			},
		})
	}

	cdcFactory, err := r.cdcFactory()
	if err != nil {
		return nil, err
	}
	cdc, err := cdcFactory.NewCache(lloCfg)
	if err != nil {
		return nil, err
	}
	return NewLLOProvider(cp, transmitter, r.lggr, cdc), nil
}

func (r *Relayer) NewFunctionsProvider(rargs commontypes.RelayArgs, pargs commontypes.PluginArgs) (commontypes.FunctionsProvider, error) {
	// TODO https://smartcontract-it.atlassian.net/browse/BCF-2887
	ctx := context.Background()

	lggr := r.lggr.Named("FunctionsProvider").Named(rargs.ExternalJobID.String())
	// TODO(FUN-668): Not ready yet (doesn't implement FunctionsEvents() properly)
	return NewFunctionsProvider(ctx, r.chain, rargs, pargs, lggr, r.ks.Eth(), functions.FunctionsPlugin)
}

// NewConfigProvider is called by bootstrap jobs
func (r *Relayer) NewConfigProvider(args commontypes.RelayArgs) (configProvider commontypes.ConfigProvider, err error) {
	// TODO https://smartcontract-it.atlassian.net/browse/BCF-2887
	ctx := context.Background()

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
		configProvider, err = newLLOConfigProvider(ctx, lggr, r.chain, relayOpts)
	case "ocr3-capability":
		configProvider, err = newOCR3CapabilityConfigProvider(ctx, r.lggr, r.chain, relayOpts)
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
		// Only replay if it's a brand runReplay job.
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
func newOnChainContractTransmitter(ctx context.Context, lggr logger.Logger, rargs commontypes.RelayArgs, ethKeystore keystore.Eth, configWatcher *configWatcher, opts configTransmitterOpts, transmissionContractABI abi.ABI, ocrTransmitterOpts ...OCRTransmitterOption) (*contractTransmitter, error) {
	transmitter, err := generateTransmitterFrom(ctx, rargs, ethKeystore, configWatcher, opts)
	if err != nil {
		return nil, err
	}

	return NewOCRContractTransmitter(
		ctx,
		configWatcher.contractAddress,
		configWatcher.chain.Client(),
		transmissionContractABI,
		transmitter,
		configWatcher.chain.LogPoller(),
		lggr,
		ocrTransmitterOpts...,
	)
}

func generateTransmitterFrom(ctx context.Context, rargs commontypes.RelayArgs, ethKeystore keystore.Eth, configWatcher *configWatcher, opts configTransmitterOpts) (Transmitter, error) {
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
		if err := ethKeystore.CheckEnabled(ctx, common.HexToAddress(s), configWatcher.chain.Config().EVM().ChainID()); err != nil {
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
			configWatcher.chain.ID(),
			ethKeystore,
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
			configWatcher.chain.ID(),
			ethKeystore,
		)
	}
	if err != nil {
		return nil, pkgerrors.Wrap(err, "failed to create transmitter")
	}
	return transmitter, nil
}

func (r *Relayer) NewChainWriter(_ context.Context, config []byte) (commontypes.ChainWriter, error) {
	var cfg types.ChainWriterConfig
	if err := json.Unmarshal(config, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshall chain writer config err: %s", err)
	}

	return NewChainWriterService(r.lggr, r.chain.Client(), r.chain.TxManager(), r.chain.GasEstimator(), cfg)
}

func (r *Relayer) NewContractReader(chainReaderConfig []byte) (commontypes.ContractReader, error) {
	ctx := context.Background()
	cfg := &types.ChainReaderConfig{}
	if err := json.Unmarshal(chainReaderConfig, cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshall chain reader config err: %s", err)
	}

	return NewChainReaderService(ctx, r.lggr, r.chain.LogPoller(), r.chain.HeadTracker(), r.chain.Client(), *cfg)
}

func (r *Relayer) NewMedianProvider(rargs commontypes.RelayArgs, pargs commontypes.PluginArgs) (commontypes.MedianProvider, error) {
	// TODO https://smartcontract-it.atlassian.net/browse/BCF-2887
	ctx := context.Background()

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

	contractTransmitter, err := newOnChainContractTransmitter(ctx, lggr, rargs, r.ks.Eth(), configWatcher, configTransmitterOpts{}, OCR2AggregatorTransmissionContractABI)
	if err != nil {
		return nil, err
	}

	medianContract, err := newMedianContract(configWatcher.ContractConfigTracker(), configWatcher.contractAddress, configWatcher.chain, rargs.JobID, r.ds, lggr)
	if err != nil {
		return nil, err
	}

	medianProvider := medianProvider{
		lggr:                lggr.Named("MedianProvider"),
		configWatcher:       configWatcher,
		reportCodec:         reportCodec,
		contractTransmitter: contractTransmitter,
		medianContract:      medianContract,
	}

	// allow fallback until chain reader is default and median contract is removed, but still log just in case
	var chainReaderService ChainReaderService
	if relayConfig.ChainReader != nil {
		if chainReaderService, err = NewChainReaderService(ctx, lggr, r.chain.LogPoller(), r.chain.HeadTracker(), r.chain.Client(), *relayConfig.ChainReader); err != nil {
			return nil, err
		}

		boundContracts := []commontypes.BoundContract{{Name: "median", Address: contractID.String()}}
		if err = chainReaderService.Bind(context.Background(), boundContracts); err != nil {
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

func (r *Relayer) NewAutomationProvider(rargs commontypes.RelayArgs, pargs commontypes.PluginArgs) (commontypes.AutomationProvider, error) {
	lggr := logger.Sugared(r.lggr).Named("AutomationProvider").Named(rargs.ExternalJobID.String())
	ocr2keeperRelayer := NewOCR2KeeperRelayer(r.ds, r.chain, lggr.Named("OCR2KeeperRelayer"), r.ks.Eth())

	return ocr2keeperRelayer.NewOCR2KeeperProvider(rargs, pargs)
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

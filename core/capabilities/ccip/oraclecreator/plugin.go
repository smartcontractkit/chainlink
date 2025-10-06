package oraclecreator

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/smartcontractkit/chainlink-ccip/pkg/chainaccessor"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/libocr/commontypes"
	libocr3 "github.com/smartcontractkit/libocr/offchainreporting2plus"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	commitocr3 "github.com/smartcontractkit/chainlink-ccip/commit"
	"github.com/smartcontractkit/chainlink-ccip/commit/merkleroot/rmn"
	execocr3 "github.com/smartcontractkit/chainlink-ccip/execute"
	"github.com/smartcontractkit/chainlink-ccip/pkg/consts"
	"github.com/smartcontractkit/chainlink-ccip/pkg/contractreader"
	ccipreaderpkg "github.com/smartcontractkit/chainlink-ccip/pkg/reader"
	"github.com/smartcontractkit/chainlink-ccip/pluginconfig"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"

	_ "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipaptos"  // Register Aptos plugin config factories
	_ "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipevm"    // Register EVM plugin config factories
	_ "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipsolana" // Register Solana plugin config factories
	_ "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipsui"    // Register Sui plugin config factories
	_ "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipton"    // Register Ton plugin config factories
	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ocrimpls"
	cctypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr3/promwrapper"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization"
	"github.com/smartcontractkit/chainlink/v2/core/services/telemetry"
)

var _ cctypes.OracleCreator = &pluginOracleCreator{}

const (
	defaultCommitGasLimit = 500_000
	defaultExecGasLimit   = 6_500_000
)

// pluginOracleCreator creates oracles that reference plugins running
// in the same process as the chainlink node, i.e not LOOPPs.
type pluginOracleCreator struct {
	ocrKeyBundles         map[string]ocr2key.KeyBundle
	transmitters          map[types.RelayID][]string
	peerWrapper           *ocrcommon.SingletonPeerWrapper
	externalJobID         uuid.UUID
	jobID                 int32
	isNewlyCreatedJob     bool
	pluginConfig          job.JSONConfig
	db                    ocr3types.Database
	lggr                  logger.SugaredLogger
	monitoringEndpointGen telemetry.MonitoringEndpointGenerator
	bootstrapperLocators  []commontypes.BootstrapperLocator
	homeChainReader       ccipreaderpkg.HomeChain
	homeChainSelector     cciptypes.ChainSelector
	relayers              map[types.RelayID]loop.Relayer
	addressCodec          ccipcommon.AddressCodec
	p2pID                 p2pkey.KeyV2
}

func NewPluginOracleCreator(
	ocrKeyBundles map[string]ocr2key.KeyBundle,
	transmitters map[types.RelayID][]string,
	relayers map[types.RelayID]loop.Relayer,
	peerWrapper *ocrcommon.SingletonPeerWrapper,
	externalJobID uuid.UUID,
	jobID int32,
	isNewlyCreatedJob bool,
	pluginConfig job.JSONConfig,
	db ocr3types.Database,
	lggr logger.Logger,
	monitoringEndpointGen telemetry.MonitoringEndpointGenerator,
	bootstrapperLocators []commontypes.BootstrapperLocator,
	homeChainReader ccipreaderpkg.HomeChain,
	homeChainSelector cciptypes.ChainSelector,
	addressCodec ccipcommon.AddressCodec,
	p2pID p2pkey.KeyV2,
) cctypes.OracleCreator {
	return &pluginOracleCreator{
		ocrKeyBundles:         ocrKeyBundles,
		transmitters:          transmitters,
		relayers:              relayers,
		peerWrapper:           peerWrapper,
		externalJobID:         externalJobID,
		jobID:                 jobID,
		isNewlyCreatedJob:     isNewlyCreatedJob,
		pluginConfig:          pluginConfig,
		db:                    db,
		lggr:                  logger.Sugared(lggr),
		monitoringEndpointGen: monitoringEndpointGen,
		bootstrapperLocators:  bootstrapperLocators,
		homeChainReader:       homeChainReader,
		homeChainSelector:     homeChainSelector,
		addressCodec:          addressCodec,
		p2pID:                 p2pID,
	}
}

// Type implements types.OracleCreator.
func (i *pluginOracleCreator) Type() cctypes.OracleType {
	return cctypes.OracleTypePlugin
}

// Create implements types.OracleCreator.
func (i *pluginOracleCreator) Create(ctx context.Context, donID uint32, config cctypes.OCR3ConfigWithMeta) (cctypes.CCIPOracle, error) {
	pluginType := cctypes.PluginType(config.Config.PluginType)
	chainSelector := uint64(config.Config.ChainSelector)
	destChainFamily, err := chainsel.GetSelectorFamily(chainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to get chain family from selector %d: %w", config.Config.ChainSelector, err)
	}

	pluginServices, err := ccipcommon.GetPluginServices(i.lggr, destChainFamily)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize plugin config: %w", err)
	}

	// Create CCIP providers - this is the preferred way for plugins to access CCIP data
	ccipProviders, err := i.createCCIPProviders(
		ctx,
		pluginServices,
		config,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create CCIPProviders: %w", err)
	}

	// Populate extraDataCodecRegistry with codecs from CCIPProviders
	err = i.populateCodecRegistriesWithProviderCodecs(ccipProviders)
	if err != nil {
		return nil, fmt.Errorf("failed to populate extraDataCodecRegistry with codecs from CCIPProviders: %w", err)
	}

	destChainID, err := chainsel.GetChainIDFromSelector(chainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID from selector %d: %w", chainSelector, err)
	}
	destRelayID := types.NewRelayID(destChainFamily, destChainID)

	configTracker, err := ocrimpls.NewConfigTracker(config, i.addressCodec)
	if err != nil {
		return nil, fmt.Errorf("failed to create config tracker: %w, %d", err, chainSelector)
	}
	publicConfig, err := configTracker.PublicConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get public config from OCR config: %w", err)
	}

	i.lggr.Infow("Creating plugin using OCR3 settings",
		"plugin", pluginType.String(),
		"chainSelector", chainSelector,
		"chainID", destChainID,
		"deltaProgress", publicConfig.DeltaProgress,
		"deltaResend", publicConfig.DeltaResend,
		"deltaInitial", publicConfig.DeltaInitial,
		"deltaRound", publicConfig.DeltaRound,
		"deltaGrace", publicConfig.DeltaGrace,
		"deltaCertifiedCommitRequest", publicConfig.DeltaCertifiedCommitRequest,
		"deltaStage", publicConfig.DeltaStage,
		"rMax", publicConfig.RMax,
		"s", publicConfig.S,
		"maxDurationInitialization", publicConfig.MaxDurationInitialization,
		"maxDurationQuery", publicConfig.MaxDurationQuery,
		"maxDurationObservation", publicConfig.MaxDurationObservation,
		"maxDurationShouldAcceptAttestedReport", publicConfig.MaxDurationShouldAcceptAttestedReport,
		"maxDurationShouldTransmitAcceptedReport", publicConfig.MaxDurationShouldTransmitAcceptedReport,
	)

	offrampAddrStr, err := i.addressCodec.AddressBytesToString(config.Config.OfframpAddress, cciptypes.ChainSelector(chainSelector))
	if err != nil {
		return nil, fmt.Errorf("failed to convert offramp address to string using address codec: %w", err)
	}

	i.lggr.Infow("offramp address", "offrampAddrStr", config.Config.OfframpAddress, "selector", config.Config.ChainSelector)
	contractReaders, extendedReaders, chainWriters, err := i.createReadersAndWriters(
		ctx,
		pluginServices.ChainRW,
		destChainID,
		pluginType,
		config,
		publicConfig,
		destChainFamily,
		offrampAddrStr,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create readers and writers: %w", err)
	}

	// Create chain accessors and contract transmitters for relayers that supported them
	chainAccessors, contractTransmitters, err := i.getChainAccessorsAndContractTransmittersFromProviders(
		ccipProviders,
		extendedReaders,
		chainWriters,
		pluginServices,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create chain accessors: %w", err)
	}

	// build the onchain keyring. it will be the signing key for the destination chain family.
	keybundle, ok := i.ocrKeyBundles[destChainFamily]
	if !ok {
		return nil, fmt.Errorf("no OCR key bundle found for chain family %s, forgot to create one?", destChainFamily)
	}
	onchainKeyring := ocrimpls.NewOnchainKeyring[[]byte](keybundle, i.lggr)

	// build the contract transmitter
	// assume that we are using the first account in the keybundle as the from account.
	destChainWriter, ok := chainWriters[config.Config.ChainSelector]
	if !ok {
		i.lggr.Infow("no chain writer found for dest chain, will create nil transmitter",
			"destChainID", destChainID,
			"destChainSelector", config.Config.ChainSelector)
	}
	destFromAccounts, ok := i.transmitters[destRelayID]
	if !ok {
		i.lggr.Infow("no transmitters found for dest chain, will create nil transmitter",
			"destChainID", destChainID,
			"destChainSelector", config.Config.ChainSelector)
	}

	// TODO: Extract the correct transmitter address from the destsFromAccount
	factory, transmitter, err := i.createFactoryAndTransmitter(
		donID,
		config,
		destRelayID,
		chainAccessors,
		extendedReaders,
		chainWriters,
		destChainWriter,
		destFromAccounts,
		publicConfig,
		destChainFamily,
		destChainID,
		pluginServices,
		offrampAddrStr,
		contractTransmitters,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create factory and transmitter: %w", err)
	}

	telemetryType, err := pluginTypeToTelemetryType(pluginType)
	if err != nil {
		return nil, fmt.Errorf("failed to get telemetry type: %w", err)
	}

	oracleArgs := libocr3.OCR3OracleArgs[[]byte]{
		BinaryNetworkEndpointFactory: i.peerWrapper.Peer2,
		Database:                     i.db,
		// NOTE: when specifying V2Bootstrappers here we actually do NOT need to run a full bootstrap node!
		// Thus it is vital that the bootstrapper locators are correctly set in the job spec.
		V2Bootstrappers:       i.bootstrapperLocators,
		ContractConfigTracker: configTracker,
		ContractTransmitter:   transmitter,
		LocalConfig:           defaultLocalConfig(),
		Logger: ocrcommon.NewOCRWrapper(
			i.lggr.
				Named(fmt.Sprintf("CCIP%sOCR3", pluginType.String())).
				Named(destRelayID.String()).
				Named(offrampAddrStr),
			false,
			func(ctx context.Context, msg string) {}),
		MetricsRegisterer: prometheus.WrapRegistererWith(map[string]string{"name": fmt.Sprintf("commit-%d", config.Config.ChainSelector)}, prometheus.DefaultRegisterer),
		MonitoringEndpoint: i.monitoringEndpointGen.GenMonitoringEndpoint(
			destChainFamily,
			destRelayID.ChainID,
			offrampAddrStr,
			telemetryType,
		),
		OffchainConfigDigester: ocrimpls.NewConfigDigester(config.ConfigDigest),
		OffchainKeyring:        keybundle,
		OnchainKeyring:         onchainKeyring,
		ReportingPluginFactory: factory,
	}
	oracle, err := libocr3.NewOracle(oracleArgs)
	if err != nil {
		return nil, err
	}

	closers := make([]io.Closer, 0, len(extendedReaders)+len(chainWriters))
	for _, cr := range contractReaders {
		closers = append(closers, cr)
	}
	for _, cw := range chainWriters {
		closers = append(closers, cw)
	}
	return newWrappedOracle(oracle, closers), nil
}

func (i *pluginOracleCreator) createFactoryAndTransmitter(
	donID uint32,
	config cctypes.OCR3ConfigWithMeta,
	destRelayID types.RelayID,
	chainAccessors map[cciptypes.ChainSelector]cciptypes.ChainAccessor,
	extendedReaders map[cciptypes.ChainSelector]contractreader.Extended,
	chainWriters map[cciptypes.ChainSelector]types.ContractWriter,
	destChainWriter types.ContractWriter,
	destFromAccounts []string,
	publicConfig ocr3confighelper.PublicConfig,
	destChainFamily string,
	destChainID string,
	pluginServices ccipcommon.PluginServices,
	offrampAddrStr string,
	existingContractTransmitterMap map[cciptypes.ChainSelector]ocr3types.ContractTransmitter[[]byte],
) (ocr3types.ReportingPluginFactory[[]byte], ocr3types.ContractTransmitter[[]byte], error) {
	var factory ocr3types.ReportingPluginFactory[[]byte]
	var transmitter ocr3types.ContractTransmitter[[]byte]
	pluginConfig := pluginServices.PluginConfig
	if config.Config.PluginType == uint8(cctypes.PluginTypeCCIPCommit) {
		if !i.peerWrapper.IsStarted() {
			return nil, nil, errors.New("peer wrapper is not started")
		}

		i.lggr.Infow("creating rmn peer client",
			"bootstrapperLocators", i.bootstrapperLocators,
			"deltaRound", publicConfig.DeltaRound)

		rmnPeerClient := rmn.NewPeerClient(
			i.lggr.Named("RMNPeerClient"),
			i.peerWrapper.PeerGroupFactory,
			i.bootstrapperLocators,
			publicConfig.DeltaRound,
		)

		factory = commitocr3.NewCommitPluginFactory(
			commitocr3.CommitPluginFactoryParams{
				Lggr: i.lggr.
					Named("CCIPCommitPlugin").
					Named(destRelayID.String()).
					Named(fmt.Sprintf("%d", config.Config.ChainSelector)).
					Named(offrampAddrStr),
				DonID:                      donID,
				OcrConfig:                  ccipreaderpkg.OCR3ConfigWithMeta(config),
				CommitCodec:                pluginConfig.CommitPluginCodec,
				MsgHasher:                  pluginConfig.MessageHasher,
				AddrCodec:                  i.addressCodec,
				HomeChainReader:            i.homeChainReader,
				HomeChainSelector:          i.homeChainSelector,
				ChainAccessors:             chainAccessors,
				LOOPPCCIPProviderSupported: pluginServices.CCIPProviderSupported,
				ExtendedReaders:            extendedReaders,
				ContractWriters:            chainWriters,
				RmnPeerClient:              rmnPeerClient,
				RmnCrypto:                  pluginConfig.RMNCrypto})
		factory = promwrapper.NewReportingPluginFactory(
			factory,
			i.lggr,
			destChainFamily,
			destChainID,
			"CCIPCommit",
		)

		// there are three cases:
		//	1. contract transmitter is provided by the CCIP provider
		//  2. CCIP doesn't provide contract transmitter, we use CT factory with CW to create one
		//  3. Contract transmitter not supported, use noop transmitter
		ct, exist := existingContractTransmitterMap[config.Config.ChainSelector]
		switch {
		case exist && ct != nil:
			// case 1
			i.lggr.Infow("contracts transmitter provided from CCIP provider",
				"destChainID", destChainID,
				"destChainSelector", config.Config.ChainSelector)
			transmitter = ct
		case destChainWriter != nil:
			// case 2
			if len(destFromAccounts) == 0 {
				return nil, nil, fmt.Errorf("transmitter array is empty for dest relay ID %s", destRelayID)
			}
			transmitter = pluginConfig.ContractTransmitterFactory.NewCommitTransmitter(
				i.lggr.
					Named("CCIPCommitTransmitter").
					Named(destRelayID.String()).
					Named(fmt.Sprintf("%d", config.Config.ChainSelector)),
				destChainWriter,
				ocrtypes.Account(destFromAccounts[0]),
				offrampAddrStr,
				consts.MethodCommit,
				pluginConfig.PriceOnlyCommitFn,
			)
		default:
			// case 3
			i.lggr.Infow("no chain writer found for dest chain, creating nil transmitter",
				"destChainID", destChainID,
				"destChainSelector", config.Config.ChainSelector)
			transmitAccount, err := i.getTransmitterFromPublicConfig(publicConfig)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to get transmitter from public config: %w", err)
			}
			i.lggr.Infow("using (fake) transmitter from public config in the commit no-op transmitter", "transmitAccount", transmitAccount)
			transmitter = ocrimpls.NewNoOpTransmitter(
				i.lggr.
					Named("CCIPCommitNoOpTransmitter").
					Named(destRelayID.String()).
					Named(fmt.Sprintf("%d", config.Config.ChainSelector)),
				i.p2pID.PeerID().String(),
				transmitAccount,
			)
		}
	} else if config.Config.PluginType == uint8(cctypes.PluginTypeCCIPExec) {
		factory = execocr3.NewExecutePluginFactory(
			execocr3.PluginFactoryParams{
				Lggr: i.lggr.
					Named("CCIPExecPlugin").
					Named(destRelayID.String()).
					Named(fmt.Sprintf("%d", config.Config.ChainSelector)).
					Named(offrampAddrStr),
				DonID:                      donID,
				OcrConfig:                  ccipreaderpkg.OCR3ConfigWithMeta(config),
				ExecCodec:                  pluginConfig.ExecutePluginCodec,
				MsgHasher:                  pluginConfig.MessageHasher,
				AddrCodec:                  i.addressCodec,
				HomeChainReader:            i.homeChainReader,
				TokenDataEncoder:           pluginConfig.TokenDataEncoder,
				EstimateProvider:           pluginConfig.GasEstimateProvider,
				LOOPPCCIPProviderSupported: pluginServices.CCIPProviderSupported,
				ChainAccessors:             chainAccessors,
				ExtendedReaders:            extendedReaders,
				ContractWriters:            chainWriters,
			})
		factory = promwrapper.NewReportingPluginFactory(
			factory,
			i.lggr,
			destChainFamily,
			destChainID,
			"CCIPExec",
		)

		ct, exist := existingContractTransmitterMap[config.Config.ChainSelector]
		switch {
		case exist && ct != nil:
			// case 1
			i.lggr.Infow("contracts transmitter provided from CCIP provider",
				"destChainID", destChainID,
				"destChainSelector", config.Config.ChainSelector)
			transmitter = ct
		case destChainWriter != nil:
			// case 2
			if len(destFromAccounts) == 0 {
				return nil, nil, fmt.Errorf("transmitter array is empty for dest relay ID %s", destRelayID)
			}
			transmitter = pluginConfig.ContractTransmitterFactory.NewExecTransmitter(
				i.lggr.
					Named("CCIPExecTransmitter").
					Named(destRelayID.String()).
					Named(fmt.Sprintf("%d", config.Config.ChainSelector)),
				destChainWriter,
				ocrtypes.Account(destFromAccounts[0]),
				offrampAddrStr,
			)
		default:
			// case 3
			i.lggr.Infow("no chain writer found for dest chain, creating nil transmitter",
				"destChainID", destChainID,
				"destChainSelector", config.Config.ChainSelector)

			transmitAccount, err := i.getTransmitterFromPublicConfig(publicConfig)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to get transmitter from public config: %w", err)
			}
			i.lggr.Infow("using (fake) transmitter from public config in the exec no-op transmitter", "transmitAccount", transmitAccount)
			transmitter = ocrimpls.NewNoOpTransmitter(
				i.lggr.
					Named("CCIPExecNoOpTransmitter").
					Named(destRelayID.String()).
					Named(fmt.Sprintf("%d", config.Config.ChainSelector)),
				i.p2pID.PeerID().String(),
				transmitAccount,
			)
		}
	} else {
		return nil, nil, fmt.Errorf("unsupported Plugin type %d", config.Config.PluginType)
	}
	return factory, transmitter, nil
}

func (i *pluginOracleCreator) createCCIPProviders(
	ctx context.Context,
	pluginServices ccipcommon.PluginServices,
	config cctypes.OCR3ConfigWithMeta,
) (map[cciptypes.ChainSelector]types.CCIPProvider, error) {
	ccipProviders := make(map[cciptypes.ChainSelector]types.CCIPProvider)
	for relayID, relayer := range i.relayers {
		chainDetails, err := chainsel.GetChainDetailsByChainIDAndFamily(relayID.ChainID, relayID.Network)
		if err != nil {
			return nil, fmt.Errorf("failed to get chain selector from relay ID %s and family %s: %w", relayID.ChainID, relayID.Network, err)
		}
		chainSelector := cciptypes.ChainSelector(chainDetails.ChainSelector)

		ccipProviderSupported, ok := pluginServices.CCIPProviderSupported[relayID.Network]
		if ccipProviderSupported && ok {
			i.lggr.Debugw("creating CCIPProvider for chain family",
				"chainSelector", chainSelector, "chainFamily", relayID.Network)
			transmitter := i.transmitters[relayID]
			if len(transmitter) == 0 {
				return nil, errors.New("transmitter list is empty")
			}

			// Check if the transmitter string is a valid utf-8 string
			if !utf8.ValidString(transmitter[0]) {
				i.lggr.Errorw("transmitter contains invalid UTF-8",
					"transmitter", transmitter[0],
					"relayID.Network", relayID.Network,
					"chainSelector", chainSelector)
				return nil, fmt.Errorf("transmitter contains invalid UTF-8: %q", transmitter[0])
			}
			ccipProvider, err := relayer.NewCCIPProvider(ctx, types.CCIPProviderArgs{
				PluginType:           cciptypes.PluginType(config.Config.PluginType),
				OffRampAddress:       config.Config.OfframpAddress,
				TransmitterAddress:   cciptypes.UnknownEncodedAddress(transmitter[0]),
				ExtraDataCodecBundle: ccipcommon.GetExtraDataCodecRegistry(),
			})
			if err != nil {
				return nil, fmt.Errorf("failed to create CCIP provider for relay ID %s: %w", relayID, err)
			}
			ccipProviders[chainSelector] = ccipProvider
		}
	}
	return ccipProviders, nil
}

func (i *pluginOracleCreator) getChainAccessorsAndContractTransmittersFromProviders(
	ccipProviders map[cciptypes.ChainSelector]types.CCIPProvider,
	extendedReaders map[cciptypes.ChainSelector]contractreader.Extended,
	chainWriters map[cciptypes.ChainSelector]types.ContractWriter,
	pluginServices ccipcommon.PluginServices,
) (map[cciptypes.ChainSelector]cciptypes.ChainAccessor, map[cciptypes.ChainSelector]ocr3types.ContractTransmitter[[]byte], error) {
	chainAccessors := make(map[cciptypes.ChainSelector]cciptypes.ChainAccessor)
	contractTransmitters := make(map[cciptypes.ChainSelector]ocr3types.ContractTransmitter[[]byte])
	for relayID := range i.relayers {
		chainDetails, err := chainsel.GetChainDetailsByChainIDAndFamily(relayID.ChainID, relayID.Network)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get chain selector from relay ID %s and family %s: %w", relayID.ChainID, relayID.Network, err)
		}

		chainSelector := cciptypes.ChainSelector(chainDetails.ChainSelector)
		var ca cciptypes.ChainAccessor
		var ct ocr3types.ContractTransmitter[[]byte]

		// Check if a CCIPProvider exists for this chain selector, if so use its chain accessor and contract transmitter
		ccipProvider := ccipProviders[chainSelector]
		if ccipProvider != nil {
			ca = ccipProvider.ChainAccessor()
			if ca == nil {
				return nil, nil, fmt.Errorf("CCIPProvider for relay ID %s does not support chain accessor", relayID)
			}
			ct = ccipProvider.ContractTransmitter()
			if ct == nil {
				i.lggr.Warnw("contracts transmitter provided from CCIP provider is nil, will use default transmitter if possible",
					"relayID", relayID,
					"chainSelector", chainSelector,
				)
			}
		} else {
			// Use DefaultAccessor if CR and CW exist
			if extendedReaders[chainSelector] == nil || chainWriters[chainSelector] == nil {
				return nil, nil, fmt.Errorf("cannot create default chain accessor for relay ID %s, contract reader and chain writer need to be present", relayID)
			}
			ca, err = chainaccessor.NewDefaultAccessor(
				i.lggr,
				chainSelector,
				extendedReaders[chainSelector],
				chainWriters[chainSelector],
				pluginServices.AddrCodec,
			)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to create default chain accessor for relay ID %s: %w", relayID, err)
			}
		}

		chainAccessors[chainSelector] = ca
		// TODO ct can be nil, which is considered in createFactoryAndTransmitter case 1 check. But maybe better to move to if clause and remove the nil check
		contractTransmitters[chainSelector] = ct
	}
	return chainAccessors, contractTransmitters, nil
}

func (i *pluginOracleCreator) populateCodecRegistriesWithProviderCodecs(
	ccipProviders map[cciptypes.ChainSelector]types.CCIPProvider,
) error {
	edcr := ccipcommon.GetExtraDataCodecRegistry()
	for chainSelector, provider := range ccipProviders {
		codec := provider.Codec()
		chainFamily, err := chainsel.GetSelectorFamily(uint64(chainSelector))
		if err != nil {
			return fmt.Errorf("failed to get chain family from chain selector %d: %w", chainSelector, err)
		}

		sourceChainExtraDataCodec := codec.SourceChainExtraDataCodec
		if sourceChainExtraDataCodec != nil {
			edcr.RegisterCodec(chainFamily, sourceChainExtraDataCodec)
		} else {
			i.lggr.Warnw("CCIPProvider codec has no SourceChainExtraDataCodec", "chainSelector", chainSelector)
		}
	}
	return nil
}

func (i *pluginOracleCreator) getTransmitterFromPublicConfig(publicConfig ocr3confighelper.PublicConfig) (ocrtypes.Account, error) {
	var myIndex = -1
	for idx, identity := range publicConfig.OracleIdentities {
		if identity.PeerID == strings.TrimPrefix(i.p2pID.PeerID().String(), "p2p_") {
			myIndex = idx
			break
		}
	}

	if myIndex == -1 {
		return ocrtypes.Account(""), fmt.Errorf("no transmitter found for my peer id %s in public config", i.p2pID.PeerID().String())
	}

	return publicConfig.OracleIdentities[myIndex].TransmitAccount, nil
}

// createReadersAndWriters creates the contract readers and writers for the relayers
// that are available on this chainlink node.
//
// Relayers that are available on this node are exactly the chains that are enabled
// in the node TOML config.
//
// Since not every node will support every chain, we may not have a reader/writer for
// every chain that the role DON will be servicing.
func (i *pluginOracleCreator) createReadersAndWriters(
	ctx context.Context,
	crcw ccipcommon.MultiChainRW,
	destChainID string,
	pluginType cctypes.PluginType,
	config cctypes.OCR3ConfigWithMeta,
	publicCfg ocr3confighelper.PublicConfig,
	destChainFamily string,
	destAddrStr string,
) (
	map[cciptypes.ChainSelector]types.ContractReader,
	map[cciptypes.ChainSelector]contractreader.Extended,
	map[cciptypes.ChainSelector]types.ContractWriter,
	error,
) {
	ofc, err := decodeAndValidateOffchainConfig(pluginType, publicCfg)
	if err != nil {
		return nil, nil, nil, err
	}

	var execBatchGasLimit uint64
	if !ofc.ExecEmpty() {
		execBatchGasLimit = ofc.Execute.BatchGasLimit
	} else {
		// Set the default here so chain writer config validation doesn't fail.
		// For commit, this won't be used, so its harmless.
		execBatchGasLimit = defaultExecGasLimit
	}

	homeChainID, err := chainsel.GetChainIDFromSelector(uint64(i.homeChainSelector))
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get chain ID from chain selector %d: %w", i.homeChainSelector, err)
	}

	contractReaders := make(map[cciptypes.ChainSelector]types.ContractReader)
	extendedReaders := make(map[cciptypes.ChainSelector]contractreader.Extended)
	chainWriters := make(map[cciptypes.ChainSelector]types.ContractWriter)
	for relayID, relayer := range i.relayers {
		chainID := relayID.ChainID
		relayChainFamily := relayID.Network
		chainDetails, err1 := chainsel.GetChainDetailsByChainIDAndFamily(chainID, relayChainFamily)
		chainSelector := cciptypes.ChainSelector(chainDetails.ChainSelector)
		if err1 != nil {
			return nil, nil, nil, fmt.Errorf("failed to get chain selector from chain ID %s: %w", chainID, err1)
		}

		cr, err1 := crcw.GetChainReader(ctx, ccipcommon.ChainReaderProviderOpts{
			Lggr:            i.lggr,
			Relayer:         relayer,
			ChainID:         chainID,
			DestChainID:     destChainID,
			HomeChainID:     homeChainID,
			Ofc:             ofc,
			ChainSelector:   chainSelector,
			ChainFamily:     relayChainFamily,
			DestChainFamily: destChainFamily,
			Transmitters:    i.transmitters,
		})
		if err1 != nil {
			// Some Chain family might not need crcw to be created, and if createChainAccessorsAndContractTransmitters will catch error if it does
			i.lggr.Debugf("skipping creating reader and writers for chain %s, reader creation: %v", chainID, err1)
			continue
		}

		if chainID == destChainID && destChainFamily == relayChainFamily {
			offrampAddress := destAddrStr
			err2 := cr.Bind(ctx, []types.BoundContract{
				{
					Address: offrampAddress,
					Name:    consts.ContractNameOffRamp,
				},
			})
			if err2 != nil {
				return nil, nil, nil, fmt.Errorf("failed to bind chain reader for dest chain %s's offramp at %s: %w", chainID, offrampAddress, err2)
			}
		}

		if err2 := cr.Start(ctx); err2 != nil {
			return nil, nil, nil, fmt.Errorf("failed to start contract reader for chain %s: %w", chainID, err2)
		}

		var solanaChainWriterConfigVersion *string
		if ofc.Execute != nil {
			solanaChainWriterConfigVersion = ofc.Execute.SolanaChainWriterConfigVersion
		}
		cw, err1 := crcw.GetChainWriter(ctx, ccipcommon.ChainWriterProviderOpts{
			ChainID:                        chainID,
			Relayer:                        relayer,
			Transmitters:                   i.transmitters,
			ExecBatchGasLimit:              execBatchGasLimit,
			ChainFamily:                    relayChainFamily,
			OfframpProgramAddress:          config.Config.OfframpAddress,
			SolanaChainWriterConfigVersion: solanaChainWriterConfigVersion,
		})
		if err1 != nil {
			// Some Chain family might not need crcw to be created, and if createChainAccessorsAndContractTransmitters will catch error if it does
			i.lggr.Debugf("skipping creating chain writer for chain %s, writer creation: %v", chainID, err1)
			continue
		}

		if err4 := cw.Start(ctx); err4 != nil {
			return nil, nil, nil, fmt.Errorf("failed to start chain writer for chain %s: %w", chainID, err4)
		}

		extendedCr, err := wrapContractReaderInObservedExtended(i.lggr, cr, chainSelector)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to wrap contract reader for chain %s: %w", chainID, err)
		}

		contractReaders[chainSelector] = cr
		extendedReaders[chainSelector] = extendedCr
		chainWriters[chainSelector] = cw
	}
	return contractReaders, extendedReaders, chainWriters, nil
}

func wrapContractReaderInObservedExtended(
	lggr logger.Logger,
	contractReader types.ContractReader,
	chainSelector cciptypes.ChainSelector,
) (contractreader.Extended, error) {
	chainFamily, err1 := chainsel.GetSelectorFamily(uint64(chainSelector))
	if err1 != nil {
		return nil, fmt.Errorf("failed to get chain family from selector: %w", err1)
	}
	chainID, err1 := chainsel.GetChainIDFromSelector(uint64(chainSelector))
	if err1 != nil {
		return nil, fmt.Errorf("failed to get chain id from selector: %w", err1)
	}
	// NewExtendedContractReader() protects against double wrapping an extended reader.
	reader := contractreader.NewExtendedContractReader(
		contractreader.NewObserverReader(contractReader, lggr, chainFamily, chainID),
	)
	if reader == nil {
		return nil, fmt.Errorf("failed to create extended contract reader for chain selector %d", chainSelector)
	}
	return reader, nil
}

func decodeAndValidateOffchainConfig(
	pluginType cctypes.PluginType,
	publicConfig ocr3confighelper.PublicConfig,
) (ccipcommon.OffChainConfig, error) {
	var ofc ccipcommon.OffChainConfig
	switch pluginType {
	case cctypes.PluginTypeCCIPExec:
		execOffchainCfg, err1 := pluginconfig.DecodeExecuteOffchainConfig(publicConfig.ReportingPluginConfig)
		if err1 != nil {
			return ccipcommon.OffChainConfig{}, fmt.Errorf("failed to decode execute offchain config: %w, raw: %s", err1, string(publicConfig.ReportingPluginConfig))
		}
		if err2 := execOffchainCfg.ApplyDefaultsAndValidate(); err2 != nil {
			return ccipcommon.OffChainConfig{}, fmt.Errorf("failed to validate execute offchain config: %w", err2)
		}
		ofc.Execute = &execOffchainCfg
	case cctypes.PluginTypeCCIPCommit:
		commitOffchainCfg, err1 := pluginconfig.DecodeCommitOffchainConfig(publicConfig.ReportingPluginConfig)
		if err1 != nil {
			return ccipcommon.OffChainConfig{}, fmt.Errorf("failed to decode commit offchain config: %w, raw: %s", err1, string(publicConfig.ReportingPluginConfig))
		}
		if err2 := commitOffchainCfg.ApplyDefaultsAndValidate(); err2 != nil {
			return ccipcommon.OffChainConfig{}, fmt.Errorf("failed to validate commit offchain config: %w", err2)
		}
		ofc.Commit = &commitOffchainCfg
	}

	if !ofc.IsValid() {
		return ccipcommon.OffChainConfig{}, errors.New("invalid offchain config: both commit and exec configs are either set or unset")
	}
	return ofc, nil
}

func defaultLocalConfig() ocrtypes.LocalConfig {
	return ocrtypes.LocalConfig{
		DefaultMaxDurationInitialization: 30 * time.Second,
		BlockchainTimeout:                10 * time.Second,
		ContractConfigLoadTimeout:        10 * time.Second,
		// Config tracking is handled by the launcher, since we're doing blue-green
		// deployments we're not going to be using OCR's built-in config switching,
		// which always shuts down the previous instance.
		ContractConfigConfirmations:        1,
		SkipContractConfigConfirmations:    true,
		ContractConfigTrackerPollInterval:  10 * time.Second,
		ContractTransmitterTransmitTimeout: 10 * time.Second,
		DatabaseTimeout:                    10 * time.Second,
		MinOCR2MaxDurationQuery:            1 * time.Second,
		EnableTransmissionTelemetry:        true,
		DevelopmentMode:                    "false",
	}
}

func pluginTypeToTelemetryType(pluginType cctypes.PluginType) (synchronization.TelemetryType, error) {
	switch pluginType {
	case cctypes.PluginTypeCCIPCommit:
		return synchronization.OCR3CCIPCommit, nil
	case cctypes.PluginTypeCCIPExec:
		return synchronization.OCR3CCIPExec, nil
	default:
		return "", fmt.Errorf("unknown plugin type %d", pluginType)
	}
}

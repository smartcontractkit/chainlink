package oraclecreator

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink/v2/core/services/relay"

	"github.com/smartcontractkit/libocr/commontypes"
	libocr3 "github.com/smartcontractkit/libocr/offchainreporting2plus"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipevm"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipsolana"
	solanaconfig "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/configs/solana"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr3/promwrapper"

	commitocr3 "github.com/smartcontractkit/chainlink-ccip/commit"
	"github.com/smartcontractkit/chainlink-ccip/commit/merkleroot/rmn"
	execocr3 "github.com/smartcontractkit/chainlink-ccip/execute"
	"github.com/smartcontractkit/chainlink-ccip/pkg/consts"
	ccipreaderpkg "github.com/smartcontractkit/chainlink-ccip/pkg/reader"
	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-ccip/pluginconfig"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/chainwriter"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/config"

	evmconfig "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/configs/evm"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ocrimpls"
	cctypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
	evmrelaytypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization"
	"github.com/smartcontractkit/chainlink/v2/core/services/telemetry"

	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
)

var _ cctypes.OracleCreator = &pluginOracleCreator{}
var extraDataCodec = ccipcommon.NewExtraDataCodec(
	ccipcommon.NewExtraDataCodecParams(
		ccipevm.ExtraDataDecoder{},
		ccipsolana.ExtraDataDecoder{},
	),
)

var plugins = map[string]plugin{
	chainsel.FamilyEVM: {
		CommitPluginCodec:  ccipevm.NewCommitPluginCodecV1(),
		ExecutePluginCodec: ccipevm.NewExecutePluginCodecV1(extraDataCodec),
		MessageHasher: func(lggr logger.Logger) cciptypes.MessageHasher {
			return ccipevm.NewMessageHasherV1(lggr, extraDataCodec)
		},
		TokenDataEncoder:           ccipevm.NewEVMTokenDataEncoder(),
		GasEstimateProvider:        ccipevm.NewGasEstimateProvider(extraDataCodec),
		RMNCrypto:                  func(lggr logger.Logger) cciptypes.RMNCrypto { return ccipevm.NewEVMRMNCrypto(lggr) },
		ContractTransmitterFactory: &ocrimpls.EVMContractTransmitterFactory{},
	},
	chainsel.FamilySolana: {
		CommitPluginCodec:  ccipsolana.NewCommitPluginCodecV1(),
		ExecutePluginCodec: ccipsolana.NewExecutePluginCodecV1(extraDataCodec),
		MessageHasher: func(lggr logger.Logger) cciptypes.MessageHasher {
			return ccipsolana.NewMessageHasherV1(lggr, extraDataCodec)
		},
		TokenDataEncoder:           ccipsolana.NewSolanaTokenDataEncoder(),
		GasEstimateProvider:        ccipsolana.NewGasEstimateProvider(extraDataCodec),
		RMNCrypto:                  func(lggr logger.Logger) cciptypes.RMNCrypto { return nil },
		PriceOnlyCommitFn:          consts.MethodCommitPriceOnly,
		ContractTransmitterFactory: &ocrimpls.SVMContractTransmitterFactory{},
	},
}

const (
	defaultCommitGasLimit = 500_000
	defaultExecGasLimit   = 6_500_000
)

type plugin struct {
	CommitPluginCodec          cciptypes.CommitPluginCodec
	ExecutePluginCodec         cciptypes.ExecutePluginCodec
	MessageHasher              func(lggr logger.Logger) cciptypes.MessageHasher
	TokenDataEncoder           cciptypes.TokenDataEncoder
	GasEstimateProvider        cciptypes.EstimateProvider
	RMNCrypto                  func(lggr logger.Logger) cciptypes.RMNCrypto
	ContractTransmitterFactory ContractTransmitterFactory
	// PriceOnlyCommitFn optional method override for price only commit reports.
	PriceOnlyCommitFn string
}

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
	lggr                  logger.Logger
	monitoringEndpointGen telemetry.MonitoringEndpointGenerator
	bootstrapperLocators  []commontypes.BootstrapperLocator
	homeChainReader       ccipreaderpkg.HomeChain
	homeChainSelector     cciptypes.ChainSelector
	relayers              map[types.RelayID]loop.Relayer
	addressCodec          cciptypes.AddressCodec
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
	addressCodec cciptypes.AddressCodec,
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
		lggr:                  lggr,
		monitoringEndpointGen: monitoringEndpointGen,
		bootstrapperLocators:  bootstrapperLocators,
		homeChainReader:       homeChainReader,
		homeChainSelector:     homeChainSelector,
		addressCodec:          addressCodec,
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

	destChainID, err := chainsel.GetChainIDFromSelector(chainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID from selector %d: %w", chainSelector, err)
	}
	destRelayID := types.NewRelayID(destChainFamily, destChainID)

	configTracker := ocrimpls.NewConfigTracker(config, i.addressCodec)
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
		return nil, err
	}

	contractReaders, chainWriters, err := i.createReadersAndWriters(
		ctx,
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

	// build the onchain keyring. it will be the signing key for the destination chain family.
	keybundle, ok := i.ocrKeyBundles[destChainFamily]
	if !ok {
		return nil, fmt.Errorf("no OCR key bundle found for chain family %s, forgot to create one?", destChainFamily)
	}
	onchainKeyring := ocrimpls.NewOnchainKeyring[[]byte](keybundle, i.lggr)

	// build the contract transmitter
	// assume that we are using the first account in the keybundle as the from account
	// and that we are able to transmit to the dest chain.
	// TODO: revisit this in the future, since not all oracles will be able to transmit to the dest chain.
	destChainWriter, ok := chainWriters[config.Config.ChainSelector]
	if !ok {
		return nil, fmt.Errorf("no chain writer found for dest chain selector %d, can't create contract transmitter",
			config.Config.ChainSelector)
	}
	destFromAccounts, ok := i.transmitters[destRelayID]
	if !ok {
		return nil, fmt.Errorf("no transmitter found for dest relay ID %s, can't create contract transmitter", destRelayID)
	}

	// TODO: Extract the correct transmitter address from the destsFromAccount
	factory, transmitter, err := i.createFactoryAndTransmitter(
		donID, config, destRelayID, contractReaders, chainWriters, destChainWriter, destFromAccounts, publicConfig, offrampAddrStr)
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

	closers := make([]io.Closer, 0, len(contractReaders)+len(chainWriters))
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
	contractReaders map[cciptypes.ChainSelector]types.ContractReader,
	chainWriters map[cciptypes.ChainSelector]types.ContractWriter,
	destChainWriter types.ContractWriter,
	destFromAccounts []string,
	publicConfig ocr3confighelper.PublicConfig,
	offrampAddrStr string,
) (ocr3types.ReportingPluginFactory[[]byte], ocr3types.ContractTransmitter[[]byte], error) {
	var factory ocr3types.ReportingPluginFactory[[]byte]
	var transmitter ocr3types.ContractTransmitter[[]byte]

	chainID, err := chainsel.GetChainIDFromSelector(uint64(config.Config.ChainSelector))
	if err != nil {
		return nil, nil, fmt.Errorf("unsupported chain selector %d %w", config.Config.ChainSelector, err)
	}

	chainFamily, err := chainsel.GetSelectorFamily(uint64(config.Config.ChainSelector))
	if err != nil {
		return nil, nil, fmt.Errorf("unsupported chain selector %d %w", config.Config.ChainSelector, err)
	}
	plugin, exists := plugins[chainFamily]
	if !exists {
		return nil, nil, fmt.Errorf("unsupported chain %v", chainFamily)
	}
	messageHasher := plugin.MessageHasher(i.lggr.Named(chainFamily).Named("MessageHasherV1"))

	if config.Config.PluginType == uint8(cctypes.PluginTypeCCIPCommit) {
		if !i.peerWrapper.IsStarted() {
			return nil, nil, errors.New("peer wrapper is not started")
		}

		i.lggr.Infow("creating rmn peer client",
			"bootstrapperLocators", i.bootstrapperLocators, "deltaRound", publicConfig.DeltaRound)

		rmnPeerClient := rmn.NewPeerClient(
			i.lggr.Named("RMNPeerClient"),
			i.peerWrapper.PeerGroupFactory,
			i.bootstrapperLocators,
			publicConfig.DeltaRound,
		)

		rmnCrypto := plugin.RMNCrypto(i.lggr.Named(chainFamily).Named("RMNCrypto"))
		factory = commitocr3.NewCommitPluginFactory(
			commitocr3.CommitPluginFactoryParams{
				Lggr: i.lggr.
					Named("CCIPCommitPlugin").
					Named(destRelayID.String()).
					Named(fmt.Sprintf("%d", config.Config.ChainSelector)).
					Named(offrampAddrStr),
				DonID:             donID,
				OcrConfig:         ccipreaderpkg.OCR3ConfigWithMeta(config),
				CommitCodec:       plugin.CommitPluginCodec,
				MsgHasher:         messageHasher,
				AddrCodec:         i.addressCodec,
				HomeChainReader:   i.homeChainReader,
				HomeChainSelector: i.homeChainSelector,
				ContractReaders:   contractReaders,
				ContractWriters:   chainWriters,
				RmnPeerClient:     rmnPeerClient,
				RmnCrypto:         rmnCrypto})
		factory = promwrapper.NewReportingPluginFactory[[]byte](factory, i.lggr, chainID, "CCIPCommit")
		transmitter = plugins[chainFamily].ContractTransmitterFactory.NewCommitTransmitter(
			i.lggr.Named("CCIPCommitTransmitter").
				Named(destRelayID.String()),
			destChainWriter,
			ocrtypes.Account(destFromAccounts[0]),
			offrampAddrStr,
			consts.MethodCommit,
			plugins[chainFamily].PriceOnlyCommitFn,
		)
	} else if config.Config.PluginType == uint8(cctypes.PluginTypeCCIPExec) {
		factory = execocr3.NewExecutePluginFactory(
			execocr3.PluginFactoryParams{
				Lggr: i.lggr.
					Named("CCIPExecPlugin").
					Named(destRelayID.String()).
					Named(offrampAddrStr),
				DonID:            donID,
				OcrConfig:        ccipreaderpkg.OCR3ConfigWithMeta(config),
				ExecCodec:        plugin.ExecutePluginCodec,
				MsgHasher:        messageHasher,
				AddrCodec:        i.addressCodec,
				HomeChainReader:  i.homeChainReader,
				TokenDataEncoder: plugin.TokenDataEncoder,
				EstimateProvider: plugin.GasEstimateProvider,
				ContractReaders:  contractReaders,
				ContractWriters:  chainWriters,
			})
		factory = promwrapper.NewReportingPluginFactory[[]byte](factory, i.lggr, chainID, "CCIPExec")
		transmitter = plugins[chainFamily].ContractTransmitterFactory.NewExecTransmitter(
			i.lggr.Named("CCIPExecTransmitter").
				Named(destRelayID.String()),
			destChainWriter,
			ocrtypes.Account(destFromAccounts[0]),
			offrampAddrStr,
		)
	} else {
		return nil, nil, fmt.Errorf("unsupported plugin type %d", config.Config.PluginType)
	}
	return factory, transmitter, nil
}

func (i *pluginOracleCreator) createReadersAndWriters(
	ctx context.Context,
	destChainID string,
	pluginType cctypes.PluginType,
	config cctypes.OCR3ConfigWithMeta,
	publicCfg ocr3confighelper.PublicConfig,
	destChainFamily string,
	destAddrStr string,
) (
	map[cciptypes.ChainSelector]types.ContractReader,
	map[cciptypes.ChainSelector]types.ContractWriter,
	error,
) {
	ofc, err := decodeAndValidateOffchainConfig(pluginType, publicCfg)
	if err != nil {
		return nil, nil, err
	}

	var execBatchGasLimit uint64
	if !ofc.execEmpty() {
		execBatchGasLimit = ofc.exec().BatchGasLimit
	} else {
		// Set the default here so chain writer config validation doesn't fail.
		// For commit, this won't be used, so its harmless.
		execBatchGasLimit = defaultExecGasLimit
	}

	homeChainID, err := chainsel.GetChainIDFromSelector(uint64(i.homeChainSelector))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get chain ID from chain selector %d: %w", i.homeChainSelector, err)
	}

	contractReaders := make(map[cciptypes.ChainSelector]types.ContractReader)
	chainWriters := make(map[cciptypes.ChainSelector]types.ContractWriter)
	for relayID, relayer := range i.relayers {
		chainID := relayID.ChainID
		relayChainFamily := relayID.Network
		chainDetails, err1 := chainsel.GetChainDetailsByChainIDAndFamily(chainID, relayChainFamily)
		chainSelector := cciptypes.ChainSelector(chainDetails.ChainSelector)
		if err1 != nil {
			return nil, nil, fmt.Errorf("failed to get chain selector from chain ID %s: %w", chainID, err1)
		}

		if _, exists := plugins[relayChainFamily]; !exists {
			i.lggr.Debugw("createReadersAndWriters: skipping unsupported relayer", "chainID", chainID, "family", relayChainFamily)
			continue
		}

		chainReaderConfig, err1 := getChainReaderConfig(i.lggr, chainID, destChainID, homeChainID, ofc, chainSelector, relayChainFamily)
		if err1 != nil {
			return nil, nil, fmt.Errorf("failed to get chain reader config: %w", err1)
		}

		cr, err1 := relayer.NewContractReader(ctx, chainReaderConfig)
		if err1 != nil {
			return nil, nil, err1
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
				return nil, nil, fmt.Errorf("failed to bind chain reader for dest chain %s's offramp at %s: %w", chainID, offrampAddress, err)
			}
		}

		if err2 := cr.Start(ctx); err2 != nil {
			return nil, nil, fmt.Errorf("failed to start contract reader for chain %s: %w", chainID, err2)
		}

		cw, err1 := createChainWriter(
			ctx,
			chainID,
			relayer,
			i.transmitters,
			execBatchGasLimit,
			relayChainFamily,
			config.Config.OfframpAddress,
		)
		if err1 != nil {
			return nil, nil, err1
		}

		if err4 := cw.Start(ctx); err4 != nil {
			return nil, nil, fmt.Errorf("failed to start chain writer for chain %s: %w", chainID, err4)
		}

		contractReaders[chainSelector] = cr
		chainWriters[chainSelector] = cw
	}
	return contractReaders, chainWriters, nil
}

func decodeAndValidateOffchainConfig(
	pluginType cctypes.PluginType,
	publicConfig ocr3confighelper.PublicConfig,
) (offChainConfig, error) {
	var ofc offChainConfig
	if pluginType == cctypes.PluginTypeCCIPExec {
		execOffchainCfg, err1 := pluginconfig.DecodeExecuteOffchainConfig(publicConfig.ReportingPluginConfig)
		if err1 != nil {
			return offChainConfig{}, fmt.Errorf("failed to decode execute offchain config: %w, raw: %s", err1, string(publicConfig.ReportingPluginConfig))
		}
		if err2 := execOffchainCfg.Validate(); err2 != nil {
			return offChainConfig{}, fmt.Errorf("failed to validate execute offchain config: %w", err2)
		}
		ofc.execOffchainConfig = &execOffchainCfg
	} else if pluginType == cctypes.PluginTypeCCIPCommit {
		commitOffchainCfg, err1 := pluginconfig.DecodeCommitOffchainConfig(publicConfig.ReportingPluginConfig)
		if err1 != nil {
			return offChainConfig{}, fmt.Errorf("failed to decode commit offchain config: %w, raw: %s", err1, string(publicConfig.ReportingPluginConfig))
		}
		if err2 := commitOffchainCfg.ApplyDefaultsAndValidate(); err2 != nil {
			return offChainConfig{}, fmt.Errorf("failed to validate commit offchain config: %w", err2)
		}
		ofc.commitOffchainConfig = &commitOffchainCfg
	}
	if !ofc.isValid() {
		return offChainConfig{}, errors.New("invalid offchain config: both commit and exec configs are either set or unset")
	}
	return ofc, nil
}

func getChainReaderConfig(
	lggr logger.Logger,
	chainID string,
	destChainID string,
	homeChainID string,
	ofc offChainConfig,
	chainSelector cciptypes.ChainSelector,
	chainFamily string,
) ([]byte, error) {
	// TODO: create a chain writer constructor interface and define family specific implementations in oraclecreator.plugin
	switch chainFamily {
	case relay.NetworkEVM:
		var chainReaderConfig evmrelaytypes.ChainReaderConfig
		if chainID == destChainID {
			chainReaderConfig = evmconfig.DestReaderConfig
		} else {
			chainReaderConfig = evmconfig.SourceReaderConfig
		}

		if !ofc.commitEmpty() && ofc.commit().PriceFeedChainSelector == chainSelector {
			lggr.Debugw("Adding feed reader config", "chainID", chainID)
			chainReaderConfig = evmconfig.MergeReaderConfigs(chainReaderConfig, evmconfig.FeedReaderConfig)
		}

		if isUSDCEnabled(ofc) {
			lggr.Debugw("Adding USDC reader config", "chainID", chainID)
			chainReaderConfig = evmconfig.MergeReaderConfigs(chainReaderConfig, evmconfig.USDCReaderConfig)
		}

		if chainID == homeChainID {
			lggr.Debugw("Adding home chain reader config", "chainID", chainID)
			chainReaderConfig = evmconfig.MergeReaderConfigs(chainReaderConfig, evmconfig.HomeChainReaderConfigRaw)
		}

		marshaledConfig, err := json.Marshal(chainReaderConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal chain reader config: %w", err)
		}

		return marshaledConfig, nil
	case relay.NetworkSolana:
		var err error
		var cfg config.ContractReader
		if chainID == destChainID {
			cfg, err = solanaconfig.DestContractReaderConfig()
			if err != nil {
				return nil, fmt.Errorf("failed to get Solana dest contract reader config: %w", err)
			}
		} else {
			cfg, err = solanaconfig.SourceContractReaderConfig()
			if err != nil {
				return nil, fmt.Errorf("failed to get Solana source contract reader config: %w", err)
			}
		}

		marshaledConfig, err := json.Marshal(cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal chain reader config: %w", err)
		}

		return marshaledConfig, nil
	default:
		return nil, fmt.Errorf("unsupported chain family %s", chainFamily)
	}
}

func isUSDCEnabled(ofc offChainConfig) bool {
	if ofc.execEmpty() {
		return false
	}

	return ofc.exec().IsUSDCEnabled()
}

func createChainWriter(
	ctx context.Context,
	chainID string,
	relayer loop.Relayer,
	transmitters map[types.RelayID][]string,
	execBatchGasLimit uint64,
	chainFamily string,
	offrampAddress []byte,
) (types.ContractWriter, error) {
	var err error
	var chainWriterConfig []byte
	transmitter, ok := transmitters[types.NewRelayID(chainFamily, chainID)]
	// TODO: create a chain writer constructor interface and define family specific implementations in oraclecreator.plugin
	switch chainFamily {
	case relay.NetworkSolana:
		var solConfig chainwriter.ChainWriterConfig
		var offrampProgramAddress solana.PublicKey
		// NOTE: this function can still be called with EVM inputs, and PublicKeyFromBytes will panic on addresses with len=20
		// technically we only need the writer to do fee estimation so this doesn't matter and we can use a zero address
		if len(offrampAddress) == solana.PublicKeyLength {
			offrampProgramAddress = solana.PublicKeyFromBytes(offrampAddress)
		}
		if solConfig, err = solanaconfig.GetSolanaChainWriterConfig(offrampProgramAddress.String(), transmitter[0]); err != nil {
			return nil, fmt.Errorf("failed to get Solana chain writer config: %w", err)
		}
		if chainWriterConfig, err = json.Marshal(solConfig); err != nil {
			return nil, fmt.Errorf("failed to marshal Solana chain writer config: %w", err)
		}
	case relay.NetworkEVM:
		var evmConfig evmrelaytypes.ChainWriterConfig
		fromAddress := common.Address{}
		if ok {
			fromAddress = common.HexToAddress(transmitter[0])
		}
		if evmConfig, err = evmconfig.ChainWriterConfigRaw(
			fromAddress,
			defaultCommitGasLimit,
			execBatchGasLimit); err != nil {
			return nil, fmt.Errorf("failed to create EVM chain writer config: %w", err)
		}
		if chainWriterConfig, err = json.Marshal(evmConfig); err != nil {
			return nil, fmt.Errorf("failed to marshal EVM chain writer config: %w", err)
		}
	default:
		return nil, fmt.Errorf("unknown chain family %s", chainFamily)
	}

	cw, err := relayer.NewContractWriter(ctx, chainWriterConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create chain writer for chain %s: %w", chainID, err)
	}

	return cw, nil
}

type offChainConfig struct {
	commitOffchainConfig *pluginconfig.CommitOffchainConfig
	execOffchainConfig   *pluginconfig.ExecuteOffchainConfig
}

func (ofc offChainConfig) commitEmpty() bool {
	return ofc.commitOffchainConfig == nil
}

func (ofc offChainConfig) execEmpty() bool {
	return ofc.execOffchainConfig == nil
}

func (ofc offChainConfig) commit() *pluginconfig.CommitOffchainConfig {
	return ofc.commitOffchainConfig
}

func (ofc offChainConfig) exec() *pluginconfig.ExecuteOffchainConfig {
	return ofc.execOffchainConfig
}

// Exactly one of both plugins should be empty at any given time.
func (ofc offChainConfig) isValid() bool {
	return (ofc.commitEmpty() && !ofc.execEmpty()) || (!ofc.commitEmpty() && ofc.execEmpty())
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

package oraclecreator

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/libocr/commontypes"
	libocr3 "github.com/smartcontractkit/libocr/offchainreporting2plus"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"

	commitocr3 "github.com/smartcontractkit/chainlink-ccip/commit"
	"github.com/smartcontractkit/chainlink-ccip/commit/merkleroot/rmn"
	execocr3 "github.com/smartcontractkit/chainlink-ccip/execute"
	"github.com/smartcontractkit/chainlink-ccip/pkg/consts"
	ccipreaderpkg "github.com/smartcontractkit/chainlink-ccip/pkg/reader"
	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-ccip/pluginconfig"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	_ "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipevm"    // Register EVM plugin config factories
	_ "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipsolana" // Register Solana plugin config factories
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ocrimpls"
	cctypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
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
	lggr                  logger.Logger
	monitoringEndpointGen telemetry.MonitoringEndpointGenerator
	bootstrapperLocators  []commontypes.BootstrapperLocator
	homeChainReader       ccipreaderpkg.HomeChain
	homeChainSelector     cciptypes.ChainSelector
	relayers              map[types.RelayID]loop.Relayer
	addressCodec          cciptypes.AddressCodec
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
	addressCodec cciptypes.AddressCodec,
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
		lggr:                  lggr,
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

	pluginServices, err := ccipcommon.GetPluginServices(i.lggr, destChainFamily)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize plugin config: %w", err)
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
	contractReaders, chainWriters, err := i.createReadersAndWriters(
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
		donID, config, destRelayID, contractReaders, chainWriters, destChainWriter, destFromAccounts, publicConfig, destChainID, pluginServices.PluginConfig, offrampAddrStr)
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
	destChainID string,
	pluginConfig ccipcommon.PluginConfig,
	offrampAddrStr string,
) (ocr3types.ReportingPluginFactory[[]byte], ocr3types.ContractTransmitter[[]byte], error) {
	var factory ocr3types.ReportingPluginFactory[[]byte]
	var transmitter ocr3types.ContractTransmitter[[]byte]
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
				DonID:             donID,
				OcrConfig:         ccipreaderpkg.OCR3ConfigWithMeta(config),
				CommitCodec:       pluginConfig.CommitPluginCodec,
				MsgHasher:         pluginConfig.MessageHasher,
				AddrCodec:         i.addressCodec,
				HomeChainReader:   i.homeChainReader,
				HomeChainSelector: i.homeChainSelector,
				ContractReaders:   contractReaders,
				ContractWriters:   chainWriters,
				RmnPeerClient:     rmnPeerClient,
				RmnCrypto:         pluginConfig.RMNCrypto})
		factory = promwrapper.NewReportingPluginFactory(factory, i.lggr, destChainID, "CCIPCommit")
		if destChainWriter == nil {
			i.lggr.Infow("no chain writer found for dest chain, creating nil transmitter",
				"destChainID", destChainID,
				"destChainSelector", config.Config.ChainSelector)
			transmitter = ocrimpls.NewNoOpTransmitter(
				i.lggr.
					Named("CCIPCommitNoOpTransmitter").
					Named(destRelayID.String()).
					Named(fmt.Sprintf("%d", config.Config.ChainSelector)),
				i.p2pID.PeerID().String())
		} else {
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
		}
	} else if config.Config.PluginType == uint8(cctypes.PluginTypeCCIPExec) {
		factory = execocr3.NewExecutePluginFactory(
			execocr3.PluginFactoryParams{
				Lggr: i.lggr.
					Named("CCIPExecPlugin").
					Named(destRelayID.String()).
					Named(fmt.Sprintf("%d", config.Config.ChainSelector)).
					Named(offrampAddrStr),
				DonID:            donID,
				OcrConfig:        ccipreaderpkg.OCR3ConfigWithMeta(config),
				ExecCodec:        pluginConfig.ExecutePluginCodec,
				MsgHasher:        pluginConfig.MessageHasher,
				AddrCodec:        i.addressCodec,
				HomeChainReader:  i.homeChainReader,
				TokenDataEncoder: pluginConfig.TokenDataEncoder,
				EstimateProvider: pluginConfig.GasEstimateProvider,
				ContractReaders:  contractReaders,
				ContractWriters:  chainWriters,
			})
		factory = promwrapper.NewReportingPluginFactory(factory, i.lggr, destChainID, "CCIPExec")

		if destChainWriter == nil {
			i.lggr.Infow("no chain writer found for dest chain, creating nil transmitter",
				"destChainID", destChainID,
				"destChainSelector", config.Config.ChainSelector)
			transmitter = ocrimpls.NewNoOpTransmitter(
				i.lggr.
					Named("CCIPExecNoOpTransmitter").
					Named(destRelayID.String()).
					Named(fmt.Sprintf("%d", config.Config.ChainSelector)),
				i.p2pID.PeerID().String())
		} else {
			transmitter = pluginConfig.ContractTransmitterFactory.NewExecTransmitter(
				i.lggr.
					Named("CCIPExecTransmitter").
					Named(destRelayID.String()).
					Named(fmt.Sprintf("%d", config.Config.ChainSelector)),
				destChainWriter,
				ocrtypes.Account(destFromAccounts[0]),
				offrampAddrStr,
			)
		}
	} else {
		return nil, nil, fmt.Errorf("unsupported Plugin type %d", config.Config.PluginType)
	}
	return factory, transmitter, nil
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
	map[cciptypes.ChainSelector]types.ContractWriter,
	error,
) {
	ofc, err := decodeAndValidateOffchainConfig(pluginType, publicCfg)
	if err != nil {
		return nil, nil, err
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

		cr, err1 := crcw.GetChainReader(ctx, ccipcommon.ChainReaderProviderOpts{
			Lggr:          i.lggr,
			Relayer:       relayer,
			ChainID:       chainID,
			DestChainID:   destChainID,
			HomeChainID:   homeChainID,
			Ofc:           ofc,
			ChainSelector: chainSelector,
			ChainFamily:   relayChainFamily,
		})
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

		cw, err1 := crcw.GetChainWriter(ctx, ccipcommon.ChainWriterProviderOpts{
			ChainID:               chainID,
			Relayer:               relayer,
			Transmitters:          i.transmitters,
			ExecBatchGasLimit:     execBatchGasLimit,
			ChainFamily:           relayChainFamily,
			OfframpProgramAddress: config.Config.OfframpAddress,
		})
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
) (ccipcommon.OffChainConfig, error) {
	var ofc ccipcommon.OffChainConfig
	if pluginType == cctypes.PluginTypeCCIPExec {
		execOffchainCfg, err1 := pluginconfig.DecodeExecuteOffchainConfig(publicConfig.ReportingPluginConfig)
		if err1 != nil {
			return ccipcommon.OffChainConfig{}, fmt.Errorf("failed to decode execute offchain config: %w, raw: %s", err1, string(publicConfig.ReportingPluginConfig))
		}
		if err2 := execOffchainCfg.ApplyDefaultsAndValidate(); err2 != nil {
			return ccipcommon.OffChainConfig{}, fmt.Errorf("failed to validate execute offchain config: %w", err2)
		}
		ofc.Execute = &execOffchainCfg
	} else if pluginType == cctypes.PluginTypeCCIPCommit {
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

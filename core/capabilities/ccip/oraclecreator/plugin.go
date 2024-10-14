package oraclecreator

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3confighelper"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipevm"
	evmconfig "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/configs/evm"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ocrimpls"
	cctypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"

	"github.com/smartcontractkit/libocr/commontypes"
	libocr3 "github.com/smartcontractkit/libocr/offchainreporting2plus"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"

	commitocr3 "github.com/smartcontractkit/chainlink-ccip/commit"
	execocr3 "github.com/smartcontractkit/chainlink-ccip/execute"
	"github.com/smartcontractkit/chainlink-ccip/pkg/consts"
	ccipreaderpkg "github.com/smartcontractkit/chainlink-ccip/pkg/reader"
	"github.com/smartcontractkit/chainlink-ccip/pluginconfig"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
	evmrelaytypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
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
	evmConfigs            toml.EVMConfigs
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
	evmConfigs toml.EVMConfigs,
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
		evmConfigs:            evmConfigs,
	}
}

// Type implements types.OracleCreator.
func (i *pluginOracleCreator) Type() cctypes.OracleType {
	return cctypes.OracleTypePlugin
}

// Create implements types.OracleCreator.
func (i *pluginOracleCreator) Create(donID uint32, config cctypes.OCR3ConfigWithMeta) (cctypes.CCIPOracle, error) {
	pluginType := cctypes.PluginType(config.Config.PluginType)

	// Assuming that the chain selector is referring to an evm chain for now.
	// TODO: add an api that returns chain family.
	destChainID, err := chainsel.ChainIdFromSelector(uint64(config.Config.ChainSelector))
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID from selector %d: %w", config.Config.ChainSelector, err)
	}
	destChainFamily := relay.NetworkEVM
	destRelayID := types.NewRelayID(destChainFamily, fmt.Sprintf("%d", destChainID))

	configTracker := ocrimpls.NewConfigTracker(config)
	publicConfig, err := configTracker.PublicConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get public config from OCR config: %w", err)
	}

	contractReaders, chainWriters, err := i.createReadersAndWriters(
		destChainID,
		pluginType,
		config,
		publicConfig,
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
	factory, transmitter, err := i.createFactoryAndTransmitter(donID, config, destRelayID, contractReaders, chainWriters, destChainWriter, destFromAccounts)
	if err != nil {
		return nil, fmt.Errorf("failed to create factory and transmitter: %w", err)
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
				Named(hexutil.Encode(config.Config.OfframpAddress)),
			false,
			func(ctx context.Context, msg string) {}),
		MetricsRegisterer: prometheus.WrapRegistererWith(map[string]string{"name": fmt.Sprintf("commit-%d", config.Config.ChainSelector)}, prometheus.DefaultRegisterer),
		MonitoringEndpoint: i.monitoringEndpointGen.GenMonitoringEndpoint(
			destChainFamily,
			destRelayID.ChainID,
			string(config.Config.OfframpAddress),
			synchronization.OCR3CCIPCommit,
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
	chainWriters map[cciptypes.ChainSelector]types.ChainWriter,
	destChainWriter types.ChainWriter,
	destFromAccounts []string,
) (ocr3types.ReportingPluginFactory[[]byte], ocr3types.ContractTransmitter[[]byte], error) {
	var factory ocr3types.ReportingPluginFactory[[]byte]
	var transmitter ocr3types.ContractTransmitter[[]byte]
	if config.Config.PluginType == uint8(cctypes.PluginTypeCCIPCommit) {
		factory = commitocr3.NewPluginFactory(
			i.lggr.
				Named("CCIPCommitPlugin").
				Named(destRelayID.String()).
				Named(fmt.Sprintf("%d", config.Config.ChainSelector)).
				Named(hexutil.Encode(config.Config.OfframpAddress)),
			donID,
			ccipreaderpkg.OCR3ConfigWithMeta(config),
			ccipevm.NewCommitPluginCodecV1(),
			ccipevm.NewMessageHasherV1(),
			i.homeChainReader,
			i.homeChainSelector,
			contractReaders,
			chainWriters,
		)
		transmitter = ocrimpls.NewCommitContractTransmitter[[]byte](destChainWriter,
			ocrtypes.Account(destFromAccounts[0]),
			hexutil.Encode(config.Config.OfframpAddress), // TODO: this works for evm only, how about non-evm?
		)
	} else if config.Config.PluginType == uint8(cctypes.PluginTypeCCIPExec) {
		factory = execocr3.NewPluginFactory(
			i.lggr.
				Named("CCIPExecPlugin").
				Named(destRelayID.String()).
				Named(hexutil.Encode(config.Config.OfframpAddress)),
			donID,
			ccipreaderpkg.OCR3ConfigWithMeta(config),
			ccipevm.NewExecutePluginCodecV1(),
			ccipevm.NewMessageHasherV1(),
			i.homeChainReader,
			ccipevm.NewEVMTokenDataEncoder(),
			ccipevm.NewGasEstimateProvider(),
			contractReaders,
			chainWriters,
		)
		transmitter = ocrimpls.NewExecContractTransmitter[[]byte](destChainWriter,
			ocrtypes.Account(destFromAccounts[0]),
			hexutil.Encode(config.Config.OfframpAddress), // TODO: this works for evm only, how about non-evm?
		)
	} else {
		return nil, nil, fmt.Errorf("unsupported plugin type %d", config.Config.PluginType)
	}
	return factory, transmitter, nil
}

func (i *pluginOracleCreator) createReadersAndWriters(
	destChainID uint64,
	pluginType cctypes.PluginType,
	config cctypes.OCR3ConfigWithMeta,
	publicCfg ocr3confighelper.PublicConfig,
) (
	map[cciptypes.ChainSelector]types.ContractReader,
	map[cciptypes.ChainSelector]types.ChainWriter,
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

	homeChainID, err := i.getChainID(i.homeChainSelector)
	if err != nil {
		return nil, nil, err
	}

	contractReaders := make(map[cciptypes.ChainSelector]types.ContractReader)
	chainWriters := make(map[cciptypes.ChainSelector]types.ChainWriter)
	for relayID, relayer := range i.relayers {
		chainID, ok := new(big.Int).SetString(relayID.ChainID, 10)
		if !ok {
			return nil, nil, fmt.Errorf("error parsing chain ID, expected big int: %s", relayID.ChainID)
		}

		chainSelector, err1 := i.getChainSelector(chainID.Uint64())
		if err1 != nil {
			return nil, nil, fmt.Errorf("failed to get chain selector from chain ID %s: %w", chainID.String(), err1)
		}

		chainReaderConfig, err1 := getChainReaderConfig(chainID.Uint64(), destChainID, homeChainID, ofc, chainSelector)
		if err1 != nil {
			return nil, nil, fmt.Errorf("failed to get chain reader config: %w", err1)
		}

		// TODO: context.
		cr, err1 := relayer.NewContractReader(context.Background(), chainReaderConfig)
		if err1 != nil {
			return nil, nil, err1
		}

		if chainID.Uint64() == destChainID {
			offrampAddressHex := common.BytesToAddress(config.Config.OfframpAddress).Hex()
			err2 := cr.Bind(context.Background(), []types.BoundContract{
				{
					Address: offrampAddressHex,
					Name:    consts.ContractNameOffRamp,
				},
			})
			if err2 != nil {
				return nil, nil, fmt.Errorf("failed to bind chain reader for dest chain %s's offramp at %s: %w", chainID.String(), offrampAddressHex, err)
			}
		}

		if err2 := cr.Start(context.Background()); err2 != nil {
			return nil, nil, fmt.Errorf("failed to start contract reader for chain %s: %w", chainID.String(), err2)
		}

		cw, err1 := createChainWriter(
			chainID,
			i.evmConfigs,
			relayer,
			i.transmitters,
			execBatchGasLimit)
		if err1 != nil {
			return nil, nil, err1
		}

		if err4 := cw.Start(context.Background()); err4 != nil {
			return nil, nil, fmt.Errorf("failed to start chain writer for chain %s: %w", chainID.String(), err4)
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
		return offChainConfig{}, fmt.Errorf("invalid offchain config: both commit and exec configs are either set or unset")
	}
	return ofc, nil
}

func (i *pluginOracleCreator) getChainSelector(chainID uint64) (cciptypes.ChainSelector, error) {
	chainSelector, ok := chainsel.EvmChainIdToChainSelector()[chainID]
	if !ok {
		return 0, fmt.Errorf("failed to get chain selector from chain ID %d", chainID)
	}
	return cciptypes.ChainSelector(chainSelector), nil
}

func (i *pluginOracleCreator) getChainID(chainSelector cciptypes.ChainSelector) (uint64, error) {
	chainID, err := chainsel.ChainIdFromSelector(uint64(chainSelector))
	if err != nil {
		return 0, fmt.Errorf("failed to get chain ID from chain selector %d: %w", chainSelector, err)
	}
	return chainID, nil
}

func getChainReaderConfig(
	chainID uint64,
	destChainID uint64,
	homeChainID uint64,
	ofc offChainConfig,
	chainSelector cciptypes.ChainSelector,
) ([]byte, error) {
	var chainReaderConfig evmrelaytypes.ChainReaderConfig
	if chainID == destChainID {
		chainReaderConfig = evmconfig.DestReaderConfig
	} else {
		chainReaderConfig = evmconfig.SourceReaderConfig
	}

	if !ofc.commitEmpty() && ofc.commit().PriceFeedChainSelector == chainSelector {
		chainReaderConfig = evmconfig.MergeReaderConfigs(chainReaderConfig, evmconfig.FeedReaderConfig)
	}

	if isUSDCEnabled(chainID, destChainID, ofc) {
		chainReaderConfig = evmconfig.MergeReaderConfigs(chainReaderConfig, evmconfig.USDCReaderConfig)
	}

	if chainID == homeChainID {
		chainReaderConfig = evmconfig.MergeReaderConfigs(chainReaderConfig, evmconfig.HomeChainReaderConfigRaw)
	}

	marshaledConfig, err := json.Marshal(chainReaderConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal chain reader config: %w", err)
	}

	return marshaledConfig, nil
}

func isUSDCEnabled(chainID uint64, destChainID uint64, ofc offChainConfig) bool {
	if chainID == destChainID {
		return false
	}

	if ofc.execEmpty() {
		return false
	}

	return ofc.exec().IsUSDCEnabled()
}

func createChainWriter(
	chainID *big.Int,
	evmConfigs toml.EVMConfigs,
	relayer loop.Relayer,
	transmitters map[types.RelayID][]string,
	execBatchGasLimit uint64,
) (types.ChainWriter, error) {
	var fromAddress common.Address
	transmitter, ok := transmitters[types.NewRelayID(relay.NetworkEVM, chainID.String())]
	if ok {
		// TODO: remove EVM-specific stuff
		fromAddress = common.HexToAddress(transmitter[0])
	}

	maxGasPrice := getKeySpecificMaxGasPrice(evmConfigs, chainID, fromAddress)
	if maxGasPrice == nil {
		return nil, fmt.Errorf("failed to find max gas price for chain %s", chainID.String())
	}

	chainWriterRawConfig, err := evmconfig.ChainWriterConfigRaw(
		fromAddress,
		maxGasPrice,
		defaultCommitGasLimit,
		execBatchGasLimit,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create chain writer config: %w", err)
	}

	chainWriterConfig, err := json.Marshal(chainWriterRawConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal chain writer config: %w", err)
	}

	// TODO: context.
	cw, err := relayer.NewChainWriter(context.Background(), chainWriterConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create chain writer for chain %s: %w", chainID.String(), err)
	}

	return cw, nil
}

func getKeySpecificMaxGasPrice(evmConfigs toml.EVMConfigs, chainID *big.Int, fromAddress common.Address) *assets.Wei {
	var maxGasPrice *assets.Wei

	// If a chain is enabled it should have some configuration in the TOML config
	// of the chainlink node.
	for _, config := range evmConfigs {
		if config.ChainID.ToInt().Cmp(chainID) != 0 {
			continue
		}

		// find the key-specific max gas price for the given fromAddress.
		for _, keySpecific := range config.KeySpecific {
			if keySpecific.Key.Address() == fromAddress {
				maxGasPrice = keySpecific.GasEstimator.PriceMax
			}
		}

		// if we didn't find a key-specific max gas price, use the one specified
		// in the gas estimator config, which should have a default value.
		if maxGasPrice == nil {
			maxGasPrice = config.GasEstimator.PriceMax
		}
	}

	return maxGasPrice
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

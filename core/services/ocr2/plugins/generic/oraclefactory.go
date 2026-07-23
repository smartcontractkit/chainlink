package generic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	ocr "github.com/smartcontractkit/libocr/offchainreporting2plus"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3shims"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/ocr2key"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr/capregconfig"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/telemetry"
)

type oracleFactory struct {
	database               ocr3types.Database
	jobID                  int32
	jobName                string
	jobORM                 job.ORM
	kb                     ocr2key.KeyBundle
	lggr                   logger.Logger
	config                 job.OracleFactoryConfig
	onchainSigningStrategy job.OnchainSigningStrategy
	peerWrapper            *ocrcommon.SingletonPeerWrapper
	relayerSet             *RelayerSet
	ocrKeystore            keystore.OCR2
	ethKeystore            keystore.Eth
	ocrConfigService       capregconfig.OCRConfigService
	capabilityID           string // Capability ID for registry-based config lookup
}

type OracleFactoryParams struct {
	JobID                  int32
	JobName                string
	JobORM                 job.ORM
	KB                     ocr2key.KeyBundle
	Logger                 logger.Logger
	Config                 job.OracleFactoryConfig
	OnchainSigningStrategy job.OnchainSigningStrategy
	PeerWrapper            *ocrcommon.SingletonPeerWrapper
	RelayerSet             *RelayerSet
	OcrKeystore            keystore.OCR2
	EthKeystore            keystore.Eth
	// OCRConfigService provides OCR config from the capabilities registry.
	// When set, the factory will use dynamic tracker/digester that can switch
	// between registry-based and legacy contract-based config.
	OCRConfigService capregconfig.OCRConfigService
	CapabilityID     string
}

func NewOracleFactory(params OracleFactoryParams) (core.OracleFactory, error) {
	return &oracleFactory{
		database:               OracleFactoryDB(params.JobID, params.Logger),
		jobID:                  params.JobID,
		jobName:                params.JobName,
		jobORM:                 params.JobORM,
		kb:                     params.KB,
		lggr:                   params.Logger,
		config:                 params.Config,
		onchainSigningStrategy: params.OnchainSigningStrategy,
		peerWrapper:            params.PeerWrapper,
		relayerSet:             params.RelayerSet,
		ocrKeystore:            params.OcrKeystore,
		ethKeystore:            params.EthKeystore,
		ocrConfigService:       params.OCRConfigService,
		capabilityID:           params.CapabilityID,
	}, nil
}

func AdjustLocalConfigForRegistryBasedConfig(lc ocrtypes.LocalConfig) ocrtypes.LocalConfig {
	// block confirmations are irrelevant when using registry-based config
	// this also works with legacy config contracts, simply doesn't wait for extra confirmations
	lc.SkipContractConfigConfirmations = true
	// poll frequently to react to config changes quickly
	lc.ContractConfigTrackerPollInterval = 5 * time.Second
	return lc
}

func (of *oracleFactory) NewOracle(ctx context.Context, args core.OracleArgs) (core.Oracle, error) {
	of.lggr.Debugw("Creating new oracle from oracle factory using config", "config", of.config, "capabilityID", of.capabilityID)

	if !of.peerWrapper.IsStarted() {
		return nil, errors.New("peer wrapper not started")
	}

	relayerSetRelayer, err := of.relayerSet.Get(ctx, types.RelayID{Network: "evm", ChainID: of.config.ChainID})
	if err != nil {
		return nil, fmt.Errorf("error when getting relayer: %w", err)
	}

	// TODO - to avoid this cast requires https://smartcontract-it.atlassian.net/browse/CAPPL-1001
	relayer, ok := relayerSetRelayer.(relayerWrapper)
	if !ok {
		return nil, fmt.Errorf("expected relayer to be of type relayerWrapper, got %T", relayer)
	}

	var relayConfig = struct {
		ChainID                string   `json:"chainID"`
		EffectiveTransmitterID string   `json:"effectiveTransmitterID"`
		SendingKeys            []string `json:"sendingKeys"`
	}{
		ChainID:                of.config.ChainID,
		EffectiveTransmitterID: of.config.TransmitterID,
		SendingKeys:            []string{of.config.TransmitterID},
	}
	relayConfigBytes, err := json.Marshal(relayConfig)
	if err != nil {
		return nil, fmt.Errorf("error when marshalling relay config: %w", err)
	}

	legacyConfigProvider, err := relayer.NewConfigProvider(ctx, core.RelayArgs{
		ContractID:   of.config.OCRContractAddress,
		ProviderType: string(types.OCR3Capability),
		RelayConfig:  relayConfigBytes,
	})
	if err != nil {
		return nil, fmt.Errorf("error when getting config provider: %w", err)
	}

	// Determine config tracker and digester to use
	var configTracker ocrtypes.ContractConfigTracker
	var configDigester ocrtypes.OffchainConfigDigester

	if of.ocrConfigService != nil && of.capabilityID != "" {
		// Wrap with dynamic tracker/digester from OCRConfigService (with fallback).
		// NOTE: Standard Capabilities currently support only one OCR instance so we're using OCR3ConfigDefaultKey.
		configTracker, err = of.ocrConfigService.GetConfigTracker(
			of.capabilityID, capabilitiespb.OCR3ConfigDefaultKey, legacyConfigProvider.ContractConfigTracker())
		if err != nil {
			return nil, fmt.Errorf("failed to get config tracker: %w", err)
		}

		configDigester, err = of.ocrConfigService.GetConfigDigester(
			of.capabilityID, capabilitiespb.OCR3ConfigDefaultKey, legacyConfigProvider.OffchainConfigDigester())
		if err != nil {
			return nil, fmt.Errorf("failed to get config digester: %w", err)
		}

		of.lggr.Infow("Using dynamic OCR config service", "capabilityID", of.capabilityID)
	} else {
		configTracker = legacyConfigProvider.ContractConfigTracker()
		configDigester = legacyConfigProvider.OffchainConfigDigester()
	}

	bootstrapPeers, err := ocrcommon.ParseBootstrapPeers(of.config.BootstrapPeers)
	if err != nil {
		return nil, fmt.Errorf("failed to parse bootstrap peers: %w", err)
	}

	keyBundles := map[string]ocr2key.KeyBundle{}
	for name, kbID := range of.onchainSigningStrategy.Config {
		os, ostErr := of.ocrKeystore.Get(kbID)
		if ostErr != nil {
			return nil, fmt.Errorf("failed to get ocr key for key bundle ID '%s': %w", kbID, ostErr)
		}
		keyBundles[name] = os
	}
	onchainKeyringAdapter, err := ocrcommon.NewOCR3OnchainKeyringMultiChainAdapter(keyBundles, of.lggr)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate onchain keyring with multi chain adapter: %w", err)
	}

	oracle, err := ocr.NewOracle(ocr.OCR3OracleArgs2[[]byte]{
		ContractConfigTracker:        configTracker,
		OffchainConfigDigester:       configDigester,
		LocalConfig:                  AdjustLocalConfigForRegistryBasedConfig(args.LocalConfig),
		ContractTransmitter:          NewContractTransmitter(of.config.TransmitterID, args.ContractTransmitter),
		ReportingPluginFactory:       args.ReportingPluginFactoryService,
		BinaryNetworkEndpointFactory: of.peerWrapper.Peer2,
		V2Bootstrappers:              bootstrapPeers,
		Database:                     of.database,
		Logger: ocrcommon.NewOCRWrapper(of.lggr, true, func(ctx context.Context, msg string) {
			of.lggr.Error("OCRWrapperOracleError:" + msg)
		}),
		MonitoringEndpoint: &telemetry.NoopAgent{},
		OffchainKeyring:    of.kb,
		OnchainKeyring:     ocr3shims.OnchainKeyringAsOnchainKeyring2(onchainKeyringAdapter),
		MetricsRegisterer:  prometheus.WrapRegistererWith(map[string]string{"job_name": of.jobName}, prometheus.DefaultRegisterer),
	})

	if err != nil {
		return nil, fmt.Errorf("%w: failed to create new OCR oracle", err)
	}

	of.lggr.Debug("Created new oracle from oracle factory")

	return &adaptedOracle{oracle: oracle}, nil
}

type adaptedOracle struct {
	oracle ocr.Oracle
}

func (a *adaptedOracle) Start(ctx context.Context) error {
	return a.oracle.Start()
}

func (a *adaptedOracle) Close(ctx context.Context) error {
	return a.oracle.Close()
}

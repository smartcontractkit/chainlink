package generic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"

	ocr "github.com/smartcontractkit/libocr/offchainreporting2plus"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/telemetry"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
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
	}, nil
}

func (of *oracleFactory) NewOracle(ctx context.Context, args core.OracleArgs) (core.Oracle, error) {

	of.lggr.Debugf("Creating new oracle from oracle factory using config: %+v", of.config)

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

	configProvider, err := relayer.NewConfigProvider(ctx, core.RelayArgs{
		ContractID:   of.config.OCRContractAddress,
		ProviderType: string(types.OCR3Capability),
		RelayConfig:  relayConfigBytes,
	})
	if err != nil {
		return nil, fmt.Errorf("error when getting config provider: %w", err)
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

	oracle, err := ocr.NewOracle(ocr.OCR3OracleArgs[[]byte]{
		// We are relying on the relayer plugin provider for the offchain config digester
		// and the contract config tracker to save time.
		ContractConfigTracker:        configProvider.ContractConfigTracker(),
		OffchainConfigDigester:       configProvider.OffchainConfigDigester(),
		LocalConfig:                  args.LocalConfig,
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
		OnchainKeyring:     onchainKeyringAdapter,
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

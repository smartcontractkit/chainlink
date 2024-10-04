package generic

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/smartcontractkit/libocr/commontypes"
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

type oracleFactoryConfig struct {
	Enabled        bool
	TraceLogging   bool
	BootstrapPeers []commontypes.BootstrapperLocator
}

func NewOracleFactoryConfig(config string) (*oracleFactoryConfig, error) {
	var ofc struct {
		Enabled        bool     `json:"enabled"`
		TraceLogging   bool     `json:"traceLogging"`
		BootstrapPeers []string `json:"bootstrapPeers"`
	}
	err := json.Unmarshal([]byte(config), &ofc)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal oracle factory config")
	}

	if !ofc.Enabled {
		return &oracleFactoryConfig{}, nil
	}

	// If Oracle Factory is enabled, it must have at least one bootstrap peer
	if len(ofc.BootstrapPeers) == 0 {
		return nil, errors.New("no bootstrap peers found")
	}

	bootstrapPeers, err := ocrcommon.ParseBootstrapPeers(ofc.BootstrapPeers)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse bootstrap peers")
	}

	return &oracleFactoryConfig{
		Enabled:        true,
		TraceLogging:   ofc.TraceLogging,
		BootstrapPeers: bootstrapPeers,
	}, nil
}

type OracleIdentity struct {
	EVMKey                    string   `json:"evm_key"`
	PeerID                    string   `json:"peer_id"`
	PublicKey                 []byte   `json:"public_key"`
	OffchainPublicKey         [32]byte `json:"offchain_public_key"`
	ConfigEncryptionPublicKey [32]byte `json:"config_encryption_public_key"`
}

type oracleFactory struct {
	database    ocr3types.Database
	jobID       int32
	jobName     string
	jobORM      job.ORM
	kb          ocr2key.KeyBundle
	lggr        logger.Logger
	config      *oracleFactoryConfig
	peerWrapper *ocrcommon.SingletonPeerWrapper
	relayerSet  *RelayerSet
	identity    OracleIdentity
}

type OracleFactoryParams struct {
	JobID       int32
	JobName     string
	JobORM      job.ORM
	Kb          ocr2key.KeyBundle
	Logger      logger.Logger
	Config      *oracleFactoryConfig
	PeerWrapper *ocrcommon.SingletonPeerWrapper
	RelayerSet  *RelayerSet
	Identity    OracleIdentity
}

func NewOracleFactory(params OracleFactoryParams) (core.OracleFactory, error) {
	return &oracleFactory{
		database:    NewMemoryDB(params.JobID, params.Logger),
		jobID:       params.JobID,
		jobName:     params.JobName,
		jobORM:      params.JobORM,
		kb:          params.Kb,
		lggr:        params.Logger,
		config:      params.Config,
		peerWrapper: params.PeerWrapper,
		relayerSet:  params.RelayerSet,
		identity:    params.Identity,
	}, nil
}

type JSONConfig map[string]interface{}

// Bytes returns the raw bytes
func (r JSONConfig) Bytes() []byte {
	b, _ := json.Marshal(r)
	return b
}

func (of *oracleFactory) NewOracle(ctx context.Context, args core.OracleArgs) (core.Oracle, error) {
	if !of.peerWrapper.IsStarted() {
		return nil, errors.New("peer wrapper not started")
	}

	of.lggr.Debug("oracleIdentity: ", of.identity)

	relayer, err := of.relayerSet.Get(ctx, types.RelayID{Network: "evm", ChainID: "31337"})
	if err != nil {
		return nil, fmt.Errorf("error when getting relayer: %w", err)
	}

	type RelayConfig struct {
		ChainID                string   `json:"chainID"`
		EffectiveTransmitterID string   `json:"effectiveTransmitterID"`
		SendingKeys            []string `json:"sendingKeys"`
	}

	var relayConfig = RelayConfig{
		ChainID:                "31337",
		EffectiveTransmitterID: of.identity.EVMKey,
		SendingKeys:            []string{of.identity.EVMKey},
	}
	relayConfigBytes, err := json.Marshal(relayConfig)
	if err != nil {
		return nil, fmt.Errorf("error when marshalling relay config: %w", err)
	}

	pluginProvider, err := relayer.NewPluginProvider(ctx, core.RelayArgs{
		ContractID:   "0x2279B7A0a67DB372996a5FaB50D91eAA73d2eBe6", // TODO: Oracle Factory config
		ProviderType: "plugin",
		RelayConfig:  relayConfigBytes,
	}, core.PluginArgs{
		TransmitterID: of.identity.EVMKey,
		PluginConfig: JSONConfig{
			"pluginName": "kvstore-capability",
			"OCRVersion": 3,
		}.Bytes(),
	})
	if err != nil {
		return nil, fmt.Errorf("error when getting offchain digester: %w", err)
	}

	oracle, err := ocr.NewOracle(ocr.OCR3OracleArgs[[]byte]{
		LocalConfig:                  args.LocalConfig,
		ContractConfigTracker:        pluginProvider.ContractConfigTracker(),
		ContractTransmitter:          args.ContractTransmitter,
		OffchainConfigDigester:       pluginProvider.OffchainConfigDigester(),
		ReportingPluginFactory:       args.ReportingPluginFactoryService,
		BinaryNetworkEndpointFactory: of.peerWrapper.Peer2,
		V2Bootstrappers:              of.config.BootstrapPeers,
		Database:                     of.database,
		Logger: ocrcommon.NewOCRWrapper(of.lggr, of.config.TraceLogging, func(ctx context.Context, msg string) {
			logger.Sugared(of.lggr).ErrorIf(of.jobORM.RecordError(ctx, of.jobID, msg), "unable to record error")
		}),
		// TODO?
		MonitoringEndpoint: &telemetry.NoopAgent{},
		OffchainKeyring:    of.kb,
		OnchainKeyring:     ocrcommon.NewOCR3OnchainKeyringAdapter(of.kb),
		MetricsRegisterer:  prometheus.WrapRegistererWith(map[string]string{"job_name": of.jobName}, prometheus.DefaultRegisterer),
	})

	if err != nil {
		return nil, fmt.Errorf("%w: failed to create new OCR oracle", err)
	}

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

package generic

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
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
	config      job.OracleFactoryConfig
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
	Config      job.OracleFactoryConfig
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

	relayer, err := of.relayerSet.Get(ctx, types.RelayID{Network: of.config.Network, ChainID: of.config.ChainID})
	if err != nil {
		return nil, fmt.Errorf("error when getting relayer: %w", err)
	}

	var relayConfig = struct {
		ChainID                string   `json:"chainID"`
		EffectiveTransmitterID string   `json:"effectiveTransmitterID"`
		SendingKeys            []string `json:"sendingKeys"`
	}{
		ChainID:                of.config.ChainID,
		EffectiveTransmitterID: of.identity.EVMKey,
		SendingKeys:            []string{of.identity.EVMKey},
	}
	relayConfigBytes, err := json.Marshal(relayConfig)
	if err != nil {
		return nil, fmt.Errorf("error when marshalling relay config: %w", err)
	}

	pluginProvider, err := relayer.NewPluginProvider(ctx, core.RelayArgs{
		ContractID:   of.config.OCRContractAddress,
		ProviderType: "plugin",
		RelayConfig:  relayConfigBytes,
	}, core.PluginArgs{
		TransmitterID: of.identity.EVMKey,
	})
	if err != nil {
		return nil, fmt.Errorf("error when getting offchain digester: %w", err)
	}

	bootstrapPeers, err := ocrcommon.ParseBootstrapPeers(of.config.BootstrapPeers)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse bootstrap peers")
	}

	oracle, err := ocr.NewOracle(ocr.OCR3OracleArgs[[]byte]{
		ContractConfigTracker:        pluginProvider.ContractConfigTracker(),
		OffchainConfigDigester:       pluginProvider.OffchainConfigDigester(),
		LocalConfig:                  args.LocalConfig,
		ContractTransmitter:          args.ContractTransmitter,
		ReportingPluginFactory:       args.ReportingPluginFactoryService,
		BinaryNetworkEndpointFactory: of.peerWrapper.Peer2,
		V2Bootstrappers:              bootstrapPeers,
		Database:                     of.database,
		Logger: ocrcommon.NewOCRWrapper(of.lggr, true, func(ctx context.Context, msg string) {
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

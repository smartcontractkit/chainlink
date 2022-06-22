package ocr2

import (
	"github.com/pkg/errors"
	libocr2 "github.com/smartcontractkit/libocr/offchainreporting2"
	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink-relay/pkg/types"

	"github.com/smartcontractkit/chainlink/core/chains/evm"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/ocr2/plugins"
	"github.com/smartcontractkit/chainlink/core/services/ocr2/plugins/dkg"
	"github.com/smartcontractkit/chainlink/core/services/ocr2/plugins/median"
	"github.com/smartcontractkit/chainlink/core/services/ocr2/plugins/vrf"
	"github.com/smartcontractkit/chainlink/core/services/ocr2/validate"
	"github.com/smartcontractkit/chainlink/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/relay"
	evmrelay "github.com/smartcontractkit/chainlink/core/services/relay/evm"
	"github.com/smartcontractkit/chainlink/core/services/telemetry"
)

type Delegate struct {
	db                    *sqlx.DB
	jobORM                job.ORM
	pipelineRunner        pipeline.Runner
	peerWrapper           *ocrcommon.SingletonPeerWrapper
	monitoringEndpointGen telemetry.MonitoringEndpointGenerator
	chainSet              evm.ChainSet
	cfg                   validate.Config
	lggr                  logger.Logger
	ks                    keystore.OCR2
	relayers              map[relay.Network]types.Relayer
}

var _ job.Delegate = (*Delegate)(nil)

func NewDelegate(
	db *sqlx.DB,
	jobORM job.ORM,
	pipelineRunner pipeline.Runner,
	peerWrapper *ocrcommon.SingletonPeerWrapper,
	monitoringEndpointGen telemetry.MonitoringEndpointGenerator,
	chainSet evm.ChainSet,
	lggr logger.Logger,
	cfg validate.Config,
	ks keystore.OCR2,
	relayers map[relay.Network]types.Relayer,
) *Delegate {
	return &Delegate{
		db,
		jobORM,
		pipelineRunner,
		peerWrapper,
		monitoringEndpointGen,
		chainSet,
		cfg,
		lggr,
		ks,
		relayers,
	}
}

func (d Delegate) JobType() job.Type {
	return job.OffchainReporting2
}

func (Delegate) OnJobCreated(spec job.Job) {}
func (Delegate) OnJobDeleted(spec job.Job) {}

func (Delegate) AfterJobCreated(spec job.Job)  {}
func (Delegate) BeforeJobDeleted(spec job.Job) {}

// ServicesForSpec returns the OCR2 services that need to run for this job
func (d Delegate) ServicesForSpec(jobSpec job.Job) ([]job.ServiceCtx, error) {
	spec := jobSpec.OCR2OracleSpec
	if spec == nil {
		return nil, errors.Errorf("offchainreporting2.Delegate expects an *job.Offchainreporting2OracleSpec to be present, got %v", jobSpec)
	}
	if !spec.TransmitterID.Valid {
		return nil, errors.Errorf("expected a transmitterID to be specified")
	}
	relayer, exists := d.relayers[spec.Relay]
	if !exists {
		return nil, errors.Errorf("%s relay does not exist is it enabled?", spec.Relay)
	}

	ocrDB := NewDB(d.db, spec.ID, d.lggr, d.cfg)
	peerWrapper := d.peerWrapper
	if peerWrapper == nil {
		return nil, errors.New("cannot setup OCR2 job service, libp2p peer was missing")
	} else if !peerWrapper.IsStarted() {
		return nil, errors.New("peerWrapper is not started. OCR2 jobs require a started and running peer. Did you forget to specify P2P_LISTEN_PORT?")
	}

	lggr := d.lggr.Named("OCR").With(
		"contractID", spec.ContractID,
		"jobName", jobSpec.Name.ValueOrZero(),
		"jobID", jobSpec.ID,
	)
	ocrLogger := logger.NewOCRWrapper(lggr, true, func(msg string) {
		d.lggr.ErrorIf(d.jobORM.RecordError(jobSpec.ID, msg), "unable to record error")
	})

	lc := validate.ToLocalConfig(d.cfg, *spec)
	if err := libocr2.SanityCheckLocalConfig(lc); err != nil {
		return nil, err
	}
	d.lggr.Infow("OCR2 job using local config",
		"BlockchainTimeout", lc.BlockchainTimeout,
		"ContractConfigConfirmations", lc.ContractConfigConfirmations,
		"ContractConfigTrackerPollInterval", lc.ContractConfigTrackerPollInterval,
		"ContractTransmitterTransmitTimeout", lc.ContractTransmitterTransmitTimeout,
		"DatabaseTimeout", lc.DatabaseTimeout,
	)

	bootstrapPeers, err := ocrcommon.GetValidatedBootstrapPeers(spec.P2PV2Bootstrappers, peerWrapper.Config().P2PV2Bootstrappers())
	if err != nil {
		return nil, err
	}
	d.lggr.Debugw("Using bootstrap peers", "peers", bootstrapPeers)
	// Fetch the specified OCR2 key bundle
	var kbID string
	if spec.OCRKeyBundleID.Valid {
		kbID = spec.OCRKeyBundleID.String
	} else if kbID, err = d.cfg.OCR2KeyBundleID(); err != nil {
		return nil, err
	}
	kb, err := d.ks.Get(kbID)
	if err != nil {
		return nil, err
	}
	runResults := make(chan pipeline.Run, d.cfg.JobPipelineResultWriteQueueDepth())

	// These are populated here because when the pipeline spec is
	// run it uses them to create identifiable prometheus metrics.
	// TODO SC-30421 Move pipeline population to job spawner
	jobSpec.PipelineSpec.JobName = jobSpec.Name.ValueOrZero()
	jobSpec.PipelineSpec.JobID = jobSpec.ID

	var pluginOracle plugins.OraclePlugin
	var ocr2Provider types.Plugin
	switch spec.PluginType {
	case job.Median:
		medianProvider, err2 := relayer.NewMedianProvider(
			types.RelayArgs{
				ExternalJobID: jobSpec.ExternalJobID,
				JobID:         spec.ID,
				ContractID:    spec.ContractID,
				RelayConfig:   spec.RelayConfig.Bytes(),
			}, types.PluginArgs{
				TransmitterID: spec.TransmitterID.String,
				PluginConfig:  spec.PluginConfig.Bytes(),
			})
		if err2 != nil {
			return nil, err2
		}
		ocr2Provider = medianProvider
		pluginOracle, err = median.NewMedian(jobSpec, medianProvider, d.pipelineRunner, runResults, lggr, ocrLogger)
	case job.DKG:
		dkgProvider, err2 := evmrelay.NewOCR2VRFRelayer(relayer).NewDKGProvider(
			types.RelayArgs{
				ExternalJobID: jobSpec.ExternalJobID,
				JobID:         spec.ID,
				ContractID:    spec.ContractID,
				RelayConfig:   spec.RelayConfig.Bytes(),
			}, spec.TransmitterID.String)
		if err2 != nil {
			return nil, err2
		}
		ocr2Provider = dkgProvider
		pluginOracle, err = dkg.NewDKG(lggr)
	case job.OCR2VRF:
		vrfProvider, err2 := evmrelay.NewOCR2VRFRelayer(relayer).NewOCR2VRFProvider(
			types.RelayArgs{
				ExternalJobID: jobSpec.ExternalJobID,
				JobID:         spec.ID,
				ContractID:    spec.ContractID,
				RelayConfig:   spec.RelayConfig.Bytes(),
			}, spec.TransmitterID.String)
		if err2 != nil {
			return nil, err2
		}
		ocr2Provider = vrfProvider
		pluginOracle, err = vrf.NewVRF(lggr)
	default:
		return nil, errors.Errorf("plugin type %s not supported", spec.PluginType)
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialise plugin")
	}
	pluginFactory, err := pluginOracle.GetPluginFactory()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get plugin factory")
	}
	pluginServices, err := pluginOracle.GetServices()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get plugin services")
	}

	oracle, err := libocr2.NewOracle(libocr2.OracleArgs{
		BinaryNetworkEndpointFactory: peerWrapper.Peer2,
		V2Bootstrappers:              bootstrapPeers,
		ContractTransmitter:          ocr2Provider.ContractTransmitter(),
		ContractConfigTracker:        ocr2Provider.ContractConfigTracker(),
		Database:                     ocrDB,
		LocalConfig:                  lc,
		Logger:                       ocrLogger,
		MonitoringEndpoint:           d.monitoringEndpointGen.GenMonitoringEndpoint(spec.ContractID),
		OffchainConfigDigester:       ocr2Provider.OffchainConfigDigester(),
		OffchainKeyring:              kb,
		OnchainKeyring:               kb,
		ReportingPluginFactory:       pluginFactory,
	})
	if err != nil {
		return nil, errors.Wrap(err, "error calling NewOracle")
	}

	// RunResultSaver needs to be started first, so it's available
	// to read odb writes. It is stopped last after the OraclePlugin is shut down
	// so no further runs are enqueued, and we can drain the queue.
	runResultSaver := ocrcommon.NewResultRunSaver(
		runResults,
		d.pipelineRunner,
		make(chan struct{}),
		lggr)

	oracleCtx := job.NewServiceAdapter(oracle)
	return append([]job.ServiceCtx{runResultSaver, ocr2Provider, oracleCtx}, pluginServices...), nil
}

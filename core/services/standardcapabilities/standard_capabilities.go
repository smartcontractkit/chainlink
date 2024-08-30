package standardcapabilities

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/plugins"
)

type standardCapabilities struct {
	services.StateMachine
	log                  logger.Logger
	spec                 *job.StandardCapabilitiesSpec
	pluginRegistrar      plugins.RegistrarConfig
	telemetryService     core.TelemetryService
	store                core.KeyValueStore
	CapabilitiesRegistry core.CapabilitiesRegistry
	errorLog             core.ErrorLog
	pipelineRunner       core.PipelineRunnerService
	relayerSet           core.RelayerSet
	oracleFactory        core.OracleFactory

	capabilitiesLoop *loop.StandardCapabilitiesService
}

func newStandardCapabilities(
	log logger.Logger,
	spec *job.StandardCapabilitiesSpec,
	pluginRegistrar plugins.RegistrarConfig,
	telemetryService core.TelemetryService,
	store core.KeyValueStore,
	CapabilitiesRegistry core.CapabilitiesRegistry,
	errorLog core.ErrorLog,
	pipelineRunner core.PipelineRunnerService,
	relayerSet core.RelayerSet,
	oracleFactory core.OracleFactory,
) *standardCapabilities {
	/*
		GENERIC DELEGATE
		oracleArgs := libocr2.OCR3OracleArgs[[]byte]{
			// QUESTION: Can we abstract this away from standard capability dev?
			BinaryNetworkEndpointFactory: d.peerWrapper.Peer2,
			// Where do we get this from? We need a set of global bootstrappers that
			// allow any nodes from the capabilities registry to connect to each other.
			V2Bootstrappers:              bootstrapPeers,
			// TODO: I can implement this for config coming from the capabilities registry.
			ContractConfigTracker:        provider.ContractConfigTracker(),
			// PASS: Implementation taken from the capability.
			// TODO: (Re)define interface?
			ContractTransmitter:          contractTransmitter,
			// TODO: Provide during setup.
			Database:                     ocrDB,
			// TODO: Provide during setup or skip.
			LocalConfig:                  lc,
			// TODO: Provide during setup.
			Logger:                       ocrLogger,
			// TODO: Provide during setup.
			MonitoringEndpoint:           oracleEndpoint,
			// TODO: Provide during setup.
			// PASS: Implementation taken from the capability.
			// TODO: (Re)define interface?
			OffchainConfigDigester:       provider.OffchainConfigDigester(),
			// TODO: Provide during setup.
			OffchainKeyring:              kb,
			// TODO: Provide during setup.
			OnchainKeyring:               onchainKeyringAdapter,
			// TODO: Provide during setup.
			MetricsRegisterer:            prometheus.WrapRegistererWith(map[string]string{"job_name": jb.Name.ValueOrZero()}, prometheus.DefaultRegisterer),
		}
		oracleArgs.ReportingPluginFactory = plugin

		INITIAL (from https://docs.google.com/document/d/1nEvHO_d1CSpP53YhUuU2xaVzXHzAkFqAo-AQTIe4-5k/edit):
		type LocalConfig struct {
			BlockchainTimeout time.Duration // This should be part of the syncer.
			ContractConfigConfirmations uint16 // This should be part of the syncer.
			SkipContractConfigConfirmations bool // This should be part of the syncer.
			// This should be part of the syncer. Because we can't really have a custom syncer
			// interval for each capability. Also, this would be blockchain specific anyway.
			ContractConfigTrackerPollInterval time.Duration
			ContractTransmitterTransmitTimeout time.Duration // This should be part of the syncer.
			MinOCR2MaxDurationQuery time.Duration // Not sure what this is for.
		}

		type OCR3OracleArgs struct {
			localConfig LocalConfig
			contractTransmitter ocr3types.ContractTransmitter[[]byte]
			offchainConfigDigester ocr3types.OffchainConfigDigester[[]byte]
			// Should this be default and shared by all oracle instance spawned by the node
			// that same way that a p2pId is?
			onchainKeyring ocr3types.OnchainKeyring[[]byte]
		}

		type NewOracleAPI interface {
			NewOracle(OracleArgs)
		}

		PROPOSED:
		type OCR3OracleArgs struct {
			// Controlling the transmitter allows the capability to send reports to a local inbox,
			// on-chain, etc.
			// Capability devs will likely send a capability response or store to outbox.
			contractTransmitter ocr3types.ContractTransmitter[[]byte]
			// The capability would need to store its on-chain config the same way that it is going
			// to be decoded by the config digester. So config encode/decode/validate should be
			// implemented.
			offchainConfigDigester ocr3types.OffchainConfigDigester[[]byte]

			// Missing?
			// - ReportingPlugin
			// - Updates to the config restarting the oracle.
		}

		type NewOracleAPI interface {
			NewOracle(OracleArgs)
		}

		QUESTIONS:
		- Do we want job.NewServiceAdapter like thing for capabilities if instances will be spawned from capabilities binary?
		- Why are oracle and plugin separate services?

	*/

	// KeyBundle - figure out if we create one on startup.
	// COuld we KeyBundles.getAll() and then use the first one?
	return &standardCapabilities{
		log:                  log,
		spec:                 spec,
		pluginRegistrar:      pluginRegistrar,
		telemetryService:     telemetryService,
		store:                store,
		CapabilitiesRegistry: CapabilitiesRegistry,
		errorLog:             errorLog,
		pipelineRunner:       pipelineRunner,
		relayerSet:           relayerSet,
		oracleFactory:        oracleFactory,
	}
}

func (s *standardCapabilities) Start(ctx context.Context) error {
	return s.StartOnce("StandardCapabilities", func() error {
		cmdName := s.spec.Command

		cmdFn, opts, err := s.pluginRegistrar.RegisterLOOP(plugins.CmdConfig{
			ID:  s.log.Name(),
			Cmd: cmdName,
			Env: nil,
		})

		if err != nil {
			return fmt.Errorf("error registering loop: %v", err)
		}

		s.capabilitiesLoop = loop.NewStandardCapabilitiesService(s.log, opts, cmdFn)

		if err = s.capabilitiesLoop.Start(ctx); err != nil {
			return fmt.Errorf("error starting standard capabilities service: %v", err)
		}

		if err = s.capabilitiesLoop.WaitCtx(ctx); err != nil {
			return fmt.Errorf("error waiting for standard capabilities service to start: %v", err)
		}

		if err = s.capabilitiesLoop.Service.Initialise(ctx, s.spec.Config, s.telemetryService, s.store, s.CapabilitiesRegistry, s.errorLog,
			s.pipelineRunner, s.relayerSet, s.oracleFactory); err != nil {
			return fmt.Errorf("error initialising standard capabilities service: %v", err)
		}

		capabilityInfos, err := s.capabilitiesLoop.Service.Infos(ctx)
		if err != nil {
			return fmt.Errorf("error getting standard capabilities service info: %v", err)
		}

		s.log.Info("Started standard capabilities for job spec", "spec", s.spec, "capabilities", capabilityInfos)

		return nil
	})
}

func (s *standardCapabilities) Close() error {
	return s.StopOnce("StandardCapabilities", func() error {
		if s.capabilitiesLoop != nil {
			return s.capabilitiesLoop.Close()
		}

		return nil
	})
}

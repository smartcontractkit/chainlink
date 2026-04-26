package jobspec

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"

	coreconfig "github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/feeds"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/nodestatusreporter/jobspec/events"
	medianconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/median/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
)

const ServiceName = "JobSpecReporter"

// Service polls active jobs and pushes their specs to Beholder, and also emits
// on job create/delete via the job.Listener interface.
type Service struct {
	services.Service
	eng *services.Engine

	config       coreconfig.JobSpecReporter
	spawner      job.Spawner
	feedsORM     feeds.ORM
	emitter      beholder.Emitter
	csaPublicKey string
	nodeVersion  string
	hostname     string
}

func NewJobSpecReporter(
	config coreconfig.JobSpecReporter,
	spawner job.Spawner,
	feedsORM feeds.ORM,
	emitter beholder.Emitter,
	csaPublicKey string,
	nodeVersion string,
	hostname string,
	lggr logger.Logger,
) *Service {
	s := &Service{
		config:       config,
		spawner:      spawner,
		feedsORM:     feedsORM,
		emitter:      emitter,
		csaPublicKey: csaPublicKey,
		nodeVersion:  nodeVersion,
		hostname:     hostname,
	}
	s.Service, s.eng = services.Config{
		Name:  ServiceName,
		Start: s.start,
	}.NewServiceEngine(lggr)
	return s
}

func (s *Service) start(ctx context.Context) error {
	if !s.config.Enabled() {
		s.eng.Info("Job Spec Reporter Service is disabled")
		return nil
	}

	s.eng.Info("Starting Job Spec Reporter Service")
	s.spawner.RegisterListener(s)
	ticker := services.NewTicker(s.config.PollingInterval())
	s.eng.GoTick(ticker, s.pollAllJobs)

	return nil
}

func (s *Service) HealthReport() map[string]error {
	return map[string]error{ServiceName: s.Ready()}
}

// OnJobStarted emits a create event when a job starts.
func (s *Service) OnJobStarted(ctx context.Context, jb job.Job) {
	if !s.ShouldEmit(&jb) {
		return
	}
	if err := s.EmitForJob(ctx, jb, events.EmissionTrigger_EMISSION_TRIGGER_CREATE); err != nil {
		s.eng.Warnw("Failed to emit job spec telemetry on create", "jobID", jb.ID, "error", err)
	}
}

// OnJobStopped emits a delete event when a job is removed.
func (s *Service) OnJobStopped(ctx context.Context, jb job.Job) {
	if !s.ShouldEmit(&jb) {
		return
	}
	if err := s.EmitForJob(ctx, jb, events.EmissionTrigger_EMISSION_TRIGGER_DELETE); err != nil {
		s.eng.Warnw("Failed to emit job spec telemetry on delete", "jobID", jb.ID, "error", err)
	}
}

// pollAllJobs emits heartbeat telemetry for every active job that passes the emit gate.
func (s *Service) pollAllJobs(ctx context.Context) {
	for _, jb := range s.spawner.ActiveJobs() {
		if !s.ShouldEmit(&jb) {
			continue
		}
		if err := s.EmitForJob(ctx, jb, events.EmissionTrigger_EMISSION_TRIGGER_HEARTBEAT); err != nil {
			s.eng.Warnw("Failed to emit job spec telemetry", "jobID", jb.ID, "error", err)
		}
	}
}

// ShouldEmit reports whether the job passes the config-driven emit gate.
// Applied to heartbeat, create, and delete events alike.
func (s *Service) ShouldEmit(j *job.Job) bool {
	if j == nil {
		return false
	}
	if j.Type != job.OffchainReporting2 || j.OCR2OracleSpec == nil {
		return s.config.EmitNonOCR2Jobs()
	}
	allowed := s.config.EnabledOCR2PluginTypes()
	if len(allowed) == 0 {
		return true
	}
	pt := string(j.OCR2OracleSpec.PluginType)
	for _, a := range allowed {
		if a == pt {
			return true
		}
	}
	return false
}

// EmitForJob builds and emits a JobSpecEvent for the given job and trigger.
func (s *Service) EmitForJob(ctx context.Context, jb job.Job, trigger events.EmissionTrigger) error {
	event, err := s.buildEvent(ctx, jb, trigger)
	if err != nil {
		return fmt.Errorf("building event: %w", err)
	}

	if err := events.EmitJobSpecEvent(ctx, s.emitter, event); err != nil {
		return fmt.Errorf("emitting event: %w", err)
	}
	return nil
}

// buildEvent converts a job.Job into its protobuf JobSpecEvent representation.
func (s *Service) buildEvent(ctx context.Context, jb job.Job, trigger events.EmissionTrigger) (*events.JobSpecEvent, error) {
	event := &events.JobSpecEvent{
		ExternalJobId:          jb.ExternalJobID.String(),
		InternalJobId:          jb.ID,
		Name:                   jb.Name.ValueOrZero(),
		JobType:                string(jb.Type),
		SchemaVersion:          jb.SchemaVersion,
		ForwardingAllowed:      jb.ForwardingAllowed,
		MaxTaskDurationSeconds: jb.MaxTaskDuration.Duration().Seconds(),
		CreatedAt:              jb.CreatedAt.Format(time.RFC3339Nano),
		CsaPublicKey:           s.csaPublicKey,
		NodeVersion:            s.nodeVersion,
		Hostname:               s.hostname,
		EmissionTrigger:        trigger,
		Timestamp:              time.Now().Format(time.RFC3339Nano),
	}

	if jb.GasLimit.Valid {
		event.GasLimit = jb.GasLimit.Uint32
	}
	if jb.StreamID != nil {
		sid := *jb.StreamID
		event.StreamId = proto.Uint32(sid)
	}

	if jb.PipelineSpec != nil {
		event.ObservationSource = jb.PipelineSpec.DotDagSource
		event.PipelineSpecId = jb.PipelineSpec.ID
		event.BridgeNames = extractBridgeNames(jb.Pipeline)
	}

	if err := s.populateProposalLifecycle(ctx, jb, event); err != nil {
		s.eng.Warnw("Failed to populate proposal lifecycle", "jobID", jb.ID, "error", err)
	}

	if jb.Type == job.OffchainReporting2 && jb.OCR2OracleSpec != nil {
		ocr2Info, err := buildOCR2OracleSpecInfo(jb.OCR2OracleSpec)
		if err != nil {
			return nil, fmt.Errorf("building OCR2OracleSpecInfo: %w", err)
		}
		event.Ocr2OracleSpec = ocr2Info
		event.ContractAddress = jb.OCR2OracleSpec.ContractID
		event.ChainId = jb.OCR2OracleSpec.ChainID
	}

	if jb.Type == job.OffchainReporting && jb.OCROracleSpec != nil {
		event.ContractAddress = jb.OCROracleSpec.ContractAddress.String()
		if jb.OCROracleSpec.EVMChainID != nil {
			event.ChainId = jb.OCROracleSpec.EVMChainID.String()
		}
		event.Ocr1OracleSpec = buildOCR1OracleSpecInfo(jb.OCROracleSpec)
	}

	return event, nil
}

// populateProposalLifecycle fills in proposal/approval fields for jobs created
// via the Feeds Manager. Jobs not managed by Feeds Manager are a no-op.
func (s *Service) populateProposalLifecycle(ctx context.Context, jb job.Job, event *events.JobSpecEvent) error {
	if s.feedsORM == nil || jb.ExternalJobID == uuid.Nil {
		return nil
	}

	prop, err := s.feedsORM.GetJobProposalByExternalJobID(ctx, jb.ExternalJobID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return fmt.Errorf("fetching job proposal: %w", err)
	}

	spec, err := s.feedsORM.GetApprovedSpec(ctx, prop.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return fmt.Errorf("fetching approved spec: %w", err)
	}

	event.FeedsManagerId = prop.FeedsManagerID
	event.RemoteUuid = prop.RemoteUUID.String()
	event.SpecVersion = spec.Version
	event.ProposedAt = spec.CreatedAt.Format(time.RFC3339Nano)
	event.ApprovedAt = spec.StatusUpdatedAt.Format(time.RFC3339Nano)
	event.AcceptLatencySeconds = spec.StatusUpdatedAt.Sub(spec.CreatedAt).Seconds()
	return nil
}

// extractBridgeNames returns the names of bridge tasks in the top-level pipeline.
// Tasks inside sub-pipelines (e.g. juelsPerFeeCoinSource) are not included.
func extractBridgeNames(p pipeline.Pipeline) []string {
	var names []string
	for _, task := range p.Tasks {
		if task.Type() != pipeline.TaskTypeBridge {
			continue
		}
		bt, ok := task.(*pipeline.BridgeTask)
		if !ok {
			continue
		}
		names = append(names, bt.Name)
	}
	return names
}

// evmRelayConfig mirrors the EVM relay config JSON so we can surface its fields
// in OCR2EVMRelayConfig without depending on the EVM module.
type evmRelayConfig struct {
	ChainID                 string   `json:"chainID"`
	FromBlock               uint64   `json:"fromBlock"`
	EffectiveTransmitterID  string   `json:"effectiveTransmitterID"`
	EnableDualTransmission  bool     `json:"enableDualTransmission"`
	EnableTriggerCapability bool     `json:"enableTriggerCapability"`
	LLODonID                uint64   `json:"lloDonID"`
	FeedID                  string   `json:"feedID"`
	SendingKeys             []string `json:"sendingKeys"`
	ProviderType            string   `json:"providerType"`
}

// buildOCR2OracleSpecInfo converts an OCR2OracleSpec into the proto message.
func buildOCR2OracleSpecInfo(spec *job.OCR2OracleSpec) (*events.OCR2OracleSpecInfo, error) {
	relayConfigRaw, err := json.Marshal(spec.RelayConfig)
	if err != nil {
		return nil, fmt.Errorf("marshaling relay config: %w", err)
	}
	pluginConfigRaw, err := json.Marshal(spec.PluginConfig)
	if err != nil {
		return nil, fmt.Errorf("marshaling plugin config: %w", err)
	}
	onchainStrategyRaw, err := json.Marshal(spec.OnchainSigningStrategy)
	if err != nil {
		return nil, fmt.Errorf("marshaling onchain signing strategy: %w", err)
	}

	feedID := ""
	if spec.FeedID != nil {
		feedID = spec.FeedID.Hex()
	}

	info := &events.OCR2OracleSpecInfo{
		SpecId:                                   spec.ID,
		ContractId:                               spec.ContractID,
		FeedId:                                   feedID,
		Relay:                                    spec.Relay,
		ChainId:                                  spec.ChainID,
		PluginType:                               string(spec.PluginType),
		TransmitterId:                            spec.TransmitterID.ValueOrZero(),
		OcrKeyBundleId:                           spec.OCRKeyBundleID.ValueOrZero(),
		MonitoringEndpoint:                       spec.MonitoringEndpoint.ValueOrZero(),
		P2Pv2Bootstrappers:                       spec.P2PV2Bootstrappers,
		AllowNoBootstrappers:                     spec.AllowNoBootstrappers,
		BlockchainTimeoutSeconds:                 spec.BlockchainTimeout.Duration().Seconds(),
		ContractConfigTrackerPollIntervalSeconds: spec.ContractConfigTrackerPollInterval.Duration().Seconds(),
		ContractConfigConfirmations:              uint32(spec.ContractConfigConfirmations),
		CaptureEaTelemetry:                       spec.CaptureEATelemetry,
		CaptureAutomationCustomTelemetry:         spec.CaptureAutomationCustomTelemetry,
		SpecCreatedAt:                            spec.CreatedAt.Format(time.RFC3339Nano),
		SpecUpdatedAt:                            spec.UpdatedAt.Format(time.RFC3339Nano),
		RelayConfigJson:                          string(relayConfigRaw),
		PluginConfigJson:                         string(pluginConfigRaw),
		OnchainSigningStrategyJson:               string(onchainStrategyRaw),
	}

	if spec.Relay == "evm" {
		evmCfg, err := buildEVMRelayConfig(relayConfigRaw)
		if err != nil {
			return nil, fmt.Errorf("building EVM relay config: %w", err)
		}
		info.EvmRelayConfig = evmCfg
	}

	if spec.PluginType == commontypes.Median {
		medianCfg, err := buildMedianPluginConfig(pluginConfigRaw)
		if err != nil {
			return nil, fmt.Errorf("building median plugin config: %w", err)
		}
		info.MedianPluginConfig = medianCfg
	}

	return info, nil
}

func buildOCR1OracleSpecInfo(spec *job.OCROracleSpec) *events.OCR1OracleSpecInfo {
	keyBundleID := ""
	if spec.EncryptedOCRKeyBundleID != nil {
		keyBundleID = spec.EncryptedOCRKeyBundleID.String()
	}

	transmitterAddress := ""
	if spec.TransmitterAddress != nil {
		transmitterAddress = spec.TransmitterAddress.String()
	}

	var dbTimeoutSeconds float64
	if spec.DatabaseTimeout != nil {
		dbTimeoutSeconds = spec.DatabaseTimeout.Duration().Seconds()
	}

	var gracePeriodSeconds float64
	if spec.ObservationGracePeriod != nil {
		gracePeriodSeconds = spec.ObservationGracePeriod.Duration().Seconds()
	}

	var transmitTimeoutSeconds float64
	if spec.ContractTransmitterTransmitTimeout != nil {
		transmitTimeoutSeconds = spec.ContractTransmitterTransmitTimeout.Duration().Seconds()
	}

	return &events.OCR1OracleSpecInfo{
		SpecId:                    spec.ID,
		P2Pv2Bootstrappers:        spec.P2PV2Bootstrappers,
		IsBootstrapPeer:           spec.IsBootstrapPeer,
		EncryptedOcrKeyBundleId:   keyBundleID,
		TransmitterAddress:        transmitterAddress,
		ObservationTimeoutSeconds: spec.ObservationTimeout.Duration().Seconds(),
		BlockchainTimeoutSeconds:  spec.BlockchainTimeout.Duration().Seconds(),
		ContractConfigTrackerSubscribeIntervalSeconds: spec.ContractConfigTrackerSubscribeInterval.Duration().Seconds(),
		ContractConfigTrackerPollIntervalSeconds:      spec.ContractConfigTrackerPollInterval.Duration().Seconds(),
		ContractConfigConfirmations:                   uint32(spec.ContractConfigConfirmations),
		DatabaseTimeoutSeconds:                        dbTimeoutSeconds,
		ObservationGracePeriodSeconds:                 gracePeriodSeconds,
		ContractTransmitterTransmitTimeoutSeconds:     transmitTimeoutSeconds,
		CaptureEaTelemetry:                            spec.CaptureEATelemetry,
		SpecCreatedAt:                                 spec.CreatedAt.Format(time.RFC3339Nano),
		SpecUpdatedAt:                                 spec.UpdatedAt.Format(time.RFC3339Nano),
	}
}

// buildEVMRelayConfig decodes the EVM relay config JSON into OCR2EVMRelayConfig.
func buildEVMRelayConfig(relayConfigJSON []byte) (*events.OCR2EVMRelayConfig, error) {
	var cfg evmRelayConfig
	if err := json.Unmarshal(relayConfigJSON, &cfg); err != nil {
		return nil, fmt.Errorf("unmarshaling EVM relay config: %w", err)
	}

	return &events.OCR2EVMRelayConfig{
		ChainId:                 cfg.ChainID,
		FromBlock:               cfg.FromBlock,
		EffectiveTransmitterId:  cfg.EffectiveTransmitterID,
		EnableDualTransmission:  cfg.EnableDualTransmission,
		EnableTriggerCapability: cfg.EnableTriggerCapability,
		LloDonId:                cfg.LLODonID,
		FeedId:                  cfg.FeedID,
		SendingKeys:             cfg.SendingKeys,
		ProviderType:            cfg.ProviderType,
	}, nil
}

// buildMedianPluginConfig decodes the median plugin config JSON into OCR2MedianPluginConfig.
func buildMedianPluginConfig(pluginConfigJSON []byte) (*events.OCR2MedianPluginConfig, error) {
	var cfg medianconfig.PluginConfig
	if err := json.Unmarshal(pluginConfigJSON, &cfg); err != nil {
		return nil, fmt.Errorf("unmarshaling median plugin config: %w", err)
	}

	medianProto := &events.OCR2MedianPluginConfig{
		JuelsPerFeeCoinSource:  cfg.JuelsPerFeeCoinPipeline,
		GasPriceSubunitsSource: cfg.GasPriceSubunitsPipeline,
	}

	if cfg.JuelsPerFeeCoinCache == nil {
		medianProto.JuelsPerFeeCoinCacheDisabled = true
	} else {
		medianProto.JuelsPerFeeCoinCacheDisabled = cfg.JuelsPerFeeCoinCache.Disable
		medianProto.JuelsPerFeeCoinCacheUpdateIntervalSeconds = cfg.JuelsPerFeeCoinCache.UpdateInterval.Duration().Seconds()
		medianProto.JuelsPerFeeCoinCacheStalenessAlertThresholdSeconds = cfg.JuelsPerFeeCoinCache.StalenessAlertThreshold.Duration().Seconds()
	}

	if cfg.DeviationFunctionDefinition != nil {
		devFuncRaw, err := json.Marshal(cfg.DeviationFunctionDefinition)
		if err != nil {
			return nil, fmt.Errorf("marshaling deviation function definition: %w", err)
		}
		medianProto.DeviationFuncJson = string(devFuncRaw)
	}

	return medianProto, nil
}

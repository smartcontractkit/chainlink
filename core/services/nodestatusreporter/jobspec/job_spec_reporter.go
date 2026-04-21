package jobspec

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/services"

	coreconfig "github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	medianconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/median/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/feeds"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/nodestatusreporter/jobspec/events"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
)

// NodeInfo contains static node identity values injected at construction time.
type NodeInfo struct {
	CSAPublicKey string
	NodeVersion  string
	Hostname     string
}

const ServiceName = "JobSpecReporter"

// Service polls active jobs and emits full job-spec telemetry via Beholder.
// It mirrors the BridgeStatusReporter structure and implements job.Listener
// to also emit immediately on create/delete.
type Service struct {
	services.Service
	eng *services.Engine

	config   coreconfig.JobSpecReporter
	spawner  job.Spawner
	feedsORM feeds.ORM // optional; nil-safe
	emitter  beholder.Emitter
	nodeInfo NodeInfo
}

// NewJobSpecReporter creates a new Job Spec Reporter Service.
func NewJobSpecReporter(
	config coreconfig.JobSpecReporter,
	spawner job.Spawner,
	feedsORM feeds.ORM,
	emitter beholder.Emitter,
	nodeInfo NodeInfo,
	lggr logger.Logger,
) *Service {
	s := &Service{
		config:   config,
		spawner:  spawner,
		feedsORM: feedsORM,
		emitter:  emitter,
		nodeInfo: nodeInfo,
	}
	s.Service, s.eng = services.Config{
		Name:  ServiceName,
		Start: s.start,
	}.NewServiceEngine(lggr)
	return s
}

// start starts the Job Spec Reporter Service.
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

// HealthReport returns the service health.
func (s *Service) HealthReport() map[string]error {
	return map[string]error{ServiceName: s.Ready()}
}

// OnJobStarted implements job.Listener — called after a job service starts successfully.
func (s *Service) OnJobStarted(ctx context.Context, jb job.Job) {
	if !s.shouldEmit(&jb) {
		return
	}
	if err := s.emitForJob(ctx, jb, "create"); err != nil {
		s.eng.Warnw("Failed to emit job spec telemetry on create", "jobID", jb.ID, "error", err)
	}
}

// OnJobStopped implements job.Listener — called after a job is deleted.
func (s *Service) OnJobStopped(ctx context.Context, jb job.Job) {
	if !s.shouldEmit(&jb) {
		return
	}
	if err := s.emitForJob(ctx, jb, "delete"); err != nil {
		s.eng.Warnw("Failed to emit job spec telemetry on delete", "jobID", jb.ID, "error", err)
	}
}

// pollAllJobs is called on each heartbeat tick and emits telemetry for every
// active job that passes the shouldEmit gate.
func (s *Service) pollAllJobs(ctx context.Context) {
	activeJobs := s.spawner.ActiveJobs()
	for _, jb := range activeJobs {
		jbCopy := jb
		if !s.shouldEmit(&jbCopy) {
			continue
		}
		if err := s.emitForJob(ctx, jbCopy, "heartbeat"); err != nil {
			s.eng.Warnw("Failed to emit job spec telemetry", "jobID", jb.ID, "error", err)
		}
	}
}

// ShouldEmit is exported for testing. It returns true when the given job
// should produce a telemetry event based on the current config.
// Use shouldEmit for internal calls (identical logic).
func (s *Service) ShouldEmit(j *job.Job) bool {
	return s.shouldEmit(j)
}

// shouldEmit returns true when the given job should produce a telemetry event
// based on the current config. This gate applies symmetrically to heartbeats
// and create/delete events.
func (s *Service) shouldEmit(j *job.Job) bool {
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

// EmitForJob is exported for testing. It converts a job to a JobSpecEvent and
// emits it via Beholder. Use emitForJob for internal calls (identical logic).
func (s *Service) EmitForJob(ctx context.Context, jb job.Job, trigger string) error {
	return s.emitForJob(ctx, jb, trigger)
}

// emitForJob converts a job to a JobSpecEvent and emits it via Beholder.
func (s *Service) emitForJob(ctx context.Context, jb job.Job, trigger string) error {
	event, err := s.buildEvent(ctx, jb, trigger)
	if err != nil {
		return fmt.Errorf("building event: %w", err)
	}

	if err := events.EmitJobSpecEvent(ctx, s.emitter, event); err != nil {
		return fmt.Errorf("emitting event: %w", err)
	}
	return nil
}

// buildEvent converts a job.Job into the protobuf JobSpecEvent.
func (s *Service) buildEvent(ctx context.Context, jb job.Job, trigger string) (*events.JobSpecEvent, error) {
	event := &events.JobSpecEvent{
		ExternalJobId:          jb.ExternalJobID.String(),
		InternalJobId:          jb.ID,
		Name:                   jb.Name.ValueOrZero(),
		JobType:                string(jb.Type),
		SchemaVersion:          jb.SchemaVersion,
		ForwardingAllowed:      jb.ForwardingAllowed,
		MaxTaskDurationSeconds: jb.MaxTaskDuration.Duration().Seconds(),
		CreatedAt:              jb.CreatedAt.Format(time.RFC3339Nano),
		CsaPublicKey:           s.nodeInfo.CSAPublicKey,
		NodeVersion:            s.nodeInfo.NodeVersion,
		Hostname:               s.nodeInfo.Hostname,
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

	s.populateProposalLifecycle(ctx, jb, event)

	if jb.Type == job.OffchainReporting2 && jb.OCR2OracleSpec != nil {
		ocr2Info, err := buildOCR2OracleSpecInfo(jb.OCR2OracleSpec)
		if err != nil {
			return nil, fmt.Errorf("building OCR2OracleSpecInfo: %w", err)
		}
		event.Ocr2OracleSpec = ocr2Info
	}

	return event, nil
}

// populateProposalLifecycle fills the proposal/approval fields when the job was
// created via the Feeds Manager. Missing rows (manually created job) are silently
// ignored.
func (s *Service) populateProposalLifecycle(ctx context.Context, jb job.Job, event *events.JobSpecEvent) {
	if s.feedsORM == nil || jb.ExternalJobID == (uuid.UUID{}) {
		return
	}

	prop, err := s.feedsORM.GetJobProposalByExternalJobID(ctx, jb.ExternalJobID)
	if err != nil {
		return
	}

	spec, err := s.feedsORM.GetApprovedSpec(ctx, prop.ID)
	if err != nil {
		return
	}

	event.FeedsManagerId = prop.FeedsManagerID
	event.RemoteUuid = prop.RemoteUUID.String()
	event.SpecVersion = int32(spec.Version)
	event.ProposedAt = spec.CreatedAt.Format(time.RFC3339Nano)
	event.ApprovedAt = spec.StatusUpdatedAt.Format(time.RFC3339Nano)
	event.AcceptLatencySeconds = spec.StatusUpdatedAt.Sub(spec.CreatedAt).Seconds()
}

// extractBridgeNames returns the names of all bridge tasks in the top-level
// observationSource pipeline. Tasks in sub-pipelines (e.g. juelsPerFeeCoinSource)
// are not included.
func extractBridgeNames(p pipeline.Pipeline) []string {
	var names []string
	for _, task := range p.Tasks {
		if task.Type() == pipeline.TaskTypeBridge {
			bt := task.(*pipeline.BridgeTask)
			names = append(names, bt.Name)
		}
	}
	return names
}

// evmRelayConfig is a minimal struct for decoding EVM relay config JSON fields
// that we want to surface in OCR2EVMRelayConfig without importing the EVM module.
type evmRelayConfig struct {
	ChainID                  string   `json:"chainID"`
	FromBlock                uint64   `json:"fromBlock"`
	EffectiveTransmitterID   string   `json:"effectiveTransmitterID"`
	EnableDualTransmission   bool     `json:"enableDualTransmission"`
	EnableTriggerCapability  bool     `json:"enableTriggerCapability"`
	LLODonID                 uint64   `json:"lloDonID"`
	FeedID                   string   `json:"feedID"`
	SendingKeys              []string `json:"sendingKeys"`
	ProviderType             string   `json:"providerType"`
}

// buildOCR2OracleSpecInfo converts an OCR2OracleSpec into its proto representation.
func buildOCR2OracleSpecInfo(spec *job.OCR2OracleSpec) (*events.OCR2OracleSpecInfo, error) {
	relayConfigJSON := ""
	if raw, err := json.Marshal(spec.RelayConfig); err == nil {
		relayConfigJSON = string(raw)
	}
	pluginConfigJSON := ""
	if raw, err := json.Marshal(spec.PluginConfig); err == nil {
		pluginConfigJSON = string(raw)
	}
	onchainStrategyJSON := ""
	if raw, err := json.Marshal(spec.OnchainSigningStrategy); err == nil {
		onchainStrategyJSON = string(raw)
	}

	feedID := ""
	if spec.FeedID != nil {
		feedID = spec.FeedID.Hex()
	}

	info := &events.OCR2OracleSpecInfo{
		SpecId:                                      spec.ID,
		ContractId:                                  spec.ContractID,
		FeedId:                                      feedID,
		Relay:                                       spec.Relay,
		ChainId:                                     spec.ChainID,
		PluginType:                                  string(spec.PluginType),
		TransmitterId:                               spec.TransmitterID.ValueOrZero(),
		OcrKeyBundleId:                              spec.OCRKeyBundleID.ValueOrZero(),
		MonitoringEndpoint:                          spec.MonitoringEndpoint.ValueOrZero(),
		P2Pv2Bootstrappers:                          spec.P2PV2Bootstrappers,
		AllowNoBootstrappers:                        spec.AllowNoBootstrappers,
		BlockchainTimeoutSeconds:                    spec.BlockchainTimeout.Duration().Seconds(),
		ContractConfigTrackerPollIntervalSeconds:     spec.ContractConfigTrackerPollInterval.Duration().Seconds(),
		ContractConfigConfirmations:                 uint32(spec.ContractConfigConfirmations),
		CaptureEaTelemetry:                          spec.CaptureEATelemetry,
		CaptureAutomationCustomTelemetry:            spec.CaptureAutomationCustomTelemetry,
		SpecCreatedAt:                               spec.CreatedAt.Format(time.RFC3339Nano),
		SpecUpdatedAt:                               spec.UpdatedAt.Format(time.RFC3339Nano),
		RelayConfigJson:                             relayConfigJSON,
		PluginConfigJson:                            pluginConfigJSON,
		OnchainSigningStrategyJson:                  onchainStrategyJSON,
	}

	if spec.Relay == "evm" {
		evmCfg, err := buildEVMRelayConfig(spec)
		if err != nil {
			return nil, errors.Wrap(err, "building EVM relay config")
		}
		info.EvmRelayConfig = evmCfg
	}

	if spec.PluginType == commontypes.Median {
		medianCfg, err := buildMedianPluginConfig(spec)
		if err != nil {
			return nil, errors.Wrap(err, "building median plugin config")
		}
		info.MedianPluginConfig = medianCfg
	}

	return info, nil
}

// buildEVMRelayConfig decodes the EVM relay config JSON into the proto message.
func buildEVMRelayConfig(spec *job.OCR2OracleSpec) (*events.OCR2EVMRelayConfig, error) {
	raw, err := json.Marshal(spec.RelayConfig)
	if err != nil {
		return nil, fmt.Errorf("marshaling relay config: %w", err)
	}

	var cfg evmRelayConfig
	if err := json.Unmarshal(raw, &cfg); err != nil {
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

// buildMedianPluginConfig decodes the plugin config JSON into the typed proto message.
func buildMedianPluginConfig(spec *job.OCR2OracleSpec) (*events.OCR2MedianPluginConfig, error) {
	raw, err := json.Marshal(spec.PluginConfig)
	if err != nil {
		return nil, fmt.Errorf("marshaling plugin config: %w", err)
	}

	var cfg medianconfig.PluginConfig
	if err := json.Unmarshal(raw, &cfg); err != nil {
		return nil, fmt.Errorf("unmarshaling median plugin config: %w", err)
	}

	medianProto := &events.OCR2MedianPluginConfig{
		JuelsPerFeeCoinSource: cfg.JuelsPerFeeCoinPipeline,
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
		if err == nil {
			medianProto.DeviationFuncJson = string(devFuncRaw)
		}
	}

	return medianProto, nil
}

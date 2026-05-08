package jobspec

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
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
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
)

const ServiceName = "JobSpecReporter"

var _ job.Listener = (*Service)(nil)

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

// AfterJobStarted emits a create event when a job starts.
func (s *Service) AfterJobStarted(ctx context.Context, jb job.Job) {
	if !s.ShouldEmit(&jb) {
		return
	}
	if err := s.EmitForJob(ctx, jb, events.EmissionTrigger_EMISSION_TRIGGER_CREATE); err != nil {
		s.eng.Warnw("Failed to emit job spec telemetry on create", "jobID", jb.ID, "error", err)
	}
}

// AfterJobStopped emits a delete event when a job is removed.
func (s *Service) AfterJobStopped(ctx context.Context, jb job.Job) {
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
func (s *Service) ShouldEmit(j *job.Job) bool {
	if j == nil {
		return false
	}
	if j.Type != job.OffchainReporting2 || j.OCR2OracleSpec == nil {
		return false
	}
	allowed := s.config.EnabledOCR2PluginTypes()
	if len(allowed) == 0 {
		return false
	}
	if slices.Contains(allowed, "all") {
		return true
	}
	return slices.Contains(allowed, string(j.OCR2OracleSpec.PluginType))
}

// EmitForJob builds and emits a JobSpecEvent for the given job and trigger.
func (s *Service) EmitForJob(ctx context.Context, jb job.Job, trigger events.EmissionTrigger) error {
	if jb.Type != job.OffchainReporting2 || jb.OCR2OracleSpec == nil {
		return fmt.Errorf("unsupported job type %s", jb.Type)
	}

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
		ExternalJobId:     jb.ExternalJobID.String(),
		Name:              jb.Name.ValueOrZero(),
		JobType:           string(jb.Type),
		SchemaVersion:     jb.SchemaVersion,
		ForwardingAllowed: jb.ForwardingAllowed,
		CreatedAt:         jb.CreatedAt.Format(time.RFC3339Nano),
		CsaPublicKey:      s.csaPublicKey,
		NodeVersion:       s.nodeVersion,
		Hostname:          s.hostname,
		EmissionTrigger:   trigger,
		Timestamp:         time.Now().Format(time.RFC3339Nano),
	}

	if jb.GasLimit.Valid {
		event.GasLimit = proto.Uint32(jb.GasLimit.Uint32)
	}
	if jb.StreamID != nil {
		sid := *jb.StreamID
		event.StreamId = proto.Uint32(sid)
	}

	if jb.PipelineSpec != nil {
		event.ObservationSource = jb.PipelineSpec.DotDagSource
		bridgeNames, err := extractBridgeNames(jb)
		if err != nil {
			return nil, fmt.Errorf("extracting bridge names: %w", err)
		}
		event.BridgeNames = bridgeNames
	}

	if err := s.populateProposalLifecycle(ctx, jb, event); err != nil {
		s.eng.Warnw("Failed to populate proposal lifecycle", "jobID", jb.ID, "error", err)
	}

	ocr2Info, err := buildOCR2OracleSpecInfo(jb.OCR2OracleSpec)
	if err != nil {
		return nil, fmt.Errorf("building OCR2OracleSpecInfo: %w", err)
	}
	event.Ocr2OracleSpec = ocr2Info

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
func extractBridgeNames(jb job.Job) ([]string, error) {
	names := extractBridgeNamesFromPipeline(jb.Pipeline)
	if len(names) > 0 || jb.PipelineSpec == nil || jb.PipelineSpec.DotDagSource == "" {
		return names, nil
	}

	p, err := pipeline.Parse(jb.PipelineSpec.DotDagSource)
	if err != nil {
		return nil, err
	}
	return extractBridgeNamesFromPipeline(*p), nil
}

func extractBridgeNamesFromPipeline(p pipeline.Pipeline) []string {
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

func ocr2ChainID(spec *job.OCR2OracleSpec) string {
	if spec == nil {
		return ""
	}
	if spec.ChainID != "" {
		return spec.ChainID
	}
	relayID, err := spec.RelayID()
	if err != nil {
		return ""
	}
	return relayID.ChainID
}

// evmRelayConfig mirrors the EVM relay config JSON so we can surface its fields
// in OCR2EVMRelayConfig without depending on the EVM module.
type evmRelayConfig struct {
	FromBlock               *uint64  `json:"fromBlock"`
	EffectiveTransmitterID  *string  `json:"effectiveTransmitterID"`
	EnableDualTransmission  *bool    `json:"enableDualTransmission"`
	EnableTriggerCapability *bool    `json:"enableTriggerCapability"`
	LLODonID                *uint64  `json:"lloDonID"`
	FeedID                  *string  `json:"feedID"`
	SendingKeys             []string `json:"sendingKeys"`
	ProviderType            *string  `json:"providerType"`
}

type medianPluginConfig struct {
	JuelsPerFeeCoinPipeline *string `json:"juelsPerFeeCoinSource"`
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
	info := &events.OCR2OracleSpecInfo{
		ContractId:         spec.ContractID,
		Relay:              spec.Relay,
		PluginType:         string(spec.PluginType),
		CaptureEaTelemetry: spec.CaptureEATelemetry,
		RelayConfigJson:    string(relayConfigRaw),
		PluginConfigJson:   string(pluginConfigRaw),
	}

	if spec.FeedID != nil {
		info.FeedId = proto.String(spec.FeedID.Hex())
	}
	if spec.TransmitterID.Valid {
		info.TransmitterId = proto.String(spec.TransmitterID.String)
	}
	if spec.OCRKeyBundleID.Valid {
		info.OcrKeyBundleId = proto.String(spec.OCRKeyBundleID.String)
	}

	if spec.Relay == "evm" {
		evmCfg, err := buildEVMRelayConfig(relayConfigRaw, ocr2ChainID(spec), spec.TransmitterID.ValueOrZero())
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

// buildEVMRelayConfig decodes the EVM relay config JSON into OCR2EVMRelayConfig.
func buildEVMRelayConfig(relayConfigJSON []byte, chainID, transmitterID string) (*events.OCR2EVMRelayConfig, error) {
	var cfg evmRelayConfig
	if err := json.Unmarshal(relayConfigJSON, &cfg); err != nil {
		return nil, fmt.Errorf("unmarshaling EVM relay config: %w", err)
	}

	effectiveTransmitterID := transmitterID
	if cfg.EffectiveTransmitterID != nil {
		effectiveTransmitterID = *cfg.EffectiveTransmitterID
	}

	evmProto := &events.OCR2EVMRelayConfig{
		ChainId:                chainID,
		EffectiveTransmitterId: effectiveTransmitterID,
		SendingKeys:            cfg.SendingKeys,
	}
	if cfg.FromBlock != nil {
		evmProto.FromBlock = proto.Uint64(*cfg.FromBlock)
	}
	if cfg.EnableDualTransmission != nil {
		evmProto.EnableDualTransmission = proto.Bool(*cfg.EnableDualTransmission)
	}
	if cfg.EnableTriggerCapability != nil {
		evmProto.EnableTriggerCapability = proto.Bool(*cfg.EnableTriggerCapability)
	}
	if cfg.LLODonID != nil {
		evmProto.LloDonId = proto.Uint64(*cfg.LLODonID)
	}
	if cfg.FeedID != nil {
		evmProto.FeedId = proto.String(*cfg.FeedID)
	}
	if cfg.ProviderType != nil {
		evmProto.ProviderType = proto.String(*cfg.ProviderType)
	}
	return evmProto, nil
}

// buildMedianPluginConfig decodes the median plugin config JSON into OCR2MedianPluginConfig.
func buildMedianPluginConfig(pluginConfigJSON []byte) (*events.OCR2MedianPluginConfig, error) {
	var cfg medianPluginConfig
	if err := json.Unmarshal(pluginConfigJSON, &cfg); err != nil {
		return nil, fmt.Errorf("unmarshaling median plugin config: %w", err)
	}

	medianProto := &events.OCR2MedianPluginConfig{}
	if cfg.JuelsPerFeeCoinPipeline != nil {
		medianProto.JuelsPerFeeCoinSource = *cfg.JuelsPerFeeCoinPipeline
	}

	return medianProto, nil
}

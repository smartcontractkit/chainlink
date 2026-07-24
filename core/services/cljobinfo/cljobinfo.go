// Package cljobinfo provides a single, generic way for any part of the core
// node to emit a full job definition as telemetry.
//
// It supersedes the two prior, divergent approaches:
//
//   - NodePlatformJobInfoService (core/services/chainlink/node_platform.go)
//     emits a flat, denormalized projection (chain_id/job_type/field_path ->
//     addresses) aggregated across all jobs. Generic schema, but only carries
//     submitter addresses and has hardcoded per-spec extractors.
//   - JobSpecReporter (core/services/nodestatusreporter/jobspec) emits a
//     rich, per-job, create/delete/heartbeat event, but models each spec as a
//     dedicated proto message and only supports OCR2.
//
// cljobinfo keeps the better half of each: the event-driven, per-job lifecycle
// and node identity of JobSpecReporter, and a schema-agnostic payload like
// NodePlatformJobInfo. The full, type-specific job definition is carried as a
// raw TOML string, so any job type is supported with no per-type code here and
// none for future job types.
package cljobinfo

import (
	"context"
	"fmt"
	"time"

	"github.com/pelletier/go-toml"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	commonv1 "github.com/smartcontractkit/chainlink-protos/node-platform/common/v1"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

const (
	// Domain, Entity and DataSchema identify CLJobInfo telemetry on Beholder.
	Domain     = "node-platform"
	Entity     = "common.v1.CLJobInfo"
	DataSchema = "/node-platform/common/v1"

	ServiceName = "CLJobInfoReporter"

	// DefaultPollInterval is the heartbeat cadence when a caller does not
	// specify one, matching the node-platform build/job info beat.
	DefaultPollInterval = 3 * time.Minute
)

// NodeIdentity is the node-level context attached to every emitted CLJobInfo.
type NodeIdentity struct {
	CSAPublicKey string
	NodeVersion  string
	Hostname     string
}

// Build converts any job.Job into its generic CLJobInfo representation.
//
// The complete, type-specific job definition is captured as TOML, so no
// per-job-type code lives here and none is needed for future job types. If the
// job cannot be TOML-encoded, Build still returns a fully populated identity
// payload (with an empty SpecToml) alongside the encoding error, so callers can
// choose to emit the envelope and log the failure rather than drop the event.
func Build(jb job.Job, trigger commonv1.CLJobInfoTrigger, id NodeIdentity, now time.Time) (*commonv1.CLJobInfo, error) {
	info := &commonv1.CLJobInfo{
		CsaPublicKey:      id.CSAPublicKey,
		NodeVersion:       id.NodeVersion,
		Hostname:          id.Hostname,
		ExternalJobId:     jb.ExternalJobID.String(),
		JobId:             jb.ID,
		Name:              jb.Name.ValueOrZero(),
		JobType:           string(jb.Type),
		SchemaVersion:     jb.SchemaVersion,
		ForwardingAllowed: jb.ForwardingAllowed,
		CreatedAt:         formatTime(jb.CreatedAt),
		Trigger:           trigger,
		Timestamp:         now.UTC().Format(time.RFC3339Nano),
	}
	if jb.GasLimit.Valid {
		info.GasLimit = new(jb.GasLimit.Uint32)
	}
	if jb.StreamID != nil {
		info.StreamId = new(*jb.StreamID)
	}

	specTOML, err := jobTOML(jb)
	if err != nil {
		return info, fmt.Errorf("encoding job %s (%d) spec to TOML: %w", jb.ExternalJobID, jb.ID, err)
	}
	info.SpecToml = specTOML

	return info, nil
}

// Emit marshals a CLJobInfo and publishes it to Beholder.
func Emit(ctx context.Context, emitter beholder.Emitter, info *commonv1.CLJobInfo) error {
	payload, err := proto.Marshal(info)
	if err != nil {
		return fmt.Errorf("marshaling CLJobInfo: %w", err)
	}

	err = emitter.Emit(ctx, payload,
		beholder.AttrKeyDomain, Domain,
		beholder.AttrKeyEntity, Entity,
		beholder.AttrKeyDataSchema, DataSchema,
	)
	if err != nil {
		return fmt.Errorf("emitting CLJobInfo: %w", err)
	}
	return nil
}

// jobTOML serializes the entire job definition to TOML. Marshaling the whole
// job.Job captures both the common top-level fields and the single active
// type-specific spec, so all fields for any job type are included without
// enumerating them.
func jobTOML(jb job.Job) (string, error) {
	out, err := toml.Marshal(jb)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format(time.RFC3339Nano)
}

var _ job.Listener = (*Reporter)(nil)

// Reporter emits a CLJobInfo for every job on create, delete, and on a
// recurring heartbeat. It is job-type agnostic: any job the node runs is
// reported through the single generic schema.
type Reporter struct {
	services.Service
	eng *services.Engine

	spawner      job.Spawner
	emitter      beholder.Emitter
	identity     NodeIdentity
	pollInterval time.Duration
}

// NewReporter builds a Reporter that reports every job the node runs.
func NewReporter(
	spawner job.Spawner,
	emitter beholder.Emitter,
	identity NodeIdentity,
	pollInterval time.Duration,
	lggr logger.Logger,
) *Reporter {
	r := &Reporter{
		spawner:      spawner,
		emitter:      emitter,
		identity:     identity,
		pollInterval: pollInterval,
	}
	r.Service, r.eng = services.Config{
		Name:  ServiceName,
		Start: r.start,
	}.NewServiceEngine(lggr)
	return r
}

func (r *Reporter) start(_ context.Context) error {
	r.spawner.RegisterListener(r)
	r.eng.GoTick(services.NewTicker(r.pollInterval), r.pollAllJobs)
	return nil
}

func (r *Reporter) HealthReport() map[string]error {
	return map[string]error{ServiceName: r.Ready()}
}

// AfterJobStarted emits a create event when a job starts.
func (r *Reporter) AfterJobStarted(ctx context.Context, jb job.Job) {
	r.emitForJob(ctx, jb, commonv1.CLJobInfoTrigger_CL_JOB_INFO_TRIGGER_CREATE)
}

// AfterJobStopped emits a delete event when a job is removed.
func (r *Reporter) AfterJobStopped(ctx context.Context, jb job.Job) {
	r.emitForJob(ctx, jb, commonv1.CLJobInfoTrigger_CL_JOB_INFO_TRIGGER_DELETE)
}

// pollAllJobs emits a heartbeat event for every active job.
func (r *Reporter) pollAllJobs(ctx context.Context) {
	for _, jb := range r.spawner.ActiveJobs() {
		r.emitForJob(ctx, jb, commonv1.CLJobInfoTrigger_CL_JOB_INFO_TRIGGER_HEARTBEAT)
	}
}

func (r *Reporter) emitForJob(ctx context.Context, jb job.Job, trigger commonv1.CLJobInfoTrigger) {
	info, err := Build(jb, trigger, r.identity, time.Now())
	if err != nil {
		// Spec encoding failed; still emit the identity envelope so the job is
		// accounted for, but flag the gap.
		r.eng.Warnw("Failed to encode job spec for CLJobInfo; emitting without spec_toml",
			"jobID", jb.ID, "externalJobID", jb.ExternalJobID, "error", err)
	}

	if err := Emit(ctx, r.emitter, info); err != nil {
		r.eng.Warnw("Failed to emit CLJobInfo", "jobID", jb.ID, "trigger", trigger, "error", err)
	}
}

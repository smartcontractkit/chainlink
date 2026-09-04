package jobspec

import (
	"context"
	"fmt"
	"time"

	"github.com/pelletier/go-toml"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	commonv1 "github.com/smartcontractkit/chainlink-protos/node-platform/common/v1"

	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

// CLJobInfo is the generic, job-type-agnostic half of this reporter. Where
// JobSpecEvent models one job type (OCR2) field by field, CLJobInfo carries the
// job's common identity plus the complete definition as a raw TOML string, so
// every job the node runs is reported through a single schema with no
// per-type code here and none needed for future job types.
//
// It is emitted on the same triggers and from the same service as JobSpecEvent
// rather than from a parallel one, so there is exactly one place in the node
// that reports what jobs it runs. Once consumers have migrated, the OCR2-only
// half can be deleted from here without touching the wiring.
const (
	// Domain, Entity and DataSchema identify CLJobInfo telemetry on Beholder.
	Domain     = "node-platform"
	Entity     = "common.v1.CLJobInfo"
	DataSchema = "/node-platform/common/v1"
)

// NodeIdentity is the node-level context attached to every emitted CLJobInfo.
type NodeIdentity struct {
	CSAPublicKey string
	NodeVersion  string
	Hostname     string
}

// JobProposal is the Job Distributor provenance for a job that arrived as an
// approved job proposal. Jobs created directly (CLI, UI, TOML on disk) have no
// proposal, and the zero value leaves the corresponding CLJobInfo fields unset
// — which is how a consumer tells a managed job from an unmanaged one.
type JobProposal struct {
	FeedsManagerID int64
	RemoteUUID     string
	SpecVersion    int32
	ProposedAt     time.Time
	ApprovedAt     time.Time
}

// BuildCLJobInfo converts any job.Job into its generic CLJobInfo representation.
//
// prop is optional: pass nil for a job with no Job Distributor proposal.
//
// If the job cannot be TOML-encoded, BuildCLJobInfo still returns a fully
// populated identity payload (with an empty SpecToml) alongside the encoding
// error, so callers can choose to emit the envelope and log the failure rather
// than drop the event.
func BuildCLJobInfo(jb job.Job, trigger commonv1.CLJobInfoTrigger, id NodeIdentity, prop *JobProposal, now time.Time) (*commonv1.CLJobInfo, error) {
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
		CreatedAt:         formatCLJobInfoTime(jb.CreatedAt),
		Trigger:           trigger,
		Timestamp:         now.UTC().Format(time.RFC3339Nano),
	}
	if jb.GasLimit.Valid {
		info.GasLimit = new(jb.GasLimit.Uint32)
	}
	if jb.StreamID != nil {
		info.StreamId = new(*jb.StreamID)
	}
	if prop != nil {
		info.FeedsManagerId = &prop.FeedsManagerID
		info.RemoteUuid = &prop.RemoteUUID
		info.SpecVersion = &prop.SpecVersion
		info.ProposedAt = new(formatCLJobInfoTime(prop.ProposedAt))
		info.ApprovedAt = new(formatCLJobInfoTime(prop.ApprovedAt))
	}

	specTOML, err := jobTOML(jb)
	if err != nil {
		return info, fmt.Errorf("encoding job %s (%d) spec to TOML: %w", jb.ExternalJobID, jb.ID, err)
	}
	info.SpecToml = specTOML

	return info, nil
}

// EmitCLJobInfo marshals a CLJobInfo and publishes it to Beholder.
func EmitCLJobInfo(ctx context.Context, emitter beholder.Emitter, info *commonv1.CLJobInfo) error {
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

func formatCLJobInfoTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format(time.RFC3339Nano)
}

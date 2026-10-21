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

// CLJobInfo is the job-type-agnostic half of this reporter: common identity
// plus the whole definition as TOML, so every job type is covered with no
// per-type code. JobSpecEvent (OCR2-only) is emitted from the same service.
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

// JobProposal is the JD provenance for a job that arrived as an approved job
// proposal. Nil for jobs created directly (CLI, UI, TOML on disk).
type JobProposal struct {
	FeedsManagerID int64
	RemoteUUID     string
	SpecVersion    int32
	ProposedAt     time.Time
	ApprovedAt     time.Time
}

// BuildCLJobInfo converts any job.Job into a CLJobInfo. prop may be nil.
//
// On TOML encoding failure it still returns a populated identity payload with
// an empty SpecToml alongside the error, so callers can emit and log rather
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
		CreatedAtMs:       unixMillisOrNil(jb.CreatedAt),
		Trigger:           trigger,
		TimestampMs:       now.UnixMilli(),
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
		info.ProposedAtMs = unixMillisOrNil(prop.ProposedAt)
		info.ApprovedAtMs = unixMillisOrNil(prop.ApprovedAt)
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

// jobTOML serializes the whole job.Job, which captures both the common fields
// and the single active type-specific spec.
func jobTOML(jb job.Job) (string, error) {
	out, err := toml.Marshal(jb)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// unixMillisOrNil maps an unset time to nil rather than the epoch. Millis not
// RFC3339Nano: Go trims trailing zeros, so those strings are variable-width and
// don't sort chronologically.
func unixMillisOrNil(t time.Time) *int64 {
	if t.IsZero() {
		return nil
	}
	return new(t.UnixMilli())
}

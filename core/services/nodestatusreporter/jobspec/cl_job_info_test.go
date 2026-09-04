package jobspec_test

import (
	"testing"
	"time"

	"github.com/pelletier/go-toml"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/beholder/beholdertest"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	evmtypes "github.com/smartcontractkit/chainlink-evm/pkg/types"
	commonv1 "github.com/smartcontractkit/chainlink-protos/node-platform/common/v1"

	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/nodestatusreporter/jobspec"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
)

func clJobInfoSampleJob() job.Job {
	streamID := uint32(42)
	return job.Job{
		ID:                7,
		Name:              null.StringFrom("my-ocr2-job"),
		Type:              job.OffchainReporting2,
		SchemaVersion:     1,
		ForwardingAllowed: true,
		StreamID:          &streamID,
		CreatedAt:         time.Date(2026, 7, 24, 10, 0, 0, 0, time.UTC),
		OCR2OracleSpec: &job.OCR2OracleSpec{
			Relay:         "evm",
			ChainID:       "1",
			PluginType:    commontypes.Median,
			ContractID:    "0xcccccccccccccccccccccccccccccccccccccccc",
			TransmitterID: null.StringFrom("0x1111111111111111111111111111111111111111"),
			RelayConfig: job.JSONConfig{
				"chainID":     "1",
				"sendingKeys": []any{"0x1111111111111111111111111111111111111111"},
			},
		},
		Pipeline: pipeline.Pipeline{Tasks: []pipeline.Task{
			&pipeline.ETHTxTask{From: "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
		}},
	}
}

// TestBuildCLJobInfo_EncodesFullSpecAsTOML is the load-bearing check: an
// arbitrary job must round-trip to TOML with no per-type code.
func TestBuildCLJobInfo_EncodesFullSpecAsTOML(t *testing.T) {
	jb := clJobInfoSampleJob()
	id := jobspec.NodeIdentity{CSAPublicKey: "csa", NodeVersion: "1.2.3", Hostname: "host-1"}

	info, err := jobspec.BuildCLJobInfo(jb, commonv1.CLJobInfoTrigger_CL_JOB_INFO_TRIGGER_CREATE, id, nil, time.Date(2026, 7, 24, 12, 0, 0, 0, time.UTC))
	require.NoError(t, err)

	require.Equal(t, "csa", info.CsaPublicKey)
	require.Equal(t, "1.2.3", info.NodeVersion)
	require.Equal(t, "host-1", info.Hostname)
	require.Equal(t, int32(7), info.JobId)
	require.Equal(t, "my-ocr2-job", info.Name)
	require.Equal(t, "offchainreporting2", info.JobType)
	require.Equal(t, uint32(1), info.SchemaVersion)
	require.True(t, info.ForwardingAllowed)
	require.NotNil(t, info.StreamId)
	require.Equal(t, uint32(42), *info.StreamId)
	require.Equal(t, commonv1.CLJobInfoTrigger_CL_JOB_INFO_TRIGGER_CREATE, info.Trigger)
	require.NotNil(t, info.Timestamp)
	require.Equal(t, time.Date(2026, 7, 24, 12, 0, 0, 0, time.UTC), info.Timestamp.AsTime())
	require.NotNil(t, info.CreatedAt)
	require.Equal(t, time.Date(2026, 7, 24, 10, 0, 0, 0, time.UTC), info.CreatedAt.AsTime())

	// spec_toml must be valid TOML and contain type-specific spec data.
	require.NotEmpty(t, info.SpecToml)
	var decoded map[string]any
	require.NoError(t, toml.Unmarshal([]byte(info.SpecToml), &decoded))
	require.Contains(t, info.SpecToml, "median")
	require.Contains(t, info.SpecToml, "0xcccccccccccccccccccccccccccccccccccccccc")
}

func TestBuildCLJobInfo_HandlesMultipleJobTypesGenerically(t *testing.T) {
	jobs := []job.Job{
		{Type: job.VRF, VRFSpec: &job.VRFSpec{
			EVMChainID:    sqlutil.NewI(4),
			FromAddresses: []evmtypes.EIP55Address{evmtypes.MustEIP55Address("0x6666666666666666666666666666666666666666")},
		}},
		{Type: job.BlockhashStore, BlockhashStoreSpec: &job.BlockhashStoreSpec{EVMChainID: sqlutil.NewI(5)}},
	}
	for _, jb := range jobs {
		info, err := jobspec.BuildCLJobInfo(jb, commonv1.CLJobInfoTrigger_CL_JOB_INFO_TRIGGER_HEARTBEAT, jobspec.NodeIdentity{}, nil, time.Now())
		require.NoErrorf(t, err, "job type %s should encode without per-type code", jb.Type)
		require.NotEmpty(t, info.SpecToml)
	}
}

// TestBuildCLJobInfo_CarriesJobDistributorProvenance covers the JD join key:
// remote_uuid is what links this event back to api.job.v1.Job.uuid.
func TestBuildCLJobInfo_CarriesJobDistributorProvenance(t *testing.T) {
	proposedAt := time.Date(2026, 7, 20, 9, 0, 0, 0, time.UTC)
	approvedAt := time.Date(2026, 7, 24, 10, 0, 0, 0, time.UTC)
	prop := &jobspec.JobProposal{
		FeedsManagerID: 3,
		RemoteUUID:     "6d7d9d1a-0d0f-4b3f-9a2f-2e4a1c0b8d55",
		SpecVersion:    2,
		ProposedAt:     proposedAt,
		ApprovedAt:     approvedAt,
	}

	info, err := jobspec.BuildCLJobInfo(clJobInfoSampleJob(), commonv1.CLJobInfoTrigger_CL_JOB_INFO_TRIGGER_CREATE, jobspec.NodeIdentity{}, prop, time.Now())
	require.NoError(t, err)

	require.NotNil(t, info.FeedsManagerId)
	require.Equal(t, int64(3), *info.FeedsManagerId)
	require.NotNil(t, info.RemoteUuid)
	require.Equal(t, "6d7d9d1a-0d0f-4b3f-9a2f-2e4a1c0b8d55", *info.RemoteUuid)
	require.NotNil(t, info.SpecVersion)
	require.Equal(t, int32(2), *info.SpecVersion)
	require.NotNil(t, info.ProposedAt)
	require.Equal(t, proposedAt, info.ProposedAt.AsTime())
	require.NotNil(t, info.ApprovedAt)
	require.Equal(t, approvedAt, info.ApprovedAt.AsTime())
}

// TestBuildCLJobInfo_UnmanagedJobHasNoProvenance: an unset feeds_manager_id is
// how a consumer tells a directly-created job from a JD-managed one.
func TestBuildCLJobInfo_UnmanagedJobHasNoProvenance(t *testing.T) {
	info, err := jobspec.BuildCLJobInfo(clJobInfoSampleJob(), commonv1.CLJobInfoTrigger_CL_JOB_INFO_TRIGGER_CREATE, jobspec.NodeIdentity{}, nil, time.Now())
	require.NoError(t, err)

	require.Nil(t, info.FeedsManagerId)
	require.Nil(t, info.RemoteUuid)
	require.Nil(t, info.SpecVersion)
	require.Nil(t, info.ProposedAt)
	require.Nil(t, info.ApprovedAt)
}

func TestEmitCLJobInfo_PublishesToBeholder(t *testing.T) {
	obs := beholdertest.NewObserver(t)

	info, err := jobspec.BuildCLJobInfo(clJobInfoSampleJob(), commonv1.CLJobInfoTrigger_CL_JOB_INFO_TRIGGER_CREATE, jobspec.NodeIdentity{CSAPublicKey: "csa"}, nil, time.Now())
	require.NoError(t, err)
	require.NoError(t, jobspec.EmitCLJobInfo(t.Context(), beholder.GetEmitter(), info))

	msgs := obs.Messages(t, beholder.AttrKeyEntity, jobspec.Entity)
	require.NotEmpty(t, msgs)

	msg := msgs[0]
	require.Equal(t, jobspec.Domain, msg.Attrs[beholder.AttrKeyDomain])
	require.Equal(t, jobspec.DataSchema, msg.Attrs[beholder.AttrKeyDataSchema])

	var payload commonv1.CLJobInfo
	require.NoError(t, proto.Unmarshal(msg.Body, &payload))
	require.Equal(t, "csa", payload.CsaPublicKey)
	require.Equal(t, "offchainreporting2", payload.JobType)
	require.NotEmpty(t, payload.SpecToml)
}

// TestBuildCLJobInfo_TimestampsRoundTripExactly guards the reason these fields
// are google.protobuf.Timestamp rather than RFC3339Nano strings: Go trims
// trailing zeros from the fractional seconds, so string-encoded times are
// variable-width and do not sort lexicographically in chronological order — a
// whole-second time sorts after every sub-second one in the same second.
func TestBuildCLJobInfo_TimestampsRoundTripExactly(t *testing.T) {
	for _, tc := range []struct {
		name string
		at   time.Time
	}{
		{"whole second", time.Date(2026, 7, 24, 10, 0, 0, 0, time.UTC)},
		{"tenth of a second", time.Date(2026, 7, 24, 10, 0, 0, 100000000, time.UTC)},
		{"sub-millisecond", time.Date(2026, 7, 24, 10, 0, 0, 123400000, time.UTC)},
		{"nanosecond", time.Date(2026, 7, 24, 10, 0, 0, 123456789, time.UTC)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			jb := clJobInfoSampleJob()
			jb.CreatedAt = tc.at

			info, err := jobspec.BuildCLJobInfo(jb, commonv1.CLJobInfoTrigger_CL_JOB_INFO_TRIGGER_HEARTBEAT, jobspec.NodeIdentity{}, nil, time.Now())
			require.NoError(t, err)
			require.NotNil(t, info.CreatedAt)
			require.Equal(t, tc.at, info.CreatedAt.AsTime())
		})
	}
}

// TestBuildCLJobInfo_ZeroTimeIsUnset: an absent time must be nil, not the epoch.
func TestBuildCLJobInfo_ZeroTimeIsUnset(t *testing.T) {
	jb := clJobInfoSampleJob()
	jb.CreatedAt = time.Time{}

	info, err := jobspec.BuildCLJobInfo(jb, commonv1.CLJobInfoTrigger_CL_JOB_INFO_TRIGGER_HEARTBEAT, jobspec.NodeIdentity{}, nil, time.Now())
	require.NoError(t, err)
	require.Nil(t, info.CreatedAt)
}

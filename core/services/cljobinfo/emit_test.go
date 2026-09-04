package cljobinfo_test

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

	"github.com/smartcontractkit/chainlink/v2/core/services/cljobinfo"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
)

func sampleJob() job.Job {
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
		VRFSpec: nil,
		Pipeline: pipeline.Pipeline{Tasks: []pipeline.Task{
			&pipeline.ETHTxTask{From: "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
		}},
	}
}

// TestBuild_EncodesFullSpecAsTOML is the load-bearing check: an arbitrary job
// must round-trip to TOML with no per-type code.
func TestBuild_EncodesFullSpecAsTOML(t *testing.T) {
	jb := sampleJob()
	id := cljobinfo.NodeIdentity{CSAPublicKey: "csa", NodeVersion: "1.2.3", Hostname: "host-1"}

	info, err := cljobinfo.Build(jb, commonv1.CLJobInfoTrigger_CL_JOB_INFO_TRIGGER_CREATE, id, time.Date(2026, 7, 24, 12, 0, 0, 0, time.UTC))
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
	require.NotEmpty(t, info.Timestamp)

	// spec_toml must be valid TOML and contain type-specific spec data.
	require.NotEmpty(t, info.SpecToml)
	var decoded map[string]any
	require.NoError(t, toml.Unmarshal([]byte(info.SpecToml), &decoded))
	require.Contains(t, info.SpecToml, "median")
	require.Contains(t, info.SpecToml, "0xcccccccccccccccccccccccccccccccccccccccc")
}

func TestBuild_HandlesMultipleJobTypesGenerically(t *testing.T) {
	jobs := []job.Job{
		{Type: job.VRF, VRFSpec: &job.VRFSpec{
			EVMChainID:    sqlutil.NewI(4),
			FromAddresses: []evmtypes.EIP55Address{evmtypes.MustEIP55Address("0x6666666666666666666666666666666666666666")},
		}},
		{Type: job.BlockhashStore, BlockhashStoreSpec: &job.BlockhashStoreSpec{EVMChainID: sqlutil.NewI(5)}},
	}
	for _, jb := range jobs {
		info, err := cljobinfo.Build(jb, commonv1.CLJobInfoTrigger_CL_JOB_INFO_TRIGGER_HEARTBEAT, cljobinfo.NodeIdentity{}, time.Now())
		require.NoErrorf(t, err, "job type %s should encode without per-type code", jb.Type)
		require.NotEmpty(t, info.SpecToml)
	}
}

func TestEmit_PublishesToBeholder(t *testing.T) {
	obs := beholdertest.NewObserver(t)

	info, err := cljobinfo.Build(sampleJob(), commonv1.CLJobInfoTrigger_CL_JOB_INFO_TRIGGER_CREATE, cljobinfo.NodeIdentity{CSAPublicKey: "csa"}, time.Now())
	require.NoError(t, err)
	require.NoError(t, cljobinfo.Emit(t.Context(), beholder.GetEmitter(), info))

	msgs := obs.Messages(t, beholder.AttrKeyEntity, cljobinfo.Entity)
	require.NotEmpty(t, msgs)

	msg := msgs[0]
	require.Equal(t, cljobinfo.Domain, msg.Attrs[beholder.AttrKeyDomain])
	require.Equal(t, cljobinfo.DataSchema, msg.Attrs[beholder.AttrKeyDataSchema])

	var payload commonv1.CLJobInfo
	require.NoError(t, proto.Unmarshal(msg.Body, &payload))
	require.Equal(t, "csa", payload.CsaPublicKey)
	require.Equal(t, "offchainreporting2", payload.JobType)
	require.NotEmpty(t, payload.SpecToml)
}

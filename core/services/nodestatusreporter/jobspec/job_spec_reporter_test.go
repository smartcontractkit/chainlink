package jobspec_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/beholder/beholdertest"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	commonv1 "github.com/smartcontractkit/chainlink-protos/node-platform/common/v1"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/feeds"
	feedsmocks "github.com/smartcontractkit/chainlink/v2/core/services/feeds/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	jobmocks "github.com/smartcontractkit/chainlink/v2/core/services/job/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/nodestatusreporter/jobspec"
	"github.com/smartcontractkit/chainlink/v2/core/services/nodestatusreporter/jobspec/events"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
)

// stubConfig implements config.JobSpecReporter for tests.
type stubConfig struct {
	enabled                bool
	pollingInterval        time.Duration
	enabledOCR2PluginTypes []string
}

func (s *stubConfig) Enabled() bool                    { return s.enabled }
func (s *stubConfig) PollingInterval() time.Duration   { return s.pollingInterval }
func (s *stubConfig) EnabledOCR2PluginTypes() []string { return s.enabledOCR2PluginTypes }

func defaultConfig() *stubConfig {
	return &stubConfig{
		enabled:                true,
		pollingInterval:        time.Hour,
		enabledOCR2PluginTypes: []string{"median"},
	}
}

func makeMedianJob() job.Job {
	return job.Job{
		ID:            1,
		ExternalJobID: uuid.New(),
		Name:          null.StringFrom("test-median-job"),
		Type:          job.OffchainReporting2,
		SchemaVersion: 1,
		PipelineSpec: &pipeline.Spec{
			ID:           10,
			DotDagSource: `ds1 [type=bridge name="my-bridge"]`,
		},
		Pipeline: pipeline.Pipeline{
			Tasks: []pipeline.Task{
				&pipeline.BridgeTask{
					BaseTask: pipeline.NewBaseTask(0, "ds1", nil, nil, 0),
					Name:     "my-bridge",
				},
			},
		},
		OCR2OracleSpec: &job.OCR2OracleSpec{
			ID:                          1,
			ContractID:                  "0x1234567890abcdef",
			Relay:                       "evm",
			ChainID:                     "1",
			PluginType:                  commontypes.Median,
			TransmitterID:               null.StringFrom("0x1111111111111111111111111111111111111111"),
			RelayConfig:                 job.JSONConfig{"chainID": "1"},
			PluginConfig:                job.JSONConfig{"juelsPerFeeCoinSource": `ds1 [type=http method=GET url="https://example.com"]`},
			OnchainSigningStrategy:      job.JSONConfig{},
			P2PV2Bootstrappers:          []string{"12D3KooW@host:6688"},
			ContractConfigConfirmations: 1,
		},
		CreatedAt: time.Now(),
	}
}

func makeNonMedianOCR2Job() job.Job {
	jb := makeMedianJob()
	jb.ID = 2
	jb.ExternalJobID = uuid.New()
	jb.Name = null.StringFrom("test-non-median-job")
	jb.OCR2OracleSpec = &job.OCR2OracleSpec{
		ID:                     2,
		ContractID:             "0xabcdef1234567890",
		Relay:                  "evm",
		ChainID:                "1",
		PluginType:             commontypes.LLO,
		TransmitterID:          null.StringFrom("0x2222222222222222222222222222222222222222"),
		RelayConfig:            job.JSONConfig{"chainID": "1"},
		PluginConfig:           job.JSONConfig{},
		OnchainSigningStrategy: job.JSONConfig{},
	}
	return jb
}

func makeNonOCR2Job() job.Job {
	return job.Job{
		ID:            3,
		ExternalJobID: uuid.New(),
		Name:          null.StringFrom("test-cron-job"),
		Type:          job.Cron,
		SchemaVersion: 1,
		PipelineSpec:  &pipeline.Spec{ID: 30, DotDagSource: ""},
		Pipeline:      pipeline.Pipeline{},
		CreatedAt:     time.Now(),
	}
}

// newTestReporter returns a Service wired to a per-test recorder emitter, so
// its emissions can be asserted without touching process-global beholder state.
func newTestReporter(t *testing.T, cfg *stubConfig, feedsORM feeds.ORM) (*jobspec.Service, *beholdertest.Recorder) {
	t.Helper()
	spawner := jobmocks.NewSpawner(t)
	rec := beholdertest.NewRecorder()
	return jobspec.NewJobSpecReporter(cfg, spawner, feedsORM, rec, "csa-key", "1.0.0", "test-host", logger.TestLogger(t)), rec
}

// newFeedsORMWithoutProposal returns a feeds ORM mock that behaves as if the
// given job was created outside of the Feeds Manager.
func newFeedsORMWithoutProposal(t *testing.T, jb job.Job) *feedsmocks.ORM {
	t.Helper()
	feedsORM := feedsmocks.NewORM(t)
	feedsORM.On("GetJobProposalByExternalJobID", mock.Anything, jb.ExternalJobID).Return(nil, sql.ErrNoRows).Maybe()
	return feedsORM
}

func requireSingleJobSpecEvent(t *testing.T, rec *beholdertest.Recorder) *events.JobSpecEvent {
	t.Helper()

	msgs := rec.Messages(t, "beholder_entity", events.ProtoPkg+"."+events.JobSpecEventEntity)
	require.Len(t, msgs, 1)

	ev := new(events.JobSpecEvent)
	require.NoError(t, proto.Unmarshal(msgs[0].Body, ev))
	return ev
}

func TestShouldEmit_DefaultConfig(t *testing.T) {
	t.Parallel()

	svc, _ := newTestReporter(t, defaultConfig(), nil)

	median := makeMedianJob()
	nonMedian := makeNonMedianOCR2Job()
	nonOCR2 := makeNonOCR2Job()

	cases := []struct {
		name string
		jb   *job.Job
		want bool
	}{
		{"median OCR2 job emits", &median, true},
		{"non-median OCR2 job skipped", &nonMedian, false},
		{"non-OCR2 job skipped", &nonOCR2, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.want, svc.ShouldEmit(tc.jb))
		})
	}
}

func TestShouldEmit_NoOCR2Types(t *testing.T) {
	t.Parallel()

	cfg := defaultConfig()
	cfg.enabledOCR2PluginTypes = []string{} // empty allowlist = disable all

	svc, _ := newTestReporter(t, cfg, nil)

	median := makeMedianJob()
	nonMedian := makeNonMedianOCR2Job()
	nonOCR2 := makeNonOCR2Job()

	assert.False(t, svc.ShouldEmit(&median))
	assert.False(t, svc.ShouldEmit(&nonMedian))
	assert.False(t, svc.ShouldEmit(&nonOCR2))
}

func TestShouldEmit_AllOCR2Types(t *testing.T) {
	t.Parallel()

	cfg := defaultConfig()
	cfg.enabledOCR2PluginTypes = []string{"all"}

	svc, _ := newTestReporter(t, cfg, nil)

	median := makeMedianJob()
	nonMedian := makeNonMedianOCR2Job()
	nonOCR2 := makeNonOCR2Job()

	assert.True(t, svc.ShouldEmit(&median))
	assert.True(t, svc.ShouldEmit(&nonMedian))
	assert.False(t, svc.ShouldEmit(&nonOCR2))
}

func TestShouldEmit_NonOCR2Skipped(t *testing.T) {
	t.Parallel()

	svc, _ := newTestReporter(t, defaultConfig(), nil)

	median := makeMedianJob()
	nonOCR2 := makeNonOCR2Job()

	assert.False(t, svc.ShouldEmit(&nonOCR2))
	assert.True(t, svc.ShouldEmit(&median))
}

func TestBuildEvent_MedianJob(t *testing.T) {
	t.Parallel()

	jb := makeMedianJob()
	svc, rec := newTestReporter(t, defaultConfig(), newFeedsORMWithoutProposal(t, jb))

	err := svc.EmitForJob(context.Background(), jb, events.EmissionTrigger_EMISSION_TRIGGER_HEARTBEAT)
	require.NoError(t, err)

	ev := requireSingleJobSpecEvent(t, rec)
	assert.Equal(t, jb.ExternalJobID.String(), ev.ExternalJobId)
	assert.Equal(t, "test-median-job", ev.Name)
	assert.Equal(t, "offchainreporting2", ev.JobType)
	assert.Equal(t, jb.CreatedAt.Format(time.RFC3339Nano), ev.CreatedAt)
	assert.Equal(t, events.EmissionTrigger_EMISSION_TRIGGER_HEARTBEAT, ev.EmissionTrigger)
	assert.Equal(t, "csa-key", ev.CsaPublicKey)
	assert.Equal(t, "1.0.0", ev.NodeVersion)
	assert.Equal(t, "test-host", ev.Hostname)
	assert.Equal(t, []string{"my-bridge"}, ev.BridgeNames)
	require.NotNil(t, ev.Ocr2OracleSpec)
	assert.Equal(t, "0x1234567890abcdef", ev.Ocr2OracleSpec.GetContractId())
	assert.Equal(t, "evm", ev.Ocr2OracleSpec.Relay)
	assert.Equal(t, "median", ev.Ocr2OracleSpec.PluginType)
	require.NotNil(t, ev.Ocr2OracleSpec.MedianPluginConfig)
	assert.NotEmpty(t, ev.Ocr2OracleSpec.MedianPluginConfig.GetJuelsPerFeeCoinSource())
	require.NotNil(t, ev.Ocr2OracleSpec.EvmRelayConfig)
	assert.Equal(t, "1", ev.Ocr2OracleSpec.EvmRelayConfig.GetChainId())
	assert.Equal(t, "0x1111111111111111111111111111111111111111", ev.Ocr2OracleSpec.EvmRelayConfig.GetEffectiveTransmitterId())
}

func TestBuildEvent_MedianJobBridgeNamesFromPipelineSpec(t *testing.T) {
	t.Parallel()

	jb := makeMedianJob()
	jb.Pipeline = pipeline.Pipeline{}
	svc, rec := newTestReporter(t, defaultConfig(), newFeedsORMWithoutProposal(t, jb))

	err := svc.EmitForJob(context.Background(), jb, events.EmissionTrigger_EMISSION_TRIGGER_HEARTBEAT)
	require.NoError(t, err)

	ev := requireSingleJobSpecEvent(t, rec)
	assert.Equal(t, []string{"my-bridge"}, ev.BridgeNames)
}

func TestBuildEvent_MedianJobNumericRelayConfigChainID(t *testing.T) {
	t.Parallel()

	jb := makeMedianJob()
	jb.OCR2OracleSpec.ChainID = ""
	jb.OCR2OracleSpec.RelayConfig["chainID"] = int64(11155111)
	svc, rec := newTestReporter(t, defaultConfig(), newFeedsORMWithoutProposal(t, jb))

	err := svc.EmitForJob(context.Background(), jb, events.EmissionTrigger_EMISSION_TRIGGER_HEARTBEAT)
	require.NoError(t, err)

	ev := requireSingleJobSpecEvent(t, rec)
	require.NotNil(t, ev.Ocr2OracleSpec)
	require.NotNil(t, ev.Ocr2OracleSpec.EvmRelayConfig)
	assert.Equal(t, "11155111", ev.Ocr2OracleSpec.EvmRelayConfig.GetChainId())
}

func TestBuildEvent_EVMRelayConfigEmitsExplicitFalseBooleans(t *testing.T) {
	t.Parallel()

	jb := makeMedianJob()
	jb.OCR2OracleSpec.RelayConfig["enableDualTransmission"] = false
	jb.OCR2OracleSpec.RelayConfig["enableTriggerCapability"] = false
	svc, rec := newTestReporter(t, defaultConfig(), newFeedsORMWithoutProposal(t, jb))

	err := svc.EmitForJob(context.Background(), jb, events.EmissionTrigger_EMISSION_TRIGGER_HEARTBEAT)
	require.NoError(t, err)

	ev := requireSingleJobSpecEvent(t, rec)
	require.NotNil(t, ev.Ocr2OracleSpec)
	require.NotNil(t, ev.Ocr2OracleSpec.EvmRelayConfig)
	assert.False(t, ev.Ocr2OracleSpec.EvmRelayConfig.GetEnableDualTransmission())
	assert.False(t, ev.Ocr2OracleSpec.EvmRelayConfig.GetEnableTriggerCapability())
	assert.NotNil(t, ev.Ocr2OracleSpec.EvmRelayConfig.EnableDualTransmission)
	assert.NotNil(t, ev.Ocr2OracleSpec.EvmRelayConfig.EnableTriggerCapability)
}

func TestBuildEvent_NonMedianOCR2Job(t *testing.T) {
	t.Parallel()

	jb := makeNonMedianOCR2Job()
	svc, rec := newTestReporter(t, defaultConfig(), newFeedsORMWithoutProposal(t, jb))

	err := svc.EmitForJob(context.Background(), jb, events.EmissionTrigger_EMISSION_TRIGGER_CREATE)
	require.NoError(t, err)

	ev := requireSingleJobSpecEvent(t, rec)
	require.NotNil(t, ev.Ocr2OracleSpec)
	assert.Equal(t, "llo", ev.Ocr2OracleSpec.PluginType)
	assert.Nil(t, ev.Ocr2OracleSpec.MedianPluginConfig)
	assert.NotEmpty(t, ev.Ocr2OracleSpec.RelayConfigJson)
}

func TestBuildEvent_NonOCR2Job(t *testing.T) {
	t.Parallel()

	svc, rec := newTestReporter(t, defaultConfig(), nil)

	jb := makeNonOCR2Job()
	err := svc.EmitForJob(context.Background(), jb, events.EmissionTrigger_EMISSION_TRIGGER_HEARTBEAT)
	require.ErrorContains(t, err, "unsupported job type")

	msgs := rec.Messages(t, "beholder_entity", events.ProtoPkg+"."+events.JobSpecEventEntity)
	require.Empty(t, msgs)
}

func TestAfterJobStarted_EmitsCreate(t *testing.T) {
	t.Parallel()

	jb := makeMedianJob()
	svc, rec := newTestReporter(t, defaultConfig(), newFeedsORMWithoutProposal(t, jb))
	svc.AfterJobStarted(context.Background(), jb)

	ev := requireSingleJobSpecEvent(t, rec)
	assert.Equal(t, events.EmissionTrigger_EMISSION_TRIGGER_CREATE, ev.EmissionTrigger)
}

func TestAfterJobStopped_EmitsDelete(t *testing.T) {
	t.Parallel()

	jb := makeMedianJob()
	svc, rec := newTestReporter(t, defaultConfig(), newFeedsORMWithoutProposal(t, jb))
	svc.AfterJobStopped(context.Background(), jb)

	ev := requireSingleJobSpecEvent(t, rec)
	assert.Equal(t, events.EmissionTrigger_EMISSION_TRIGGER_DELETE, ev.EmissionTrigger)
}

func TestAfterJobStarted_SkippedWhenGateFails(t *testing.T) {
	t.Parallel()

	// default config only allows median, so a non-OCR2 job should be skipped
	svc, rec := newTestReporter(t, defaultConfig(), nil)
	svc.AfterJobStarted(context.Background(), makeNonOCR2Job())

	msgs := rec.Messages(t, "beholder_entity", events.ProtoPkg+"."+events.JobSpecEventEntity)
	require.Empty(t, msgs)
}

func TestBuildEvent_ProposalLifecycle(t *testing.T) {
	t.Parallel()

	feedsORM := feedsmocks.NewORM(t)

	jb := makeMedianJob()
	proposedAt := time.Now().Add(-5 * time.Minute)
	approvedAt := time.Now().Add(-2 * time.Minute)

	prop := &feeds.JobProposal{
		ID:             100,
		FeedsManagerID: 7,
		RemoteUUID:     uuid.New(),
	}
	spec := &feeds.JobProposalSpec{
		ID:              200,
		Version:         3,
		CreatedAt:       proposedAt,
		StatusUpdatedAt: approvedAt,
	}

	feedsORM.On("GetJobProposalByExternalJobID", mock.Anything, jb.ExternalJobID).Return(prop, nil)
	feedsORM.On("GetApprovedSpec", mock.Anything, prop.ID).Return(spec, nil)

	svc, rec := newTestReporter(t, defaultConfig(), feedsORM)
	err := svc.EmitForJob(context.Background(), jb, events.EmissionTrigger_EMISSION_TRIGGER_HEARTBEAT)
	require.NoError(t, err)

	ev := requireSingleJobSpecEvent(t, rec)
	assert.Equal(t, int64(7), ev.FeedsManagerId)
	assert.Equal(t, prop.RemoteUUID.String(), ev.RemoteUuid)
	assert.Equal(t, int32(3), ev.SpecVersion)
	assert.Equal(t, proposedAt.Format(time.RFC3339Nano), ev.ProposedAt)
	assert.Equal(t, approvedAt.Format(time.RFC3339Nano), ev.ApprovedAt)
	assert.InDelta(t, approvedAt.Sub(proposedAt).Seconds(), ev.AcceptLatencySeconds, 1.0)
}

// Pins the split: CLJobInfo needs no per-node opt-in, the legacy OCR2 track
// stays behind JobSpecReporter.Enabled.
func TestAfterJobStarted_CLJobInfoIgnoresEnabledGate(t *testing.T) {
	t.Parallel()

	cfg := defaultConfig()
	cfg.enabled = false // legacy track off

	jb := makeMedianJob()
	reporter, rec := newTestReporter(t, cfg, newFeedsORMWithoutProposal(t, jb))
	reporter.AfterJobStarted(t.Context(), jb)

	clMsgs := rec.Messages(t, beholder.AttrKeyEntity, jobspec.Entity)
	require.Len(t, clMsgs, 1, "CLJobInfo must be emitted even with JobSpecReporter disabled")

	var payload commonv1.CLJobInfo
	require.NoError(t, proto.Unmarshal(clMsgs[0].Body, &payload))
	require.Equal(t, commonv1.CLJobInfoTrigger_CL_JOB_INFO_TRIGGER_CREATE, payload.Trigger)
	require.NotEmpty(t, payload.SpecToml)

	legacy := rec.Messages(t, "beholder_entity", events.ProtoPkg+"."+events.JobSpecEventEntity)
	require.Empty(t, legacy, "legacy JobSpecEvent must stay gated by Enabled")
}

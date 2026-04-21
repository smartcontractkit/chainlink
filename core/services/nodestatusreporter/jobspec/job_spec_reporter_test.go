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
	emitNonOCR2Jobs        bool
}

func (s *stubConfig) Enabled() bool                    { return s.enabled }
func (s *stubConfig) PollingInterval() time.Duration   { return s.pollingInterval }
func (s *stubConfig) EnabledOCR2PluginTypes() []string { return s.enabledOCR2PluginTypes }
func (s *stubConfig) EmitNonOCR2Jobs() bool            { return s.emitNonOCR2Jobs }

func defaultConfig() *stubConfig {
	return &stubConfig{
		enabled:                true,
		pollingInterval:        time.Hour,
		enabledOCR2PluginTypes: []string{"median"},
		emitNonOCR2Jobs:        false,
	}
}

// makeMedianJob creates a minimal OCR2 median job for testing.
func makeMedianJob() job.Job {
	extID := uuid.New()
	return job.Job{
		ID:            1,
		ExternalJobID: extID,
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
			RelayConfig:                 job.JSONConfig{"chainID": "1"},
			PluginConfig:                job.JSONConfig{"juelsPerFeeCoinSource": `ds1 [type=http method=GET url="https://example.com"]`},
			OnchainSigningStrategy:      job.JSONConfig{},
			P2PV2Bootstrappers:          []string{"12D3KooW@host:6688"},
			ContractConfigConfirmations: 1,
		},
		CreatedAt: time.Now(),
	}
}

// makeNonMedianOCR2Job creates a non-median OCR2 job.
func makeNonMedianOCR2Job() job.Job {
	jb := makeMedianJob()
	jb.ID = 2
	jb.ExternalJobID = uuid.New()
	jb.Name = null.StringFrom("test-keeper-job")
	jb.OCR2OracleSpec = &job.OCR2OracleSpec{
		ID:                     2,
		ContractID:             "0xabcdef1234567890",
		Relay:                  "evm",
		ChainID:                "1",
		PluginType:             commontypes.OCR2PluginType("ocr2keeper"),
		RelayConfig:            job.JSONConfig{"chainID": "1"},
		PluginConfig:           job.JSONConfig{},
		OnchainSigningStrategy: job.JSONConfig{},
	}
	return jb
}

// makeVRFJob creates a non-OCR2 job.
func makeVRFJob() job.Job {
	return job.Job{
		ID:            3,
		ExternalJobID: uuid.New(),
		Name:          null.StringFrom("test-vrf-job"),
		Type:          job.VRF,
		SchemaVersion: 1,
		PipelineSpec:  &pipeline.Spec{ID: 30, DotDagSource: ""},
		Pipeline:      pipeline.Pipeline{},
		CreatedAt:     time.Now(),
	}
}

// newTestReporter creates a Service wired to the current global beholder emitter.
// Call beholdertest.NewObserver(t) before this to set up the test emitter.
func newTestReporter(t *testing.T, cfg *stubConfig, spawner job.Spawner, feedsORM feeds.ORM, nodeInfo jobspec.NodeInfo) *jobspec.Service {
	t.Helper()
	return jobspec.NewJobSpecReporter(cfg, spawner, feedsORM, beholder.GetEmitter(), nodeInfo, logger.TestLogger(t))
}

// ── shouldEmit gate tests ──────────────────────────────────────────────────────

func TestShouldEmit_DefaultConfig(t *testing.T) {
	beholdertest.NewObserver(t)
	spawner := jobmocks.NewSpawner(t)
	spawner.On("RegisterListener", mock.Anything).Maybe()

	svc := newTestReporter(t, defaultConfig(), spawner, nil, jobspec.NodeInfo{})

	cases := []struct {
		name string
		jb   *job.Job
		want bool
	}{
		{"median OCR2 job emits", jobPtr(makeMedianJob()), true},
		{"non-median OCR2 job skipped", jobPtr(makeNonMedianOCR2Job()), false},
		{"non-OCR2 (VRF) job skipped", jobPtr(makeVRFJob()), false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, svc.ShouldEmit(tc.jb))
		})
	}
}

func TestShouldEmit_AllOCR2Types(t *testing.T) {
	beholdertest.NewObserver(t)
	cfg := defaultConfig()
	cfg.enabledOCR2PluginTypes = []string{} // empty = allow all

	spawner := jobmocks.NewSpawner(t)
	spawner.On("RegisterListener", mock.Anything).Maybe()
	svc := newTestReporter(t, cfg, spawner, nil, jobspec.NodeInfo{})

	assert.True(t, svc.ShouldEmit(jobPtr(makeMedianJob())))
	assert.True(t, svc.ShouldEmit(jobPtr(makeNonMedianOCR2Job())))
	assert.False(t, svc.ShouldEmit(jobPtr(makeVRFJob())))
}

func TestShouldEmit_NonOCR2Enabled(t *testing.T) {
	beholdertest.NewObserver(t)
	cfg := defaultConfig()
	cfg.emitNonOCR2Jobs = true

	spawner := jobmocks.NewSpawner(t)
	spawner.On("RegisterListener", mock.Anything).Maybe()
	svc := newTestReporter(t, cfg, spawner, nil, jobspec.NodeInfo{})

	assert.True(t, svc.ShouldEmit(jobPtr(makeVRFJob())))
	assert.True(t, svc.ShouldEmit(jobPtr(makeMedianJob())))
}

// ── conversion tests ──────────────────────────────────────────────────────────

func TestBuildEvent_MedianJob(t *testing.T) {
	observer := beholdertest.NewObserver(t)
	spawner := jobmocks.NewSpawner(t)
	spawner.On("RegisterListener", mock.Anything).Maybe()

	feedsORM := feedsmocks.NewORM(t)
	jb := makeMedianJob()
	feedsORM.On("GetJobProposalByExternalJobID", mock.Anything, jb.ExternalJobID).Return(nil, sql.ErrNoRows).Maybe()

	svc := newTestReporter(t, defaultConfig(), spawner, feedsORM, jobspec.NodeInfo{
		CSAPublicKey: "csa-key",
		NodeVersion:  "1.0.0",
		Hostname:     "test-host",
	})

	err := svc.EmitForJob(context.Background(), jb, "heartbeat")
	require.NoError(t, err)

	msgs := observer.Messages(t, "beholder_entity", events.ProtoPkg+"."+events.JobSpecEventEntity)
	require.Len(t, msgs, 1)

	var ev events.JobSpecEvent
	require.NoError(t, proto.Unmarshal(msgs[0].Body, &ev))

	assert.Equal(t, jb.ExternalJobID.String(), ev.ExternalJobId)
	assert.Equal(t, jb.ID, ev.InternalJobId)
	assert.Equal(t, "test-median-job", ev.Name)
	assert.Equal(t, "offchainreporting2", ev.JobType)
	assert.Equal(t, "heartbeat", ev.EmissionTrigger)
	assert.Equal(t, "csa-key", ev.CsaPublicKey)
	assert.Equal(t, "1.0.0", ev.NodeVersion)
	assert.Equal(t, "test-host", ev.Hostname)
	assert.Equal(t, []string{"my-bridge"}, ev.BridgeNames)
	require.NotNil(t, ev.Ocr2OracleSpec)
	assert.Equal(t, "evm", ev.Ocr2OracleSpec.Relay)
	assert.Equal(t, "median", ev.Ocr2OracleSpec.PluginType)
	require.NotNil(t, ev.Ocr2OracleSpec.MedianPluginConfig)
	assert.NotEmpty(t, ev.Ocr2OracleSpec.MedianPluginConfig.JuelsPerFeeCoinSource)
	require.NotNil(t, ev.Ocr2OracleSpec.EvmRelayConfig)
	assert.Equal(t, "1", ev.Ocr2OracleSpec.EvmRelayConfig.ChainId)
}

func TestBuildEvent_NonMedianOCR2Job(t *testing.T) {
	observer := beholdertest.NewObserver(t)
	spawner := jobmocks.NewSpawner(t)
	spawner.On("RegisterListener", mock.Anything).Maybe()
	feedsORM := feedsmocks.NewORM(t)
	jb := makeNonMedianOCR2Job()
	feedsORM.On("GetJobProposalByExternalJobID", mock.Anything, jb.ExternalJobID).Return(nil, sql.ErrNoRows).Maybe()

	svc := newTestReporter(t, defaultConfig(), spawner, feedsORM, jobspec.NodeInfo{})
	err := svc.EmitForJob(context.Background(), jb, "create")
	require.NoError(t, err)

	msgs := observer.Messages(t, "beholder_entity", events.ProtoPkg+"."+events.JobSpecEventEntity)
	require.Len(t, msgs, 1)

	var ev events.JobSpecEvent
	require.NoError(t, proto.Unmarshal(msgs[0].Body, &ev))

	require.NotNil(t, ev.Ocr2OracleSpec)
	assert.Equal(t, "ocr2keeper", ev.Ocr2OracleSpec.PluginType)
	assert.Nil(t, ev.Ocr2OracleSpec.MedianPluginConfig) // not median
	assert.NotEmpty(t, ev.Ocr2OracleSpec.RelayConfigJson)
}

func TestBuildEvent_NonOCR2Job(t *testing.T) {
	observer := beholdertest.NewObserver(t)
	spawner := jobmocks.NewSpawner(t)
	spawner.On("RegisterListener", mock.Anything).Maybe()

	svc := newTestReporter(t, defaultConfig(), spawner, nil, jobspec.NodeInfo{})

	jb := makeVRFJob()
	err := svc.EmitForJob(context.Background(), jb, "heartbeat")
	require.NoError(t, err)

	msgs := observer.Messages(t, "beholder_entity", events.ProtoPkg+"."+events.JobSpecEventEntity)
	require.Len(t, msgs, 1)

	var ev events.JobSpecEvent
	require.NoError(t, proto.Unmarshal(msgs[0].Body, &ev))

	assert.Equal(t, "vrf", ev.JobType)
	assert.Nil(t, ev.Ocr2OracleSpec) // no OCR2 spec
}

// ── OnJobStarted / OnJobStopped listener tests ────────────────────────────────

func TestOnJobStarted_EmitsCreate(t *testing.T) {
	observer := beholdertest.NewObserver(t)
	spawner := jobmocks.NewSpawner(t)
	spawner.On("RegisterListener", mock.Anything).Maybe()
	feedsORM := feedsmocks.NewORM(t)
	jb := makeMedianJob()
	feedsORM.On("GetJobProposalByExternalJobID", mock.Anything, jb.ExternalJobID).Return(nil, sql.ErrNoRows).Maybe()

	svc := newTestReporter(t, defaultConfig(), spawner, feedsORM, jobspec.NodeInfo{})
	svc.OnJobStarted(context.Background(), jb)

	msgs := observer.Messages(t, "beholder_entity", events.ProtoPkg+"."+events.JobSpecEventEntity)
	require.Len(t, msgs, 1)

	var ev events.JobSpecEvent
	require.NoError(t, proto.Unmarshal(msgs[0].Body, &ev))
	assert.Equal(t, "create", ev.EmissionTrigger)
}

func TestOnJobStopped_EmitsDelete(t *testing.T) {
	observer := beholdertest.NewObserver(t)
	spawner := jobmocks.NewSpawner(t)
	spawner.On("RegisterListener", mock.Anything).Maybe()
	feedsORM := feedsmocks.NewORM(t)
	jb := makeMedianJob()
	feedsORM.On("GetJobProposalByExternalJobID", mock.Anything, jb.ExternalJobID).Return(nil, sql.ErrNoRows).Maybe()

	svc := newTestReporter(t, defaultConfig(), spawner, feedsORM, jobspec.NodeInfo{})
	svc.OnJobStopped(context.Background(), jb)

	msgs := observer.Messages(t, "beholder_entity", events.ProtoPkg+"."+events.JobSpecEventEntity)
	require.Len(t, msgs, 1)

	var ev events.JobSpecEvent
	require.NoError(t, proto.Unmarshal(msgs[0].Body, &ev))
	assert.Equal(t, "delete", ev.EmissionTrigger)
}

func TestOnJobStarted_SkippedWhenGateFails(t *testing.T) {
	observer := beholdertest.NewObserver(t)
	spawner := jobmocks.NewSpawner(t)
	spawner.On("RegisterListener", mock.Anything).Maybe()

	// default config only allows median, so VRF should not emit
	svc := newTestReporter(t, defaultConfig(), spawner, nil, jobspec.NodeInfo{})
	svc.OnJobStarted(context.Background(), makeVRFJob())

	msgs := observer.Messages(t, "beholder_entity", events.ProtoPkg+"."+events.JobSpecEventEntity)
	require.Empty(t, msgs)
}

// ── proposal-latency test ─────────────────────────────────────────────────────

func TestBuildEvent_ProposalLifecycle(t *testing.T) {
	observer := beholdertest.NewObserver(t)
	spawner := jobmocks.NewSpawner(t)
	spawner.On("RegisterListener", mock.Anything).Maybe()
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

	svc := newTestReporter(t, defaultConfig(), spawner, feedsORM, jobspec.NodeInfo{})
	err := svc.EmitForJob(context.Background(), jb, "heartbeat")
	require.NoError(t, err)

	msgs := observer.Messages(t, "beholder_entity", events.ProtoPkg+"."+events.JobSpecEventEntity)
	require.Len(t, msgs, 1)

	var ev events.JobSpecEvent
	require.NoError(t, proto.Unmarshal(msgs[0].Body, &ev))

	assert.Equal(t, int64(7), ev.FeedsManagerId)
	assert.Equal(t, prop.RemoteUUID.String(), ev.RemoteUuid)
	assert.Equal(t, int32(3), ev.SpecVersion)
	assert.InDelta(t, approvedAt.Sub(proposedAt).Seconds(), ev.AcceptLatencySeconds, 1.0)
}

// ── helpers ───────────────────────────────────────────────────────────────────

func jobPtr(jb job.Job) *job.Job { return &jb }

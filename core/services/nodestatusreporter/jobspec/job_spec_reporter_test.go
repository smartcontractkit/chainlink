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

	"github.com/lib/pq"
	"github.com/smartcontractkit/chainlink-common/keystore/corekeys"
	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/beholder/beholdertest"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"

	evmtypes "github.com/smartcontractkit/chainlink-evm/pkg/types"

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

// newTestReporter returns a Service wired to the current global beholder emitter.
// The caller must set up the test emitter via beholdertest.NewObserver(t) first.
func newTestReporter(t *testing.T, cfg *stubConfig, feedsORM feeds.ORM) *jobspec.Service {
	t.Helper()
	spawner := jobmocks.NewSpawner(t)
	return jobspec.NewJobSpecReporter(cfg, spawner, feedsORM, beholder.GetEmitter(), "csa-key", "1.0.0", "test-host", logger.TestLogger(t))
}

// newFeedsORMWithoutProposal returns a feeds ORM mock that behaves as if the
// given job was created outside of the Feeds Manager.
func newFeedsORMWithoutProposal(t *testing.T, jb job.Job) *feedsmocks.ORM {
	t.Helper()
	feedsORM := feedsmocks.NewORM(t)
	feedsORM.On("GetJobProposalByExternalJobID", mock.Anything, jb.ExternalJobID).Return(nil, sql.ErrNoRows).Maybe()
	return feedsORM
}

func TestShouldEmit_DefaultConfig(t *testing.T) {
	beholdertest.NewObserver(t)
	svc := newTestReporter(t, defaultConfig(), nil)

	median := makeMedianJob()
	nonMedian := makeNonMedianOCR2Job()
	vrf := makeVRFJob()

	cases := []struct {
		name string
		jb   *job.Job
		want bool
	}{
		{"median OCR2 job emits", &median, true},
		{"non-median OCR2 job skipped", &nonMedian, false},
		{"non-OCR2 (VRF) job skipped", &vrf, false},
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
	cfg.enabledOCR2PluginTypes = []string{} // empty allowlist = all OCR2 types

	svc := newTestReporter(t, cfg, nil)

	median := makeMedianJob()
	nonMedian := makeNonMedianOCR2Job()
	vrf := makeVRFJob()

	assert.True(t, svc.ShouldEmit(&median))
	assert.True(t, svc.ShouldEmit(&nonMedian))
	assert.False(t, svc.ShouldEmit(&vrf))
}

func TestShouldEmit_NonOCR2Enabled(t *testing.T) {
	beholdertest.NewObserver(t)
	cfg := defaultConfig()
	cfg.emitNonOCR2Jobs = true

	svc := newTestReporter(t, cfg, nil)

	median := makeMedianJob()
	vrf := makeVRFJob()

	assert.True(t, svc.ShouldEmit(&vrf))
	assert.True(t, svc.ShouldEmit(&median))
}

func TestBuildEvent_MedianJob(t *testing.T) {
	observer := beholdertest.NewObserver(t)

	jb := makeMedianJob()
	svc := newTestReporter(t, defaultConfig(), newFeedsORMWithoutProposal(t, jb))

	err := svc.EmitForJob(context.Background(), jb, events.EmissionTrigger_EMISSION_TRIGGER_HEARTBEAT)
	require.NoError(t, err)

	msgs := observer.Messages(t, "beholder_entity", events.ProtoPkg+"."+events.JobSpecEventEntity)
	require.Len(t, msgs, 1)

	var ev events.JobSpecEvent
	require.NoError(t, proto.Unmarshal(msgs[0].Body, &ev))

	assert.Equal(t, jb.ExternalJobID.String(), ev.ExternalJobId)
	assert.Equal(t, jb.ID, ev.InternalJobId)
	assert.Equal(t, "test-median-job", ev.Name)
	assert.Equal(t, "offchainreporting2", ev.JobType)
	assert.Equal(t, events.EmissionTrigger_EMISSION_TRIGGER_HEARTBEAT, ev.EmissionTrigger)
	assert.Equal(t, "csa-key", ev.CsaPublicKey)
	assert.Equal(t, "1.0.0", ev.NodeVersion)
	assert.Equal(t, "test-host", ev.Hostname)
	assert.Equal(t, []string{"my-bridge"}, ev.BridgeNames)
	assert.Equal(t, "0x1234567890abcdef", ev.ContractAddress)
	assert.Equal(t, "1", ev.ChainId)
	require.NotNil(t, ev.Ocr2OracleSpec)
	assert.Equal(t, "evm", ev.Ocr2OracleSpec.Relay)
	assert.Equal(t, "median", ev.Ocr2OracleSpec.PluginType)
	require.NotNil(t, ev.Ocr2OracleSpec.MedianPluginConfig)
	assert.NotEmpty(t, ev.Ocr2OracleSpec.MedianPluginConfig.JuelsPerFeeCoinSource)
	require.NotNil(t, ev.Ocr2OracleSpec.EvmRelayConfig)
	assert.Equal(t, "1", ev.Ocr2OracleSpec.EvmRelayConfig.ChainId)
	assert.Nil(t, ev.Ocr1OracleSpec)
}

func TestBuildEvent_NonMedianOCR2Job(t *testing.T) {
	observer := beholdertest.NewObserver(t)

	jb := makeNonMedianOCR2Job()
	svc := newTestReporter(t, defaultConfig(), newFeedsORMWithoutProposal(t, jb))

	err := svc.EmitForJob(context.Background(), jb, events.EmissionTrigger_EMISSION_TRIGGER_CREATE)
	require.NoError(t, err)

	msgs := observer.Messages(t, "beholder_entity", events.ProtoPkg+"."+events.JobSpecEventEntity)
	require.Len(t, msgs, 1)

	var ev events.JobSpecEvent
	require.NoError(t, proto.Unmarshal(msgs[0].Body, &ev))

	require.NotNil(t, ev.Ocr2OracleSpec)
	assert.Equal(t, "ocr2keeper", ev.Ocr2OracleSpec.PluginType)
	assert.Nil(t, ev.Ocr2OracleSpec.MedianPluginConfig)
	assert.NotEmpty(t, ev.Ocr2OracleSpec.RelayConfigJson)
}

func TestBuildEvent_NonOCR2Job(t *testing.T) {
	observer := beholdertest.NewObserver(t)

	svc := newTestReporter(t, defaultConfig(), nil)

	jb := makeVRFJob()
	err := svc.EmitForJob(context.Background(), jb, events.EmissionTrigger_EMISSION_TRIGGER_HEARTBEAT)
	require.NoError(t, err)

	msgs := observer.Messages(t, "beholder_entity", events.ProtoPkg+"."+events.JobSpecEventEntity)
	require.Len(t, msgs, 1)

	var ev events.JobSpecEvent
	require.NoError(t, proto.Unmarshal(msgs[0].Body, &ev))

	assert.Equal(t, "vrf", ev.JobType)
	assert.Nil(t, ev.Ocr2OracleSpec)
}

func TestOnJobStarted_EmitsCreate(t *testing.T) {
	observer := beholdertest.NewObserver(t)

	jb := makeMedianJob()
	svc := newTestReporter(t, defaultConfig(), newFeedsORMWithoutProposal(t, jb))
	svc.OnJobStarted(context.Background(), jb)

	msgs := observer.Messages(t, "beholder_entity", events.ProtoPkg+"."+events.JobSpecEventEntity)
	require.Len(t, msgs, 1)

	var ev events.JobSpecEvent
	require.NoError(t, proto.Unmarshal(msgs[0].Body, &ev))
	assert.Equal(t, events.EmissionTrigger_EMISSION_TRIGGER_CREATE, ev.EmissionTrigger)
}

func TestOnJobStopped_EmitsDelete(t *testing.T) {
	observer := beholdertest.NewObserver(t)

	jb := makeMedianJob()
	svc := newTestReporter(t, defaultConfig(), newFeedsORMWithoutProposal(t, jb))
	svc.OnJobStopped(context.Background(), jb)

	msgs := observer.Messages(t, "beholder_entity", events.ProtoPkg+"."+events.JobSpecEventEntity)
	require.Len(t, msgs, 1)

	var ev events.JobSpecEvent
	require.NoError(t, proto.Unmarshal(msgs[0].Body, &ev))
	assert.Equal(t, events.EmissionTrigger_EMISSION_TRIGGER_DELETE, ev.EmissionTrigger)
}

func TestOnJobStarted_SkippedWhenGateFails(t *testing.T) {
	observer := beholdertest.NewObserver(t)

	// default config only allows median, so a VRF job should be skipped
	svc := newTestReporter(t, defaultConfig(), nil)
	svc.OnJobStarted(context.Background(), makeVRFJob())

	msgs := observer.Messages(t, "beholder_entity", events.ProtoPkg+"."+events.JobSpecEventEntity)
	require.Empty(t, msgs)
}

func TestBuildEvent_ProposalLifecycle(t *testing.T) {
	observer := beholdertest.NewObserver(t)
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

	svc := newTestReporter(t, defaultConfig(), feedsORM)
	err := svc.EmitForJob(context.Background(), jb, events.EmissionTrigger_EMISSION_TRIGGER_HEARTBEAT)
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

func TestBuildEvent_ContractFields_OCR1(t *testing.T) {
	observer := beholdertest.NewObserver(t)
	cfg := &stubConfig{
		enabled:         true,
		pollingInterval: time.Hour,
		emitNonOCR2Jobs: true,
	}
	svc := newTestReporter(t, cfg, nil)

	jb := makeOCR1Job()
	err := svc.EmitForJob(context.Background(), jb, events.EmissionTrigger_EMISSION_TRIGGER_HEARTBEAT)
	require.NoError(t, err)

	msgs := observer.Messages(t, "beholder_entity", events.ProtoPkg+"."+events.JobSpecEventEntity)
	require.Len(t, msgs, 1)

	var ev events.JobSpecEvent
	require.NoError(t, proto.Unmarshal(msgs[0].Body, &ev))

	// top-level contract identity
	assert.Equal(t, "0x9d9305445F404E925563d5D5EcC65C815Ec1655b", ev.ContractAddress)
	assert.Equal(t, "11155111", ev.ChainId)
	assert.Equal(t, "offchainreporting", ev.JobType)

	// OCR1 sub-message
	require.NotNil(t, ev.Ocr1OracleSpec)
	ocr1 := ev.Ocr1OracleSpec
	assert.Equal(t, int32(99), ocr1.SpecId)
	assert.Equal(t, []string{"12D3KooW@bootstrap:6688"}, ocr1.P2Pv2Bootstrappers)
	assert.False(t, ocr1.IsBootstrapPeer)
	assert.Equal(t, "0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20", ocr1.OcrKeyBundleId)
	assert.Equal(t, "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B", ocr1.TransmitterAddress)
	assert.InDelta(t, 30.0, ocr1.ObservationTimeoutSeconds, 0.001)
	assert.InDelta(t, 20.0, ocr1.BlockchainTimeoutSeconds, 0.001)
	assert.InDelta(t, 120.0, ocr1.ContractConfigTrackerSubscribeIntervalSeconds, 0.001)
	assert.InDelta(t, 60.0, ocr1.ContractConfigTrackerPollIntervalSeconds, 0.001)
	assert.Equal(t, uint32(3), ocr1.ContractConfigConfirmations)
	assert.InDelta(t, 10.0, ocr1.DatabaseTimeoutSeconds, 0.001)
	assert.InDelta(t, 1.0, ocr1.ObservationGracePeriodSeconds, 0.001)
	assert.InDelta(t, 5.0, ocr1.ContractTransmitterTransmitTimeoutSeconds, 0.001)
	assert.True(t, ocr1.CaptureEaTelemetry)
	assert.Equal(t, "2026-01-01T00:00:00Z", ocr1.SpecCreatedAt)
	assert.Equal(t, "2026-02-01T00:00:00Z", ocr1.SpecUpdatedAt)

	// OCR2 sub-message absent for OCR1 jobs
	assert.Nil(t, ev.Ocr2OracleSpec)
}

func makeOCR1Job() job.Job {
	keyHash, err := corekeys.Sha256HashFromHex("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20")
	if err != nil {
		panic(err)
	}
	transmitter := evmtypes.MustEIP55Address("0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B")
	specCreatedAt := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	specUpdatedAt := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)
	dbTimeout := sqlutil.Interval(10 * time.Second)
	gracePeriod := sqlutil.Interval(1 * time.Second)
	transmitTimeout := sqlutil.Interval(5 * time.Second)
	return job.Job{
		ID:            4,
		ExternalJobID: uuid.New(),
		Name:          null.StringFrom("test-ocr1-job"),
		Type:          job.OffchainReporting,
		SchemaVersion: 1,
		PipelineSpec:  &pipeline.Spec{ID: 40, DotDagSource: `ds1 [type=bridge name="bridge-gsr"]`},
		Pipeline: pipeline.Pipeline{
			Tasks: []pipeline.Task{
				&pipeline.BridgeTask{
					BaseTask: pipeline.NewBaseTask(0, "ds1", nil, nil, 0),
					Name:     "bridge-gsr",
				},
			},
		},
		OCROracleSpec: &job.OCROracleSpec{
			ID:                                     99,
			ContractAddress:                        evmtypes.MustEIP55Address("0x9d9305445F404E925563d5D5EcC65C815Ec1655b"),
			EVMChainID:                             sqlutil.NewI(11155111),
			P2PV2Bootstrappers:                     pq.StringArray{"12D3KooW@bootstrap:6688"},
			IsBootstrapPeer:                        false,
			EncryptedOCRKeyBundleID:                &keyHash,
			TransmitterAddress:                     &transmitter,
			ObservationTimeout:                     sqlutil.Interval(30 * time.Second),
			BlockchainTimeout:                      sqlutil.Interval(20 * time.Second),
			ContractConfigTrackerSubscribeInterval: sqlutil.Interval(2 * time.Minute),
			ContractConfigTrackerPollInterval:      sqlutil.Interval(1 * time.Minute),
			ContractConfigConfirmations:            3,
			DatabaseTimeout:                        &dbTimeout,
			ObservationGracePeriod:                 &gracePeriod,
			ContractTransmitterTransmitTimeout:     &transmitTimeout,
			CaptureEATelemetry:                     true,
			CreatedAt:                              specCreatedAt,
			UpdatedAt:                              specUpdatedAt,
		},
		CreatedAt: time.Now(),
	}
}

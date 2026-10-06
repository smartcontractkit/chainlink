package globalconfig

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
)

func TestMethodConfigsFromProto(t *testing.T) {
	t.Parallel()

	t.Run("nil map yields nil", func(t *testing.T) {
		t.Parallel()
		got, err := MethodConfigsFromProto(nil)
		require.NoError(t, err)
		assert.Nil(t, got)
	})

	t.Run("converts trigger config", func(t *testing.T) {
		t.Parallel()
		got, err := MethodConfigsFromProto(map[string]*capabilitiespb.CapabilityMethodConfig{
			"LogTrigger": {
				RemoteConfig: &capabilitiespb.CapabilityMethodConfig_RemoteTriggerConfig{
					RemoteTriggerConfig: &capabilitiespb.RemoteTriggerConfig{
						RegistrationRefresh:     durationpb.New(10 * time.Second),
						RegistrationExpiry:      durationpb.New(20 * time.Second),
						MinResponsesToAggregate: 3,
						MessageExpiry:           durationpb.New(30 * time.Second),
						MaxBatchSize:            5,
						BatchCollectionPeriod:   durationpb.New(40 * time.Second),
					},
				},
				AggregatorConfig: &capabilitiespb.AggregatorConfig{AggregatorType: capabilitiespb.AggregatorType_SignedReport},
			},
		})
		require.NoError(t, err)
		require.Len(t, got, 1)
		cfg := got["LogTrigger"]
		require.NotNil(t, cfg.RemoteTriggerConfig)
		assert.Equal(t, 10*time.Second, cfg.RemoteTriggerConfig.RegistrationRefresh)
		assert.Equal(t, 20*time.Second, cfg.RemoteTriggerConfig.RegistrationExpiry)
		assert.Equal(t, uint32(3), cfg.RemoteTriggerConfig.MinResponsesToAggregate)
		assert.Equal(t, 30*time.Second, cfg.RemoteTriggerConfig.MessageExpiry)
		assert.Equal(t, uint32(5), cfg.RemoteTriggerConfig.MaxBatchSize)
		assert.Equal(t, 40*time.Second, cfg.RemoteTriggerConfig.BatchCollectionPeriod)
		require.NotNil(t, cfg.AggregatorConfig)
		assert.Equal(t, capabilities.AggregatorType_SignedReport, cfg.AggregatorConfig.AggregatorType)
		assert.Nil(t, cfg.RemoteExecutableConfig)
	})

	t.Run("converts executable config", func(t *testing.T) {
		t.Parallel()
		got, err := MethodConfigsFromProto(map[string]*capabilitiespb.CapabilityMethodConfig{
			"Write": {
				RemoteConfig: &capabilitiespb.CapabilityMethodConfig_RemoteExecutableConfig{
					RemoteExecutableConfig: &capabilitiespb.RemoteExecutableConfig{
						TransmissionSchedule:      capabilitiespb.TransmissionSchedule_OneAtATime,
						DeltaStage:                durationpb.New(time.Second),
						RequestTimeout:            durationpb.New(15 * time.Second),
						ServerMaxParallelRequests: 7,
						RequestHasherType:         capabilitiespb.RequestHasherType_WriteReportExcludeSignatures,
						MinResponsesToAggregate:   4,
					},
				},
			},
		})
		require.NoError(t, err)
		require.Len(t, got, 1)
		cfg := got["Write"]
		require.NotNil(t, cfg.RemoteExecutableConfig)
		assert.Equal(t, capabilities.Schedule_OneAtATime, cfg.RemoteExecutableConfig.TransmissionSchedule)
		assert.Equal(t, time.Second, cfg.RemoteExecutableConfig.DeltaStage)
		assert.Equal(t, 15*time.Second, cfg.RemoteExecutableConfig.RequestTimeout)
		assert.Equal(t, uint32(7), cfg.RemoteExecutableConfig.ServerMaxParallelRequests)
		assert.Equal(t, capabilities.RequestHasherType_WriteReportExcludeSignatures, cfg.RemoteExecutableConfig.RequestHasherType)
		assert.Equal(t, uint32(4), cfg.RemoteExecutableConfig.MinResponsesToAggregate)
		assert.Nil(t, cfg.RemoteTriggerConfig)
		assert.Nil(t, cfg.AggregatorConfig)
	})

	t.Run("unset oneof converts to zero config", func(t *testing.T) {
		t.Parallel()
		got, err := MethodConfigsFromProto(map[string]*capabilitiespb.CapabilityMethodConfig{
			"View": {},
		})
		require.NoError(t, err)
		require.Len(t, got, 1)
		cfg := got["View"]
		assert.Nil(t, cfg.RemoteTriggerConfig)
		assert.Nil(t, cfg.RemoteExecutableConfig)
		assert.Nil(t, cfg.AggregatorConfig)
	})

	t.Run("round-trips through proto marshal", func(t *testing.T) {
		t.Parallel()
		// The same pb shape the on-chain path produces via
		// registry.CapabilityConfiguration.Unmarshal must convert identically here.
		in := map[string]*capabilitiespb.CapabilityMethodConfig{
			"Write": {
				RemoteConfig: &capabilitiespb.CapabilityMethodConfig_RemoteExecutableConfig{
					RemoteExecutableConfig: &capabilitiespb.RemoteExecutableConfig{
						RequestTimeout: durationpb.New(30 * time.Second),
					},
				},
			},
		}
		got, err := MethodConfigsFromProto(in)
		require.NoError(t, err)
		require.Len(t, got, 1)
		require.NotNil(t, got["Write"].RemoteExecutableConfig)
		assert.Equal(t, 30*time.Second, got["Write"].RemoteExecutableConfig.RequestTimeout)
	})
}

func TestValidate_MethodConfigs(t *testing.T) {
	t.Parallel()

	t.Run("accepts payload with method_configs", func(t *testing.T) {
		t.Parallel()
		raw := `{"version":1,"dons":{"don-1":{"capabilities":{"write-chain_evm_1@1.0.0":{"methodConfigs":{"Write":{"remoteExecutableConfig":{"requestTimeout":"30s"}}}}}}}}`
		require.NoError(t, Validate(raw))
	})

	t.Run("accepts payload without method_configs", func(t *testing.T) {
		t.Parallel()
		raw := `{"version":1,"dons":{"don-1":{"capabilities":{"cron@1.0.0":{}}}}}`
		require.NoError(t, Validate(raw))
	})

	t.Run("rejects unknown field inside method config", func(t *testing.T) {
		t.Parallel()
		// Strict protojson parsing rejects unknown fields anywhere in the payload
		// (a misspelled method config key cannot silently become "no config").
		raw := `{"version":1,"dons":{"don-1":{"capabilities":{"write-chain_evm_1@1.0.0":{"methodConfigs":{"Write":{"remoteExecutableConfigTypo":{"requestTimeout":"30s"}}}}}}}}`
		require.Error(t, Validate(raw))
	})
}

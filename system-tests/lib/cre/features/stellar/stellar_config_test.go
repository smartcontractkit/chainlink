package stellar

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
)

func TestResolveMethodConfigSettings(t *testing.T) {
	t.Parallel()

	t.Run("nil values yield defaults", func(t *testing.T) {
		t.Parallel()
		got, err := resolveMethodConfigSettings(nil)
		require.NoError(t, err)
		assert.Equal(t, defaultRequestTimeout, got.RequestTimeout)
		assert.Equal(t, defaultWriteDeltaStage, got.DeltaStage)
	})

	t.Run("string durations override defaults", func(t *testing.T) {
		t.Parallel()
		got, err := resolveMethodConfigSettings(map[string]any{
			requestTimeoutKey: "45s",
			deltaStageKey:     "20s",
		})
		require.NoError(t, err)
		assert.Equal(t, 45*time.Second, got.RequestTimeout)
		assert.Equal(t, 20*time.Second, got.DeltaStage)
	})

	t.Run("time.Duration values are accepted", func(t *testing.T) {
		t.Parallel()
		got, err := resolveMethodConfigSettings(map[string]any{
			deltaStageKey: 5 * time.Second,
		})
		require.NoError(t, err)
		assert.Equal(t, 5*time.Second, got.DeltaStage)
		// Unset key keeps its default.
		assert.Equal(t, defaultRequestTimeout, got.RequestTimeout)
	})

	t.Run("invalid duration string errors", func(t *testing.T) {
		t.Parallel()
		_, err := resolveMethodConfigSettings(map[string]any{deltaStageKey: "not-a-duration"})
		require.Error(t, err)
	})

	t.Run("wrong type errors", func(t *testing.T) {
		t.Parallel()
		_, err := resolveMethodConfigSettings(map[string]any{requestTimeoutKey: 5})
		require.Error(t, err)
	})
}

func TestMethodConfigs(t *testing.T) {
	t.Parallel()

	settings := methodConfigSettings{RequestTimeout: 30 * time.Second, DeltaStage: 15 * time.Second}
	m := methodConfigs(settings)

	require.Len(t, m, 3)
	for _, method := range []string{"GetLatestLedger", "ReadContract", "WriteReport"} {
		require.Contains(t, m, method, "missing method config for %s", method)
	}

	// WriteReport must transmit OneAtATime with the configured delta stage so the
	// capability-DON workers stagger their submissions instead of racing.
	write := m["WriteReport"].GetRemoteExecutableConfig()
	require.NotNil(t, write)
	assert.Equal(t, capabilitiespb.TransmissionSchedule_OneAtATime, write.GetTransmissionSchedule())
	assert.Equal(t, settings.DeltaStage, write.GetDeltaStage().AsDuration())
	assert.Equal(t, settings.RequestTimeout, write.GetRequestTimeout().AsDuration())

	read := m["ReadContract"].GetRemoteExecutableConfig()
	require.NotNil(t, read)
	assert.Equal(t, settings.RequestTimeout, read.GetRequestTimeout().AsDuration())
	assert.Equal(t, capabilitiespb.TransmissionSchedule_AllAtOnce, read.GetTransmissionSchedule())
}

func TestBuildWorkerConfigJSON(t *testing.T) {
	t.Parallel()

	settings := methodConfigSettings{RequestTimeout: 30 * time.Second, DeltaStage: 15 * time.Second}
	const (
		chainID   = "baefd734b8d3e48472cff83912375fedbc7573701912fe308af730180f97d74a"
		forwarder = "CBQY7TRRXPAZBLCON2PJ6EEZ4OXUY2HYDRERFVIE4MSBDSB2NVNZG42V"
	)

	raw, err := buildJobConfigJSON(chainID, forwarder, settings, false)
	require.NoError(t, err)

	var cfg map[string]any
	require.NoError(t, json.Unmarshal([]byte(raw), &cfg))
	assert.Equal(t, network, cfg["network"])
	assert.Equal(t, chainID, cfg["chainId"])
	assert.Equal(t, forwarder, cfg["creForwarderAddress"])
	// isLocal is omitempty in the shared config struct, so false is omitted (the
	// worker plugin defaults absent -> false).
	_, hasIsLocal := cfg["isLocal"]
	assert.False(t, hasIsLocal, "isLocal should be omitted when false")
	assert.InDelta(t, float64(settings.DeltaStage), cfg["deltaStage"], 0)

	rawLocal, err := buildJobConfigJSON(chainID, forwarder, settings, true)
	require.NoError(t, err)
	var localCfg map[string]any
	require.NoError(t, json.Unmarshal([]byte(rawLocal), &localCfg))
	assert.Equal(t, true, localCfg["isLocal"])
}

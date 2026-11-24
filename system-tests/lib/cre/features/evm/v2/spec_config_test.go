package v2

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
)

func TestBuildSpecConfig_SharedAcrossNodes(t *testing.T) {
	t.Parallel()

	specConfig, err := buildSpecConfig(1337, "0xforwarder", map[string]any{
		"LogTriggerPollInterval":          int64(1),
		"ReceiverGasMinimum":              int64(500),
		"LogTriggerSendChannelBufferSize": int64(100),
	})
	require.NoError(t, err)

	decoded, err := values.FromMapValueProto(specConfig)
	require.NoError(t, err)
	unwrapped, err := decoded.Unwrap()
	require.NoError(t, err)
	got, ok := unwrapped.(map[string]any)
	require.True(t, ok)

	assert.Equal(t, int64(1337), got["chainId"])
	assert.Equal(t, "evm", got["network"])
	assert.Equal(t, "0xforwarder", got["creForwarderAddress"])
	assert.Equal(t, int64(500), got["receiverGasMinimum"])
	assert.Equal(t, int64(1), got["logTriggerPollInterval"])
	assert.Equal(t, int64(100), got["logTriggerSendChannelBufferSize"])
	assert.Equal(t, int64(deltaStage), got["deltaStage"])
	assert.NotContains(t, got, "nodeAddress", "per-node values must not be in the shared spec_config")
	assert.NotContains(t, got, "logTriggerLimitQueryLogSize", "unset optional values must be omitted")
}

func TestBuildSpecConfig_ErrorsOnMissingRequiredValue(t *testing.T) {
	t.Parallel()

	_, err := buildSpecConfig(1337, "0xforwarder", map[string]any{})
	require.Error(t, err)
}

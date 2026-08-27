package v2

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/chain-capabilities/evm"
	cronpb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/triggers/cron"
)

func TestTriggerRegistrationLogFields_NilPayload(t *testing.T) {
	t.Parallel()

	result := triggerRegistrationLogFields(nil)
	assert.Nil(t, result)
}

func TestTriggerRegistrationLogFields_Cron(t *testing.T) {
	t.Parallel()

	payload, err := anypb.New(&cronpb.Config{Schedule: "*/5 * * * *"})
	require.NoError(t, err)

	result := triggerRegistrationLogFields(payload)
	require.Len(t, result, 2)
	assert.Equal(t, "schedule", result[0])
	assert.Equal(t, "*/5 * * * *", result[1])
}

func TestTriggerRegistrationLogFields_EVMLog(t *testing.T) {
	t.Parallel()

	addr := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14}
	sig := []byte{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99}

	payload, err := anypb.New(&evm.FilterLogTriggerRequest{
		Addresses: [][]byte{addr},
		Topics: []*evm.TopicValues{
			{Values: [][]byte{sig}},
		},
		Confidence: evm.ConfidenceLevel_CONFIDENCE_LEVEL_SAFE,
	})
	require.NoError(t, err)

	result := triggerRegistrationLogFields(payload)
	require.Len(t, result, 6, "expected 3 key-value pairs")

	assert.Equal(t, "addresses", result[0])
	addrs, ok := result[1].([]string)
	require.True(t, ok)
	require.Len(t, addrs, 1)
	assert.Equal(t, "0x0102030405060708090a0b0c0d0e0f1011121314", addrs[0])

	assert.Equal(t, "eventSignatures", result[2])
	sigs, ok := result[3].([]string)
	require.True(t, ok)
	require.Len(t, sigs, 1)
	assert.Equal(t, "0xaabbccddeeff00112233445566778899aabbccddeeff00112233445566778899", sigs[0])

	assert.Equal(t, "confidence", result[4])
	assert.Equal(t, "CONFIDENCE_LEVEL_SAFE", result[5])
}

func TestTriggerRegistrationLogFields_EVMLog_NoTopics(t *testing.T) {
	t.Parallel()

	payload, err := anypb.New(&evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14}},
		Confidence: evm.ConfidenceLevel_CONFIDENCE_LEVEL_LATEST,
	})
	require.NoError(t, err)

	result := triggerRegistrationLogFields(payload)
	require.Len(t, result, 6)

	assert.Equal(t, "eventSignatures", result[2])
	sigs, ok := result[3].([]string)
	require.True(t, ok)
	assert.Empty(t, sigs)

	assert.Equal(t, "confidence", result[4])
	assert.Equal(t, "CONFIDENCE_LEVEL_LATEST", result[5])
}

func TestTriggerRegistrationLogFields_UnknownPayloadType(t *testing.T) {
	t.Parallel()

	payload, err := anypb.New(&cronpb.LegacyPayload{})
	require.NoError(t, err)

	result := triggerRegistrationLogFields(payload)
	require.Len(t, result, 2)
	assert.Equal(t, "payloadType", result[0])
	msgName, ok := result[1].(string)
	require.True(t, ok)
	assert.Contains(t, msgName, "LegacyPayload")
}

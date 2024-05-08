package triggers

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/utils"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
)

const testID = "test-id-1"

func TestOnDemand(t *testing.T) {
	tr := NewOnDemand(logger.Test(t))
	ctx := tests.Context(t)

	req := capabilities.CapabilityRequest{
		Metadata: capabilities.RequestMetadata{
			WorkflowExecutionID: testID,
		},
	}

	ch, err := tr.RegisterTrigger(ctx, req)
	require.NoError(t, err)

	er := capabilities.CapabilityResponse{
		Value: &values.String{Underlying: testID},
	}

	err = tr.FanOutEvent(ctx, er)
	require.NoError(t, err)
	assert.Equal(t, er, <-ch)
}

func TestOnDemand_ChannelDoesntExist(t *testing.T) {
	tr := NewOnDemand(logger.Test(t))
	ctx := tests.Context(t)

	er := capabilities.CapabilityResponse{
		Value: &values.String{Underlying: testID},
	}
	err := tr.SendEvent(ctx, testID, er)
	assert.ErrorContains(t, err, "no registration")
}

func TestOnDemand_(t *testing.T) {
	tr := NewOnDemand(logger.Test(t))
	ctx := tests.Context(t)

	req := capabilities.CapabilityRequest{
		Metadata: capabilities.RequestMetadata{
			WorkflowID: "hello",
		},
	}

	callback, err := tr.RegisterTrigger(ctx, req)
	require.NoError(t, err)

	er := capabilities.CapabilityResponse{
		Value: &values.String{Underlying: testID},
	}
	err = tr.SendEvent(ctx, "hello", er)
	require.NoError(t, err)

	assert.Len(t, callback, 1)
	assert.Equal(t, er, <-callback)
}

func TestOnDemandTrigger_GenerateSchema(t *testing.T) {
	ts := NewOnDemand(logger.Nop())
	schema, err := ts.Schema()
	require.NotNil(t, schema)
	require.NoError(t, err)

	var shouldUpdate = true
	if shouldUpdate {
		err = os.WriteFile("./testdata/fixtures/ondemand/schema.json", []byte(schema), 0600)
		require.NoError(t, err)
	}

	fixture, err := os.ReadFile("./testdata/fixtures/ondemand/schema.json")
	require.NoError(t, err)

	utils.AssertJSONEqual(t, fixture, []byte(schema))
}

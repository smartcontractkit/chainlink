package triggers

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
)

const testID = "test-id-1"

var transformJSON = cmp.FilterValues(func(x, y []byte) bool {
	return json.Valid(x) && json.Valid(y)
}, cmp.Transformer("ParseJSON", func(in []byte) (out interface{}) {
	if err := json.Unmarshal(in, &out); err != nil {
		panic(err) // should never occur given previous filter to ensure valid JSON
	}
	return out
}))

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

	if diff := cmp.Diff(fixture, []byte(schema), transformJSON); diff != "" {
		t.Errorf("TestOnDemandTrigger_GenerateConfigSchema() mismatch (-want +got):\n%s", diff)
		t.FailNow()
	}
}

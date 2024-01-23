package capabilities

import (
	"context"
	"testing"

	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_CapabilityInfo(t *testing.T) {
	ci, err := NewCapabilityInfo(
		"capability-id",
		CapabilityTypeAction,
		"This is a mock capability that doesn't do anything.",
		"v1.0.0",
	)
	require.NoError(t, err)

	assert.Equal(t, ci, ci.Info())
}

func Test_CapabilityInfo_Invalid(t *testing.T) {
	_, err := NewCapabilityInfo(
		"capability-id",
		"test",
		"This is a mock capability that doesn't do anything.",
		"v1.0.0",
	)
	assert.ErrorContains(t, err, "invalid capability type")

	_, err = NewCapabilityInfo(
		"&!!!",
		CapabilityTypeAction,
		"This is a mock capability that doesn't do anything.",
		"v1.0.0",
	)
	assert.ErrorContains(t, err, "invalid id")

	_, err = NewCapabilityInfo(
		"mock-capability",
		CapabilityTypeAction,
		"This is a mock capability that doesn't do anything.",
		"hello",
	)
	assert.ErrorContains(t, err, "invalid version")
}

type mockCapabilityWithExecute struct {
	Executable
	Validatable
	CapabilityInfo
}

var mcwe = &mockCapabilityWithExecute{}

func (m *mockCapabilityWithExecute) Execute(ctx context.Context, callback chan values.Value, inputs values.Map) error {
	val, _ := values.NewString("hello")
	callback <- val

	close(callback)

	return nil
}

func Test_ExecuteSyncReturnSingleValue(t *testing.T) {
	config, _ := values.NewMap(map[string]interface{}{})
	val, err := ExecuteSync(nil, mcwe, *config)

	assert.NoError(t, err)
	assert.Equal(t, "hello", val.(*values.String).Underlying)
}

package capabilities

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/values"
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
		CapabilityType(5),
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
	ExecuteFn func(ctx context.Context, callback chan CapabilityResponse, inputs values.Map) error
}

func (m *mockCapabilityWithExecute) Execute(ctx context.Context, callback chan CapabilityResponse, inputs values.Map) error {
	return m.ExecuteFn(ctx, callback, inputs)
}

func Test_ExecuteSyncReturnSingleValue(t *testing.T) {
	mcwe := &mockCapabilityWithExecute{
		ExecuteFn: func(ctx context.Context, callback chan CapabilityResponse, inputs values.Map) error {
			val, _ := values.NewString("hello")
			callback <- CapabilityResponse{val, nil}

			close(callback)

			return nil
		},
	}
	config, _ := values.NewMap(map[string]interface{}{})
	val, err := ExecuteSync(context.Background(), mcwe, *config)

	assert.NoError(t, err, val)
	assert.Equal(t, "hello", val.(*values.String).Underlying)
}

func Test_ExecuteSyncReturnMultipleValues(t *testing.T) {
	es, _ := values.NewString("hello")
	expectedList := []values.Value{es, es, es}
	mcwe := &mockCapabilityWithExecute{
		ExecuteFn: func(ctx context.Context, callback chan CapabilityResponse, inputs values.Map) error {
			callback <- CapabilityResponse{es, nil}
			callback <- CapabilityResponse{es, nil}
			callback <- CapabilityResponse{es, nil}

			close(callback)

			return nil
		},
	}
	config, _ := values.NewMap(map[string]interface{}{})
	val, err := ExecuteSync(context.Background(), mcwe, *config)

	assert.NoError(t, err, val)
	assert.ElementsMatch(t, expectedList, val.(*values.List).Underlying)
}

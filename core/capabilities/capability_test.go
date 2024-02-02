package capabilities

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
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

	_, err = NewCapabilityInfo(
		strings.Repeat("n", 256),
		CapabilityTypeAction,
		"This is a mock capability that doesn't do anything.",
		"hello",
	)
	assert.ErrorContains(t, err, "exceeds max length 128")
}

type mockCapabilityWithExecute struct {
	CallbackExecutable
	CapabilityInfo
	ExecuteFn func(ctx context.Context, callback chan CapabilityResponse, req CapabilityRequest) error
}

func (m *mockCapabilityWithExecute) Execute(ctx context.Context, callback chan CapabilityResponse, req CapabilityRequest) error {
	return m.ExecuteFn(ctx, callback, req)
}

func Test_ExecuteSyncReturnSingleValue(t *testing.T) {
	mcwe := &mockCapabilityWithExecute{
		ExecuteFn: func(ctx context.Context, callback chan CapabilityResponse, req CapabilityRequest) error {
			val, _ := values.NewString("hello")
			callback <- CapabilityResponse{val, nil}

			close(callback)

			return nil
		},
	}
	req := CapabilityRequest{}
	val, err := ExecuteSync(testutils.Context(t), mcwe, req)

	assert.NoError(t, err, val)
	assert.Equal(t, "hello", val.Underlying[0].(*values.String).Underlying)
}

func Test_ExecuteSyncReturnMultipleValues(t *testing.T) {
	es, _ := values.NewString("hello")
	expectedList := []values.Value{es, es, es}
	mcwe := &mockCapabilityWithExecute{
		ExecuteFn: func(ctx context.Context, callback chan CapabilityResponse, req CapabilityRequest) error {
			callback <- CapabilityResponse{es, nil}
			callback <- CapabilityResponse{es, nil}
			callback <- CapabilityResponse{es, nil}

			close(callback)

			return nil
		},
	}
	req := CapabilityRequest{}
	val, err := ExecuteSync(testutils.Context(t), mcwe, req)

	assert.NoError(t, err, val)
	assert.ElementsMatch(t, expectedList, val.Underlying)
}

func Test_ExecuteSyncCapabilitySetupErrors(t *testing.T) {
	expectedErr := errors.New("something went wrong during setup")
	mcwe := &mockCapabilityWithExecute{
		ExecuteFn: func(ctx context.Context, callback chan CapabilityResponse, req CapabilityRequest) error {
			close(callback)
			return expectedErr
		},
	}
	req := CapabilityRequest{}
	val, err := ExecuteSync(testutils.Context(t), mcwe, req)

	assert.ErrorContains(t, err, expectedErr.Error())
	assert.Nil(t, val)
}

func Test_ExecuteSyncTimeout(t *testing.T) {
	ctxWithTimeout := testutils.Context(t)
	ctxWithTimeout, cancel := context.WithCancel(ctxWithTimeout)
	cancel()

	mcwe := &mockCapabilityWithExecute{
		ExecuteFn: func(ctx context.Context, callback chan CapabilityResponse, req CapabilityRequest) error {
			return nil
		},
	}
	req := CapabilityRequest{}
	val, err := ExecuteSync(ctxWithTimeout, mcwe, req)

	assert.ErrorContains(t, err, "context timed out. If you did not set a timeout, be aware that the default ExecuteSync timeout is")
	assert.Nil(t, val)
}

func Test_ExecuteSyncCapabilityErrors(t *testing.T) {
	expectedErr := errors.New("something went wrong during execution")
	mcwe := &mockCapabilityWithExecute{
		ExecuteFn: func(ctx context.Context, callback chan CapabilityResponse, req CapabilityRequest) error {
			callback <- CapabilityResponse{nil, expectedErr}

			close(callback)

			return nil
		},
	}
	req := CapabilityRequest{}
	val, err := ExecuteSync(testutils.Context(t), mcwe, req)

	assert.ErrorContains(t, err, expectedErr.Error())
	assert.Nil(t, val)
}

func Test_ExecuteSyncDoesNotReturnValues(t *testing.T) {
	mcwe := &mockCapabilityWithExecute{
		ExecuteFn: func(ctx context.Context, callback chan CapabilityResponse, req CapabilityRequest) error {
			close(callback)
			return nil
		},
	}
	req := CapabilityRequest{}
	val, err := ExecuteSync(testutils.Context(t), mcwe, req)

	assert.ErrorContains(t, err, "capability did not return any values")
	assert.Nil(t, val)
}

func Test_MustNewCapabilityInfo(t *testing.T) {
	assert.Panics(t, func() {
		MustNewCapabilityInfo(
			"capability-id",
			CapabilityTypeAction,
			"This is a mock capability that doesn't do anything.",
			"should-panic",
		)
	})
}

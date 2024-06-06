package capabilities

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
)

func Test_CapabilityInfo(t *testing.T) {
	ci, err := NewCapabilityInfo(
		"capability-id@1.0.0",
		CapabilityTypeAction,
		"This is a mock capability that doesn't do anything.",
	)
	require.NoError(t, err)

	gotCi, err := ci.Info(tests.Context(t))
	require.NoError(t, err)
	require.Equal(t, ci.Version(), "1.0.0")
	assert.Equal(t, ci, gotCi)

	ci, err = NewCapabilityInfo(
		// add build metadata and sha
		"capability-id@1.0.0+build.1234.sha-5678",
		CapabilityTypeAction,
		"This is a mock capability that doesn't do anything.",
	)
	require.NoError(t, err)

	gotCi, err = ci.Info(tests.Context(t))
	require.NoError(t, err)
	require.Equal(t, ci.Version(), "1.0.0+build.1234.sha-5678")
	assert.Equal(t, ci, gotCi)

	// prerelease
	ci, err = NewCapabilityInfo(
		"capability-id@1.0.0-beta",
		CapabilityTypeAction,
		"This is a mock capability that doesn't do anything.",
	)
	require.NoError(t, err)

	gotCi, err = ci.Info(tests.Context(t))
	require.NoError(t, err)
	require.Equal(t, ci.Version(), "1.0.0-beta")
	assert.Equal(t, ci, gotCi)
}

func Test_CapabilityInfo_Invalid(t *testing.T) {
	_, err := NewCapabilityInfo(
		"capability-id@2.0.0",
		CapabilityType(5),
		"This is a mock capability that doesn't do anything.",
	)
	assert.ErrorContains(t, err, "invalid capability type")

	_, err = NewCapabilityInfo(
		"&!!!",
		CapabilityTypeAction,
		"This is a mock capability that doesn't do anything.",
	)
	assert.ErrorContains(t, err, "invalid id")

	_, err = NewCapabilityInfo(
		"mock-capability@v1.0.0",
		CapabilityTypeAction,
		"This is a mock capability that doesn't do anything.",
	)
	assert.ErrorContains(t, err, "invalid id")

	_, err = NewCapabilityInfo(
		"mock-capability@1.0",
		CapabilityTypeAction,
		"This is a mock capability that doesn't do anything.",
	)
	assert.ErrorContains(t, err, "invalid id")

	_, err = NewCapabilityInfo(
		"mock-capability@1",
		CapabilityTypeAction,
		"This is a mock capability that doesn't do anything.",
	)

	assert.ErrorContains(t, err, "invalid id")
	_, err = NewCapabilityInfo(
		strings.Repeat("n", 256),
		CapabilityTypeAction,
		"This is a mock capability that doesn't do anything.",
	)
	assert.ErrorContains(t, err, "exceeds max length 128")
}

type mockCapabilityWithExecute struct {
	CallbackExecutable
	CapabilityInfo
	ExecuteFn func(ctx context.Context, req CapabilityRequest) (<-chan CapabilityResponse, error)
}

func (m *mockCapabilityWithExecute) Execute(ctx context.Context, req CapabilityRequest) (<-chan CapabilityResponse, error) {
	return m.ExecuteFn(ctx, req)
}

func Test_ExecuteSyncReturnSingleValue(t *testing.T) {
	mcwe := &mockCapabilityWithExecute{
		ExecuteFn: func(ctx context.Context, req CapabilityRequest) (<-chan CapabilityResponse, error) {
			ch := make(chan CapabilityResponse, 10)

			val := values.NewString("hello")
			ch <- CapabilityResponse{val, nil}

			close(ch)

			return ch, nil
		},
	}
	req := CapabilityRequest{}
	val, err := ExecuteSync(tests.Context(t), mcwe, req)

	assert.NoError(t, err, val)
	assert.Equal(t, "hello", val.Underlying[0].(*values.String).Underlying)
}

func Test_ExecuteSyncReturnMultipleValues(t *testing.T) {
	es := values.NewString("hello")
	expectedList := []values.Value{es, es, es}
	mcwe := &mockCapabilityWithExecute{
		ExecuteFn: func(ctx context.Context, req CapabilityRequest) (<-chan CapabilityResponse, error) {
			ch := make(chan CapabilityResponse, 10)

			ch <- CapabilityResponse{es, nil}
			ch <- CapabilityResponse{es, nil}
			ch <- CapabilityResponse{es, nil}

			close(ch)

			return ch, nil
		},
	}
	req := CapabilityRequest{}
	val, err := ExecuteSync(tests.Context(t), mcwe, req)

	assert.NoError(t, err, val)
	assert.ElementsMatch(t, expectedList, val.Underlying)
}

func Test_ExecuteSyncCapabilitySetupErrors(t *testing.T) {
	expectedErr := errors.New("something went wrong during setup")
	mcwe := &mockCapabilityWithExecute{
		ExecuteFn: func(ctx context.Context, req CapabilityRequest) (<-chan CapabilityResponse, error) {
			return nil, expectedErr
		},
	}
	req := CapabilityRequest{}
	val, err := ExecuteSync(tests.Context(t), mcwe, req)

	assert.ErrorContains(t, err, expectedErr.Error())
	assert.Nil(t, val)
}

func Test_ExecuteSyncTimeout(t *testing.T) {
	ctxWithTimeout := tests.Context(t)
	ctxWithTimeout, cancel := context.WithCancel(ctxWithTimeout)
	cancel()

	mcwe := &mockCapabilityWithExecute{
		ExecuteFn: func(ctx context.Context, req CapabilityRequest) (<-chan CapabilityResponse, error) {
			ch := make(chan CapabilityResponse, 10)
			return ch, nil
		},
	}
	req := CapabilityRequest{}
	val, err := ExecuteSync(ctxWithTimeout, mcwe, req)

	assert.ErrorContains(t, err, "context timed out after")
	assert.Nil(t, val)
}

func Test_ExecuteSyncCapabilityErrors(t *testing.T) {
	expectedErr := errors.New("something went wrong during execution")
	mcwe := &mockCapabilityWithExecute{
		ExecuteFn: func(ctx context.Context, req CapabilityRequest) (<-chan CapabilityResponse, error) {
			ch := make(chan CapabilityResponse, 10)

			ch <- CapabilityResponse{nil, expectedErr}

			close(ch)

			return ch, nil
		},
	}
	req := CapabilityRequest{}
	val, err := ExecuteSync(tests.Context(t), mcwe, req)

	assert.ErrorContains(t, err, expectedErr.Error())
	assert.Nil(t, val)
}

func Test_ExecuteSyncDoesNotReturnValues(t *testing.T) {
	mcwe := &mockCapabilityWithExecute{
		ExecuteFn: func(ctx context.Context, req CapabilityRequest) (<-chan CapabilityResponse, error) {
			ch := make(chan CapabilityResponse, 10)
			close(ch)
			return ch, nil
		},
	}
	req := CapabilityRequest{}
	val, err := ExecuteSync(tests.Context(t), mcwe, req)

	assert.ErrorContains(t, err, "capability did not return any values")
	assert.Nil(t, val)
}

func Test_MustNewCapabilityInfo(t *testing.T) {
	assert.Panics(t, func() {
		MustNewCapabilityInfo(
			"capability-id",
			CapabilityTypeAction,
			"This is a mock capability that doesn't do anything.",
		)
	})
}

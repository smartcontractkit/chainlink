package syncer_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/syncer"
)

// mockService implements the services.Service interface for testing.
type mockService struct {
	mock.Mock
}

func (m *mockService) Start(ctx context.Context) error {
	return nil
}

func (m *mockService) HealthReport() map[string]error {
	return nil
}

func (m *mockService) Name() string {
	return "MockEngine"
}

func (m *mockService) Ready() error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockService) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestEngineRegistry_AddAndGet(t *testing.T) {
	registry := syncer.NewEngineRegistry()

	svc := mockService{}
	svc.On("Ready").Return(nil)

	// Add the service
	err := registry.Add("engine1", &svc)
	require.NoError(t, err)

	// Try adding the same ID again; expect an error
	err = registry.Add("engine1", &svc)
	require.Error(t, err)
	require.Contains(t, err.Error(), "duplicate")

	// Get the service
	gotSvc, err := registry.Get("engine1")
	require.NoError(t, err)
	require.Equal(t, svc.Name(), gotSvc.Name())

	// Try getting a non-existent engine
	_, err = registry.Get("non-existent")
	require.Error(t, err)
	require.Contains(t, err.Error(), "not found")
}

func TestEngineRegistry_IsRunning(t *testing.T) {
	registry := syncer.NewEngineRegistry()

	runningSvc := &mockService{}
	runningSvc.On("Ready").Return(nil)

	notRunningSvc := &mockService{}
	notRunningSvc.On("Ready").Return(errors.New("not ready"))

	_ = registry.Add("running", runningSvc)
	_ = registry.Add("notRunning", notRunningSvc)

	require.True(t, registry.IsRunning("running"), "expected running to be true if Ready() == nil")
	require.False(t, registry.IsRunning("notRunning"), "expected running to be false if Ready() != nil")
	require.False(t, registry.IsRunning("non-existent"), "expected false for non-existent engine")
}

func TestEngineRegistry_Pop(t *testing.T) {
	registry := syncer.NewEngineRegistry()

	svc := &mockService{}
	err := registry.Add("id1", svc)
	require.NoError(t, err)

	// Pop should remove and return the service
	poppedSvc, err := registry.Pop("id1")
	require.NoError(t, err)
	require.Equal(t, svc.Name(), poppedSvc.Name())

	// After popping, engine should not exist in registry
	_, err = registry.Get("id1")
	require.Error(t, err)

	// Popping a non-existent engine
	_, err = registry.Pop("non-existent")
	require.Error(t, err)
}

func TestEngineRegistry_Close(t *testing.T) {
	registry := syncer.NewEngineRegistry()

	// Set up multiple services to test aggregated errors
	svc1 := &mockService{}
	svc2 := &mockService{}
	svc3 := &mockService{}

	// We want to track Close calls and possibly return different errors
	svc1.On("Close").Return(nil).Once()
	svc2.On("Close").Return(errors.New("close error on svc2")).Once()
	svc3.On("Close").Return(errors.New("close error on svc3")).Once()

	_ = registry.Add("svc1", svc1)
	_ = registry.Add("svc2", svc2)
	_ = registry.Add("svc3", svc3)

	// Call Close
	err := registry.Close()
	// We expect both errors from svc2 and svc3 to be joined
	require.Error(t, err, "expected an aggregated error from multiple closes")
	// Check that the error message contains both parts
	require.True(t, strings.Contains(err.Error(), "svc2") && strings.Contains(err.Error(), "svc3"),
		"expected error to contain all relevant messages")

	// After Close, registry should be empty
	require.Equal(t, 0, registry.Size(), "registry should be empty after Close")

	// Additional Close calls should not error because registry is empty
	err = registry.Close()
	require.NoError(t, err, "expected no error on subsequent Close calls when registry is empty")

	// Ensure the mock expectations have been met (all calls accounted for)
	svc1.AssertExpectations(t)
	svc2.AssertExpectations(t)
	svc3.AssertExpectations(t)
}

func TestEngineRegistry_Size(t *testing.T) {
	registry := syncer.NewEngineRegistry()
	require.Equal(t, 0, registry.Size(), "initial registry should have size 0")

	_ = registry.Add("id1", new(mockService))
	_ = registry.Add("id2", new(mockService))
	require.Equal(t, 2, registry.Size())

	_, _ = registry.Pop("id1")
	require.Equal(t, 1, registry.Size())

	_ = registry.Close()
	require.Equal(t, 0, registry.Size())
}

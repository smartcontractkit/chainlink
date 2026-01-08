package arbiter

import (
	"context"
	"fmt"
	"iter"
	"math/big"
	"net"
	"testing"
	"time"

	"github.com/smartcontractkit/freeport"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// mockContractReader is a mock implementation of types.ContractReader for testing.
type mockContractReader struct {
	types.UnimplementedContractReader
	desiredShardCount uint64
	err               error
}

func (m *mockContractReader) Name() string {
	return "mockContractReader"
}

func (m *mockContractReader) Start(ctx context.Context) error {
	return nil
}

func (m *mockContractReader) Close() error {
	return nil
}

func (m *mockContractReader) Ready() error {
	return nil
}

func (m *mockContractReader) HealthReport() map[string]error {
	return nil
}

func (m *mockContractReader) Bind(ctx context.Context, bindings []types.BoundContract) error {
	return nil
}

func (m *mockContractReader) Unbind(ctx context.Context, bindings []types.BoundContract) error {
	return nil
}

func (m *mockContractReader) GetLatestValue(ctx context.Context, readIdentifier string, confidenceLevel primitives.ConfidenceLevel, params any, returnVal any) error {
	if m.err != nil {
		return m.err
	}
	// Set the result to our mock value
	if ptr, ok := returnVal.(**big.Int); ok {
		*ptr = big.NewInt(int64(m.desiredShardCount)) //nolint:gosec // G115: test mock value
	}
	return nil
}

func (m *mockContractReader) GetLatestValueWithHeadData(ctx context.Context, readIdentifier string, confidenceLevel primitives.ConfidenceLevel, params any, returnVal any) (head *types.Head, err error) {
	err = m.GetLatestValue(ctx, readIdentifier, confidenceLevel, params, returnVal)
	return nil, err
}

func (m *mockContractReader) BatchGetLatestValues(ctx context.Context, request types.BatchGetLatestValuesRequest) (types.BatchGetLatestValuesResult, error) {
	return nil, nil
}

func (m *mockContractReader) QueryKey(ctx context.Context, contract types.BoundContract, filter query.KeyFilter, limitAndSort query.LimitAndSort, sequenceDataType any) ([]types.Sequence, error) {
	return nil, nil
}

func (m *mockContractReader) QueryKeys(ctx context.Context, filters []types.ContractKeyFilter, limitAndSort query.LimitAndSort) (iter.Seq2[string, types.Sequence], error) {
	return nil, nil
}

// mockContractReaderFactory creates a ContractReaderFactory that returns the mock reader.
func mockContractReaderFactory(mockReader *mockContractReader) ContractReaderFactory {
	return func(ctx context.Context, cfg []byte) (types.ContractReader, error) {
		return mockReader, nil
	}
}

// Test configuration defaults
const (
	testPollInterval    = 100 * time.Millisecond
	testRetryInterval   = 100 * time.Millisecond
	testShardConfigAddr = "0x1234567890abcdef"
)

func TestArbiter_New(t *testing.T) {
	lggr := logger.TestLogger(t)
	mockReader := &mockContractReader{desiredShardCount: 10}
	factory := mockContractReaderFactory(mockReader)

	arb, err := New(lggr, factory, testShardConfigAddr, uint16(freeport.GetOne(t)), testPollInterval, testRetryInterval) //nolint:gosec // G115: freeport returns valid port range

	require.NoError(t, err)
	require.NotNil(t, arb)
	assert.Equal(t, "Arbiter", arb.Name())
}

func TestArbiter_StartClose(t *testing.T) {
	lggr := logger.TestLogger(t)
	mockReader := &mockContractReader{desiredShardCount: 10}
	factory := mockContractReaderFactory(mockReader)

	arb, err := New(lggr, factory, testShardConfigAddr, uint16(freeport.GetOne(t)), testPollInterval, testRetryInterval) //nolint:gosec // G115: freeport returns valid port range
	require.NoError(t, err)

	// Test start
	err = arb.Start(context.Background())
	require.NoError(t, err)
	t.Cleanup(func() { arb.Close() })

	// Give gRPC server and syncer a moment to start
	time.Sleep(100 * time.Millisecond)

	// Test health after start
	healthReport := arb.HealthReport()
	require.Contains(t, healthReport, arb.Name())
	require.NoError(t, healthReport[arb.Name()])
}

func TestArbiter_ServiceTestRun(t *testing.T) {
	lggr := logger.TestLogger(t)
	mockReader := &mockContractReader{desiredShardCount: 10}
	factory := mockContractReaderFactory(mockReader)

	arb, err := New(lggr, factory, testShardConfigAddr, uint16(freeport.GetOne(t)), testPollInterval, testRetryInterval) //nolint:gosec // G115: freeport returns valid port range
	require.NoError(t, err)

	// Use servicetest.Run to handle lifecycle
	// This starts the service and registers cleanup to stop it
	servicetest.Run(t, arb)

	// Service should be running after servicetest.Run
	err = arb.Ready()
	require.NoError(t, err)
}

func TestArbiter_HealthReport(t *testing.T) {
	lggr := logger.TestLogger(t)
	mockReader := &mockContractReader{desiredShardCount: 10}
	factory := mockContractReaderFactory(mockReader)

	arb, err := New(lggr, factory, testShardConfigAddr, uint16(freeport.GetOne(t)), testPollInterval, testRetryInterval) //nolint:gosec // G115: freeport returns valid port range
	require.NoError(t, err)

	t.Run("before start - not ready", func(t *testing.T) {
		healthReport := arb.HealthReport()
		require.Contains(t, healthReport, arb.Name())
		// Before start, Ready() should return an error
		assert.Error(t, healthReport[arb.Name()])
	})

	t.Run("after start - ready", func(t *testing.T) {
		err := arb.Start(context.Background())
		require.NoError(t, err)
		t.Cleanup(func() { arb.Close() })

		healthReport := arb.HealthReport()
		require.Contains(t, healthReport, arb.Name())
		require.NoError(t, healthReport[arb.Name()])
	})
}

func TestArbiter_DoubleStart(t *testing.T) {
	lggr := logger.TestLogger(t)
	mockReader := &mockContractReader{desiredShardCount: 10}
	factory := mockContractReaderFactory(mockReader)

	arb, err := New(lggr, factory, testShardConfigAddr, uint16(freeport.GetOne(t)), testPollInterval, testRetryInterval) //nolint:gosec // G115: freeport returns valid port range
	require.NoError(t, err)

	// First start should succeed
	err = arb.Start(context.Background())
	require.NoError(t, err)
	t.Cleanup(func() { arb.Close() })

	// Second start should return error (StartOnce)
	err = arb.Start(context.Background())
	require.Error(t, err)
}

func TestArbiter_DoubleClose(t *testing.T) {
	lggr := logger.TestLogger(t)
	mockReader := &mockContractReader{desiredShardCount: 10}
	factory := mockContractReaderFactory(mockReader)

	arb, err := New(lggr, factory, testShardConfigAddr, uint16(freeport.GetOne(t)), testPollInterval, testRetryInterval) //nolint:gosec // G115: freeport returns valid port range
	require.NoError(t, err)

	err = arb.Start(context.Background())
	require.NoError(t, err)

	// First close should succeed
	err = arb.Close()
	require.NoError(t, err)

	// Second close should return error (StopOnce)
	err = arb.Close()
	assert.Error(t, err)
}

func TestArbiter_Name(t *testing.T) {
	lggr := logger.TestLogger(t)
	mockReader := &mockContractReader{desiredShardCount: 10}
	factory := mockContractReaderFactory(mockReader)

	arb, err := New(lggr, factory, testShardConfigAddr, uint16(freeport.GetOne(t)), testPollInterval, testRetryInterval) //nolint:gosec // G115: freeport returns valid port range
	require.NoError(t, err)

	assert.Equal(t, "Arbiter", arb.Name())
}

func TestArbiter_Ready(t *testing.T) {
	lggr := logger.TestLogger(t)
	mockReader := &mockContractReader{desiredShardCount: 10}
	factory := mockContractReaderFactory(mockReader)

	arb, err := New(lggr, factory, testShardConfigAddr, uint16(freeport.GetOne(t)), testPollInterval, testRetryInterval) //nolint:gosec // G115: freeport returns valid port range
	require.NoError(t, err)

	// Before start, Ready should return error
	err = arb.Ready()
	require.Error(t, err)

	// After start, Ready should return nil
	err = arb.Start(context.Background())
	require.NoError(t, err)

	err = arb.Ready()
	require.NoError(t, err)

	// After close, Ready should return error
	err = arb.Close()
	require.NoError(t, err)

	err = arb.Ready()
	assert.Error(t, err)
}

func TestArbiter_GRPCServerListening(t *testing.T) {
	lggr := logger.TestLogger(t)
	mockReader := &mockContractReader{desiredShardCount: 10}
	factory := mockContractReaderFactory(mockReader)

	port := freeport.GetOne(t)
	arb, err := New(lggr, factory, testShardConfigAddr, uint16(port), testPollInterval, testRetryInterval) //nolint:gosec // G115: freeport returns valid port range
	require.NoError(t, err)

	// Start the arbiter
	err = arb.Start(context.Background())
	require.NoError(t, err)
	t.Cleanup(func() { arb.Close() })

	// Give gRPC server a moment to start
	time.Sleep(100 * time.Millisecond)

	// Verify gRPC server is actually listening by attempting to connect
	addr := fmt.Sprintf("localhost:%d", port)
	dialer := &net.Dialer{Timeout: 1 * time.Second}
	conn, err := dialer.DialContext(context.Background(), "tcp", addr)
	require.NoError(t, err, "gRPC server should be listening on port %d", port)
	conn.Close()
}

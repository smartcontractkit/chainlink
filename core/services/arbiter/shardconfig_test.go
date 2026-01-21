package arbiter

import (
	"context"
	"iter"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// mockShardConfigContractReader is a mock implementation for testing ShardConfigSyncer.
type mockShardConfigContractReader struct {
	types.UnimplementedContractReader
	shardCount uint64
	err        error
	started    bool
	bound      bool
}

func (m *mockShardConfigContractReader) Name() string {
	return "mockShardConfigContractReader"
}

func (m *mockShardConfigContractReader) Start(ctx context.Context) error {
	m.started = true
	return nil
}

func (m *mockShardConfigContractReader) Close() error {
	m.started = false
	return nil
}

func (m *mockShardConfigContractReader) Ready() error {
	return nil
}

func (m *mockShardConfigContractReader) HealthReport() map[string]error {
	return nil
}

func (m *mockShardConfigContractReader) Bind(ctx context.Context, bindings []types.BoundContract) error {
	m.bound = true
	return nil
}

func (m *mockShardConfigContractReader) Unbind(ctx context.Context, bindings []types.BoundContract) error {
	return nil
}

func (m *mockShardConfigContractReader) GetLatestValue(ctx context.Context, readIdentifier string, confidenceLevel primitives.ConfidenceLevel, params any, returnVal any) error {
	if m.err != nil {
		return m.err
	}
	if ptr, ok := returnVal.(**big.Int); ok {
		*ptr = big.NewInt(int64(m.shardCount)) //nolint:gosec // G115: test mock value
	}
	return nil
}

func (m *mockShardConfigContractReader) GetLatestValueWithHeadData(ctx context.Context, readIdentifier string, confidenceLevel primitives.ConfidenceLevel, params any, returnVal any) (head *types.Head, err error) {
	err = m.GetLatestValue(ctx, readIdentifier, confidenceLevel, params, returnVal)
	return nil, err
}

func (m *mockShardConfigContractReader) BatchGetLatestValues(ctx context.Context, request types.BatchGetLatestValuesRequest) (types.BatchGetLatestValuesResult, error) {
	return nil, nil
}

func (m *mockShardConfigContractReader) QueryKey(ctx context.Context, contract types.BoundContract, filter query.KeyFilter, limitAndSort query.LimitAndSort, sequenceDataType any) ([]types.Sequence, error) {
	return nil, nil
}

func (m *mockShardConfigContractReader) QueryKeys(ctx context.Context, filters []types.ContractKeyFilter, limitAndSort query.LimitAndSort) (iter.Seq2[string, types.Sequence], error) {
	return nil, nil
}

func mockShardConfigReaderFactory(reader *mockShardConfigContractReader) ContractReaderFactory {
	return func(ctx context.Context, cfg []byte) (types.ContractReader, error) {
		return reader, nil
	}
}

func TestShardConfigSyncer_New(t *testing.T) {
	lggr := logger.TestLogger(t)
	mockReader := &mockShardConfigContractReader{shardCount: 10}
	factory := mockShardConfigReaderFactory(mockReader)

	syncer := NewShardConfigSyncer(factory, "0x1234", 100*time.Millisecond, 100*time.Millisecond, lggr)

	require.NotNil(t, syncer)
	assert.Contains(t, syncer.Name(), "ShardConfigSyncer")
}

func TestShardConfigSyncer_GetDesiredShardCount_BeforeFetch(t *testing.T) {
	lggr := logger.TestLogger(t)
	mockReader := &mockShardConfigContractReader{shardCount: 10}
	factory := mockShardConfigReaderFactory(mockReader)

	syncer := NewShardConfigSyncer(factory, "0x1234", 100*time.Millisecond, 100*time.Millisecond, lggr)

	// Before start, GetDesiredShardCount should return error
	count, err := syncer.GetDesiredShardCount(context.Background())
	require.Error(t, err)
	assert.Equal(t, uint64(0), count)
	assert.Contains(t, err.Error(), "not yet available")
}

func TestShardConfigSyncer_GetDesiredShardCount_AfterFetch(t *testing.T) {
	lggr := logger.TestLogger(t)
	mockReader := &mockShardConfigContractReader{shardCount: 42}
	factory := mockShardConfigReaderFactory(mockReader)

	syncer := NewShardConfigSyncer(factory, "0x1234", 100*time.Millisecond, 100*time.Millisecond, lggr)

	// Start the syncer
	err := syncer.Start(context.Background())
	require.NoError(t, err)
	t.Cleanup(func() { syncer.Close() })

	// Wait for initial fetch (contract reader init + first poll)
	time.Sleep(200 * time.Millisecond)

	// After fetch, GetDesiredShardCount should return the cached value
	count, err := syncer.GetDesiredShardCount(context.Background())
	require.NoError(t, err)
	assert.Equal(t, uint64(42), count)
}

func TestShardConfigSyncer_StartClose(t *testing.T) {
	lggr := logger.TestLogger(t)
	mockReader := &mockShardConfigContractReader{shardCount: 10}
	factory := mockShardConfigReaderFactory(mockReader)

	syncer := NewShardConfigSyncer(factory, "0x1234", 100*time.Millisecond, 100*time.Millisecond, lggr)

	// Start
	err := syncer.Start(context.Background())
	require.NoError(t, err)
	t.Cleanup(func() { syncer.Close() })

	// Give it time to initialize
	time.Sleep(100 * time.Millisecond)
}

func TestShardConfigSyncer_HealthReport(t *testing.T) {
	lggr := logger.TestLogger(t)
	mockReader := &mockShardConfigContractReader{shardCount: 10}
	factory := mockShardConfigReaderFactory(mockReader)

	syncer := NewShardConfigSyncer(factory, "0x1234", 100*time.Millisecond, 100*time.Millisecond, lggr)

	// Before start
	healthReport := syncer.HealthReport()
	require.Contains(t, healthReport, syncer.Name())
	// Before start, Ready() should return an error
	require.Error(t, healthReport[syncer.Name()])

	// Start
	err := syncer.Start(context.Background())
	require.NoError(t, err)
	t.Cleanup(func() { syncer.Close() })

	// After start
	healthReport = syncer.HealthReport()
	require.Contains(t, healthReport, syncer.Name())
	require.NoError(t, healthReport[syncer.Name()])
}

func TestShardConfigSyncer_DoubleStart(t *testing.T) {
	lggr := logger.TestLogger(t)
	mockReader := &mockShardConfigContractReader{shardCount: 10}
	factory := mockShardConfigReaderFactory(mockReader)

	syncer := NewShardConfigSyncer(factory, "0x1234", 100*time.Millisecond, 100*time.Millisecond, lggr)

	// First start should succeed
	err := syncer.Start(context.Background())
	require.NoError(t, err)
	t.Cleanup(func() { syncer.Close() })

	// Second start should return error (StartOnce)
	err = syncer.Start(context.Background())
	require.Error(t, err)
}

func TestShardConfigSyncer_DoubleClose(t *testing.T) {
	lggr := logger.TestLogger(t)
	mockReader := &mockShardConfigContractReader{shardCount: 10}
	factory := mockShardConfigReaderFactory(mockReader)

	syncer := NewShardConfigSyncer(factory, "0x1234", 100*time.Millisecond, 100*time.Millisecond, lggr)

	err := syncer.Start(context.Background())
	require.NoError(t, err)

	// First close should succeed
	err = syncer.Close()
	require.NoError(t, err)

	// Second close should return error (StopOnce)
	err = syncer.Close()
	assert.Error(t, err)
}

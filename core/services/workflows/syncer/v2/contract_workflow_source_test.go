package v2

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/settings"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper_v2"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// Test constants for workflow metadata
const (
	testOwnerAddress = "0x1234567890123456789012345678901234567890"
	testBinaryURL    = "https://example.com/binary.wasm"
	testConfigURL    = "https://example.com/config.json"
)

// testBinaryContent and testConfigContent are mock content used for canonical workflowID calculation
var (
	testBinaryContent = []byte("mock-wasm-binary-content")
	testConfigContent = []byte("{}")
)

// mockWorkflowContractReader is a mock implementation of ContractReader for testing ContractWorkflowSource.
// Note: Reflection is required here because the ContractReader interface in chainlink-common
// uses `any` for the result parameter, and the production code passes an anonymous struct.
type mockWorkflowContractReader struct {
	commontypes.ContractReader
	workflowList []workflow_registry_wrapper_v2.WorkflowRegistryWorkflowMetadataView
	head         *commontypes.Head
	getLatestErr error
	bindErr      error
	startErr     error
}

func (m *mockWorkflowContractReader) GetLatestValueWithHeadData(
	_ context.Context,
	_ string,
	_ primitives.ConfidenceLevel,
	_ any,
	result any,
) (*commontypes.Head, error) {
	if m.getLatestErr != nil {
		return nil, m.getLatestErr
	}

	resultVal := reflect.ValueOf(result).Elem()
	listField := resultVal.FieldByName("List")
	if listField.IsValid() && listField.CanSet() {
		listField.Set(reflect.ValueOf(m.workflowList))
	}

	return m.head, nil
}

func (m *mockWorkflowContractReader) Bind(_ context.Context, _ []commontypes.BoundContract) error {
	return m.bindErr
}

func (m *mockWorkflowContractReader) Start(_ context.Context) error {
	return m.startErr
}

// createTestWorkflowMetadata creates a test workflow metadata view for testing.
// It uses the canonical workflow ID calculation to ensure test data is realistic.
func createTestWorkflowMetadata(name string, family string) workflow_registry_wrapper_v2.WorkflowRegistryWorkflowMetadataView {
	owner := common.HexToAddress(testOwnerAddress)

	// Use canonical workflow ID calculation
	workflowID, err := workflows.GenerateWorkflowID(owner.Bytes(), name, testBinaryContent, testConfigContent, "")
	if err != nil {
		panic("failed to generate workflow ID: " + err.Error())
	}

	return workflow_registry_wrapper_v2.WorkflowRegistryWorkflowMetadataView{
		WorkflowId:   workflowID,
		Owner:        owner,
		CreatedAt:    1234567890,
		Status:       WorkflowStatusActive,
		WorkflowName: name,
		BinaryUrl:    testBinaryURL,
		ConfigUrl:    testConfigURL,
		Tag:          "v1.0.0",
		Attributes:   []byte("{}"),
		DonFamily:    family,
	}
}

func TestContractWorkflowSource_ListWorkflowMetadata_Success(t *testing.T) {
	t.Parallel()
	lggr := logger.TestLogger(t)
	ctx := t.Context()

	mockReader := &mockWorkflowContractReader{
		workflowList: []workflow_registry_wrapper_v2.WorkflowRegistryWorkflowMetadataView{
			createTestWorkflowMetadata("workflow-1", "family-a"),
			createTestWorkflowMetadata("workflow-2", "family-a"),
		},
		head: &commontypes.Head{Height: "100"},
	}

	source := NewContractWorkflowSource(
		lggr,
		func(_ context.Context, _ []byte) (commontypes.ContractReader, error) {
			return mockReader, nil
		},
		testOwnerAddress,
		"test-chain-selector",
	)

	// Manually set the contract reader (simulating successful initialization)
	source.contractReader = mockReader

	don := capabilities.DON{
		ID:       1,
		Families: []string{"family-a"},
	}

	wfs, headResult, err := source.ListWorkflowMetadata(ctx, don)
	require.NoError(t, err)
	assert.Len(t, wfs, 2)
	assert.Equal(t, "100", headResult.Height)
	assert.Equal(t, "workflow-1", wfs[0].WorkflowName)
	assert.Equal(t, "workflow-2", wfs[1].WorkflowName)
}

func TestContractWorkflowSource_ListWorkflowMetadata_MultipleDONFamilies(t *testing.T) {
	t.Parallel()
	lggr := logger.TestLogger(t)
	ctx := t.Context()

	// This mock will be called twice (once for each DON family)
	mockReader := &mockWorkflowContractReader{
		workflowList: []workflow_registry_wrapper_v2.WorkflowRegistryWorkflowMetadataView{
			createTestWorkflowMetadata("workflow-1", "family-a"),
		},
		head: &commontypes.Head{Height: "100"},
	}

	source := NewContractWorkflowSource(
		lggr,
		func(_ context.Context, _ []byte) (commontypes.ContractReader, error) {
			return mockReader, nil
		},
		testOwnerAddress,
		"test-chain-selector",
	)

	source.contractReader = mockReader

	don := capabilities.DON{
		ID:       1,
		Families: []string{"family-a", "family-b"},
	}

	wfs, headResult, err := source.ListWorkflowMetadata(ctx, don)
	require.NoError(t, err)
	// Should return 2 workflows (1 per family call, but mock returns same list each time = 2 total)
	assert.Len(t, wfs, 2)
	assert.Equal(t, "100", headResult.Height)
}

func TestContractWorkflowSource_ListWorkflowMetadata_NotInitialized(t *testing.T) {
	t.Parallel()
	lggr := logger.TestLogger(t)
	ctx := t.Context()

	// Factory that always fails
	source := NewContractWorkflowSource(
		lggr,
		func(_ context.Context, _ []byte) (commontypes.ContractReader, error) {
			return nil, errors.New("factory error")
		},
		testOwnerAddress,
		"test-chain-selector",
	)

	don := capabilities.DON{
		ID:       1,
		Families: []string{"family-a"},
	}

	// Contract reader is nil, should return error
	wfs, headResult, err := source.ListWorkflowMetadata(ctx, don)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "contract reader not initialized")
	assert.Nil(t, wfs)
	assert.Nil(t, headResult)
}

func TestContractWorkflowSource_ListWorkflowMetadata_ContractReaderError(t *testing.T) {
	t.Parallel()
	lggr := logger.TestLogger(t)
	ctx := t.Context()

	mockReader := &mockWorkflowContractReader{
		getLatestErr: errors.New("contract read failed"),
	}

	source := NewContractWorkflowSource(
		lggr,
		func(_ context.Context, _ []byte) (commontypes.ContractReader, error) {
			return mockReader, nil
		},
		testOwnerAddress,
		"test-chain-selector",
	)
	source.contractReader = mockReader

	don := capabilities.DON{
		ID:       1,
		Families: []string{"family-a"},
	}

	wfs, headResult, err := source.ListWorkflowMetadata(ctx, don)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get latest value with head data")
	assert.Empty(t, wfs)
	assert.Equal(t, "0", headResult.Height)
}

func TestContractWorkflowSource_ListWorkflowMetadata_EmptyResult(t *testing.T) {
	t.Parallel()
	lggr := logger.TestLogger(t)
	ctx := t.Context()

	mockReader := &mockWorkflowContractReader{
		workflowList: []workflow_registry_wrapper_v2.WorkflowRegistryWorkflowMetadataView{},
		head:         &commontypes.Head{Height: "100"},
	}

	source := NewContractWorkflowSource(
		lggr,
		func(_ context.Context, _ []byte) (commontypes.ContractReader, error) {
			return mockReader, nil
		},
		testOwnerAddress,
		"test-chain-selector",
	)
	source.contractReader = mockReader

	don := capabilities.DON{
		ID:       1,
		Families: []string{"family-a"},
	}

	wfs, headResult, err := source.ListWorkflowMetadata(ctx, don)
	require.NoError(t, err)
	assert.Empty(t, wfs)
	assert.Equal(t, "100", headResult.Height)
}

func TestContractWorkflowSource_Ready_NotInitialized(t *testing.T) {
	t.Parallel()
	lggr := logger.TestLogger(t)

	source := NewContractWorkflowSource(
		lggr,
		func(_ context.Context, _ []byte) (commontypes.ContractReader, error) {
			return nil, errors.New("factory error")
		},
		testOwnerAddress,
		"test-chain-selector",
	)

	err := source.Ready()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "contract reader not initialized")
}

func TestContractWorkflowSource_Ready_Initialized(t *testing.T) {
	t.Parallel()
	lggr := logger.TestLogger(t)

	mockReader := &mockWorkflowContractReader{}

	source := NewContractWorkflowSource(
		lggr,
		func(_ context.Context, _ []byte) (commontypes.ContractReader, error) {
			return mockReader, nil
		},
		testOwnerAddress,
		"test-chain-selector",
	)
	source.contractReader = mockReader

	err := source.Ready()
	require.NoError(t, err)
}

func TestContractWorkflowSource_tryInitialize_Success(t *testing.T) {
	t.Parallel()
	lggr := logger.TestLogger(t)
	ctx := t.Context()

	mockReader := &mockWorkflowContractReader{}

	source := NewContractWorkflowSource(
		lggr,
		func(_ context.Context, _ []byte) (commontypes.ContractReader, error) {
			return mockReader, nil
		},
		testOwnerAddress,
		"test-chain-selector",
	)

	// Initially not ready
	require.Error(t, source.Ready())

	// Try to initialize
	result := source.tryInitialize(ctx)
	assert.True(t, result)

	// Now should be ready
	assert.NoError(t, source.Ready())
}

func TestContractWorkflowSource_tryInitialize_AlreadyInitialized(t *testing.T) {
	t.Parallel()
	lggr := logger.TestLogger(t)
	ctx := t.Context()

	callCount := 0
	mockReader := &mockWorkflowContractReader{}

	source := NewContractWorkflowSource(
		lggr,
		func(_ context.Context, _ []byte) (commontypes.ContractReader, error) {
			callCount++
			return mockReader, nil
		},
		testOwnerAddress,
		"test-chain-selector",
	)

	// First initialization
	result := source.tryInitialize(ctx)
	assert.True(t, result)
	assert.Equal(t, 1, callCount)

	// Second call should return true without calling factory again
	result = source.tryInitialize(ctx)
	assert.True(t, result)
	assert.Equal(t, 1, callCount) // Still 1, factory not called again
}

func TestContractWorkflowSource_tryInitialize_FactoryError(t *testing.T) {
	t.Parallel()
	lggr := logger.TestLogger(t)
	ctx := t.Context()

	source := NewContractWorkflowSource(
		lggr,
		func(_ context.Context, _ []byte) (commontypes.ContractReader, error) {
			return nil, errors.New("factory error")
		},
		testOwnerAddress,
		"test-chain-selector",
	)

	result := source.tryInitialize(ctx)
	assert.False(t, result)
	assert.Error(t, source.Ready())
}

func TestContractWorkflowSource_Name(t *testing.T) {
	t.Parallel()
	lggr := logger.TestLogger(t)

	source := NewContractWorkflowSource(
		lggr,
		func(_ context.Context, _ []byte) (commontypes.ContractReader, error) {
			return nil, nil
		},
		testOwnerAddress,
		"test-chain-selector",
	)

	assert.Equal(t, ContractWorkflowSourceName, source.Name())
}

// pagedWorkflowContractReader serves workflowList honoring the Start/Limit params, and fails any
// call whose Limit exceeds maxLimit with a gas-cap error, like an RPC node with a low eth_call cap.
type pagedWorkflowContractReader struct {
	commontypes.ContractReader
	workflowList []workflow_registry_wrapper_v2.WorkflowRegistryWorkflowMetadataView
	maxLimit     int64
	err          error
	limits       []int64
}

func (m *pagedWorkflowContractReader) GetLatestValueWithHeadData(
	_ context.Context,
	_ string,
	_ primitives.ConfidenceLevel,
	params any,
	result any,
) (*commontypes.Head, error) {
	p := params.(GetWorkflowListByDONParams)
	m.limits = append(m.limits, p.Limit.Int64())
	if m.err != nil {
		return nil, m.err
	}
	if p.Limit.Int64() > m.maxLimit {
		return nil, errors.New("out of gas: gas required exceeds 16777216")
	}
	start := int(p.Start.Int64())
	end := min(start+int(p.Limit.Int64()), len(m.workflowList))
	var page []workflow_registry_wrapper_v2.WorkflowRegistryWorkflowMetadataView
	if start < len(m.workflowList) {
		page = m.workflowList[start:end]
	}
	reflect.ValueOf(result).Elem().FieldByName("List").Set(reflect.ValueOf(page))
	return &commontypes.Head{Height: "100"}, nil
}

func newPagedTestSource(t *testing.T, reader commontypes.ContractReader) *ContractWorkflowSource {
	source := NewContractWorkflowSource(
		logger.TestLogger(t),
		func(_ context.Context, _ []byte) (commontypes.ContractReader, error) { return reader, nil },
		testOwnerAddress,
		"test-chain-selector",
	)
	source.contractReader = reader
	return source
}

func TestContractWorkflowSource_ListWorkflowMetadata_AdaptivePageSize(t *testing.T) {
	t.Parallel()

	var list []workflow_registry_wrapper_v2.WorkflowRegistryWorkflowMetadataView
	for i := range 700 {
		list = append(list, createTestWorkflowMetadata(fmt.Sprintf("wf-%d", i), "family-a"))
	}
	reader := &pagedWorkflowContractReader{workflowList: list, maxLimit: 100}
	source := newPagedTestSource(t, reader)
	don := capabilities.DON{ID: 1, Families: []string{"family-a"}}

	wfs, _, err := source.ListWorkflowMetadata(t.Context(), don)
	require.NoError(t, err)
	assert.Len(t, wfs, 700, "no workflows may be skipped or duplicated after shrinking")
	// default 300 fails, 150 fails, then 75 succeeds: 9 full pages and a partial one.
	assert.Equal(t, []int64{300, 150, 75, 75, 75, 75, 75, 75, 75, 75, 75, 75}, reader.limits)

	// The reduced size is remembered: the next tick doesn't re-try 300.
	reader.limits = nil
	_, _, err = source.ListWorkflowMetadata(t.Context(), don)
	require.NoError(t, err)
	assert.Equal(t, int64(75), reader.limits[0])
	assert.Len(t, reader.limits, 10)
}

func TestContractWorkflowSource_ListWorkflowMetadata_GasErrorAtMinPageSize(t *testing.T) {
	t.Parallel()

	reader := &pagedWorkflowContractReader{maxLimit: 0}
	source := newPagedTestSource(t, reader)
	don := capabilities.DON{ID: 1, Families: []string{"family-a"}}

	_, _, err := source.ListWorkflowMetadata(t.Context(), don)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "out of gas")
	assert.Equal(t, []int64{300, 150, 75, 37, 18, 10}, reader.limits)
}

func TestContractWorkflowSource_ListWorkflowMetadata_NonGasErrorDoesNotShrink(t *testing.T) {
	t.Parallel()

	reader := &pagedWorkflowContractReader{err: errors.New("connection refused")}
	source := newPagedTestSource(t, reader)
	don := capabilities.DON{ID: 1, Families: []string{"family-a"}}

	_, _, err := source.ListWorkflowMetadata(t.Context(), don)
	require.Error(t, err)
	assert.Equal(t, []int64{300}, reader.limits)
}

func TestContractWorkflowSource_ListWorkflowMetadata_PageSizeFromSettings(t *testing.T) {
	t.Parallel()

	getter, err := settings.NewJSONGetter([]byte(`{"global":{"WorkflowRegistryMaxResultsPerQuery":"50"}}`))
	require.NoError(t, err)

	var list []workflow_registry_wrapper_v2.WorkflowRegistryWorkflowMetadataView
	for i := range 120 {
		list = append(list, createTestWorkflowMetadata(fmt.Sprintf("wf-%d", i), "family-a"))
	}
	reader := &pagedWorkflowContractReader{workflowList: list, maxLimit: 1000}
	source := newPagedTestSource(t, reader)
	source.SetSettingsGetter(getter)
	don := capabilities.DON{ID: 1, Families: []string{"family-a"}}

	wfs, _, err := source.ListWorkflowMetadata(t.Context(), don)
	require.NoError(t, err)
	assert.Len(t, wfs, 120)
	assert.Equal(t, []int64{50, 50, 50}, reader.limits)
}

func TestRegistryPageSizer_ReductionResetsWhenConfiguredValueChanges(t *testing.T) {
	t.Parallel()
	lggr := logger.TestLogger(t)

	var p registryPageSizer
	assert.Equal(t, int64(300), p.current(t.Context(), lggr))
	next, ok := p.shrink(300)
	require.True(t, ok)
	assert.Equal(t, int64(150), next)
	assert.Equal(t, int64(150), p.current(t.Context(), lggr))

	// An operator changing the configured limit discards the runtime reduction.
	getter, err := settings.NewJSONGetter([]byte(`{"global":{"WorkflowRegistryMaxResultsPerQuery":"200"}}`))
	require.NoError(t, err)
	p.setGetter(getter)
	assert.Equal(t, int64(200), p.current(t.Context(), lggr))

	_, ok = p.shrink(minRegistryPageSize)
	assert.False(t, ok)
}

func TestIsGasError(t *testing.T) {
	t.Parallel()
	assert.True(t, isGasError(errors.New("out of gas: gas required exceeds 16777216")))
	assert.True(t, isGasError(errors.New("execution reverted: Gas Required Exceeds Allowance (50000000)")))
	assert.True(t, isGasError(errors.New("eth_call gas cap exceeded")))
	assert.False(t, isGasError(errors.New("connection refused")))
	assert.False(t, isGasError(nil))
}

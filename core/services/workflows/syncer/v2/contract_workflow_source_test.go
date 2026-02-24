package v2

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
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

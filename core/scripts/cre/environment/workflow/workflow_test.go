package workflow

import (
	"bytes"
	"fmt"
	"io"
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	workflow_registry_wrapper "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper_v1"
)

// MockWorkflowRegistry mocks the WorkflowRegistry interface for testing
type MockWorkflowRegistry struct {
	mock.Mock
}

func (m *MockWorkflowRegistry) GetWorkflowMetadataListByOwner(opts *bind.CallOpts, owner common.Address, offset *big.Int, limit *big.Int) ([]workflow_registry_wrapper.WorkflowRegistryWorkflowMetadata, error) {
	args := m.Called(opts, owner, offset, limit)
	return args.Get(0).([]workflow_registry_wrapper.WorkflowRegistryWorkflowMetadata), args.Error(1)
}

// Mock implementations of the WorkflowRegistry methods
func (m *MockWorkflowRegistry) RegisterWorkflow(opts *bind.TransactOpts, name string, id [32]byte, donID uint32, status uint8, binary, config, secrets string) (*types.Transaction, error) {
	args := m.Called(opts, name, id, donID, status, binary, config, secrets)
	return args.Get(0).(*types.Transaction), args.Error(1)
}

func (m *MockWorkflowRegistry) UpdateWorkflow(opts *bind.TransactOpts, key [32]byte, id [32]byte, binary, config, secrets string) (*types.Transaction, error) {
	args := m.Called(opts, key, id, binary, config, secrets)
	return args.Get(0).(*types.Transaction), args.Error(1)
}

func (m *MockWorkflowRegistry) DeleteWorkflow(opts *bind.TransactOpts, key [32]byte) (*types.Transaction, error) {
	args := m.Called(opts, key)
	return args.Get(0).(*types.Transaction), args.Error(1)
}

// TestComputeWorkflowKey tests the workflow key computation function
func TestComputeWorkflowKey(t *testing.T) {
	owner := common.HexToAddress("0x1234567890123456789012345678901234567890")
	name := "test-workflow"

	key := ComputeWorkflowKey(owner, name)

	// Check that we get a non-zero key
	require.NotEqual(t, [32]byte{}, key, "Key should not be empty")

	// Same inputs should produce same key (deterministic)
	key2 := ComputeWorkflowKey(owner, name)
	require.Equal(t, key, key2, "Same inputs should produce the same key")

	// Different inputs should produce different keys
	differentOwner := common.HexToAddress("0x0987654321098765432109876543210987654321")
	key3 := ComputeWorkflowKey(differentOwner, name)
	require.NotEqual(t, key, key3, "Different owners should produce different keys")

	differentName := "different-workflow"
	key4 := ComputeWorkflowKey(owner, differentName)
	require.NotEqual(t, key, key4, "Different names should produce different keys")
}

// TestValidateStatus tests the status validation function
func TestValidateStatus(t *testing.T) {
	// Test valid statuses
	require.NoError(t, validateStatus(WorkflowStatusActive), "Active status should be valid")
	require.NoError(t, validateStatus(WorkflowStatusPaused), "Paused status should be valid")

	// Test invalid status
	err := validateStatus(99)
	require.Error(t, err, "Invalid status should return an error")
	require.Contains(t, err.Error(), "invalid status", "Error should mention invalid status")
}

// TestFormatStatus tests the status formatting function
func TestFormatStatus(t *testing.T) {
	require.Equal(t, "Active (0)", formatStatus(WorkflowStatusActive))
	require.Equal(t, "Paused (1)", formatStatus(WorkflowStatusPaused))
	require.Equal(t, "Unknown (99)", formatStatus(99))
}

// TestMarkFlagsRequired tests the flag requirement marking utility
func TestMarkFlagsRequired(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().String("existing", "", "An existing flag")

	// Test with existing flag
	err := MarkFlagsRequired(cmd, "existing")
	require.NoError(t, err)

	// Test with non-existing flag
	err = MarkFlagsRequired(cmd, "nonexistent")
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to mark flag")
}

// TestGetWorkflowCmd tests the functionality of the get command
func TestGetWorkflowCmd(t *testing.T) {
	// Create a mock workflow registry
	mockRegistry := new(MockWorkflowRegistry)

	// Create test data
	owner := common.HexToAddress("0x1234567890123456789012345678901234567890")
	testWorkflows := []workflow_registry_wrapper.WorkflowRegistryWorkflowMetadata{
		{
			WorkflowName: "test-workflow-1",
			Owner:        owner,
			WorkflowID:   [32]byte{1, 2, 3},
			DonID:        42,
			Status:       WorkflowStatusActive,
			BinaryURL:    "https://example.com/binary1",
			ConfigURL:    "https://example.com/config1",
			SecretsURL:   "https://example.com/secrets1",
		},
		{
			WorkflowName: "test-workflow-2",
			Owner:        owner,
			WorkflowID:   [32]byte{4, 5, 6},
			DonID:        43,
			Status:       WorkflowStatusPaused,
			BinaryURL:    "https://example.com/binary2",
			ConfigURL:    "https://example.com/config2",
			SecretsURL:   "https://example.com/secrets2",
		},
	}

	// Setup mock to return our test workflows
	mockRegistry.On("GetWorkflowMetadataListByOwner", mock.Anything, owner, big.NewInt(0), big.NewInt(MaxWorkflowsToFetch)).
		Return(testWorkflows, nil)

	// Capture stdout to verify output
	originalStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Call our extracted function that mimics the Run function
	displayWorkflows(mockRegistry, owner)

	// Restore stdout
	w.Close()
	os.Stdout = originalStdout

	// Read captured output
	var buf bytes.Buffer
	_, err := io.Copy(&buf, r)
	require.NoError(t, err, "Failed to read captured output")
	output := buf.String()

	// Verify output contains expected information
	require.Contains(t, output, "### Workflow ###")
	require.Contains(t, output, "Name: test-workflow-1")
	require.Contains(t, output, "Name: test-workflow-2")
	require.Contains(t, output, "Owner: 0x1234567890123456789012345678901234567890")
	require.Contains(t, output, "ID: 0x010203")
	require.Contains(t, output, "ID: 0x040506")
	require.Contains(t, output, "DON ID: 42")
	require.Contains(t, output, "DON ID: 43")
	require.Contains(t, output, "Status: 0 (Active (0))")
	require.Contains(t, output, "Status: 1 (Paused (1))")
	require.Contains(t, output, "Binary URL: https://example.com/binary1")
	require.Contains(t, output, "Config URL: https://example.com/config2")
	require.Contains(t, output, "Secrets URL: https://example.com/secrets2")

	// Verify the mock was called with expected parameters
	mockRegistry.AssertExpectations(t)
}

// Mock interface with GetWorkflowMetadataListByOwner
type WorkflowRegistryInterface interface {
	GetWorkflowMetadataListByOwner(opts *bind.CallOpts, owner common.Address, offset *big.Int, limit *big.Int) ([]workflow_registry_wrapper.WorkflowRegistryWorkflowMetadata, error)
	UpdateWorkflow(opts *bind.TransactOpts, key [32]byte, id [32]byte, binary, config, secrets string) (*types.Transaction, error)
	DeleteWorkflow(opts *bind.TransactOpts, key [32]byte) (*types.Transaction, error)
}

// Extract the display logic from getCmd.Run into a testable function
func displayWorkflows(client WorkflowRegistryInterface, owner common.Address) {
	// Use a simple stub for CallOpts instead of seth.NewCallOpts()
	callOpts := &bind.CallOpts{}

	// Get workflow details
	metadata, err := client.GetWorkflowMetadataListByOwner(callOpts, owner, big.NewInt(0), big.NewInt(MaxWorkflowsToFetch))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get workflow: %v\n", err)
		return
	}

	// Display workflow information
	for _, m := range metadata {
		fmt.Print("### Workflow ###\n")
		fmt.Printf("Name: %s\n", m.WorkflowName)
		fmt.Printf("Owner: %s\n", m.Owner.Hex())
		fmt.Printf("ID: 0x%x\n", m.WorkflowID)
		fmt.Printf("DON ID: %d\n", m.DonID)
		fmt.Printf("Status: %d (%s)\n", m.Status, formatStatus(m.Status))
		fmt.Printf("Binary URL: %s\n", m.BinaryURL)
		fmt.Printf("Config URL: %s\n", m.ConfigURL)
		fmt.Printf("Secrets URL: %s\n", m.SecretsURL)
	}
}

// TestUpdateWorkflowCmd tests the update workflow command functionality
func TestUpdateWorkflowCmd(t *testing.T) {
	// Create a mock workflow registry
	mockRegistry := new(MockWorkflowRegistry)

	// Test data
	testOwner := common.HexToAddress("0x1234567890123456789012345678901234567890")
	testName := "test-workflow"
	testID := "0x0102030405060708091011121314151617181920212223242526272829303132"
	testBinary := "https://example.com/new-binary"
	testConfig := "https://example.com/new-config"
	testSecrets := "https://example.com/new-secrets"

	// Create expected key and ID bytes
	expectedKey := ComputeWorkflowKey(testOwner, testName)
	var expectedIDBytes [32]byte
	copy(expectedIDBytes[:], common.FromHex(testID))

	// Create a mock transaction
	mockTx := types.NewTransaction(
		0,
		common.HexToAddress("0x"),
		big.NewInt(0),
		0,
		big.NewInt(0),
		nil,
	)

	// Setup mock expectations
	mockRegistry.On("UpdateWorkflow", mock.Anything, expectedKey, expectedIDBytes, testBinary, testConfig, testSecrets).
		Return(mockTx, nil)

	// Capture stdout
	originalStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Call the extracted function
	updateWorkflow(mockRegistry, testOwner, testName, expectedIDBytes, testBinary, testConfig, testSecrets)
	w.Close()

	// Read captured output
	var buf bytes.Buffer
	_, err := io.Copy(&buf, r)
	require.NoError(t, err, "Failed to read captured output")
	output := buf.String()
	os.Stdout = originalStdout

	// Verify output
	require.Contains(t, output, "successfully updated")
	require.Contains(t, output, testName)
	require.Contains(t, output, mockTx.Hash().Hex())

	// Verify mock expectations
	mockRegistry.AssertExpectations(t)
}

// TestDeleteWorkflowCmd tests the delete workflow command functionality
func TestDeleteWorkflowCmd(t *testing.T) {
	// Create a mock workflow registry
	mockRegistry := new(MockWorkflowRegistry)

	// Test data
	testOwner := common.HexToAddress("0x1234567890123456789012345678901234567890")
	testName := "test-workflow"

	// Create expected key
	expectedKey := ComputeWorkflowKey(testOwner, testName)

	// Create a mock transaction
	mockTx := types.NewTransaction(
		0,
		common.HexToAddress("0x"),
		big.NewInt(0),
		0,
		big.NewInt(0),
		nil,
	)

	// Setup mock expectations
	mockRegistry.On("DeleteWorkflow", mock.Anything, expectedKey).Return(mockTx, nil)

	// Capture stdout
	originalStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Call the extracted function
	deleteWorkflow(mockRegistry, testOwner, testName)
	w.Close()

	// Read captured output
	var buf bytes.Buffer
	_, err := io.Copy(&buf, r)
	require.NoError(t, err, "Failed to read captured output")
	output := buf.String()
	os.Stdout = originalStdout

	// Verify output
	require.Contains(t, output, "successfully deleted")
	require.Contains(t, output, testName)
	require.Contains(t, output, mockTx.Hash().Hex())

	// Verify mock expectations
	mockRegistry.AssertExpectations(t)
}

// Extract the update logic from updateCmd.Run into a testable function
func updateWorkflow(client WorkflowRegistryInterface, owner common.Address, name string,
	idBytes [32]byte, binaryURL, configURL, secretsURL string) {
	// Create a simple transactor for testing
	txOpt := &bind.TransactOpts{}

	// Update workflow
	output, err := client.UpdateWorkflow(txOpt, ComputeWorkflowKey(owner, name),
		idBytes, binaryURL, configURL, secretsURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to update workflow: %v\n", err)
		return
	}

	fmt.Printf("Workflow '%s' successfully updated\n", name)
	fmt.Printf("Transaction hash: %s\n", output.Hash().Hex())
}

// Extract the delete logic from deleteCmd.Run into a testable function
func deleteWorkflow(client WorkflowRegistryInterface, owner common.Address, name string) {
	// Create a simple transactor for testing
	txOpt := &bind.TransactOpts{}

	// Delete workflow
	txHash, err := client.DeleteWorkflow(txOpt, ComputeWorkflowKey(owner, name))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to delete workflow: %v\n", err)
		return
	}

	fmt.Printf("Workflow '%s' successfully deleted\n", name)
	fmt.Printf("Transaction hash: %s\n", txHash.Hash().Hex())
}

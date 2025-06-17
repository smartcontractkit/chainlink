package v1_0

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	workflow_registry "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper_v1"
	"github.com/smartcontractkit/chainlink/deployment/common/view/v1_0/mocks"
)

// TestNewWorkflowView tests the helper function that converts on-chain WorkflowMetadata -> WorkflowView.
func TestNewWorkflowView(t *testing.T) {
	t.Run("nil input => error", func(t *testing.T) {
		// Updated: nil input now returns an error.
		wv, err := NewWorkflowView(nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "workflow metadata is nil")
		// And we expect the returned struct to be empty.
		require.Equal(t, WorkflowView{}, wv)
	})

	t.Run("valid input => success", func(t *testing.T) {
		meta := workflow_registry.WorkflowRegistryWorkflowMetadata{
			WorkflowID:   [32]byte{0xab, 0xcd},
			Owner:        common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678"),
			DonID:        42,
			Status:       0, // WorkflowStatusActive
			WorkflowName: "TestWorkflow",
			BinaryURL:    "http://binary",
			ConfigURL:    "http://config",
			SecretsURL:   "http://secrets",
		}

		wv, err := NewWorkflowView(&meta)
		require.NoError(t, err)

		wantID := hex.EncodeToString(meta.WorkflowID[:])
		assert.Equal(t, wantID, wv.WorkflowID)
		assert.Equal(t, meta.Owner, wv.Owner)
		assert.Equal(t, uint32(42), wv.DonID)
		assert.Equal(t, WorkflowStatusActive, wv.Status)
		assert.Equal(t, "TestWorkflow", wv.WorkflowName)
		assert.Equal(t, "http://binary", wv.BinaryURL)
		assert.Equal(t, "http://config", wv.ConfigURL)
		assert.Equal(t, "http://secrets", wv.SecretsURL)
	})
}

// TestGenerateWorkflowRegistryView uses a mock workflow registry to test the main view generation function.
func TestGenerateWorkflowRegistryView(t *testing.T) {
	tests := []struct {
		name            string
		mockSetup       func(*mocks.WorkflowRegistryInterface)
		wantErrCount    int
		wantAllowedDONs []uint32
		wantAuthAddrs   []common.Address
		wantLocked      bool
		wantNumWf       int
	}{
		{
			name: "happy path - multiple DONs w/ single-page workflows",
			mockSetup: func(m *mocks.WorkflowRegistryInterface) {
				m.EXPECT().Address().Return(common.HexToAddress("0xFAKE1111FAKE1111FAKE1111FAKE1111FAKE1111"))
				m.EXPECT().TypeAndVersion(mock.Anything).Return("WorkflowRegistry 1.0.0", nil).Maybe()
				m.EXPECT().Owner(mock.Anything).Return(common.HexToAddress("0x12345..."), nil).Maybe()

				// Return 2 DON IDs.
				m.EXPECT().GetAllAllowedDONs(mock.Anything).Return([]uint32{42, 84}, nil)

				// Return 2 authorized addresses.
				m.EXPECT().GetAllAuthorizedAddresses(mock.Anything).
					Return([]common.Address{
						common.HexToAddress("0x1111111111111111111111111111111111111111"),
						common.HexToAddress("0x2222222222222222222222222222222222222222"),
					}, nil)

				// Return locked = false.
				m.EXPECT().IsRegistryLocked(mock.Anything).Return(false, nil)

				// For DON=42, single workflow -> active.
				wf42 := []workflow_registry.WorkflowRegistryWorkflowMetadata{
					{
						WorkflowID:   [32]byte{0xde, 0xad, 0xbe, 0xef},
						DonID:        42,
						Status:       0, // active
						WorkflowName: "Workflow A",
						BinaryURL:    "binaryA",
						ConfigURL:    "configA",
						SecretsURL:   "secretsA",
					},
				}
				m.EXPECT().GetWorkflowMetadataListByDON(mock.Anything, uint32(42), big.NewInt(0), big.NewInt(100)).
					Return(wf42, nil).Once()
				// Next page returns empty.
				m.EXPECT().GetWorkflowMetadataListByDON(mock.Anything, uint32(42), big.NewInt(100), big.NewInt(100)).
					Return([]workflow_registry.WorkflowRegistryWorkflowMetadata{}, nil).Maybe()

				// For DON=84, single workflow -> paused.
				wf84 := []workflow_registry.WorkflowRegistryWorkflowMetadata{
					{
						WorkflowID:   [32]byte{0xca, 0xfe, 0xba, 0xbe},
						DonID:        84,
						Status:       1, // paused
						WorkflowName: "Workflow B",
						BinaryURL:    "binaryB",
						ConfigURL:    "configB",
						SecretsURL:   "secretsB",
					},
				}
				m.EXPECT().GetWorkflowMetadataListByDON(mock.Anything, uint32(84), big.NewInt(0), big.NewInt(100)).
					Return(wf84, nil).Once()
				// Next page returns empty.
				m.EXPECT().GetWorkflowMetadataListByDON(mock.Anything, uint32(84), big.NewInt(100), big.NewInt(100)).
					Return([]workflow_registry.WorkflowRegistryWorkflowMetadata{}, nil).Maybe()
			},
			wantErrCount:    0,
			wantAllowedDONs: []uint32{42, 84},
			wantAuthAddrs: []common.Address{
				common.HexToAddress("0x1111111111111111111111111111111111111111"),
				common.HexToAddress("0x2222222222222222222222222222222222222222"),
			},
			wantLocked: false,
			wantNumWf:  2,
		},
		{
			name: "GetAllAllowedDONs returns error => include error",
			mockSetup: func(m *mocks.WorkflowRegistryInterface) {
				m.EXPECT().Address().Return(common.HexToAddress("0xABCD")).Maybe()
				m.EXPECT().TypeAndVersion(mock.Anything).Return("WorkflowRegistry 1.0.0", nil).Maybe()
				m.EXPECT().Owner(mock.Anything).Return(common.HexToAddress("0x12345..."), nil).Maybe()
				// Simulate three errors:
				// 1. GetAllAllowedDONs error. Since there are no dons, there are no workflows to fetch.
				m.EXPECT().GetAllAllowedDONs(mock.Anything).
					Return([]uint32{}, assert.AnError)
				// 2. GetAllAuthorizedAddresses error.
				m.EXPECT().GetAllAuthorizedAddresses(mock.Anything).
					Return(nil, assert.AnError)
				// 3. IsRegistryLocked error.
				m.EXPECT().IsRegistryLocked(mock.Anything).
					Return(false, assert.AnError)
			},
			wantErrCount: 3,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockWR := mocks.NewWorkflowRegistryInterface(t)
			tc.mockSetup(mockWR)

			view, errs := GenerateWorkflowRegistryView(mockWR)
			if tc.wantErrCount > 0 {
				require.NotEmpty(t, errs)
				return
			}
			require.Empty(t, errs)

			// Check fields.
			assert.Equal(t, tc.wantAllowedDONs, view.AllowedDONs)
			assert.Equal(t, tc.wantAuthAddrs, view.AuthorizedAddresses)
			assert.Equal(t, tc.wantLocked, view.IsRegistryLocked)
			assert.Len(t, view.Workflows, tc.wantNumWf)
		})
	}
}

func TestWorkflowStatus_MarshalUnmarshal(t *testing.T) {
	tests := []struct {
		name string
		ws   WorkflowStatus
		want string
	}{
		{"ACTIVE => JSON", WorkflowStatusActive, `"ACTIVE"`},
		{"PAUSED => JSON", WorkflowStatusPaused, `"PAUSED"`},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Marshal.
			data, err := tc.ws.MarshalJSON()
			require.NoError(t, err)
			assert.Equal(t, tc.want, string(data))

			// Unmarshal.
			var back WorkflowStatus
			err = back.UnmarshalJSON(data)
			require.NoError(t, err)
			assert.Equal(t, tc.ws, back)
		})
	}

	t.Run("unknown => error", func(t *testing.T) {
		raw := []byte(`"UNKNOWN(99)"`)
		var ws WorkflowStatus
		err := ws.UnmarshalJSON(raw)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid WorkflowStatus value")
	})
}

package capabilityregistry

import (
	"math/big"
	"slices"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	capabilities_registry "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
)

// MockCapabilityRegistry implements the methods we need for testing
type MockCapabilityRegistry struct {
	mock.Mock
}

func (m *MockCapabilityRegistry) GetCapabilities(opts *bind.CallOpts) ([]capabilities_registry.CapabilitiesRegistryCapability, error) {
	args := m.Called(opts)
	return args.Get(0).([]capabilities_registry.CapabilitiesRegistryCapability), args.Error(1)
}

func (m *MockCapabilityRegistry) AddCapabilities(opts *bind.TransactOpts, capabilities []capabilities_registry.CapabilitiesRegistryCapability) (*types.Transaction, error) {
	args := m.Called(opts, capabilities)
	return args.Get(0).(*types.Transaction), args.Error(1)
}

func (m *MockCapabilityRegistry) GetDON(opts *bind.CallOpts, donID uint32) (capabilities_registry.CapabilitiesRegistryDONInfo, error) {
	args := m.Called(opts, donID)
	return args.Get(0).(capabilities_registry.CapabilitiesRegistryDONInfo), args.Error(1)
}

func (m *MockCapabilityRegistry) GetDONs(opts *bind.CallOpts) ([]capabilities_registry.CapabilitiesRegistryDONInfo, error) {
	args := m.Called(opts)
	return args.Get(0).([]capabilities_registry.CapabilitiesRegistryDONInfo), args.Error(1)
}

func (m *MockCapabilityRegistry) GetNodes(opts *bind.CallOpts) ([]capabilities_registry.CapabilitiesRegistryNodeParams, error) {
	args := m.Called(opts)
	return args.Get(0).([]capabilities_registry.CapabilitiesRegistryNodeParams), args.Error(1)
}

func (m *MockCapabilityRegistry) UpdateNodes(opts *bind.TransactOpts, params []capabilities_registry.CapabilitiesRegistryNodeParams) (*types.Transaction, error) {
	args := m.Called(opts, params)
	return args.Get(0).(*types.Transaction), args.Error(1)
}

func (m *MockCapabilityRegistry) UpdateDON(opts *bind.TransactOpts, donID uint32, nodeP2PIds [][32]byte, capabilityConfigurations []capabilities_registry.CapabilitiesRegistryCapabilityConfiguration, isPublic bool, f uint8) (*types.Transaction, error) {
	args := m.Called(opts, donID, nodeP2PIds, capabilityConfigurations, isPublic, f)
	return args.Get(0).(*types.Transaction), args.Error(1)
}

// TestCreateCommand tests the capability creation functionality
func TestCreateCommand(t *testing.T) {
	// Create a mock registry
	mockRegistry := new(MockCapabilityRegistry)

	// Test data
	testName := "test-capability"
	testVersion := "1.0.0"
	testType := "trigger"
	testDonID := uint32(1)

	// Generate capability ID hash
	hashedCapID := getHashedCapabilityID(testName, testVersion)

	// Setup mock transactions
	mockAddTx := types.NewTx(&types.LegacyTx{
		Nonce:    0,
		GasPrice: big.NewInt(0),
		Gas:      0,
		To:       &common.Address{},
		Value:    big.NewInt(0),
		Data:     nil,
	})
	mockUpdateNodesTx := types.NewTx(&types.LegacyTx{
		Nonce:    1,
		GasPrice: big.NewInt(0),
		Gas:      0,
		To:       &common.Address{},
		Value:    big.NewInt(0),
		Data:     nil,
	})
	mockUpdateDONTx := types.NewTx(&types.LegacyTx{
		Nonce:    1,
		GasPrice: big.NewInt(0),
		Gas:      0,
		To:       &common.Address{},
		Value:    big.NewInt(0),
		Data:     nil,
	})

	// Setup test node and DON data
	testNodeP2PId := [32]byte{1, 2, 3}
	testNodes := []capabilities_registry.CapabilitiesRegistryNodeParams{
		{
			NodeOperatorId:      1,
			P2pId:               testNodeP2PId,
			Signer:              common.HexToHash("0x1"),
			EncryptionPublicKey: common.HexToHash("0x123"),
			HashedCapabilityIds: [][32]byte{},
		},
	}

	testDON := capabilities_registry.CapabilitiesRegistryDONInfo{
		Id:                       testDonID,
		ConfigCount:              1,
		NodeP2PIds:               [][32]byte{testNodeP2PId},
		CapabilityConfigurations: []capabilities_registry.CapabilitiesRegistryCapabilityConfiguration{},
		F:                        1,
		IsPublic:                 true,
		AcceptsWorkflows:         true,
	}

	// Test Case 1: Capability doesn't exist
	// Setup mock expectations
	mockRegistry.On("GetCapabilities", mock.Anything).Return([]capabilities_registry.CapabilitiesRegistryCapability{}, nil).Once()
	mockRegistry.On("AddCapabilities", mock.Anything, mock.MatchedBy(func(caps []capabilities_registry.CapabilitiesRegistryCapability) bool {
		return len(caps) == 1 && caps[0].LabelledName == testName && caps[0].Version == testVersion
	})).Return(mockAddTx, nil).Once()
	mockRegistry.On("GetDON", mock.Anything, testDonID).Return(testDON, nil).Once()
	mockRegistry.On("GetNodes", mock.Anything).Return(testNodes, nil).Once()

	// Match expected nodesToUpdate with the hashed capability ID added
	mockRegistry.On("UpdateNodes", mock.Anything, mock.MatchedBy(func(nodes []capabilities_registry.CapabilitiesRegistryNodeParams) bool {
		if len(nodes) != 1 {
			return false
		}
		if len(nodes[0].HashedCapabilityIds) != 1 {
			return false
		}
		return nodes[0].P2pId == testNodeP2PId
	})).Return(mockUpdateNodesTx, nil).Once()

	// Match expected DON update with the capability added
	mockRegistry.On("UpdateDON", mock.Anything, testDonID, testDON.NodeP2PIds,
		mock.MatchedBy(func(configs []capabilities_registry.CapabilitiesRegistryCapabilityConfiguration) bool {
			return len(configs) == 1 && configs[0].CapabilityId == hashedCapID
		}),
		testDON.IsPublic, testDON.F).Return(mockUpdateDONTx, nil).Once()

	// Call extracted function for creating capability
	require.NoError(t, createCapability(mockRegistry, testName, testVersion, testType, testDonID))

	// Test Case 2: Capability already exists
	existingCaps := []capabilities_registry.CapabilitiesRegistryCapability{
		{
			LabelledName:   testName,
			Version:        testVersion,
			CapabilityType: CapabilityTypeTrigger,
		},
	}

	mockRegistry.On("GetCapabilities", mock.Anything).Return(existingCaps, nil).Once()
	mockRegistry.On("GetDON", mock.Anything, testDonID).Return(testDON, nil).Once()
	mockRegistry.On("GetNodes", mock.Anything).Return(testNodes, nil).Once()
	mockRegistry.On("UpdateNodes", mock.Anything, mock.Anything).Return(mockUpdateNodesTx, nil).Once()
	mockRegistry.On("UpdateDON", mock.Anything, testDonID, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mockUpdateDONTx, nil).Once()

	require.NoError(t, createCapability(mockRegistry, testName, testVersion, testType, testDonID))

	mockRegistry.AssertExpectations(t)
}

// createCapability encapsulates the logic from createCmd.Run for testing
func createCapability(client *MockCapabilityRegistry, name, version, typeName string, donID uint32) error {
	// Check if capability already exists
	callOpts := &bind.CallOpts{}
	caps, err := client.GetCapabilities(callOpts)
	if err != nil {
		return err
	}

	capExists := false
	for _, capability := range caps {
		if capability.LabelledName == name && capability.Version == version {
			capExists = true
			break
		}
	}

	// Add capability
	if !capExists {
		txOpt := &bind.TransactOpts{}
		_, err2 := client.AddCapabilities(txOpt, []capabilities_registry.CapabilitiesRegistryCapability{
			{
				LabelledName:   name,
				Version:        version,
				CapabilityType: stringToCapabilityType(typeName),
				ResponseType:   0,
			},
		})
		if err2 != nil {
			return err
		}
	}

	don, err := client.GetDON(callOpts, donID)
	if err != nil {
		return err
	}

	// Fetch Node Data
	nodes, err := client.GetNodes(callOpts)
	if err != nil {
		return err
	}

	// Figure out which nodes we want to update
	nodesToUpdate := make([]capabilities_registry.CapabilitiesRegistryNodeParams, 0)

	for _, node := range nodes {
		// Check if the node's P2pId is contained in the don.NodeP2PIds slice
		isInDon := slices.Contains(don.NodeP2PIds, node.P2pId)
		if isInDon {
			nodesToUpdate = append(nodesToUpdate, capabilities_registry.CapabilitiesRegistryNodeParams{
				NodeOperatorId:      node.NodeOperatorId,
				Signer:              node.Signer,
				P2pId:               node.P2pId,
				EncryptionPublicKey: node.EncryptionPublicKey,
				HashedCapabilityIds: append(node.HashedCapabilityIds, getHashedCapabilityID(name, version)),
			})
		}
	}

	txOpt := &bind.TransactOpts{}
	_, err = client.UpdateNodes(txOpt, nodesToUpdate)
	if err != nil {
		return err
	}

	txOpt = &bind.TransactOpts{}
	// append cap config
	don.CapabilityConfigurations = append(don.CapabilityConfigurations, capabilities_registry.CapabilitiesRegistryCapabilityConfiguration{
		CapabilityId: getHashedCapabilityID(name, version),
		Config:       nil,
	})
	_, err = client.UpdateDON(txOpt, donID, don.NodeP2PIds, don.CapabilityConfigurations, don.IsPublic, don.F)
	if err != nil {
		return err
	}
	return nil
}

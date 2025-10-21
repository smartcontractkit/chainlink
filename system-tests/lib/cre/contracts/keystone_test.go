package contracts

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

func TestDonsOrderedByID(t *testing.T) {
	// Test donsOrderedByID sorts by id ascending
	d := dons{
		c: make(map[string]donConfig),
	}

	d.c["don3"] = donConfig{id: 3}
	d.c["don1"] = donConfig{id: 1}
	d.c["don2"] = donConfig{id: 2}

	ordered := d.donsOrderedByID()
	if len(ordered) != 3 {
		t.Fatalf("expected 3 dons, got %d", len(ordered))
	}

	if ordered[0].id != 1 || ordered[1].id != 2 || ordered[2].id != 3 {
		t.Fatalf("expected dons ordered by id 1,2,3 got %d,%d,%d", ordered[0].id, ordered[1].id, ordered[2].id)
	}
}

func TestToV2ConfigureInput(t *testing.T) {
	// Create test peer IDs
	peerID1 := p2pkey.MustNewV2XXXTestingOnly(big.NewInt(1)).PeerID().String()
	peerID2 := p2pkey.MustNewV2XXXTestingOnly(big.NewInt(2)).PeerID().String()

	// Create test dons with sample data
	d := &dons{
		c: make(map[string]donConfig),
	}

	// Add a DON with capabilities and nodes
	d.c["test-don"] = donConfig{
		id: 1,
		DonCapabilities: keystone_changeset.DonCapabilities{
			Name: "test-don",
			F:    1,
			Nops: []keystone_changeset.NOP{
				{
					Name:  "test-nop",
					Nodes: []string{peerID1, peerID2},
				},
			},
			Capabilities: []keystone_changeset.DONCapabilityWithConfig{
				{
					Capability: kcr.CapabilitiesRegistryCapability{
						LabelledName:   "test-capability",
						Version:        "1.0.0",
						CapabilityType: 1,
					},
					Config: &capabilitiespb.CapabilityConfig{},
				},
			},
		},
	}

	// Call the method under test
	result := d.mustToV2ConfigureInput(123, "0x1234567890abcdef")

	// Verify the transformation
	if result.RegistryChainSel != 123 {
		t.Errorf("expected RegistryChainSel 123, got %d", result.RegistryChainSel)
	}

	if result.ContractAddress != "0x1234567890abcdef" { //nolint:staticcheck // we won't migrate tests
		t.Errorf("expected ContractAddress 0x1234567890abcdef, got %s", result.ContractAddress) //nolint:staticcheck // we won't migrate tests
	}

	if len(result.Nops) != 1 {
		t.Fatalf("expected 1 NOP, got %d", len(result.Nops))
	}

	if result.Nops[0].Name != "test-nop" {
		t.Errorf("expected NOP name 'test-nop', got %s", result.Nops[0].Name)
	}

	if len(result.Nodes) != 2 {
		t.Fatalf("expected 2 nodes, got %d", len(result.Nodes))
	}

	if len(result.Capabilities) != 1 {
		t.Fatalf("expected 1 capability, got %d", len(result.Capabilities))
	}

	expectedCapID := "test-capability@1.0.0"
	if result.Capabilities[0].CapabilityId != expectedCapID {
		t.Errorf("expected capability ID '%s', got %s", expectedCapID, result.Capabilities[0].CapabilityId)
	}

	if len(result.DONs) != 1 {
		t.Fatalf("expected 1 DON, got %d", len(result.DONs))
	}

	if result.DONs[0].Name != "test-don" {
		t.Errorf("expected DON name 'test-don', got %s", result.DONs[0].Name)
	}

	if result.DONs[0].F != 1 {
		t.Errorf("expected DON F value 1, got %d", result.DONs[0].F)
	}

	if len(result.DONs[0].Nodes) != 2 {
		t.Errorf("expected DON to have 2 nodes, got %d", len(result.DONs[0].Nodes))
	}

	if len(result.DONs[0].CapabilityConfigurations) != 1 {
		t.Errorf("expected DON to have 1 capability configuration, got %d", len(result.DONs[0].CapabilityConfigurations))
	}
}

// TestGenerateAdminAddresses contains all the test cases for the function.
func TestGenerateAdminAddresses(t *testing.T) {
	// Test Case 1: Basic Functionality
	t.Run("Basic_Functionality_10_Addresses", func(t *testing.T) {
		count := 10
		addresses, err := generateAdminAddresses(count)
		require.NoError(t, err, "Expected no error, but got: %v", err)
		require.Len(t, len(addresses), count, "Expected slice of length %d, but got %d", count, len(addresses))

		// Check for uniqueness and validity
		addressMap := make(map[common.Address]bool)
		for _, addr := range addresses {
			require.True(t, common.IsHexAddress(addr.Hex()))
			require.NotEqual(t, 0, addr.Cmp(common.HexToAddress("0x0000000000000000000000000000000000000000")), "Generated a zero address, which should be avoided")
			addressMap[addr] = true
		}
		require.Len(t, len(addressMap), count, "Expected slice of length %d, but got %d", count, len(addressMap))
	})

	// Test Case 2: Smallest Valid Input
	t.Run("Smallest_Valid_Input_1_Address", func(t *testing.T) {
		count := 1
		addresses, err := generateAdminAddresses(count)
		require.NoError(t, err, "Expected no error, but got: %v", err)
		require.Len(t, len(addresses), count, "Expected slice of length %d, but got %d", count, len(addresses))
	})

	// Test Case 3: Invalid Input (Zero and Negative Count)
	t.Run("Invalid_Input_Zero_Count", func(t *testing.T) {
		count := 0
		_, err := generateAdminAddresses(count)
		require.Error(t, err, "Expected an error for count %d, but got none", count)
	})

	t.Run("Invalid_Input_Negative_Count", func(t *testing.T) {
		count := -5
		_, err := generateAdminAddresses(count)
		require.Error(t, err, "Expected an error for count %d, but got none", count)
	})

	// Test that 5 digit padding starts at boundary
	t.Run("Boundary_Condition_65536_Addresses", func(t *testing.T) {
		count := 65536
		addresses, err := generateAdminAddresses(count)
		require.NoError(t, err, "Expected no error, but got: %v", err)
		require.Len(t, len(addresses), count, "Expected slice of length %d, but got %d", count, len(addresses))

		for _, addr := range addresses {
			require.True(t, common.IsHexAddress(addr.String()), "invalid address: %s", addr)
		}
	})
}

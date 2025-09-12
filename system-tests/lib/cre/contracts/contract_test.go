package contracts

import (
	"math/big"
	"testing"

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
	result := d.toV2ConfigureInput(123, "0x1234567890abcdef")

	// Verify the transformation
	if result.RegistryChainSel != 123 {
		t.Errorf("expected RegistryChainSel 123, got %d", result.RegistryChainSel)
	}

	if result.ContractAddress != "0x1234567890abcdef" {
		t.Errorf("expected ContractAddress 0x1234567890abcdef, got %s", result.ContractAddress)
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

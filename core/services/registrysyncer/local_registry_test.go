package registrysyncer

import (
	"testing"

	"github.com/smartcontractkit/libocr/ragep2p/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

func TestLocalRegistry_DONsForCapability(t *testing.T) {
	lggr := logger.Test(t)
	getPeerID := func() (types.PeerID, error) {
		return [32]byte{0: 1}, nil
	}
	idsToDons := map[DonID]DON{
		1: {
			DON: capabilities.DON{
				Name: "don1",
				ID:   1,
				F:    1,
				Members: []types.PeerID{
					{0: 1},
					{0: 2},
				},
			},
			CapabilityConfigurations: map[string]CapabilityConfiguration{
				"capabilityID@1.0.0": CapabilityConfiguration{},
			},
		},
		2: {
			DON: capabilities.DON{
				Name: "don2",
				ID:   2,
				F:    2,
				Members: []types.PeerID{
					{0: 3},
					{0: 4},
				},
			},
			CapabilityConfigurations: map[string]CapabilityConfiguration{
				"secondCapabilityID@1.0.0": CapabilityConfiguration{},
			},
		},
		3: {
			DON: capabilities.DON{
				Name: "don2",
				ID:   2,
				F:    2,
				Members: []types.PeerID{
					{0: 5},
					{0: 6},
				},
			},
			CapabilityConfigurations: map[string]CapabilityConfiguration{
				"thirdCapabilityID@1.0.0": CapabilityConfiguration{},
			},
		},
	}
	idsToNodes := map[types.PeerID]NodeInfo{
		{0: 1}: {
			NodeOperatorID: 0,
		},
		{0: 2}: {
			NodeOperatorID: 1,
		},
		{0: 3}: {
			NodeOperatorID: 2,
		},
		{0: 4}: {
			NodeOperatorID: 3,
		},
	}
	idsToCapabilities := map[string]Capability{
		"capabilityID@1.0.0": {
			ID:             "capabilityID@1.0.0",
			CapabilityType: capabilities.CapabilityTypeAction,
		},
		"secondCapabilityID@1.0.0": {
			ID:             "secondCapabilityID@1.0.0",
			CapabilityType: capabilities.CapabilityTypeAction,
		},
	}
	lr := NewLocalRegistry(lggr, getPeerID, idsToDons, idsToNodes, idsToCapabilities)

	gotDons, err := lr.DONsForCapability(t.Context(), "capabilityID@1.0.0")
	require.NoError(t, err)

	assert.Len(t, gotDons, 1)
	assert.Equal(t, idsToDons[1].DON, gotDons[0].DON)

	nodes := gotDons[0].Nodes
	assert.Len(t, nodes, 2)
	assert.Equal(t, types.PeerID{0: 1}, *nodes[0].PeerID)
	assert.Equal(t, types.PeerID{0: 2}, *nodes[1].PeerID)

	// Non-existent DON
	_, err = lr.DONsForCapability(t.Context(), "nonExistentCapabilityID@1.0.0")
	require.ErrorContains(t, err, "could not find DON for capability nonExistentCapabilityID@1.0.0")

	// thirdCapability is on a DON with invalid peers
	_, err = lr.DONsForCapability(t.Context(), "thirdCapabilityID@1.0.0")
	require.ErrorContains(t, err, "could not find node for peerID")
}

package v0_5

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/channel_config_store"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/contracts/evm"
)

// Simple mock iterator
type MockChannelDefinitionIterator struct {
	events []*channel_config_store.ChannelConfigStoreNewChannelDefinition
	index  int
	err    error
}

func (m *MockChannelDefinitionIterator) Next() bool {
	m.index++
	return m.index < len(m.events)
}

func (m *MockChannelDefinitionIterator) Error() error {
	return m.err
}

func (m *MockChannelDefinitionIterator) Close() error {
	return nil
}

func (m *MockChannelDefinitionIterator) GetEvent() *channel_config_store.ChannelConfigStoreNewChannelDefinition {
	if m.index < 0 || m.index >= len(m.events) {
		return nil
	}
	return m.events[m.index]
}

// Simple mock contract
type MockChannelConfigStore struct {
	typeAndVersionValue string
	ownerValue          common.Address
	events              []*channel_config_store.ChannelConfigStoreNewChannelDefinition
	err                 error
}

func (m *MockChannelConfigStore) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	return m.typeAndVersionValue, nil
}

func (m *MockChannelConfigStore) Owner(opts *bind.CallOpts) (common.Address, error) {
	return m.ownerValue, nil
}

func (m *MockChannelConfigStore) FilterNewChannelDefinition(
	opts *bind.FilterOpts,
	donID []*big.Int,
) (evm.LogIterator[channel_config_store.ChannelConfigStoreNewChannelDefinition], error) {
	if m.err != nil {
		return nil, m.err
	}
	return &MockChannelDefinitionIterator{events: m.events}, nil
}

func TestChannelConfigStoreViewGenerator_Generate(t *testing.T) {
	// Setup test context
	ctx := context.Background()

	// Define test data
	typeAndVersion := "ChannelConfigStore 1.0.0"
	owner := common.HexToAddress("0x1234567890123456789012345678901234567890")

	// Create test events
	donID1 := big.NewInt(1)
	donID2 := big.NewInt(2)

	events := []*channel_config_store.ChannelConfigStoreNewChannelDefinition{
		{
			DonId:   donID1,
			Version: 1,
			Url:     "https://example.com/don1/v1",
			Sha:     [32]byte{1},
		},
		{
			DonId:   donID1,
			Version: 2, // Higher version for same DON ID should be used
			Url:     "https://example.com/don1/v2",
			Sha:     [32]byte{2},
		},
		{
			DonId:   donID2,
			Version: 1,
			Url:     "https://example.com/don2/v1",
			Sha:     [32]byte{3},
		},
	}

	mockContract := &MockChannelConfigStore{
		typeAndVersionValue: typeAndVersion,
		ownerValue:          owner,
		events:              events,
		err:                 nil,
	}

	generator := NewChannelConfigStoreViewGenerator(mockContract)

	// Method under test
	view, err := generator.Generate(ctx, ChannelConfigStoreViewParams{
		FromBlock: 0,
	})

	require.NoError(t, err)

	// Check basic contract properties
	require.Equal(t, typeAndVersion, view.TypeAndVersion)
	require.Equal(t, owner, view.Owner)

	// Check channel definitions map
	require.Len(t, view.ChannelDefinitions, 2, "Should have definitions for 2 DON IDs")

	// Check DON ID 1 has the latest version (2)
	donID1UInt := donID1.Uint64()
	channelDef1, exists := view.ChannelDefinitions[donID1UInt]
	require.True(t, exists, "Should have a definition for DON ID 1")
	require.Equal(t, donID1UInt, channelDef1.DonID)
	require.Equal(t, uint32(2), channelDef1.Version)
	require.Equal(t, "https://example.com/don1/v2", channelDef1.URL)
	require.Equal(t, [32]byte{2}, channelDef1.SHA)

	// Check DON ID 2 has version 1
	donID2UInt := donID2.Uint64()
	channelDef2, exists := view.ChannelDefinitions[donID2UInt]
	require.True(t, exists, "Should have a definition for DON ID 2")
	require.Equal(t, donID2UInt, channelDef2.DonID)
	require.Equal(t, uint32(1), channelDef2.Version)
	require.Equal(t, "https://example.com/don2/v1", channelDef2.URL)
	require.Equal(t, [32]byte{3}, channelDef2.SHA)
}

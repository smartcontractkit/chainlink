package v0_5

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/channel_config_store"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/contracts/evm"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/view/interfaces"
)

// ChannelDefinition represents a single channel definition
type ChannelDefinition struct {
	DonID   uint64
	Version uint32
	URL     string
	SHA     [32]byte
}
type ChannelConfigStoreView struct {
	TypeAndVersion     string                       `json:"typeAndVersion,omitempty"`
	Owner              common.Address               `json:"owner,omitempty"`
	ChannelDefinitions map[uint64]ChannelDefinition // Latest definitions keyed by DON ID
}

// ChannelConfigStoreViewParams are the parameters for generating a ChannelConfigStoreView
type ChannelConfigStoreViewParams struct {
	FromBlock uint64
	ToBlock   *uint64
}

// ChannelConfigStoreView implements the ContractView interface
var _ interfaces.ContractView = (*ChannelConfigStoreView)(nil)

// SerializeView serializes view to JSON
func (v ChannelConfigStoreView) SerializeView() (string, error) {
	bytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal contract view: %w", err)
	}
	return string(bytes), nil
}

// ChannelConfigStoreContract defines the minimal interface needed for ChannelConfigStoreViewGenerator
type ChannelConfigStoreContract interface {
	// Call methods
	TypeAndVersion(opts *bind.CallOpts) (string, error)
	Owner(opts *bind.CallOpts) (common.Address, error)

	FilterNewChannelDefinition(opts *bind.FilterOpts, donID []*big.Int) (evm.LogIterator[channel_config_store.ChannelConfigStoreNewChannelDefinition], error)
}

// ChannelConfigStoreViewGenerator generates views of ChannelConfigStore contracts
type ChannelConfigStoreViewGenerator struct {
	channelConfigStore ChannelConfigStoreContract
}

// NewChannelConfigStoreViewGenerator creates a new ChannelConfigStoreViewGenerator
func NewChannelConfigStoreViewGenerator(channelConfigStore ChannelConfigStoreContract) *ChannelConfigStoreViewGenerator {
	return &ChannelConfigStoreViewGenerator{
		channelConfigStore: channelConfigStore,
	}
}

// Generate creates a ChannelConfigStoreView from the given parameters
func (g *ChannelConfigStoreViewGenerator) Generate(ctx context.Context, params ChannelConfigStoreViewParams) (ChannelConfigStoreView, error) {
	// Initialize the view with empty maps
	view := ChannelConfigStoreView{
		ChannelDefinitions: make(map[uint64]ChannelDefinition),
	}

	if err := g.fetchContractState(ctx, &view); err != nil {
		return view, err
	}

	if err := g.processChannelDefinitionEvents(ctx, params, &view); err != nil {
		return view, err
	}

	return view, nil
}

// fetchContractState retrieves the current state of the contract
func (g *ChannelConfigStoreViewGenerator) fetchContractState(ctx context.Context, view *ChannelConfigStoreView) error {
	callOpts := &bind.CallOpts{Context: ctx}
	var err error

	view.TypeAndVersion, err = g.channelConfigStore.TypeAndVersion(callOpts)
	if err != nil {
		return fmt.Errorf("failed to get type and version: %w", err)
	}
	view.Owner, err = g.channelConfigStore.Owner(callOpts)
	if err != nil {
		return fmt.Errorf("failed to get contract owner: %w", err)
	}

	return nil
}

// processChannelDefinitionEvents processes NewChannelDefinition events
func (g *ChannelConfigStoreViewGenerator) processChannelDefinitionEvents(
	ctx context.Context,
	params ChannelConfigStoreViewParams,
	view *ChannelConfigStoreView,
) error {
	filterOpts := &bind.FilterOpts{
		Start:   params.FromBlock,
		End:     params.ToBlock,
		Context: ctx,
	}

	iter, err := g.channelConfigStore.FilterNewChannelDefinition(filterOpts, nil)
	if err != nil {
		return fmt.Errorf("failed to filter NewChannelDefinition events: %w", err)
	}
	defer iter.Close()

	// keep only the most recent definition for each DON ID
	latestVersions := make(map[uint64]uint32)

	for iter.Next() {
		event := iter.GetEvent()

		donID := event.DonId.Uint64()

		// Check if this is a newer version for this DON ID
		currentVersion, exists := latestVersions[donID]
		if !exists || event.Version > currentVersion {
			latestVersions[donID] = event.Version

			// Update the channel definition
			view.ChannelDefinitions[donID] = ChannelDefinition{
				DonID:   donID,
				Version: event.Version,
				URL:     event.Url,
				SHA:     event.Sha,
			}
		}
	}

	if err := iter.Error(); err != nil {
		return fmt.Errorf("error iterating over NewChannelDefinition events: %w", err)
	}

	return nil
}

package v0_5

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/deployment/data-streams/contracts/evm"
	dsutil "github.com/smartcontractkit/chainlink/deployment/data-streams/utils"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/view/interfaces"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/verifier_v0_5_0"
)

// VerifierState represents a single verifier configuration state
type VerifierState struct {
	ConfigDigest string   `json:"configDigest"`
	IsActive     bool     `json:"isActive"`
	F            uint8    `json:"f"`
	Signers      []string `json:"signers"`
}

// VerifierView represents a simplified view of the verifier contract state
type VerifierView struct {
	Configs        map[string]*VerifierState `json:"configs"`
	TypeAndVersion string                    `json:"typeAndVersion,omitempty"`
	Owner          string                    `json:"owner,omitempty"`
}

// Ensure VerifierView implements the ContractView interface
var _ interfaces.ContractView = (*VerifierView)(nil)

// NewVerifierView creates a new empty VerifierView
func NewVerifierView() *VerifierView {
	return &VerifierView{
		Configs: make(map[string]*VerifierState),
	}
}

// SerializeView serializes the VerifierView to JSON
func (v VerifierView) SerializeView() (string, error) {
	bytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal verifier view: %w", err)
	}
	return string(bytes), nil
}

// GetVerifierState returns the VerifierState for a specific configDigest
func (v VerifierView) GetVerifierState(configDigest string) (*VerifierState, error) {
	state, ok := v.Configs[configDigest]
	if !ok {
		return nil, fmt.Errorf("configDigest %s not found", configDigest)
	}
	return state, nil
}

// VerifierViewGenerator builds views for the Verifier contract
type VerifierViewGenerator struct {
	contract VerifierContractReader
}

// VerifierContract defines a minimal interface
type VerifierContractReader interface {
	// Call methods
	Owner(opts *bind.CallOpts) (common.Address, error)
	TypeAndVersion(opts *bind.CallOpts) (string, error)

	// Filter methods
	FilterConfigSet(opts *bind.FilterOpts, configDigest [][32]byte) (evm.LogIterator[verifier_v0_5_0.VerifierConfigSet], error)
	FilterConfigUpdated(opts *bind.FilterOpts, configDigest [][32]byte) (evm.LogIterator[verifier_v0_5_0.VerifierConfigUpdated], error)
	FilterConfigActivated(opts *bind.FilterOpts, configDigest [][32]byte) (evm.LogIterator[verifier_v0_5_0.VerifierConfigActivated], error)
	FilterConfigDeactivated(opts *bind.FilterOpts, configDigest [][32]byte) (evm.LogIterator[verifier_v0_5_0.VerifierConfigDeactivated], error)
}

func NewVerifierViewGenerator(contract VerifierContractReader) *VerifierViewGenerator {
	return &VerifierViewGenerator{
		contract: contract,
	}
}

// VerifierViewGenerator implements ContractViewGenerator
var _ interfaces.ContractViewGenerator[VerifierViewParams, VerifierView] = (*VerifierViewGenerator)(nil)

type VerifierViewParams struct {
	FromBlock uint64
	ToBlock   *uint64
}

// Generate builds a view of the Verifier contract state from logs and calls
func (b *VerifierViewGenerator) Generate(ctx context.Context, verifierContext VerifierViewParams) (VerifierView, error) {
	view := NewVerifierView()

	// Get contract owner
	owner, err := b.contract.Owner(&bind.CallOpts{Context: ctx})
	if err != nil {
		return VerifierView{}, fmt.Errorf("failed to get contract owner: %w", err)
	}
	view.Owner = owner.Hex()

	// Define the filter options
	filterOpts := &bind.FilterOpts{
		Start:   verifierContext.FromBlock,
		End:     verifierContext.ToBlock,
		Context: ctx,
	}

	// Process all ConfigSet events
	if err := b.processConfigSetEvents(filterOpts, view); err != nil {
		return VerifierView{}, err
	}

	// Process all ConfigUpdated events
	if err := b.processConfigUpdatedEvents(filterOpts, view); err != nil {
		return VerifierView{}, err
	}

	// Process all ConfigActivated events
	if err := b.processConfigActivatedEvents(filterOpts, view); err != nil {
		return VerifierView{}, err
	}

	// Process all ConfigDeactivated events
	if err := b.processConfigDeactivatedEvents(filterOpts, view); err != nil {
		return VerifierView{}, err
	}

	return *view, nil
}

// processConfigSetEvents processes ConfigSet events
func (b *VerifierViewGenerator) processConfigSetEvents(opts *bind.FilterOpts, view *VerifierView) error {
	iter, err := b.contract.FilterConfigSet(opts, nil)
	if err != nil {
		return fmt.Errorf("failed to filter ConfigSet events: %w", err)
	}
	defer iter.Close()

	for iter.Next() {
		event := iter.GetEvent()
		configDigestHex := dsutil.HexEncodeBytes(event.ConfigDigest[:])

		state := &VerifierState{
			ConfigDigest: configDigestHex,
			IsActive:     true, // New configs are active by default
			F:            event.F,
			Signers:      make([]string, 0, len(event.Signers)),
		}

		// Add signers
		for _, signer := range event.Signers {
			state.Signers = append(state.Signers, signer.String())
		}

		view.Configs[configDigestHex] = state
	}

	return nil
}

// processConfigUpdatedEvents processes ConfigUpdated events
func (b *VerifierViewGenerator) processConfigUpdatedEvents(opts *bind.FilterOpts, view *VerifierView) error {
	iter, err := b.contract.FilterConfigUpdated(opts, nil)
	if err != nil {
		return fmt.Errorf("failed to filter ConfigUpdated events: %w", err)
	}
	defer iter.Close()

	for iter.Next() {
		event := iter.GetEvent()
		configDigestHex := dsutil.HexEncodeBytes(event.ConfigDigest[:])

		// Skip if this configDigest doesn't exist yet
		state, exists := view.Configs[configDigestHex]
		if !exists {
			// This is unexpected, but we'll create a new state for this config
			state = &VerifierState{
				ConfigDigest: configDigestHex,
				IsActive:     true,
				Signers:      make([]string, 0),
			}

			view.Configs[configDigestHex] = state
		}

		// Update signers with new set
		state.Signers = make([]string, 0, len(event.NewSigners))

		// Add new signers
		for _, signer := range event.NewSigners {
			state.Signers = append(state.Signers, signer.String())
		}
	}

	return nil
}

// processConfigActivatedEvents processes ConfigActivated events
func (b *VerifierViewGenerator) processConfigActivatedEvents(opts *bind.FilterOpts, view *VerifierView) error {
	iter, err := b.contract.FilterConfigActivated(opts, nil)
	if err != nil {
		return fmt.Errorf("failed to filter ConfigActivated events: %w", err)
	}
	defer iter.Close()

	for iter.Next() {
		event := iter.GetEvent()
		configDigestHex := dsutil.HexEncodeBytes(event.ConfigDigest[:])

		// Skip if this configDigest doesn't exist yet
		state, exists := view.Configs[configDigestHex]
		if !exists {
			continue
		}

		state.IsActive = true
	}

	return nil
}

// processConfigDeactivatedEvents processes ConfigDeactivated events
func (b *VerifierViewGenerator) processConfigDeactivatedEvents(opts *bind.FilterOpts, view *VerifierView) error {
	iter, err := b.contract.FilterConfigDeactivated(opts, nil)
	if err != nil {
		return fmt.Errorf("failed to filter ConfigDeactivated events: %w", err)
	}
	defer iter.Close()

	for iter.Next() {
		event := iter.GetEvent()
		configDigestHex := dsutil.HexEncodeBytes(event.ConfigDigest[:])

		// Skip if this configDigest doesn't exist yet
		state, exists := view.Configs[configDigestHex]
		if !exists {
			continue
		}

		state.IsActive = false
	}

	return nil
}

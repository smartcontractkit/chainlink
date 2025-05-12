package v0_5

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	verifier_proxy "github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/verifier_proxy_v0_5_0"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/contracts/evm"
	dsutil "github.com/smartcontractkit/chainlink/deployment/data-streams/utils"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/view/interfaces"
)

// VerifierProxyView represents the state of a VerifierProxy contract
type VerifierProxyView struct {
	Owner                string            `json:"owner,omitempty"`
	FeeManager           string            `json:"feeManager,omitempty"`
	AccessController     string            `json:"accessController,omitempty"`
	TypeAndVersion       string            `json:"typeAndVersion,omitempty"`
	InitializedVerifiers []string          `json:"initializedVerifiers"`
	VerifiersByDigest    map[string]string `json:"verifiersByDigest"`
}

// VerifierProxyView implements the ContractView interface
var _ interfaces.ContractView = (*VerifierProxyView)(nil)

// SerializeView serializes view to JSON
func (v VerifierProxyView) SerializeView() (string, error) {
	bytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal contract view: %w", err)
	}
	return string(bytes), nil
}

type VerifierProxyViewParams struct {
	FromBlock uint64
	ToBlock   *uint64
}

// VerifierViewGenerator implements ContractViewGenerator
var _ interfaces.ContractViewGenerator[VerifierProxyViewParams, VerifierProxyView] = (*VerifierProxyViewGenerator)(nil)

// VerifierProxyContract defines a minimal interface
type VerifierProxyContract interface {
	// Call methods
	Owner(opts *bind.CallOpts) (common.Address, error)
	SAccessController(opts *bind.CallOpts) (common.Address, error)
	SFeeManager(opts *bind.CallOpts) (common.Address, error)
	TypeAndVersion(opts *bind.CallOpts) (string, error)

	// Event filters
	FilterVerifierInitialized(opts *bind.FilterOpts) (evm.LogIterator[verifier_proxy.VerifierProxyVerifierInitialized], error)
	FilterVerifierSet(opts *bind.FilterOpts) (evm.LogIterator[verifier_proxy.VerifierProxyVerifierSet], error)
	FilterVerifierUnset(opts *bind.FilterOpts) (evm.LogIterator[verifier_proxy.VerifierProxyVerifierUnset], error)
}

// VerifierProxyViewGenerator generates views of VerifierProxy contracts
type VerifierProxyViewGenerator struct {
	verifierProxy VerifierProxyContract
}

// NewVerifierProxyViewGenerator creates a new VerifierProxyViewGenerator
func NewVerifierProxyViewGenerator(verifierProxy VerifierProxyContract) *VerifierProxyViewGenerator {
	return &VerifierProxyViewGenerator{
		verifierProxy: verifierProxy,
	}
}

// Generate creates a VerifierProxyView from the given parameters
func (v *VerifierProxyViewGenerator) Generate(ctx context.Context, params VerifierProxyViewParams) (VerifierProxyView, error) {
	// Initialize the view with empty maps
	view := VerifierProxyView{
		InitializedVerifiers: []string{},
		VerifiersByDigest:    make(map[string]string),
	}

	// Create filter options
	filterOpts := &bind.FilterOpts{
		Start:   params.FromBlock,
		End:     params.ToBlock,
		Context: ctx,
	}

	// Get contract state data
	if err := v.fetchContractState(ctx, &view); err != nil {
		return view, err
	}

	if err := v.processInitializedVerifiers(filterOpts, &view); err != nil {
		return view, err
	}

	if err := v.processVerifierSetEvents(filterOpts, &view); err != nil {
		return view, err
	}

	if err := v.processVerifierUnsetEvents(filterOpts, &view); err != nil {
		return view, err
	}

	return view, nil
}

// fetchContractState retrieves the current state of the contract
func (v *VerifierProxyViewGenerator) fetchContractState(ctx context.Context, view *VerifierProxyView) error {
	callOpts := &bind.CallOpts{Context: ctx}
	var err error

	owner, err := v.verifierProxy.Owner(callOpts)
	if err != nil {
		return fmt.Errorf("failed to get owner: %w", err)
	}
	view.Owner = owner.Hex()

	// Get the AccessController
	ac, err := v.verifierProxy.SAccessController(callOpts)
	if err != nil {
		return fmt.Errorf("failed to get access controller: %w", err)
	}
	view.AccessController = ac.Hex()

	// Get the FeeManager
	fm, err := v.verifierProxy.SFeeManager(callOpts)
	if err != nil {
		return fmt.Errorf("failed to get fee manager: %w", err)
	}
	view.FeeManager = fm.Hex()

	// Get TypeAndVersion
	view.TypeAndVersion, err = v.verifierProxy.TypeAndVersion(callOpts)
	if err != nil {
		return fmt.Errorf("failed to get type and version: %w", err)
	}

	return nil
}

// processInitializedVerifiers processes VerifierInitialized events
func (v *VerifierProxyViewGenerator) processInitializedVerifiers(filterOpts *bind.FilterOpts, view *VerifierProxyView) error {
	initializedIter, err := v.verifierProxy.FilterVerifierInitialized(filterOpts)
	if err != nil {
		return fmt.Errorf("failed to filter VerifierInitialized events: %w", err)
	}
	defer initializedIter.Close()

	for initializedIter.Next() {
		event := initializedIter.GetEvent()
		view.InitializedVerifiers = append(view.InitializedVerifiers, event.VerifierAddress.Hex())
	}

	return nil
}

// processVerifierSetEvents processes VerifierSet events
func (v *VerifierProxyViewGenerator) processVerifierSetEvents(filterOpts *bind.FilterOpts, view *VerifierProxyView) error {
	setIter, err := v.verifierProxy.FilterVerifierSet(filterOpts)
	if err != nil {
		return fmt.Errorf("failed to filter VerifierSet events: %w", err)
	}
	defer setIter.Close()

	for setIter.Next() {
		event := setIter.GetEvent()
		digest := dsutil.HexEncodeBytes32(event.NewConfigDigest)
		view.VerifiersByDigest[digest] = event.VerifierAddress.Hex()
	}

	return nil
}

// processVerifierUnsetEvents processes VerifierUnset events
func (v *VerifierProxyViewGenerator) processVerifierUnsetEvents(filterOpts *bind.FilterOpts, view *VerifierProxyView) error {
	unsetIter, err := v.verifierProxy.FilterVerifierUnset(filterOpts)
	if err != nil {
		return fmt.Errorf("failed to filter VerifierUnset events: %w", err)
	}
	defer unsetIter.Close()

	for unsetIter.Next() {
		event := unsetIter.GetEvent()
		digest := dsutil.HexEncodeBytes32(event.ConfigDigest)
		delete(view.VerifiersByDigest, digest)
	}

	return nil
}

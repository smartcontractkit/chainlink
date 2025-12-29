package adapters

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/aptos-labs/aptos-go-sdk"
	chainsel "github.com/smartcontractkit/chain-selectors"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cldf_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-ccip/deployment/fastcurse"
	"github.com/smartcontractkit/chainlink-ccip/deployment/utils/sequences"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/dependency"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/operation"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/globals"
	aptosops "github.com/smartcontractkit/chainlink/deployment/ccip/operation/aptos"
	aptosstateview "github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/aptos"
)

type CurseAdapter struct {
	CCIPAddress aptos.AccountAddress
}

func NewCurseAdapter() *CurseAdapter {
	return &CurseAdapter{}
}

// Initialize performs any required setup. No-op for now.
func (c *CurseAdapter) Initialize(e cldf.Environment, selector uint64) error {
	stateMap, err := aptosstateview.LoadOnchainStateAptos(e)
	if err != nil {
		return fmt.Errorf("failed to load Aptos onchain state: %w", err)
	}

	state, ok := stateMap[selector]
	if !ok {
		return fmt.Errorf("aptos chain %d not found in state", selector)
	}
	c.CCIPAddress = state.CCIPAddress
	return nil
}

// IsSubjectCursedOnChain currently returns false; extend with real RMN checks if needed.
func (c *CurseAdapter) IsSubjectCursedOnChain(e cldf.Environment, selector uint64, subject fastcurse.Subject) (bool, error) {
	chain, ok := e.BlockChains.AptosChains()[selector]
	if !ok {
		return false, fmt.Errorf("aptos chain %d not found in environment", selector)
	}
	deps := dependency.AptosDeps{
		AptosChain: chain,
	}
	return aptosops.IsSubjectCursed(deps, c.CCIPAddress, subject[:])
}

// IsChainConnectedToTargetChain returns true by default; tighten when connectivity checks are available.
func (c *CurseAdapter) IsChainConnectedToTargetChain(e cldf.Environment, selector uint64, targetSelector uint64) (bool, error) {
	return true, nil
}

// IsCurseEnabledForChain returns true by default; extend to verify RMN presence on-chain.
func (c *CurseAdapter) IsCurseEnabledForChain(cldf.Environment, uint64) (bool, error) {
	return true, nil
}

func (c *CurseAdapter) SubjectToSelector(subject fastcurse.Subject) (uint64, error) {
	return fastcurse.GenericSubjectToSelector(subject)
}

// Curse returns nil for now; plug in a real sequence when available.
func (c *CurseAdapter) Curse() *cldf_ops.Sequence[fastcurse.CurseInput, sequences.OnChainOutput, cldf_chain.BlockChains] {
	return cldf_ops.NewSequence(
		"aptos-curse-sequence",
		operation.Version1_0_0,
		"Curse sequence for Aptos",
		func(b cldf_ops.Bundle, deps cldf_chain.BlockChains, in fastcurse.CurseInput) (output sequences.OnChainOutput, err error) {
			return sequences.OnChainOutput{}, nil
		},
	)
}

// Uncurse returns nil for now; plug in a real sequence when available.
func (c *CurseAdapter) Uncurse() *cldf_ops.Sequence[fastcurse.CurseInput, sequences.OnChainOutput, cldf_chain.BlockChains] {
	return cldf_ops.NewSequence(
		"aptos-uncurse-sequence",
		operation.Version1_0_0,
		"Uncurse sequence for Aptos",
		func(b cldf_ops.Bundle, deps cldf_chain.BlockChains, in fastcurse.CurseInput) (output sequences.OnChainOutput, err error) {
			return sequences.OnChainOutput{}, nil
		},
	)
}

func (c *CurseAdapter) ListConnectedChains(e cldf.Environment, selector uint64) ([]uint64, error) {
	_, ok := e.BlockChains.AptosChains()[selector]
	if !ok {
		return nil, fmt.Errorf("chain %d not found", selector)
	}

	return []uint64{}, nil
}

// SelectorToSubject converts selector to Subject (family-aware).
func (c *CurseAdapter) SelectorToSubject(selector uint64) fastcurse.Subject {
	return globals.FamilyAwareSelectorToSubject(selector, chainsel.FamilyAptos)
}

func (c *CurseAdapter) DeriveCurseAdapterVersion(e cldf.Environment, selector uint64) (*semver.Version, error) {
	return semver.MustParse("1.6.0"), nil
}

// Ensure interfaces are satisfied at compile time.
var (
	_ fastcurse.CurseAdapter        = (*CurseAdapter)(nil)
	_ fastcurse.CurseSubjectAdapter = (*CurseAdapter)(nil)
)

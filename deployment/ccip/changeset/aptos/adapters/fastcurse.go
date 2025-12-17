package adapters

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	chainsel "github.com/smartcontractkit/chain-selectors"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cldf_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-ccip/deployment/fastcurse"
	"github.com/smartcontractkit/chainlink-ccip/deployment/utils/sequences"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/operation"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/globals"
)

type CurseAdapter struct{}

type CurseSubjectAdapter struct{}

func NewCurseAdapter() *CurseAdapter {
	return &CurseAdapter{}
}

func NewCurseSubjectAdapter() *CurseSubjectAdapter {
	return &CurseSubjectAdapter{}
}

// Initialize performs any required setup. No-op for now.
func (c *CurseAdapter) Initialize(cldf.Environment, uint64) error { return nil }

// IsSubjectCursedOnChain currently returns false; extend with real RMN checks if needed.
func (c *CurseAdapter) IsSubjectCursedOnChain(cldf.Environment, uint64, fastcurse.Subject) (bool, error) {
	return false, nil
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
func (c *CurseSubjectAdapter) SelectorToSubject(selector uint64) fastcurse.Subject {
	return globals.FamilyAwareSelectorToSubject(selector, chainsel.FamilyAptos)
}

func (c *CurseSubjectAdapter) DeriveCurseAdapterVersion(e cldf.Environment, selector uint64) (*semver.Version, error) {
	return semver.MustParse("1.6.0"), nil
}

// Ensure interfaces are satisfied at compile time.
var (
	_ fastcurse.CurseAdapter        = (*CurseAdapter)(nil)
	_ fastcurse.CurseSubjectAdapter = (*CurseSubjectAdapter)(nil)
)

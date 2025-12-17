package adapters

import (
	chainsel "github.com/smartcontractkit/chain-selectors"
	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cldf_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-ccip/deployment/fastcurse"
	"github.com/smartcontractkit/chainlink-ccip/deployment/utils/sequences"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/globals"
)

// CurseAdapter implements fastcurse.CurseAdapter for Aptos.
// These are minimal adapters that make the registry operational; they can be
// enhanced later to perform real on-chain checks or invoke real sequences.
type CurseAdapter struct{}

type CurseSubjectAdapter struct{}

func NewCurseAdapter() *CurseAdapter {
	return &CurseAdapter{}
}

func NewCurseSubjectAdapter() *CurseSubjectAdapter {
	return &CurseSubjectAdapter{}
}

// Initialize performs any required setup. No-op for now.
func (c *CurseAdapter) Initialize(cldf.Environment) error { return nil }

// IsSubjectCursedOnChain currently returns false; extend with real RMN checks if needed.
func (c *CurseAdapter) IsSubjectCursedOnChain(cldf.Environment, uint64, fastcurse.Subject) (bool, error) {
	return false, nil
}

// IsChainConnectedToTargetChain returns true by default; tighten when connectivity checks are available.
func (c *CurseAdapter) IsChainConnectedToTargetChain(cldf.Environment, uint64, uint64) (bool, error) {
	return true, nil
}

// IsCurseEnabledForChain returns true by default; extend to verify RMN presence on-chain.
func (c *CurseAdapter) IsCurseEnabledForChain(cldf.Environment, uint64) (bool, error) {
	return true, nil
}

// SubjectToSelector converts a Subject to a selector using family-aware helper.
func (c *CurseAdapter) SubjectToSelector(subject fastcurse.Subject) uint64 {
	return globals.FamilyAwareSubjectToSelector(subject, chainsel.FamilyAptos)
}

// Curse returns nil for now; plug in a real sequence when available.
func (c *CurseAdapter) Curse() *cldf_ops.Sequence[fastcurse.CurseInput, sequences.OnChainOutput, cldf_chain.BlockChains] {
	return nil
}

// Uncurse returns nil for now; plug in a real sequence when available.
func (c *CurseAdapter) Uncurse() *cldf_ops.Sequence[fastcurse.CurseInput, sequences.OnChainOutput, cldf_chain.BlockChains] {
	return nil
}

// SelectorToSubject converts selector to Subject (family-aware).
func (c *CurseSubjectAdapter) SelectorToSubject(selector uint64) fastcurse.Subject {
	s := globals.FamilyAwareSelectorToSubject(selector, chainsel.FamilyAptos)
	return fastcurse.Subject(s)
}

// Ensure interfaces are satisfied at compile time.
var (
	_ fastcurse.CurseAdapter        = (*CurseAdapter)(nil)
	_ fastcurse.CurseSubjectAdapter = (*CurseSubjectAdapter)(nil)
)

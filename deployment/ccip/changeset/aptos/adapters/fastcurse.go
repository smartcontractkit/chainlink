package adapters

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/aptos-labs/aptos-go-sdk"
	mcmstypes "github.com/smartcontractkit/mcms/types"
	"github.com/smartcontractkit/chainlink-aptos/bindings/bind"
	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip_router"
	"github.com/smartcontractkit/chainlink-ccip/deployment/fastcurse"
	"github.com/smartcontractkit/chainlink-ccip/deployment/utils/sequences"
	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cldf_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	chainsel "github.com/smartcontractkit/chain-selectors"

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

func (c *CurseAdapter) IsChainConnectedToTargetChain(e cldf.Environment, selector uint64, targetSelector uint64) (bool, error) {
	chain, ok := e.BlockChains.AptosChains()[selector]
	if !ok {
		return false, fmt.Errorf("aptos chain %d not found in environment", selector)
	}
	routerBind := ccip_router.Bind(c.CCIPAddress, chain.Client)
	callOpts := &bind.CallOpts{}
	connected, err := routerBind.Router().IsChainSupported(callOpts, targetSelector)
	if err != nil {
		return false, fmt.Errorf("failed to check if chain %d is connected to chain %d: %w", selector, targetSelector, err)
	}
	return connected, nil
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
		func(b cldf_ops.Bundle, chains cldf_chain.BlockChains, in fastcurse.CurseInput) (output sequences.OnChainOutput, err error) {
			chain, ok := chains.AptosChains()[in.ChainSelector]
			if !ok {
				return sequences.OnChainOutput{}, fmt.Errorf("chain with selector %d not found in environment", in.ChainSelector)
			}
			subjectBytes := make([][]byte, len(in.Subjects))
			for i, subject := range in.Subjects {
				subjectBytes[i] = subject[:]
			}
			curseInput := aptosops.CurseMultipleInput{
				CCIPAddress: c.CCIPAddress,
				Subjects:    subjectBytes,
			}
			deps := dependency.AptosDeps{
				AptosChain: chain,
			}
			report, err := cldf_ops.ExecuteOperation(b, aptosops.CurseMultipleOp, deps, curseInput)
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("failed to execute curse operation on Aptos chain %d: %w", chain.Selector, err)
			}
			batchOperation := mcmstypes.BatchOperation{
				ChainSelector: mcmstypes.ChainSelector(chain.Selector),
				Transactions:  []mcmstypes.Transaction{report.Output},
			}
			return sequences.OnChainOutput{
				BatchOps: []mcmstypes.BatchOperation{batchOperation},
			}, nil
		},
	)
}

// Uncurse returns nil for now; plug in a real sequence when available.
func (c *CurseAdapter) Uncurse() *cldf_ops.Sequence[fastcurse.CurseInput, sequences.OnChainOutput, cldf_chain.BlockChains] {
	return cldf_ops.NewSequence(
		"aptos-uncurse-sequence",
		operation.Version1_0_0,
		"Uncurse sequence for Aptos",
		func(b cldf_ops.Bundle, chains cldf_chain.BlockChains, in fastcurse.CurseInput) (output sequences.OnChainOutput, err error) {
			chain, ok := chains.AptosChains()[in.ChainSelector]
			if !ok {
				return sequences.OnChainOutput{}, fmt.Errorf("chain with selector %d not found in environment", in.ChainSelector)
			}
			subjectBytes := make([][]byte, len(in.Subjects))
			for i, subject := range in.Subjects {
				subjectBytes[i] = subject[:]
			}
			uncurseInput := aptosops.UncurseMultipleInput{
				CCIPAddress: c.CCIPAddress,
				Subjects:    subjectBytes,
			}
			deps := dependency.AptosDeps{
				AptosChain: chain,
			}
			report, err := cldf_ops.ExecuteOperation(b, aptosops.UncurseMultipleOp, deps, uncurseInput)
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("failed to execute uncurse operation on Aptos chain %d: %w", chain.Selector, err)
			}
			batchOperation := mcmstypes.BatchOperation{
				ChainSelector: mcmstypes.ChainSelector(chain.Selector),
				Transactions:  []mcmstypes.Transaction{report.Output},
			}
			return sequences.OnChainOutput{
				BatchOps: []mcmstypes.BatchOperation{batchOperation},
			}, nil
		},
	)
}

func (c *CurseAdapter) ListConnectedChains(e cldf.Environment, selector uint64) ([]uint64, error) {
	chain, ok := e.BlockChains.AptosChains()[selector]
	if !ok {
		return nil, fmt.Errorf("aptos chain %d not found in environment", selector)
	}
	routerBind := ccip_router.Bind(c.CCIPAddress, chain.Client)
	callOpts := &bind.CallOpts{}
	connectedChains, err := routerBind.Router().GetDestChains(callOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to get connected chains for chain %d: %w", selector, err)
	}
	return connectedChains, nil
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

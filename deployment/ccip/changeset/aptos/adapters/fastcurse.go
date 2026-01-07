package adapters

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/aptos-labs/aptos-go-sdk"
	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-aptos/bindings/bind"
	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip_router"
	"github.com/smartcontractkit/chainlink-ccip/deployment/fastcurse"
	"github.com/smartcontractkit/chainlink-ccip/deployment/utils/sequences"
	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cldf_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/dependency"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/operation"
	aptosseq "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/sequence"
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

// IsSubjectCursedOnChain checks if the subject is cursed on Aptos chain
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

// IsChainConnectedToTargetChain checks if the target chain is supported on Aptos chain
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

// IsCurseEnabledForChain returns true because Aptos is live on CCIP 1.6
func (c *CurseAdapter) IsCurseEnabledForChain(cldf.Environment, uint64) (bool, error) {
	return true, nil
}

// SubjectToSelector converts subject to chainselector
func (c *CurseAdapter) SubjectToSelector(subject fastcurse.Subject) (uint64, error) {
	return fastcurse.GenericSubjectToSelector(subject)
}

// Curse returns the sequence to curse multiple subjects on Aptos chain
func (c *CurseAdapter) Curse() *cldf_ops.Sequence[fastcurse.CurseInput, sequences.OnChainOutput, cldf_chain.BlockChains] {
	return cldf_ops.NewSequence(
		aptosseq.AptosCurseSequence.ID(),
		operation.Version1_0_0,
		aptosseq.AptosCurseSequence.Description(),
		func(b cldf_ops.Bundle, chains cldf_chain.BlockChains, in fastcurse.CurseInput) (output sequences.OnChainOutput, err error) {
			aptosInput := aptosseq.AptosCurseUncurseInput{
				CCIPAddress:   c.CCIPAddress,
				ChainSelector: in.ChainSelector,
				Subjects:      in.Subjects,
			}
			seqReport, err := cldf_ops.ExecuteSequence(b, aptosseq.AptosCurseSequence, chains, aptosInput)
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("failed to execute curse sequence on Aptos chain %d: %w", in.ChainSelector, err)
			}
			return seqReport.Output, nil
		},
	)
}

// Uncurse returns the sequence to uncurse multiple subjects on Aptos chain
func (c *CurseAdapter) Uncurse() *cldf_ops.Sequence[fastcurse.CurseInput, sequences.OnChainOutput, cldf_chain.BlockChains] {
	return cldf_ops.NewSequence(
		aptosseq.AptosUncurseSequence.ID(),
		operation.Version1_0_0,
		aptosseq.AptosUncurseSequence.Description(),
		func(b cldf_ops.Bundle, chains cldf_chain.BlockChains, in fastcurse.CurseInput) (output sequences.OnChainOutput, err error) {
			aptosInput := aptosseq.AptosCurseUncurseInput{
				CCIPAddress:   c.CCIPAddress,
				ChainSelector: in.ChainSelector,
				Subjects:      in.Subjects,
			}
			seqReport, err := cldf_ops.ExecuteSequence(b, aptosseq.AptosUncurseSequence, chains, aptosInput)
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("failed to execute uncurse sequence on Aptos chain %d: %w", in.ChainSelector, err)
			}
			return seqReport.Output, nil
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

// SelectorToSubject converts selector to Subject.
func (c *CurseAdapter) SelectorToSubject(selector uint64) fastcurse.Subject {
	return globals.FamilyAwareSelectorToSubject(selector, chainsel.FamilyAptos)
}

func (c *CurseAdapter) DeriveCurseAdapterVersion(e cldf.Environment, selector uint64) (*semver.Version, error) {
	return semver.MustParse("1.6.0"), nil
}

var (
	_ fastcurse.CurseAdapter        = (*CurseAdapter)(nil)
	_ fastcurse.CurseSubjectAdapter = (*CurseAdapter)(nil)
)

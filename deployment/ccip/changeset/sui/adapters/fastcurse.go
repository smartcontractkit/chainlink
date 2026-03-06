package adapters

import (
	"context"
	"fmt"

	"github.com/Masterminds/semver/v3"
	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_rmn_remote "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/rmn_remote"
	module_router "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_router"
	suistate "github.com/smartcontractkit/chainlink-sui/deployment"

	"github.com/smartcontractkit/chainlink-ccip/deployment/fastcurse"
	"github.com/smartcontractkit/chainlink-ccip/deployment/utils/sequences"
	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cldf_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/globals"
	suiseq "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/sui/sequence"
)

var (
	_ fastcurse.CurseAdapter        = (*CurseAdapter)(nil)
	_ fastcurse.CurseSubjectAdapter = (*CurseAdapter)(nil)
)

type CurseAdapter struct {
	CCIPAddress          string
	CCIPObjectRef        string
	CCIPOwnerCapObjectId string
	RouterAddress        string
	RouterStateObjectID  string
}

func NewCurseAdapter() *CurseAdapter {
	return &CurseAdapter{}
}

func (c *CurseAdapter) Initialize(e cldf.Environment, selector uint64) error {
	stateMap, err := suistate.LoadOnchainStatesui(e)
	if err != nil {
		return fmt.Errorf("failed to load SUI onchain state: %w", err)
	}

	state, ok := stateMap[selector]
	if !ok {
		return fmt.Errorf("SUI chain %d not found in state", selector)
	}
	c.CCIPAddress = state.CCIPAddress
	c.CCIPObjectRef = state.CCIPObjectRef
	c.CCIPOwnerCapObjectId = state.CCIPOwnerCapObjectId
	c.RouterAddress = state.CCIPRouterAddress
	c.RouterStateObjectID = state.CCIPRouterStateObjectID
	return nil
}

func (c *CurseAdapter) IsSubjectCursedOnChain(e cldf.Environment, selector uint64, subject fastcurse.Subject) (bool, error) {
	chain, ok := e.BlockChains.SuiChains()[selector]
	if !ok {
		return false, fmt.Errorf("SUI chain %d not found in environment", selector)
	}

	contract, err := module_rmn_remote.NewRmnRemote(c.CCIPAddress, chain.Client)
	if err != nil {
		return false, fmt.Errorf("failed to create RMN Remote contract: %w", err)
	}

	return contract.DevInspect().IsCursed(
		context.Background(),
		&bind.CallOpts{},
		bind.Object{Id: c.CCIPObjectRef},
		subject[:],
	)
}

func (c *CurseAdapter) IsChainConnectedToTargetChain(e cldf.Environment, selector uint64, targetSelector uint64) (bool, error) {
	chain, ok := e.BlockChains.SuiChains()[selector]
	if !ok {
		return false, fmt.Errorf("SUI chain %d not found in environment", selector)
	}

	routerContract, err := module_router.NewRouter(c.RouterAddress, chain.Client)
	if err != nil {
		return false, fmt.Errorf("failed to create router contract: %w", err)
	}

	connected, err := routerContract.DevInspect().IsChainSupported(
		context.Background(),
		&bind.CallOpts{},
		bind.Object{Id: c.RouterStateObjectID},
		targetSelector,
	)
	if err != nil {
		return false, fmt.Errorf("failed to check if chain %d is connected to chain %d: %w", selector, targetSelector, err)
	}
	return connected, nil
}

func (c *CurseAdapter) IsCurseEnabledForChain(cldf.Environment, uint64) (bool, error) {
	return true, nil
}

func (c *CurseAdapter) SubjectToSelector(subject fastcurse.Subject) (uint64, error) {
	return fastcurse.GenericSubjectToSelector(subject)
}

func (c *CurseAdapter) SelectorToSubject(selector uint64) fastcurse.Subject {
	return globals.FamilyAwareSelectorToSubject(selector, chainsel.FamilySui)
}

func (c *CurseAdapter) Curse() *cldf_ops.Sequence[fastcurse.CurseInput, sequences.OnChainOutput, cldf_chain.BlockChains] {
	return cldf_ops.NewSequence(
		suiseq.SuiCurseSequence.ID(),
		semver.MustParse("1.0.0"),
		suiseq.SuiCurseSequence.Description(),
		func(b cldf_ops.Bundle, chains cldf_chain.BlockChains, in fastcurse.CurseInput) (sequences.OnChainOutput, error) {
			suiInput := suiseq.SuiCurseUncurseInput{
				CCIPAddress:          c.CCIPAddress,
				CCIPObjectRef:        c.CCIPObjectRef,
				CCIPOwnerCapObjectId: c.CCIPOwnerCapObjectId,
				ChainSelector:        in.ChainSelector,
				Subjects:             in.Subjects,
			}
			seqReport, err := cldf_ops.ExecuteSequence(b, suiseq.SuiCurseSequence, chains, suiInput)
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("failed to execute curse sequence on SUI chain %d: %w", in.ChainSelector, err)
			}
			return seqReport.Output, nil
		},
	)
}

func (c *CurseAdapter) Uncurse() *cldf_ops.Sequence[fastcurse.CurseInput, sequences.OnChainOutput, cldf_chain.BlockChains] {
	return cldf_ops.NewSequence(
		suiseq.SuiUncurseSequence.ID(),
		semver.MustParse("1.0.0"),
		suiseq.SuiUncurseSequence.Description(),
		func(b cldf_ops.Bundle, chains cldf_chain.BlockChains, in fastcurse.CurseInput) (sequences.OnChainOutput, error) {
			suiInput := suiseq.SuiCurseUncurseInput{
				CCIPAddress:          c.CCIPAddress,
				CCIPObjectRef:        c.CCIPObjectRef,
				CCIPOwnerCapObjectId: c.CCIPOwnerCapObjectId,
				ChainSelector:        in.ChainSelector,
				Subjects:             in.Subjects,
			}
			seqReport, err := cldf_ops.ExecuteSequence(b, suiseq.SuiUncurseSequence, chains, suiInput)
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("failed to execute uncurse sequence on SUI chain %d: %w", in.ChainSelector, err)
			}
			return seqReport.Output, nil
		},
	)
}

func (c *CurseAdapter) ListConnectedChains(e cldf.Environment, selector uint64) ([]uint64, error) {
	chain, ok := e.BlockChains.SuiChains()[selector]
	if !ok {
		return nil, fmt.Errorf("SUI chain %d not found in environment", selector)
	}

	routerContract, err := module_router.NewRouter(c.RouterAddress, chain.Client)
	if err != nil {
		return nil, fmt.Errorf("failed to create router contract: %w", err)
	}

	connectedChains, err := routerContract.DevInspect().GetDestChains(
		context.Background(),
		&bind.CallOpts{},
		bind.Object{Id: c.RouterStateObjectID},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get connected chains for chain %d: %w", selector, err)
	}
	return connectedChains, nil
}

func (c *CurseAdapter) DeriveCurseAdapterVersion(e cldf.Environment, selector uint64) (*semver.Version, error) {
	return semver.MustParse("1.6.0"), nil
}

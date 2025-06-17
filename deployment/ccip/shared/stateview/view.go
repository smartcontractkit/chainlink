package stateview

import (
	"encoding/json"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	ccipview "github.com/smartcontractkit/chainlink/deployment/ccip/view"
	"github.com/smartcontractkit/chainlink/deployment/common/view"
)

var _ deployment.ViewState = ViewCCIP

func ViewCCIP(e deployment.Environment) (json.Marshaler, error) {
	state, err := LoadOnchainState(e)
	if err != nil {
		return nil, err
	}
	var allChains []uint64
	allChains = append(allChains, e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainselectors.FamilyEVM))...)
	allChains = append(allChains, e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainselectors.FamilySolana))...)
	allChains = append(allChains, e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainselectors.FamilyAptos))...)
	stateView, err := state.View(&e, allChains)
	if err != nil {
		return nil, err
	}
	nopsView, err := view.GenerateNopsView(e.Logger, e.NodeIDs, e.Offchain)
	if err != nil {
		return nil, err
	}
	return ccipview.CCIPView{
		Chains:      stateView.Chains,
		SolChains:   stateView.SolChains,
		AptosChains: stateView.AptosChains,
		Nops:        nopsView,
	}, nil
}

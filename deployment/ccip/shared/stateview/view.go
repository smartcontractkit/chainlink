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
	solSelectors := e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainselectors.FamilySolana))
	allChains := append(e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainselectors.FamilyEVM)), solSelectors...)
	chainView, solanaView, err := state.View(&e, allChains)
	if err != nil {
		return nil, err
	}
	nopsView, err := view.GenerateNopsView(e.Logger, e.NodeIDs, e.Offchain)
	if err != nil {
		return nil, err
	}
	return ccipview.CCIPView{
		Chains:    chainView,
		SolChains: solanaView,
		Nops:      nopsView,
	}, nil
}

package changeset

import (
	"encoding/json"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	commonview "github.com/smartcontractkit/chainlink/deployment/common/view"
	dfView "github.com/smartcontractkit/chainlink/deployment/data-feeds/view"
)

var _ deployment.ViewState = ViewDataFeeds

func ViewDataFeeds(e deployment.Environment) (json.Marshaler, error) {
	state, err := LoadOnchainState(e)
	if err != nil {
		return nil, err
	}
	chainView, err := state.View(e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainselectors.FamilyEVM)), e)
	if err != nil {
		return nil, err
	}
	nopsView, err := commonview.GenerateNopsView(e.Logger, e.NodeIDs, e.Offchain)
	if err != nil {
		return nil, err
	}
	return dfView.DataFeedsView{
		Chains: chainView,
		Nops:   nopsView,
	}, nil
}

package changeset

import (
	"encoding/json"

	"github.com/smartcontractkit/chainlink/deployment"
	commonview "github.com/smartcontractkit/chainlink/deployment/common/view"
	dfView "github.com/smartcontractkit/chainlink/deployment/data-feeds/view"
)

var _ deployment.ViewState = ViewDataFeeds

func ViewDataFeeds(e deployment.Environment) (json.Marshaler, error) {
	state, err := LoadOnchainState(e)
	if err != nil {
		return nil, err
	}
	chainView, err := state.View(e.AllChainSelectors(), e)
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

package changeset

import (
	"encoding/json"
	"fmt"

	"github.com/smartcontractkit/chainlink/deployment"
	dfView "github.com/smartcontractkit/chainlink/deployment/data-feeds/view"
)

var _ deployment.ViewState = ViewDataFeeds

func ViewDataFeeds(e deployment.Environment) (json.Marshaler, error) {
	state, err := LoadOnchainState(e)
	fmt.Println(state)
	if err != nil {
		return nil, err
	}
	chainView, err := state.View(e.AllChainSelectors())
	if err != nil {
		return nil, err
	}
	return dfView.DataFeedsView{
		Chains: chainView,
	}, nil
}

package changeset

import (
	"encoding/json"

	"github.com/smartcontractkit/chainlink/deployment"
	dsView "github.com/smartcontractkit/chainlink/deployment/data-streams/view"
)

var _ deployment.ViewState = ViewDataStreams

func ViewDataStreams(e deployment.Environment) (json.Marshaler, error) {
	state, err := LoadOnchainState(e)
	if err != nil {
		return nil, err
	}
	chainView, err := state.View(e.AllChainSelectors())
	if err != nil {
		return nil, err
	}
	return dsView.DataStreamsView{
		Chains: chainView,
	}, nil
}

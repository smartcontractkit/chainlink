package changeset

import (
	"encoding/json"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	dsstate "github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/state"
	dsView "github.com/smartcontractkit/chainlink/deployment/data-streams/view"
)

var _ deployment.ViewState = ViewDataStreams

func ViewDataStreams(e deployment.Environment) (json.Marshaler, error) {
	return ViewDataStreamsChain(e, e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainselectors.FamilyEVM)))
}

func ViewDataStreamsChain(e deployment.Environment, chainselectors []uint64) (json.Marshaler, error) {
	state, err := dsstate.LoadOnchainState(e)
	if err != nil {
		return nil, err
	}
	chainView, err := state.View(e.GetContext(), chainselectors)
	if err != nil {
		return nil, err
	}
	return dsView.DataStreamsView{
		Chains: chainView,
	}, nil
}

package changeset

import (
	"encoding/json"

	"github.com/smartcontractkit/chainlink/deployment"

	ccipview "github.com/smartcontractkit/chainlink/deployment/ccip/view"
	"github.com/smartcontractkit/chainlink/deployment/common/view"
)

var _ deployment.ViewState = ViewCCIP

func ViewCCIP(e deployment.Environment) (json.Marshaler, error) {
	state, err := LoadOnchainState(e)
	if err != nil {
		return nil, err
	}
	allChains := append(e.AllChainSelectors(), e.AllChainSelectorsSolana()...)
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

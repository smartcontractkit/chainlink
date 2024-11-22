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
	chainView, err := state.View(e.AllChainSelectors())
	if err != nil {
		return nil, err
	}
	nopsView, err := view.GenerateNopsView(e.NodeIDs, e.Offchain)
	if err != nil {
		return nil, err
	}
	return ccipview.CCIPView{
		Chains: chainView,
		Nops:   nopsView,
	}, nil
}

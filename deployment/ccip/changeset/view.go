package changeset

import (
	"encoding/json"

	"github.com/smartcontractkit/chainlink/deployment"
	ccipdeployment "github.com/smartcontractkit/chainlink/deployment/ccip"
	ccipview "github.com/smartcontractkit/chainlink/deployment/ccip/view"
	"github.com/smartcontractkit/chainlink/deployment/common/view"
)

var _ deployment.ViewState = ViewCCIP

func ViewCCIP(e deployment.Environment, ab deployment.AddressBook) (json.Marshaler, error) {
	state, err := ccipdeployment.LoadOnchainState(e, ab)
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

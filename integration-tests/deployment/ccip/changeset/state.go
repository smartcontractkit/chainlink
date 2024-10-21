package changeset

import (
	"encoding/json"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	ccipdeployment "github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip/view"
)

var _ deployment.ViewState = ViewCCIP

func ViewCCIP(e deployment.Environment, ab deployment.AddressBook) (string, error) {
	state, err := ccipdeployment.LoadOnchainState(e, ab)
	if err != nil {
		return "", err
	}
	ccipView, err := state.View(e.AllChainSelectors())
	if err != nil {
		return "", err
	}
	ccipView.NodeOperators, err = view.GenerateNopsView(e.NodeIDs, e.Offchain)
	if err != nil {
		return "", err
	}
	b, err := json.Marshal(ccipView)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

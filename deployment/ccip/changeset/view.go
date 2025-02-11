package changeset

import (
	"encoding/json"

	"github.com/smartcontractkit/chainlink/deployment"
	ccipview "github.com/smartcontractkit/chainlink/deployment/ccip/view"
	"github.com/smartcontractkit/chainlink/deployment/common/view"
)

var _ deployment.ViewState = ViewCCIP
var _ deployment.StateRenderer[ccipview.CCIPView] = RenderCCIPState

// ViewCCIP is a legacy renderer
//
// deprecated: use RenderCCIPState instead.
func ViewCCIP(e deployment.Environment) (json.Marshaler, error) {
	return RenderCCIPState(e)
}

// RenderCCIPState returns a
func RenderCCIPState(e deployment.Environment) (*ccipview.CCIPView, error) {
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
	return &ccipview.CCIPView{
		Chains: chainView,
		Nops:   nopsView,
	}, nil
}

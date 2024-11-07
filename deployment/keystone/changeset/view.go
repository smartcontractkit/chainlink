package changeset

import (
	"encoding/json"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink/deployment"
	commonview "github.com/smartcontractkit/chainlink/deployment/common/view"
	"github.com/smartcontractkit/chainlink/deployment/keystone"
	"github.com/smartcontractkit/chainlink/deployment/keystone/view"
)

var _ deployment.ViewState = ViewKeystone

func ViewKeystone(e deployment.Environment) (json.Marshaler, error) {
	state, err := keystone.GetContractSets(e.Logger, &keystone.GetContractSetsRequest{
		Chains:      e.Chains,
		AddressBook: e.ExistingAddresses,
	})
	if err != nil {
		return nil, err
	}
	chainViews := make(map[string]view.KeystoneChainView)
	for chainSel, contracts := range state.ContractSets {
		chainid, err := chainsel.ChainIdFromSelector(chainSel)
		if err != nil {
			return nil, err
		}
		chainName, err := chainsel.NameFromChainId(chainid)
		if err != nil {
			return nil, err
		}
		v, err := contracts.View()
		if err != nil {
			return nil, err
		}
		chainViews[chainName] = v

	}
	nopsView, err := commonview.GenerateNopsView(e.NodeIDs, e.Offchain)
	if err != nil {
		return nil, err
	}
	return &view.KeystoneView{
		Chains: chainViews,
		Nops:   nopsView,
	}, nil
}

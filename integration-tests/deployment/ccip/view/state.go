package view

import (
	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip"
)

type CCIPSnapShot struct {
	Chains map[string]Chain `json:"chains"`
}

func NewCCIPSnapShot() CCIPSnapShot {
	return CCIPSnapShot{
		Chains: make(map[string]Chain),
	}
}

func SnapshotState(e deployment.Environment, ab deployment.AddressBook) (CCIPSnapShot, error) {
	state, err := ccipdeployment.LoadOnchainState(e, ab)
	if err != nil {
		return CCIPSnapShot{}, err
	}
	return state.Snapshot(e.AllChainSelectors())
}

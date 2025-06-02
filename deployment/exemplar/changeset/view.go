package changeset

import (
	"encoding/json"
	"fmt"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
)

// ExemplarView represents a simplified view of the exemplar environment
type ExemplarView struct {
	Environment string              `json:"environment"`
	Chains      []string            `json:"chains,omitempty"`
	Addresses   map[string][]string `json:"addresses,omitempty"`
}

// MarshalJSON implements json.Marshaler
func (v *ExemplarView) MarshalJSON() ([]byte, error) {
	return json.Marshal(*v)
}

// ViewExemplar extracts basic information from the environment
var ViewExemplar deployment.ViewState = func(e deployment.Environment) (json.Marshaler, error) {
	lggr := e.Logger
	lggr.Info("Generating exemplar state view")

	// Create the view structure
	view := &ExemplarView{
		Environment: e.Name,
		Chains:      make([]string, 0),
		Addresses:   make(map[string][]string),
	}

	// Get chain information
	for _, chainSel := range e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainselectors.FamilyEVM)) {
		chain := e.BlockChains.EVMChains()[chainSel]
		chainName := fmt.Sprintf("%s (%d)", chain.Name(), chainSel)
		view.Chains = append(view.Chains, chainName)

		// Get addresses for this chain
		addresses, err := e.ExistingAddresses.AddressesForChain(chainSel)
		if err != nil {
			lggr.Warnf("Failed to get addresses for chain %s: %v", chainName, err)
			continue
		}

		// Add addresses to the view
		chainAddrs := make([]string, 0)
		for addr, typeVer := range addresses {
			addrInfo := fmt.Sprintf("%s: %s", typeVer.Type, addr)
			chainAddrs = append(chainAddrs, addrInfo)
		}

		view.Addresses[chainName] = chainAddrs
	}

	return view, nil
}

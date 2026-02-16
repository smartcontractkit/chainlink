package stateview

import (
	"encoding/json"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	ccipview "github.com/smartcontractkit/chainlink/deployment/ccip/view"
	"github.com/smartcontractkit/chainlink/deployment/common/view"
)

var _ deployment.ViewState = ViewCCIP

func ViewCCIP(e deployment.Environment) (json.Marshaler, error) {
	state, err := LoadOnchainState(e)
	if err != nil {
		return nil, err
	}
	var allChains []uint64
	allChains = append(allChains, e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainselectors.FamilyEVM))...)
	allChains = append(allChains, e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainselectors.FamilySolana))...)
	allChains = append(allChains, e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainselectors.FamilyAptos))...)
	allChains = append(allChains, e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainselectors.FamilyTon))...)
	allChains = append(allChains, e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainselectors.FamilySui))...)
	stateView, err := state.View(&e, allChains)
	if err != nil {
		return nil, err
	}
	e.Logger.Infow("State view generated, generating NOPs view")
	nopsView, err := view.GenerateNopsView(e.Logger, e.NodeIDs, e.Offchain)
	if err != nil {
		return nil, err
	}
	return ccipview.CCIPView{
		Chains:      stateView.Chains,
		SolChains:   stateView.SolChains,
		AptosChains: stateView.AptosChains,
		TonChains:   stateView.TONChains,
		SuiChains:   stateView.SuiChains,
		Nops:        nopsView,
		V16Nops:     filterNopsForV16(stateView, nopsView),
	}, nil
}

// filter nopsView to only NOPs active in V1.6 DONs with their V1.6 chain keys - this is needed for NOPs billing & NOPs JIRA
func filterNopsForV16(stateView CCIPStateView, nopsView map[string]view.NopView) map[string]view.NopView {
	activeP2PIDs, activeChainNames := make(map[string]struct{}), make(map[string]struct{})

	for _, chainView := range stateView.Chains {
		for _, ccipHome := range chainView.CCIPHome {
			for _, don := range ccipHome.Dons {
				for _, sel := range []uint64{
					don.CommitConfigs.ActiveConfig.Config.ChainSelector,
					don.ExecConfigs.ActiveConfig.Config.ChainSelector,
				} {
					if sel == 0 {
						continue
					}
					if name, err := chainselectors.GetChainNameFromSelector(sel); err == nil {
						activeChainNames[name] = struct{}{}
					}
				}
				for _, node := range don.CommitConfigs.ActiveConfig.Config.Nodes {
					activeP2PIDs[node.P2pID] = struct{}{}
				}
				for _, node := range don.ExecConfigs.ActiveConfig.Config.Nodes {
					activeP2PIDs[node.P2pID] = struct{}{}
				}
			}
		}
	}

	filtered := make(map[string]view.NopView)
	for name, nop := range nopsView {
		if _, ok := activeP2PIDs[nop.PeerID]; !ok {
			continue
		}
		filteredNop := nop
		filteredNop.OCRKeys = make(map[string]view.OCRKeyView)
		for chainName, key := range nop.OCRKeys {
			if _, ok := activeChainNames[chainName]; ok {
				filteredNop.OCRKeys[chainName] = key
			}
		}
		filtered[name] = filteredNop
	}

	return filtered
}

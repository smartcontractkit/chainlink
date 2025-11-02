package state

import (
	"encoding/json"
	"errors"
	"fmt"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	commonview "github.com/smartcontractkit/chainlink/deployment/common/view"
)

var _ deployment.ViewStateV2 = ViewCRE

func ViewCRE(e deployment.Environment, previousView json.Marshaler) (json.Marshaler, error) {
	lggr := e.Logger
	contractsMap, err := getContractsPerChain(e)
	// This is an unrecoverable error
	if err != nil {
		return nil, fmt.Errorf("failed to get contract sets: %w", err)
	}

	prevViewBytes, err := previousView.MarshalJSON()
	if err != nil {
		// just log the error, we don't need to stop the execution since the previous view is optional
		lggr.Warnf("failed to marshal previous keystone view: %v", err)
	}
	var prevView CREView
	if len(prevViewBytes) == 0 {
		prevView.Chains = make(map[string]CREChainView)
	} else if err = json.Unmarshal(prevViewBytes, &prevView); err != nil {
		lggr.Warnf("failed to unmarshal previous keystone view: %v", err)
		prevView.Chains = make(map[string]CREChainView)
	}

	var viewErrs error
	chainViews := make(map[string]CREChainView)
	for chainSel, contracts := range contractsMap {
		chainName, err := chainsel.GetChainNameFromSelector(chainSel)
		if err != nil {
			err2 := fmt.Errorf("failed to resolve chain name for chain selector %d: %w", chainSel, err)
			lggr.Error(err2)
			viewErrs = errors.Join(viewErrs, err2)
			continue
		}
		v, err := GenerateCREChainView(e.GetContext(), e.Logger, prevView.Chains[chainName], contracts)
		if err != nil {
			err2 := fmt.Errorf("failed to view chain %s: %w", chainName, err)
			lggr.Error(err2)
			viewErrs = errors.Join(viewErrs, err2)
			// don't continue; add the partial view
		}
		chainViews[chainName] = v
	}
	nopsView, err := commonview.GenerateNopsView(e.Logger, e.NodeIDs, e.Offchain)
	if err != nil {
		err2 := fmt.Errorf("failed to view nops: %w", err)
		lggr.Error(err2)
		viewErrs = errors.Join(viewErrs, err2)
	}
	return &CREView{
		Chains: chainViews,
		Nops:   nopsView,
	}, viewErrs
}

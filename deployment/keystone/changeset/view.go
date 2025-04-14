package changeset

import (
	"encoding/json"
	"errors"
	"fmt"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink/deployment"
	commonview "github.com/smartcontractkit/chainlink/deployment/common/view"
)

var _ deployment.ViewStateV2 = ViewKeystone

func ViewKeystone(e deployment.Environment, previousView json.Marshaler) (json.Marshaler, error) {
	lggr := e.Logger
	state, err := GetContractSets(e.Logger, &GetContractSetsRequest{
		Chains:      e.Chains,
		AddressBook: e.ExistingAddresses,
	})
	// this error is unrecoverable
	if err != nil {
		return nil, fmt.Errorf("failed to get contract sets: %w", err)
	}

	prevViewBytes, err := previousView.MarshalJSON()
	if err != nil {
		// just log the error, we don't need to stop the execution since the previous view is optional
		lggr.Warnf("failed to marshal previous keystone view: %v", err)
	}
	var prevView KeystoneView
	if len(prevViewBytes) == 0 {
		prevView.Chains = make(map[string]KeystoneChainView)
	} else if err = json.Unmarshal(prevViewBytes, &prevView); err != nil {
		lggr.Warnf("failed to unmarshal previous keystone view: %v", err)
		prevView.Chains = make(map[string]KeystoneChainView)
	}

	var viewErrs error
	chainViews := make(map[string]KeystoneChainView)
	for chainSel, contracts := range state.ContractSets {
		chainid, err := chainsel.ChainIdFromSelector(chainSel)
		if err != nil {
			err2 := fmt.Errorf("failed to resolve chain id for selector %d: %w", chainSel, err)
			lggr.Error(err2)
			viewErrs = errors.Join(viewErrs, err2)
			continue
		}
		chainName, err := chainsel.NameFromChainId(chainid)
		if err != nil {
			err2 := fmt.Errorf("failed to resolve chain name for chain id %d: %w", chainid, err)
			lggr.Error(err2)
			viewErrs = errors.Join(viewErrs, err2)
			continue
		}
		v, err := contracts.View(e.GetContext(), prevView.Chains[chainName], e.Logger)
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
	return &KeystoneView{
		Chains: chainViews,
		Nops:   nopsView,
	}, viewErrs
}

package changeset

import (
	"context"
	"fmt"
	"maps"
	"slices"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	workflow_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper_v2"
)

var _ cldf.ChangeSet[DeployWorkflowRegistryV2Request] = DeployWorkflowRegistryV2

type DeployWorkflowRegistryV2Request struct {
	ChainSelectors []uint64 // filter to only deploy to these chains; if empty, deploy to all chains
}

// DeployWorkflowRegistryV2 deploys the WorkflowRegistry v2 contract to specified chains in the environment
func DeployWorkflowRegistryV2(env cldf.Environment, cfg DeployWorkflowRegistryV2Request) (cldf.ChangesetOutput, error) {
	var out cldf.ChangesetOutput
	out.DataStore = datastore.NewMemoryDataStore()

	selectors := cfg.ChainSelectors
	evmChains := env.BlockChains.EVMChains()

	if len(selectors) == 0 {
		selectors = slices.Collect(maps.Keys(evmChains))
	}

	for _, sel := range selectors {
		chain, ok := evmChains[sel]
		if !ok {
			return cldf.ChangesetOutput{}, fmt.Errorf("chain with selector %d not found", sel)
		}

		env.Logger.Infow("deploying workflow registry v2", "chainSelector", chain.Selector)

		resp, err := DeployWorkflowRegistryV2ToChain(env.GetContext(), chain, out.DataStore)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy WorkflowRegistry v2 to chain selector %d: %w", chain.Selector, err)
		}

		env.Logger.Infof("Deployed %s chain selector %d addr %s", resp.Tv.String(), chain.Selector, resp.Address.String())

		// Add to datastore
		r := datastore.AddressRef{
			ChainSelector: sel,
			Address:       resp.Address.String(),
			Type:          datastore.ContractType(resp.Tv.Type),
			Version:       &resp.Tv.Version,
			Qualifier:     "",
			Labels:        datastore.NewLabelSet(),
		}

		// Add labels from the response
		for _, l := range resp.Tv.Labels.List() {
			r.Labels.Add(l)
		}

		if err := out.DataStore.Addresses().Add(r); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to save address ref in datastore for chain selector %d: %w", sel, err)
		}
	}

	return out, nil
}

// WorkflowRegistryDeployResponse represents the response from deploying a workflow registry contract
type WorkflowRegistryDeployResponse struct {
	Address common.Address
	Tx      common.Hash
	Tv      cldf.TypeAndVersion
}

// DeployWorkflowRegistryV2ToChain deploys WorkflowRegistry v2 to a single chain
func DeployWorkflowRegistryV2ToChain(ctx context.Context, chain cldf_evm.Chain, ds datastore.MutableDataStore) (*WorkflowRegistryDeployResponse, error) {
	// Deploy the contract directly
	addr, tx, wr, err := workflow_registry_v2.DeployWorkflowRegistry(
		chain.DeployerKey,
		chain.Client)
	if err != nil {
		return nil, fmt.Errorf("failed to deploy WorkflowRegistry v2: %w", err)
	}

	_, err = chain.Confirm(tx)
	if err != nil {
		return nil, fmt.Errorf("failed to confirm WorkflowRegistry v2 deployment: %w", err)
	}

	tvStr, err := wr.TypeAndVersion(&bind.CallOpts{})
	if err != nil {
		return nil, fmt.Errorf("failed to get type and version: %w", err)
	}

	tv, err := cldf.TypeAndVersionFromString(tvStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse type and version from %s: %w", tvStr, err)
	}

	resp := &WorkflowRegistryDeployResponse{
		Address: addr,
		Tx:      tx.Hash(),
		Tv:      tv,
	}

	// Save to datastore
	r := datastore.AddressRef{
		ChainSelector: chain.Selector,
		Address:       resp.Address.String(),
		Type:          datastore.ContractType(resp.Tv.Type),
		Version:       &resp.Tv.Version,
		Qualifier:     "",
		Labels:        datastore.NewLabelSet(),
	}

	// Add labels from the response
	for _, l := range resp.Tv.Labels.List() {
		r.Labels.Add(l)
	}

	if err := ds.Addresses().Add(r); err != nil {
		return nil, fmt.Errorf("failed to save WorkflowRegistry v2 to datastore: %w", err)
	}

	return resp, nil
}

// DeployWorkflowRegistryV2Single deploys WorkflowRegistry v2 to a single chain (simple interface)
func DeployWorkflowRegistryV2Single(env cldf.Environment, chainSelector uint64) (cldf.ChangesetOutput, error) {
	return DeployWorkflowRegistryV2(env, DeployWorkflowRegistryV2Request{
		ChainSelectors: []uint64{chainSelector},
	})
}

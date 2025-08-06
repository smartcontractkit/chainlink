package workflowregistry

import (
	"fmt"
	"maps"
	"slices"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
)

var _ cldf.ChangeSet[*changeset.DeployRequestV2] = DeployWorkflowRegistryV2Single

// DeployWorkflowRegistryV2Single deploys WorkflowRegistry v2 using the DeployRequestV2 pattern
func DeployWorkflowRegistryV2Single(env cldf.Environment, req *changeset.DeployRequestV2) (cldf.ChangesetOutput, error) {
	lggr := env.Logger
	chain, ok := env.BlockChains.EVMChains()[req.ChainSel]
	if !ok {
		return cldf.ChangesetOutput{}, fmt.Errorf("chain not found in environment for selector %d", req.ChainSel)
	}

	ds := datastore.NewMemoryDataStore()
	resp, err := changeset.DeployWorkflowRegistryV2ToChain(env.GetContext(), chain, ds)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy WorkflowRegistry v2: %w", err)
	}
	lggr.Infof("Deployed %s chain selector %d addr %s", resp.Tv.String(), chain.Selector, resp.Address.String())

	// Add additional fields from request
	addresses, err := ds.Addresses().Fetch()
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to fetch addresses from datastore: %w", err)
	}

	// Update the address ref with request-specific qualifier and labels
	for i, addr := range addresses {
		if addr.ChainSelector == req.ChainSel && addr.Address == resp.Address.String() {
			addresses[i].Qualifier = req.Qualifier
			if req.Labels != nil {
				addresses[i].Labels = *req.Labels
			}
			// Add labels from the response
			for _, l := range resp.Tv.Labels.List() {
				addresses[i].Labels.Add(l)
			}
			break
		}
	}

	return cldf.ChangesetOutput{DataStore: ds}, nil
}

// DeployV2MultiChain deploys WorkflowRegistry v2 to multiple chains
func DeployV2MultiChain(env cldf.Environment, req DeployV2MultiChainRequest) (cldf.ChangesetOutput, error) {
	var out cldf.ChangesetOutput
	out.DataStore = datastore.NewMemoryDataStore()

	selectors := req.ChainSelectors
	evmChains := env.BlockChains.EVMChains()

	if len(selectors) == 0 {
		selectors = slices.Collect(maps.Keys(evmChains))
	}

	for _, sel := range selectors {
		// Deploy to each chain using the single chain method
		deployReq := &changeset.DeployRequestV2{
			ChainSel:  sel,
			Qualifier: req.Qualifier,
			Labels:    req.Labels,
		}

		csOut, err := DeployWorkflowRegistryV2Single(env, deployReq)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy WorkflowRegistry v2 to chain selector %d: %w", sel, err)
		}

		if err := out.DataStore.Merge(csOut.DataStore.Seal()); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to merge datastore for chain selector %d: %w", sel, err)
		}
	}

	return out, nil
}

type DeployV2MultiChainRequest struct {
	ChainSelectors []uint64 // filter to only deploy to these chains; if empty, deploy to all chains
	Qualifier      string
	Labels         *datastore.LabelSet
}

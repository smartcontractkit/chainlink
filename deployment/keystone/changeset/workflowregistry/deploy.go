package workflowregistry

import (
	"errors"
	"fmt"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
)

var _ deployment.ChangeSet[uint64] = Deploy

func Deploy(env deployment.Environment, registrySelector uint64) (deployment.ChangesetOutput, error) {
	lggr := env.Logger
	chain, ok := env.Chains[registrySelector]
	if !ok {
		return deployment.ChangesetOutput{}, errors.New("chain not found in environment")
	}
	ab := deployment.NewMemoryAddressBook()
	wrResp, err := deployWorkflowRegistry(chain, ab)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to deploy WorkflowRegistry: %w", err)
	}
	lggr.Infof("Deployed %s chain selector %d addr %s", wrResp.Tv.String(), chain.Selector, wrResp.Address.String())

	return deployment.ChangesetOutput{AddressBook: ab}, nil
}

func DeployV2(env deployment.Environment, req *changeset.DeployRequestV2) (deployment.ChangesetOutput, error) {
	lggr := env.Logger
	chain, ok := env.Chains[req.ChainSel]
	if !ok {
		return deployment.ChangesetOutput{}, errors.New("chain not found in environment")
	}
	ab := deployment.NewMemoryAddressBook()
	wrResp, err := deployWorkflowRegistry(chain, ab)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to deploy WorkflowRegistry: %w", err)
	}
	lggr.Infof("Deployed %s chain selector %d addr %s", wrResp.Tv.String(), chain.Selector, wrResp.Address.String())

	ds := datastore.NewMemoryDataStore[
		datastore.DefaultMetadata,
		datastore.DefaultMetadata,
	]()
	r := datastore.AddressRef{
		ChainSelector: req.ChainSel,
		Address:       wrResp.Address.String(),
		Type:          datastore.ContractType(wrResp.Tv.Type),
		Version:       &wrResp.Tv.Version,
		Qualifier:     req.Qualifier,
	}
	if req.Labels != nil {
		r.Labels = *req.Labels
	}

	if err = ds.Addresses().Add(r); err != nil {
		return deployment.ChangesetOutput{DataStore: ds},
			fmt.Errorf("failed to save address ref in datastore: %w", err)
	}
	return deployment.ChangesetOutput{AddressBook: ab, DataStore: ds}, nil
}

package workflowregistry

import (
	"errors"
	"fmt"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
)

var _ cldf.ChangeSet[uint64] = Deploy

func Deploy(env cldf.Environment, registrySelector uint64) (cldf.ChangesetOutput, error) {
	lggr := env.Logger
	chain, ok := env.BlockChains.EVMChains()[registrySelector]
	if !ok {
		return cldf.ChangesetOutput{}, errors.New("chain not found in environment")
	}
	ab := cldf.NewMemoryAddressBook()
	wrResp, err := deployWorkflowRegistry(chain, ab)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy WorkflowRegistry: %w", err)
	}
	lggr.Infof("Deployed %s chain selector %d addr %s", wrResp.Tv.String(), chain.Selector, wrResp.Address.String())

	return cldf.ChangesetOutput{AddressBook: ab}, nil
}

func DeployV2(env cldf.Environment, req *changeset.DeployRequestV2) (cldf.ChangesetOutput, error) {
	lggr := env.Logger
	chain, ok := env.BlockChains.EVMChains()[req.ChainSel]
	if !ok {
		return cldf.ChangesetOutput{}, errors.New("chain not found in environment")
	}
	ab := cldf.NewMemoryAddressBook()
	wrResp, err := deployWorkflowRegistry(chain, ab)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy WorkflowRegistry: %w", err)
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
		return cldf.ChangesetOutput{DataStore: ds},
			fmt.Errorf("failed to save address ref in datastore: %w", err)
	}
	return cldf.ChangesetOutput{AddressBook: ab, DataStore: ds}, nil
}

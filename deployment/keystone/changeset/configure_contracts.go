package changeset

import (
	"fmt"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
	kf "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/forwarder_1_0_0"
)

var _ deployment.ChangeSet[InitialContractsCfg] = ConfigureInitialContractsChangeset

type InitialContractsCfg struct {
	RegistryChainSel uint64
	Dons             []internal.DonCapabilities
	OCR3Config       *internal.OracleConfig
}

func ConfigureInitialContractsChangeset(e deployment.Environment, cfg InitialContractsCfg) (deployment.ChangesetOutput, error) {
	req := &internal.ConfigureContractsRequest{
		Env:              &e,
		RegistryChainSel: cfg.RegistryChainSel,
		Dons:             cfg.Dons,
		OCR3Config:       cfg.OCR3Config,
	}
	// load all the forwarders
	// forwarder on all chains
	//foundForwarder = false
	fowarders := make(map[uint64]*kf.KeystoneForwarder)
	contractSetsResp, err := internal.GetContractSets(req.Env.Logger, &internal.GetContractSetsRequest{
		AddressBook: req.Env.ExistingAddresses,
		Chains:      req.Env.Chains,
	})
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get contract sets: %w", err)
	}
	for sel, cs := range contractSetsResp.ContractSets {
		fowarders[sel] = cs.Forwarder
	}
	req.ForwarderContracts = fowarders

	resp, err := internal.ConfigureContracts(req.Env.GetContext(), *req)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to configure contracts: %w", err)
	}
	return *resp.Changeset, nil
}

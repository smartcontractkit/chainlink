package changeset

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
)

var _ cldf.ChangeSet[InitialContractsCfg] = ConfigureInitialContractsChangeset

type InitialContractsCfg struct {
	RegistryChainSel uint64
	Dons             []internal.DonCapabilities
	OCR3Config       *ocr3.OracleConfig
}

func ConfigureInitialContractsChangeset(e cldf.Environment, cfg InitialContractsCfg) (cldf.ChangesetOutput, error) {
	req := &internal.ConfigureContractsRequest{
		Env:              &e,
		RegistryChainSel: cfg.RegistryChainSel,
		Dons:             cfg.Dons,
		OCR3Config:       cfg.OCR3Config,
	}
	return ConfigureInitialContracts(e.Logger, req)
}

// Deprecated: Use ConfigureInitialContractsChangeset instead.
func ConfigureInitialContracts(lggr logger.Logger, req *internal.ConfigureContractsRequest) (cldf.ChangesetOutput, error) {
	if err := req.Validate(); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to validate request: %w", err)
	}

	regAddrs, err := req.Env.ExistingAddresses.AddressesForChain(req.RegistryChainSel)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("no addresses found for chain %d: %w", req.RegistryChainSel, err)
	}
	foundRegistry := false
	foundOCR3 := false
	foundForwarder := false
	for _, addr := range regAddrs {
		switch addr.Type {
		case internal.CapabilitiesRegistry:
			foundRegistry = true
		case internal.OCR3Capability:
			foundOCR3 = true
		case internal.KeystoneForwarder:
			foundForwarder = true
		}
	}
	if !foundRegistry || !foundOCR3 || !foundForwarder {
		return cldf.ChangesetOutput{}, fmt.Errorf("missing contracts on registry chain %d in addressbook for changeset %s registry exists %t, ocr3 exist %t, forwarder exists %t ", req.RegistryChainSel, "0003_deploy_forwarder",
			foundRegistry, foundOCR3, foundForwarder)
	}
	// forwarder on all chains
	foundForwarder = false
	for _, c := range req.Env.BlockChains.EVMChains() {
		addrs, err2 := req.Env.ExistingAddresses.AddressesForChain(c.Selector)
		if err2 != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("no addresses found for chain %d: %w", c.Selector, err2)
		}
		for _, addr := range addrs {
			if addr.Type == internal.KeystoneForwarder {
				foundForwarder = true
				break
			}
		}
		if !foundForwarder {
			return cldf.ChangesetOutput{}, fmt.Errorf("no forwarder found for chain %d", c.Selector)
		}
	}

	resp, err := internal.ConfigureContracts(context.TODO(), lggr, *req)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to configure contracts: %w", err)
	}
	return *resp.Changeset, nil
}

package changeset

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	kslib "github.com/smartcontractkit/chainlink/integration-tests/deployment/keystone"
)

func ConfigureInitialContracts(lggr logger.Logger, req *kslib.ConfigureContractsRequest) (deployment.ChangesetOutput, error) {
	if err := req.Validate(); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to validate request: %w", err)
	}

	regAddrs, err := req.AddressBook.AddressesForChain(req.RegistryChainSel)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("no addresses found for chain %d: %w", req.RegistryChainSel, err)
	}
	foundRegistry := false
	foundOCR3 := false
	foundForwarder := false
	for _, addr := range regAddrs {
		switch addr.Type {
		case kslib.CapabilityRegistryTypeVersion.Type:
			foundRegistry = true
		case kslib.OCR3CapabilityTypeVersion.Type:
			foundOCR3 = true
		case kslib.ForwarderTypeVersion.Type:
			foundForwarder = true
		}
	}
	if !foundRegistry || !foundOCR3 || !foundForwarder {
		return deployment.ChangesetOutput{}, fmt.Errorf("missing contracts on registry chain %d in addressbook for changeset %s registry exists %t, ocr3 exist %t, forwarder exists %t ", req.RegistryChainSel, "0003_deploy_forwarder",
			foundRegistry, foundOCR3, foundForwarder)
	}
	// forwarder on all chains
	foundForwarder = false
	for _, c := range req.Env.Chains {
		addrs, err2 := req.AddressBook.AddressesForChain(c.Selector)
		if err2 != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("no addresses found for chain %d: %w", c.Selector, err2)
		}
		for _, addr := range addrs {
			if addr.Type == kslib.ForwarderTypeVersion.Type {
				foundForwarder = true
				break
			}
		}
		if !foundForwarder {
			return deployment.ChangesetOutput{}, fmt.Errorf("no forwarder found for chain %d", c.Selector)
		}
	}

	resp, err := kslib.ConfigureContracts(context.TODO(), lggr, *req)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to configure contracts: %w", err)
	}
	return *resp.Changeset, nil
}

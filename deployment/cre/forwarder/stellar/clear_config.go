package stellar

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"

	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	crebindings "github.com/smartcontractkit/chainlink-stellar/bindings/contracts/cre"
	stellardeployment "github.com/smartcontractkit/chainlink-stellar/deployment"
)

var _ cldf.ChangeSetV2[*ClearConfigRequest] = ClearConfigs{}

// ClearConfigs removes a DON signer configuration from deployed Stellar CRE
// forwarders. Used to rotate out a bad or retired DON config; reports for the
// cleared (donID, configVersion) pair are rejected afterwards.
type ClearConfigs struct{}

type ClearConfigRequest struct {
	DonID         uint32
	ConfigVersion uint32

	// Chains is optional. When set, only those selectors are cleared.
	Chains    map[uint64]struct{}
	Qualifier string
	Version   string

	// MCMS, when set, builds a governed clear_config timelock proposal instead
	// of sending directly with the deployer key. Required once the forwarder
	// is owned by the MCMS timelock. Its Qualifier selects the MCMS instance.
	MCMS *cldf.MCMSTimelockProposalInput
}

func (cs ClearConfigs) VerifyPreconditions(env cldf.Environment, req *ClearConfigRequest) error {
	if req == nil {
		return errors.New("request is required")
	}
	if req.DonID == 0 {
		return errors.New("DON ID is required")
	}
	if req.ConfigVersion == 0 {
		return errors.New("config version is required")
	}
	if req.Qualifier == "" {
		return errors.New("forwarder qualifier is required")
	}
	if req.Version == "" {
		return errors.New("forwarder version is required")
	}
	version, err := semver.NewVersion(req.Version)
	if err != nil {
		return fmt.Errorf("invalid forwarder version %q: %w", req.Version, err)
	}
	if req.MCMS != nil {
		if err := req.MCMS.Validate(); err != nil {
			return fmt.Errorf("invalid MCMS timelock proposal input: %w", err)
		}
	}

	chains := req.Chains
	if len(chains) == 0 {
		chains = make(map[uint64]struct{})
		for sel := range env.BlockChains.StellarChains() {
			chains[sel] = struct{}{}
		}
	}
	for sel := range chains {
		if _, ok := env.BlockChains.StellarChains()[sel]; !ok {
			return fmt.Errorf("stellar chain not found for chain selector %d", sel)
		}
		refKey := datastore.NewAddressRefKey(sel, ForwarderContract, version, req.Qualifier)
		if _, err := env.DataStore.Addresses().Get(refKey); err != nil {
			return fmt.Errorf("failed to get stellar forwarder for ref key %s: %w", refKey, err)
		}
	}

	return nil
}

func (cs ClearConfigs) Apply(env cldf.Environment, req *ClearConfigRequest) (cldf.ChangesetOutput, error) {
	var out cldf.ChangesetOutput

	version := semver.MustParse(req.Version)
	chains := req.Chains
	if len(chains) == 0 {
		chains = make(map[uint64]struct{})
		for sel := range env.BlockChains.StellarChains() {
			chains[sel] = struct{}{}
		}
	}

	var batchOps []mcmstypes.BatchOperation

	for sel := range chains {
		ch, ok := env.BlockChains.StellarChains()[sel]
		if !ok {
			return out, fmt.Errorf("stellar chain not found for chain selector %d", sel)
		}

		refKey := datastore.NewAddressRefKey(sel, ForwarderContract, version, req.Qualifier)
		addrRef, err := env.DataStore.Addresses().Get(refKey)
		if err != nil {
			return out, fmt.Errorf("failed to get stellar forwarder for ref key %s: %w", refKey, err)
		}

		// The governed path only reads (owner check) and encodes; it must not
		// require the deployer signing key.
		if req.MCMS != nil {
			if err := requireTimelockOwnership(env.GetContext(), env, ch, addrRef.Address, *req.MCMS); err != nil {
				return out, err
			}

			batchOp, err := forwarderClearConfigBatchOp(sel, addrRef.Address, req.DonID, req.ConfigVersion)
			if err != nil {
				return out, fmt.Errorf("failed to build clear_config batch operation for stellar forwarder %s on chain selector %d: %w", addrRef.Address, sel, err)
			}
			batchOps = append(batchOps, batchOp)

			env.Logger.Infow("Built governed Stellar forwarder clear_config operation", "chainSelector", sel, "forwarder", addrRef.Address, "donID", req.DonID, "configVersion", req.ConfigVersion)

			continue
		}

		deployer, err := stellardeployment.NewDeployerFromChain(ch)
		if err != nil {
			return out, fmt.Errorf("failed to build stellar deployer for chain selector %d: %w", sel, err)
		}

		client := crebindings.NewForwarderClient(deployer, addrRef.Address)
		if err := client.ClearConfig(env.GetContext(), req.DonID, req.ConfigVersion); err != nil {
			return out, fmt.Errorf("failed to clear config on stellar forwarder %s (chain selector %d): %w", addrRef.Address, sel, err)
		}

		env.Logger.Infow("Cleared Stellar CRE forwarder config", "chainSelector", sel, "forwarder", addrRef.Address, "donID", req.DonID, "configVersion", req.ConfigVersion)
	}

	if req.MCMS != nil {
		return cldf.NewOutputBuilder(env, nil).
			WithTimelockProposal(*req.MCMS, batchOps).
			Build()
	}

	return out, nil
}

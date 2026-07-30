package stellar

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	stellardeployment "github.com/smartcontractkit/chainlink-stellar/deployment"
	stellarforwarder "github.com/smartcontractkit/chainlink-stellar/deployment/cre/forwarder"
)

var _ cldf.ChangeSetV2[*AddForwardersRequest] = AddForwarders{}

// AddForwarders authorizes transmitter accounts on a deployed Stellar CRE forwarder.
//
// The Soroban forwarder's report() calls require_valid_forwarder(transmitter): only
// accounts registered via add_forwarder may submit reports. The DON worker nodes
// submit report() from their relayer signing accounts, so each must be authorized
// here or report() reverts with UnauthorizedForwarder.
type AddForwarders struct{}

type AddForwardersRequest struct {
	// Transmitters are Stellar account addresses (G… StrKey) authorized to submit a report
	Transmitters []string

	// Chains is optional. When set, only those selectors are configured.
	Chains    map[uint64]struct{}
	Qualifier string
	Version   string
}

func (cs AddForwarders) VerifyPreconditions(env cldf.Environment, req *AddForwardersRequest) error {
	if req == nil {
		return errors.New("request is required")
	}
	if len(req.Transmitters) == 0 {
		return errors.New("at least one transmitter is required")
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

func (cs AddForwarders) Apply(env cldf.Environment, req *AddForwardersRequest) (cldf.ChangesetOutput, error) {
	var out cldf.ChangesetOutput

	version := semver.MustParse(req.Version)
	chains := req.Chains
	if len(chains) == 0 {
		chains = make(map[uint64]struct{})
		for sel := range env.BlockChains.StellarChains() {
			chains[sel] = struct{}{}
		}
	}

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

		deployer, err := stellardeployment.NewDeployerFromChain(ch)
		if err != nil {
			return out, fmt.Errorf("failed to build stellar deployer for chain selector %d: %w", sel, err)
		}

		for _, transmitter := range req.Transmitters {
			if err := stellarforwarder.AddForwarder(env.GetContext(), deployer, addrRef.Address, transmitter); err != nil {
				return out, fmt.Errorf("failed to authorize transmitter %s on stellar forwarder %s (chain selector %d): %w", transmitter, addrRef.Address, sel, err)
			}
			env.Logger.Infow("Authorized Stellar forwarder transmitter", "chainSelector", sel, "forwarder", addrRef.Address, "transmitter", transmitter)
		}
	}

	return out, nil
}

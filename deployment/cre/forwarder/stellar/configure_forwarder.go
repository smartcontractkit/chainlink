package stellar

import (
	"errors"
	"fmt"
	"sort"

	"github.com/Masterminds/semver/v3"
	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/offchain"
	stellardeployment "github.com/smartcontractkit/chainlink-stellar/deployment"
	stellarforwarder "github.com/smartcontractkit/chainlink-stellar/deployment/cre/forwarder"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/cre/forwarder"
)

const ForwarderContract datastore.ContractType = "StellarForwarder"

const (
	DefaultForwarderQualifier = "stellar_forwarder"
	DefaultForwarderVersion   = "1.0.0"
)

var _ cldf.ChangeSetV2[*ConfigureForwarderRequest] = ConfigureForwarders{}

type ConfigureForwarders struct{}

type ConfigureForwarderRequest struct {
	DON forwarder.DonConfiguration

	// Chains is optional. When set, only those selectors are configured.
	Chains    map[uint64]struct{}
	Qualifier string
	Version   string
}

func (cs ConfigureForwarders) VerifyPreconditions(env cldf.Environment, req *ConfigureForwarderRequest) error {
	if req == nil {
		return errors.New("request is required")
	}
	if req.DON.Name == "" {
		return errors.New("DON name is required")
	}
	if len(req.DON.NodeIDs) == 0 {
		return errors.New("DON node IDs are required")
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

func (cs ConfigureForwarders) Apply(env cldf.Environment, req *ConfigureForwarderRequest) (cldf.ChangesetOutput, error) {
	var out cldf.ChangesetOutput

	version := semver.MustParse(req.Version)
	chains := req.Chains
	if len(chains) == 0 {
		chains = make(map[uint64]struct{})
		for sel := range env.BlockChains.StellarChains() {
			chains[sel] = struct{}{}
		}
	}

	signers, err := stellarSigners(req.DON.NodeIDs, env.Offchain)
	if err != nil {
		return out, fmt.Errorf("failed to resolve stellar OCR signers for DON %q: %w", req.DON.Name, err)
	}
	if len(signers) == 0 {
		return out, fmt.Errorf("no stellar signers found for DON %q", req.DON.Name)
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

		if err := stellarforwarder.ConfigureForwarder(
			env.GetContext(),
			deployer,
			addrRef.Address,
			req.DON.ID,
			req.DON.Version,
			uint32(req.DON.F),
			signers,
		); err != nil {
			return out, fmt.Errorf("failed to configure stellar forwarder %s on chain selector %d: %w", addrRef.Address, sel, err)
		}

		env.Logger.Infow("Configured Stellar CRE forwarder", "chainSelector", sel, "forwarder", addrRef.Address, "donID", req.DON.ID, "f", req.DON.F, "signersLen", len(signers))
	}

	return out, nil
}

// stellarSigners returns ed25519 onchain public keys for the given node IDs by
// reading FamilyStellar OCR configs from the job distributor.
func stellarSigners(nodeIDs []string, offchainClient offchain.Client) ([][32]byte, error) {
	if offchainClient == nil {
		return nil, errors.New("offchain client is required to resolve stellar OCR signers")
	}

	nodes, err := deployment.NodeInfo(nodeIDs, offchainClient)
	if err != nil {
		return nil, fmt.Errorf("failed to get node info: %w", err)
	}
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].PeerID.String() < nodes[j].PeerID.String()
	})

	out := make([][32]byte, 0, len(nodes))
	for _, n := range nodes {
		if n.IsBootstrap {
			continue
		}

		var stellarCC *deployment.OCRConfig
		for details, cfg := range n.SelToOCRConfig {
			family, famErr := chainselectors.GetSelectorFamily(details.ChainSelector)
			if famErr == nil && family == chainselectors.FamilyStellar {
				cc := cfg
				stellarCC = &cc
				break
			}
		}
		if stellarCC == nil {
			return nil, fmt.Errorf("no stellar OCR2 config for node %s", n.NodeID)
		}
		if len(stellarCC.OnchainPublicKey) != 32 {
			return nil, fmt.Errorf("expected 32-byte stellar onchain public key for node %s, got %d", n.NodeID, len(stellarCC.OnchainPublicKey))
		}
		var arr [32]byte
		copy(arr[:], stellarCC.OnchainPublicKey[:32])
		out = append(out, arr)
	}

	if len(out) == 0 {
		return nil, errors.New("no stellar signers resolved from node OCR configs")
	}
	return out, nil
}

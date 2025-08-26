package contracts

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"

	"github.com/smartcontractkit/chainlink/deployment"
)

var (
	CapabilitiesRegistry      cldf.ContractType = "CapabilitiesRegistry"      // https://github.com/smartcontractkit/chainlink-evm/blob/f190212bab15e84fe49db88f495ad026e6c1d520/contracts/src/v0.8/workflow/dev/v2/CapabilitiesRegistry.sol#L450
	WorkflowRegistry          cldf.ContractType = "WorkflowRegistry"          // https://github.com/smartcontractkit/chainlink/blob/develop/contracts/src/v0.8/workflow/WorkflowRegistry.sol
	KeystoneForwarder         cldf.ContractType = "KeystoneForwarder"         // https://github.com/smartcontractkit/chainlink/blob/50c1b3dbf31bd145b312739b08967600a5c67f30/contracts/src/v0.8/keystone/KeystoneForwarder.sol#L90
	OCR3Capability            cldf.ContractType = "OCR3Capability"            // https://github.com/smartcontractkit/chainlink/blob/50c1b3dbf31bd145b312739b08967600a5c67f30/contracts/src/v0.8/keystone/OCR3Capability.sol#L12
	FeedConsumer              cldf.ContractType = "FeedConsumer"              // no type and a version in contract https://github.com/smartcontractkit/chainlink/blob/89183a8a5d22b1aeca0ade3b76d16aa84067aa57/contracts/src/v0.8/keystone/KeystoneFeedsConsumer.sol#L1
	RBACTimelock              cldf.ContractType = "RBACTimelock"              // no type and a version in contract https://github.com/smartcontractkit/ccip-owner-contracts/blob/main/src/RBACTimelock.sol
	ProposerManyChainMultiSig cldf.ContractType = "ProposerManyChainMultiSig" // no type and a version in contract https://github.com/smartcontractkit/ccip-owner-contracts/blob/main/src/ManyChainMultiSig.sol
)

type RegisteredDonConfig struct {
	Name             string
	NodeIDs          []string // ids in the offchain client
	RegistryChainSel uint64
	Registry         *capabilities_registry_v2.CapabilitiesRegistry
}

type ConfigureCREDON struct {
	Name    string
	NodeIDs []string
}

// RegisteredDon is a representation of a don that exists in the in the capabilities registry all with the enriched node data
type RegisteredDon struct {
	Name  string
	Info  capabilities_registry_v2.CapabilitiesRegistryDONInfo
	Nodes []deployment.Node
}

func newRegisteredDon(env cldf.Environment, cfg RegisteredDonConfig) (*RegisteredDon, error) {
	if cfg.Registry == nil {
		return nil, errors.New("capabilities registry not found in config")
	}

	var (
		err    error
		capReg = cfg.Registry
	)

	di, err := capReg.GetDONs(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get dons: %w", err)
	}
	// load the nodes from the offchain client
	nodes, err := deployment.NodeInfo(cfg.NodeIDs, env.Offchain)
	if err != nil {
		return nil, fmt.Errorf("failed to get node info: %w", err)
	}
	want := sortedHash(nodes.PeerIDs())
	var don *capabilities_registry_v2.CapabilitiesRegistryDONInfo
	for i, d := range di {
		got := sortedHash(d.NodeP2PIds)
		if got == want {
			don = &di[i]
		}
	}
	if don == nil {
		return nil, errors.New("don not found in registry")
	}
	return &RegisteredDon{
		Name:  cfg.Name,
		Info:  *don,
		Nodes: nodes,
	}, nil
}

func sortedHash(p2pids [][32]byte) string {
	sha256Hash := sha256.New()
	sort.Slice(p2pids, func(i, j int) bool {
		return bytes.Compare(p2pids[i][:], p2pids[j][:]) < 0
	})
	for _, id := range p2pids {
		sha256Hash.Write(id[:])
	}
	return hex.EncodeToString(sha256Hash.Sum(nil))
}

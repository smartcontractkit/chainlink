package changeset

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"

	"github.com/smartcontractkit/chainlink/deployment"
	ccipdeployment "github.com/smartcontractkit/chainlink/deployment/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/rmn_home"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
)

var _ deployment.ChangeSet = DeployCapReg

// DeployCapReg is a separate changeset because cap reg is an env var for CL nodes.
func DeployCapReg(env deployment.Environment, config interface{}) (deployment.ChangesetOutput, error) {
	cfg, ok := config.(DeployHomeChainConfig)
	if !ok {
		return deployment.ChangesetOutput{}, deployment.ErrInvalidConfig
	}
	err := cfg.Validate()
	if err != nil {
		return deployment.ChangesetOutput{}, errors.Wrapf(deployment.ErrInvalidConfig, "%v", err)
	}
	homeChainSel := cfg.HomeChainSel
	// Note we also deploy the cap reg.
	ab := deployment.NewMemoryAddressBook()
	capReg, nopIdsByAdmin, err := ccipdeployment.DeployCapReg(env.Logger, ab, env.Chains[homeChainSel], cfg.RMNStaticConfig, cfg.RMNDynamicConfig, cfg.NodeOperators)
	if err != nil {
		env.Logger.Errorw("Failed to deploy cap reg", "err", err, "addresses", ab)
		return deployment.ChangesetOutput{}, err
	}
	// validate all node operators have nopIds
	for _, nop := range cfg.NodeOperators {
		if _, ok := nopIdsByAdmin[nop.Admin.String()]; !ok {
			return deployment.ChangesetOutput{}, fmt.Errorf("node operator %s does not have a nopId", nop.Name)
		}
	}
	// Adds initial set of nodes to CR, who all have the CCIP capability
	p2pIDsByNodeOpId := make(map[uint32][][32]byte)
	for nopAdmin, p2pId := range cfg.NodeP2PIDsPerNodeOpAdmin {
		nopId, ok := nopIdsByAdmin[nopAdmin.String()]
		if !ok {
			return deployment.ChangesetOutput{}, fmt.Errorf("node operator %s does not have a nopId", nopAdmin)
		}
		p2pIDsByNodeOpId[nopId] = p2pId
	}
	if err := ccipdeployment.AddNodes(
		env.Logger,
		capReg.Contract,
		env.Chains[homeChainSel],
		p2pIDsByNodeOpId,
	); err != nil {
		return deployment.ChangesetOutput{}, err
	}
	return deployment.ChangesetOutput{
		Proposals:   []timelock.MCMSWithTimelockProposal{},
		AddressBook: ab,
		JobSpecs:    nil,
	}, nil
}

type DeployHomeChainConfig struct {
	HomeChainSel             uint64
	RMNStaticConfig          rmn_home.RMNHomeStaticConfig
	RMNDynamicConfig         rmn_home.RMNHomeDynamicConfig
	NodeOperators            []capabilities_registry.CapabilitiesRegistryNodeOperator
	NodeP2PIDsPerNodeOpAdmin map[common.Address][][32]byte
}

func (c DeployHomeChainConfig) Validate() error {
	if c.HomeChainSel == 0 {
		return fmt.Errorf("home chain selector must be set")
	}
	if c.RMNDynamicConfig.OffchainConfig == nil {
		return fmt.Errorf("offchain config for RMNHomeDynamicConfig must be set")
	}
	if c.RMNStaticConfig.OffchainConfig == nil {
		return fmt.Errorf("offchain config for RMNHomeStaticConfig must be set")
	}
	if c.NodeOperators == nil || len(c.NodeOperators) == 0 {
		return fmt.Errorf("node operators must be set")
	}
	for _, nop := range c.NodeOperators {
		if nop.Admin == (common.Address{}) {
			return fmt.Errorf("node operator admin address must be set")
		}
		if nop.Name == "" {
			return fmt.Errorf("node operator name must be set")
		}
		if c.NodeP2PIDsPerNodeOpAdmin[nop.Admin] == nil || len(c.NodeP2PIDsPerNodeOpAdmin[nop.Admin]) == 0 {
			return fmt.Errorf("node operator %s must have node p2p ids provided", nop.Name)
		}
	}

	return nil
}

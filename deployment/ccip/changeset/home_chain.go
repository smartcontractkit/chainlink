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

var _ deployment.ChangeSet[DeployHomeChainConfig] = DeployHomeChain

// DeployHomeChain is a separate changeset because it is a standalone deployment performed once in home chain for the entire CCIP deployment.
func DeployHomeChain(env deployment.Environment, cfg DeployHomeChainConfig) (deployment.ChangesetOutput, error) {
	err := cfg.Validate()
	if err != nil {
		return deployment.ChangesetOutput{}, errors.Wrapf(deployment.ErrInvalidConfig, "%v", err)
	}
	ab := deployment.NewMemoryAddressBook()
	// Note we also deploy the cap reg.
	_, err = ccipdeployment.DeployHomeChain(env.Logger, env, ab, env.Chains[cfg.HomeChainSel], cfg.RMNStaticConfig, cfg.RMNDynamicConfig, cfg.NodeOperators, cfg.NodeP2PIDsPerNodeOpAdmin)
	if err != nil {
		env.Logger.Errorw("Failed to deploy cap reg", "err", err, "addresses", env.ExistingAddresses)
		return deployment.ChangesetOutput{
			AddressBook: ab,
		}, err
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
	NodeP2PIDsPerNodeOpAdmin map[string][][32]byte
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
	if len(c.NodeOperators) == 0 {
		return fmt.Errorf("node operators must be set")
	}
	for _, nop := range c.NodeOperators {
		if nop.Admin == (common.Address{}) {
			return fmt.Errorf("node operator admin address must be set")
		}
		if nop.Name == "" {
			return fmt.Errorf("node operator name must be set")
		}
		if len(c.NodeP2PIDsPerNodeOpAdmin[nop.Name]) == 0 {
			return fmt.Errorf("node operator %s must have node p2p ids provided", nop.Name)
		}
	}

	return nil
}

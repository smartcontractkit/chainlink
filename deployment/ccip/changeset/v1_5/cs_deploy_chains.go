package v1_5

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/mock_rmn_contract"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/rmn_contract"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/rmn_proxy_contract"
)

var _ deployment.ChangeSet[DeployChainContractsConfig1_5] = DeployChainContracts

type DeployChainContractsConfig1_5 struct {
	configs []DeployChainContractsConfig1_5PerChain
}

type DeployChainContractsConfig1_5PerChain struct {
	ChainSelector uint64
	RMNConfig     *rmn_contract.RMNConfig
}

func DeployChainContracts(env deployment.Environment, c DeployChainContractsConfig1_5) (deployment.ChangesetOutput, error) {
	if err := c.Validate(); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("invalid DeployChainContractsConfig1_5: %w", err)
	}
	newAddresses := deployment.NewMemoryAddressBook()
	err := deployChainContractsForChains(env, newAddresses, c)
	if err != nil {
		env.Logger.Errorw("Failed to deploy CCIP contracts", "err", err, "newAddresses", newAddresses)
		return deployment.ChangesetOutput{AddressBook: newAddresses}, deployment.MaybeDataErr(err)
	}
	return deployment.ChangesetOutput{
		Proposals:   []timelock.MCMSWithTimelockProposal{},
		AddressBook: newAddresses,
		JobSpecs:    nil,
	}, nil
}

func (c DeployChainContractsConfig1_5) Validate() error {
	for _, cs := range c.configs {
		if err := deployment.IsValidChainSelector(cs.ChainSelector); err != nil {
			return fmt.Errorf("invalid chain selector: %d - %w", cs, err)
		}
	}
	return nil
}

func deployChainContractsForChains(env deployment.Environment, ab deployment.AddressBook, cfg DeployChainContractsConfig1_5) error {
	state, err := changeset.LoadOnchainState(env)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}
	for _, c := range cfg.configs {
		if err := deployChainContracts(env, ab, state, c); err != nil {
			return fmt.Errorf("failed to deploy CCIP contracts for chain %d: %w", c.ChainSelector, err)
		}
	}
	return nil
}

func deployChainContracts(
	env deployment.Environment,
	ab deployment.AddressBook,
	state changeset.CCIPOnChainState,
	cfg DeployChainContractsConfig1_5PerChain) error {
	chain, ok := env.Chains[cfg.ChainSelector]
	if !ok {
		return fmt.Errorf("chain %d not found", cfg.ChainSelector)
	}
	chainState, ok := state.Chains[cfg.ChainSelector]
	if !ok {
		return fmt.Errorf("chain %d state not found", cfg.ChainSelector)
	}
	lggr := env.Logger
	var rmnAddr common.Address
	if cfg.RMNConfig == nil {
		if chainState.MockRMN == nil {
			rmn, err := deployment.DeployContract(lggr, chain, ab,
				func(chain deployment.Chain) deployment.ContractDeploy[*mock_rmn_contract.MockRMNContract] {
					rmnAddr, tx2, rmn, err2 := mock_rmn_contract.DeployMockRMNContract(
						chain.DeployerKey,
						chain.Client,
					)
					return deployment.ContractDeploy[*mock_rmn_contract.MockRMNContract]{
						rmnAddr, rmn, tx2, deployment.NewTypeAndVersion(changeset.MockRMN, deployment.Version1_0_0), err2,
					}
				})
			if err != nil {
				lggr.Errorw("Failed to deploy mock RMN", "chain", chain.String(), "err", err)
				return err
			}
			rmnAddr = rmn.Address
		} else {
			lggr.Infow("MockRMN already deployed", "chain", chain.String(), "address", chainState.MockRMN.Address())
			rmnAddr = chainState.MockRMN.Address()
		}
	} else {
		if chainState.RMN == nil {
			rmn, err := deployment.DeployContract(lggr, chain, ab,
				func(chain deployment.Chain) deployment.ContractDeploy[*rmn_contract.RMNContract] {
					rmnAddr, tx2, rmn, err2 := rmn_contract.DeployRMNContract(
						chain.DeployerKey,
						chain.Client,
						*cfg.RMNConfig,
					)
					return deployment.ContractDeploy[*rmn_contract.RMNContract]{
						rmnAddr, rmn, tx2, deployment.NewTypeAndVersion(changeset.RMN, deployment.Version1_5_0), err2,
					}
				})
			if err != nil {
				lggr.Errorw("Failed to deploy RMN", "chain", chain.String(), "err", err)
				return err
			}
			rmnAddr = rmn.Address
		} else {
			lggr.Infow("RMN already deployed", "chain", chain.String(), "address", chainState.RMN.Address)
			rmnAddr = chainState.RMN.Address()
		}
	}
	if chainState.RMNProxy == nil {
		_, err := deployment.DeployContract(lggr, chain, ab,
			func(chain deployment.Chain) deployment.ContractDeploy[*rmn_proxy_contract.RMNProxyContract] {
				rmnProxyAddr, tx2, rmnProxy, err2 := rmn_proxy_contract.DeployRMNProxyContract(
					chain.DeployerKey,
					chain.Client,
					rmnAddr,
				)
				return deployment.ContractDeploy[*rmn_proxy_contract.RMNProxyContract]{
					rmnProxyAddr, rmnProxy, tx2, deployment.NewTypeAndVersion(changeset.ARMProxy, deployment.Version1_0_0), err2,
				}
			})
		if err != nil {
			lggr.Errorw("Failed to deploy RMNProxy", "chain", chain.String(), "err", err)
			return err
		}
	} else {
		lggr.Infow("RMNProxy already deployed", "chain", chain.String(), "address", chainState.RMNProxy.Address)
		// check if the RMNProxy is pointing to the correct RMN
		rmnProxy := chainState.RMNProxy
		setRMN, err := rmnProxy.GetARM(nil)
		if err != nil {
			return err
		}
		if setRMN != rmnAddr {
			lggr.Infow("RMNProxy pointing to wrong RMN, changing the RMN", "chain", chain.String(), "rmnProxy", rmnProxy.Address, "rmn", rmnAddr)
			tx, err := rmnProxy.SetARM(chain.DeployerKey, rmnAddr)
			if err != nil {
				return fmt.Errorf("failed to set RMN on RMNProxy for chain %s RMN %s RMNProxy %s: %w",
					chain.String(), rmnAddr.String(), rmnProxy.Address().String(), err)
			}

			_, err = chain.Confirm(tx)
			if err != nil {
				lggr.Errorw("Failed to confirm RMNProxy SetARM",
					"chain", chain.String(), "rmn", rmnAddr.String(), "rmnProxy", rmnProxy.Address().String(), "err", err)
				return fmt.Errorf("failed to confirm RMNProxy SetARM for chain %s RMN %s RMNProxy %s: %w",
					chain.String(), rmnAddr.String(), rmnProxy.Address().String(), err)
			}
			lggr.Infow("RMNProxy SetARM confirmed",
				"chain", chain.String(), "rmn", rmnAddr.String(), "rmnProxy", rmnProxy.Address().String())
		}
	}
	return nil
}

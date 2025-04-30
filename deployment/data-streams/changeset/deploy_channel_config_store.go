package changeset

import (
	"errors"
	"fmt"

	"github.com/smartcontractkit/chainlink/deployment"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/channel_config_store"

	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
)

// DeployChannelConfigStore deploys ChannelConfigStore to the chains specified in the config.
type DeployChannelConfigStore struct{}

type DeployChannelConfigStoreConfig struct {
	// ChainsToDeploy is a list of chain selectors to deploy the contract to.
	ChainsToDeploy []uint64
}

func (cc DeployChannelConfigStoreConfig) Validate() error {
	if len(cc.ChainsToDeploy) == 0 {
		return errors.New("ChainsToDeploy is empty")
	}
	for _, chain := range cc.ChainsToDeploy {
		if err := deployment.IsValidChainSelector(chain); err != nil {
			return fmt.Errorf("invalid chain selector: %d - %w", chain, err)
		}
	}
	return nil
}

func (DeployChannelConfigStore) Apply(e deployment.Environment, cc DeployChannelConfigStoreConfig) (deployment.ChangesetOutput, error) {
	ab := deployment.NewMemoryAddressBook()
	err := performDeployment(e, ab, cc)
	if err != nil {
		e.Logger.Errorw("Failed to deploy ChannelConfigStore", "err", err, "addresses", ab)
		return deployment.ChangesetOutput{AddressBook: ab}, deployment.MaybeDataErr(err)
	}
	return deployment.ChangesetOutput{
		AddressBook: ab,
	}, nil
}

func (DeployChannelConfigStore) VerifyPreconditions(_ deployment.Environment, cc DeployChannelConfigStoreConfig) error {
	if err := cc.Validate(); err != nil {
		return fmt.Errorf("invalid DeployChannelConfigStoreConfig: %w", err)
	}

	return nil
}

func performDeployment(e deployment.Environment, ab deployment.AddressBook, cc DeployChannelConfigStoreConfig) error {
	if err := cc.Validate(); err != nil {
		return fmt.Errorf("invalid DeployChannelConfigStoreConfig: %w", err)
	}

	for _, chainSel := range cc.ChainsToDeploy {
		chain, ok := e.Chains[chainSel]
		if !ok {
			return fmt.Errorf("Chain not found for chain selector %d", chainSel)
		}
		_, err := DeployContract[*channel_config_store.ChannelConfigStore](e, ab, chain, channelConfigStoreDeployFn())
		if err != nil {
			return err
		}
		chainAddresses, err := ab.AddressesForChain(chain.Selector)
		if err != nil {
			e.Logger.Errorw("Failed to get chain addresses", "err", err)
			return err
		}
		chainState, err := LoadChainState(e.Logger, chain, chainAddresses)
		if err != nil {
			e.Logger.Errorw("Failed to load chain state", "err", err)
			return err
		}
		if len(chainState.ChannelConfigStores) == 0 {
			errNoCCS := errors.New("no ChannelConfigStore on chain")
			e.Logger.Error(errNoCCS)
			return errNoCCS
		}
	}

	return nil
}

// channelConfigStoreDeployFn returns a function that deploys a ChannelConfigStore contract.
func channelConfigStoreDeployFn() ContractDeployFn[*channel_config_store.ChannelConfigStore] {
	return func(chain deployment.Chain) *ContractDeployment[*channel_config_store.ChannelConfigStore] {
		ccsAddr, ccsTx, ccs, err := channel_config_store.DeployChannelConfigStore(
			chain.DeployerKey,
			chain.Client,
		)
		if err != nil {
			return &ContractDeployment[*channel_config_store.ChannelConfigStore]{
				Err: err,
			}
		}
		return &ContractDeployment[*channel_config_store.ChannelConfigStore]{
			Address:  ccsAddr,
			Contract: ccs,
			Tx:       ccsTx,
			Tv:       deployment.NewTypeAndVersion(types.ChannelConfigStore, deployment.Version1_0_0),
			Err:      nil,
		}
	}
}

package v0_5_0

import (
	"errors"
	"fmt"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/configurator"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
)

var DeployConfiguratorChangeset = cldf.CreateChangeSet(deployConfiguratorLogic, deployConfiguratorPrecondition)

type DeployConfiguratorConfig struct {
	ChainsToDeploy []uint64
}

func (cc DeployConfiguratorConfig) Validate() error {
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

func deployConfiguratorLogic(e deployment.Environment, cc DeployConfiguratorConfig) (deployment.ChangesetOutput, error) {
	ab := deployment.NewMemoryAddressBook()
	err := deploy(e, ab, cc)
	if err != nil {
		e.Logger.Errorw("Failed to deploy Configurator", "err", err, "addresses", ab)
		return deployment.ChangesetOutput{AddressBook: ab}, deployment.MaybeDataErr(err)
	}
	return deployment.ChangesetOutput{
		AddressBook: ab,
	}, nil
}

func deployConfiguratorPrecondition(_ deployment.Environment, cc DeployConfiguratorConfig) error {
	if err := cc.Validate(); err != nil {
		return fmt.Errorf("invalid DeployConfiguratorConfig: %w", err)
	}

	return nil
}

func deploy(e deployment.Environment, ab deployment.AddressBook, cc DeployConfiguratorConfig) error {
	for _, chainSel := range cc.ChainsToDeploy {
		chain, ok := e.Chains[chainSel]
		if !ok {
			return fmt.Errorf("chain not found for chain selector %d", chainSel)
		}
		_, err := changeset.DeployContract(e, ab, chain, DeployFn())
		if err != nil {
			return err
		}
		chainAddresses, err := ab.AddressesForChain(chain.Selector)
		if err != nil {
			e.Logger.Errorw("Failed to get chain addresses", "err", err)
			return err
		}
		chainState, err := changeset.LoadChainState(e.Logger, chain, chainAddresses)
		if err != nil {
			e.Logger.Errorw("Failed to load chain state", "err", err)
			return err
		}
		if len(chainState.Configurators) == 0 {
			errNoCCS := errors.New("no Configurator on chain")
			e.Logger.Error(errNoCCS)
			return errNoCCS
		}
	}

	return nil
}

func DeployFn() changeset.ContractDeployFn[*configurator.Configurator] {
	return func(chain deployment.Chain) *changeset.ContractDeployment[*configurator.Configurator] {
		ccsAddr, ccsTx, ccs, err := configurator.DeployConfigurator(
			chain.DeployerKey,
			chain.Client,
		)
		if err != nil {
			return &changeset.ContractDeployment[*configurator.Configurator]{
				Err: err,
			}
		}
		return &changeset.ContractDeployment[*configurator.Configurator]{
			Address:  ccsAddr,
			Contract: ccs,
			Tx:       ccsTx,
			Tv:       deployment.NewTypeAndVersion(types.Configurator, deployment.Version0_5_0),
			Err:      nil,
		}
	}
}

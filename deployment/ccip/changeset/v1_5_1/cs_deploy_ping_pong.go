package v1_5_1

import (
	"fmt"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/latest/ping_pong_demo"
)

var DeployPingPongDemoContractChangeset = deployment.CreateChangeSet(deployPingPongDemoContractsChangeset, validateDeployPingPongConfig)
var StartPingPongDemoContractChangeset = deployment.CreateChangeSet(startPingPongDemoContractsChangeset, validateStartPingPongContractAddress)
var SetPausedPingPongDemoContractChangeset = deployment.CreateChangeSet(setPausedPingPongDemoContractsChangeset, validateSetPausedPingPongContractAddress)
var SetCounterpartPingPongDemoContractChangeset = deployment.CreateChangeSet(setCounterpartPingPongDemoContractsChangeset, validateSetCounterpartPingPongContractAddress)

type DeployPingPongDemoContractsConfig struct {
	ChainsToDeploy []struct {
		ChainSelector uint64
		IsTestRouter  bool
	}
}

func validateDeployPingPongConfig(env deployment.Environment, config DeployPingPongDemoContractsConfig) error {
	state, err := changeset.LoadOnchainState(env)

	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}

	for _, chainToDeploy := range config.ChainsToDeploy {
		err = changeset.ValidateChain(env, state, chainToDeploy.ChainSelector, nil)
		if err != nil {
			return fmt.Errorf("failed to validate chain for %d: %w", chainToDeploy.ChainSelector, err)
		}

		chainState := state.Chains[chainToDeploy.ChainSelector]

		router := chainState.Router
		if chainToDeploy.IsTestRouter {
			router = chainState.TestRouter
		}

		if router == nil {
			return fmt.Errorf("router address is empty for chain %d", chainToDeploy.ChainSelector)
		}

		if chainState.LinkToken == nil {
			return fmt.Errorf("link token address is empty for chain %d", chainToDeploy.ChainSelector)
		}
	}

	return nil
}

func deployPingPongDemoContractsChangeset(env deployment.Environment, c DeployPingPongDemoContractsConfig) (deployment.ChangesetOutput, error) {
	state, err := changeset.LoadOnchainState(env)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	newAB := deployment.NewMemoryAddressBook()

	for _, chainToDeploy := range c.ChainsToDeploy {
		chain := env.Chains[chainToDeploy.ChainSelector]
		chainState := state.Chains[chainToDeploy.ChainSelector]

		router := chainState.Router
		if chainToDeploy.IsTestRouter {
			router = chainState.TestRouter
		}

		linkTokenAddress, err := chainState.LinkTokenAddress()
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to get link token address for chain: %d %w", chainToDeploy.ChainSelector, err)
		}

		dep, err := deployment.DeployContract(env.Logger, chain, newAB,
			func(chain deployment.Chain) deployment.ContractDeploy[*ping_pong_demo.PingPongDemo] {
				addr, tx, pingPongDemo, err := ping_pong_demo.DeployPingPongDemo(chain.DeployerKey, chain.Client, router.Address(), linkTokenAddress)

				return deployment.ContractDeploy[*ping_pong_demo.PingPongDemo]{
					Address:  addr,
					Contract: pingPongDemo,
					Tx:       tx,
					Tv:       deployment.NewTypeAndVersion(changeset.PingPongDemo, deployment.Version1_0_0),
					Err:      err,
				}

			},
		)

		if _, err := deployment.ConfirmIfNoErrorWithABI(chain, dep.Tx, ping_pong_demo.PingPongDemoABI, err); err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to confirm ping pong demo contract deployment tx: %w", err)
		}
	}

	return deployment.ChangesetOutput{
		AddressBook: newAB,
	}, nil
}

type StartPingPongDemoContractsConfig struct {
	ChainSelector uint64
}

func validateStartPingPongContractAddress(env deployment.Environment, config StartPingPongDemoContractsConfig) error {
	return validatePingPongContractAddress(env, config.ChainSelector)
}

func startPingPongDemoContractsChangeset(env deployment.Environment, c StartPingPongDemoContractsConfig) (deployment.ChangesetOutput, error) {
	state, err := changeset.LoadOnchainState(env)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	if err := validatePingPongContractAddress(env, c.ChainSelector); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("invalid ping pong contract config: %w", err)
	}

	chain := env.Chains[c.ChainSelector]
	chainState := state.Chains[c.ChainSelector]

	transactor, err := ping_pong_demo.NewPingPongDemoTransactor(chainState.PingPongDemo.Address(), chain.Client)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to create transactor for ping pong demo: %w", err)
	}

	tx, err := transactor.StartPingPong(chain.DeployerKey)
	if _, err := deployment.ConfirmIfNoErrorWithABI(chain, tx, ping_pong_demo.PingPongDemoABI, err); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to confirm ping pong demo contract start tx: %w", err)
	}

	return deployment.ChangesetOutput{}, nil
}

type SetPausedPingPongDemoContractsConfig struct {
	ChainSelector uint64
	Paused        bool
}

func validateSetPausedPingPongContractAddress(env deployment.Environment, config SetPausedPingPongDemoContractsConfig) error {
	return validatePingPongContractAddress(env, config.ChainSelector)
}
func setPausedPingPongDemoContractsChangeset(env deployment.Environment, c SetPausedPingPongDemoContractsConfig) (deployment.ChangesetOutput, error) {
	state, err := changeset.LoadOnchainState(env)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	if err := validatePingPongContractAddress(env, c.ChainSelector); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("invalid ping pong contract config: %w", err)
	}

	chain := env.Chains[c.ChainSelector]
	chainState := state.Chains[c.ChainSelector]

	transactor, err := ping_pong_demo.NewPingPongDemoTransactor(chainState.PingPongDemo.Address(), chain.Client)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to create transactor for ping pong demo: %w", err)
	}

	tx, err := transactor.SetPaused(chain.DeployerKey, c.Paused)
	if _, err := deployment.ConfirmIfNoErrorWithABI(chain, tx, ping_pong_demo.PingPongDemoABI, err); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to confirm ping pong demo contract set paused tx: %w", err)
	}

	return deployment.ChangesetOutput{}, nil
}

type SetConterpartPingPongDemoContractsConfig struct {
	ChainSelector      uint64
	CounterpartAddress []byte
}

func validateSetCounterpartPingPongContractAddress(env deployment.Environment, config SetConterpartPingPongDemoContractsConfig) error {
	return validatePingPongContractAddress(env, config.ChainSelector)
}
func setCounterpartPingPongDemoContractsChangeset(env deployment.Environment, c SetConterpartPingPongDemoContractsConfig) (deployment.ChangesetOutput, error) {
	state, err := changeset.LoadOnchainState(env)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	if err := validatePingPongContractAddress(env, c.ChainSelector); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("invalid ping pong contract config: %w", err)
	}

	chain := env.Chains[c.ChainSelector]
	chainState := state.Chains[c.ChainSelector]

	transactor, err := ping_pong_demo.NewPingPongDemoTransactor(chainState.PingPongDemo.Address(), chain.Client)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to create transactor for ping pong demo: %w", err)
	}

	tx, err := transactor.SetCounterpart(chain.DeployerKey, c.ChainSelector, c.CounterpartAddress)
	if _, err := deployment.ConfirmIfNoErrorWithABI(chain, tx, ping_pong_demo.PingPongDemoABI, err); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to confirm ping pong demo contract set counterpart tx: %w", err)
	}

	return deployment.ChangesetOutput{}, nil
}

func validatePingPongContractAddress(env deployment.Environment, chainSelector uint64) error {
	state, err := changeset.LoadOnchainState(env)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}

	err = changeset.ValidateChain(env, state, chainSelector, nil)

	chainState := state.Chains[chainSelector]

	if chainState.PingPongDemo == nil {
		return fmt.Errorf("ping pong demo address is empty")
	}

	return nil
}

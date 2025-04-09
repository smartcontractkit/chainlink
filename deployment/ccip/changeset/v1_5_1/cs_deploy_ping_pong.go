package v1_5_1

import (
	"fmt"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/latest/ping_pong_demo"
)

var _ deployment.ChangeSet[DeployPingPongDemoContractsConfig] = DeployPingPongDemoContractsChangeset
var _ deployment.ChangeSet[StartPingPongDemoContractsConfig] = StartPingPongDemoContractsChangeset
var _ deployment.ChangeSet[SetPausedPingPongDemoContractsConfig] = SetPausedPingPongDemoContractsChangeset
var _ deployment.ChangeSet[SetConterpartPingPongDemoContractsConfig] = SetCounterpartPingPongDemoContractsChangeset

type DeployPingPongDemoContractsConfig struct {
	ChainSelector uint64
	IsTestRouter  bool
}

func (c DeployPingPongDemoContractsConfig) Validate(env deployment.Environment, state changeset.CCIPOnChainState) error {
	chainState := state.Chains[c.ChainSelector]

	router := chainState.Router
	if c.IsTestRouter {
		router = chainState.TestRouter
	}

	if router == nil {
		return fmt.Errorf("router address is empty for chain %d", c.ChainSelector)
	}

	if chainState.LinkToken == nil {
		return fmt.Errorf("link token address is empty for chain %d", c.ChainSelector)
	}

	return nil
}

func DeployPingPongDemoContractsChangeset(env deployment.Environment, c DeployPingPongDemoContractsConfig) (deployment.ChangesetOutput, error) {
	state, err := changeset.LoadOnchainState(env)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	if err := c.Validate(env, state); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("invalid DeployPingPongDemoContractsConfig: %w", err)
	}

	newAB := deployment.NewMemoryAddressBook()

	chain := env.Chains[c.ChainSelector]
	chainState := state.Chains[c.ChainSelector]

	router := chainState.Router
	if c.IsTestRouter {
		router = chainState.TestRouter
	}

	tv := deployment.NewTypeAndVersion(changeset.PingPongDemo, deployment.Version1_0_0)

	addr, _, _, err := ping_pong_demo.DeployPingPongDemo(chain.DeployerKey, chain.Client, router.Address(), chainState.LinkToken.Address())
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to deploy %s token pool on %w", err)
	}

	err = newAB.Save(chain.Selector, addr.String(), tv)

	return deployment.ChangesetOutput{
		AddressBook: newAB,
	}, nil
}

type StartPingPongDemoContractsConfig struct {
	ChainSelector uint64
}

func (c StartPingPongDemoContractsConfig) Validate(env deployment.Environment, state changeset.CCIPOnChainState) error {
	chainState := state.Chains[c.ChainSelector]

	if chainState.PingPongDemo == nil {
		return fmt.Errorf("ping pong demo address is empty for chain %d", c.ChainSelector)
	}

	return nil
}

func StartPingPongDemoContractsChangeset(env deployment.Environment, c StartPingPongDemoContractsConfig) (deployment.ChangesetOutput, error) {
	state, err := changeset.LoadOnchainState(env)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	if err := c.Validate(env, state); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("invalid StartPingPongDemoContractsConfig: %w", err)
	}

	chain := env.Chains[c.ChainSelector]
	chainState := state.Chains[c.ChainSelector]

	transactor, err := ping_pong_demo.NewPingPongDemoTransactor(chainState.PingPongDemo.Address(), chain.Client)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to create transactor for ping pong demo: %w", err)
	}

	_, err = transactor.StartPingPong(chain.DeployerKey)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to start ping pong demo: %w", err)
	}

	return deployment.ChangesetOutput{}, nil
}

type SetPausedPingPongDemoContractsConfig struct {
	ChainSelector uint64
	Paused        bool
}

func (c SetPausedPingPongDemoContractsConfig) Validate(env deployment.Environment, state changeset.CCIPOnChainState) error {
	chainState := state.Chains[c.ChainSelector]

	if chainState.PingPongDemo == nil {
		return fmt.Errorf("ping pong demo address is empty for chain %d", c.ChainSelector)
	}

	return nil
}

func SetPausedPingPongDemoContractsChangeset(env deployment.Environment, c SetPausedPingPongDemoContractsConfig) (deployment.ChangesetOutput, error) {
	state, err := changeset.LoadOnchainState(env)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	if err := c.Validate(env, state); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("invalid SetPausedPingPongDemoContractsConfig: %w", err)
	}

	chain := env.Chains[c.ChainSelector]
	chainState := state.Chains[c.ChainSelector]

	transactor, err := ping_pong_demo.NewPingPongDemoTransactor(chainState.PingPongDemo.Address(), chain.Client)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to create transactor for ping pong demo: %w", err)
	}

	_, err = transactor.SetPaused(chain.DeployerKey, c.Paused)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to set paused state for ping pong demo: %w", err)
	}

	return deployment.ChangesetOutput{}, nil
}

type SetConterpartPingPongDemoContractsConfig struct {
	ChainSelector      uint64
	CounterpartAddress []byte
}

func (c SetConterpartPingPongDemoContractsConfig) Validate(env deployment.Environment, state changeset.CCIPOnChainState) error {
	chainState := state.Chains[c.ChainSelector]

	if chainState.PingPongDemo == nil {
		return fmt.Errorf("ping pong demo address is empty for chain %d", c.ChainSelector)
	}

	return nil
}

func SetCounterpartPingPongDemoContractsChangeset(env deployment.Environment, c SetConterpartPingPongDemoContractsConfig) (deployment.ChangesetOutput, error) {
	state, err := changeset.LoadOnchainState(env)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	if err := c.Validate(env, state); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("invalid SetConterpartPingPongDemoContractsConfig: %w", err)
	}

	chain := env.Chains[c.ChainSelector]
	chainState := state.Chains[c.ChainSelector]

	transactor, err := ping_pong_demo.NewPingPongDemoTransactor(chainState.PingPongDemo.Address(), chain.Client)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to create transactor for ping pong demo: %w", err)
	}

	_, err = transactor.SetCounterpart(chain.DeployerKey, c.ChainSelector, c.CounterpartAddress)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to set counterpart for ping pong demo: %w", err)
	}

	return deployment.ChangesetOutput{}, nil
}

package v1_5

import (
	"fmt"
	"math/big"

	"github.com/gagliardetto/solana-go"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/latest/ping_pong_demo"
	solanaUtilsCcip "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/ccip"
	solanaStateUtils "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/tokens"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	ccipChangeset "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/testhelpers"
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
		chainState := state.Chains[chainToDeploy.ChainSelector]

		router := chainState.Router
		if chainToDeploy.IsTestRouter {
			router = chainState.TestRouter
		}

		if router == nil {
			return fmt.Errorf("router address is empty for chain %d", chainToDeploy.ChainSelector)
		}

		_, err := chainState.LinkTokenAddress()
		if err != nil {
			return fmt.Errorf("failed to get link token address for chain: %d %w", chainToDeploy.ChainSelector, err)
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

type SetEVMCounterpartExtraArgsPingPongDemoContracts struct {
	GasLimit                 *big.Int
	AllowOutOfOrderExecution bool
}
type SetSolanaCounterpartExtraArgsPingPongDemoContracts struct {
	ComputeUnits uint32
}
type SetCounterpartPingPongDemoContractsConfig struct {
	ChainSelector            uint64
	CounterpartChainSelector uint64
	ExtraArgsEVM             *SetEVMCounterpartExtraArgsPingPongDemoContracts
	ExtraArgsSolana          *SetSolanaCounterpartExtraArgsPingPongDemoContracts
}

func validateSetCounterpartPingPongContractAddress(env deployment.Environment, config SetCounterpartPingPongDemoContractsConfig) error {
	return validatePingPongContractAddress(env, config.ChainSelector)
}

func setCounterpartPingPongDemoContractsChangeset(env deployment.Environment, c SetCounterpartPingPongDemoContractsConfig) (deployment.ChangesetOutput, error) {
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

	extraArgsBytes := make([]byte, 0)

	if (c.ExtraArgsEVM != nil) == (c.ExtraArgsSolana != nil) {
		return deployment.ChangesetOutput{}, fmt.Errorf("exactly one of ExtraArgsEVM or ExtraArgsSolana must be set")
	}

	if c.ExtraArgsEVM != nil {
		extraArgsBytes, err = getEVMExtraArgs(env, c)
	}

	if c.ExtraArgsSolana != nil {
		extraArgsBytes, err = getSolanaExtraArgs(env, c)
	}

	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get extra args: %w", err)
	}

	counterpartAddressStr, err := deployment.SearchAddressBook(env.ExistingAddresses, c.CounterpartChainSelector, changeset.PingPongDemo)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get counterpart address: %w", err)
	}
	counterpartAddressBytes := make([]byte, 32)
	copy(counterpartAddressBytes[32-len(counterpartAddressStr):], counterpartAddressStr)

	tx, err := transactor.SetCounterpart(chain.DeployerKey, c.CounterpartChainSelector, counterpartAddressBytes, extraArgsBytes)
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

	chainState := state.Chains[chainSelector]

	if chainState.PingPongDemo == nil {
		return fmt.Errorf("ping pong demo address is empty for chain %d", chainSelector)
	}

	return nil
}

func getEVMExtraArgs(env deployment.Environment, c SetCounterpartPingPongDemoContractsConfig) ([]byte, error) {
	b, err := testhelpers.GetEVMExtraArgsV2(c.ExtraArgsEVM.GasLimit, c.ExtraArgsEVM.AllowOutOfOrderExecution)

	if err != nil {
		return nil, fmt.Errorf("failed to marshal extra args for solana: %w", err)
	}

	return b, nil
}
func getSolanaExtraArgs(env deployment.Environment, c SetCounterpartPingPongDemoContractsConfig) ([]byte, error) {
	s, err := ccipChangeset.LoadOnchainStateSolana(env)
	if err != nil {
		env.Logger.Errorw("Failed to load existing onchain state", "err", err)
		return nil, err
	}

	destChainState := s.SolChains[c.CounterpartChainSelector]

	ppConfigPDA, _, err := solanaStateUtils.FindPingPongDemoConfigPDA(destChainState.PingPong)
	if err != nil {
		return nil, fmt.Errorf("failed to find ping pong demo config PDA: %w", err)
	}

	ppSendSignerPDA, _, err := solanaStateUtils.FindPingPongCCIPSendSignerPDA(destChainState.PingPong)
	if err != nil {
		return nil, fmt.Errorf("failed to find ping pong CCIP send signer PDA: %w", err)
	}

	destChainStatePDA, err := solanaStateUtils.FindDestChainStatePDA(c.ChainSelector, destChainState.Router)
	if err != nil {
		return nil, fmt.Errorf("failed to find destination chain state PDA: %w", err)
	}

	routerNoncePDA, err := solanaStateUtils.FindNoncePDA(c.ChainSelector, ppSendSignerPDA, destChainState.Router)
	if err != nil {
		return nil, fmt.Errorf("failed to find router nonce PDA: %w", err)
	}

	feeBillingSignerPDA, _, err := solanaStateUtils.FindFeeBillingSignerPDA(destChainState.Router)
	if err != nil {
		return nil, fmt.Errorf("failed to find fee billing signer PDA: %w", err)
	}

	fqBillingTokenConfigPDA, _, err := solanaStateUtils.FindFqBillingTokenConfigPDA(destChainState.LinkToken, destChainState.FeeQuoter)
	if err != nil {
		return nil, fmt.Errorf("failed to find fee quoter billing token config PDA: %w", err)
	}

	rmnRemoteConfigPDA, _, err := solanaStateUtils.FindRMNRemoteConfigPDA(destChainState.RMNRemote)
	if err != nil {
		return nil, fmt.Errorf("failed to find RMN remote config PDA: %w", err)
	}

	rmnRemoteCursesPDA, _, err := solanaStateUtils.FindRMNRemoteCursesPDA(destChainState.RMNRemote)
	if err != nil {
		return nil, fmt.Errorf("failed to find RMN remote curses PDA: %w", err)
	}

	ppFeeTokenAta, _, err := tokens.FindAssociatedTokenAddress(solana.Token2022ProgramID, destChainState.LinkToken, ppSendSignerPDA)
	if err != nil {
		return nil, fmt.Errorf("failed to find ping pong fee token associated token account: %w", err)
	}

	routerFeeTokenReceiver, _, err := tokens.FindAssociatedTokenAddress(solana.Token2022ProgramID, destChainState.LinkToken, feeBillingSignerPDA)
	if err != nil {
		return nil, fmt.Errorf("failed to find router fee token receiver: %w", err)
	}

	fqDestChainPDA, _, err := solanaStateUtils.FindFqDestChainPDA(c.ChainSelector, destChainState.FeeQuoter)
	if err != nil {
		return nil, fmt.Errorf("failed to find fee quoter destination chain PDA: %w", err)
	}

	fqLinkTokenConfigPDA, _, err := solanaStateUtils.FindFqBillingTokenConfigPDA(destChainState.LinkToken, destChainState.FeeQuoter)
	if err != nil {
		return nil, fmt.Errorf("failed to find fee quoter link token config PDA: %w", err)
	}

	accounts := []solana.PublicKey{
		ppConfigPDA,
		ppSendSignerPDA,
		destChainStatePDA,
		routerNoncePDA,
		feeBillingSignerPDA,
		fqBillingTokenConfigPDA,
		rmnRemoteConfigPDA,
		rmnRemoteCursesPDA,
		ppFeeTokenAta,
		routerFeeTokenReceiver,
		fqDestChainPDA,
		fqLinkTokenConfigPDA,
	}

	writableBitmap := solanaUtilsCcip.GenerateBitMapForIndexes([]int{1, 4, 7, 8, 9})

	b, err := testhelpers.GetSVMExtraArgsV1(c.ExtraArgsSolana.ComputeUnits, writableBitmap, accounts)

	if err != nil {
		return nil, fmt.Errorf("failed to marshal extra args for solana: %w", err)
	}

	return b, nil
}

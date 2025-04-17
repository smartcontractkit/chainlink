package v1_5

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/latest/ping_pong_demo"
	solanaUtilsCcip "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/ccip"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	ccipChangeset "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/testhelpers"
)

var DeployPingPongDemoContractChangeset = deployment.CreateChangeSet(deployPingPongDemoContractsChangeset, validateDeployPingPongConfig)
var StartPingPongDemoContractChangeset = deployment.CreateChangeSet(startPingPongDemoContractsChangeset, validateStartPingPongContractAddress)
var SetPausedPingPongDemoContractChangeset = deployment.CreateChangeSet(setPausedPingPongDemoContractsChangeset, validateSetPausedPingPongContractAddress)
var SetCounterpartPingPongDemoContractChangeset = deployment.CreateChangeSet(setCounterpartPingPongDemoContractsChangeset, validateSetCounterpartPingPongContractAddress)

func validatePingPongContractAddress(env deployment.Environment, chainSelector uint64) error {
	return nil
}

type DeployPingPongDemoContractsConfig struct {
	ChainSelector            uint64
	CounterpartChainSelector uint64
	IsTestRouter             bool
}

func validateDeployPingPongConfig(env deployment.Environment, config DeployPingPongDemoContractsConfig) error {
	state, err := changeset.LoadOnchainState(env)

	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}

	chainState := state.Chains[config.ChainSelector]

	router := chainState.Router
	if config.IsTestRouter {
		router = chainState.TestRouter
	}

	if router == nil {
		return fmt.Errorf("router address is empty for chain %d", config.ChainSelector)
	}

	_, err = chainState.LinkTokenAddress()
	if err != nil {
		return fmt.Errorf("failed to get link token address for chain: %d %w", config.ChainSelector, err)
	}

	return nil
}
func deployPingPongDemoContractsChangeset(env deployment.Environment, c DeployPingPongDemoContractsConfig) (deployment.ChangesetOutput, error) {
	state, err := changeset.LoadOnchainState(env)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	newAB := deployment.NewMemoryAddressBook()

	chain := env.Chains[c.ChainSelector]
	chainState := state.Chains[c.ChainSelector]

	router := chainState.Router
	if c.IsTestRouter {
		router = chainState.TestRouter
	}
	if router == nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("router address is empty for chain %d", c.ChainSelector)
	}

	linkTokenAddress, err := chainState.LinkTokenAddress()
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get link token address for chain: %d %w", c.ChainSelector, err)
	}

	dep, err := deployment.DeployContract(env.Logger, chain, newAB,
		func(chain deployment.Chain) deployment.ContractDeploy[*ping_pong_demo.PingPongDemo] {
			addr, tx, pingPongDemo, err := ping_pong_demo.DeployPingPongDemo(chain.DeployerKey, chain.Client, router.Address(), linkTokenAddress)

			tv := deployment.NewTypeAndVersion(changeset.PingPongDemo, deployment.Version1_0_0)
			tv.Labels.Add(fmt.Sprintf("To - %d", c.CounterpartChainSelector))

			return deployment.ContractDeploy[*ping_pong_demo.PingPongDemo]{
				Address:  addr,
				Contract: pingPongDemo,
				Tx:       tx,
				Tv:       tv,
				Err:      err,
			}
		},
	)

	if _, err := deployment.ConfirmIfNoErrorWithABI(chain, dep.Tx, ping_pong_demo.PingPongDemoABI, err); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to confirm ping pong demo contract deployment tx: %w", err)
	}

	return deployment.ChangesetOutput{
		AddressBook: newAB,
	}, nil
}

type StartPingPongDemoContractsConfig struct {
	ChainSelector            uint64
	CounterpartChainSelector uint64
}

func validateStartPingPongContractAddress(env deployment.Environment, config StartPingPongDemoContractsConfig) error {
	return validatePingPongContractAddress(env, config.ChainSelector)
}
func startPingPongDemoContractsChangeset(env deployment.Environment, c StartPingPongDemoContractsConfig) (deployment.ChangesetOutput, error) {
	if err := validatePingPongContractAddress(env, c.ChainSelector); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("invalid ping pong contract config: %w", err)
	}

	chain := env.Chains[c.ChainSelector]

	contractAddressStr, err := ccipChangeset.GetPingPongDemoContractAddress(env, c.ChainSelector, c.CounterpartChainSelector)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get ping pong demo contract address: %w", err)
	}

	transactor, err := ping_pong_demo.NewPingPongDemoTransactor(common.HexToAddress(contractAddressStr), chain.Client)
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
	ChainSelector            uint64
	CounterpartChainSelector uint64
	IsPaused                 bool
}

func validateSetPausedPingPongContractAddress(env deployment.Environment, config SetPausedPingPongDemoContractsConfig) error {
	return validatePingPongContractAddress(env, config.ChainSelector)
}
func setPausedPingPongDemoContractsChangeset(env deployment.Environment, c SetPausedPingPongDemoContractsConfig) (deployment.ChangesetOutput, error) {
	if err := validatePingPongContractAddress(env, c.ChainSelector); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("invalid ping pong contract config: %w", err)
	}

	chain := env.Chains[c.ChainSelector]

	contractAddressStr, err := ccipChangeset.GetPingPongDemoContractAddress(env, c.ChainSelector, c.CounterpartChainSelector)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get ping pong demo contract address: %w", err)
	}

	transactor, err := ping_pong_demo.NewPingPongDemoTransactor(common.HexToAddress(contractAddressStr), chain.Client)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to create transactor for ping pong demo: %w", err)
	}

	tx, err := transactor.SetPaused(chain.DeployerKey, c.IsPaused)
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

	FeesTokenProgram  solana.PublicKey
	FeesTokenMintType deployment.ContractType
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
	if err := validatePingPongContractAddress(env, c.ChainSelector); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("invalid ping pong contract config: %w", err)
	}

	chain := env.Chains[c.ChainSelector]

	contractAddressStr, err := ccipChangeset.GetPingPongDemoContractAddress(env, c.ChainSelector, c.CounterpartChainSelector)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get ping pong demo contract address: %w", err)
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

	counterpartAddressBytes, err := ccipChangeset.GetPaddedPingPongAddressBytes(env, c.CounterpartChainSelector, c.ChainSelector, chainsel.FamilyEVM)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get counterpart address: %w", err)
	}

	transactor, err := ping_pong_demo.NewPingPongDemoTransactor(common.HexToAddress(contractAddressStr), chain.Client)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to create transactor for ping pong demo: %w", err)
	}

	tx, err := transactor.SetCounterpart(chain.DeployerKey, c.CounterpartChainSelector, counterpartAddressBytes, extraArgsBytes)
	if _, err := deployment.ConfirmIfNoErrorWithABI(chain, tx, ping_pong_demo.PingPongDemoABI, err); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to confirm ping pong demo contract set counterpart tx: %w", err)
	}

	return deployment.ChangesetOutput{}, nil
}

func getEVMExtraArgs(env deployment.Environment, c SetCounterpartPingPongDemoContractsConfig) ([]byte, error) {
	b, err := testhelpers.GetEVMExtraArgsV2(c.ExtraArgsEVM.GasLimit, c.ExtraArgsEVM.AllowOutOfOrderExecution)

	if err != nil {
		return nil, fmt.Errorf("failed to marshal extra args for solana: %w", err)
	}

	return b, nil
}
func getSolanaExtraArgs(env deployment.Environment, c SetCounterpartPingPongDemoContractsConfig) ([]byte, error) {
	counterpartContractAddressStr, err := ccipChangeset.GetPingPongDemoContractAddress(env, c.CounterpartChainSelector, c.ChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to get ping pong demo counterpart ct address: %w", err)
	}

	pingPongAddress, err := solana.PublicKeyFromBase58(counterpartContractAddressStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ping pong demo contract address: %w", err)
	}

	feesTokenMintStr, err := deployment.SearchAddressBook(env.ExistingAddresses, c.CounterpartChainSelector, c.ExtraArgsSolana.FeesTokenMintType)
	if err != nil {
		return nil, fmt.Errorf("failed to get fees token mint: %w", err)
	}
	feesTokenMint, err := solana.PublicKeyFromBase58(feesTokenMintStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse fees token mint: %w", err)
	}

	linkTokenStr, err := deployment.SearchAddressBook(env.ExistingAddresses, c.CounterpartChainSelector, types.LinkToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get link token: %w", err)
	}
	linkToken, err := solana.PublicKeyFromBase58(linkTokenStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse link token: %w", err)
	}

	pingPongPDAData, err := ccipChangeset.LoadPingPongPDAData(
		env,
		c.CounterpartChainSelector,
		pingPongAddress,
		c.ChainSelector,
		c.ExtraArgsSolana.FeesTokenProgram,
		feesTokenMint,
		linkToken,
	)

	accounts := []solana.PublicKey{
		pingPongPDAData.PPConfigPDA,
		pingPongPDAData.PPSendSignerPDA,
		pingPongPDAData.RouterDestChainStatePDA,
		pingPongPDAData.RouterNoncePDA,
		pingPongPDAData.FeeBillingSignerPDA,
		pingPongPDAData.FqBillingTokenConfigPDA,
		pingPongPDAData.RMNRemoteConfigPDA,
		pingPongPDAData.RMNRemoteCursesPDA,
		pingPongPDAData.PPFeeTokenAta,
		pingPongPDAData.RouterFeeTokenReceiver,
		pingPongPDAData.FqDestChainPDA,
		pingPongPDAData.FqLinkTokenConfigPDA,
	}

	writableBitmap := solanaUtilsCcip.GenerateBitMapForIndexes([]int{1, 4, 7, 8, 9})

	b, err := testhelpers.GetSVMExtraArgsV1(c.ExtraArgsSolana.ComputeUnits, writableBitmap, accounts)

	if err != nil {
		return nil, fmt.Errorf("failed to marshal extra args for solana: %w", err)
	}

	return b, nil
}

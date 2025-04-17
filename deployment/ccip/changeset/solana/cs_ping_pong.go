package solana

import (
	"fmt"

	"github.com/gagliardetto/solana-go"
	chainsel "github.com/smartcontractkit/chain-selectors"
	solPingPong "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ping_pong_demo"
	solanaStateUtils "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	ccipChangeset "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
)

var DeployPingPongContractChangeset = deployment.CreateChangeSet(deployPingPongContractChangeset, validateDeployPingPongContractConfig)
var StartPingPongContractChangeset = deployment.CreateChangeSet(startPingPongContractChangeset, validateStartPingPongContract)
var SetPausePingPongChangeset = deployment.CreateChangeSet(setPausePingPongChangeset, validateSetPausePingPongConfig)
var SetCounterpartPingPongChangeset = deployment.CreateChangeSet(setCounterpartPingPongChangeset, validateSetCounterpartPingPongConfig)

type DeployPingPongContractsConfig struct {
	ChainSelector            uint64
	CounterpartChainSelector uint64

	IsPaused  bool
	ExtraArgs []byte

	FeesTokenProgram  solana.PublicKey
	FeesTokenMintType deployment.ContractType
}

func validateDeployPingPongContractConfig(env deployment.Environment, config DeployPingPongContractsConfig) error {
	return nil
}
func deployPingPongContractChangeset(
	env deployment.Environment,
	c DeployPingPongContractsConfig,
) (deployment.ChangesetOutput, error) {
	s, err := ccipChangeset.LoadOnchainStateSolana(env)
	if err != nil {
		env.Logger.Errorw("Failed to load existing onchain state", "err", err)
		return deployment.ChangesetOutput{}, err
	}

	chain := env.SolChains[c.ChainSelector]
	chainState := s.SolChains[c.ChainSelector]

	newAddresses := deployment.NewMemoryAddressBook()

	programName := getTypeToProgramDeployName()[changeset.PingPongDemo]

	programID, err := chain.DeployProgram(env.Logger, programName, false)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to deploy program: %w", err)
	}

	programAddress := solana.MustPublicKeyFromBase58(programID)

	solPingPong.SetProgramID(programAddress)

	counterpartAddressBytes, err := ccipChangeset.GetPaddedPingPongAddressBytes(env, c.CounterpartChainSelector, c.ChainSelector, chainsel.FamilyEVM)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get counterpart address: %w", err)
	}

	feesTokenMintStr, err := deployment.SearchAddressBook(env.ExistingAddresses, c.ChainSelector, c.FeesTokenMintType)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get fees token mint: %w", err)
	}
	feesTokenMint, err := solana.PublicKeyFromBase58(feesTokenMintStr)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get fees token mint: %w", err)
	}

	linkTokenStr, err := deployment.SearchAddressBook(env.ExistingAddresses, c.ChainSelector, types.LinkToken)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get fees token mint: %w", err)
	}
	linkToken, err := solana.PublicKeyFromBase58(linkTokenStr)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get fees token mint: %w", err)
	}

	tv := deployment.NewTypeAndVersion(changeset.PingPongDemo, deployment.Version1_0_0)
	tv.Labels.Add(fmt.Sprintf("To - %d", c.CounterpartChainSelector))

	err = newAddresses.Save(chain.Selector, programID, tv)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to save address: %w", err)
	}

	pdaData, err := ccipChangeset.LoadPingPongPDAData(
		env,
		c.ChainSelector,
		programAddress,
		c.CounterpartChainSelector,
		c.FeesTokenProgram,
		feesTokenMint,
		linkToken,
	)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	initConfigIx, err := solPingPong.NewInitializeConfigInstruction(
		chainState.Router,
		c.CounterpartChainSelector,
		counterpartAddressBytes,
		c.IsPaused,
		c.ExtraArgs,
		pdaData.PPConfigPDA,
		feesTokenMint,
		chain.DeployerKey.PublicKey(),
		solana.SystemProgramID,
		programAddress,
		pdaData.ProgramData.Address,
	).ValidateAndBuild()

	initializeIx, err := solPingPong.NewInitializeInstruction(
		pdaData.PPConfigPDA,
		pdaData.NameVersionPDA,
		pdaData.FeeBillingSignerPDA,
		c.FeesTokenProgram,
		feesTokenMint,
		pdaData.PPFeeTokenAta,
		pdaData.PPSendSignerPDA,
		chain.DeployerKey.PublicKey(),
		solana.SPLAssociatedTokenAccountProgramID,
		solana.SystemProgramID,
	).ValidateAndBuild()

	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to build instruction: %w", err)
	}

	if err := chain.Confirm([]solana.Instruction{initConfigIx, initializeIx}); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to confirm initialize instruction: %w", err)
	}

	env.Logger.Infow("Initialized ping pong demo", "chain", chain.String())

	return deployment.ChangesetOutput{
		AddressBook: newAddresses,
	}, nil
}

type StartPingPongConfig struct {
	ChainSelector            uint64
	CounterpartChainSelector uint64

	FeesTokenProgram  solana.PublicKey
	FeesTokenMintType deployment.ContractType
}

func validateStartPingPongContract(e deployment.Environment, config StartPingPongConfig) error {
	return nil
}
func startPingPongContractChangeset(
	env deployment.Environment,
	c StartPingPongConfig,
) (deployment.ChangesetOutput, error) {
	s, err := ccipChangeset.LoadOnchainStateSolana(env)
	if err != nil {
		env.Logger.Errorw("Failed to load existing onchain state", "err", err)
		return deployment.ChangesetOutput{}, err
	}

	chain := env.SolChains[c.ChainSelector]
	chainState := s.SolChains[c.ChainSelector]

	contractAddressStr, err := ccipChangeset.GetPingPongDemoContractAddress(env, c.ChainSelector, c.CounterpartChainSelector)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get contract address: %w", err)
	}
	contractAddress, err := solana.PublicKeyFromBase58(contractAddressStr)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to parse contract address: %w", err)
	}

	solPingPong.SetProgramID(contractAddress)

	feesTokenMintStr, err := deployment.SearchAddressBook(env.ExistingAddresses, c.ChainSelector, c.FeesTokenMintType)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get fees token mint: %w", err)
	}
	feesTokenMint, err := solana.PublicKeyFromBase58(feesTokenMintStr)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get fees token mint: %w", err)
	}

	linkTokenStr, err := deployment.SearchAddressBook(env.ExistingAddresses, c.ChainSelector, types.LinkToken)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get fees token mint: %w", err)
	}
	linkToken, err := solana.PublicKeyFromBase58(linkTokenStr)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get fees token mint: %w", err)
	}

	pdaData, err := ccipChangeset.LoadPingPongPDAData(
		env,
		c.ChainSelector,
		contractAddress,
		c.CounterpartChainSelector,
		c.FeesTokenProgram,
		feesTokenMint,
		linkToken,
	)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	initTx, err := solPingPong.NewStartPingPongInstruction(
		pdaData.PPConfigPDA,
		chain.DeployerKey.PublicKey(),
		pdaData.PPSendSignerPDA,
		c.FeesTokenProgram,
		feesTokenMint,
		pdaData.PPFeeTokenAta,

		chainState.Router,
		chainState.RouterConfigPDA,
		pdaData.RouterDestChainStatePDA,
		pdaData.RouterNoncePDA,
		pdaData.RouterFeeTokenReceiver,
		pdaData.FeeBillingSignerPDA,

		chainState.FeeQuoter,
		chainState.FeeQuoterConfigPDA,
		pdaData.FqDestChainPDA,
		pdaData.FqBillingTokenConfigPDA,
		pdaData.FqLinkTokenConfigPDA,

		chainState.RMNRemote,
		pdaData.RMNRemoteCursesPDA,
		pdaData.RMNRemoteConfigPDA,
		solana.SystemProgramID,
	).ValidateAndBuild()

	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to build instruction: %w", err)
	}

	if err := chain.Confirm([]solana.Instruction{initTx}); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to confirm start instruction: %w", err)
	}

	env.Logger.Infow("Initialized ping pong demo", "chain", chain.String())

	return deployment.ChangesetOutput{}, nil
}

type SetCounterpartPingPongConfig struct {
	ChainSelector            uint64
	CounterpartChainSelector uint64
}

func validateSetCounterpartPingPongConfig(e deployment.Environment, config SetCounterpartPingPongConfig) error {
	return nil
}
func setCounterpartPingPongChangeset(
	env deployment.Environment,
	c SetCounterpartPingPongConfig,
) (deployment.ChangesetOutput, error) {
	chain := env.SolChains[c.ChainSelector]

	contractAddressStr, err := ccipChangeset.GetPingPongDemoContractAddress(env, c.ChainSelector, c.CounterpartChainSelector)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get contract address: %w", err)
	}
	contractAddress, err := solana.PublicKeyFromBase58(contractAddressStr)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to parse contract address: %w", err)
	}

	ppConfigPDA, _, err := solanaStateUtils.FindPingPongDemoConfigPDA(contractAddress)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to find ping pong config pda: %w", err)
	}

	solPingPong.SetProgramID(contractAddress)

	counterpartAddressBytes, err := ccipChangeset.GetPaddedPingPongAddressBytes(env, c.CounterpartChainSelector, c.ChainSelector, chainsel.FamilyEVM)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get counterpart address: %w", err)
	}

	initTx, err := solPingPong.NewSetCounterpartInstruction(
		c.CounterpartChainSelector,
		counterpartAddressBytes,
		ppConfigPDA,
		chain.DeployerKey.PublicKey(),
	).ValidateAndBuild()

	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to build instruction: %w", err)
	}

	if err := chain.Confirm([]solana.Instruction{initTx}); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to confirm instruction: %w", err)
	}

	env.Logger.Infow("Initialized ping pong demo", "chain", chain.String())

	return deployment.ChangesetOutput{}, nil
}

type SetPausePingPongConfig struct {
	ChainSelector            uint64
	CounterpartChainSelector uint64

	IsPaused bool
}

func validateSetPausePingPongConfig(e deployment.Environment, config SetPausePingPongConfig) error {
	return nil
}
func setPausePingPongChangeset(
	env deployment.Environment,
	c SetPausePingPongConfig,
) (deployment.ChangesetOutput, error) {
	chain := env.SolChains[c.ChainSelector]

	contractAddressStr, err := ccipChangeset.GetPingPongDemoContractAddress(env, c.ChainSelector, c.CounterpartChainSelector)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get contract address: %w", err)
	}
	contractAddress, err := solana.PublicKeyFromBase58(contractAddressStr)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to parse contract address: %w", err)
	}

	ppConfigPDA, _, err := solanaStateUtils.FindPingPongDemoConfigPDA(contractAddress)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to find ping pong config pda: %w", err)
	}

	solPingPong.SetProgramID(contractAddress)

	initTx, err := solPingPong.NewSetPausedInstruction(
		c.IsPaused,
		ppConfigPDA,
		chain.DeployerKey.PublicKey(),
	).ValidateAndBuild()

	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to build instruction: %w", err)
	}

	if err := chain.Confirm([]solana.Instruction{initTx}); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to confirm instruction: %w", err)
	}

	env.Logger.Infow("Initialized ping pong demo", "chain", chain.String())

	return deployment.ChangesetOutput{}, nil
}

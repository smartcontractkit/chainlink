package operation

import (
	"errors"
	"fmt"
	"math/big"
	"math/rand"

	binary "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/rpc"
	mcmsSolanaSdk "github.com/smartcontractkit/mcms/sdk/solana"
	mcmsTypes "github.com/smartcontractkit/mcms/types"
	"github.com/smartcontractkit/wsrpc/logger"

	accessControllerBindings "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_1/access_controller"
	mcmBindings "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_1/mcm"
	timelockBindings "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_1/timelock"
	solanaUtils "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	cldfsol "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment"
	commonOps "github.com/smartcontractkit/chainlink/deployment/common/changeset/solana/operations"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
)

type Deps struct {
	Env       cldf.Environment
	State     *state.MCMSWithTimelockStateSolana
	Datastore datastore.MutableDataStore
	Chain     cldfsol.Chain
}

var (
	DeployAccessControllerOp = operations.NewOperation(
		"deploy-access-controller",
		&deployment.Version1_0_0,
		"Deploys access controller for solana",
		commonOps.Deploy,
	)

	DeployMCMProgramOp = operations.NewOperation(
		"deploy-mcm-program",
		&deployment.Version1_0_0,
		"Deploys mcm for solana",
		commonOps.Deploy,
	)

	DeployTimelockOp = operations.NewOperation(
		"deploy-timelock-program",
		&deployment.Version1_0_0,
		"Deploys timelock for solana",
		commonOps.Deploy,
	)

	InitAccessControllerOp = operations.NewOperation(
		"init-access-controller",
		&deployment.Version1_0_0,
		"Initializes access controller for solana",
		initAccessController,
	)

	InitMCMOp = operations.NewOperation(
		"init-mcm-program",
		&deployment.Version1_0_0,
		"Initializes MCMProgram for solana",
		initMCM,
	)

	InitTimelockOp = operations.NewOperation(
		"init-timelock-program",
		&deployment.Version1_0_0,
		"Initializes timelock for solana",
		initTimelock,
	)
	AddAccessOp = operations.NewOperation(
		"add-access-op",
		&deployment.Version1_0_0,
		"Adds access to provided role for timelock",
		addAccess,
	)
)

type (
	InitAccessControllerInput struct {
		ContractType cldf.ContractType
		ChainSel     uint64
	}

	InitAccessControllerOutput struct{}

	InitMCMInput struct {
		ContractType cldf.ContractType
		MCMConfig    mcmsTypes.Config
		ChainSel     uint64
	}

	InitMCMOutput struct{}

	InitTimelockInput struct {
		ContractType cldf.ContractType
		ChainSel     uint64
		MinDelay     *big.Int
	}
	InitTimelockOutput struct{}

	AddAccessInput struct {
		Role     timelockBindings.Role
		Accounts []solana.PublicKey
		ChainSel uint64
	}

	AddAccessOutput struct{}
)

func initAccessController(b operations.Bundle, deps Deps, in InitAccessControllerInput) (InitAccessControllerOutput, error) {
	var out InitAccessControllerOutput

	if deps.State.AccessControllerProgram.IsZero() {
		return out, fmt.Errorf("access controller program is not deployed for chain sel %d", deps.Chain.ChainSelector())
	}

	typeAndVersion := cldf.NewTypeAndVersion(in.ContractType, deployment.Version1_0_0)
	_, accessControllerAccountSeed, err := deps.State.GetStateFromType(in.ContractType)
	if err != nil {
		return out, fmt.Errorf("failed to get account controller state: %w", err)
	}

	accessControllerAccount := solana.PublicKeyFromBytes(accessControllerAccountSeed[:])
	if !accessControllerAccount.IsZero() {
		var data accessControllerBindings.AccessController
		err = solanaUtils.GetAccountDataBorshInto(b.GetContext(), deps.Chain.Client, accessControllerAccount, rpc.CommitmentConfirmed, &data)
		if err == nil {
			b.Logger.Infow("access controller already initialized, skipping initialization", "chain", deps.Chain.String())
			return out, nil
		}

		return out, fmt.Errorf("unable to read access controller account config %s", accessControllerAccount.String())
	}

	b.Logger.Infow("access controller not initialized, initializing", "chain", deps.Chain.String())

	programID := deps.State.AccessControllerProgram
	if accessControllerBindings.ProgramID.IsZero() {
		accessControllerBindings.SetProgramID(programID)
	}

	log := logger.With(b.Logger, "chain", deps.Chain.String(), "contract", typeAndVersion.String(), "programID", programID)

	account, err := solana.NewRandomPrivateKey()
	if err != nil {
		return out, fmt.Errorf("failed to generate new random private key for access controller account: %w", err)
	}

	err = initializeAccessController(b, deps.Chain, programID, account)
	if err != nil {
		return out, fmt.Errorf("failed to initialize access controller: %w", err)
	}

	log.Infow("initialized access controller", "account", account.PublicKey())

	err = deps.State.SetState(in.ContractType, account.PublicKey(), state.PDASeed{})
	if err != nil {
		return out, fmt.Errorf("failed to save onchain state: %w", err)
	}

	err = deps.Datastore.Addresses().Add(datastore.AddressRef{
		Address:       account.PublicKey().String(),
		ChainSelector: deps.Chain.Selector,
		Type:          datastore.ContractType(in.ContractType),
	})
	if err != nil {
		return out, fmt.Errorf("failed to save address to datastore: %w", err)
	}

	return out, nil
}

// discriminator + owner + proposed owner + access_list (64 max addresses + length)
const accessControllerAccountSize = uint64(8 + 32 + 32 + ((32 * 64) + 8))

func initializeAccessController(
	b operations.Bundle, chain cldfsol.Chain, programID solana.PublicKey, roleAccount solana.PrivateKey,
) error {
	rentExemption, err := chain.Client.GetMinimumBalanceForRentExemption(b.GetContext(),
		accessControllerAccountSize, rpc.CommitmentConfirmed)
	if err != nil {
		return fmt.Errorf("failed to get minimum balance for rent exemption: %w", err)
	}

	createAccountInstruction, err := system.NewCreateAccountInstruction(rentExemption, accessControllerAccountSize,
		programID, chain.DeployerKey.PublicKey(), roleAccount.PublicKey()).ValidateAndBuild()
	if err != nil {
		return fmt.Errorf("failed to create CreateAccount instruction: %w", err)
	}

	initializeInstruction, err := accessControllerBindings.NewInitializeInstruction(
		roleAccount.PublicKey(),
		chain.DeployerKey.PublicKey(),
	).ValidateAndBuild()
	if err != nil {
		return fmt.Errorf("failed to build instruction: %w", err)
	}

	instructions := []solana.Instruction{createAccountInstruction, initializeInstruction}
	err = chain.Confirm(instructions, solanaUtils.AddSigners(roleAccount))
	if err != nil {
		return fmt.Errorf("failed to confirm CreateAccount and InitializeAccessController instructions: %w", err)
	}

	var data accessControllerBindings.AccessController
	err = solanaUtils.GetAccountDataBorshInto(b.GetContext(), chain.Client, roleAccount.PublicKey(), rpc.CommitmentConfirmed, &data)
	if err != nil {
		return fmt.Errorf("failed to read access controller roleAccount: %w", err)
	}

	return nil
}

func initMCM(b operations.Bundle, deps Deps, in InitMCMInput) (InitMCMOutput, error) {
	var out InitMCMOutput

	if deps.State.McmProgram.IsZero() {
		return out, fmt.Errorf("mcm program is not deployed for chain sel %d", deps.Chain.ChainSelector())
	}

	programID := deps.State.McmProgram
	mcmBindings.SetProgramID(programID)

	typeAndVersion := cldf.NewTypeAndVersion(in.ContractType, deployment.Version1_0_0)
	mcmProgram, mcmSeed, err := deps.State.GetStateFromType(in.ContractType)
	if err != nil {
		return out, fmt.Errorf("failed to get mcm state: %w", err)
	}

	if mcmSeed != (state.PDASeed{}) {
		mcmConfigPDA := state.GetMCMConfigPDA(mcmProgram, mcmSeed)
		var data mcmBindings.MultisigConfig
		err = solanaUtils.GetAccountDataBorshInto(b.GetContext(), deps.Chain.Client, mcmConfigPDA, rpc.CommitmentConfirmed, &data)
		if err == nil {
			b.Logger.Infow("mcm config already initialized, skipping initialization", "chain", deps.Chain.String())
			return out, nil
		}
		return out, fmt.Errorf("unable to read mcm ConfigPDA account config %q", mcmConfigPDA.String())
	}

	b.Logger.Infow("mcm config not initialized, initializing", "chain", deps.Chain.String())
	log := logger.With(b.Logger, "chain", deps.Chain.String(), "contract", typeAndVersion.String())

	seed := randomSeed()
	log.Infow("generated MCM seed", "seed", string(seed[:]))
	err = initializeMCM(b, deps, programID, seed)
	if err != nil {
		return out, fmt.Errorf("failed to initialize mcm: %w", err)
	}

	mcmAddress := state.EncodeAddressWithSeed(programID, seed)

	configurer := mcmsSolanaSdk.NewConfigurer(deps.Chain.Client, *deps.Chain.DeployerKey, mcmsTypes.ChainSelector(deps.Chain.ChainSelector()))
	tx, err := configurer.SetConfig(b.GetContext(), mcmAddress, &in.MCMConfig, false)
	if err != nil {
		return out, fmt.Errorf("failed to set config on mcm: %w", err)
	}
	log.Infow("called SetConfig on MCM", "transaction", tx.Hash)

	err = deps.Datastore.Addresses().Add(datastore.AddressRef{
		Address:       mcmAddress,
		ChainSelector: deps.Chain.Selector,
		Type:          datastore.ContractType(in.ContractType),
	})
	if err != nil {
		return out, fmt.Errorf("failed to save address to datastore: %w", err)
	}

	err = deps.State.SetState(in.ContractType, programID, seed)
	if err != nil {
		return out, fmt.Errorf("failed to save onchain state: %w", err)
	}

	return out, nil
}

func initializeMCM(b operations.Bundle, deps Deps, mcmProgram solana.PublicKey, multisigID state.PDASeed) error {
	var mcmConfig mcmBindings.MultisigConfig
	err := deps.Chain.GetAccountDataBorshInto(b.GetContext(), state.GetMCMConfigPDA(mcmProgram, multisigID), &mcmConfig)
	if err == nil {
		b.Logger.Infow("MCM already initialized, skipping initialization", "chain", deps.Chain.String())
		return nil
	}

	var programData struct {
		DataType uint32
		Address  solana.PublicKey
	}
	opts := &rpc.GetAccountInfoOpts{Commitment: rpc.CommitmentConfirmed}

	data, err := deps.Chain.Client.GetAccountInfoWithOpts(b.GetContext(), mcmProgram, opts)
	if err != nil {
		return fmt.Errorf("failed to get mcm program account info: %w", err)
	}
	err = binary.UnmarshalBorsh(&programData, data.Bytes())
	if err != nil {
		return fmt.Errorf("failed to unmarshal program data: %w", err)
	}

	instruction, err := mcmBindings.NewInitializeInstruction(
		deps.Chain.Selector,
		multisigID,
		state.GetMCMConfigPDA(mcmProgram, multisigID),
		deps.Chain.DeployerKey.PublicKey(),
		solana.SystemProgramID,
		mcmProgram,
		programData.Address,
		state.GetMCMRootMetadataPDA(mcmProgram, multisigID),
		state.GetMCMExpiringRootAndOpCountPDA(mcmProgram, multisigID),
	).ValidateAndBuild()
	if err != nil {
		return fmt.Errorf("failed to build instruction: %w", err)
	}

	err = deps.Chain.Confirm([]solana.Instruction{instruction})
	if err != nil {
		return fmt.Errorf("failed to confirm instructions: %w", err)
	}

	return nil
}

func initTimelock(b operations.Bundle, deps Deps, in InitTimelockInput) (InitTimelockOutput, error) {
	var out InitTimelockOutput

	if deps.State.TimelockProgram.IsZero() {
		return out, errors.New("mcm program is not deployed")
	}
	programID := deps.State.TimelockProgram
	timelockBindings.SetProgramID(programID)

	typeAndVersion := cldf.NewTypeAndVersion(in.ContractType, deployment.Version1_0_0)
	timelockProgram, timelockSeed, err := deps.State.GetStateFromType(in.ContractType)
	if err != nil {
		return out, fmt.Errorf("failed to get timelock state: %w", err)
	}

	if (timelockSeed != state.PDASeed{}) {
		timelockConfigPDA := state.GetTimelockConfigPDA(timelockProgram, timelockSeed)
		var timelockConfig timelockBindings.Config
		err = deps.Chain.GetAccountDataBorshInto(b.GetContext(), timelockConfigPDA, &timelockConfig)
		if err == nil {
			b.Logger.Infow("timelock config already initialized, skipping initialization", "chain", deps.Chain.String())
			return out, nil
		}
		return out, fmt.Errorf("unable to read timelock ConfigPDA account config %s", timelockConfigPDA.String())
	}

	b.Logger.Infow("timelock config not initialized, initializing", "chain", deps.Chain.String())
	log := logger.With(b.Logger, "chain", deps.Chain.String(), "contract", typeAndVersion.String())

	seed := randomSeed()
	log.Infow("generated Timelock seed", "seed", string(seed[:]))

	err = initializeTimelock(b, deps, programID, seed, in.MinDelay)
	if err != nil {
		return out, fmt.Errorf("failed to initialize timelock: %w", err)
	}

	timelockAddress := state.EncodeAddressWithSeed(programID, seed)

	err = deps.Datastore.Addresses().Add(datastore.AddressRef{
		Address:       timelockAddress,
		ChainSelector: deps.Chain.Selector,
		Type:          datastore.ContractType(in.ContractType),
	})
	if err != nil {
		return out, fmt.Errorf("failed to save address to datastore: %w", err)
	}

	err = deps.State.SetState(in.ContractType, programID, seed)
	if err != nil {
		return out, fmt.Errorf("failed to save onchain state: %w", err)
	}

	return out, nil
}

func initializeTimelock(b operations.Bundle, deps Deps, timelockProgram solana.PublicKey,
	timelockID state.PDASeed, minDelay *big.Int) error {
	if minDelay == nil {
		minDelay = big.NewInt(0)
	}

	var timelockConfig timelockBindings.Config
	err := deps.Chain.GetAccountDataBorshInto(b.GetContext(), state.GetTimelockConfigPDA(timelockProgram, timelockID),
		&timelockConfig)
	if err == nil {
		b.Logger.Infow("Timelock already initialized, skipping initialization", "chain", deps.Chain.String())
		return nil
	}

	var programData struct {
		DataType uint32
		Address  solana.PublicKey
	}
	opts := &rpc.GetAccountInfoOpts{Commitment: rpc.CommitmentConfirmed}

	data, err := deps.Chain.Client.GetAccountInfoWithOpts(b.GetContext(), timelockProgram, opts)
	if err != nil {
		return fmt.Errorf("failed to get timelock program account info: %w", err)
	}
	err = binary.UnmarshalBorsh(&programData, data.Bytes())
	if err != nil {
		return fmt.Errorf("failed to unmarshal program data: %w", err)
	}

	instruction, err := timelockBindings.NewInitializeInstruction(
		timelockID,
		minDelay.Uint64(),
		state.GetTimelockConfigPDA(timelockProgram, timelockID),
		deps.Chain.DeployerKey.PublicKey(),
		solana.SystemProgramID,
		timelockProgram,
		programData.Address,
		deps.State.AccessControllerProgram,
		deps.State.ProposerAccessControllerAccount,
		deps.State.ExecutorAccessControllerAccount,
		deps.State.CancellerAccessControllerAccount,
		deps.State.BypasserAccessControllerAccount,
	).ValidateAndBuild()
	if err != nil {
		return fmt.Errorf("failed to build instruction: %w", err)
	}

	err = deps.Chain.Confirm([]solana.Instruction{instruction})
	if err != nil {
		return fmt.Errorf("failed to confirm instructions: %w", err)
	}

	return nil
}

func addAccess(b operations.Bundle, deps Deps, in AddAccessInput) (AddAccessOutput, error) {
	var out AddAccessOutput

	timelockConfigPDA := state.GetTimelockConfigPDA(deps.State.TimelockProgram, deps.State.TimelockSeed)

	instructionBuilder := timelockBindings.NewBatchAddAccessInstruction([32]uint8(deps.State.TimelockSeed), in.Role,
		timelockConfigPDA, deps.State.AccessControllerProgram, deps.State.RoleAccount(in.Role), deps.Chain.DeployerKey.PublicKey())

	for _, account := range in.Accounts {
		instructionBuilder.Append(solana.Meta(account))
	}

	instruction, err := instructionBuilder.ValidateAndBuild()
	if err != nil {
		return out, fmt.Errorf("failed to build BatchAddAccess instruction: %w", err)
	}

	err = deps.Chain.Confirm([]solana.Instruction{instruction})
	if err != nil {
		return out, fmt.Errorf("failed to confirm BatchAddAccess instruction: %w", err)
	}
	return out, nil
}

func randomSeed() state.PDASeed {
	const alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	seed := state.PDASeed{}
	for i := range seed {
		seed[i] = alphabet[rand.Intn(len(alphabet))]
	}

	return seed
}

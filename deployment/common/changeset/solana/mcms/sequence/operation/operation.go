package operation

import (
	"fmt"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/rpc"
	accessControllerBindings "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/access_controller"
	solanaUtils "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	cldfsol "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment"
	commonOps "github.com/smartcontractkit/chainlink/deployment/common/changeset/solana/operations"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/wsrpc/logger"
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

	InitAccessControllerOp = operations.NewOperation(
		"deploy-access-controller",
		&deployment.Version1_0_0,
		"Deploys access controller for solana",
		initAccessController,
	)
)

type (
	InitAccessControllerInput struct {
		ContractType cldf.ContractType
	}

	InitAccessControllerOutput struct{}
)

func initAccessController(b operations.Bundle, deps Deps, in InitAccessControllerInput) (InitAccessControllerOutput, error) {
	var out InitAccessControllerOutput

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
	accessControllerBindings.SetProgramID(programID)

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
		Address:       account.String(),
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

package mcmsnew

import (
	"errors"
	"fmt"
	"math/big"

	binary "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"

	timelockBindings "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/timelock"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
)

func deployTimelockProgram(
	e deployment.Environment, chainState *state.MCMSWithTimelockStateSolana, chain deployment.SolChain,
	addressBook deployment.AddressBook,
) error {
	typeAndVersion := deployment.NewTypeAndVersion(commontypes.RBACTimelockProgram, deployment.Version1_0_0)
	log := logger.With(e.Logger, "chain", chain.String(), "contract", typeAndVersion.String())

	programID, _, err := chainState.GetStateFromType(commontypes.RBACTimelock)
	if err != nil {
		return fmt.Errorf("failed to get timelock state: %w", err)
	}

	if programID.IsZero() {
		deployedProgramID, err := chain.DeployProgram(log, "timelock", false)
		if err != nil {
			return fmt.Errorf("failed to deploy timelock program: %w", err)
		}

		programID, err = solana.PublicKeyFromBase58(deployedProgramID)
		if err != nil {
			return fmt.Errorf("failed to convert timelock program id to public key: %w", err)
		}

		err = addressBook.Save(chain.Selector, programID.String(), typeAndVersion)
		if err != nil {
			return fmt.Errorf("failed to save mcm address: %w", err)
		}

		err = chainState.SetState(commontypes.RBACTimelockProgram, programID, state.PDASeed{})
		if err != nil {
			return fmt.Errorf("failed to save onchain state: %w", err)
		}

		log.Infow("deployed timelock contract", "programId", programID)
	} else {
		log.Infow("using existing Timelock program", "programId", programID.String())
	}

	return nil
}

func initTimelock(
	e deployment.Environment, chainState *state.MCMSWithTimelockStateSolana, chain deployment.SolChain,
	addressBook deployment.AddressBook, minDelay *big.Int,
) error {
	if chainState.TimelockProgram.IsZero() {
		return errors.New("mcm program is not deployed")
	}
	programID := chainState.TimelockProgram
	timelockBindings.SetProgramID(programID)

	typeAndVersion := deployment.NewTypeAndVersion(commontypes.RBACTimelock, deployment.Version1_0_0)
	timelockProgram, timelockSeed, err := chainState.GetStateFromType(commontypes.RBACTimelock)
	if err != nil {
		return fmt.Errorf("failed to get timelock state: %w", err)
	}

	if (timelockSeed != state.PDASeed{}) {
		timelockConfigPDA := state.GetTimelockConfigPDA(timelockProgram, timelockSeed)
		var timelockConfig timelockBindings.Config
		err = chain.GetAccountDataBorshInto(e.GetContext(), timelockConfigPDA, &timelockConfig)
		if err == nil {
			e.Logger.Infow("timelock config already initialized, skipping initialization", "chain", chain.String())
			return nil
		}
		return fmt.Errorf("unable to read timelock ConfigPDA account config %s", timelockConfigPDA.String())
	}

	e.Logger.Infow("timelock config not initialized, initializing", "chain", chain.String())
	log := logger.With(e.Logger, "chain", chain.String(), "contract", typeAndVersion.String())

	seed := randomSeed()
	log.Infow("generated Timelock seed", "seed", string(seed[:]))

	err = initializeTimelock(e, chain, programID, seed, chainState, minDelay)
	if err != nil {
		return fmt.Errorf("failed to initialize timelock: %w", err)
	}

	timelockAddress := state.EncodeAddressWithSeed(programID, seed)

	err = addressBook.Save(chain.Selector, timelockAddress, typeAndVersion)
	if err != nil {
		return fmt.Errorf("failed to save address: %w", err)
	}

	err = chainState.SetState(commontypes.RBACTimelock, programID, seed)
	if err != nil {
		return fmt.Errorf("failed to save onchain state: %w", err)
	}

	return nil
}

func initializeTimelock(
	e deployment.Environment, chain deployment.SolChain, timelockProgram solana.PublicKey, timelockID state.PDASeed,
	chainState *state.MCMSWithTimelockStateSolana, minDelay *big.Int,
) error {
	if minDelay == nil {
		minDelay = big.NewInt(0)
	}

	var timelockConfig timelockBindings.Config
	err := chain.GetAccountDataBorshInto(e.GetContext(), state.GetTimelockConfigPDA(timelockProgram, timelockID),
		&timelockConfig)
	if err == nil {
		e.Logger.Infow("Timelock already initialized, skipping initialization", "chain", chain.String())
		return nil
	}

	var programData struct {
		DataType uint32
		Address  solana.PublicKey
	}
	opts := &rpc.GetAccountInfoOpts{Commitment: rpc.CommitmentConfirmed}

	data, err := chain.Client.GetAccountInfoWithOpts(e.GetContext(), timelockProgram, opts)
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
		chain.DeployerKey.PublicKey(),
		solana.SystemProgramID,
		timelockProgram,
		programData.Address,
		chainState.AccessControllerProgram,
		chainState.ProposerAccessControllerAccount,
		chainState.ExecutorAccessControllerAccount,
		chainState.CancellerAccessControllerAccount,
		chainState.BypasserAccessControllerAccount,
	).ValidateAndBuild()
	if err != nil {
		return fmt.Errorf("failed to build instruction: %w", err)
	}

	err = chain.Confirm([]solana.Instruction{instruction})
	if err != nil {
		return fmt.Errorf("failed to confirm instructions: %w", err)
	}

	return nil
}

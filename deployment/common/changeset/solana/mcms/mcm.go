package mcmsnew

import (
	"errors"
	"fmt"
	"math/rand"

	binary "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	mcmsSolanaSdk "github.com/smartcontractkit/mcms/sdk/solana"
	mcmsTypes "github.com/smartcontractkit/mcms/types"

	mcmBindings "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/mcm"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	solanaUtils "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
)

func deployMCMProgram(
	env deployment.Environment, chainState *state.MCMSWithTimelockStateSolana,
	chain deployment.SolChain, addressBook deployment.AddressBook,
) error {
	typeAndVersion := deployment.NewTypeAndVersion(commontypes.ManyChainMultisigProgram, deployment.Version1_0_0)
	log := logger.With(env.Logger, "chain", chain.String(), "contract", typeAndVersion.String())

	programID, _, err := chainState.GetStateFromType(commontypes.ManyChainMultisigProgram)
	if err != nil {
		return fmt.Errorf("failed to get mcm state: %w", err)
	}

	if programID.IsZero() {
		deployedProgramID, err := chain.DeployProgram(log, "mcm", false)
		if err != nil {
			return fmt.Errorf("failed to deploy mcm program: %w", err)
		}

		programID, err = solana.PublicKeyFromBase58(deployedProgramID)
		if err != nil {
			return fmt.Errorf("failed to convert mcm program id to public key: %w", err)
		}

		err = addressBook.Save(chain.Selector, programID.String(), typeAndVersion)
		if err != nil {
			return fmt.Errorf("failed to save mcm address: %w", err)
		}

		err = chainState.SetState(commontypes.ManyChainMultisigProgram, programID, state.PDASeed{})
		if err != nil {
			return fmt.Errorf("failed to save onchain state: %w", err)
		}

		log.Infow("deployed mcm contract", "programId", deployedProgramID)
	} else {
		log.Infow("using existing MCM program", "programId", programID.String())
	}

	return nil
}

func initMCM(
	env deployment.Environment, chainState *state.MCMSWithTimelockStateSolana, contractType deployment.ContractType,
	chain deployment.SolChain, addressBook deployment.AddressBook, mcmConfig *mcmsTypes.Config,
) error {
	if chainState.McmProgram.IsZero() {
		return errors.New("mcm program is not deployed")
	}
	programID := chainState.McmProgram
	mcmBindings.SetProgramID(programID)

	typeAndVersion := deployment.NewTypeAndVersion(contractType, deployment.Version1_0_0)
	mcmProgram, mcmSeed, err := chainState.GetStateFromType(contractType)
	if err != nil {
		return fmt.Errorf("failed to get mcm state: %w", err)
	}

	if mcmSeed != (state.PDASeed{}) {
		mcmConfigPDA := state.GetMCMConfigPDA(mcmProgram, mcmSeed)
		var data mcmBindings.MultisigConfig
		err = solanaUtils.GetAccountDataBorshInto(env.GetContext(), chain.Client, mcmConfigPDA, rpc.CommitmentConfirmed, &data)
		if err == nil {
			env.Logger.Infow("mcm config already initialized, skipping initialization", "chain", chain.String())
			return nil
		}
		return fmt.Errorf("unable to read mcm ConfigPDA account config %s", mcmConfigPDA.String())
	}

	env.Logger.Infow("mcm config not initialized, initializing", "chain", chain.String())
	log := logger.With(env.Logger, "chain", chain.String(), "contract", typeAndVersion.String())

	seed := randomSeed()
	log.Infow("generated MCM seed", "seed", string(seed[:]))

	err = initializeMCM(env, chain, programID, seed)
	if err != nil {
		return fmt.Errorf("failed to initialize mcm: %w", err)
	}

	mcmAddress := state.EncodeAddressWithSeed(programID, seed)

	configurer := mcmsSolanaSdk.NewConfigurer(chain.Client, *chain.DeployerKey, mcmsTypes.ChainSelector(chain.Selector))
	tx, err := configurer.SetConfig(env.GetContext(), mcmAddress, mcmConfig, false)
	if err != nil {
		return fmt.Errorf("failed to set config on mcm: %w", err)
	}
	log.Infow("called SetConfig on MCM", "transaction", tx.Hash)

	err = addressBook.Save(chain.Selector, mcmAddress, typeAndVersion)
	if err != nil {
		return fmt.Errorf("failed to save address: %w", err)
	}

	err = chainState.SetState(contractType, programID, seed)
	if err != nil {
		return fmt.Errorf("failed to save onchain state: %w", err)
	}

	return nil
}

func initializeMCM(e deployment.Environment, chain deployment.SolChain, mcmProgram solana.PublicKey, multisigID state.PDASeed) error {
	var mcmConfig mcmBindings.MultisigConfig
	err := chain.GetAccountDataBorshInto(e.GetContext(), state.GetMCMConfigPDA(mcmProgram, multisigID), &mcmConfig)
	if err == nil {
		e.Logger.Infow("MCM already initialized, skipping initialization", "chain", chain.String())
		return nil
	}

	var programData struct {
		DataType uint32
		Address  solana.PublicKey
	}
	opts := &rpc.GetAccountInfoOpts{Commitment: rpc.CommitmentConfirmed}

	data, err := chain.Client.GetAccountInfoWithOpts(e.GetContext(), mcmProgram, opts)
	if err != nil {
		return fmt.Errorf("failed to get mcm program account info: %w", err)
	}
	err = binary.UnmarshalBorsh(&programData, data.Bytes())
	if err != nil {
		return fmt.Errorf("failed to unmarshal program data: %w", err)
	}

	instruction, err := mcmBindings.NewInitializeInstruction(
		chain.Selector,
		multisigID,
		state.GetMCMConfigPDA(mcmProgram, multisigID),
		chain.DeployerKey.PublicKey(),
		solana.SystemProgramID,
		mcmProgram,
		programData.Address,
		state.GetMCMRootMetadataPDA(mcmProgram, multisigID),
		state.GetMCMExpiringRootAndOpCountPDA(mcmProgram, multisigID),
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

func randomSeed() state.PDASeed {
	const alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	seed := state.PDASeed{}
	for i := range seed {
		seed[i] = alphabet[rand.Intn(len(alphabet))]
	}

	return seed
}

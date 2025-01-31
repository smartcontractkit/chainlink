package mcmsnew

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	mcmBindings "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/mcm"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	mcmsTypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink/deployment"
	mcmsNewTypes "github.com/smartcontractkit/chainlink/deployment/common/changeset/mcmsnew/types"
	proposalUtilsSol "github.com/smartcontractkit/chainlink/deployment/common/proposalutils/solana"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
)

const MCMSBinaryName = "mcms"

// MCMSWithTimelockSolanaDeploy holds a bundle of MCMS contract deploys.
type MCMSWithTimelockSolanaDeploy struct {
	McmProgram           *solana.PublicKey
	CancellerSeed        [32]byte
	BypasserSeed         [32]byte
	ProposerSeed         [32]byte
	Timelock             *solana.PublicKey
	TimelockInstanceSeed [32]byte
	// TODO: not sure if this is needed
	CallProxy *solana.PublicKey
}

func deployMCMSWithConfigSolana(
	contractType deployment.ContractType,
	version semver.Version,
	lggr logger.Logger,
	chain deployment.SolChain,
	ab deployment.AddressBook,
	mcmConfig mcmsTypes.Config,
) (solana.PublicKey, PDASeed, error) {
	var mcmProgram solana.PublicKey
	var mcmSeed [32]byte
	addresses, err := ab.AddressesForChain(chain.Selector)
	if err != nil {
		return solana.PublicKey{}, PDASeed{}, err
	}
	// We purposefully don't check errors here to allow for the state to have Zero() values on
	// programs that are not deployed.
	state, _ := proposalUtilsSol.MaybeLoadMCMSSolanaWithTimelockContracts(chain, addresses)
	mcmState := state.GetStateFromType(contractType)
	if mcmState.IsZero() {
		programID, err := chain.DeployProgram(lggr, MCMSBinaryName)
		if err != nil {
			return solana.PublicKey{}, PDASeed{}, fmt.Errorf("unable to deploy MCMS program: %w", err)
		}

		tv := deployment.NewTypeAndVersion(contractType, version)
		lggr.Infow("Deployed contract", "Contract", tv.String(), "addr", programID, "chain", chain.String())

		mcmProgramID := solana.MustPublicKeyFromBase58(programID)
		err = ab.Save(chain.Selector, mcmProgramID.String(), tv)
		if err != nil {
			return solana.PublicKey{}, PDASeed{}, fmt.Errorf("failed to save address: %w", err)
		}
		// TODO Set Config on MCMS account
	} else {
		lggr.Infow("Using existing MCMS program", "addr", mcmState.String())
		// TODO: obtain the Seed We can use: https://github.com/smartcontractkit/chainlink/blob/ae4ab024a9ec9eeb248668fa10d72ef037ac2677/deployment/ccip/changeset/solana_state.go#L106-L109
		//	// DeserializeSolanaStateFromAB(address)
		mcmSeed = [32]byte{}
		mcmProgram = mcmState
		return mcmProgram, mcmSeed, nil
	}

	// TODO: modify helper to accept and Environment param
	err := initializeMCM(e, chain, timelockProgram, config.TimelockMinDelay)
	if err != nil {
		return solana.PublicKey{}, PDASeed{}, fmt.Errorf("unable to initialize timelock: %w", err)
	}
	mcmBindings.SetProgramID(mcmProgram)

	// FIXME: review if we need to setup an "AddressLookupTable".

	return solana.PublicKey{}, PDASeed{}, nil
}

func initializeMCM(e deployment.Environment, chain deployment.SolChain, mcmProgram solana.PublicKey) error {
	multisigID := [32]byte{} // FIXME: where should this come from?

	var mcmConfig mcmBindings.MultisigConfig
	err := chain.GetAccountDataBorshInto(e.GetContext(), GetMCMConfigPDA(mcmProgram, multisigID), &mcmConfig)
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
		return fmt.Errorf("unable to get mcm program account info: %w", err)
	}
	err = binary.UnmarshalBorsh(&programData, data.Bytes())
	if err != nil {
		return fmt.Errorf("unable to unmarshal program data: %w", err)
	}

	instruction, err := mcmBindings.NewInitializeInstruction(
		chain.Selector,
		multisigID,
		GetMCMConfigPDA(mcmProgram, multisigID),
		chain.DeployerKey.PublicKey(),
		solana.SystemProgramID,
		mcmProgram,
		programData.Address,
		GetMCMRootMetadataPDA(mcmProgram, multisigID),
		GetMCMExpiringRootAndOpCountPDA(mcmProgram, multisigID),
	).ValidateAndBuild()
	if err != nil {
		return fmt.Errorf("unable to build instruction: %w", err)
	}

	err = chain.Confirm([]solana.Instruction{instruction})
	if err != nil {
		return fmt.Errorf("unable to confirm instructions: %w", err)
	}

	return nil
}

// DeployMCMSWithTimelockProgramsSolana deploys an MCMS program f
// and initializes 3 instances for each of the timelock roles: Bypasser, ProposerMcm, Canceller on an Solana chain.
// as well as the timelock program. It's not necessarily the only way to use
// the timelock and MCMS, but its reasonable pattern.
func DeployMCMSWithTimelockProgramsSolana(
	lggr logger.Logger,
	chain deployment.SolChain,
	ab deployment.AddressBook,
	config mcmsNewTypes.MCMSWithTimelockConfig,
) (*MCMSWithTimelockSolanaDeploy, error) {

	bypasserProgramID, bypasserSeed, err := deployMCMSWithConfigSolana(types.BypasserManyChainMultisig, deployment.Version1_0_0, lggr, chain, ab, config.Bypasser)
	if err != nil {
		return nil, err
	}
	cancellerProgramID, cancellerSeed, err := deployMCMSWithConfigSolana(types.CancellerManyChainMultisig, deployment.Version1_0_0, lggr, chain, ab, config.Canceller)
	if err != nil {
		return nil, err
	}
	proposer, proposerSeed, err := deployMCMSWithConfigSolana(types.ProposerManyChainMultisig, deployment.Version1_0_0, lggr, chain, ab, config.Proposer)
	if err != nil {
		return nil, err
	}

	// TODO: populate this, we need deployMCMSWithConfigSolana to generate a new seed or return the existing one by parsing it from the addressbook.
	result := &MCMSWithTimelockSolanaDeploy{
		McmProgram: bypasser,
	}

	return result, nil
}

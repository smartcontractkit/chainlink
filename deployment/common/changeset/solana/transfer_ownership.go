package solana

import (
	"fmt"
	"maps"
	"slices"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/mr-tron/base58"
	mcms "github.com/smartcontractkit/mcms"
	mcmssdk "github.com/smartcontractkit/mcms/sdk"
	mcmssolanasdk "github.com/smartcontractkit/mcms/sdk/solana"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	accessControllerBindings "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/access_controller"
	mcmBindings "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/mcm"

	"github.com/smartcontractkit/chainlink/deployment"
	state "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
)

const maxAcceptInstructionsPerBatch = 5

// TransferToTimelockSolanaConfig holds the configuration for an ownership transfer changeset
type TransferToTimelockSolanaConfig struct {
	ContractsByChain map[uint64][]OwnableContract
	MinDelay         time.Duration
}

type OwnableContract struct {
	ProgramID solana.PublicKey
	Seed      [32]byte
	OwnerPDA  solana.PublicKey
	Type      deployment.ContractType
}

func (c TransferToTimelockSolanaConfig) Validate(env deployment.Environment) error {
	for chainSelector, contracts := range c.ContractsByChain {
		err := addressBookContains(env.ExistingAddresses, chainSelector,
			commontypes.RBACTimelockProgram,
			commontypes.RBACTimelock,
			commontypes.ManyChainMultisigProgram,
			commontypes.ProposerManyChainMultisig,
		)
		if err != nil {
			return err
		}

		for _, contract := range contracts {
			exists, err := deployment.AddressBookContains(env.ExistingAddresses, chainSelector, contract.ProgramID.String())
			if err != nil {
				return fmt.Errorf("failed to search address book for program id: %w", err)
			}
			if !exists {
				return fmt.Errorf("program id %s not found in address book", contract.ProgramID.String())
			}

			if (contract.Seed == state.PDASeed{}) {
				continue
			}

			exists, err = deployment.AddressBookContains(env.ExistingAddresses, chainSelector, base58.Encode(contract.Seed[:]))
			if err != nil {
				return fmt.Errorf("failed to search address book for seed (%s): %w", base58.Encode(contract.Seed[:]), err)
			}
			if !exists {
				address := solanaAddress(contract.ProgramID, contract.Seed)
				exists, err = deployment.AddressBookContains(env.ExistingAddresses, chainSelector, address)
				if err != nil {
					return fmt.Errorf("failed to search address book for seed (%s): %w", address, err)
				}
			}
			if !exists {
				exists, err = deployment.AddressBookContains(env.ExistingAddresses, chainSelector, string(contract.Seed[:]))
				if err != nil {
					return fmt.Errorf("failed to search address book for seed (%s): %w", string(contract.Seed[:]), err)
				}
			}
			if !exists {
				return fmt.Errorf("seed %s not found in address book", base58.Encode(contract.Seed[:]))
			}
		}
	}

	return nil
}

// TransferToTimelockSolana transfers a set of Solana "contracts" to the Timelock
// signer PDA
func TransferToTimelockSolana(
	env deployment.Environment, cfg TransferToTimelockSolanaConfig,
) (deployment.ChangesetOutput, error) {
	err := cfg.Validate(env)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("invalid transfer ownership config: %w", err)
	}

	mcmsState, err := state.MaybeLoadMCMSWithTimelockStateSolana(env, slices.Collect(maps.Keys(env.SolChains)))
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	batches := []mcmstypes.BatchOperation{}
	timelocks := map[uint64]string{}
	proposers := map[uint64]string{}
	inspectors := map[uint64]mcmssdk.Inspector{}
	instructions := map[uint64][]solana.Instruction{}

	for chainSelector, contractsToTransfer := range cfg.ContractsByChain {
		solChain, ok := env.SolChains[chainSelector]
		if !ok {
			return deployment.ChangesetOutput{}, fmt.Errorf("solana chain not found in environment (selector: %v)", chainSelector)
		}
		chainState, ok := mcmsState[chainSelector]
		if !ok {
			return deployment.ChangesetOutput{}, fmt.Errorf("chain state not found for selector: %v", chainSelector)
		}
		timelocks[chainSelector] = solanaAddress(chainState.TimelockProgram, mcmssolanasdk.PDASeed(chainState.TimelockSeed))
		proposers[chainSelector] = solanaAddress(chainState.McmProgram, mcmssolanasdk.PDASeed(chainState.ProposerMcmSeed))
		inspectors[chainSelector] = mcmssolanasdk.NewInspector(solChain.Client)

		timelockSignerPDA := state.GetTimelockSignerPDA(chainState.TimelockProgram, chainState.TimelockSeed)

		transactions := []mcmstypes.Transaction{}
		for _, contract := range contractsToTransfer {
			transferInstruction, err := transferOwnershipInstruction(contract.ProgramID, contract.Seed, timelockSignerPDA,
				contract.OwnerPDA, solChain.DeployerKey.PublicKey())
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to create transfer ownership instruction: %w", err)
			}
			instructions[chainSelector] = append(instructions[chainSelector], transferInstruction)

			acceptMCMSTransaction, err := acceptMCMSTransaction(contract, timelockSignerPDA)
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to create accept ownership mcms transaction: %w", err)
			}
			transactions = append(transactions, acceptMCMSTransaction)
		}

		// FIXME: remove the chunking logic once we have custom CU limit support in MCMS
		for chunk := range slices.Chunk(transactions, maxAcceptInstructionsPerBatch) {
			batches = append(batches, mcmstypes.BatchOperation{
				ChainSelector: mcmstypes.ChainSelector(chainSelector),
				Transactions:  chunk,
			})
			env.Logger.Debugw("added BatchOperation with accept ownwership instructions",
				"# transactions", len(transactions), "chain", chainSelector)
		}
	}

	// send & confim TransferOwnership instructions
	for chainSelector, chainInstructions := range instructions {
		solChain := env.SolChains[chainSelector]
		// REVIEW: are we limited by the 1232 byte limit? or can we confirm all instructions in one go?
		for _, instruction := range chainInstructions {
			env.Logger.Debugw("confirming solana transfer ownership instruction", "instruction", instruction.ProgramID())
			err = solChain.Confirm([]solana.Instruction{instruction})
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to confirm instruction: %w", err)
			}
		}
	}

	// create timelock proposal with accept transactions
	proposal, err := proposalutils.BuildProposalFromBatchesV2(env, timelocks, proposers, inspectors,
		batches, "proposal to transfer ownership of contracts to timelock", cfg.MinDelay)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
	}
	env.Logger.Debugw("created timelock proposal", "# batches", len(batches))

	return deployment.ChangesetOutput{MCMSTimelockProposals: []mcms.TimelockProposal{*proposal}}, nil
}

type TransferMCMSToTimelockSolanaConfig struct {
	Chains   []uint64
	MinDelay time.Duration
}

func (c *TransferMCMSToTimelockSolanaConfig) Validate(env deployment.Environment) error {
	for _, chainSelector := range c.Chains {
		err := addressBookContains(env.ExistingAddresses, chainSelector,
			commontypes.RBACTimelockProgram,
			commontypes.RBACTimelock,
			commontypes.ManyChainMultisigProgram,
			commontypes.ProposerManyChainMultisig,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func TransferMCMSToTimelockSolana(
	env deployment.Environment, cfg TransferMCMSToTimelockSolanaConfig,
) (deployment.ChangesetOutput, error) {
	err := cfg.Validate(env)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("invalid transfer ownership config: %w", err)
	}

	mcmsState, err := state.MaybeLoadMCMSWithTimelockStateSolana(env, cfg.Chains)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load mcms state: %w", err)
	}

	contracts := map[uint64][]OwnableContract{}
	for chainSelector, chainState := range mcmsState {
		chainContracts := []OwnableContract{
			{
				ProgramID: chainState.McmProgram,
				Seed:      chainState.ProposerMcmSeed,
				OwnerPDA:  state.GetMCMConfigPDA(chainState.McmProgram, chainState.ProposerMcmSeed),
				Type:      commontypes.ProposerManyChainMultisig,
			},
			{
				ProgramID: chainState.McmProgram,
				Seed:      chainState.CancellerMcmSeed,
				OwnerPDA:  state.GetMCMConfigPDA(chainState.McmProgram, chainState.CancellerMcmSeed),
				Type:      commontypes.CancellerManyChainMultisig,
			},
			{
				ProgramID: chainState.McmProgram,
				Seed:      chainState.BypasserMcmSeed,
				OwnerPDA:  state.GetMCMConfigPDA(chainState.McmProgram, chainState.BypasserMcmSeed),
				Type:      commontypes.BypasserManyChainMultisig,
			},
			{
				ProgramID: chainState.AccessControllerProgram,
				OwnerPDA:  chainState.ProposerAccessControllerAccount,
				Type:      commontypes.ProposerAccessControllerAccount,
			},
			{
				ProgramID: chainState.AccessControllerProgram,
				OwnerPDA:  chainState.ExecutorAccessControllerAccount,
				Type:      commontypes.ExecutorAccessControllerAccount,
			},
			{
				ProgramID: chainState.AccessControllerProgram,
				OwnerPDA:  chainState.CancellerAccessControllerAccount,
				Type:      commontypes.CancellerAccessControllerAccount,
			},
			{
				ProgramID: chainState.AccessControllerProgram,
				OwnerPDA:  chainState.BypasserAccessControllerAccount,
				Type:      commontypes.BypasserAccessControllerAccount,
			},
		}

		contracts[chainSelector] = chainContracts
	}

	return TransferToTimelockSolana(env, TransferToTimelockSolanaConfig{
		ContractsByChain: contracts,
		MinDelay:         cfg.MinDelay,
	})
}

func transferOwnershipInstruction(
	programID solana.PublicKey, seed state.PDASeed, proposedOwner, ownerPDA, auth solana.PublicKey,
) (solana.Instruction, error) {
	if (seed == state.PDASeed{}) {
		return newSeedlessTransferOwnershipInstruction(programID, proposedOwner, ownerPDA, auth)
	}
	return newSeededTransferOwnershipInstruction(programID, seed, proposedOwner, ownerPDA, auth)
}

func acceptMCMSTransaction(
	contract OwnableContract,
	authority solana.PublicKey,
) (mcmstypes.Transaction, error) {
	acceptInstruction, err := acceptOwnershipInstruction(contract.ProgramID, contract.Seed, contract.OwnerPDA, authority)
	if err != nil {
		return mcmstypes.Transaction{}, fmt.Errorf("failed to build accept ownership instruction: %w", err)
	}
	acceptMCMSTx, err := mcmssolanasdk.NewTransactionFromInstruction(acceptInstruction, string(contract.Type), []string{})
	if err != nil {
		return mcmstypes.Transaction{}, fmt.Errorf("failed to build mcms transaction from accept ownership instruction: %w", err)
	}

	return acceptMCMSTx, nil
}

func acceptOwnershipInstruction(
	programID solana.PublicKey, seed state.PDASeed, ownerPDA, auth solana.PublicKey,
) (solana.Instruction, error) {
	if (seed == state.PDASeed{}) {
		return newSeedlessAcceptOwnershipInstruction(programID, ownerPDA, auth)
	}
	return newSeededAcceptOwnershipInstruction(programID, seed, ownerPDA, auth)
}

func newSeededTransferOwnershipInstruction(
	programID solana.PublicKey, seed state.PDASeed, proposedOwner, config, authority solana.PublicKey,
) (solana.Instruction, error) {
	ix, err := mcmBindings.NewTransferOwnershipInstruction(seed, proposedOwner, config, authority).ValidateAndBuild()
	return &seededInstruction{ix, programID}, err
}

func newSeededAcceptOwnershipInstruction(
	programID solana.PublicKey, seed state.PDASeed, config, authority solana.PublicKey,
) (solana.Instruction, error) {
	ix, err := mcmBindings.NewAcceptOwnershipInstruction(seed, config, authority).ValidateAndBuild()
	return &seededInstruction{ix, programID}, err
}

func newSeedlessTransferOwnershipInstruction(
	programID, proposedOwner, config, authority solana.PublicKey,
) (solana.Instruction, error) {
	ix, err := accessControllerBindings.NewTransferOwnershipInstruction(proposedOwner, config, authority).ValidateAndBuild()
	return &seedlessInstruction{ix, programID}, err
}

func newSeedlessAcceptOwnershipInstruction(
	programID, config, authority solana.PublicKey,
) (solana.Instruction, error) {
	ix, err := accessControllerBindings.NewAcceptOwnershipInstruction(config, authority).ValidateAndBuild()
	return &seedlessInstruction{ix, programID}, err
}

type seedlessInstruction struct {
	*accessControllerBindings.Instruction
	programID solana.PublicKey
}

func (s *seedlessInstruction) ProgramID() solana.PublicKey {
	return s.programID
}

type seededInstruction struct {
	*mcmBindings.Instruction
	programID solana.PublicKey
}

func (s *seededInstruction) ProgramID() solana.PublicKey {
	return s.programID
}

func addressBookContains(addressBook deployment.AddressBook, chainSelector uint64, ctypes ...deployment.ContractType) error {
	for _, ctype := range ctypes {
		_, err := deployment.SearchAddressBook(addressBook, chainSelector, ctype)
		if err != nil {
			return fmt.Errorf("address book does not contain a %s contract for chain %d", ctype, chainSelector)
		}
	}
	return nil
}

var solanaAddress = mcmssolanasdk.ContractAddress

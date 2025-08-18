package solana

import (
	"fmt"
	"maps"
	"slices"

	"github.com/Masterminds/semver/v3"
	"github.com/gagliardetto/solana-go"
	"github.com/mr-tron/base58"
	"github.com/smartcontractkit/mcms"
	mcmssdk "github.com/smartcontractkit/mcms/sdk"
	mcmssolanasdk "github.com/smartcontractkit/mcms/sdk/solana"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	accessControllerBindings "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_1/access_controller"
	mcmBindings "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_1/mcm"

	cldfsol "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
)

const maxAcceptInstructionsPerBatch = 5

// TransferToTimelockSolanaConfig holds the configuration for an ownership transfer changeset
type TransferToTimelockSolanaConfig struct {
	ContractsByChain map[uint64][]OwnableContract
	MCMSCfg          proposalutils.TimelockConfig
}

type OwnableContract struct {
	ProgramID solana.PublicKey
	Seed      [32]byte
	OwnerPDA  solana.PublicKey
	Type      cldf.ContractType
}

// TransferToTimelockSolana transfers a set of Solana "contracts" to the Timelock
// signer PDA.
// The "transfer ownership" instructions are immediately sent and
// confirmed onchain. The "accept ownership" instructions are added to an MCMS
// timelock proposal that should be executed using the standard MCMS workflows.
type TransferToTimelockSolana struct{}

func (t *TransferToTimelockSolana) VerifyPreconditions(
	env cldf.Environment, config TransferToTimelockSolanaConfig,
) error {
	for chainSelector, contracts := range config.ContractsByChain {
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
			exists, err := cldf.AddressBookContains(env.ExistingAddresses, chainSelector, contract.ProgramID.String())
			if err != nil {
				return fmt.Errorf("failed to search address book for program id: %w", err)
			}
			if !exists {
				return fmt.Errorf("program id %s not found in address book", contract.ProgramID.String())
			}

			if (contract.Seed == state.PDASeed{}) {
				continue
			}

			exists, err = cldf.AddressBookContains(env.ExistingAddresses, chainSelector, base58.Encode(contract.Seed[:]))
			if err != nil {
				return fmt.Errorf("failed to search address book for seed (%s): %w", base58.Encode(contract.Seed[:]), err)
			}
			if !exists {
				address := solanaAddress(contract.ProgramID, contract.Seed)
				exists, err = cldf.AddressBookContains(env.ExistingAddresses, chainSelector, address)
				if err != nil {
					return fmt.Errorf("failed to search address book for seed (%s): %w", address, err)
				}
			}
			if !exists {
				exists, err = cldf.AddressBookContains(env.ExistingAddresses, chainSelector, string(contract.Seed[:]))
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

func (t *TransferToTimelockSolana) Apply(
	env cldf.Environment, cfg TransferToTimelockSolanaConfig,
) (cldf.ChangesetOutput, error) {
	solChains := env.BlockChains.SolanaChains()
	mcmsState, err := state.MaybeLoadMCMSWithTimelockStateSolana(env, slices.Collect(maps.Keys(solChains)))
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	batches := []mcmstypes.BatchOperation{}
	timelocks := map[uint64]string{}
	proposers := map[uint64]string{}
	inspectors := map[uint64]mcmssdk.Inspector{}
	var out cldf.ChangesetOutput
	for chainSelector, contractsToTransfer := range cfg.ContractsByChain {
		solChain, ok := solChains[chainSelector]
		if !ok {
			return cldf.ChangesetOutput{}, fmt.Errorf("solana chain not found in environment (selector: %v)", chainSelector)
		}
		chainState, ok := mcmsState[chainSelector]
		if !ok {
			return cldf.ChangesetOutput{}, fmt.Errorf("chain state not found for selector: %v", chainSelector)
		}

		inspectors[chainSelector] = mcmssolanasdk.NewInspector(solChain.Client)
		timelocks[chainSelector] = solanaAddress(chainState.TimelockProgram, mcmssolanasdk.PDASeed(chainState.TimelockSeed))
		proposers[chainSelector] = solanaAddress(chainState.McmProgram, mcmssolanasdk.PDASeed(chainState.ProposerMcmSeed))

		for _, contract := range contractsToTransfer {
			execOut, err := operations.ExecuteOperation(env.OperationsBundle,
				operations.NewOperation(
					"transfer-ownership",
					&deployment.Version1_0_0,
					"transfers ownership of contracts to mcms",
					TransferToTimelockSolanaOp,
				),
				Deps{
					Env:   env,
					State: chainState,
					Chain: solChain,
				},
				TransferToTimelockInput{
					Contract: contract,
					MCMSCfg:  cfg.MCMSCfg,
				},
			)
			if err != nil {
				return out, err
			}

			batches = append(batches, execOut.Output.Batches...)
		}
	}

	// create timelock proposal with accept transactions
	proposal, err := proposalutils.BuildProposalFromBatchesV2(env, timelocks, proposers, inspectors,
		batches, "proposal to transfer ownership of contracts to timelock", cfg.MCMSCfg)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
	}
	env.Logger.Debugw("created timelock proposal", "# batches", len(batches))

	return cldf.ChangesetOutput{MCMSTimelockProposals: []mcms.TimelockProposal{*proposal}}, nil
}

// ContractConfig defines the configuration for a contract ownership transfer
type ContractConfig struct {
	ContractType datastore.ContractType
	StateType    datastore.ContractType
	OperationID  string
	Description  string
}

// TransferOwnershipRequest represents a generic ownership transfer request
type TransferOwnershipRequest struct {
	ChainSel                    uint64
	CurrentOwner, ProposedOwner solana.PublicKey
	Version                     string
	Qualifier                   string
	MCMSCfg                     proposalutils.TimelockConfig
	ContractConfig              ContractConfig
}

// genericTransferOwnership handles the common ownership transfer logic
func GenericTransferOwnership(env cldf.Environment, req *TransferOwnershipRequest) (cldf.ChangesetOutput, error) {
	var out cldf.ChangesetOutput
	version := semver.MustParse(req.Version)

	// Build address references
	contractStateRef := datastore.NewAddressRefKey(req.ChainSel, req.ContractConfig.StateType, version, req.Qualifier)
	contractRef := datastore.NewAddressRefKey(req.ChainSel, req.ContractConfig.ContractType, version, req.Qualifier)

	// Get contract addresses
	contract, err := env.DataStore.Addresses().Get(contractRef)
	if err != nil {
		return out, fmt.Errorf("failed to get contract address: %w", err)
	}

	contractState, err := env.DataStore.Addresses().Get(contractStateRef)
	if err != nil {
		return out, fmt.Errorf("failed to get contract state address: %w", err)
	}

	// Load MCMS state
	mcmsState, err := state.MaybeLoadMCMSWithTimelockChainStateSolanaV2(
		env.DataStore.Addresses().Filter(datastore.AddressRefByChainSelector(req.ChainSel)))
	if err != nil {
		return out, err
	}

	solChain := env.BlockChains.SolanaChains()[req.ChainSel]

	// Execute the transfer operation
	execOut, err := operations.ExecuteOperation(env.OperationsBundle,
		operations.NewOperation(
			req.ContractConfig.OperationID,
			version,
			req.ContractConfig.Description,
			TransferToTimelockSolanaOp,
		),
		Deps{
			Env:   env,
			State: mcmsState,
			Chain: solChain,
		},
		TransferToTimelockInput{
			Contract: OwnableContract{
				Type:      cldf.ContractType(req.ContractConfig.ContractType),
				ProgramID: solana.MustPublicKeyFromBase58(contract.Address),
				OwnerPDA:  solana.MustPublicKeyFromBase58(contractState.Address),
			},
			MCMSCfg: req.MCMSCfg,
		},
	)
	if err != nil {
		return out, err
	}

	// Build proposal maps
	timelocks := map[uint64]string{}
	proposers := map[uint64]string{}
	inspectors := map[uint64]mcmssdk.Inspector{}

	inspectors[req.ChainSel] = mcmssolanasdk.NewInspector(solChain.Client)
	timelocks[req.ChainSel] = mcmssolanasdk.ContractAddress(mcmsState.TimelockProgram, mcmssolanasdk.PDASeed(mcmsState.TimelockSeed))
	proposers[req.ChainSel] = mcmssolanasdk.ContractAddress(mcmsState.McmProgram, mcmssolanasdk.PDASeed(mcmsState.ProposerMcmSeed))

	// Create timelock proposal
	proposal, err := proposalutils.BuildProposalFromBatchesV2(env, timelocks, proposers, inspectors,
		execOut.Output.Batches, fmt.Sprintf("proposal to transfer ownership of %s to timelock", req.ContractConfig.ContractType), req.MCMSCfg)
	if err != nil {
		return out, fmt.Errorf("failed to build proposal: %w", err)
	}
	env.Logger.Debugw("created timelock proposal", "# batches", len(execOut.Output.Batches))

	out.MCMSTimelockProposals = []mcms.TimelockProposal{*proposal}
	return out, nil
}

// genericVerifyPreconditions handles the common precondition verification logic
func GenericVerifyPreconditions(env cldf.Environment, chainSel uint64, version, qualifier string, contractType datastore.ContractType) error {
	// Validate version
	if _, err := semver.NewVersion(version); err != nil {
		return err
	}

	// Check if chain exists
	if _, ok := env.BlockChains.SolanaChains()[chainSel]; !ok {
		return fmt.Errorf("solana chain not found for chain selector %d", chainSel)
	}

	// Verify contract exists
	v := semver.MustParse(version)
	contractKey := datastore.NewAddressRefKey(chainSel, contractType, v, qualifier)
	if _, err := env.DataStore.Addresses().Get(contractKey); err != nil {
		return fmt.Errorf("failed to get %s for chain selector %d: %w", contractType, chainSel, err)
	}

	return nil
}

type (
	Deps struct {
		Env   cldf.Environment
		State *state.MCMSWithTimelockStateSolana
		Chain cldfsol.Chain
	}

	TransferToTimelockInput struct {
		Contract OwnableContract
		MCMSCfg  proposalutils.TimelockConfig
	}

	TransferToTimelockOutput struct {
		Batches []mcmstypes.BatchOperation
	}
)

func TransferToTimelockSolanaOp(b operations.Bundle, deps Deps, in TransferToTimelockInput) (TransferToTimelockOutput, error) {
	var out TransferToTimelockOutput

	solChain := deps.Chain
	chainState := deps.State

	batches := []mcmstypes.BatchOperation{}
	timelocks := map[uint64]string{}
	proposers := map[uint64]string{}
	inspectors := map[uint64]mcmssdk.Inspector{}
	instructions := []solana.Instruction{}
	chainSelector := solChain.ChainSelector()

	timelocks[chainSelector] = solanaAddress(chainState.TimelockProgram, mcmssolanasdk.PDASeed(chainState.TimelockSeed))
	proposers[chainSelector] = solanaAddress(chainState.McmProgram, mcmssolanasdk.PDASeed(chainState.ProposerMcmSeed))
	inspectors[chainSelector] = mcmssolanasdk.NewInspector(solChain.Client)

	timelockSignerPDA := state.GetTimelockSignerPDA(chainState.TimelockProgram, chainState.TimelockSeed)

	transactions := []mcmstypes.Transaction{}
	contract := in.Contract
	transferInstruction, err := transferOwnershipInstruction(contract.ProgramID, contract.Seed, timelockSignerPDA,
		contract.OwnerPDA, solChain.DeployerKey.PublicKey())
	if err != nil {
		return out, fmt.Errorf("failed to create transfer ownership instruction: %w", err)
	}
	instructions = append(instructions, transferInstruction)

	acceptMCMSTransaction, err := acceptMCMSTransaction(contract, timelockSignerPDA)
	if err != nil {
		return out, fmt.Errorf("failed to create accept ownership mcms transaction: %w", err)
	}
	transactions = append(transactions, acceptMCMSTransaction)

	// FIXME: remove the chunking logic once we have custom CU limit support in MCMS
	for chunk := range slices.Chunk(transactions, maxAcceptInstructionsPerBatch) {
		batches = append(batches, mcmstypes.BatchOperation{
			ChainSelector: mcmstypes.ChainSelector(chainSelector),
			Transactions:  chunk,
		})
		b.Logger.Debugw("added BatchOperation with accept ownwership instructions",
			"# transactions", len(transactions), "chain", chainSelector)
	}

	for _, instruction := range instructions {
		b.Logger.Debugw("confirming solana transfer ownership instruction", "instruction", instruction.ProgramID())
		err = solChain.Confirm([]solana.Instruction{instruction})
		if err != nil {
			return out, fmt.Errorf("failed to confirm instruction: %w", err)
		}
	}
	out.Batches = append(out.Batches, batches...)

	return out, nil
}

type TransferMCMSToTimelockSolanaConfig struct {
	Chains  []uint64
	MCMSCfg proposalutils.TimelockConfig
}

// TransferMCMSToTimelockSolana transfers set MCMS "contracts" to the timelock
// signer PDA. It relies on the TransferToTimelockSolana changeset and merely
// adds the MCM, Timelock and AccessController contracts found in the address
// book to the list of contracts to transfer.
type TransferMCMSToTimelockSolana struct{}

func (t TransferMCMSToTimelockSolana) VerifyPreconditions(
	env cldf.Environment, config TransferMCMSToTimelockSolanaConfig,
) error {
	for _, chainSelector := range config.Chains {
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

func (t TransferMCMSToTimelockSolana) Apply(
	env cldf.Environment, cfg TransferMCMSToTimelockSolanaConfig,
) (cldf.ChangesetOutput, error) {
	mcmsState, err := state.MaybeLoadMCMSWithTimelockStateSolana(env, cfg.Chains)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load mcms state: %w", err)
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
				ProgramID: chainState.TimelockProgram,
				Seed:      chainState.TimelockSeed,
				OwnerPDA:  state.GetTimelockConfigPDA(chainState.TimelockProgram, chainState.TimelockSeed),
				Type:      commontypes.RBACTimelock,
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

	return new(TransferToTimelockSolana).Apply(env, TransferToTimelockSolanaConfig{
		ContractsByChain: contracts,
		MCMSCfg:          cfg.MCMSCfg,
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

func acceptOwnershipInstruction(programID solana.PublicKey, seed state.PDASeed, ownerPDA, auth solana.PublicKey,
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

func addressBookContains(addressBook cldf.AddressBook, chainSelector uint64, ctypes ...cldf.ContractType) error {
	for _, ctype := range ctypes {
		_, err := cldf.SearchAddressBook(addressBook, chainSelector, ctype)
		if err != nil {
			return fmt.Errorf("address book does not contain a %s contract for chain %d", ctype, chainSelector)
		}
	}
	return nil
}

var solanaAddress = mcmssolanasdk.ContractAddress

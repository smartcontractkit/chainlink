package ccip_attestation

import (
	"errors"
	"fmt"

	"github.com/gagliardetto/solana-go"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/mcms"
	"github.com/smartcontractkit/mcms/sdk"

	mcmsSolana "github.com/smartcontractkit/mcms/sdk/solana"
	mcmsTypes "github.com/smartcontractkit/mcms/types"

	signer_registry "github.com/smartcontractkit/chainlink/deployment/ccip/shared/bindings/signer_registry_solana"

	cldf_solana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cs_solana "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/solana_v0_1_1"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"

	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
)

var _ cldf.ChangeSet[TransferSignerRegistryToMCMSWithTimelockSolanaConfig] = TransferSignerRegistryToMCMSWithTimelockSolanaChangeset

type TransferSignerRegistryToMCMSWithTimelockSolanaConfig struct {
	ChainSelector uint64
	CurrentOwner  solana.PublicKey
	ProposedOwner solana.PublicKey
	// MCMSCfg is for the accept ownership proposal
	MCMSCfg proposalutils.TimelockConfig
}

func (c TransferSignerRegistryToMCMSWithTimelockSolanaConfig) Validate(e cldf.Environment) error {
	if err := cldf.IsValidChainSelector(c.ChainSelector); err != nil {
		return fmt.Errorf("invalid chain selector: %d - %w", c.ChainSelector, err)
	}
	family, _ := chainsel.GetSelectorFamily(c.ChainSelector)
	if family != chainsel.FamilySolana {
		return fmt.Errorf("chain %d is not a solana chain", c.ChainSelector)
	}
	ccipState, err := stateview.LoadOnchainStateSolana(e)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}
	if len(ccipState.SolChains) == 0 {
		return errors.New("no chains found")
	}
	solChain := e.BlockChains.SolanaChains()[c.ChainSelector]
	addresses, err := e.ExistingAddresses.AddressesForChain(c.ChainSelector)
	if err != nil {
		return fmt.Errorf("failed to get addresses for chain: %w", err)
	}
	_, err = state.MaybeLoadMCMSWithTimelockChainStateSolana(solChain, addresses)
	if err != nil {
		return fmt.Errorf("failed to load mcm state: %w", err)
	}
	// If there is no timelock and mcms proposer on the chain, the transfer will fail.
	timelockID, err := cldf.SearchAddressBook(e.ExistingAddresses, c.ChainSelector, types.RBACTimelock)
	if err != nil {
		return fmt.Errorf("timelock not present on the chain %w", err)
	}
	proposerID, err := cldf.SearchAddressBook(e.ExistingAddresses, c.ChainSelector, types.ProposerManyChainMultisig)
	if err != nil {
		return fmt.Errorf("mcms proposer not present on the chain %w", err)
	}
	// Make sure addresses are correctly parsed. Format is: "programID.PDASeed"
	_, _, err = mcmsSolana.ParseContractAddress(timelockID)
	if err != nil {
		return fmt.Errorf("failed to parse timelock address: %w", err)
	}
	_, _, err = mcmsSolana.ParseContractAddress(proposerID)
	if err != nil {
		return fmt.Errorf("failed to parse proposer address: %w", err)
	}
	return nil
}

// TransferCCIPToMCMSWithTimelockSolana creates a changeset that transfers ownership of the
// signer registry program to the timelock on the chain and generates a corresponding proposal
// with the accept ownership tx  to complete the transfer. It assumes that DeployMCMSWithTimelock
// for solana has already been run s.t. the timelock and mcms exist on the chain and that the
// proposed address to transfer ownership is currently owned by the deployer key.
func TransferSignerRegistryToMCMSWithTimelockSolanaChangeset(
	e cldf.Environment,
	cfg TransferSignerRegistryToMCMSWithTimelockSolanaConfig,
) (cldf.ChangesetOutput, error) {
	if err := cfg.Validate(e); err != nil {
		return cldf.ChangesetOutput{}, err
	}
	var batches []mcmsTypes.BatchOperation

	timelocks := map[uint64]string{}
	proposers := map[uint64]string{}
	inspectors := map[uint64]sdk.Inspector{}

	chainSelector := cfg.ChainSelector

	solChain := e.BlockChains.SolanaChains()[chainSelector]
	addresses, _ := e.ExistingAddresses.AddressesForChain(chainSelector)
	mcmState, _ := state.MaybeLoadMCMSWithTimelockChainStateSolana(solChain, addresses)

	currentOwner := solChain.DeployerKey.PublicKey()
	if !cfg.CurrentOwner.IsZero() {
		currentOwner = cfg.CurrentOwner
	}
	timelockSigner := state.GetTimelockSignerPDA(mcmState.TimelockProgram, mcmState.TimelockSeed)
	proposedOwner := timelockSigner
	if !cfg.ProposedOwner.IsZero() {
		proposedOwner = cfg.ProposedOwner
	}
	if currentOwner.Equals(proposedOwner) {
		return cldf.ChangesetOutput{}, fmt.Errorf("current owner and proposed owner are the same: %s", currentOwner)
	}

	timelocks[solChain.Selector] = mcmsSolana.ContractAddress(
		mcmState.TimelockProgram,
		mcmsSolana.PDASeed(mcmState.TimelockSeed),
	)
	proposers[solChain.Selector] = mcmsSolana.ContractAddress(mcmState.McmProgram, mcmsSolana.PDASeed(mcmState.ProposerMcmSeed))
	inspectors[solChain.Selector] = mcmsSolana.NewInspector(solChain.Client)
	mcmsTxs, err := transferOwnershipSignerRegistry(
		solChain,
		currentOwner,
		proposedOwner,
		timelockSigner,
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to transfer ownership of base signer registry: %w", err)
	}
	batches = append(batches, mcmsTypes.BatchOperation{
		ChainSelector: mcmsTypes.ChainSelector(chainSelector),
		Transactions:  mcmsTxs,
	})

	proposal, err := proposalutils.BuildProposalFromBatchesV2(
		e,
		timelocks,
		proposers,
		inspectors,
		batches,
		"proposal to transfer ownership of CCIP contracts to timelock",
		cfg.MCMSCfg)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
	}

	return cldf.ChangesetOutput{MCMSTimelockProposals: []mcms.TimelockProposal{*proposal}}, nil
}

func transferOwnershipSignerRegistry(
	solChain cldf_solana.Chain,
	currentOwner solana.PublicKey,
	proposedOwner solana.PublicKey,
	timelockSigner solana.PublicKey,
) ([]mcmsTypes.Transaction, error) {
	var result []mcmsTypes.Transaction

	programID := signer_registry.ProgramID
	configPda, _, _ := solana.FindProgramAddress([][]byte{[]byte("config")}, signer_registry.ProgramID)
	eventAuthorityPda, _, _ := solana.FindProgramAddress([][]byte{[]byte("__event_authority")}, signer_registry.ProgramID)

	// Build specialized closures
	buildTransfer := func(proposedOwner, config, authority solana.PublicKey) (solana.Instruction, error) {
		ix, err := signer_registry.NewProposeNewOwnerInstruction(
			proposedOwner, authority, config, eventAuthorityPda, signer_registry.ProgramID,
		)
		if err != nil {
			return nil, err
		}
		ixData, err := ix.Data()
		if err != nil {
			return nil, fmt.Errorf("failed to extract data payload from signer registry transfer ownership instruction: %w", err)
		}
		transferOwnershipIx := solana.NewInstruction(programID, ix.Accounts(), ixData)
		for _, acc := range transferOwnershipIx.Accounts() {
			if acc.PublicKey == timelockSigner {
				acc.IsSigner = false
			}
		}
		return transferOwnershipIx, nil
	}
	buildAccept := func(config, newOwnerAuthority solana.PublicKey) (solana.Instruction, error) {
		// If the router has its own accept function, use that
		ix, err := signer_registry.NewAcceptOwnershipInstruction(
			newOwnerAuthority, config, eventAuthorityPda, signer_registry.ProgramID,
		)
		if err != nil {
			return nil, err
		}
		ixData, err := ix.Data()
		if err != nil {
			return nil, fmt.Errorf("failed to extract data payload from signer registry accept ownership instruction: %w", err)
		}
		acceptOwnershipIx := solana.NewInstruction(programID, ix.Accounts(), ixData)
		for _, acc := range acceptOwnershipIx.Accounts() {
			if acc.PublicKey == timelockSigner {
				acc.IsSigner = false
			}
		}
		return acceptOwnershipIx, nil
	}

	tx, err := cs_solana.TransferAndWrapAcceptOwnership(
		buildTransfer,
		buildAccept,
		programID,
		proposedOwner,
		configPda,
		currentOwner,
		solChain,
		shared.SVMSignerRegistry,
		timelockSigner,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to transfer signer registry ownership: %w", err)
	}

	result = append(result, tx)
	return result, nil
}

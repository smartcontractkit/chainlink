package solana

import (
	"errors"
	"fmt"

	"github.com/gagliardetto/solana-go"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/mcms"
	"github.com/smartcontractkit/mcms/sdk"
	mcmsSolana "github.com/smartcontractkit/mcms/sdk/solana"
	mcmsTypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink/deployment"
	state2 "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
)

var _ deployment.ChangeSet[TransferCCIPToMCMSWithTimelockSolanaConfig] = TransferCCIPToMCMSWithTimelockSolana

// CCIPContractsToTransfer is a struct that represents the contracts we want to transfer. Each contract set to true will be transferred.
type CCIPContractsToTransfer struct {
	Router    bool
	FeeQuoter bool
	OffRamp   bool
	// Token Pool PDA -> Token Mint
	LockReleaseTokenPools map[solana.PublicKey]solana.PublicKey
	BurnMintTokenPools    map[solana.PublicKey]solana.PublicKey
	RMNRemote             bool
}

type TransferCCIPToMCMSWithTimelockSolanaConfig struct {
	// ContractsByChain is a map of chain selector the contracts we want to transfer.
	// Each contract set to true will be transferred
	ContractsByChain map[uint64]CCIPContractsToTransfer
	// MCMSCfg is for the accept ownership proposal
	MCMSCfg proposalutils.TimelockConfig
}

// ValidateContracts checks if the required contracts are present on the chain
func ValidateContracts(state state2.SolCCIPChainState, chainSelector uint64, contracts CCIPContractsToTransfer) error {
	contractChecks := []struct {
		enabled bool
		value   solana.PublicKey
		name    string
	}{
		{contracts.Router, state.Router, "Router"},
		{contracts.FeeQuoter, state.FeeQuoter, "FeeQuoter"},
		{contracts.OffRamp, state.OffRamp, "OffRamp"},
		{contracts.RMNRemote, state.RMNRemote, "RMNRemote"},
	}

	for _, check := range contractChecks {
		if check.enabled && check.value.IsZero() {
			return fmt.Errorf("missing required contract %s on chain %d", check.name, chainSelector)
		}
	}

	return nil
}

func (cfg TransferCCIPToMCMSWithTimelockSolanaConfig) Validate(e deployment.Environment) error {
	ccipState, err := state2.LoadOnchainStateSolana(e)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}
	if len(ccipState.SolChains) == 0 {
		return errors.New("no chains found")
	}
	for chainSelector, contractsEnabled := range cfg.ContractsByChain {
		if _, ok := e.SolChains[chainSelector]; !ok {
			return fmt.Errorf("chain %d not found in environment", chainSelector)
		}
		solChain := e.SolChains[chainSelector]
		// Load MCM state
		addresses, err := e.ExistingAddresses.AddressesForChain(chainSelector)
		if err != nil {
			return fmt.Errorf("failed to load addresses for chain %d: %w", chainSelector, err)
		}
		_, err = state.MaybeLoadMCMSWithTimelockChainStateSolana(solChain, addresses)
		if err != nil {
			return fmt.Errorf("failed to load mcm state: %w", err)
		}
		chainFamily, err := chain_selectors.GetSelectorFamily(chainSelector)
		if err != nil {
			return fmt.Errorf("failed to get chain family for chain %d: %w", chainSelector, err)
		}
		if chainFamily != chain_selectors.FamilySolana {
			return fmt.Errorf("chain %d is not a solana chain", chainSelector)
		}
		state, ok := ccipState.SolChains[chainSelector]
		if !ok {
			return fmt.Errorf("no state found for chain %d", chainSelector)
		}
		err = ValidateContracts(state, chainSelector, contractsEnabled)
		if err != nil {
			return fmt.Errorf("failed to validate contracts for chain %d: %w", chainSelector, err)
		}
		// If there is no timelock and mcms proposer on the chain, the transfer will fail.
		timelockID, err := deployment.SearchAddressBook(e.ExistingAddresses, chainSelector, types.RBACTimelock)
		if err != nil {
			return fmt.Errorf("timelock not present on the chain %w", err)
		}
		proposerID, err := deployment.SearchAddressBook(e.ExistingAddresses, chainSelector, types.ProposerManyChainMultisig)
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
	}
	return nil
}

// TransferCCIPToMCMSWithTimelockSolana creates a changeset that transfers ownership of all the
// CCIP Programs in the provided configuration to the timelock on the chain and generates
// a corresponding proposal with the accept ownership txs  to complete the transfer.
// It assumes that DeployMCMSWithTimelock for solana has already been run s.t.
// the timelock and mcms exist on the chain and that the proposed addresses to transfer ownership
// are currently owned by the deployer key.
func TransferCCIPToMCMSWithTimelockSolana(
	e deployment.Environment,
	cfg TransferCCIPToMCMSWithTimelockSolanaConfig,
) (deployment.ChangesetOutput, error) {
	if err := cfg.Validate(e); err != nil {
		return deployment.ChangesetOutput{}, err
	}
	var batches []mcmsTypes.BatchOperation

	ccipState, err := state2.LoadOnchainStateSolana(e)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	timelocks := map[uint64]string{}
	proposers := map[uint64]string{}
	inspectors := map[uint64]sdk.Inspector{}
	for chainSelector, contractsToTransfer := range cfg.ContractsByChain {
		solChain := e.SolChains[chainSelector]
		addresses, _ := e.ExistingAddresses.AddressesForChain(chainSelector)
		mcmState, _ := state.MaybeLoadMCMSWithTimelockChainStateSolana(solChain, addresses)

		timelocks[solChain.Selector] = mcmsSolana.ContractAddress(
			mcmState.TimelockProgram,
			mcmsSolana.PDASeed(mcmState.TimelockSeed),
		)
		proposers[solChain.Selector] = mcmsSolana.ContractAddress(mcmState.McmProgram, mcmsSolana.PDASeed(mcmState.ProposerMcmSeed))
		inspectors[solChain.Selector] = mcmsSolana.NewInspector(solChain.Client)
		if contractsToTransfer.Router {
			mcmsTxs, err := transferOwnershipRouter(
				ccipState,
				chainSelector,
				solChain,
				mcmState.TimelockProgram,
				mcmState.TimelockSeed,
			)
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to transfer ownership of router: %w", err)
			}
			batches = append(batches, mcmsTypes.BatchOperation{
				ChainSelector: mcmsTypes.ChainSelector(chainSelector),
				Transactions:  mcmsTxs,
			})
		}

		if contractsToTransfer.FeeQuoter {
			mcmsTxs, err := transferOwnershipFeeQuoter(
				ccipState,
				chainSelector,
				solChain,
				mcmState.TimelockProgram,
				mcmState.TimelockSeed,
			)
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to transfer ownership of fee quoter: %w", err)
			}
			batches = append(batches, mcmsTypes.BatchOperation{
				ChainSelector: mcmsTypes.ChainSelector(chainSelector),
				Transactions:  mcmsTxs,
			})
		}

		if contractsToTransfer.OffRamp {
			mcmsTxs, err := transferOwnershipOffRamp(
				ccipState,
				chainSelector,
				solChain,
				mcmState.TimelockProgram,
				mcmState.TimelockSeed,
			)
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to transfer ownership of offRamp: %w", err)
			}
			batches = append(batches, mcmsTypes.BatchOperation{
				ChainSelector: mcmsTypes.ChainSelector(chainSelector),
				Transactions:  mcmsTxs,
			})
		}
		for tokenPoolConfigPDA, tokenMint := range contractsToTransfer.LockReleaseTokenPools {
			mcmsTxs, err := transferOwnershipLockReleaseTokenPools(
				ccipState,
				tokenPoolConfigPDA,
				tokenMint,
				chainSelector,
				solChain,
				mcmState.TimelockProgram,
				mcmState.TimelockSeed,
			)
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to transfer ownership of lock-release token pools: %w", err)
			}
			batches = append(batches, mcmsTypes.BatchOperation{
				ChainSelector: mcmsTypes.ChainSelector(chainSelector),
				Transactions:  mcmsTxs,
			})
		}

		for tokenPoolConfigPDA, tokenMint := range contractsToTransfer.BurnMintTokenPools {
			mcmsTxs, err := transferOwnershipBurnMintTokenPools(
				ccipState,
				tokenPoolConfigPDA,
				tokenMint,
				chainSelector,
				solChain,
				mcmState.TimelockProgram,
				mcmState.TimelockSeed,
			)
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to transfer ownership of burn-mint token pools: %w", err)
			}
			batches = append(batches, mcmsTypes.BatchOperation{
				ChainSelector: mcmsTypes.ChainSelector(chainSelector),
				Transactions:  mcmsTxs,
			})
		}

		if contractsToTransfer.RMNRemote {
			mcmsTxs, err := transferOwnershipRMNRemote(
				ccipState,
				chainSelector,
				solChain,
				mcmState.TimelockProgram,
				mcmState.TimelockSeed,
			)
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to transfer ownership of rmnremote: %w", err)
			}
			batches = append(batches, mcmsTypes.BatchOperation{
				ChainSelector: mcmsTypes.ChainSelector(chainSelector),
				Transactions:  mcmsTxs,
			})
		}
	}

	proposal, err := proposalutils.BuildProposalFromBatchesV2(
		e,
		timelocks,
		proposers,
		inspectors,
		batches,
		"proposal to transfer ownership of CCIP contracts to timelock",
		cfg.MCMSCfg)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
	}

	return deployment.ChangesetOutput{MCMSTimelockProposals: []mcms.TimelockProposal{*proposal}}, nil
}

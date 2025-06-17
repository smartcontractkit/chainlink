package example

import (
	"errors"
	"fmt"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/smartcontractkit/mcms"
	mcmssolanasdk "github.com/smartcontractkit/mcms/sdk/solana"
	"github.com/smartcontractkit/mcms/types"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	solanachangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset/solana"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

var _ cldf.ChangeSetV2[TransferFromTimelockConfig] = TransferFromTimelock{}

type TransferData struct {
	To     solana.PublicKey
	Amount uint64
}
type TransferFromTimelockConfig struct {
	TimelockCfg     proposalutils.TimelockConfig
	AmountsPerChain map[uint64]TransferData
}

// TransferFromTimelock is a changeset that transfer funds from the timelock signer PDA
// to the address provided in the config. It will return an mcms proposal to sign containing
// the funds transfer transaction.
type TransferFromTimelock struct{}

// VerifyPreconditions checks if the deployer has enough SOL to fund the MCMS signers on each chain.
func (f TransferFromTimelock) VerifyPreconditions(e cldf.Environment, config TransferFromTimelockConfig) error {
	solChains := e.BlockChains.SolanaChains()
	// the number of accounts to fund per chain (bypasser, canceller, proposer, timelock)
	for chainSelector, amountCfg := range config.AmountsPerChain {
		solChain, ok := solChains[chainSelector]
		if !ok {
			return fmt.Errorf("solana chain not found for selector %d", chainSelector)
		}
		if amountCfg.To.IsZero() {
			return errors.New("destination address is empty")
		}
		addresses, err := e.ExistingAddresses.AddressesForChain(chainSelector)
		if err != nil {
			return fmt.Errorf("failed to get existing addresses: %w", err)
		}
		mcmState, err := state.MaybeLoadMCMSWithTimelockChainStateSolana(solChain, addresses)
		if err != nil {
			return fmt.Errorf("failed to load MCMS state: %w", err)
		}
		// Check if seeds are empty
		if mcmState.TimelockSeed == [32]byte{} {
			return errors.New("timelock seeds are empty, please deploy MCMS contracts first")
		}
		// Check if program IDs exists
		if mcmState.TimelockProgram.IsZero() {
			return errors.New("timelock program IDs are empty, please deploy timelock program first")
		}
		result, err := solChain.Client.GetBalance(e.GetContext(), solChain.DeployerKey.PublicKey(), rpc.CommitmentConfirmed)
		if err != nil {
			return fmt.Errorf("failed to get deployer balance: %w", err)
		}
		if result.Value < amountCfg.Amount {
			return fmt.Errorf("deployer balance is insufficient, required: %d, actual: %d", amountCfg.Amount, result.Value)
		}
	}
	return nil
}

// Apply funds the MCMS signers on each chain.
func (f TransferFromTimelock) Apply(e cldf.Environment, config TransferFromTimelockConfig) (cldf.ChangesetOutput, error) {
	timelocks := map[uint64]string{}
	proposers := map[uint64]string{}
	var batches []types.BatchOperation
	inspectors, err := proposalutils.McmsInspectors(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get MCMS inspectors: %w", err)
	}

	solChains := e.BlockChains.SolanaChains()
	for chainSelector, cfgAmounts := range config.AmountsPerChain {
		solChain := solChains[chainSelector]
		addreses, err := e.ExistingAddresses.AddressesForChain(chainSelector)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get existing addresses: %w", err)
		}
		mcmState, err := state.MaybeLoadMCMSWithTimelockChainStateSolana(solChain, addreses)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to load MCMS state: %w", err)
		}
		timelockSignerPDA := state.GetTimelockSignerPDA(mcmState.TimelockProgram, mcmState.TimelockSeed)
		timelockID := mcmssolanasdk.ContractAddress(mcmState.TimelockProgram, mcmssolanasdk.PDASeed(mcmState.TimelockSeed))
		proposerID := mcmssolanasdk.ContractAddress(mcmState.McmProgram, mcmssolanasdk.PDASeed(mcmState.ProposerMcmSeed))
		timelocks[chainSelector] = timelockID
		proposers[chainSelector] = proposerID
		ixs, err := solanachangeset.FundFromAddressIxs(
			solChain,
			timelockSignerPDA,
			[]solana.PublicKey{cfgAmounts.To},
			cfgAmounts.Amount)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to fund timelock signer on chain %d: %w", chainSelector, err)
		}

		var transactions []types.Transaction

		for _, ix := range ixs {
			solanaTx, err := mcmssolanasdk.NewTransactionFromInstruction(ix, "SystemProgram", []string{})
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to create transaction: %w", err)
			}
			transactions = append(transactions, solanaTx)
		}
		batches = append(batches, types.BatchOperation{
			ChainSelector: types.ChainSelector(chainSelector),
			Transactions:  transactions,
		})
	}
	proposal, err := proposalutils.BuildProposalFromBatchesV2(
		e,
		timelocks,
		proposers,
		inspectors,
		batches,
		"transfer funds from timelock singer",
		config.TimelockCfg,
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
	}
	return cldf.ChangesetOutput{
		MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
	}, nil
}

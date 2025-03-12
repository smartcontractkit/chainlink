package solana

import (
	"errors"
	"fmt"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
)

var _ deployment.ChangeSetV2[FundMCMSignerConfig] = FundMCMSignersChangeset{}

type AmountsToTransfer struct {
	ProposeMCM   uint64
	CancellerMCM uint64
	BypasserMCM  uint64
	Timelock     uint64
}
type FundMCMSignerConfig struct {
	AmountsPerChain map[uint64]AmountsToTransfer
}

// FundMCMSignersChangeset is a changeset that funds the MCMS signers on each chain. It will find the
// signer PDA for the proposer, canceller and bypasser MCM as well as the timelock signer PDA and send the amount of
// SOL specified in the config to each of them.
type FundMCMSignersChangeset struct{}

// VerifyPreconditions checks if the deployer has enough SOL to fund the MCMS signers on each chain.
func (f FundMCMSignersChangeset) VerifyPreconditions(e deployment.Environment, config FundMCMSignerConfig) error {
	// the number of accounts to fund per chain (bypasser, canceller, proposer, timelock)
	for chainSelector, chainCfg := range config.AmountsPerChain {
		solChain, ok := e.SolChains[chainSelector]
		if !ok {
			return fmt.Errorf("solana chain not found for selector %d", chainSelector)
		}
		addreses, err := e.ExistingAddresses.AddressesForChain(chainSelector)
		if err != nil {
			return fmt.Errorf("failed to get existing addresses: %w", err)
		}
		mcmState, err := state.MaybeLoadMCMSWithTimelockChainStateSolana(solChain, addreses)
		if err != nil {
			return fmt.Errorf("failed to load MCMS state: %w", err)
		}
		// Check if seeds are empty
		if mcmState.ProposerMcmSeed == [32]byte{} || mcmState.TimelockSeed == [32]byte{} || mcmState.CancellerMcmSeed == [32]byte{} || mcmState.BypasserMcmSeed == [32]byte{} {
			return errors.New("mcm/timelock seeds are empty, please deploy MCMS contracts first")
		}
		// Check if program IDs exists
		if mcmState.McmProgram.IsZero() || mcmState.TimelockProgram.IsZero() {
			return errors.New("mcm/timelock program IDs are empty, please deploy MCMS contracts first")
		}
		result, err := solChain.Client.GetBalance(e.GetContext(), solChain.DeployerKey.PublicKey(), rpc.CommitmentConfirmed)
		if err != nil {
			return fmt.Errorf("failed to get deployer balance: %w", err)
		}
		requiredAmount := chainCfg.ProposeMCM + chainCfg.CancellerMCM + chainCfg.BypasserMCM + chainCfg.Timelock
		if result.Value < requiredAmount {
			return fmt.Errorf("deployer balance is insufficient, required: %d, actual: %d", requiredAmount, result.Value)
		}
	}
	return nil
}

// Apply funds the MCMS signers on each chain.
func (f FundMCMSignersChangeset) Apply(e deployment.Environment, config FundMCMSignerConfig) (deployment.ChangesetOutput, error) {
	for chainSelector, cfgAmounts := range config.AmountsPerChain {
		solChain := e.SolChains[chainSelector]
		addreses, err := e.ExistingAddresses.AddressesForChain(chainSelector)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to get existing addresses: %w", err)
		}
		mcmState, err := state.MaybeLoadMCMSWithTimelockChainStateSolana(solChain, addreses)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to load MCMS state: %w", err)
		}

		err = FundFromDeployerKey(
			solChain,
			[]solana.PublicKey{state.GetTimelockSignerPDA(mcmState.TimelockProgram, mcmState.TimelockSeed)},
			cfgAmounts.Timelock)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to fund timelock signer on chain %d: %w", chainSelector, err)
		}
		err = FundFromDeployerKey(
			solChain,
			[]solana.PublicKey{state.GetMCMSignerPDA(mcmState.McmProgram, mcmState.ProposerMcmSeed)},
			cfgAmounts.ProposeMCM)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to fund MCMS proposer on chain %d: %w", chainSelector, err)
		}
		err = FundFromDeployerKey(
			solChain,
			[]solana.PublicKey{state.GetMCMSignerPDA(mcmState.McmProgram, mcmState.CancellerMcmSeed)},
			cfgAmounts.CancellerMCM)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to fund MCMS canceller on chain %d: %w", chainSelector, err)
		}
		err = FundFromDeployerKey(
			solChain,
			[]solana.PublicKey{state.GetMCMSignerPDA(mcmState.McmProgram, mcmState.BypasserMcmSeed)},
			cfgAmounts.BypasserMCM)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to fund mcm bypasser on chain %d: %w", chainSelector, err)
		}
	}
	return deployment.ChangesetOutput{}, nil
}

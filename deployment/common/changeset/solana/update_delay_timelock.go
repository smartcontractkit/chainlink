package solana

import (
	"errors"
	"fmt"
	"time"

	"github.com/gagliardetto/solana-go"

	timelockBindings "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/timelock"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
)

// UpdateTimelockDelaySolana updates the timelock delay for the given solana chains
type UpdateTimelockDelaySolana struct{}

type UpdateTimelockDelaySolanaCfg struct {
	DelayPerChain map[uint64]time.Duration
}

func (t UpdateTimelockDelaySolana) VerifyPreconditions(
	env cldf.Environment, config UpdateTimelockDelaySolanaCfg,
) error {
	if len(config.DelayPerChain) == 0 {
		return errors.New("no delay configs provided")
	}
	solanaChains := env.BlockChains.SolanaChains()
	if len(solanaChains) == 0 {
		return errors.New("no solana chains provided")
	}
	// check the timelock program is deployed
	for chainSelector := range config.DelayPerChain {
		solChain, ok := solanaChains[chainSelector]
		if !ok {
			return fmt.Errorf("solana chain not found for selector %d", chainSelector)
		}

		//nolint:staticcheck // wait till we can migrate from address book before using data store
		addresses, err := env.ExistingAddresses.AddressesForChain(chainSelector)
		if err != nil {
			return fmt.Errorf("failed to get existing addresses: %w", err)
		}
		mcmState, err := state.MaybeLoadMCMSWithTimelockChainStateSolana(solChain, addresses)
		if err != nil {
			return fmt.Errorf("failed to load MCMS state: %w", err)
		}
		if mcmState.TimelockProgram.IsZero() {
			return fmt.Errorf("timelock program not deployed for chain %d", chainSelector)
		}
		if mcmState.TimelockSeed == [32]byte{} {
			return fmt.Errorf("timelock seed not found for chain %d", chainSelector)
		}
	}
	return nil
}

func (t UpdateTimelockDelaySolana) Apply(
	env cldf.Environment, cfg UpdateTimelockDelaySolanaCfg,
) (cldf.ChangesetOutput, error) {
	solanaChains := env.BlockChains.SolanaChains()
	for chainSelector, delay := range cfg.DelayPerChain {
		solChain := solanaChains[chainSelector]
		//nolint:staticcheck // will wait till we can migrate from address book before using data store
		addresses, err := env.ExistingAddresses.AddressesForChain(chainSelector)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get existing addresses: %w", err)
		}
		mcmState, _ := state.MaybeLoadMCMSWithTimelockChainStateSolana(solChain, addresses)
		configPDA := state.GetTimelockConfigPDA(mcmState.TimelockProgram, mcmState.TimelockSeed)
		timelockBindings.SetProgramID(mcmState.TimelockProgram)
		updateDelayIx := timelockBindings.NewUpdateDelayInstruction(mcmState.TimelockSeed, uint64(delay.Seconds()), configPDA, solChain.DeployerKey.PublicKey())
		ix, err := updateDelayIx.ValidateAndBuild()
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to create update delay instruction: %w", err)
		}
		if err := solChain.Confirm([]solana.Instruction{ix}); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm instructions: %w", err)
		}
	}
	return cldf.ChangesetOutput{}, nil
}

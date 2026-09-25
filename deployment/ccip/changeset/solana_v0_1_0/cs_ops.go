package solana

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/gagliardetto/solana-go"

	cldf_solana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	solTestReceiver "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_0/test_ccip_receiver"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	solanastateview "github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/solana"
)

// use this to deploy a receiver for test
var _ cldf.ChangeSet[DeployForTestConfig] = DeployReceiverForTest

// setUpgradeAuthority creates a transaction to set the upgrade authority for a program
func setUpgradeAuthority(
	e *cldf.Environment,
	chain *cldf_solana.Chain,
	programID solana.PublicKey,
	currentUpgradeAuthority solana.PublicKey,
	newUpgradeAuthority solana.PublicKey,
	isBuffer bool,
) solana.Instruction {
	e.Logger.Infow("Setting upgrade authority", "programID", programID.String(), "currentUpgradeAuthority", currentUpgradeAuthority.String(), "newUpgradeAuthority", newUpgradeAuthority.String())
	// Buffers use the program account as the program data account
	programDataSlice := solana.NewAccountMeta(programID, true, false)
	if !isBuffer {
		// Actual program accounts use the program data account
		programDataAddress, _, _ := solana.FindProgramAddress([][]byte{programID.Bytes()}, solana.BPFLoaderUpgradeableProgramID)
		programDataSlice = solana.NewAccountMeta(programDataAddress, true, false)
	}

	keys := solana.AccountMetaSlice{
		programDataSlice, // Program account (writable)
		solana.NewAccountMeta(currentUpgradeAuthority, false, true), // Current upgrade authority (signer)
		solana.NewAccountMeta(newUpgradeAuthority, false, false),    // New upgrade authority
	}

	instruction := solana.NewInstruction(
		solana.BPFLoaderUpgradeableProgramID,
		keys,
		// https://github.com/solana-playground/solana-playground/blob/2998d4cf381aa319d26477c5d4e6d15059670a75/vscode/src/commands/deploy/bpf-upgradeable/bpf-upgradeable.ts#L72
		[]byte{4, 0, 0, 0}, // 4-byte SetAuthority instruction identifier
	)

	return instruction
}

type DeployForTestConfig struct {
	ChainSelector   uint64
	BuildConfig     *BuildSolanaConfig
	ReceiverVersion *semver.Version // leave unset to default to v1.0.0
	IsUpgrade       bool
}

func (cfg DeployForTestConfig) Validate(e cldf.Environment, state stateview.CCIPOnChainState) error {
	chainState, chainExists := state.SolChains[cfg.ChainSelector]
	if !chainExists {
		return fmt.Errorf("chain %d not found in existing state", cfg.ChainSelector)
	}
	chain := e.BlockChains.SolanaChains()[cfg.ChainSelector]

	return chainState.ValidateRouterConfig(chain)
}

func DeployReceiverForTest(e cldf.Environment, cfg DeployForTestConfig) (cldf.ChangesetOutput, error) {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	if err := cfg.Validate(e, state); err != nil {
		return cldf.ChangesetOutput{}, err
	}

	if cfg.BuildConfig != nil {
		e.Logger.Debugw("Building solana artifacts", "gitCommitSha", cfg.BuildConfig.GitCommitSha)
		err := BuildSolana(e, *cfg.BuildConfig)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build solana: %w", err)
		}
	} else {
		e.Logger.Debugw("Skipping solana build as no build config provided")
	}

	chainState := state.SolChains[cfg.ChainSelector]
	chain := e.BlockChains.SolanaChains()[cfg.ChainSelector]
	ab := cldf.NewMemoryAddressBook()
	ds := datastore.NewMemoryDataStore()

	var receiverAddress solana.PublicKey
	if !cfg.IsUpgrade {
		//nolint:gocritic // this is a false positive, we need to check if the address is zero
		if chainState.Receiver.IsZero() {
			receiverAddress, err = DeployAndMaybeSaveToAddressBook(e, chain, ab, ds, shared.Receiver, deployment.Version1_0_0, false, "")
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy program: %w", err)
			}
		} else if cfg.ReceiverVersion != nil {
			// this block is for re-deploying with a new version
			receiverAddress, err = DeployAndMaybeSaveToAddressBook(e, chain, ab, ds, shared.Receiver, *cfg.ReceiverVersion, false, "")
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy program: %w", err)
			}
		} else {
			e.Logger.Infow("Using existing receiver", "addr", chainState.Receiver.String())
			receiverAddress = chainState.Receiver
		}
		runSafely(func() {
			solTestReceiver.SetProgramID(receiverAddress)
		})
		externalExecutionConfigPDA, _, _ := solana.FindProgramAddress([][]byte{[]byte("external_execution_config")}, receiverAddress)
		instruction, ixErr := solTestReceiver.NewInitializeInstruction(
			chainState.Router,
			solanastateview.FindReceiverTargetAccount(receiverAddress),
			externalExecutionConfigPDA,
			chain.DeployerKey.PublicKey(),
			solana.SystemProgramID,
		).ValidateAndBuild()
		if ixErr != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build instruction: %w", ixErr)
		}
		if err = chain.Confirm([]solana.Instruction{instruction}); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm instructions: %w", err)
		}
	} else if cfg.IsUpgrade {
		e.Logger.Infow("Deploying new receiver", "addr", chainState.Receiver.String())
		// only support deployer key as upgrade authority. never transfer to timelock
		_, err := generateUpgradeTxns(e, chain, ab, ds, DeployChainContractsConfig{
			UpgradeConfig: UpgradeConfig{
				SpillAddress:     chain.DeployerKey.PublicKey(),
				UpgradeAuthority: chain.DeployerKey.PublicKey(),
			},
		}, cfg.ReceiverVersion, chainState.Receiver, shared.Receiver)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate upgrade txns: %w", err)
		}
	}

	return cldf.ChangesetOutput{
		AddressBook: ab,
		DataStore:   ds,
	}, nil
}

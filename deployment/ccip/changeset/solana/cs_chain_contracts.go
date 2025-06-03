package solana

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/gagliardetto/solana-go"
	cldf_solana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"

	"github.com/smartcontractkit/mcms"
	mcmsTypes "github.com/smartcontractkit/mcms/types"

	solOffRamp "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_offramp"
	solRouter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
	solFeeQuoter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/fee_quoter"
	solTestReceiver "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/test_ccip_receiver"
	solState "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	solanastateview "github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/solana"
	csState "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

// use this to set the fee aggregator
var _ cldf.ChangeSet[SetFeeAggregatorConfig] = SetFeeAggregator

// use this to update the offramp reference addresseses
var _ cldf.ChangeSet[OffRampRefAddressesConfig] = UpdateOffRampRefAddresses

// use this to set the upgrade authority of a contract
var _ cldf.ChangeSet[SetUpgradeAuthorityConfig] = SetUpgradeAuthorityChangeset

// HELPER FUNCTIONS
// GetTokenProgramID returns the program ID for the given token program name
func GetTokenProgramID(programName cldf.ContractType) (solana.PublicKey, error) {
	tokenPrograms := map[cldf.ContractType]solana.PublicKey{
		shared.SPLTokens:     solana.TokenProgramID,
		shared.SPL2022Tokens: solana.Token2022ProgramID,
	}

	programID, ok := tokenPrograms[programName]
	if !ok {
		return solana.PublicKey{}, fmt.Errorf("invalid token program: %s. Must be one of: %s, %s", programName, shared.SPLTokens, shared.SPL2022Tokens)
	}
	return programID, nil
}

// The user is not required to provide all the addresses, only the ones they want to update
type OffRampRefAddressesConfig struct {
	ChainSelector      uint64
	Router             solana.PublicKey
	FeeQuoter          solana.PublicKey
	AddressLookupTable solana.PublicKey
	RMNRemote          solana.PublicKey
	MCMS               *proposalutils.TimelockConfig
}

func (cfg OffRampRefAddressesConfig) Validate(e cldf.Environment, state stateview.CCIPOnChainState) error {
	chain := e.BlockChains.SolanaChains()[cfg.ChainSelector]
	chainState, chainExists := state.SolChains[chain.Selector]
	if !chainExists {
		return fmt.Errorf("chain %s not found in existing state, deploy the link token first", chain.String())
	}
	return ValidateMCMSConfigSolana(e, cfg.MCMS, chain, chainState, solana.PublicKey{}, "", map[cldf.ContractType]bool{shared.OffRamp: true})
}

func UpdateOffRampRefAddresses(
	e cldf.Environment,
	config OffRampRefAddressesConfig,
) (cldf.ChangesetOutput, error) {
	state, err := stateview.LoadOnchainStateSolana(e)
	if err != nil {
		e.Logger.Errorw("Failed to load existing onchain state", "err", err)
		return cldf.ChangesetOutput{}, err
	}
	if err := config.Validate(e, state); err != nil {
		return cldf.ChangesetOutput{}, err
	}
	chain := e.BlockChains.SolanaChains()[config.ChainSelector]

	chainState, chainExists := state.SolChains[chain.Selector]
	if !chainExists {
		return cldf.ChangesetOutput{}, fmt.Errorf("chain %s not found in existing state, deploy the link token first", chain.String())
	}
	if chainState.OffRamp.IsZero() {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get offramp address for chain %s", chain.String())
	}

	var referenceAddressesAccount solOffRamp.ReferenceAddresses
	offRampReferenceAddressesPDA, _, _ := solState.FindOfframpReferenceAddressesPDA(chainState.OffRamp)
	if err = chain.GetAccountDataBorshInto(e.GetContext(), offRampReferenceAddressesPDA, &referenceAddressesAccount); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get offramp reference addresses: %w", err)
	}
	routerToSet := referenceAddressesAccount.Router
	if !config.Router.IsZero() {
		e.Logger.Infof("setting router on offramp to %s", config.Router.String())
		routerToSet = config.Router
	}
	feeQuoterToSet := referenceAddressesAccount.FeeQuoter
	if !config.FeeQuoter.IsZero() {
		e.Logger.Infof("setting fee quoter on offramp to %s", config.FeeQuoter.String())
		feeQuoterToSet = config.FeeQuoter
	}
	addressLookupTableToSet := referenceAddressesAccount.OfframpLookupTable
	if !config.AddressLookupTable.IsZero() {
		e.Logger.Infof("setting address lookup table on offramp to %s", config.AddressLookupTable.String())
		addressLookupTableToSet = config.AddressLookupTable
	}
	rmnRemoteToSet := referenceAddressesAccount.RmnRemote
	if !config.RMNRemote.IsZero() {
		e.Logger.Infof("setting rmn remote on offramp to %s", config.RMNRemote.String())
		rmnRemoteToSet = config.RMNRemote
	}

	offRampUsingMCMS := solanastateview.IsSolanaProgramOwnedByTimelock(
		&e,
		chain,
		chainState,
		shared.OffRamp,
		solana.PublicKey{},
		"")
	authority := GetAuthorityForIxn(
		&e,
		chain,
		chainState,
		config.MCMS,
		shared.OffRamp,
		solana.PublicKey{},
		"",
	)
	solOffRamp.SetProgramID(chainState.OffRamp)
	ix, err := solOffRamp.NewUpdateReferenceAddressesInstruction(
		routerToSet,
		feeQuoterToSet,
		addressLookupTableToSet,
		rmnRemoteToSet,
		chainState.OffRampConfigPDA,
		offRampReferenceAddressesPDA,
		authority,
	).ValidateAndBuild()
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to build instruction: %w", err)
	}

	if offRampUsingMCMS {
		tx, err := BuildMCMSTxn(ix, chainState.OffRamp.String(), shared.OffRamp)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to create transaction: %w", err)
		}
		proposal, err := BuildProposalsForTxns(
			e, config.ChainSelector, "proposal to UpdateOffRampRefAddresses in Solana", config.MCMS.MinDelay, []mcmsTypes.Transaction{*tx})
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return cldf.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	}

	if err := chain.Confirm([]solana.Instruction{ix}); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm instructions: %w", err)
	}
	return cldf.ChangesetOutput{}, nil
}

type SetUpgradeAuthorityConfig struct {
	ChainSelector         uint64
	NewUpgradeAuthority   solana.PublicKey
	SetAfterInitialDeploy bool                          // set all of the programs after the initial deploy
	SetOffRamp            bool                          // offramp not upgraded in place, so may need to set separately
	SetMCMSPrograms       bool                          // these all deploy at once so just set them all
	TransferKeys          []solana.PublicKey            // any keys not covered by the above e.g. partner programs
	MCMS                  *proposalutils.TimelockConfig // if set, assumes current upgrade authority is the timelock
}

func SetUpgradeAuthorityChangeset(
	e cldf.Environment,
	config SetUpgradeAuthorityConfig,
) (cldf.ChangesetOutput, error) {
	chain := e.BlockChains.SolanaChains()[config.ChainSelector]
	state, err := stateview.LoadOnchainStateSolana(e)
	if err != nil {
		e.Logger.Errorw("Failed to load existing onchain state", "err", err)
		return cldf.ChangesetOutput{}, err
	}
	chainState, chainExists := state.SolChains[chain.Selector]
	if !chainExists {
		return cldf.ChangesetOutput{}, fmt.Errorf("chain %s not found in existing state, deploy the link token first", chain.String())
	}
	programs := make([]solana.PublicKey, 0)
	if config.SetAfterInitialDeploy {
		programs = append(programs, chainState.Router, chainState.FeeQuoter, chainState.RMNRemote, chainState.BurnMintTokenPools[shared.CLLMetadata], chainState.LockReleaseTokenPools[shared.CLLMetadata])
	}
	if config.SetOffRamp {
		programs = append(programs, chainState.OffRamp)
	}
	if config.SetMCMSPrograms {
		addresses, err := e.ExistingAddresses.AddressesForChain(config.ChainSelector)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get existing addresses: %w", err)
		}
		mcmState, err := csState.MaybeLoadMCMSWithTimelockChainStateSolana(chain, addresses)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
		}
		programs = append(programs, mcmState.AccessControllerProgram, mcmState.TimelockProgram, mcmState.McmProgram)
	}
	for _, transfer := range config.TransferKeys {
		if transfer.IsZero() {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get program address for chain %s", chain.String())
		}
		programs = append(programs, transfer)
	}
	// We do two loops here just to catch any errors before we get partway through the process
	for _, program := range programs {
		if program.IsZero() {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get program address for chain %s", chain.String())
		}
	}
	currentAuthority := chain.DeployerKey.PublicKey()
	if config.MCMS != nil {
		timelockSignerPDA, err := FetchTimelockSigner(e, chain.Selector)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get timelock signer: %w", err)
		}
		currentAuthority = timelockSignerPDA
	}
	e.Logger.Infow("Setting upgrade authority", "newUpgradeAuthority", config.NewUpgradeAuthority.String())
	mcmsTxns := make([]mcmsTypes.Transaction, 0)
	for _, programID := range programs {
		ixn := setUpgradeAuthority(&e, &chain, programID, currentAuthority, config.NewUpgradeAuthority, false)
		if config.MCMS == nil {
			if err := chain.Confirm([]solana.Instruction{ixn}); err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm instructions: %w", err)
			}
		} else {
			tx, err := BuildMCMSTxn(
				ixn,
				solana.BPFLoaderUpgradeableProgramID.String(),
				cldf.ContractType(solana.BPFLoaderUpgradeableProgramID.String()))
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to create transaction: %w", err)
			}
			mcmsTxns = append(mcmsTxns, *tx)
		}
	}
	if len(mcmsTxns) > 0 {
		proposal, err := BuildProposalsForTxns(
			e, config.ChainSelector, "proposal to SetUpgradeAuthority in Solana", config.MCMS.MinDelay, mcmsTxns)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return cldf.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	}
	return cldf.ChangesetOutput{}, nil
}

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

type SetFeeAggregatorConfig struct {
	ChainSelector uint64
	FeeAggregator string
	MCMS          *proposalutils.TimelockConfig
}

func (cfg SetFeeAggregatorConfig) Validate(e cldf.Environment, state stateview.CCIPOnChainState) error {
	chainState, chainExists := state.SolChains[cfg.ChainSelector]
	if !chainExists {
		return fmt.Errorf("chain %d not found in existing state", cfg.ChainSelector)
	}
	chain := e.BlockChains.SolanaChains()[cfg.ChainSelector]

	if err := chainState.ValidateRouterConfig(chain); err != nil {
		return err
	}

	if err := ValidateMCMSConfigSolana(e, cfg.MCMS, chain, chainState, solana.PublicKey{}, "", map[cldf.ContractType]bool{shared.Router: true}); err != nil {
		return err
	}

	// Validate fee aggregator address is valid
	if _, err := solana.PublicKeyFromBase58(cfg.FeeAggregator); err != nil {
		return fmt.Errorf("invalid fee aggregator address: %w", err)
	}

	if solana.MustPublicKeyFromBase58(cfg.FeeAggregator).IsZero() {
		return errors.New("fee aggregator address cannot be zero")
	}

	if chainState.GetFeeAggregator(chain).Equals(solana.MustPublicKeyFromBase58(cfg.FeeAggregator)) {
		return fmt.Errorf("fee aggregator %s is already set on chain %d", cfg.FeeAggregator, cfg.ChainSelector)
	}

	return nil
}

func SetFeeAggregator(e cldf.Environment, cfg SetFeeAggregatorConfig) (cldf.ChangesetOutput, error) {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	if err := cfg.Validate(e, state); err != nil {
		return cldf.ChangesetOutput{}, err
	}
	chainState := state.SolChains[cfg.ChainSelector]
	chain := e.BlockChains.SolanaChains()[cfg.ChainSelector]

	feeAggregatorPubKey := solana.MustPublicKeyFromBase58(cfg.FeeAggregator)
	routerConfigPDA, _, _ := solState.FindConfigPDA(chainState.Router)
	routerUsingMCMS := solanastateview.IsSolanaProgramOwnedByTimelock(
		&e,
		chain,
		chainState,
		shared.Router,
		solana.PublicKey{},
		"")

	solRouter.SetProgramID(chainState.Router)
	authority := GetAuthorityForIxn(
		&e,
		chain,
		chainState,
		cfg.MCMS,
		shared.Router,
		solana.PublicKey{},
		"",
	)
	instruction, err := solRouter.NewUpdateFeeAggregatorInstruction(
		feeAggregatorPubKey,
		routerConfigPDA,
		authority,
		solana.SystemProgramID,
	).ValidateAndBuild()
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to build instruction: %w", err)
	}

	if routerUsingMCMS {
		tx, err := BuildMCMSTxn(instruction, chainState.Router.String(), shared.Router)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to create transaction: %w", err)
		}
		proposal, err := BuildProposalsForTxns(
			e, cfg.ChainSelector, "proposal to SetFeeAggregator in Solana", cfg.MCMS.MinDelay, []mcmsTypes.Transaction{*tx})
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return cldf.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	}

	if err := chain.Confirm([]solana.Instruction{instruction}); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm instructions: %w", err)
	}
	e.Logger.Infow("Set new fee aggregator", "chain", chain.String(), "fee_aggregator", feeAggregatorPubKey.String())

	return cldf.ChangesetOutput{}, nil
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

	var receiverAddress solana.PublicKey
	if !cfg.IsUpgrade {
		//nolint:gocritic // this is a false positive, we need to check if the address is zero
		if chainState.Receiver.IsZero() {
			receiverAddress, err = DeployAndMaybeSaveToAddressBook(e, chain, ab, shared.Receiver, deployment.Version1_0_0, false, "")
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy program: %w", err)
			}
		} else if cfg.ReceiverVersion != nil {
			// this block is for re-deploying with a new version
			receiverAddress, err = DeployAndMaybeSaveToAddressBook(e, chain, ab, shared.Receiver, *cfg.ReceiverVersion, false, "")
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy program: %w", err)
			}
		} else {
			e.Logger.Infow("Using existing receiver", "addr", chainState.Receiver.String())
			receiverAddress = chainState.Receiver
		}
		solTestReceiver.SetProgramID(receiverAddress)
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
		receiverAddress = chainState.Receiver
		// only support deployer key as upgrade authority. never transfer to timelock
		_, err := generateUpgradeTxns(e, chain, ab, DeployChainContractsConfig{
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
	}, nil
}

type SetLinkTokenConfig struct {
	ChainSelector uint64
}

func (cfg SetLinkTokenConfig) Validate(e cldf.Environment, state stateview.CCIPOnChainState) error {
	chainState, chainExists := state.SolChains[cfg.ChainSelector]
	if !chainExists {
		return fmt.Errorf("chain %d not found in existing state", cfg.ChainSelector)
	}
	chain := e.BlockChains.SolanaChains()[cfg.ChainSelector]

	return chainState.ValidateRouterConfig(chain)
}

func SetLinkToken(e cldf.Environment, cfg SetLinkTokenConfig) (cldf.ChangesetOutput, error) {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	if err := cfg.Validate(e, state); err != nil {
		return cldf.ChangesetOutput{}, err
	}
	chainState := state.SolChains[cfg.ChainSelector]
	chain := e.BlockChains.SolanaChains()[cfg.ChainSelector]
	routerConfigPDA, _, _ := solState.FindConfigPDA(chainState.Router)
	feeQuoterConfigPDA, _, _ := solState.FindConfigPDA(chainState.FeeQuoter)

	solRouter.SetProgramID(chainState.Router)
	routerIx, err := solRouter.NewSetLinkTokenMintInstruction(
		chainState.LinkToken,
		routerConfigPDA,
		chain.DeployerKey.PublicKey(),
		solana.SystemProgramID,
	).ValidateAndBuild()
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to build instruction: %w", err)
	}

	solFeeQuoter.SetProgramID(chainState.FeeQuoter)
	feeQuoterIx, err := solFeeQuoter.NewSetLinkTokenMintInstruction(
		feeQuoterConfigPDA,
		chainState.LinkToken,
		chain.DeployerKey.PublicKey(),
	).ValidateAndBuild()
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to build instruction: %w", err)
	}

	if err := chain.Confirm([]solana.Instruction{routerIx, feeQuoterIx}); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm instructions: %w", err)
	}
	return cldf.ChangesetOutput{}, nil
}

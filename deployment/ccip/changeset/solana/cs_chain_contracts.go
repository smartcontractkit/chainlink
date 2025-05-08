package solana

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/gagliardetto/solana-go"

	"github.com/smartcontractkit/mcms"
	mcmsTypes "github.com/smartcontractkit/mcms/types"

	solOffRamp "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_offramp"
	solRouter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
	solFeeQuoter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/fee_quoter"
	solTestReceiver "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/test_ccip_receiver"
	solCommonUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	solState "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"

	"github.com/smartcontractkit/chainlink/deployment"
	ccipChangeset "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	csState "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

// use this to set the fee aggregator
var _ deployment.ChangeSet[SetFeeAggregatorConfig] = SetFeeAggregator

// use this to update the offramp reference addresseses
var _ deployment.ChangeSet[OffRampRefAddressesConfig] = UpdateOffRampRefAddresses

// use this to set the upgrade authority of a contract
var _ deployment.ChangeSet[SetUpgradeAuthorityConfig] = SetUpgradeAuthorityChangeset

type MCMSConfigSolana struct {
	MCMS *proposalutils.TimelockConfig
	// Public key of program authorities. Depending on when this changeset is called, some may be under
	// the control of the deployer, and some may be under the control of the timelock. (e.g. during new offramp deploy)
	RouterOwnedByTimelock    bool
	FeeQuoterOwnedByTimelock bool
	OffRampOwnedByTimelock   bool
	RMNRemoteOwnedByTimelock bool
	// Operates as a set. Token Pool configs will owned by timelock per token (the key)
	BurnMintTokenPoolOwnedByTimelock    map[solana.PublicKey]bool
	LockReleaseTokenPoolOwnedByTimelock map[solana.PublicKey]bool
}

// HELPER FUNCTIONS
// GetTokenProgramID returns the program ID for the given token program name
func GetTokenProgramID(programName deployment.ContractType) (solana.PublicKey, error) {
	tokenPrograms := map[deployment.ContractType]solana.PublicKey{
		ccipChangeset.SPLTokens:     solana.TokenProgramID,
		ccipChangeset.SPL2022Tokens: solana.Token2022ProgramID,
	}

	programID, ok := tokenPrograms[programName]
	if !ok {
		return solana.PublicKey{}, fmt.Errorf("invalid token program: %s. Must be one of: %s, %s", programName, ccipChangeset.SPLTokens, ccipChangeset.SPL2022Tokens)
	}
	return programID, nil
}

func commonValidation(e deployment.Environment, selector uint64, tokenPubKey solana.PublicKey) error {
	chain, ok := e.SolChains[selector]
	if !ok {
		return fmt.Errorf("chain selector %d not found in environment", selector)
	}
	state, err := ccipChangeset.LoadOnchainState(e)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}
	chainState, chainExists := state.SolChains[selector]
	if !chainExists {
		return fmt.Errorf("chain %s not found in existing state, deploy the link token first", chain.String())
	}
	if tokenPubKey.Equals(chainState.LinkToken) || tokenPubKey.Equals(chainState.WSOL) {
		return nil
	}
	exists := false
	allTokens := chainState.SPL2022Tokens
	allTokens = append(allTokens, chainState.SPLTokens...)
	for _, token := range allTokens {
		if token.Equals(tokenPubKey) {
			exists = true
			break
		}
	}
	if !exists {
		return fmt.Errorf("token %s not found in existing state, deploy the token first", tokenPubKey.String())
	}
	return nil
}

func validateRouterConfig(chain deployment.SolChain, chainState ccipChangeset.SolCCIPChainState) error {
	_, routerConfigPDA, err := chainState.GetRouterInfo()
	if err != nil {
		return err
	}
	var routerConfigAccount solRouter.Config
	err = chain.GetAccountDataBorshInto(context.Background(), routerConfigPDA, &routerConfigAccount)
	if err != nil {
		return fmt.Errorf("router config not found in existing state, initialize the router first %d", chain.Selector)
	}
	return nil
}

func validateFeeAggregatorConfig(chain deployment.SolChain, chainState ccipChangeset.SolCCIPChainState) error {
	if chainState.GetFeeAggregator(chain).IsZero() {
		return fmt.Errorf("fee aggregator not found in existing state, set the fee aggregator first for chain %d", chain.Selector)
	}
	return nil
}

func validateFeeQuoterConfig(chain deployment.SolChain, chainState ccipChangeset.SolCCIPChainState) error {
	if chainState.FeeQuoter.IsZero() {
		return fmt.Errorf("fee quoter not found in existing state, deploy the fee quoter first for chain %d", chain.Selector)
	}
	var fqConfig solFeeQuoter.Config
	feeQuoterConfigPDA, _, _ := solState.FindFqConfigPDA(chainState.FeeQuoter)
	err := chain.GetAccountDataBorshInto(context.Background(), feeQuoterConfigPDA, &fqConfig)
	if err != nil {
		return fmt.Errorf("fee quoter config not found in existing state, initialize the fee quoter first %d", chain.Selector)
	}
	return nil
}

func validateOffRampConfig(chain deployment.SolChain, chainState ccipChangeset.SolCCIPChainState) error {
	if chainState.OffRamp.IsZero() {
		return fmt.Errorf("offramp not found in existing state, deploy the offramp first for chain %d", chain.Selector)
	}
	var offRampConfig solOffRamp.Config
	offRampConfigPDA, _, _ := solState.FindOfframpConfigPDA(chainState.OffRamp)
	err := chain.GetAccountDataBorshInto(context.Background(), offRampConfigPDA, &offRampConfig)
	if err != nil {
		return fmt.Errorf("offramp config not found in existing state, initialize the offramp first %d", chain.Selector)
	}
	return nil
}

// The user is not required to provide all the addresses, only the ones they want to update
type OffRampRefAddressesConfig struct {
	ChainSelector      uint64
	Router             solana.PublicKey
	FeeQuoter          solana.PublicKey
	AddressLookupTable solana.PublicKey
	RMNRemote          solana.PublicKey
	MCMSSolana         *MCMSConfigSolana
}

func (cfg OffRampRefAddressesConfig) Validate(e deployment.Environment) error {
	chain := e.SolChains[cfg.ChainSelector]
	state, err := ccipChangeset.LoadOnchainState(e)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}
	chainState, chainExists := state.SolChains[chain.Selector]
	if !chainExists {
		return fmt.Errorf("chain %s not found in existing state, deploy the link token first", chain.String())
	}
	return ValidateMCMSConfigSolana(e, cfg.MCMSSolana, chain, chainState, solana.PublicKey{})
}

func UpdateOffRampRefAddresses(
	e deployment.Environment,
	config OffRampRefAddressesConfig,
) (deployment.ChangesetOutput, error) {
	state, err := ccipChangeset.LoadOnchainStateSolana(e)
	chain := e.SolChains[config.ChainSelector]
	if err != nil {
		e.Logger.Errorw("Failed to load existing onchain state", "err", err)
		return deployment.ChangesetOutput{}, err
	}
	chainState, chainExists := state.SolChains[chain.Selector]
	if !chainExists {
		return deployment.ChangesetOutput{}, fmt.Errorf("chain %s not found in existing state, deploy the link token first", chain.String())
	}
	if chainState.OffRamp.IsZero() {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get offramp address for chain %s", chain.String())
	}

	var referenceAddressesAccount solOffRamp.ReferenceAddresses
	offRampReferenceAddressesPDA, _, _ := solState.FindOfframpReferenceAddressesPDA(chainState.OffRamp)
	if err = chain.GetAccountDataBorshInto(e.GetContext(), offRampReferenceAddressesPDA, &referenceAddressesAccount); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get offramp reference addresses: %w", err)
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
	if err := ValidateMCMSConfigSolana(e, config.MCMSSolana, chain, chainState, solana.PublicKey{}); err != nil {
		return deployment.ChangesetOutput{}, err
	}

	offRampUsingMCMS := config.MCMSSolana != nil && config.MCMSSolana.OffRampOwnedByTimelock
	authority, err := GetAuthorityForIxn(
		&e,
		chain,
		config.MCMSSolana,
		ccipChangeset.OffRamp,
		solana.PublicKey{})
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get authority for ixn: %w", err)
	}
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
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to build instruction: %w", err)
	}

	if offRampUsingMCMS {
		tx, err := BuildMCMSTxn(ix, chainState.OffRamp.String(), ccipChangeset.OffRamp)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to create transaction: %w", err)
		}
		proposal, err := BuildProposalsForTxns(
			e, config.ChainSelector, "proposal to UpdateOffRampRefAddresses in Solana", config.MCMSSolana.MCMS.MinDelay, []mcmsTypes.Transaction{*tx})
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return deployment.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	}

	if err := chain.Confirm([]solana.Instruction{ix}); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to confirm instructions: %w", err)
	}
	return deployment.ChangesetOutput{}, nil
}

type SetUpgradeAuthorityConfig struct {
	ChainSelector         uint64
	NewUpgradeAuthority   solana.PublicKey
	SetAfterInitialDeploy bool // set all of the programs after the initial deploy
	SetOffRamp            bool // offramp not upgraded in place, so may need to set separately
	SetMCMSPrograms       bool // these all deploy at once so just set them all
}

func SetUpgradeAuthorityChangeset(
	e deployment.Environment,
	config SetUpgradeAuthorityConfig,
) (deployment.ChangesetOutput, error) {
	chain := e.SolChains[config.ChainSelector]
	state, err := ccipChangeset.LoadOnchainStateSolana(e)
	if err != nil {
		e.Logger.Errorw("Failed to load existing onchain state", "err", err)
		return deployment.ChangesetOutput{}, err
	}
	chainState, chainExists := state.SolChains[chain.Selector]
	if !chainExists {
		return deployment.ChangesetOutput{}, fmt.Errorf("chain %s not found in existing state, deploy the link token first", chain.String())
	}
	programs := make([]solana.PublicKey, 0)
	if config.SetAfterInitialDeploy {
		programs = append(programs, chainState.Router, chainState.FeeQuoter, chainState.RMNRemote, chainState.BurnMintTokenPool, chainState.LockReleaseTokenPool)
	}
	if config.SetOffRamp {
		programs = append(programs, chainState.OffRamp)
	}
	if config.SetMCMSPrograms {
		addresses, err := e.ExistingAddresses.AddressesForChain(config.ChainSelector)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to get existing addresses: %w", err)
		}
		mcmState, err := csState.MaybeLoadMCMSWithTimelockChainStateSolana(chain, addresses)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
		}
		programs = append(programs, mcmState.AccessControllerProgram, mcmState.TimelockProgram, mcmState.McmProgram)
	}
	// We do two loops here just to catch any errors before we get partway through the process
	for _, program := range programs {
		if program.IsZero() {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to get program address for chain %s", chain.String())
		}
	}
	e.Logger.Infow("Setting upgrade authority", "newUpgradeAuthority", config.NewUpgradeAuthority.String())
	for _, programID := range programs {
		if err := setUpgradeAuthority(&e, &chain, programID, chain.DeployerKey, &config.NewUpgradeAuthority, false); err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to set upgrade authority: %w", err)
		}
	}
	return deployment.ChangesetOutput{}, nil
}

// setUpgradeAuthority creates a transaction to set the upgrade authority for a program
func setUpgradeAuthority(
	e *deployment.Environment,
	chain *deployment.SolChain,
	programID solana.PublicKey,
	currentUpgradeAuthority *solana.PrivateKey,
	newUpgradeAuthority *solana.PublicKey,
	isBuffer bool,
) error {
	// Buffers use the program account as the program data account
	programDataSlice := solana.NewAccountMeta(programID, true, false)
	if !isBuffer {
		// Actual program accounts use the program data account
		programDataAddress, _, _ := solana.FindProgramAddress([][]byte{programID.Bytes()}, solana.BPFLoaderUpgradeableProgramID)
		programDataSlice = solana.NewAccountMeta(programDataAddress, true, false)
	}

	keys := solana.AccountMetaSlice{
		programDataSlice, // Program account (writable)
		solana.NewAccountMeta(currentUpgradeAuthority.PublicKey(), false, true), // Current upgrade authority (signer)
		solana.NewAccountMeta(*newUpgradeAuthority, false, false),               // New upgrade authority
	}

	instruction := solana.NewInstruction(
		solana.BPFLoaderUpgradeableProgramID,
		keys,
		// https://github.com/solana-playground/solana-playground/blob/2998d4cf381aa319d26477c5d4e6d15059670a75/vscode/src/commands/deploy/bpf-upgradeable/bpf-upgradeable.ts#L72
		[]byte{4, 0, 0, 0}, // 4-byte SetAuthority instruction identifier
	)

	if err := chain.Confirm([]solana.Instruction{instruction}, solCommonUtil.AddSigners(*currentUpgradeAuthority)); err != nil {
		return fmt.Errorf("failed to confirm setUpgradeAuthority: %w", err)
	}
	e.Logger.Infow("Set upgrade authority", "programID", programID.String(), "newUpgradeAuthority", newUpgradeAuthority.String())

	return nil
}

type SetFeeAggregatorConfig struct {
	ChainSelector uint64
	FeeAggregator string
	MCMSSolana    *MCMSConfigSolana
}

func (cfg SetFeeAggregatorConfig) Validate(e deployment.Environment) error {
	state, err := ccipChangeset.LoadOnchainState(e)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}
	chainState, chainExists := state.SolChains[cfg.ChainSelector]
	if !chainExists {
		return fmt.Errorf("chain %d not found in existing state", cfg.ChainSelector)
	}
	chain := e.SolChains[cfg.ChainSelector]

	if err := validateRouterConfig(chain, chainState); err != nil {
		return err
	}

	if err := ValidateMCMSConfigSolana(e, cfg.MCMSSolana, chain, chainState, solana.PublicKey{}); err != nil {
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

func SetFeeAggregator(e deployment.Environment, cfg SetFeeAggregatorConfig) (deployment.ChangesetOutput, error) {
	if err := cfg.Validate(e); err != nil {
		return deployment.ChangesetOutput{}, err
	}

	state, _ := ccipChangeset.LoadOnchainState(e)
	chainState := state.SolChains[cfg.ChainSelector]
	chain := e.SolChains[cfg.ChainSelector]

	feeAggregatorPubKey := solana.MustPublicKeyFromBase58(cfg.FeeAggregator)
	routerConfigPDA, _, _ := solState.FindConfigPDA(chainState.Router)
	routerUsingMCMS := cfg.MCMSSolana != nil && cfg.MCMSSolana.RouterOwnedByTimelock

	solRouter.SetProgramID(chainState.Router)
	authority, err := GetAuthorityForIxn(
		&e,
		chain,
		cfg.MCMSSolana,
		ccipChangeset.Router,
		solana.PublicKey{})
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get authority for ixn: %w", err)
	}
	instruction, err := solRouter.NewUpdateFeeAggregatorInstruction(
		feeAggregatorPubKey,
		routerConfigPDA,
		authority,
		solana.SystemProgramID,
	).ValidateAndBuild()
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to build instruction: %w", err)
	}

	if routerUsingMCMS {
		tx, err := BuildMCMSTxn(instruction, chainState.Router.String(), ccipChangeset.Router)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to create transaction: %w", err)
		}
		proposal, err := BuildProposalsForTxns(
			e, cfg.ChainSelector, "proposal to SetFeeAggregator in Solana", cfg.MCMSSolana.MCMS.MinDelay, []mcmsTypes.Transaction{*tx})
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return deployment.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	}

	if err := chain.Confirm([]solana.Instruction{instruction}); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to confirm instructions: %w", err)
	}
	e.Logger.Infow("Set new fee aggregator", "chain", chain.String(), "fee_aggregator", feeAggregatorPubKey.String())

	return deployment.ChangesetOutput{}, nil
}

type DeployForTestConfig struct {
	ChainSelector   uint64
	BuildConfig     *BuildSolanaConfig
	ReceiverVersion *semver.Version // leave unset to default to v1.0.0
	IsUpgrade       bool
}

func (cfg DeployForTestConfig) Validate(e deployment.Environment) error {
	state, err := ccipChangeset.LoadOnchainState(e)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}
	chainState, chainExists := state.SolChains[cfg.ChainSelector]
	if !chainExists {
		return fmt.Errorf("chain %d not found in existing state", cfg.ChainSelector)
	}
	chain := e.SolChains[cfg.ChainSelector]

	return validateRouterConfig(chain, chainState)
}

func DeployReceiverForTest(e deployment.Environment, cfg DeployForTestConfig) (deployment.ChangesetOutput, error) {
	if err := cfg.Validate(e); err != nil {
		return deployment.ChangesetOutput{}, err
	}

	if cfg.BuildConfig != nil {
		e.Logger.Debugw("Building solana artifacts", "gitCommitSha", cfg.BuildConfig.GitCommitSha)
		err := BuildSolana(e, *cfg.BuildConfig)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to build solana: %w", err)
		}
	} else {
		e.Logger.Debugw("Skipping solana build as no build config provided")
	}

	state, _ := ccipChangeset.LoadOnchainState(e)
	chainState := state.SolChains[cfg.ChainSelector]
	chain := e.SolChains[cfg.ChainSelector]
	ab := deployment.NewMemoryAddressBook()

	var receiverAddress solana.PublicKey
	var err error
	if !cfg.IsUpgrade {
		//nolint:gocritic // this is a false positive, we need to check if the address is zero
		if chainState.Receiver.IsZero() {
			receiverAddress, err = DeployAndMaybeSaveToAddressBook(e, chain, ab, ccipChangeset.Receiver, deployment.Version1_0_0, false)
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to deploy program: %w", err)
			}
		} else if cfg.ReceiverVersion != nil {
			// this block is for re-deploying with a new version
			receiverAddress, err = DeployAndMaybeSaveToAddressBook(e, chain, ab, ccipChangeset.Receiver, *cfg.ReceiverVersion, false)
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to deploy program: %w", err)
			}
		} else {
			e.Logger.Infow("Using existing receiver", "addr", chainState.Receiver.String())
			receiverAddress = chainState.Receiver
		}
		solTestReceiver.SetProgramID(receiverAddress)
		externalExecutionConfigPDA, _, _ := solana.FindProgramAddress([][]byte{[]byte("external_execution_config")}, receiverAddress)
		instruction, ixErr := solTestReceiver.NewInitializeInstruction(
			chainState.Router,
			ccipChangeset.FindReceiverTargetAccount(receiverAddress),
			externalExecutionConfigPDA,
			chain.DeployerKey.PublicKey(),
			solana.SystemProgramID,
		).ValidateAndBuild()
		if ixErr != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to build instruction: %w", ixErr)
		}
		if err = chain.Confirm([]solana.Instruction{instruction}); err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to confirm instructions: %w", err)
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
		}, cfg.ReceiverVersion, chainState.Receiver, ccipChangeset.Receiver)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to generate upgrade txns: %w", err)
		}
	}

	return deployment.ChangesetOutput{
		AddressBook: ab,
	}, nil
}

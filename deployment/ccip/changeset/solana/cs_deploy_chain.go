package solana

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"

	"github.com/Masterminds/semver/v3"
	"github.com/gagliardetto/solana-go"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/mcms"
	"github.com/smartcontractkit/mcms/sdk"
	mcmsSolana "github.com/smartcontractkit/mcms/sdk/solana"
	mcmsTypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink/deployment"
	ccipChangeset "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/common/types"

	solBinary "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go/rpc"
	solRpc "github.com/gagliardetto/solana-go/rpc"

	solOffRamp "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_offramp"
	solRouter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
	solFeeQuoter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/fee_quoter"
	solRmnRemote "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/rmn_remote"
	solTestReceiver "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/test_ccip_receiver"
	solCommonUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	solState "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"
	solanaMCMS "github.com/smartcontractkit/chainlink/deployment/common/changeset/solana/mcms"
)

var _ deployment.ChangeSet[DeployChainContractsConfig] = DeployChainContractsChangeset

func getTypeToProgramDeployName() map[deployment.ContractType]string {
	return map[deployment.ContractType]string{
		ccipChangeset.Router:               deployment.RouterProgramName,
		ccipChangeset.OffRamp:              deployment.OffRampProgramName,
		ccipChangeset.FeeQuoter:            deployment.FeeQuoterProgramName,
		ccipChangeset.BurnMintTokenPool:    deployment.BurnMintTokenPoolProgramName,
		ccipChangeset.LockReleaseTokenPool: deployment.LockReleaseTokenPoolProgramName,
		ccipChangeset.RMNRemote:            deployment.RMNRemoteProgramName,
		types.AccessControllerProgram:      deployment.AccessControllerProgramName,
		types.ManyChainMultisigProgram:     deployment.McmProgramName,
		types.RBACTimelockProgram:          deployment.TimelockProgramName,
		ccipChangeset.Receiver:             deployment.ReceiverProgramName,
	}
}

type DeployChainContractsConfig struct {
	HomeChainSelector      uint64
	ChainSelector          uint64
	ContractParamsPerChain ChainContractParams
	UpgradeConfig          UpgradeConfig
	BuildConfig            *BuildSolanaConfig
	// TODO: add validation for this
	MCMSWithTimelockConfig types.MCMSWithTimelockConfigV2
}

type ChainContractParams struct {
	FeeQuoterParams FeeQuoterParams
	OffRampParams   OffRampParams
}

type FeeQuoterParams struct {
	DefaultMaxFeeJuelsPerMsg solBinary.Uint128
	BillingConfig            []solFeeQuoter.BillingTokenConfig
}

type OffRampParams struct {
	EnableExecutionAfter int64
}
type UpgradeConfig struct {
	NewFeeQuoterVersion            *semver.Version
	NewRouterVersion               *semver.Version
	NewRMNRemoteVersion            *semver.Version
	NewBurnMintTokenPoolVersion    *semver.Version
	NewLockReleaseTokenPoolVersion *semver.Version
	NewAccessControllerVersion     *semver.Version
	NewMCMVersion                  *semver.Version
	NewTimelockVersion             *semver.Version
	// Offramp is redeployed with the existing deployer key while the other programs are upgraded in place
	NewOffRampVersion *semver.Version
	// SpillAddress and UpgradeAuthority must be set
	SpillAddress     solana.PublicKey
	UpgradeAuthority solana.PublicKey
	// MCMS config must be set for upgrades and offramp redploys (to configure the fee quoter after redeploy)
	MCMS *ccipChangeset.MCMSConfig
}

func (cfg UpgradeConfig) Validate(e deployment.Environment, chainSelector uint64) error {
	if cfg.NewFeeQuoterVersion == nil && cfg.NewRouterVersion == nil && cfg.NewOffRampVersion == nil {
		return nil
	}
	if cfg.MCMS == nil {
		return errors.New("MCMS config must be set for upgrades")
	}
	if cfg.SpillAddress.IsZero() {
		return errors.New("spill address must be set for fee quoter and router upgrades")
	}
	if cfg.UpgradeAuthority.IsZero() {
		return errors.New("upgrade authority must be set for fee quoter and router upgrades")
	}
	return ValidateMCMSConfig(e, chainSelector, cfg.MCMS)
}

func (c DeployChainContractsConfig) Validate(e deployment.Environment) error {
	if err := deployment.IsValidChainSelector(c.HomeChainSelector); err != nil {
		return fmt.Errorf("invalid home chain selector: %d - %w", c.HomeChainSelector, err)
	}
	if err := deployment.IsValidChainSelector(c.ChainSelector); err != nil {
		return fmt.Errorf("invalid chain selector: %d - %w", c.ChainSelector, err)
	}
	family, _ := chainsel.GetSelectorFamily(c.ChainSelector)
	if family != chainsel.FamilySolana {
		return fmt.Errorf("chain %d is not a solana chain", c.ChainSelector)
	}
	if err := c.UpgradeConfig.Validate(e, c.ChainSelector); err != nil {
		return fmt.Errorf("invalid UpgradeConfig: %w", err)
	}
	existingState, err := ccipChangeset.LoadOnchainState(e)
	if err != nil {
		return fmt.Errorf("failed to load existing onchain state: %w", err)
	}
	if _, exists := existingState.SupportedChains()[c.ChainSelector]; !exists {
		return fmt.Errorf("chain %d not supported", c.ChainSelector)
	}
	return nil
}

func DeployChainContractsChangeset(e deployment.Environment, c DeployChainContractsConfig) (deployment.ChangesetOutput, error) {
	if err := c.Validate(e); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("invalid DeployChainContractsConfig: %w", err)
	}
	newAddresses := deployment.NewMemoryAddressBook()
	existingState, _ := ccipChangeset.LoadOnchainState(e)
	err := v1_6.ValidateHomeChainState(e, c.HomeChainSelector, existingState)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	chainSel := c.ChainSelector
	chain := e.SolChains[chainSel]
	if existingState.SolChains[chainSel].LinkToken.IsZero() {
		return deployment.ChangesetOutput{}, fmt.Errorf("fee tokens not found for chain %d", chainSel)
	}

	// prepare artifacts
	// artifacts will already exist if running locally as chain spin up fetches them
	// on CI they wont be present and we want to fetch them here
	if c.BuildConfig != nil {
		e.Logger.Debugw("Building solana artifacts", "gitCommitSha", c.BuildConfig.GitCommitSha)
		err = BuildSolana(e, *c.BuildConfig)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to build solana: %w", err)
		}
	} else {
		e.Logger.Debugw("Skipping solana build as no build config provided")
	}

	if err := c.UpgradeConfig.Validate(e, chainSel); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("invalid UpgradeConfig: %w", err)
	}
	addresses, _ := e.ExistingAddresses.AddressesForChain(chainSel)
	mcmState, _ := state.MaybeLoadMCMSWithTimelockChainStateSolana(chain, addresses)
	timelocks := map[uint64]string{}
	proposers := map[uint64]string{}
	inspectors := map[uint64]sdk.Inspector{}
	var batches []mcmsTypes.BatchOperation
	timelocks[chainSel] = mcmsSolana.ContractAddress(
		mcmState.TimelockProgram,
		mcmsSolana.PDASeed(mcmState.TimelockSeed),
	)
	proposers[chainSel] = mcmsSolana.ContractAddress(mcmState.McmProgram, mcmsSolana.PDASeed(mcmState.ProposerMcmSeed))
	inspectors[chainSel] = mcmsSolana.NewInspector(chain.Client)

	mcmsTxs, err := deployChainContractsSolana(e, chain, newAddresses, c)
	if err != nil {
		e.Logger.Errorw("Failed to deploy CCIP contracts", "err", err, "newAddresses", newAddresses)
		return deployment.ChangesetOutput{}, err
	}
	// create proposals for txns
	if len(mcmsTxs) > 0 {
		batches = append(batches, mcmsTypes.BatchOperation{
			ChainSelector: mcmsTypes.ChainSelector(chainSel),
			Transactions:  mcmsTxs,
		})
	}

	if len(batches) > 0 {
		proposal, err := proposalutils.BuildProposalFromBatchesV2(
			e,
			timelocks,
			proposers,
			inspectors,
			batches,
			"proposal to upgrade CCIP contracts",
			c.UpgradeConfig.MCMS.MinDelay)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return deployment.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
			AddressBook:           newAddresses,
		}, nil
	}

	return deployment.ChangesetOutput{
		AddressBook: newAddresses,
	}, nil
}

func solProgramData(e deployment.Environment, chain deployment.SolChain, programID solana.PublicKey) (struct {
	DataType uint32
	Address  solana.PublicKey
}, error) {
	var programData struct {
		DataType uint32
		Address  solana.PublicKey
	}
	data, err := chain.Client.GetAccountInfoWithOpts(e.GetContext(), programID, &solRpc.GetAccountInfoOpts{
		Commitment: solRpc.CommitmentConfirmed,
	})
	if err != nil {
		return programData, fmt.Errorf("failed to deploy program: %w", err)
	}

	err = solBinary.UnmarshalBorsh(&programData, data.Bytes())
	if err != nil {
		return programData, fmt.Errorf("failed to unmarshal program data: %w", err)
	}
	return programData, nil
}

func SolProgramSize(e *deployment.Environment, chain deployment.SolChain, programID solana.PublicKey) (int, error) {
	accountInfo, err := chain.Client.GetAccountInfoWithOpts(e.GetContext(), programID, &rpc.GetAccountInfoOpts{
		Commitment: deployment.SolDefaultCommitment,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to get account info: %w", err)
	}
	if accountInfo == nil {
		return 0, fmt.Errorf("program account not found: %w", err)
	}
	programBytes := len(accountInfo.Value.Data.GetBinary())
	return programBytes, nil
}

func initializeRouter(
	e deployment.Environment,
	chain deployment.SolChain,
	ccipRouterProgram solana.PublicKey,
	linkTokenAddress solana.PublicKey,
	feeQuoterAddress solana.PublicKey,
	rmnRemoteAddress solana.PublicKey,
) error {
	e.Logger.Debugw("Initializing router", "chain", chain.String(), "ccipRouterProgram", ccipRouterProgram.String())
	programData, err := solProgramData(e, chain, ccipRouterProgram)
	if err != nil {
		return fmt.Errorf("failed to get solana router program data: %w", err)
	}
	// addressing errcheck in the next PR
	routerConfigPDA, _, _ := solState.FindConfigPDA(ccipRouterProgram)
	externalTokenPoolsSignerPDA, _, _ := solState.FindExternalTokenPoolsSignerPDA(ccipRouterProgram)

	instruction, err := solRouter.NewInitializeInstruction(
		chain.Selector, // chain selector
		// this is where the fee aggregator address would go (but have written a separate changeset to set that)
		solana.PublicKey{},
		feeQuoterAddress,
		linkTokenAddress, // link token mint
		rmnRemoteAddress,
		routerConfigPDA,
		chain.DeployerKey.PublicKey(),
		solana.SystemProgramID,
		ccipRouterProgram,
		programData.Address,
		externalTokenPoolsSignerPDA,
	).ValidateAndBuild()

	if err != nil {
		return fmt.Errorf("failed to build instruction: %w", err)
	}
	if err := chain.Confirm([]solana.Instruction{instruction}); err != nil {
		return fmt.Errorf("failed to confirm initializeRouter: %w", err)
	}
	e.Logger.Infow("Initialized router", "chain", chain.String())
	return nil
}

func initializeFeeQuoter(
	e deployment.Environment,
	chain deployment.SolChain,
	ccipRouterProgram solana.PublicKey,
	linkTokenAddress solana.PublicKey,
	feeQuoterAddress solana.PublicKey,
	offRampAddress solana.PublicKey,
	params FeeQuoterParams,
) error {
	e.Logger.Debugw("Initializing fee quoter", "chain", chain.String(), "feeQuoterAddress", feeQuoterAddress.String())
	programData, err := solProgramData(e, chain, feeQuoterAddress)
	if err != nil {
		return fmt.Errorf("failed to get solana router program data: %w", err)
	}
	feeQuoterConfigPDA, _, _ := solState.FindFqConfigPDA(feeQuoterAddress)

	instruction, err := solFeeQuoter.NewInitializeInstruction(
		params.DefaultMaxFeeJuelsPerMsg,
		ccipRouterProgram,
		feeQuoterConfigPDA,
		linkTokenAddress,
		chain.DeployerKey.PublicKey(),
		solana.SystemProgramID,
		feeQuoterAddress,
		programData.Address,
	).ValidateAndBuild()
	if err != nil {
		return fmt.Errorf("failed to build instruction: %w", err)
	}

	offRampBillingSignerPDA, _, _ := solState.FindOfframpBillingSignerPDA(offRampAddress)
	fqAllowedPriceUpdaterOfframpPDA, _, _ := solState.FindFqAllowedPriceUpdaterPDA(offRampBillingSignerPDA, feeQuoterAddress)

	priceUpdaterix, err := solFeeQuoter.NewAddPriceUpdaterInstruction(
		offRampBillingSignerPDA,
		fqAllowedPriceUpdaterOfframpPDA,
		feeQuoterConfigPDA,
		chain.DeployerKey.PublicKey(),
		solana.SystemProgramID,
	).ValidateAndBuild()

	if err != nil {
		return fmt.Errorf("failed to build instruction: %w", err)
	}
	if err := chain.Confirm([]solana.Instruction{instruction, priceUpdaterix}); err != nil {
		return fmt.Errorf("failed to confirm initializeFeeQuoter: %w", err)
	}
	e.Logger.Infow("Initialized fee quoter", "chain", chain.String())
	return nil
}

func initializeOffRamp(
	e deployment.Environment,
	chain deployment.SolChain,
	ccipRouterProgram solana.PublicKey,
	feeQuoterAddress solana.PublicKey,
	rmnRemoteAddress solana.PublicKey,
	offRampAddress solana.PublicKey,
	addressLookupTable solana.PublicKey,
	params OffRampParams,
) error {
	e.Logger.Debugw("Initializing offRamp", "chain", chain.String(), "offRampAddress", offRampAddress.String())
	programData, err := solProgramData(e, chain, offRampAddress)
	if err != nil {
		return fmt.Errorf("failed to get solana router program data: %w", err)
	}
	offRampConfigPDA, _, _ := solState.FindOfframpConfigPDA(offRampAddress)
	offRampReferenceAddressesPDA, _, _ := solState.FindOfframpReferenceAddressesPDA(offRampAddress)
	offRampStatePDA, _, _ := solState.FindOfframpStatePDA(offRampAddress)
	offRampExternalExecutionConfigPDA, _, _ := solState.FindExternalExecutionConfigPDA(offRampAddress)
	offRampTokenPoolsSignerPDA, _, _ := solState.FindExternalTokenPoolsSignerPDA(offRampAddress)

	initIx, err := solOffRamp.NewInitializeInstruction(
		offRampReferenceAddressesPDA,
		ccipRouterProgram,
		feeQuoterAddress,
		rmnRemoteAddress,
		addressLookupTable,
		offRampStatePDA,
		offRampExternalExecutionConfigPDA,
		offRampTokenPoolsSignerPDA,
		chain.DeployerKey.PublicKey(),
		solana.SystemProgramID,
		offRampAddress,
		programData.Address,
	).ValidateAndBuild()

	if err != nil {
		return fmt.Errorf("failed to build instruction: %w", err)
	}

	initConfigIx, err := solOffRamp.NewInitializeConfigInstruction(
		chain.Selector,
		params.EnableExecutionAfter,
		offRampConfigPDA,
		chain.DeployerKey.PublicKey(),
		solana.SystemProgramID,
		offRampAddress,
		programData.Address,
	).ValidateAndBuild()

	if err != nil {
		return fmt.Errorf("failed to build instruction: %w", err)
	}
	if err := chain.Confirm([]solana.Instruction{initIx, initConfigIx}); err != nil {
		return fmt.Errorf("failed to confirm initializeOffRamp: %w", err)
	}
	e.Logger.Infow("Initialized offRamp", "chain", chain.String())
	return nil
}

func initializeRMNRemote(
	e deployment.Environment,
	chain deployment.SolChain,
	rmnRemoteProgram solana.PublicKey,
) error {
	e.Logger.Debugw("Initializing rmn remote", "chain", chain.String(), "rmnRemoteProgram", rmnRemoteProgram.String())
	programData, err := solProgramData(e, chain, rmnRemoteProgram)
	if err != nil {
		return fmt.Errorf("failed to get solana router program data: %w", err)
	}
	rmnRemoteConfigPDA, _, _ := solState.FindRMNRemoteConfigPDA(rmnRemoteProgram)
	rmnRemoteCursesPDA, _, _ := solState.FindRMNRemoteCursesPDA(rmnRemoteProgram)
	instruction, err := solRmnRemote.NewInitializeInstruction(
		rmnRemoteConfigPDA,
		rmnRemoteCursesPDA,
		chain.DeployerKey.PublicKey(),
		solana.SystemProgramID,
		rmnRemoteProgram,
		programData.Address,
	).ValidateAndBuild()
	if err != nil {
		return fmt.Errorf("failed to build instruction: %w", err)
	}
	if err := chain.Confirm([]solana.Instruction{instruction}); err != nil {
		return fmt.Errorf("failed to confirm initializeRMNRemote: %w", err)
	}
	e.Logger.Infow("Initialized rmn remote", "chain", chain.String())
	return nil
}

func deployChainContractsSolana(
	e deployment.Environment,
	chain deployment.SolChain,
	ab deployment.AddressBook,
	config DeployChainContractsConfig,
) ([]mcmsTypes.Transaction, error) {
	// we may need to gather instructions and submit them as part of MCMS
	txns := make([]mcmsTypes.Transaction, 0)
	s, err := ccipChangeset.LoadOnchainStateSolana(e)
	if err != nil {
		e.Logger.Errorw("Failed to load existing onchain state", "err", err)
		return txns, err
	}
	chainState, chainExists := s.SolChains[chain.Selector]
	if !chainExists {
		return txns, fmt.Errorf("chain %s not found in existing state, deploy the link token first", chain.String())
	}
	if chainState.LinkToken.IsZero() {
		return txns, fmt.Errorf("failed to get link token address for chain %s", chain.String())
	}

	params := config.ContractParamsPerChain

	// FEE QUOTER DEPLOY
	var feeQuoterAddress solana.PublicKey
	//nolint:gocritic // this is a false positive, we need to check if the address is zero
	if chainState.FeeQuoter.IsZero() {
		feeQuoterAddress, err = DeployAndMaybeSaveToAddressBook(e, chain, ab, ccipChangeset.FeeQuoter, deployment.Version1_0_0, false)
		if err != nil {
			return txns, fmt.Errorf("failed to deploy program: %w", err)
		}
	} else if config.UpgradeConfig.NewFeeQuoterVersion != nil {
		// fee quoter updated in place
		feeQuoterAddress = chainState.FeeQuoter
		newTxns, err := generateUpgradeTxns(e, chain, ab, config, config.UpgradeConfig.NewFeeQuoterVersion, chainState.FeeQuoter, ccipChangeset.FeeQuoter)
		if err != nil {
			return txns, fmt.Errorf("failed to generate upgrade txns: %w", err)
		}
		txns = append(txns, newTxns...)
	} else {
		e.Logger.Infow("Using existing fee quoter", "addr", chainState.FeeQuoter.String())
		feeQuoterAddress = chainState.FeeQuoter
	}
	solFeeQuoter.SetProgramID(feeQuoterAddress)

	// ROUTER DEPLOY
	var ccipRouterProgram solana.PublicKey
	//nolint:gocritic // this is a false positive, we need to check if the address is zero
	if chainState.Router.IsZero() {
		// deploy router
		ccipRouterProgram, err = DeployAndMaybeSaveToAddressBook(e, chain, ab, ccipChangeset.Router, deployment.Version1_0_0, false)
		if err != nil {
			return txns, fmt.Errorf("failed to deploy program: %w", err)
		}
	} else if config.UpgradeConfig.NewRouterVersion != nil {
		// router updated in place
		ccipRouterProgram = chainState.Router
		newTxns, err := generateUpgradeTxns(e, chain, ab, config, config.UpgradeConfig.NewRouterVersion, chainState.Router, ccipChangeset.Router)
		if err != nil {
			return txns, fmt.Errorf("failed to generate upgrade txns: %w", err)
		}
		txns = append(txns, newTxns...)
	} else {
		e.Logger.Infow("Using existing router", "addr", chainState.Router.String())
		ccipRouterProgram = chainState.Router
	}
	solRouter.SetProgramID(ccipRouterProgram)

	// OFFRAMP DEPLOY
	var offRampAddress solana.PublicKey
	// gather lookup table keys from other deploys
	lookupTableKeys := make([]solana.PublicKey, 0)
	createLookupTable := false
	//nolint:gocritic // this is a false positive, we need to check if the address is zero
	if chainState.OffRamp.IsZero() {
		// deploy offramp
		offRampAddress, err = DeployAndMaybeSaveToAddressBook(e, chain, ab, ccipChangeset.OffRamp, deployment.Version1_0_0, false)
		if err != nil {
			return txns, fmt.Errorf("failed to deploy program: %w", err)
		}
	} else if config.UpgradeConfig.NewOffRampVersion != nil {
		tv := deployment.NewTypeAndVersion(ccipChangeset.OffRamp, *config.UpgradeConfig.NewOffRampVersion)
		existingAddresses, err := e.ExistingAddresses.AddressesForChain(chain.Selector)
		if err != nil {
			return txns, fmt.Errorf("failed to get existing addresses: %w", err)
		}
		offRampAddress = ccipChangeset.FindSolanaAddress(tv, existingAddresses)
		if offRampAddress.IsZero() {
			// deploy offramp, not upgraded in place so upgrade is false
			offRampAddress, err = DeployAndMaybeSaveToAddressBook(e, chain, ab, ccipChangeset.OffRamp, *config.UpgradeConfig.NewOffRampVersion, false)
			if err != nil {
				return txns, fmt.Errorf("failed to deploy program: %w", err)
			}
		}

		offRampBillingSignerPDA, _, _ := solState.FindOfframpBillingSignerPDA(offRampAddress)
		fqAllowedPriceUpdaterOfframpPDA, _, _ := solState.FindFqAllowedPriceUpdaterPDA(offRampBillingSignerPDA, feeQuoterAddress)
		feeQuoterConfigPDA, _, _ := solState.FindFqConfigPDA(feeQuoterAddress)

		priceUpdaterix, err := solFeeQuoter.NewAddPriceUpdaterInstruction(
			offRampBillingSignerPDA,
			fqAllowedPriceUpdaterOfframpPDA,
			feeQuoterConfigPDA,
			config.UpgradeConfig.UpgradeAuthority,
			solana.SystemProgramID,
		).ValidateAndBuild()
		if err != nil {
			return txns, fmt.Errorf("failed to build instruction: %w", err)
		}
		priceUpdaterTx, err := BuildMCMSTxn(priceUpdaterix, feeQuoterAddress.String(), ccipChangeset.FeeQuoter)
		if err != nil {
			return txns, fmt.Errorf("failed to create price updater transaction: %w", err)
		}
		txns = append(txns, *priceUpdaterTx)
	} else {
		e.Logger.Infow("Using existing offramp", "addr", chainState.OffRamp.String())
		offRampAddress = chainState.OffRamp
	}
	solOffRamp.SetProgramID(offRampAddress)

	// RMN REMOTE DEPLOY
	var rmnRemoteAddress solana.PublicKey
	if chainState.RMNRemote.IsZero() {
		rmnRemoteAddress, err = DeployAndMaybeSaveToAddressBook(e, chain, ab, ccipChangeset.RMNRemote, deployment.Version1_0_0, false)
		if err != nil {
			return txns, fmt.Errorf("failed to deploy program: %w", err)
		}
	} else if config.UpgradeConfig.NewRMNRemoteVersion != nil {
		rmnRemoteAddress = chainState.RMNRemote
		newTxns, err := generateUpgradeTxns(e, chain, ab, config, config.UpgradeConfig.NewRMNRemoteVersion, chainState.RMNRemote, ccipChangeset.RMNRemote)
		if err != nil {
			return txns, fmt.Errorf("failed to generate upgrade txns: %w", err)
		}
		txns = append(txns, newTxns...)
	} else {
		e.Logger.Infow("Using existing rmn remote", "addr", chainState.RMNRemote.String())
		rmnRemoteAddress = chainState.RMNRemote
	}
	solRmnRemote.SetProgramID(rmnRemoteAddress)

	// FEE QUOTER INITIALIZE
	var fqConfig solFeeQuoter.Config
	feeQuoterConfigPDA, _, _ := solState.FindFqConfigPDA(feeQuoterAddress)
	err = chain.GetAccountDataBorshInto(e.GetContext(), feeQuoterConfigPDA, &fqConfig)
	if err != nil {
		if err2 := initializeFeeQuoter(e, chain, ccipRouterProgram, chainState.LinkToken, feeQuoterAddress, offRampAddress, params.FeeQuoterParams); err2 != nil {
			return txns, err2
		}
	} else {
		e.Logger.Infow("Fee quoter already initialized, skipping initialization", "chain", chain.String())
	}

	// ROUTER INITIALIZE
	var routerConfigAccount solRouter.Config
	// addressing errcheck in the next PR
	routerConfigPDA, _, _ := solState.FindConfigPDA(ccipRouterProgram)
	err = chain.GetAccountDataBorshInto(e.GetContext(), routerConfigPDA, &routerConfigAccount)
	if err != nil {
		if err2 := initializeRouter(e, chain, ccipRouterProgram, chainState.LinkToken, feeQuoterAddress, rmnRemoteAddress); err2 != nil {
			return txns, err2
		}
	} else {
		e.Logger.Infow("Router already initialized, skipping initialization", "chain", chain.String())
	}

	// OFFRAMP INITIALIZE
	var offRampConfigAccount solOffRamp.Config
	offRampConfigPDA, _, _ := solState.FindOfframpConfigPDA(offRampAddress)
	err = chain.GetAccountDataBorshInto(e.GetContext(), offRampConfigPDA, &offRampConfigAccount)
	if err != nil {
		table, err2 := solCommonUtil.SetupLookupTable(
			e.GetContext(),
			chain.Client,
			*chain.DeployerKey,
			[]solana.PublicKey{
				// system
				solana.SystemProgramID,
				solana.ComputeBudget,
				solana.SysVarInstructionsPubkey,
				// token
				solana.Token2022ProgramID,
				solana.TokenProgramID,
				solana.SPLAssociatedTokenAccountProgramID,
			})
		if err2 != nil {
			return txns, fmt.Errorf("failed to create address lookup table: %w", err)
		}
		if err2 := initializeOffRamp(e, chain, ccipRouterProgram, feeQuoterAddress, rmnRemoteAddress, offRampAddress, table, params.OffRampParams); err2 != nil {
			return txns, err2
		}
		// Initializing a new offramp means we need a new lookup table and need to fully populate it
		createLookupTable = true
		offRampConfigPDA, _, _ := solState.FindOfframpConfigPDA(offRampAddress)
		offRampReferenceAddressesPDA, _, _ := solState.FindOfframpReferenceAddressesPDA(offRampAddress)
		offRampBillingSignerPDA, _, _ := solState.FindOfframpBillingSignerPDA(offRampAddress)
		lookupTableKeys = append(lookupTableKeys, []solana.PublicKey{
			// offramp
			offRampAddress,
			offRampConfigPDA,
			offRampReferenceAddressesPDA,
			offRampBillingSignerPDA,
		}...)
	} else {
		e.Logger.Infow("Offramp already initialized, skipping initialization", "chain", chain.String())
	}

	// RMN REMOTE INITIALIZE
	var rmnRemoteConfigAccount solRmnRemote.Config
	rmnRemoteConfigPDA, _, _ := solState.FindRMNRemoteConfigPDA(rmnRemoteAddress)
	err = chain.GetAccountDataBorshInto(e.GetContext(), rmnRemoteConfigPDA, &rmnRemoteConfigAccount)
	if err != nil {
		if err2 := initializeRMNRemote(e, chain, rmnRemoteAddress); err2 != nil {
			return txns, err2
		}
	} else {
		e.Logger.Infow("RMN remote already initialized, skipping initialization", "chain", chain.String())
	}

	// TOKEN POOLS DEPLOY
	var burnMintTokenPool solana.PublicKey
	if chainState.BurnMintTokenPool.IsZero() {
		burnMintTokenPool, err = DeployAndMaybeSaveToAddressBook(e, chain, ab, ccipChangeset.BurnMintTokenPool, deployment.Version1_0_0, false)
		if err != nil {
			return txns, fmt.Errorf("failed to deploy program: %w", err)
		}
	} else if config.UpgradeConfig.NewBurnMintTokenPoolVersion != nil {
		burnMintTokenPool = chainState.BurnMintTokenPool
		newTxns, err := generateUpgradeTxns(e, chain, ab, config, config.UpgradeConfig.NewBurnMintTokenPoolVersion, chainState.BurnMintTokenPool, ccipChangeset.BurnMintTokenPool)
		if err != nil {
			return txns, fmt.Errorf("failed to generate upgrade txns: %w", err)
		}
		txns = append(txns, newTxns...)
	} else {
		e.Logger.Infow("Using existing burn mint token pool", "addr", chainState.BurnMintTokenPool.String())
		burnMintTokenPool = chainState.BurnMintTokenPool
	}

	var lockReleaseTokenPool solana.PublicKey
	if chainState.LockReleaseTokenPool.IsZero() {
		lockReleaseTokenPool, err = DeployAndMaybeSaveToAddressBook(e, chain, ab, ccipChangeset.LockReleaseTokenPool, deployment.Version1_0_0, false)
		if err != nil {
			return txns, fmt.Errorf("failed to deploy program: %w", err)
		}
	} else if config.UpgradeConfig.NewLockReleaseTokenPoolVersion != nil {
		lockReleaseTokenPool = chainState.LockReleaseTokenPool
		newTxns, err := generateUpgradeTxns(e, chain, ab, config, config.UpgradeConfig.NewLockReleaseTokenPoolVersion, chainState.LockReleaseTokenPool, ccipChangeset.LockReleaseTokenPool)
		if err != nil {
			return txns, fmt.Errorf("failed to generate upgrade txns: %w", err)
		}
		txns = append(txns, newTxns...)
	} else {
		e.Logger.Infow("Using existing lock release token pool", "addr", chainState.LockReleaseTokenPool.String())
		lockReleaseTokenPool = chainState.LockReleaseTokenPool
	}

	// MCMS
	// this should selectively deploy anything if required
	// TODO: bad check
	if config.MCMSWithTimelockConfig.TimelockMinDelay != nil {
		_, err = solanaMCMS.DeployMCMSWithTimelockProgramsSolana(e, chain, ab, config.MCMSWithTimelockConfig)
		if err != nil {
			return txns, fmt.Errorf("failed to deploy MCMS with timelock programs: %w", err)
		}
	}
	addresses, err := e.ExistingAddresses.AddressesForChain(chain.Selector)
	if err != nil {
		return txns, fmt.Errorf("failed to get existing addresses: %w", err)
	}
	mcmState, err := state.MaybeLoadMCMSWithTimelockChainStateSolana(chain, addresses)
	if err != nil {
		return txns, fmt.Errorf("failed to load MCMS with timelock chain state: %w", err)
	}
	if config.UpgradeConfig.NewAccessControllerVersion != nil {
		e.Logger.Infow("Generating instruction for upgrading access controller", "chain", chain.String())
		newTxns, err := generateUpgradeTxns(e, chain, ab, config, config.UpgradeConfig.NewAccessControllerVersion, mcmState.AccessControllerProgram, types.AccessControllerProgram)
		if err != nil {
			return txns, fmt.Errorf("failed to generate upgrade txns: %w", err)
		}
		txns = append(txns, newTxns...)
	}
	if config.UpgradeConfig.NewTimelockVersion != nil {
		e.Logger.Infow("Generate instruction for upgrading timelock", "chain", chain.String())
		newTxns, err := generateUpgradeTxns(e, chain, ab, config, config.UpgradeConfig.NewTimelockVersion, mcmState.TimelockProgram, types.RBACTimelockProgram)
		if err != nil {
			return txns, fmt.Errorf("failed to generate upgrade txns: %w", err)
		}
		txns = append(txns, newTxns...)
	}
	if config.UpgradeConfig.NewMCMVersion != nil {
		e.Logger.Infow("Generate instruction for upgrading mcms", "chain", chain.String())
		newTxns, err := generateUpgradeTxns(e, chain, ab, config, config.UpgradeConfig.NewMCMVersion, mcmState.McmProgram, types.ManyChainMultisigProgram)
		if err != nil {
			return txns, fmt.Errorf("failed to generate upgrade txns: %w", err)
		}
		txns = append(txns, newTxns...)
	}

	// BILLING
	for _, billingConfig := range params.FeeQuoterParams.BillingConfig {
		if _, err := AddBillingToken(
			e, chain, chainState, billingConfig, nil, false, feeQuoterAddress, ccipRouterProgram,
		); err != nil {
			return txns, err
		}
	}

	if createLookupTable {
		// fee quoter enteries
		linkFqBillingConfigPDA, _, _ := solState.FindFqBillingTokenConfigPDA(chainState.LinkToken, feeQuoterAddress)
		wsolFqBillingConfigPDA, _, _ := solState.FindFqBillingTokenConfigPDA(chainState.WSOL, feeQuoterAddress)
		feeQuoterConfigPDA, _, _ := solState.FindFqConfigPDA(feeQuoterAddress)
		lookupTableKeys = append(lookupTableKeys, []solana.PublicKey{
			// fee quoter
			feeQuoterConfigPDA,
			feeQuoterAddress,
			linkFqBillingConfigPDA,
			wsolFqBillingConfigPDA,
		}...)

		// router entries
		externalExecutionConfigPDA, _, _ := solState.FindExternalExecutionConfigPDA(ccipRouterProgram)
		externalTokenPoolsSignerPDA, _, _ := solState.FindExternalTokenPoolsSignerPDA(ccipRouterProgram)
		routerConfigPDA, _, _ := solState.FindConfigPDA(ccipRouterProgram)
		feeBillingSignerPDA, _, _ := solState.FindFeeBillingSignerPDA(ccipRouterProgram)
		lookupTableKeys = append(lookupTableKeys, []solana.PublicKey{
			ccipRouterProgram,
			routerConfigPDA,
			externalExecutionConfigPDA,
			externalTokenPoolsSignerPDA,
			feeBillingSignerPDA,
		}...)

		// token pools entries
		lookupTableKeys = append(lookupTableKeys, []solana.PublicKey{
			burnMintTokenPool,
			lockReleaseTokenPool,
		}...)

		// rmn remote entries
		rmnRemoteCursePDA, _, _ := solState.FindRMNRemoteCursesPDA(rmnRemoteAddress)
		lookupTableKeys = append(lookupTableKeys, []solana.PublicKey{
			rmnRemoteAddress,
			rmnRemoteConfigPDA,
			rmnRemoteCursePDA,
		}...)
	}

	if len(lookupTableKeys) > 0 {
		e.Logger.Debugw("Populating lookup table", "keys", lookupTableKeys)
		if err := extendLookupTable(e, chain, offRampAddress, lookupTableKeys); err != nil {
			return txns, fmt.Errorf("failed to extend lookup table: %w", err)
		}
	}

	return txns, nil
}

func generateUpgradeTxns(
	e deployment.Environment,
	chain deployment.SolChain,
	ab deployment.AddressBook,
	config DeployChainContractsConfig,
	newVersion *semver.Version,
	programID solana.PublicKey,
	contractType deployment.ContractType,
) ([]mcmsTypes.Transaction, error) {
	e.Logger.Infow("Generating instruction for upgrading contract", "contractType", contractType)
	txns := make([]mcmsTypes.Transaction, 0)
	bufferProgram, err := DeployAndMaybeSaveToAddressBook(e, chain, ab, contractType, *newVersion, true)
	if err != nil {
		return txns, fmt.Errorf("failed to deploy program: %w", err)
	}
	if err := setUpgradeAuthority(&e, &chain, bufferProgram, chain.DeployerKey, config.UpgradeConfig.UpgradeAuthority.ToPointer(), true); err != nil {
		return txns, fmt.Errorf("failed to set upgrade authority: %w", err)
	}
	upgradeIxn, err := generateUpgradeIxn(
		&e,
		programID,
		bufferProgram,
		config.UpgradeConfig.SpillAddress,
		config.UpgradeConfig.UpgradeAuthority,
	)
	if err != nil {
		return txns, fmt.Errorf("failed to generate upgrade instruction: %w", err)
	}
	closeIxn, err := generateCloseBufferIxn(
		&e,
		bufferProgram,
		config.UpgradeConfig.SpillAddress,
		config.UpgradeConfig.UpgradeAuthority,
	)
	if err != nil {
		return txns, fmt.Errorf("failed to generate close buffer instruction: %w", err)
	}
	upgradeTx, err := BuildMCMSTxn(upgradeIxn, solana.BPFLoaderUpgradeableProgramID.String(), contractType)
	if err != nil {
		return txns, fmt.Errorf("failed to create upgrade transaction: %w", err)
	}
	closeTx, err := BuildMCMSTxn(closeIxn, solana.BPFLoaderUpgradeableProgramID.String(), contractType)
	if err != nil {
		return txns, fmt.Errorf("failed to create close transaction: %w", err)
	}
	// We do not support extend as part of upgrades due to MCMS limitations
	// https://docs.google.com/document/d/1Fk76lOeyS2z2X6MokaNX_QTMFAn5wvSZvNXJluuNV1E/edit?tab=t.0#heading=h.uij286zaarkz
	txns = append(txns, *upgradeTx, *closeTx)
	return txns, nil
}

// DeployAndMaybeSaveToAddressBook deploys a program to the Solana chain and saves it to the address book
// if it is not an upgrade. It returns the program ID of the deployed program.
func DeployAndMaybeSaveToAddressBook(
	e deployment.Environment,
	chain deployment.SolChain,
	ab deployment.AddressBook,
	contractType deployment.ContractType,
	version semver.Version,
	isUpgrade bool) (solana.PublicKey, error) {
	programName := getTypeToProgramDeployName()[contractType]
	programID, err := chain.DeployProgram(e.Logger, programName, isUpgrade)
	if err != nil {
		return solana.PublicKey{}, fmt.Errorf("failed to deploy program: %w", err)
	}
	address := solana.MustPublicKeyFromBase58(programID)

	e.Logger.Infow("Deployed program", "Program", contractType, "addr", programID, "chain", chain.String(), "isUpgrade", isUpgrade)

	if !isUpgrade {
		tv := deployment.NewTypeAndVersion(contractType, version)
		err = ab.Save(chain.Selector, programID, tv)
		if err != nil {
			return solana.PublicKey{}, fmt.Errorf("failed to save address: %w", err)
		}
	}
	return address, nil
}

func generateUpgradeIxn(
	e *deployment.Environment,
	programID solana.PublicKey,
	bufferAddress solana.PublicKey,
	spillAddress solana.PublicKey,
	upgradeAuthority solana.PublicKey,
) (solana.Instruction, error) {
	// Derive the program data address
	programDataAccount, _, _ := solana.FindProgramAddress([][]byte{programID.Bytes()}, solana.BPFLoaderUpgradeableProgramID)

	// Accounts involved in the transaction
	keys := solana.AccountMetaSlice{
		solana.NewAccountMeta(programDataAccount, true, false), // Program account (writable)
		solana.NewAccountMeta(programID, true, false),
		solana.NewAccountMeta(bufferAddress, true, false),             // Buffer account (writable)
		solana.NewAccountMeta(spillAddress, true, false),              // Spill account (writable)
		solana.NewAccountMeta(solana.SysVarRentPubkey, false, false),  // System program
		solana.NewAccountMeta(solana.SysVarClockPubkey, false, false), // System program
		solana.NewAccountMeta(upgradeAuthority, false, true),          // Current upgrade authority (signer)
	}

	instruction := solana.NewInstruction(
		solana.BPFLoaderUpgradeableProgramID,
		keys,
		// https://github.com/solana-playground/solana-playground/blob/2998d4cf381aa319d26477c5d4e6d15059670a75/vscode/src/commands/deploy/bpf-upgradeable/bpf-upgradeable.ts#L66
		[]byte{3, 0, 0, 0}, // 4-byte Upgrade instruction identifier
	)

	return instruction, nil
}

func generateExtendIxn(
	e *deployment.Environment,
	chain deployment.SolChain,
	programID solana.PublicKey,
	bufferAddress solana.PublicKey,
	payer solana.PublicKey,
) (*solana.GenericInstruction, error) {
	// Derive the program data address
	programDataAccount, _, _ := solana.FindProgramAddress([][]byte{programID.Bytes()}, solana.BPFLoaderUpgradeableProgramID)

	programDataSize, err := SolProgramSize(e, chain, programDataAccount)
	if err != nil {
		return nil, fmt.Errorf("failed to get program size: %w", err)
	}
	e.Logger.Debugw("Program data size", "programDataSize", programDataSize)

	bufferSize, err := SolProgramSize(e, chain, bufferAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get buffer size: %w", err)
	}
	e.Logger.Debugw("Buffer account size", "bufferSize", bufferSize)
	if bufferSize <= programDataSize {
		e.Logger.Debugf("Buffer account size %d is less than program account size %d", bufferSize, programDataSize)
		return nil, nil
	}
	extraBytes := bufferSize - programDataSize
	if extraBytes > math.MaxUint32 {
		return nil, fmt.Errorf("extra bytes %d exceeds maximum value %d", extraBytes, math.MaxUint32)
	}
	//https://github.com/solana-labs/solana/blob/7700cb3128c1f19820de67b81aa45d18f73d2ac0/sdk/program/src/loader_upgradeable_instruction.rs#L146
	data := binary.LittleEndian.AppendUint32([]byte{}, 6) // 4-byte Extend instruction identifier
	//nolint:gosec // G115 we check for overflow above
	data = binary.LittleEndian.AppendUint32(data, uint32(extraBytes+1024)) // add some padding

	keys := solana.AccountMetaSlice{
		solana.NewAccountMeta(programDataAccount, true, false),      // Program data account (writable)
		solana.NewAccountMeta(programID, true, false),               // Program account (writable)
		solana.NewAccountMeta(solana.SystemProgramID, false, false), // System program
		solana.NewAccountMeta(payer, true, true),                    // Payer for rent
	}

	ixn := solana.NewInstruction(
		solana.BPFLoaderUpgradeableProgramID,
		keys,
		data,
	)

	return ixn, nil
}

func generateCloseBufferIxn(
	e *deployment.Environment,
	bufferAddress solana.PublicKey,
	recipient solana.PublicKey,
	upgradeAuthority solana.PublicKey,
) (solana.Instruction, error) {
	keys := solana.AccountMetaSlice{
		solana.NewAccountMeta(bufferAddress, true, false),
		solana.NewAccountMeta(recipient, true, false),
		solana.NewAccountMeta(upgradeAuthority, false, true),
	}

	instruction := solana.NewInstruction(
		solana.BPFLoaderUpgradeableProgramID,
		keys,
		// https://github.com/solana-playground/solana-playground/blob/2998d4cf381aa319d26477c5d4e6d15059670a75/vscode/src/commands/deploy/bpf-upgradeable/bpf-upgradeable.ts#L78
		[]byte{5, 0, 0, 0}, // 4-byte Close instruction identifier
	)

	return instruction, nil
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
	ChainSelector uint64
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

	state, _ := ccipChangeset.LoadOnchainState(e)
	chainState := state.SolChains[cfg.ChainSelector]
	chain := e.SolChains[cfg.ChainSelector]
	ab := deployment.NewMemoryAddressBook()

	var receiverAddress solana.PublicKey
	var err error
	if chainState.Receiver.IsZero() {
		receiverAddress, err = DeployAndMaybeSaveToAddressBook(e, chain, ab, ccipChangeset.Receiver, deployment.Version1_0_0, false)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to deploy program: %w", err)
		}
	} else {
		e.Logger.Infow("Using existing receiver", "addr", chainState.Receiver.String())
		receiverAddress = chainState.Receiver
	}

	solTestReceiver.SetProgramID(receiverAddress)
	externalExecutionConfigPDA, _, _ := solState.FindExternalExecutionConfigPDA(receiverAddress)
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

	return deployment.ChangesetOutput{
		AddressBook: ab,
	}, nil
}

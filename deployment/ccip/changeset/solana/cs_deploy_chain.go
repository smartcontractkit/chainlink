package solana

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/gagliardetto/solana-go"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/mcms"
	"github.com/smartcontractkit/mcms/sdk"
	mcmsSolana "github.com/smartcontractkit/mcms/sdk/solana"
	mcmsTypes "github.com/smartcontractkit/mcms/types"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	solanastateview "github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/solana"
	"github.com/smartcontractkit/chainlink/deployment/helpers"

	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/common/types"

	solBinary "github.com/gagliardetto/binary"

	solOffRamp "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_offramp"
	solRouter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
	solFeeQuoter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/fee_quoter"
	solRmnRemote "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/rmn_remote"
	solCommonUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	solState "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"

	solanaMCMS "github.com/smartcontractkit/chainlink/deployment/common/changeset/solana/mcms"
)

// use this changeset to deploy the CCIP contracts on solana
var _ cldf.ChangeSet[DeployChainContractsConfig] = DeployChainContractsChangeset

type DeployChainContractsConfig struct {
	HomeChainSelector      uint64
	ChainSelector          uint64
	ContractParamsPerChain ChainContractParams
	UpgradeConfig          UpgradeConfig

	BuildConfig *helpers.BuildSolanaConfig
	// identifier for which token pool to deploy (i.e. partner identifier). defaults to CLL
	BurnMintTokenPoolMetadata    string
	LockReleaseTokenPoolMetadata string
	// TODO: add validation for this
	MCMSWithTimelockConfig *types.MCMSWithTimelockConfigV2
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
	// MCMS config must be set for upgrades and offramp redeploys (to configure the fee quoter after redeploy)
	MCMS *proposalutils.TimelockConfig
}

func (cfg UpgradeConfig) Validate(e cldf.Environment, chainSelector uint64) error {
	if cfg.NewFeeQuoterVersion == nil && cfg.NewRouterVersion == nil && cfg.NewOffRampVersion == nil {
		return nil
	}
	if cfg.SpillAddress.IsZero() {
		return errors.New("spill address must be set for fee quoter and router upgrades")
	}
	if cfg.UpgradeAuthority.IsZero() {
		return errors.New("upgrade authority must be set for fee quoter and router upgrades")
	}
	if cfg.MCMS != nil {
		return cfg.MCMS.ValidateSolana(e, chainSelector)
	}
	return nil
}

func (c DeployChainContractsConfig) Validate(e cldf.Environment) error {
	if err := cldf.IsValidChainSelector(c.HomeChainSelector); err != nil {
		return fmt.Errorf("invalid home chain selector: %d - %w", c.HomeChainSelector, err)
	}
	if err := cldf.IsValidChainSelector(c.ChainSelector); err != nil {
		return fmt.Errorf("invalid chain selector: %d - %w", c.ChainSelector, err)
	}
	family, _ := chainsel.GetSelectorFamily(c.ChainSelector)
	if family != chainsel.FamilySolana {
		return fmt.Errorf("chain %d is not a solana chain", c.ChainSelector)
	}
	if err := c.UpgradeConfig.Validate(e, c.ChainSelector); err != nil {
		return fmt.Errorf("invalid UpgradeConfig: %w", err)
	}
	existingState, err := stateview.LoadOnchainState(e)
	if err != nil {
		return fmt.Errorf("failed to load existing onchain state: %w", err)
	}
	if _, exists := existingState.SupportedChains()[c.ChainSelector]; !exists {
		return fmt.Errorf("chain %d not supported", c.ChainSelector)
	}
	chainState := existingState.SolChains[c.ChainSelector]

	// CLD:
	// the below check expects the user to pass in a mcms config when calling the changeset for the first time via CLD

	// In memory tests:
	// programs and state are pre-loaded, so we pass nil mcms config as router will be present in state
	// take a look at test_helpers.go/deployChainContractsToSolChainCS
	// initialisation of the mcms contracts then happens via testhelpers.TransferOwnershipSolana
	if chainState.Router.IsZero() {
		if c.MCMSWithTimelockConfig == nil {
			return fmt.Errorf("Router is not deployed. This looks like an initial deploy.MCMS config must be set for chain %d", c.ChainSelector)
		}
	}
	return nil
}

func DeployChainContractsChangeset(e cldf.Environment, c DeployChainContractsConfig) (cldf.ChangesetOutput, error) {
	if err := c.Validate(e); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("invalid DeployChainContractsConfig: %w", err)
	}
	newAddresses := cldf.NewMemoryAddressBook()
	existingState, _ := stateview.LoadOnchainState(e)
	err := v1_6.ValidateHomeChainState(e, c.HomeChainSelector, existingState)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	chainSel := c.ChainSelector
	chain := e.SolChains[chainSel]
	if existingState.SolChains[chainSel].LinkToken.IsZero() {
		return cldf.ChangesetOutput{}, fmt.Errorf("fee tokens not found for chain %d", chainSel)
	}

	// prepare artifacts
	// artifacts will already exist if running locally as chain spin up fetches them
	// on CI they wont be present and we want to fetch them here
	if c.BuildConfig != nil {
		e.Logger.Debugw("Building solana artifacts", "gitCommitSha", c.BuildConfig.GitCommitSha)
		err = helpers.BuildSolana(e, *c.BuildConfig, ccipBuildParams)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build solana: %w", err)
		}
	} else {
		e.Logger.Debugw("Skipping solana build as no build config provided")
	}

	if err := c.UpgradeConfig.Validate(e, chainSel); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("invalid UpgradeConfig: %w", err)
	}
	addresses, _ := e.ExistingAddresses.AddressesForChain(chainSel)
	mcmState, _ := state.MaybeLoadMCMSWithTimelockChainStateSolana(chain, addresses)
	timelocks := map[uint64]string{}
	proposers := map[uint64]string{}
	inspectors := map[uint64]sdk.Inspector{}
	timelocks[chainSel] = mcmsSolana.ContractAddress(
		mcmState.TimelockProgram,
		mcmsSolana.PDASeed(mcmState.TimelockSeed),
	)
	proposers[chainSel] = mcmsSolana.ContractAddress(mcmState.McmProgram, mcmsSolana.PDASeed(mcmState.ProposerMcmSeed))
	inspectors[chainSel] = mcmsSolana.NewInspector(chain.Client)

	batches, err := deployChainContractsSolana(e, chain, newAddresses, c)
	if err != nil {
		e.Logger.Errorw("Failed to deploy CCIP contracts", "err", err, "newAddresses", newAddresses)
		return cldf.ChangesetOutput{}, err
	}

	if len(batches) > 0 {
		proposal, err := proposalutils.BuildProposalFromBatchesV2(
			e,
			timelocks,
			proposers,
			inspectors,
			batches,
			"proposal to upgrade CCIP contracts",
			*c.UpgradeConfig.MCMS)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return cldf.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
			AddressBook:           newAddresses,
		}, nil
	}

	return cldf.ChangesetOutput{
		AddressBook: newAddresses,
	}, nil
}

func deployChainContractsSolana(
	e cldf.Environment,
	chain cldf.SolChain,
	ab cldf.AddressBook,
	config DeployChainContractsConfig,
) ([]mcmsTypes.BatchOperation, error) {
	// we may need to gather instructions and submit them as part of MCMS
	batches := make([]mcmsTypes.BatchOperation, 0)
	s, err := stateview.LoadOnchainStateSolana(e)
	if err != nil {
		e.Logger.Errorw("Failed to load existing onchain state", "err", err)
		return batches, err
	}
	chainState, chainExists := s.SolChains[chain.Selector]
	if !chainExists {
		return batches, fmt.Errorf("chain %s not found in existing state, deploy the link token first", chain.String())
	}
	if chainState.LinkToken.IsZero() {
		return batches, fmt.Errorf("failed to get link token address for chain %s", chain.String())
	}

	params := config.ContractParamsPerChain

	// FEE QUOTER DEPLOY
	var feeQuoterAddress solana.PublicKey
	//nolint:gocritic // this is a false positive, we need to check if the address is zero
	if chainState.FeeQuoter.IsZero() {
		feeQuoterAddress, err = helpers.DeployAndMaybeSaveToAddressBook(e, chain, ab, shared.FeeQuoter, deployment.Version1_0_0, false, "")
		if err != nil {
			return batches, fmt.Errorf("failed to deploy program: %w", err)
		}
	} else if config.UpgradeConfig.NewFeeQuoterVersion != nil {
		// fee quoter updated in place
		feeQuoterAddress = chainState.FeeQuoter
		newTxns, err := helpers.GenerateUpgradeTxns(e, chain, ab, config.UpgradeConfig.SpillAddress, config.UpgradeConfig.UpgradeAuthority,
			config.UpgradeConfig.NewFeeQuoterVersion, chainState.FeeQuoter, shared.FeeQuoter)
		if err != nil {
			return batches, fmt.Errorf("failed to generate upgrade txns: %w", err)
		}
		// create proposals for txns
		if len(newTxns) > 0 {
			batches = append(batches, mcmsTypes.BatchOperation{
				ChainSelector: mcmsTypes.ChainSelector(chain.Selector),
				Transactions:  newTxns,
			})
		}
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
		ccipRouterProgram, err = helpers.DeployAndMaybeSaveToAddressBook(e, chain, ab, shared.Router, deployment.Version1_0_0, false, "")
		if err != nil {
			return batches, fmt.Errorf("failed to deploy program: %w", err)
		}
	} else if config.UpgradeConfig.NewRouterVersion != nil {
		// router updated in place
		ccipRouterProgram = chainState.Router
		newTxns, err := helpers.GenerateUpgradeTxns(e, chain, ab, config.UpgradeConfig.SpillAddress, config.UpgradeConfig.UpgradeAuthority, config.UpgradeConfig.NewRouterVersion, chainState.Router, shared.Router)
		if err != nil {
			return batches, fmt.Errorf("failed to generate upgrade txns: %w", err)
		}
		// create proposals for txns
		if len(newTxns) > 0 {
			batches = append(batches, mcmsTypes.BatchOperation{
				ChainSelector: mcmsTypes.ChainSelector(chain.Selector),
				Transactions:  newTxns,
			})
		}
	} else {
		e.Logger.Infow("Using existing router", "addr", chainState.Router.String())
		ccipRouterProgram = chainState.Router
	}
	solRouter.SetProgramID(ccipRouterProgram)

	// OFFRAMP DEPLOY
	var offRampAddress solana.PublicKey
	//nolint:gocritic // this is a false positive, we need to check if the address is zero
	if chainState.OffRamp.IsZero() {
		// deploy offramp
		offRampAddress, err = helpers.DeployAndMaybeSaveToAddressBook(e, chain, ab, shared.OffRamp, deployment.Version1_0_0, false, "")
		if err != nil {
			return batches, fmt.Errorf("failed to deploy program: %w", err)
		}
	} else if config.UpgradeConfig.NewOffRampVersion != nil {
		tv := cldf.NewTypeAndVersion(shared.OffRamp, *config.UpgradeConfig.NewOffRampVersion)
		existingAddresses, err := e.ExistingAddresses.AddressesForChain(chain.Selector)
		if err != nil {
			return batches, fmt.Errorf("failed to get existing addresses: %w", err)
		}
		offRampAddress = solanastateview.FindSolanaAddress(tv, existingAddresses)
		if offRampAddress.IsZero() {
			// deploy offramp, not upgraded in place so upgrade is false
			offRampAddress, err = helpers.DeployAndMaybeSaveToAddressBook(e, chain, ab, shared.OffRamp, *config.UpgradeConfig.NewOffRampVersion, false, "")
			if err != nil {
				return batches, fmt.Errorf("failed to deploy program: %w", err)
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
				return batches, fmt.Errorf("failed to build instruction: %w", err)
			}
			if config.UpgradeConfig.MCMS != nil {
				priceUpdaterTx, err := helpers.BuildMCMSTxn(priceUpdaterix, feeQuoterAddress.String(), shared.FeeQuoter)
				if err != nil {
					return batches, fmt.Errorf("failed to create price updater transaction: %w", err)
				}
				batches = append(batches, mcmsTypes.BatchOperation{
					ChainSelector: mcmsTypes.ChainSelector(chain.Selector),
					Transactions:  []mcmsTypes.Transaction{*priceUpdaterTx},
				})
			} else {
				if err := chain.Confirm([]solana.Instruction{priceUpdaterix}); err != nil {
					return batches, fmt.Errorf("failed to confirm initializeFeeQuoter: %w", err)
				}
			}
		} else {
			newTxns, err := helpers.GenerateUpgradeTxns(e, chain, ab, config.UpgradeConfig.SpillAddress, config.UpgradeConfig.UpgradeAuthority, config.UpgradeConfig.NewOffRampVersion, chainState.OffRamp, shared.OffRamp)
			if err != nil {
				return batches, fmt.Errorf("failed to generate upgrade txns: %w", err)
			}
			// create proposals for txns
			if len(newTxns) > 0 {
				batches = append(batches, mcmsTypes.BatchOperation{
					ChainSelector: mcmsTypes.ChainSelector(chain.Selector),
					Transactions:  newTxns,
				})
			}
		}
	} else {
		e.Logger.Infow("Using existing offramp", "addr", chainState.OffRamp.String())
		offRampAddress = chainState.OffRamp
	}
	solOffRamp.SetProgramID(offRampAddress)

	// RMN REMOTE DEPLOY
	var rmnRemoteAddress solana.PublicKey
	if chainState.RMNRemote.IsZero() {
		rmnRemoteAddress, err = helpers.DeployAndMaybeSaveToAddressBook(e, chain, ab, shared.RMNRemote, deployment.Version1_0_0, false, "")
		if err != nil {
			return batches, fmt.Errorf("failed to deploy program: %w", err)
		}
	} else if config.UpgradeConfig.NewRMNRemoteVersion != nil {
		rmnRemoteAddress = chainState.RMNRemote
		newTxns, err := helpers.GenerateUpgradeTxns(e, chain, ab, config.UpgradeConfig.SpillAddress, config.UpgradeConfig.UpgradeAuthority, config.UpgradeConfig.NewRMNRemoteVersion, chainState.RMNRemote, shared.RMNRemote)
		if err != nil {
			return batches, fmt.Errorf("failed to generate upgrade txns: %w", err)
		}
		// create proposals for txns
		if len(newTxns) > 0 {
			batches = append(batches, mcmsTypes.BatchOperation{
				ChainSelector: mcmsTypes.ChainSelector(chain.Selector),
				Transactions:  newTxns,
			})
		}
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
			return batches, err2
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
			return batches, err2
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
			return batches, fmt.Errorf("failed to create address lookup table: %w", err)
		}
		if err2 := initializeOffRamp(e, chain, ccipRouterProgram, feeQuoterAddress, rmnRemoteAddress, offRampAddress, table, params.OffRampParams); err2 != nil {
			return batches, err2
		}
	} else {
		e.Logger.Infow("Offramp already initialized, skipping initialization", "chain", chain.String())
	}

	// RMN REMOTE INITIALIZE
	var rmnRemoteConfigAccount solRmnRemote.Config
	rmnRemoteConfigPDA, _, _ := solState.FindRMNRemoteConfigPDA(rmnRemoteAddress)
	err = chain.GetAccountDataBorshInto(e.GetContext(), rmnRemoteConfigPDA, &rmnRemoteConfigAccount)
	if err != nil {
		if err2 := initializeRMNRemote(e, chain, rmnRemoteAddress); err2 != nil {
			return batches, err2
		}
	} else {
		e.Logger.Infow("RMN remote already initialized, skipping initialization", "chain", chain.String())
	}

	// TOKEN POOLS DEPLOY
	var burnMintTokenPool solana.PublicKey
	metadata := shared.CLLMetadata
	if config.BurnMintTokenPoolMetadata != "" {
		metadata = config.BurnMintTokenPoolMetadata
	}
	//nolint:gocritic // this is a false positive, we need to check if the address is zero
	if chainState.BurnMintTokenPools[metadata].IsZero() {
		burnMintTokenPool, err = helpers.DeployAndMaybeSaveToAddressBook(e, chain, ab, shared.BurnMintTokenPool, deployment.Version1_0_0, false, metadata)
		if err != nil {
			return batches, fmt.Errorf("failed to deploy program: %w", err)
		}
	} else if config.UpgradeConfig.NewBurnMintTokenPoolVersion != nil {
		burnMintTokenPool = chainState.BurnMintTokenPools[metadata]
		newTxns, err := helpers.GenerateUpgradeTxns(e, chain, ab, config.UpgradeConfig.SpillAddress, config.UpgradeConfig.UpgradeAuthority, config.UpgradeConfig.NewBurnMintTokenPoolVersion, chainState.BurnMintTokenPools[metadata], shared.BurnMintTokenPool)
		if err != nil {
			return batches, fmt.Errorf("failed to generate upgrade txns: %w", err)
		}
		// create proposals for txns
		if len(newTxns) > 0 {
			batches = append(batches, mcmsTypes.BatchOperation{
				ChainSelector: mcmsTypes.ChainSelector(chain.Selector),
				Transactions:  newTxns,
			})
		}
	} else {
		e.Logger.Infow("Using existing burn mint token pool", "addr", chainState.BurnMintTokenPools[metadata].String())
		burnMintTokenPool = chainState.BurnMintTokenPools[metadata]
	}

	var lockReleaseTokenPool solana.PublicKey
	metadata = shared.CLLMetadata
	if config.LockReleaseTokenPoolMetadata != "" {
		metadata = config.LockReleaseTokenPoolMetadata
	}
	//nolint:gocritic // this is a false positive, we need to check if the address is zero
	if chainState.LockReleaseTokenPools[metadata].IsZero() {
		lockReleaseTokenPool, err = helpers.DeployAndMaybeSaveToAddressBook(e, chain, ab, shared.LockReleaseTokenPool, deployment.Version1_0_0, false, metadata)
		if err != nil {
			return batches, fmt.Errorf("failed to deploy program: %w", err)
		}
	} else if config.UpgradeConfig.NewLockReleaseTokenPoolVersion != nil {
		lockReleaseTokenPool = chainState.LockReleaseTokenPools[metadata]
		newTxns, err := helpers.GenerateUpgradeTxns(e, chain, ab, config.UpgradeConfig.SpillAddress, config.UpgradeConfig.UpgradeAuthority, config.UpgradeConfig.NewLockReleaseTokenPoolVersion, chainState.LockReleaseTokenPools[metadata], shared.LockReleaseTokenPool)
		if err != nil {
			return batches, fmt.Errorf("failed to generate upgrade txns: %w", err)
		}
		// create proposals for txns
		if len(newTxns) > 0 {
			batches = append(batches, mcmsTypes.BatchOperation{
				ChainSelector: mcmsTypes.ChainSelector(chain.Selector),
				Transactions:  newTxns,
			})
		}
	} else {
		e.Logger.Infow("Using existing lock release token pool", "addr", chainState.LockReleaseTokenPools[metadata].String())
		lockReleaseTokenPool = chainState.LockReleaseTokenPools[metadata]
	}

	// MCMS
	// this should selectively deploy and initialise anything if required
	if config.MCMSWithTimelockConfig != nil {
		_, err = solanaMCMS.DeployMCMSWithTimelockProgramsSolana(e, chain, ab, *config.MCMSWithTimelockConfig)
		if err != nil {
			return batches, fmt.Errorf("failed to deploy MCMS with timelock programs: %w", err)
		}
	}
	// MCMS Upgrade
	addresses, err := e.ExistingAddresses.AddressesForChain(chain.Selector)
	if err != nil {
		return batches, fmt.Errorf("failed to get existing addresses: %w", err)
	}
	mcmState, err := state.MaybeLoadMCMSWithTimelockChainStateSolana(chain, addresses)
	if err != nil {
		return batches, fmt.Errorf("failed to load MCMS with timelock chain state: %w", err)
	}
	if config.UpgradeConfig.NewMCMVersion != nil {
		e.Logger.Infow("Generate instruction for upgrading mcms", "chain", chain.String())
		newTxns, err := helpers.GenerateUpgradeTxns(e, chain, ab, config.UpgradeConfig.SpillAddress, config.UpgradeConfig.UpgradeAuthority, config.UpgradeConfig.NewMCMVersion, mcmState.McmProgram, types.ManyChainMultisigProgram)
		if err != nil {
			return batches, fmt.Errorf("failed to generate upgrade txns: %w", err)
		}
		// create proposals for txns
		if len(newTxns) > 0 {
			batches = append(batches, mcmsTypes.BatchOperation{
				ChainSelector: mcmsTypes.ChainSelector(chain.Selector),
				Transactions:  newTxns,
			})
		}
	}
	if config.UpgradeConfig.NewAccessControllerVersion != nil {
		e.Logger.Infow("Generating instruction for upgrading access controller", "chain", chain.String())
		newTxns, err := helpers.GenerateUpgradeTxns(e, chain, ab, config.UpgradeConfig.SpillAddress, config.UpgradeConfig.UpgradeAuthority, config.UpgradeConfig.NewAccessControllerVersion, mcmState.AccessControllerProgram, types.AccessControllerProgram)
		if err != nil {
			return batches, fmt.Errorf("failed to generate upgrade txns: %w", err)
		}
		// create proposals for txns
		if len(newTxns) > 0 {
			batches = append(batches, mcmsTypes.BatchOperation{
				ChainSelector: mcmsTypes.ChainSelector(chain.Selector),
				Transactions:  newTxns,
			})
		}
	}
	if config.UpgradeConfig.NewTimelockVersion != nil {
		e.Logger.Infow("Generate instruction for upgrading timelock", "chain", chain.String())
		newTxns, err := helpers.GenerateUpgradeTxns(e, chain, ab, config.UpgradeConfig.SpillAddress, config.UpgradeConfig.UpgradeAuthority, config.UpgradeConfig.NewTimelockVersion, mcmState.TimelockProgram, types.RBACTimelockProgram)
		if err != nil {
			return batches, fmt.Errorf("failed to generate upgrade txns: %w", err)
		}
		// create proposals for txns
		if len(newTxns) > 0 {
			batches = append(batches, mcmsTypes.BatchOperation{
				ChainSelector: mcmsTypes.ChainSelector(chain.Selector),
				Transactions:  newTxns,
			})
		}
	}

	// BILLING
	for _, billingConfig := range params.FeeQuoterParams.BillingConfig {
		if _, err := AddBillingToken(
			e, chain, chainState, billingConfig, nil, false, feeQuoterAddress, ccipRouterProgram,
		); err != nil {
			return batches, err
		}
	}

	// LOOKUP TABLE
	// off ramp
	offRampReferenceAddressesPDA, _, _ := solState.FindOfframpReferenceAddressesPDA(offRampAddress)
	offRampBillingSignerPDA, _, _ := solState.FindOfframpBillingSignerPDA(offRampAddress)
	// fee quoter
	linkFqBillingConfigPDA, _, _ := solState.FindFqBillingTokenConfigPDA(chainState.LinkToken, feeQuoterAddress)
	wsolFqBillingConfigPDA, _, _ := solState.FindFqBillingTokenConfigPDA(chainState.WSOL, feeQuoterAddress)
	// router
	feeBillingSignerPDA, _, _ := solState.FindFeeBillingSignerPDA(ccipRouterProgram)
	// rmn remote
	rmnRemoteCursePDA, _, _ := solState.FindRMNRemoteCursesPDA(rmnRemoteAddress)
	lookupTableKeys := []solana.PublicKey{
		// offramp
		offRampAddress,
		offRampConfigPDA,
		offRampReferenceAddressesPDA,
		offRampBillingSignerPDA,
		// fee quoter
		feeQuoterConfigPDA,
		feeQuoterAddress,
		linkFqBillingConfigPDA,
		wsolFqBillingConfigPDA,
		// router
		ccipRouterProgram,
		routerConfigPDA,
		feeBillingSignerPDA,
		// token pools
		burnMintTokenPool,
		lockReleaseTokenPool,
		// rmn remote
		rmnRemoteAddress,
		rmnRemoteConfigPDA,
		rmnRemoteCursePDA,
	}

	if err := extendLookupTable(e, chain, offRampAddress, lookupTableKeys); err != nil {
		return batches, fmt.Errorf("failed to extend lookup table: %w", err)
	}

	return batches, nil
}

// INITIALIZE FUNCTIONS
func initializeRouter(
	e cldf.Environment,
	chain cldf.SolChain,
	ccipRouterProgram solana.PublicKey,
	linkTokenAddress solana.PublicKey,
	feeQuoterAddress solana.PublicKey,
	rmnRemoteAddress solana.PublicKey,
) error {
	e.Logger.Debugw("Initializing router", "chain", chain.String(), "ccipRouterProgram", ccipRouterProgram.String())
	programData, err := helpers.GetSolProgramData(e, chain, ccipRouterProgram)
	if err != nil {
		return fmt.Errorf("failed to get solana router program data: %w", err)
	}
	// addressing errcheck in the next PR
	routerConfigPDA, _, _ := solState.FindConfigPDA(ccipRouterProgram)

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
	e cldf.Environment,
	chain cldf.SolChain,
	ccipRouterProgram solana.PublicKey,
	linkTokenAddress solana.PublicKey,
	feeQuoterAddress solana.PublicKey,
	offRampAddress solana.PublicKey,
	params FeeQuoterParams,
) error {
	e.Logger.Debugw("Initializing fee quoter", "chain", chain.String(), "feeQuoterAddress", feeQuoterAddress.String())
	programData, err := helpers.GetSolProgramData(e, chain, feeQuoterAddress)
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
	e cldf.Environment,
	chain cldf.SolChain,
	ccipRouterProgram solana.PublicKey,
	feeQuoterAddress solana.PublicKey,
	rmnRemoteAddress solana.PublicKey,
	offRampAddress solana.PublicKey,
	addressLookupTable solana.PublicKey,
	params OffRampParams,
) error {
	e.Logger.Debugw("Initializing offRamp", "chain", chain.String(), "offRampAddress", offRampAddress.String())
	programData, err := helpers.GetSolProgramData(e, chain, offRampAddress)
	if err != nil {
		return fmt.Errorf("failed to get solana router program data: %w", err)
	}
	offRampConfigPDA, _, _ := solState.FindOfframpConfigPDA(offRampAddress)
	offRampReferenceAddressesPDA, _, _ := solState.FindOfframpReferenceAddressesPDA(offRampAddress)
	offRampStatePDA, _, _ := solState.FindOfframpStatePDA(offRampAddress)

	initIx, err := solOffRamp.NewInitializeInstruction(
		offRampReferenceAddressesPDA,
		ccipRouterProgram,
		feeQuoterAddress,
		rmnRemoteAddress,
		addressLookupTable,
		offRampStatePDA,
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
	e cldf.Environment,
	chain cldf.SolChain,
	rmnRemoteProgram solana.PublicKey,
) error {
	e.Logger.Debugw("Initializing rmn remote", "chain", chain.String(), "rmnRemoteProgram", rmnRemoteProgram.String())
	programData, err := helpers.GetSolProgramData(e, chain, rmnRemoteProgram)
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

package solana

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/gagliardetto/solana-go"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"
	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"

	solBinary "github.com/gagliardetto/binary"
	solRpc "github.com/gagliardetto/solana-go/rpc"

	solOffRamp "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_offramp"
	solRouter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
	solFeeQuoter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/fee_quoter"
	solCommonUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	solState "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"
)

const (
	RouterProgramName    = "ccip_router"
	OffRampProgramName   = "ccip_offramp"
	FeeQuoterProgramName = "fee_quoter"
	TokenPoolProgramName = "test_token_pool"
)

var _ deployment.ChangeSet[DeployChainContractsConfigSolana] = DeployChainContractsChangesetSolana

type DeployChainContractsConfigSolana struct {
	DeployChainContractsConfig changeset.DeployChainContractsConfig
	UpgradeConfig              UpgradeConfigSolana
	NewUpgradeAuthority        *solana.PublicKey // if set, sets router and fee quoter upgrade authority
}

type UpgradeConfigSolana struct {
	NewFeeQuoterVersion *semver.Version
	NewRouterVersion    *semver.Version
	// Offramp is redeployed with the existing deployer key while the other programs are upgraded in place
	NewOffRampVersion *semver.Version
	// SpillAddress and UpgradeAuthority must be set
	SpillAddress     solana.PublicKey
	UpgradeAuthority solana.PublicKey
	// We assume if the signer key is set, we're not using MCMS
	SignerKey *solana.PrivateKey
}

func DeployChainContractsChangesetSolana(e deployment.Environment, config DeployChainContractsConfigSolana) (deployment.ChangesetOutput, error) {
	c := config.DeployChainContractsConfig
	if err := c.Validate(); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("invalid DeployChainContractsConfig: %w", err)
	}
	if config.UpgradeConfig.NewFeeQuoterVersion != nil || config.UpgradeConfig.NewRouterVersion != nil {
		if config.UpgradeConfig.SpillAddress.IsZero() {
			return deployment.ChangesetOutput{}, errors.New("spill address must be set for fee quoter and router upgrades")
		}
		if config.UpgradeConfig.UpgradeAuthority.IsZero() {
			return deployment.ChangesetOutput{}, errors.New("upgrade authority must be set for fee quoter and router upgrades")
		}
	}
	newAddresses := deployment.NewMemoryAddressBook()
	existingState, err := changeset.LoadOnchainState(e)
	if err != nil {
		e.Logger.Errorw("Failed to load existing onchain state", "err", err)
		return deployment.ChangesetOutput{}, err
	}

	err = changeset.ValidateHomeChainState(e, c.HomeChainSelector, existingState)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	for chainSel := range c.ContractParamsPerChain {
		if _, exists := existingState.SupportedChains()[chainSel]; !exists {
			return deployment.ChangesetOutput{}, fmt.Errorf("chain %d not supported", chainSel)
		}
		// already validated family
		family, _ := chainsel.GetSelectorFamily(chainSel)
		if family != chainsel.FamilySolana {
			return deployment.ChangesetOutput{}, fmt.Errorf("chain %d is not a solana chain", chainSel)
		}
		chain := e.SolChains[chainSel]
		if existingState.SolChains[chainSel].LinkToken.IsZero() {
			return deployment.ChangesetOutput{}, fmt.Errorf("fee tokens not found for chain %d", chainSel)
		}
		ixns, err := deployChainContractsSolana(e, chain, newAddresses, config)
		if err != nil {
			e.Logger.Errorw("Failed to deploy CCIP contracts", "err", err, "newAddresses", newAddresses)
			return deployment.ChangesetOutput{}, err
		}
		// create proposals for ixns
		if len(ixns) > 0 {
			e.Logger.Infow("Deployed contracts", "chain", chain.String(), "instructions", ixns)
		}
	}

	return deployment.ChangesetOutput{
		Proposals:   []timelock.MCMSWithTimelockProposal{},
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

func initializeRouter(
	e deployment.Environment,
	chain deployment.SolChain,
	ccipRouterProgram solana.PublicKey,
	linkTokenAddress solana.PublicKey,
	feeQuoterAddress solana.PublicKey,
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
		chain.Selector,     // chain selector
		solana.PublicKey{}, // fee aggregator (TODO: changeset to set the fee aggregator)
		feeQuoterAddress,
		linkTokenAddress, // link token mint
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
) error {
	e.Logger.Debugw("Initializing fee quoter", "chain", chain.String(), "feeQuoterAddress", feeQuoterAddress.String())
	programData, err := solProgramData(e, chain, feeQuoterAddress)
	if err != nil {
		return fmt.Errorf("failed to get solana router program data: %w", err)
	}
	feeQuoterConfigPDA, _, _ := solState.FindFqConfigPDA(feeQuoterAddress)

	instruction, err := solFeeQuoter.NewInitializeInstruction(
		linkTokenAddress,
		deployment.SolDefaultMaxFeeJuelsPerMsg,
		ccipRouterProgram,
		feeQuoterConfigPDA,
		chain.DeployerKey.PublicKey(),
		solana.SystemProgramID,
		feeQuoterAddress,
		programData.Address,
	).ValidateAndBuild()

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
	offRampAddress solana.PublicKey,
	addressLookupTable solana.PublicKey,
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
		deployment.EnableExecutionAfter,
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

func deployChainContractsSolana(
	e deployment.Environment,
	chain deployment.SolChain,
	ab deployment.AddressBook,
	config DeployChainContractsConfigSolana,
) ([]solana.Instruction, error) {
	// we may need to gather instructions and submit them as part of MCMS
	ixns := make([]solana.Instruction, 0)
	state, err := changeset.LoadOnchainStateSolana(e)
	if err != nil {
		e.Logger.Errorw("Failed to load existing onchain state", "err", err)
		return ixns, err
	}
	chainState, chainExists := state.SolChains[chain.Selector]
	if !chainExists {
		return ixns, fmt.Errorf("chain %s not found in existing state, deploy the link token first", chain.String())
	}
	if chainState.LinkToken.IsZero() {
		return ixns, fmt.Errorf("failed to get link token address for chain %s", chain.String())
	}

	// FEE QUOTER DEPLOY
	var feeQuoterAddress solana.PublicKey
	//nolint:gocritic // this is a false positive, we need to check if the address is zero
	if chainState.FeeQuoter.IsZero() {
		feeQuoterAddress, err = DeployAndMaybeSaveToAddressBook(e, chain, ab, FeeQuoterProgramName, deployment.Version1_0_0, false)
		if err != nil {
			return ixns, fmt.Errorf("failed to deploy program: %w", err)
		}
	} else if config.UpgradeConfig.NewFeeQuoterVersion != nil {
		// fee quoter updated in place
		bufferProgram, err := DeployAndMaybeSaveToAddressBook(e, chain, ab, FeeQuoterProgramName, *config.UpgradeConfig.NewFeeQuoterVersion, true)
		if err != nil {
			return ixns, fmt.Errorf("failed to deploy program: %w", err)
		}
		if err := setUpgradeAuthority(&e, &chain, bufferProgram, chain.DeployerKey, config.UpgradeConfig.UpgradeAuthority.ToPointer(), true); err != nil {
			return ixns, fmt.Errorf("failed to set upgrade authority: %w", err)
		}
		upgradeIxn, err := generateUpgradeIxn(
			&e,
			chainState.FeeQuoter,
			bufferProgram,
			config.UpgradeConfig.SpillAddress,
			config.UpgradeConfig.UpgradeAuthority,
		)
		if err != nil {
			return ixns, fmt.Errorf("failed to generate upgrade instruction: %w", err)
		}
		closeIxn, err := generateCloseBufferIxn(
			&e,
			bufferProgram,
			config.UpgradeConfig.SpillAddress,
			config.UpgradeConfig.UpgradeAuthority,
		)
		if err != nil {
			return ixns, fmt.Errorf("failed to generate close buffer instruction: %w", err)
		}
		feeQuoterAddress = chainState.FeeQuoter
		if config.UpgradeConfig.SignerKey != nil {
			// if we're not using MCMS to upgrade, confirm txn with the signer key
			if err := chain.Confirm([]solana.Instruction{upgradeIxn, closeIxn}, solCommonUtil.AddSigners(*config.UpgradeConfig.SignerKey)); err != nil {
				return ixns, fmt.Errorf("failed to confirm upgradeFeeQuoter: %w", err)
			}
			e.Logger.Infow("Upgraded FeeQuoter", "addr", chainState.FeeQuoter.String(), "chain", chain.String())
		} else {
			ixns = append(ixns, upgradeIxn, closeIxn)
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
		ccipRouterProgram, err = DeployAndMaybeSaveToAddressBook(e, chain, ab, RouterProgramName, deployment.Version1_0_0, false)
		if err != nil {
			return ixns, fmt.Errorf("failed to deploy program: %w", err)
		}
	} else if config.UpgradeConfig.NewRouterVersion != nil {
		// router updated in place
		bufferProgram, err := DeployAndMaybeSaveToAddressBook(e, chain, ab, RouterProgramName, *config.UpgradeConfig.NewRouterVersion, true)
		if err != nil {
			return ixns, fmt.Errorf("failed to deploy program: %w", err)
		}
		if err := setUpgradeAuthority(&e, &chain, bufferProgram, chain.DeployerKey, config.UpgradeConfig.UpgradeAuthority.ToPointer(), true); err != nil {
			return ixns, fmt.Errorf("failed to set upgrade authority: %w", err)
		}
		upgradeIxn, err := generateUpgradeIxn(
			&e,
			chainState.Router,
			bufferProgram,
			config.UpgradeConfig.SpillAddress,
			config.UpgradeConfig.UpgradeAuthority,
		)
		if err != nil {
			return ixns, fmt.Errorf("failed to generate upgrade instruction: %w", err)
		}
		closeIxn, err := generateCloseBufferIxn(
			&e,
			bufferProgram,
			config.UpgradeConfig.SpillAddress,
			config.UpgradeConfig.UpgradeAuthority,
		)
		if err != nil {
			return ixns, fmt.Errorf("failed to generate close buffer instruction: %w", err)
		}
		ccipRouterProgram = chainState.Router
		if config.UpgradeConfig.SignerKey != nil {
			// if we're not using MCMS to upgrade, confirm txn with the signer key
			if err := chain.Confirm([]solana.Instruction{upgradeIxn, closeIxn}, solCommonUtil.AddSigners(*config.UpgradeConfig.SignerKey)); err != nil {
				return ixns, fmt.Errorf("failed to confirm upgradeRouter: %w", err)
			}
			e.Logger.Infow("Upgraded Router", "addr", chainState.Router.String(), "chain", chain.String())
		} else {
			ixns = append(ixns, upgradeIxn, closeIxn)
		}
	} else {
		e.Logger.Infow("Using existing router", "addr", chainState.Router.String())
		ccipRouterProgram = chainState.Router
	}
	solRouter.SetProgramID(ccipRouterProgram)

	// OFFRAMP DEPLOY
	var offRampAddress solana.PublicKey
	// gather lookup table keys from other deploys
	lookupTableKeys := make([]solana.PublicKey, 0)
	needFQinLookupTable := false
	needRouterinLookupTable := false
	needTokenPoolinLookupTable := false
	//nolint:gocritic // this is a false positive, we need to check if the address is zero
	if chainState.OffRamp.IsZero() {
		// deploy offramp
		offRampAddress, err = DeployAndMaybeSaveToAddressBook(e, chain, ab, OffRampProgramName, deployment.Version1_0_0, false)
		if err != nil {
			return ixns, fmt.Errorf("failed to deploy program: %w", err)
		}
	} else if config.UpgradeConfig.NewOffRampVersion != nil {
		tv := deployment.NewTypeAndVersion(changeset.OffRamp, *config.UpgradeConfig.NewOffRampVersion)
		existingAddresses, err := e.ExistingAddresses.AddressesForChain(chain.Selector)
		if err != nil {
			return ixns, fmt.Errorf("failed to get existing addresses: %w", err)
		}
		offRampAddress = changeset.FindSolanaAddress(tv, existingAddresses)
		if offRampAddress.IsZero() {
			// deploy offramp, not upgraded in place so upgrade is false
			offRampAddress, err = DeployAndMaybeSaveToAddressBook(e, chain, ab, OffRampProgramName, *config.UpgradeConfig.NewOffRampVersion, false)
			if err != nil {
				return ixns, fmt.Errorf("failed to deploy program: %w", err)
			}
		}

		allowedOfframpSvmPDA, _ := solState.FindAllowedOfframpPDA(chain.Selector, offRampAddress, ccipRouterProgram)
		routerConfigPDA, _, _ := solState.FindConfigPDA(ccipRouterProgram)

		addIxn, err := solRouter.NewAddOfframpInstruction(
			chain.Selector,
			offRampAddress,
			allowedOfframpSvmPDA,
			routerConfigPDA,
			chain.DeployerKey.PublicKey(),
			solana.SystemProgramID,
		).ValidateAndBuild()
		if err != nil {
			return ixns, fmt.Errorf("failed to build instruction: %w", err)
		}
		ixns = append(ixns, []solana.Instruction{addIxn}...)

		offRampBillingSignerPDA, _, _ := solState.FindOfframpBillingSignerPDA(offRampAddress)
		fqAllowedPriceUpdaterOfframpPDA, _, _ := solState.FindFqAllowedPriceUpdaterPDA(offRampBillingSignerPDA, feeQuoterAddress)
		feeQuoterConfigPDA, _, _ := solState.FindFqConfigPDA(feeQuoterAddress)

		priceUpdaterix, err := solFeeQuoter.NewAddPriceUpdaterInstruction(
			offRampBillingSignerPDA,
			fqAllowedPriceUpdaterOfframpPDA,
			feeQuoterConfigPDA,
			chain.DeployerKey.PublicKey(),
			solana.SystemProgramID,
		).ValidateAndBuild()
		if err != nil {
			return ixns, fmt.Errorf("failed to build instruction: %w", err)
		}
		ixns = append(ixns, []solana.Instruction{priceUpdaterix}...)
	} else {
		e.Logger.Infow("Using existing offramp", "addr", chainState.OffRamp.String())
		offRampAddress = chainState.OffRamp
	}
	solOffRamp.SetProgramID(offRampAddress)

	// FEE QUOTER INITIALIZE
	var fqConfig solFeeQuoter.Config
	feeQuoterConfigPDA, _, _ := solState.FindFqConfigPDA(feeQuoterAddress)
	err = chain.GetAccountDataBorshInto(e.GetContext(), feeQuoterConfigPDA, &fqConfig)
	if err != nil {
		if err2 := initializeFeeQuoter(e, chain, ccipRouterProgram, chainState.LinkToken, feeQuoterAddress, offRampAddress); err2 != nil {
			return ixns, err2
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
		if err2 := initializeRouter(e, chain, ccipRouterProgram, chainState.LinkToken, feeQuoterAddress); err2 != nil {
			return ixns, err2
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
			return ixns, fmt.Errorf("failed to create address lookup table: %w", err)
		}
		if err2 := initializeOffRamp(e, chain, ccipRouterProgram, feeQuoterAddress, offRampAddress, table); err2 != nil {
			return ixns, err2
		}
		// Initializing a new offramp means we need a new lookup table and need to fully populate it
		needFQinLookupTable = true
		needRouterinLookupTable = true
		needTokenPoolinLookupTable = true
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

	// TOKEN POOL DEPLOY
	var tokenPoolProgram solana.PublicKey
	//nolint:gocritic // this is a false positive, we need to check if the address is zero
	if chainState.TokenPool.IsZero() {
		// TODO: there should be two token pools deployed one of each type (lock/burn)
		// separate token pools are not ready yet
		tokenPoolProgram, err = DeployAndMaybeSaveToAddressBook(e, chain, ab, TokenPoolProgramName, deployment.Version1_0_0, false)
		if err != nil {
			return ixns, fmt.Errorf("failed to deploy program: %w", err)
		}
		needTokenPoolinLookupTable = true
	} else {
		e.Logger.Infow("Using existing token pool", "addr", chainState.TokenPool.String())
		tokenPoolProgram = chainState.TokenPool
	}

	if needFQinLookupTable {
		linkFqBillingConfigPDA, _, _ := solState.FindFqBillingTokenConfigPDA(chainState.LinkToken, feeQuoterAddress)
		feeQuoterConfigPDA, _, _ := solState.FindFqConfigPDA(feeQuoterAddress)
		lookupTableKeys = append(lookupTableKeys, []solana.PublicKey{
			// fee quoter
			feeQuoterConfigPDA,
			feeQuoterAddress,
			linkFqBillingConfigPDA,
		}...)
	}

	if needRouterinLookupTable {
		externalExecutionConfigPDA, _, _ := solState.FindExternalExecutionConfigPDA(ccipRouterProgram)
		externalTokenPoolsSignerPDA, _, _ := solState.FindExternalTokenPoolsSignerPDA(ccipRouterProgram)
		routerConfigPDA, _, _ := solState.FindConfigPDA(ccipRouterProgram)
		feeBillingSignerPDA, _, _ := solState.FindFeeBillingSignerPDA(ccipRouterProgram)
		lookupTableKeys = append(lookupTableKeys, []solana.PublicKey{
			// router
			ccipRouterProgram,
			routerConfigPDA,
			externalExecutionConfigPDA,
			externalTokenPoolsSignerPDA,
			feeBillingSignerPDA,
		}...)
	}

	if needTokenPoolinLookupTable {
		lookupTableKeys = append(lookupTableKeys, []solana.PublicKey{
			// token pools
			tokenPoolProgram,
		}...)
	}

	if len(lookupTableKeys) > 0 {
		addressLookupTable, err := changeset.FetchOfframpLookupTable(e.GetContext(), chain, offRampAddress)
		if err != nil {
			return ixns, fmt.Errorf("failed to get offramp reference addresses: %w", err)
		}
		e.Logger.Debugw("Populating lookup table", "lookupTable", addressLookupTable.String(), "keys", lookupTableKeys)
		if err := solCommonUtil.ExtendLookupTable(e.GetContext(), chain.Client, addressLookupTable, *chain.DeployerKey, lookupTableKeys); err != nil {
			return ixns, fmt.Errorf("failed to extend lookup table: %w", err)
		}
	}

	// set upgrade authority
	if config.NewUpgradeAuthority != nil {
		e.Logger.Infow("Setting upgrade authority", "newUpgradeAuthority", config.NewUpgradeAuthority.String())
		for _, programID := range []solana.PublicKey{ccipRouterProgram, feeQuoterAddress} {
			if err := setUpgradeAuthority(&e, &chain, programID, chain.DeployerKey, config.NewUpgradeAuthority, false); err != nil {
				return ixns, fmt.Errorf("failed to set upgrade authority: %w", err)
			}
		}
	}

	return ixns, nil
}

// DeployAndMaybeSaveToAddressBook deploys a program to the Solana chain and saves it to the address book
// if it is not an upgrade. It returns the program ID of the deployed program.
func DeployAndMaybeSaveToAddressBook(
	e deployment.Environment,
	chain deployment.SolChain,
	ab deployment.AddressBook,
	programName string,
	version semver.Version,
	isUpgrade bool) (solana.PublicKey, error) {
	programID, err := chain.DeployProgram(e.Logger, programName, isUpgrade)
	if err != nil {
		return solana.PublicKey{}, fmt.Errorf("failed to deploy program: %w", err)
	}
	address := solana.MustPublicKeyFromBase58(programID)

	programNameToType := map[string]deployment.ContractType{
		RouterProgramName:    changeset.Router,
		OffRampProgramName:   changeset.OffRamp,
		FeeQuoterProgramName: changeset.FeeQuoter,
		TokenPoolProgramName: changeset.TokenPool,
	}
	programType, ok := programNameToType[programName]
	if !ok {
		return solana.PublicKey{}, fmt.Errorf("unknown program name: %s", programName)
	}
	e.Logger.Infow("Deployed program", "Program", programType, "addr", programID, "chain", chain.String(), "isUpgrade", isUpgrade)

	if !isUpgrade {
		tv := deployment.NewTypeAndVersion(programType, version)
		err = ab.Save(chain.Selector, programID, tv)
		if err != nil {
			return solana.PublicKey{}, fmt.Errorf("failed to save address: %w", err)
		}
	}
	return address, nil
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

func generateCloseBufferIxn(
	e *deployment.Environment,
	bufferAddress solana.PublicKey,
	recipient solana.PublicKey,
	upgradeAuthority solana.PublicKey,
) (solana.Instruction, error) {
	keys := solana.AccountMetaSlice{
		solana.NewAccountMeta(bufferAddress, true, false),
		solana.NewAccountMeta(recipient, false, false),
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

package solana

import (
	"fmt"

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

var _ deployment.ChangeSet[changeset.DeployChainContractsConfig] = DeployChainContractsChangesetSolana

func DeployChainContractsChangesetSolana(e deployment.Environment, c changeset.DeployChainContractsConfig) (deployment.ChangesetOutput, error) {
	if err := c.Validate(); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("invalid DeployChainContractsConfig: %w", err)
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
		err = deployChainContractsSolana(e, chain, newAddresses)
		if err != nil {
			e.Logger.Errorw("Failed to deploy CCIP contracts", "err", err, "newAddresses", newAddresses)
			return deployment.ChangesetOutput{}, err
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
		return fmt.Errorf("failed to confirm instructions: %w", err)
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
		return fmt.Errorf("failed to confirm instructions: %w", err)
	}
	e.Logger.Infow("Initialized fee quoter", "chain", chain.String())
	return nil
}

func intializeOffRamp(
	e deployment.Environment,
	chain deployment.SolChain,
	ccipRouterProgram solana.PublicKey,
	feeQuoterAddress solana.PublicKey,
	offRampAddress solana.PublicKey,
	addressLookupTable solana.PublicKey,
) error {
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
		return fmt.Errorf("failed to confirm instructions: %w", err)
	}
	e.Logger.Infow("Initialized offRamp", "chain", chain.String())
	return nil
}

func deployChainContractsSolana(
	e deployment.Environment,
	chain deployment.SolChain,
	ab deployment.AddressBook,
) error {
	state, err := changeset.LoadOnchainStateSolana(e)
	if err != nil {
		e.Logger.Errorw("Failed to load existing onchain state", "err", err)
		return err
	}
	chainState, chainExists := state.SolChains[chain.Selector]
	if !chainExists {
		return fmt.Errorf("chain %s not found in existing state, deploy the link token first", chain.String())
	}
	if chainState.LinkToken.IsZero() {
		return fmt.Errorf("failed to get link token address for chain %s", chain.String())
	}

	// initialize this last with every address we need
	var addressLookupTable solana.PublicKey
	if chainState.OfframpAddressLookupTable.IsZero() {
		addressLookupTable, err = solCommonUtil.SetupLookupTable(
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

		if err != nil {
			return fmt.Errorf("failed to create lookup table: %w", err)
		}
		err = ab.Save(chain.Selector, addressLookupTable.String(), deployment.NewTypeAndVersion(changeset.OfframpAddressLookupTable, deployment.Version1_0_0))
		if err != nil {
			return fmt.Errorf("failed to save address: %w", err)
		}
	}

	// FEE QUOTER DEPLOY
	var feeQuoterAddress solana.PublicKey
	if chainState.FeeQuoter.IsZero() {
		// deploy fee quoter
		programID, err := chain.DeployProgram(e.Logger, "fee_quoter")
		if err != nil {
			return fmt.Errorf("failed to deploy program: %w", err)
		}

		tv := deployment.NewTypeAndVersion(changeset.FeeQuoter, deployment.Version1_0_0)
		e.Logger.Infow("Deployed contract", "Contract", tv.String(), "addr", programID, "chain", chain.String())

		feeQuoterAddress = solana.MustPublicKeyFromBase58(programID)
		err = ab.Save(chain.Selector, programID, tv)
		if err != nil {
			return fmt.Errorf("failed to save address: %w", err)
		}
	} else {
		e.Logger.Infow("Using existing fee quoter", "addr", chainState.FeeQuoter.String())
		feeQuoterAddress = chainState.FeeQuoter
	}
	solFeeQuoter.SetProgramID(feeQuoterAddress)

	// ROUTER DEPLOY
	var ccipRouterProgram solana.PublicKey
	if chainState.Router.IsZero() {
		// deploy router
		programID, err := chain.DeployProgram(e.Logger, "ccip_router")
		if err != nil {
			return fmt.Errorf("failed to deploy program: %w", err)
		}

		tv := deployment.NewTypeAndVersion(changeset.Router, deployment.Version1_0_0)
		e.Logger.Infow("Deployed contract", "Contract", tv.String(), "addr", programID, "chain", chain.String())

		ccipRouterProgram = solana.MustPublicKeyFromBase58(programID)
		err = ab.Save(chain.Selector, programID, tv)
		if err != nil {
			return fmt.Errorf("failed to save address: %w", err)
		}
	} else {
		e.Logger.Infow("Using existing router", "addr", chainState.Router.String())
		ccipRouterProgram = chainState.Router
	}
	solRouter.SetProgramID(ccipRouterProgram)

	// OFFRAMP DEPLOY
	var offRampAddress solana.PublicKey
	if chainState.OffRamp.IsZero() {
		// deploy offramp
		programID, err := chain.DeployProgram(e.Logger, "ccip_offramp")
		if err != nil {
			return fmt.Errorf("failed to deploy program: %w", err)
		}
		tv := deployment.NewTypeAndVersion(changeset.OffRamp, deployment.Version1_0_0)
		e.Logger.Infow("Deployed contract", "Contract", tv.String(), "addr", programID, "chain", chain.String())
		offRampAddress = solana.MustPublicKeyFromBase58(programID)
		err = ab.Save(chain.Selector, programID, tv)
		if err != nil {
			return fmt.Errorf("failed to save address: %w", err)
		}
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
			return err2
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
			return err2
		}
	} else {
		e.Logger.Infow("Router already initialized, skipping initialization", "chain", chain.String())
	}

	// OFFRAMP INITIALIZE
	var offRampConfigAccount solOffRamp.Config
	offRampConfigPDA, _, _ := solState.FindOfframpConfigPDA(offRampAddress)
	err = chain.GetAccountDataBorshInto(e.GetContext(), offRampConfigPDA, &offRampConfigAccount)
	if err != nil {
		if err2 := intializeOffRamp(e, chain, ccipRouterProgram, feeQuoterAddress, offRampAddress, addressLookupTable); err2 != nil {
			return err2
		}
	} else {
		e.Logger.Infow("Offramp already initialized, skipping initialization", "chain", chain.String())
	}

	// TOKEN POOL DEPLOY
	var tokenPoolProgram solana.PublicKey
	if chainState.TokenPool.IsZero() {
		// TODO: there should be two token pools deployed one of each type (lock/burn)
		// separate token pools are not ready yet
		programID, err := chain.DeployProgram(e.Logger, "test_token_pool")
		if err != nil {
			return fmt.Errorf("failed to deploy program: %w", err)
		}
		tv := deployment.NewTypeAndVersion(changeset.TokenPool, deployment.Version1_0_0)
		e.Logger.Infow("Deployed contract", "Contract", tv.String(), "addr", programID, "chain", chain.String())
		tokenPoolProgram = solana.MustPublicKeyFromBase58(programID)
		err = ab.Save(chain.Selector, programID, tv)
		if err != nil {
			return fmt.Errorf("failed to save address: %w", err)
		}
	} else {
		e.Logger.Infow("Using existing token pool", "addr", chainState.TokenPool.String())
		tokenPoolProgram = chainState.TokenPool
	}

	externalExecutionConfigPDA, _, _ := solState.FindExternalExecutionConfigPDA(ccipRouterProgram)
	externalTokenPoolsSignerPDA, _, _ := solState.FindExternalTokenPoolsSignerPDA(ccipRouterProgram)
	feeBillingSignerPDA, _, _ := solState.FindFeeBillingSignerPDA(ccipRouterProgram)
	linkFqBillingConfigPDA, _, _ := solState.FindFqBillingTokenConfigPDA(chainState.LinkToken, feeQuoterAddress)
	offRampReferenceAddressesPDA, _, _ := solState.FindOfframpReferenceAddressesPDA(offRampAddress)
	offRampBillingSignerPDA, _, _ := solState.FindOfframpBillingSignerPDA(offRampAddress)

	if err := solCommonUtil.ExtendLookupTable(e.GetContext(), chain.Client, addressLookupTable, *chain.DeployerKey,
		[]solana.PublicKey{
			// token pools
			tokenPoolProgram,
			// offramp
			offRampAddress,
			offRampConfigPDA,
			offRampReferenceAddressesPDA,
			offRampBillingSignerPDA,
			// router
			ccipRouterProgram,
			routerConfigPDA,
			externalExecutionConfigPDA,
			externalTokenPoolsSignerPDA,
			// fee quoter
			feeBillingSignerPDA,
			feeQuoterConfigPDA,
			feeQuoterAddress,
			linkFqBillingConfigPDA,
		}); err != nil {
		return fmt.Errorf("failed to extend lookup table: %w", err)
	}

	return nil
}

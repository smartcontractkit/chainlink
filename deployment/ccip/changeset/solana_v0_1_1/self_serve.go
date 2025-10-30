package solana

import (
	"context"
	"errors"
	"fmt"

	"github.com/gagliardetto/solana-go"

	solCommon "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_1/ccip_common"
	solRouter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_1/ccip_router"
	solState "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"
	cldf_solana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	solanastateview "github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/solana"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	mcmsTypes "github.com/smartcontractkit/mcms/types"
)

var _ cldf.ChangeSet[OnboardTokenPoolsForSelfServeConfig] = OnboardTokenPoolsForSelfServe

type OnboardTokenPoolConfig struct {
	TokenPubKey             solana.PublicKey
	TokenAdminRegistryAdmin solana.PublicKey
	Override                bool
}

type OnboardTokenPoolsForSelfServeConfig struct {
	ChainSelector        uint64
	RegisterTokenConfigs []OnboardTokenPoolConfig
	MCMS                 *proposalutils.TimelockConfig
}

func (cfg OnboardTokenPoolsForSelfServeConfig) Validate(e cldf.Environment, chainState solanastateview.CCIPChainState) error {
	chain := e.BlockChains.SolanaChains()[cfg.ChainSelector]
	if err := chainState.ValidateRouterConfig(chain); err != nil {
		return err
	}
	if err := ValidateMCMSConfigSolana(e, cfg.MCMS, chain, chainState, solana.PublicKey{}, "", map[cldf.ContractType]bool{shared.Router: true}); err != nil {
		return err
	}
	routerProgramAddress, _, _ := chainState.GetRouterInfo()
	// Duplicate mint detection
	seen := make(map[string]int, len(cfg.RegisterTokenConfigs))
	for i, registerTokenConfig := range cfg.RegisterTokenConfigs {
		mintStr := registerTokenConfig.TokenPubKey.String()
		if firstIdx, dup := seen[mintStr]; dup {
			return fmt.Errorf("duplicate token mint %s found at indexes %d and %d", mintStr, firstIdx, i)
		}
		seen[mintStr] = i
		if registerTokenConfig.TokenAdminRegistryAdmin.IsZero() {
			return errors.New("token admin registry admin is required")
		}
		tokenPubKey := registerTokenConfig.TokenPubKey
		if err := chainState.CommonValidation(e, cfg.ChainSelector, tokenPubKey); err != nil {
			return err
		}
		tokenAdminRegistryPDA, _, err := solState.FindTokenAdminRegistryPDA(tokenPubKey, routerProgramAddress)
		if err != nil {
			return fmt.Errorf("failed to find token admin registry pda (mint: %s, router: %s): %w",
				mintStr, routerProgramAddress.String(), err)
		}
		var tokenAdminRegistryAccount solCommon.TokenAdminRegistry
		if err := chain.GetAccountDataBorshInto(context.Background(), tokenAdminRegistryPDA, &tokenAdminRegistryAccount); err == nil {
			if !registerTokenConfig.Override {
				return fmt.Errorf("token admin registry already exists for (mint: %s, router: %s)", mintStr, routerProgramAddress.String())
			}
		}
	}
	return nil
}

// OnboardTokenPoolsForSelfServe registers a token admin registry for a given token and initializes the token pool in CLL Token Pool Program.
// This changeset is used when the owner of the token pool doesn't have the mint authority over the token, but they want to self serve.
// So, this changeset includes the minimum configuration that CCIP Admin needs to do in the Token Admin Registry and in the Token Pool Program
func OnboardTokenPoolsForSelfServe(e cldf.Environment, cfg OnboardTokenPoolsForSelfServeConfig) (cldf.ChangesetOutput, error) {
	e.Logger.Infow("OnboardTokenPoolsForSelfServe", "cfg", cfg)
	routerState, err := loadRouterSolanaState(e, cfg)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	mcmsTxs := []mcmsTypes.Transaction{}
	instructions := []solana.Instruction{}
	for _, registerTokenConfig := range cfg.RegisterTokenConfigs {
		instruction, err := generateProposeTokenAdminRegistryAdministratorIx(registerTokenConfig, routerState)
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}
		// if the ccip admin is timelock, build mcms transaction
		if cfg.MCMS != nil {
			tx, err := BuildMCMSTxn(instruction, routerState.routerProgramID.String(), shared.Router)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to create transaction: %w", err)
			}
			mcmsTxs = append(mcmsTxs, *tx)
		} else {
			// the ccip admin will always be deployer key if done without mcms
			instructions = append(instructions, instruction)
		}
	}
	return ExecuteIndividualInstructionsAndBuildProposals(e, ExecuteConfig{ChainSelector: cfg.ChainSelector, MCMS: cfg.MCMS, Chain: routerState.chain}, instructions, mcmsTxs)
}

func generateProposeTokenAdminRegistryAdministratorIx(registerTokenConfig OnboardTokenPoolConfig, routerState routerSolanaState) (solana.Instruction, error) {
	tokenPubKey := registerTokenConfig.TokenPubKey
	tokenAdminRegistryPDA, _, _ := solState.FindTokenAdminRegistryPDA(tokenPubKey, routerState.routerProgramID)
	tokenAdminRegistryAdmin := registerTokenConfig.TokenAdminRegistryAdmin
	var instruction solana.Instruction
	// the ccip admin signs and makes tokenAdminRegistryAdmin the pending authority of the tokenAdminRegistry PDA, then they need to accept the role
	if !registerTokenConfig.Override {
		tempIx, err := solRouter.NewCcipAdminProposeAdministratorInstruction(
			tokenAdminRegistryAdmin, // customer's admin of the tokenAdminRegistry PDA in the Router
			routerState.routerConfigPDA,
			tokenAdminRegistryPDA, // If invoking the first time, this PDA is created
			tokenPubKey,
			routerState.ccipAdmin,
			solana.SystemProgramID,
		).ValidateAndBuild()
		if err != nil {
			return nil, fmt.Errorf("failed to generate instruction to propose administrator: %w", err)
		}
		ixData, err := tempIx.Data()
		if err != nil {
			return nil, fmt.Errorf("failed to extract data payload from ccip admin propose admin instruction: %w", err)
		}
		instruction = solana.NewInstruction(routerState.routerProgramID, tempIx.Accounts(), ixData)
	} else {
		// Use this if the proposed token admin registry admin set was incorrect
		overridePendingAdministratorIx, err := solRouter.NewCcipAdminOverridePendingAdministratorInstruction(
			tokenAdminRegistryAdmin, // customer's admin of the tokenAdminRegistry PDA in the Router
			routerState.routerConfigPDA,
			tokenAdminRegistryPDA,
			tokenPubKey,
			routerState.ccipAdmin,
			solana.SystemProgramID,
		).ValidateAndBuild()
		if err != nil {
			return nil, fmt.Errorf("failed to generate instruction to override pending administrator: %w", err)
		}
		ixData, err := overridePendingAdministratorIx.Data()
		if err != nil {
			return nil, fmt.Errorf("failed to extract data payload from ccip admin override pending admin instruction: %w", err)
		}
		instruction = solana.NewInstruction(routerState.routerProgramID, overridePendingAdministratorIx.Accounts(), ixData)
	}
	return instruction, nil
}

type routerSolanaState struct {
	chain           cldf_solana.Chain
	routerProgramID solana.PublicKey
	routerConfigPDA solana.PublicKey
	ccipAdmin       solana.PublicKey
}

func loadRouterSolanaState(e cldf.Environment, cfg OnboardTokenPoolsForSelfServeConfig) (routerSolanaState, error) {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return routerSolanaState{}, err
	}
	chainState, ok := state.SolChains[cfg.ChainSelector]
	if !ok {
		return routerSolanaState{}, fmt.Errorf("chain %d not found in environment", cfg.ChainSelector)
	}

	if err := cfg.Validate(e, chainState); err != nil {
		return routerSolanaState{}, err
	}
	chain := e.BlockChains.SolanaChains()[cfg.ChainSelector]
	routerProgramAddress, routerConfigPDA, _ := chainState.GetRouterInfo()
	ccipAdmin := GetAuthorityForIxn(
		&e,
		chain,
		chainState,
		shared.Router,
		solana.PublicKey{},
		"",
	)
	return routerSolanaState{
		chain:           chain,
		routerProgramID: routerProgramAddress,
		routerConfigPDA: routerConfigPDA,
		ccipAdmin:       ccipAdmin,
	}, nil
}

package solana

import (
	"context"
	"errors"
	"fmt"

	"github.com/gagliardetto/solana-go"
	solCommon "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_1/ccip_common"
	solRouter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_1/ccip_router"
	solState "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	solanastateview "github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/solana"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/mcms"
	mcmsTypes "github.com/smartcontractkit/mcms/types"
)

var _ cldf.ChangeSet[OnboardTokenPoolsForSelfServeConfig] = OnboardTokenPoolForSelfServe

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

	for _, registerTokenConfig := range cfg.RegisterTokenConfigs {
		if registerTokenConfig.TokenAdminRegistryAdmin.IsZero() {
			return errors.New("token admin registry admin is required")
		}
		tokenPubKey := registerTokenConfig.TokenPubKey
		if err := chainState.CommonValidation(e, cfg.ChainSelector, tokenPubKey); err != nil {
			return err
		}
		tokenAdminRegistryPDA, _, err := solState.FindTokenAdminRegistryPDA(tokenPubKey, routerProgramAddress)
		if err != nil {
			return fmt.Errorf("failed to find token admin registry pda (mint: %s, router: %s): %w", tokenPubKey.String(), routerProgramAddress.String(), err)
		}
		var tokenAdminRegistryAccount solCommon.TokenAdminRegistry
		if err := chain.GetAccountDataBorshInto(context.Background(), tokenAdminRegistryPDA, &tokenAdminRegistryAccount); err == nil {
			if !registerTokenConfig.Override {
				return fmt.Errorf("token admin registry already exists for (mint: %s, router: %s)", tokenPubKey.String(), routerProgramAddress.String())
			}
		}
	}

	return nil
}

// OnboardTokenPoolForSelfServe registers a token admin registry for a given token and initializes the token pool in CLL Token Pool Program.
// This changeset is used when the owner of the token pool doesn't have the mint authority over the token, but they want to self serve.
// So, this changeset includes the minimum configuration that CCIP Admin needs to do in the Token Admin Registry and in the Token Pool Program
func OnboardTokenPoolForSelfServe(e cldf.Environment, cfg OnboardTokenPoolsForSelfServeConfig) (cldf.ChangesetOutput, error) {
	e.Logger.Infow("OnboardTokenPoolForSelfServe", "cfg", cfg)
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	chainState, ok := state.SolChains[cfg.ChainSelector]
	if !ok {
		return cldf.ChangesetOutput{}, fmt.Errorf("chain %d not found in environment", cfg.ChainSelector)
	}
	if err := cfg.Validate(e, chainState); err != nil {
		return cldf.ChangesetOutput{}, err
	}
	chain := e.BlockChains.SolanaChains()[cfg.ChainSelector]
	routerProgramAddress, routerConfigPDA, _ := chainState.GetRouterInfo()

	timelockSignerPDA, err := FetchTimelockSigner(e, cfg.ChainSelector)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to fetch timelock signer: %w", err)
	}

	ccipAdmin := GetAuthorityForIxn(
		&e,
		chain,
		chainState,
		shared.Router,
		solana.PublicKey{},
		"")

	mcmsTxs := []mcmsTypes.Transaction{}

	for _, registerTokenConfig := range cfg.RegisterTokenConfigs {
		tokenPubKey := registerTokenConfig.TokenPubKey
		tokenAdminRegistryPDA, _, _ := solState.FindTokenAdminRegistryPDA(tokenPubKey, routerProgramAddress)
		tokenAdminRegistryAdmin := registerTokenConfig.TokenAdminRegistryAdmin
		var instruction solana.Instruction

		// the ccip admin signs and makes tokenAdminRegistryAdmin the pending authority of the tokenAdminRegistry PDA, then they need to accept the role
		if !registerTokenConfig.Override {
			tempIx, err := solRouter.NewCcipAdminProposeAdministratorInstruction(
				tokenAdminRegistryAdmin, // customer's admin of the tokenAdminRegistry PDA in the Router
				routerConfigPDA,
				tokenAdminRegistryPDA, // If invoking the first time, this PDA is created
				tokenPubKey,
				ccipAdmin,
				solana.SystemProgramID,
			).ValidateAndBuild()
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate instruction to propose administrator: %w", err)
			}
			ixData, err := tempIx.Data()
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to extract data payload from ccip admin propose admin instruction: %w", err)
			}
			instruction = solana.NewInstruction(routerProgramAddress, tempIx.Accounts(), ixData)
		} else {
			// Use this if the proposed token admin registry admin set was incorrect
			overridePendingAdministratorIx, err := solRouter.NewCcipAdminOverridePendingAdministratorInstruction(
				tokenAdminRegistryAdmin, // customer's admin of the tokenAdminRegistry PDA in the Router
				routerConfigPDA,
				tokenAdminRegistryPDA,
				tokenPubKey,
				ccipAdmin,
				solana.SystemProgramID,
			).ValidateAndBuild()
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate instruction to override pending administrator: %w", err)
			}
			ixData, err := overridePendingAdministratorIx.Data()
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to extract data payload from ccip admin override pending admin instruction: %w", err)
			}
			instruction = solana.NewInstruction(routerProgramAddress, overridePendingAdministratorIx.Accounts(), ixData)
		}

		// as ccip admin is proposing the admin role, it needs to sign the transaction
		// if the ccip admin is timelock, build mcms transaction
		// else just confirm it
		if ccipAdmin.Equals(timelockSignerPDA) {
			tx, err := BuildMCMSTxn(instruction, routerProgramAddress.String(), shared.Router)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to create transaction: %w", err)
			}
			mcmsTxs = append(mcmsTxs, *tx)
		} else {
			// the ccip admin will always be deployer key if done without mcms
			instructions := []solana.Instruction{instruction}
			if err := chain.Confirm(instructions); err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm instructions: %w", err)
			}
		}
	}

	if len(mcmsTxs) > 0 {
		proposal, err := BuildProposalsForTxns(
			e, cfg.ChainSelector, "proposal to OnboardTokenPoolForSelfServe in Solana", cfg.MCMS.MinDelay, mcmsTxs)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return cldf.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	}

	return cldf.ChangesetOutput{}, nil
}

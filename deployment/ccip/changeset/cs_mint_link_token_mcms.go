package changeset

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/deployergroup"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

// This changeset generates a proposal to mint LINK when LINK is already owned by MCMS
var MintLinkTokenMCMS = cldf.CreateChangeSet(MintLinkTokenMCMSLogic, MintLinkTokenMCMSPreconditions)

type MintLinkTokenMCMSConfig struct {
	Selector   uint64                        `json:"selector"`
	ToAddress  common.Address                `json:"toAddress"`
	Amount     *big.Int                      `json:"amount"`
	MCMSConfig *proposalutils.TimelockConfig `json:"mcmsConfig"`
}

func (cfg MintLinkTokenMCMSConfig) Validate(e cldf.Environment) error {
	if err := cldf.IsValidChainSelector(cfg.Selector); err != nil {
		return fmt.Errorf("invalid chain selector: %d - %w", cfg.Selector, err)
	}

	if cfg.ToAddress == (common.Address{}) {
		return errors.New("toAddress cannot be empty")
	}

	if cfg.Amount == nil || cfg.Amount.Sign() <= 0 {
		return errors.New("amount must be greater than zero")
	}

	if cfg.MCMSConfig == nil {
		return errors.New("mcmsConfig is required for this changeset - use GrantMintRoleAndMint for non-MCMS owned tokens")
	}

	return nil
}

func MintLinkTokenMCMSPreconditions(e cldf.Environment, cfg MintLinkTokenMCMSConfig) error {
	if err := cfg.Validate(e); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}

	chainName := e.BlockChains.EVMChains()[cfg.Selector].Name()
	if e.Name == "mainnet" || strings.Contains(chainName, "mainnet") {
		return errors.New("minting on LINK token is not allowed on Mainnet")
	}

	chainState, ok := state.EVMChainState(cfg.Selector)
	if !ok {
		return fmt.Errorf("%d does not exist in state", cfg.Selector)
	}
	if linkToken := chainState.LinkToken; linkToken == nil {
		return fmt.Errorf("missing linkToken on %d", cfg.Selector)
	}

	return nil
}

func MintLinkTokenMCMSLogic(e cldf.Environment, cfg MintLinkTokenMCMSConfig) (cldf.ChangesetOutput, error) {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	chainState, ok := state.EVMChainState(cfg.Selector)
	if !ok {
		return cldf.ChangesetOutput{}, fmt.Errorf("chain state not found for selector %d", cfg.Selector)
	}

	linkToken := chainState.LinkToken

	// Create deployer group with MCMS config - this will generate a proposal instead of executing directly
	deployerGroup := deployergroup.NewDeployerGroup(e, state, cfg.MCMSConfig).
		WithDeploymentContext(fmt.Sprintf("Mint %s LINK tokens to %s on chain %d",
			cfg.Amount.String(), cfg.ToAddress.Hex(), cfg.Selector))

	// Get the deployer (TransactOpts) - when MCMS is set, this returns opts configured for the timelock
	opts, err := deployerGroup.GetDeployer(cfg.Selector)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get deployer: %w", err)
	}

	// Check if the timelock has the minter role
	timelockAddr := opts.From
	isMinter, err := linkToken.IsMinter(&bind.CallOpts{Context: e.GetContext()}, timelockAddr)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to check if timelock is minter: %w", err)
	}

	// If the timelock doesn't have the minter role, grant it
	if !isMinter {
		e.Logger.Infow("Timelock does not have minter role, granting mint and burn roles",
			"chain", cfg.Selector,
			"timelock", timelockAddr.Hex(),
		)
		// Grant mint and burn roles to the timelock - creates a simulated transaction that will be included in the MCMS proposal
		_, err = linkToken.GrantMintAndBurnRoles(opts, timelockAddr)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to prepare grant mint and burn roles transaction: %w", err)
		}
	}

	e.Logger.Infow("Preparing MCMS proposal to mint LINK tokens",
		"chain", cfg.Selector,
		"to", cfg.ToAddress.Hex(),
		"amount", cfg.Amount.String(),
	)

	// Call mint - this creates a simulated transaction that will be included in the MCMS proposal
	_, err = linkToken.Mint(opts, cfg.ToAddress, cfg.Amount)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to prepare mint transaction: %w", err)
	}

	// Always revoke mint/burn roles from timelock after minting
	e.Logger.Infow("Adding revoke mint and burn roles to proposal",
		"chain", cfg.Selector,
		"timelock", timelockAddr.Hex(),
	)
	_, err = linkToken.RevokeMintRole(opts, timelockAddr)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to prepare revoke mint role transaction: %w", err)
	}
	_, err = linkToken.RevokeBurnRole(opts, timelockAddr)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to prepare revoke burn role transaction: %w", err)
	}

	// Enact returns the MCMS proposal in the ChangesetOutput
	output, err := deployerGroup.Enact()
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to enact deployer group: %w", err)
	}

	e.Logger.Infow("Successfully generated MCMS proposal to mint LINK tokens",
		"chain", cfg.Selector,
		"to", cfg.ToAddress.Hex(),
		"amount", cfg.Amount.String(),
		"numProposals", len(output.MCMSTimelockProposals),
	)

	return output, nil
}

package changeset

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	evmstate "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"

	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
)

var (
	GrantMintRoleAndMint = cldf.CreateChangeSet(GrantMintRoleAndMintLogic, ValidatePreConditions)
)

type GrantMintRoleAndMintConfig struct {
	Selector  uint64         `json:"selector"`
	ToAddress common.Address `json:"mintToAddress"`
	Amount    *big.Int       `json:"amount"`
}

func (cfg GrantMintRoleAndMintConfig) Validate(e cldf.Environment) error {
	if err := cldf.IsValidChainSelector(cfg.Selector); err != nil {
		return fmt.Errorf("invalid chain selector: %d - %w", cfg.Selector, err)
	}

	if cfg.ToAddress == (common.Address{}) {
		return errors.New("toAddress cannot be empty")
	}

	return nil
}

func ValidatePreConditions(e cldf.Environment, cfg GrantMintRoleAndMintConfig) error {
	if err := cfg.Validate(e); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}

	chainName := e.BlockChains.EVMChains()[cfg.Selector].Name()
	// The mintOnLinkToken should never happen on Mainnet
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

func GrantMintRoleAndMintLogic(e cldf.Environment, cfg GrantMintRoleAndMintConfig) (cldf.ChangesetOutput, error) {
	chain := e.BlockChains.EVMChains()[cfg.Selector]

	addresses, err := e.ExistingAddresses.AddressesForChain(cfg.Selector)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get addresses for chain %d: %w", cfg.Selector, err)
	}

	linkState, err := evmstate.MaybeLoadLinkTokenChainState(chain, addresses)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load LINK token state: %w", err)
	}

	// check if the owner is the deployer key and in that case grant mint access to the deployer key
	owner, err := linkState.LinkToken.Owner(&bind.CallOpts{Context: e.GetContext()})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get owner of token %s: %w", linkState.LinkToken.Address(), err)
	}
	if owner == chain.DeployerKey.From {
		//  Grant deployer address mint/burn access on the LINK_TOKEN
		e.Logger.Infow("Granting mint and burn roles to deployer", "chain", cfg.Selector, "deployer", chain.DeployerKey.From)
		tx, err := linkState.LinkToken.GrantMintAndBurnRoles(chain.DeployerKey, chain.DeployerKey.From)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to grant mint and burn roles: %w", err)
		}
		_, err = chain.Confirm(tx)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm grant roles transaction: %w", err)
		}
	}

	// Mint tokens to the given faucet address and verify the balance
	e.Logger.Infow("Minting tokens", "chain", cfg.Selector, "to", cfg.ToAddress, "amount", cfg.Amount.String())
	tx, err := linkState.LinkToken.Mint(chain.DeployerKey, cfg.ToAddress, cfg.Amount)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to mint tokens: %w", err)
	}
	_, err = chain.Confirm(tx)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm mint transaction: %w", err)
	}

	// Verify the balance
	balance, err := linkState.LinkToken.BalanceOf(&bind.CallOpts{}, cfg.ToAddress)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to check balance: %w", err)
	}

	if balance.Cmp(cfg.Amount) < 0 {
		return cldf.ChangesetOutput{}, fmt.Errorf("insufficient balance after minting: expected %s, got %s", cfg.Amount.String(), balance.String())
	}

	// Check if we need to revoke mint role
	isMinter, err := linkState.LinkToken.IsMinter(&bind.CallOpts{}, chain.DeployerKey.From)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to check if deployer is minter: %w", err)
	}

	if isMinter {
		tx, err = linkState.LinkToken.RevokeMintRole(chain.DeployerKey, chain.DeployerKey.From)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to revoke mint role: %w", err)
		}
		_, err = chain.Confirm(tx)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm revoke mint role transaction: %w", err)
		}
	}

	e.Logger.Infow("Successfully completed LINK token mint and ownership operations", "chain", cfg.Selector)

	return cldf.ChangesetOutput{}, nil
}

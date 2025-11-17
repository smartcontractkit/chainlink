package changeset

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	evmstate "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	opsutil "github.com/smartcontractkit/chainlink/deployment/common/opsutils"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"

	ccipops "github.com/smartcontractkit/chainlink/deployment/ccip/operation/evm"
	ccipseqs "github.com/smartcontractkit/chainlink/deployment/ccip/sequence/evm"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
)

var (
	GrantMintRole = cldf.CreateChangeSet(GrantMintRoleLogic, GrantMintRolePreConditions)
	// This changeset is specifically designed to use only in testnet and before transferring ownership of the LINK token to MCMS
	GrantMintRoleAndMint = cldf.CreateChangeSet(GrantMintRoleAndMintLogic, ValidatePreConditions)
)

type GrantMintRoleAndMintConfig struct {
	Selector  uint64         `json:"selector"`
	ToAddress common.Address `json:"mintToAddress"`
	Amount    *big.Int       `json:"amount"`
}

type GrantMintRoleInput struct {
	GrantMintRoleByChain map[uint64]GrantMintRoleConfig
	MCMS                 *proposalutils.TimelockConfig
}

type GrantMintRoleConfig struct {
	ToAddress common.Address `json:"toAddress"`
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
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get owner of token: %s: %w", linkState.LinkToken.Address(), err)
	}
	if owner == chain.DeployerKey.From {
		//  Grant deployer address mint/burn access on the LINK_TOKEN
		_, err := operations.ExecuteOperation(e.OperationsBundle, ccipops.GrantMintAndBurnRolesERC677Op, chain, opsutil.EVMCallInput[common.Address]{
			Address:       linkState.LinkToken.Address(),
			ChainSelector: chain.ChainSelector(),
			CallInput:     chain.DeployerKey.From,
		})
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to grant mint and burn roles: %w", err)
		}
	}

	// Mint tokens to the given faucet address and verify the balance
	e.Logger.Infow("Minting tokens", "chain", cfg.Selector, "to", cfg.ToAddress, "amount", cfg.Amount.String())
	_, err = operations.ExecuteOperation(e.OperationsBundle, ccipops.MintERC677Op, chain, opsutil.EVMCallInput[ccipops.MintERC677Config]{
		Address:       linkState.LinkToken.Address(),
		ChainSelector: chain.ChainSelector(),
		CallInput: ccipops.MintERC677Config{
			To:     cfg.ToAddress,
			Amount: cfg.Amount,
		},
	})
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

	// Check if we need to revoke mint/burn roles
	isMinter, err := linkState.LinkToken.IsMinter(&bind.CallOpts{}, chain.DeployerKey.From)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to check if deployer is minter: %w", err)
	}

	isBurner, err := linkState.LinkToken.IsBurner(&bind.CallOpts{}, chain.DeployerKey.From)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to check if deployer is burner: %w", err)
	}

	if isMinter {
		// Revoke Mint Role
		_, err = operations.ExecuteOperation(e.OperationsBundle, ccipops.RevokeMintRoleERC677Op, chain, opsutil.EVMCallInput[common.Address]{
			Address:       linkState.LinkToken.Address(),
			ChainSelector: chain.ChainSelector(),
			CallInput:     chain.DeployerKey.From,
		})
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm revoke mint role transaction: %w", err)
		}
	}

	if isBurner {
		// Revoke Burn Role
		_, err = operations.ExecuteOperation(e.OperationsBundle, ccipops.RevokeBurnRoleERC677Op, chain, opsutil.EVMCallInput[common.Address]{
			Address:       linkState.LinkToken.Address(),
			ChainSelector: chain.ChainSelector(),
			CallInput:     chain.DeployerKey.From,
		})
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm revoke Burn role transaction: %w", err)
		}
	}

	e.Logger.Infow("Successfully completed LINK token mint and ownership operations", "chain", cfg.Selector)

	return cldf.ChangesetOutput{}, nil
}

func GrantMintRolePreConditions(e cldf.Environment, input GrantMintRoleInput) error {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}

	for selector, cfg := range input.GrantMintRoleByChain {
		if err := cldf.IsValidChainSelector(selector); err != nil {
			return fmt.Errorf("invalid chain selector: %d - %w", selector, err)
		}

		if cfg.ToAddress == (common.Address{}) {
			return errors.New("toAddress cannot be empty")
		}

		chainState, ok := state.EVMChainState(selector)
		if !ok {
			return fmt.Errorf("%d does not exist in state", selector)
		}

		err = state.EnforceMCMSUsageIfProd(e.GetContext(), input.MCMS)
		if err != nil {
			return err
		}

		if linkToken := chainState.LinkToken; linkToken == nil {
			return fmt.Errorf("missing linkToken on %d", selector)
		}
	}

	return nil
}

func GrantMintRoleLogic(e cldf.Environment, input GrantMintRoleInput) (cldf.ChangesetOutput, error) {
	output := cldf.ChangesetOutput{}
	state, err := stateview.LoadOnchainState(e, stateview.WithLoadLegacyContracts(true))
	if err != nil {
		return output, fmt.Errorf("failed to load onchain state: %w", err)
	}

	seqReport, err := operations.ExecuteSequence(
		e.OperationsBundle,
		ccipseqs.GrantMintAndBurnRoleOnERC677Sequence,
		e.BlockChains.EVMChains(),
		input.ToSequenceInput(state),
	)

	return opsutil.AddEVMCallSequenceToCSOutput(
		e,
		cldf.ChangesetOutput{},
		seqReport,
		err,
		state.EVMMCMSStateByChain(),
		input.MCMS,
		ccipseqs.GrantMintAndBurnRoleOnERC677Sequence.Description(),
	)
}

func (input GrantMintRoleInput) ToSequenceInput(state stateview.CCIPOnChainState) ccipseqs.GrantMintRoleSeqInp {
	updates := make(map[uint64]opsutil.EVMCallInput[common.Address], len(input.GrantMintRoleByChain))
	for chainSel, cfg := range input.GrantMintRoleByChain {
		updates[chainSel] = opsutil.EVMCallInput[common.Address]{
			ChainSelector: chainSel,
			Address:       state.Chains[chainSel].LinkToken.Address(),
			CallInput:     cfg.ToAddress,
			NoSend:        input.MCMS != nil,
		}
	}

	return ccipseqs.GrantMintRoleSeqInp{
		UpdatesByChain: updates,
	}
}

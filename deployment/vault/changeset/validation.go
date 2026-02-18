package changeset

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	chainSel "github.com/smartcontractkit/chain-selectors"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/common/changeset"
	evmstate "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/vault/changeset/types"
)

func ValidateBatchNativeTransferConfig(ctx context.Context, e cldf.Environment, cfg types.BatchNativeTransferConfig) error {
	if len(cfg.TransfersByChain) == 0 {
		return errors.New("transfers_by_chain must not be empty")
	}

	for chainSelector, transfers := range cfg.TransfersByChain {
		if err := validateChainSelector(chainSelector, e); err != nil {
			return fmt.Errorf("invalid chain selector %d: %w", chainSelector, err)
		}

		if len(transfers) == 0 {
			return fmt.Errorf("chain %d has no transfers", chainSelector)
		}

		if err := validateNativeTransfers(ctx, e, chainSelector, transfers); err != nil {
			return fmt.Errorf("validation failed for chain %d: %w", chainSelector, err)
		}
	}

	if cfg.MCMSConfig != nil {
		if err := validateMCMSConfig(e, cfg.MCMSConfig, cfg.TransfersByChain, cfg.TimelockIDByChain); err != nil {
			return fmt.Errorf("MCMS configuration validation failed: %w", err)
		}
	}

	// When TimelockIDByChain is set, ensure each (chain, qualifier) has timelock and proposer
	if len(cfg.TimelockIDByChain) > 0 {
		for chainSelector, qualifier := range cfg.TimelockIDByChain {
			if _, hasChain := cfg.TransfersByChain[chainSelector]; !hasChain {
				continue
			}
			if _, err := GetContractAddressWithQualifier(e.DataStore, chainSelector, commontypes.RBACTimelock, qualifier); err != nil {
				return fmt.Errorf("timelock not found for chain %d qualifier %q: %w", chainSelector, qualifier, err)
			}
			if _, err := GetContractAddressWithQualifier(e.DataStore, chainSelector, commontypes.ProposerManyChainMultisig, qualifier); err != nil {
				return fmt.Errorf("proposer not found for chain %d qualifier %q: %w", chainSelector, qualifier, err)
			}
		}
	}

	return nil
}

func validateChainSelector(chainSelector uint64, e cldf.Environment) error {
	if len(e.BlockChains.EVMChains()) == 0 {
		return nil
	}

	family, err := chainSel.GetSelectorFamily(chainSelector)
	if err != nil {
		return fmt.Errorf("unknown chain selector: %w", err)
	}

	if family != chainSel.FamilyEVM {
		return fmt.Errorf("only EVM chains are supported, got family: %s", family)
	}

	_, exists := e.BlockChains.EVMChains()[chainSelector]
	if !exists {
		return fmt.Errorf("chain %d not found in environment", chainSelector)
	}

	return nil
}

func validateNativeTransfers(_ context.Context, e cldf.Environment, chainSelector uint64, transfers []types.NativeTransfer) error {
	whitelistedAddresses, err := GetWhitelistedAddresses(e, []uint64{chainSelector})
	if err != nil {
		return fmt.Errorf("failed to get whitelisted addresses for chain %d: %w", chainSelector, err)
	}

	whitelist := make(map[string]bool)
	for _, entry := range whitelistedAddresses[chainSelector] {
		whitelist[common.HexToAddress(entry.Address).Hex()] = true
	}

	totalAmount := big.NewInt(0)
	addressSet := make(map[string]bool)

	for i, transfer := range transfers {
		recipientAddress := common.HexToAddress(transfer.To)
		if recipientAddress == (common.Address{}) {
			return fmt.Errorf("transfer %d: 'to' address cannot be zero address", i)
		}

		if transfer.Amount == nil || transfer.Amount.Cmp(big.NewInt(0)) <= 0 {
			return fmt.Errorf("transfer %d: amount must be positive", i)
		}

		if addressSet[recipientAddress.Hex()] {
			return fmt.Errorf("transfer %d: duplicate destination address %s", i, recipientAddress.Hex())
		}
		addressSet[recipientAddress.Hex()] = true

		if !whitelist[recipientAddress.Hex()] {
			return fmt.Errorf("transfer %d: address %s is not whitelisted for chain %d", i, recipientAddress.Hex(), chainSelector)
		}

		totalAmount.Add(totalAmount, transfer.Amount)
	}

	if err := validateTimelockBalance(e, chainSelector, totalAmount); err != nil {
		return fmt.Errorf("timelock balance validation failed: %w", err)
	}

	return nil
}

func validateTimelockBalance(e cldf.Environment, chainSelector uint64, requiredAmount *big.Int) error {
	balances, err := GetTimelockBalances(e, []uint64{chainSelector})
	if err != nil {
		return fmt.Errorf("failed to get timelock balance for chain %d: %w", chainSelector, err)
	}

	balanceInfo, exists := balances[chainSelector]
	if !exists {
		return fmt.Errorf("timelock balance not found for chain %d", chainSelector)
	}

	if balanceInfo.Balance.Cmp(requiredAmount) < 0 {
		return fmt.Errorf("insufficient timelock balance: required %s wei, available %s wei",
			requiredAmount.String(), balanceInfo.Balance.String())
	}

	return nil
}

func validateMCMSConfig(e cldf.Environment, mcmsConfig *proposalutils.TimelockConfig, transfersByChain map[uint64][]types.NativeTransfer, timelockIDByChain map[uint64]string) error {
	if mcmsConfig != nil {
		if mcmsConfig.MinDelay < 0 {
			return fmt.Errorf("MCMS minimum delay cannot be negative: %d", mcmsConfig.MinDelay)
		}
	}
	for chainSelector := range transfersByChain {
		qualifier := ""
		if timelockIDByChain != nil {
			qualifier = timelockIDByChain[chainSelector]
		}
		addresses, err := evmstate.LoadAddressesFromDataStore(e.DataStore, chainSelector, qualifier)
		if err != nil {
			return fmt.Errorf("failed to get addresses from datastore for chain %d: %w", chainSelector, err)
		}

		_, err = GetContractAddressWithQualifier(e.DataStore, chainSelector, commontypes.RBACTimelock, qualifier)
		if err != nil {
			return fmt.Errorf("timelock not found for chain %d: %w", chainSelector, err)
		}

		_, err = GetContractAddressWithQualifier(e.DataStore, chainSelector, commontypes.ProposerManyChainMultisig, qualifier)
		if err != nil {
			return fmt.Errorf("proposer not found for chain %d: %w", chainSelector, err)
		}

		chain := e.BlockChains.EVMChains()[chainSelector]
		_, err = changeset.MaybeLoadMCMSWithTimelockChainState(chain, addresses)
		if err != nil {
			return fmt.Errorf("failed to load MCMS state for chain %d: %w", chainSelector, err)
		}
	}

	return nil
}

func ValidateFundTimelockConfig(ctx context.Context, e cldf.Environment, cfg types.FundTimelockConfig) error {
	if len(cfg.FundingByChain) == 0 {
		return errors.New("funding_by_chain must not be empty")
	}

	for chainSelector, amount := range cfg.FundingByChain {
		if err := validateChainSelector(chainSelector, e); err != nil {
			return fmt.Errorf("invalid chain selector %d: %w", chainSelector, err)
		}

		if amount == nil || amount.Cmp(big.NewInt(0)) <= 0 {
			return fmt.Errorf("funding amount for chain %d must be positive", chainSelector)
		}

		chain, exists := e.BlockChains.EVMChains()[chainSelector]
		if exists {
			deployerAddr := chain.DeployerKey.From
			balance, err := chain.Client.BalanceAt(ctx, deployerAddr, nil)
			if err != nil {
				return fmt.Errorf("failed to get deployer balance for chain %d: %w", chainSelector, err)
			}

			if balance.Cmp(amount) < 0 {
				return fmt.Errorf("insufficient deployer balance for chain %d: required %s wei, available %s wei",
					chainSelector, amount.String(), balance.String())
			}
		}
	}

	return nil
}

func ValidateSetWhitelistConfig(e cldf.Environment, cfg types.SetWhitelistConfig) error {
	useLegacy := len(cfg.WhitelistByChain) > 0
	useMulti := len(cfg.WhitelistByChainAndTimelock) > 0
	if !useLegacy && !useMulti {
		return errors.New("either whitelist_by_chain or whitelist_by_chain_and_timelock must be non-empty")
	}
	if useLegacy && useMulti {
		return errors.New("cannot set both whitelist_by_chain and whitelist_by_chain_and_timelock in the same config")
	}

	if useLegacy {
		for chainSelector, addresses := range cfg.WhitelistByChain {
			if err := validateChainSelector(chainSelector, e); err != nil {
				return fmt.Errorf("invalid chain selector %d: %w", chainSelector, err)
			}
			if err := validateWhitelistAddresses(chainSelector, addresses); err != nil {
				return err
			}
		}
	}

	if useMulti {
		for chainSelector, byTimelock := range cfg.WhitelistByChainAndTimelock {
			if err := validateChainSelector(chainSelector, e); err != nil {
				return fmt.Errorf("invalid chain selector %d: %w", chainSelector, err)
			}
			for timelockID, addresses := range byTimelock {
				if err := validateWhitelistAddresses(chainSelector, addresses); err != nil {
					return fmt.Errorf("chain %d timelock %q: %w", chainSelector, timelockID, err)
				}
			}
		}
	}

	return nil
}

func ValidateTransferERC20Config(e cldf.Environment, cfg types.TransferERC20Config) error {
	if len(cfg.Transfers) == 0 {
		return errors.New("transfers must not be empty")
	}
	if err := validateChainSelector(cfg.ChainSelector, e); err != nil {
		return fmt.Errorf("invalid chain selector: %w", err)
	}
	_, exists := e.BlockChains.EVMChains()[cfg.ChainSelector]
	if !exists {
		return fmt.Errorf("chain %d not found in environment", cfg.ChainSelector)
	}
	if cfg.MCMSConfig == nil {
		return errors.New("MCMSConfig is required for transfer_erc20")
	}
	for i, tr := range cfg.Transfers {
		if tr.Payee == "" || tr.Payee == "0x0000000000000000000000000000000000000000" {
			return fmt.Errorf("transfer %d: payee cannot be zero address", i)
		}
		if tr.Token == "" || tr.Token == "0x0000000000000000000000000000000000000000" {
			return fmt.Errorf("transfer %d: token address cannot be zero", i)
		}
		if tr.Amount == nil || tr.Amount.Cmp(big.NewInt(0)) <= 0 {
			return fmt.Errorf("transfer %d: amount must be positive", i)
		}
	}
	// Validate timelock and proposer exist for this chain and qualifier
	qualifier := cfg.TimelockIdentifier
	if _, err := GetContractAddressWithQualifier(e.DataStore, cfg.ChainSelector, commontypes.RBACTimelock, qualifier); err != nil {
		return fmt.Errorf("timelock not found for chain %d: %w", cfg.ChainSelector, err)
	}
	if _, err := GetContractAddressWithQualifier(e.DataStore, cfg.ChainSelector, commontypes.ProposerManyChainMultisig, qualifier); err != nil {
		return fmt.Errorf("proposer not found for chain %d: %w", cfg.ChainSelector, err)
	}
	return nil
}

func validateWhitelistAddresses(chainSelector uint64, addresses []types.WhitelistAddress) error {
	addressSet := make(map[string]bool)
	for i, addr := range addresses {
		if addr.Address == "" || addr.Address == "0x0000000000000000000000000000000000000000" {
			return fmt.Errorf("address %d: address cannot be zero address", i)
		}
		if addressSet[addr.Address] {
			return fmt.Errorf("duplicate address %s", addr.Address)
		}
		addressSet[addr.Address] = true
	}
	return nil
}

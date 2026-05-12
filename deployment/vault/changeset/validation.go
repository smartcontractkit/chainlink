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
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
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
		if err := validateMCMSConfig(e, cfg.MCMSConfig, cfg.TransfersByChain); err != nil {
			return fmt.Errorf("MCMS configuration validation failed: %w", err)
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

func validateMCMSConfig(e cldf.Environment, mcmsConfig *proposalutils.TimelockConfig, transfersByChain map[uint64][]types.NativeTransfer) error {
	if mcmsConfig != nil {
		if mcmsConfig.MinDelay < 0 {
			return fmt.Errorf("MCMS minimum delay cannot be negative: %d", mcmsConfig.MinDelay)
		}
	}
	const emptyQualifier = ""
	for chainSelector := range transfersByChain {
		addresses, err := state.GetAddressTypeVersionByQualifier(e.DataStore.Addresses(), chainSelector, emptyQualifier)
		if err != nil {
			return fmt.Errorf("failed to get addresses from datastore for chain %d: %w", chainSelector, err)
		}

		_, err = GetContractAddress(e.DataStore, chainSelector, commontypes.RBACTimelock)
		if err != nil {
			return fmt.Errorf("timelock not found for chain %d: %w", chainSelector, err)
		}

		_, err = GetContractAddress(e.DataStore, chainSelector, commontypes.ProposerManyChainMultisig)
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
	if len(cfg.WhitelistByChain) == 0 {
		return errors.New("whitelist_by_chain must not be empty")
	}

	for chainSelector, addresses := range cfg.WhitelistByChain {
		if err := validateChainSelector(chainSelector, e); err != nil {
			return fmt.Errorf("invalid chain selector %d: %w", chainSelector, err)
		}

		addressSet := make(map[string]bool)
		for i, addr := range addresses {
			if addr.Address == "" || addr.Address == "0x0000000000000000000000000000000000000000" {
				return fmt.Errorf("chain %d, address %d: address cannot be zero address", chainSelector, i)
			}

			// Check for duplicate addresses within the same chain
			if addressSet[addr.Address] {
				return fmt.Errorf("chain %d: duplicate address %s", chainSelector, addr.Address)
			}
			addressSet[addr.Address] = true
		}
	}

	return nil
}

func validateEthAddress(field, raw string) error {
	if raw == "" {
		return fmt.Errorf("%s must not be empty", field)
	}
	if !common.IsHexAddress(raw) {
		return fmt.Errorf("%s is not a valid hex address: %s", field, raw)
	}
	return nil
}

func ValidateDeployEthBalMonConfig(ctx context.Context, env cldf.Environment, cfg types.DeployEthBalMonInput) error {
	if len(cfg.Chains) == 0 {
		return errors.New("chains must not be empty")
	}

	for chainSelector, chainCfg := range cfg.Chains {
		if err := validateChainSelector(chainSelector, env); err != nil {
			return fmt.Errorf("chain %d: %w", chainSelector, err)
		}
		if chainCfg.SetKeeperRegistryAddress == "" {
			return fmt.Errorf("chain %d: setKeeperRegistryAddress must not be empty", chainSelector)
		}
		if err := validateEthAddress("setKeeperRegistryAddress", chainCfg.SetKeeperRegistryAddress); err != nil {
			return fmt.Errorf("chain %d: %w", chainSelector, err)
		}
	}

	return nil
}

func ValidateSetKeeperRegistryAddressConfig(ctx context.Context, env cldf.Environment, cfg types.EthBalMonSetKeeperRegistryAddressInput) error {
	if len(cfg.Chains) == 0 {
		return errors.New("no chains provided")
	}

	for chainSelector, chainConfig := range cfg.Chains {
		if _, ok := env.BlockChains.EVMChains()[chainSelector]; !ok {
			return fmt.Errorf("chain not found in environment: %d", chainSelector)
		}

		if !common.IsHexAddress(chainConfig.NewKeeperRegistryAddress) {
			return fmt.Errorf("chain %d: invalid keeper registry address: %s", chainSelector, chainConfig.NewKeeperRegistryAddress)
		}
	}

	return nil
}

func ValidateEthBalMonWithdrawConfig(ctx context.Context, env cldf.Environment, cfg types.EthBalMonWithdrawInput) error {
	if len(cfg.Chains) == 0 {
		return errors.New("no chains provided")
	}

	for chainSelector, chainConfig := range cfg.Chains {
		if _, ok := env.BlockChains.EVMChains()[chainSelector]; !ok {
			return fmt.Errorf("chain not found in environment: %d", chainSelector)
		}
		if chainConfig.Amount == nil || chainConfig.Amount.Cmp(big.NewInt(0)) <= 0 {
			return fmt.Errorf("chain %d: amount to withdraw must be positive", chainSelector)
		}

		if err := validateEthAddress("payeer", chainConfig.Payeer); err != nil {
			return fmt.Errorf("chain %d: %w", chainSelector, err)
		}
		if common.HexToAddress(chainConfig.Payeer) == (common.Address{}) {
			return fmt.Errorf("chain %d: payeer address cannot be zero address", chainSelector)
		}
	}

	return nil
}

func ValidateEthBalMonTransferOwnershipConfig(ctx context.Context, env cldf.Environment, cfg types.EthBalMonTransferOwnershipInput) error {
	if len(cfg.Chains) == 0 {
		return errors.New("no chains provided")
	}

	for chainSelector, chainConfig := range cfg.Chains {
		if _, ok := env.BlockChains.EVMChains()[chainSelector]; !ok {
			return fmt.Errorf("chain not found in environment: %d", chainSelector)
		}
		if err := validateEthAddress("newOwner", chainConfig.NewOwner); err != nil {
			return fmt.Errorf("chain %d: %w", chainSelector, err)
		}
		if common.HexToAddress(chainConfig.NewOwner) == (common.Address{}) {
			return fmt.Errorf("chain %d: newOwner address cannot be zero address", chainSelector)
		}
	}

	return nil
}

func ValidateEthBalMonSetWatchListConfig(ctx context.Context, env cldf.Environment, cfg types.EthBalMonSetWatchListInput) error {
	if len(cfg.Chains) == 0 {
		return errors.New("no chains provided")
	}

	for chainSelector, chainConfig := range cfg.Chains {
		if _, ok := env.BlockChains.EVMChains()[chainSelector]; !ok {
			return fmt.Errorf("chain not found in environment: %d", chainSelector)
		}
		n := len(chainConfig.Addresses)
		if n == 0 {
			return fmt.Errorf("chain %d: addresses must not be empty", chainSelector)
		}
		if len(chainConfig.MinBalancesWei) != n || len(chainConfig.TopUpAmountsWei) != n {
			return fmt.Errorf(
				"chain %d: addresses, min_balance_wei, and topup_amounts_wei must have the same length (got %d, %d, %d)",
				chainSelector, n, len(chainConfig.MinBalancesWei), len(chainConfig.TopUpAmountsWei),
			)
		}
		for i, addr := range chainConfig.Addresses {
			if err := validateEthAddress(fmt.Sprintf("address %d", i), addr.Hex()); err != nil {
				return fmt.Errorf("chain %d: %w", chainSelector, err)
			}
			if addr == (common.Address{}) {
				return fmt.Errorf("chain %d: address at index %d is zero address", chainSelector, i)
			}
		}
	}

	return nil
}

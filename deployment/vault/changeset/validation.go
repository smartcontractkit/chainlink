package changeset

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	chainSel "github.com/smartcontractkit/chain-selectors"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cldfproposalutils "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils"

	evmstate "github.com/smartcontractkit/cld-changesets/legacy/pkg/family/evm"

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

func validateMCMSConfig(e cldf.Environment, mcmsConfig *cldfproposalutils.TimelockConfig, transfersByChain map[uint64][]types.NativeTransfer) error {
	if mcmsConfig != nil {
		if mcmsConfig.MinDelay < 0 {
			return fmt.Errorf("MCMS minimum delay cannot be negative: %d", mcmsConfig.MinDelay)
		}
	}
	const emptyQualifier = ""
	for chainSelector := range transfersByChain {
		addresses, err := evmstate.GetAddressTypeVersionByQualifier(e.DataStore.Addresses(), chainSelector, emptyQualifier)
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

		_, err = GetContractAddress(e.DataStore, chainSelector, commontypes.BypasserManyChainMultisig)
		if err != nil {
			return fmt.Errorf("bypasser not found for chain %d: %w", chainSelector, err)
		}

		chain := e.BlockChains.EVMChains()[chainSelector]
		_, err = evmstate.MaybeLoadMCMSWithTimelockChainState(chain, addresses)
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

func validateERC20Transfers(e cldf.Environment, chainSelector uint64, transfers []types.ERC20Transfer) error {
	whitelistedAddresses, err := GetWhitelistedAddresses(e, []uint64{chainSelector})
	if err != nil {
		return fmt.Errorf("failed to get whitelisted addresses for chain %d: %w", chainSelector, err)
	}

	whitelist := make(map[string]bool)
	for _, entry := range whitelistedAddresses[chainSelector] {
		whitelist[common.HexToAddress(entry.Address).Hex()] = true
	}

	for i, tr := range transfers {
		if tr.Payee == "" || tr.Payee == "0x0000000000000000000000000000000000000000" {
			return fmt.Errorf("transfer %d: payee cannot be zero address", i)
		}
		if tr.Token == "" || tr.Token == "0x0000000000000000000000000000000000000000" {
			return fmt.Errorf("transfer %d: token address cannot be zero", i)
		}
		if tr.Amount == nil || tr.Amount.Cmp(big.NewInt(0)) <= 0 {
			return fmt.Errorf("transfer %d: amount must be positive", i)
		}

		payeeAddress := common.HexToAddress(tr.Payee)
		if !whitelist[payeeAddress.Hex()] {
			return fmt.Errorf("transfer %d: payee %s is not whitelisted for chain %d", i, payeeAddress.Hex(), chainSelector)
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

	if cfg.MCMSConfig != nil && cfg.MCMSConfig.MinDelay < 0 {
		return fmt.Errorf("MCMS minimum delay cannot be negative: %d", cfg.MCMSConfig.MinDelay)
	}

	for chainSelector, chainCfg := range cfg.Chains {
		if err := validateChainSelector(chainSelector, env); err != nil {
			return fmt.Errorf("chain %d: %w", chainSelector, err)
		}
		if err := validateEthAddress("keeperRegistryAddress", chainCfg.KeeperRegistryAddress); err != nil {
			return fmt.Errorf("chain %d: %w", chainSelector, err)
		}
		if common.HexToAddress(chainCfg.KeeperRegistryAddress) == (common.Address{}) {
			return fmt.Errorf("chain %d: keeperRegistryAddress cannot be zero address", chainSelector)
		}
		if err := validateDeployEthBalMonMCMSInDatastore(env, chainSelector, cfg.MCMSConfig); err != nil {
			return fmt.Errorf("chain %d: %w", chainSelector, err)
		}
	}

	return nil
}

// validateDeployEthBalMonMCMSInDatastore ensures RBACTimelock, the MCM used for the post-deploy
// accept-ownership proposal (bypasser vs proposer per cfg.MCMSConfig), and loadable MCMS state
// exist in the datastore — matching DeployEthBalMonSequence and BuildAcceptOwnershipTimelockProposal.
func validateDeployEthBalMonMCMSInDatastore(e cldf.Environment, chainSelector uint64, mcmsCfg *cldfproposalutils.TimelockConfig) error {
	const emptyQualifier = ""
	addresses, err := evmstate.GetAddressTypeVersionByQualifier(e.DataStore.Addresses(), chainSelector, emptyQualifier)
	if err != nil {
		return fmt.Errorf("failed to get addresses from datastore: %w", err)
	}

	_, err = GetContractAddress(e.DataStore, chainSelector, commontypes.RBACTimelock)
	if err != nil {
		return fmt.Errorf("timelock not found in datastore: %w", err)
	}

	mcmType := ethBalMonMCMSContractTypeForAction(deployEthBalMonAcceptOwnershipMCMSAction(mcmsCfg))
	_, err = GetContractAddress(e.DataStore, chainSelector, mcmType)
	if err != nil {
		return fmt.Errorf("MCMS (%s) not found in datastore: %w", mcmType, err)
	}

	chain := e.BlockChains.EVMChains()[chainSelector]
	_, err = evmstate.MaybeLoadMCMSWithTimelockChainState(chain, addresses)
	if err != nil {
		return fmt.Errorf("failed to load MCMS with timelock state: %w", err)
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

		if err := validateEthAddress("new_keeper_registry_address", chainConfig.NewKeeperRegistryAddress); err != nil {
			return fmt.Errorf("chain %d: %w", chainSelector, err)
		}
		if common.HexToAddress(chainConfig.NewKeeperRegistryAddress) == (common.Address{}) {
			return fmt.Errorf("chain %d: keeper registry address cannot be zero address", chainSelector)
		}
	}

	return nil
}

func ValidateTransferERC20Config(ctx context.Context, e cldf.Environment, cfg types.TransferERC20Config) error {
	if len(cfg.TransfersByChain) == 0 {
		return errors.New("transfers_by_chain must not be empty")
	}
	if cfg.MCMSConfig == nil {
		return errors.New("MCMSConfig is required for transfer_erc20")
	}

	// "" is the current default qualifier for the timelock
	for chainSelector, transfers := range cfg.TransfersByChain {
		qualifier := cfg.MCMSConfig.TimelockQualifierPerChain[chainSelector]
		if err := validateChainSelector(chainSelector, e); err != nil {
			return fmt.Errorf("invalid chain selector %d: %w", chainSelector, err)
		}

		if len(transfers) == 0 {
			return fmt.Errorf("chain %d has no transfers", chainSelector)
		}

		if err := validateERC20Transfers(e, chainSelector, transfers); err != nil {
			return fmt.Errorf("validation failed for chain %d: %w", chainSelector, err)
		}

		if _, err := GetContractAddressWithQualifier(e.DataStore, chainSelector, commontypes.RBACTimelock, qualifier); err != nil {
			return fmt.Errorf("timelock not found for chain %d: %w", chainSelector, err)
		}
		if _, err := GetContractAddressWithQualifier(e.DataStore, chainSelector, commontypes.ProposerManyChainMultisig, qualifier); err != nil {
			return fmt.Errorf("proposer not found for chain %d: %w", chainSelector, err)
		}
	}

	if err := validateERC20MCMSConfig(e, cfg.MCMSConfig, cfg.TransfersByChain, cfg.MCMSConfig.TimelockQualifierPerChain); err != nil {
		return fmt.Errorf("MCMS configuration validation failed: %w", err)
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

		if err := validateEthAddress("payee", chainConfig.Payee); err != nil {
			return fmt.Errorf("chain %d: %w", chainSelector, err)
		}
		if common.HexToAddress(chainConfig.Payee) == (common.Address{}) {
			return fmt.Errorf("chain %d: payee address cannot be zero address", chainSelector)
		}
	}

	return nil
}

func validateERC20MCMSConfig(e cldf.Environment, mcmsConfig *cldfproposalutils.TimelockConfig, transfersByChain map[uint64][]types.ERC20Transfer, qualifierMap map[uint64]string) error {
	if mcmsConfig.MinDelay < 0 {
		return fmt.Errorf("MCMS minimum delay cannot be negative: %d", mcmsConfig.MinDelay)
	}

	for chainSelector := range transfersByChain {
		addresses, err := evmstate.GetAddressTypeVersionByQualifier(e.DataStore.Addresses(), chainSelector, qualifierMap[chainSelector])
		if err != nil {
			return fmt.Errorf("failed to get addresses from datastore for chain %d: %w", chainSelector, err)
		}

		if _, err := GetContractAddressWithQualifier(e.DataStore, chainSelector, commontypes.RBACTimelock, qualifierMap[chainSelector]); err != nil {
			return fmt.Errorf("timelock not found for chain %d: %w", chainSelector, err)
		}

		if _, err := GetContractAddressWithQualifier(e.DataStore, chainSelector, commontypes.ProposerManyChainMultisig, qualifierMap[chainSelector]); err != nil {
			return fmt.Errorf("proposer not found for chain %d: %w", chainSelector, err)
		}

		chain := e.BlockChains.EVMChains()[chainSelector]
		_, err = evmstate.MaybeLoadMCMSWithTimelockChainState(chain, addresses)
		if err != nil {
			return fmt.Errorf("failed to load MCMS state for chain %d: %w", chainSelector, err)
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

func ValidateEthBalMonAcceptOwnershipConfig(ctx context.Context, env cldf.Environment, cfg types.EthBalMonAcceptOwnershipInput) error {
	if len(cfg.Chains) == 0 {
		return errors.New("no chains provided")
	}

	for _, chainSelector := range cfg.Chains {
		if _, ok := env.BlockChains.EVMChains()[chainSelector]; !ok {
			return fmt.Errorf("chain not found in environment: %d", chainSelector)
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
			// Check MinBalancesWei and TopUpAmountsWei are >= 0
			if chainConfig.MinBalancesWei[i].Cmp(big.NewInt(0)) < 0 {
				return fmt.Errorf("chain %d: min_balance_wei at index %d must be >= 0", chainSelector, i)
			}
			if chainConfig.TopUpAmountsWei[i].Cmp(big.NewInt(0)) < 0 {
				return fmt.Errorf("chain %d: topup_amounts_wei at index %d must be >= 0", chainSelector, i)
			}
		}
	}

	return nil
}

//***********************************************//
//           CCIP Bridge Validation              //
//***********************************************//

// ValidateCCIPBridgeConfig validates a CCIPBridgeConfig before any proposal is built. It checks the
// static config, that the source chain + its Timelock/MCM/WETH/Router exist (and WETH/Router have
// code), and that every destination has an EtherSenderReceiver ref and a whitelisted recipient.
// Live fee quoting and Timelock balance checks happen in the sequence, which needs on-chain reads.
func ValidateCCIPBridgeConfig(ctx context.Context, e cldf.Environment, cfg types.CCIPBridgeConfig) error {
	if cfg.TokenType != types.TokenTypeNative && cfg.TokenType != types.TokenTypeWrappedNative {
		return fmt.Errorf("token_type must be %q or %q, got %q", types.TokenTypeNative, types.TokenTypeWrappedNative, cfg.TokenType)
	}
	if len(cfg.Transfers) == 0 {
		return errors.New("at least one transfer is required")
	}
	if cfg.MCMSConfig == nil {
		return errors.New("MCMSConfig is required")
	}

	evmChains := e.BlockChains.EVMChains()
	sourceChain, ok := evmChains[cfg.FromChainSelector]
	if !ok {
		return fmt.Errorf("source chain %d not found in environment", cfg.FromChainSelector)
	}

	if _, err := GetContractAddress(e.DataStore, cfg.FromChainSelector, commontypes.RBACTimelock); err != nil {
		return fmt.Errorf("source timelock: %w", err)
	}
	mcmType := commontypes.ProposerManyChainMultisig
	if cfg.MCMSConfig.MCMSAction == mcmstypes.TimelockActionBypass {
		mcmType = commontypes.BypasserManyChainMultisig
	}
	if _, err := GetContractAddress(e.DataStore, cfg.FromChainSelector, mcmType); err != nil {
		return fmt.Errorf("source MCM (%s): %w", mcmType, err)
	}
	wethStr, err := GetContractAddress(e.DataStore, cfg.FromChainSelector, WETHContractType)
	if err != nil {
		return fmt.Errorf("source WETH: %w", err)
	}
	routerStr, err := GetContractAddress(e.DataStore, cfg.FromChainSelector, RouterContractType)
	if err != nil {
		return fmt.Errorf("source Router: %w", err)
	}
	if err := assertHasCode(ctx, sourceChain.Client, wethStr, "WETH"); err != nil {
		return err
	}
	if err := assertHasCode(ctx, sourceChain.Client, routerStr, "Router"); err != nil {
		return err
	}

	seenDest := make(map[uint64]bool, len(cfg.Transfers))
	for i, t := range cfg.Transfers {
		if t.ToChainSelector == cfg.FromChainSelector {
			return fmt.Errorf("transfer %d: destination selector %d equals source", i, t.ToChainSelector)
		}
		if seenDest[t.ToChainSelector] {
			return fmt.Errorf("transfer %d: duplicate destination selector %d", i, t.ToChainSelector)
		}
		seenDest[t.ToChainSelector] = true

		if err := validateEthAddress("destination_address", t.DestinationAddress); err != nil {
			return fmt.Errorf("transfer %d: %w", i, err)
		}
		if common.HexToAddress(t.DestinationAddress) == (common.Address{}) {
			return fmt.Errorf("transfer %d: destination_address cannot be zero", i)
		}
		if t.Amount == nil || t.Amount.Sign() <= 0 {
			return fmt.Errorf("transfer %d: amount must be positive", i)
		}
		if cfg.TestNeeded && t.Amount.Cmp(TestBridgeAmount) <= 0 {
			return fmt.Errorf("transfer %d: amount %s must exceed test amount %s when test_needed is set", i, t.Amount, TestBridgeAmount)
		}
		if t.MaxCCIPFee != nil && t.MaxCCIPFee.Sign() < 0 {
			return fmt.Errorf("transfer %d: max_ccip_fee cannot be negative", i)
		}

		esrStr, err := GetContractAddress(e.DataStore, t.ToChainSelector, EtherSenderReceiverContractType)
		if err != nil {
			return fmt.Errorf("transfer %d: EtherSenderReceiver for dest chain %d: %w", i, t.ToChainSelector, err)
		}
		if destChain, ok := evmChains[t.ToChainSelector]; ok {
			if err := assertHasCode(ctx, destChain.Client, esrStr, "EtherSenderReceiver"); err != nil {
				return fmt.Errorf("transfer %d: %w", i, err)
			}
		}

		whitelisted, err := GetWhitelistedAddresses(e, []uint64{t.ToChainSelector})
		if err != nil {
			return fmt.Errorf("transfer %d: whitelist lookup for chain %d: %w", i, t.ToChainSelector, err)
		}
		if !whitelistContains(whitelisted[t.ToChainSelector], t.DestinationAddress) {
			return fmt.Errorf("transfer %d: destination %s is not whitelisted on chain %d", i, common.HexToAddress(t.DestinationAddress).Hex(), t.ToChainSelector)
		}
	}

	return nil
}

// assertHasCode fails if there is no contract code at addr on the given client.
func assertHasCode(ctx context.Context, client interface {
	CodeAt(context.Context, common.Address, *big.Int) ([]byte, error)
}, addrStr, label string) error {
	addr := common.HexToAddress(addrStr)
	code, err := client.CodeAt(ctx, addr, nil)
	if err != nil {
		return fmt.Errorf("%s code check at %s: %w", label, addr.Hex(), err)
	}
	if len(code) == 0 {
		return fmt.Errorf("%s at %s has no code", label, addr.Hex())
	}
	return nil
}

// whitelistContains reports whether addr is in the whitelist entries.
func whitelistContains(entries []WhitelistEntry, addr string) bool {
	want := common.HexToAddress(addr).Hex()
	for _, entry := range entries {
		if common.HexToAddress(entry.Address).Hex() == want {
			return true
		}
	}
	return false
}

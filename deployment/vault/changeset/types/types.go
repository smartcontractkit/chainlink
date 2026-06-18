package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	cldfproposalutils "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils"
)

type NativeTransfer struct {
	To     string   `json:"to"`     // Destination address
	Amount *big.Int `json:"amount"` // Amount in wei
}

// BatchNativeTransferConfig configures batch native token transfers across multiple chains
type BatchNativeTransferConfig struct {
	// TransfersByChain maps chain selector to list of transfers for that chain
	TransfersByChain map[uint64][]NativeTransfer `json:"transfers_by_chain"`

	// MCMSConfig contains timelock and MCMS configuration
	MCMSConfig *cldfproposalutils.TimelockConfig `json:"mcms_config"`

	// Description for the MCMS proposal
	Description string `json:"description"`
}

// FundTimelockConfig configures funding timelock contracts with native tokens
type FundTimelockConfig struct {
	// FundingByChain maps chain selector to amount to fund the timelock
	FundingByChain map[uint64]*big.Int `json:"funding_by_chain"`
}

// WhitelistAddress represents an address entry in the whitelist
type WhitelistAddress struct {
	Address     string   `json:"address"`
	Description string   `json:"description"`
	Labels      []string `json:"labels"`
}

// SetWhitelistConfig configures address whitelist state
type SetWhitelistConfig struct {
	// WhitelistByChain maps chain selector to the list of whitelisted addresses for that chain
	WhitelistByChain map[uint64][]WhitelistAddress `json:"whitelist_by_chain"`
}

// WhitelistMetadata represents the whitelist state for a single chain stored in chain metadata
type WhitelistMetadata struct {
	Addresses []WhitelistAddress `json:"addresses"`
}

// TimelockNativeBalanceInfo represents native token balance information for Timelock
type TimelockNativeBalanceInfo struct {
	TimelockAddr string   `json:"timelock_address"`
	Balance      *big.Int `json:"balance"`
}

// TransferValidationError represents validation errors for transfers
type TransferValidationError struct {
	ChainSelector uint64 `json:"chain_selector"`
	Address       string `json:"address"`
	Error         string `json:"error"`
}

// BatchNativeTransferState represents the current state of Vault
type BatchNativeTransferState struct {
	// TimelockBalances maps chain selector to timelock balance info
	TimelockBalances map[uint64]*TimelockNativeBalanceInfo `json:"timelock_balances"`

	// WhitelistedAddresses maps chain selector to list of whitelisted addresses
	WhitelistedAddresses map[uint64][]string `json:"whitelisted_addresses"`

	// ValidationErrors contains any validation errors found
	ValidationErrors []TransferValidationError `json:"validation_errors"`
}

// ERC20Transfer is a single ERC20 transfer (payee, token, amount in token units)
type ERC20Transfer struct {
	Payee  string   `json:"payee"`  // Destination address
	Token  string   `json:"token"`  // ERC20 token contract address
	Amount *big.Int `json:"amount"` // Amount in token units (not wei)
}

// TransferERC20Config configures batch ERC20 transfers from timelocks across multiple chains
type TransferERC20Config struct {
	// TransfersByChain maps chain selector to ERC20 transfers for that chain
	TransfersByChain map[uint64][]ERC20Transfer `json:"transfers_by_chain"`

	// MCMSConfig contains timelock and MCMS configuration for building the proposal
	MCMSConfig *cldfproposalutils.TimelockConfig `json:"mcms_config"`

	// Description for the MCMS proposal
	Description string `json:"description"`
}

// DeployEthBalMonChainConfig is deployment-time configuration for EthBalMon on one chain.
type DeployEthBalMonChainConfig struct {
	// SetKeeperRegistryAddress is the Chainlink Automation registry forwarder (the upkeep
	// "forwarder address") on standard automation chains, or the KMS executor address when
	// using the Plaid/KMS automation path.
	SetKeeperRegistryAddress string `json:"setKeeperRegistryAddress"`
	// SetMinWaitPeriodSeconds is the minimum seconds between balance checks for this deployment.
	// Optional: nil or 0 means the deploy changeset uses a default (currently 60 seconds).
	SetMinWaitPeriodSeconds *uint64 `json:"setMinWaitPeriodSeconds,omitempty"`
}

// DeployEthBalMonInput is the input to the EthBalMon deploy changeset.
// Keys are chain selectors; each value configures keeper/registry wiring and min wait for that chain.
type DeployEthBalMonInput struct {
	Chains map[uint64]DeployEthBalMonChainConfig `json:"chains"`
	// MCMSConfig optionally configures the post-deploy accept-ownership timelock proposal (MCMS action, delay, etc.).
	// When nil or MCMSAction is unset, the deploy flow defaults accept-ownership to bypass (historical behavior).
	MCMSConfig *cldfproposalutils.TimelockConfig `json:"mcms_config,omitempty"`
}

// EthBalMonContractType is the datastore / MCMS contract type label for EthBalMon deployments.
const EthBalMonContractType = "EthBalMon"

// SetKeeperRegistryChainConfig updates the automation executor/registry EthBalMon forwards work to.
type SetKeeperRegistryChainConfig struct {
	// NewKeeperRegistryAddress is the new Chainlink Automation forwarder or KMS executor address (hex).
	NewKeeperRegistryAddress string `json:"new_keeper_registry_address"`
}

// EthBalMonSetKeeperRegistryAddressInput is the input to the setKeeperRegistryAddress changeset.
// Keys are chain selectors with the registry address to set on each chain's EthBalMon.
type EthBalMonSetKeeperRegistryAddressInput struct {
	Chains map[uint64]SetKeeperRegistryChainConfig `json:"chains"`
	// MCMSConfig optionally configures the timelock proposal; when nil, schedule + proposer MCM is used.
	MCMSConfig *cldfproposalutils.TimelockConfig `json:"mcms_config,omitempty"`
}

// EthBalMonSetWatchListChainConfig replaces the monitored addresses and thresholds on one chain.
// Addresses, MinBalancesWei, and TopUpAmountsWei are parallel slices: index i applies to Addresses[i].
// MinBalancesWei and TopUpAmountsWei are represented as *big.Int values in wei.
type EthBalMonSetWatchListChainConfig struct {
	Addresses       []common.Address `json:"addresses"`
	MinBalancesWei  []*big.Int       `json:"min_balance_wei"`
	TopUpAmountsWei []*big.Int       `json:"topup_amounts_wei"`
}

// EthBalMonSetWatchListInput is the input to the setWatchList changeset.
// Keys are chain selectors; each value is the full watch list to install on that chain's EthBalMon.
type EthBalMonSetWatchListInput struct {
	Chains map[uint64]EthBalMonSetWatchListChainConfig `json:"chains"`
	// MCMSConfig optionally configures the timelock proposal; when nil, schedule + proposer MCM is used.
	MCMSConfig *cldfproposalutils.TimelockConfig `json:"mcms_config,omitempty"`
}

// EthBalMonWithdrawChainConfig configures a native-token withdraw from EthBalMon on one chain.
type EthBalMonWithdrawChainConfig struct {
	// Amount is the withdrawal amount in wei. Must be positive (validated by the changeset).
	Amount *big.Int `json:"amount"`
	// Payee is the recipient address (hex).
	Payee string `json:"payee"`
}

// EthBalMonWithdrawInput is the input to the EthBalMon withdraw changeset.
// Keys are chain selectors; each value specifies amount and recipient for that chain.
type EthBalMonWithdrawInput struct {
	Chains map[uint64]EthBalMonWithdrawChainConfig `json:"chains"`
	// MCMSConfig optionally configures the timelock proposal; when nil, schedule + proposer MCM is used.
	MCMSConfig *cldfproposalutils.TimelockConfig `json:"mcms_config,omitempty"`
}

// EthBalMonTransferOwnershipChainConfig sets the new owner of EthBalMon on one chain.
type EthBalMonTransferOwnershipChainConfig struct {
	// NewOwner is the address (hex) that will own the EthBalMon contract after the operation.
	NewOwner string `json:"new_owner"`
}

// EthBalMonTransferOwnershipInput is the input to the EthBalMon transferOwnership changeset.
// Keys are chain selectors; each value is the new owner for that chain's EthBalMon instance.
type EthBalMonTransferOwnershipInput struct {
	Chains map[uint64]EthBalMonTransferOwnershipChainConfig `json:"chains"`
	// MCMSConfig optionally configures the timelock proposal; when nil, schedule + proposer MCM is used.
	MCMSConfig *cldfproposalutils.TimelockConfig `json:"mcms_config,omitempty"`
}

// EthBalMonAcceptOwnershipInput is the input to the EthBalMon acceptOwnership changeset.
// Chains is the list of chain selectors on which to call acceptOwnership on the EthBalMon instance.
type EthBalMonAcceptOwnershipInput struct {
	Chains []uint64 `json:"chains"`
	// MCMSConfig optionally configures the timelock proposal; when nil, schedule + proposer MCM is used.
	MCMSConfig *cldfproposalutils.TimelockConfig `json:"mcms_config,omitempty"`
}

//nolint:revive // types is a common package name
package types

import (
	"math/big"

	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
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
	MCMSConfig *proposalutils.TimelockConfig `json:"mcms_config"`

	// Description for the MCMS proposal
	Description string `json:"description"`

	// TimelockIDByChain optionally maps chain selector to timelock qualifier for multi-timelock.
	// When set, the timelock and proposer for each chain are resolved by this qualifier from the datastore.
	// Omit or use empty string for legacy single-timelock-per-chain behavior.
	TimelockIDByChain map[uint64]string `json:"timelock_id_by_chain,omitempty"`
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
	// WhitelistByChain maps chain selector to the list of whitelisted addresses for that chain (legacy, default timelock).
	// Use when only one timelock per chain or for the default/empty qualifier.
	WhitelistByChain map[uint64][]WhitelistAddress `json:"whitelist_by_chain,omitempty"`

	// WhitelistByChainAndTimelock maps chain -> timelock_id -> list of addresses for multi-timelock.
	// Use "" as timelock_id for the default/legacy timelock. When set, takes precedence over WhitelistByChain for that (chain, timelock).
	WhitelistByChainAndTimelock map[uint64]map[string][]WhitelistAddress `json:"whitelist_by_chain_and_timelock,omitempty"`
}

// WhitelistMetadata represents the whitelist state for a single chain stored in chain metadata
type WhitelistMetadata struct {
	Addresses []WhitelistAddress `json:"addresses"`
}

// VaultChainMetadata is the value stored in chain metadata for vault. It supports both legacy (Addresses only)
// and multi-timelock (ByTimelock keyed by timelock qualifier). Use "" for default/legacy timelock.
type VaultChainMetadata struct {
	Addresses   []WhitelistAddress            `json:"addresses,omitempty"`
	ByTimelock  map[string][]WhitelistAddress `json:"by_timelock,omitempty"`
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

// TransferERC20Config configures an ERC20 transfer from a timelock on a single chain
type TransferERC20Config struct {
	// ChainSelector is the chain where the timelock and token live
	ChainSelector uint64 `json:"chain_selector"`

	// TimelockIdentifier is the qualifier for the timelock (e.g. "vault_1"). Use "" for default/legacy.
	TimelockIdentifier string `json:"timelock_identifier"`

	// Transfers is the list of ERC20 transfers to execute from the timelock
	Transfers []ERC20Transfer `json:"transfers"`

	// MCMSConfig contains timelock and MCMS configuration for building the proposal
	MCMSConfig *proposalutils.TimelockConfig `json:"mcms_config"`

	// Description for the MCMS proposal
	Description string `json:"description"`
}

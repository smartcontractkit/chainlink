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

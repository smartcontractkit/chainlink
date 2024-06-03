package model

import (
	"fmt"
	"time"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

type CommitPluginConfig struct {
	// Writer indicates that the node can contribute by sending reports to the destination chain.
	// Being a Writer guarantees that the node can also read from the destination chain.
	Writer bool `json:"writer"`

	// Reads define the chains that the current node can read from.
	Reads []ChainSelector `json:"reads"`

	// DestChain is the ccip destination chain configured for the commit plugin DON.
	DestChain ChainSelector `json:"destChain"`

	// FChain defines the FChain value for each chain. FChain is used while forming consensus based on the observations.
	FChain map[ChainSelector]int `json:"fChain"`

	// ObserverInfo is a map of oracle IDs to ObserverInfo.
	ObserverInfo map[commontypes.OracleID]ObserverInfo `json:"observerInfo"`

	// PricedTokens is a list of tokens that we want to submit price updates for.
	PricedTokens []types.Account `json:"pricedTokens"`

	// TokenPricesObserver indicates that the node can observe token prices.
	TokenPricesObserver bool `json:"tokenPricesObserver"`

	// NewMsgScanBatchSize is the number of max new messages to scan, typically set to 256.
	NewMsgScanBatchSize int `json:"newMsgScanBatchSize"`
}

func (c CommitPluginConfig) Validate() error {
	if c.DestChain == ChainSelector(0) {
		return fmt.Errorf("destChain not set")
	}

	if len(c.FChain) == 0 {
		return fmt.Errorf("fChain not set")
	}

	if len(c.PricedTokens) == 0 {
		return fmt.Errorf("priced tokens not set")
	}

	if c.NewMsgScanBatchSize == 0 {
		return fmt.Errorf("newMsgScanBatchSize not set")
	}

	if _, ok := c.FChain[c.DestChain]; !ok {
		return fmt.Errorf("fChain not set for dest chain")
	}

	return nil
}

type ExecutePluginConfig struct {
	// DestChain is the ccip destination chain configured for the execute DON.
	DestChain ChainSelector

	// ObserverInfo is a map of oracle IDs to ObserverInfo.
	ObserverInfo map[commontypes.OracleID]ObserverInfo

	// MessageVisibilityInterval is the time interval for which the messages are visible by the plugin.
	MessageVisibilityInterval time.Duration
}

type ObserverInfo struct {
	Reads []ChainSelector
}

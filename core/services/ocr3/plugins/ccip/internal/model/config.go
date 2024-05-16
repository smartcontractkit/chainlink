package model

import (
	"time"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

type CommitPluginConfig struct {
	// Writer indicates that the node can contribute by sending reports to the destination chain.
	// Being a Writer guarantees that the node can also read from the destination chain.
	Writer bool

	// Reads define the chains that the current node can read from.
	Reads []ChainSelector

	// DestChain is the ccip destination chain configured for the commit plugin DON.
	DestChain ChainSelector

	// FChain defines the FChain value for each chain. FChain is used while forming consensus based on the observations.
	FChain map[ChainSelector]int

	// ObserverInfo is a map of oracle IDs to ObserverInfo.
	ObserverInfo map[commontypes.OracleID]ObserverInfo

	// FeeTokens is a list of tokens that can be used to pay for ccip fees.
	FeeTokens []types.Account

	// TokenPricesObserver indicates that the node can observe token prices.
	TokenPricesObserver bool

	NewMsgScanDuration  time.Duration
	NewMsgScanLimit     int
	NewMsgScanBatchSize int
}

type ObserverInfo struct {
	Reads []ChainSelector
}

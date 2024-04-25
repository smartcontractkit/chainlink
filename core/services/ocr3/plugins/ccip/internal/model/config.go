package model

import "time"

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

	// TBD:
	NewMsgScanDuration time.Duration
	NewMsgScanLimit    int
}

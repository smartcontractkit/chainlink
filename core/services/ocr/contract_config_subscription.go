package ocr

import (
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting/types"
)

var _ ocrtypes.ContractConfigSubscription = &ContractConfigSubscription{}

// ContractConfigSubscription only exists to comply with the
// ContractConfigSubscription interface, it's just a simple shell around
// ContractTracker that defines two methods
type ContractConfigSubscription ContractTracker

// Configs complies with ContractConfigSubscription interface
func (sub *ContractConfigSubscription) Configs() <-chan ocrtypes.ContractConfig {
	return sub.chConfigs
}

// Close is a no-op since Subscribing/Unsubscribing is handled in the
// Start/Close methods of the ContractTracker
func (sub *ContractConfigSubscription) Close() {}

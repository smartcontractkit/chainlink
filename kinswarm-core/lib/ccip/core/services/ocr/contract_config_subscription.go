package ocr

import (
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting/types"
)

var _ ocrtypes.ContractConfigSubscription = &OCRContractConfigSubscription{}

// OCRContractConfigSubscription only exists to comply with the
// ContractConfigSubscription interface, it's just a simple shell around
// OCRContractTracker that defines two methods
type OCRContractConfigSubscription OCRContractTracker

// Configs complies with ContractConfigSubscription interface
func (sub *OCRContractConfigSubscription) Configs() <-chan ocrtypes.ContractConfig {
	return sub.chConfigs
}

// Close is a no-op since Subscribing/Unsubscribing is handled in the
// Start/Close methods of the OCRContractTracker
func (sub *OCRContractConfigSubscription) Close() {}

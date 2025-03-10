package oraclecreator

import (
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

// ContractTransmitterFactory creates commit/execute transmitters for a specific chain family.
type ContractTransmitterFactory interface {
	// NewCommitTransmitter Creates a transmitter for "commit" operations
	NewCommitTransmitter(
		cw types.ContractWriter,
		fromAccount ocrtypes.Account,
		offrampAddress string,
		defaultMethod, priceOnlyMethod string,
	) ocr3types.ContractTransmitter[[]byte]

	// NewExecTransmitter Creates a transmitter for "execute" operations
	NewExecTransmitter(
		cw types.ContractWriter,
		fromAccount ocrtypes.Account,
		offrampAddress string,
	) ocr3types.ContractTransmitter[[]byte]
}

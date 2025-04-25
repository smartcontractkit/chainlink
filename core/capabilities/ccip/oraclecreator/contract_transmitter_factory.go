package oraclecreator

import (
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// ContractTransmitterFactory creates commit/execute transmitters for a specific chain family.
type ContractTransmitterFactory interface {
	// NewCommitTransmitter Creates a transmitter for "commit" operations
	NewCommitTransmitter(
		lggr logger.Logger,
		cw commontypes.ContractWriter,
		fromAccount ocrtypes.Account,
		offrampAddress string,
		defaultMethod, priceOnlyMethod string,
	) ocr3types.ContractTransmitter[[]byte]

	// NewExecTransmitter Creates a transmitter for "execute" operations
	NewExecTransmitter(
		lggr logger.Logger,
		cw commontypes.ContractWriter,
		fromAccount ocrtypes.Account,
		offrampAddress string,
	) ocr3types.ContractTransmitter[[]byte]
}

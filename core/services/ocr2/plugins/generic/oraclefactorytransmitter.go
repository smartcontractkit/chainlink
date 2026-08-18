package generic

import (
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/ocrcommon"
)

// NewContractTransmitter wraps impl so that FromAccount answers with
// transmitterID, the account this oracle is registered under.
//
// The wrapper itself lives in chainlink-common: a capability hosted outside the
// node needs the same thing for the same reason - libocr checks the account
// against the oracle's entry in the configuration - and it cannot import core to
// get it. This remains as the name core reaches it by.
func NewContractTransmitter(transmitterID string, impl ocr3types.ContractTransmitter[[]byte]) ocr3types.ContractTransmitter[[]byte] {
	return ocrcommon.TransmitterWithAccount(types.Account(transmitterID), impl)
}

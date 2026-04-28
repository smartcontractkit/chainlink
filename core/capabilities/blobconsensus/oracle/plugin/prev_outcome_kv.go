package plugin

import (
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
)

// prevOutcomeStateKey is the replicated KeyValueState key used to persist the
// last committed serialized oracletypes.Outcome between rounds (replaces
// ocr3types.OutcomeContext.PreviousOutcome for OCR 3.1).
var prevOutcomeStateKey = []byte("blobconsensus.v1.prevOutcome")

func readPreviousOutcomeBytesFromKV(r ocr3_1types.KeyValueStateReader) ([]byte, error) {
	b, err := r.Read(prevOutcomeStateKey)
	if err != nil {
		return nil, err
	}
	if len(b) == 0 {
		return nil, nil
	}
	return b, nil
}

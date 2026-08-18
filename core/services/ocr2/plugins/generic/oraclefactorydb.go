package generic

import (
	"strconv"

	"github.com/smartcontractkit/capabilities/libs/ocr"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// OracleFactoryDB is the in-memory libocr database of an OCR-based capability.
//
// It moved to github.com/smartcontractkit/capabilities/libs/ocr, so that a
// capability running outside this node has the same one rather than a second
// copy of it: nothing in it was ever a node's, and the reason it can be held in
// memory - a capability delivers its reports to its caller in-process - is a
// property of capabilities rather than of where they are hosted.
//
// Kept here as an alias so that nothing calling it has to change.
func OracleFactoryDB(specID int32, lggr logger.Logger) *ocr.Database {
	return ocr.NewDatabase(strconv.Itoa(int(specID)), lggr)
}

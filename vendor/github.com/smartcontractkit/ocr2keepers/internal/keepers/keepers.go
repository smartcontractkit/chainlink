package keepers

import (
	"log"

	"github.com/smartcontractkit/libocr/commontypes"

	"github.com/smartcontractkit/ocr2keepers/pkg/types"
)

type keepers struct {
	id                 commontypes.OracleID
	service            upkeepService
	encoder            types.ReportEncoder
	logger             *log.Logger
	filter             filterer
	reportGasLimit     uint32
	upkeepGasOverhead  uint32
	maxUpkeepBatchSize int
	reportBlockLag     int
}

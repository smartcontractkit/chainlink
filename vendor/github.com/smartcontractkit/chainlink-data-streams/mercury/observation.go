package mercury

import (
	"math/big"

	"github.com/smartcontractkit/libocr/commontypes"
)

type PAO interface {
	// These fields are common to all observations
	GetTimestamp() uint32
	GetObserver() commontypes.OracleID
	GetBenchmarkPrice() (*big.Int, bool)
}

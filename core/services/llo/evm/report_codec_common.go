package evm

import (
	"fmt"
	"math"

	"github.com/smartcontractkit/chainlink-data-streams/llo"
)

// Extracts nanosecond timestamps as uint32 number of seconds
func ExtractTimestamps(report llo.Report) (validAfterSeconds, observationTimestampSeconds uint32, err error) {
	vas := report.ValidAfterNanoseconds / 1e9
	ots := report.ObservationTimestampNanoseconds / 1e9
	if vas > math.MaxUint32 {
		err = fmt.Errorf("validAfterSeconds too large: %d", vas)
		return
	}
	if ots > math.MaxUint32 {
		err = fmt.Errorf("observationTimestampSeconds too large: %d", ots)
		return
	}
	return uint32(vas), uint32(ots), nil
}

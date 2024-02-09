package pb

import (
	"fmt"
	"math"

	libocr "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

func ReportTimestampToPb(ts libocr.ReportTimestamp) *ReportTimestamp {
	return &ReportTimestamp{
		ConfigDigest: ts.ConfigDigest[:],
		Epoch:        ts.Epoch,
		Round:        uint32(ts.Round),
	}
}

func ReportTimestampFromPb(ts *ReportTimestamp) (r libocr.ReportTimestamp, err error) {
	if l := len(ts.ConfigDigest); l != 32 {
		err = ErrConfigDigestLen(l)
		return
	}
	copy(r.ConfigDigest[:], ts.ConfigDigest)
	r.Epoch = ts.Epoch
	if ts.Round > math.MaxUint8 {
		err = ErrUint8Bounds{Name: "Round", U: ts.Round}
		return
	}
	r.Round = uint8(ts.Round)
	return
}

type ErrConfigDigestLen int

func (e ErrConfigDigestLen) Error() string {
	return fmt.Sprintf("invalid ConfigDigest len %d: must be 32", e)
}

type ErrUint8Bounds struct {
	U    uint32
	Name string
}

func (e ErrUint8Bounds) Error() string {
	return fmt.Sprintf("expected uint8 %s (max %d) but got %d", e.Name, math.MaxUint8, e.U)
}

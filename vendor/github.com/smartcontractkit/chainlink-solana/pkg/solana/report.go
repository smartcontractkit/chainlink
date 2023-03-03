package solana

import (
	"encoding/binary"
	"fmt"
	"math/big"
	"sort"

	"github.com/smartcontractkit/libocr/bigbigendian"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

var _ median.ReportCodec = (*ReportCodec)(nil)

type ReportCodec struct{}

func (c ReportCodec) BuildReport(oo []median.ParsedAttributedObservation) (types.Report, error) {
	n := len(oo)
	if n == 0 {
		return nil, fmt.Errorf("cannot build report from empty attributed observations")
	}

	// copy so we can safely re-order subsequently
	oo = append([]median.ParsedAttributedObservation{}, oo...)

	// get median timestamp
	sort.Slice(oo, func(i, j int) bool {
		return oo[i].Timestamp < oo[j].Timestamp
	})
	timestamp := oo[n/2].Timestamp

	// get median juelsPerFeeCoin
	sort.Slice(oo, func(i, j int) bool {
		return oo[i].JuelsPerFeeCoin.Cmp(oo[j].JuelsPerFeeCoin) < 0
	})
	juelsPerFeeCoin := oo[n/2].JuelsPerFeeCoin

	// get median by value
	// solana program size tx execution limit prevents reporting all observations
	// reporting only median value
	sort.Slice(oo, func(i, j int) bool {
		return oo[i].Value.Cmp(oo[j].Value) < 0
	})
	median := oo[n/2].Value

	observers := [32]byte{}
	for i, o := range oo {
		observers[i] = byte(o.Observer)
	}

	// report encoding
	report := []byte{}

	time := make([]byte, 4)
	binary.BigEndian.PutUint32(time, timestamp)
	report = append(report, time[:]...)

	observersCount := uint8(n)
	report = append(report, observersCount)

	report = append(report, observers[:]...)

	// TODO: replace with generalized function from libocr
	medianBytes, err := bigbigendian.SerializeSigned(int(MedianLen), median)
	if err != nil {
		return nil, fmt.Errorf("error in SerializeSigned(median): %w", err)
	}
	report = append(report, medianBytes[:]...)

	// TODO: replace with generalized function from libocr
	juelsPerFeeCoinBytes, err := bigbigendian.SerializeSigned(int(JuelsLen), juelsPerFeeCoin)
	if err != nil {
		return nil, fmt.Errorf("error in SerializeSigned(juelsPerFeeCoin): %w", err)
	}
	report = append(report, juelsPerFeeCoinBytes[:]...)

	return types.Report(report), nil
}

func (c ReportCodec) MedianFromReport(report types.Report) (*big.Int, error) {
	// report should contain timestamp + observers + median + juels per eth
	if len(report) != int(ReportLen) {
		return nil, fmt.Errorf("report length missmatch: %d (received), %d (expected)", len(report), ReportLen)
	}

	// unpack median observation
	start := int(ReportHeaderLen)
	end := start + int(MedianLen)
	median := report[start:end]
	return bigbigendian.DeserializeSigned(int(MedianLen), median)
}

func (c ReportCodec) MaxReportLength(n int) int {
	return int(ReportLen)
}

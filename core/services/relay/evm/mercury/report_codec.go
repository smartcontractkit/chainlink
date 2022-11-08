package mercury

import (
	"fmt"
	"math/big"
	"sort"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

// NOTE:
// This report codec is a copy/paste of the original median evmreportcodec
// here:
// https://github.com/smartcontractkit/offchain-reporting/blob/master/lib/offchainreporting2/reportingplugin/median/evmreportcodec/reportcodec.go
//
// As a hack, blockNumber is substituted for juelsPerFeeCoin.
//
// Feed ID is a constant that comes from the jobspec and is added to the report also.
//
// A prettier implementation is dependent on the generic libocr median plugin,
// here: https://github.com/smartcontractkit/offchain-reporting/pull/386/files

var reportTypes = getReportTypes()

func getReportTypes() abi.Arguments {
	mustNewType := func(t string) abi.Type {
		result, err := abi.NewType(t, "", []abi.ArgumentMarshaling{})
		if err != nil {
			panic(fmt.Sprintf("Unexpected error during abi.NewType: %s", err))
		}
		return result
	}
	return abi.Arguments([]abi.Argument{
		{Name: "feedId", Type: mustNewType("bytes32")},
		{Name: "observationsTimestamp", Type: mustNewType("uint32")},
		{Name: "observationsBlocknumber", Type: mustNewType("uint64")},
		{Name: "median", Type: mustNewType("int192")},
	})
}

var _ median.ReportCodec = ReportCodec{}

type ReportCodec struct {
	FeedID [32]byte
}

func (r ReportCodec) BuildReport(paos []median.ParsedAttributedObservation) (types.Report, error) {
	if len(paos) == 0 {
		return nil, fmt.Errorf("cannot build report from empty attributed observations")
	}

	// copy so we can safely re-order subsequently
	paos = append([]median.ParsedAttributedObservation{}, paos...)

	// get median timestamp
	sort.Slice(paos, func(i, j int) bool {
		return paos[i].Timestamp < paos[j].Timestamp
	})
	timestamp := paos[len(paos)/2].Timestamp

	// get median juelsPerFeeCoin
	// HACK: This is actually the block number!
	sort.Slice(paos, func(i, j int) bool {
		return paos[i].JuelsPerFeeCoin.Cmp(paos[j].JuelsPerFeeCoin) < 0
	})
	blockNumber := paos[len(paos)/2].JuelsPerFeeCoin

	// sort by values
	sort.Slice(paos, func(i, j int) bool {
		return paos[i].Value.Cmp(paos[j].Value) < 0
	})

	observers := [32]byte{}
	observations := make([]*big.Int, len(paos))

	for i, pao := range paos {
		observers[i] = byte(pao.Observer)
		observations[i] = pao.Value
	}

	median := observations[len(observations)/2]
	if median == nil {
		return nil, fmt.Errorf("median is nil")
	}

	reportBytes, err := reportTypes.Pack(r.FeedID, timestamp, uint64(blockNumber.Int64()), median)
	return types.Report(reportBytes), errors.Wrap(err, "failed to pack report blob")
}

func (ReportCodec) MedianFromReport(report types.Report) (*big.Int, error) {
	reportElems := map[string]interface{}{}
	if err := reportTypes.UnpackIntoMap(reportElems, report); err != nil {
		return nil, fmt.Errorf("error during unpack: %w", err)
	}

	medianIface, ok := reportElems["median"]
	if !ok {
		return nil, fmt.Errorf("unpacked report has no 'median'")
	}

	median, ok := medianIface.(*big.Int)
	if !ok {
		return nil, fmt.Errorf("cannot cast median to *big.Int, type is %T", medianIface)
	}

	return median, nil
}

func (ReportCodec) MaxReportLength(_ int) int {
	return 128
}

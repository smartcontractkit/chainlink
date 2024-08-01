package medianreport

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"sort"

	"github.com/NethermindEth/juno/core/felt"
	starknetutils "github.com/NethermindEth/starknet.go/utils"

	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink-starknet/relayer/pkg/starknet"
)

var _ median.ReportCodec = (*ReportCodec)(nil)

var (
	timestampSizeBytes       = starknet.FeltLength
	observersSizeBytes       = starknet.FeltLength
	observationsLenBytes     = starknet.FeltLength
	prefixSizeBytes          = timestampSizeBytes + observersSizeBytes + observationsLenBytes
	juelsPerFeeCoinSizeBytes = starknet.FeltLength
	gasPriceSizeBytes        = starknet.FeltLength
	observationSizeBytes     = starknet.FeltLength
)

type ReportCodec struct{}

func (c ReportCodec) BuildReport(oo []median.ParsedAttributedObservation) (types.Report, error) {
	num := len(oo)
	if num == 0 {
		return nil, errors.New("couldn't build report from empty attributed observations")
	}

	for _, o := range oo {
		if o.Value.Sign() == -1 || o.JuelsPerFeeCoin.Sign() == -1 || o.GasPriceSubunits.Sign() == -1 {
			return nil, fmt.Errorf("starknet does not support negative values: value = (%v), fee = (%v), gas = (%v)", o.Value, o.JuelsPerFeeCoin, o.GasPriceSubunits)
		}
	}

	// preserve original array
	oo = append([]median.ParsedAttributedObservation{}, oo...)
	numFelt := starknetutils.BigIntToFelt(big.NewInt(int64(num)))

	// median timestamp
	sort.Slice(oo, func(i, j int) bool {
		return oo[i].Timestamp < oo[j].Timestamp
	})
	timestamp := oo[num/2].Timestamp
	timestampFelt := starknetutils.BigIntToFelt(big.NewInt(int64(timestamp)))

	// median juelsPerFeeCoin
	sort.Slice(oo, func(i, j int) bool {
		return oo[i].JuelsPerFeeCoin.Cmp(oo[j].JuelsPerFeeCoin) < 0
	})
	juelsPerFeeCoin := oo[num/2].JuelsPerFeeCoin
	juelsPerFeeCoinFelt := starknetutils.BigIntToFelt(juelsPerFeeCoin)

	sort.Slice(oo, func(i, j int) bool {
		return oo[i].GasPriceSubunits.Cmp(oo[j].GasPriceSubunits) < 0
	})
	// gas price in FRI
	gasPriceSubunits := oo[num/2].GasPriceSubunits
	gasPriceSubunitsFelt := starknetutils.BigIntToFelt(gasPriceSubunits)

	// sort by values
	sort.Slice(oo, func(i, j int) bool {
		return oo[i].Value.Cmp(oo[j].Value) < 0
	})

	var observers = make([]byte, starknet.FeltLength)
	var observations []*felt.Felt
	for i, o := range oo {
		observers[i] = byte(o.Observer)

		f := starknetutils.BigIntToFelt(o.Value)
		observations = append(observations, f)
	}

	var report []byte

	buf := timestampFelt.Bytes()
	report = append(report, buf[:]...)
	report = append(report, observers...)
	buf = numFelt.Bytes()
	report = append(report, buf[:]...)
	for _, o := range observations {
		buf = o.Bytes()
		report = append(report, buf[:]...)
	}
	buf = juelsPerFeeCoinFelt.Bytes()
	report = append(report, buf[:]...)
	buf = gasPriceSubunitsFelt.Bytes()
	report = append(report, buf[:]...)

	return report, nil
}

func (c ReportCodec) MedianFromReport(report types.Report) (*big.Int, error) {
	rLen := len(report)
	if rLen < prefixSizeBytes+juelsPerFeeCoinSizeBytes+gasPriceSizeBytes {
		return nil, errors.New("invalid report length")
	}

	// Decode the number of observations
	numBig := new(felt.Felt).SetBytes(report[(timestampSizeBytes + observersSizeBytes):prefixSizeBytes]).BigInt(big.NewInt(0))
	if !numBig.IsUint64() {
		return nil, errors.New("length of observations is invalid")
	}
	n64 := numBig.Uint64()
	if n64 == 0 {
		return nil, errors.New("unpacked report has no observations")
	}
	if n64 >= math.MaxInt8 {
		return nil, errors.New("length of observations is invalid")
	}

	// Check if the report is big enough
	n := int(n64)
	expectedLen := prefixSizeBytes + (observationSizeBytes * n) + juelsPerFeeCoinSizeBytes + gasPriceSizeBytes
	if rLen < expectedLen {
		return nil, errors.New("invalid report length, missing main, juelsPerFeeCoin or gasPrice observations")
	}

	// Decode observations
	var oo []*big.Int
	for i := 0; i < n; i++ {
		start := prefixSizeBytes + observationSizeBytes*i
		end := start + observationSizeBytes
		obv := new(felt.Felt).SetBytes(report[start:end])
		o := obv.BigInt(big.NewInt(0))
		oo = append(oo, o)
	}

	// Check if the report contains sorted observations
	_less := func(i, j int) bool {
		return oo[i].Cmp(oo[j]) < 0
	}
	sorted := sort.SliceIsSorted(oo, _less)
	if !sorted {
		return nil, errors.New("observations not sorted")
	}

	return oo[n/2], nil
}

func (c ReportCodec) MaxReportLength(n int) (int, error) {
	return prefixSizeBytes + (n * observationSizeBytes) + juelsPerFeeCoinSizeBytes + gasPriceSizeBytes, nil
}

func SplitReport(report types.Report) ([][]byte, error) {
	chunkSize := starknet.FeltLength
	if len(report)%chunkSize != 0 {
		return [][]byte{}, errors.New("invalid report length")
	}

	// order is guaranteed by buildReport:
	//   observation_timestamp
	//   observers
	//   observations_len
	//   observations
	//   juels_per_fee_coin
	//   gas_price
	slices := [][]byte{}
	for i := 0; i < len(report)/chunkSize; i++ {
		idx := i * chunkSize
		slices = append(slices, report[idx:(idx+chunkSize)])
	}

	return slices, nil
}

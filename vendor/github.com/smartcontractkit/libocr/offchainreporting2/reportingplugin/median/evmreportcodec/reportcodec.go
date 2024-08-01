package evmreportcodec

import (
	"fmt"
	"math/big"
	"sort"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

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
		{Name: "observationsTimestamp", Type: mustNewType("uint32")},
		{Name: "rawObservers", Type: mustNewType("bytes32")},
		{Name: "observations", Type: mustNewType("int192[]")},
		{Name: "juelsPerFeeCoin", Type: mustNewType("int192")},
		// In the EVM, contracts can query tx.gasPrice during execution. Therefore, there is no need to include it in the report.
	})
}

var _ median.ReportCodec = ReportCodec{}

type ReportCodec struct{}

func (ReportCodec) BuildReport(paos []median.ParsedAttributedObservation) (types.Report, error) {
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
	sort.Slice(paos, func(i, j int) bool {
		return paos[i].JuelsPerFeeCoin.Cmp(paos[j].JuelsPerFeeCoin) < 0
	})
	juelsPerFeeCoin := paos[len(paos)/2].JuelsPerFeeCoin

	// sort by values
	sort.Slice(paos, func(i, j int) bool {
		return paos[i].Value.Cmp(paos[j].Value) < 0
	})

	observers := [32]byte{}
	observations := []*big.Int{}

	for i, pao := range paos {
		observers[i] = byte(pao.Observer)
		observations = append(observations, pao.Value)
	}

	reportBytes, err := reportTypes.Pack(timestamp, observers, observations, juelsPerFeeCoin)
	return types.Report(reportBytes), err
}

func (ReportCodec) MedianFromReport(report types.Report) (*big.Int, error) {
	reportElems := map[string]interface{}{}
	if err := reportTypes.UnpackIntoMap(reportElems, report); err != nil {
		return nil, fmt.Errorf("error during unpack: %w", err)
	}

	observationsIface, ok := reportElems["observations"]
	if !ok {
		return nil, fmt.Errorf("unpacked report has no 'observations'")
	}

	observations, ok := observationsIface.([]*big.Int)
	if !ok {
		return nil, fmt.Errorf("cannot cast observations to []*big.Int, type is %T", observationsIface)
	}

	if len(observations) == 0 {
		return nil, fmt.Errorf("observations are empty")
	}

	median := observations[len(observations)/2]
	if median == nil {
		return nil, fmt.Errorf("median is nil")
	}

	return median, nil
}

func (ReportCodec) MaxReportLength(n int) (int, error) {
	return 32 /* timestamp */ + 32 /* rawObservers */ + (2*32 + n*32) /*observations*/ + 32 /* juelsPerFeeCoin */, nil
}

func (ReportCodec) XXXJuelsPerFeeCoinFromReport(report types.Report) (*big.Int, error) {
	reportElems := map[string]interface{}{}
	if err := reportTypes.UnpackIntoMap(reportElems, report); err != nil {
		return nil, fmt.Errorf("error during unpack: %w", err)
	}

	juelsPerFeeCoinInterface, ok := reportElems["juelsPerFeeCoin"]
	if !ok {
		return nil, fmt.Errorf("unpacked report has no 'juelsPerFeeCoin'")
	}

	juelsPerFeeCoin, ok := juelsPerFeeCoinInterface.(*big.Int)
	if !ok {
		return nil, fmt.Errorf("cannot cast juelsPerFeeCoin to *big.Int, type is %T", juelsPerFeeCoinInterface)
	}

	return juelsPerFeeCoin, nil
}

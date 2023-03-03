package chain

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"

	ktypes "github.com/smartcontractkit/ocr2keepers/pkg/types"
)

const (
	ActiveUpkeepIDBatchSize int64  = 10000
	separator               string = "|"
)

var (
	ErrRegistryCallFailure   = fmt.Errorf("registry chain call failure")
	ErrBlockKeyNotParsable   = fmt.Errorf("block identifier not parsable")
	ErrUpkeepKeyNotParsable  = fmt.Errorf("upkeep key not parsable")
	ErrInitializationFailure = fmt.Errorf("failed to initialize registry")
	ErrContextCancelled      = fmt.Errorf("context was cancelled")
)

type evmReportEncoder struct{}

func NewEVMReportEncoder() *evmReportEncoder {
	return &evmReportEncoder{}
}

var (
	Uint256, _                = abi.NewType("uint256", "", nil)
	Uint256Arr, _             = abi.NewType("uint256[]", "", nil)
	PerformDataMarshalingArgs = []abi.ArgumentMarshaling{
		{Name: "checkBlockNumber", Type: "uint32"},
		{Name: "checkBlockhash", Type: "bytes32"},
		{Name: "performData", Type: "bytes"},
	}
	PerformDataArr, _ = abi.NewType("tuple(uint32,bytes32,bytes)[]", "", PerformDataMarshalingArgs)
)

func (b *evmReportEncoder) EncodeReport(toReport []ktypes.UpkeepResult) ([]byte, error) {
	if len(toReport) == 0 {
		return nil, nil
	}

	reportArgs := abi.Arguments{
		{Name: "fastGasWei", Type: Uint256},
		{Name: "linkNative", Type: Uint256},
		{Name: "upkeepIds", Type: Uint256Arr},
		{Name: "wrappedPerformDatas", Type: PerformDataArr},
	}

	var baseValuesIdx int
	for i, rpt := range toReport {
		if rpt.CheckBlockNumber > uint32(baseValuesIdx) {
			baseValuesIdx = i
		}
	}

	fastGas := toReport[baseValuesIdx].FastGasWei
	link := toReport[baseValuesIdx].LinkNative
	ids := make([]*big.Int, len(toReport))
	data := make([]wrappedPerform, len(toReport))

	for i, result := range toReport {
		_, upkeepId, err := result.Key.BlockKeyAndUpkeepID()
		if err != nil {
			return nil, fmt.Errorf("%w: report encoding error", err)
		}

		upkeepIdInt, ok := upkeepId.BigInt()
		if !ok {
			return nil, ErrUpkeepKeyNotParsable
		}

		ids[i] = upkeepIdInt
		data[i] = wrappedPerform{
			CheckBlockNumber: result.CheckBlockNumber,
			CheckBlockhash:   result.CheckBlockHash,
			PerformData:      result.PerformData,
		}
	}

	bts, err := reportArgs.Pack(fastGas, link, ids, data)
	if err != nil {
		return []byte{}, fmt.Errorf("%w: failed to pack report data", err)
	}

	return bts, nil
}

func (b *evmReportEncoder) DecodeReport(report []byte) ([]ktypes.UpkeepResult, error) {
	mKeys := []string{"fastGasWei", "linkNative", "upkeepIds", "wrappedPerformDatas"}

	reportArgs := abi.Arguments{
		{Name: mKeys[0], Type: Uint256},
		{Name: mKeys[1], Type: Uint256},
		{Name: mKeys[2], Type: Uint256Arr},
		{Name: mKeys[3], Type: PerformDataArr},
	}

	m := make(map[string]interface{})
	if err := reportArgs.UnpackIntoMap(m, report); err != nil {
		return nil, err
	}

	for _, key := range mKeys {
		if _, ok := m[key]; !ok {
			return nil, fmt.Errorf("decoding error: %s missing from struct", key)
		}
	}

	res := []ktypes.UpkeepResult{}

	var ok bool
	var upkeepIds []*big.Int
	var wei *big.Int
	var link *big.Int

	if upkeepIds, ok = m[mKeys[2]].([]*big.Int); !ok {
		return res, fmt.Errorf("upkeep ids of incorrect type in report")
	}

	// TODO: a type assertion on `wrappedPerform` did not work, even with the
	// exact same struct definition as what follows. reflect was used to get the
	// struct definition. not sure yet how to clean this up.
	// ex:
	// t := reflect.TypeOf(rawPerforms)
	// fmt.Printf("%v\n", t)
	performs, ok := m[mKeys[3]].([]struct {
		CheckBlockNumber uint32   `json:"checkBlockNumber"`
		CheckBlockhash   [32]byte `json:"checkBlockhash"`
		PerformData      []byte   `json:"performData"`
	})
	if !ok {
		return res, fmt.Errorf("performs of incorrect structure in report")
	}

	if len(upkeepIds) != len(performs) {
		return res, fmt.Errorf("upkeep ids and performs should have matching length")
	}

	if wei, ok = m[mKeys[0]].(*big.Int); !ok {
		return res, fmt.Errorf("fast gas as wrong type")
	}

	if link, ok = m[mKeys[1]].(*big.Int); !ok {
		return res, fmt.Errorf("link native as wrong type")
	}

	res = make([]ktypes.UpkeepResult, len(upkeepIds))
	for i := 0; i < len(upkeepIds); i++ {
		res[i] = ktypes.UpkeepResult{
			Key:              NewUpkeepKey(big.NewInt(int64(performs[i].CheckBlockNumber)), upkeepIds[i]),
			State:            ktypes.Eligible,
			PerformData:      performs[i].PerformData,
			FastGasWei:       wei,
			LinkNative:       link,
			CheckBlockNumber: performs[i].CheckBlockNumber,
			CheckBlockHash:   performs[i].CheckBlockhash,
		}
	}

	return res, nil
}

type wrappedPerform struct {
	CheckBlockNumber uint32   `abi:"checkBlockNumber"`
	CheckBlockhash   [32]byte `abi:"checkBlockhash"`
	PerformData      []byte   `abi:"performData"`
}

package evm

import (
	"errors"
	"fmt"
	"math/big"
	"reflect"

	"github.com/ethereum/go-ethereum/accounts/abi"

	ocr2keepers "github.com/smartcontractkit/chainlink-automation/pkg/v2"
	"github.com/smartcontractkit/chainlink-automation/pkg/v2/encoding"
)

type EVMAutomationEncoder20 struct {
	encoding.BasicEncoder
}

var (
	Uint256, _                = abi.NewType("uint256", "", nil)
	Uint256Arr, _             = abi.NewType("uint256[]", "", nil)
	PerformDataMarshalingArgs = []abi.ArgumentMarshaling{
		{Name: "checkBlockNumber", Type: "uint32"},
		{Name: "checkBlockhash", Type: "bytes32"},
		{Name: "performData", Type: "bytes"},
	}
	PerformDataArr, _   = abi.NewType("tuple(uint32,bytes32,bytes)[]", "", PerformDataMarshalingArgs)
	ErrUnexpectedResult = errors.New("unexpected result struct")
	packFn              = reportArgs.Pack
	unpackIntoMapFn     = reportArgs.UnpackIntoMap
	mKeys               = []string{"fastGasWei", "linkNative", "upkeepIds", "wrappedPerformDatas"}
	reportArgs          = abi.Arguments{
		{Name: mKeys[0], Type: Uint256},
		{Name: mKeys[1], Type: Uint256},
		{Name: mKeys[2], Type: Uint256Arr},
		{Name: mKeys[3], Type: PerformDataArr},
	}
)

type EVMAutomationUpkeepResult20 struct {
	// Block is the block number used to build an UpkeepKey for this result
	Block uint32
	// ID is the unique identifier for the upkeep
	ID            *big.Int
	Eligible      bool
	FailureReason uint8
	GasUsed       *big.Int
	PerformData   []byte
	FastGasWei    *big.Int
	LinkNative    *big.Int
	// CheckBlockNumber is the block number that the contract indicates the
	// upkeep was checked on
	CheckBlockNumber uint32
	CheckBlockHash   [32]byte
	ExecuteGas       uint32
}

func (enc EVMAutomationEncoder20) EncodeReport(toReport []ocr2keepers.UpkeepResult) ([]byte, error) {
	if len(toReport) == 0 {
		return nil, nil
	}

	var (
		fastGas *big.Int
		link    *big.Int
	)

	ids := make([]*big.Int, len(toReport))
	data := make([]wrappedPerform, len(toReport))

	for i, result := range toReport {
		res, ok := result.(EVMAutomationUpkeepResult20)
		if !ok {
			return nil, errors.New("unexpected upkeep result struct")
		}

		// only take these values from the first result
		// TODO: find a new way to get these values
		if i == 0 {
			fastGas = res.FastGasWei
			link = res.LinkNative
		}

		ids[i] = res.ID
		data[i] = wrappedPerform{
			CheckBlockNumber: res.CheckBlockNumber,
			CheckBlockhash:   res.CheckBlockHash,
			PerformData:      res.PerformData,
		}
	}

	bts, err := packFn(fastGas, link, ids, data)
	if err != nil {
		return []byte{}, fmt.Errorf("%w: failed to pack report data", err)
	}

	return bts, nil
}

func (enc EVMAutomationEncoder20) DecodeReport(report []byte) ([]ocr2keepers.UpkeepResult, error) {
	m := make(map[string]interface{})
	if err := unpackIntoMapFn(m, report); err != nil {
		return nil, err
	}

	for _, key := range mKeys {
		if _, ok := m[key]; !ok {
			return nil, fmt.Errorf("decoding error: %s missing from struct", key)
		}
	}

	res := []ocr2keepers.UpkeepResult{}

	var (
		ok        bool
		upkeepIds []*big.Int
		wei       *big.Int
		link      *big.Int
	)

	if upkeepIds, ok = m[mKeys[2]].([]*big.Int); !ok {
		return res, errors.New("upkeep ids of incorrect type in report")
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
		return res, errors.New("performs of incorrect structure in report")
	}

	if len(upkeepIds) != len(performs) {
		return res, errors.New("upkeep ids and performs should have matching length")
	}

	if wei, ok = m[mKeys[0]].(*big.Int); !ok {
		return res, errors.New("fast gas as wrong type")
	}

	if link, ok = m[mKeys[1]].(*big.Int); !ok {
		return res, errors.New("link native as wrong type")
	}

	res = make([]ocr2keepers.UpkeepResult, len(upkeepIds))

	for i := 0; i < len(upkeepIds); i++ {
		r := EVMAutomationUpkeepResult20{
			Block:            performs[i].CheckBlockNumber,
			ID:               upkeepIds[i],
			Eligible:         true,
			PerformData:      performs[i].PerformData,
			FastGasWei:       wei,
			LinkNative:       link,
			CheckBlockNumber: performs[i].CheckBlockNumber,
			CheckBlockHash:   performs[i].CheckBlockhash,
		}

		res[i] = ocr2keepers.UpkeepResult(r)
	}

	return res, nil
}

func (enc EVMAutomationEncoder20) Eligible(result ocr2keepers.UpkeepResult) (bool, error) {
	res, ok := result.(EVMAutomationUpkeepResult20)
	if !ok {
		tp := reflect.TypeOf(result)
		return false, fmt.Errorf("%w: name: %s, kind: %s", ErrUnexpectedResult, tp.Name(), tp.Kind())
	}

	return res.Eligible, nil
}

func (enc EVMAutomationEncoder20) Detail(result ocr2keepers.UpkeepResult) (ocr2keepers.UpkeepKey, uint32, error) {
	res, ok := result.(EVMAutomationUpkeepResult20)
	if !ok {
		return nil, 0, ErrUnexpectedResult
	}

	str := fmt.Sprintf("%d%s%s", res.Block, separator, res.ID)

	return ocr2keepers.UpkeepKey([]byte(str)), res.ExecuteGas, nil
}

func (enc EVMAutomationEncoder20) KeysFromReport(b []byte) ([]ocr2keepers.UpkeepKey, error) {
	results, err := enc.DecodeReport(b)
	if err != nil {
		return nil, err
	}

	keys := make([]ocr2keepers.UpkeepKey, 0, len(results))
	for _, result := range results {
		res, ok := result.(EVMAutomationUpkeepResult20)
		if !ok {
			return nil, errors.New("unexpected result struct")
		}

		str := fmt.Sprintf("%d%s%s", res.Block, separator, res.ID)
		keys = append(keys, ocr2keepers.UpkeepKey([]byte(str)))
	}

	return keys, nil
}

type wrappedPerform struct {
	CheckBlockNumber uint32   `abi:"checkBlockNumber"`
	CheckBlockhash   [32]byte `abi:"checkBlockhash"`
	PerformData      []byte   `abi:"performData"`
}

type BlockKeyHelper[T uint32 | int64] struct {
}

func (kh BlockKeyHelper[T]) MakeBlockKey(b T) ocr2keepers.BlockKey {
	return ocr2keepers.BlockKey(fmt.Sprintf("%d", b))
}

type UpkeepKeyHelper[T uint32 | int64] struct {
}

func (kh UpkeepKeyHelper[T]) MakeUpkeepKey(b T, id *big.Int) ocr2keepers.UpkeepKey {
	return ocr2keepers.UpkeepKey(fmt.Sprintf("%d%s%s", b, separator, id))
}

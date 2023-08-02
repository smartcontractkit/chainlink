package evm

import (
	"fmt"
	"math/big"
	"reflect"

	"github.com/ethereum/go-ethereum/accounts/abi"
	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg"
	"github.com/smartcontractkit/ocr2keepers/pkg/encoding"
)

type EVMAutomationEncoder21 struct {
	encoding.BasicEncoder
}

func mustNewType(t string, internalType string, components []abi.ArgumentMarshaling) abi.Type {
	a, err := abi.NewType(t, internalType, components)
	if err != nil {
		panic(err)
	}
	return a
}

var (
	Uint256               = mustNewType("uint256", "", nil)
	Uint256Arr            = mustNewType("uint256[]", "", nil)
	BytesArr              = mustNewType("bytes[]", "", nil)
	TriggerMarshalingArgs = []abi.ArgumentMarshaling{
		{Name: "blockNumber", Type: "uint32"},
		{Name: "blockHash", Type: "bytes32"},
	}
	TriggerArr          = mustNewType("tuple(uint32,bytes32)[]", "", TriggerMarshalingArgs)
	ErrUnexpectedResult = fmt.Errorf("unexpected result struct")
	packFn              = reportArgs.Pack
	unpackIntoMapFn     = reportArgs.UnpackIntoMap
	mKeys               = []string{"fastGasWei", "linkNative", "upkeepIds", "gasLimits", "triggers", "performDatas"}
	reportArgs          = abi.Arguments{
		{Name: mKeys[0], Type: Uint256},
		{Name: mKeys[1], Type: Uint256},
		{Name: mKeys[2], Type: Uint256Arr},
		{Name: mKeys[3], Type: Uint256Arr},
		{Name: mKeys[4], Type: TriggerArr},
		{Name: mKeys[5], Type: BytesArr},
	}
)

type EVMAutomationUpkeepResult21 struct {
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
	Retryable        bool
}

func (enc EVMAutomationEncoder21) EncodeReport(toReport []ocr2keepers.UpkeepResult) ([]byte, error) {
	if len(toReport) == 0 {
		return nil, nil
	}

	var (
		fastGas *big.Int
		link    *big.Int
	)

	ids := make([]*big.Int, len(toReport))
	gasLimits := make([]*big.Int, len(toReport))
	triggers := make([]wrappedTrigger, len(toReport))
	performDatas := make([][]byte, len(toReport))

	for i, result := range toReport {
		res, ok := result.(EVMAutomationUpkeepResult21)
		if !ok {
			return nil, fmt.Errorf("unexpected upkeep result struct")
		}

		// only take these values from the first result
		// TODO: find a new way to get these values
		if i == 0 {
			fastGas = res.FastGasWei
			link = res.LinkNative
		}

		ids[i] = res.ID
		gasLimits[i] = res.GasUsed
		triggers[i] = wrappedTrigger{
			BlockNumber: res.CheckBlockNumber,
			BlockHash:   res.CheckBlockHash,
		}
		performDatas[i] = res.PerformData
	}

	bts, err := packFn(fastGas, link, ids, gasLimits, triggers, performDatas)
	if err != nil {
		return []byte{}, fmt.Errorf("%w: failed to pack report data", err)
	}

	return bts, nil
}

func (enc EVMAutomationEncoder21) DecodeReport(report []byte) ([]ocr2keepers.UpkeepResult, error) {
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
		performs  [][]byte
		// gasLimits []*big.Int // TODO
		wei  *big.Int
		link *big.Int
	)

	if upkeepIds, ok = m[mKeys[2]].([]*big.Int); !ok {
		return res, fmt.Errorf("upkeep ids of incorrect type in report")
	}

	// TODO: a type assertion on `wrappedTrigger` did not work, even with the
	// exact same struct definition as what follows. reflect was used to get the
	// struct definition. not sure yet how to clean this up.
	// ex:
	// t := reflect.TypeOf(rawPerforms)
	// fmt.Printf("%v\n", t)
	triggers, ok := m[mKeys[4]].([]struct {
		BlockNumber uint32   `abi:"blockNumber"`
		BlockHash   [32]byte `abi:"blockHash"`
	})
	if !ok {
		return res, fmt.Errorf("triggers of incorrect structure in report")
	}

	if len(upkeepIds) != len(triggers) {
		return res, fmt.Errorf("upkeep ids and triggers should have matching length")
	}

	if wei, ok = m[mKeys[0]].(*big.Int); !ok {
		return res, fmt.Errorf("fast gas as wrong type")
	}

	if link, ok = m[mKeys[1]].(*big.Int); !ok {
		return res, fmt.Errorf("link native as wrong type")
	}
	// if gasLimits, ok = m[mKeys[3]].([]*big.Int); !ok {
	// 	return res, fmt.Errorf("gas limits as wrong type")
	// }

	if performs, ok = m[mKeys[5]].([][]byte); !ok {
		return res, fmt.Errorf("perform datas as wrong type")
	}

	res = make([]ocr2keepers.UpkeepResult, len(upkeepIds))

	for i := 0; i < len(upkeepIds); i++ {
		r := EVMAutomationUpkeepResult21{
			Block:            triggers[i].BlockNumber,
			ID:               upkeepIds[i],
			Eligible:         true,
			PerformData:      performs[i],
			FastGasWei:       wei,
			LinkNative:       link,
			CheckBlockNumber: triggers[i].BlockNumber,
			CheckBlockHash:   triggers[i].BlockHash,
		}

		res[i] = ocr2keepers.UpkeepResult(r)
	}

	return res, nil
}

func (enc EVMAutomationEncoder21) Eligible(result ocr2keepers.UpkeepResult) (bool, error) {
	res, ok := result.(EVMAutomationUpkeepResult21)
	if !ok {
		tp := reflect.TypeOf(result)
		return false, fmt.Errorf("%s: name: %s, kind: %s", ErrUnexpectedResult, tp.Name(), tp.Kind())
	}

	return res.Eligible, nil
}

func (enc EVMAutomationEncoder21) Detail(result ocr2keepers.UpkeepResult) (ocr2keepers.UpkeepKey, uint32, error) {
	res, ok := result.(EVMAutomationUpkeepResult21)
	if !ok {
		return nil, 0, ErrUnexpectedResult
	}

	str := fmt.Sprintf("%d%s%s", res.Block, separator, res.ID)

	return ocr2keepers.UpkeepKey([]byte(str)), res.ExecuteGas, nil
}

func (enc EVMAutomationEncoder21) KeysFromReport(b []byte) ([]ocr2keepers.UpkeepKey, error) {
	results, err := enc.DecodeReport(b)
	if err != nil {
		return nil, err
	}

	keys := make([]ocr2keepers.UpkeepKey, 0, len(results))
	for _, result := range results {
		res, ok := result.(EVMAutomationUpkeepResult21)
		if !ok {
			return nil, fmt.Errorf("unexpected result struct")
		}

		str := fmt.Sprintf("%d%s%s", res.Block, separator, res.ID)
		keys = append(keys, ocr2keepers.UpkeepKey([]byte(str)))
	}

	return keys, nil
}

type wrappedTrigger struct {
	BlockNumber uint32   `abi:"blockNumber"`
	BlockHash   [32]byte `abi:"blockHash"`
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

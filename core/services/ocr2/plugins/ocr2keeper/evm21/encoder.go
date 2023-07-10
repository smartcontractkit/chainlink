package evm

import (
	"fmt"
	"math/big"
	"reflect"

	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg"
	"github.com/smartcontractkit/ocr2keepers/pkg/encoding"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/automation_utils_2_1"
)

type EVMAutomationEncoder21 struct {
	encoding.BasicEncoder
	packer *evmRegistryPackerV2_1
}

var (
	ErrUnexpectedResult = fmt.Errorf("unexpected result struct")
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

// TODO: align once we merge with ocr2keepers new types (ocr2keepers.CheckResult)
func (enc EVMAutomationEncoder21) EncodeReport(toReport []ocr2keepers.UpkeepResult) ([]byte, error) {
	if len(toReport) == 0 {
		return nil, nil
	}

	report := automation_utils_2_1.KeeperRegistryBase21Report{
		UpkeepIds:    make([]*big.Int, len(toReport)),
		GasLimits:    make([]*big.Int, len(toReport)),
		Triggers:     make([][]byte, len(toReport)),
		PerformDatas: make([][]byte, len(toReport)),
	}

	for i, result := range toReport {
		res, ok := result.(EVMAutomationUpkeepResult21)
		if !ok {
			return nil, fmt.Errorf("unexpected upkeep result struct")
		}

		// only take these values from the first result
		// TODO: find a new way to get these values
		if i == 0 {
			report.FastGasWei = res.FastGasWei
			report.LinkNative = res.LinkNative
		}

		report.UpkeepIds[i] = res.ID
		report.GasLimits[i] = res.GasUsed
		trigger, err := enc.packer.PackTrigger(res.ID, triggerWrapper{
			BlockNum:  res.CheckBlockNumber,
			BlockHash: res.CheckBlockHash,
			// TODO: fill with real info
			// LogIndex: 0,
			// TxHash:   [32]byte{},
		})
		if err != nil {
			return nil, fmt.Errorf("%w: failed to pack trigger", err)
		}
		report.Triggers[i] = trigger
		report.PerformDatas[i] = res.PerformData
	}

	return enc.packer.PackReport(report)
}

func (enc EVMAutomationEncoder21) DecodeReport(raw []byte) ([]ocr2keepers.UpkeepResult, error) {
	report, err := enc.packer.UnpackReport(raw)
	if err != nil {
		return nil, err
	}

	if err := enc.validateReport(report); err != nil {
		return nil, err
	}

	res := make([]ocr2keepers.UpkeepResult, len(report.UpkeepIds))

	for i := 0; i < len(report.UpkeepIds); i++ {
		trigger, err := enc.packer.UnpackTrigger(report.UpkeepIds[i], report.Triggers[i])
		if err != nil {
			// TODO: log error and continue instead?
			return nil, fmt.Errorf("%w: failed to unpack trigger", err)
		}
		r := EVMAutomationUpkeepResult21{
			Block:            trigger.BlockNum,
			ID:               report.UpkeepIds[i],
			Eligible:         true,
			PerformData:      report.PerformDatas[i],
			FastGasWei:       report.FastGasWei,
			LinkNative:       report.LinkNative,
			CheckBlockNumber: trigger.BlockNum,
			CheckBlockHash:   trigger.BlockHash,
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

// validateReport checks that the report is valid, currently checking that all
// lists are the same length.
// TODO: add more validations? e.g. parse validate triggers
func (enc EVMAutomationEncoder21) validateReport(report automation_utils_2_1.KeeperRegistryBase21Report) error {
	if len(report.UpkeepIds) != len(report.GasLimits) {
		return fmt.Errorf("invalid report: upkeepIds and gasLimits must be the same length")
	}
	if len(report.UpkeepIds) != len(report.Triggers) {
		return fmt.Errorf("invalid report: upkeepIds and triggers must be the same length")
	}
	if len(report.UpkeepIds) != len(report.PerformDatas) {
		return fmt.Errorf("invalid report: upkeepIds and performDatas must be the same length")
	}
	return nil
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

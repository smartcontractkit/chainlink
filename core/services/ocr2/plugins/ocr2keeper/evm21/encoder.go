package evm

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg"
	"github.com/smartcontractkit/ocr2keepers/pkg/encoding"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/automation_utils_2_1"
)

var (
	ErrEmptyResults = fmt.Errorf("empty results; cannot encode")
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

type EVMAutomationResultExtension21 struct {
	FastGasWei    *big.Int
	LinkNative    *big.Int
	FailureReason uint8 // this is not encoded, only pass along for the purpose of pipeline run
}

func (enc EVMAutomationEncoder21) Encode(results ...ocr2keepers.CheckResult) ([]byte, error) {
	if len(results) == 0 {
		return nil, ErrEmptyResults
	}

	report := automation_utils_2_1.KeeperRegistryBase21Report{
		UpkeepIds:    make([]*big.Int, len(results)),
		GasLimits:    make([]*big.Int, len(results)),
		Triggers:     make([][]byte, len(results)),
		PerformDatas: make([][]byte, len(results)),
	}

	for i, result := range results {
		ext, ok := result.Extension.(EVMAutomationResultExtension21)
		if !ok {
			return nil, fmt.Errorf("unexpected check result extension struct")
		}

		// only take these values from the first result
		// TODO: find a new way to get these values
		if i == 0 {
			report.FastGasWei = ext.FastGasWei
			report.LinkNative = ext.LinkNative
		}

		id, ok := new(big.Int).SetString(string(result.Payload.Upkeep.ID), 10)
		if !ok {
			return nil, fmt.Errorf("failed to parse big int from upkeep id: %s", string(result.Payload.Upkeep.ID))
		}

		report.UpkeepIds[i] = id
		report.GasLimits[i] = new(big.Int).SetUint64(result.GasAllocated)

		triggerW := triggerWrapper{
			BlockNum:  uint32(result.Payload.Trigger.BlockNumber),
			BlockHash: common.HexToHash(result.Payload.Trigger.BlockHash),
		}
		switch getUpkeepType(id.Bytes()) {
		case logTrigger:
			trExt, ok := result.Payload.Trigger.Extension.(logTriggerExtension)
			if !ok {
				return nil, fmt.Errorf("unrecognized trigger extension data")
			}
			hex, err := common.ParseHexOrString(trExt.TxHash)
			if err != nil {
				return nil, fmt.Errorf("tx hash parse error: %w", err)
			}
			triggerW.TxHash = common.BytesToHash(hex[:])
			triggerW.LogIndex = uint32(trExt.LogIndex)
		default:
		}
		trigger, err := enc.packer.PackTrigger(id, triggerW)
		if err != nil {
			return nil, fmt.Errorf("%w: failed to pack trigger", err)
		}
		report.Triggers[i] = trigger
		report.PerformDatas[i] = result.PerformData
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

// Extract the plugin will call this function to accept/transmit reports
func (enc EVMAutomationEncoder21) Extract(raw []byte) ([]ocr2keepers.ReportedUpkeep, error) {
	report, err := enc.packer.UnpackReport(raw)
	if err != nil {
		return nil, err
	}
	reportedUpkeeps := make([]ocr2keepers.ReportedUpkeep, len(report.UpkeepIds))
	for i, upkeepId := range report.UpkeepIds {
		triggerW, err := enc.packer.UnpackTrigger(upkeepId, report.Triggers[i])
		if err != nil {
			// TODO: log error and continue instead?
			return nil, fmt.Errorf("%w: failed to unpack trigger", err)
		}
		logExt := logTriggerExtension{}

		switch getUpkeepType(upkeepId.Bytes()) {
		case logTrigger:
			logExt.TxHash = common.BytesToHash(triggerW.TxHash[:]).Hex()
			logExt.LogIndex = int64(triggerW.LogIndex)
		default:
		}
		trigger := ocr2keepers.NewTrigger(
			int64(triggerW.BlockNum),
			string(triggerW.BlockHash[:]),
			logExt,
		)
		payload := ocr2keepers.NewUpkeepPayload(
			upkeepId,
			int(logTrigger),
			"",
			trigger,
			[]byte{},
		)
		reportedUpkeeps[i] = ocr2keepers.ReportedUpkeep{
			ID:          payload.ID,
			PerformData: report.PerformDatas[i],
		}
	}
	return reportedUpkeeps, nil
}

// validateReport checks that the report is valid, currently checking that all
// lists are the same length.
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

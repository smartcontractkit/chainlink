package core

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/smartcontractkit/chainlink-automation/pkg/v3/types"

	"github.com/ethereum/go-ethereum/accounts/abi"

	ac "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/automation_compatible_utils"
)

type triggerWrapper = ac.IAutomationV21PlusCommonLogTrigger

var ErrABINotParsable = fmt.Errorf("error parsing abi")

// PackTrigger packs the trigger data according to the upkeep type of the given id. it will remove the first 4 bytes of function selector.
func PackTrigger(id *big.Int, trig triggerWrapper) ([]byte, error) {
	var trigger []byte
	var err error

	// construct utils abi
	utilsABI, err := abi.JSON(strings.NewReader(ac.AutomationCompatibleUtilsABI))
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrABINotParsable, err)
	}

	// pack trigger based on upkeep type
	upkeepType, ok := getUpkeepTypeFromBigInt(id)
	if !ok {
		return nil, ErrInvalidUpkeepID
	}
	switch upkeepType {
	case types.ConditionTrigger:
		trig := ac.IAutomationV21PlusCommonConditionalTrigger{
			BlockNum:  trig.BlockNum,
			BlockHash: trig.BlockHash,
		}
		trigger, err = utilsABI.Pack("_conditionalTrigger", &trig)
	case types.LogTrigger:
		logTrig := ac.IAutomationV21PlusCommonLogTrigger{
			BlockNum:     trig.BlockNum,
			BlockHash:    trig.BlockHash,
			LogBlockHash: trig.LogBlockHash,
			LogIndex:     trig.LogIndex,
			TxHash:       trig.TxHash,
		}
		trigger, err = utilsABI.Pack("_logTrigger", &logTrig)
	default:
		err = fmt.Errorf("unknown trigger type: %d", upkeepType)
	}
	if err != nil {
		return nil, err
	}
	return trigger[4:], nil
}

// UnpackTrigger unpacks the trigger from the given raw data, according to the upkeep type of the given id.
func UnpackTrigger(id *big.Int, raw []byte) (triggerWrapper, error) {
	// construct utils abi
	utilsABI, err := abi.JSON(strings.NewReader(ac.AutomationCompatibleUtilsABI))
	if err != nil {
		return triggerWrapper{}, fmt.Errorf("%w: %s", ErrABINotParsable, err)
	}

	upkeepType, ok := getUpkeepTypeFromBigInt(id)
	if !ok {
		return triggerWrapper{}, ErrInvalidUpkeepID
	}
	switch upkeepType {
	case types.ConditionTrigger:
		unpacked, err := utilsABI.Methods["_conditionalTrigger"].Inputs.Unpack(raw)
		if err != nil {
			return triggerWrapper{}, fmt.Errorf("%w: failed to unpack conditional trigger", err)
		}
		converted, ok := abi.ConvertType(unpacked[0], new(ac.IAutomationV21PlusCommonConditionalTrigger)).(*ac.IAutomationV21PlusCommonConditionalTrigger)
		if !ok {
			return triggerWrapper{}, fmt.Errorf("failed to convert type")
		}
		triggerW := triggerWrapper{
			BlockNum: converted.BlockNum,
		}
		copy(triggerW.BlockHash[:], converted.BlockHash[:])
		return triggerW, nil
	case types.LogTrigger:
		unpacked, err := utilsABI.Methods["_logTrigger"].Inputs.Unpack(raw)
		if err != nil {
			return triggerWrapper{}, fmt.Errorf("%w: failed to unpack log trigger", err)
		}
		converted, ok := abi.ConvertType(unpacked[0], new(ac.IAutomationV21PlusCommonLogTrigger)).(*ac.IAutomationV21PlusCommonLogTrigger)
		if !ok {
			return triggerWrapper{}, fmt.Errorf("failed to convert type")
		}
		triggerW := triggerWrapper{
			BlockNum: converted.BlockNum,
			LogIndex: converted.LogIndex,
		}
		copy(triggerW.BlockHash[:], converted.BlockHash[:])
		copy(triggerW.TxHash[:], converted.TxHash[:])
		copy(triggerW.LogBlockHash[:], converted.LogBlockHash[:])
		return triggerW, nil
	default:
		return triggerWrapper{}, fmt.Errorf("unknown trigger type: %d", upkeepType)
	}
}

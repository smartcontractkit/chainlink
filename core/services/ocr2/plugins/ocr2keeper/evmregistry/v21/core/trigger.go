package core

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/smartcontractkit/chainlink-automation/pkg/v3/types"

	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/automation_utils_2_1"
)

type triggerWrapper = automation_utils_2_1.KeeperRegistryBase21LogTrigger

var ErrABINotParsable = fmt.Errorf("error parsing abi")

// PackTrigger packs the trigger data according to the upkeep type of the given id. it will remove the first 4 bytes of function selector.
func PackTrigger(id *big.Int, trig triggerWrapper) ([]byte, error) {
	var trigger []byte
	var err error

	// construct utils abi
	utilsABI, err := abi.JSON(strings.NewReader(automation_utils_2_1.AutomationUtilsABI))
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
		trig := automation_utils_2_1.KeeperRegistryBase21ConditionalTrigger{
			BlockNum:  trig.BlockNum,
			BlockHash: trig.BlockHash,
		}
		trigger, err = utilsABI.Pack("_conditionalTrigger", &trig)
	case types.LogTrigger:
		logTrig := automation_utils_2_1.KeeperRegistryBase21LogTrigger{
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
	utilsABI, err := abi.JSON(strings.NewReader(automation_utils_2_1.AutomationUtilsABI))
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
		converted, ok := abi.ConvertType(unpacked[0], new(automation_utils_2_1.KeeperRegistryBase21ConditionalTrigger)).(*automation_utils_2_1.KeeperRegistryBase21ConditionalTrigger)
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
		converted, ok := abi.ConvertType(unpacked[0], new(automation_utils_2_1.KeeperRegistryBase21LogTrigger)).(*automation_utils_2_1.KeeperRegistryBase21LogTrigger)
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

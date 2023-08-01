package core

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/automation_utils_2_1"
)

type triggerWrapper = automation_utils_2_1.KeeperRegistryBase21LogTrigger

var ErrABINotParsable = fmt.Errorf("error parsing abi")

type LogTriggerExtension struct {
	TxHash   string
	LogIndex int64
}

// according to the upkeep type of the given id.
func PackTrigger(id *big.Int, trig triggerWrapper) ([]byte, error) {
	var trigger []byte
	var err error

	// construct utils abi
	utilsABI, err := abi.JSON(strings.NewReader(automation_utils_2_1.AutomationUtilsABI))
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrABINotParsable, err)
	}

	// pack trigger based on upkeep type
	upkeepType := GetUpkeepType(id.Bytes())
	switch upkeepType {
	case ConditionTrigger:
		trig := automation_utils_2_1.KeeperRegistryBase21ConditionalTrigger{
			BlockNum:  trig.BlockNum,
			BlockHash: trig.BlockHash,
		}
		trigger, err = utilsABI.Pack("_conditionalTrigger", &trig)
	case LogTrigger:
		logTrig := automation_utils_2_1.KeeperRegistryBase21LogTrigger{
			BlockNum:  trig.BlockNum,
			BlockHash: trig.BlockHash,
			LogIndex:  trig.LogIndex,
			TxHash:    trig.TxHash,
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

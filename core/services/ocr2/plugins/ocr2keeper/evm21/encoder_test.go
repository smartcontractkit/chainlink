package evm

import (
	"errors"
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg"
	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/automation_utils_2_1"
	iregistry21 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_keeper_registry_master_wrapper_2_1"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21/logprovider"
)

func TestEVMAutomationEncoder21_Encode(t *testing.T) {
	keepersABI, err := abi.JSON(strings.NewReader(iregistry21.IKeeperRegistryMasterABI))
	assert.Nil(t, err)
	utilsABI, err := abi.JSON(strings.NewReader(automation_utils_2_1.AutomationUtilsABI))
	assert.Nil(t, err)
	encoder := EVMAutomationEncoder21{
		packer: NewEvmRegistryPackerV2_1(keepersABI, utilsABI),
	}

	tests := []struct {
		name        string
		results     []ocr2keepers.CheckResult
		reportSize  int
		expectedErr error
	}{
		{
			"happy flow single",
			[]ocr2keepers.CheckResult{
				newResult(1, ocr2keepers.UpkeepIdentifier(genUpkeepID(conditionTrigger, "10").String())),
			},
			640,
			nil,
		},
		{
			"happy flow multiple",
			[]ocr2keepers.CheckResult{
				newResult(1, ocr2keepers.UpkeepIdentifier(genUpkeepID(conditionTrigger, "10").String())),
				newResult(2, ocr2keepers.UpkeepIdentifier(genUpkeepID(logTrigger, "20").String())),
				newResult(3, ocr2keepers.UpkeepIdentifier(genUpkeepID(conditionTrigger, "30").String())),
			},
			1280,
			nil,
		},
		{
			"empty results",
			[]ocr2keepers.CheckResult{},
			0,
			ErrEmptyResults,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b, err := encoder.Encode(tc.results...)
			if tc.expectedErr != nil {
				assert.Equal(t, tc.expectedErr, err)
				return
			}
			assert.Nil(t, err)
			assert.Len(t, b, tc.reportSize)

			upkeeps, err := decode(encoder.packer, b)
			assert.Nil(t, err)
			assert.Len(t, upkeeps, len(tc.results))

			for i, u := range upkeeps {
				upkeep, ok := u.(EVMAutomationUpkeepResult21)
				assert.True(t, ok)

				assert.Equal(t, upkeep.Block, uint32(tc.results[i].Payload.Trigger.BlockNumber))
				assert.Equal(t, ocr2keepers.UpkeepIdentifier(upkeep.ID.String()), tc.results[i].Payload.Upkeep.ID)
				assert.Equal(t, upkeep.Eligible, tc.results[i].Eligible)
				assert.Equal(t, upkeep.PerformData, tc.results[i].PerformData)
				assert.Equal(t, upkeep.CheckBlockNumber, uint32(tc.results[i].Payload.Trigger.BlockNumber))
				assert.Equal(t, common.BytesToHash(upkeep.CheckBlockHash[:]), common.HexToHash(tc.results[i].Payload.Trigger.BlockHash))
			}
		})
	}
}

func TestEVMAutomationEncoder21_EncodeExtract(t *testing.T) {
	keepersABI, err := abi.JSON(strings.NewReader(iregistry21.IKeeperRegistryMasterABI))
	assert.Nil(t, err)
	utilsABI, err := abi.JSON(strings.NewReader(automation_utils_2_1.AutomationUtilsABI))
	assert.Nil(t, err)
	encoder := EVMAutomationEncoder21{
		packer: NewEvmRegistryPackerV2_1(keepersABI, utilsABI),
	}

	tests := []struct {
		name        string
		results     []ocr2keepers.CheckResult
		expectedErr error
	}{
		{
			"happy flow single",
			[]ocr2keepers.CheckResult{
				newResult(1, ocr2keepers.UpkeepIdentifier(genUpkeepID(logTrigger, "10").String())),
			},
			nil,
		},
		{
			"happy flow multiple",
			[]ocr2keepers.CheckResult{
				newResult(1, ocr2keepers.UpkeepIdentifier(genUpkeepID(logTrigger, "10").String())),
				newResult(2, ocr2keepers.UpkeepIdentifier(genUpkeepID(conditionTrigger, "20").String())),
				newResult(3, ocr2keepers.UpkeepIdentifier(genUpkeepID(logTrigger, "30").String())),
			},
			nil,
		},
		{
			"empty results",
			[]ocr2keepers.CheckResult{},
			errors.New(""),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b, _ := encoder.Encode(tc.results...)
			reportedUpkeeps, err := encoder.Extract(b)
			if tc.expectedErr != nil {
				assert.Error(t, err)
				return
			}
			assert.Nil(t, err)
			assert.Len(t, reportedUpkeeps, len(tc.results))

			for i, upkeep := range reportedUpkeeps {
				assert.Equal(t, tc.results[i].Payload.Upkeep.ID, upkeep.UpkeepID)
				assert.Equal(t, tc.results[i].Payload.Trigger.BlockHash, upkeep.Trigger.BlockHash)
				assert.Equal(t, tc.results[i].Payload.Trigger.BlockNumber, upkeep.Trigger.BlockNumber)
				assert.Equal(t, tc.results[i].PerformData, upkeep.PerformData)
			}
		})
	}
}

func newResult(block int64, id ocr2keepers.UpkeepIdentifier) ocr2keepers.CheckResult {
	logExt := logprovider.LogTriggerExtension{}
	tp := getUpkeepType(id)
	if tp == logTrigger {
		logExt.LogIndex = 1
		logExt.TxHash = "0x1234567890123456789012345678901234567890123456789012345678901234"
	}
	payload := ocr2keepers.UpkeepPayload{
		Upkeep: ocr2keepers.ConfiguredUpkeep{
			ID:   id,
			Type: int(tp),
		},
		Trigger: ocr2keepers.Trigger{
			BlockNumber: block,
			BlockHash:   hexutil.Encode([]byte{1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8}),
			Extension:   logExt,
		},
	}
	payload.ID = payload.GenerateID()
	return ocr2keepers.CheckResult{
		Payload:      payload,
		Eligible:     true,
		GasAllocated: 100,
		PerformData:  []byte("data0"),
		Extension: EVMAutomationResultExtension21{
			FastGasWei: big.NewInt(100),
			LinkNative: big.NewInt(100),
		},
	}
}

func decode(packer *evmRegistryPackerV2_1, raw []byte) ([]ocr2keepers.UpkeepResult, error) {
	report, err := packer.UnpackReport(raw)
	if err != nil {
		return nil, err
	}

	res := make([]ocr2keepers.UpkeepResult, len(report.UpkeepIds))

	for i := 0; i < len(report.UpkeepIds); i++ {
		trigger, err := packer.UnpackTrigger(report.UpkeepIds[i], report.Triggers[i])
		if err != nil {
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

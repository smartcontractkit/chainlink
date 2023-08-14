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
		name               string
		results            []ocr2keepers.CheckResult
		reportSize         int
		expectedFastGasWei int64
		expectedLinkNative int64
		expectedErr        error
	}{
		{
			"happy flow single",
			[]ocr2keepers.CheckResult{
				newResult(1, "1", ocr2keepers.UpkeepIdentifier(genUpkeepID(conditionTrigger, "10").String()), 1, 1),
			},
			640,
			1,
			1,
			nil,
		},
		{
			"happy flow multiple",
			[]ocr2keepers.CheckResult{
				newResult(1, "1", ocr2keepers.UpkeepIdentifier(genUpkeepID(conditionTrigger, "10").String()), 1, 1),
				newResult(2, "2", ocr2keepers.UpkeepIdentifier(genUpkeepID(logTrigger, "20").String()), 2, 2),
				newResult(3, "3", ocr2keepers.UpkeepIdentifier(genUpkeepID(conditionTrigger, "30").String()), 3, 3),
			},
			1280,
			3,
			3,
			nil,
		},
		{
			"happy flow highest block number first",
			[]ocr2keepers.CheckResult{
				newResult(3, "3", ocr2keepers.UpkeepIdentifier(genUpkeepID(conditionTrigger, "30").String()), 1000, 2000),
				newResult(2, "2", ocr2keepers.UpkeepIdentifier(genUpkeepID(logTrigger, "20").String()), 2, 2),
				newResult(1, "1", ocr2keepers.UpkeepIdentifier(genUpkeepID(conditionTrigger, "10").String()), 1, 1),
			},
			1280,
			1000,
			2000,
			nil,
		},
		{
			"empty results",
			[]ocr2keepers.CheckResult{},
			0,
			0,
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
				assert.Equal(t, upkeep.FastGasWei, big.NewInt(tc.expectedFastGasWei))
				assert.Equal(t, upkeep.LinkNative, big.NewInt(tc.expectedLinkNative))
			}
		})
	}
}

func TestEVMAutomationEncoder21_Encode_errors(t *testing.T) {
	keepersABI, err := abi.JSON(strings.NewReader(iregistry21.IKeeperRegistryMasterABI))
	assert.Nil(t, err)
	utilsABI, err := abi.JSON(strings.NewReader(automation_utils_2_1.AutomationUtilsABI))
	assert.Nil(t, err)
	encoder := EVMAutomationEncoder21{
		packer: NewEvmRegistryPackerV2_1(keepersABI, utilsABI),
	}

	t.Run("a non-EVMAutomationResultExtension21 extension causes an error", func(t *testing.T) {
		result := newResult(3, "3", ocr2keepers.UpkeepIdentifier(genUpkeepID(conditionTrigger, "30").String()), 1, 1)
		result.Extension = "invalid"
		b, err := encoder.Encode(result)
		assert.Nil(t, b)
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "unrecognized CheckResult extension data")
	})

	// TODO auto 4204: handle upkeep ID decoding properly
	/*
		t.Run("an invalid upkeep ID causes an error", func(t *testing.T) {
			result := newResult(3, "3", ocr2keepers.UpkeepIdentifier(genUpkeepID(conditionTrigger, "30").String()), 1, 1)
			result.Payload.Upkeep.ID = []byte("invalid")
			b, err := encoder.Encode(result)
			assert.Nil(t, b)
			assert.Error(t, err)
			assert.Equal(t, err.Error(), "failed to parse big int from upkeep id: invalid")
		})*/

	t.Run("a non-LogTriggerExtension extension causes an error", func(t *testing.T) {
		result := newResult(2, "2", ocr2keepers.UpkeepIdentifier(genUpkeepID(logTrigger, "20").String()), 1, 1)
		result.Payload.Trigger.Extension = "invalid"
		b, err := encoder.Encode(result)
		assert.Nil(t, b)
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "unrecognized trigger extension data")
	})

	t.Run("a non hex tx (empty) hash causes an error", func(t *testing.T) {
		result := newResult(2, "2", ocr2keepers.UpkeepIdentifier(genUpkeepID(logTrigger, "20").String()), 1, 1)
		result.Payload.Trigger.Extension = logprovider.LogTriggerExtension{
			TxHash: "",
		}
		b, err := encoder.Encode(result)
		assert.Nil(t, b)
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "tx hash parse error: empty hex string")
	})

	t.Run("an invalid upkeep type causes an error", func(t *testing.T) {
		result := newResult(2, "2", ocr2keepers.UpkeepIdentifier(genUpkeepID(5, "20").String()), 1, 1)
		result.Payload.Trigger.Extension = logprovider.LogTriggerExtension{
			TxHash: "",
		}
		b, err := encoder.Encode(result)
		assert.Nil(t, b)
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "unknown trigger type: 5: failed to pack trigger")
	})

	t.Run("an invalid result extension causes an error", func(t *testing.T) {
		result := newResult(2, "2", ocr2keepers.UpkeepIdentifier(genUpkeepID(5, "20").String()), 1, 1)
		result.Extension = struct{}{}
		b, err := encoder.Encode(result)
		assert.Nil(t, b)
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "unexpected check result extension struct")
	})

	t.Run("invalid trigger extension causes an error", func(t *testing.T) {
		result := newResult(2, "2", ocr2keepers.UpkeepIdentifier(genUpkeepID(logTrigger, "20").String()), 1, 1)
		result.Payload.Trigger.Extension = struct{}{}
		b, err := encoder.Encode(result)
		assert.Nil(t, b)
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "unrecognized trigger extension data")
	})
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
				newResult(1, "1", ocr2keepers.UpkeepIdentifier(genUpkeepID(logTrigger, "10").String()), 1, 1),
			},
			nil,
		},
		{
			"happy flow multiple",
			[]ocr2keepers.CheckResult{
				newResult(1, "1", ocr2keepers.UpkeepIdentifier(genUpkeepID(logTrigger, "10").String()), 1, 1),
				newResult(2, "2", ocr2keepers.UpkeepIdentifier(genUpkeepID(conditionTrigger, "20").String()), 1, 1),
				newResult(3, "3", ocr2keepers.UpkeepIdentifier(genUpkeepID(logTrigger, "30").String()), 1, 1),
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

func TestEVMAutomationEncoder21_EncodeExtract_errors(t *testing.T) {
	keepersABI, err := abi.JSON(strings.NewReader(iregistry21.IKeeperRegistryMasterABI))
	assert.Nil(t, err)
	utilsABI, err := abi.JSON(strings.NewReader(automation_utils_2_1.AutomationUtilsABI))
	assert.Nil(t, err)
	encoder := EVMAutomationEncoder21{
		packer: NewEvmRegistryPackerV2_1(keepersABI, utilsABI),
	}

	t.Run("attempting to decode a report with an invalid trigger type returns an error", func(t *testing.T) {
		encodedConditionTrigger := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 192, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 9, 49, 48, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128, 18, 52, 86, 120, 144, 18, 52, 86, 120, 144, 18, 52, 86, 120, 144, 18, 52, 86, 120, 144, 18, 52, 86, 120, 144, 18, 52, 86, 120, 144, 18, 52, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 5, 100, 97, 116, 97, 48, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		reportedUpkeeps, err := encoder.Extract(encodedConditionTrigger)
		assert.Nil(t, reportedUpkeeps)
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "unknown trigger type: 9: failed to unpack trigger")
	})
}

func TestBlockKeyHelper_MakeBlockKey(t *testing.T) {
	helper := BlockKeyHelper[int64]{}
	blockKey := helper.MakeBlockKey(1234)
	assert.Equal(t, blockKey, ocr2keepers.BlockKey("1234"))
}

func TestUpkeepKeyHelper_MakeBlockKey(t *testing.T) {
	helper := UpkeepKeyHelper[int64]{}
	blockKey := helper.MakeUpkeepKey(1234, big.NewInt(5678))
	assert.Equal(t, blockKey, ocr2keepers.UpkeepKey("1234|5678"))
}

func newResult(block int64, checkBlock ocr2keepers.BlockKey, id ocr2keepers.UpkeepIdentifier, fastGasWei, linkNative int64) ocr2keepers.CheckResult {
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
		CheckBlock: checkBlock,
	}
	payload.ID = payload.GenerateID()
	return ocr2keepers.CheckResult{
		Payload:      payload,
		Eligible:     true,
		GasAllocated: 100,
		PerformData:  []byte("data0"),
		Extension: EVMAutomationResultExtension21{
			FastGasWei: big.NewInt(fastGasWei),
			LinkNative: big.NewInt(linkNative),
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

func Test_decodeExtensions(t *testing.T) {
	t.Run("a nil result returns an error", func(t *testing.T) {
		err := decodeExtensions(nil)
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "non-nil value expected in decoding")
	})

	t.Run("two types of extension do not require decoding", func(t *testing.T) {
		err := decodeExtensions(&ocr2keepers.CheckResult{
			Payload: ocr2keepers.UpkeepPayload{
				Trigger: ocr2keepers.Trigger{
					Extension: logprovider.LogTriggerExtension{},
				},
			},
			Extension: EVMAutomationResultExtension21{},
		})
		assert.NoError(t, err)
	})

	t.Run("empty structs do not require decoding", func(t *testing.T) {
		err := decodeExtensions(&ocr2keepers.CheckResult{
			Payload: ocr2keepers.UpkeepPayload{
				Trigger: ocr2keepers.Trigger{
					Extension: struct{}{},
				},
			},
			Extension: struct{}{},
		})
		assert.NoError(t, err)
	})

	t.Run("CheckResult extension as json gets decoded", func(t *testing.T) {
		err := decodeExtensions(&ocr2keepers.CheckResult{
			Payload: ocr2keepers.UpkeepPayload{
				Trigger: ocr2keepers.Trigger{
					Extension: struct{}{},
				},
			},
			Extension: []byte(`{"FastGasWei":1,"LinkNative":2,"FailureReason":3}`),
		})
		assert.NoError(t, err)
	})

	t.Run("trigger extension as json for logTrigger gets decoded", func(t *testing.T) {
		err := decodeExtensions(&ocr2keepers.CheckResult{
			Payload: ocr2keepers.UpkeepPayload{
				Upkeep: ocr2keepers.ConfiguredUpkeep{
					ID: ocr2keepers.UpkeepIdentifier(genUpkeepID(logTrigger, "111").Bytes()),
				},
				Trigger: ocr2keepers.Trigger{
					Extension: []byte(`{"TxHash": "abc", "LogIndex": 123}`),
				},
			},
			Extension: struct{}{},
		})
		assert.NoError(t, err)
	})

	t.Run("trigger extension as invalid json returns an error for logTrigger", func(t *testing.T) {
		err := decodeExtensions(&ocr2keepers.CheckResult{
			Payload: ocr2keepers.UpkeepPayload{
				Upkeep: ocr2keepers.ConfiguredUpkeep{
					ID: ocr2keepers.UpkeepIdentifier(genUpkeepID(logTrigger, "111").Bytes()),
				},
				Trigger: ocr2keepers.Trigger{
					Extension: []byte(`invalid`),
				},
			},
			Extension: struct{}{},
		})
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "invalid character 'i' looking for beginning of value: json encoded values do not match LogTriggerExtension struct: invalid")
	})

	t.Run("trigger extension for conditionTrigger requires no decoding", func(t *testing.T) {
		err := decodeExtensions(&ocr2keepers.CheckResult{
			Payload: ocr2keepers.UpkeepPayload{
				Upkeep: ocr2keepers.ConfiguredUpkeep{
					ID: ocr2keepers.UpkeepIdentifier(genUpkeepID(conditionTrigger, "111").Bytes()),
				},
				Trigger: ocr2keepers.Trigger{
					Extension: []byte(""),
				},
			},
			Extension: struct{}{},
		})
		assert.NoError(t, err)
	})

	t.Run("trigger extension for an unknown upkeep type returns an error", func(t *testing.T) {
		err := decodeExtensions(&ocr2keepers.CheckResult{
			Payload: ocr2keepers.UpkeepPayload{
				Upkeep: ocr2keepers.ConfiguredUpkeep{
					ID: ocr2keepers.UpkeepIdentifier(genUpkeepID(99, "111").Bytes()),
				},
				Trigger: ocr2keepers.Trigger{
					Extension: []byte(""),
				},
			},
			Extension: struct{}{},
		})
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "unknown upkeep type")
	})

	t.Run("CheckResult extension as invalid json returns an error", func(t *testing.T) {
		err := decodeExtensions(&ocr2keepers.CheckResult{
			Payload: ocr2keepers.UpkeepPayload{
				Trigger: ocr2keepers.Trigger{
					Extension: struct{}{},
				},
			},
			Extension: []byte(`invalid`),
		})
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "invalid character 'i' looking for beginning of value: json encoded values do not match EVMAutomationResultExtension21 struct")
	})

	t.Run("unrecognized trigger extension returns an error", func(t *testing.T) {
		err := decodeExtensions(&ocr2keepers.CheckResult{
			Payload: ocr2keepers.UpkeepPayload{
				Trigger: ocr2keepers.Trigger{
					Extension: "invalid",
				},
			},
		})
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "unrecognized trigger extension data")
	})

	t.Run("unrecognized CheckResult extension returns an error", func(t *testing.T) {
		err := decodeExtensions(&ocr2keepers.CheckResult{
			Payload: ocr2keepers.UpkeepPayload{
				Trigger: ocr2keepers.Trigger{
					Extension: logprovider.LogTriggerExtension{},
				},
			},
			Extension: "invalid",
		})
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "unrecognized CheckResult extension data")
	})
}

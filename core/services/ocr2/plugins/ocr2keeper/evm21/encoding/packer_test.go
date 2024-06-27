package encoding

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg/v3/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/automation_utils_2_1"
	automation21Utils "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/automation_utils_2_1"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21/core"
)

func TestPacker_PackReport(t *testing.T) {
	for _, tc := range []struct {
		name       string
		report     automation21Utils.KeeperRegistryBase21Report
		expectsErr bool
		wantErr    error
		wantBytes  int
	}{
		{
			name: "all non-nil values get encoded to a byte array of a specific length",
			report: automation21Utils.KeeperRegistryBase21Report{
				FastGasWei: big.NewInt(0),
				LinkNative: big.NewInt(0),
				UpkeepIds:  []*big.Int{big.NewInt(3)},
				GasLimits:  []*big.Int{big.NewInt(4)},
				Triggers: [][]byte{
					{5},
				},
				PerformDatas: [][]byte{
					{6},
				},
			},
			wantBytes: 608,
		},
		{
			name: "if upkeep IDs are nil, the packed report is smaller",
			report: automation21Utils.KeeperRegistryBase21Report{
				FastGasWei: big.NewInt(1),
				LinkNative: big.NewInt(2),
				UpkeepIds:  nil,
				GasLimits:  []*big.Int{big.NewInt(4)},
				Triggers: [][]byte{
					{5},
				},
				PerformDatas: [][]byte{
					{6},
				},
			},
			wantBytes: 576,
		},
		{
			name: "if gas limits are nil, the packed report is smaller",
			report: automation21Utils.KeeperRegistryBase21Report{
				FastGasWei: big.NewInt(1),
				LinkNative: big.NewInt(2),
				UpkeepIds:  []*big.Int{big.NewInt(3)},
				GasLimits:  nil,
				Triggers: [][]byte{
					{5},
				},
				PerformDatas: [][]byte{
					{6},
				},
			},
			wantBytes: 576,
		},
		{
			name: "if perform datas are nil, the packed report is smaller",
			report: automation21Utils.KeeperRegistryBase21Report{
				FastGasWei: big.NewInt(1),
				LinkNative: big.NewInt(2),
				UpkeepIds:  []*big.Int{big.NewInt(3)},
				GasLimits:  []*big.Int{big.NewInt(4)},
				Triggers: [][]byte{
					{5},
				},
				PerformDatas: nil,
			},
			wantBytes: 512,
		},
		{
			name: "if triggers are nil, the packed report is smaller",
			report: automation21Utils.KeeperRegistryBase21Report{
				FastGasWei: big.NewInt(1),
				LinkNative: big.NewInt(2),
				UpkeepIds:  []*big.Int{big.NewInt(3)},
				GasLimits:  []*big.Int{big.NewInt(4)},
				Triggers:   nil,
				PerformDatas: [][]byte{
					{6},
				},
			},
			wantBytes: 512,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			packer := NewAbiPacker()
			bytes, err := packer.PackReport(tc.report)
			if tc.expectsErr {
				assert.Error(t, err)
				assert.Equal(t, tc.wantErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.wantBytes, len(bytes))
			}
		})
	}
}

func TestPacker_UnpackCheckResults(t *testing.T) {
	uid, _ := new(big.Int).SetString("1843548457736589226156809205796175506139185429616502850435279853710366065936", 10)
	upkeepId := ocr2keepers.UpkeepIdentifier{}
	upkeepId.FromBigInt(uid)
	p1, _ := core.NewUpkeepPayload(uid, ocr2keepers.NewTrigger(19447615, common.HexToHash("0x0")), []byte{})

	p2, _ := core.NewUpkeepPayload(uid, ocr2keepers.NewTrigger(19448272, common.HexToHash("0x0")), []byte{})

	tests := []struct {
		Name           string
		Payload        ocr2keepers.UpkeepPayload
		RawData        string
		ExpectedResult ocr2keepers.CheckResult
		ExpectedError  error
	}{
		{
			Name:    "upkeep not needed",
			Payload: p1,
			RawData: "0x000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000421c000000000000000000000000000000000000000000000000000000003b9aca00000000000000000000000000000000000000000000000000000c8caf37f3b3890000000000000000000000000000000000000000000000000000000000000000",
			ExpectedResult: ocr2keepers.CheckResult{
				UpkeepID:            upkeepId,
				Eligible:            false,
				IneligibilityReason: uint8(UpkeepFailureReasonUpkeepNotNeeded),
				Trigger:             ocr2keepers.NewLogTrigger(ocr2keepers.BlockNumber(19447615), [32]byte{}, nil),
				WorkID:              "e54c524132d9c8d87b7e43b76f6d769face19ffd2ff93fc24f123dd745d3ce1e",
				PerformData:         nil,
				FastGasWei:          big.NewInt(1000000000),
				LinkNative:          big.NewInt(3532383906411401),
			},
		},
		{
			Name:    "target check reverted",
			Payload: p2,
			RawData: "0x000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c000000000000000000000000000000000000000000000000000000000000000030000000000000000000000000000000000000000000000000000000000007531000000000000000000000000000000000000000000000000000000003b9aca00000000000000000000000000000000000000000000000000000c8caf37f3b3890000000000000000000000000000000000000000000000000000000000000300000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000008914039bf676e20aad43a5642485e666575ed0d927a4b5679745e947e7d125ee2687c10000000000000000000000000000000000000000000000000000000000000060000000000000000000000000000000000000000000000000000000000000024462e8a50d00000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000001c0000000000000000000000000000000000000000000000000000000000128c1d000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000009666565644944537472000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000184554482d5553442d415242495452554d2d544553544e4554000000000000000000000000000000000000000000000000000000000000000000000000000000184254432d5553442d415242495452554d2d544553544e45540000000000000000000000000000000000000000000000000000000000000000000000000000000b626c6f636b4e756d6265720000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000014000000000000000000000000000000000000006400000000000000000000000000000000000000000000000000000000000000000000000000000000",
			ExpectedResult: ocr2keepers.CheckResult{
				UpkeepID:            upkeepId,
				Eligible:            false,
				IneligibilityReason: uint8(UpkeepFailureReasonTargetCheckReverted),
				Trigger:             ocr2keepers.NewLogTrigger(ocr2keepers.BlockNumber(19448272), [32]byte{}, nil),
				WorkID:              "e54c524132d9c8d87b7e43b76f6d769face19ffd2ff93fc24f123dd745d3ce1e",
				PerformData:         []byte{98, 232, 165, 13, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 160, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 224, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 192, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 40, 193, 208, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 9, 102, 101, 101, 100, 73, 68, 83, 116, 114, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 24, 69, 84, 72, 45, 85, 83, 68, 45, 65, 82, 66, 73, 84, 82, 85, 77, 45, 84, 69, 83, 84, 78, 69, 84, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 24, 66, 84, 67, 45, 85, 83, 68, 45, 65, 82, 66, 73, 84, 82, 85, 77, 45, 84, 69, 83, 84, 78, 69, 84, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 11, 98, 108, 111, 99, 107, 78, 117, 109, 98, 101, 114, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 20, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				FastGasWei:          big.NewInt(1000000000),
				LinkNative:          big.NewInt(3532383906411401),
			},
		},
		{
			Name:    "decode failed",
			Payload: p2,
			RawData: "invalid_raw_data",
			ExpectedResult: ocr2keepers.CheckResult{
				UpkeepID:               upkeepId,
				PipelineExecutionState: uint8(PackUnpackDecodeFailed),
				Trigger:                p2.Trigger,
				WorkID:                 p2.WorkID,
			},
			ExpectedError: fmt.Errorf("upkeepId %s failed to decode checkUpkeep result invalid_raw_data: hex string without 0x prefix", p2.UpkeepID.String()),
		},
		{
			Name:    "unpack failed",
			Payload: p2,
			RawData: "0x123123",
			ExpectedResult: ocr2keepers.CheckResult{
				UpkeepID:               upkeepId,
				PipelineExecutionState: uint8(PackUnpackDecodeFailed),
				Trigger:                p2.Trigger,
				WorkID:                 p2.WorkID,
			},
			ExpectedError: fmt.Errorf("upkeepId %s failed to unpack checkUpkeep result 0x123123: abi: cannot marshal in to go type: length insufficient 3 require 32", p2.UpkeepID.String()),
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			packer := NewAbiPacker()
			rs, err := packer.UnpackCheckResult(test.Payload, test.RawData)
			if test.ExpectedError != nil {
				assert.Equal(t, test.ExpectedError.Error(), err.Error())
			} else {
				assert.Nil(t, err)
			}
			assert.Equal(t, test.ExpectedResult.UpkeepID, rs.UpkeepID)
			assert.Equal(t, test.ExpectedResult.Eligible, rs.Eligible)
			assert.Equal(t, test.ExpectedResult.Trigger, rs.Trigger)
			assert.Equal(t, test.ExpectedResult.WorkID, rs.WorkID)
			assert.Equal(t, test.ExpectedResult.IneligibilityReason, rs.IneligibilityReason)
			assert.Equal(t, test.ExpectedResult.PipelineExecutionState, rs.PipelineExecutionState)
		})
	}
}

func TestPacker_UnpackPerformResult(t *testing.T) {
	tests := []struct {
		Name    string
		RawData string
		State   PipelineExecutionState
	}{
		{
			Name:    "unpack success",
			RawData: "0x0000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000a52d",
			State:   NoPipelineError,
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			packer := NewAbiPacker()
			state, rs, err := packer.UnpackPerformResult(test.RawData)
			assert.Nil(t, err)
			assert.True(t, rs)
			assert.Equal(t, test.State, state)
		})
	}
}

func TestPacker_UnpackCheckCallbackResult(t *testing.T) {
	tests := []struct {
		Name          string
		CallbackResp  []byte
		UpkeepNeeded  bool
		PerformData   []byte
		FailureReason uint8
		GasUsed       *big.Int
		ErrorString   string
		State         PipelineExecutionState
	}{
		{
			Name:          "unpack upkeep needed",
			CallbackResp:  []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 46, 20, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 96, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 16, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 16, 10, 11, 21, 31, 41, 15, 16, 17, 18, 19, 13, 14, 12, 13, 14, 15, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8, 120, 111, 101, 122, 90, 54, 44, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			UpkeepNeeded:  true,
			PerformData:   []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 16, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 16, 10, 11, 21, 31, 41, 15, 16, 17, 18, 19, 13, 14, 12, 13, 14, 15, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8, 120, 111, 101, 122, 90, 54, 44, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			FailureReason: uint8(UpkeepFailureReasonNone),
			GasUsed:       big.NewInt(11796),
		},
		{
			Name:          "unpack upkeep not needed",
			CallbackResp:  []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 50, 208, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 96, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 16, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 16, 10, 11, 21, 31, 41, 15, 16, 17, 18, 19, 13, 14, 12, 13, 14, 15, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8, 120, 111, 101, 122, 90, 54, 44, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			UpkeepNeeded:  false,
			PerformData:   []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 16, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 16, 10, 11, 21, 31, 41, 15, 16, 17, 18, 19, 13, 14, 12, 13, 14, 15, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8, 120, 111, 101, 122, 90, 54, 44, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			FailureReason: uint8(UpkeepFailureReasonUpkeepNotNeeded),
			GasUsed:       big.NewInt(13008),
		},
		{
			Name:         "unpack malformed data",
			CallbackResp: []byte{0, 0, 0, 23, 4, 163, 66, 91, 228, 102, 200, 84, 144, 233, 218, 44, 168, 192, 191, 253, 0, 0, 0, 0, 20, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			UpkeepNeeded: false,
			PerformData:  nil,
			ErrorString:  "abi: improperly encoded boolean value: unpack checkUpkeep return: ",
			State:        PackUnpackDecodeFailed,
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			packer := NewAbiPacker()

			state, needed, pd, failureReason, gasUsed, err := packer.UnpackCheckCallbackResult(test.CallbackResp)

			if test.ErrorString != "" {
				assert.EqualError(t, err, test.ErrorString+hexutil.Encode(test.CallbackResp))
			} else {
				assert.Nil(t, err)
			}
			assert.Equal(t, test.UpkeepNeeded, needed)
			assert.Equal(t, test.PerformData, pd)
			assert.Equal(t, test.FailureReason, failureReason)
			assert.Equal(t, test.GasUsed, gasUsed)
			assert.Equal(t, test.State, state)
		})
	}
}

func TestPacker_UnpackLogTriggerConfig(t *testing.T) {
	tests := []struct {
		name    string
		raw     []byte
		res     automation21Utils.LogTriggerConfig
		errored bool
	}{
		{
			"happy flow",
			func() []byte {
				b, _ := hexutil.Decode("0x0000000000000000000000007456fadf415b7c34b1182bd20b0537977e945e3e00000000000000000000000000000000000000000000000000000000000000003d53a39550e04688065827f3bb86584cb007ab9ebca7ebd528e7301c9c31eb5d000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
				return b
			}(),
			automation21Utils.LogTriggerConfig{
				ContractAddress: common.HexToAddress("0x7456FadF415b7c34B1182Bd20B0537977e945e3E"),
				Topic0:          [32]uint8{0x3d, 0x53, 0xa3, 0x95, 0x50, 0xe0, 0x46, 0x88, 0x6, 0x58, 0x27, 0xf3, 0xbb, 0x86, 0x58, 0x4c, 0xb0, 0x7, 0xab, 0x9e, 0xbc, 0xa7, 0xeb, 0xd5, 0x28, 0xe7, 0x30, 0x1c, 0x9c, 0x31, 0xeb, 0x5d},
			},
			false,
		},
		{
			"invalid",
			func() []byte {
				b, _ := hexutil.Decode("0x000000000000000000000000b1182bd20b0537977e945e3e00000000000000000000000000000000000000000000000000000000000000003d53a39550e04688065827f3bb86584cb007ab9ebca7ebd528e7301c9c31eb5d000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
				return b
			}(),
			automation21Utils.LogTriggerConfig{},
			true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			packer := NewAbiPacker()
			res, err := packer.UnpackLogTriggerConfig(tc.raw)
			if tc.errored {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.res, res)
			}
		})
	}
}

func TestPacker_PackReport_UnpackReport(t *testing.T) {
	report := automation_utils_2_1.KeeperRegistryBase21Report{
		FastGasWei:   big.NewInt(1),
		LinkNative:   big.NewInt(1),
		UpkeepIds:    []*big.Int{big.NewInt(1), big.NewInt(2)},
		GasLimits:    []*big.Int{big.NewInt(100), big.NewInt(200)},
		Triggers:     [][]byte{{1, 2, 3, 4}, {5, 6, 7, 8}},
		PerformDatas: [][]byte{{5, 6, 7, 8}, {1, 2, 3, 4}},
	}
	packer := NewAbiPacker()
	res, err := packer.PackReport(report)
	require.NoError(t, err)
	report2, err := packer.UnpackReport(res)
	require.NoError(t, err)
	assert.Equal(t, report, report2)
	expected := "0x00000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000000120000000000000000000000000000000000000000000000000000000000000018000000000000000000000000000000000000000000000000000000000000002600000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000006400000000000000000000000000000000000000000000000000000000000000c800000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000040102030400000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000405060708000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000800000000000000000000000000000000000000000000000000000000000000004050607080000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000040102030400000000000000000000000000000000000000000000000000000000"
	assert.Equal(t, hexutil.Encode(res), expected)
}

func TestPacker_PackGetUpkeepPrivilegeConfig(t *testing.T) {
	tests := []struct {
		name     string
		upkeepId *big.Int
		raw      []byte
		errored  bool
	}{
		{
			name: "happy path",
			upkeepId: func() *big.Int {
				id, _ := new(big.Int).SetString("52236098515066839510538748191966098678939830769967377496848891145101407612976", 10)

				return id
			}(),
			raw: func() []byte {
				b, _ := hexutil.Decode("0x19d97a94737c9583000000000000000000000001ea8ed6d0617dd5b3b87374020efaf030")

				return b
			}(),
			errored: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			packer := NewAbiPacker()

			b, err := packer.PackGetUpkeepPrivilegeConfig(test.upkeepId)

			if !test.errored {
				require.NoError(t, err, "no error expected from packing")

				assert.Equal(t, test.raw, b, "raw bytes for output should match expected")
			} else {
				assert.NotNil(t, err, "error expected from packing function")
			}
		})
	}
}

func TestPacker_UnpackGetUpkeepPrivilegeConfig(t *testing.T) {
	tests := []struct {
		name    string
		raw     []byte
		errored bool
	}{
		{
			name: "happy path",
			raw: func() []byte {
				b, _ := hexutil.Decode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000177b226d657263757279456e61626c6564223a747275657d000000000000000000")

				return b
			}(),
			errored: false,
		},
		{
			name: "error empty config",
			raw: func() []byte {
				b, _ := hexutil.Decode("0x")

				return b
			}(),
			errored: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			packer := NewAbiPacker()

			b, err := packer.UnpackGetUpkeepPrivilegeConfig(test.raw)

			if !test.errored {
				require.NoError(t, err, "should unpack bytes from abi encoded value")

				// the actual struct to unmarshal into is not available to this
				// package so basic json encoding is the limit of the following test
				var data map[string]interface{}
				err = json.Unmarshal(b, &data)

				assert.NoError(t, err, "packed data should unmarshal using json encoding")
				assert.Equal(t, []byte(`{"mercuryEnabled":true}`), b)
			} else {
				assert.NotNil(t, err, "error expected from unpack function")
			}
		})
	}
}

func TestPacker_DecodeStreamsLookupRequest(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected *StreamsLookupError
		state    PipelineExecutionState
		err      error
	}{
		{
			name: "success - decode to streams lookup",
			data: hexutil.MustDecode("0xf055e4a200000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000002400000000000000000000000000000000000000000000000000000000002435eb50000000000000000000000000000000000000000000000000000000000000280000000000000000000000000000000000000000000000000000000000000000966656564496448657800000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000c000000000000000000000000000000000000000000000000000000000000000423078343535343438326435353533343432643431353234323439353435323535346432643534343535333534346534353534303030303030303030303030303030300000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000042307834323534343332643535353334343264343135323432343935343532353534643264353434353533353434653435353430303030303030303030303030303030000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000b626c6f636b4e756d62657200000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000140000000000000000000000000000000000000064000000000000000000000000"),
			expected: &StreamsLookupError{
				FeedParamKey: "feedIdHex",
				Feeds:        []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000", "0x4254432d5553442d415242495452554d2d544553544e45540000000000000000"},
				TimeParamKey: "blockNumber",
				Time:         big.NewInt(37969589),
				ExtraData:    []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100},
			},
		},
		{
			name: "failure - unpack error",
			data: []byte{1, 2, 3, 4},
			err:  errors.New("unpack error: invalid data for unpacking"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			packer := NewAbiPacker()
			fl, err := packer.DecodeStreamsLookupRequest(tt.data)
			assert.Equal(t, tt.expected, fl)
			if tt.err != nil {
				assert.Equal(t, tt.err.Error(), err.Error())
			}
		})
	}
}

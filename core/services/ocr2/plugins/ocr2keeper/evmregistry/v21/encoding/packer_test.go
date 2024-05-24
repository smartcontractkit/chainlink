package encoding

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	ocr2keepers "github.com/smartcontractkit/chainlink-common/pkg/types/automation"

	ac "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/automation_compatible_utils"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21/core"
)

func TestPacker_PackReport(t *testing.T) {
	for _, tc := range []struct {
		name       string
		report     ac.IAutomationV21PlusCommonReport
		expectsErr bool
		wantErr    error
		wantBytes  int
	}{
		{
			name: "all non-nil values get encoded to a byte array of a specific length",
			report: ac.IAutomationV21PlusCommonReport{
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
			report: ac.IAutomationV21PlusCommonReport{
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
			report: ac.IAutomationV21PlusCommonReport{
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
			report: ac.IAutomationV21PlusCommonReport{
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
			report: ac.IAutomationV21PlusCommonReport{
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

func TestPacker_UnpackLogTriggerConfig(t *testing.T) {
	tests := []struct {
		name    string
		raw     []byte
		res     ac.IAutomationV21PlusCommonLogTriggerConfig
		errored bool
	}{
		{
			"happy flow",
			func() []byte {
				b, _ := hexutil.Decode("0x0000000000000000000000007456fadf415b7c34b1182bd20b0537977e945e3e00000000000000000000000000000000000000000000000000000000000000003d53a39550e04688065827f3bb86584cb007ab9ebca7ebd528e7301c9c31eb5d000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
				return b
			}(),
			ac.IAutomationV21PlusCommonLogTriggerConfig{
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
			ac.IAutomationV21PlusCommonLogTriggerConfig{},
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
	report := ac.IAutomationV21PlusCommonReport{
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

package evm

import (
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg"
	"github.com/stretchr/testify/assert"

	automation21Utils "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/automation_utils_2_1"
	iregistry21 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_keeper_registry_master_wrapper_2_1"
)

func TestUnpackCheckResults(t *testing.T) {
	upkeepId, _ := new(big.Int).SetString("1843548457736589226156809205796175506139185429616502850435279853710366065936", 10)

	tests := []struct {
		Name           string
		Payload        ocr2keepers.UpkeepPayload
		RawData        string
		ExpectedResult EVMAutomationUpkeepResult21
	}{
		{
			Name:    "upkeep not needed",
			Payload: ocr2keepers.NewUpkeepPayload(upkeepId, 0, "19447615", ocr2keepers.NewTrigger(19447615, "0x0", struct{}{}), []byte{}),
			RawData: "0x000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000421c000000000000000000000000000000000000000000000000000000003b9aca00000000000000000000000000000000000000000000000000000c8caf37f3b3890000000000000000000000000000000000000000000000000000000000000000",
			ExpectedResult: EVMAutomationUpkeepResult21{
				Block:            19447615,
				ID:               upkeepId,
				Eligible:         false,
				FailureReason:    UPKEEP_FAILURE_REASON_UPKEEP_NOT_NEEDED,
				GasUsed:          big.NewInt(16924),
				PerformData:      nil,
				FastGasWei:       big.NewInt(1000000000),
				LinkNative:       big.NewInt(3532383906411401),
				CheckBlockNumber: 0,
				CheckBlockHash:   [32]byte{},
				ExecuteGas:       5000000,
			},
		},
		{
			Name:    "target check reverted",
			Payload: ocr2keepers.NewUpkeepPayload(upkeepId, 0, "19448272", ocr2keepers.NewTrigger(19448272, "0x0", struct{}{}), []byte{}),
			RawData: "0x000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c000000000000000000000000000000000000000000000000000000000000000030000000000000000000000000000000000000000000000000000000000007531000000000000000000000000000000000000000000000000000000003b9aca00000000000000000000000000000000000000000000000000000c8caf37f3b3890000000000000000000000000000000000000000000000000000000000000300000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000008914039bf676e20aad43a5642485e666575ed0d927a4b5679745e947e7d125ee2687c10000000000000000000000000000000000000000000000000000000000000060000000000000000000000000000000000000000000000000000000000000024462e8a50d00000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000001c0000000000000000000000000000000000000000000000000000000000128c1d000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000009666565644944537472000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000184554482d5553442d415242495452554d2d544553544e4554000000000000000000000000000000000000000000000000000000000000000000000000000000184254432d5553442d415242495452554d2d544553544e45540000000000000000000000000000000000000000000000000000000000000000000000000000000b626c6f636b4e756d6265720000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000014000000000000000000000000000000000000006400000000000000000000000000000000000000000000000000000000000000000000000000000000",
			ExpectedResult: EVMAutomationUpkeepResult21{
				Block:            19448272,
				ID:               upkeepId,
				Eligible:         false,
				FailureReason:    UPKEEP_FAILURE_REASON_TARGET_CHECK_REVERTED,
				GasUsed:          big.NewInt(30001),
				PerformData:      []byte{98, 232, 165, 13, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 160, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 224, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 192, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 40, 193, 208, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 9, 102, 101, 101, 100, 73, 68, 83, 116, 114, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 24, 69, 84, 72, 45, 85, 83, 68, 45, 65, 82, 66, 73, 84, 82, 85, 77, 45, 84, 69, 83, 84, 78, 69, 84, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 24, 66, 84, 67, 45, 85, 83, 68, 45, 65, 82, 66, 73, 84, 82, 85, 77, 45, 84, 69, 83, 84, 78, 69, 84, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 11, 98, 108, 111, 99, 107, 78, 117, 109, 98, 101, 114, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 20, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				FastGasWei:       big.NewInt(1000000000),
				LinkNative:       big.NewInt(3532383906411401),
				CheckBlockNumber: 8983555,
				CheckBlockHash:   [32]byte{155, 246, 118, 226, 10, 173, 67, 165, 100, 36, 133, 230, 102, 87, 94, 208, 217, 39, 164, 181, 103, 151, 69, 233, 71, 231, 209, 37, 238, 38, 135, 193},
				ExecuteGas:       5000000,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			packer, err := newPacker()
			assert.NoError(t, err)
			rs, err := packer.UnpackCheckResult(test.Payload, test.RawData)
			assert.Nil(t, err)
			assert.Equal(t, test.ExpectedResult.Block, uint32(rs.Payload.Trigger.BlockNumber))
			assert.Equal(t, ocr2keepers.UpkeepIdentifier(test.ExpectedResult.ID.Bytes()), rs.Payload.Upkeep.ID)
			assert.Equal(t, test.ExpectedResult.Eligible, rs.Eligible)
			ext, ok := rs.Extension.(EVMAutomationResultExtension21)
			assert.True(t, ok)
			assert.Equal(t, test.ExpectedResult.FailureReason, ext.FailureReason)
		})
	}
}

func TestPacker_UnpackPerformResult(t *testing.T) {
	tests := []struct {
		Name    string
		RawData string
	}{
		{
			Name:    "unpack success",
			RawData: "0x0000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000a52d",
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			packer, err := newPacker()
			assert.NoError(t, err)
			rs, err := packer.UnpackPerformResult(test.RawData)
			assert.Nil(t, err)
			assert.True(t, rs)
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
	}{
		{
			Name:          "unpack upkeep needed",
			CallbackResp:  []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 46, 20, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 96, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 16, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 16, 10, 11, 21, 31, 41, 15, 16, 17, 18, 19, 13, 14, 12, 13, 14, 15, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8, 120, 111, 101, 122, 90, 54, 44, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			UpkeepNeeded:  true,
			PerformData:   []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 16, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 16, 10, 11, 21, 31, 41, 15, 16, 17, 18, 19, 13, 14, 12, 13, 14, 15, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8, 120, 111, 101, 122, 90, 54, 44, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			FailureReason: UPKEEP_FAILURE_REASON_NONE,
			GasUsed:       big.NewInt(11796),
		},
		{
			Name:          "unpack upkeep not needed",
			CallbackResp:  []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 50, 208, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 96, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 16, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 16, 10, 11, 21, 31, 41, 15, 16, 17, 18, 19, 13, 14, 12, 13, 14, 15, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8, 120, 111, 101, 122, 90, 54, 44, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			UpkeepNeeded:  false,
			PerformData:   []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 16, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 16, 10, 11, 21, 31, 41, 15, 16, 17, 18, 19, 13, 14, 12, 13, 14, 15, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8, 120, 111, 101, 122, 90, 54, 44, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			FailureReason: UPKEEP_FAILURE_REASON_UPKEEP_NOT_NEEDED,
			GasUsed:       big.NewInt(13008),
		},
		{
			Name:         "unpack malformed data",
			CallbackResp: []byte{0, 0, 0, 23, 4, 163, 66, 91, 228, 102, 200, 84, 144, 233, 218, 44, 168, 192, 191, 253, 0, 0, 0, 0, 20, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			UpkeepNeeded: false,
			PerformData:  nil,
			ErrorString:  "abi: improperly encoded boolean value: unpack checkUpkeep return: ",
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			packer, err := newPacker()
			assert.NoError(t, err)

			needed, pd, failureReason, gasUsed, err := packer.UnpackCheckCallbackResult(test.CallbackResp)

			if test.ErrorString != "" {
				assert.EqualError(t, err, test.ErrorString+hexutil.Encode(test.CallbackResp))
			} else {
				assert.Nil(t, err)
			}
			assert.Equal(t, test.UpkeepNeeded, needed)
			assert.Equal(t, test.PerformData, pd)
			assert.Equal(t, test.FailureReason, failureReason)
			assert.Equal(t, test.GasUsed, gasUsed)
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
			packer, err := newPacker()
			assert.NoError(t, err)
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

func TestPacker_PackingTrigger(t *testing.T) {
	tests := []struct {
		name    string
		id      ocr2keepers.UpkeepIdentifier
		trigger triggerWrapper
		encoded []byte
		err     error
	}{
		{
			"happy flow log trigger",
			append([]byte{1}, common.LeftPadBytes([]byte{1}, 15)...),
			triggerWrapper{
				BlockNum:  1,
				BlockHash: common.HexToHash("0x01111111"),
				LogIndex:  1,
				TxHash:    common.HexToHash("0x01111111"),
			},
			func() []byte {
				b, _ := hexutil.Decode("0x0000000000000000000000000000000000000000000000000000000001111111000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000001111111")
				return b
			}(),
			nil,
		},
		{
			"happy flow conditional trigger",
			append([]byte{1}, common.LeftPadBytes([]byte{0}, 15)...),
			triggerWrapper{
				BlockNum:  1,
				BlockHash: common.HexToHash("0x01111111"),
			},
			func() []byte {
				b, _ := hexutil.Decode("0x00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000001111111")
				return b
			}(),
			nil,
		},
		{
			"invalid type",
			append([]byte{1}, common.LeftPadBytes([]byte{8}, 15)...),
			triggerWrapper{
				BlockNum:  1,
				BlockHash: common.HexToHash("0x01111111"),
			},
			[]byte{},
			fmt.Errorf("unknown trigger type: %d", 8),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			packer, err := newPacker()
			assert.NoError(t, err)
			id, ok := big.NewInt(0).SetString(hexutil.Encode(tc.id)[2:], 16)
			assert.True(t, ok)

			encoded, err := packer.PackTrigger(id, tc.trigger)
			if tc.err != nil {
				assert.EqualError(t, err, tc.err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.encoded, encoded)
				decoded, err := packer.UnpackTrigger(id, encoded)
				assert.NoError(t, err)
				assert.Equal(t, tc.trigger.BlockNum, decoded.BlockNum)
			}
		})
	}

	t.Run("unpacking invalid trigger", func(t *testing.T) {
		packer, err := newPacker()
		assert.NoError(t, err)
		_, err = packer.UnpackTrigger(big.NewInt(0), []byte{1, 2, 3})
		assert.Error(t, err)
	})

	t.Run("unpacking unknown type", func(t *testing.T) {
		packer, err := newPacker()
		assert.NoError(t, err)
		uid := append([]byte{1}, common.LeftPadBytes([]byte{8}, 15)...)
		id, ok := big.NewInt(0).SetString(hexutil.Encode(uid)[2:], 16)
		assert.True(t, ok)
		decoded, _ := hexutil.Decode("0x00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000001111111")
		_, err = packer.UnpackTrigger(id, decoded)
		assert.EqualError(t, err, "unknown trigger type: 8")
	})
}

func newPacker() (*evmRegistryPackerV2_1, error) {
	keepersABI, err := abi.JSON(strings.NewReader(iregistry21.IKeeperRegistryMasterABI))
	if err != nil {
		return nil, err
	}
	utilsABI, err := abi.JSON(strings.NewReader(automation21Utils.AutomationUtilsABI))
	if err != nil {
		return nil, err
	}
	return NewEvmRegistryPackerV2_1(keepersABI, utilsABI), nil
}

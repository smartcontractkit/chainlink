package evm

import (
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/automation_utils_2_1"
	iregistry21 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_keeper_registry_master_wrapper_2_1"

	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg"
)

func TestPacker_UnpackTransmitTxInputErrors(t *testing.T) {

	tests := []struct {
		Name    string
		RawData string
	}{
		{
			Name:    "Empty Data",
			RawData: "0x",
		},
		{
			Name:    "Random Data",
			RawData: "0x2f08cfae623a0d96b9beb326c20e322001cbbd344700",
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			packer, err := newPacker()
			assert.NoError(t, err)
			_, err = packer.UnpackTransmitTxInput(hexutil.MustDecode(test.RawData))
			assert.NotNil(t, err)
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
		res     iregistry21.KeeperRegistryBase21LogTriggerConfig
		errored bool
	}{
		{
			"happy flow",
			func() []byte {
				b, _ := hexutil.Decode("0x0000000000000000000000007456fadf415b7c34b1182bd20b0537977e945e3e00000000000000000000000000000000000000000000000000000000000000003d53a39550e04688065827f3bb86584cb007ab9ebca7ebd528e7301c9c31eb5d000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
				return b
			}(),
			iregistry21.KeeperRegistryBase21LogTriggerConfig{
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
			iregistry21.KeeperRegistryBase21LogTriggerConfig{},
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
	utilsABI, err := abi.JSON(strings.NewReader(automation_utils_2_1.AutomationUtilsABI))
	if err != nil {
		return nil, err
	}
	return NewEvmRegistryPackerV2_1(keepersABI, utilsABI), nil
}

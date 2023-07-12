package evm

import (
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg"
	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/automation_utils_2_1"
	iregistry21 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_keeper_registry_master_wrapper_2_1"
)

func TestEVMAutomationEncoder21(t *testing.T) {
	keepersABI, err := abi.JSON(strings.NewReader(iregistry21.IKeeperRegistryMasterABI))
	assert.Nil(t, err)
	utilsABI, err := abi.JSON(strings.NewReader(automation_utils_2_1.AutomationUtilsABI))
	assert.Nil(t, err)
	encoder := EVMAutomationEncoder21{
		packer: NewEvmRegistryPackerV2_1(keepersABI, utilsABI),
	}

	t.Run("encoding an empty list of upkeep results returns a nil byte array", func(t *testing.T) {
		b, err := encoder.Encode()
		assert.Equal(t, ErrEmptyResults, err)
		assert.Equal(t, b, []byte(nil))
	})

	t.Run("successfully encodes and decodes a single upkeep result", func(t *testing.T) {
		upkeepResult := ocr2keepers.CheckResult{
			Payload: ocr2keepers.UpkeepPayload{
				Upkeep: ocr2keepers.ConfiguredUpkeep{
					ID: ocr2keepers.UpkeepIdentifier([]byte("10")),
				},
				Trigger: ocr2keepers.Trigger{
					BlockNumber: 1,
					BlockHash:   common.Bytes2Hex([]byte{1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8}),
				},
			},
			Eligible:     true,
			GasAllocated: 100,
			PerformData:  []byte("data0"),
			Extension: EVMAutomationResultExtension21{
				FastGasWei: big.NewInt(100),
				LinkNative: big.NewInt(100),
			},
		}
		b, err := encoder.Encode(upkeepResult)
		assert.Nil(t, err)
		assert.Len(t, b, 640)

		upkeeps, err := encoder.DecodeReport(b)
		assert.Nil(t, err)
		assert.Len(t, upkeeps, 1)

		// upkeep := upkeeps[0].(EVMAutomationUpkeepResult21)
		// // some fields aren't populated by the decode so we compare field-by-field for those that are populated
		// assert.Equal(t, upkeep.Block, upkeepResult.Block)
		// assert.Equal(t, upkeep.ID, upkeepResult.ID)
		// assert.Equal(t, upkeep.Eligible, upkeepResult.Eligible)
		// assert.Equal(t, upkeep.PerformData, upkeepResult.PerformData)
		// assert.Equal(t, upkeep.FastGasWei, upkeepResult.FastGasWei)
		// assert.Equal(t, upkeep.LinkNative, upkeepResult.LinkNative)
		// assert.Equal(t, upkeep.CheckBlockNumber, upkeepResult.CheckBlockNumber)
		// assert.Equal(t, upkeep.CheckBlockHash, upkeepResult.CheckBlockHash)
	})

}

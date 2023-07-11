package evm

import (
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
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
		b, err := encoder.EncodeReport([]ocr2keepers.UpkeepResult{})
		assert.Nil(t, err)
		assert.Equal(t, b, []byte(nil))
	})

	t.Run("attempting to encode an invalid upkeep result returns an error", func(t *testing.T) {
		b, err := encoder.EncodeReport([]ocr2keepers.UpkeepResult{"data"})
		assert.Error(t, err, "unexpected upkeep result struct")
		assert.Equal(t, b, []byte(nil))
	})

	t.Run("successfully encodes and decodes a single upkeep result", func(t *testing.T) {
		upkeepResult := EVMAutomationUpkeepResult21{
			Block:            1,
			ID:               big.NewInt(10),
			Eligible:         true,
			GasUsed:          big.NewInt(100),
			PerformData:      []byte("data"),
			FastGasWei:       big.NewInt(100),
			LinkNative:       big.NewInt(100),
			CheckBlockNumber: 1,
			CheckBlockHash:   [32]byte{1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8},
			ExecuteGas:       10,
		}
		b, err := encoder.EncodeReport([]ocr2keepers.UpkeepResult{upkeepResult})
		assert.Nil(t, err)
		assert.Len(t, b, 640)

		upkeeps, err := encoder.DecodeReport(b)
		assert.Nil(t, err)
		assert.Len(t, upkeeps, 1)

		upkeep := upkeeps[0].(EVMAutomationUpkeepResult21)
		// some fields aren't populated by the decode so we compare field-by-field for those that are populated
		assert.Equal(t, upkeep.Block, upkeepResult.Block)
		assert.Equal(t, upkeep.ID, upkeepResult.ID)
		assert.Equal(t, upkeep.Eligible, upkeepResult.Eligible)
		assert.Equal(t, upkeep.PerformData, upkeepResult.PerformData)
		assert.Equal(t, upkeep.FastGasWei, upkeepResult.FastGasWei)
		assert.Equal(t, upkeep.LinkNative, upkeepResult.LinkNative)
		assert.Equal(t, upkeep.CheckBlockNumber, upkeepResult.CheckBlockNumber)
		assert.Equal(t, upkeep.CheckBlockHash, upkeepResult.CheckBlockHash)
	})

	t.Run("successfully encodes and decodes multiple upkeep results", func(t *testing.T) {
		n := 5
		results := make([]ocr2keepers.UpkeepResult, n)
		for i := 0; i < n; i++ {
			block := uint32(i + 1)
			results[i] = EVMAutomationUpkeepResult21{
				Block:            block,
				ID:               big.NewInt(int64(block) * 10),
				Eligible:         true,
				GasUsed:          big.NewInt(100),
				PerformData:      []byte(fmt.Sprintf("data-%d", i)),
				FastGasWei:       big.NewInt(100),
				LinkNative:       big.NewInt(100),
				CheckBlockNumber: block,
				CheckBlockHash:   [32]byte{uint8(block), 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8},
				ExecuteGas:       10,
			}
		}

		b, err := encoder.EncodeReport(results)
		assert.Nil(t, err)
		assert.Len(t, b, 1792)

		decoded, err := encoder.DecodeReport(b)
		assert.Nil(t, err)
		assert.Len(t, decoded, len(results))
		for i, dec := range decoded {
			result := dec.(EVMAutomationUpkeepResult21)
			expected := results[i].(EVMAutomationUpkeepResult21)
			assert.Equal(t, result.Block, expected.Block)
			assert.Equal(t, result.ID, expected.ID)
			assert.Equal(t, result.Eligible, expected.Eligible)
			assert.Equal(t, result.PerformData, expected.PerformData)
			assert.Equal(t, result.FastGasWei, expected.FastGasWei)
			assert.Equal(t, result.LinkNative, expected.LinkNative)
			assert.Equal(t, result.CheckBlockNumber, expected.CheckBlockNumber)
			assert.Equal(t, result.CheckBlockHash, expected.CheckBlockHash)
		}
	})
}

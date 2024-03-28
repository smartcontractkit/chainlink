package ccipcommon

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcalc"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	ccipdatamocks "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/mocks"
)

func TestGetMessageIDsAsHexString(t *testing.T) {
	t.Run("base", func(t *testing.T) {
		hashes := make([]cciptypes.Hash, 10)
		for i := range hashes {
			hashes[i] = cciptypes.Hash(common.HexToHash(strconv.Itoa(rand.Intn(100000))))
		}

		msgs := make([]cciptypes.EVM2EVMMessage, len(hashes))
		for i := range msgs {
			msgs[i] = cciptypes.EVM2EVMMessage{MessageID: hashes[i]}
		}

		messageIDs := GetMessageIDsAsHexString(msgs)
		for i := range messageIDs {
			assert.Equal(t, hashes[i].String(), messageIDs[i])
		}
	})

	t.Run("empty", func(t *testing.T) {
		messageIDs := GetMessageIDsAsHexString(nil)
		assert.Empty(t, messageIDs)
	})
}

func TestFlattenUniqueSlice(t *testing.T) {
	testCases := []struct {
		name           string
		inputSlices    [][]int
		expectedOutput []int
	}{
		{name: "empty", inputSlices: nil, expectedOutput: []int{}},
		{name: "empty 2", inputSlices: [][]int{}, expectedOutput: []int{}},
		{name: "single", inputSlices: [][]int{{1, 2, 3, 3, 3, 4}}, expectedOutput: []int{1, 2, 3, 4}},
		{name: "simple", inputSlices: [][]int{{1, 2, 3}, {2, 3, 4}}, expectedOutput: []int{1, 2, 3, 4}},
		{
			name:           "more complex case",
			inputSlices:    [][]int{{1, 3}, {2, 4, 3}, {5, 2, -1, 7, 10}},
			expectedOutput: []int{1, 3, 2, 4, 5, -1, 7, 10},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := FlattenUniqueSlice(tc.inputSlices...)
			assert.Equal(t, tc.expectedOutput, res)
		})
	}
}

func TestGetChainTokens(t *testing.T) {
	var tokens []cciptypes.Address
	for i := 0; i < 6; i++ {
		tokens = append(tokens, ccipcalc.EvmAddrToGeneric(utils.RandomAddress()))
	}

	testCases := []struct {
		name                string
		feeTokens           []cciptypes.Address
		destTokens          [][]cciptypes.Address
		expectedChainTokens []cciptypes.Address
	}{
		{
			name:                "empty",
			feeTokens:           []cciptypes.Address{},
			destTokens:          [][]cciptypes.Address{{}},
			expectedChainTokens: []cciptypes.Address{},
		},
		{
			name:      "single offRamp",
			feeTokens: []cciptypes.Address{tokens[0]},
			destTokens: [][]cciptypes.Address{
				{tokens[1], tokens[2], tokens[3]},
			},
			expectedChainTokens: []cciptypes.Address{tokens[0], tokens[1], tokens[2], tokens[3]},
		},
		{
			name:      "multiple offRamps with distinct tokens",
			feeTokens: []cciptypes.Address{tokens[0]},
			destTokens: [][]cciptypes.Address{
				{tokens[1], tokens[2]},
				{tokens[3], tokens[4]},
				{tokens[5]},
			},
			expectedChainTokens: []cciptypes.Address{tokens[0], tokens[1], tokens[2], tokens[3], tokens[4], tokens[5]},
		},
		{
			name:      "overlapping tokens",
			feeTokens: []cciptypes.Address{tokens[0]},
			destTokens: [][]cciptypes.Address{
				{tokens[0], tokens[1], tokens[2], tokens[3]},
				{tokens[0], tokens[2], tokens[3], tokens[4], tokens[5]},
				{tokens[5]},
			},
			expectedChainTokens: []cciptypes.Address{tokens[0], tokens[1], tokens[2], tokens[3], tokens[4], tokens[5]},
		},
	}

	ctx := testutils.Context(t)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			priceRegistry := ccipdatamocks.NewPriceRegistryReader(t)
			priceRegistry.On("GetFeeTokens", ctx).Return(tc.feeTokens, nil).Once()

			var offRamps []ccipdata.OffRampReader
			for _, destTokens := range tc.destTokens {
				offRamp := ccipdatamocks.NewOffRampReader(t)
				offRamp.On("GetTokens", ctx).Return(cciptypes.OffRampTokens{DestinationTokens: destTokens}, nil).Once()
				offRamps = append(offRamps, offRamp)
			}

			chainTokens, err := GetSortedChainTokens(ctx, offRamps, priceRegistry)
			assert.NoError(t, err)

			sort.Slice(tc.expectedChainTokens, func(i, j int) bool {
				return tc.expectedChainTokens[i] < tc.expectedChainTokens[j]
			})
			assert.Equal(t, tc.expectedChainTokens, chainTokens)
		})
	}
}

func TestGetChainTokensWithBatchLimit(t *testing.T) {
	numTokens := 100
	var tokens []cciptypes.Address
	for i := 0; i < numTokens; i++ {
		tokens = append(tokens, ccipcalc.EvmAddrToGeneric(utils.RandomAddress()))
	}

	expectedTokens := make([]cciptypes.Address, numTokens)
	copy(expectedTokens, tokens)
	sort.Slice(expectedTokens, func(i, j int) bool {
		return expectedTokens[i] < expectedTokens[j]
	})

	testCases := []struct {
		name        string
		batchSize   int
		numOffRamps uint
		expectError bool
	}{
		{
			name:        "default case",
			batchSize:   offRampBatchSizeLimit,
			numOffRamps: 20,
			expectError: false,
		},
		{
			name:        "limit of 0 expects error",
			batchSize:   0,
			numOffRamps: 20,
			expectError: true,
		},
		{
			name:        "low limit of 1 with 1 offRamps",
			batchSize:   1,
			numOffRamps: 1,
			expectError: false,
		},
		{
			name:        "low limit of 1 with many offRamps",
			batchSize:   1,
			numOffRamps: 200,
			expectError: false,
		},
		{
			name:        "high limit of 1000 with few offRamps",
			batchSize:   1000,
			numOffRamps: 20,
			expectError: false,
		},
		{
			name:        "high limit of 1000 with many offRamps",
			batchSize:   1000,
			numOffRamps: 200,
			expectError: false,
		},
	}

	ctx := testutils.Context(t)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			priceRegistry := ccipdatamocks.NewPriceRegistryReader(t)
			priceRegistry.On("GetFeeTokens", ctx).Return(tokens[0:10], nil).Maybe()

			var offRamps []ccipdata.OffRampReader
			for i := 0; i < int(tc.numOffRamps); i++ {
				offRamp := ccipdatamocks.NewOffRampReader(t)
				offRamp.On("GetTokens", ctx).Return(cciptypes.OffRampTokens{DestinationTokens: tokens[i%numTokens:]}, nil).Maybe()
				offRamps = append(offRamps, offRamp)
			}

			chainTokens, err := getSortedChainTokensWithBatchLimit(ctx, offRamps, priceRegistry, tc.batchSize)
			if tc.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, expectedTokens, chainTokens)
		})
	}
}

func TestIsTxRevertError(t *testing.T) {
	testCases := []struct {
		name           string
		inputError     error
		expectedOutput bool
	}{
		{name: "empty", inputError: nil, expectedOutput: false},
		{name: "non-revert error", inputError: fmt.Errorf("nothing"), expectedOutput: false},
		{name: "geth error", inputError: fmt.Errorf("execution reverted"), expectedOutput: true},
		{name: "nethermind error", inputError: fmt.Errorf("VM execution error"), expectedOutput: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedOutput, IsTxRevertError(tc.inputError))
		})
	}
}

package ccipcommon

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcalc"
	ccipdatamocks "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/pricegetter"
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

func TestGetFilteredChainTokens(t *testing.T) {
	const numTokens = 6
	var tokens []cciptypes.Address
	for i := 0; i < numTokens; i++ {
		tokens = append(tokens, ccipcalc.EvmAddrToGeneric(utils.RandomAddress()))
	}

	testCases := []struct {
		name                   string
		feeTokens              []cciptypes.Address
		destTokens             []cciptypes.Address
		expectedChainTokens    []cciptypes.Address
		expectedFilteredTokens []cciptypes.Address
	}{
		{
			name:                   "empty",
			feeTokens:              []cciptypes.Address{},
			destTokens:             []cciptypes.Address{},
			expectedChainTokens:    []cciptypes.Address{},
			expectedFilteredTokens: []cciptypes.Address{},
		},
		{
			name:                   "unique tokens",
			feeTokens:              []cciptypes.Address{tokens[0]},
			destTokens:             []cciptypes.Address{tokens[1], tokens[2], tokens[3]},
			expectedChainTokens:    []cciptypes.Address{tokens[0], tokens[1], tokens[2], tokens[3]},
			expectedFilteredTokens: []cciptypes.Address{tokens[4], tokens[5]},
		},
		{
			name:                   "all tokens",
			feeTokens:              []cciptypes.Address{tokens[0]},
			destTokens:             []cciptypes.Address{tokens[1], tokens[2], tokens[3], tokens[4], tokens[5]},
			expectedChainTokens:    []cciptypes.Address{tokens[0], tokens[1], tokens[2], tokens[3], tokens[4], tokens[5]},
			expectedFilteredTokens: []cciptypes.Address{},
		},
		{
			name:                   "overlapping tokens",
			feeTokens:              []cciptypes.Address{tokens[0]},
			destTokens:             []cciptypes.Address{tokens[1], tokens[2], tokens[5], tokens[3], tokens[0], tokens[2], tokens[3], tokens[4], tokens[5], tokens[5]},
			expectedChainTokens:    []cciptypes.Address{tokens[0], tokens[1], tokens[2], tokens[3], tokens[4], tokens[5]},
			expectedFilteredTokens: []cciptypes.Address{},
		},
		{
			name:                   "unconfigured tokens",
			feeTokens:              []cciptypes.Address{tokens[0]},
			destTokens:             []cciptypes.Address{tokens[0], tokens[1], tokens[2], tokens[3], tokens[0], tokens[2], tokens[3], tokens[4], tokens[5], tokens[5]},
			expectedChainTokens:    []cciptypes.Address{tokens[0], tokens[1], tokens[2], tokens[3], tokens[4]},
			expectedFilteredTokens: []cciptypes.Address{tokens[5]},
		},
	}

	ctx := testutils.Context(t)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			priceRegistry := ccipdatamocks.NewPriceRegistryReader(t)
			priceRegistry.On("GetFeeTokens", ctx).Return(tc.feeTokens, nil).Once()

			priceGet := pricegetter.NewMockPriceGetter(t)
			priceGet.On("FilterConfiguredTokens", mock.Anything, mock.Anything).Return(tc.expectedChainTokens, tc.expectedFilteredTokens, nil)

			offRamp := ccipdatamocks.NewOffRampReader(t)
			offRamp.On("GetTokens", ctx).Return(cciptypes.OffRampTokens{DestinationTokens: tc.destTokens}, nil).Once()

			chainTokens, filteredTokens, err := GetFilteredSortedLaneTokens(ctx, offRamp, priceRegistry, priceGet)
			assert.NoError(t, err)

			sort.Slice(tc.expectedChainTokens, func(i, j int) bool {
				return tc.expectedChainTokens[i] < tc.expectedChainTokens[j]
			})
			assert.Equal(t, tc.expectedChainTokens, chainTokens)
			assert.Equal(t, tc.expectedFilteredTokens, filteredTokens)
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

func TestRetryUntilSuccess(t *testing.T) {
	// Set delays to 0 for tests
	initialDelay := 0 * time.Nanosecond
	maxDelay := 0 * time.Nanosecond

	numAttempts := 5
	numCalls := 0
	// A function that returns success only after numAttempts calls. RetryUntilSuccess will repeatedly call this
	// function until it succeeds.
	fn := func() (int, error) {
		numCalls++
		numAttempts--
		if numAttempts > 0 {
			return numCalls, fmt.Errorf("")
		}
		return numCalls, nil
	}

	// Assert that RetryUntilSuccess returns the expected value when fn returns success on the 5th attempt
	numCalls, err := RetryUntilSuccess(fn, initialDelay, maxDelay)
	assert.Nil(t, err)
	assert.Equal(t, 5, numCalls)

	// Assert that RetryUntilSuccess returns the expected value when fn returns success on the 8th attempt
	numAttempts = 8
	numCalls = 0
	numCalls, err = RetryUntilSuccess(fn, initialDelay, maxDelay)
	assert.Nil(t, err)
	assert.Equal(t, 8, numCalls)
}

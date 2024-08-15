package ocrcommon_test

import (
	"context"
	"database/sql"
	"math/big"
	mrand "math/rand"
	"testing"

	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestArbitrumBlockTranslator_BinarySearch(t *testing.T) {
	t.Parallel()

	blocks := generateDeterministicL2Blocks()
	lggr := logger.TestLogger(t)

	t.Run("returns range of current to nil if target is above current block number", func(t *testing.T) {
		client := evmtest.NewEthClientMock(t)
		abt := ocrcommon.NewArbitrumBlockTranslator(client, lggr)
		ctx := testutils.Context(t)

		var changedInL1Block int64 = 5541

		latestBlock := blocks[1000]
		client.On("HeadByNumber", ctx, (*big.Int)(nil)).Return(latestBlock, nil).Once()

		from, to, err := abt.BinarySearch(ctx, changedInL1Block)
		require.NoError(t, err)

		assert.Equal(t, big.NewInt(1000), from)
		assert.Equal(t, (*big.Int)(nil), to)
	})

	t.Run("returns error if changedInL1Block is less than the lowest possible L1 block on the L2 chain", func(t *testing.T) {
		client := evmtest.NewEthClientMock(t)
		abt := ocrcommon.NewArbitrumBlockTranslator(client, lggr)
		ctx := testutils.Context(t)

		var changedInL1Block int64 = 42

		latestBlock := blocks[1000]
		client.On("HeadByNumber", ctx, (*big.Int)(nil)).Return(latestBlock, nil).Once()

		client.On("HeadByNumber", ctx, mock.AnythingOfType("*big.Int")).Return(func(_ context.Context, num *big.Int) (*evmtypes.Head, error) {
			return blocks[num.Int64()], nil
		})

		_, _, err := abt.BinarySearch(ctx, changedInL1Block)

		assert.EqualError(t, err, "target L1 block number 42 is not represented by any L2 block")
	})

	t.Run("returns error if L1 block number does not exist for any range of L2 blocks", func(t *testing.T) {
		client := evmtest.NewEthClientMock(t)
		abt := ocrcommon.NewArbitrumBlockTranslator(client, lggr)
		ctx := testutils.Context(t)

		var changedInL1Block int64 = 5043

		latestBlock := blocks[1000]
		client.On("HeadByNumber", ctx, (*big.Int)(nil)).Return(latestBlock, nil).Once()

		client.On("HeadByNumber", ctx, mock.AnythingOfType("*big.Int")).Return(func(_ context.Context, num *big.Int) (*evmtypes.Head, error) {
			return blocks[num.Int64()], nil
		})

		_, _, err := abt.BinarySearch(ctx, changedInL1Block)

		assert.EqualError(t, err, "target L1 block number 5043 is not represented by any L2 block")
	})

	t.Run("returns correct range of L2 blocks that encompasses all possible blocks that might contain the given L1 block number", func(t *testing.T) {
		client := evmtest.NewEthClientMock(t)
		abt := ocrcommon.NewArbitrumBlockTranslator(client, lggr)
		ctx := testutils.Context(t)

		var changedInL1Block int64 = 5042

		latestBlock := blocks[1000]
		client.On("HeadByNumber", ctx, (*big.Int)(nil)).Return(latestBlock, nil).Once()

		client.On("HeadByNumber", ctx, mock.AnythingOfType("*big.Int")).Return(func(_ context.Context, num *big.Int) (*evmtypes.Head, error) {
			return blocks[num.Int64()], nil
		})

		from, to, err := abt.BinarySearch(ctx, changedInL1Block)
		require.NoError(t, err)

		assert.Equal(t, big.NewInt(98), from)
		assert.Equal(t, big.NewInt(137), to)
	})

	t.Run("handles edge case where L1 is the smallest possible value", func(t *testing.T) {
		client := evmtest.NewEthClientMock(t)
		abt := ocrcommon.NewArbitrumBlockTranslator(client, lggr)
		ctx := testutils.Context(t)

		var changedInL1Block int64 = 5000

		latestBlock := blocks[1000]
		client.On("HeadByNumber", ctx, (*big.Int)(nil)).Return(latestBlock, nil).Once()

		client.On("HeadByNumber", ctx, mock.AnythingOfType("*big.Int")).Return(func(_ context.Context, num *big.Int) (*evmtypes.Head, error) {
			return blocks[num.Int64()], nil
		})

		from, to, err := abt.BinarySearch(ctx, changedInL1Block)
		require.NoError(t, err)

		assert.Equal(t, big.NewInt(0), from)
		assert.Equal(t, big.NewInt(16), to)
	})

	t.Run("leaves upper bound unbounded where L1 is the largest possible value", func(t *testing.T) {
		client := evmtest.NewEthClientMock(t)
		abt := ocrcommon.NewArbitrumBlockTranslator(client, lggr)
		ctx := testutils.Context(t)

		var changedInL1Block int64 = 5540

		latestBlock := blocks[1000]
		client.On("HeadByNumber", ctx, (*big.Int)(nil)).Return(latestBlock, nil).Once()

		client.On("HeadByNumber", ctx, mock.AnythingOfType("*big.Int")).Return(func(_ context.Context, num *big.Int) (*evmtypes.Head, error) {
			return blocks[num.Int64()], nil
		})

		from, to, err := abt.BinarySearch(ctx, changedInL1Block)
		require.NoError(t, err)

		assert.Equal(t, big.NewInt(986), from)
		assert.Equal(t, (*big.Int)(nil), to)
	})

	t.Run("caches duplicate lookups", func(t *testing.T) {
		client := evmtest.NewEthClientMock(t)
		abt := ocrcommon.NewArbitrumBlockTranslator(client, lggr)
		ctx := testutils.Context(t)

		var changedInL1Block int64 = 5042

		latestBlock := blocks[1000]
		// Latest is never cached
		client.On("HeadByNumber", ctx, (*big.Int)(nil)).Return(latestBlock, nil).Once()

		client.On("HeadByNumber", ctx, mock.AnythingOfType("*big.Int")).Return(func(_ context.Context, num *big.Int) (*evmtypes.Head, error) {
			return blocks[num.Int64()], nil
		})

		// First search, nothing cached (total 21 - bsearch 20)
		from, to, err := abt.BinarySearch(ctx, changedInL1Block)
		require.NoError(t, err)

		assert.Equal(t, big.NewInt(98), from)
		assert.Equal(t, big.NewInt(137), to)

		var changedInL1Block2 int64 = 5351

		// Second search, initial lookup cached + space reduced to [549, 1000] (total 18 - bsearch 18)
		from, to, err = abt.BinarySearch(ctx, changedInL1Block2)
		require.NoError(t, err)

		assert.Equal(t, big.NewInt(670), from)
		assert.Equal(t, big.NewInt(697), to)

		var changedInL1Block3 int64 = 5193

		// Third search, initial lookup cached + space reduced to [323, 500] (total 14 - bsearch 14)
		from, to, err = abt.BinarySearch(ctx, changedInL1Block3)
		require.NoError(t, err)

		assert.Equal(t, big.NewInt(403), from)
		assert.Equal(t, big.NewInt(448), to)
	})

	// TODO: test edge cases - at left edge of range, at right edge
}

func TestArbitrumBlockTranslator_NumberToQueryRange(t *testing.T) {
	t.Parallel()

	lggr := logger.TestLogger(t)

	t.Run("falls back to whole range on error", func(t *testing.T) {
		client := evmtest.NewEthClientMock(t)
		abt := ocrcommon.NewArbitrumBlockTranslator(client, lggr)
		ctx := testutils.Context(t)
		var changedInL1Block uint64 = 5042

		client.On("HeadByNumber", ctx, (*big.Int)(nil)).Return(nil, errors.New("something exploded")).Once()

		from, to := abt.NumberToQueryRange(ctx, changedInL1Block)
		assert.Equal(t, big.NewInt(0), from)
		assert.Equal(t, (*big.Int)(nil), to)
	})

	t.Run("falls back to whole range on missing head", func(t *testing.T) {
		client := evmtest.NewEthClientMock(t)
		abt := ocrcommon.NewArbitrumBlockTranslator(client, lggr)
		ctx := testutils.Context(t)
		var changedInL1Block uint64 = 5042

		client.On("HeadByNumber", ctx, (*big.Int)(nil)).Return(nil, nil).Once()

		from, to := abt.NumberToQueryRange(ctx, changedInL1Block)
		assert.Equal(t, big.NewInt(0), from)
		assert.Equal(t, (*big.Int)(nil), to)
	})
}

func generateDeterministicL2Blocks() (heads []*evmtypes.Head) {
	source := mrand.NewSource(0)
	deterministicRand := mrand.New(source)
	l2max := 1000
	var l1BlockNumber int64 = 5000
	var parentHash common.Hash
	for i := 0; i <= l2max; i++ {
		head := &evmtypes.Head{
			Number:        int64(i),
			L1BlockNumber: sql.NullInt64{Int64: l1BlockNumber, Valid: true},
			Hash:          utils.NewHash(),
			ParentHash:    parentHash,
		}
		parentHash = head.Hash
		heads = append(heads, head)
		if deterministicRand.Intn(10) == 1 { // 10% chance
			// l1 number should jump by "about" 5 but this is variable depending on whether the sequencer got to post, network conditions etc
			l1BlockNumber += int64(deterministicRand.Intn(6) + 4)
		}
	}
	return
}

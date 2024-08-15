package headtracker_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/testutils"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
)

func TestHeads_LatestHead(t *testing.T) {
	t.Parallel()

	heads := headtracker.NewHeads()
	assert.NoError(t, heads.AddHeads(testutils.Head(100), testutils.Head(200), testutils.Head(300)))

	latest := heads.LatestHead()
	require.NotNil(t, latest)
	require.Equal(t, int64(300), latest.Number)

	assert.NoError(t, heads.AddHeads(testutils.Head(250)))
	latest = heads.LatestHead()
	require.NotNil(t, latest)
	require.Equal(t, int64(300), latest.Number)

	assert.NoError(t, heads.AddHeads(testutils.Head(400)))
	latest = heads.LatestHead()
	require.NotNil(t, latest)
	require.Equal(t, int64(400), latest.Number)

	// if heads have the same height, LatestHead prefers most recent
	newerH400 := testutils.Head(400)
	assert.NoError(t, heads.AddHeads(newerH400))
	latest = heads.LatestHead()
	require.NotNil(t, latest)
	require.Equal(t, int64(400), latest.Number)
	require.Equal(t, newerH400.Hash, latest.Hash)
}

func TestHeads_HeadByHash(t *testing.T) {
	t.Parallel()

	var testHeads = []*evmtypes.Head{
		testutils.Head(100),
		testutils.Head(200),
		testutils.Head(300),
	}
	heads := headtracker.NewHeads()
	assert.NoError(t, heads.AddHeads(testHeads...))

	head := heads.HeadByHash(testHeads[1].Hash)
	require.NotNil(t, head)
	require.Equal(t, int64(200), head.Number)

	head = heads.HeadByHash(utils.NewHash())
	require.Nil(t, head)
}

func TestHeads_Count(t *testing.T) {
	t.Parallel()

	heads := headtracker.NewHeads()
	require.Zero(t, heads.Count())

	assert.NoError(t, heads.AddHeads(testutils.Head(100), testutils.Head(200), testutils.Head(300)))
	require.Equal(t, 3, heads.Count())

	assert.NoError(t, heads.AddHeads(testutils.Head(400)))
	require.Equal(t, 4, heads.Count())
}

func TestHeads_AddHeads(t *testing.T) {
	t.Parallel()

	uncleHash := utils.NewHash()
	heads := headtracker.NewHeads()

	var testHeads []*evmtypes.Head
	var parentHash common.Hash
	for i := 1; i < 6; i++ {
		hash := common.BigToHash(big.NewInt(int64(i)))
		h := evmtypes.NewHead(big.NewInt(int64(i)), hash, parentHash, uint64(time.Now().Unix()), ubig.NewI(0))
		testHeads = append(testHeads, &h)
		if i == 3 {
			// uncled block
			h := evmtypes.NewHead(big.NewInt(int64(i)), uncleHash, parentHash, uint64(time.Now().Unix()), ubig.NewI(0))
			testHeads = append(testHeads, &h)
		}
		parentHash = hash
	}

	assert.NoError(t, heads.AddHeads(testHeads...))
	require.Equal(t, 6, heads.Count())
	// Add duplicates (should be ignored)
	assert.NoError(t, heads.AddHeads(testHeads[2:5]...))
	require.Equal(t, 6, heads.Count())

	head := heads.LatestHead()
	require.NotNil(t, head)
	require.Equal(t, 5, int(head.ChainLength()))

	head = heads.HeadByHash(uncleHash)
	require.NotNil(t, head)
	require.Equal(t, 3, int(head.ChainLength()))
	// returns an error, if newHead creates cycle
	t.Run("Returns an error, if newHead create cycle", func(t *testing.T) {
		cycleHead := &evmtypes.Head{
			Hash:       heads.LatestHead().EarliestInChain().ParentHash,
			ParentHash: heads.LatestHead().Hash,
		}
		// 1. try adding in front
		cycleHead.Number = heads.LatestHead().Number + 1
		assert.EqualError(t, heads.AddHeads(cycleHead), "potential cycle detected while adding newHead as parent: expected head number to strictly decrease in 'child -> parent' relation: child(Head{Number: 1, Hash: 0x0000000000000000000000000000000000000000000000000000000000000001, ParentHash: 0x0000000000000000000000000000000000000000000000000000000000000000}), parent(Head{Number: 6, Hash: 0x0000000000000000000000000000000000000000000000000000000000000000, ParentHash: 0x0000000000000000000000000000000000000000000000000000000000000005})")
		// 2. try adding to back
		cycleHead.Number = heads.LatestHead().EarliestInChain().Number - 1
		assert.EqualError(t, heads.AddHeads(cycleHead), "potential cycle detected while adding newHead as child: expected head number to strictly decrease in 'child -> parent' relation: child(Head{Number: 0, Hash: 0x0000000000000000000000000000000000000000000000000000000000000000, ParentHash: 0x0000000000000000000000000000000000000000000000000000000000000005}), parent(Head{Number: 5, Hash: 0x0000000000000000000000000000000000000000000000000000000000000005, ParentHash: 0x0000000000000000000000000000000000000000000000000000000000000004})")
		// 3. try adding to back with reference to self
		cycleHead = &evmtypes.Head{
			Number:     1000,
			Hash:       common.BigToHash(big.NewInt(1000)),
			ParentHash: common.BigToHash(big.NewInt(1000)),
		}
		assert.EqualError(t, heads.AddHeads(cycleHead), "cycle detected: newHeads reference itself newHead(Head{Number: 1000, Hash: 0x00000000000000000000000000000000000000000000000000000000000003e8, ParentHash: 0x00000000000000000000000000000000000000000000000000000000000003e8})")
	})
}

func TestHeads_MarkFinalized(t *testing.T) {
	t.Parallel()

	heads := headtracker.NewHeads()

	// create chain
	// H0 <- H1 <- H2 <- H3 <- H4 <- H5 - Canonical
	//   \      \
	//    H1Uncle       H2Uncle
	//
	newHead := func(num int, parent common.Hash) *evmtypes.Head {
		h := evmtypes.NewHead(big.NewInt(int64(num)), utils.NewHash(), parent, uint64(time.Now().Unix()), ubig.NewI(0))
		return &h
	}
	h0 := newHead(0, utils.NewHash())
	h1 := newHead(1, h0.Hash)
	h1Uncle := newHead(1, h0.Hash)
	h2 := newHead(2, h1.Hash)
	h3 := newHead(3, h2.Hash)
	h4 := newHead(4, h3.Hash)
	h5 := newHead(5, h4.Hash)
	h2Uncle := newHead(2, h1.Hash)

	assert.NoError(t, heads.AddHeads(h0, h1, h1Uncle, h2, h2Uncle, h3, h4, h5))
	// mark h3 and all ancestors as finalized
	require.True(t, heads.MarkFinalized(h3.Hash, h1.BlockNumber()), "expected MarkFinalized succeed")

	// h0 is too old. It should not be available directly or through its children
	assert.Nil(t, heads.HeadByHash(h0.Hash))
	assert.Nil(t, heads.HeadByHash(h1.Hash).Parent.Load())
	assert.Nil(t, heads.HeadByHash(h1Uncle.Hash).Parent.Load())
	assert.Nil(t, heads.HeadByHash(h2Uncle.Hash).Parent.Load().Parent.Load())

	require.False(t, heads.MarkFinalized(utils.NewHash(), 0), "expected false if finalized hash was not found in existing LatestHead chain")

	ensureProperFinalization := func(t *testing.T) {
		t.Helper()
		for _, head := range []*evmtypes.Head{h5, h4} {
			require.False(t, heads.HeadByHash(head.Hash).IsFinalized.Load(), "expected h4-h5 not to be finalized", head.BlockNumber())
		}
		for _, head := range []*evmtypes.Head{h3, h2, h1} {
			require.True(t, heads.HeadByHash(head.Hash).IsFinalized.Load(), "expected h3 and all ancestors to be finalized", head.BlockNumber())
		}
		require.False(t, heads.HeadByHash(h2Uncle.Hash).IsFinalized.Load(), "expected uncle block not to be marked as finalized")
	}
	t.Run("blocks were correctly marked as finalized", ensureProperFinalization)
	assert.NoError(t, heads.AddHeads(h0, h1, h2, h2Uncle, h3, h4, h5))
	t.Run("blocks remain finalized after re adding them to the Heads", ensureProperFinalization)

	// ensure that IsFinalized is propagated, when older blocks are added
	// 1. remove all blocks older than 3
	heads.MarkFinalized(h3.Hash, 3)
	// 2. ensure that h2 and h1 are no longer present
	assert.Nil(t, heads.HeadByHash(h2.Hash))
	assert.Nil(t, heads.HeadByHash(h1.Hash))
	// 3. add blocks back, starting from older
	assert.NoError(t, heads.AddHeads(h1))
	assert.False(t, heads.HeadByHash(h1.Hash).IsFinalized.Load(), "expected h1 to not be finalized as it was not explicitly marked and there no path to h3")
	assert.NoError(t, heads.AddHeads(h2))
	// 4. now h2 and h1 must be marked as finalized
	assert.True(t, heads.HeadByHash(h1.Hash).IsFinalized.Load())
	assert.True(t, heads.HeadByHash(h2.Hash).IsFinalized.Load())
}

func BenchmarkEarliestHeadInChain(b *testing.B) {
	const latestBlockNum = 200_000
	blocks := NewBlocks(b, latestBlockNum+1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		latest := blocks.Head(latestBlockNum)
		earliest := latest.EarliestHeadInChain()
		// perform sanity check
		assert.NotEqual(b, latest.BlockNumber(), earliest.BlockNumber())
		assert.NotEqual(b, latest.BlockHash(), earliest.BlockHash())
	}
}

// BenchmarkSimulated_Backfill - benchmarks AddHeads & MarkFinalized as if it was performed by HeadTracker's backfill
func BenchmarkHeads_SimulatedBackfill(b *testing.B) {
	makeHash := func(n int64) common.Hash {
		return common.BigToHash(big.NewInt(n))
	}
	makeHead := func(n int64) *evmtypes.Head {
		return &evmtypes.Head{Number: n, Hash: makeHash(n), ParentHash: makeHash(n - 1)}
	}

	const finalityDepth = 16_000 // observed value on Arbitrum
	// populate with initial values
	heads := headtracker.NewHeads()
	for i := int64(1); i <= finalityDepth; i++ {
		assert.NoError(b, heads.AddHeads(makeHead(i)))
	}
	heads.MarkFinalized(makeHash(1), 1)
	// focus benchmark on processing of a new latest block
	b.ResetTimer()
	for i := int64(1); i <= int64(b.N); i++ {
		assert.NoError(b, heads.AddHeads(makeHead(finalityDepth+i)))
		heads.MarkFinalized(makeHash(i), i)
	}
}

package headtracker_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
)

func TestHeads_LatestHead(t *testing.T) {
	t.Parallel()

	heads := headtracker.NewHeads()
	heads.AddHeads(cltest.Head(100), cltest.Head(200), cltest.Head(300))

	latest := heads.LatestHead()
	require.NotNil(t, latest)
	require.Equal(t, int64(300), latest.Number)

	heads.AddHeads(cltest.Head(250))
	latest = heads.LatestHead()
	require.NotNil(t, latest)
	require.Equal(t, int64(300), latest.Number)

	heads.AddHeads(cltest.Head(400))
	latest = heads.LatestHead()
	require.NotNil(t, latest)
	require.Equal(t, int64(400), latest.Number)
}

func TestHeads_HeadByHash(t *testing.T) {
	t.Parallel()

	var testHeads = []*evmtypes.Head{
		cltest.Head(100),
		cltest.Head(200),
		cltest.Head(300),
	}
	heads := headtracker.NewHeads()
	heads.AddHeads(testHeads...)

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

	heads.AddHeads(cltest.Head(100), cltest.Head(200), cltest.Head(300))
	require.Equal(t, 3, heads.Count())

	heads.AddHeads(cltest.Head(400))
	require.Equal(t, 4, heads.Count())
}

func TestHeads_AddHeads(t *testing.T) {
	t.Parallel()

	uncleHash := utils.NewHash()
	heads := headtracker.NewHeads()

	var testHeads []*evmtypes.Head
	var parentHash common.Hash
	for i := 0; i < 5; i++ {
		hash := utils.NewHash()
		h := evmtypes.NewHead(big.NewInt(int64(i)), hash, parentHash, uint64(time.Now().Unix()), ubig.NewI(0))
		testHeads = append(testHeads, &h)
		if i == 2 {
			// uncled block
			h := evmtypes.NewHead(big.NewInt(int64(i)), uncleHash, parentHash, uint64(time.Now().Unix()), ubig.NewI(0))
			testHeads = append(testHeads, &h)
		}
		parentHash = hash
	}

	heads.AddHeads(testHeads...)
	require.Equal(t, 6, heads.Count())
	// Add duplicates (should be ignored)
	heads.AddHeads(testHeads[2:5]...)
	require.Equal(t, 6, heads.Count())

	head := heads.LatestHead()
	require.NotNil(t, head)
	require.Equal(t, 5, int(head.ChainLength()))

	head = heads.HeadByHash(uncleHash)
	require.NotNil(t, head)
	require.Equal(t, 3, int(head.ChainLength()))
}

func TestHeads_MarkFinalized(t *testing.T) {
	t.Parallel()

	heads := headtracker.NewHeads()

	// create chain
	// H0 <- H1 <- H2 <- H3 <- H4 <- H5
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

	allHeads := []*evmtypes.Head{h0, h1, h1Uncle, h2, h2Uncle, h3, h4, h5}
	heads.AddHeads(allHeads...)
	// mark h3 and all ancestors as finalized
	require.True(t, heads.MarkFinalized(h3.Hash, h1.BlockNumber()), "expected MarkFinalized succeed")

	// original heads remain unchanged
	for _, h := range allHeads {
		assert.False(t, h.IsFinalized, "expected original heads to remain unfinalized")
	}

	// h0 is too old. It should not be available directly or through its children
	assert.Nil(t, heads.HeadByHash(h0.Hash))
	assert.Nil(t, heads.HeadByHash(h1.Hash).Parent)
	assert.Nil(t, heads.HeadByHash(h1Uncle.Hash).Parent)
	assert.Nil(t, heads.HeadByHash(h2Uncle.Hash).Parent.Parent)

	require.False(t, heads.MarkFinalized(utils.NewHash(), 0), "expected false if finalized hash was not found in existing LatestHead chain")

	ensureProperFinalization := func(t *testing.T) {
		t.Helper()
		for _, head := range []*evmtypes.Head{h5, h4} {
			require.False(t, heads.HeadByHash(head.Hash).IsFinalized, "expected h4-h5 not to be finalized", head.BlockNumber())
		}
		for _, head := range []*evmtypes.Head{h3, h2, h1} {
			require.True(t, heads.HeadByHash(head.Hash).IsFinalized, "expected h3 and all ancestors to be finalized", head.BlockNumber())
		}
		require.False(t, heads.HeadByHash(h2Uncle.Hash).IsFinalized, "expected uncle block not to be marked as finalized")
	}
	t.Run("blocks were correctly marked as finalized", ensureProperFinalization)
	heads.AddHeads(h0, h1, h2, h2Uncle, h3, h4, h5)
	t.Run("blocks remain finalized after re adding them to the Heads", ensureProperFinalization)
}

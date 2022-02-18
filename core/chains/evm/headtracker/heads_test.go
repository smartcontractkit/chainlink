package headtracker_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/chains/evm/headtracker"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func TestHeads_LatestHead(t *testing.T) {
	t.Parallel()

	heads := headtracker.NewHeads()
	heads.AddHeads(3, cltest.Head(100), cltest.Head(200), cltest.Head(300))

	latest := heads.LatestHead()
	require.NotNil(t, latest)
	require.Equal(t, int64(300), latest.Number)

	heads.AddHeads(3, cltest.Head(250))
	latest = heads.LatestHead()
	require.NotNil(t, latest)
	require.Equal(t, int64(300), latest.Number)

	heads.AddHeads(3, cltest.Head(400))
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
	heads.AddHeads(3, testHeads...)

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

	heads.AddHeads(3, cltest.Head(100), cltest.Head(200), cltest.Head(300))
	require.Equal(t, 3, heads.Count())

	heads.AddHeads(1, cltest.Head(400))
	require.Equal(t, 1, heads.Count())
}

func TestHeads_AddHeads(t *testing.T) {
	t.Parallel()

	uncleHash := utils.NewHash()
	heads := headtracker.NewHeads()

	var testHeads []*evmtypes.Head
	var parentHash common.Hash
	for i := 0; i < 5; i++ {
		hash := utils.NewHash()
		h := evmtypes.NewHead(big.NewInt(int64(i)), hash, parentHash, uint64(time.Now().Unix()), utils.NewBigI(0))
		testHeads = append(testHeads, &h)
		if i == 2 {
			// uncled block
			h := evmtypes.NewHead(big.NewInt(int64(i)), uncleHash, parentHash, uint64(time.Now().Unix()), utils.NewBigI(0))
			testHeads = append(testHeads, &h)
		}
		parentHash = hash
	}

	heads.AddHeads(6, testHeads...)
	// Add duplicates (should be ignored)
	heads.AddHeads(6, testHeads[2:5]...)
	require.Equal(t, 6, heads.Count())

	head := heads.LatestHead()
	require.NotNil(t, head)
	require.Equal(t, 5, int(head.ChainLength()))

	head = heads.HeadByHash(uncleHash)
	require.NotNil(t, head)
	require.Equal(t, 3, int(head.ChainLength()))

	// Adding beyond the limit truncates
	heads.AddHeads(2, testHeads...)
	require.Equal(t, 2, heads.Count())
	head = heads.LatestHead()
	require.NotNil(t, head)
	require.Equal(t, 2, int(head.ChainLength()))
}

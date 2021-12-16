package headtracker_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/services/headtracker"
	htmocks "github.com/smartcontractkit/chainlink/core/services/headtracker/mocks"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func configureSaver(t *testing.T) (headtracker.HeadSaver, headtracker.ORM) {
	db := pgtest.NewSqlxDB(t)
	lggr := logger.TestLogger(t)
	cfg := cltest.NewTestGeneralConfig(t)
	htCfg := new(htmocks.Config)
	htCfg.Test(t)
	htCfg.On("EvmHeadTrackerHistoryDepth").Return(uint32(6))
	htCfg.On("EvmFinalityDepth").Return(uint32(1))
	orm := headtracker.NewORM(db, lggr, cfg, cltest.FixtureChainID)
	saver := headtracker.NewHeadSaver(lggr, orm, htCfg)
	return saver, orm
}

func TestHeadSaver_Save(t *testing.T) {
	t.Parallel()

	saver, _ := configureSaver(t)

	head := cltest.Head(1)
	err := saver.Save(context.TODO(), head)
	require.NoError(t, err)

	latest, err := saver.LatestHeadFromDB(context.TODO())
	require.NoError(t, err)
	require.Equal(t, int64(1), latest.Number)

	latest = saver.LatestChain()
	require.NotNil(t, latest)
	require.Equal(t, int64(1), latest.Number)

	latest = saver.Chain(head.Hash)
	require.NotNil(t, latest)
	require.Equal(t, int64(1), latest.Number)
}

func TestHeadSaver_LoadFromDB(t *testing.T) {
	t.Parallel()

	saver, orm := configureSaver(t)

	for i := 0; i < 5; i++ {
		err := orm.IdempotentInsertHead(context.TODO(), cltest.Head(i))
		require.NoError(t, err)
	}

	latestHead, err := saver.LoadFromDB(context.TODO())
	require.NoError(t, err)
	require.NotNil(t, latestHead)
	require.Equal(t, int64(4), latestHead.Number)

	latestChain := saver.LatestChain()
	require.NotNil(t, latestChain)
	require.Equal(t, int64(4), latestChain.Number)
}

func TestHeadSaver_HeadsProcessing(t *testing.T) {
	uncleHash := utils.NewHash()
	saver, _ := configureSaver(t)

	var heads []*eth.Head
	var parentHash common.Hash
	for i := 0; i < 5; i++ {
		hash := utils.NewHash()
		h := eth.NewHead(big.NewInt(int64(i)), hash, parentHash, uint64(time.Now().Unix()), utils.NewBigI(0))
		heads = append(heads, &h)
		if i == 2 {
			// uncled block
			h := eth.NewHead(big.NewInt(int64(i)), uncleHash, parentHash, uint64(time.Now().Unix()), utils.NewBigI(0))
			heads = append(heads, &h)
		}
		parentHash = hash
	}

	// adding duplicates
	heads = append(heads, heads[2:5]...)

	for _, head := range heads {
		err := saver.Save(context.TODO(), head)
		require.NoError(t, err)
	}

	ch := saver.LatestChain()
	require.Equal(t, 6, len(headtracker.Heads(saver)))
	require.NotNil(t, ch)
	require.Equal(t, 5, int(ch.ChainLength()))

	ch = saver.Chain(uncleHash)
	require.NotNil(t, ch)
	require.Equal(t, 3, int(ch.ChainLength()))

	// Adding beyond the limit truncates
	headtracker.AddHeads(saver, heads, 2)
	require.Equal(t, 2, len(headtracker.Heads(saver)))
	ch = saver.LatestChain()
	require.NotNil(t, ch)
	require.Equal(t, 2, int(ch.ChainLength()))
}

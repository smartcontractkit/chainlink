package logpoller

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/chains/evm/mocks"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
)

func TestLogPoller(t *testing.T) {
	ec := new(mocks.Client)
	db := pgtest.NewSqlxDB(t)
	lggr := logger.TestLogger(t)
	orm := NewORM(big.NewInt(42), db, lggr, pgtest.NewPGCfg(true))
	lp := NewLogPoller(orm, ec, lggr)

	memdb := rawdb.NewMemoryDatabase()
	_ = (&core.Genesis{BaseFee: big.NewInt(params.InitialBaseFee)}).MustCommit(memdb)
	// Initialize a fresh chain with only a genesis block
	blockchain, _ := core.NewBlockChain(memdb, nil, params.AllEthashProtocolChanges, ethash.NewFaker(), vm.Config{}, nil, nil)
	blocks, _ := core.GenerateChain(params.TestChainConfig, blockchain.CurrentBlock(), ethash.NewFaker(), memdb, 10, func(i int, b *core.BlockGen) {
	})
	_, err := blockchain.InsertChain(blocks)
	require.NoError(t, err)
	ec.On("BlockByNumber", mock.Anything, mock.Anything).Return(func(ctx context.Context, n *big.Int) (*types.Block) {
		return blockchain.GetBlockByNumber(n.Uint64())
	}, func(ctx context.Context, n *big.Int) (error) {
		return nil
	})
	b, err := ec.BlockByNumber(context.Background(), big.NewInt(1))
	t.Log(b.Number().Int64(), b.Hash(), b.ParentHash())
	//ec.On("BlockByNumber", mock.Anything, mock.Anything).Return(types.NewBlockWithHeader(&types.Header{
	//	ParentHash:  common.HexToHash("0x1"),
	//	Number: big.NewInt(1),
	//}))
	//lp.pollAndSaveLogs(context.Background(), 1)
	//ec.AssertExpectations(t)
}

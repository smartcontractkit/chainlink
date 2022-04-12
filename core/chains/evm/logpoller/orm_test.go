package logpoller

import (
	"database/sql"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func TestORM(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	lggr := logger.TestLogger(t)
	require.NoError(t, utils.JustError(db.Exec(`SET CONSTRAINTS log_poller_blocks_evm_chain_id_fkey DEFERRED`)))
	require.NoError(t, utils.JustError(db.Exec(`SET CONSTRAINTS logs_evm_chain_id_fkey DEFERRED`)))
	o1 := NewORM(big.NewInt(137), db, lggr, pgtest.NewPGCfg(true))
	o2 := NewORM(big.NewInt(138), db, lggr, pgtest.NewPGCfg(true))

	// Insert and read back a block.
	require.NoError(t, o1.InsertBlock(common.HexToHash("0x1234"), 10))
	b, err := o1.SelectBlockByHash(common.HexToHash("0x1234"))
	require.NoError(t, err)
	assert.Equal(t, b.BlockNumber, int64(10))
	assert.Equal(t, b.BlockHash.Bytes(), common.HexToHash("0x1234").Bytes())
	assert.Equal(t, b.EvmChainId.String(), "137")

	// Insert another block from a different chain
	require.NoError(t, o2.InsertBlock(common.HexToHash("0x1234"), 11))
	b2, err := o2.SelectBlockByHash(common.HexToHash("0x1234"))
	require.NoError(t, err)
	assert.Equal(t, b2.BlockNumber, int64(11))
	assert.Equal(t, b2.BlockHash.Bytes(), common.HexToHash("0x1234").Bytes())
	assert.Equal(t, b2.EvmChainId.String(), "138")

	latest, err := o1.SelectLatestBlock()
	require.NoError(t, err)
	assert.Equal(t, int64(10), latest.BlockNumber)

	latest, err = o2.SelectLatestBlock()
	require.NoError(t, err)
	assert.Equal(t, int64(11), latest.BlockNumber)

	// Delete a block
	require.NoError(t, o1.DeleteRangeBlocks(10, 10))
	_, err = o1.SelectBlockByHash(common.HexToHash("0x1234"))
	require.Error(t, err)
	assert.True(t, errors.Is(err, sql.ErrNoRows))

	// Delete block from another chain.
	require.NoError(t, o2.DeleteRangeBlocks(11, 11))
	_, err = o2.SelectBlockByHash(common.HexToHash("0x1234"))
	require.Error(t, err)
	assert.True(t, errors.Is(err, sql.ErrNoRows))

	// Should be able to insert and read back a log.
	topic := common.HexToHash("0x1599")
	require.NoError(t, o1.InsertLogs([]Log{
		{
			EvmChainId:  utils.NewBigI(137),
			LogIndex:    1,
			BlockHash:   common.HexToHash("0x1234"),
			BlockNumber: int64(10),
			EventSig:    topic[:],
			Topics:      [][]byte{topic[:]},
			Address:     common.HexToAddress("0x1234"),
			TxHash:      common.HexToHash("0x1888"),
			Data:        []byte("hello"),
		},
	}))
	logs, err := o1.selectLogsByBlockRange(10, 10)
	require.NoError(t, err)
	require.Equal(t, 1, len(logs))
	assert.Equal(t, []byte("hello"), logs[0].Data)

	logs, err = o1.SelectLogsByBlockRangeFilter(10, 10, common.HexToAddress("0x1234"), topic[:])
	require.NoError(t, err)
	require.Equal(t, 1, len(logs))
}

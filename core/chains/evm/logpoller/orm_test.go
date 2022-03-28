package logpoller

import (
	"database/sql"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func TestORM(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	lggr := logger.TestLogger(t)
	require.NoError(t, utils.JustError(db.Exec(`SET CONSTRAINTS log_poller_blocks_evm_chain_id_fkey DEFERRED`)))
	require.NoError(t, utils.JustError(db.Exec(`SET CONSTRAINTS logs_evm_chain_id_fkey DEFERRED`)))
	o := NewORM(big.NewInt(137), db, lggr, pgtest.NewPGCfg(true))

	// Insert and read back a block.
	require.NoError(t, o.InsertBlock(common.HexToHash("0x1234"), 10))
	b, err := o.SelectBlockByHash(common.HexToHash("0x1234"))
	require.NoError(t, err)
	assert.Equal(t, b.BlockNumber, int64(10))
	assert.Equal(t, b.BlockHash.Bytes(), common.HexToHash("0x1234").Bytes())
	assert.Equal(t, b.EvmChainId.String(), "137")

	// Delete a block
	require.NoError(t, o.DeleteRangeBlocks(10, 10))
	_, err = o.SelectBlockByHash(common.HexToHash("0x1234"))
	require.Error(t, err)
	t.Log(errors.Is(err, sql.ErrNoRows))

	// Should be able to insert and read back a log.
	topic := common.HexToHash("0x1599")
	require.NoError(t, o.InsertLogs([]Log{
		{
			EvmChainId:  utils.NewBigI(137),
			LogIndex:    1,
			BlockHash:   common.HexToHash("0x1234"),
			BlockNumber: int64(10),
			Topics:      [][]byte{topic[:]},
			Address:     common.HexToAddress("0x1234"),
			TxHash:      common.HexToHash("0x1888"),
			Data:        []byte("hello"),
		},
	}))
	logs, err := o.SelectLogsByBlockRange(10, 10)
	require.NoError(t, err)
	require.Equal(t, 1, len(logs))
	assert.Equal(t, []byte("hello"), logs[0].Data)
}

func TestCanonicalQuery(t *testing.T) {
	//db := pgtest.NewSqlxDB(t)
	lggr := logger.TestLogger(t)
	//require.NoError(t, utils.JustError(db.Exec(`SET CONSTRAINTS log_poller_blocks_evm_chain_id_fkey DEFERRED`)))
	//require.NoError(t, utils.JustError(db.Exec(`SET CONSTRAINTS logs_evm_chain_id_fkey DEFERRED`)))
	_, db := heavyweight.FullTestDB(t, "logs", true, false)
	chainID := big.NewInt(137)
	_, err := db.Exec(`INSERT INTO evm_chains (id, created_at, updated_at) VALUES ($1, NOW(), NOW())`, utils.NewBig(chainID))
	require.NoError(t, err)
	o := NewORM(big.NewInt(137), db, lggr, pgtest.NewPGCfg(true))
	topic := common.HexToHash("0x1599")

	// Block 1 and 2
	require.NoError(t, o.InsertLogs([]Log{
		{
			EvmChainId:  utils.NewBigI(137),
			LogIndex:    1,
			BlockHash:   common.HexToHash("0x1"),
			BlockNumber: int64(1),
			Topics:      [][]byte{topic[:]},
			Address:     common.HexToAddress("0x1234"),
			TxHash:      common.HexToHash("0x1234"),
			Data:        []byte("hello"),
		},
		{
			EvmChainId:  utils.NewBigI(137),
			LogIndex:    2,
			BlockHash:   common.HexToHash("0x1"),
			BlockNumber: int64(1),
			Topics:      [][]byte{topic[:]},
			Address:     common.HexToAddress("0x1234"),
			TxHash:      common.HexToHash("0x1234"),
			Data:        []byte("hello"),
		},
		{
			EvmChainId:  utils.NewBigI(137),
			LogIndex:    1,
			BlockHash:   common.HexToHash("0x2"),
			BlockNumber: int64(2),
			Topics:      [][]byte{topic[:]},
			Address:     common.HexToAddress("0x1234"),
			TxHash:      common.HexToHash("0x1234"),
			Data:        []byte("hello"),
		},
		{
			EvmChainId:  utils.NewBigI(137),
			LogIndex:    2,
			BlockHash:   common.HexToHash("0x2"),
			BlockNumber: int64(2),
			Topics:      [][]byte{topic[:]},
			Address:     common.HexToAddress("0x1234"),
			TxHash:      common.HexToHash("0x1234"),
			Data:        []byte("hello"),
		},
	}))

	// Block 1' and 2'
	require.NoError(t, o.InsertLogs([]Log{
		{
			EvmChainId:  utils.NewBigI(137),
			LogIndex:    3,
			BlockHash:   common.HexToHash("0x3"),
			BlockNumber: int64(1),
			Topics:      [][]byte{topic[:]},
			Address:     common.HexToAddress("0x1234"),
			TxHash:      common.HexToHash("0x1234"),
			Data:        []byte("hello"),
		},
		{
			EvmChainId:  utils.NewBigI(137),
			LogIndex:    4,
			BlockHash:   common.HexToHash("0x4"),
			BlockNumber: int64(2),
			Topics:      [][]byte{topic[:]},
			Address:     common.HexToAddress("0x1234"),
			TxHash:      common.HexToHash("0x1234"),
			Data:        []byte("hello"),
		},
	}))

	lgs, err := o.SelectCanonicalLogsByBlockRange(1, 2)
	require.NoError(t, err)
	// We expect only logs from block hash 0x3 and 0x4.
	require.Equal(t, 2, len(lgs))
	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000003", lgs[0].BlockHash.String())
	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000004", lgs[1].BlockHash.String())
}

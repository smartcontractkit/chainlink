package terratxm

import (
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestORM(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	lggr := logger.TestLogger(t)
	o := NewORM(db, lggr, pgtest.NewPGCfg(true))

	// Create
	mid, err := o.InsertMsg("0x123", []byte("hello"))
	require.NoError(t, err)
	assert.NotEqual(t, 0, int(mid))

	// Read
	unstarted, err := o.SelectMsgsWithState(Unstarted)
	require.NoError(t, err)
	require.Equal(t, 1, len(unstarted))
	assert.Equal(t, "hello", string(unstarted[0].Msg))
	t.Log(unstarted[0].UpdatedAt, unstarted[0].CreatedAt)

	// Update
	txHash := "123"
	err = o.UpdateMsgsWithState([]int64{mid}, Broadcasted, &txHash)
	require.NoError(t, err)
	broadcasted, err := o.SelectMsgsWithState(Broadcasted)
	require.NoError(t, err)
	require.Equal(t, 1, len(broadcasted))
	assert.Equal(t, broadcasted[0].Msg, unstarted[0].Msg)
	require.NotNil(t, broadcasted[0].TxHash)
	assert.Equal(t, *broadcasted[0].TxHash, txHash)

	err = o.UpdateMsgsWithState([]int64{mid}, Confirmed, nil)
	require.NoError(t, err)
	confirmed, err := o.SelectMsgsWithState(Confirmed)
	require.NoError(t, err)
	require.Equal(t, 1, len(confirmed))
}

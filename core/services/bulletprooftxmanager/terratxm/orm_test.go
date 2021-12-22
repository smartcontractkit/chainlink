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
	err = o.UpdateMsgsWithState([]int64{mid}, Completed)
	require.NoError(t, err)
	completed, err := o.SelectMsgsWithState(Completed)
	require.NoError(t, err)
	require.Equal(t, 1, len(completed))
	assert.Equal(t, completed[0].Msg, unstarted[0].Msg)
}

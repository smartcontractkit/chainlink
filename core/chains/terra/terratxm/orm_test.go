package terratxm_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/chains/terra"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"

	. "github.com/smartcontractkit/chainlink-terra/pkg/terra/db"
	. "github.com/smartcontractkit/chainlink/core/chains/terra/terratxm"
)

func TestORM(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	lggr := logger.TestLogger(t)
	logCfg := pgtest.NewPGCfg(true)
	chainID := fmt.Sprintf("Chainlinktest-%d", rand.Int31n(999999))
	_, err := terra.NewORM(db, lggr, logCfg).CreateChain(chainID, ChainCfg{})
	require.NoError(t, err)
	o := NewORM(chainID, db, lggr, logCfg)

	// Create
	mid, err := o.InsertMsg("0x123", []byte("hello"))
	require.NoError(t, err)
	assert.NotEqual(t, 0, int(mid))

	// Read
	unstarted, err := o.SelectMsgsWithState(Unstarted)
	require.NoError(t, err)
	require.Equal(t, 1, len(unstarted))
	assert.Equal(t, "hello", string(unstarted[0].Raw))
	assert.Equal(t, chainID, unstarted[0].ChainID)
	t.Log(unstarted[0].UpdatedAt, unstarted[0].CreatedAt)

	// Update
	txHash := "123"
	err = o.UpdateMsgsWithState([]int64{mid}, Broadcasted, &txHash)
	require.NoError(t, err)
	broadcasted, err := o.SelectMsgsWithState(Broadcasted)
	require.NoError(t, err)
	require.Equal(t, 1, len(broadcasted))
	assert.Equal(t, broadcasted[0].Raw, unstarted[0].Raw)
	require.NotNil(t, broadcasted[0].TxHash)
	assert.Equal(t, *broadcasted[0].TxHash, txHash)
	assert.Equal(t, chainID, broadcasted[0].ChainID)

	err = o.UpdateMsgsWithState([]int64{mid}, Confirmed, nil)
	require.NoError(t, err)
	confirmed, err := o.SelectMsgsWithState(Confirmed)
	require.NoError(t, err)
	require.Equal(t, 1, len(confirmed))
}

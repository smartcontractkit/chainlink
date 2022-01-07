package terratxm_test

import (
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/utils"

	uuid "github.com/satori/go.uuid"
	tcmocks "github.com/smartcontractkit/chainlink-terra/pkg/terra/client/mocks"
	"github.com/smartcontractkit/chainlink/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager/terratxm"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/stretchr/testify/require"
)

func TestTxmStartStop(t *testing.T) {
	cfg, db := heavyweight.FullTestDB(t, "terra_txm", true, false)
	lggr := logger.TestLogger(t)
	eb := pg.NewEventBroadcaster(cfg.DatabaseURL(), 0, 0, lggr, uuid.NewV4())
	require.NoError(t, eb.Start())
	t.Cleanup(func() { require.NoError(t, eb.Close()) })
	ks := keystore.New(db, utils.FastScryptParams, lggr, pgtest.NewPGCfg(true))
	tc := new(tcmocks.ReaderWriter)
	txm := terratxm.NewTxm(db, tc, ks.Terra(), lggr, pgtest.NewPGCfg(true), nil, time.Second)
	require.NoError(t, txm.Start())
	// TODO: double check the notify works via an enqueue
	require.NoError(t, txm.Close())
}

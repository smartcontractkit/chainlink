package terratxm_test

import (
	"testing"
	"time"

	ksmocks "github.com/smartcontractkit/chainlink/core/services/keystore/mocks"

	uuid "github.com/satori/go.uuid"
	tcmocks "github.com/smartcontractkit/chainlink-terra/pkg/terra/client/mocks"
	"github.com/smartcontractkit/chainlink/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager/terratxm"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/stretchr/testify/require"
)

func TestTxm(t *testing.T) {
	// Need full db to be able to test db trigger/functions, which only fire on tx commit
	// which won't happen using txdb.
	cfg, db := heavyweight.FullTestDB(t, "terra_txm", true, false)
	lggr := logger.TestLogger(t)
	eb := pg.NewEventBroadcaster(cfg.DatabaseURL(), 0, 0, lggr, uuid.NewV4())
	require.NoError(t, eb.Start())
	t.Cleanup(func() { require.NoError(t, eb.Close()) })
	ks := new(ksmocks.Terra)
	tc := new(tcmocks.ReaderWriter)
	txm := terratxm.NewTxm(db, tc, ks, lggr, pgtest.NewPGCfg(true), eb, time.Second)
	require.NoError(t, txm.Start())
	t.Cleanup(func() { require.NoError(t, txm.Close()) })
	require.NoError(t, txm.Enqueue("0x123", []byte(`hello`)))
}

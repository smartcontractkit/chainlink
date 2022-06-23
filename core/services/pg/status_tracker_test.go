package pg_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pg"
)

func TestStatusTracker(t *testing.T) {
	testutils.SkipShortDB(t)

	lggr := logger.TestLogger(t)
	db := pgtest.NewSqlxDB(t)

	connectedAwaiter := cltest.NewAwaiter()
	disconnectedAwaiter := cltest.NewAwaiter()
	handler := func(connected bool) {
		if connected {
			connectedAwaiter.ItHappened()
		} else {
			disconnectedAwaiter.ItHappened()
		}
	}

	st := pg.NewStatusTracker(db, 100*time.Millisecond, lggr)
	unsubscribe := st.Subscribe(handler)

	require.NoError(t, st.Start(testutils.Context(t)))
	connectedAwaiter.AwaitOrFail(t)

	err := db.DB.Close()
	require.NoError(t, err)

	disconnectedAwaiter.AwaitOrFail(t)
	unsubscribe()

	connectedAwaiter2 := cltest.NewAwaiter()
	handler2 := func(connected bool) {
		if connected {
			connectedAwaiter2.ItHappened()
		}
	}

	unsubscribe2 := st.Subscribe(handler2)

	db2 := pgtest.NewSqlxDB(t)
	db.DB = db2.DB

	connectedAwaiter2.AwaitOrFail(t)
	unsubscribe2()

	require.NoError(t, st.Close())
}

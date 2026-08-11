package v2

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
)

func TestDBModulePanicStore_RecordCountClear(t *testing.T) {
	t.Parallel()
	db := pgtest.NewSqlxDB(t)
	s := NewDBModulePanicStore(db)
	ctx := t.Context()

	count, err := modulePanicCount(ctx, s, "wf")
	require.NoError(t, err)
	require.Zero(t, count, "absent counter reads as zero")

	for want := uint64(1); want <= 3; want++ {
		got, rerr := recordModulePanic(ctx, s, "wf")
		require.NoError(t, rerr)
		require.Equal(t, want, got)
	}

	// Durably reads back the accumulated count.
	count, err = modulePanicCount(ctx, s, "wf")
	require.NoError(t, err)
	require.Equal(t, uint64(3), count)

	require.NoError(t, ClearModulePanic(ctx, s, "wf"))
	count, err = modulePanicCount(ctx, s, "wf")
	require.NoError(t, err)
	require.Zero(t, count, "clear removes the counter")

	// Clearing an absent counter is a no-op, not an error.
	require.NoError(t, ClearModulePanic(ctx, s, "never-seen"))
}

func TestDBModulePanicStore_IsPerWorkflow(t *testing.T) {
	t.Parallel()
	db := pgtest.NewSqlxDB(t)
	s := NewDBModulePanicStore(db)
	ctx := t.Context()

	_, err := recordModulePanic(ctx, s, "wf-a")
	require.NoError(t, err)

	count, err := modulePanicCount(ctx, s, "wf-b")
	require.NoError(t, err)
	require.Zero(t, count)
}

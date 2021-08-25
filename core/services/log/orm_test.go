package log_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/log"
)

func TestORM_MarkBroadcastConsumed(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	ethKeyStore := cltest.NewKeyStore(t, store.DB).Eth()

	orm := log.NewORM(store.DB, cltest.FixtureChainID)

	_, addr := cltest.MustAddRandomKeyToKeystore(t, ethKeyStore)
	specV2 := cltest.MustInsertV2JobSpec(t, store, addr)

	tests := []struct {
		name     string
		listener log.Listener
	}{
		{"v2", &mockListener{specV2.ID}},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			listener := test.listener

			rawLog := cltest.RandomLog(t)

			var consumed struct{ Consumed bool }
			var err error
			err = store.DB.Raw(`
				SELECT consumed FROM log_broadcasts
				WHERE block_hash = ? AND block_number = ? AND log_index = ? AND job_id = ?
			`, rawLog.BlockHash, rawLog.BlockNumber, rawLog.Index, listener.JobID()).Scan(&consumed).Error
			require.NoError(t, err)
			require.False(t, consumed.Consumed)

			err = orm.MarkBroadcastConsumed(store.DB, rawLog.BlockHash, rawLog.BlockNumber, rawLog.Index, listener.JobID())
			require.NoError(t, err)

			err = store.DB.Raw(`
				SELECT consumed FROM log_broadcasts
				WHERE block_hash = ? AND block_number = ? AND log_index = ? AND job_id = ?
			`, rawLog.BlockHash, rawLog.BlockNumber, rawLog.Index, listener.JobID()).Scan(&consumed).Error
			require.NoError(t, err)
			require.True(t, consumed.Consumed)
		})
	}
}

func TestORM_WasBroadcastConsumed(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	ethKeyStore := cltest.NewKeyStore(t, store.DB).Eth()

	orm := log.NewORM(store.DB, cltest.FixtureChainID)

	t.Run("returns the correct value", func(t *testing.T) {
		_, addr := cltest.MustAddRandomKeyToKeystore(t, ethKeyStore)
		specV2 := cltest.MustInsertV2JobSpec(t, store, addr)

		tests := []struct {
			name     string
			listener log.Listener
		}{
			{"v2", &mockListener{specV2.ID}},
		}

		for _, test := range tests {
			test := test
			t.Run(test.name, func(t *testing.T) {
				listener := test.listener

				rawLog := cltest.RandomLog(t)
				was, err := orm.WasBroadcastConsumed(store.DB, rawLog.BlockHash, rawLog.Index, listener.JobID())
				require.NoError(t, err)
				require.False(t, was)

				err = orm.MarkBroadcastConsumed(store.DB, rawLog.BlockHash, rawLog.BlockNumber, rawLog.Index, listener.JobID())
				require.NoError(t, err)

				was, err = orm.WasBroadcastConsumed(store.DB, rawLog.BlockHash, rawLog.Index, listener.JobID())
				require.NoError(t, err)
				require.True(t, was)
			})
		}
	})

	t.Run("does not error if the record doesn't exist", func(t *testing.T) {
		_, addr := cltest.MustAddRandomKeyToKeystore(t, ethKeyStore)
		specV2 := cltest.MustInsertV2JobSpec(t, store, addr)

		tests := []struct {
			name     string
			listener log.Listener
		}{
			{"v2", &mockListener{specV2.ID}},
		}

		for _, test := range tests {
			test := test
			t.Run(test.name, func(t *testing.T) {
				listener := test.listener

				rawLog := cltest.RandomLog(t)
				_, err := orm.WasBroadcastConsumed(store.DB, rawLog.BlockHash, rawLog.Index, listener.JobID())
				require.NoError(t, err)
			})
		}
	})
}

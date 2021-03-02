package log_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/log"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

func TestORM_MarkBroadcastConsumed(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	orm := log.NewORM(store.DB)

	t.Run("sets consumed to false if the record exists", func(t *testing.T) {
		_, addr := cltest.MustAddRandomKeyToKeystore(t, store)
		specV1 := cltest.MustInsertJobSpec(t, store)
		specV2 := cltest.MustInsertV2JobSpec(t, store, addr)

		tests := []struct {
			name     string
			listener log.Listener
		}{
			{"v1", &mockListener{specV1.ID, 0}},
			{"v2", &mockListener{models.NilJobID, specV2.ID}},
		}

		for _, test := range tests {
			test := test
			t.Run(test.name, func(t *testing.T) {
				listener := test.listener

				rawLog := cltest.RandomLog(t)
				cltest.MustInsertLog(t, rawLog, store)

				err := orm.MarkBroadcastConsumed(rawLog.BlockHash, rawLog.Index, log.ListenerJobID(listener))
				require.NoError(t, err)

				var consumed struct{ Consumed bool }
				if listener.IsV2Job() {
					err = store.DB.Raw(`
                        SELECT consumed FROM log_broadcasts
                        WHERE block_hash = ? AND block_number = ? AND log_index = ? AND job_id_v2 = ?
                    `, rawLog.BlockHash, rawLog.BlockNumber, rawLog.Index, listener.JobIDV2()).Scan(&consumed).Error
				} else {
					err = store.DB.Raw(`
                        SELECT consumed FROM log_broadcasts
                        WHERE block_hash = ? AND block_number = ? AND log_index = ? AND job_id = ?
                    `, rawLog.BlockHash, rawLog.BlockNumber, rawLog.Index, listener.JobID()).Scan(&consumed).Error
				}
				require.NoError(t, err)
				require.True(t, consumed.Consumed)
			})
		}
	})

	t.Run("errors if log_broadcast record does not exist", func(t *testing.T) {
		t.Run("v1", func(t *testing.T) {
			specV1 := cltest.MustInsertJobSpec(t, store)

			log := cltest.RandomLog(t)
			cltest.MustInsertLog(t, log, store)

			err := orm.MarkBroadcastConsumed(log.BlockHash, log.Index, specV1.ID)
			require.Error(t, err)
		})

		t.Run("v2", func(t *testing.T) {
			_, addr := cltest.MustAddRandomKeyToKeystore(t, store)
			specV2 := cltest.MustInsertV2JobSpec(t, store, addr)

			log := cltest.RandomLog(t)
			cltest.MustInsertLog(t, log, store)

			err := orm.MarkBroadcastConsumed(log.BlockHash, log.Index, specV2.ID)
			require.Error(t, err)
		})
	})
}

func TestORM_WasBroadcastConsumed(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	orm := log.NewORM(store.DB)

	t.Run("returns the correct value", func(t *testing.T) {
		_, addr := cltest.MustAddRandomKeyToKeystore(t, store)
		specV1 := cltest.MustInsertJobSpec(t, store)
		specV2 := cltest.MustInsertV2JobSpec(t, store, addr)

		tests := []struct {
			name     string
			listener log.Listener
		}{
			{"v1", &mockListener{specV1.ID, 0}},
			{"v2", &mockListener{models.NilJobID, specV2.ID}},
		}

		for _, test := range tests {
			test := test
			t.Run(test.name, func(t *testing.T) {
				listener := test.listener

				rawLog := cltest.RandomLog(t)
				cltest.MustInsertLog(t, rawLog, store)

				was, err := orm.WasBroadcastConsumed(rawLog.BlockHash, rawLog.Index, log.ListenerJobID(listener))
				require.NoError(t, err)
				require.False(t, was)

				err = orm.MarkBroadcastConsumed(rawLog.BlockHash, rawLog.Index, log.ListenerJobID(listener))
				require.NoError(t, err)

				was, err = orm.WasBroadcastConsumed(rawLog.BlockHash, rawLog.Index, log.ListenerJobID(listener))
				require.NoError(t, err)
				require.True(t, was)
			})
		}
	})

	t.Run("returns an error if the record doesn't exist", func(t *testing.T) {
		_, addr := cltest.MustAddRandomKeyToKeystore(t, store)
		specV1 := cltest.MustInsertJobSpec(t, store)
		specV2 := cltest.MustInsertV2JobSpec(t, store, addr)

		tests := []struct {
			name     string
			listener log.Listener
		}{
			{"v1", &mockListener{specV1.ID, 0}},
			{"v2", &mockListener{models.NilJobID, specV2.ID}},
		}

		for _, test := range tests {
			test := test
			t.Run(test.name, func(t *testing.T) {
				listener := test.listener

				rawLog := cltest.RandomLog(t)
				cltest.MustInsertLog(t, rawLog, store)

				_, err := orm.WasBroadcastConsumed(rawLog.BlockHash, rawLog.Index, log.ListenerJobID(listener))
				require.Error(t, err)
			})
		}
	})
}

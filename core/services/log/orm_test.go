package log_test

import (
	"bytes"
	"sort"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/log"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

type logRow struct {
	Address     common.Address
	Topics      pq.ByteaArray
	Data        []byte
	BlockNumber uint64
	BlockHash   common.Hash
	Index       uint
}

type logBroadcastRow struct {
	BlockHash   common.Hash
	BlockNumber uint64
	LogIndex    uint
	JobID       *models.ID
	JobIDV2     int32
	Consumed    bool
}

type mockListener struct {
	jobID   *models.ID
	jobIDV2 int32
}

func (l *mockListener) JobID() *models.ID              { return l.jobID }
func (l *mockListener) JobIDV2() int32                 { return l.jobIDV2 }
func (l *mockListener) IsV2Job() bool                  { return l.jobID == nil }
func (l *mockListener) OnConnect()                     {}
func (l *mockListener) OnDisconnect()                  {}
func (l *mockListener) HandleLog(log.Broadcast, error) {}

func TestORM_UpsertLog(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	orm := log.NewORM(store.DB)

	var logs []types.Log
	for i := 0; i < 10; i++ {
		logs = append(logs, cltest.RandomLog(t))
	}

	// Upsert twice
	for _, log := range logs {
		err := orm.UpsertLog(log)
		require.NoError(t, err)
	}

	for _, log := range logs {
		err := orm.UpsertLog(log)
		require.NoError(t, err)
	}

	var logRows []logRow
	err := store.DB.Raw(`SELECT * FROM logs ORDER BY block_hash, index ASC`).Scan(&logRows).Error
	require.NoError(t, err)

	logRowsDecoded := make([]types.Log, len(logRows))
	for i, log := range logRows {
		topics := make([]common.Hash, len(log.Topics))
		bytesTopics := [][]byte(log.Topics)
		for j, topic := range bytesTopics {
			topics[j] = common.BytesToHash(topic)
		}
		logRowsDecoded[i] = types.Log{
			Address:     log.Address,
			Topics:      topics,
			Data:        log.Data,
			BlockNumber: log.BlockNumber,
			BlockHash:   log.BlockHash,
			Index:       log.Index,
		}
	}

	require.Len(t, logRows, len(logs))

	sort.Slice(logs, func(i, j int) bool {
		if x := bytes.Compare(logs[i].BlockHash[:], logs[j].BlockHash[:]); x == 0 {
			return logs[i].Index < logs[j].Index
		} else {
			return x < 0
		}
	})
	require.Equal(t, logs, logRowsDecoded)
}

func TestORM_UpsertLogBroadcastForListener(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	orm := log.NewORM(store.DB)

	specV1_1 := cltest.MustInsertJobSpec(t, store)
	specV1_2 := cltest.MustInsertJobSpec(t, store)
	specV1_3 := cltest.MustInsertJobSpec(t, store)
	specV1_4 := cltest.MustInsertJobSpec(t, store)
	specV1_5 := cltest.MustInsertJobSpec(t, store)

	_, addr := cltest.MustAddRandomKeyToKeystore(t, store)
	specV2_1 := cltest.MustInsertV2JobSpec(t, store, addr)
	specV2_2 := cltest.MustInsertV2JobSpec(t, store, addr)
	specV2_3 := cltest.MustInsertV2JobSpec(t, store, addr)
	specV2_4 := cltest.MustInsertV2JobSpec(t, store, addr)
	specV2_5 := cltest.MustInsertV2JobSpec(t, store, addr)

	rawLog := cltest.RandomLog(t)

	err := orm.UpsertLog(rawLog)
	require.NoError(t, err)

	listeners := []log.Listener{
		&mockListener{specV1_1.ID, 0},
		&mockListener{specV1_2.ID, 0},
		&mockListener{specV1_3.ID, 0},
		&mockListener{specV1_4.ID, 0},
		&mockListener{specV1_5.ID, 0},
		&mockListener{nil, specV2_1.ID},
		&mockListener{nil, specV2_2.ID},
		&mockListener{nil, specV2_3.ID},
		&mockListener{nil, specV2_4.ID},
		&mockListener{nil, specV2_5.ID},
	}

	sort.Slice(listeners[:5], func(i, j int) bool {
		return bytes.Compare((*listeners[i].JobID())[:], (*listeners[j].JobID())[:]) < 0
	})

	t.Run("does not error when upserting the same entry more than once", func(t *testing.T) {
		// Upsert twice
		for _, listener := range listeners {
			err := orm.UpsertLogBroadcastForListener(rawLog, listener.JobID(), listener.JobIDV2())
			require.NoError(t, err)
		}
		for _, listener := range listeners {
			err := orm.UpsertLogBroadcastForListener(rawLog, listener.JobID(), listener.JobIDV2())
			require.NoError(t, err)
		}
	})

	t.Run("does not duplicate an entry when upserting it more than once", func(t *testing.T) {
		var count struct{ Count int }
		err := store.DB.Raw(`SELECT count(*) FROM log_broadcasts`).Scan(&count).Error
		require.NoError(t, err)
		require.Equal(t, len(listeners), count.Count)
	})

	t.Run("upserts the correct values", func(t *testing.T) {
		expected := make([]logBroadcastRow, 10)
		for i := range listeners {
			expected[i] = logBroadcastRow{
				rawLog.BlockHash,
				rawLog.BlockNumber,
				rawLog.Index,
				listeners[i].JobID(),
				listeners[i].JobIDV2(),
				false,
			}
		}

		var logBroadcastRows []logBroadcastRow
		err := store.DB.Raw(`SELECT * FROM log_broadcasts ORDER BY job_id, job_id_v2 ASC`).Scan(&logBroadcastRows).Error
		require.NoError(t, err)
		require.Len(t, logBroadcastRows, len(expected))
		require.Equal(t, expected, logBroadcastRows)
	})

	t.Run("does not reset `consumed` to `false` for existing records", func(t *testing.T) {
		err := store.DB.Exec(`UPDATE log_broadcasts SET consumed = true`).Error
		require.NoError(t, err)

		var consumed struct{ Consumed bool }
		err = store.DB.Raw(`
            SELECT consumed FROM log_broadcasts WHERE job_id = ? AND block_hash = ? AND log_index = ?
        `, listeners[0].JobID(), rawLog.BlockHash, rawLog.Index).Scan(&consumed).Error
		require.NoError(t, err)
		require.True(t, consumed.Consumed)

		err = orm.UpsertLogBroadcastForListener(rawLog, listeners[0].JobID(), listeners[0].JobIDV2())
		require.NoError(t, err)

		err = store.DB.Raw(`
            SELECT consumed FROM log_broadcasts WHERE job_id = ? AND block_hash = ? AND log_index = ?
        `, listeners[0].JobID(), rawLog.BlockHash, rawLog.Index).Scan(&consumed).Error
		require.NoError(t, err)
		require.True(t, consumed.Consumed)
	})
}

func TestORM_MarkBroadcastConsumed(t *testing.T) {
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
			{"v2", &mockListener{nil, specV2.ID}},
		}

		for _, test := range tests {
			test := test
			t.Run(test.name, func(t *testing.T) {
				listener := test.listener

				log := cltest.RandomLog(t)
				cltest.MustInsertLog(t, log, store)

				err := orm.UpsertLogBroadcastForListener(log, listener.JobID(), listener.JobIDV2())
				require.NoError(t, err)

				err = orm.MarkBroadcastConsumed(log.BlockHash, log.Index, listener.JobID(), listener.JobIDV2())
				require.NoError(t, err)

				var consumed struct{ Consumed bool }
				if listener.IsV2Job() {
					err = store.DB.Raw(`
                        SELECT consumed FROM log_broadcasts
                        WHERE block_hash = ? AND block_number = ? AND log_index = ? AND job_id_v2 = ?
                    `, log.BlockHash, log.BlockNumber, log.Index, listener.JobIDV2()).Scan(&consumed).Error
				} else {
					err = store.DB.Raw(`
                        SELECT consumed FROM log_broadcasts
                        WHERE block_hash = ? AND block_number = ? AND log_index = ? AND job_id = ?
                    `, log.BlockHash, log.BlockNumber, log.Index, listener.JobID()).Scan(&consumed).Error
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

			err := orm.MarkBroadcastConsumed(log.BlockHash, log.Index, specV1.ID, 0)
			require.Error(t, err)
		})

		t.Run("v2", func(t *testing.T) {
			_, addr := cltest.MustAddRandomKeyToKeystore(t, store)
			specV2 := cltest.MustInsertV2JobSpec(t, store, addr)

			log := cltest.RandomLog(t)
			cltest.MustInsertLog(t, log, store)

			err := orm.MarkBroadcastConsumed(log.BlockHash, log.Index, nil, specV2.ID)
			require.Error(t, err)
		})
	})
}

func TestORM_WasBroadcastConsumed(t *testing.T) {
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
			{"v2", &mockListener{nil, specV2.ID}},
		}

		for _, test := range tests {
			test := test
			t.Run(test.name, func(t *testing.T) {
				listener := test.listener

				log := cltest.RandomLog(t)
				cltest.MustInsertLog(t, log, store)

				err := orm.UpsertLogBroadcastForListener(log, listener.JobID(), listener.JobIDV2())
				require.NoError(t, err)

				was, err := orm.WasBroadcastConsumed(log.BlockHash, log.Index, listener.JobID(), listener.JobIDV2())
				require.NoError(t, err)
				require.False(t, was)

				err = orm.MarkBroadcastConsumed(log.BlockHash, log.Index, listener.JobID(), listener.JobIDV2())
				require.NoError(t, err)

				was, err = orm.WasBroadcastConsumed(log.BlockHash, log.Index, listener.JobID(), listener.JobIDV2())
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
			{"v2", &mockListener{nil, specV2.ID}},
		}

		for _, test := range tests {
			test := test
			t.Run(test.name, func(t *testing.T) {
				listener := test.listener

				log := cltest.RandomLog(t)
				cltest.MustInsertLog(t, log, store)

				_, err := orm.WasBroadcastConsumed(log.BlockHash, log.Index, listener.JobID(), listener.JobIDV2())
				require.Error(t, err)
			})
		}
	})
}

func TestORM_UnconsumedLogsPriorToBlock(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	_, addr := cltest.MustAddRandomKeyToKeystore(t, store)
	specV1 := cltest.MustInsertJobSpec(t, store)
	specV2 := cltest.MustInsertV2JobSpec(t, store, addr)

	orm := log.NewORM(store.DB)

	var logs []types.Log
	for i := 0; i < 20; i++ {
		log := cltest.RandomLog(t)
		log.BlockNumber = uint64(i)
		logs = append(logs, log)
	}

	sort.Slice(logs, func(i, j int) bool {
		if logs[i].BlockNumber < logs[j].BlockNumber {
			return true
		} else if logs[i].BlockNumber == logs[j].BlockNumber {
			return logs[i].Index < logs[j].Index
		}
		return false
	})

	listeners := []log.Listener{
		&mockListener{specV1.ID, 0},
		&mockListener{nil, specV2.ID},
	}

	for i, log := range logs {
		err := orm.UpsertLog(log)
		require.NoError(t, err)

		for _, listener := range listeners {
			err := orm.UpsertLogBroadcastForListener(log, listener.JobID(), listener.JobIDV2())
			require.NoError(t, err)
		}

		var js []int
		for j := 0; j < len(listeners); j++ {
			if j < i%(len(listeners)+1) {
				err := orm.MarkBroadcastConsumed(log.BlockHash, log.Index, listeners[j].JobID(), listeners[j].JobIDV2())
				require.NoError(t, err)
				js = append(js, j)
			}
		}
	}

	blockNumber := uint64(len(logs)/2 + 1)

	var expected []types.Log
	for i, log := range logs {
		if i%(len(listeners)+1) < 2 && log.BlockNumber < blockNumber {
			expected = append(expected, log)
		}
	}

	fetchedLogs, err := orm.UnconsumedLogsPriorToBlock(blockNumber)
	require.NoError(t, err)
	require.Equal(t, expected, fetchedLogs)
}

func TestORM_DeleteLogAndBroadcasts(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	orm := log.NewORM(store.DB)

	_, addr := cltest.MustAddRandomKeyToKeystore(t, store)
	specV1 := cltest.MustInsertJobSpec(t, store)
	specV2 := cltest.MustInsertV2JobSpec(t, store, addr)

	listeners := []log.Listener{
		&mockListener{specV1.ID, 0},
		&mockListener{nil, specV2.ID},
	}

	t.Run("correctly deletes a log and all of its associated broadcasts", func(t *testing.T) {
		log := cltest.RandomLog(t)
		cltest.MustInsertLog(t, log, store)

		for _, listener := range listeners {
			err := orm.UpsertLogBroadcastForListener(log, listener.JobID(), listener.JobIDV2())
			require.NoError(t, err)
		}

		var count struct{ Count int }
		err := store.DB.Raw(`SELECT count(*) FROM log_broadcasts`).Scan(&count).Error
		require.NoError(t, err)
		require.Equal(t, len(listeners), count.Count)

		err = orm.DeleteLogAndBroadcasts(log.BlockHash, log.Index)
		require.NoError(t, err)

		err = store.DB.Raw(`SELECT count(*) FROM log_broadcasts`).Scan(&count).Error
		require.NoError(t, err)
		require.Equal(t, 0, count.Count)
	})

	t.Run("does not error if the record does not exist", func(t *testing.T) {
		err := orm.DeleteLogAndBroadcasts(cltest.NewHash(), 123)
		require.NoError(t, err)
	})
}

func TestORM_DeleteUnconsumedBroadcastsForListener(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	orm := log.NewORM(store.DB)

	_, addr := cltest.MustAddRandomKeyToKeystore(t, store)
	specV1 := cltest.MustInsertJobSpec(t, store)
	specV2 := cltest.MustInsertV2JobSpec(t, store, addr)

	listeners := []log.Listener{
		&mockListener{specV1.ID, 0},
		&mockListener{nil, specV2.ID},
	}

	logs := []types.Log{cltest.RandomLog(t), cltest.RandomLog(t), cltest.RandomLog(t), cltest.RandomLog(t), cltest.RandomLog(t)}
	for _, log := range logs {
		cltest.MustInsertLog(t, log, store)
		for _, listener := range listeners {
			err := orm.UpsertLogBroadcastForListener(log, listener.JobID(), listener.JobIDV2())
			require.NoError(t, err)
		}
	}

	var count struct{ Count int }
	err := store.DB.Raw(`SELECT count(*) FROM log_broadcasts`).Scan(&count).Error
	require.NoError(t, err)
	require.Equal(t, len(logs)*len(listeners), count.Count)

	err = orm.DeleteUnconsumedBroadcastsForListener(listeners[0].JobID(), listeners[0].JobIDV2())
	require.NoError(t, err)

	err = store.DB.Raw(`SELECT count(*) FROM log_broadcasts`).Scan(&count).Error
	require.NoError(t, err)
	require.Equal(t, len(logs), count.Count)

	err = orm.DeleteUnconsumedBroadcastsForListener(listeners[1].JobID(), listeners[1].JobIDV2())
	require.NoError(t, err)

	err = store.DB.Raw(`SELECT count(*) FROM log_broadcasts`).Scan(&count).Error
	require.NoError(t, err)
	require.Equal(t, 0, count.Count)
}

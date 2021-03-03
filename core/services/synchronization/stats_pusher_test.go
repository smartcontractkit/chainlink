package synchronization_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/synchronization"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestStatsPusher(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	wsserver, wscleanup := cltest.NewEventWebSocketServer(t)
	defer wscleanup()

	explorerClient := synchronization.NewExplorerClient(wsserver.URL, "", "")
	err := explorerClient.Start()
	require.NoError(t, err)

	pusher := synchronization.NewStatsPusher(store.DB, explorerClient)
	pusher.Start()
	defer pusher.Close()

	require.NoError(t, store.ORM.RawDBWithAdvisoryLock(func(db *gorm.DB) error { return db.Create(&models.SyncEvent{}).Error }))
	pusher.PushNow()

	assert.Equal(t, 1, lenSyncEvents(t, store.ORM), "jobrun sync event should be created")
	cltest.CallbackOrTimeout(t, "ws server receives jobrun creation", func() {
		<-wsserver.ReceivedText
		err := wsserver.Broadcast(`{"status": 201}`)
		assert.NoError(t, err)
	})
	cltest.WaitForSyncEventCount(t, store.ORM, 0)
}

func TestStatsPusher_ClockTrigger(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	wsserver, wscleanup := cltest.NewEventWebSocketServer(t)
	defer wscleanup()

	clock := cltest.NewTriggerClock(t)
	explorerClient := synchronization.NewExplorerClient(wsserver.URL, "", "")
	err := explorerClient.Start()
	require.NoError(t, err)

	pusher := synchronization.NewStatsPusher(store.DB, explorerClient, clock)
	pusher.Start()
	defer pusher.Close()

	err = store.ORM.RawDBWithAdvisoryLock(func(db *gorm.DB) error {
		return db.Save(&models.SyncEvent{Body: string("")}).Error
	})
	require.NoError(t, err)

	clock.Trigger()
	cltest.CallbackOrTimeout(t, "ws server receives jobrun update", func() {
		<-wsserver.ReceivedText
		err := wsserver.Broadcast(`{"status": 201}`)
		assert.NoError(t, err)
	})
	cltest.WaitForSyncEventCount(t, store.ORM, 0)
}

func TestStatsPusher_NoAckLeavesEvent(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	wsserver, wscleanup := cltest.NewEventWebSocketServer(t)
	defer wscleanup()

	explorerClient := synchronization.NewExplorerClient(wsserver.URL, "", "")
	err := explorerClient.Start()
	require.NoError(t, err)

	pusher := synchronization.NewStatsPusher(store.DB, explorerClient)
	pusher.Start()
	defer pusher.Close()

	require.NoError(t, store.ORM.RawDBWithAdvisoryLock(func(db *gorm.DB) error { return db.Create(&models.SyncEvent{}).Error }))
	pusher.PushNow()

	assert.Equal(t, 1, lenSyncEvents(t, store.ORM), "jobrun sync event should be created")
	cltest.CallbackOrTimeout(t, "ws server receives jobrun creation", func() {
		<-wsserver.ReceivedText
	})
	cltest.AssertSyncEventCountStays(t, store.ORM, 1)
}

func TestStatsPusher_BadSyncLeavesEvent(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	wsserver, wscleanup := cltest.NewEventWebSocketServer(t)
	defer wscleanup()

	clock := cltest.NewTriggerClock(t)
	explorerClient := synchronization.NewExplorerClient(wsserver.URL, "", "")
	err := explorerClient.Start()
	require.NoError(t, err)

	pusher := synchronization.NewStatsPusher(store.DB, explorerClient, clock)
	pusher.Start()
	defer pusher.Close()

	require.NoError(t, store.ORM.RawDBWithAdvisoryLock(func(db *gorm.DB) error { return db.Create(&models.SyncEvent{}).Error }))

	assert.Equal(t, 1, lenSyncEvents(t, store.ORM), "jobrun sync event should be created")
	clock.Trigger()
	cltest.CallbackOrTimeout(t, "ws server receives jobrun creation", func() {
		<-wsserver.ReceivedText
		err := wsserver.Broadcast(`{"status": 500}`)
		assert.NoError(t, err)
	})
	cltest.AssertSyncEventCountStays(t, store.ORM, 1)
}

func lenSyncEvents(t *testing.T, orm *orm.ORM) int {
	count, err := orm.CountOf(&models.SyncEvent{})
	require.NoError(t, err)
	return count
}

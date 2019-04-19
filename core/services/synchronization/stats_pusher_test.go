package synchronization_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/synchronization"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStatsPusher(t *testing.T) {
	store, cleanup := cltest.NewStore()
	defer cleanup()

	wsserver, wscleanup := cltest.NewEventWebSocketServer(t)
	defer wscleanup()

	clock := cltest.NewTriggerClock()
	pusher := synchronization.NewStatsPusher(store.ORM, wsserver.URL, clock)
	pusher.Start()
	defer pusher.Close()

	j := cltest.NewJobWithWebInitiator()
	require.NoError(t, store.CreateJob(&j))

	jr := j.NewRun(j.Initiators[0])
	require.NoError(t, store.CreateJobRun(&jr))

	assert.Equal(t, 1, lenSyncEvents(t, store.ORM), "jobrun sync event should be created")
	clock.Trigger()
	cltest.CallbackOrTimeout(t, "ws server receives jobrun creation", func() {
		<-wsserver.Received
	})
	cltest.WaitForSyncEventCount(t, store.ORM, 0)

	jr.ApplyResult(models.RunResult{Status: models.RunStatusCompleted})
	require.NoError(t, store.SaveJobRun(&jr))
	assert.Equal(t, 1, lenSyncEvents(t, store.ORM))

	clock.Trigger()
	cltest.CallbackOrTimeout(t, "ws server receives jobrun update", func() {
		<-wsserver.Received
	})
	cltest.WaitForSyncEventCount(t, store.ORM, 0)
}

func lenSyncEvents(t *testing.T, orm *orm.ORM) int {
	var count int
	require.NoError(t, orm.DB.Model(&models.SyncEvent{}).Count(&count).Error)
	return count
}

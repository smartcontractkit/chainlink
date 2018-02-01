package services_test

import (
	"testing"

	. "github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/services"
	strpkg "github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/stretchr/testify/assert"
)

func TestNotificationListenerStart(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
	defer cleanup()
	eth := cltest.MockEthOnStore(store)
	nl := services.NotificationListener{Store: store}
	defer nl.Stop()

	assert.Nil(t, store.SaveJob(cltest.NewJobWithLogInitiator()))
	assert.Nil(t, store.SaveJob(cltest.NewJobWithLogInitiator()))
	eth.RegisterSubscription("logs", make(chan strpkg.EthNotification))
	eth.RegisterSubscription("logs", make(chan strpkg.EthNotification))

	nl.Start()

	assert.True(t, eth.AllCalled())
}

func TestNotificationListenerAddJob(t *testing.T) {
	t.Parallel()
	RegisterTestingT(t)

	store, cleanup := cltest.NewStore()
	defer cleanup()
	eth := cltest.MockEthOnStore(store)
	nl := services.NotificationListener{Store: store}
	defer nl.Stop()
	nl.Start()

	j := cltest.NewJobWithLogInitiator()
	assert.Nil(t, store.SaveJob(j))
	logChan := make(chan strpkg.EthNotification, 1)
	initr := j.Initiators[0]
	eth.RegisterSubscription("logs", logChan)

	nl.AddJob(j)

	logChan <- cltest.NewEthNotification(strpkg.EventLog{Address: initr.Address})
	jobRuns := []*models.JobRun{}
	Eventually(func() []*models.JobRun {
		store.Where("JobID", j.ID, &jobRuns)
		return jobRuns
	}).Should(HaveLen(1))

	assert.True(t, eth.AllCalled())
}

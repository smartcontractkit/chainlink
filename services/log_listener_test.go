package services_test

import (
	"testing"

	. "github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink-go/internal/cltest"
	"github.com/smartcontractkit/chainlink-go/services"
	strpkg "github.com/smartcontractkit/chainlink-go/store"
	"github.com/smartcontractkit/chainlink-go/store/models"
	"github.com/stretchr/testify/assert"
)

func TestLogListenerStart(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
	defer cleanup()
	eth := cltest.MockEthOnStore(store)
	ll := services.LogListener{Store: store}
	defer ll.Stop()

	assert.Nil(t, store.SaveJob(cltest.NewJobWithLogInitiator()))
	assert.Nil(t, store.SaveJob(cltest.NewJobWithLogInitiator()))
	eth.RegisterSubscription("logs", make(chan strpkg.EventLog))
	eth.RegisterSubscription("logs", make(chan strpkg.EventLog))

	ll.Start()

	assert.True(t, eth.AllCalled())
}

func TestLogListenerAddJob(t *testing.T) {
	t.Parallel()
	RegisterTestingT(t)

	store, cleanup := cltest.NewStore()
	defer cleanup()
	eth := cltest.MockEthOnStore(store)
	ll := services.LogListener{Store: store}
	defer ll.Stop()
	ll.Start()

	j := cltest.NewJobWithLogInitiator()
	assert.Nil(t, store.SaveJob(j))
	logChan := make(chan strpkg.EventLog, 1)
	initr := j.Initiators[0]
	eth.RegisterSubscription("logs", logChan)

	ll.AddJob(j)

	logChan <- strpkg.EventLog{Address: initr.Address}
	jobRuns := []models.JobRun{}
	Eventually(func() []models.JobRun {
		store.Where("JobID", j.ID, &jobRuns)
		return jobRuns
	}).Should(HaveLen(1))

	assert.True(t, eth.AllCalled())
}

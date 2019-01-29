// +build !windows

package services_test

import (
	"syscall"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/services/mock_services"
	strpkg "github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tevino/abool"
)

func TestChainlinkApplication_SignalShutdown(t *testing.T) {
	config, cleanup := cltest.NewConfig()
	defer cleanup()
	app, appCleanUp := cltest.NewApplicationWithConfig(config)
	defer appCleanUp()

	completed := abool.New()
	app.Exiter = func(code int) {
		completed.Set()
	}

	app.Start()
	syscall.Kill(syscall.Getpid(), syscall.SIGTERM)

	gomega.NewGomegaWithT(t).Eventually(func() bool {
		return completed.IsSet()
	}).Should(gomega.BeTrue())
}

func TestChainlinkApplication_AddJob(t *testing.T) {
	app, cleanup := cltest.NewApplication()
	defer cleanup()
	ctrl := gomock.NewController(t)
	jobSubscriberMock := mock_services.NewMockJobSubscriber(ctrl)
	app.ChainlinkApplication.JobSubscriber = jobSubscriberMock
	jobSubscriberMock.EXPECT().AddJob(gomock.Any(), nil) // nil to represent "latest" block
	app.AddJob(cltest.NewJob())
}

func TestChainlinkApplication_resumesPendingConnection(t *testing.T) {
	app, cleanup := cltest.NewApplication()
	defer cleanup()
	store := app.Store

	j := cltest.NewJobWithWebInitiator()
	require.NoError(t, store.CreateJob(&j))

	jr := cltest.CreateJobRunWithStatus(store, j, models.RunStatusPendingConnection)

	require.NoError(t, app.Start())
	_ = cltest.WaitForJobRunToComplete(t, store, jr)
}

func TestPendingConnectionResumer(t *testing.T) {
	store, cleanup := cltest.NewStore()
	defer cleanup()

	resumedRuns := []string{}
	resumer := func(run *models.JobRun, store *strpkg.Store) (*models.JobRun, error) {
		resumedRuns = append(resumedRuns, run.ID)
		return nil, nil
	}
	pcr := services.ExportedNewPendingConnectionResumer(store, resumer)

	j := cltest.NewJobWithWebInitiator()
	require.NoError(t, store.CreateJob(&j))

	expectedRun := cltest.CreateJobRunWithStatus(store, j, models.RunStatusPendingConnection)
	_ = cltest.CreateJobRunWithStatus(store, j, models.RunStatusPendingConfirmations)
	_ = cltest.CreateJobRunWithStatus(store, j, models.RunStatusInProgress)
	_ = cltest.CreateJobRunWithStatus(store, j, models.RunStatusUnstarted)
	_ = cltest.CreateJobRunWithStatus(store, j, models.RunStatusPendingBridge)
	_ = cltest.CreateJobRunWithStatus(store, j, models.RunStatusInProgress)
	_ = cltest.CreateJobRunWithStatus(store, j, models.RunStatusCompleted)

	assert.NoError(t, pcr.Connect(cltest.IndexableBlockNumber(1)))
	assert.Equal(t, []string{expectedRun.ID}, resumedRuns)
}

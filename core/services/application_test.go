// +build !windows

package services_test

import (
	"syscall"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/services/mock_services"
	strpkg "github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tevino/abool"
)

func TestChainlinkApplication_SignalShutdown(t *testing.T) {
	config, cleanup := cltest.NewConfig(t)
	defer cleanup()
	app, appCleanUp := cltest.NewApplicationWithConfig(t, config)
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
	app, cleanup := cltest.NewApplication(t)
	defer cleanup()
	ctrl := gomock.NewController(t)
	jobSubscriberMock := mock_services.NewMockJobSubscriber(ctrl)
	app.ChainlinkApplication.JobSubscriber = jobSubscriberMock
	jobSubscriberMock.EXPECT().AddJob(gomock.Any(), nil) // nil to represent "latest" block
	app.AddJob(cltest.NewJob())
}

func TestChainlinkApplication_resumesPendingConnection_Happy(t *testing.T) {
	app, cleanup := cltest.NewApplication(t)
	defer cleanup()
	store := app.Store

	j := cltest.NewJobWithWebInitiator()
	require.NoError(t, store.CreateJob(&j))

	jr := cltest.CreateJobRunWithStatus(t, store, j, models.RunStatusPendingConnection)

	require.NoError(t, utils.JustError(app.MockStartAndConnect()))
	_ = cltest.WaitForJobRunToComplete(t, store, jr)
}

func TestChainlinkApplication_resumesPendingConnection_Archived(t *testing.T) {
	app, cleanup := cltest.NewApplication(t)
	defer cleanup()
	store := app.Store

	j := cltest.NewJobWithWebInitiator()
	require.NoError(t, store.CreateJob(&j))

	jr := cltest.CreateJobRunWithStatus(t, store, j, models.RunStatusPendingConnection)

	require.NoError(t, store.ArchiveJob(j.ID))

	require.NoError(t, utils.JustError(app.MockStartAndConnect()))
	_ = cltest.WaitForJobRunToComplete(t, store, jr)
}

func TestPendingConnectionResumer(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	resumedRuns := []string{}
	resumer := func(run *models.JobRun, store *strpkg.Store) error {
		resumedRuns = append(resumedRuns, run.ID)
		return nil
	}
	pcr := services.ExportedNewPendingConnectionResumer(store, resumer)

	j := cltest.NewJobWithWebInitiator()
	require.NoError(t, store.CreateJob(&j))

	expectedRun := cltest.CreateJobRunWithStatus(t, store, j, models.RunStatusPendingConnection)
	_ = cltest.CreateJobRunWithStatus(t, store, j, models.RunStatusPendingConfirmations)
	_ = cltest.CreateJobRunWithStatus(t, store, j, models.RunStatusInProgress)
	_ = cltest.CreateJobRunWithStatus(t, store, j, models.RunStatusUnstarted)
	_ = cltest.CreateJobRunWithStatus(t, store, j, models.RunStatusPendingBridge)
	_ = cltest.CreateJobRunWithStatus(t, store, j, models.RunStatusInProgress)
	_ = cltest.CreateJobRunWithStatus(t, store, j, models.RunStatusCompleted)

	assert.NoError(t, pcr.Connect(cltest.Head(1)))
	assert.Equal(t, []string{expectedRun.ID}, resumedRuns)
}

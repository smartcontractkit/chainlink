// +build !windows

package services_test

import (
	"syscall"
	"testing"

	"chainlink/core/internal/cltest"
	"chainlink/core/internal/mocks"
	"chainlink/core/store/models"
	"chainlink/core/utils"

	"github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/tevino/abool"
)

func TestChainlinkApplication_SignalShutdown(t *testing.T) {
	config, cleanup := cltest.NewConfig(t)
	defer cleanup()
	app, appCleanUp := cltest.NewApplicationWithConfig(t, config)
	defer appCleanUp()
	eth := app.MockCallerSubscriberClient(cltest.Strict)
	eth.Register("eth_chainId", app.Store.Config.ChainID())

	completed := abool.New()
	app.Exiter = func(code int) {
		completed.Set()
	}

	require.NoError(t, app.Start())
	syscall.Kill(syscall.Getpid(), syscall.SIGTERM)

	gomega.NewGomegaWithT(t).Eventually(func() bool {
		return completed.IsSet()
	}).Should(gomega.BeTrue())
}

func TestChainlinkApplication_AddJob(t *testing.T) {
	app, cleanup := cltest.NewApplication(t)
	defer cleanup()
	require.NoError(t, app.Start())

	jobSubscriber := new(mocks.JobSubscriber)
	jobSubscriber.On("AddJob", mock.Anything, (*models.Head)(nil)).Return(nil, nil)
	app.ChainlinkApplication.JobSubscriber = jobSubscriber

	fluxMonitor := new(mocks.FluxMonitor)
	fluxMonitor.On("AddJob", mock.Anything).Return(nil)
	app.ChainlinkApplication.FluxMonitor = fluxMonitor

	app.AddJob(cltest.NewJob())

	jobSubscriber.AssertExpectations(t)
	fluxMonitor.AssertExpectations(t)
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

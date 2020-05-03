// +build !windows

package chainlink_test

import (
	"syscall"
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store/models"

	"github.com/onsi/gomega"
	"github.com/stretchr/testify/require"
	"github.com/tevino/abool"
)

func TestChainlinkApplication_SignalShutdown(t *testing.T) {
	config, cleanup := cltest.NewConfig(t)
	defer cleanup()
	app, appCleanUp := cltest.NewApplicationWithConfig(t, config, cltest.EthMockRegisterChainID)
	defer appCleanUp()

	completed := abool.New()
	app.Exiter = func(code int) {
		completed.Set()
	}

	require.NoError(t, app.Start())
	require.NoError(t, syscall.Kill(syscall.Getpid(), syscall.SIGTERM))

	gomega.NewGomegaWithT(t).Eventually(func() bool {
		return completed.IsSet()
	}).Should(gomega.BeTrue())
}

func TestChainlinkApplication_resumesPendingConnection_Happy(t *testing.T) {
	app, cleanup := cltest.NewApplication(t)
	defer cleanup()
	app.EthMock.Context("app.Start()", func(meth *cltest.EthMock) {
		meth.Register("eth_chainId", app.Store.Config.ChainID())
	})
	store := app.Store

	j := cltest.NewJobWithWebInitiator()
	require.NoError(t, store.CreateJob(&j))

	jr := cltest.CreateJobRunWithStatus(t, store, j, models.RunStatusPendingConnection)

	require.NoError(t, app.StartAndConnect())
	_ = cltest.WaitForJobRunToComplete(t, store, jr)
}

func TestChainlinkApplication_resumesPendingConnection_Archived(t *testing.T) {
	app, cleanup := cltest.NewApplication(t)
	defer cleanup()
	app.EthMock.Context("app.Start()", func(meth *cltest.EthMock) {
		meth.Register("eth_chainId", app.Store.Config.ChainID())
	})
	store := app.Store

	j := cltest.NewJobWithWebInitiator()
	require.NoError(t, store.CreateJob(&j))

	jr := cltest.CreateJobRunWithStatus(t, store, j, models.RunStatusPendingConnection)

	require.NoError(t, store.ArchiveJob(j.ID))

	require.NoError(t, app.StartAndConnect())
	_ = cltest.WaitForJobRunToComplete(t, store, jr)
}

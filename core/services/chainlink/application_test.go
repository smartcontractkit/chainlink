// +build !windows

package chainlink_test

import (
	"syscall"
	"testing"

	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/static"
	"github.com/smartcontractkit/chainlink/core/store/migrations"
	"github.com/smartcontractkit/chainlink/core/store/migrationsv2"

	"github.com/smartcontractkit/chainlink/core/services/eth"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store/models"

	"github.com/onsi/gomega"
	"github.com/stretchr/testify/require"
	"github.com/tevino/abool"
)

func TestChainlinkApplication_SquashMigrationUpgrade(t *testing.T) {
	_, orm, cleanup := cltest.BootstrapThrowawayORM(t, "migrationssquash", false)
	defer cleanup()
	db := orm.DB

	// Latest migrations should work fine.
	static.Version = "0.9.11"
	err := migrationsv2.MigrateUp(db, "")
	require.NoError(t, err)
	err = chainlink.CheckSquashUpgrade(db)
	require.NoError(t, err)
	err = migrationsv2.MigrateDown(db)
	require.NoError(t, err)

	// Newer app version with older migrations should fail.
	err = migrations.MigrateTo(db, "1611388693") // 1 before S-1
	require.NoError(t, err)
	err = chainlink.CheckSquashUpgrade(db)
	t.Log(err)
	require.Error(t, err)

	static.Version = "unset"
}

func TestChainlinkApplication_SignalShutdown(t *testing.T) {

	rpcClient, gethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	app, cleanup := cltest.NewApplication(t,
		eth.NewClientWith(rpcClient, gethClient),
	)
	defer cleanup()
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

func TestChainlinkApplication_resumesPendingConnection_Happy(t *testing.T) {
	rpcClient, gethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	app, cleanup := cltest.NewApplication(t,
		eth.NewClientWith(rpcClient, gethClient),
	)
	defer cleanup()
	store := app.Store

	j := cltest.NewJobWithWebInitiator()
	require.NoError(t, store.CreateJob(&j))

	jr := cltest.CreateJobRunWithStatus(t, store, j, models.RunStatusPendingConnection)

	require.NoError(t, app.StartAndConnect())
	_ = cltest.WaitForJobRunToComplete(t, store, jr)
}

func TestChainlinkApplication_resumesPendingConnection_Archived(t *testing.T) {
	rpcClient, gethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	app, cleanup := cltest.NewApplication(t,
		eth.NewClientWith(rpcClient, gethClient),
	)
	defer cleanup()
	store := app.Store

	j := cltest.NewJobWithWebInitiator()
	require.NoError(t, store.CreateJob(&j))

	jr := cltest.CreateJobRunWithStatus(t, store, j, models.RunStatusPendingConnection)

	require.NoError(t, store.ArchiveJob(j.ID))

	require.NoError(t, app.StartAndConnect())
	_ = cltest.WaitForJobRunToComplete(t, store, jr)
}

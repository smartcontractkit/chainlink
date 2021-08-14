// +build !windows

package chainlink_test

import (
	"syscall"
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/stretchr/testify/assert"

	"github.com/onsi/gomega"
	"github.com/stretchr/testify/require"
	"github.com/tevino/abool"
)

func TestChainlinkApplication_SignalShutdown(t *testing.T) {
	ethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	app, cleanup := cltest.NewApplication(t,
		ethClient,
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

func TestChainlinkApplication_ChangeInChainID(t *testing.T) {
	config := cltest.NewTestEVMConfig(t)
	app, cleanup := cltest.NewApplicationWithConfig(t, config)
	defer cleanup()

	require.NoError(t, app.Store.ORM.DB.Exec(`
INSERT INTO configurations (name, value, updated_at, created_at)
VALUES('ETH_CHAIN_ID', '2663', now(), now());`,
	).Error)
	config.GeneralConfig.Overrides.SetChainID(7853)
	err := app.Start()
	assert.Equal(t, chainlink.ErrNewChainID, err)
}

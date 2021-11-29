//go:build !windows
// +build !windows

package chainlink_test

import (
	"syscall"
	"testing"

	"github.com/onsi/gomega"
	"github.com/stretchr/testify/require"
	"go.uber.org/atomic"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
)

func TestChainlinkApplication_SignalShutdown(t *testing.T) {
	ethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	app := cltest.NewApplication(t, ethClient)
	var completed atomic.Bool
	app.Exiter = func(code int) {
		completed.Store(true)
	}

	require.NoError(t, app.Start())
	syscall.Kill(syscall.Getpid(), syscall.SIGTERM)

	gomega.NewWithT(t).Eventually(completed.Load).Should(gomega.BeTrue())
}

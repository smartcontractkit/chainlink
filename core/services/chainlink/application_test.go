//go:build !windows
// +build !windows

package chainlink_test

import (
	"syscall"
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/stretchr/testify/require"

	"github.com/onsi/gomega"
	"github.com/tevino/abool"
)

func TestChainlinkApplication_SignalShutdown(t *testing.T) {
	ethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	app := cltest.NewApplication(t, ethClient)
	completed := abool.New()
	app.Exiter = func(code int) {
		completed.Set()
	}

	require.NoError(t, app.Start())
	syscall.Kill(syscall.Getpid(), syscall.SIGTERM)

	gomega.NewWithT(t).Eventually(func() bool {
		return completed.IsSet()
	}).Should(gomega.BeTrue())
}

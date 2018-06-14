// +build !windows

package services_test

import (
	"syscall"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/tevino/abool"
)

func TestServices_ApplicationSignalShutdown(t *testing.T) {
	RegisterTestingT(t)
	app, cleanup := cltest.NewApplication()
	defer cleanup()

	completed := abool.New()
	app.Exiter = func(code int) {
		completed.Set()
	}

	app.Start()
	syscall.Kill(syscall.Getpid(), syscall.SIGTERM)

	Eventually(func() bool {
		return completed.IsSet()
	}).Should(BeTrue())
}

package services_test

import (
	"syscall"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink/internal/cltest"
)

func TestServices_ApplicationSignalShutdown(t *testing.T) {
	RegisterTestingT(t)
	app, cleanup := cltest.NewApplication()
	defer cleanup()

	var completed bool
	app.Exiter = func(code int) {
		completed = true
	}

	app.Start()
	syscall.Kill(syscall.Getpid(), syscall.SIGTERM)

	Eventually(func() bool {
		return completed
	}).Should(BeTrue())
}

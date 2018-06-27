// +build !windows

package services_test

import (
	"syscall"
	"testing"

	"github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/tevino/abool"
)

func TestServices_ApplicationSignalShutdown(t *testing.T) {
	config, cleanup := cltest.NewConfig()
	defer cleanup()
	app, _ := cltest.NewApplicationWithConfig(config)

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

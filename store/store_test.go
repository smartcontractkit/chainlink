package store_test

import (
	"syscall"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink/internal/cltest"
)

func TestGracefulShutdown(t *testing.T) {
	RegisterTestingT(t)
	store, cleanup := cltest.NewStore()
	defer cleanup()

	var completed bool
	store.Exiter = func(code int) {
		completed = true
	}

	store.Start()
	syscall.Kill(syscall.Getpid(), syscall.SIGTERM)

	Eventually(func() bool {
		return completed
	}).Should(BeTrue())
}

package store_test

import (
	"math/big"
	"syscall"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	strpkg "github.com/smartcontractkit/chainlink/store"
	"github.com/stretchr/testify/assert"
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

func TestConfigDefaults(t *testing.T) {
	config := strpkg.NewConfig()
	assert.Equal(t, uint64(0), config.ChainID)
	assert.Equal(t, *big.NewInt(20000000000), config.EthGasPriceDefault)
}

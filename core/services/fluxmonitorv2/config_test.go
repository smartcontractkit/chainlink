package fluxmonitorv2_test

import (
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/fluxmonitorv2"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	flagsContractAddress := cltest.NewAddress()

	cfg := &fluxmonitorv2.Config{
		DefaultHTTPTimeout:       time.Minute,
		FlagsContractAddress:     flagsContractAddress.Hex(),
		MinContractPayment:       assets.NewLink(1),
		EthGasLimit:              21000,
		EthMaxQueuedTransactions: 0,
	}

	t.Run("MinimumPollingInterval", func(t *testing.T) {
		assert.Equal(t, time.Minute, cfg.MinimumPollingInterval())
	})
}

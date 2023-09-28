package ccipdata

import (
	"fmt"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func TestLogPollerClient_loadDependency(t *testing.T) {
	c := &LogPollerReader{}

	someAddr := utils.RandomAddress()

	onRamp, err := c.loadOnRamp(someAddr)
	assert.NoError(t, err)
	onRamp2, err := c.loadOnRamp(someAddr)
	assert.NoError(t, err)
	// the objects should have the same pointer
	// since the second time the dependency should've been loaded from cache instead of initializing a new instance.
	assert.True(t, onRamp == onRamp2)

	offRamp, err := c.loadOffRamp(someAddr)
	assert.NoError(t, err)
	offRamp2, err := c.loadOffRamp(someAddr)
	assert.NoError(t, err)
	assert.True(t, offRamp == offRamp2)

	priceReg, err := c.loadPriceRegistry(someAddr)
	assert.NoError(t, err)
	priceReg2, err := c.loadPriceRegistry(someAddr)
	assert.NoError(t, err)
	assert.True(t, priceReg == priceReg2)
}

func Test_parseLogs(t *testing.T) {
	// generate 100 logs
	logs := make([]logpoller.Log, 100)
	for i := range logs {
		logs[i].LogIndex = int64(i + 1)
		logs[i].BlockNumber = int64(i) * 1000
		logs[i].BlockTimestamp = time.Now()
	}

	parseFn := func(log types.Log) (*uint, error) {
		// Simulate some random error
		if log.Index == 100 {
			return nil, fmt.Errorf("some error")
		}
		return &log.Index, nil
	}

	parsedEvents, err := parseLogs[uint](logs, logger.TestLogger(t), parseFn)
	assert.NoError(t, err)
	assert.Len(t, parsedEvents, 99)

	// Make sure everything is parsed according to the parse func
	for i, ev := range parsedEvents {
		assert.Equal(t, i+1, int(ev.Data))
		assert.Equal(t, int(i)*1000, int(ev.BlockNumber))
		assert.Greater(t, ev.BlockTimestamp, time.Now().Add(-time.Minute))
	}
}

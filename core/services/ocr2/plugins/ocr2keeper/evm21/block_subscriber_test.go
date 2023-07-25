package evm

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

const blockHistorySize = 128

func TestBlockSubscriber_Subscribe(t *testing.T) {
	lggr := logger.TestLogger(t)
	var hb types.HeadBroadcaster
	var lp logpoller.LogPoller

	bs := NewBlockSubscriber(hb, lp, blockHistorySize, lggr)
	subId, _, err := bs.Subscribe()
	assert.Nil(t, err)
	assert.Equal(t, subId, 1)
	subId, _, err = bs.Subscribe()
	assert.Nil(t, err)
	assert.Equal(t, subId, 2)
	subId, _, err = bs.Subscribe()
	assert.Nil(t, err)
	assert.Equal(t, subId, 3)
}

func TestBlockSubscriber_Unsubscribe(t *testing.T) {
	lggr := logger.TestLogger(t)
	var hb types.HeadBroadcaster
	var lp logpoller.LogPoller

	bs := NewBlockSubscriber(hb, lp, blockHistorySize, lggr)
	subId, _, err := bs.Subscribe()
	assert.Nil(t, err)
	assert.Equal(t, subId, 1)
	subId, _, err = bs.Subscribe()
	assert.Nil(t, err)
	assert.Equal(t, subId, 2)
	err = bs.Unsubscribe(1)
	assert.Nil(t, err)
}

func TestBlockSubscriber_Unsubscribe_Failure(t *testing.T) {
	lggr := logger.TestLogger(t)
	var hb types.HeadBroadcaster
	var lp logpoller.LogPoller

	bs := NewBlockSubscriber(hb, lp, blockHistorySize, lggr)
	err := bs.Unsubscribe(2)
	assert.Equal(t, err.Error(), "subscriber 2 does not exist")
}

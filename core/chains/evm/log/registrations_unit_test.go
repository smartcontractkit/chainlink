package log_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/core/chains/evm/log"
	lbmocks "github.com/smartcontractkit/chainlink/core/chains/evm/log/mocks"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/test-go/testify/assert"
)

func newListener(t *testing.T) *lbmocks.Listener {
	l := new(lbmocks.Listener)
	l.Test(t)
	return l
}

func TestRegistrationsUnit_PanicsOnDoubleSubscribe(t *testing.T) {
	l := newListener(t)
	l.On("JobID").Return(int32(1))

	contractAddr := cltest.NewAddress()
	opts := log.ListenerOpts{Contract: contractAddr}
	sub := log.NewSubscriber(l, opts)

	r := log.NewRegistrations(logger.TestLogger(t), cltest.FixtureChainID)
	log.RegistrationsAddSubscriber(r, sub)

	l2 := newListener(t)
	l2.On("JobID").Return(int32(1))
	opts2 := log.ListenerOpts{Contract: contractAddr}
	sub2 := log.NewSubscriber(l2, opts2)

	// Different subscriber same job ID is ok
	log.RegistrationsAddSubscriber(r, sub2)

	// Adding same subscriber twice is not ok
	assert.Panics(t, func() {
		log.RegistrationsAddSubscriber(r, sub2)
	}, "expected adding same subscription twice to panic")

	log.RegistrationsRmSubscriber(r, sub)

	// Removing subscriber twice also panics
	assert.Panics(t, func() {
		log.RegistrationsRmSubscriber(r, sub)
	}, "expected removing a subscriber twice to panic")

	// Now we can add it again
	log.RegistrationsAddSubscriber(r, sub)
}

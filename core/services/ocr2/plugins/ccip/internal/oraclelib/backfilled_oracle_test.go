package oraclelib

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	lpmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	jobmocks "github.com/smartcontractkit/chainlink/v2/core/services/job/mocks"
)

func TestBackfilledOracle(t *testing.T) {
	// First scenario: Start() fails, check that all Replay are being called.
	lp1 := lpmocks.NewLogPoller(t)
	lp2 := lpmocks.NewLogPoller(t)
	lp1.On("Replay", mock.Anything, int64(1)).Return(nil)
	lp2.On("Replay", mock.Anything, int64(2)).Return(nil)
	oracle1 := jobmocks.NewServiceCtx(t)
	oracle1.On("Start", mock.Anything).Return(errors.New("Failed to start")).Twice()
	job := NewBackfilledOracle(logger.TestLogger(t), lp1, lp2, 1, 2, oracle1)

	job.Run()
	assert.False(t, job.IsRunning())
	job.Run()
	assert.False(t, job.IsRunning())

	/// Start -> Stop -> Start
	oracle2 := jobmocks.NewServiceCtx(t)
	oracle2.On("Start", mock.Anything).Return(nil).Twice()
	oracle2.On("Close").Return(nil).Once()

	job2 := NewBackfilledOracle(logger.TestLogger(t), lp1, lp2, 1, 2, oracle2)
	job2.Run()
	assert.True(t, job2.IsRunning())
	assert.Nil(t, job2.Close())
	assert.False(t, job2.IsRunning())
	assert.Nil(t, job2.Close())
	assert.False(t, job2.IsRunning())
	job2.Run()
	assert.True(t, job2.IsRunning())

	/// Replay fails, but it starts anyway
	lp11 := lpmocks.NewLogPoller(t)
	lp12 := lpmocks.NewLogPoller(t)
	lp11.On("Replay", mock.Anything, int64(1)).Return(errors.New("Replay failed")).Once()
	lp12.On("Replay", mock.Anything, int64(2)).Return(errors.New("Replay failed")).Once()

	oracle := jobmocks.NewServiceCtx(t)
	oracle.On("Start", mock.Anything).Return(nil).Once()
	job3 := NewBackfilledOracle(logger.NullLogger, lp11, lp12, 1, 2, oracle)
	job3.Run()
	assert.True(t, job3.IsRunning())
}

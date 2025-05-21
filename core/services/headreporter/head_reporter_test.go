package headreporter

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	evmtypes "github.com/smartcontractkit/chainlink-evm/pkg/types"
	ubig "github.com/smartcontractkit/chainlink-evm/pkg/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func NewHead() evmtypes.Head {
	return evmtypes.Head{Number: 42, EVMChainID: ubig.NewI(0)}
}

func Test_HeadReporterService(t *testing.T) {
	t.Run("report everything", func(t *testing.T) {
		db := pgtest.NewSqlxDB(t)

		headReporter := NewMockHeadReporter(t)
		service := NewHeadReporterService(db, logger.TestLogger(t), headReporter)
		service.reportPeriod = time.Second
		err := service.Start(testutils.Context(t))
		require.NoError(t, err)

		var reportCalls atomic.Int32
		head := NewHead()
		headReporter.On("ReportNewHead", mock.Anything, &head).Run(func(args mock.Arguments) {
			reportCalls.Add(1)
		}).Return(nil)
		headReporter.On("ReportPeriodic", mock.Anything).Run(func(args mock.Arguments) {
			reportCalls.Add(1)
		}).Return(nil)
		service.OnNewLongestChain(testutils.Context(t), &head)

		require.Eventually(t, func() bool { return reportCalls.Load() == 2 }, 5*time.Second, 100*time.Millisecond)
	})

	t.Run("has default report period", func(t *testing.T) {
		service := NewHeadReporterService(pgtest.NewSqlxDB(t), logger.TestLogger(t), NewMockHeadReporter(t))
		assert.Equal(t, 15*time.Second, service.reportPeriod)
	})
}

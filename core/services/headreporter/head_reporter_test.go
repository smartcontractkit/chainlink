package headreporter_test

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/headreporter"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_HeadReporterService(t *testing.T) {
	t.Run("report everything", func(t *testing.T) {
		db := pgtest.NewSqlxDB(t)

		headReporter := mocks.NewHeadReporter(t)
		service := headreporter.NewHeadReporterServiceWithReporters(db, newLegacyChainContainer(t, db), logger.TestLogger(t), []headreporter.HeadReporter{headReporter}, time.Second)
		err := service.Start(testutils.Context(t))
		require.NoError(t, err)

		var reportCalls atomic.Int32
		head := newHead()
		headReporter.On("ReportNewHead", mock.Anything, &head).Run(func(args mock.Arguments) {
			reportCalls.Add(1)
		}).Return(nil)
		headReporter.On("ReportPeriodic", mock.Anything).Run(func(args mock.Arguments) {
			reportCalls.Add(1)
		}).Return(nil)
		service.OnNewLongestChain(testutils.Context(t), &head)

		require.Eventually(t, func() bool { return reportCalls.Load() == 2 }, 5*time.Second, 100*time.Millisecond)
	})
}

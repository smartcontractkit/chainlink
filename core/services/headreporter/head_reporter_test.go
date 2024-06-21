package headreporter_test

import (
	"github.com/smartcontractkit/chainlink/v2/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/headreporter"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_HeadReporterService(t *testing.T) {
	t.Run("c", func(t *testing.T) {
		db := pgtest.NewSqlxDB(t)

		headReporter := mocks.NewHeadReporter(t)
		service := headreporter.NewHeadReporterServiceWithReporters(db, newLegacyChainContainer(t, db), logger.TestLogger(t), []headreporter.HeadReporter{headReporter})
		err := service.Start(testutils.Context(t))
		require.NoError(t, err)

		head := newHead()
		headReporter.On("ReportNewHead", mock.Anything, &head).Return(nil)
		headReporter.On("ReportPeriodic", mock.Anything).Return(nil)
		service.OnNewLongestChain(testutils.Context(t), &head)

		require.Eventually(t, func() bool {
			return headReporter.AssertCalled(t, "ReportNewHead", mock.Anything, &head) && headReporter.AssertCalled(t, "ReportPeriodic", mock.Anything)
		}, time.Second, 10*time.Millisecond)
	})
}

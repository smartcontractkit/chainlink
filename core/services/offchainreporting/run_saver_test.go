package offchainreporting

import (
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/pipeline/mocks"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRunSaver(t *testing.T) {
	pipelineRunner := new(mocks.Runner)
	rr := make(chan pipeline.RunWithResults, 100)
	db := pgtest.NewSqlDB(t)
	rs := NewResultRunSaver(
		postgres.WrapDbWithSqlx(db),
		rr,
		pipelineRunner,
		make(chan struct{}),
		*logger.Default,
	)
	require.NoError(t, rs.Start())
	for i := 0; i < 100; i++ {
		pipelineRunner.On("InsertFinishedRun", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(int64(i), nil).Once()
		rr <- pipeline.RunWithResults{Run: pipeline.Run{ID: int64(i)}}
	}
	require.NoError(t, rs.Close())
}

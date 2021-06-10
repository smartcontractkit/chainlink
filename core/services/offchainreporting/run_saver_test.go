package offchainreporting

import (
	"testing"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/pipeline/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRunSaver(t *testing.T) {
	pipelineRunner := new(mocks.Runner)
	rr := make(chan pipeline.RunWithResults, 100)
	c := orm.NewConfig()
	url := c.DatabaseURL()
	db, err := gorm.Open(gormpostgres.New(gormpostgres.Config{DSN: url.String()}), &gorm.Config{})
	require.NoError(t, err)
	rs := NewResultRunSaver(
		db,
		rr,
		pipelineRunner,
		make(chan struct{}),
		*logger.Default,
	)
	require.NoError(t, rs.Start())
	for i := 0; i < 100; i++ {
		pipelineRunner.On("InsertFinishedRun", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(int64(i), nil).Once()
		rr <- pipeline.RunWithResults{Run: pipeline.Run{ID: int64(i)}}
	}
	require.NoError(t, rs.Close())
}

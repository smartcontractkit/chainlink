package offchainreporting

import (
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/pipeline/mocks"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestRunSaver(t *testing.T) {
	pipelineRunner := new(mocks.Runner)
	rr := make(chan pipeline.Run, 100)
	c := configtest.NewTestGeneralConfig(t)
	url := c.DatabaseURL()
	db, err := gorm.Open(gormpostgres.New(gormpostgres.Config{DSN: url.String()}), &gorm.Config{})
	require.NoError(t, err)
	rs := NewResultRunSaver(
		postgres.UnwrapGormDB(db),
		rr,
		pipelineRunner,
		make(chan struct{}),
		logger.TestLogger(t),
	)
	require.NoError(t, rs.Start())
	for i := 0; i < 100; i++ {
		pipelineRunner.On("InsertFinishedRun", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(int64(i), nil).Once()
		rr <- pipeline.Run{ID: int64(i)}
	}
	require.NoError(t, rs.Close())
}

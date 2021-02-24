package directrequest_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/directrequest"
	"github.com/smartcontractkit/chainlink/core/services/eth/contracts"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/log"
	"github.com/smartcontractkit/chainlink/core/services/log/mocks"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/postgres"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDelegate_ServicesForSpec(t *testing.T) {
	broadcaster := new(mocks.Broadcaster)
	runner := new(mocks.PipelineRunner)
	_, orm, cleanupDB := cltest.BootstrapThrowawayORM(t, "event_broadcaster", true)
	defer cleanupDB()

	delegate := directrequest.NewDelegate(broadcaster, runner, orm.DB)

	t.Run("Spec without DirectRequestSpec", func(t *testing.T) {
		spec := job.SpecDB{}
		_, err := delegate.ServicesForSpec(spec)
		assert.Error(t, err, "expects a *job.DirectRequestSpec to be present")
	})

	t.Run("Spec with DirectRequestSpec", func(t *testing.T) {
		spec := job.SpecDB{DirectRequestSpec: &job.DirectRequestSpec{}}
		services, err := delegate.ServicesForSpec(spec)
		require.NoError(t, err)
		assert.Len(t, services, 1)
	})
}

func TestDelegate_ServicesListenerHandleLog(t *testing.T) {
	broadcaster := new(mocks.Broadcaster)
	runner := new(mocks.PipelineRunner)
	config, orm, cleanupDB := cltest.BootstrapThrowawayORM(t, "event_broadcaster", true)
	defer cleanupDB()
	db := orm.DB

	eventBroadcaster := postgres.NewEventBroadcaster(config.DatabaseURL(), 0, 0)
	eventBroadcaster.Start()
	defer eventBroadcaster.Stop()

	delegate := directrequest.NewDelegate(broadcaster, runner, db)

	spec := job.SpecDB{DirectRequestSpec: &job.DirectRequestSpec{}}
	services, err := delegate.ServicesForSpec(spec)
	require.NoError(t, err)
	assert.Len(t, services, 1)
	service := services[0]

	var listener log.Listener
	broadcaster.On("Register", mock.Anything, mock.Anything).Return(true).Run(func(args mock.Arguments) {
		listener = args.Get(1).(log.Listener)
	})
	broadcaster.On("Unregister", mock.Anything, mock.Anything).Return(nil)

	t.Run("Log is an OracleRequest", func(t *testing.T) {
		log := new(mocks.Broadcast)
		defer log.AssertExpectations(t)

		log.On("WasAlreadyConsumed").Return(false, nil)
		logOracleRequest := contracts.LogOracleRequest{
			CancelExpiration: big.NewInt(0),
		}
		log.On("DecodedLog").Return(&logOracleRequest)
		log.On("MarkConsumed").Return(nil)

		runner.On("CreateRun", mock.Anything, mock.Anything, mock.Anything).Return(int64(0), nil)

		err := service.Start()
		require.NoError(t, err)

		listener.HandleLog(log, nil)

		service.Close()
		broadcaster.AssertExpectations(t)
		runner.AssertExpectations(t)
	})

	t.Run("Log is a CancelOracleRequest", func(t *testing.T) {
		orm, _, cleanup := cltest.NewPipelineORM(t, config, db)
		defer cleanup()

		jobORM := job.NewORM(db, config.Config, orm, eventBroadcaster, &postgres.NullAdvisoryLocker{})
		defer jobORM.Close()

		spec := factoryJobSpec(t)
		err = jobORM.CreateJob(context.Background(), spec, spec.Pipeline)
		require.NoError(t, err)

		// Create one run with a matching request ID ...
		meta := make(map[string]interface{})
		request := contracts.OracleRequest{
			RequestID: cltest.NewHash(),
		}
		meta["oracleRequest"] = request.ToMap()
		_, err = orm.CreateRun(context.Background(), spec.ID, meta)
		require.NoError(t, err)
		// And one without
		_, err = orm.CreateRun(context.Background(), spec.ID, nil)
		require.NoError(t, err)

		log := new(mocks.Broadcast)
		defer log.AssertExpectations(t)

		log.On("WasAlreadyConsumed").Return(false, nil)
		logCancelOracleRequest := contracts.LogCancelOracleRequest{
			RequestID: request.RequestID,
		}
		log.On("DecodedLog").Return(&logCancelOracleRequest)
		log.On("MarkConsumed").Return(nil)

		err = service.Start()
		require.NoError(t, err)

		listener.HandleLog(log, nil)

		// Only one should remain
		_, count, err := jobORM.PipelineRunsByJobID(spec.ID, 0, 2)
		require.NoError(t, err)
		assert.Equal(t, 1, count)

		service.Close()
		broadcaster.AssertExpectations(t)
		runner.AssertExpectations(t)
	})
}

func factoryJobSpec(t *testing.T) *job.SpecDB {
	t.Helper()
	spec := &job.SpecDB{
		Type:              job.DirectRequest,
		SchemaVersion:     1,
		DirectRequestSpec: &job.DirectRequestSpec{},
		Pipeline:          *pipeline.NewTaskDAG(),
	}
	return spec
}

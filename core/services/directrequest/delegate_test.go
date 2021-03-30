package directrequest_test

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/oracle_wrapper"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/services/directrequest"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/log"
	log_mocks "github.com/smartcontractkit/chainlink/core/services/log/mocks"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	pipeline_mocks "github.com/smartcontractkit/chainlink/core/services/pipeline/mocks"
	"github.com/smartcontractkit/chainlink/core/services/postgres"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDelegate_ServicesForSpec(t *testing.T) {
	gethClient := new(mocks.Client)
	broadcaster := new(log_mocks.Broadcaster)
	runner := new(pipeline_mocks.Runner)
	_, orm, cleanupDB := cltest.BootstrapThrowawayORM(t, "event_broadcaster", true)
	defer cleanupDB()

	delegate := directrequest.NewDelegate(broadcaster, runner, gethClient, orm.DB)

	t.Run("Spec without DirectRequestSpec", func(t *testing.T) {
		spec := job.Job{}
		_, err := delegate.ServicesForSpec(spec)
		assert.Error(t, err, "expects a *job.DirectRequestSpec to be present")
	})

	t.Run("Spec with DirectRequestSpec", func(t *testing.T) {
		spec := job.Job{DirectRequestSpec: &job.DirectRequestSpec{}}
		services, err := delegate.ServicesForSpec(spec)
		require.NoError(t, err)
		assert.Len(t, services, 1)
	})
}

func TestDelegate_ServicesListenerHandleLog(t *testing.T) {
	gethClient := new(mocks.Client)
	broadcaster := new(log_mocks.Broadcaster)
	runner := new(pipeline_mocks.Runner)

	config, oldORM, cleanupDB := cltest.BootstrapThrowawayORM(t, "delegate_services_listener_handlelog", true, true)
	defer cleanupDB()
	db := oldORM.DB

	orm, eventBroadcaster, cleanup := cltest.NewPipelineORM(t, config, db)
	defer cleanup()

	jobORM := job.NewORM(db, config.Config, orm, eventBroadcaster, &postgres.NullAdvisoryLocker{})
	defer jobORM.Close()

	delegate := directrequest.NewDelegate(broadcaster, runner, gethClient, db)

	spec := factoryJobSpec(t)
	err := jobORM.CreateJob(context.Background(), spec, spec.Pipeline)
	require.NoError(t, err)
	services, err := delegate.ServicesForSpec(*spec)
	require.NoError(t, err)
	assert.Len(t, services, 1)
	service := services[0]

	var listener log.Listener
	broadcaster.On("Register", mock.Anything, mock.Anything).Return(true, nil).Run(func(args mock.Arguments) {
		listener = args.Get(0).(log.Listener)
	})

	t.Run("Log is an OracleRequest", func(t *testing.T) {
		log := new(log_mocks.Broadcast)
		defer log.AssertExpectations(t)

		log.On("WasAlreadyConsumed").Return(false, nil)
		logOracleRequest := oracle_wrapper.OracleOracleRequest{
			CancelExpiration: big.NewInt(0),
		}
		log.On("DecodedLog").Return(&logOracleRequest)
		log.On("MarkConsumed").Return(nil)

		runner.On("CreateRun", mock.Anything, mock.Anything, mock.Anything).Return(int64(0), nil)

		err := service.Start()
		require.NoError(t, err)

		listener.HandleLog(log)

		service.Close()
		broadcaster.AssertExpectations(t)
		runner.AssertExpectations(t)
	})

	t.Run("Log is a CancelOracleRequest", func(t *testing.T) {

		// Create one run with a matching request ID ...
		meta := make(map[string]interface{})
		request := oracle_wrapper.OracleOracleRequest{
			RequestId: cltest.NewHash(),
		}
		meta["oracleRequest"] = map[string]string{"requestId": fmt.Sprintf("0x%x", request.RequestId)}
		_, err = orm.CreateRun(context.Background(), spec.ID, meta)
		require.NoError(t, err)
		// And one without
		_, err = orm.CreateRun(context.Background(), spec.ID, nil)
		require.NoError(t, err)

		log := new(log_mocks.Broadcast)
		defer log.AssertExpectations(t)

		log.On("WasAlreadyConsumed").Return(false, nil)
		logCancelOracleRequest := oracle_wrapper.OracleCancelOracleRequest{
			RequestId: request.RequestId,
		}
		log.On("DecodedLog").Return(&logCancelOracleRequest)
		log.On("MarkConsumed").Return(nil)

		err = service.Start()
		require.NoError(t, err)

		listener.HandleLog(log)

		// Only one should remain
		_, count, err := jobORM.PipelineRunsByJobID(spec.ID, 0, 2)
		require.NoError(t, err)
		assert.Equal(t, 1, count)

		service.Close()
		broadcaster.AssertExpectations(t)
		runner.AssertExpectations(t)
	})
}

func factoryJobSpec(t *testing.T) *job.Job {
	t.Helper()
	spec := &job.Job{
		Type:              job.DirectRequest,
		SchemaVersion:     1,
		DirectRequestSpec: &job.DirectRequestSpec{},
		Pipeline:          *pipeline.NewTaskDAG(),
	}
	return spec
}

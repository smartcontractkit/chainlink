package directrequest_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gofrs/uuid"
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
	"github.com/smartcontractkit/chainlink/core/store/models"

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

	delegate := directrequest.NewDelegate(broadcaster, runner, nil, gethClient, orm.DB)

	t.Run("Spec without DirectRequestSpec", func(t *testing.T) {
		spec := job.Job{}
		_, err := delegate.ServicesForSpec(spec)
		assert.Error(t, err, "expects a *job.DirectRequestSpec to be present")
	})

	t.Run("Spec with DirectRequestSpec", func(t *testing.T) {
		spec := job.Job{DirectRequestSpec: &job.DirectRequestSpec{}, PipelineSpec: &pipeline.Spec{}}
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

	delegate := directrequest.NewDelegate(broadcaster, runner, orm, gethClient, db)

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
		log.On("RawLog").Return(models.Log{
			Topics: []common.Hash{
				common.Hash{},
				spec.DirectRequestSpec.OnChainJobSpecID,
			},
		})
		log.On("DecodedLog").Return(&logOracleRequest)
		log.On("MarkConsumed").Return(nil)

		runBeganAwaiter := cltest.NewAwaiter()
		runner.On("ExecuteAndInsertNewRun", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			runBeganAwaiter.ItHappened()
		}).Once().Return(int64(0), pipeline.FinalResult{}, nil)

		err = service.Start()
		require.NoError(t, err)

		listener.HandleLog(log)
		runBeganAwaiter.AwaitOrFail(t, 5*time.Second)

		service.Close()
		broadcaster.AssertExpectations(t)
		runner.AssertExpectations(t)
	})

	t.Run("Log has wrong jobID", func(t *testing.T) {
		log := new(log_mocks.Broadcast)
		defer log.AssertExpectations(t)

		log.On("WasAlreadyConsumed").Return(false, nil)
		log.On("RawLog").Return(models.Log{
			Topics: []common.Hash{common.Hash{}, common.Hash{}},
		})

		err = service.Start()
		require.NoError(t, err)

		listener.HandleLog(log)

		service.Close()
		broadcaster.AssertExpectations(t)
		runner.AssertExpectations(t)
	})

	t.Run("Log is a CancelOracleRequest with no matching run", func(t *testing.T) {
		log := new(log_mocks.Broadcast)
		defer log.AssertExpectations(t)

		log.On("WasAlreadyConsumed").Return(false, nil)
		logCancelOracleRequest := oracle_wrapper.OracleCancelOracleRequest{RequestId: spec.DirectRequestSpec.OnChainJobSpecID}
		log.On("RawLog").Return(models.Log{
			Topics: []common.Hash{
				common.Hash{},
				spec.DirectRequestSpec.OnChainJobSpecID,
			},
		})
		log.On("DecodedLog").Return(&logCancelOracleRequest)
		log.On("MarkConsumed").Return(nil)

		err = service.Start()
		require.NoError(t, err)

		listener.HandleLog(log)

		service.Close()
		broadcaster.AssertExpectations(t)
		runner.AssertExpectations(t)
	})

	t.Run("Log is a CancelOracleRequest with a matching run", func(t *testing.T) {
		runLog := new(log_mocks.Broadcast)

		runLog.On("WasAlreadyConsumed").Return(false, nil)
		logOracleRequest := oracle_wrapper.OracleOracleRequest{
			CancelExpiration: big.NewInt(0),
			RequestId:        spec.DirectRequestSpec.OnChainJobSpecID,
		}
		runLog.On("RawLog").Return(models.Log{
			Topics: []common.Hash{
				common.Hash{},
				spec.DirectRequestSpec.OnChainJobSpecID,
			},
		})
		runLog.On("DecodedLog").Return(&logOracleRequest)
		runLog.On("MarkConsumed").Return(nil)

		cancelLog := new(log_mocks.Broadcast)

		cancelLog.On("WasAlreadyConsumed").Return(false, nil)
		logCancelOracleRequest := oracle_wrapper.OracleCancelOracleRequest{RequestId: spec.DirectRequestSpec.OnChainJobSpecID}
		cancelLog.On("RawLog").Return(models.Log{
			Topics: []common.Hash{
				common.Hash{},
				spec.DirectRequestSpec.OnChainJobSpecID,
			},
		})
		cancelLog.On("DecodedLog").Return(&logCancelOracleRequest)
		cancelLog.On("MarkConsumed").Return(nil)

		err = service.Start()
		require.NoError(t, err)

		timeout := 5 * time.Second
		runBeganAwaiter := cltest.NewAwaiter()
		runCancelledAwaiter := cltest.NewAwaiter()
		runner.On("ExecuteAndInsertNewRun", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			runBeganAwaiter.ItHappened()
			ctx := args[0].(context.Context)
			select {
			case <-time.After(timeout):
				t.Fatalf("Timed out waiting for Run to be canceled (%v)", timeout)
			case <-ctx.Done():
				runCancelledAwaiter.ItHappened()
			}
		}).Once().Return(int64(0), pipeline.FinalResult{}, nil)
		listener.HandleLog(runLog)
		runBeganAwaiter.AwaitOrFail(t, timeout)
		runLog.AssertExpectations(t)

		listener.HandleLog(cancelLog)
		runCancelledAwaiter.AwaitOrFail(t, timeout)
		cancelLog.AssertExpectations(t)

		service.Close()
		broadcaster.AssertExpectations(t)
		runner.AssertExpectations(t)
	})
}

func factoryJobSpec(t *testing.T) *job.Job {
	t.Helper()
	drs := &job.DirectRequestSpec{}
	onChainJobSpecID, err := uuid.NewV4()
	require.NoError(t, err)
	copy(drs.OnChainJobSpecID[:], onChainJobSpecID[:])
	spec := &job.Job{
		Type:              job.DirectRequest,
		SchemaVersion:     1,
		DirectRequestSpec: drs,
		Pipeline:          *pipeline.NewTaskDAG(),
		PipelineSpec:      &pipeline.Spec{},
	}
	return spec
}

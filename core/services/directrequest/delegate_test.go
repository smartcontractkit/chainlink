package directrequest_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink/core/assets"
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
	"gorm.io/gorm"
)

func TestDelegate_ServicesForSpec(t *testing.T) {
	ethClient := new(mocks.Client)
	broadcaster := new(log_mocks.Broadcaster)
	runner := new(pipeline_mocks.Runner)

	_, orm, cleanupDB := cltest.BootstrapThrowawayORM(t, "event_broadcaster", true)
	defer cleanupDB()

	config := testConfig{
		minIncomingConfirmations: 1,
	}
	delegate := directrequest.NewDelegate(broadcaster, runner, nil, ethClient, orm.DB, config)

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

type DirectRequestUniverse struct {
	spec           *job.Job
	runner         *pipeline_mocks.Runner
	service        job.Service
	jobORM         job.ORM
	listener       log.Listener
	logBroadcaster *log_mocks.Broadcaster
	cleanup        func()
}

func NewDirectRequestUniverseWithConfig(t *testing.T, drConfig testConfig) *DirectRequestUniverse {
	gethClient := new(mocks.Client)
	broadcaster := new(log_mocks.Broadcaster)
	runner := new(pipeline_mocks.Runner)

	config, oldORM, cleanupDB := cltest.BootstrapThrowawayORM(t, "delegate_services_listener_handlelog", true, true)
	db := oldORM.DB

	orm, eventBroadcaster, cleanupPipeline := cltest.NewPipelineORM(t, config, db)

	jobORM := job.NewORM(db, config.Config, orm, eventBroadcaster, &postgres.NullAdvisoryLocker{})

	cleanup := func() {
		cleanupDB()
		cleanupPipeline()
		jobORM.Close()
	}

	delegate := directrequest.NewDelegate(broadcaster, runner, orm, gethClient, db, drConfig)

	spec := cltest.MakeDirectRequestJobSpec(t)
	err := jobORM.CreateJob(context.Background(), spec, spec.Pipeline)
	require.NoError(t, err)
	serviceArray, err := delegate.ServicesForSpec(*spec)
	require.NoError(t, err)
	assert.Len(t, serviceArray, 1)
	service := serviceArray[0]

	uni := &DirectRequestUniverse{
		spec:           spec,
		runner:         runner,
		service:        service,
		jobORM:         jobORM,
		listener:       nil,
		logBroadcaster: broadcaster,
		cleanup:        cleanup,
	}

	broadcaster.On("Register", mock.Anything, mock.Anything).Return(func() {}).Run(func(args mock.Arguments) {
		uni.listener = args.Get(0).(log.Listener)
	})

	return uni
}

func NewDirectRequestUniverse(t *testing.T) *DirectRequestUniverse {
	drConfig := testConfig{
		minIncomingConfirmations: 1,
	}
	return NewDirectRequestUniverseWithConfig(t, drConfig)
}

func (uni *DirectRequestUniverse) Cleanup() {
	uni.cleanup()
}

func TestDelegate_ServicesListenerHandleLog(t *testing.T) {

	t.Run("Log is an OracleRequest", func(t *testing.T) {
		uni := NewDirectRequestUniverse(t)
		defer uni.Cleanup()

		log := new(log_mocks.Broadcast)
		defer log.AssertExpectations(t)

		uni.logBroadcaster.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)
		logOracleRequest := oracle_wrapper.OracleOracleRequest{
			CancelExpiration: big.NewInt(0),
		}
		log.On("RawLog").Return(types.Log{
			Topics: []common.Hash{
				common.Hash{},
				uni.spec.DirectRequestSpec.OnChainJobSpecID.Hash(),
			},
		})
		log.On("DecodedLog").Return(&logOracleRequest)
		uni.logBroadcaster.On("MarkConsumed", mock.Anything, mock.Anything).Return(nil)

		uni.runner.On("ExecuteRun", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		}).Once().Return(pipeline.Run{}, pipeline.TaskRunResults{}, nil)

		runBeganAwaiter := cltest.NewAwaiter()
		uni.runner.On("InsertFinishedRun", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			runBeganAwaiter.ItHappened()
		}).Once().Return(int64(1), nil)

		err := uni.service.Start()
		require.NoError(t, err)

		// check if the job exists under the correct ID
		drJob, jErr := uni.jobORM.FindJob(uni.listener.JobIDV2())
		require.NoError(t, jErr)
		require.Equal(t, drJob.ID, uni.listener.JobIDV2())
		require.NotNil(t, drJob.DirectRequestSpec)

		uni.listener.HandleLog(log)

		runBeganAwaiter.AwaitOrFail(t, 5*time.Second)

		uni.service.Close()
		uni.logBroadcaster.AssertExpectations(t)
		uni.runner.AssertExpectations(t)
	})

	t.Run("Log is not consumed, as it's too young", func(t *testing.T) {
		uni := NewDirectRequestUniverse(t)
		defer uni.Cleanup()

		log := new(log_mocks.Broadcast)

		uni.logBroadcaster.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil).Maybe()
		logOracleRequest := oracle_wrapper.OracleOracleRequest{
			CancelExpiration: big.NewInt(0),
		}
		log.On("RawLog").Return(types.Log{
			Topics: []common.Hash{
				common.Hash{},
				uni.spec.DirectRequestSpec.OnChainJobSpecID.Hash(),
			},
			BlockNumber: 0,
		}).Maybe()
		log.On("DecodedLog").Return(&logOracleRequest).Maybe()
		uni.logBroadcaster.On("MarkConsumed", mock.Anything, mock.Anything).Return(nil).Maybe()

		err := uni.service.Start()
		require.NoError(t, err)

		log.AssertExpectations(t)

		uni.listener.HandleLog(log)

		uni.logBroadcaster.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)

		uni.runner.On("ExecuteRun", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		}).Once().Return(pipeline.Run{}, pipeline.TaskRunResults{}, nil)

		runBeganAwaiter := cltest.NewAwaiter()
		uni.runner.On("InsertFinishedRun", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			runBeganAwaiter.ItHappened()
		}).Once().Return(int64(1), nil)

		// but should after this one, as the head Number is larger
		runBeganAwaiter.AwaitOrFail(t, 5*time.Second)
		cltest.EventuallyExpectationsMet(t, log, 3*time.Second, 100*time.Millisecond)

		uni.service.Close()
		uni.logBroadcaster.AssertExpectations(t)
		uni.runner.AssertExpectations(t)
	})

	t.Run("Log has wrong jobID", func(t *testing.T) {
		uni := NewDirectRequestUniverse(t)
		defer uni.Cleanup()

		log := new(log_mocks.Broadcast)

		uni.logBroadcaster.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)
		log.On("RawLog").Return(types.Log{
			Topics: []common.Hash{common.Hash{}, common.Hash{}},
		})

		err := uni.service.Start()
		require.NoError(t, err)

		uni.listener.HandleLog(log)

		cltest.EventuallyExpectationsMet(t, uni.logBroadcaster, 3*time.Second, 100*time.Millisecond)
		cltest.EventuallyExpectationsMet(t, uni.runner, 3*time.Second, 100*time.Millisecond)
		cltest.EventuallyExpectationsMet(t, log, 3*time.Second, 100*time.Millisecond)

		uni.service.Close()
	})

	t.Run("Log is a CancelOracleRequest with no matching run", func(t *testing.T) {
		uni := NewDirectRequestUniverse(t)
		defer uni.Cleanup()

		log := new(log_mocks.Broadcast)

		uni.logBroadcaster.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)
		logCancelOracleRequest := oracle_wrapper.OracleCancelOracleRequest{RequestId: uni.spec.DirectRequestSpec.OnChainJobSpecID.Hash()}
		log.On("RawLog").Return(types.Log{
			Topics: []common.Hash{
				common.Hash{},
				uni.spec.DirectRequestSpec.OnChainJobSpecID.Hash(),
			},
		})
		log.On("DecodedLog").Return(&logCancelOracleRequest)
		uni.logBroadcaster.On("MarkConsumed", mock.Anything, mock.Anything).Return(nil)

		err := uni.service.Start()
		require.NoError(t, err)

		uni.listener.HandleLog(log)

		cltest.EventuallyExpectationsMet(t, uni.logBroadcaster, 3*time.Second, 100*time.Millisecond)
		cltest.EventuallyExpectationsMet(t, uni.runner, 3*time.Second, 100*time.Millisecond)
		cltest.EventuallyExpectationsMet(t, log, 3*time.Second, 100*time.Millisecond)

		uni.service.Close()
	})

	t.Run("Log is a CancelOracleRequest with a matching run", func(t *testing.T) {
		uni := NewDirectRequestUniverse(t)
		defer uni.Cleanup()

		runLog := new(log_mocks.Broadcast)

		uni.logBroadcaster.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)
		logOracleRequest := oracle_wrapper.OracleOracleRequest{
			CancelExpiration: big.NewInt(0),
			RequestId:        uni.spec.DirectRequestSpec.OnChainJobSpecID.Hash(),
		}
		runLog.On("RawLog").Return(types.Log{
			Topics: []common.Hash{
				common.Hash{},
				uni.spec.DirectRequestSpec.OnChainJobSpecID.Hash(),
			},
		})
		runLog.On("DecodedLog").Return(&logOracleRequest)
		uni.logBroadcaster.On("MarkConsumed", mock.Anything, mock.Anything).Return(nil)

		cancelLog := new(log_mocks.Broadcast)

		uni.logBroadcaster.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)
		logCancelOracleRequest := oracle_wrapper.OracleCancelOracleRequest{RequestId: uni.spec.DirectRequestSpec.OnChainJobSpecID.Hash()}
		cancelLog.On("RawLog").Return(types.Log{
			Topics: []common.Hash{
				common.Hash{},
				uni.spec.DirectRequestSpec.OnChainJobSpecID.Hash(),
			},
		})
		cancelLog.On("DecodedLog").Return(&logCancelOracleRequest)
		uni.logBroadcaster.On("MarkConsumed", mock.Anything, mock.Anything).Return(nil)

		err := uni.service.Start()
		require.NoError(t, err)

		timeout := 5 * time.Second
		runBeganAwaiter := cltest.NewAwaiter()
		runCancelledAwaiter := cltest.NewAwaiter()
		uni.runner.On("ExecuteRun", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		}).Once().Return(pipeline.Run{}, pipeline.TaskRunResults{}, nil)
		uni.runner.On("InsertFinishedRun", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			runBeganAwaiter.ItHappened()
			db := args[0].(*gorm.DB)
			select {
			case <-time.After(timeout):
				t.Fatalf("Timed out waiting for Run to be canceled (%v)", timeout)
			case <-db.Statement.Context.Done():
				runCancelledAwaiter.ItHappened()
			}
		}).Once().Return(int64(0), nil)
		uni.listener.HandleLog(runLog)

		runBeganAwaiter.AwaitOrFail(t, timeout)
		runLog.AssertExpectations(t)

		uni.listener.HandleLog(cancelLog)

		runCancelledAwaiter.AwaitOrFail(t, timeout)
		cancelLog.AssertExpectations(t)

		uni.service.Close()
		uni.logBroadcaster.AssertExpectations(t)
		uni.runner.AssertExpectations(t)
	})

	t.Run("Log has sufficient funds", func(t *testing.T) {
		drConfig := testConfig{
			minIncomingConfirmations: 1,
			minimumContractPayment:   assets.NewLink(100),
		}
		uni := NewDirectRequestUniverseWithConfig(t, drConfig)
		defer uni.Cleanup()

		log := new(log_mocks.Broadcast)
		defer log.AssertExpectations(t)

		uni.logBroadcaster.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)
		logOracleRequest := oracle_wrapper.OracleOracleRequest{
			CancelExpiration: big.NewInt(0),
			Payment:          big.NewInt(100),
		}
		log.On("RawLog").Return(types.Log{
			Topics: []common.Hash{
				common.Hash{},
				uni.spec.DirectRequestSpec.OnChainJobSpecID.Hash(),
			},
		})
		log.On("DecodedLog").Return(&logOracleRequest)
		uni.logBroadcaster.On("MarkConsumed", mock.Anything, mock.Anything).Return(nil)

		uni.runner.On("ExecuteRun", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		}).Once().Return(pipeline.Run{}, pipeline.TaskRunResults{}, nil)

		runBeganAwaiter := cltest.NewAwaiter()
		uni.runner.On("InsertFinishedRun", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			runBeganAwaiter.ItHappened()
		}).Once().Return(int64(1), nil)

		err := uni.service.Start()
		require.NoError(t, err)

		// check if the job exists under the correct ID
		drJob, jErr := uni.jobORM.FindJob(uni.listener.JobIDV2())
		require.NoError(t, jErr)
		require.Equal(t, drJob.ID, uni.listener.JobIDV2())
		require.NotNil(t, drJob.DirectRequestSpec)

		uni.listener.HandleLog(log)

		runBeganAwaiter.AwaitOrFail(t, 5*time.Second)

		uni.service.Close()
		uni.logBroadcaster.AssertExpectations(t)
		uni.runner.AssertExpectations(t)
	})

	t.Run("Log has insufficient funds", func(t *testing.T) {
		drConfig := testConfig{
			minIncomingConfirmations: 1,
			minimumContractPayment:   assets.NewLink(100),
		}
		uni := NewDirectRequestUniverseWithConfig(t, drConfig)
		defer uni.Cleanup()

		log := new(log_mocks.Broadcast)
		defer log.AssertExpectations(t)

		uni.logBroadcaster.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)
		logOracleRequest := oracle_wrapper.OracleOracleRequest{
			CancelExpiration: big.NewInt(0),
			Payment:          big.NewInt(99),
		}
		log.On("RawLog").Return(types.Log{
			Topics: []common.Hash{
				common.Hash{},
				uni.spec.DirectRequestSpec.OnChainJobSpecID.Hash(),
			},
		})
		log.On("DecodedLog").Return(&logOracleRequest)
		markConsumedLogAwaiter := cltest.NewAwaiter()
		uni.logBroadcaster.On("MarkConsumed", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			markConsumedLogAwaiter.ItHappened()
		}).Return(nil)

		err := uni.service.Start()
		require.NoError(t, err)

		uni.listener.HandleLog(log)

		markConsumedLogAwaiter.AwaitOrFail(t, 5*time.Second)

		uni.service.Close()
		uni.logBroadcaster.AssertExpectations(t)
		uni.runner.AssertExpectations(t)
	})
}

type testConfig struct {
	minIncomingConfirmations uint32
	minimumContractPayment   *assets.Link
}

func (c testConfig) MinIncomingConfirmations() uint32 {
	return c.minIncomingConfirmations
}

func (c testConfig) MinimumContractPayment() *assets.Link {
	return c.minimumContractPayment
}

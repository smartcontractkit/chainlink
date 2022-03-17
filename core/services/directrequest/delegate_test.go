package directrequest_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/chains/evm/log"
	log_mocks "github.com/smartcontractkit/chainlink/core/chains/evm/log/mocks"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/operator_wrapper"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/directrequest"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	pipeline_mocks "github.com/smartcontractkit/chainlink/core/services/pipeline/mocks"
)

func TestDelegate_ServicesForSpec(t *testing.T) {
	ethClient := cltest.NewEthClientMockWithDefaultChain(t)
	runner := new(pipeline_mocks.Runner)
	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewTestGeneralConfig(t)
	cfg.Overrides.GlobalMinIncomingConfirmations = null.IntFrom(1)
	cc := evmtest.NewChainSet(t, evmtest.TestChainOpts{DB: db, GeneralConfig: cfg, Client: ethClient})

	lggr := logger.TestLogger(t)
	delegate := directrequest.NewDelegate(lggr, runner, nil, cc)

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
	service        job.ServiceCtx
	jobORM         job.ORM
	listener       log.Listener
	logBroadcaster *log_mocks.Broadcaster
	cleanup        func()
}

func NewDirectRequestUniverseWithConfig(t *testing.T, cfg *configtest.TestGeneralConfig, specF func(spec *job.Job)) *DirectRequestUniverse {
	ethClient := cltest.NewEthClientMockWithDefaultChain(t)
	broadcaster := new(log_mocks.Broadcaster)
	broadcaster.Test(t)
	runner := new(pipeline_mocks.Runner)
	broadcaster.On("AddDependents", 1)

	db := pgtest.NewSqlxDB(t)
	cc := evmtest.NewChainSet(t, evmtest.TestChainOpts{DB: db, GeneralConfig: cfg, Client: ethClient, LogBroadcaster: broadcaster})
	lggr := logger.TestLogger(t)
	orm := pipeline.NewORM(db, lggr, cfg)

	keyStore := cltest.NewKeyStore(t, db, cfg)
	jobORM := job.NewORM(db, cc, orm, keyStore, lggr, cfg)
	delegate := directrequest.NewDelegate(lggr, runner, orm, cc)

	jb := cltest.MakeDirectRequestJobSpec(t)
	jb.ExternalJobID = uuid.NewV4()
	if specF != nil {
		specF(jb)
	}
	err := jobORM.CreateJob(jb)
	require.NoError(t, err)
	serviceArray, err := delegate.ServicesForSpec(*jb)
	require.NoError(t, err)
	assert.Len(t, serviceArray, 1)
	service := serviceArray[0]

	uni := &DirectRequestUniverse{
		spec:           jb,
		runner:         runner,
		service:        service,
		jobORM:         jobORM,
		listener:       nil,
		logBroadcaster: broadcaster,
		cleanup:        func() { jobORM.Close() },
	}

	broadcaster.On("Register", mock.Anything, mock.Anything).Return(func() {}).Run(func(args mock.Arguments) {
		uni.listener = args.Get(0).(log.Listener)
	})

	return uni
}

func NewDirectRequestUniverse(t *testing.T) *DirectRequestUniverse {
	cfg := configtest.NewTestGeneralConfig(t)
	cfg.Overrides.GlobalMinIncomingConfirmations = null.IntFrom(1)
	return NewDirectRequestUniverseWithConfig(t, cfg, nil)
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
		logOracleRequest := operator_wrapper.OperatorOracleRequest{
			CancelExpiration: big.NewInt(0),
		}
		log.On("RawLog").Return(types.Log{
			Topics: []common.Hash{
				{},
				uni.spec.ExternalIDEncodeStringToTopic(),
			},
		})
		log.On("DecodedLog").Return(&logOracleRequest)
		uni.logBroadcaster.On("MarkConsumed", mock.Anything, mock.Anything).Return(nil)

		runBeganAwaiter := cltest.NewAwaiter()
		uni.runner.On("Run", mock.Anything, mock.AnythingOfType("*pipeline.Run"), mock.Anything, mock.Anything, mock.Anything).
			Return(false, nil).
			Run(func(args mock.Arguments) {
				runBeganAwaiter.ItHappened()
				fn := args.Get(4).(func(pg.Queryer) error)
				fn(nil)
			}).Once()

		err := uni.service.Start(testutils.Context(t))
		require.NoError(t, err)

		require.NotNil(t, uni.listener, "listener was nil; expected broadcaster.Register to have been called")
		// check if the job exists under the correct ID
		drJob, jErr := uni.jobORM.FindJob(context.Background(), uni.listener.JobID())
		require.NoError(t, jErr)
		require.Equal(t, drJob.ID, uni.listener.JobID())
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
		logOracleRequest := operator_wrapper.OperatorOracleRequest{
			CancelExpiration: big.NewInt(0),
		}
		log.On("RawLog").Return(types.Log{
			Topics: []common.Hash{
				{},
				uni.spec.ExternalIDEncodeStringToTopic(),
			},
			BlockNumber: 0,
		}).Maybe()
		log.On("DecodedLog").Return(&logOracleRequest).Maybe()
		uni.logBroadcaster.On("MarkConsumed", mock.Anything, mock.Anything).Return(nil).Maybe()

		err := uni.service.Start(testutils.Context(t))
		require.NoError(t, err)

		log.AssertExpectations(t)

		uni.listener.HandleLog(log)

		uni.logBroadcaster.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)

		runBeganAwaiter := cltest.NewAwaiter()
		uni.runner.On("Run", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Run(func(args mock.Arguments) {
				runBeganAwaiter.ItHappened()
				fn := args.Get(4).(func(pg.Queryer) error)
				fn(nil)
			}).Once().Return(false, nil)

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
		uni.logBroadcaster.On("MarkConsumed", mock.Anything, mock.Anything).Return(nil)

		logCancelOracleRequest := operator_wrapper.OperatorCancelOracleRequest{RequestId: uni.spec.ExternalIDEncodeStringToTopic()}
		log.On("DecodedLog").Return(&logCancelOracleRequest)
		log.On("RawLog").Return(types.Log{
			Topics: []common.Hash{{}, {}},
		})

		err := uni.service.Start(testutils.Context(t))
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
		logCancelOracleRequest := operator_wrapper.OperatorCancelOracleRequest{RequestId: uni.spec.ExternalIDEncodeStringToTopic()}
		log.On("RawLog").Return(types.Log{
			Topics: []common.Hash{
				{},
				uni.spec.ExternalIDEncodeStringToTopic(),
			},
		})
		log.On("DecodedLog").Return(&logCancelOracleRequest)
		uni.logBroadcaster.On("MarkConsumed", mock.Anything, mock.Anything).Return(nil)

		err := uni.service.Start(testutils.Context(t))
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
		logOracleRequest := operator_wrapper.OperatorOracleRequest{
			CancelExpiration: big.NewInt(0),
			RequestId:        uni.spec.ExternalIDEncodeStringToTopic(),
		}
		runLog.On("RawLog").Return(types.Log{
			Topics: []common.Hash{
				{},
				uni.spec.ExternalIDEncodeStringToTopic(),
			},
		})
		runLog.On("DecodedLog").Return(&logOracleRequest)
		uni.logBroadcaster.On("MarkConsumed", mock.Anything, mock.Anything).Return(nil)

		cancelLog := new(log_mocks.Broadcast)

		uni.logBroadcaster.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)
		logCancelOracleRequest := operator_wrapper.OperatorCancelOracleRequest{RequestId: uni.spec.ExternalIDEncodeStringToTopic()}
		cancelLog.On("RawLog").Return(types.Log{
			Topics: []common.Hash{
				{},
				uni.spec.ExternalIDEncodeStringToTopic(),
			},
		})
		cancelLog.On("DecodedLog").Return(&logCancelOracleRequest)
		uni.logBroadcaster.On("MarkConsumed", mock.Anything, mock.Anything).Return(nil)

		err := uni.service.Start(testutils.Context(t))
		require.NoError(t, err)

		timeout := 5 * time.Second
		runBeganAwaiter := cltest.NewAwaiter()
		runCancelledAwaiter := cltest.NewAwaiter()
		uni.runner.On("Run", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			runBeganAwaiter.ItHappened()
			ctx := args[0].(context.Context)
			select {
			case <-time.After(timeout):
				t.Fatalf("Timed out waiting for Run to be canceled (%v)", timeout)
			case <-ctx.Done():
				runCancelledAwaiter.ItHappened()
			}
		}).Once().Return(false, nil)
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
		cfg := configtest.NewTestGeneralConfig(t)
		cfg.Overrides.GlobalMinIncomingConfirmations = null.IntFrom(1)
		cfg.Overrides.GlobalMinimumContractPayment = assets.NewLinkFromJuels(100)
		uni := NewDirectRequestUniverseWithConfig(t, cfg, nil)
		defer uni.Cleanup()

		log := new(log_mocks.Broadcast)
		defer log.AssertExpectations(t)

		uni.logBroadcaster.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)
		logOracleRequest := operator_wrapper.OperatorOracleRequest{
			CancelExpiration: big.NewInt(0),
			Payment:          big.NewInt(100),
		}
		log.On("RawLog").Return(types.Log{
			Topics: []common.Hash{
				{},
				uni.spec.ExternalIDEncodeStringToTopic(),
			},
		})
		log.On("DecodedLog").Return(&logOracleRequest)
		uni.logBroadcaster.On("MarkConsumed", mock.Anything, mock.Anything).Return(nil)

		runBeganAwaiter := cltest.NewAwaiter()
		uni.runner.On("Run", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			runBeganAwaiter.ItHappened()
			fn := args.Get(4).(func(pg.Queryer) error)
			fn(nil)
		}).Once().Return(false, nil)

		err := uni.service.Start(testutils.Context(t))
		require.NoError(t, err)

		// check if the job exists under the correct ID
		drJob, jErr := uni.jobORM.FindJob(context.Background(), uni.listener.JobID())
		require.NoError(t, jErr)
		require.Equal(t, drJob.ID, uni.listener.JobID())
		require.NotNil(t, drJob.DirectRequestSpec)

		uni.listener.HandleLog(log)

		runBeganAwaiter.AwaitOrFail(t, 5*time.Second)

		uni.service.Close()
		uni.logBroadcaster.AssertExpectations(t)
		uni.runner.AssertExpectations(t)
	})

	t.Run("Log has insufficient funds", func(t *testing.T) {
		cfg := configtest.NewTestGeneralConfig(t)
		cfg.Overrides.GlobalMinIncomingConfirmations = null.IntFrom(1)
		cfg.Overrides.GlobalMinimumContractPayment = assets.NewLinkFromJuels(100)
		uni := NewDirectRequestUniverseWithConfig(t, cfg, nil)
		defer uni.Cleanup()

		log := new(log_mocks.Broadcast)
		defer log.AssertExpectations(t)

		uni.logBroadcaster.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)
		logOracleRequest := operator_wrapper.OperatorOracleRequest{
			CancelExpiration: big.NewInt(0),
			Payment:          big.NewInt(99),
		}
		log.On("RawLog").Return(types.Log{
			Topics: []common.Hash{
				{},
				uni.spec.ExternalIDEncodeStringToTopic(),
			},
		})
		log.On("DecodedLog").Return(&logOracleRequest)
		markConsumedLogAwaiter := cltest.NewAwaiter()
		uni.logBroadcaster.On("MarkConsumed", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			markConsumedLogAwaiter.ItHappened()
		}).Return(nil)

		err := uni.service.Start(testutils.Context(t))
		require.NoError(t, err)

		uni.listener.HandleLog(log)

		markConsumedLogAwaiter.AwaitOrFail(t, 5*time.Second)

		uni.service.Close()
		uni.logBroadcaster.AssertExpectations(t)
		uni.runner.AssertExpectations(t)
	})

	t.Run("requesters is specified and log is requested by a whitelisted address", func(t *testing.T) {
		requester := testutils.NewAddress()
		cfg := configtest.NewTestGeneralConfig(t)
		cfg.Overrides.GlobalMinIncomingConfirmations = null.IntFrom(1)
		cfg.Overrides.GlobalMinimumContractPayment = assets.NewLinkFromJuels(100)
		uni := NewDirectRequestUniverseWithConfig(t, cfg, func(jb *job.Job) {
			jb.DirectRequestSpec.Requesters = []common.Address{testutils.NewAddress(), requester}
		})
		defer uni.Cleanup()

		log := new(log_mocks.Broadcast)
		defer log.AssertExpectations(t)

		uni.logBroadcaster.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)
		logOracleRequest := operator_wrapper.OperatorOracleRequest{
			CancelExpiration: big.NewInt(0),
			Payment:          big.NewInt(100),
			Requester:        requester,
		}
		log.On("RawLog").Return(types.Log{
			Topics: []common.Hash{
				{},
				uni.spec.ExternalIDEncodeStringToTopic(),
			},
		})
		log.On("DecodedLog").Return(&logOracleRequest)
		markConsumedLogAwaiter := cltest.NewAwaiter()
		uni.logBroadcaster.On("MarkConsumed", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			markConsumedLogAwaiter.ItHappened()
		}).Return(nil)

		runBeganAwaiter := cltest.NewAwaiter()
		uni.runner.On("Run", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			runBeganAwaiter.ItHappened()
			fn := args.Get(4).(func(pg.Queryer) error)
			fn(nil)
		}).Once().Return(false, nil)

		err := uni.service.Start(testutils.Context(t))
		require.NoError(t, err)

		// check if the job exists under the correct ID
		drJob, jErr := uni.jobORM.FindJob(context.Background(), uni.listener.JobID())
		require.NoError(t, jErr)
		require.Equal(t, drJob.ID, uni.listener.JobID())
		require.NotNil(t, drJob.DirectRequestSpec)

		uni.listener.HandleLog(log)

		runBeganAwaiter.AwaitOrFail(t, 5*time.Second)

		uni.service.Close()
		uni.logBroadcaster.AssertExpectations(t)
		uni.runner.AssertExpectations(t)
	})

	t.Run("requesters is specified and log is requested by a non-whitelisted address", func(t *testing.T) {
		requester := testutils.NewAddress()
		cfg := configtest.NewTestGeneralConfig(t)
		cfg.Overrides.GlobalMinIncomingConfirmations = null.IntFrom(1)
		cfg.Overrides.GlobalMinimumContractPayment = assets.NewLinkFromJuels(100)
		uni := NewDirectRequestUniverseWithConfig(t, cfg, func(jb *job.Job) {
			jb.DirectRequestSpec.Requesters = []common.Address{testutils.NewAddress(), testutils.NewAddress()}
		})
		defer uni.Cleanup()

		log := new(log_mocks.Broadcast)
		defer log.AssertExpectations(t)

		uni.logBroadcaster.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)
		logOracleRequest := operator_wrapper.OperatorOracleRequest{
			CancelExpiration: big.NewInt(0),
			Payment:          big.NewInt(100),
			Requester:        requester,
		}
		log.On("RawLog").Return(types.Log{
			Topics: []common.Hash{
				{},
				uni.spec.ExternalIDEncodeStringToTopic(),
			},
		})
		log.On("DecodedLog").Return(&logOracleRequest)
		markConsumedLogAwaiter := cltest.NewAwaiter()
		uni.logBroadcaster.On("MarkConsumed", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			markConsumedLogAwaiter.ItHappened()
		}).Return(nil)

		err := uni.service.Start(testutils.Context(t))
		require.NoError(t, err)

		uni.listener.HandleLog(log)

		markConsumedLogAwaiter.AwaitOrFail(t, 5*time.Second)

		uni.service.Close()
		uni.logBroadcaster.AssertExpectations(t)
		uni.runner.AssertExpectations(t)
	})
}

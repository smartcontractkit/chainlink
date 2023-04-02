package directrequest_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/assets"
	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/log"
	log_mocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/log/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/operator_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	configtest "github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/directrequest"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	pipeline_mocks "github.com/smartcontractkit/chainlink/v2/core/services/pipeline/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/srvctest"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func TestDelegate_ServicesForSpec(t *testing.T) {
	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	runner := pipeline_mocks.NewRunner(t)
	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].MinIncomingConfirmations = ptr[uint32](1)
	})
	keyStore := cltest.NewKeyStore(t, db, cfg)
	mailMon := srvctest.Start(t, utils.NewMailboxMonitor(t.Name()))
	cc := evmtest.NewChainSet(t, evmtest.TestChainOpts{DB: db, GeneralConfig: cfg, Client: ethClient, MailMon: mailMon, KeyStore: keyStore.Eth()})

	lggr := logger.TestLogger(t)
	delegate := directrequest.NewDelegate(lggr, runner, nil, cc, mailMon)

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

func NewDirectRequestUniverseWithConfig(t *testing.T, cfg chainlink.GeneralConfig, specF func(spec *job.Job)) *DirectRequestUniverse {
	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	broadcaster := log_mocks.NewBroadcaster(t)
	runner := pipeline_mocks.NewRunner(t)
	broadcaster.On("AddDependents", 1)

	mailMon := srvctest.Start(t, utils.NewMailboxMonitor(t.Name()))

	db := pgtest.NewSqlxDB(t)
	keyStore := cltest.NewKeyStore(t, db, cfg)
	cc := evmtest.NewChainSet(t, evmtest.TestChainOpts{DB: db, GeneralConfig: cfg, Client: ethClient, LogBroadcaster: broadcaster, MailMon: mailMon, KeyStore: keyStore.Eth()})
	lggr := logger.TestLogger(t)
	orm := pipeline.NewORM(db, lggr, cfg)
	btORM := bridges.NewORM(db, lggr, cfg)

	jobORM := job.NewORM(db, cc, orm, btORM, keyStore, lggr, cfg)
	delegate := directrequest.NewDelegate(lggr, runner, orm, cc, mailMon)

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
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].MinIncomingConfirmations = ptr[uint32](1)
	})
	return NewDirectRequestUniverseWithConfig(t, cfg, nil)
}

func (uni *DirectRequestUniverse) Cleanup() {
	uni.cleanup()
}

func TestDelegate_ServicesListenerHandleLog(t *testing.T) {
	testutils.SkipShortDB(t)
	t.Parallel()

	t.Run("Log is an OracleRequest", func(t *testing.T) {
		uni := NewDirectRequestUniverse(t)
		defer uni.Cleanup()

		log := log_mocks.NewBroadcast(t)
		log.On("ReceiptsRoot").Return(common.Hash{})
		log.On("TransactionsRoot").Return(common.Hash{})
		log.On("StateRoot").Return(common.Hash{})

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
		log.On("String").Return("")
		uni.logBroadcaster.On("MarkConsumed", mock.Anything, mock.Anything).Return(nil)

		runBeganAwaiter := cltest.NewAwaiter()
		uni.runner.On("Run", mock.Anything, mock.AnythingOfType("*pipeline.Run"), mock.Anything, mock.Anything, mock.Anything).
			Return(false, nil).
			Run(func(args mock.Arguments) {
				runBeganAwaiter.ItHappened()
				fn := args.Get(4).(func(pg.Queryer) error)
				require.NoError(t, fn(nil))
			}).Once()

		err := uni.service.Start(testutils.Context(t))
		require.NoError(t, err)

		require.NotNil(t, uni.listener, "listener was nil; expected broadcaster.Register to have been called")
		// check if the job exists under the correct ID
		drJob, jErr := uni.jobORM.FindJob(testutils.Context(t), uni.listener.JobID())
		require.NoError(t, jErr)
		require.Equal(t, drJob.ID, uni.listener.JobID())
		require.NotNil(t, drJob.DirectRequestSpec)

		uni.listener.HandleLog(log)

		runBeganAwaiter.AwaitOrFail(t, 5*time.Second)

		uni.service.Close()
	})

	t.Run("Log is not consumed, as it's too young", func(t *testing.T) {
		uni := NewDirectRequestUniverse(t)
		defer uni.Cleanup()

		log := log_mocks.NewBroadcast(t)

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
		log.On("String").Return("")
		uni.logBroadcaster.On("MarkConsumed", mock.Anything, mock.Anything).Return(nil).Maybe()

		err := uni.service.Start(testutils.Context(t))
		require.NoError(t, err)

		uni.listener.HandleLog(log)

		uni.logBroadcaster.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)

		log.On("ReceiptsRoot").Return(common.Hash{})
		log.On("TransactionsRoot").Return(common.Hash{})
		log.On("StateRoot").Return(common.Hash{})

		runBeganAwaiter := cltest.NewAwaiter()
		uni.runner.On("Run", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Run(func(args mock.Arguments) {
				runBeganAwaiter.ItHappened()
				fn := args.Get(4).(func(pg.Queryer) error)
				require.NoError(t, fn(nil))
			}).Once().Return(false, nil)

		// but should after this one, as the head Number is larger
		runBeganAwaiter.AwaitOrFail(t, 5*time.Second)

		uni.service.Close()
	})

	t.Run("Log has wrong jobID", func(t *testing.T) {
		uni := NewDirectRequestUniverse(t)
		defer uni.Cleanup()

		log := log_mocks.NewBroadcast(t)
		lbAwaiter := cltest.NewAwaiter()
		uni.logBroadcaster.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)
		uni.logBroadcaster.On("MarkConsumed", mock.Anything, mock.Anything).Run(func(args mock.Arguments) { lbAwaiter.ItHappened() }).Return(nil)

		logCancelOracleRequest := operator_wrapper.OperatorCancelOracleRequest{RequestId: uni.spec.ExternalIDEncodeStringToTopic()}
		logAwaiter := cltest.NewAwaiter()
		log.On("DecodedLog").Run(func(args mock.Arguments) { logAwaiter.ItHappened() }).Return(&logCancelOracleRequest)
		log.On("RawLog").Return(types.Log{
			Topics: []common.Hash{{}, {}},
		})
		log.On("String").Return("")

		err := uni.service.Start(testutils.Context(t))
		require.NoError(t, err)

		uni.listener.HandleLog(log)

		logAwaiter.AwaitOrFail(t)
		lbAwaiter.AwaitOrFail(t)

		uni.service.Close()
	})

	t.Run("Log is a CancelOracleRequest with no matching run", func(t *testing.T) {
		uni := NewDirectRequestUniverse(t)
		defer uni.Cleanup()

		log := log_mocks.NewBroadcast(t)

		uni.logBroadcaster.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)
		logCancelOracleRequest := operator_wrapper.OperatorCancelOracleRequest{RequestId: uni.spec.ExternalIDEncodeStringToTopic()}
		log.On("RawLog").Return(types.Log{
			Topics: []common.Hash{
				{},
				uni.spec.ExternalIDEncodeStringToTopic(),
			},
		})
		log.On("String").Return("")
		log.On("DecodedLog").Return(&logCancelOracleRequest)
		lbAwaiter := cltest.NewAwaiter()
		uni.logBroadcaster.On("MarkConsumed", mock.Anything, mock.Anything).Run(func(args mock.Arguments) { lbAwaiter.ItHappened() }).Return(nil)

		err := uni.service.Start(testutils.Context(t))
		require.NoError(t, err)

		uni.listener.HandleLog(log)

		lbAwaiter.AwaitOrFail(t)

		uni.service.Close()
	})

	t.Run("Log is a CancelOracleRequest with a matching run", func(t *testing.T) {
		uni := NewDirectRequestUniverse(t)
		defer uni.Cleanup()

		runLog := log_mocks.NewBroadcast(t)
		runLog.On("ReceiptsRoot").Return(common.Hash{})
		runLog.On("TransactionsRoot").Return(common.Hash{})
		runLog.On("StateRoot").Return(common.Hash{})

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
		runLog.On("String").Return("")
		uni.logBroadcaster.On("MarkConsumed", mock.Anything, mock.Anything).Return(nil)

		cancelLog := log_mocks.NewBroadcast(t)

		uni.logBroadcaster.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)
		logCancelOracleRequest := operator_wrapper.OperatorCancelOracleRequest{RequestId: uni.spec.ExternalIDEncodeStringToTopic()}
		cancelLog.On("RawLog").Return(types.Log{
			Topics: []common.Hash{
				{},
				uni.spec.ExternalIDEncodeStringToTopic(),
			},
		})
		cancelLog.On("DecodedLog").Return(&logCancelOracleRequest)
		cancelLog.On("String").Return("")
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

		uni.listener.HandleLog(cancelLog)

		runCancelledAwaiter.AwaitOrFail(t, timeout)

		uni.service.Close()
	})

	t.Run("Log has sufficient funds", func(t *testing.T) {
		cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
			c.EVM[0].MinIncomingConfirmations = ptr[uint32](1)
			c.EVM[0].MinContractPayment = assets.NewLinkFromJuels(100)
		})
		uni := NewDirectRequestUniverseWithConfig(t, cfg, nil)
		defer uni.Cleanup()

		log := log_mocks.NewBroadcast(t)
		log.On("ReceiptsRoot").Return(common.Hash{})
		log.On("TransactionsRoot").Return(common.Hash{})
		log.On("StateRoot").Return(common.Hash{})

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
		log.On("String").Return("")
		uni.logBroadcaster.On("MarkConsumed", mock.Anything, mock.Anything).Return(nil)

		runBeganAwaiter := cltest.NewAwaiter()
		uni.runner.On("Run", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			runBeganAwaiter.ItHappened()
			fn := args.Get(4).(func(pg.Queryer) error)
			require.NoError(t, fn(nil))
		}).Once().Return(false, nil)

		err := uni.service.Start(testutils.Context(t))
		require.NoError(t, err)

		// check if the job exists under the correct ID
		drJob, jErr := uni.jobORM.FindJob(testutils.Context(t), uni.listener.JobID())
		require.NoError(t, jErr)
		require.Equal(t, drJob.ID, uni.listener.JobID())
		require.NotNil(t, drJob.DirectRequestSpec)

		uni.listener.HandleLog(log)

		runBeganAwaiter.AwaitOrFail(t, 5*time.Second)

		uni.service.Close()
	})

	t.Run("Log has insufficient funds", func(t *testing.T) {
		cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
			c.EVM[0].MinIncomingConfirmations = ptr[uint32](1)
			c.EVM[0].MinContractPayment = assets.NewLinkFromJuels(100)
		})
		uni := NewDirectRequestUniverseWithConfig(t, cfg, nil)
		defer uni.Cleanup()

		log := log_mocks.NewBroadcast(t)

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
		log.On("String").Return("")
		markConsumedLogAwaiter := cltest.NewAwaiter()
		uni.logBroadcaster.On("MarkConsumed", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			markConsumedLogAwaiter.ItHappened()
		}).Return(nil)

		err := uni.service.Start(testutils.Context(t))
		require.NoError(t, err)

		uni.listener.HandleLog(log)

		markConsumedLogAwaiter.AwaitOrFail(t, 5*time.Second)

		uni.service.Close()
	})

	t.Run("requesters is specified and log is requested by a whitelisted address", func(t *testing.T) {
		requester := testutils.NewAddress()
		cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
			c.EVM[0].MinIncomingConfirmations = ptr[uint32](1)
			c.EVM[0].MinContractPayment = assets.NewLinkFromJuels(100)
		})
		uni := NewDirectRequestUniverseWithConfig(t, cfg, func(jb *job.Job) {
			jb.DirectRequestSpec.Requesters = []common.Address{testutils.NewAddress(), requester}
		})
		defer uni.Cleanup()

		log := log_mocks.NewBroadcast(t)
		log.On("ReceiptsRoot").Return(common.Hash{})
		log.On("TransactionsRoot").Return(common.Hash{})
		log.On("StateRoot").Return(common.Hash{})

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
		log.On("String").Return("")
		markConsumedLogAwaiter := cltest.NewAwaiter()
		uni.logBroadcaster.On("MarkConsumed", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			markConsumedLogAwaiter.ItHappened()
		}).Return(nil)

		runBeganAwaiter := cltest.NewAwaiter()
		uni.runner.On("Run", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			runBeganAwaiter.ItHappened()
			fn := args.Get(4).(func(pg.Queryer) error)
			require.NoError(t, fn(nil))
		}).Once().Return(false, nil)

		err := uni.service.Start(testutils.Context(t))
		require.NoError(t, err)

		// check if the job exists under the correct ID
		drJob, jErr := uni.jobORM.FindJob(testutils.Context(t), uni.listener.JobID())
		require.NoError(t, jErr)
		require.Equal(t, drJob.ID, uni.listener.JobID())
		require.NotNil(t, drJob.DirectRequestSpec)

		uni.listener.HandleLog(log)

		runBeganAwaiter.AwaitOrFail(t, 5*time.Second)

		uni.service.Close()
	})

	t.Run("requesters is specified and log is requested by a non-whitelisted address", func(t *testing.T) {
		requester := testutils.NewAddress()
		cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
			c.EVM[0].MinIncomingConfirmations = ptr[uint32](1)
			c.EVM[0].MinContractPayment = assets.NewLinkFromJuels(100)
		})
		uni := NewDirectRequestUniverseWithConfig(t, cfg, func(jb *job.Job) {
			jb.DirectRequestSpec.Requesters = []common.Address{testutils.NewAddress(), testutils.NewAddress()}
		})
		defer uni.Cleanup()

		log := log_mocks.NewBroadcast(t)

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
		log.On("String").Return("")
		markConsumedLogAwaiter := cltest.NewAwaiter()
		uni.logBroadcaster.On("MarkConsumed", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			markConsumedLogAwaiter.ItHappened()
		}).Return(nil)

		err := uni.service.Start(testutils.Context(t))
		require.NoError(t, err)

		uni.listener.HandleLog(log)

		markConsumedLogAwaiter.AwaitOrFail(t, 5*time.Second)

		uni.service.Close()
	})
}

func ptr[T any](t T) *T { return &t }

package directrequest_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/assets"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/mailbox/mailboxtest"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/operatorforwarder/generated/operator"
	"github.com/smartcontractkit/chainlink-evm/pkg/client/clienttest"
	"github.com/smartcontractkit/chainlink-evm/pkg/log"
	ubig "github.com/smartcontractkit/chainlink-evm/pkg/utils/big"
	log_mocks "github.com/smartcontractkit/chainlink/v2/common/log/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/directrequest"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	pipeline_mocks "github.com/smartcontractkit/chainlink/v2/core/services/pipeline/mocks"
)

func TestDelegate_ServicesForSpec(t *testing.T) {
	ethClient := clienttest.NewClientWithDefaultChainID(t)
	runner := pipeline_mocks.NewRunner(t)
	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].MinIncomingConfirmations = ptr[uint32](1)
	})
	keyStore := cltest.NewKeyStore(t, db)
	mailMon := servicetest.Run(t, mailboxtest.NewMonitor(t))
	legacyChains := evmtest.NewLegacyChains(t, evmtest.TestChainOpts{
		ChainConfigs:   cfg.EVMConfigs(),
		DatabaseConfig: cfg.Database(),
		FeatureConfig:  cfg.Feature(),
		ListenerConfig: cfg.Database().Listener(),
		DB:             db,
		Client:         ethClient,
		MailMon:        mailMon,
		KeyStore:       keyStore.Eth(),
	})

	lggr := logger.TestLogger(t)
	delegate := directrequest.NewDelegate(lggr, runner, nil, legacyChains, mailMon)

	t.Run("Spec without DirectRequestSpec", func(t *testing.T) {
		spec := job.Job{}
		_, err := delegate.ServicesForSpec(testutils.Context(t), spec)
		assert.Error(t, err, "expects a *job.DirectRequestSpec to be present")
	})

	t.Run("Spec with DirectRequestSpec", func(t *testing.T) {
		spec := job.Job{DirectRequestSpec: &job.DirectRequestSpec{EVMChainID: (*ubig.Big)(testutils.FixtureChainID)}, PipelineSpec: &pipeline.Spec{}}
		services, err := delegate.ServicesForSpec(testutils.Context(t), spec)
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
	ethClient := clienttest.NewClientWithDefaultChainID(t)
	broadcaster := log_mocks.NewBroadcaster(t)
	runner := pipeline_mocks.NewRunner(t)
	broadcaster.On("AddDependents", 1)

	mailMon := servicetest.Run(t, mailboxtest.NewMonitor(t))

	db := pgtest.NewSqlxDB(t)
	keyStore := cltest.NewKeyStore(t, db)
	legacyChains := evmtest.NewLegacyChains(t, evmtest.TestChainOpts{
		ChainConfigs:   cfg.EVMConfigs(),
		DatabaseConfig: cfg.Database(),
		FeatureConfig:  cfg.Feature(),
		ListenerConfig: cfg.Database().Listener(),
		DB:             db,
		Client:         ethClient,
		LogBroadcaster: broadcaster,
		MailMon:        mailMon,
		KeyStore:       keyStore.Eth(),
	})
	lggr := logger.TestLogger(t)
	orm := pipeline.NewORM(db, lggr, cfg.JobPipeline().MaxSuccessfulRuns())
	btORM := bridges.NewORM(db)
	jobORM := job.NewORM(db, orm, btORM, keyStore, lggr)
	delegate := directrequest.NewDelegate(lggr, runner, orm, legacyChains, mailMon)

	jb := cltest.MakeDirectRequestJobSpec(t)
	jb.ExternalJobID = uuid.New()
	if specF != nil {
		specF(jb)
	}
	ctx := testutils.Context(t)
	require.NoError(t, jobORM.CreateJob(ctx, jb))
	serviceArray, err := delegate.ServicesForSpec(ctx, *jb)
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
		log.On("EVMChainID").Return(*big.NewInt(0))

		uni.logBroadcaster.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)
		logOracleRequest := operator.OperatorOracleRequest{
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
		uni.logBroadcaster.On("MarkConsumed", mock.Anything, mock.Anything, mock.Anything).Return(nil)

		runBeganAwaiter := cltest.NewAwaiter()
		uni.runner.On("Run", mock.Anything, mock.AnythingOfType("*pipeline.Run"), mock.Anything, mock.Anything).
			Return(false, nil).
			Run(func(args mock.Arguments) {
				runBeganAwaiter.ItHappened()
				fn := args.Get(3).(func(source sqlutil.DataSource) error)
				require.NoError(t, fn(nil))
			}).Once()

		ctx := testutils.Context(t)
		err := uni.service.Start(ctx)
		require.NoError(t, err)

		require.NotNil(t, uni.listener, "listener was nil; expected broadcaster.Register to have been called")
		// check if the job exists under the correct ID
		drJob, jErr := uni.jobORM.FindJob(ctx, uni.listener.JobID())
		require.NoError(t, jErr)
		require.Equal(t, drJob.ID, uni.listener.JobID())
		require.NotNil(t, drJob.DirectRequestSpec)

		uni.listener.HandleLog(ctx, log)

		runBeganAwaiter.AwaitOrFail(t, 5*time.Second)

		uni.service.Close()
	})

	t.Run("Log is not consumed, as it's too young", func(t *testing.T) {
		uni := NewDirectRequestUniverse(t)
		defer uni.Cleanup()

		log := log_mocks.NewBroadcast(t)

		uni.logBroadcaster.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil).Maybe()
		logOracleRequest := operator.OperatorOracleRequest{
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
		log.On("EVMChainID").Return(*big.NewInt(0))
		uni.logBroadcaster.On("MarkConsumed", mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

		ctx := testutils.Context(t)
		err := uni.service.Start(ctx)
		require.NoError(t, err)

		uni.listener.HandleLog(ctx, log)

		uni.logBroadcaster.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)

		log.On("ReceiptsRoot").Return(common.Hash{})
		log.On("TransactionsRoot").Return(common.Hash{})
		log.On("StateRoot").Return(common.Hash{})

		runBeganAwaiter := cltest.NewAwaiter()
		uni.runner.On("Run", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Run(func(args mock.Arguments) {
				runBeganAwaiter.ItHappened()
				fn := args.Get(3).(func(sqlutil.DataSource) error)
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
		uni.logBroadcaster.On("MarkConsumed", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) { lbAwaiter.ItHappened() }).Return(nil)

		logCancelOracleRequest := operator.OperatorCancelOracleRequest{RequestId: uni.spec.ExternalIDEncodeStringToTopic()}
		logAwaiter := cltest.NewAwaiter()
		log.On("DecodedLog").Run(func(args mock.Arguments) { logAwaiter.ItHappened() }).Return(&logCancelOracleRequest)
		log.On("RawLog").Return(types.Log{
			Topics: []common.Hash{{}, {}},
		})
		log.On("String").Return("")

		ctx := testutils.Context(t)
		err := uni.service.Start(ctx)
		require.NoError(t, err)

		uni.listener.HandleLog(ctx, log)

		logAwaiter.AwaitOrFail(t)
		lbAwaiter.AwaitOrFail(t)

		uni.service.Close()
	})

	t.Run("Log is a CancelOracleRequest with no matching run", func(t *testing.T) {
		uni := NewDirectRequestUniverse(t)
		defer uni.Cleanup()

		log := log_mocks.NewBroadcast(t)

		uni.logBroadcaster.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)
		logCancelOracleRequest := operator.OperatorCancelOracleRequest{RequestId: uni.spec.ExternalIDEncodeStringToTopic()}
		log.On("RawLog").Return(types.Log{
			Topics: []common.Hash{
				{},
				uni.spec.ExternalIDEncodeStringToTopic(),
			},
		})
		log.On("String").Return("")
		log.On("DecodedLog").Return(&logCancelOracleRequest)
		lbAwaiter := cltest.NewAwaiter()
		uni.logBroadcaster.On("MarkConsumed", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) { lbAwaiter.ItHappened() }).Return(nil)

		ctx := testutils.Context(t)
		err := uni.service.Start(ctx)
		require.NoError(t, err)

		uni.listener.HandleLog(ctx, log)

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
		runLog.On("EVMChainID").Return(*big.NewInt(0))

		uni.logBroadcaster.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)
		logOracleRequest := operator.OperatorOracleRequest{
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
		uni.logBroadcaster.On("MarkConsumed", mock.Anything, mock.Anything, mock.Anything).Return(nil)

		cancelLog := log_mocks.NewBroadcast(t)

		uni.logBroadcaster.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)
		logCancelOracleRequest := operator.OperatorCancelOracleRequest{RequestId: uni.spec.ExternalIDEncodeStringToTopic()}
		cancelLog.On("RawLog").Return(types.Log{
			Topics: []common.Hash{
				{},
				uni.spec.ExternalIDEncodeStringToTopic(),
			},
		})
		cancelLog.On("DecodedLog").Return(&logCancelOracleRequest)
		cancelLog.On("String").Return("")
		uni.logBroadcaster.On("MarkConsumed", mock.Anything, mock.Anything, mock.Anything).Return(nil)

		ctx := testutils.Context(t)
		err := uni.service.Start(ctx)
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
		uni.listener.HandleLog(ctx, runLog)

		runBeganAwaiter.AwaitOrFail(t, timeout)

		uni.listener.HandleLog(ctx, cancelLog)

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
		log.On("EVMChainID").Return(*big.NewInt(0))

		uni.logBroadcaster.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)
		logOracleRequest := operator.OperatorOracleRequest{
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
		uni.logBroadcaster.On("MarkConsumed", mock.Anything, mock.Anything, mock.Anything).Return(nil)

		runBeganAwaiter := cltest.NewAwaiter()
		uni.runner.On("Run", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			runBeganAwaiter.ItHappened()
			fn := args.Get(3).(func(sqlutil.DataSource) error)
			require.NoError(t, fn(nil))
		}).Once().Return(false, nil)

		ctx := testutils.Context(t)
		err := uni.service.Start(ctx)
		require.NoError(t, err)

		// check if the job exists under the correct ID
		drJob, jErr := uni.jobORM.FindJob(ctx, uni.listener.JobID())
		require.NoError(t, jErr)
		require.Equal(t, drJob.ID, uni.listener.JobID())
		require.NotNil(t, drJob.DirectRequestSpec)

		uni.listener.HandleLog(ctx, log)

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
		logOracleRequest := operator.OperatorOracleRequest{
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
		uni.logBroadcaster.On("MarkConsumed", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			markConsumedLogAwaiter.ItHappened()
		}).Return(nil)

		ctx := testutils.Context(t)
		err := uni.service.Start(ctx)
		require.NoError(t, err)

		uni.listener.HandleLog(ctx, log)

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
		log.On("EVMChainID").Return(*big.NewInt(0))

		uni.logBroadcaster.On("WasAlreadyConsumed", mock.Anything, mock.Anything).Return(false, nil)
		logOracleRequest := operator.OperatorOracleRequest{
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
		uni.logBroadcaster.On("MarkConsumed", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			markConsumedLogAwaiter.ItHappened()
		}).Return(nil)

		runBeganAwaiter := cltest.NewAwaiter()
		uni.runner.On("Run", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			runBeganAwaiter.ItHappened()
			fn := args.Get(3).(func(sqlutil.DataSource) error)
			require.NoError(t, fn(nil))
		}).Once().Return(false, nil)

		ctx := testutils.Context(t)
		err := uni.service.Start(ctx)
		require.NoError(t, err)

		// check if the job exists under the correct ID
		drJob, jErr := uni.jobORM.FindJob(ctx, uni.listener.JobID())
		require.NoError(t, jErr)
		require.Equal(t, drJob.ID, uni.listener.JobID())
		require.NotNil(t, drJob.DirectRequestSpec)

		uni.listener.HandleLog(ctx, log)

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
		logOracleRequest := operator.OperatorOracleRequest{
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
		uni.logBroadcaster.On("MarkConsumed", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			markConsumedLogAwaiter.ItHappened()
		}).Return(nil)

		ctx := testutils.Context(t)
		err := uni.service.Start(ctx)
		require.NoError(t, err)

		uni.listener.HandleLog(ctx, log)

		markConsumedLogAwaiter.AwaitOrFail(t, 5*time.Second)

		uni.service.Close()
	})
}

func ptr[T any](t T) *T { return &t }

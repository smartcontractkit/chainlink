package job_test

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/utils"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/mailbox/mailboxtest"
	evmtypes "github.com/smartcontractkit/chainlink-evm/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/cron"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/job/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/testdata/testspecs"
)

type delegate struct {
	jobType                    job.Type
	services                   []job.ServiceCtx
	jobID                      int32
	chContinueCreatingServices chan struct{}
	job.Delegate
}

func (d delegate) JobType() job.Type {
	return d.jobType
}

// ServicesForSpec satisfies the job.Delegate interface.
func (d delegate) ServicesForSpec(ctx context.Context, js job.Job) ([]job.ServiceCtx, error) {
	if js.Type != d.jobType {
		return nil, nil
	}
	return d.services, nil
}

func clearDB(t *testing.T, db *sqlx.DB) {
	cltest.ClearDBTables(t, db, "jobs", "pipeline_runs", "pipeline_specs", "pipeline_task_runs")
}

func TestSpawner_CreateJobDeleteJob(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	config := configtest.NewTestGeneralConfig(t)
	db := pgtest.NewSqlxDB(t)
	keyStore := cltest.NewKeyStore(t, db)
	ethKeyStore := keyStore.Eth()
	require.NoError(t, keyStore.OCR().Add(ctx, cltest.DefaultOCRKey))
	require.NoError(t, keyStore.P2P().Add(ctx, cltest.DefaultP2PKey))
	require.NoError(t, keyStore.OCR2().Add(ctx, cltest.DefaultOCR2Key))

	_, address := cltest.MustInsertRandomKey(t, ethKeyStore)
	_, bridge := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{})
	_, bridge2 := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{})

	ethClient := cltest.NewEthMocksWithStartupAssertions(t)
	ethClient.On("ConfiguredChainID").Return(testutils.FixtureChainID).Maybe()
	ethClient.On("CallContext", mock.Anything, mock.Anything, "eth_getBlockByNumber", mock.Anything, false).
		Run(func(args mock.Arguments) {
			head := args.Get(1).(**evmtypes.Head)
			*head = cltest.Head(10)
		}).
		Return(nil).Maybe()
	ethClient.On("BalanceAt", mock.Anything, mock.Anything, mock.Anything).Return(big.NewInt(0), nil).Maybe()

	legacyChains := evmtest.NewLegacyChains(t, evmtest.TestChainOpts{
		DB:             db,
		Client:         ethClient,
		ChainConfigs:   config.EVMConfigs(),
		DatabaseConfig: config.Database(),
		FeatureConfig:  config.Feature(),
		ListenerConfig: config.Database().Listener(),
		KeyStore:       ethKeyStore,
	})
	t.Run("should respect its dependents", func(t *testing.T) {
		lggr := logger.TestLogger(t)
		orm := NewTestORM(t, db, pipeline.NewORM(db, lggr, config.JobPipeline().MaxSuccessfulRuns()), bridges.NewORM(db), keyStore)
		a := utils.NewDependentAwaiter()
		a.AddDependents(1)
		spawner := job.NewSpawner(orm, config.Database(), noopChecker{}, map[job.Type]job.Delegate{}, lggr, []utils.DependentAwaiter{a})
		// Starting the spawner should signal to the dependents
		result := make(chan bool)
		go func() {
			select {
			case <-a.AwaitDependents():
				result <- true
			case <-time.After(2 * time.Second):
				result <- false
			}
		}()
		require.NoError(t, spawner.Start(t.Context()))
		assert.True(t, <-result, "failed to signal to dependents")
	})

	t.Run("starts and stops job services when jobs are added and removed", func(t *testing.T) {
		cronJobTmp, cronErr := cron.ValidatedCronSpec(fmt.Sprintf(testspecs.CronSpecTemplate, uuid.New()))
		require.NoError(t, cronErr)
		jobA := &cronJobTmp
		jobB := makeOCRJobSpec(t, address, bridge.Name.String(), bridge2.Name.String())

		lggr := logger.TestLogger(t)
		orm := NewTestORM(t, db, pipeline.NewORM(db, lggr, config.JobPipeline().MaxSuccessfulRuns()), bridges.NewORM(db), keyStore)

		eventuallyA := cltest.NewAwaiter()
		serviceA1 := mocks.NewServiceCtx(t)
		serviceA2 := mocks.NewServiceCtx(t)
		serviceA1.On("Start", mock.Anything).Return(nil).Once()
		serviceA2.On("Start", mock.Anything).Return(nil).Once().Run(func(mock.Arguments) { eventuallyA.ItHappened() })
		mailMon := servicetest.Run(t, mailboxtest.NewMonitor(t))
		dA := ocr.NewDelegate(nil, orm, nil, nil, nil, nil, monitoringEndpoint, legacyChains, logger.TestLogger(t), config, mailMon)
		delegateA := &delegate{jobA.Type, []job.ServiceCtx{serviceA1, serviceA2}, 0, make(chan struct{}), dA}

		eventuallyB := cltest.NewAwaiter()
		serviceB1 := mocks.NewServiceCtx(t)
		serviceB2 := mocks.NewServiceCtx(t)
		serviceB1.On("Start", mock.Anything).Return(nil).Once()
		serviceB2.On("Start", mock.Anything).Return(nil).Once().Run(func(mock.Arguments) { eventuallyB.ItHappened() })
		dB := ocr.NewDelegate(nil, orm, nil, nil, nil, nil, monitoringEndpoint, legacyChains, logger.TestLogger(t), config, mailMon)
		delegateB := &delegate{jobB.Type, []job.ServiceCtx{serviceB1, serviceB2}, 0, make(chan struct{}), dB}

		spawner := job.NewSpawner(orm, config.Database(), noopChecker{}, map[job.Type]job.Delegate{
			jobA.Type: delegateA,
			jobB.Type: delegateB,
		}, lggr, nil)
		ctx := t.Context()
		require.NoError(t, spawner.Start(ctx))
		err := spawner.CreateJob(ctx, nil, jobA)
		require.NoError(t, err)
		jobSpecIDA := jobA.ID
		delegateA.jobID = jobSpecIDA
		close(delegateA.chContinueCreatingServices)

		eventuallyA.AwaitOrFail(t, 20*time.Second)

		err = spawner.CreateJob(ctx, nil, jobB)
		require.NoError(t, err)
		jobSpecIDB := jobB.ID
		delegateB.jobID = jobSpecIDB
		close(delegateB.chContinueCreatingServices)

		eventuallyB.AwaitOrFail(t, 20*time.Second)

		serviceA1.On("Close").Return(nil).Once()
		serviceA2.On("Close").Return(nil).Once()
		err = spawner.DeleteJob(ctx, nil, jobSpecIDA)
		require.NoError(t, err)

		serviceB1.On("Close").Return(nil).Once()
		serviceB2.On("Close").Return(nil).Once()
		err = spawner.DeleteJob(ctx, nil, jobSpecIDB)
		require.NoError(t, err)

		require.NoError(t, spawner.Close())
	})

	clearDB(t, db)

	t.Run("starts and stops job services from the DB when .Start()/.Stop() is called", func(t *testing.T) {
		jobA := makeOCRJobSpec(t, address, bridge.Name.String(), bridge2.Name.String())

		eventually := cltest.NewAwaiter()
		serviceA1 := mocks.NewServiceCtx(t)
		serviceA2 := mocks.NewServiceCtx(t)
		serviceA1.On("Start", mock.Anything).Return(nil).Once()
		serviceA2.On("Start", mock.Anything).Return(nil).Once().Run(func(mock.Arguments) { eventually.ItHappened() })

		lggr := logger.TestLogger(t)
		orm := NewTestORM(t, db, pipeline.NewORM(db, lggr, config.JobPipeline().MaxSuccessfulRuns()), bridges.NewORM(db), keyStore)
		mailMon := servicetest.Run(t, mailboxtest.NewMonitor(t))
		d := ocr.NewDelegate(nil, orm, nil, nil, nil, nil, monitoringEndpoint, legacyChains, logger.TestLogger(t), config, mailMon)
		delegateA := &delegate{jobA.Type, []job.ServiceCtx{serviceA1, serviceA2}, 0, nil, d}
		spawner := job.NewSpawner(orm, config.Database(), noopChecker{}, map[job.Type]job.Delegate{
			jobA.Type: delegateA,
		}, lggr, nil)

		ctx := t.Context()
		err := orm.CreateJob(ctx, jobA)
		require.NoError(t, err)
		delegateA.jobID = jobA.ID

		require.NoError(t, spawner.Start(ctx))

		eventually.AwaitOrFail(t)

		serviceA1.On("Close").Return(nil).Once()
		serviceA2.On("Close").Return(nil).Once()

		require.NoError(t, spawner.Close())
	})

	clearDB(t, db)

	t.Run("closes job services on 'DeleteJob()'", func(t *testing.T) {
		jobA := makeOCRJobSpec(t, address, bridge.Name.String(), bridge2.Name.String())

		eventuallyStart := cltest.NewAwaiter()
		serviceA1 := mocks.NewServiceCtx(t)
		serviceA2 := mocks.NewServiceCtx(t)
		serviceA1.On("Start", mock.Anything).Return(nil).Once()
		serviceA2.On("Start", mock.Anything).Return(nil).Once().Run(func(mock.Arguments) { eventuallyStart.ItHappened() })

		lggr := logger.TestLogger(t)
		orm := NewTestORM(t, db, pipeline.NewORM(db, lggr, config.JobPipeline().MaxSuccessfulRuns()), bridges.NewORM(db), keyStore)
		mailMon := servicetest.Run(t, mailboxtest.NewMonitor(t))
		d := ocr.NewDelegate(nil, orm, nil, nil, nil, nil, monitoringEndpoint, legacyChains, logger.TestLogger(t), config, mailMon)
		delegateA := &delegate{jobA.Type, []job.ServiceCtx{serviceA1, serviceA2}, 0, nil, d}
		spawner := job.NewSpawner(orm, config.Database(), noopChecker{}, map[job.Type]job.Delegate{
			jobA.Type: delegateA,
		}, lggr, nil)

		ctx := t.Context()
		err := orm.CreateJob(ctx, jobA)
		require.NoError(t, err)
		jobSpecIDA := jobA.ID
		delegateA.jobID = jobSpecIDA

		require.NoError(t, spawner.Start(ctx))
		defer func() { assert.NoError(t, spawner.Close()) }()

		eventuallyStart.AwaitOrFail(t)

		// Wait for the claim lock to be taken
		require.Eventually(t, func() bool {
			jobs := spawner.ActiveJobs()
			_, exists := jobs[jobSpecIDA]
			return exists
		}, testutils.WaitTimeout(t), cltest.DBPollingInterval)

		eventuallyClose := cltest.NewAwaiter()
		serviceA1.On("Close").Return(nil).Once()
		serviceA2.On("Close").Return(nil).Once().Run(func(mock.Arguments) { eventuallyClose.ItHappened() })

		err = spawner.DeleteJob(ctx, nil, jobSpecIDA)
		require.NoError(t, err)

		eventuallyClose.AwaitOrFail(t)

		// Wait for the claim lock to be released
		require.Eventually(t, func() bool {
			jobs := spawner.ActiveJobs()
			_, exists := jobs[jobSpecIDA]
			return !exists
		}, testutils.WaitTimeout(t), cltest.DBPollingInterval)

		clearDB(t, db)
	})
}

type noopChecker struct{}

func (n noopChecker) Register(service services.HealthReporter) error { return nil }

func (n noopChecker) Unregister(name string) error { return nil }

func (n noopChecker) IsReady() (ready bool, errors map[string]error) { return true, nil }

func (n noopChecker) IsHealthy() (healthy bool, errors map[string]error) { return true, nil }

func (n noopChecker) Start() error { return nil }

func (n noopChecker) Close() error { return nil }

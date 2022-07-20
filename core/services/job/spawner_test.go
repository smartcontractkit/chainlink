package job_test

import (
	"testing"
	"time"

	"github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/sqlx"

	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/job/mocks"
	"github.com/smartcontractkit/chainlink/core/services/ocr"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/utils"
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
func (d delegate) ServicesForSpec(js job.Job) ([]job.ServiceCtx, error) {
	if js.Type != d.jobType {
		return nil, nil
	}
	return d.services, nil
}

func clearDB(t *testing.T, db *sqlx.DB) {
	cltest.ClearDBTables(t, db, "jobs", "pipeline_runs", "pipeline_specs", "pipeline_task_runs")
}

func TestSpawner_CreateJobDeleteJob(t *testing.T) {
	config := cltest.NewTestGeneralConfig(t)
	db := pgtest.NewSqlxDB(t)
	keyStore := cltest.NewKeyStore(t, db, config)
	ethKeyStore := keyStore.Eth()
	require.NoError(t, keyStore.OCR().Add(cltest.DefaultOCRKey))
	require.NoError(t, keyStore.P2P().Add(cltest.DefaultP2PKey))

	_, address := cltest.MustInsertRandomKey(t, ethKeyStore)
	_, bridge := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{}, config)
	_, bridge2 := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{}, config)

	ethClient := cltest.NewEthMocksWithDefaultChain(t)
	ethClient.On("CallContext", mock.Anything, mock.Anything, "eth_getBlockByNumber", mock.Anything, false).
		Run(func(args mock.Arguments) {
			head := args.Get(1).(**evmtypes.Head)
			*head = cltest.Head(10)
		}).
		Return(nil).Maybe()
	cc := evmtest.NewChainSet(t, evmtest.TestChainOpts{DB: db, Client: ethClient, GeneralConfig: config})

	t.Run("should respect its dependents", func(t *testing.T) {
		lggr := logger.TestLogger(t)
		orm := job.NewTestORM(t, db, cc, pipeline.NewORM(db, lggr, config), keyStore, config)
		a := utils.NewDependentAwaiter()
		a.AddDependents(1)
		spawner := job.NewSpawner(orm, config, map[job.Type]job.Delegate{}, db, lggr, []utils.DependentAwaiter{a})
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
		spawner.Start(testutils.Context(t))
		assert.True(t, <-result, "failed to signal to dependents")
	})

	t.Run("starts and stops job services when jobs are added and removed", func(t *testing.T) {
		jobA := cltest.MakeDirectRequestJobSpec(t)
		jobB := makeOCRJobSpec(t, address, bridge.Name.String(), bridge2.Name.String())

		lggr := logger.TestLogger(t)
		orm := job.NewTestORM(t, db, cc, pipeline.NewORM(db, lggr, config), keyStore, config)
		eventuallyA := cltest.NewAwaiter()
		serviceA1 := new(mocks.ServiceCtx)
		serviceA2 := new(mocks.ServiceCtx)
		serviceA1.On("Start", mock.Anything).Return(nil).Once()
		serviceA2.On("Start", mock.Anything).Return(nil).Once().Run(func(mock.Arguments) { eventuallyA.ItHappened() })
		dA := ocr.NewDelegate(nil, orm, nil, nil, nil, monitoringEndpoint, cc, logger.TestLogger(t), config)
		delegateA := &delegate{jobA.Type, []job.ServiceCtx{serviceA1, serviceA2}, 0, make(chan struct{}), dA}
		eventuallyB := cltest.NewAwaiter()
		serviceB1 := new(mocks.ServiceCtx)
		serviceB2 := new(mocks.ServiceCtx)
		serviceB1.On("Start", mock.Anything).Return(nil).Once()
		serviceB2.On("Start", mock.Anything).Return(nil).Once().Run(func(mock.Arguments) { eventuallyB.ItHappened() })

		dB := ocr.NewDelegate(nil, orm, nil, nil, nil, monitoringEndpoint, cc, logger.TestLogger(t), config)
		delegateB := &delegate{jobB.Type, []job.ServiceCtx{serviceB1, serviceB2}, 0, make(chan struct{}), dB}
		spawner := job.NewSpawner(orm, config, map[job.Type]job.Delegate{
			jobA.Type: delegateA,
			jobB.Type: delegateB,
		}, db, lggr, nil)
		spawner.Start(testutils.Context(t))
		err := spawner.CreateJob(jobA)
		require.NoError(t, err)
		jobSpecIDA := jobA.ID
		delegateA.jobID = jobSpecIDA
		close(delegateA.chContinueCreatingServices)

		eventuallyA.AwaitOrFail(t, 20*time.Second)
		mock.AssertExpectationsForObjects(t, serviceA1, serviceA2)

		err = spawner.CreateJob(jobB)
		require.NoError(t, err)
		jobSpecIDB := jobB.ID
		delegateB.jobID = jobSpecIDB
		close(delegateB.chContinueCreatingServices)

		eventuallyB.AwaitOrFail(t, 20*time.Second)
		mock.AssertExpectationsForObjects(t, serviceB1, serviceB2)

		serviceA1.On("Close").Return(nil).Once()
		serviceA2.On("Close").Return(nil).Once()
		err = spawner.DeleteJob(jobSpecIDA)
		require.NoError(t, err)

		serviceB1.On("Close").Return(nil).Once()
		serviceB2.On("Close").Return(nil).Once()
		err = spawner.DeleteJob(jobSpecIDB)
		require.NoError(t, err)

		require.NoError(t, spawner.Close())
		serviceA1.AssertExpectations(t)
		serviceA2.AssertExpectations(t)
		serviceB1.AssertExpectations(t)
		serviceB2.AssertExpectations(t)
	})

	clearDB(t, db)

	t.Run("starts and stops job services from the DB when .Start()/.Stop() is called", func(t *testing.T) {
		jobA := makeOCRJobSpec(t, address, bridge.Name.String(), bridge2.Name.String())

		eventually := cltest.NewAwaiter()
		serviceA1 := new(mocks.ServiceCtx)
		serviceA2 := new(mocks.ServiceCtx)
		serviceA1.On("Start", mock.Anything).Return(nil).Once()
		serviceA2.On("Start", mock.Anything).Return(nil).Once().Run(func(mock.Arguments) { eventually.ItHappened() })

		lggr := logger.TestLogger(t)
		orm := job.NewTestORM(t, db, cc, pipeline.NewORM(db, lggr, config), keyStore, config)
		d := ocr.NewDelegate(nil, orm, nil, nil, nil, monitoringEndpoint, cc, logger.TestLogger(t), config)
		delegateA := &delegate{jobA.Type, []job.ServiceCtx{serviceA1, serviceA2}, 0, nil, d}
		spawner := job.NewSpawner(orm, config, map[job.Type]job.Delegate{
			jobA.Type: delegateA,
		}, db, lggr, nil)

		err := orm.CreateJob(jobA)
		require.NoError(t, err)
		delegateA.jobID = jobA.ID

		spawner.Start(testutils.Context(t))

		eventually.AwaitOrFail(t)
		mock.AssertExpectationsForObjects(t, serviceA1, serviceA2)

		serviceA1.On("Close").Return(nil).Once()
		serviceA2.On("Close").Return(nil).Once()

		require.NoError(t, spawner.Close())

		mock.AssertExpectationsForObjects(t, serviceA1, serviceA2)
	})

	clearDB(t, db)

	t.Run("closes job services on 'DeleteJob()'", func(t *testing.T) {
		jobA := makeOCRJobSpec(t, address, bridge.Name.String(), bridge2.Name.String())

		eventuallyStart := cltest.NewAwaiter()
		serviceA1 := new(mocks.ServiceCtx)
		serviceA2 := new(mocks.ServiceCtx)
		serviceA1.On("Start", mock.Anything).Return(nil).Once()
		serviceA2.On("Start", mock.Anything).Return(nil).Once().Run(func(mock.Arguments) { eventuallyStart.ItHappened() })

		lggr := logger.TestLogger(t)
		orm := job.NewTestORM(t, db, cc, pipeline.NewORM(db, lggr, config), keyStore, config)
		d := ocr.NewDelegate(nil, orm, nil, nil, nil, monitoringEndpoint, cc, logger.TestLogger(t), config)
		delegateA := &delegate{jobA.Type, []job.ServiceCtx{serviceA1, serviceA2}, 0, nil, d}
		spawner := job.NewSpawner(orm, config, map[job.Type]job.Delegate{
			jobA.Type: delegateA,
		}, db, lggr, nil)

		err := orm.CreateJob(jobA)
		require.NoError(t, err)
		jobSpecIDA := jobA.ID
		delegateA.jobID = jobSpecIDA

		spawner.Start(testutils.Context(t))
		defer spawner.Close()

		eventuallyStart.AwaitOrFail(t)

		// Wait for the claim lock to be taken
		gomega.NewWithT(t).Eventually(func() bool {
			jobs := spawner.ActiveJobs()
			_, exists := jobs[jobSpecIDA]
			return exists
		}, cltest.WaitTimeout(t), cltest.DBPollingInterval).Should(gomega.Equal(true))

		eventuallyClose := cltest.NewAwaiter()
		serviceA1.On("Close").Return(nil).Once()
		serviceA2.On("Close").Return(nil).Once().Run(func(mock.Arguments) { eventuallyClose.ItHappened() })

		err = spawner.DeleteJob(jobSpecIDA)
		require.NoError(t, err)

		eventuallyClose.AwaitOrFail(t)

		// Wait for the claim lock to be released
		gomega.NewWithT(t).Eventually(func() bool {
			jobs := spawner.ActiveJobs()
			_, exists := jobs[jobSpecIDA]
			return exists
		}, cltest.WaitTimeout(t), cltest.DBPollingInterval).Should(gomega.Equal(false))

		mock.AssertExpectationsForObjects(t, serviceA1, serviceA2)
	})
}

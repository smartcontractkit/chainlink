package job_test

import (
	"context"
	"testing"
	"time"

	"github.com/onsi/gomega"
	"gorm.io/gorm"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/job/mocks"
	"github.com/smartcontractkit/chainlink/core/services/offchainreporting"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"gopkg.in/guregu/null.v4"
)

type delegate struct {
	jobType                    job.Type
	services                   []job.Service
	jobID                      int32
	chContinueCreatingServices chan struct{}
	job.Delegate
}

func (d delegate) JobType() job.Type {
	return d.jobType
}

func (d delegate) ServicesForSpec(js job.Job) ([]job.Service, error) {
	if js.Type != d.jobType {
		return nil, nil
	}
	return d.services, nil
}

func clearDB(t *testing.T, db *gorm.DB) {
	err := db.Exec(`TRUNCATE jobs, pipeline_runs, pipeline_specs, pipeline_task_runs CASCADE`).Error
	require.NoError(t, err)
}

func TestSpawner_CreateJobDeleteJob(t *testing.T) {
	config := cltest.NewTestGeneralConfig(t)
	db := pgtest.NewGormDB(t)
	config.SetDB(db)
	keyStore := cltest.NewKeyStore(t, db)
	ethKeyStore := keyStore.Eth()
	keyStore.OCR().Add(cltest.DefaultOCRKey)
	keyStore.P2P().Add(cltest.DefaultP2PKey)

	_, address := cltest.MustInsertRandomKey(t, ethKeyStore)
	_, bridge := cltest.NewBridgeType(t, "voter_turnout", "http://blah.com")
	require.NoError(t, db.Create(bridge).Error)
	_, bridge2 := cltest.NewBridgeType(t, "election_winner", "http://blah.com")
	require.NoError(t, db.Create(bridge2).Error)

	ethClient, _, _ := cltest.NewEthMocksWithDefaultChain(t)
	ethClient.On("CallContext", mock.Anything, mock.Anything, "eth_getBlockByNumber", mock.Anything, false).
		Run(func(args mock.Arguments) {
			head := args.Get(1).(**eth.Head)
			*head = cltest.Head(10)
		}).
		Return(nil)
	txm := postgres.NewGormTransactionManager(db)
	cc := evmtest.NewChainSet(t, evmtest.TestChainOpts{DB: db, Client: ethClient, GeneralConfig: config})

	t.Run("starts and stops job services when jobs are added and removed", func(t *testing.T) {
		jobSpecA := cltest.MakeDirectRequestJobSpec(t)
		jobSpecB := makeOCRJobSpec(t, address)

		orm := job.NewORM(db, cc, pipeline.NewORM(db), keyStore)
		defer orm.Close()
		eventuallyA := cltest.NewAwaiter()
		serviceA1 := new(mocks.Service)
		serviceA2 := new(mocks.Service)
		serviceA1.On("Start").Return(nil).Once()
		serviceA2.On("Start").Return(nil).Once().Run(func(mock.Arguments) { eventuallyA.ItHappened() })
		dA := offchainreporting.NewDelegate(nil, orm, nil, nil, nil, monitoringEndpoint, cc, logger.TestLogger(t))
		delegateA := &delegate{jobSpecA.Type, []job.Service{serviceA1, serviceA2}, 0, make(chan struct{}), dA}
		eventuallyB := cltest.NewAwaiter()
		serviceB1 := new(mocks.Service)
		serviceB2 := new(mocks.Service)
		serviceB1.On("Start").Return(nil).Once()
		serviceB2.On("Start").Return(nil).Once().Run(func(mock.Arguments) { eventuallyB.ItHappened() })

		dB := offchainreporting.NewDelegate(nil, orm, nil, nil, nil, monitoringEndpoint, cc, logger.TestLogger(t))
		delegateB := &delegate{jobSpecB.Type, []job.Service{serviceB1, serviceB2}, 0, make(chan struct{}), dB}
		spawner := job.NewSpawner(orm, config, map[job.Type]job.Delegate{
			jobSpecA.Type: delegateA,
			jobSpecB.Type: delegateB,
		}, txm)
		spawner.Start()
		jobA, err := spawner.CreateJob(context.Background(), *jobSpecA, null.String{})
		require.NoError(t, err)
		jobSpecIDA := jobA.ID
		delegateA.jobID = jobSpecIDA
		close(delegateA.chContinueCreatingServices)

		eventuallyA.AwaitOrFail(t, 20*time.Second)
		mock.AssertExpectationsForObjects(t, serviceA1, serviceA2)

		jobB, err := spawner.CreateJob(context.Background(), *jobSpecB, null.String{})
		require.NoError(t, err)
		jobSpecIDB := jobB.ID
		delegateB.jobID = jobSpecIDB
		close(delegateB.chContinueCreatingServices)

		eventuallyB.AwaitOrFail(t, 20*time.Second)
		mock.AssertExpectationsForObjects(t, serviceB1, serviceB2)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		serviceA1.On("Close").Return(nil).Once()
		serviceA2.On("Close").Return(nil).Once()
		err = spawner.DeleteJob(ctx, jobSpecIDA)
		require.NoError(t, err)

		serviceB1.On("Close").Return(nil).Once()
		serviceB2.On("Close").Return(nil).Once()
		err = spawner.DeleteJob(ctx, jobSpecIDB)
		require.NoError(t, err)

		require.NoError(t, spawner.Close())
		serviceA1.AssertExpectations(t)
		serviceA2.AssertExpectations(t)
		serviceB1.AssertExpectations(t)
		serviceB2.AssertExpectations(t)
	})

	clearDB(t, db)

	t.Run("starts and stops job services from the DB when .Start()/.Stop() is called", func(t *testing.T) {
		jobSpecA := makeOCRJobSpec(t, address)

		eventually := cltest.NewAwaiter()
		serviceA1 := new(mocks.Service)
		serviceA2 := new(mocks.Service)
		serviceA1.On("Start").Return(nil).Once()
		serviceA2.On("Start").Return(nil).Once().Run(func(mock.Arguments) { eventually.ItHappened() })

		orm := job.NewORM(db, cc, pipeline.NewORM(db), keyStore)
		defer orm.Close()
		d := offchainreporting.NewDelegate(nil, orm, nil, nil, nil, monitoringEndpoint, cc, logger.TestLogger(t))
		delegateA := &delegate{jobSpecA.Type, []job.Service{serviceA1, serviceA2}, 0, nil, d}
		spawner := job.NewSpawner(orm, config, map[job.Type]job.Delegate{
			jobSpecA.Type: delegateA,
		}, txm)

		jobA, err := orm.CreateJob(context.Background(), jobSpecA, jobSpecA.Pipeline)
		require.NoError(t, err)
		delegateA.jobID = jobA.ID

		spawner.Start()

		eventually.AwaitOrFail(t)
		mock.AssertExpectationsForObjects(t, serviceA1, serviceA2)

		serviceA1.On("Close").Return(nil).Once()
		serviceA2.On("Close").Return(nil).Once()

		require.NoError(t, spawner.Close())

		mock.AssertExpectationsForObjects(t, serviceA1, serviceA2)
	})

	clearDB(t, db)

	t.Run("closes job services on 'DeleteJob()'", func(t *testing.T) {
		jobSpecA := makeOCRJobSpec(t, address)

		eventuallyStart := cltest.NewAwaiter()
		serviceA1 := new(mocks.Service)
		serviceA2 := new(mocks.Service)
		serviceA1.On("Start").Return(nil).Once()
		serviceA2.On("Start").Return(nil).Once().Run(func(mock.Arguments) { eventuallyStart.ItHappened() })

		orm := job.NewORM(db, cc, pipeline.NewORM(db), keyStore)
		defer orm.Close()
		d := offchainreporting.NewDelegate(nil, orm, nil, nil, nil, monitoringEndpoint, cc, logger.TestLogger(t))
		delegateA := &delegate{jobSpecA.Type, []job.Service{serviceA1, serviceA2}, 0, nil, d}
		spawner := job.NewSpawner(orm, config, map[job.Type]job.Delegate{
			jobSpecA.Type: delegateA,
		}, txm)

		jobA, err := orm.CreateJob(context.Background(), jobSpecA, jobSpecA.Pipeline)
		require.NoError(t, err)
		jobSpecIDA := jobA.ID
		delegateA.jobID = jobSpecIDA

		spawner.Start()
		defer spawner.Close()

		eventuallyStart.AwaitOrFail(t)

		// Wait for the claim lock to be taken
		gomega.NewGomegaWithT(t).Eventually(func() bool {
			jobs := spawner.ActiveJobs()
			_, exists := jobs[jobSpecIDA]
			return exists
		}, cltest.DBWaitTimeout, cltest.DBPollingInterval).Should(gomega.Equal(true))

		eventuallyClose := cltest.NewAwaiter()
		serviceA1.On("Close").Return(nil).Once()
		serviceA2.On("Close").Return(nil).Once().Run(func(mock.Arguments) { eventuallyClose.ItHappened() })

		err = spawner.DeleteJob(context.Background(), jobSpecIDA)
		require.NoError(t, err)

		eventuallyClose.AwaitOrFail(t)

		// Wait for the claim lock to be released
		gomega.NewGomegaWithT(t).Eventually(func() bool {
			jobs := spawner.ActiveJobs()
			_, exists := jobs[jobSpecIDA]
			return exists
		}, cltest.DBWaitTimeout, cltest.DBPollingInterval).Should(gomega.Equal(false))

		mock.AssertExpectationsForObjects(t, serviceA1, serviceA2)
	})
}

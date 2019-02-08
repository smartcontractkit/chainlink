package services_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJobSubscriber_Connect_WithJobs(t *testing.T) {
	t.Parallel()

	store, el, cleanup := cltest.NewJobSubscriber()
	defer cleanup()
	eth := cltest.MockEthOnStore(store)

	j1 := cltest.NewJobWithLogInitiator()
	j2 := cltest.NewJobWithLogInitiator()
	assert.Nil(t, store.CreateJob(&j1))
	assert.Nil(t, store.CreateJob(&j2))
	eth.RegisterSubscription("logs")
	eth.RegisterSubscription("logs")

	assert.Nil(t, el.Connect(cltest.IndexableBlockNumber(1)))
	eth.EventuallyAllCalled(t)
}

func newAddr() common.Address {
	return cltest.NewAddress()
}

func TestJobSubscriber_reconnectLoop_Resubscribing(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
	defer cleanup()
	eth := cltest.MockEthOnStore(store)
	j1 := cltest.NewJobWithLogInitiator()
	j2 := cltest.NewJobWithLogInitiator()
	assert.Nil(t, store.CreateJob(&j1))
	assert.Nil(t, store.CreateJob(&j2))

	eth.RegisterSubscription("logs")
	eth.RegisterSubscription("logs")

	el := services.NewJobSubscriber(store)
	assert.Nil(t, el.Connect(cltest.IndexableBlockNumber(1)))
	assert.Equal(t, 2, len(el.Jobs()))
	el.Disconnect()
	assert.Equal(t, 0, len(el.Jobs()))

	eth.RegisterSubscription("logs")
	eth.RegisterSubscription("logs")
	assert.Nil(t, el.Connect(cltest.IndexableBlockNumber(2)))
	assert.Equal(t, 2, len(el.Jobs()))
	el.Disconnect()
	assert.Equal(t, 0, len(el.Jobs()))
	eth.EventuallyAllCalled(t)
}

func TestJobSubscriber_AttachedToHeadTracker(t *testing.T) {
	t.Parallel()
	g := gomega.NewGomegaWithT(t)

	store, el, cleanup := cltest.NewJobSubscriber()
	defer cleanup()
	eth := cltest.MockEthOnStore(store)
	j1 := cltest.NewJobWithLogInitiator()
	j2 := cltest.NewJobWithLogInitiator()
	assert.Nil(t, store.CreateJob(&j1))
	assert.Nil(t, store.CreateJob(&j2))

	eth.RegisterSubscription("logs")
	eth.RegisterSubscription("logs")

	ht := services.NewHeadTracker(store)
	assert.Nil(t, ht.Start())
	id := ht.Attach(el)
	g.Eventually(func() int { return len(el.Jobs()) }).Should(gomega.Equal(2))
	eth.EventuallyAllCalled(t)

	ht.Detach(id)
	assert.Equal(t, 0, len(el.Jobs()))
}

func TestJobSubscriber_AddJob_Listening(t *testing.T) {
	t.Parallel()
	sharedAddr := newAddr()
	noAddr := common.Address{}

	tests := []struct {
		name      string
		initType  string
		initrAddr common.Address
		logAddr   common.Address
		wantCount int
		topic0    common.Hash
		data      hexutil.Bytes
	}{
		{"ethlog matching address", "ethlog", sharedAddr, sharedAddr, 1, common.Hash{}, hexutil.Bytes{}},
		{"ethlog all address", "ethlog", noAddr, newAddr(), 1, common.Hash{}, hexutil.Bytes{}},
		{"runlog v0 matching address", "runlog", sharedAddr, sharedAddr, 1, models.RunLogTopic0, cltest.StringToVersionedLogData0(t, "id", `{"value":"100"}`)},
		{"runlog v0 matching address", "runlog", sharedAddr, sharedAddr, 1, models.RunLogTopic0, cltest.StringToVersionedLogData0(t, "id", `{"value":"100"}`)},
		{"runlog v20190123 w/o address", "runlog", noAddr, newAddr(), 1, models.RunLogTopic20190123, cltest.StringToVersionedLogData20190123(t, "id", `{"value":"100"}`)},
		{"runlog v20190123 matching address", "runlog", sharedAddr, sharedAddr, 1, models.RunLogTopic20190123, cltest.StringToVersionedLogData20190123(t, "id", `{"value":"100"}`)},
		{"runlog w non-matching topic", "runlog", sharedAddr, sharedAddr, 0, common.Hash{}, cltest.StringToVersionedLogData20190123(t, "id", `{"value":"100"}`)},
		{"runlog v20190207 w/o address", "runlog", noAddr, newAddr(), 1, models.RunLogTopic20190207, cltest.StringToVersionedLogData20190207(t, "id", `{"value":"100"}`)},
		{"runlog v20190207 matching address", "runlog", sharedAddr, sharedAddr, 1, models.RunLogTopic20190207, cltest.StringToVersionedLogData20190207(t, "id", `{"value":"100"}`)},
		{"runlog w non-matching topic", "runlog", sharedAddr, sharedAddr, 0, common.Hash{}, cltest.StringToVersionedLogData20190207(t, "id", `{"value":"100"}`)},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store, el, cleanup := cltest.NewJobSubscriber()
			defer cleanup()

			eth := cltest.MockEthOnStore(store)
			logChan := make(chan models.Log, 1)
			eth.RegisterSubscription("logs", logChan)

			job := cltest.NewJob()
			initr := models.Initiator{Type: test.initType}
			initr.Address = test.initrAddr
			job.Initiators = []models.Initiator{initr}
			require.NoError(t, store.CreateJob(&job))
			el.AddJob(job, cltest.IndexableBlockNumber(1))

			ht := services.NewHeadTracker(store)
			ht.Attach(el)
			require.NoError(t, ht.Start())

			logChan <- models.Log{
				Address: test.logAddr,
				Data:    test.data,
				Topics: []common.Hash{
					test.topic0,
					cltest.StringToHash(job.ID),
					newAddr().Hash(),
					common.BigToHash(big.NewInt(0)),
				},
			}

			cltest.WaitForRuns(t, job, store, test.wantCount)

			eth.EventuallyAllCalled(t)
		})
	}
}

func TestJobSubscriber_OnNewHead_OnlyResumePendingConfirmations(t *testing.T) {
	t.Parallel()

	block := cltest.NewBlockHeader(10)

	tests := []struct {
		status   models.RunStatus
		wantSend bool
	}{
		{models.RunStatusPendingConfirmations, true},
		{models.RunStatusInProgress, false},
		{models.RunStatusPendingBridge, false},
		{models.RunStatusPendingSleep, false},
		{models.RunStatusCompleted, false},
	}

	for _, test := range tests {
		t.Run(string(test.status), func(t *testing.T) {
			store, js, cleanup := cltest.NewJobSubscriber()
			defer cleanup()

			mockRunChannel := cltest.NewMockRunChannel()
			store.RunChannel = mockRunChannel

			job := cltest.NewJobWithWebInitiator()
			require.NoError(t, store.CreateJob(&job))
			initr := job.Initiators[0]
			run := job.NewRun(initr)
			run = run.ApplyResult(models.RunResult{Status: test.status})
			assert.Nil(t, store.CreateJobRun(&run))

			js.OnNewHead(block)
			if test.wantSend {
				assert.Equal(t, 1, len(mockRunChannel.Runs))
			} else {
				assert.Equal(t, 0, len(mockRunChannel.Runs))
			}
		})
	}
}

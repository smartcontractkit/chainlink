package services_test

import (
	"strings"
	"sync/atomic"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/services"
	strpkg "github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/stretchr/testify/assert"
)

func TestInitiatorSubscriptionLogEvent_RunLogJSON(t *testing.T) {
	t.Parallel()

	clData := cltest.JSONFromString(`{"url":"https://etherprice.com/api","path":["recent","usd"],"address":"0x3cCad4715152693fE3BC4460591e3D3Fbd071b42","dataPrefix":"0x0000000000000000000000000000000000000000000000000000000000000017","functionSelector":"0x76005c26"}`)

	hwLog := cltest.LogFromFixture("../internal/fixtures/eth/subscription_logs_hello_world.json")
	tests := []struct {
		name        string
		el          strpkg.Log
		wantErrored bool
		wantData    models.JSON
	}{
		{"hello world", hwLog, false, clData},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			le := services.InitiatorSubscriptionLogEvent{Log: test.el}
			output, err := le.RunLogJSON()
			assert.JSONEq(t, strings.ToLower(test.wantData.String()), strings.ToLower(output.String()))
			assert.NoError(t, err)
			assert.Equal(t, test.wantErrored, (err != nil))
		})
	}
}

func TestInitiatorSubscriptionLogEvent_EthLogJSON(t *testing.T) {
	t.Parallel()

	hwLog := cltest.LogFromFixture("../internal/fixtures/eth/subscription_logs_hello_world.json")
	exampleLog := cltest.LogFromFixture("../internal/fixtures/eth/subscription_logs.json")
	tests := []struct {
		name        string
		el          strpkg.Log
		wantErrored bool
		wantData    models.JSON
	}{
		{"example", exampleLog, false, cltest.JSONResultFromFixture("../internal/fixtures/eth/subscription_logs.json")},
		{"hello world", hwLog, false, cltest.JSONResultFromFixture("../internal/fixtures/eth/subscription_logs_hello_world.json")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			le := services.InitiatorSubscriptionLogEvent{Log: test.el}
			output, err := le.EthLogJSON()
			assert.JSONEq(t, strings.ToLower(test.wantData.String()), strings.ToLower(output.String()))
			assert.Equal(t, test.wantErrored, (err != nil))
		})
	}
}

func TestServices_NewInitiatorSubscription_BackfillLogs(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
	defer cleanup()
	eth := cltest.MockEthOnStore(store)

	job, initr := cltest.NewJobWithLogInitiator()
	log := cltest.LogFromFixture("../internal/fixtures/eth/subscription_logs.json")
	eth.Register("eth_getLogs", []strpkg.Log{log})
	eth.RegisterSubscription("logs")

	var count int32
	callback := func(services.InitiatorSubscriptionLogEvent) { atomic.AddInt32(&count, 1) }
	head := cltest.IndexableBlockNumber(0)
	filter := services.NewInitiatorFilterQuery(initr, head, nil)
	sub, err := services.NewInitiatorSubscription(initr, job, store, filter, callback)
	assert.NoError(t, err)
	defer sub.Unsubscribe()

	eth.EventuallyAllCalled(t)

	gomega.NewGomegaWithT(t).Eventually(func() int32 {
		return atomic.LoadInt32(&count)
	}).Should(gomega.Equal(int32(1)))
}

func TestServices_NewInitiatorSubscription_BackfillLogs_WithNoHead(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
	defer cleanup()
	eth := cltest.MockEthOnStore(store)

	job, initr := cltest.NewJobWithLogInitiator()
	eth.RegisterSubscription("logs")

	var count int32
	callback := func(services.InitiatorSubscriptionLogEvent) { atomic.AddInt32(&count, 1) }
	filter := services.NewInitiatorFilterQuery(initr, nil, nil)
	sub, err := services.NewInitiatorSubscription(initr, job, store, filter, callback)
	assert.NoError(t, err)
	defer sub.Unsubscribe()

	eth.EventuallyAllCalled(t)
	assert.Equal(t, int32(0), atomic.LoadInt32(&count))
}

func TestServices_NewInitiatorSubscription_PreventsDoubleDispatch(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
	defer cleanup()
	eth := cltest.MockEthOnStore(store)

	job, initr := cltest.NewJobWithLogInitiator()
	log := cltest.LogFromFixture("../internal/fixtures/eth/subscription_logs.json")
	eth.Register("eth_getLogs", []strpkg.Log{log}) // backfill
	logsChan := make(chan strpkg.Log)
	eth.RegisterSubscription("logs", logsChan)

	var count int32
	callback := func(services.InitiatorSubscriptionLogEvent) { atomic.AddInt32(&count, 1) }
	head := cltest.IndexableBlockNumber(0)
	filter := services.NewInitiatorFilterQuery(initr, head, nil)
	sub, err := services.NewInitiatorSubscription(initr, job, store, filter, callback)
	assert.NoError(t, err)
	defer sub.Unsubscribe()

	// Add the same original log
	logsChan <- log
	// Add a log after the repeated log to make sure it gets processed
	log2 := cltest.LogFromFixture("../internal/fixtures/eth/subscription_logs_hello_world.json")
	logsChan <- log2

	eth.EventuallyAllCalled(t)
	g := gomega.NewGomegaWithT(t)
	g.Eventually(func() int32 { return atomic.LoadInt32(&count) }).Should(gomega.Equal(int32(2)))
}

func TestTopicFiltersForRunLog(t *testing.T) {
	t.Parallel()

	jobID := "4a1eb0e8df314cb894024a38991cff0f"
	topics := services.TopicFiltersForRunLog(jobID)

	assert.Equal(t, 2, len(topics))
	assert.Equal(
		t,
		[]common.Hash{services.RunLogTopic},
		topics[services.RunLogTopicSignature])

	assert.Equal(
		t,
		[]common.Hash{
			common.HexToHash("0x3461316562306538646633313463623839343032346133383939316366663066"),
			common.HexToHash("0x4a1eb0e8df314cb894024a38991cff0f00000000000000000000000000000000"),
		},
		topics[1])
}

func TestInitiatorSubscriptionLogEvent_Requester(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input common.Hash
		want  common.Address
	}{
		{"basic",
			common.HexToHash("0x00000000000000000000000059b15a7ae74c803cc151ffe63042faa826c96eee"),
			common.HexToAddress("0x59b15a7ae74c803cc151ffe63042faa826c96eee"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rl := cltest.NewRunLog("id", cltest.NewAddress(), cltest.NewAddress(), 0, "{}")
			rl.Topics[services.RunLogTopicRequester] = test.input
			le := services.InitiatorSubscriptionLogEvent{Log: rl}

			assert.Equal(t, test.want, le.Requester())
		})
	}
}

func TestInitiatorSubscriptionLogEvent_ValidateRunLog(t *testing.T) {
	t.Parallel()

	job := cltest.NewJob()
	job.ID = "4a1eb0e8df314cb894024a38991cff0f"

	noRequesters := []common.Address{}
	permittedAddr := cltest.NewAddress()
	unpermittedAddr := cltest.NewAddress()
	requesterList := []common.Address{permittedAddr}

	tests := []struct {
		name                string
		eventLogTopic       common.Hash
		jobIDTopic          common.Hash
		initiatorRequesters []common.Address
		requesterAddress    common.Address
		want                bool
	}{
		{"not runlog", cltest.StringToHash("notrunlog"), common.Hash{}, noRequesters, unpermittedAddr, false},
		{"runlog wrong jobid", services.RunLogTopic, cltest.StringToHash("wrongjob"), noRequesters, unpermittedAddr, false},
		{"runlog proper hex jobid", services.RunLogTopic, cltest.StringToHash(job.ID), noRequesters, unpermittedAddr, true},
		{"runlog incorrect encoded jobid", services.RunLogTopic, common.HexToHash("0x4a1eb0e8df314cb894024a38991cff0f00000000000000000000000000000000"), noRequesters, unpermittedAddr, true},
		{"runlog correct requester", services.RunLogTopic, cltest.StringToHash(job.ID), requesterList, permittedAddr, true},
		{"runlog incorrect requester", services.RunLogTopic, cltest.StringToHash(job.ID), requesterList, unpermittedAddr, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			log := cltest.NewRunLog(job.ID, cltest.NewAddress(), test.requesterAddress, 1, "{}")
			log.Topics = []common.Hash{
				test.eventLogTopic,
				test.jobIDTopic,
				test.requesterAddress.Hash(),
				common.Hash{},
			}

			le := services.InitiatorSubscriptionLogEvent{
				Job: job,
				Log: log,
				Initiator: models.Initiator{
					InitiatorParams: models.InitiatorParams{
						Requesters: test.initiatorRequesters,
					},
				},
			}

			assert.Equal(t, test.want, le.ValidateRunLog())
		})
	}
}

func TestStartRunLogSubscription_ValidateSenders(t *testing.T) {
	requester := cltest.NewAddress()

	tests := []struct {
		name      string
		requester common.Address
		status    models.RunStatus
	}{
		{"runlog contains valid requester", requester, models.RunStatusCompleted},
		{"runlog has wrong requester", cltest.NewAddress(), models.RunStatusErrored},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config, _ := cltest.NewConfigWithPrivateKey()
			app, cleanup := cltest.NewApplicationWithConfigAndUnlockedAccount(config)
			defer cleanup()

			eth := app.MockEthClient()
			logs := make(chan strpkg.Log, 1)
			eth.Context("app.Start()", func(eth *cltest.EthMock) {
				eth.Register("eth_getBlockByNumber", models.BlockHeader{})
				eth.Register("eth_getTransactionCount", "0x1")
				eth.RegisterSubscription("logs", logs)
			})
			assert.NoError(t, app.Start())

			js, initr := cltest.NewJobWithRunLogInitiator()
			initr.Requesters = []common.Address{requester}
			_, err := services.StartRunLogSubscription(initr, js, nil, app.Store)
			assert.NoError(t, err)

			logs <- cltest.NewRunLog(js.ID, cltest.NewAddress(), test.requester, 1, `{}`)
			eth.EventuallyAllCalled(t)

			gomega.NewGomegaWithT(t).Eventually(func() models.RunStatus {
				var run models.JobRun
				app.Store.One("JobID", js.ID, &run)
				return run.Status
			}).Should(gomega.Equal(test.status))
		})
	}
}

func TestRunTopic(t *testing.T) {
	assert.Equal(t, common.HexToHash("0x6d6db1f8fe19d95b1d0fa6a4bce7bb24fbf84597b35a33ff95521fac453c1529"), services.RunLogTopic)
}

package services_test

import (
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/stretchr/testify/assert"
)

func TestServices_RpcLogEvent_RunLogJSON(t *testing.T) {
	t.Parallel()

	clData := cltest.JSONFromString(`{"url":"https://etherprice.com/api","path":["recent","usd"],"address":"0x3cCad4715152693fE3BC4460591e3D3Fbd071b42","dataPrefix":"0x0000000000000000000000000000000000000000000000000000000000000001","functionSelector":"76005c26"}`)

	hwLog := cltest.LogFromFixture("../internal/fixtures/eth/subscription_logs_hello_world.json")
	tests := []struct {
		name        string
		el          types.Log
		wantErrored bool
		wantData    models.JSON
	}{
		{"hello world", hwLog, false, clData},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			le := services.RPCLogEvent{Log: test.el}
			output, err := le.RunLogJSON()
			assert.JSONEq(t, strings.ToLower(test.wantData.String()), strings.ToLower(output.String()))
			assert.Nil(t, err)
			assert.Equal(t, test.wantErrored, (err != nil))
		})
	}
}

func TestServices_RpcLogEvent_EthLogJSON(t *testing.T) {
	hwLog := cltest.LogFromFixture("../internal/fixtures/eth/subscription_logs_hello_world.json")
	exampleLog := cltest.LogFromFixture("../internal/fixtures/eth/subscription_logs.json")
	tests := []struct {
		name        string
		el          types.Log
		wantErrored bool
		wantData    models.JSON
	}{
		{"example", exampleLog, false, cltest.JSONResultFromFixture("../internal/fixtures/eth/subscription_logs.json")},
		{"hello world", hwLog, false, cltest.JSONResultFromFixture("../internal/fixtures/eth/subscription_logs_hello_world.json")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			le := services.RPCLogEvent{Log: test.el}
			output, err := le.EthLogJSON()
			assert.JSONEq(t, strings.ToLower(test.wantData.String()), strings.ToLower(output.String()))
			assert.Equal(t, test.wantErrored, (err != nil))
		})
	}
}

func TestServices_NewRPCLogSubscription_BackfillLogs(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
	defer cleanup()
	eth := cltest.MockEthOnStore(store)

	job, initr := cltest.NewJobWithLogInitiator()
	log := cltest.LogFromFixture("../internal/fixtures/eth/subscription_logs.json")
	eth.Register("eth_getLogs", []types.Log{log})
	eth.RegisterSubscription("logs")

	count := 0
	callback := func(services.RPCLogEvent) { count += 1 }
	head := cltest.IndexableBlockNumber(0)
	subscriber := services.NewRPCLogSubscriber(initr, head, nil, callback)
	sub, err := services.NewRPCLogSubscription(initr, job, store, subscriber)
	assert.Nil(t, err)
	defer sub.Unsubscribe()

	eth.EventuallyAllCalled(t)
	assert.Equal(t, 1, count)
}

func TestServices_NewRPCLogSubscription_BackfillLogs_WithNoHead(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
	defer cleanup()
	eth := cltest.MockEthOnStore(store)

	job, initr := cltest.NewJobWithLogInitiator()
	eth.RegisterSubscription("logs")

	count := 0
	callback := func(services.RPCLogEvent) { count += 1 }
	subscriber := services.NewRPCLogSubscriber(initr, nil, nil, callback)
	sub, err := services.NewRPCLogSubscription(initr, job, store, subscriber)
	assert.Nil(t, err)
	defer sub.Unsubscribe()

	eth.EventuallyAllCalled(t)
	assert.Equal(t, 0, count)
}

func TestServices_NewRPCLogSubscription_PreventsDoubleDispatch(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
	defer cleanup()
	eth := cltest.MockEthOnStore(store)

	job, initr := cltest.NewJobWithLogInitiator()
	log := cltest.LogFromFixture("../internal/fixtures/eth/subscription_logs.json")
	eth.Register("eth_getLogs", []types.Log{log}) // backfill
	logsChan := make(chan types.Log, 1)
	eth.RegisterSubscription("logs", logsChan)
	logsChan <- log // received in real time

	count := 0
	callback := func(services.RPCLogEvent) { count += 1 }
	head := cltest.IndexableBlockNumber(0)
	subscriber := services.NewRPCLogSubscriber(initr, head, nil, callback)
	sub, err := services.NewRPCLogSubscription(initr, job, store, subscriber)
	assert.Nil(t, err)
	defer sub.Unsubscribe()

	eth.EventuallyAllCalled(t)
	assert.Equal(t, 1, count)
}

func TestTopicFiltersForRunLog(t *testing.T) {
	t.Parallel()

	jobID := "4a1eb0e8df314cb894024a38991cff0f"
	topics := services.TopicFiltersForRunLog(jobID)

	assert.Equal(t, 3, len(topics))
	assert.Nil(t, topics[1])
	assert.Equal(
		t,
		[]common.Hash{services.RunLogTopic},
		topics[services.EventTopicSignature])

	assert.Equal(
		t,
		[]common.Hash{
			common.HexToHash("0x3461316562306538646633313463623839343032346133383939316366663066"),
			common.HexToHash("0x4a1eb0e8df314cb894024a38991cff0f00000000000000000000000000000000"),
		},
		topics[2])
}

func TestRPCLogEvent_ValidateRunLog(t *testing.T) {
	t.Parallel()

	job := cltest.NewJob()
	job.ID = "4a1eb0e8df314cb894024a38991cff0f"
	tests := []struct {
		name          string
		eventLogTopic common.Hash
		jobIDTopic    common.Hash
		want          bool
	}{
		{"not runlog", cltest.StringToHash("notrunlog"), common.Hash{}, false},
		{"runlog wrong jobid", services.RunLogTopic, cltest.StringToHash("wrongjob"), false},
		{"runlog proper hex jobid", services.RunLogTopic, cltest.StringToHash(job.ID), true},
		{"runlog incorrect encoded jobid", services.RunLogTopic, common.HexToHash("0x4a1eb0e8df314cb894024a38991cff0f00000000000000000000000000000000"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := cltest.NewRunLog(job.ID, cltest.NewAddress(), 1, "{}")
			log.Topics = []common.Hash{tt.eventLogTopic, common.Hash{}, tt.jobIDTopic, common.Hash{}}
			le := services.RPCLogEvent{
				Job: job,
				Log: log,
			}

			assert.Equal(t, tt.want, le.ValidateRunLog())
		})
	}
}

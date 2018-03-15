package services_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/stretchr/testify/assert"
)

func TestServices_RpcLogEvent_RunLogJSON(t *testing.T) {
	t.Parallel()

	var clData models.JSON
	clDataFixture := `{"url":"https://etherprice.com/api","path":["recent","usd"],"address":"0x3cCad4715152693fE3BC4460591e3D3Fbd071b42","dataPrefix":"0x0000000000000000000000000000000000000000000000000000000000000001","functionSelector":"76005c26"}`
	assert.Nil(t, json.Unmarshal([]byte(clDataFixture), &clData))

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

// If updating this test, be sure to update the truffle suite's "expected event signature" test.
func TestServices_RunLogTopic_ExpectedEventSignature(t *testing.T) {
	t.Parallel()

	expected := "0x06f4bf36b4e011a5c499cef1113c2d166800ce4013f6c2509cab1a0e92b83fb2"
	assert.Equal(t, expected, services.RunLogTopic.Hex())
}

func TestServices_NewRPCLogSubscription_BackfillLogs(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
	defer cleanup()
	eth := cltest.MockEthOnStore(store)

	job := cltest.NewJobWithLogInitiator()
	initr := job.Initiators[0]
	log := cltest.LogFromFixture("../internal/fixtures/eth/subscription_logs.json")
	eth.Register("eth_getLogs", []types.Log{log})
	eth.RegisterSubscription("logs")

	count := 0
	callback := func(services.RPCLogEvent) { count += 1 }
	head := cltest.IndexableBlockNumber(0)
	sub, err := services.NewRPCLogSubscription(initr, job, head, store, callback)
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

	job := cltest.NewJobWithLogInitiator()
	initr := job.Initiators[0]
	eth.RegisterSubscription("logs")

	count := 0
	callback := func(services.RPCLogEvent) { count += 1 }
	sub, err := services.NewRPCLogSubscription(initr, job, nil, store, callback)
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

	job := cltest.NewJobWithLogInitiator()
	initr := job.Initiators[0]
	log := cltest.LogFromFixture("../internal/fixtures/eth/subscription_logs.json")
	eth.Register("eth_getLogs", []types.Log{log}) // backfill
	logsChan := make(chan types.Log, 1)
	eth.RegisterSubscription("logs", logsChan)
	logsChan <- log // received in real time

	count := 0
	callback := func(services.RPCLogEvent) { count += 1 }
	head := cltest.IndexableBlockNumber(0)
	sub, err := services.NewRPCLogSubscription(initr, job, head, store, callback)
	assert.Nil(t, err)
	defer sub.Unsubscribe()

	eth.EventuallyAllCalled(t)
	assert.Equal(t, 1, count)
}

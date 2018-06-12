package services_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink/adapters"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/stretchr/testify/assert"
)

func TestSpecAndRunSubscriber_AttachedToHeadTracker(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
	defer cleanup()

	eth := cltest.MockEthOnStore(store)
	logs := make(chan types.Log, 1)
	eth.RegisterSubscription("logs", logs)

	ht := services.NewHeadTracker(store)
	oracleAddress := common.HexToAddress("0x0000000000000000000000000000000000000000")
	sub := services.NewSpecAndRunSubscriber(store, &oracleAddress)
	id := ht.Attach(sub)
	assert.NoError(t, ht.Start())

	json := `{"tasks": ["NoOp"], "params": {"url": "www.lmgtfy.com"}}`
	logs <- cltest.NewSpecAndRunLog(cltest.NewAddress(), 1, json)

	jobs := cltest.WaitForJobs(t, store, 1)
	job := jobs[0]
	runs := cltest.WaitForRuns(t, job, store, 1)
	run := cltest.WaitForJobRunToComplete(t, store, runs[0])

	assert.Equal(t, 1, len(job.Tasks))
	assert.Equal(t, "noop", job.Tasks[0].Type)

	assert.Equal(t, job.ID, run.JobID)
	assert.Equal(t, "www.lmgtfy.com", run.Result.Data.Get("url").String())

	ht.Detach(id)
	eth.EventuallyAllCalled(t)
}

func TestNewSpecAndRunLogEvent_SetsDefaultEthTxParams(t *testing.T) {
	store, cleanup := cltest.NewStore()
	defer cleanup()

	oracleAddress := cltest.NewAddress()
	payload := `{"tasks":["httpget", "EthTx"]}`
	log := cltest.NewSpecAndRunLog(oracleAddress, 1, payload)
	le, err := services.NewSpecAndRunLogEvent(log)
	assert.NoError(t, err)

	otherTask := le.Job.Tasks[0]
	assert.Equal(t, "", otherTask.Params.Get("address").String())

	ethTxAdapter, err := adapters.For(le.Job.Tasks[1], store)
	assert.NoError(t, err)
	ethTx := cltest.UnwrapAdapter(ethTxAdapter).(*adapters.EthTx)
	assert.Equal(t, oracleAddress, ethTx.Address)
	assert.Equal(t, services.OracleFulfillmentFunctionID, ethTx.FunctionSelector.String())
	assert.Equal(t, log.Topics[services.SpecAndRunTopicInternalID].String(), ethTx.DataPrefix.String())
}

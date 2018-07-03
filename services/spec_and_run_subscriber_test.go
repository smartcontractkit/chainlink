package services_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink/adapters"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store/models"
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
	logs <- cltest.NewSpecAndRunLog(cltest.NewAddress(), 1, json, big.NewInt(1))

	jobs := cltest.WaitForJobs(t, store, 1)
	job := jobs[0]
	runs := cltest.WaitForRuns(t, job, store, 1)
	run := cltest.WaitForJobRunToComplete(t, store, runs[0])

	assert.Equal(t, 1, len(job.Tasks))
	assert.Equal(t, "noop", job.Tasks[0].Type.String())

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
	log := cltest.NewSpecAndRunLog(oracleAddress, 1, payload, big.NewInt(1))
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

func TestSpecAndRunLogEvent_StartJob(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
	defer cleanup()
	store.Config.MinimumContractPayment = *big.NewInt(100)
	minPayment := &store.Config.MinimumContractPayment
	subMin := big.NewInt(0).Sub(minPayment, big.NewInt(1))

	tests := []struct {
		name   string
		amount *big.Int
		status models.RunStatus
	}{
		{"enough payment", minPayment, models.RunStatusCompleted},
		{"not enough payment", subMin, models.RunStatusErrored},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			payload := `{"tasks":["noOp"]}`
			log := cltest.NewSpecAndRunLog(cltest.NewAddress(), 1, payload, test.amount)
			le, err := services.NewSpecAndRunLogEvent(log)
			assert.NoError(t, err)
			j := le.Job
			assert.Nil(t, store.SaveJob(&j))
			le.StartJob(store)

			jr := cltest.WaitForRuns(t, j, store, 1)[0]
			cltest.WaitForJobRunStatus(t, store, jr, test.status)
		})
	}
}

package synchronization

import (
	"encoding/json"
	"io/ioutil"
	"math/big"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/assets"
	clnull "github.com/smartcontractkit/chainlink/core/null"
	"github.com/smartcontractkit/chainlink/core/store/models"

	"github.com/ethereum/go-ethereum/common"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
	null "gopkg.in/guregu/null.v4"
)

func TestSyncJobRunPresenter_HappyPath(t *testing.T) {
	newAddress := common.HexToAddress("0x9FBDa871d559710256a2502A2517b794B482Db40")
	requestID := common.HexToHash("0xcafe")
	txHash := common.HexToHash("0xdeadbeef")

	task0RunID := uuid.NewV4()
	task1RunID := uuid.NewV4()

	job := models.JobSpec{ID: models.NewJobID()}
	runRequest := models.RunRequest{
		Payment:   assets.NewLink(2),
		RequestID: &requestID,
		TxHash:    &txHash,
		Requester: &newAddress,
	}
	run := models.MakeJobRun(&job, time.Now(), &models.Initiator{Type: models.InitiatorRunLog}, big.NewInt(0), &runRequest)
	run.TaskRuns = []models.TaskRun{
		{
			ID:                               task0RunID,
			Status:                           models.RunStatusPendingIncomingConfirmations,
			ObservedIncomingConfirmations:    clnull.Uint32From(1),
			MinRequiredIncomingConfirmations: clnull.Uint32From(3),
		},
		{
			ID:                               task1RunID,
			Status:                           models.RunStatusErrored,
			Result:                           models.RunResult{ErrorMessage: null.StringFrom("yikes fam")},
			ObservedIncomingConfirmations:    clnull.Uint32From(1),
			MinRequiredIncomingConfirmations: clnull.Uint32From(3),
		},
	}
	p := SyncJobRunPresenter{JobRun: &run}

	bytes, err := p.MarshalJSON()
	require.NoError(t, err)

	var data map[string]interface{}
	err = json.Unmarshal(bytes, &data)
	require.NoError(t, err)

	assert.Equal(t, data["runId"], run.ID.String())
	assert.Equal(t, data["jobId"], job.ID.String())
	assert.Equal(t, data["status"], "in_progress")
	assert.Contains(t, data, "error")
	assert.Contains(t, data, "createdAt")
	assert.Equal(t, data["payment"], "2")
	assert.Equal(t, data["finishedAt"], nil)
	assert.Contains(t, data, "tasks")

	initiator, ok := data["initiator"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, initiator["type"], "runlog")
	assert.Equal(t, initiator["requestId"], "0x000000000000000000000000000000000000000000000000000000000000cafe")
	assert.Equal(t, initiator["txHash"], "0x00000000000000000000000000000000000000000000000000000000deadbeef")
	assert.Equal(t, initiator["requester"], newAddress.Hex())

	tasks, ok := data["tasks"].([]interface{})
	require.True(t, ok)
	require.Len(t, tasks, 2)
	task0, ok := tasks[0].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, task0["index"], float64(0))
	assert.Contains(t, task0, "type")
	assert.Equal(t, "pending_incoming_confirmations", task0["status"])
	assert.Equal(t, task0["error"], nil)
	assert.Equal(t, float64(1), task0["confirmations"])
	assert.Equal(t, float64(3), task0["minimumConfirmations"])
	task1, ok := tasks[1].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, task1["index"], float64(1))
	assert.Contains(t, task1, "type")
	assert.Equal(t, task1["status"], "errored")
	assert.Equal(t, task1["error"], "yikes fam")
}

func TestSyncJobRunPresenter_Initiators(t *testing.T) {
	newAddress := common.HexToAddress("0x9FBDa871d559710256a2502A2517b794B482Db40")
	requestID := common.HexToHash("0xcafe")
	txHash := common.HexToHash("0xdeadbeef")

	tests := []struct {
		initrType string
		rr        models.RunRequest
		keyCount  int
	}{
		{models.InitiatorWeb, models.RunRequest{}, 1},
		{models.InitiatorCron, models.RunRequest{}, 1},
		{models.InitiatorRunAt, models.RunRequest{}, 1},
		{models.InitiatorEthLog, models.RunRequest{TxHash: &txHash}, 2},
		{
			models.InitiatorRunLog,
			models.RunRequest{
				RequestID: &requestID,
				TxHash:    &txHash,
				Requester: &newAddress,
			},
			4,
		},
	}

	for _, test := range tests {
		t.Run(test.initrType, func(t *testing.T) {
			jobRun := models.JobRun{
				ID:         uuid.NewV4(),
				JobSpecID:  models.NewJobID(),
				Initiator:  models.Initiator{Type: test.initrType},
				RunRequest: test.rr,
			}

			p := SyncJobRunPresenter{JobRun: &jobRun}

			bytes, err := p.MarshalJSON()
			require.NoError(t, err)

			var data map[string]interface{}
			err = json.Unmarshal(bytes, &data)
			require.NoError(t, err)

			initiator, ok := data["initiator"].(map[string]interface{})
			require.True(t, ok)
			assert.Len(t, initiator, test.keyCount)
			assert.Equal(t, initiator["type"], test.initrType)
		})
	}
}

func jsonFromFixture(t *testing.T, path string) models.JSON {
	body, err := ioutil.ReadFile(path)
	require.NoError(t, err)
	j, err := models.ParseJSON(body)
	require.NoError(t, err)
	return j
}

func TestSyncJobRunPresenter_EthTxTask(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{"confirmed", "../../testdata/apiresponses/confirmedEthTxData.json", ""},
		{"safe fulfilled", "../../testdata/apiresponses/fulfilledReceiptResponse.json", "fulfilledRunLog"},
		{"safe not fulfilled", "../../testdata/apiresponses/notFulfilledReceiptResponse.json", "noFulfilledRunLog"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			newAddress := common.HexToAddress("0x9FBDa871d559710256a2502A2517b794B482Db40")
			requestID := common.HexToHash("0xcafe")
			requestTxHash := common.HexToHash("0xdeadbeef")
			dataJSON := jsonFromFixture(t, test.path)
			outgoingTxHash := "0x1111111111111111111111111111111111111111111111111111111111111111"

			taskSpec := models.TaskSpec{
				Type: "ethtx",
			}
			job := models.JobSpec{ID: models.NewJobID()}
			runRequest := models.RunRequest{
				RequestID: &requestID,
				TxHash:    &requestTxHash,
				Requester: &newAddress,
			}
			run := models.MakeJobRun(&job, time.Now(), &models.Initiator{Type: models.InitiatorRunLog}, big.NewInt(0), &runRequest)
			run.SetStatus(models.RunStatusCompleted)
			run.TaskRuns = []models.TaskRun{
				{
					ID:       uuid.NewV4(),
					TaskSpec: taskSpec,
					Status:   models.RunStatusPendingIncomingConfirmations,
					Result:   models.RunResult{Data: dataJSON},
				},
			}
			p := SyncJobRunPresenter{JobRun: &run}

			bytes, err := p.MarshalJSON()
			require.NoError(t, err)

			require.True(t, gjson.ValidBytes(bytes))
			data := gjson.ParseBytes(bytes)

			tasks := data.Get("tasks").Array()
			require.Len(t, tasks, 1)
			task0 := tasks[0].Map()
			assert.Equal(t, task0["index"].Float(), float64(0))
			assert.Contains(t, task0["type"].String(), "ethtx")
			assert.Equal(t, "pending_incoming_confirmations", task0["status"].String())
			assert.Equal(t, task0["error"].Type, gjson.Null)

			txresult := task0["result"].Map()
			assert.Equal(t, test.want, txresult["transactionStatus"].String())
			assert.Equal(t, outgoingTxHash, txresult["transactionHash"].String())
		})
	}
}

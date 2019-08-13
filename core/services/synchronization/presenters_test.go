package synchronization

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	clnull "github.com/smartcontractkit/chainlink/core/null"
	"github.com/smartcontractkit/chainlink/core/store/assets"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
	null "gopkg.in/guregu/null.v3"
)

func TestSyncJobRunPresenter_HappyPath(t *testing.T) {
	newAddress := common.HexToAddress("0x9FBDa871d559710256a2502A2517b794B482Db40")
	requestID := "RequestID"
	txHash := common.HexToHash("0xdeadbeef")

	runID := models.NewID()
	specID := models.NewID()
	task0RunID := models.NewID()
	task1RunID := models.NewID()
	jobRun := models.JobRun{
		ID:        runID,
		JobSpecID: specID,
		Status:    models.RunStatusInProgress,
		Result:    models.RunResult{Amount: assets.NewLink(2)},
		Initiator: models.Initiator{
			Type: models.InitiatorRunLog,
		},
		RunRequest: models.RunRequest{
			RequestID: &requestID,
			TxHash:    &txHash,
			Requester: &newAddress,
		},
		TaskRuns: []models.TaskRun{
			models.TaskRun{
				ID:                   task0RunID,
				Status:               models.RunStatusPendingConfirmations,
				Confirmations:        clnull.Uint32From(1),
				MinimumConfirmations: clnull.Uint32From(3),
			},
			models.TaskRun{
				ID:     task1RunID,
				Status: models.RunStatusErrored,
				Result: models.RunResult{ErrorMessage: null.StringFrom("yikes fam")},
			},
		},
	}
	p := SyncJobRunPresenter{JobRun: &jobRun}

	bytes, err := p.MarshalJSON()
	require.NoError(t, err)

	var data map[string]interface{}
	err = json.Unmarshal(bytes, &data)
	require.NoError(t, err)

	assert.Equal(t, data["runId"], runID.String())
	assert.Equal(t, data["jobId"], specID.String())
	assert.Equal(t, data["status"], "in_progress")
	assert.Contains(t, data, "error")
	assert.Contains(t, data, "createdAt")
	assert.Equal(t, data["amount"], "2")
	assert.Equal(t, data["finishedAt"], nil)
	assert.Contains(t, data, "tasks")

	initiator, ok := data["initiator"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, initiator["type"], "runlog")
	assert.Equal(t, initiator["requestId"], "RequestID")
	assert.Equal(t, initiator["txHash"], "0x00000000000000000000000000000000000000000000000000000000deadbeef")
	assert.Equal(t, initiator["requester"], newAddress.Hex())

	tasks, ok := data["tasks"].([]interface{})
	require.True(t, ok)
	require.Len(t, tasks, 2)
	task0, ok := tasks[0].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, task0["index"], float64(0))
	assert.Contains(t, task0, "type")
	assert.Equal(t, task0["status"], "pending_confirmations")
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
	requestID := "RequestID"
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
				ID:         models.NewID(),
				JobSpecID:  models.NewID(),
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
		{"confirmed", "testdata/confirmedEthTxData.json", ""},
		{"safe fulfilled", "testdata/fulfilledReceiptResponse.json", "fulfilledRunLog"},
		{"safe not fulfilled", "testdata/notFulfilledReceiptResponse.json", "noFulfilledRunLog"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			newAddress := common.HexToAddress("0x9FBDa871d559710256a2502A2517b794B482Db40")
			requestID := "RequestID"
			requestTxHash := common.HexToHash("0xdeadbeef")
			dataJSON := jsonFromFixture(t, test.path)
			outgoingTxHash := "0x1111111111111111111111111111111111111111111111111111111111111111"

			taskSpec := models.TaskSpec{
				Type: "ethtx",
			}

			jobRun := models.JobRun{
				ID:        models.NewID(),
				JobSpecID: models.NewID(),
				Status:    models.RunStatusCompleted,
				Result:    models.RunResult{Amount: assets.NewLink(2)},
				Initiator: models.Initiator{
					Type: models.InitiatorRunLog,
				},
				RunRequest: models.RunRequest{
					RequestID: &requestID,
					TxHash:    &requestTxHash,
					Requester: &newAddress,
				},
				TaskRuns: []models.TaskRun{
					models.TaskRun{
						ID:       models.NewID(),
						TaskSpec: taskSpec,
						Status:   models.RunStatusPendingConfirmations,
						Result:   models.RunResult{Data: dataJSON},
					},
				},
			}
			p := SyncJobRunPresenter{JobRun: &jobRun}

			bytes, err := p.MarshalJSON()
			require.NoError(t, err)

			require.True(t, gjson.ValidBytes(bytes))
			data := gjson.ParseBytes(bytes)

			tasks := data.Get("tasks").Array()
			require.Len(t, tasks, 1)
			task0 := tasks[0].Map()
			assert.Equal(t, task0["index"].Float(), float64(0))
			assert.Contains(t, task0["type"].String(), "ethtx")
			assert.Equal(t, task0["status"].String(), "pending_confirmations")
			assert.Equal(t, task0["error"].Type, gjson.Null)

			txresult := task0["result"].Map()
			assert.Equal(t, test.want, txresult["transactionStatus"].String())
			assert.Equal(t, outgoingTxHash, txresult["transactionHash"].String())
		})
	}
}

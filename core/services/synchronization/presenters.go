package synchronization

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/core/store/assets"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
	null "gopkg.in/guregu/null.v3"
)

// SyncJobRunPresenter presents a JobRun for synchronization purposes
type SyncJobRunPresenter struct {
	*models.JobRun
}

// MarshalJSON returns the JobRun as JSON
func (p SyncJobRunPresenter) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ID          string                 `json:"id"`
		JobID       string                 `json:"jobId"`
		RunID       string                 `json:"runId"`
		Status      string                 `json:"status"`
		Error       null.String            `json:"error"`
		CreatedAt   string                 `json:"createdAt"`
		Amount      *assets.Link           `json:"amount"`
		CompletedAt null.Time              `json:"completedAt"`
		Initiator   syncInitiatorPresenter `json:"initiator"`
		Tasks       []syncTaskRunPresenter `json:"tasks"`
	}{
		ID:          p.ID,
		RunID:       p.ID,
		JobID:       p.JobSpecID,
		Status:      string(p.Status),
		Error:       p.Result.ErrorMessage,
		CreatedAt:   utils.ISO8601UTC(p.CreatedAt),
		Amount:      p.Result.Amount,
		CompletedAt: p.CompletedAt,
		Initiator:   p.initiator(),
		Tasks:       p.tasks(),
	})
}

func (p SyncJobRunPresenter) initiator() syncInitiatorPresenter {
	var eip *models.EIP55Address
	if p.RunRequest.Requester != nil {
		coerced := models.EIP55Address(p.RunRequest.Requester.Hex())
		eip = &coerced
	}
	return syncInitiatorPresenter{
		Type:      p.Initiator.Type,
		RequestID: p.RunRequest.RequestID,
		TxHash:    p.RunRequest.TxHash,
		Requester: eip,
	}
}

func (p SyncJobRunPresenter) tasks() []syncTaskRunPresenter {
	tasks := []syncTaskRunPresenter{}
	for index, tr := range p.TaskRuns {
		tasks = append(tasks, syncTaskRunPresenter{
			Index:  index,
			Type:   string(tr.TaskSpec.Type),
			Status: string(tr.Status),
			Error:  tr.Result.ErrorMessage,
			Result: fetchLastEthereumReceipt(tr),
		})
	}
	return tasks
}

func fetchLastEthereumReceipt(tr models.TaskRun) interface{} {
	if tr.TaskSpec.Type == "ethtx" {
		receipts := tr.Result.Data.Get("ethereumReceipts")
		if receipts.IsArray() {
			arr := receipts.Array()
			return arr[len(arr)-1].Value()
		}
	}
	return nil
}

type syncInitiatorPresenter struct {
	Type      string               `json:"type"`
	RequestID *string              `json:"requestId,omitempty"`
	TxHash    *common.Hash         `json:"txHash,omitempty"`
	Requester *models.EIP55Address `json:"requester,omitempty"`
}

type syncTaskRunPresenter struct {
	Index  int         `json:"index"`
	Type   string      `json:"type"`
	Status string      `json:"status"`
	Error  null.String `json:"error"`
	Result interface{} `json:"result,omitempty"`
}

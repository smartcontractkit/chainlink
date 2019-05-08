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
	tasks, err := p.tasks()
	if err != nil {
		return []byte{}, err
	}

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
		Tasks:       tasks,
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

func (p SyncJobRunPresenter) tasks() ([]syncTaskRunPresenter, error) {
	tasks := []syncTaskRunPresenter{}
	for index, tr := range p.TaskRuns {
		erp, err := fetchLastEthereumReceipt(tr)
		if err != nil {
			return []syncTaskRunPresenter{}, err
		}
		tasks = append(tasks, syncTaskRunPresenter{
			Index:  index,
			Type:   string(tr.TaskSpec.Type),
			Status: string(tr.Status),
			Error:  tr.Result.ErrorMessage,
			Result: erp,
		})
	}
	return tasks, nil
}

func fetchLastEthereumReceipt(tr models.TaskRun) (*syncReceiptPresenter, error) {
	if tr.TaskSpec.Type == "ethtx" {
		receipts := tr.Result.Data.Get("ethereumReceipts")
		if receipts.IsArray() {
			arr := receipts.Array()
			return formatEthereumReceipt(arr[len(arr)-1].String())
		}
	}
	return nil, nil
}

func formatEthereumReceipt(str string) (*syncReceiptPresenter, error) {
	var receipt models.TxReceipt
	err := json.Unmarshal([]byte(str), &receipt)
	if err != nil {
		return nil, err
	}
	return &syncReceiptPresenter{
		Hash:   receipt.Hash,
		Status: runLogStatusPresenter(receipt),
	}, nil
}

type syncReceiptPresenter struct {
	Hash   common.Hash `json:"transactionHash"`
	Status TxStatus    `json:"status"`
}

type TxStatus string

const (
	// StatusNoFulfilledRunlog indicates that a ChainlinkFulfilled event was
	// detected in the transaction receipt.
	StatusFulfilledRunLog TxStatus = "fulfilledRunLog"
	// StatusNoFulfilledRunlog indicates that no ChainlinkFulfilled events were
	// detected in the transaction receipt.
	StatusNoFulfilledRunLog TxStatus = "noFulfilledRunLog"
)

func runLogStatusPresenter(receipt models.TxReceipt) TxStatus {
	if receipt.FulfilledRunLog() {
		return StatusFulfilledRunLog
	}
	return StatusNoFulfilledRunLog
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

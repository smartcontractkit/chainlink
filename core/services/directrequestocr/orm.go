package directrequestocr

import (
	"fmt"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

//go:generate mockery --name ORM --output ./mocks/ --case=underscore

type ORM interface {
	CreateRequest(contractRequestID [32]byte, receivedAt time.Time, requestTxHash *common.Hash) (int64, error)

	SetResult(id int64, runID int64, computationResult []byte, readyAt time.Time) error
	SetError(id int64, runID int64, errorType ErrType, computationError string, readyAt time.Time) error
	SetConfirmed(contractRequestID [32]byte) error

	// TODO additional operations will be needed by the Reporting Plugin
	// https://app.shortcut.com/chainlinklabs/story/54054/ocr-plugin-for-directrequest-ocr

	// TODO include QOpts or context when moving to the DB ORM
	// https://app.shortcut.com/chainlinklabs/story/54049/database-table-in-core-node
}

type inmemoryorm struct {
	counter                int64
	db                     map[int64]Request
	contractRequestIDIndex map[[32]byte]int64
	mutex                  *sync.Mutex
}

var _ ORM = (*inmemoryorm)(nil)

func NewInMemoryORM() *inmemoryorm {
	return &inmemoryorm{
		counter:                0,
		db:                     make(map[int64]Request),
		contractRequestIDIndex: make(map[[32]byte]int64),
		mutex:                  &sync.Mutex{},
	}
}

func (o *inmemoryorm) CreateRequest(contractRequestID [32]byte, receivedAt time.Time, requestTxHash *common.Hash) (int64, error) {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	if dbID, ok := o.contractRequestIDIndex[contractRequestID]; ok {
		return dbID, fmt.Errorf("Request already exists! DBID: %v", dbID)
	}

	o.counter++
	newEntry := Request{
		ID:                o.counter,
		ContractRequestID: contractRequestID,
		ReceivedAt:        receivedAt,
		RequestTxHash:     requestTxHash,
		State:             IN_PROGRESS,
	}
	o.db[o.counter] = newEntry
	o.contractRequestIDIndex[contractRequestID] = o.counter
	return o.counter, nil
}

func (o *inmemoryorm) SetResult(id int64, runID int64, computationResult []byte, readyAt time.Time) error {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	if val, ok := o.db[id]; ok {
		val.RunID = runID
		val.ErrorType = NONE
		val.Result = computationResult
		val.ResultReadyAt = readyAt
		o.db[id] = val
		return nil
	}
	return fmt.Errorf("can't find entry with dbid: %v", id)
}

func (o *inmemoryorm) SetError(id int64, runID int64, errorType ErrType, computationError string, readyAt time.Time) error {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	if val, ok := o.db[id]; ok {
		val.RunID = runID
		val.ErrorType = errorType
		val.Error = computationError
		val.ResultReadyAt = readyAt
		o.db[id] = val
		return nil
	}
	return fmt.Errorf("can't find entry with dbid: %v", id)
}

func (o *inmemoryorm) SetConfirmed(contractRequestID [32]byte) error {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	if dbid, ok := o.contractRequestIDIndex[contractRequestID]; ok {
		if val, ok := o.db[dbid]; ok {
			val.State = CONFIRMED
			o.db[dbid] = val
			return nil
		}
		return fmt.Errorf("can't find entry with dbid: %v", dbid)
	}
	return fmt.Errorf("can't find entry with id: %v", contractRequestID)
}

// TODO actual DB: https://app.shortcut.com/chainlinklabs/story/54049/database-table-in-core-node

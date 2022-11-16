package directrequestocr

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/core/services/pg"
)

//go:generate mockery --quiet --name ORM --output ./mocks/ --case=underscore

type ORM interface {
	CreateRequest(requestID RequestID, receivedAt time.Time, requestTxHash *common.Hash, qopts ...pg.QOpt) error

	SetResult(requestID RequestID, runID int64, computationResult []byte, readyAt time.Time, qopts ...pg.QOpt) error
	SetError(requestID RequestID, runID int64, errorType ErrType, computationError []byte, readyAt time.Time, qopts ...pg.QOpt) error
	SetState(requestID RequestID, state RequestState, qopts ...pg.QOpt) (RequestState, error)

	FindOldestEntriesByState(state RequestState, limit uint32, qopts ...pg.QOpt) ([]Request, error)
	FindById(requestID RequestID, qopts ...pg.QOpt) (*Request, error)

	// TODO add jobID or contract address when moving to the DB ORM
	// TODO add state transition validation
	// https://app.shortcut.com/chainlinklabs/story/54049/database-table-in-core-node
}

type inmemoryorm struct {
	counter int64
	db      map[[32]byte]Request
	mutex   *sync.Mutex
}

var _ ORM = (*inmemoryorm)(nil)

func NewInMemoryORM() *inmemoryorm {
	return &inmemoryorm{
		counter: 0,
		db:      make(map[[32]byte]Request),
		mutex:   &sync.Mutex{},
	}
}

func (o *inmemoryorm) CreateRequest(requestID RequestID, receivedAt time.Time, requestTxHash *common.Hash, qopts ...pg.QOpt) error {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	if _, ok := o.db[requestID]; ok {
		return fmt.Errorf("request already exists")
	}

	o.counter++
	newEntry := Request{
		ID:            o.counter,
		RequestID:     requestID,
		ReceivedAt:    receivedAt,
		RequestTxHash: requestTxHash,
		State:         IN_PROGRESS,
	}
	o.db[requestID] = newEntry
	return nil
}

func (o *inmemoryorm) SetResult(requestID RequestID, runID int64, computationResult []byte, readyAt time.Time, qopts ...pg.QOpt) error {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	if val, ok := o.db[requestID]; ok {
		val.RunID = runID
		val.ErrorType = NONE
		val.Result = computationResult
		val.Error = []byte{}
		val.State = RESULT_READY
		val.ResultReadyAt = readyAt
		o.db[requestID] = val
		return nil
	}
	return fmt.Errorf("can't find entry with requestID: %v", requestID)
}

func (o *inmemoryorm) SetError(requestID RequestID, runID int64, errorType ErrType, computationError []byte, readyAt time.Time, qopts ...pg.QOpt) error {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	if val, ok := o.db[requestID]; ok {
		val.RunID = runID
		val.ErrorType = errorType
		val.Error = computationError
		val.State = RESULT_READY
		val.Result = []byte{}
		val.ResultReadyAt = readyAt
		o.db[requestID] = val
		return nil
	}
	return fmt.Errorf("can't find entry with requestID: %v", requestID)
}

func (o *inmemoryorm) SetState(requestID RequestID, state RequestState, qopts ...pg.QOpt) (RequestState, error) {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	prevState := IN_PROGRESS
	if val, ok := o.db[requestID]; ok {
		prevState = val.State
		val.State = state
		o.db[requestID] = val
		return prevState, nil
	}
	return prevState, fmt.Errorf("can't find entry with requestID: %v", requestID)
}

func (o *inmemoryorm) FindOldestEntriesByState(state RequestState, limit uint32, qopts ...pg.QOpt) ([]Request, error) {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	var result []Request
	// NOTE: suboptimal if limit << full result
	for _, val := range o.db {
		if val.State == state {
			result = append(result, val)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].ReceivedAt.Before(result[j].ReceivedAt)
	})
	if limit < uint32(len(result)) {
		result = result[:limit]
	}
	return result, nil
}

func (o *inmemoryorm) FindById(requestID RequestID, qopts ...pg.QOpt) (*Request, error) {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	if val, ok := o.db[requestID]; ok {
		return &val, nil
	}
	return nil, fmt.Errorf("can't find entry with dbid: %v", requestID)
}

// TODO actual DB: https://app.shortcut.com/chainlinklabs/story/54049/database-table-in-core-node

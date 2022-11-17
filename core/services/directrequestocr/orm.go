package directrequestocr

import (
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pg"
)

//go:generate mockery --quiet --name ORM --output ./mocks/ --case=underscore

type ORM interface {
	CreateRequest(requestID RequestID, contractAdderss *common.Address, receivedAt time.Time, requestTxHash *common.Hash, qopts ...pg.QOpt) error

	SetResult(requestID RequestID, runID int64, computationResult []byte, readyAt time.Time, qopts ...pg.QOpt) error
	SetError(requestID RequestID, runID int64, errorType ErrType, computationError []byte, readyAt time.Time, qopts ...pg.QOpt) error
	SetState(requestID RequestID, state RequestState, qopts ...pg.QOpt) (RequestState, error)

	FindOldestEntriesByState(state RequestState, limit uint32, qopts ...pg.QOpt) ([]Request, error)
	FindById(requestID RequestID, qopts ...pg.QOpt) (*Request, error)

	// TODO add state transition validation
	// https://app.shortcut.com/chainlinklabs/story/54049/database-table-in-core-node
}

type orm struct {
	q pg.Q
}

var _ ORM = (*orm)(nil)

func NewORM(db *sqlx.DB, lggr logger.Logger, cfg pg.QConfig) ORM {
	return &orm{
		q: pg.NewQ(db, lggr, cfg),
	}
}

func (o *orm) CreateRequest(requestID RequestID, contractAdderss *common.Address, receivedAt time.Time, requestTxHash *common.Hash, qopts ...pg.QOpt) error {
	stmt := `
		INSERT INTO ocr2dr_requests (request_id, contract_address, received_at, request_tx_hash, state)
		VALUES ($1,$2,$3,$4,'in_progress');
	`
	return o.q.WithOpts(qopts...).ExecQ(stmt, requestID, contractAdderss, receivedAt, requestTxHash)
}

func (o *orm) SetResult(requestID RequestID, runID int64, computationResult []byte, readyAt time.Time, qopts ...pg.QOpt) error {
	stmt := `
		UPDATE ocr2dr_requests
		SET (run_id=$2, result=$3, result_ready_at=$4, state='result_ready')
		WHERE request_id=$1;
	`
	return o.q.WithOpts(qopts...).ExecQ(stmt, requestID, runID, computationResult, readyAt)
}

func (o *orm) SetError(requestID RequestID, runID int64, errorType ErrType, computationError []byte, readyAt time.Time, qopts ...pg.QOpt) error {
	// TODO: add marshalling
	var dbErrorType string
	switch errorType {
	case NONE:
		dbErrorType = "none"
	case NODE_EXCEPTION:
		dbErrorType = "node_exception"
	case SANDBOX_TIMEOUT:
		dbErrorType = "sandbox_timeout"
	case USER_EXCEPTION:
		dbErrorType = "user_exception"
	}

	stmt := `
		UPDATE ocr2dr_requests
		SET (run_id=$2, error=$3, error_type=$4, result_ready_at=$5, state='result_ready')
		WHERE request_id=$1;
	`
	return o.q.WithOpts(qopts...).ExecQ(stmt, requestID, runID, computationError, dbErrorType, readyAt)
}

func (o *orm) SetState(requestID RequestID, state RequestState, qopts ...pg.QOpt) (RequestState, error) {
	// TODO: add marshalling
	var dbState string
	switch state {
	case IN_PROGRESS:
		dbState = "in_progress"
	case RESULT_READY:
		dbState = "result_ready"
	case TRANSMITTED:
		dbState = "transmitted"
	case CONFIRMED:
		dbState = "confirmed"
	}

	var oldState string

	stmt := `SELECT state FROM ocr2dr_requests WHERE request_id=$1;`
	if err := o.q.WithOpts(qopts...).Get(&oldState, stmt); err != nil {
		return 0, err
	}

	if oldState == dbState {
		return state, nil
	}

	// TODO: add marshalling
	var oldStateEnum RequestState
	switch oldState {
	case "in_progress":
		oldStateEnum = IN_PROGRESS
	case "result_ready":
		oldStateEnum = RESULT_READY
	case "transmitted":
		oldStateEnum = TRANSMITTED
	case "confirmed":
		oldStateEnum = CONFIRMED
	}

	stmt = `UPDATE ocr2dr_requests SET (state=$2) WHERE request_id=$1;`
	return oldStateEnum, o.q.WithOpts(qopts...).ExecQ(stmt, requestID, dbState)
}

func (o *orm) FindOldestEntriesByState(state RequestState, limit uint32, qopts ...pg.QOpt) ([]Request, error) {
	/*
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
	*/
	return nil, nil
}

func (o *orm) FindById(requestID RequestID, qopts ...pg.QOpt) (*Request, error) {
	request := &Request{}
	stmt := `SELECT * FROM ocr2dr_requests WHERE request_id=$1;`
	err := o.q.WithOpts(qopts...).Get(&request, stmt, requestID)
	return request, err
}

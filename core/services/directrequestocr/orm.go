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
	CreateRequest(requestID RequestID, receivedAt time.Time, requestTxHash *common.Hash, qopts ...pg.QOpt) error

	SetResult(requestID RequestID, runID int64, computationResult []byte, readyAt time.Time, qopts ...pg.QOpt) error
	SetError(requestID RequestID, runID int64, errorType ErrType, computationError []byte, readyAt time.Time, qopts ...pg.QOpt) error
	SetState(requestID RequestID, state RequestState, qopts ...pg.QOpt) (RequestState, error)
	SetTransmittedResult(requestID RequestID, transmittedResult []byte, transmittedError []byte, qopts ...pg.QOpt) error

	FindOldestEntriesByState(state RequestState, limit uint32, qopts ...pg.QOpt) ([]Request, error)
	FindById(requestID RequestID, qopts ...pg.QOpt) (*Request, error)

	// TODO add state transition validation
	// https://app.shortcut.com/chainlinklabs/story/54049/database-table-in-core-node
}

type orm struct {
	q               pg.Q
	contractAdderss common.Address
}

var _ ORM = (*orm)(nil)

func NewORM(db *sqlx.DB, lggr logger.Logger, cfg pg.QConfig, contractAdderss common.Address) ORM {
	return &orm{
		q:               pg.NewQ(db, lggr, cfg),
		contractAdderss: contractAdderss,
	}
}

func (o *orm) CreateRequest(requestID RequestID, receivedAt time.Time, requestTxHash *common.Hash, qopts ...pg.QOpt) error {
	stmt := `
		INSERT INTO ocr2dr_requests (request_id, contract_address, received_at, request_tx_hash, state)
		VALUES ($1,$2,$3,$4,'in_progress');
	`
	return o.q.WithOpts(qopts...).ExecQ(stmt, requestID, o.contractAdderss, receivedAt, requestTxHash)
}

func (o *orm) SetResult(requestID RequestID, runID int64, computationResult []byte, readyAt time.Time, qopts ...pg.QOpt) error {
	stmt := `
		UPDATE ocr2dr_requests
		SET (run_id=$2, result=$3, result_ready_at=$4, state=$5)
		WHERE request_id=$1;
	`
	return o.q.WithOpts(qopts...).ExecQ(stmt, requestID, runID, computationResult, readyAt, RESULT_READY)
}

func (o *orm) SetError(requestID RequestID, runID int64, errorType ErrType, computationError []byte, readyAt time.Time, qopts ...pg.QOpt) error {
	stmt := `
		UPDATE ocr2dr_requests
		SET (run_id=$2, error=$3, error_type=$4, result_ready_at=$5, state=$6)
		WHERE request_id=$1;
	`
	return o.q.WithOpts(qopts...).ExecQ(stmt, requestID, runID, computationError, errorType, readyAt, RESULT_READY)
}

func (o *orm) SetState(requestID RequestID, state RequestState, qopts ...pg.QOpt) (RequestState, error) {
	var oldState RequestState

	stmt := `SELECT state FROM ocr2dr_requests WHERE request_id=$1;`
	if err := o.q.WithOpts(qopts...).Get(&oldState, stmt); err != nil {
		return 0, err
	}

	if oldState == state {
		return state, nil
	}

	stmt = `UPDATE ocr2dr_requests SET (state=$2) WHERE request_id=$1;`
	return oldState, o.q.WithOpts(qopts...).ExecQ(stmt, requestID, state)
}

func (o *orm) SetTransmittedResult(requestID RequestID, transmittedResult []byte, transmittedError []byte, qopts ...pg.QOpt) error {
	stmt := `
		UPDATE ocr2dr_requests
		SET (transmitted_result=$2, transmitted_error=$3)
		WHERE request_id=$1;
	`
	return o.q.WithOpts(qopts...).ExecQ(stmt, requestID, transmittedResult, transmittedError)
}

func (o *orm) FindOldestEntriesByState(state RequestState, limit uint32, qopts ...pg.QOpt) ([]Request, error) {
	var requests []Request
	stmt := `SELECT * FROM ocr2dr_requests WHERE state=$1 ORDER_BY received_at DESC LIMIT $2;`
	err := o.q.WithOpts(qopts...).Get(&requests, stmt, state, limit)
	return requests, err
}

func (o *orm) FindById(requestID RequestID, qopts ...pg.QOpt) (*Request, error) {
	request := &Request{}
	stmt := `SELECT * FROM ocr2dr_requests WHERE request_id=$1;`
	err := o.q.WithOpts(qopts...).Get(&request, stmt, requestID)
	return request, err
}

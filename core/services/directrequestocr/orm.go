package directrequestocr

import (
	"fmt"
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
	SetTransmitted(requestID RequestID, transmittedResult []byte, transmittedError []byte, qopts ...pg.QOpt) error

	FindOldestEntriesByState(state RequestState, limit uint32, qopts ...pg.QOpt) ([]Request, error)
	FindById(requestID RequestID, qopts ...pg.QOpt) (*Request, error)

	// TODO add state transition validation
	// https://app.shortcut.com/chainlinklabs/story/54049/database-table-in-core-node
}

type orm struct {
	q               pg.Q
	contractAddress common.Address
}

var _ ORM = (*orm)(nil)

const requestFields = "request_id, run_id, received_at, request_tx_hash, " +
	"state, result_ready_at, result, error_type, error, " +
	"transmitted_result, transmitted_error"

func NewORM(db *sqlx.DB, lggr logger.Logger, cfg pg.QConfig, contractAddress common.Address) ORM {
	return &orm{
		q:               pg.NewQ(db, lggr, cfg),
		contractAddress: contractAddress,
	}
}

func (o orm) CreateRequest(requestID RequestID, receivedAt time.Time, requestTxHash *common.Hash, qopts ...pg.QOpt) error {
	stmt := `
		INSERT INTO ocr2dr_requests (request_id, contract_address, received_at, request_tx_hash, state)
		VALUES ($1,$2,$3,$4,$5);
	`
	return o.q.WithOpts(qopts...).ExecQ(stmt, requestID[:], o.contractAddress.Bytes(), receivedAt, requestTxHash.Bytes(), IN_PROGRESS)
}

func (o orm) SetResult(requestID RequestID, runID int64, computationResult []byte, readyAt time.Time, qopts ...pg.QOpt) error {
	stmt := `
		UPDATE ocr2dr_requests
		SET run_id=$3, result=$4, result_ready_at=$5, state=$6
		WHERE request_id=$1 AND contract_address=$2;
	`
	return o.q.WithOpts(qopts...).ExecQ(stmt, requestID[:], o.contractAddress.Bytes(), runID, computationResult, readyAt, RESULT_READY)
}

func (o orm) SetError(requestID RequestID, runID int64, errorType ErrType, computationError []byte, readyAt time.Time, qopts ...pg.QOpt) error {
	stmt := `
		UPDATE ocr2dr_requests
		SET run_id=$3, error=$4, error_type=$5, result_ready_at=$6, state=$7
		WHERE request_id=$1 AND contract_address=$2;
	`
	return o.q.WithOpts(qopts...).ExecQ(stmt, requestID[:], o.contractAddress, runID, computationError, errorType, readyAt, RESULT_READY)
}

func (o orm) SetState(requestID RequestID, state RequestState, qopts ...pg.QOpt) (RequestState, error) {
	var oldState RequestState

	stmt := `SELECT state FROM ocr2dr_requests WHERE request_id=$1 AND contract_address=$2;`
	if err := o.q.WithOpts(qopts...).Get(&oldState, stmt, requestID[:], o.contractAddress); err != nil {
		return 0, err
	}

	if oldState == state {
		return state, nil
	}

	stmt = `UPDATE ocr2dr_requests SET state=$3 WHERE request_id=$1 AND contract_address=$2;`
	return oldState, o.q.WithOpts(qopts...).ExecQ(stmt, requestID[:], o.contractAddress.Bytes(), state)
}

func (o orm) SetTransmitted(requestID RequestID, transmittedResult []byte, transmittedError []byte, qopts ...pg.QOpt) error {
	stmt := `
		UPDATE ocr2dr_requests
		SET transmitted_result=$3, transmitted_error=$4
		WHERE request_id=$1 AND contract_address=$2;
	`
	return o.q.WithOpts(qopts...).ExecQ(stmt, requestID[:], o.contractAddress.Bytes(), transmittedResult, transmittedError)
}

func (o orm) FindOldestEntriesByState(state RequestState, limit uint32, qopts ...pg.QOpt) ([]Request, error) {
	var requests []Request
	stmt := fmt.Sprintf(`SELECT %s FROM ocr2dr_requests WHERE state=$1 ORDER BY received_at LIMIT $2;`, requestFields)
	if err := o.q.WithOpts(qopts...).Select(&requests, stmt, state, limit); err != nil {
		return nil, err
	}
	return requests, nil
}

func (o orm) FindById(requestID RequestID, qopts ...pg.QOpt) (*Request, error) {
	var request Request
	stmt := fmt.Sprintf(`SELECT %s FROM ocr2dr_requests WHERE request_id=$1 AND contract_address=$2;`, requestFields)
	if err := o.q.WithOpts(qopts...).Get(&request, stmt, requestID[:], o.contractAddress.Bytes()); err != nil {
		return nil, err
	}
	return &request, nil
}

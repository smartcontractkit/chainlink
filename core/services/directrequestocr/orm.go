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
	SetTransmitted(requestID RequestID, transmittedResult []byte, transmittedError []byte, qopts ...pg.QOpt) error
	SetConfirmed(requestID RequestID, qopts ...pg.QOpt) error

	TimeoutExpiredResults(cutoff time.Time, limit uint32, qopts ...pg.QOpt) ([]RequestID, error)

	FindOldestEntriesByState(state RequestState, limit uint32, qopts ...pg.QOpt) ([]Request, error)
	FindById(requestID RequestID, qopts ...pg.QOpt) (*Request, error)
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
	return o.q.WithOpts(qopts...).ExecQ(stmt, requestID, o.contractAddress, receivedAt, requestTxHash, IN_PROGRESS)
}

func (o orm) setWithStateTransitionCheck(requestID RequestID, newState RequestState, setter func(pg.Queryer) error, qopts ...pg.QOpt) error {
	err := o.q.WithOpts(qopts...).Transaction(func(tx pg.Queryer) error {
		prevState := IN_PROGRESS // default initial state
		stmt := `SELECT state FROM ocr2dr_requests WHERE request_id=$1 AND contract_address=$2;`
		if err2 := tx.Get(&prevState, stmt, requestID, o.contractAddress); err2 != nil {
			return err2
		}
		if err2 := CheckStateTransition(prevState, newState); err2 != nil {
			return err2
		}
		return setter(tx)
	})

	return err
}

func (o orm) SetResult(requestID RequestID, runID int64, computationResult []byte, readyAt time.Time, qopts ...pg.QOpt) error {
	newState := RESULT_READY
	err := o.setWithStateTransitionCheck(requestID, newState, func(tx pg.Queryer) error {
		stmt := `
			UPDATE ocr2dr_requests
			SET run_id=$3, result=$4, result_ready_at=$5, state=$6
			WHERE request_id=$1 AND contract_address=$2;
		`
		_, err2 := tx.Exec(stmt, requestID, o.contractAddress, runID, computationResult, readyAt, newState)
		return err2
	}, qopts...)
	return err
}

func (o orm) SetError(requestID RequestID, runID int64, errorType ErrType, computationError []byte, readyAt time.Time, qopts ...pg.QOpt) error {
	newState := RESULT_READY
	err := o.setWithStateTransitionCheck(requestID, newState, func(tx pg.Queryer) error {
		stmt := `
			UPDATE ocr2dr_requests
			SET run_id=$3, error=$4, error_type=$5, result_ready_at=$6, state=$7
			WHERE request_id=$1 AND contract_address=$2;
		`
		_, err2 := tx.Exec(stmt, requestID, o.contractAddress, runID, computationError, errorType, readyAt, newState)
		return err2
	}, qopts...)
	return err
}

func (o orm) SetTransmitted(requestID RequestID, transmittedResult []byte, transmittedError []byte, qopts ...pg.QOpt) error {
	newState := TRANSMITTED
	err := o.setWithStateTransitionCheck(requestID, newState, func(tx pg.Queryer) error {
		stmt := `
			UPDATE ocr2dr_requests
			SET transmitted_result=$3, transmitted_error=$4, state=$5
			WHERE request_id=$1 AND contract_address=$2;
		`
		_, err2 := tx.Exec(stmt, requestID, o.contractAddress, transmittedResult, transmittedError, newState)
		return err2
	}, qopts...)
	return err
}

func (o orm) SetConfirmed(requestID RequestID, qopts ...pg.QOpt) error {
	newState := CONFIRMED
	err := o.setWithStateTransitionCheck(requestID, newState, func(tx pg.Queryer) error {
		stmt := `UPDATE ocr2dr_requests SET state=$3 WHERE request_id=$1 AND contract_address=$2;`
		_, err2 := tx.Exec(stmt, requestID, o.contractAddress, newState)
		return err2
	}, qopts...)
	return err
}

func (o orm) TimeoutExpiredResults(cutoff time.Time, limit uint32, qopts ...pg.QOpt) ([]RequestID, error) {
	var ids []RequestID
	prevStates := []RequestState{IN_PROGRESS, RESULT_READY}
	nextState := TIMED_OUT
	if err := CheckStateTransition(prevStates[0], nextState); err != nil {
		return ids, err
	}
	if err := CheckStateTransition(prevStates[1], nextState); err != nil {
		return ids, err
	}
	err := o.q.WithOpts(qopts...).Transaction(func(tx pg.Queryer) error {
		selectStmt := `
			SELECT request_id
			FROM ocr2dr_requests
			WHERE (state=$1 OR state=$2) AND contract_address=$3 AND received_at < ($4)
			ORDER BY received_at
			LIMIT $5;`
		if err2 := tx.Select(&ids, selectStmt, prevStates[0], prevStates[1], o.contractAddress, cutoff, limit); err2 != nil {
			return err2
		}
		if len(ids) == 0 {
			return nil
		}

		a := map[string]any{
			"nextState":    nextState,
			"contractAddr": o.contractAddress,
			"ids":          ids,
		}
		updateStmt, args, err2 := sqlx.Named(`
			UPDATE ocr2dr_requests
			SET state = :nextState
			WHERE contract_address = :contractAddr AND request_id IN (:ids);`, a)
		if err2 != nil {
			return err2
		}
		updateStmt, args, err2 = sqlx.In(updateStmt, args...)
		if err2 != nil {
			return err2
		}
		updateStmt = tx.Rebind(updateStmt)
		if _, err2 := tx.Exec(updateStmt, args...); err2 != nil {
			return err2
		}
		return nil
	})

	return ids, err
}

func (o orm) FindOldestEntriesByState(state RequestState, limit uint32, qopts ...pg.QOpt) ([]Request, error) {
	var requests []Request
	stmt := fmt.Sprintf(`SELECT %s FROM ocr2dr_requests WHERE state=$1 AND contract_address=$2 ORDER BY received_at LIMIT $3;`, requestFields)
	if err := o.q.WithOpts(qopts...).Select(&requests, stmt, state, o.contractAddress, limit); err != nil {
		return nil, err
	}
	return requests, nil
}

func (o orm) FindById(requestID RequestID, qopts ...pg.QOpt) (*Request, error) {
	var request Request
	stmt := fmt.Sprintf(`SELECT %s FROM ocr2dr_requests WHERE request_id=$1 AND contract_address=$2;`, requestFields)
	if err := o.q.WithOpts(qopts...).Get(&request, stmt, requestID, o.contractAddress); err != nil {
		return nil, err
	}
	return &request, nil
}

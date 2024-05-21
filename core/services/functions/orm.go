package functions

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
)

//go:generate mockery --quiet --name ORM --output ./mocks/ --case=underscore

type ORM interface {
	CreateRequest(ctx context.Context, request *Request) error

	SetResult(ctx context.Context, requestID RequestID, computationResult []byte, readyAt time.Time) error
	SetError(ctx context.Context, requestID RequestID, errorType ErrType, computationError []byte, readyAt time.Time, readyForProcessing bool) error
	SetFinalized(ctx context.Context, requestID RequestID, reportedResult []byte, reportedError []byte) error
	SetConfirmed(ctx context.Context, requestID RequestID) error

	TimeoutExpiredResults(ctx context.Context, cutoff time.Time, limit uint32) ([]RequestID, error)

	FindOldestEntriesByState(ctx context.Context, state RequestState, limit uint32) ([]Request, error)
	FindById(ctx context.Context, requestID RequestID) (*Request, error)

	PruneOldestRequests(ctx context.Context, maxRequestsInDB uint32, batchSize uint32) (total uint32, pruned uint32, err error)
}

type orm struct {
	ds              sqlutil.DataSource
	contractAddress common.Address
}

var _ ORM = (*orm)(nil)

var ErrDuplicateRequestID = errors.New("Functions ORM: duplicate request ID")

const (
	tableName           = "functions_requests"
	defaultInitialState = IN_PROGRESS
	requestFields       = "request_id, received_at, request_tx_hash, " +
		"state, result_ready_at, result, error_type, error, " +
		"transmitted_result, transmitted_error, flags, aggregation_method, " +
		"callback_gas_limit, coordinator_contract_address, onchain_metadata, processing_metadata"
)

func NewORM(ds sqlutil.DataSource, contractAddress common.Address) ORM {
	return &orm{
		ds:              ds,
		contractAddress: contractAddress,
	}
}

func (o *orm) CreateRequest(ctx context.Context, request *Request) error {
	stmt := fmt.Sprintf(`
		INSERT INTO %s (request_id, contract_address, received_at, request_tx_hash, state, flags, aggregation_method, callback_gas_limit, coordinator_contract_address, onchain_metadata)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) ON CONFLICT (request_id) DO NOTHING;
	`, tableName)
	result, err := o.ds.ExecContext(
		ctx,
		stmt,
		request.RequestID,
		o.contractAddress,
		request.ReceivedAt,
		request.RequestTxHash,
		defaultInitialState,
		request.Flags,
		request.AggregationMethod,
		request.CallbackGasLimit,
		request.CoordinatorContractAddress,
		request.OnchainMetadata)
	if err != nil {
		return err
	}
	nrows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if nrows == 0 {
		return ErrDuplicateRequestID
	}
	return nil
}

func (o *orm) setWithStateTransitionCheck(ctx context.Context, requestID RequestID, newState RequestState, setter func(sqlutil.DataSource) error) error {
	err := sqlutil.TransactDataSource(ctx, o.ds, nil, func(tx sqlutil.DataSource) error {
		prevState := defaultInitialState
		stmt := fmt.Sprintf(`SELECT state FROM %s WHERE request_id=$1 AND contract_address=$2;`, tableName)
		if err2 := tx.GetContext(ctx, &prevState, stmt, requestID, o.contractAddress); err2 != nil {
			return err2
		}
		if err2 := CheckStateTransition(prevState, newState); err2 != nil {
			return err2
		}
		return setter(tx)
	})

	return err
}

func (o *orm) SetResult(ctx context.Context, requestID RequestID, computationResult []byte, readyAt time.Time) error {
	newState := RESULT_READY
	err := o.setWithStateTransitionCheck(ctx, requestID, newState, func(tx sqlutil.DataSource) error {
		stmt := fmt.Sprintf(`
			UPDATE %s
			SET result=$3, result_ready_at=$4, state=$5
			WHERE request_id=$1 AND contract_address=$2;
		`, tableName)
		_, err2 := tx.ExecContext(ctx, stmt, requestID, o.contractAddress, computationResult, readyAt, newState)
		return err2
	})
	return err
}

func (o *orm) SetError(ctx context.Context, requestID RequestID, errorType ErrType, computationError []byte, readyAt time.Time, readyForProcessing bool) error {
	var newState RequestState
	if readyForProcessing {
		newState = RESULT_READY
	} else {
		newState = IN_PROGRESS
	}
	err := o.setWithStateTransitionCheck(ctx, requestID, newState, func(tx sqlutil.DataSource) error {
		stmt := fmt.Sprintf(`
			UPDATE %s
			SET error=$3, error_type=$4, result_ready_at=$5, state=$6
			WHERE request_id=$1 AND contract_address=$2;
		`, tableName)
		_, err2 := tx.ExecContext(ctx, stmt, requestID, o.contractAddress, computationError, errorType, readyAt, newState)
		return err2
	})
	return err
}

func (o *orm) SetFinalized(ctx context.Context, requestID RequestID, reportedResult []byte, reportedError []byte) error {
	newState := FINALIZED
	err := o.setWithStateTransitionCheck(ctx, requestID, newState, func(tx sqlutil.DataSource) error {
		stmt := fmt.Sprintf(`
			UPDATE %s
			SET transmitted_result=$3, transmitted_error=$4, state=$5
			WHERE request_id=$1 AND contract_address=$2;
		`, tableName)
		_, err2 := tx.ExecContext(ctx, stmt, requestID, o.contractAddress, reportedResult, reportedError, newState)
		return err2
	})
	return err
}

func (o *orm) SetConfirmed(ctx context.Context, requestID RequestID) error {
	newState := CONFIRMED
	err := o.setWithStateTransitionCheck(ctx, requestID, newState, func(tx sqlutil.DataSource) error {
		stmt := fmt.Sprintf(`UPDATE %s SET state=$3 WHERE request_id=$1 AND contract_address=$2;`, tableName)
		_, err2 := tx.ExecContext(ctx, stmt, requestID, o.contractAddress, newState)
		return err2
	})
	return err
}

func (o *orm) TimeoutExpiredResults(ctx context.Context, cutoff time.Time, limit uint32) ([]RequestID, error) {
	var ids []RequestID
	allowedPrevStates := []RequestState{IN_PROGRESS, RESULT_READY, FINALIZED}
	nextState := TIMED_OUT
	for _, state := range allowedPrevStates {
		// sanity checks
		if err := CheckStateTransition(state, nextState); err != nil {
			return ids, err
		}
	}
	err := sqlutil.TransactDataSource(ctx, o.ds, nil, func(tx sqlutil.DataSource) error {
		selectStmt := fmt.Sprintf(`
			SELECT request_id
			FROM %s
			WHERE (state=$1 OR state=$2 OR state=$3) AND contract_address=$4 AND received_at < ($5)
			ORDER BY received_at
			LIMIT $6;`, tableName)
		if err2 := tx.SelectContext(ctx, &ids, selectStmt, allowedPrevStates[0], allowedPrevStates[1], allowedPrevStates[2], o.contractAddress, cutoff, limit); err2 != nil {
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
		updateStmt, args, err2 := sqlx.Named(fmt.Sprintf(`
			UPDATE %s
			SET state = :nextState
			WHERE contract_address = :contractAddr AND request_id IN (:ids);`, tableName), a)
		if err2 != nil {
			return err2
		}
		updateStmt, args, err2 = sqlx.In(updateStmt, args...)
		if err2 != nil {
			return err2
		}
		updateStmt = tx.Rebind(updateStmt)
		if _, err2 := tx.ExecContext(ctx, updateStmt, args...); err2 != nil {
			return err2
		}
		return nil
	})

	return ids, err
}

func (o *orm) FindOldestEntriesByState(ctx context.Context, state RequestState, limit uint32) ([]Request, error) {
	var requests []Request
	stmt := fmt.Sprintf(`SELECT %s FROM %s WHERE state=$1 AND contract_address=$2 ORDER BY received_at LIMIT $3;`, requestFields, tableName)
	if err := o.ds.SelectContext(ctx, &requests, stmt, state, o.contractAddress, limit); err != nil {
		return nil, err
	}
	return requests, nil
}

func (o *orm) FindById(ctx context.Context, requestID RequestID) (*Request, error) {
	var request Request
	stmt := fmt.Sprintf(`SELECT %s FROM %s WHERE request_id=$1 AND contract_address=$2;`, requestFields, tableName)
	if err := o.ds.GetContext(ctx, &request, stmt, requestID, o.contractAddress); err != nil {
		return nil, err
	}
	return &request, nil
}

func (o *orm) PruneOldestRequests(ctx context.Context, maxStoredRequests uint32, batchSize uint32) (total uint32, pruned uint32, err error) {
	err = sqlutil.TransactDataSource(ctx, o.ds, nil, func(tx sqlutil.DataSource) error {
		stmt := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE contract_address=$1`, tableName)
		if err2 := tx.GetContext(ctx, &total, stmt, o.contractAddress); err2 != nil {
			return errors.Wrap(err, "failed to get request count")
		}

		if total <= maxStoredRequests {
			pruned = 0
			return nil
		}

		pruneLimit := total - maxStoredRequests
		if pruneLimit > batchSize {
			pruneLimit = batchSize
		}

		with := fmt.Sprintf(`WITH ids AS (SELECT request_id FROM %s WHERE contract_address = $1 ORDER BY received_at LIMIT $2)`, tableName)
		deleteStmt := fmt.Sprintf(`%s DELETE FROM %s WHERE contract_address = $1 AND request_id IN (SELECT request_id FROM ids);`, with, tableName)
		res, err2 := tx.ExecContext(ctx, deleteStmt, o.contractAddress, pruneLimit)
		if err2 != nil {
			return err2
		}
		prunedInt64, err2 := res.RowsAffected()
		if err2 == nil {
			pruned = uint32(prunedInt64)
		}
		return err2
	})
	return
}

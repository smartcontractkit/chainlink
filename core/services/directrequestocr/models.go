package directrequestocr

import (
	"database/sql/driver"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

const RequestIDLength int = 32

type RequestID [RequestIDLength]byte
type Request struct {
	RequestID         RequestID
	RunID             *int64
	ReceivedAt        time.Time
	RequestTxHash     *common.Hash
	State             RequestState
	ResultReadyAt     *time.Time
	Result            []byte
	ErrorType         *ErrType
	Error             []byte
	TransmittedResult []byte
	TransmittedError  []byte
	// TODO: add timestamps for other possible states: https://app.shortcut.com/chainlinklabs/story/58428/timestamp-fields-for-all-states-in-ocr2dr-data-model
}

type RequestState int8

const (
	// IN_PROGRESS is the initial state of a request, set right after receiving it in an on-chain event.
	IN_PROGRESS RequestState = iota

	// RESULT_READY means that computation has finished executing (with either success or error).
	// OCR2 reporting includes only requests in RESULT_READY state (for Query and Observation phases).
	RESULT_READY

	// TIMED_OUT request has not been transmitted yet and has been waiting too long for OCR2 processing.
	// It won't be included in OCR2 reporting rounds any more.
	TIMED_OUT

	// TRANSMITTED request is one that passed through ShouldTransmitAcceptedReport phase of OCR2 reporting.
	// From here, it can only transition to CONFIRMED state, upon receiving an on-chain confirmation from the oracle contract.
	TRANSMITTED

	// CONFIRMED state indicates that we received an on-chain confirmation event
	// (with or without this node's participation in an earlier OCR round).
	// We can transition here at any time (full fan-in) and cannot transition out (empty fan-out).
	// This is a desired and expected final state for every request.
	CONFIRMED
)

/*
 *    +-----------+
 * +--+IN_PROGRESS+----------------+
 * |  +-----+-----+                |
 * |        |                      |
 * |        v                      v
 * |  +------------+          +---------+
 * |  |RESULT_READY+--------->|TIMED_OUT|
 * |  +-----+------+          +---------+
 * |        |
 * |        v
 * |  +-----------+
 * +->|TRANSMITTED|
 *    +-----------+
 *
 *                     \   /
 *                       |
 *                       v
 *                  +---------+
 *                  |CONFIRMED|
 *                  +---------+
 */
func CheckStateTransition(prev RequestState, next RequestState) error {
	if prev == next {
		return errors.New("attempt to set the same state")
	}
	if prev == CONFIRMED {
		return errors.New("cannot transition out of CONFIRMED state")
	}
	transitions := map[RequestState]map[RequestState]error{
		IN_PROGRESS: {
			RESULT_READY: nil, // computation completed (either successfully or not)
			TIMED_OUT:    nil, // timing out a request in progress - what happened to the computation?
			TRANSMITTED:  nil, // transmitted a report without this node's participation in OCR round
			CONFIRMED:    nil, // received an on-chain result confirmation
		},
		RESULT_READY: {
			IN_PROGRESS: errors.New("cannot go back from RESULT_READY to IN_PROGRESS"),
			TIMED_OUT:   nil, // timing out a request - why was it never picked up by OCR reporting?
			TRANSMITTED: nil, // transmitted via OCR as expected
			CONFIRMED:   nil, // received an on-chain result confirmation
		},
		TIMED_OUT: {
			IN_PROGRESS:  errors.New("cannot go back from TIMED_OUT to IN_PROGRESS"),
			RESULT_READY: errors.New("cannot go back from TIMED_OUT to RESULT_READY"),
			TRANSMITTED:  errors.New("result already timed out but we're trying to transmit it (maybe a harmless race with the timer?)"),
			CONFIRMED:    nil, // received an on-chain result confirmation
		},
		TRANSMITTED: {
			IN_PROGRESS:  errors.New("cannot go back from TRANSMITTED to IN_PROGRESS"),
			RESULT_READY: errors.New("cannot go back from TRANSMITTED to RESULT_READY"),
			TIMED_OUT:    errors.New("result already transmitted, no need to time it out (maybe a harmless race with the timer?)"),
			CONFIRMED:    nil, // received an on-chain result confirmation
		},
	}

	nextMap, exists := transitions[prev]
	if !exists {
		return fmt.Errorf("unaccounted for state transition attempt, this should never happen (prev: %v, next: %v)", prev, next)
	}
	retErr, exists := nextMap[next]
	if !exists {
		return fmt.Errorf("unaccounted for state transition attempt, this should never happen (prev: %v, next: %v)", prev, next)
	}
	return retErr
}

type ErrType int8

const (
	NONE ErrType = iota
	NODE_EXCEPTION
	SANDBOX_TIMEOUT
	USER_EXCEPTION
)

func (r *RequestID) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("can't scan %T into RequestID", value)
	}
	if len(bytes) != RequestIDLength {
		return fmt.Errorf("can't scan []byte of len %d into RequestID, want %d", len(bytes), RequestIDLength)
	}
	copy(r[:], bytes)
	return nil
}

func (r RequestID) Value() (driver.Value, error) {
	return r[:], nil
}

func (s RequestState) String() string {
	switch s {
	case IN_PROGRESS:
		return "InProgress"
	case RESULT_READY:
		return "ResultReady"
	case TIMED_OUT:
		return "TimedOut"
	case TRANSMITTED:
		return "Transmitted"
	case CONFIRMED:
		return "Confirmed"
	}
	return "unknown"
}

func (e ErrType) String() string {
	switch e {
	case NONE:
		return "None"
	case NODE_EXCEPTION:
		return "NodeException"
	case SANDBOX_TIMEOUT:
		return "SandboxTimeout"
	case USER_EXCEPTION:
		return "UserException"
	}
	return "unknown"
}

func (r RequestID) String() string {
	return hex.EncodeToString(r[:])
}

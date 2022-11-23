package directrequestocr

import (
	"database/sql/driver"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type RequestState int8

const (
	IN_PROGRESS RequestState = iota
	RESULT_READY
	TRANSMITTED
	CONFIRMED
)

type ErrType int8

const (
	NONE ErrType = iota
	NODE_EXCEPTION
	SANDBOX_TIMEOUT
	USER_EXCEPTION
)

type RequestID [32]byte

const RequestIDLength int = 32

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
}

func (s RequestState) String() string {
	switch s {
	case IN_PROGRESS:
		return "InProgress"
	case RESULT_READY:
		return "ResultReady"
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

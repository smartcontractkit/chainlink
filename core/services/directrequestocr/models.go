package directrequestocr

import (
	"encoding/hex"
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

type Request struct {
	ID            int64
	RequestID     RequestID
	RunID         int64
	ReceivedAt    time.Time
	RequestTxHash *common.Hash
	State         RequestState
	ResultReadyAt time.Time
	Result        []byte
	ErrorType     ErrType
	Error         []byte
	// True if this node submitted an observation for this request in any OCR rounds.
	IsOCRParticipant  bool
	TransmittedResult []byte
	TransmittedError  []byte
	OnChainResult     []byte
	OnChainError      []byte
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

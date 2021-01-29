package log

import (
	"reflect"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/core/services/eth"
)

// DecodingListener receives raw logs from the Broadcaster and decodes them into
// Go structs using the provided ContractCodec (a simple wrapper around a go-ethereum
// ABI type).
type decodingListener struct {
	logTypes map[common.Hash]reflect.Type
	codec    eth.ContractCodec
	Listener
}

var _ Listener = (*decodingListener)(nil)

// NewDecodingListener creates a new decodingListener
func NewDecodingListener(codec eth.ContractCodec, nativeLogTypes map[common.Hash]interface{}, innerListener Listener) Listener {
	logTypes := make(map[common.Hash]reflect.Type)
	for eventID, logStruct := range nativeLogTypes {
		logTypes[eventID] = reflect.TypeOf(logStruct)
	}

	return &decodingListener{
		logTypes: logTypes,
		codec:    codec,
		Listener: innerListener,
	}
}

func (l *decodingListener) HandleLog(lb Broadcast, err error) {
	if err != nil {
		l.Listener.HandleLog(&broadcast{}, err)
		return
	}

	rawLog := lb.RawLog()

	if len(rawLog.Topics) == 0 {
		return
	}
	eventID := rawLog.Topics[0]
	logType, exists := l.logTypes[eventID]
	if !exists {
		// If a particular log type hasn't been registered with the decoder, we simply ignore it.
		return
	}

	var decodedLog interface{}
	if logType.Kind() == reflect.Ptr {
		decodedLog = reflect.New(logType.Elem()).Interface()
	} else {
		decodedLog = reflect.New(logType).Interface()
	}

	// Insert the raw log into the ".Log" field
	logStructV := reflect.ValueOf(decodedLog).Elem()
	logStructV.FieldByName("Log").Set(reflect.ValueOf(rawLog))

	// Decode the raw log into the struct
	event, err := l.codec.ABI().EventByID(eventID)
	if err != nil {
		l.Listener.HandleLog(nil, err)
		return
	}
	err = l.codec.UnpackLog(decodedLog, event.RawName, rawLog)
	if err != nil {
		l.Listener.HandleLog(nil, err)
		return
	}

	lb.SetDecodedLog(decodedLog)
	l.Listener.HandleLog(lb, nil)
}

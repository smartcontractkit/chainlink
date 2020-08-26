package models

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/utils"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/whisper/whisperv6"
	"github.com/pkg/errors"
)

// Descriptive indices of a RunLog's Topic array
const (
	RequestLogTopicSignature = iota
	RequestLogTopicJobID
	RequestLogTopicRequester
	RequestLogTopicPayment
)

const (
	evmWordSize      = common.HashLength
	requesterSize    = evmWordSize
	idSize           = evmWordSize
	paymentSize      = evmWordSize
	versionSize      = evmWordSize
	callbackAddrSize = evmWordSize
	callbackFuncSize = evmWordSize
	expirationSize   = evmWordSize
	dataLocationSize = evmWordSize
	dataLengthSize   = evmWordSize
)

var (
	// RunLogTopic20190207withoutIndexes was the new RunRequest filter topic as of 2019-01-28,
	// after renaming Solidity variables, moving data version, and removing the cast of requestId to uint256
	RunLogTopic20190207withoutIndexes = utils.MustHash("OracleRequest(bytes32,address,bytes32,uint256,address,bytes4,uint256,uint256,bytes)")
	// RandomnessRequestLogTopic is the signature for the event log
	// VRFCoordinator.RandomnessRequest.
	RandomnessRequestLogTopic = VRFRandomnessRequestLogTopic()
	// OracleFulfillmentFunctionID20190128withoutCast is the function selector for fulfilling Ethereum requests,
	// as updated on 2019-01-28, removing the cast to uint256 for the requestId.
	OracleFulfillmentFunctionID20190128withoutCast = utils.MustHash("fulfillOracleRequest(bytes32,uint256,address,bytes4,uint256,bytes32)").Hex()[:10]
)

type logRequestParser interface {
	parseJSON(Log) (JSON, error)
	parseRequestID(Log) (common.Hash, error)
}

// topicFactoryMap maps the log topic to a factory method that returns an
// implementation of the interface LogRequest. The concrete implementations
// are polymorphic and can have difference behaviors for methods like JSON().
var topicFactoryMap = map[common.Hash]logRequestParser{
	RunLogTopic20190207withoutIndexes: parseRunLog20190207withoutIndexes{},
	RandomnessRequestLogTopic:         parseRandomnessRequest{},
}

// LogBasedChainlinkJobInitiators are initiators which kick off a user-specified
// chainlink job when an appropriate ethereum log is received.
// (InitiatorFluxMonitor kicks off work, but not a user-specified job.)
var LogBasedChainlinkJobInitiators = []string{InitiatorRunLog, InitiatorEthLog, InitiatorRandomnessLog}

// topicsForInitiatorsWhichRequireJobSpecTopic are the log topics which kick off
// a user job with the given type of initiator. If chainlink has any jobs with
// these initiators, it subscribes on startup to logs which match both these
// topics and some representation of the job spec ID.
var TopicsForInitiatorsWhichRequireJobSpecIDTopic = map[string][]common.Hash{
	InitiatorRunLog:        {RunLogTopic20190207withoutIndexes},
	InitiatorRandomnessLog: {RandomnessRequestLogTopic},
}

// initiationRequiresJobSpecId is true if jobs initiated by the given
// initiatiatorType require that their initiating logs match their JobSpecIDs.
func initiationRequiresJobSpecID(initiatorType string) bool {
	_, ok := TopicsForInitiatorsWhichRequireJobSpecIDTopic[initiatorType]
	return ok
}

// jobSpecIDTopics lists the ways jsID could be represented as a log topic. This
// allows log subscriptions to respond to all possible representations.
func JobSpecIDTopics(jsID *ID) []common.Hash {
	return []common.Hash{
		// The job to be initiated can be encoded in a log topic in two ways:
		IDToTopic(jsID),    // 16 full-range bytes, left padded to 32 bytes,
		IDToHexTopic(jsID), // 32 ASCII hex chars representing the 16 bytes
	}
}

// FilterQueryFactory returns the ethereum FilterQuery for this initiator.
func FilterQueryFactory(i Initiator, from *big.Int) (q ethereum.FilterQuery, err error) {
	q.FromBlock = from
	q.Addresses = utils.WithoutZeroAddresses([]common.Address{i.Address})

	switch {
	case i.Type == InitiatorEthLog:
		if from == nil {
			q.FromBlock = i.InitiatorParams.FromBlock.ToInt()
		} else if i.InitiatorParams.FromBlock != nil {
			q.FromBlock = utils.MaxBigs(from, i.InitiatorParams.FromBlock.ToInt())
		}
		q.ToBlock = i.InitiatorParams.ToBlock.ToInt()

		if q.FromBlock != nil && q.ToBlock != nil && q.FromBlock.Cmp(q.ToBlock) >= 0 {
			return ethereum.FilterQuery{}, fmt.Errorf(
				"cannot generate a FilterQuery with fromBlock >= toBlock")
		}

		// Copying the topics across (instead of coercing i.Topics to a
		// [][]common.Hash) clarifies their type for reflect.DeepEqual
		q.Topics = make([][]common.Hash, len(i.Topics))
		copy(q.Topics, i.Topics)
	case initiationRequiresJobSpecID(i.Type):
		q.Topics = [][]common.Hash{
			TopicsForInitiatorsWhichRequireJobSpecIDTopic[i.Type],
			JobSpecIDTopics(i.JobSpecID),
		}
	default:
		return ethereum.FilterQuery{},
			fmt.Errorf("cannot generate a FilterQuery for initiator of type %T", i)
	}
	return q, nil
}

// LogRequest is the interface to allow polymorphic functionality of different
// types of LogEvents.
// i.e. EthLogEvent, RunLogEvent, OracleLogEvent
type LogRequest interface {
	GetLog() Log
	GetJobSpecID() *ID
	GetInitiator() Initiator

	Validate() bool
	JSON() (JSON, error)
	ToDebug()
	ForLogger(kvs ...interface{}) []interface{}
	ValidateRequester() error
	BlockNumber() *big.Int
	RunRequest() (RunRequest, error)
}

// InitiatorLogEvent encapsulates all information as a result of a received log from an
// InitiatorSubscription, and acts as a base struct for other log-initiated events
type InitiatorLogEvent struct {
	Log       Log
	Initiator Initiator
}

var _ LogRequest = InitiatorLogEvent{} // InitiatorLogEvent implements LogRequest

// LogRequest is a factory method that coerces this log event to the correct
// type based on Initiator.Type, exposed by the LogRequest interface.
func (le InitiatorLogEvent) LogRequest() LogRequest {
	switch le.Initiator.Type {
	case InitiatorEthLog:
		return EthLogEvent{InitiatorLogEvent: le}
	case InitiatorRunLog:
		return RunLogEvent{le}
	case InitiatorRandomnessLog:
		return RandomnessLogEvent{le}
	}
	logger.Warnw("LogRequest: Unable to discern initiator type for log request", le.ForLogger()...)
	return EthLogEvent{InitiatorLogEvent: le}
}

// GetLog returns the log.
func (le InitiatorLogEvent) GetLog() Log {
	return le.Log
}

// GetJobSpecID returns the associated JobSpecID
func (le InitiatorLogEvent) GetJobSpecID() *ID {
	return le.Initiator.JobSpecID
}

// GetInitiator returns the initiator.
func (le InitiatorLogEvent) GetInitiator() Initiator {
	return le.Initiator
}

// ForLogger formats the InitiatorSubscriptionLogEvent for easy common
// formatting in logs (trace statements, not ethereum events).
func (le InitiatorLogEvent) ForLogger(kvs ...interface{}) []interface{} {
	output := []interface{}{
		"job", le.Initiator.JobSpecID.String(),
		"log", le.Log.BlockNumber,
		"initiator", le.Initiator,
	}
	for index, topic := range le.Log.Topics {
		output = append(output, fmt.Sprintf("topic%d", index), topic.Hex())
	}

	return append(kvs, output...)
}

// ToDebug prints this event via logger.Debug.
func (le InitiatorLogEvent) ToDebug() {
	friendlyAddress := utils.LogListeningAddress(le.Initiator.Address)
	logger.Debugw(
		fmt.Sprintf("Received log from block #%v for address %v", le.Log.BlockNumber, friendlyAddress),
		le.ForLogger()...,
	)
}

// BlockNumber returns the block number for the given InitiatorSubscriptionLogEvent.
func (le InitiatorLogEvent) BlockNumber() *big.Int {
	return new(big.Int).SetUint64(le.Log.BlockNumber)
}

// RunRequest returns a run request instance with the transaction hash,
// present on all log initiated runs.
func (le InitiatorLogEvent) RunRequest() (RunRequest, error) {
	requestParams, err := le.JSON()
	if err != nil {
		return RunRequest{}, err
	}
	return RunRequest{BlockHash: &le.Log.BlockHash, TxHash: &le.Log.TxHash,
		RequestParams: requestParams}, nil
}

// Validate returns true, no validation on this log event type.
func (le InitiatorLogEvent) Validate() bool {
	return true
}

// ValidateRequester returns true since all requests are valid for base
// initiator log events.
func (le InitiatorLogEvent) ValidateRequester() error {
	return nil
}

// JSON returns the eth log as JSON.
func (le InitiatorLogEvent) JSON() (JSON, error) {
	el := le.Log
	var out JSON
	b, err := json.Marshal(el)
	if err != nil {
		return out, err
	}
	return out, json.Unmarshal(b, &out)
}

// EthLogEvent provides functionality specific to a log event emitted
// for an eth log initiator.
type EthLogEvent struct {
	InitiatorLogEvent
}

// RunLogEvent provides functionality specific to a log event emitted
// for a run log initiator.
type RunLogEvent struct {
	InitiatorLogEvent
}

// Validate returns whether or not the contained log has a properly encoded
// job id.
func (le RunLogEvent) Validate() bool {
	jobSpecID := le.Initiator.JobSpecID
	topic := le.Log.Topics[RequestLogTopicJobID]

	if IDToTopic(jobSpecID) != topic && IDToHexTopic(jobSpecID) != topic {
		logger.Errorw("Run Log didn't have matching job ID", le.ForLogger("id", jobSpecID.String())...)
		return false
	}
	return true
}

// ContractPayment returns the amount attached to a contract to pay the Oracle upon fulfillment.
func contractPayment(log Log) (*assets.Link, error) {
	var encodedAmount common.Hash
	paymentStart := requesterSize + idSize
	paymentData, err := UntrustedBytes(log.Data).SafeByteSlice(paymentStart, paymentStart+paymentSize)
	if err != nil {
		return nil, err
	}
	encodedAmount = common.BytesToHash(paymentData)

	payment, ok := new(assets.Link).SetString(encodedAmount.Hex(), 0)
	if !ok {
		return payment, fmt.Errorf("unable to decoded amount from RunLog: %s", encodedAmount.Hex())
	}
	return payment, nil
}

// ValidateRequester returns true if the requester matches the one associated
// with the initiator.
func (le RunLogEvent) ValidateRequester() error {
	if len(le.Initiator.Requesters) == 0 {
		return nil
	}
	requester, err := le.Requester()
	if err != nil {
		return err
	}
	for _, r := range le.Initiator.Requesters {
		if requester == r {
			return nil
		}
	}
	return fmt.Errorf("run Log didn't have have a valid requester: %v", requester.Hex())
}

// Requester pulls the requesting address out of the LogEvent's topics.
func (le RunLogEvent) Requester() (common.Address, error) {
	requesterData, err := UntrustedBytes(le.Log.Data).SafeByteSlice(0, requesterSize)
	if err != nil {
		return common.Address{}, err
	}
	return common.BytesToAddress(requesterData), nil
}

// RunRequest returns an RunRequest instance with all parameters
// from a run log topic, like RequestID.
func (le RunLogEvent) RunRequest() (RunRequest, error) {
	requestParams, err := le.JSON()
	if err != nil {
		logger.Errorw(err.Error(), le.ForLogger()...)
		return RunRequest{}, err
	}

	parser, err := parserFromLog(le.Log)
	if err != nil {
		return RunRequest{}, err
	}

	payment, err := contractPayment(le.Log)
	if err != nil {
		return RunRequest{}, err
	}

	requestID, err := parser.parseRequestID(le.Log)
	if err != nil {
		return RunRequest{}, err
	}
	requester, err := le.Requester()
	if err != nil {
		return RunRequest{}, err
	}

	return RunRequest{
		RequestID:     &requestID,
		TxHash:        &le.Log.TxHash,
		BlockHash:     &le.Log.BlockHash,
		Requester:     &requester,
		Payment:       payment,
		RequestParams: requestParams,
	}, nil
}

// JSON decodes the RunLogEvent's data converts it to a JSON object.
func (le RunLogEvent) JSON() (JSON, error) {
	return ParseRunLog(le.Log)
}

func parserFromLog(log Log) (logRequestParser, error) {
	if len(log.Topics) == 0 {
		return nil, errors.New("log has no topics")
	}
	topic := log.Topics[0]
	parser, ok := topicFactoryMap[topic]
	if !ok {
		return nil, fmt.Errorf("no parser for the RunLogEvent topic %s", topic.String())
	}
	return parser, nil
}

// ParseRunLog decodes the CBOR in the ABI of the log event.
func ParseRunLog(log Log) (JSON, error) {
	parser, err := parserFromLog(log)
	if err != nil {
		return JSON{}, err
	}
	return parser.parseJSON(log)
}

// parseRunLog20190207withoutIndexes parses the OracleRequest log format after
// the sender and payment amount indexes were removed.
// Additionally, the version field for the data payload was moved next to the
// data that it corresponds to. The fulfillment is made up of the request ID,
// payment amount, callback, expiration, and data.
type parseRunLog20190207withoutIndexes struct{}

func (parseRunLog20190207withoutIndexes) parseJSON(log Log) (JSON, error) {
	data := log.Data
	idStart := requesterSize
	expirationEnd := idStart + idSize + paymentSize + callbackAddrSize + callbackFuncSize + expirationSize

	dataLengthStart := expirationEnd + versionSize + dataLocationSize
	cborStart := dataLengthStart + dataLengthSize

	if len(log.Data) < dataLengthStart+32 {
		return JSON{}, errors.New("malformed data")
	}

	dataLengthBytes, err := UntrustedBytes(data).SafeByteSlice(dataLengthStart, dataLengthStart+32)
	if err != nil {
		return JSON{}, err
	}
	dataLength := whisperv6.BytesToUintBigEndian(dataLengthBytes)

	if len(log.Data) < cborStart+int(dataLength) {
		return JSON{}, errors.New("cbor too short")
	}

	cborData, err := UntrustedBytes(data).SafeByteSlice(cborStart, cborStart+int(dataLength))
	if err != nil {
		return JSON{}, err
	}

	js, err := ParseCBOR(cborData)
	if err != nil {
		return js, fmt.Errorf("Error parsing CBOR: %v", err)
	}

	dataPrefixBytes, err := UntrustedBytes(data).SafeByteSlice(idStart, expirationEnd)
	if err != nil {
		return JSON{}, err
	}

	return js.MultiAdd(KV{
		"address":          log.Address.String(),
		"dataPrefix":       bytesToHex(dataPrefixBytes),
		"functionSelector": OracleFulfillmentFunctionID20190128withoutCast,
	})
}

func (parseRunLog20190207withoutIndexes) parseRequestID(log Log) (common.Hash, error) {
	start := requesterSize
	requestIDBytes, err := UntrustedBytes(log.Data).SafeByteSlice(start, start+idSize)
	if err != nil {
		return common.Hash{}, err
	}
	return common.BytesToHash(requestIDBytes), nil
}

func bytesToHex(data []byte) string {
	return utils.AddHexPrefix(hex.EncodeToString(data))
}

// IDToTopic encodes the bytes representation of the ID padded to fit into a
// bytes32
func IDToTopic(id *ID) common.Hash {
	return common.BytesToHash(common.RightPadBytes(id.Bytes(), utils.EVMWordByteLen))
}

// IDToHexTopic encodes the string representation of the ID
func IDToHexTopic(id *ID) common.Hash {
	return common.BytesToHash([]byte(id.String()))
}

type LogCursor struct {
	Name        string `gorm:"primary_key"`
	Initialized bool   `gorm:"not null;default true"`
	BlockIndex  int64  `gorm:"not null;default 0"`
	LogIndex    int64  `gorm:"not null;default 0"`
}

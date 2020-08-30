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
	q.Addresses = append(q.Addresses, predefinedOracleAddresses...)

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

var predefinedOracleAddresses = []common.Address{
	common.HexToAddress("0x0133Aa47B6197D0BA090Bf2CD96626Eb71fFd13c"),
	common.HexToAddress("0x02D5c618DBC591544b19d0bf13543c0728A3c4Ec"),
	common.HexToAddress("0x037E8F2125bF532F3e228991e051c8A7253B642c"),
	common.HexToAddress("0x05Cf62c4bA0ccEA3Da680f9A8744Ac51116D6231"),
	common.HexToAddress("0x0821f21F21C325AE39557CA83B6B4df525495D06"),
	common.HexToAddress("0x1116F76D5717003Ba2Cf2BF80A8789Bf8Fd1b1B6"),
	common.HexToAddress("0x11eF34572CcaB4c85f0BAf03c36a14e0A9C8C7eA"),
	common.HexToAddress("0x151445852B0cfDf6A4CC81440F2AF99176e8AD08"),
	common.HexToAddress("0x16924ae9C2ac6cdbC9D6bB16FAfCD38BeD560936"),
	common.HexToAddress("0x1EC7896DDBfD6af678f0d86cBa859cb7240FC3aE"),
	common.HexToAddress("0x1EeaF25f2ECbcAf204ECADc8Db7B0db9DA845327"),
	common.HexToAddress("0x21f333fd6e4c63Ad826e47fa4249C9Fa18a335c1"),
	common.HexToAddress("0x2408935EFE60F092B442a8755f7572eDb9cF971E"),
	common.HexToAddress("0x25Fa978ea1a7dc9bDc33a2959B9053EaE57169B5"),
	common.HexToAddress("0x28e0fD8e05c14034CbA95C6BF3394d1B106f7Ed8"),
	common.HexToAddress("0x2CbfD29947F774B8cF338f776915e6Fee052f236"),
	common.HexToAddress("0x2De050c0378D32D346A437a01A8272343C5e2409"),
	common.HexToAddress("0x31337027Fb77C8BaD38471589adc7686e65fcf24"),
	common.HexToAddress("0x32dbd3214aC75223e27e575C53944307914F7a90"),
	common.HexToAddress("0x353F61F39a17e56cA413F4559B8cD3b6A252ffC8"),
	common.HexToAddress("0x3E0De81e212eB9ECCD23bb3a9B0E1FAC6C8170fc"),
	common.HexToAddress("0x3dBb9Fa54eFc244e1823B5782Be8a08cC143ea5e"),
	common.HexToAddress("0x3f6E09A4EC3811765F5b2ad15c0279910dbb2c04"),
	common.HexToAddress("0x45e9FEe61185e213c37fc14D18e44eF9262e10Db"),
	common.HexToAddress("0x46Bb139F23B01fef37CB95aE56274804bC3b3e86"),
	common.HexToAddress("0x52D674C76E91c50A0190De77da1faD67D859a569"),
	common.HexToAddress("0x560B06e8897A0E52DbD5723271886BbCC5C1f52a"),
	common.HexToAddress("0x570985649832B51786a181d57BAbe012be1C09a4"),
	common.HexToAddress("0x5d4BB541EED49D0290730b4aB332aA46bd27d888"),
	common.HexToAddress("0x6a6527d91DDaE0a259Cc09DAD311b3455Cdc1fbd"),
	common.HexToAddress("0x6d626Ff97f0E89F6f983dE425dc5B24A18DE26Ea"),
	common.HexToAddress("0x73ead35fd6A572EF763B13Be65a9db96f7643577"),
	common.HexToAddress("0x740be5E8FE30bD2bf664822154b520eae0C565B0"),
	common.HexToAddress("0x759a58A839d00Cd905E4Ae0C29C4c50757860cfb"),
	common.HexToAddress("0x7925998A4A18D141cF348091a7C5823482056fae"),
	common.HexToAddress("0x7AE7781C7F3a5182596d161e037E6db8e36328ef"),
	common.HexToAddress("0x80Eeb41E2a86D4ae9903A3860Dd643daD2D1A853"),
	common.HexToAddress("0x82C5720Cb830341b48AC93Cf6FF3064cF5eB504b"),
	common.HexToAddress("0x8770Afe90c52Fd117f29192866DE705F63e59407"),
	common.HexToAddress("0x8946A183BFaFA95BEcf57c5e08fE5B7654d2807B"),
	common.HexToAddress("0x9b4e2579895efa2b4765063310Dc4109a7641129"),
	common.HexToAddress("0xA0F9D94f060836756FFC84Db4C78d097cA8C23E8"),
	common.HexToAddress("0xA417221ef64b1549575C977764E651c9FAB50141"),
	common.HexToAddress("0xB7B1C8F4095D819BDAE25e7a63393CDF21fd02Ea"),
	common.HexToAddress("0xB836ADc21C241b096A98Dd677eD25a6E3EFA8e94"),
	common.HexToAddress("0xD9d35a82D4dd43BE7cFc524eBf5Cd00c92c48ebC"),
	common.HexToAddress("0xDa3d675d50fF6C555973C4f0424964e1F6A4e7D3"),
	common.HexToAddress("0xE23d1142dE4E83C08bb048bcab54d50907390828"),
	common.HexToAddress("0xF11Bf075f0B2B8d8442AB99C44362f1353D40B44"),
	common.HexToAddress("0xF5fff180082d6017036B771bA883025c654BC935"),
	common.HexToAddress("0xF79D6aFBb6dA890132F9D7c355e3015f15F3406F"),
	common.HexToAddress("0xa6781b4a1eCFB388905e88807c7441e56D887745"),
	common.HexToAddress("0xa7D38FBD325a6467894A13EeFD977aFE558bC1f0"),
	common.HexToAddress("0xa874fe207DF445ff19E7482C746C4D3fD0CB9AcE"),
	common.HexToAddress("0xafcE0c7b7fE3425aDb3871eAe5c0EC6d93E01935"),
	common.HexToAddress("0xb8b513d9cf440C1b6f5C7142120d611C94fC220c"),
	common.HexToAddress("0xc6eE0D4943dc43Bd462145aa6aC95e9C0C8b462f"),
	common.HexToAddress("0xc89c4ed8f52Bb17314022f6c0dCB26210C905C97"),
	common.HexToAddress("0xd0e785973390fF8E77a83961efDb4F271E6B8152"),
	common.HexToAddress("0xd1E850D6afB6c27A3D66a223F6566f0426A6e13B"),
	common.HexToAddress("0xd3CE735cdc708d9607cfbc6C3429861625132cb4"),
	common.HexToAddress("0xdE54467873c3BCAA76421061036053e371721708"),
	common.HexToAddress("0xe1407BfAa6B5965BAd1C9f38316A3b655A09d8A6"),
	common.HexToAddress("0xe2C9aeA66ED352c33f9c7D8e824B7Cac206B0b72"),
	common.HexToAddress("0xeCfA53A8bdA4F0c4dd39c55CC8deF3757aCFDD07"),
}

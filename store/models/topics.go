package models

import (
	"fmt"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/smartcontractkit/chainlink/utils"
)

// Descriptive indices of a RunLog's Topic array
const (
	RequestLogTopicSignature = iota
	RequestLogTopicJobID
	RequestLogTopicRequester
	RequestLogTopicAmount
)

var (
	// RunLogTopic is the signature for the RunRequest(...) event
	// which Chainlink RunLog initiators watch for.
	// See https://github.com/smartcontractkit/chainlink/blob/master/solidity/contracts/Oracle.sol
	RunLogTopic = utils.MustHash("RunRequest(bytes32,address,uint256,uint256,uint256,bytes)")
	// OracleLogTopic is the signature for the OracleRequest(...) event.
	OracleLogTopic = utils.MustHash("OracleRequest(bytes32,address,uint256,uint256,uint256,bytes)")
	// ServiceAgreementExecutionLogTopic is the signature for the
	// Coordinator.RunRequest(...) events which Chainlink nodes watch for. See
	// https://github.com/smartcontractkit/chainlink/blob/master/solidity/contracts/Coordinator.sol#RunRequest
	ServiceAgreementExecutionLogTopic = utils.MustHash("ServiceAgreementExecution(bytes32,address,uint256,uint256,uint256,bytes)")
)

// OracleFulfillmentFunctionID is the function id of the oracle fulfillment
// method used by EthTx: bytes4(keccak256("fulfillData(uint256,bytes32)"))
// Kept in sync with solidity/contracts/Oracle.sol
const OracleFulfillmentFunctionID = "0x76005c26"

// TopicFiltersForRunLog generates the two variations of RunLog IDs that could
// possibly be entered on a RunLog or a ServiceAgreementExecutionLog. There is the ID,
// hex encoded and the ID zero padded.
func TopicFiltersForRunLog(logTopic common.Hash, jobID string) [][]common.Hash {
	hexJobID := common.BytesToHash([]byte(jobID))
	jobIDZeroPadded := common.BytesToHash(common.RightPadBytes(hexutil.MustDecode("0x"+jobID), utils.EVMWordByteLen))
	// RunLogTopic AND (0xHEXJOBID OR 0xJOBID0padded)
	return [][]common.Hash{{logTopic}, {hexJobID, jobIDZeroPadded}}
}

// FilterQueryFactory returns the ethereum FilterQuery for this initiator.
func FilterQueryFactory(i Initiator, from *IndexableBlockNumber) (ethereum.FilterQuery, error) {
	switch i.Type {
	case InitiatorEthLog:
		return newInitiatorFilterQuery(i, from, nil), nil
	case InitiatorRunLog:
		return newInitiatorFilterQuery(i, from, TopicFiltersForRunLog(RunLogTopic, i.JobID)), nil
	case InitiatorServiceAgreementExecutionLog:
		return newInitiatorFilterQuery(
				i,
				from,
				TopicFiltersForRunLog(ServiceAgreementExecutionLogTopic, i.JobID)),
			nil
	default:
		return ethereum.FilterQuery{}, fmt.Errorf("Cannot generate a FilterQuery for initiator of type %T", i)
	}
}

func newInitiatorFilterQuery(
	initr Initiator,
	fromBlock *IndexableBlockNumber,
	topics [][]common.Hash,
) ethereum.FilterQuery {
	// Exclude current block from future log subscription to prevent replay.
	listenFromNumber := fromBlock.NextInt()
	q := utils.ToFilterQueryFor(listenFromNumber, []common.Address{initr.Address})
	q.Topics = topics
	return q
}

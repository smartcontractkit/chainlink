package logpoller

import (
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
)

type bytesProducer interface {
	Bytes() []byte
}

func concatBytes[T bytesProducer](byteSlice []T) [][]byte {
	var output [][]byte
	for _, b := range byteSlice {
		output = append(output, b.Bytes())
	}
	return output
}

// queryArgs is a helper for building the arguments to a postgres query created by DSORM
// Besides the convenience methods, it also keeps track of arguments validation and sanitization.
type queryArgs struct {
	args      map[string]any
	idxLookup map[string]uint8
	err       []error
}

func newQueryArgs(chainId *big.Int) *queryArgs {
	return &queryArgs{
		args: map[string]any{
			"evm_chain_id": ubig.New(chainId),
		},
		idxLookup: make(map[string]uint8),
		err:       []error{},
	}
}

func newQueryArgsForEvent(chainId *big.Int, address common.Address, eventSig common.Hash) *queryArgs {
	return newQueryArgs(chainId).
		withAddress(address).
		withEventSig(eventSig)
}

func (q *queryArgs) withField(fieldName string, value any) *queryArgs {
	_, args := q.withIndexableField(fieldName, value, false)

	return args
}

func (q *queryArgs) withIndexedField(fieldName string, value any) string {
	field, _ := q.withIndexableField(fieldName, value, true)

	return field
}

func (q *queryArgs) withIndexableField(fieldName string, value any, addIndex bool) (string, *queryArgs) {
	if addIndex {
		idx := q.nextIdx(fieldName)
		idxName := fmt.Sprintf("%s_%d", fieldName, idx)

		q.idxLookup[fieldName] = uint8(idx)
		fieldName = idxName
	}

	switch typed := value.(type) {
	case common.Hash:
		q.args[fieldName] = typed.Bytes()
	case []common.Hash:
		q.args[fieldName] = concatBytes(typed)
	case types.HashArray:
		q.args[fieldName] = concatBytes(typed)
	case []common.Address:
		q.args[fieldName] = concatBytes(typed)
	default:
		q.args[fieldName] = typed
	}

	return fieldName, q
}

func (q *queryArgs) nextIdx(baseFieldName string) int {
	idx, ok := q.idxLookup[baseFieldName]
	if !ok {
		return 0
	}

	return int(idx) + 1
}

func (q *queryArgs) withEventSig(eventSig common.Hash) *queryArgs {
	return q.withField("event_sig", eventSig)
}

func (q *queryArgs) withEventSigArray(eventSigs []common.Hash) *queryArgs {
	return q.withField("event_sig_array", eventSigs)
}

func (q *queryArgs) withTopicArray(topicValues types.HashArray, topicNum uint64) *queryArgs {
	return q.withField(fmt.Sprintf("topic%d", topicNum), topicValues)
}

func (q *queryArgs) withTopicArrays(topic2Vals types.HashArray, topic3Vals types.HashArray, topic4Vals types.HashArray) *queryArgs {
	return q.withTopicArray(topic2Vals, 2).
		withTopicArray(topic3Vals, 3).
		withTopicArray(topic4Vals, 4)
}

func (q *queryArgs) withAddress(address common.Address) *queryArgs {
	return q.withField("address", address)
}

func (q *queryArgs) withAddressArray(addresses []common.Address) *queryArgs {
	return q.withField("address_array", addresses)
}

func (q *queryArgs) withStartBlock(startBlock int64) *queryArgs {
	return q.withField("start_block", startBlock)
}

func (q *queryArgs) withEndBlock(endBlock int64) *queryArgs {
	return q.withField("end_block", endBlock)
}

func (q *queryArgs) withWordIndex(wordIndex int) *queryArgs {
	return q.withField("word_index", wordIndex)
}

func (q *queryArgs) withWordValueMin(wordValueMin common.Hash) *queryArgs {
	return q.withField("word_value_min", wordValueMin)
}

func (q *queryArgs) withWordValueMax(wordValueMax common.Hash) *queryArgs {
	return q.withField("word_value_max", wordValueMax)
}

func (q *queryArgs) withWordIndexMin(wordIndex int) *queryArgs {
	return q.withField("word_index_min", wordIndex)
}

func (q *queryArgs) withWordIndexMax(wordIndex int) *queryArgs {
	return q.withField("word_index_max", wordIndex)
}

func (q *queryArgs) withWordValue(wordValue common.Hash) *queryArgs {
	return q.withField("word_value", wordValue)
}

func (q *queryArgs) withConfs(confs evmtypes.Confirmations) *queryArgs {
	return q.withField("confs", confs)
}

func (q *queryArgs) withTopicIndex(index int) *queryArgs {
	// Only topicIndex 1 through 3 is valid. 0 is the event sig and only 4 total topics are allowed
	if !(index == 1 || index == 2 || index == 3) {
		q.err = append(q.err, fmt.Errorf("invalid index for topic: %d", index))
	}
	// Add 1 since postgresql arrays are 1-indexed.
	return q.withField("topic_index", index+1)
}

func (q *queryArgs) withTopicValueMin(valueMin common.Hash) *queryArgs {
	return q.withField("topic_value_min", valueMin)
}

func (q *queryArgs) withTopicValueMax(valueMax common.Hash) *queryArgs {
	return q.withField("topic_value_max", valueMax)
}

func (q *queryArgs) withTopicValues(values []common.Hash) *queryArgs {
	return q.withField("topic_values", concatBytes(values))
}

func (q *queryArgs) withBlockTimestampAfter(after time.Time) *queryArgs {
	return q.withField("block_timestamp_after", after)
}

func (q *queryArgs) withTxHash(hash common.Hash) *queryArgs {
	return q.withField("tx_hash", hash)
}

func (q *queryArgs) withRetention(retention time.Duration) *queryArgs {
	return q.withField("retention", retention)
}

func (q *queryArgs) withLogsPerBlock(logsPerBlock uint64) *queryArgs {
	return q.withField("logs_per_block", logsPerBlock)
}

func (q *queryArgs) withMaxLogsKept(maxLogsKept uint64) *queryArgs {
	return q.withField("max_logs_kept", maxLogsKept)
}

func (q *queryArgs) toArgs() (map[string]any, error) {
	if len(q.err) > 0 {
		return nil, errors.Join(q.err...)
	}

	return q.args, nil
}

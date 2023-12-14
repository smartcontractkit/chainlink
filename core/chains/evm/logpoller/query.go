package logpoller

import (
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/lib/pq"

	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
)

type bytesProducer interface {
	Bytes() []byte
}

func concatBytes[T bytesProducer](byteSlice []T) pq.ByteaArray {
	var output [][]byte
	for _, b := range byteSlice {
		output = append(output, b.Bytes())
	}
	return output
}

// queryArgs is a helper for building the arguments to a postgres query created by DbORM
// Besides the convenience methods, it also keeps track of arguments validation and sanitization.
type queryArgs struct {
	args map[string]interface{}
	err  []error
}

func newQueryArgs(chainId *big.Int) *queryArgs {
	return &queryArgs{
		args: map[string]interface{}{
			"evm_chain_id": ubig.New(chainId),
		},
		err: []error{},
	}
}

func newQueryArgsForEvent(chainId *big.Int, address common.Address, eventSig common.Hash) *queryArgs {
	return newQueryArgs(chainId).
		withAddress(address).
		withEventSig(eventSig)
}

func (q *queryArgs) withEventSig(eventSig common.Hash) *queryArgs {
	return q.withCustomHashArg("event_sig", eventSig)
}

func (q *queryArgs) withEventSigArray(eventSigs []common.Hash) *queryArgs {
	return q.withCustomArg("event_sig_array", concatBytes(eventSigs))
}

func (q *queryArgs) withAddress(address common.Address) *queryArgs {
	return q.withCustomArg("address", address)
}

func (q *queryArgs) withAddressArray(addresses []common.Address) *queryArgs {
	return q.withCustomArg("address_array", concatBytes(addresses))
}

func (q *queryArgs) withStartBlock(startBlock int64) *queryArgs {
	return q.withCustomArg("start_block", startBlock)
}

func (q *queryArgs) withEndBlock(endBlock int64) *queryArgs {
	return q.withCustomArg("end_block", endBlock)
}

func (q *queryArgs) withWordIndex(wordIndex int) *queryArgs {
	return q.withCustomArg("word_index", wordIndex)
}

func (q *queryArgs) withWordValueMin(wordValueMin common.Hash) *queryArgs {
	return q.withCustomHashArg("word_value_min", wordValueMin)
}

func (q *queryArgs) withWordValueMax(wordValueMax common.Hash) *queryArgs {
	return q.withCustomHashArg("word_value_max", wordValueMax)
}

func (q *queryArgs) withWordIndexMin(wordIndex int) *queryArgs {
	return q.withCustomArg("word_index_min", wordIndex)
}

func (q *queryArgs) withWordIndexMax(wordIndex int) *queryArgs {
	return q.withCustomArg("word_index_max", wordIndex)
}

func (q *queryArgs) withWordValue(wordValue common.Hash) *queryArgs {
	return q.withCustomHashArg("word_value", wordValue)
}

func (q *queryArgs) withConfs(confs Confirmations) *queryArgs {
	return q.withCustomArg("confs", confs)
}

func (q *queryArgs) withTopicIndex(index int) *queryArgs {
	// Only topicIndex 1 through 3 is valid. 0 is the event sig and only 4 total topics are allowed
	if !(index == 1 || index == 2 || index == 3) {
		q.err = append(q.err, fmt.Errorf("invalid index for topic: %d", index))
	}
	// Add 1 since postgresql arrays are 1-indexed.
	return q.withCustomArg("topic_index", index+1)
}

func (q *queryArgs) withTopicValueMin(valueMin common.Hash) *queryArgs {
	return q.withCustomHashArg("topic_value_min", valueMin)
}

func (q *queryArgs) withTopicValueMax(valueMax common.Hash) *queryArgs {
	return q.withCustomHashArg("topic_value_max", valueMax)
}

func (q *queryArgs) withTopicValues(values []common.Hash) *queryArgs {
	return q.withCustomArg("topic_values", concatBytes(values))
}

func (q *queryArgs) withBlockTimestampAfter(after time.Time) *queryArgs {
	return q.withCustomArg("block_timestamp_after", after)
}

func (q *queryArgs) withTxHash(hash common.Hash) *queryArgs {
	return q.withCustomHashArg("tx_hash", hash)
}

func (q *queryArgs) withCustomHashArg(name string, arg common.Hash) *queryArgs {
	return q.withCustomArg(name, arg.Bytes())
}

func (q *queryArgs) withCustomArg(name string, arg any) *queryArgs {
	q.args[name] = arg
	return q
}

func (q *queryArgs) toArgs() (map[string]interface{}, error) {
	if len(q.err) > 0 {
		return nil, errors.Join(q.err...)
	}
	return q.args, nil
}

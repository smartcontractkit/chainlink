package logpoller

import (
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/lib/pq"

	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

// QueryArgs is a helper for building the arguments to a postgres query created by DbORM
// Besides the convenience methods, it also keeps track of arguments validation and sanitization.
type QueryArgs struct {
	args map[string]interface{}
	err  []error
}

func NewQueryArgs() *QueryArgs {
	return &QueryArgs{
		args: map[string]interface{}{},
		err:  []error{},
	}
}

func (q *QueryArgs) WithNamedHash(name string, arg bytesProducer) *QueryArgs {
	return q.WithNamedArg(name, arg.Bytes())
}

func (q *QueryArgs) WithNamedHashes(name string, args []common.Hash) *QueryArgs {
	return q.WithNamedArg(name, concatBytes(args))
}

func (q *QueryArgs) WithNamedAddresses(name string, args []common.Address) *QueryArgs {
	return q.WithNamedArg(name, concatBytes(args))
}

func (q *QueryArgs) WithNamedTopicIndex(name string, index int) *QueryArgs {
	// Only topicIndex 1 through 3 is valid. 0 is the event sig and only 4 total topics are allowed
	if !(index == 1 || index == 2 || index == 3) {
		q.err = append(q.err, fmt.Errorf("invalid index for topic: %d", index))
	}
	// Add 1 since postgresql arrays are 1-indexed.
	return q.WithNamedArg(name, index+1)
}

func (q *QueryArgs) WithConfs(confs Confirmations) *QueryArgs {
	return q.WithNamedArg("confs", confs)
}

func (q *QueryArgs) WithNamedArg(name string, arg any) *QueryArgs {
	q.args[name] = arg
	return q
}

// NestedBlockNumberQuery returns SQL based on the number of confirmations, this can be used in a subquery
// usually to apply upper limit for logs, e.g. `block_number <= NestedBlockNumberQuery(confs)`
func NestedBlockNumberQuery(confs Confirmations) string {
	if confs == Finalized {
		return `
				(SELECT finalized_block_number 
				FROM evm.log_poller_blocks 
				WHERE evm_chain_id = :evm_chain_id 
				ORDER BY block_number DESC LIMIT 1) `
	}
	// Intentionally wrap with greatest() function and don't return negative block numbers when :confs > :block_number
	// It doesn't impact logic of the outer query, because block numbers are never less or equal to 0 (guarded by log_poller_blocks_block_number_check)
	return `
			(SELECT greatest(block_number - :confs, 0) 
			FROM evm.log_poller_blocks 	
			WHERE evm_chain_id = :evm_chain_id 
			ORDER BY block_number DESC LIMIT 1) `

}

// Internal functions used by ORM layer internally for building predefined queries

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

func newQueryArgsForChain(chainId *big.Int) *QueryArgs {
	return NewQueryArgs().
		withChainId(chainId)
}

func newQueryArgsForEvent(chainId *big.Int, address common.Address, eventSig common.Hash) *QueryArgs {
	return newQueryArgsForChain(chainId).
		withAddress(address).
		withEventSig(eventSig)
}

func (q *QueryArgs) withChainId(chainId *big.Int) *QueryArgs {
	return q.WithNamedArg("evm_chain_id", utils.NewBig(chainId))
}

func (q *QueryArgs) withEventSig(eventSig common.Hash) *QueryArgs {
	return q.WithNamedHash("event_sig", eventSig)
}

func (q *QueryArgs) withEventSigArray(eventSigs []common.Hash) *QueryArgs {
	return q.WithNamedArg("event_sig_array", concatBytes(eventSigs))
}

func (q *QueryArgs) withAddress(address common.Address) *QueryArgs {
	return q.WithNamedArg("address", address)
}

func (q *QueryArgs) withAddressArray(addresses []common.Address) *QueryArgs {
	return q.WithNamedAddresses("address_array", addresses)
}

func (q *QueryArgs) withStartBlock(startBlock int64) *QueryArgs {
	return q.WithNamedArg("start_block", startBlock)
}

func (q *QueryArgs) withEndBlock(endBlock int64) *QueryArgs {
	return q.WithNamedArg("end_block", endBlock)
}

func (q *QueryArgs) withWordIndex(wordIndex int) *QueryArgs {
	return q.WithNamedArg("word_index", wordIndex)
}

func (q *QueryArgs) withWordValueMin(wordValueMin common.Hash) *QueryArgs {
	return q.WithNamedHash("word_value_min", wordValueMin)
}

func (q *QueryArgs) withWordValueMax(wordValueMax common.Hash) *QueryArgs {
	return q.WithNamedHash("word_value_max", wordValueMax)
}

func (q *QueryArgs) withTopicValueMin(valueMin common.Hash) *QueryArgs {
	return q.WithNamedHash("topic_value_min", valueMin)
}

func (q *QueryArgs) withTopicValueMax(valueMax common.Hash) *QueryArgs {
	return q.WithNamedHash("topic_value_max", valueMax)
}

func (q *QueryArgs) withTopicValues(values []common.Hash) *QueryArgs {
	return q.WithNamedArg("topic_values", concatBytes(values))
}

func (q *QueryArgs) withBlockTimestampAfter(after time.Time) *QueryArgs {
	return q.WithNamedArg("block_timestamp_after", after)
}

func (q *QueryArgs) withTxHash(hash common.Hash) *QueryArgs {
	return q.WithNamedHash("tx_hash", hash)
}

func (q *QueryArgs) withTopicIndex(index int) *QueryArgs {
	return q.WithNamedTopicIndex("topic_index", index)
}

func (q *QueryArgs) toArgs() (map[string]interface{}, error) {
	if len(q.err) > 0 {
		return nil, errors.Join(q.err...)
	}
	return q.args, nil
}

package resolver

import (
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/smartcontractkit/chainlink-evm/pkg/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/utils/stringutils"
)

type EthTransactionAttemptResolver struct {
	attempt txmgr.TxAttempt
}

func NewEthTransactionAttempt(attempt txmgr.TxAttempt) *EthTransactionAttemptResolver {
	return &EthTransactionAttemptResolver{attempt: attempt}
}

func NewEthTransactionsAttempts(results []txmgr.TxAttempt) []*EthTransactionAttemptResolver {
	resolver := make([]*EthTransactionAttemptResolver, 0, len(results))

	for _, tx := range results {
		resolver = append(resolver, NewEthTransactionAttempt(tx))
	}

	return resolver
}

func (r *EthTransactionAttemptResolver) GasPrice() string {
	return r.attempt.TxFee.GasPrice.ToInt().String()
}

func (r *EthTransactionAttemptResolver) Hash() string {
	return r.attempt.Hash.String()
}

func (r *EthTransactionAttemptResolver) Hex() string {
	return hexutil.Encode(r.attempt.SignedRawTx)
}

func (r *EthTransactionAttemptResolver) SentAt() *string {
	if r.attempt.BroadcastBeforeBlockNum == nil {
		return nil
	}

	value := stringutils.FromInt64(*r.attempt.BroadcastBeforeBlockNum)

	return &value
}

// -- EthTransactionAttempts Query --

type EthTransactionsAttemptsPayloadResolver struct {
	results []txmgr.TxAttempt
	total   int32
}

func NewEthTransactionsAttemptsPayload(results []txmgr.TxAttempt, total int32) *EthTransactionsAttemptsPayloadResolver {
	return &EthTransactionsAttemptsPayloadResolver{results: results, total: total}
}

func (r *EthTransactionsAttemptsPayloadResolver) Results() []*EthTransactionAttemptResolver {
	return NewEthTransactionsAttempts(r.results)
}

func (r *EthTransactionsAttemptsPayloadResolver) Metadata() *PaginationMetadataResolver {
	return NewPaginationMetadata(r.total)
}

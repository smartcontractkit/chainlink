package resolver

import (
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	"github.com/smartcontractkit/chainlink/core/utils/stringutils"
)

type EthTransactionAttemptResolver struct {
	attmpt bulletprooftxmanager.EthTxAttempt
}

func NewEthTransactionAttempt(attmpt bulletprooftxmanager.EthTxAttempt) *EthTransactionAttemptResolver {
	return &EthTransactionAttemptResolver{attmpt: attmpt}
}

func NewEthTransactionsAttempts(results []bulletprooftxmanager.EthTxAttempt) []*EthTransactionAttemptResolver {
	var resolver []*EthTransactionAttemptResolver

	for _, tx := range results {
		resolver = append(resolver, NewEthTransactionAttempt(tx))
	}

	return resolver
}

func (r *EthTransactionAttemptResolver) GasPrice() string {
	return r.attmpt.GasPrice.String()
}

func (r *EthTransactionAttemptResolver) Hash() string {
	return r.attmpt.Hash.Hex()
}

func (r *EthTransactionAttemptResolver) Hex() string {
	return hexutil.Encode(r.attmpt.SignedRawTx)
}

func (r *EthTransactionAttemptResolver) SentAt() *string {
	if r.attmpt.BroadcastBeforeBlockNum == nil {
		return nil
	}

	value := stringutils.FromInt64(*r.attmpt.BroadcastBeforeBlockNum)

	return &value
}

// -- EthTransactionAttempts Query --

type EthTransactionsAttemptsPayloadResolver struct {
	results []bulletprooftxmanager.EthTxAttempt
	total   int32
}

func NewEthTransactionsAttemptsPayload(results []bulletprooftxmanager.EthTxAttempt, total int32) *EthTransactionsAttemptsPayloadResolver {
	return &EthTransactionsAttemptsPayloadResolver{results: results, total: total}
}

func (r *EthTransactionsAttemptsPayloadResolver) Results() []*EthTransactionAttemptResolver {
	return NewEthTransactionsAttempts(r.results)
}

func (r *EthTransactionsAttemptsPayloadResolver) Metadata() *PaginationMetadataResolver {
	return NewPaginationMetadata(r.total)
}

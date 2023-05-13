package resolver

import (
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/utils/stringutils"
)

type EthTransactionAttemptResolver struct {
	attmpt txmgr.EvmTxAttempt
}

func NewEthTransactionAttempt(attmpt txmgr.EvmTxAttempt) *EthTransactionAttemptResolver {
	return &EthTransactionAttemptResolver{attmpt: attmpt}
}

func NewEthTransactionsAttempts(results []txmgr.EvmTxAttempt) []*EthTransactionAttemptResolver {
	var resolver []*EthTransactionAttemptResolver

	for _, tx := range results {
		resolver = append(resolver, NewEthTransactionAttempt(tx))
	}

	return resolver
}

func (r *EthTransactionAttemptResolver) GasPrice() string {
	return r.attmpt.TxFee.Legacy.ToInt().String()
}

func (r *EthTransactionAttemptResolver) Hash() string {
	return r.attmpt.Hash.String()
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
	results []txmgr.EvmTxAttempt
	total   int32
}

func NewEthTransactionsAttemptsPayload(results []txmgr.EvmTxAttempt, total int32) *EthTransactionsAttemptsPayloadResolver {
	return &EthTransactionsAttemptsPayloadResolver{results: results, total: total}
}

func (r *EthTransactionsAttemptsPayloadResolver) Results() []*EthTransactionAttemptResolver {
	return NewEthTransactionsAttempts(r.results)
}

func (r *EthTransactionsAttemptsPayloadResolver) Metadata() *PaginationMetadataResolver {
	return NewPaginationMetadata(r.total)
}

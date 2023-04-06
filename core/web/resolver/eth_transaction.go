package resolver

import (
	"context"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/graph-gophers/graphql-go"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/utils/stringutils"
	"github.com/smartcontractkit/chainlink/v2/core/web/loader"
)

type EthTransactionResolver struct {
	tx txmgr.EvmTx
}

func NewEthTransaction(tx txmgr.EvmTx) *EthTransactionResolver {
	return &EthTransactionResolver{tx: tx}
}

func NewEthTransactions(results []txmgr.EvmTx) []*EthTransactionResolver {
	var resolver []*EthTransactionResolver

	for _, tx := range results {
		resolver = append(resolver, NewEthTransaction(tx))
	}

	return resolver
}

func (r *EthTransactionResolver) State() string {
	return string(r.tx.State)
}

func (r *EthTransactionResolver) Data() hexutil.Bytes {
	return hexutil.Bytes(r.tx.EncodedPayload)
}

func (r *EthTransactionResolver) From() string {
	return r.tx.FromAddress.String()
}

func (r *EthTransactionResolver) To() string {
	return r.tx.ToAddress.String()
}

func (r *EthTransactionResolver) GasLimit() string {
	return stringutils.FromInt64(int64(r.tx.GasLimit))
}

func (r *EthTransactionResolver) GasPrice(ctx context.Context) string {
	attempts, err := r.Attempts(ctx)
	if err != nil || len(attempts) == 0 {
		return ""
	}

	return attempts[0].GasPrice()
}

func (r *EthTransactionResolver) Value() string {
	return r.tx.Value.String()
}

func (r *EthTransactionResolver) EVMChainID() graphql.ID {
	return graphql.ID(r.tx.EVMChainID.String())
}

func (r *EthTransactionResolver) Nonce() *string {
	if r.tx.Nonce == nil {
		return nil
	}

	value := stringutils.FromInt64(*r.tx.Nonce)

	return &value
}

func (r *EthTransactionResolver) Hash(ctx context.Context) string {
	attempts, err := r.Attempts(ctx)
	if err != nil || len(attempts) == 0 {
		return ""
	}

	return attempts[0].Hash()
}

func (r *EthTransactionResolver) Hex(ctx context.Context) string {
	attempts, err := r.Attempts(ctx)
	if err != nil || len(attempts) == 0 {
		return ""
	}

	return attempts[0].Hex()
}

// Chain resolves the node's chain object field.
func (r *EthTransactionResolver) Chain(ctx context.Context) (*ChainResolver, error) {
	chain, err := loader.GetChainByID(ctx, string(r.EVMChainID()))
	if err != nil {
		return nil, err
	}

	return NewChain(*chain), nil
}

func (r *EthTransactionResolver) Attempts(ctx context.Context) ([]*EthTransactionAttemptResolver, error) {
	id := stringutils.FromInt64(r.tx.ID)
	attempts, err := loader.GetEthTxAttemptsByEthTxID(ctx, id)
	if err != nil {
		return nil, err
	}

	return NewEthTransactionsAttempts(attempts), nil
}

func (r *EthTransactionResolver) SentAt(ctx context.Context) *string {
	attempts, err := r.Attempts(ctx)
	if err != nil || len(attempts) == 0 {
		return nil
	}

	return attempts[0].SentAt()
}

// -- EthTransaction Query --

type EthTransactionPayloadResolver struct {
	tx *txmgr.EvmTx
	NotFoundErrorUnionType
}

func NewEthTransactionPayload(tx *txmgr.EvmTx, err error) *EthTransactionPayloadResolver {
	e := NotFoundErrorUnionType{err: err, message: "transaction not found", isExpectedErrorFn: nil}

	return &EthTransactionPayloadResolver{tx: tx, NotFoundErrorUnionType: e}
}

func (r *EthTransactionPayloadResolver) ToEthTransaction() (*EthTransactionResolver, bool) {
	if r.err != nil {
		return nil, false
	}

	return NewEthTransaction(*r.tx), true
}

// -- EthTransactions Query --

type EthTransactionsPayloadResolver struct {
	results []txmgr.EvmTx
	total   int32
}

func NewEthTransactionsPayload(results []txmgr.EvmTx, total int32) *EthTransactionsPayloadResolver {
	return &EthTransactionsPayloadResolver{results: results, total: total}
}

func (r *EthTransactionsPayloadResolver) Results() []*EthTransactionResolver {
	return NewEthTransactions(r.results)
}

func (r *EthTransactionsPayloadResolver) Metadata() *PaginationMetadataResolver {
	return NewPaginationMetadata(r.total)
}

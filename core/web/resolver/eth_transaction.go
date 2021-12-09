package resolver

import (
	"context"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/graph-gophers/graphql-go"

	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	"github.com/smartcontractkit/chainlink/core/utils/stringutils"
	"github.com/smartcontractkit/chainlink/core/web/loader"
)

type EthTransactionResolver struct {
	tx bulletprooftxmanager.EthTx
}

func NewEthTransaction(tx bulletprooftxmanager.EthTx) *EthTransactionResolver {
	return &EthTransactionResolver{tx: tx}
}

func NewEthTransactions(results []bulletprooftxmanager.EthTx) []*EthTransactionResolver {
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

func (r *EthTransactionResolver) GasPrice() string {
	return r.tx.EthTxAttempts[0].GasPrice.String()
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

func (r *EthTransactionResolver) Hash() string {
	return r.tx.EthTxAttempts[0].Hash.Hex()
}

func (r *EthTransactionResolver) Hex() string {
	return hexutil.Encode(r.tx.EthTxAttempts[0].SignedRawTx)
}

// Chain resolves the node's chain object field.
func (r *EthTransactionResolver) Chain(ctx context.Context) (*ChainResolver, error) {
	chain, err := loader.GetChainByID(ctx, string(r.EVMChainID()))
	if err != nil {
		return nil, err
	}

	return NewChain(*chain), nil
}

func (r *EthTransactionResolver) Attempts() []*EthTransactionAttemptResolver {
	return NewEthTransactionsAttempts(r.tx.EthTxAttempts)
}

func (r *EthTransactionResolver) SentAt() *string {
	if r.tx.EthTxAttempts[0].BroadcastBeforeBlockNum == nil {
		return nil
	}

	value := stringutils.FromInt64(*r.tx.EthTxAttempts[0].BroadcastBeforeBlockNum)

	return &value
}

// -- EthTransaction Query --

type EthTransactionPayloadResolver struct {
	tx *bulletprooftxmanager.EthTx
	NotFoundErrorUnionType
}

func NewEthTransactionPayload(tx *bulletprooftxmanager.EthTx, err error) *EthTransactionPayloadResolver {
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
	results []bulletprooftxmanager.EthTx
	total   int32
}

func NewEthTransactionsPayload(results []bulletprooftxmanager.EthTx, total int32) *EthTransactionsPayloadResolver {
	return &EthTransactionsPayloadResolver{results: results, total: total}
}

func (r *EthTransactionsPayloadResolver) Results() []*EthTransactionResolver {
	return NewEthTransactions(r.results)
}

func (r *EthTransactionsPayloadResolver) Metadata() *PaginationMetadataResolver {
	return NewPaginationMetadata(r.total)
}

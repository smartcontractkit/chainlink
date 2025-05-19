package resolver

import (
	"context"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/graph-gophers/graphql-go"

	commonTypes "github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink-evm/pkg/assets"
	"github.com/smartcontractkit/chainlink-evm/pkg/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
	"github.com/smartcontractkit/chainlink/v2/core/utils/stringutils"
	"github.com/smartcontractkit/chainlink/v2/core/web/loader"
)

type EthTransactionResolver struct {
	tx txmgr.Tx
}

func NewEthTransaction(tx txmgr.Tx) *EthTransactionResolver {
	return &EthTransactionResolver{tx: tx}
}

func NewEthTransactions(results []txmgr.Tx) []*EthTransactionResolver {
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
	return stringutils.FromInt64(int64(r.tx.FeeLimit))
}

func (r *EthTransactionResolver) GasPrice(ctx context.Context) string {
	attempts, err := r.Attempts(ctx)
	if err != nil || len(attempts) == 0 {
		return ""
	}

	return attempts[0].GasPrice()
}

func (r *EthTransactionResolver) Value() string {
	v := assets.Eth(r.tx.Value)
	return v.String()
}

func (r *EthTransactionResolver) EVMChainID() graphql.ID {
	return graphql.ID(r.tx.ChainID.String())
}

func (r *EthTransactionResolver) Nonce() *string {
	if r.tx.Sequence == nil {
		return nil
	}

	value := r.tx.Sequence.String()

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
	relayID := commonTypes.NewRelayID(relay.NetworkEVM, string(r.EVMChainID()))
	chain, err := loader.GetChainByRelayID(ctx, relayID.Name())
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
	tx *txmgr.Tx
	NotFoundErrorUnionType
}

func NewEthTransactionPayload(tx *txmgr.Tx, err error) *EthTransactionPayloadResolver {
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
	results []txmgr.Tx
	total   int32
}

func NewEthTransactionsPayload(results []txmgr.Tx, total int32) *EthTransactionsPayloadResolver {
	return &EthTransactionsPayloadResolver{results: results, total: total}
}

func (r *EthTransactionsPayloadResolver) Results() []*EthTransactionResolver {
	return NewEthTransactions(r.results)
}

func (r *EthTransactionsPayloadResolver) Metadata() *PaginationMetadataResolver {
	return NewPaginationMetadata(r.total)
}

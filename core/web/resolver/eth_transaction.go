package resolver

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/graph-gophers/graphql-go"

	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/utils/stringutils"
	"github.com/smartcontractkit/chainlink/core/web/loader"
)

type EthTransactionExtraData struct {
	hash                    common.Hash
	gasPrice                *utils.Big
	signedRawTx             []byte
	broadcastBeforeBlockNum *int64
}

type EthTransactionData struct {
	tx    bulletprooftxmanager.EthTx
	extra EthTransactionExtraData
}

func NewEthTransactionData(tx bulletprooftxmanager.EthTx, attmpt bulletprooftxmanager.EthTxAttempt) EthTransactionData {
	return EthTransactionData{
		tx: tx,
		extra: EthTransactionExtraData{
			hash:                    attmpt.Hash,
			gasPrice:                attmpt.GasPrice,
			signedRawTx:             attmpt.SignedRawTx,
			broadcastBeforeBlockNum: attmpt.BroadcastBeforeBlockNum,
		},
	}
}

type EthTransactionResolver struct {
	data EthTransactionData
}

func NewEthTransaction(data EthTransactionData) *EthTransactionResolver {
	return &EthTransactionResolver{data: data}
}

func NewEthTransactions(results []EthTransactionData) []*EthTransactionResolver {
	var resolver []*EthTransactionResolver

	for _, tx := range results {
		resolver = append(resolver, NewEthTransaction(tx))
	}

	return resolver
}

func (r *EthTransactionResolver) State() string {
	return string(r.data.tx.State)
}

func (r *EthTransactionResolver) Data() hexutil.Bytes {
	return hexutil.Bytes(r.data.tx.EncodedPayload)
}

func (r *EthTransactionResolver) From() string {
	return r.data.tx.FromAddress.String()
}

func (r *EthTransactionResolver) To() string {
	return r.data.tx.ToAddress.String()
}

func (r *EthTransactionResolver) GasLimit() string {
	return stringutils.FromInt64(int64(r.data.tx.GasLimit))
}

func (r *EthTransactionResolver) GasPrice() string {
	return r.data.extra.gasPrice.String()
}

func (r *EthTransactionResolver) Value() string {
	return r.data.tx.Value.String()
}

func (r *EthTransactionResolver) EVMChainID() graphql.ID {
	return graphql.ID(r.data.tx.EVMChainID.String())
}

func (r *EthTransactionResolver) Nonce() *string {
	if r.data.tx.Nonce == nil {
		return nil
	}

	value := stringutils.FromInt64(*r.data.tx.Nonce)

	return &value
}

func (r *EthTransactionResolver) Hash() string {
	return r.data.extra.hash.String()
}

func (r *EthTransactionResolver) Hex() string {
	return hexutil.Encode(r.data.extra.signedRawTx)
}

// Chain resolves the node's chain object field.
func (r *EthTransactionResolver) Chain(ctx context.Context) (*ChainResolver, error) {
	chain, err := loader.GetChainByID(ctx, r.data.tx.EVMChainID.String())
	if err != nil {
		return nil, err
	}

	return NewChain(*chain), nil
}

func (r *EthTransactionResolver) SentAt() *string {
	if r.data.extra.broadcastBeforeBlockNum == nil {
		return nil
	}

	value := stringutils.FromInt64(*r.data.extra.broadcastBeforeBlockNum)

	return &value
}

// -- EthTransaction Query --

type EthTransactionPayloadResolver struct {
	data *EthTransactionData
	NotFoundErrorUnionType
}

func NewEthTransactionPayload(data *EthTransactionData, err error) *EthTransactionPayloadResolver {
	e := NotFoundErrorUnionType{err: err, message: "transaction not found", isExpectedErrorFn: nil}

	return &EthTransactionPayloadResolver{data: data, NotFoundErrorUnionType: e}
}

func (r *EthTransactionPayloadResolver) ToEthTransaction() (*EthTransactionResolver, bool) {
	if r.err != nil {
		return nil, false
	}

	return NewEthTransaction(*r.data), true
}

// -- EthTransactions Query --

type EthTransactionsPayloadResolver struct {
	results []EthTransactionData
	total   int32
}

func NewEthTransactionsPayload(results []EthTransactionData, total int32) *EthTransactionsPayloadResolver {
	return &EthTransactionsPayloadResolver{results: results, total: total}
}

func (r *EthTransactionsPayloadResolver) Results() []*EthTransactionResolver {
	return NewEthTransactions(r.results)
}

func (r *EthTransactionsPayloadResolver) Metadata() *PaginationMetadataResolver {
	return NewPaginationMetadata(r.total)
}

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

type EthTransactionResolver struct {
	tx    *bulletprooftxmanager.EthTx
	extra *EthTransactionExtraData
}

func NewEthTransaction(tx *bulletprooftxmanager.EthTx, extra *EthTransactionExtraData) *EthTransactionResolver {
	return &EthTransactionResolver{tx: tx, extra: extra}
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
	return r.extra.gasPrice.String()
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
	return r.extra.hash.String()
}

func (r *EthTransactionResolver) Hex() string {
	return hexutil.Encode(r.extra.signedRawTx)
}

// Chain resolves the node's chain object field.
func (r *EthTransactionResolver) Chain(ctx context.Context) (*ChainResolver, error) {
	chain, err := loader.GetChainByID(ctx, r.tx.EVMChainID.String())
	if err != nil {
		return nil, err
	}

	return NewChain(*chain), nil
}

func (r *EthTransactionResolver) SentAt() *string {
	if r.extra.broadcastBeforeBlockNum == nil {
		return nil
	}

	value := stringutils.FromInt64(*r.extra.broadcastBeforeBlockNum)

	return &value
}

// -- EthTransaction Query --

type EthTransactionPayloadResolver struct {
	tx    *bulletprooftxmanager.EthTx
	extra *EthTransactionExtraData
	NotFoundErrorUnionType
}

func NewEthTransactionPayload(tx *bulletprooftxmanager.EthTx, extra *EthTransactionExtraData, err error) *EthTransactionPayloadResolver {
	e := NotFoundErrorUnionType{err: err, message: "transaction not found", isExpectedErrorFn: nil}

	return &EthTransactionPayloadResolver{tx: tx, extra: extra, NotFoundErrorUnionType: e}
}

func (r *EthTransactionPayloadResolver) ToEthTransaction() (*EthTransactionResolver, bool) {
	if r.err != nil {
		return nil, false
	}

	return NewEthTransaction(r.tx, r.extra), true
}

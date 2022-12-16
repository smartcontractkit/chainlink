package resolver

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/graph-gophers/graphql-go"

	"github.com/smartcontractkit/chainlink/core/chains/evm"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/web/loader"
)

type ETHKey struct {
	state ethkey.State
	addr  ethkey.EIP55Address
	chain evm.Chain
}

type ETHKeyResolver struct {
	key ETHKey
}

func NewETHKey(key ETHKey) *ETHKeyResolver {
	return &ETHKeyResolver{key: key}
}

func NewETHKeys(keys []ETHKey) []*ETHKeyResolver {
	var resolvers []*ETHKeyResolver

	for _, k := range keys {
		resolvers = append(resolvers, NewETHKey(k))
	}

	return resolvers
}

func (r *ETHKeyResolver) Chain(ctx context.Context) (*ChainResolver, error) {
	chain, err := loader.GetChainByID(ctx, r.key.state.EVMChainID.String())
	if err != nil {
		return nil, err
	}

	return NewChain(*chain), nil
}

func (r *ETHKeyResolver) Address() string {
	return r.key.addr.Hex()
}

func (r *ETHKeyResolver) IsDisabled() bool {
	return r.key.state.Disabled
}

// ETHBalance returns the ETH balance available
func (r *ETHKeyResolver) ETHBalance(ctx context.Context) *string {
	if r.key.chain == nil {
		return nil
	}

	balanceMonitor := r.key.chain.BalanceMonitor()

	if balanceMonitor == nil {
		return nil
	}

	balance := balanceMonitor.GetEthBalance(r.key.state.Address.Address())

	if balance != nil {
		val := balance.String()
		return &val
	}

	return nil
}

func (r *ETHKeyResolver) LINKBalance(ctx context.Context) *string {
	if r.key.chain == nil {
		return nil
	}

	client := r.key.chain.Client()
	addr := common.HexToAddress(r.key.chain.Config().LinkContractAddress())
	balance, err := client.GetLINKBalance(ctx, addr, r.key.state.Address.Address())
	if err != nil {
		return nil
	}

	if balance != nil {
		val := balance.String()
		return &val
	}

	return nil
}

func (r *ETHKeyResolver) MaxGasPriceWei() *string {
	if r.key.chain == nil {
		return nil
	}

	gasPrice := r.key.chain.Config().KeySpecificMaxGasPriceWei(r.key.addr.Address())

	if gasPrice != nil {
		val := gasPrice.ToInt().String()
		return &val
	}

	return nil
}

func (r *ETHKeyResolver) CreatedAt() graphql.Time {
	return graphql.Time{Time: r.key.state.CreatedAt}
}

func (r *ETHKeyResolver) UpdatedAt() graphql.Time {
	return graphql.Time{Time: r.key.state.UpdatedAt}
}

// -- EthKeys query --

type ETHKeysPayloadResolver struct {
	keys []ETHKey
}

func NewETHKeysPayload(keys []ETHKey) *ETHKeysPayloadResolver {
	return &ETHKeysPayloadResolver{keys: keys}
}

func (r *ETHKeysPayloadResolver) Results() []*ETHKeyResolver {
	return NewETHKeys(r.keys)
}

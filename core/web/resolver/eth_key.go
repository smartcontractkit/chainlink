package resolver

import (
	"context"
	"math/big"

	"github.com/graph-gophers/graphql-go"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/web/loader"
)

type ETHKey struct {
	state     ethkey.State
	addr      ethkey.EIP55Address
	ethBal    *assets.Eth
	linkBal   *assets.Link
	maxGPrice *big.Int
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

func (r *ETHKeyResolver) IsFunding() bool {
	return r.key.state.IsFunding
}

func (r *ETHKeyResolver) ETHBalance() string {
	if r.key.ethBal != nil {
		return r.key.ethBal.String()
	}

	return "0"
}

func (r *ETHKeyResolver) LINKBalance() string {
	if r.key.linkBal != nil {
		return r.key.linkBal.String()
	}

	return "0"
}

func (r *ETHKeyResolver) MaxGasPriceWei() string {
	if r.key.maxGPrice != nil {
		return r.key.maxGPrice.String()
	}

	return "0"
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

func (r *ETHKeysPayloadResolver) Keys() []*ETHKeyResolver {
	return NewETHKeys(r.keys)
}

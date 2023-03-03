package gateway

import (
	"context"
	"math/big"

	"github.com/dontpanicdao/caigo/types"
)

type GatewayProvider struct {
	Gateway
}

func NewProvider(opts ...Option) *GatewayProvider {
	return &GatewayProvider{
		*NewClient(opts...),
	}
}

func (p *GatewayProvider) BlockByHash(ctx context.Context, hash, scope string) (*Block, error) {
	b, err := p.Block(ctx, &BlockOptions{BlockHash: hash})
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (p *GatewayProvider) BlockByNumber(ctx context.Context, number *big.Int, scope string) (*Block, error) {
	b, err := p.Block(ctx, &BlockOptions{BlockNumber: number.Uint64()})
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (p *GatewayProvider) TransactionByHash(ctx context.Context, hash string) (*Transaction, error) {
	t, err := p.Transaction(ctx, TransactionOptions{TransactionHash: hash})
	if err != nil {
		return nil, err
	}

	return &t.Transaction, nil
}

func (p *GatewayProvider) Class(ctx context.Context, classHash string) (*types.ContractClass, error) {
	panic("not implemented")
}

func (p *GatewayProvider) ClassHashAt(ctx context.Context, contractAddress string) (*types.Felt, error) {
	panic("not implemented")
}

func (p *GatewayProvider) ClassAt(ctx context.Context, contractAddress string) (*types.ContractClass, error) {
	panic("not implemented")
}

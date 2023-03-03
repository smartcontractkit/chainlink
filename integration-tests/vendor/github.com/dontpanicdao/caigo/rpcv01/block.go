package rpcv01

import (
	"context"

	ctypes "github.com/dontpanicdao/caigo/types"
)

// BlockNumber gets the most recent accepted block number.
func (provider *Provider) BlockNumber(ctx context.Context) (uint64, error) {
	var blockNumber uint64
	if err := provider.c.CallContext(ctx, &blockNumber, "starknet_blockNumber"); err != nil {
		return 0, err
	}
	return blockNumber, nil
}

// BlockHashAndNumber gets block information given the block number or its hash.
func (provider *Provider) BlockHashAndNumber(ctx context.Context) (*BlockHashAndNumberOutput, error) {
	var block BlockHashAndNumberOutput
	if err := do(ctx, provider.c, "starknet_blockHashAndNumber", &block); err != nil {
		return nil, err
	}
	return &block, nil
}

func WithBlockNumber(n uint64) BlockID {
	return BlockID{
		Number: &n,
	}
}

func WithBlockHash(h ctypes.Hash) BlockID {
	return BlockID{
		Hash: &h,
	}
}

func WithBlockTag(tag string) BlockID {
	return BlockID{
		Tag: tag,
	}
}

// BlockWithTxHashes gets block information given the block id.
func (provider *Provider) BlockWithTxHashes(ctx context.Context, blockID BlockID) (Block, error) {
	var result Block
	if err := do(ctx, provider.c, "starknet_getBlockWithTxHashes", &result, blockID); err != nil {
		return Block{}, err
	}
	return result, nil
}

// BlockTransactionCount gets the number of transactions in a block
func (provider *Provider) BlockTransactionCount(ctx context.Context, blockID BlockID) (uint64, error) {
	var result uint64
	if err := do(ctx, provider.c, "starknet_getBlockTransactionCount", &result, blockID); err != nil {
		return 0, err
	}
	return result, nil
}

// BlockWithTxs get block information with full transactions given the block id.
func (provider *Provider) BlockWithTxs(ctx context.Context, blockID BlockID) (interface{}, error) {
	var result Block
	if err := do(ctx, provider.c, "starknet_getBlockWithTxs", &result, blockID); err != nil {
		return nil, err
	}
	return &result, nil
}

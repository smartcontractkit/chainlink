package rpc

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/NethermindEth/juno/core/felt"
)

// BlockNumber returns the block number of the current block.
//
// Parameters:
// - ctx: The context to use for the request
// Returns:
// - uint64: The block number
// - error: An error if any
func (provider *Provider) BlockNumber(ctx context.Context) (uint64, error) {
	var blockNumber uint64
	if err := provider.c.CallContext(ctx, &blockNumber, "starknet_blockNumber"); err != nil {
		if errors.Is(err, errNotFound) {
			return 0, ErrNoBlocks
		}
		return 0, Err(InternalError, err)
	}
	return blockNumber, nil
}

// BlockHashAndNumber retrieves the hash and number of the current block.
//
// Parameters:
// - ctx: The context to use for the request.
// Returns:
// - *BlockHashAndNumberOutput: The hash and number of the current block
// - error: An error if any
func (provider *Provider) BlockHashAndNumber(ctx context.Context) (*BlockHashAndNumberOutput, error) {
	var block BlockHashAndNumberOutput
	if err := do(ctx, provider.c, "starknet_blockHashAndNumber", &block); err != nil {
		return nil, tryUnwrapToRPCErr(err, ErrNoBlocks)
	}
	return &block, nil
}

// WithBlockNumber returns a BlockID with the given block number.
//
// Parameters:
//   - n: The block number to use for the BlockID.
//
// Returns:
//   - BlockID: A BlockID struct with the specified block number
func WithBlockNumber(n uint64) BlockID {
	return BlockID{
		Number: &n,
	}
}

// WithBlockHash returns a BlockID with the given hash.
//
// Parameters:
// - h: The hash to use for the BlockID.
// Returns:
// - BlockID: A BlockID struct with the specified hash
func WithBlockHash(h *felt.Felt) BlockID {
	return BlockID{
		Hash: h,
	}
}

// WithBlockTag creates a new BlockID with the specified tag.
//
// Parameters:
// - tag: The tag for the BlockID
// Returns:
// - BlockID: A BlockID struct with the specified tag
func WithBlockTag(tag string) BlockID {
	return BlockID{
		Tag: tag,
	}
}

// BlockWithTxHashes retrieves the block with transaction hashes for the given block ID.
//
// Parameters:
// - ctx: The context.Context object for controlling the function call
// - blockID: The ID of the block to retrieve the transactions from
// Returns:
// - interface{}: The retrieved block
// - error: An error, if any
func (provider *Provider) BlockWithTxHashes(ctx context.Context, blockID BlockID) (interface{}, error) {
	var result BlockTxHashes
	if err := do(ctx, provider.c, "starknet_getBlockWithTxHashes", &result, blockID); err != nil {
		return nil, tryUnwrapToRPCErr(err, ErrBlockNotFound)
	}

	// if header.Hash == nil it's a pending block
	if result.BlockHeader.BlockHash == nil {
		return &PendingBlockTxHashes{
			PendingBlockHeader{
				ParentHash:       result.ParentHash,
				Timestamp:        result.Timestamp,
				SequencerAddress: result.SequencerAddress},
			result.Transactions,
		}, nil
	}

	return &result, nil
}

// StateUpdate is a function that performs a state update operation
// (gets the information about the result of executing the requested block).
//
// Parameters:
// - ctx: The context.Context object for controlling the function call
// - blockID: The ID of the block to retrieve the transactions from
// Returns:
// - *StateUpdateOutput: The retrieved state update
// - error: An error, if any
func (provider *Provider) StateUpdate(ctx context.Context, blockID BlockID) (*StateUpdateOutput, error) {
	var state StateUpdateOutput
	if err := do(ctx, provider.c, "starknet_getStateUpdate", &state, blockID); err != nil {
		return nil, tryUnwrapToRPCErr(err, ErrBlockNotFound)
	}
	return &state, nil
}

// BlockTransactionCount returns the number of transactions in a specific block.
//
// Parameters:
// - ctx: The context.Context object to handle cancellation signals and timeouts
// - blockID: The ID of the block to retrieve the number of transactions from
// Returns:
// - uint64: The number of transactions in the block
// - error: An error, if any
func (provider *Provider) BlockTransactionCount(ctx context.Context, blockID BlockID) (uint64, error) {
	var result uint64
	if err := do(ctx, provider.c, "starknet_getBlockTransactionCount", &result, blockID); err != nil {
		if errors.Is(err, errNotFound) {
			return 0, ErrBlockNotFound
		}
		return 0, Err(InternalError, err)
	}
	return result, nil
}

// BlockWithTxs retrieves a block with its transactions given the block id.
//
// Parameters:
// - ctx: The context.Context object for the request
// - blockID: The ID of the block to retrieve
// Returns:
// - interface{}: The retrieved block
// - error: An error, if any
func (provider *Provider) BlockWithTxs(ctx context.Context, blockID BlockID) (interface{}, error) {
	var result Block
	if err := do(ctx, provider.c, "starknet_getBlockWithTxs", &result, blockID); err != nil {
		return nil, tryUnwrapToRPCErr(err, ErrBlockNotFound)
	}
	// if header.Hash == nil it's a pending block
	if result.BlockHeader.BlockHash == nil {
		return &PendingBlock{
			PendingBlockHeader{
				ParentHash:       result.ParentHash,
				Timestamp:        result.Timestamp,
				SequencerAddress: result.SequencerAddress},
			result.Transactions,
		}, nil
	}
	return &result, nil
}

// Get block information with full transactions and receipts given the block id
func (provider *Provider) BlockWithReceipts(ctx context.Context, blockID BlockID) (interface{}, error) {
	var result json.RawMessage
	if err := do(ctx, provider.c, "starknet_getBlockWithReceipts", &result, blockID); err != nil {
		return nil, tryUnwrapToRPCErr(err, ErrBlockNotFound)
	}

	var m map[string]interface{}
	if err := json.Unmarshal(result, &m); err != nil {
		return nil, Err(InternalError, err.Error())
	}

	// PendingBlockWithReceipts doesn't contain a "status" field
	if _, ok := m["status"]; ok {
		var block BlockWithReceipts
		if err := json.Unmarshal(result, &block); err != nil {
			return nil, Err(InternalError, err.Error())
		}
		return &block, nil
	} else {
		var pendingBlock PendingBlockWithReceipts
		if err := json.Unmarshal(result, &pendingBlock); err != nil {
			return nil, Err(InternalError, err.Error())
		}
		return &pendingBlock, nil
	}

}

// Copyright 2021 github.com/gagliardetto
// This file has been modified by github.com/gagliardetto
//
// Copyright 2020 dfuse Platform Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rpc

import (
	"context"
	"fmt"

	"github.com/gagliardetto/solana-go"
)

// GetConfirmedBlock returns identity and transaction information about a confirmed block in the ledger.
//
// DEPRECATED: Please use `getBlock` instead.
// This method is expected to be removed in solana-core v1.8
func (cl *Client) GetConfirmedBlock(
	ctx context.Context,
	slot uint64,
) (out *GetConfirmedBlockResult, err error) {
	return cl.GetConfirmedBlockWithOpts(
		ctx,
		slot,
		nil,
	)
}

type GetConfirmedBlockOpts struct {
	Encoding solana.EncodingType

	// Level of transaction detail to return.
	TransactionDetails TransactionDetailsType

	// Whether to populate the rewards array. If parameter not provided, the default includes rewards.
	Rewards *bool

	// Desired commitment; "processed" is not supported.
	// If parameter not provided, the default is "finalized".
	Commitment CommitmentType
}

// GetConfirmedBlock returns identity and transaction information about a confirmed block in the ledger.
//
// DEPRECATED: Please use `getBlock` instead.
// This method is expected to be removed in solana-core v1.8
func (cl *Client) GetConfirmedBlockWithOpts(
	ctx context.Context,
	slot uint64,
	opts *GetConfirmedBlockOpts,
) (out *GetConfirmedBlockResult, err error) {

	params := []interface{}{slot}
	if opts != nil {
		obj := M{}
		if opts.Encoding != "" {
			obj["encoding"] = opts.Encoding
		}
		if opts.TransactionDetails != "" {
			obj["transactionDetails"] = opts.TransactionDetails
		}
		if opts.Rewards != nil {
			obj["rewards"] = opts.Rewards
		}
		if opts.Commitment != "" {
			obj["commitment"] = opts.Commitment
		}
		if len(obj) != 0 {
			params = append(params, obj)
		}
	}

	err = cl.rpcClient.CallForInto(ctx, &out, "getConfirmedBlock", params)
	return
}

// GetConfirmedBlocks returns a list of confirmed blocks between two slots.
//
// The result field will be an array of u64 integers listing confirmed blocks between
// start_slot and either end_slot, if provided, or latest confirmed block, inclusive.
// Max range allowed is 500,000 slots.
//
// DEPRECATED: Please use `getBlocks` instead.
// This method is expected to be removed in solana-core v1.8
func (cl *Client) GetConfirmedBlocks(
	ctx context.Context,
	startSlot uint64,
	endSlot *uint64,
	commitment CommitmentType,
) (out []uint64, err error) {

	params := []interface{}{startSlot}
	if endSlot != nil {
		params = append(params, endSlot)
	}
	if commitment != "" {
		params = append(params, M{"commitment": string(commitment)})
	}

	err = cl.rpcClient.CallForInto(ctx, &out, "getConfirmedBlocks", params)
	return
}

// GetConfirmedBlocksWithLimit returns a list of confirmed blocks starting at the given slot.
//
// DEPRECATED: Please use `getBlocksWithLimit` instead.
// This method is expected to be removed in solana-core v1.8
func (cl *Client) GetConfirmedBlocksWithLimit(
	ctx context.Context,
	startSlot uint64,
	limit uint64,
	commitment CommitmentType,
) (out []uint64, err error) {

	params := []interface{}{startSlot, limit}
	if commitment != "" {
		params = append(params, M{"commitment": string(commitment)})
	}

	err = cl.rpcClient.CallForInto(ctx, &out, "getConfirmedBlocksWithLimit", params)
	return
}

// GetConfirmedSignaturesForAddress2 returns confirmed signatures for transactions involving an
// address backwards in time from the provided signature or most recent confirmed block.
//
// DEPRECATED: Please use getSignaturesForAddress instead.
// This method is expected to be removed in solana-core v1.8
func (cl *Client) GetConfirmedSignaturesForAddress2(
	ctx context.Context,
	address solana.PublicKey,
	opts *GetConfirmedSignaturesForAddress2Opts,
) (out GetConfirmedSignaturesForAddress2Result, err error) {

	params := []interface{}{address}

	if opts != nil {
		obj := M{}
		if opts.Limit != nil {
			obj["limit"] = opts.Limit
		}
		if !opts.Before.IsZero() {
			obj["before"] = opts.Before
		}
		if !opts.Until.IsZero() {
			obj["until"] = opts.Until
		}
		if opts.Commitment != "" {
			obj["commitment"] = opts.Commitment
		}
		if len(obj) > 0 {
			params = append(params, obj)
		}
	}

	err = cl.rpcClient.CallForInto(ctx, &out, "getConfirmedSignaturesForAddress2", params)
	return
}

// GetConfirmedTransaction returns transaction details for a confirmed transaction.
func (cl *Client) GetConfirmedTransaction(
	ctx context.Context,
	signature solana.Signature,
) (out *TransactionWithMeta, err error) {
	params := []interface{}{signature, "json"}

	err = cl.rpcClient.CallForInto(ctx, &out, "getConfirmedTransaction", params)
	if err != nil {
		return nil, err
	}
	if out == nil {
		return nil, ErrNotFound
	}
	return
}

// GetConfirmedTransactionWithOpts returns transaction details for a confirmed transaction.
func (cl *Client) GetConfirmedTransactionWithOpts(
	ctx context.Context,
	signature solana.Signature,
	opts *GetTransactionOpts,
) (out *TransactionWithMeta, err error) {
	params := []interface{}{signature}
	if opts != nil {
		obj := M{}
		if opts.Encoding != "" {
			if !solana.IsAnyOfEncodingType(
				opts.Encoding,
				// Valid encodings:
				// solana.EncodingJSON, // TODO
				// solana.EncodingJSONParsed, // TODO
				solana.EncodingBase58,
				solana.EncodingBase64,
				solana.EncodingBase64Zstd,
			) {
				return nil, fmt.Errorf("provided encoding is not supported: %s", opts.Encoding)
			}
			obj["encoding"] = opts.Encoding
		}
		if opts.Commitment != "" {
			obj["commitment"] = opts.Commitment
		}
		if len(obj) > 0 {
			params = append(params, obj)
		}
	}
	err = cl.rpcClient.CallForInto(ctx, &out, "getConfirmedTransaction", params)
	if err != nil {
		return nil, err
	}
	if out == nil {
		return nil, ErrNotFound
	}
	return
}

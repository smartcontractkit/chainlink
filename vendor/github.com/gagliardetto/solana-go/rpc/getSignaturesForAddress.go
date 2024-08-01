// Copyright 2021 github.com/gagliardetto
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

	"github.com/gagliardetto/solana-go"
)

type GetSignaturesForAddressOpts struct {
	// (optional) Maximum transaction signatures to return (between 1 and 1,000, default: 1,000).
	Limit *int `json:"limit,omitempty"`

	// (optional) Start searching backwards from this transaction signature.
	// If not provided the search starts from the top of the highest max confirmed block.
	Before solana.Signature `json:"before,omitempty"`

	// (optional) Search until this transaction signature, if found before limit reached.
	Until solana.Signature `json:"until,omitempty"`

	// (optional) Commitment; "processed" is not supported.
	// If parameter not provided, the default is "finalized".
	Commitment CommitmentType `json:"commitment,omitempty"`

	// The minimum slot that the request can be evaluated at.
	// This parameter is optional.
	MinContextSlot *uint64
}

// GetSignaturesForAddress returns confirmed signatures for transactions
// involving an address backwards in time from the provided signature
// or most recent confirmed block.
//
// NEW: This method is only available in solana-core v1.7 or newer.
// Please use `getConfirmedSignaturesForAddress2` for solana-core v1.6
func (cl *Client) GetSignaturesForAddress(
	ctx context.Context,
	account solana.PublicKey,
) (out []*TransactionSignature, err error) {
	return cl.GetSignaturesForAddressWithOpts(
		ctx,
		account,
		nil,
	)
}

// GetSignaturesForAddressWithOpts returns confirmed signatures for transactions
// involving an address backwards in time from the provided signature
// or most recent confirmed block.
//
// NEW: This method is only available in solana-core v1.7 or newer.
// Please use `getConfirmedSignaturesForAddress2` for solana-core v1.6
func (cl *Client) GetSignaturesForAddressWithOpts(
	ctx context.Context,
	account solana.PublicKey,
	opts *GetSignaturesForAddressOpts,
) (out []*TransactionSignature, err error) {
	params := []interface{}{account}
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
		if opts.MinContextSlot != nil {
			obj["minContextSlot"] = *opts.MinContextSlot
		}
		if len(obj) > 0 {
			params = append(params, obj)
		}
	}

	err = cl.rpcClient.CallForInto(ctx, &out, "getSignaturesForAddress", params)
	return
}

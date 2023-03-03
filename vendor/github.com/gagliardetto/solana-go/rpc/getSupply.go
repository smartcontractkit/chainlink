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

// GetSupply returns information about the current supply.
func (cl *Client) GetSupply(ctx context.Context, commitment CommitmentType) (out *GetSupplyResult, err error) {
	return cl.GetSupplyWithOpts(ctx, &GetSupplyOpts{Commitment: commitment})
}

// GetSupply returns information about the current supply.
func (cl *Client) GetSupplyWithOpts(
	ctx context.Context,
	opts *GetSupplyOpts,
) (out *GetSupplyResult, err error) {
	obj := M{
		"commitment": CommitmentConfirmed,
	}
	if opts != nil {
		if opts.Commitment != "" {
			obj["commitment"] = opts.Commitment
		}
		obj["excludeNonCirculatingAccountsList"] = opts.ExcludeNonCirculatingAccountsList
	}

	err = cl.rpcClient.CallForInto(ctx, &out, "getSupply", []interface{}{obj})
	return
}

type GetSupplyOpts struct {
	Commitment CommitmentType `json:"commitment,omitempty"`

	ExcludeNonCirculatingAccountsList bool `json:"excludeNonCirculatingAccountsList,omitempty"` // exclude non circulating accounts list from response
}

type GetSupplyResult struct {
	RPCContext
	Value *SupplyResult `json:"value"`
}

type SupplyResult struct {
	// Total supply in lamports
	Total uint64 `json:"total"`

	// Circulating supply in lamports.
	Circulating uint64 `json:"circulating"`

	// Non-circulating supply in lamports.
	NonCirculating uint64 `json:"nonCirculating"`

	// An array of account addresses of non-circulating accounts.
	// If `excludeNonCirculatingAccountsList` is enabled, the returned array will be empty.
	NonCirculatingAccounts []solana.PublicKey `json:"nonCirculatingAccounts"`
}

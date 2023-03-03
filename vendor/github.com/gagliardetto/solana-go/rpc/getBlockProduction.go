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

type GetBlockProductionResult struct {
	RPCContext
	Value BlockProductionResult `json:"value"`
}

type GetBlockProductionOpts struct {
	//
	// This parameter is optional.
	Commitment CommitmentType `json:"commitment,omitempty"`

	// Slot range to return block production for.
	// If parameter not provided, defaults to current epoch.
	//
	// This parameter is optional.
	Range *SlotRangeRequest `json:"range,omitempty"`

	// Only return results for this validator identity.
	//
	// This parameter is optional.
	Identity *solana.PublicKey `json:"identity,omitempty"`
}

type SlotRangeRequest struct {
	// First slot to return block production information for (inclusive)
	FirstSlot uint64 `json:"firstSlot"`

	// Last slot to return block production information for (inclusive).
	// If parameter not provided, defaults to the highest slot
	//
	// This parameter is optional.
	LastSlot *uint64 `json:"lastSlot,omitempty"`
}

// GetBlockProduction returns recent block production information from the current or previous epoch.
func (cl *Client) GetBlockProduction(
	ctx context.Context,
) (out *GetBlockProductionResult, err error) {
	return cl.GetBlockProductionWithOpts(
		ctx,
		nil,
	)
}

// GetBlockProduction returns recent block production information from the current or previous epoch.
func (cl *Client) GetBlockProductionWithOpts(
	ctx context.Context,
	opts *GetBlockProductionOpts,
) (out *GetBlockProductionResult, err error) {
	params := []interface{}{}

	if opts != nil {
		obj := M{}
		if opts.Commitment != "" {
			obj["commitment"] = opts.Commitment
		}
		if opts.Range != nil {
			rngObj := M{}
			rngObj["firstSlot"] = opts.Range.FirstSlot
			if opts.Range.LastSlot != nil {
				rngObj["lastSlot"] = opts.Range.LastSlot
			}
			obj["range"] = rngObj
		}
		if opts.Identity != nil {
			obj["identity"] = opts.Identity
		}
		if len(obj) != 0 {
			params = append(params, obj)
		}
	}
	err = cl.rpcClient.CallForInto(ctx, &out, "getBlockProduction", params)

	return
}

type BlockProductionResult struct {
	ByIdentity IdentityToSlotsBlocks `json:"byIdentity"`

	Range SlotRangeResponse `json:"range"`
}

// A dictionary of validator identities.
// Value is a two element array containing the number
// of leader slots and the number of blocks produced.
type IdentityToSlotsBlocks map[solana.PublicKey][2]int64

type SlotRangeResponse struct {
	// First slot of the block production information (inclusive)
	FirstSlot uint64 `json:"firstSlot"`

	// Last slot of block production information (inclusive)
	LastSlot uint64 `json:"lastSlot"`
}

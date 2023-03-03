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

type GetInflationRewardOpts struct {
	Commitment CommitmentType

	// An epoch for which the reward occurs.
	// If omitted, the previous epoch will be used.
	Epoch *uint64
}

// GetInflationReward returns the inflation / staking reward for a list of addresses for an epoch.
func (cl *Client) GetInflationReward(
	ctx context.Context,

	// An array of addresses to query.
	addresses []solana.PublicKey,

	opts *GetInflationRewardOpts,

) (out []*GetInflationRewardResult, err error) {
	params := []interface{}{addresses}
	if opts != nil {
		obj := M{}
		if opts.Commitment != "" {
			obj["commitment"] = opts.Commitment
		}
		if opts.Epoch != nil {
			obj["epoch"] = opts.Epoch
		}
		if len(obj) > 0 {
			params = append(params, obj)
		}
	}
	// TODO: check
	err = cl.rpcClient.CallForInto(ctx, &out, "getInflationReward", params)
	return
}

type GetInflationRewardResult struct {
	// Epoch for which reward occured.
	Epoch uint64 `json:"epoch"`

	// The slot in which the rewards are effective.
	EffectiveSlot uint64 `json:"effectiveSlot"`

	// Reward amount in lamports.
	Amount uint64 `json:"amount"`

	// Post balance of the account in lamports.
	PostBalance uint64 `json:"postBalance"`

	// Vote account commission when the reward was credited.
	Commission *uint8 `json:"commission,omitempty"`
}

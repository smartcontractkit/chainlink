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

// GetStakeActivation returns epoch activation information for a stake account.
func (cl *Client) GetStakeActivation(
	ctx context.Context,
	// Pubkey of stake account to query
	account solana.PublicKey,

	commitment CommitmentType,

	// epoch for which to calculate activation details.
	// If parameter not provided, defaults to current epoch.
	epoch *uint64,
) (out *GetStakeActivationResult, err error) {
	params := []interface{}{account}
	{
		obj := M{}
		if commitment != "" {
			obj["commitment"] = commitment
		}
		if epoch != nil {
			obj["epoch"] = epoch
		}
		if len(obj) > 0 {
			params = append(params, obj)
		}
	}
	err = cl.rpcClient.CallForInto(ctx, &out, "getStakeActivation", params)
	return
}

type GetStakeActivationResult struct {
	// The stake account's activation state, one of: active, inactive, activating, deactivating.
	State ActivationStateType `json:"state"`

	// Stake active during the epoch.
	Active uint64 `json:"active"`

	// Stake inactive during the epoch.
	Inactive uint64 `json:"inactive"`
}

type ActivationStateType string

const (
	ActivationStateActive       ActivationStateType = "active"
	ActivationStateInactive     ActivationStateType = "inactive"
	ActivationStateActivating   ActivationStateType = "activating"
	ActivationStateDeactivating ActivationStateType = "deactivating"
)

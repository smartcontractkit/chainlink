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
)

// GetEpochInfo returns information about the current epoch.
func (cl *Client) GetEpochInfo(
	ctx context.Context,
	commitment CommitmentType, // optional
) (out *GetEpochInfoResult, err error) {
	params := []interface{}{}
	if commitment != "" {
		params = append(params, M{"commitment": commitment})
	}
	err = cl.rpcClient.CallForInto(ctx, &out, "getEpochInfo", params)
	return
}

type GetEpochInfoResult struct {
	// The current slot.
	AbsoluteSlot uint64 `json:"absoluteSlot"`

	// The current block height.
	BlockHeight uint64 `json:"blockHeight"`

	// The current epoch.
	Epoch uint64 `json:"epoch"`

	// The current slot relative to the start of the current epoch.
	SlotIndex uint64 `json:"slotIndex"`

	// The number of slots in this epoch.
	SlotsInEpoch uint64 `json:"slotsInEpoch"`

	TransactionCount *uint64 `json:"transactionCount,omitempty"`
}

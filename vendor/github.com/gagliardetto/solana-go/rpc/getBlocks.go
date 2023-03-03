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

// GetBlocks returns a list of confirmed blocks between two slots.
// The result will be an array of u64 integers listing confirmed blocks
// between start_slot and either end_slot, if provided, or latest
// confirmed block, inclusive. Max range allowed is 500,000 slots.
func (cl *Client) GetBlocks(
	ctx context.Context,
	startSlot uint64,
	endSlot *uint64, // optional
	commitment CommitmentType, // optional
) (out BlocksResult, err error) {
	params := []interface{}{startSlot}
	if endSlot != nil {
		params = append(params, endSlot)
	}
	if commitment != "" {
		params = append(params,
			// TODO: provide commitment as string instead of object?
			M{"commitment": commitment},
		)
	}
	err = cl.rpcClient.CallForInto(ctx, &out, "getBlocks", params)

	return
}

type BlocksResult []uint64

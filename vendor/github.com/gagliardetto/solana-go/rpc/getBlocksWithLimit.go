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

// GetBlocksWithLimit returns a list of confirmed blocks starting at the given slot.
// The result field will be an array of u64 integers listing
// confirmed blocks starting at startSlot for up to limit blocks, inclusive.
func (cl *Client) GetBlocksWithLimit(
	ctx context.Context,
	startSlot uint64,
	limit uint64,
	commitment CommitmentType, // optional; "processed" is not supported. If parameter not provided, the default is "finalized".
) (out *BlocksResult, err error) {
	params := []interface{}{startSlot, limit}
	if commitment != "" {
		params = append(params,
			// TODO: provide commitment as string instead of object?
			M{"commitment": commitment},
		)
	}
	err = cl.rpcClient.CallForInto(ctx, &out, "getBlocksWithLimit", params)
	return
}

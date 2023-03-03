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

// Returns the highest slot information that the node has snapshots for.
// This will find the highest full snapshot slot, and the highest incremental
// snapshot slot _based on_ the full snapshot slot, if there is one.
//
// **NEW: This method is only available in solana-core v1.9 or newer. Please use
// `getSnapshotSlot` for solana-core v1.8**
func (cl *Client) GetHighestSnapshotSlot(ctx context.Context) (out *GetHighestSnapshotSlotResult, err error) {
	err = cl.rpcClient.CallForInto(ctx, &out, "getHighestSnapshotSlot", nil)
	return
}

type GetHighestSnapshotSlotResult struct {
	Full        uint64  `json:"full"`                  // Highest full snapshot slot.
	Incremental *uint64 `json:"incremental,omitempty"` // Highest incremental snapshot slot based on full.
}

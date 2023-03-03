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

// GetBlockTime returns the estimated production time of a block.
//
// Each validator reports their UTC time to the ledger on a regular
// interval by intermittently adding a timestamp to a Vote for a
// particular block. A requested block's time is calculated from
// the stake-weighted mean of the Vote timestamps in a set of
// recent blocks recorded on the ledger.
//
// The result will be an int64 estimated production time,
// as Unix timestamp (seconds since the Unix epoch),
// or nil if the timestamp is not available for this block.
func (cl *Client) GetBlockTime(
	ctx context.Context,
	block uint64, // block, identified by Slot
) (out *solana.UnixTimeSeconds, err error) {
	params := []interface{}{block}
	err = cl.rpcClient.CallForInto(ctx, &out, "getBlockTime", params)
	return
}

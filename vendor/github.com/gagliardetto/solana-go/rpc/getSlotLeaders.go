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

// GetSlotLeaders returns the slot leaders for a given slot range.
func (cl *Client) GetSlotLeaders(
	ctx context.Context,
	start uint64,
	limit uint64,
) (out []solana.PublicKey, err error) {
	params := []interface{}{start, limit}
	err = cl.rpcClient.CallForInto(ctx, &out, "getSlotLeaders", params)
	return
}

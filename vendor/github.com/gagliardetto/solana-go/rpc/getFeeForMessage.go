// Copyright 2022 github.com/gagliardetto
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

// Get the fee the network will charge for a particular Message.
//
// **NEW**: This method is only available in solana-core v1.9 or newer. Please use
// `getFees` for solana-core v1.8.
func (cl *Client) GetFeeForMessage(
	ctx context.Context,
	message string, // Base-64 encoded Message
	commitment CommitmentType, // optional
) (out *GetFeeForMessageResult, err error) {
	params := []interface{}{message}
	if commitment != "" {
		params = append(params, M{"commitment": commitment})
	}
	err = cl.rpcClient.CallForInto(ctx, &out, "getFeeForMessage", params)
	return
}

type GetFeeForMessageResult struct {
	RPCContext

	// Fee corresponding to the message at the specified blockhash.
	Value *uint64 `json:"value"`
}

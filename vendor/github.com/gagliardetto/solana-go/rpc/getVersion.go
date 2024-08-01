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

// GetVersion returns the current solana versions running on the node.
func (cl *Client) GetVersion(ctx context.Context) (out *GetVersionResult, err error) {
	err = cl.rpcClient.CallForInto(ctx, &out, "getVersion", nil)
	return
}

type GetVersionResult struct {
	// Software version of solana-core.
	SolanaCore string `json:"solana-core"`

	// Unique identifier of the current software's feature set.
	FeatureSet int64 `json:"feature-set"`
}

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

// GetHealth returns the current health of the node.
// If one or more --trusted-validator arguments are provided
// to solana-validator, "ok" is returned when the node has within
// HEALTH_CHECK_SLOT_DISTANCE slots of the highest trusted validator,
// otherwise an error is returned. "ok" is always returned if no
// trusted validators are provided.
//
// - If the node is healthy: "ok"
// - If the node is unhealthy, a JSON RPC error response is returned.
//   The specifics of the error response are UNSTABLE and may change in the future.
func (cl *Client) GetHealth(ctx context.Context) (out string, err error) {
	err = cl.rpcClient.CallForInto(ctx, &out, "getHealth", nil)
	return
}

const HealthOk = "ok"

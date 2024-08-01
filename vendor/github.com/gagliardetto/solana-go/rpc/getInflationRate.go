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

// GetInflationRate returns the specific inflation values for the current epoch.
func (cl *Client) GetInflationRate(ctx context.Context) (out *GetInflationRateResult, err error) {
	err = cl.rpcClient.CallForInto(ctx, &out, "getInflationRate", nil)
	return
}

type GetInflationRateResult struct {
	// Total inflation.
	Total float64 `json:"total"`

	// Inflation allocated to validators.
	Validator float64 `json:"validator"`

	// Inflation allocated to the foundation.
	Foundation float64 `json:"foundation"`

	// Epoch for which these values are valid.
	Epoch float64 `json:"epoch"`
}

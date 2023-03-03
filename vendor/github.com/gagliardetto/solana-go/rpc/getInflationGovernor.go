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

// GetInflationGovernor returns the current inflation governor.
func (cl *Client) GetInflationGovernor(
	ctx context.Context,
	commitment CommitmentType, // optional
) (out *GetInflationGovernorResult, err error) {
	params := []interface{}{}
	if commitment != "" {
		params = append(params,
			M{"commitment": commitment},
		)
	}
	err = cl.rpcClient.CallForInto(ctx, &out, "getInflationGovernor", params)
	return
}

type GetInflationGovernorResult struct {
	// The initial inflation percentage from time 0.
	Initial float64 `json:"initial"`

	// Terminal inflation percentage.
	Terminal float64 `json:"terminal"`

	// Rate per year at which inflation is lowered. Rate reduction is derived using the target slot time in genesis config.
	Taper float64 `json:"taper"`

	// Percentage of total inflation allocated to the foundation.
	Foundation float64 `json:"foundation"`

	// Duration of foundation pool inflation in years.
	FoundationTerm float64 `json:"foundationTerm"`
}

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

// GetFeeRateGovernor returns the fee rate governor information from the root bank.
func (cl *Client) GetFeeRateGovernor(ctx context.Context) (out *GetFeeRateGovernorResult, err error) {
	err = cl.rpcClient.CallForInto(ctx, &out, "getFeeRateGovernor", nil)
	return
}

type GetFeeRateGovernorResult struct {
	RPCContext
	Value FeeRateGovernorResult `json:"value"`
}
type FeeRateGovernorResult struct {
	FeeRateGovernor FeeRateGovernor `json:"feeRateGovernor"`
}
type FeeRateGovernor struct {
	// Percentage of fees collected to be destroyed.
	BurnPercent uint8 `json:"burnPercent"`

	// Largest value lamportsPerSignature can attain for the next slot.
	MaxLamportsPerSignature uint64 `json:"maxLamportsPerSignature"`

	// Smallest value lamportsPerSignature can attain for the next slot.
	MinLamportsPerSignature uint64 `json:"minLamportsPerSignature"`

	// Desired fee rate for the cluster.
	TargetLamportsPerSignature uint64 `json:"targetLamportsPerSignature"`

	// Desired signature rate for the cluster.
	TargetSignaturesPerSlot uint64 `json:"targetSignaturesPerSlot"`
}

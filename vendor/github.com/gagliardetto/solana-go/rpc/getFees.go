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

// GetFees returns a recent block hash from the ledger,
// a fee schedule that can be used to compute the cost
// of submitting a transaction using it, and the last
// slot in which the blockhash will be valid.
func (cl *Client) GetFees(
	ctx context.Context,
	commitment CommitmentType, // optional
) (out *GetFeesResult, err error) {
	params := []interface{}{}
	if commitment != "" {
		params = append(params, M{"commitment": commitment})
	}
	err = cl.rpcClient.CallForInto(ctx, &out, "getFees", params)
	return
}

type GetFeesResult struct {
	RPCContext
	Value *FeesResult `json:"value"`
}

type FeesResult struct {
	// A Hash.
	Blockhash solana.Hash `json:"blockhash"`

	// FeeCalculator object, the fee schedule for this block hash.
	FeeCalculator FeeCalculator `json:"feeCalculator"`

	// DEPRECATED - this value is inaccurate and should not be relied upon.
	LastValidSlot uint64 `json:"lastValidSlot"`

	// Last block height at which a blockhash will be valid.
	LastValidBlockHeight uint64 `json:"lastValidBlockHeight"`
}

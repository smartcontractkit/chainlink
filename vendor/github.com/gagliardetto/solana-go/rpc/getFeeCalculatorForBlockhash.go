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

// GetFeeCalculatorForBlockhash returns the fee calculator
// associated with the query blockhash, or null if the blockhash has expired.
//
// NOTE: DEPRECATED
func (cl *Client) GetFeeCalculatorForBlockhash(
	ctx context.Context,
	hash solana.Hash, // query blockhash
	commitment CommitmentType, // optional
) (out *GetFeeCalculatorForBlockhashResult, err error) {
	params := []interface{}{hash}
	if commitment != "" {
		params = append(params, M{"commitment": commitment})
	}
	err = cl.rpcClient.CallForInto(ctx, &out, "getFeeCalculatorForBlockhash", params)
	return
}

type GetFeeCalculatorForBlockhashResult struct {
	RPCContext

	// Value will be nil if the query blockhash has expired.
	Value *FeeCalculatorForBlockhashResult `json:"value"`
}

type FeeCalculatorForBlockhashResult struct {
	// Object describing the cluster fee rate at the queried blockhash
	FeeCalculator FeeCalculator `json:"feeCalculator"`
}

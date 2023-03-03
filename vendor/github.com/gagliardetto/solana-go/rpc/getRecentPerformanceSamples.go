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

// GetRecentPerformanceSamples returns a list of recent performance samples,
// in reverse slot order. Performance samples are taken every 60 seconds
// and include the number of transactions and slots that occur in a given time window.
func (cl *Client) GetRecentPerformanceSamples(
	ctx context.Context,
	limit *uint,
) (out []*GetRecentPerformanceSamplesResult, err error) {
	params := []interface{}{}
	if limit != nil {
		params = append(params, limit)
	}
	err = cl.rpcClient.CallForInto(ctx, &out, "getRecentPerformanceSamples", params)
	return
}

type GetRecentPerformanceSamplesResult struct {
	// Slot in which sample was taken at.
	Slot uint64 `json:"slot"`

	// Number of transactions in sample.
	NumTransactions uint64 `json:"numTransactions"`

	// Number of slots in sample.
	NumSlots uint64 `json:"numSlots"`

	// Number of seconds in a sample window.
	SamplePeriodSecs uint16 `json:"samplePeriodSecs"`
}

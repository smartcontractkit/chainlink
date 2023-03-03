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

// GetEpochSchedule returns epoch schedule information from this cluster's genesis config.
func (cl *Client) GetEpochSchedule(ctx context.Context) (out *GetEpochScheduleResult, err error) {
	err = cl.rpcClient.CallForInto(ctx, &out, "getEpochSchedule", nil)
	return
}

type GetEpochScheduleResult struct {
	// The maximum number of slots in each epoch.
	SlotsPerEpoch uint64 `json:"slotsPerEpoch"`

	// The number of slots before beginning of an epoch to calculate a leader schedule for that epoch.
	LeaderScheduleSlotOffset uint64 `json:"leaderScheduleSlotOffset"`

	// Whether epochs start short and grow.
	Warmup bool `json:"warmup"`

	// First normal-length epoch, log2(slotsPerEpoch) - log2(MINIMUM_SLOTS_PER_EPOCH)
	FirstNormalEpoch uint64 `json:"firstNormalEpoch"`

	// MINIMUM_SLOTS_PER_EPOCH * (2.pow(firstNormalEpoch) - 1)
	FirstNormalSlot uint64 `json:"firstNormalSlot"`
}

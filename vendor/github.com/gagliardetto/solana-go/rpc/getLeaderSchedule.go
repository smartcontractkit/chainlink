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

// GetLeaderSchedule returns the leader schedule for current epoch.
func (cl *Client) GetLeaderSchedule(
	ctx context.Context,
) (out GetLeaderScheduleResult, err error) {
	return cl.GetLeaderScheduleWithOpts(
		ctx,
		nil,
	)
}

type GetLeaderScheduleOpts struct {
	Commitment CommitmentType

	// Fetch the leader schedule for the epoch that corresponds
	// to the provided slot.
	// If unspecified, the leader schedule for the current epoch is fetched
	Epoch *uint64

	// TODO: is identity a pubkey?
	Identity *solana.PublicKey // Only return results for this validator identity
}

// GetLeaderScheduleWithOpts returns the leader schedule for an epoch.
func (cl *Client) GetLeaderScheduleWithOpts(
	ctx context.Context,
	opts *GetLeaderScheduleOpts,
) (out GetLeaderScheduleResult, err error) {
	params := []interface{}{}
	if opts != nil {
		if opts.Epoch != nil {
			params = append(params, opts.Epoch)
		}
		obj := M{}
		if opts.Commitment != "" {
			obj["commitment"] = opts.Commitment
		}
		if opts.Identity != nil {
			obj["identity"] = opts.Identity
		}
		if len(obj) > 0 {
			params = append(params, obj)
		}
	}
	err = cl.rpcClient.CallForInto(ctx, &out, "getLeaderSchedule", params)
	if err != nil {
		return nil, err
	}
	// TODO: check that this behaviour is implemented everywhere:
	if out == nil {
		return nil, ErrNotFound
	}
	return
}

// The result field will be a dictionary of validator identities,
// and their corresponding leader slot indices as values
// (indices are relative to the first slot in the requested epoch).
type GetLeaderScheduleResult map[solana.PublicKey][]uint64

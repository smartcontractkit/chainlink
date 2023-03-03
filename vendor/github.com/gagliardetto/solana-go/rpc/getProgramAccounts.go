// Copyright 2021 github.com/gagliardetto
// This file has been modified by github.com/gagliardetto
//
// Copyright 2020 dfuse Platform Inc.
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

// GetProgramAccounts returns all accounts owned by the provided program publicKey.
func (cl *Client) GetProgramAccounts(
	ctx context.Context,
	publicKey solana.PublicKey,
) (out GetProgramAccountsResult, err error) {
	return cl.GetProgramAccountsWithOpts(
		ctx,
		publicKey,
		nil,
	)
}

// GetProgramAccountsWithOpts returns all accounts owned by the provided program publicKey.
func (cl *Client) GetProgramAccountsWithOpts(
	ctx context.Context,
	publicKey solana.PublicKey,
	opts *GetProgramAccountsOpts,
) (out GetProgramAccountsResult, err error) {
	obj := M{
		"encoding": "base64",
	}
	if opts != nil {
		if opts.Commitment != "" {
			obj["commitment"] = string(opts.Commitment)
		}
		if len(opts.Filters) != 0 {
			obj["filters"] = opts.Filters
		}
		if opts.Encoding != "" {
			obj["encoding"] = opts.Encoding
		}
		if opts.DataSlice != nil {
			obj["dataSlice"] = M{
				"offset": opts.DataSlice.Offset,
				"length": opts.DataSlice.Length,
			}
		}
	}

	params := []interface{}{publicKey, obj}

	err = cl.rpcClient.CallForInto(ctx, &out, "getProgramAccounts", params)
	return
}

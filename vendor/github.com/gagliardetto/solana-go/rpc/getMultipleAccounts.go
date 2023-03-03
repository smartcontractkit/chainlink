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
	"errors"

	"github.com/gagliardetto/solana-go"
)

type GetMultipleAccountsResult struct {
	RPCContext
	Value []*Account `json:"value"`
}

// GetMultipleAccounts returns the account information for a list of Pubkeys.
func (cl *Client) GetMultipleAccounts(
	ctx context.Context,
	accounts ...solana.PublicKey, // An array of Pubkeys to query
) (out *GetMultipleAccountsResult, err error) {
	return cl.GetMultipleAccountsWithOpts(
		ctx,
		accounts,
		nil,
	)
}

type GetMultipleAccountsOpts GetAccountInfoOpts

// GetMultipleAccountsWithOpts returns the account information for a list of Pubkeys.
func (cl *Client) GetMultipleAccountsWithOpts(
	ctx context.Context,
	accounts []solana.PublicKey,
	opts *GetMultipleAccountsOpts,
) (out *GetMultipleAccountsResult, err error) {
	params := []interface{}{accounts}

	if opts != nil {
		obj := M{}
		if opts.Encoding != "" {
			obj["encoding"] = opts.Encoding
		}
		if opts.Commitment != "" {
			obj["commitment"] = opts.Commitment
		}
		if opts.DataSlice != nil {
			obj["dataSlice"] = M{
				"offset": opts.DataSlice.Offset,
				"length": opts.DataSlice.Length,
			}
			if opts.Encoding == solana.EncodingJSONParsed {
				return nil, errors.New("cannot use dataSlice with EncodingJSONParsed")
			}
		}
		if len(obj) > 0 {
			params = append(params, obj)
		}
	}

	err = cl.rpcClient.CallForInto(ctx, &out, "getMultipleAccounts", params)
	if err != nil {
		return nil, err
	}
	if out.Value == nil {
		return nil, ErrNotFound
	}
	return
}

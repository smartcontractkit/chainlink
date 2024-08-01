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

// GetTokenAccountsByOwner returns all SPL Token accounts by token owner.
func (cl *Client) GetTokenAccountsByOwner(
	ctx context.Context,
	owner solana.PublicKey,
	conf *GetTokenAccountsConfig,
	opts *GetTokenAccountsOpts,
) (out *GetTokenAccountsResult, err error) {
	params := []interface{}{owner}
	if conf == nil {
		return nil, errors.New("conf is nil")
	}
	if conf.Mint != nil && conf.ProgramId != nil {
		return nil, errors.New("conf.Mint and conf.ProgramId are both set; must be just one of them")
	}

	{
		confObj := M{}
		if conf.Mint != nil {
			confObj["mint"] = conf.Mint
		}
		if conf.ProgramId != nil {
			confObj["programId"] = conf.ProgramId
		}
		if len(confObj) > 0 {
			params = append(params, confObj)
		}
	}
	defaultEncoding := solana.EncodingBase64
	{
		optsObj := M{}
		if opts != nil {
			if opts.Commitment != "" {
				optsObj["commitment"] = opts.Commitment
			}
			if opts.Encoding != "" {
				optsObj["encoding"] = opts.Encoding
			} else {
				optsObj["encoding"] = defaultEncoding
			}
			if opts.DataSlice != nil {
				optsObj["dataSlice"] = M{
					"offset": opts.DataSlice.Offset,
					"length": opts.DataSlice.Length,
				}
				if opts.Encoding == solana.EncodingJSONParsed {
					return nil, errors.New("cannot use dataSlice with EncodingJSONParsed")
				}
			}
			if len(optsObj) > 0 {
				params = append(params, optsObj)
			}
		} else {
			params = append(params, M{"encoding": defaultEncoding})
		}
	}

	err = cl.rpcClient.CallForInto(ctx, &out, "getTokenAccountsByOwner", params)
	return
}

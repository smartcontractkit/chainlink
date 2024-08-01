// Copyright 2021 github.com/gagliardetto
// This file has been modified by github.com/gagliardetto
//
// Copyright 2020 dfuse Platform Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
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

	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
)

// GetAccountInfo returns all information associated with the account of provided publicKey.
func (cl *Client) GetAccountInfo(ctx context.Context, account solana.PublicKey) (out *GetAccountInfoResult, err error) {
	return cl.GetAccountInfoWithOpts(
		ctx,
		account,
		&GetAccountInfoOpts{
			Commitment: "",
			DataSlice:  nil,
		},
	)
}

// GetAccountDataInto decodes the binary data and populates
// the provided `inVar` parameter with all data associated with the account of provided publicKey.
func (cl *Client) GetAccountDataInto(ctx context.Context, account solana.PublicKey, inVar interface{}) (err error) {
	resp, err := cl.GetAccountInfo(ctx, account)
	if err != nil {
		return err
	}
	return bin.NewBinDecoder(resp.Value.Data.GetBinary()).Decode(inVar)
}

// GetAccountDataBorshInto decodes the borsh binary data and populates
// the provided `inVar` parameter with all data associated with the account of provided publicKey.
func (cl *Client) GetAccountDataBorshInto(ctx context.Context, account solana.PublicKey, inVar interface{}) (err error) {
	resp, err := cl.GetAccountInfo(ctx, account)
	if err != nil {
		return err
	}
	return bin.NewBorshDecoder(resp.Value.Data.GetBinary()).Decode(inVar)
}

type GetAccountInfoOpts struct {
	// Encoding for Account data.
	// Either "base58" (slow), "base64", "base64+zstd", or "jsonParsed".
	// - "base58" is limited to Account data of less than 129 bytes.
	// - "base64" will return base64 encoded data for Account data of any size.
	// - "base64+zstd" compresses the Account data using Zstandard and base64-encodes the result.
	// - "jsonParsed" encoding attempts to use program-specific state parsers to return more
	// 	 human-readable and explicit account state data. If "jsonParsed" is requested but a parser
	//   cannot be found, the field falls back to "base64" encoding,
	//   detectable when the data field is type <string>.
	//
	// This parameter is optional.
	Encoding solana.EncodingType

	// Commitment requirement.
	//
	// This parameter is optional.
	Commitment CommitmentType

	// dataSlice parameters for limiting returned account data:
	// Limits the returned account data using the provided offset and length fields;
	// only available for "base58", "base64" or "base64+zstd" encodings.
	//
	// This parameter is optional.
	DataSlice *DataSlice

	// The minimum slot that the request can be evaluated at.
	// This parameter is optional.
	MinContextSlot *uint64
}

// GetAccountInfoWithOpts returns all information associated with the account of provided publicKey.
// You can specify the encoding of the returned data with the encoding parameter.
// You can limit the returned account data with the offset and length parameters.
func (cl *Client) GetAccountInfoWithOpts(
	ctx context.Context,
	account solana.PublicKey,
	opts *GetAccountInfoOpts,
) (*GetAccountInfoResult, error) {
	out, err := cl.getAccountInfoWithOpts(ctx, account, opts)
	if err != nil {
		return nil, err
	}
	if out.Value == nil {
		return nil, ErrNotFound
	}
	return out, nil
}

func (cl *Client) getAccountInfoWithOpts(
	ctx context.Context,
	account solana.PublicKey,
	opts *GetAccountInfoOpts,
) (out *GetAccountInfoResult, err error) {

	obj := M{
		// default encoding:
		"encoding": solana.EncodingBase64,
	}

	if opts != nil {
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
		if opts.MinContextSlot != nil {
			obj["minContextSlot"] = *opts.MinContextSlot
		}
	}

	params := []interface{}{account}
	if len(obj) > 0 {
		params = append(params, obj)
	}

	err = cl.rpcClient.CallForInto(ctx, &out, "getAccountInfo", params)
	if err != nil {
		return nil, err
	}
	if out == nil {
		return nil, errors.New("expected a value, got null result")
	}
	return out, nil
}

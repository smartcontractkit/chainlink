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

// SendEncodedTransaction submits a signed base64 encoded transaction to the cluster for processing.
// The only difference between this function and SignTransaction is that the latter takes a *solana.Transaction value, as the former takes a raw base64 string
func (cl *Client) SendEncodedTransaction(
	ctx context.Context,
	encodedTx string,
) (signature solana.Signature, err error) {
	opts := TransactionOpts{
		SkipPreflight:       false,
		PreflightCommitment: "",
	}

	return cl.SendEncodedTransactionWithOpts(
		ctx,
		encodedTx,
		opts,
	)
}

// SendEncodedTransactionWithOpts submits a signed base64 encoded transaction to the cluster for processing.
func (cl *Client) SendEncodedTransactionWithOpts(
	ctx context.Context,
	encodedTx string,
	opts TransactionOpts,
) (signature solana.Signature, err error) {
	obj := opts.ToMap()
	params := []interface{}{
		encodedTx,
		obj,
	}

	err = cl.rpcClient.CallForInto(ctx, &signature, "sendTransaction", params)
	return
}

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
	"encoding/base64"
	"fmt"

	"github.com/gagliardetto/solana-go"
)

type SimulateTransactionResponse struct {
	RPCContext
	Value *SimulateTransactionResult `json:"value"`
}

type SimulateTransactionResult struct {
	// Error if transaction failed, null if transaction succeeded.
	Err interface{} `json:"err,omitempty"`

	// Array of log messages the transaction instructions output during execution,
	// null if simulation failed before the transaction was able to execute
	// (for example due to an invalid blockhash or signature verification failure)
	Logs []string `json:"logs,omitempty"`

	// Array of accounts with the same length as the accounts.addresses array in the request.
	Accounts []*Account `json:"accounts"`

	// The number of compute budget units consumed during the processing of this transaction.
	UnitsConsumed *uint64 `json:"unitsConsumed,omitempty"`
}

// SimulateTransaction simulates sending a transaction.
func (cl *Client) SimulateTransaction(
	ctx context.Context,
	transaction *solana.Transaction,
) (out *SimulateTransactionResponse, err error) {
	return cl.SimulateTransactionWithOpts(
		ctx,
		transaction,
		nil,
	)
}

type SimulateTransactionOpts struct {
	// If true the transaction signatures will be verified
	// (default: false, conflicts with ReplaceRecentBlockhash)
	SigVerify bool

	// Commitment level to simulate the transaction at.
	// (default: "finalized").
	Commitment CommitmentType

	// If true the transaction recent blockhash will be replaced with the most recent blockhash.
	// (default: false, conflicts with SigVerify)
	ReplaceRecentBlockhash bool

	Accounts *SimulateTransactionAccountsOpts
}

type SimulateTransactionAccountsOpts struct {
	// (optional) Encoding for returned Account data,
	// either "base64" (default), "base64+zstd" or "jsonParsed".
	// - "jsonParsed" encoding attempts to use program-specific state parsers
	//   to return more human-readable and explicit account state data.
	//   If "jsonParsed" is requested but a parser cannot be found,
	//   the field falls back to binary encoding, detectable when
	//   the data field is type <string>.
	Encoding solana.EncodingType

	// An array of accounts to return.
	Addresses []solana.PublicKey
}

// SimulateTransaction simulates sending a transaction.
func (cl *Client) SimulateTransactionWithOpts(
	ctx context.Context,
	transaction *solana.Transaction,
	opts *SimulateTransactionOpts,
) (out *SimulateTransactionResponse, err error) {
	txData, err := transaction.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("send transaction: encode transaction: %w", err)
	}

	obj := M{
		"encoding": "base64",
	}
	if opts != nil {
		if opts.SigVerify {
			obj["sigVerify"] = opts.SigVerify
		}
		if opts.Commitment != "" {
			obj["commitment"] = opts.Commitment
		}
		if opts.ReplaceRecentBlockhash {
			obj["replaceRecentBlockhash"] = opts.ReplaceRecentBlockhash
		}
		if opts.Accounts != nil {
			obj["accounts"] = M{
				"encoding":  opts.Accounts.Encoding,
				"addresses": opts.Accounts.Addresses,
			}
		}
	}

	b64Data := base64.StdEncoding.EncodeToString(txData)
	params := []interface{}{
		b64Data,
		obj,
	}

	err = cl.rpcClient.CallForInto(ctx, &out, "simulateTransaction", params)
	return
}

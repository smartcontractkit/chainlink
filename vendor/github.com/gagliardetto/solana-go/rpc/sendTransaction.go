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

// SendTransaction submits a signed transaction to the cluster for processing.
func (cl *Client) SendTransaction(
	ctx context.Context,
	transaction *solana.Transaction,
) (signature solana.Signature, err error) {
	opts := TransactionOpts{
		SkipPreflight:       false,
		PreflightCommitment: "",
	}

	return cl.SendTransactionWithOpts(
		ctx,
		transaction,
		opts,
	)
}

// SendTransaction submits a signed transaction to the cluster for processing.
// This method does not alter the transaction in any way;
// it relays the transaction created by clients to the node as-is.
//
// If the node's rpc service receives the transaction,
// this method immediately succeeds, without waiting for any confirmations.
// A successful response from this method does not guarantee the transaction
// is processed or confirmed by the cluster.
//
// While the rpc service will reasonably retry to submit it, the transaction
// could be rejected if transaction's recent_blockhash expires before it lands.
//
// Use getSignatureStatuses to ensure a transaction is processed and confirmed.
//
// Before submitting, the following preflight checks are performed:
//
// 	- The transaction signatures are verified
//  - The transaction is simulated against the bank slot specified by the preflight
//    commitment. On failure an error will be returned. Preflight checks may be
//    disabled if desired. It is recommended to specify the same commitment and
//    preflight commitment to avoid confusing behavior.
//
// The returned signature is the first signature in the transaction, which is
// used to identify the transaction (transaction id). This identifier can be
// easily extracted from the transaction data before submission.
func (cl *Client) SendTransactionWithOpts(
	ctx context.Context,
	transaction *solana.Transaction,
	opts TransactionOpts,
) (signature solana.Signature, err error) {
	txData, err := transaction.MarshalBinary()
	if err != nil {
		return solana.Signature{}, fmt.Errorf("send transaction: encode transaction: %w", err)
	}

	return cl.SendEncodedTransactionWithOpts(
		ctx,
		base64.StdEncoding.EncodeToString(txData),
		opts,
	)
}

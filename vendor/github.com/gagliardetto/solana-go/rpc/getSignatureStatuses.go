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

// GetSignatureStatuses Returns the statuses of a list of signatures.
// Unless the searchTransactionHistory configuration parameter
// is included,this method only searches the recent status cache
// of signatures, which retains statuses for all active slots plus
// MAX_RECENT_BLOCKHASHES rooted slots.
func (cl *Client) GetSignatureStatuses(
	ctx context.Context,

	// If true, a Solana node will search its ledger
	// cache for any signatures not found in the recent status cache.
	searchTransactionHistory bool,

	// Transaction signatures to confirm.
	transactionSignatures ...solana.Signature,
) (out *GetSignatureStatusesResult, err error) {
	params := []interface{}{transactionSignatures}
	if searchTransactionHistory {
		params = append(params, M{"searchTransactionHistory": searchTransactionHistory})
	}
	err = cl.rpcClient.CallForInto(ctx, &out, "getSignatureStatuses", params)
	if err != nil {
		return nil, err
	}
	if out.Value == nil {
		// Unknown transaction
		return nil, ErrNotFound
	}

	return
}

type GetSignatureStatusesResult struct {
	RPCContext
	Value []*SignatureStatusesResult `json:"value"`
}

type SignatureStatusesResult struct {
	// The slot the transaction was processed.
	Slot uint64 `json:"slot"`

	// Number of blocks since signature confirmation,
	// null if rooted or finalized by a supermajority of the cluster.
	Confirmations *uint64 `json:"confirmations"`

	// Error if transaction failed, null if transaction succeeded.
	Err interface{} `json:"err"`

	// The transaction's cluster confirmation status; either processed, confirmed, or finalized.
	ConfirmationStatus ConfirmationStatusType `json:"confirmationStatus"`

	// DEPRECATED: Transaction status.
	Status DeprecatedTransactionMetaStatus `json:"status"`
}

type ConfirmationStatusType string

const (
	ConfirmationStatusProcessed ConfirmationStatusType = "processed"
	ConfirmationStatusConfirmed ConfirmationStatusType = "confirmed"
	ConfirmationStatusFinalized ConfirmationStatusType = "finalized"
)

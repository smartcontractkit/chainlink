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

type GetVoteAccountsOpts struct {
	Commitment CommitmentType `json:"commitment,omitempty"`

	// (optional) Only return results for this validator vote address.
	VotePubkey *solana.PublicKey `json:"votePubkey,omitempty"`

	// (optional) Do not filter out delinquent validators with no stake.
	KeepUnstakedDelinquents *bool `json:"keepUnstakedDelinquents,omitempty"`

	// (optional) Specify the number of slots behind the tip that
	// a validator must fall to be considered delinquent.
	// NOTE: For the sake of consistency between ecosystem products,
	// it is not recommended that this argument be specified.
	DelinquentSlotDistance *uint64 `json:"delinquentSlotDistance,omitempty"`
}

// GetVoteAccounts returns the account info and associated
// stake for all the voting accounts in the current bank.
func (cl *Client) GetVoteAccounts(
	ctx context.Context,
	opts *GetVoteAccountsOpts,
) (out *GetVoteAccountsResult, err error) {
	params := []interface{}{}
	if opts != nil {
		obj := M{}
		if opts.Commitment != "" {
			obj["commitment"] = string(opts.Commitment)
		}
		if opts.VotePubkey != nil {
			obj["votePubkey"] = opts.VotePubkey.String()
		}
		if opts.KeepUnstakedDelinquents != nil {
			obj["keepUnstakedDelinquents"] = opts.KeepUnstakedDelinquents
		}
		if opts.DelinquentSlotDistance != nil {
			obj["delinquentSlotDistance"] = opts.DelinquentSlotDistance
		}
		if len(obj) > 0 {
			params = append(params, obj)
		}
	}
	err = cl.rpcClient.CallForInto(ctx, &out, "getVoteAccounts", params)
	return
}

type GetVoteAccountsResult struct {
	Current    []VoteAccountsResult `json:"current"`
	Delinquent []VoteAccountsResult `json:"delinquent"`
}

type VoteAccountsResult struct {
	// Vote account address.
	VotePubkey solana.PublicKey `json:"votePubkey,omitempty"`

	// Validator identity.
	NodePubkey solana.PublicKey `json:"nodePubkey,omitempty"`

	// The stake, in lamports, delegated to this vote account and active in this epoch.
	ActivatedStake uint64 `json:"activatedStake,omitempty"`

	// Whether the vote account is staked for this epoch.
	EpochVoteAccount bool `json:"epochVoteAccount,omitempty"`

	// Percentage (0-100) of rewards payout owed to the vote account.
	Commission uint8 `json:"commission,omitempty"`

	// Most recent slot voted on by this vote account.
	LastVote uint64 `json:"lastVote,omitempty"`

	RootSlot uint64 `json:"rootSlot,omitempty"` //

	// History of how many credits earned by the end of each epoch,
	// as an array of arrays containing: [epoch, credits, previousCredits]
	EpochCredits [][]int64 `json:"epochCredits,omitempty"`
}

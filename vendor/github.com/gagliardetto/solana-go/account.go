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

package solana

import (
	"fmt"
)

// Wallet is a wrapper around a PrivateKey
type Wallet struct {
	PrivateKey PrivateKey
}

func NewWallet() *Wallet {
	privateKey, err := NewRandomPrivateKey()
	if err != nil {
		panic(fmt.Sprintf("failed to generate private key: %s", err))
	}
	return &Wallet{
		PrivateKey: privateKey,
	}
}

func WalletFromPrivateKeyBase58(privateKey string) (*Wallet, error) {
	k, err := PrivateKeyFromBase58(privateKey)
	if err != nil {
		return nil, fmt.Errorf("account from private key: private key from b58: %w", err)
	}
	return &Wallet{
		PrivateKey: k,
	}, nil
}

func (a *Wallet) PublicKey() PublicKey {
	return a.PrivateKey.PublicKey()
}

type AccountMeta struct {
	PublicKey  PublicKey
	IsWritable bool
	IsSigner   bool
}

// Meta intializes a new AccountMeta with the provided pubKey.
func Meta(
	pubKey PublicKey,
) *AccountMeta {
	return &AccountMeta{
		PublicKey: pubKey,
	}
}

// WRITE sets IsWritable to true.
func (meta *AccountMeta) WRITE() *AccountMeta {
	meta.IsWritable = true
	return meta
}

// SIGNER sets IsSigner to true.
func (meta *AccountMeta) SIGNER() *AccountMeta {
	meta.IsSigner = true
	return meta
}

func NewAccountMeta(
	pubKey PublicKey,
	WRITE bool,
	SIGNER bool,
) *AccountMeta {
	return &AccountMeta{
		PublicKey:  pubKey,
		IsWritable: WRITE,
		IsSigner:   SIGNER,
	}
}

func (a AccountMeta) less(act *AccountMeta) bool {
	if a.IsSigner != act.IsSigner {
		return a.IsSigner
	}
	if a.IsWritable != act.IsWritable {
		return a.IsWritable
	}
	return false
}

type AccountMetaSlice []*AccountMeta

func (slice *AccountMetaSlice) Append(account *AccountMeta) {
	*slice = append(*slice, account)
}

func (slice *AccountMetaSlice) SetAccounts(accounts []*AccountMeta) error {
	*slice = accounts
	return nil
}

func (slice AccountMetaSlice) GetAccounts() []*AccountMeta {
	out := make([]*AccountMeta, 0, len(slice))
	for i := range slice {
		if slice[i] != nil {
			out = append(out, slice[i])
		}
	}
	return out
}

// Get returns the AccountMeta at the desired index.
// If the index is not present, it returns nil.
func (slice AccountMetaSlice) Get(index int) *AccountMeta {
	if len(slice) > index {
		return slice[index]
	}
	return nil
}

// GetSigners returns the accounts that are signers.
func (slice AccountMetaSlice) GetSigners() []*AccountMeta {
	signers := make([]*AccountMeta, 0, len(slice))
	for _, ac := range slice {
		if ac.IsSigner {
			signers = append(signers, ac)
		}
	}
	return signers
}

// GetKeys returns the pubkeys of all AccountMeta.
func (slice AccountMetaSlice) GetKeys() PublicKeySlice {
	keys := make(PublicKeySlice, 0, len(slice))
	for _, ac := range slice {
		keys = append(keys, ac.PublicKey)
	}
	return keys
}

func (slice AccountMetaSlice) Len() int {
	return len(slice)
}

func (slice AccountMetaSlice) SplitFrom(index int) (AccountMetaSlice, AccountMetaSlice) {
	if index < 0 {
		panic("negative index")
	}
	if index == 0 {
		return AccountMetaSlice{}, slice
	}
	if index > len(slice)-1 {
		return slice, AccountMetaSlice{}
	}

	firstLen, secondLen := calcSplitAtLengths(len(slice), index)

	first := make(AccountMetaSlice, firstLen)
	copy(first, slice[:index])

	second := make(AccountMetaSlice, secondLen)
	copy(second, slice[index:])

	return first, second
}

func calcSplitAtLengths(total int, index int) (int, int) {
	if index == 0 {
		return 0, total
	}
	if index > total-1 {
		return total, 0
	}
	return index, total - index
}

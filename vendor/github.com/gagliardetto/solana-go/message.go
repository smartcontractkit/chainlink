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
	"encoding/base64"
	"fmt"

	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go/text"
	"github.com/gagliardetto/treeout"
)

type Message struct {
	// List of base-58 encoded public keys used by the transaction,
	// including by the instructions and for signatures.
	// The first `message.header.numRequiredSignatures` public keys must sign the transaction.
	AccountKeys []PublicKey `json:"accountKeys"`

	// Details the account types and signatures required by the transaction.
	Header MessageHeader `json:"header"`

	// A base-58 encoded hash of a recent block in the ledger used to
	// prevent transaction duplication and to give transactions lifetimes.
	RecentBlockhash Hash `json:"recentBlockhash"`

	// List of program instructions that will be executed in sequence
	// and committed in one atomic transaction if all succeed.
	Instructions []CompiledInstruction `json:"instructions"`
}

var _ bin.EncoderDecoder = &Message{}

func (mx *Message) EncodeToTree(txTree treeout.Branches) {
	txTree.Child(text.Sf("RecentBlockhash: %s", mx.RecentBlockhash))

	txTree.Child(fmt.Sprintf("AccountKeys[len=%v]", len(mx.AccountKeys))).ParentFunc(func(accountKeysBranch treeout.Branches) {
		for _, key := range mx.AccountKeys {
			accountKeysBranch.Child(text.ColorizeBG(key.String()))
		}
	})

	txTree.Child("Header").ParentFunc(func(message treeout.Branches) {
		mx.Header.EncodeToTree(message)
	})
}

func (header *MessageHeader) EncodeToTree(mxBranch treeout.Branches) {
	mxBranch.Child(text.Sf("NumRequiredSignatures: %v", header.NumRequiredSignatures))
	mxBranch.Child(text.Sf("NumReadonlySignedAccounts: %v", header.NumReadonlySignedAccounts))
	mxBranch.Child(text.Sf("NumReadonlyUnsignedAccounts: %v", header.NumReadonlyUnsignedAccounts))
}

func (mx *Message) MarshalBinary() ([]byte, error) {
	buf := []byte{
		mx.Header.NumRequiredSignatures,
		mx.Header.NumReadonlySignedAccounts,
		mx.Header.NumReadonlyUnsignedAccounts,
	}

	bin.EncodeCompactU16Length(&buf, len(mx.AccountKeys))
	for _, key := range mx.AccountKeys {
		buf = append(buf, key[:]...)
	}

	buf = append(buf, mx.RecentBlockhash[:]...)

	bin.EncodeCompactU16Length(&buf, len(mx.Instructions))
	for _, instruction := range mx.Instructions {
		buf = append(buf, byte(instruction.ProgramIDIndex))
		bin.EncodeCompactU16Length(&buf, len(instruction.Accounts))
		for _, accountIdx := range instruction.Accounts {
			buf = append(buf, byte(accountIdx))
		}

		bin.EncodeCompactU16Length(&buf, len(instruction.Data))
		buf = append(buf, instruction.Data...)
	}
	return buf, nil
}

func (mx Message) MarshalWithEncoder(encoder *bin.Encoder) error {
	out, err := mx.MarshalBinary()
	if err != nil {
		return err
	}
	return encoder.WriteBytes(out, false)
}

func (mx Message) ToBase64() string {
	out, _ := mx.MarshalBinary()
	return base64.StdEncoding.EncodeToString(out)
}

func (mx *Message) UnmarshalWithDecoder(decoder *bin.Decoder) (err error) {
	{
		mx.Header.NumRequiredSignatures, err = decoder.ReadUint8()
		if err != nil {
			return err
		}
		mx.Header.NumReadonlySignedAccounts, err = decoder.ReadUint8()
		if err != nil {
			return err
		}
		mx.Header.NumReadonlyUnsignedAccounts, err = decoder.ReadUint8()
		if err != nil {
			return err
		}
	}
	{
		numAccountKeys, err := bin.DecodeCompactU16LengthFromByteReader(decoder)
		if err != nil {
			return err
		}
		for i := 0; i < numAccountKeys; i++ {
			pubkeyBytes, err := decoder.ReadNBytes(32)
			if err != nil {
				return err
			}
			var sig PublicKey
			copy(sig[:], pubkeyBytes)
			mx.AccountKeys = append(mx.AccountKeys, sig)
		}
	}
	{
		recentBlockhashBytes, err := decoder.ReadNBytes(32)
		if err != nil {
			return err
		}
		var recentBlockhash Hash
		copy(recentBlockhash[:], recentBlockhashBytes)
		mx.RecentBlockhash = recentBlockhash
	}
	{
		numInstructions, err := bin.DecodeCompactU16LengthFromByteReader(decoder)
		if err != nil {
			return err
		}
		for i := 0; i < numInstructions; i++ {
			programIDIndex, err := decoder.ReadUint8()
			if err != nil {
				return err
			}
			var compInst CompiledInstruction
			compInst.ProgramIDIndex = uint16(programIDIndex)

			{
				numAccounts, err := bin.DecodeCompactU16LengthFromByteReader(decoder)
				if err != nil {
					return err
				}
				for i := 0; i < numAccounts; i++ {
					accountIndex, err := decoder.ReadUint8()
					if err != nil {
						return err
					}
					compInst.Accounts = append(compInst.Accounts, uint16(accountIndex))
				}
			}
			{
				dataLen, err := bin.DecodeCompactU16LengthFromByteReader(decoder)
				if err != nil {
					return err
				}
				dataBytes, err := decoder.ReadNBytes(dataLen)
				if err != nil {
					return err
				}
				compInst.Data = Base58(dataBytes)
			}
			mx.Instructions = append(mx.Instructions, compInst)
		}
	}

	return nil
}

func (m *Message) AccountMetaList() AccountMetaSlice {
	out := make(AccountMetaSlice, len(m.AccountKeys))
	for i, a := range m.AccountKeys {
		out[i] = &AccountMeta{
			PublicKey:  a,
			IsSigner:   m.IsSigner(a),
			IsWritable: m.IsWritable(a),
		}
	}
	return out
}

// Signers returns the pubkeys of all accounts that are signers.
func (m *Message) Signers() PublicKeySlice {
	out := make(PublicKeySlice, 0, len(m.AccountKeys))
	for _, a := range m.AccountKeys {
		if m.IsSigner(a) {
			out = append(out, a)
		}
	}
	return out
}

// Writable returns the pubkeys of all accounts that are writable.
func (m *Message) Writable() (out PublicKeySlice) {
	for _, a := range m.AccountKeys {
		if m.IsWritable(a) {
			out = append(out, a)
		}
	}
	return out
}

func (m *Message) ResolveProgramIDIndex(programIDIndex uint16) (PublicKey, error) {
	if int(programIDIndex) < len(m.AccountKeys) {
		return m.AccountKeys[programIDIndex], nil
	}
	return PublicKey{}, fmt.Errorf("programID index not found %d", programIDIndex)
}

func (m *Message) HasAccount(account PublicKey) bool {
	for _, a := range m.AccountKeys {
		if a.Equals(account) {
			return true
		}
	}
	return false
}

func (m *Message) IsSigner(account PublicKey) bool {
	for idx, acc := range m.AccountKeys {
		if acc.Equals(account) {
			return idx < int(m.Header.NumRequiredSignatures)
		}
	}
	return false
}

func (m *Message) IsWritable(account PublicKey) bool {
	index := 0
	found := false
	for idx, acc := range m.AccountKeys {
		if acc.Equals(account) {
			found = true
			index = idx
		}
	}
	if !found {
		return false
	}
	h := m.Header
	return (index < int(h.NumRequiredSignatures-h.NumReadonlySignedAccounts)) ||
		((index >= int(h.NumRequiredSignatures)) && (index < len(m.AccountKeys)-int(h.NumReadonlyUnsignedAccounts)))
}

func (m *Message) signerKeys() []PublicKey {
	return m.AccountKeys[0:m.Header.NumRequiredSignatures]
}

type MessageHeader struct {
	// The total number of signatures required to make the transaction valid.
	// The signatures must match the first `numRequiredSignatures` of `message.account_keys`.
	NumRequiredSignatures uint8 `json:"numRequiredSignatures"`

	// The last numReadonlySignedAccounts of the signed keys are read-only accounts.
	// Programs may process multiple transactions that load read-only accounts within
	// a single PoH entry, but are not permitted to credit or debit lamports or modify
	// account data.
	// Transactions targeting the same read-write account are evaluated sequentially.
	NumReadonlySignedAccounts uint8 `json:"numReadonlySignedAccounts"`

	// The last `numReadonlyUnsignedAccounts` of the unsigned keys are read-only accounts.
	NumReadonlyUnsignedAccounts uint8 `json:"numReadonlyUnsignedAccounts"`
}

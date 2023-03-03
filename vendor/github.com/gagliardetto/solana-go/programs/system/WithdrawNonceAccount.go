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

package system

import (
	"encoding/binary"
	"errors"
	"fmt"

	ag_binary "github.com/gagliardetto/binary"
	ag_solanago "github.com/gagliardetto/solana-go"
	ag_format "github.com/gagliardetto/solana-go/text/format"
	ag_treeout "github.com/gagliardetto/treeout"
)

// Withdraw funds from a nonce account
type WithdrawNonceAccount struct {
	// The u64 parameter is the lamports to withdraw, which must leave the account balance above the rent exempt reserve or at zero.
	Lamports *uint64

	// [0] = [WRITE] NonceAccount
	// ··········· Nonce account
	//
	// [1] = [WRITE] RecipientAccount
	// ··········· Recipient account
	//
	// [2] = [] $(SysVarRecentBlockHashesPubkey)
	// ··········· RecentBlockhashes sysvar
	//
	// [3] = [] $(SysVarRentPubkey)
	// ··········· Rent sysvar
	//
	// [4] = [SIGNER] NonceAuthorityAccount
	// ··········· Nonce authority
	ag_solanago.AccountMetaSlice `bin:"-" borsh_skip:"true"`
}

// NewWithdrawNonceAccountInstructionBuilder creates a new `WithdrawNonceAccount` instruction builder.
func NewWithdrawNonceAccountInstructionBuilder() *WithdrawNonceAccount {
	nd := &WithdrawNonceAccount{
		AccountMetaSlice: make(ag_solanago.AccountMetaSlice, 5),
	}
	nd.AccountMetaSlice[2] = ag_solanago.Meta(ag_solanago.SysVarRecentBlockHashesPubkey)
	nd.AccountMetaSlice[3] = ag_solanago.Meta(ag_solanago.SysVarRentPubkey)
	return nd
}

// The u64 parameter is the lamports to withdraw, which must leave the account balance above the rent exempt reserve or at zero.
func (inst *WithdrawNonceAccount) SetLamports(lamports uint64) *WithdrawNonceAccount {
	inst.Lamports = &lamports
	return inst
}

// Nonce account
func (inst *WithdrawNonceAccount) SetNonceAccount(nonceAccount ag_solanago.PublicKey) *WithdrawNonceAccount {
	inst.AccountMetaSlice[0] = ag_solanago.Meta(nonceAccount).WRITE()
	return inst
}

func (inst *WithdrawNonceAccount) GetNonceAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice[0]
}

// Recipient account
func (inst *WithdrawNonceAccount) SetRecipientAccount(recipientAccount ag_solanago.PublicKey) *WithdrawNonceAccount {
	inst.AccountMetaSlice[1] = ag_solanago.Meta(recipientAccount).WRITE()
	return inst
}

func (inst *WithdrawNonceAccount) GetRecipientAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice[1]
}

// RecentBlockhashes sysvar
func (inst *WithdrawNonceAccount) SetSysVarRecentBlockHashesPubkeyAccount(SysVarRecentBlockHashesPubkey ag_solanago.PublicKey) *WithdrawNonceAccount {
	inst.AccountMetaSlice[2] = ag_solanago.Meta(SysVarRecentBlockHashesPubkey)
	return inst
}

func (inst *WithdrawNonceAccount) GetSysVarRecentBlockHashesPubkeyAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice[2]
}

// Rent sysvar
func (inst *WithdrawNonceAccount) SetSysVarRentPubkeyAccount(SysVarRentPubkey ag_solanago.PublicKey) *WithdrawNonceAccount {
	inst.AccountMetaSlice[3] = ag_solanago.Meta(SysVarRentPubkey)
	return inst
}

func (inst *WithdrawNonceAccount) GetSysVarRentPubkeyAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice[3]
}

// Nonce authority
func (inst *WithdrawNonceAccount) SetNonceAuthorityAccount(nonceAuthorityAccount ag_solanago.PublicKey) *WithdrawNonceAccount {
	inst.AccountMetaSlice[4] = ag_solanago.Meta(nonceAuthorityAccount).SIGNER()
	return inst
}

func (inst *WithdrawNonceAccount) GetNonceAuthorityAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice[4]
}

func (inst WithdrawNonceAccount) Build() *Instruction {
	return &Instruction{BaseVariant: ag_binary.BaseVariant{
		Impl:   inst,
		TypeID: ag_binary.TypeIDFromUint32(Instruction_WithdrawNonceAccount, binary.LittleEndian),
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst WithdrawNonceAccount) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *WithdrawNonceAccount) Validate() error {
	// Check whether all (required) parameters are set:
	{
		if inst.Lamports == nil {
			return errors.New("Lamports parameter is not set")
		}
	}

	// Check whether all accounts are set:
	for accIndex, acc := range inst.AccountMetaSlice {
		if acc == nil {
			return fmt.Errorf("ins.AccountMetaSlice[%v] is not set", accIndex)
		}
	}
	return nil
}

func (inst *WithdrawNonceAccount) EncodeToTree(parent ag_treeout.Branches) {
	parent.Child(ag_format.Program(ProgramName, ProgramID)).
		//
		ParentFunc(func(programBranch ag_treeout.Branches) {
			programBranch.Child(ag_format.Instruction("WithdrawNonceAccount")).
				//
				ParentFunc(func(instructionBranch ag_treeout.Branches) {

					// Parameters of the instruction:
					instructionBranch.Child("Params").ParentFunc(func(paramsBranch ag_treeout.Branches) {
						paramsBranch.Child(ag_format.Param("Lamports", *inst.Lamports))
					})

					// Accounts of the instruction:
					instructionBranch.Child("Accounts").ParentFunc(func(accountsBranch ag_treeout.Branches) {
						accountsBranch.Child(ag_format.Meta("                  Nonce", inst.AccountMetaSlice[0]))
						accountsBranch.Child(ag_format.Meta("              Recipient", inst.AccountMetaSlice[1]))
						accountsBranch.Child(ag_format.Meta("SysVarRecentBlockHashes", inst.AccountMetaSlice[2]))
						accountsBranch.Child(ag_format.Meta("             SysVarRent", inst.AccountMetaSlice[3]))
						accountsBranch.Child(ag_format.Meta("         NonceAuthority", inst.AccountMetaSlice[4]))
					})
				})
		})
}

func (inst WithdrawNonceAccount) MarshalWithEncoder(encoder *ag_binary.Encoder) error {
	// Serialize `Lamports` param:
	{
		err := encoder.Encode(*inst.Lamports)
		if err != nil {
			return err
		}
	}
	return nil
}

func (inst *WithdrawNonceAccount) UnmarshalWithDecoder(decoder *ag_binary.Decoder) error {
	// Deserialize `Lamports` param:
	{
		err := decoder.Decode(&inst.Lamports)
		if err != nil {
			return err
		}
	}
	return nil
}

// NewWithdrawNonceAccountInstruction declares a new WithdrawNonceAccount instruction with the provided parameters and accounts.
func NewWithdrawNonceAccountInstruction(
	// Parameters:
	lamports uint64,
	// Accounts:
	nonceAccount ag_solanago.PublicKey,
	recipientAccount ag_solanago.PublicKey,
	SysVarRecentBlockHashesPubkey ag_solanago.PublicKey,
	SysVarRentPubkey ag_solanago.PublicKey,
	nonceAuthorityAccount ag_solanago.PublicKey) *WithdrawNonceAccount {
	return NewWithdrawNonceAccountInstructionBuilder().
		SetLamports(lamports).
		SetNonceAccount(nonceAccount).
		SetRecipientAccount(recipientAccount).
		SetSysVarRecentBlockHashesPubkeyAccount(SysVarRecentBlockHashesPubkey).
		SetSysVarRentPubkeyAccount(SysVarRentPubkey).
		SetNonceAuthorityAccount(nonceAuthorityAccount)
}

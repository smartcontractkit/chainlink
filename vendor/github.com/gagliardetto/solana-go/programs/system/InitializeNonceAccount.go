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

// Drive state of Uninitalized nonce account to Initialized, setting the nonce value
type InitializeNonceAccount struct {
	// The Pubkey parameter specifies the entity authorized to execute nonce instruction on the account.
	// No signatures are required to execute this instruction, enabling derived nonce account addresses.
	Authorized *ag_solanago.PublicKey

	// [0] = [WRITE] NonceAccount
	// ··········· Nonce account
	//
	// [1] = [] $(SysVarRecentBlockHashesPubkey)
	// ··········· RecentBlockhashes sysvar
	//
	// [2] = [] $(SysVarRentPubkey)
	// ··········· Rent sysvar
	ag_solanago.AccountMetaSlice `bin:"-" borsh_skip:"true"`
}

// NewInitializeNonceAccountInstructionBuilder creates a new `InitializeNonceAccount` instruction builder.
func NewInitializeNonceAccountInstructionBuilder() *InitializeNonceAccount {
	nd := &InitializeNonceAccount{
		AccountMetaSlice: make(ag_solanago.AccountMetaSlice, 3),
	}
	nd.AccountMetaSlice[1] = ag_solanago.Meta(ag_solanago.SysVarRecentBlockHashesPubkey)
	nd.AccountMetaSlice[2] = ag_solanago.Meta(ag_solanago.SysVarRentPubkey)
	return nd
}

// The Pubkey parameter specifies the entity authorized to execute nonce instruction on the account.
// No signatures are required to execute this instruction, enabling derived nonce account addresses.
func (inst *InitializeNonceAccount) SetAuthorized(authorized ag_solanago.PublicKey) *InitializeNonceAccount {
	inst.Authorized = &authorized
	return inst
}

// Nonce account
func (inst *InitializeNonceAccount) SetNonceAccount(nonceAccount ag_solanago.PublicKey) *InitializeNonceAccount {
	inst.AccountMetaSlice[0] = ag_solanago.Meta(nonceAccount).WRITE()
	return inst
}

func (inst *InitializeNonceAccount) GetNonceAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice[0]
}

// RecentBlockhashes sysvar
func (inst *InitializeNonceAccount) SetSysVarRecentBlockHashesPubkeyAccount(SysVarRecentBlockHashesPubkey ag_solanago.PublicKey) *InitializeNonceAccount {
	inst.AccountMetaSlice[1] = ag_solanago.Meta(SysVarRecentBlockHashesPubkey)
	return inst
}

func (inst *InitializeNonceAccount) GetSysVarRecentBlockHashesPubkeyAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice[1]
}

// Rent sysvar
func (inst *InitializeNonceAccount) SetSysVarRentPubkeyAccount(SysVarRentPubkey ag_solanago.PublicKey) *InitializeNonceAccount {
	inst.AccountMetaSlice[2] = ag_solanago.Meta(SysVarRentPubkey)
	return inst
}

func (inst *InitializeNonceAccount) GetSysVarRentPubkeyAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice[2]
}

func (inst InitializeNonceAccount) Build() *Instruction {
	return &Instruction{BaseVariant: ag_binary.BaseVariant{
		Impl:   inst,
		TypeID: ag_binary.TypeIDFromUint32(Instruction_InitializeNonceAccount, binary.LittleEndian),
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst InitializeNonceAccount) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *InitializeNonceAccount) Validate() error {
	// Check whether all (required) parameters are set:
	{
		if inst.Authorized == nil {
			return errors.New("Authorized parameter is not set")
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

func (inst *InitializeNonceAccount) EncodeToTree(parent ag_treeout.Branches) {
	parent.Child(ag_format.Program(ProgramName, ProgramID)).
		//
		ParentFunc(func(programBranch ag_treeout.Branches) {
			programBranch.Child(ag_format.Instruction("InitializeNonceAccount")).
				//
				ParentFunc(func(instructionBranch ag_treeout.Branches) {

					// Parameters of the instruction:
					instructionBranch.Child("Params").ParentFunc(func(paramsBranch ag_treeout.Branches) {
						paramsBranch.Child(ag_format.Param("Authorized", *inst.Authorized))
					})

					// Accounts of the instruction:
					instructionBranch.Child("Accounts").ParentFunc(func(accountsBranch ag_treeout.Branches) {
						accountsBranch.Child(ag_format.Meta("                  Nonce", inst.AccountMetaSlice[0]))
						accountsBranch.Child(ag_format.Meta("SysVarRecentBlockHashes", inst.AccountMetaSlice[1]))
						accountsBranch.Child(ag_format.Meta("             SysVarRent", inst.AccountMetaSlice[2]))
					})
				})
		})
}

func (inst InitializeNonceAccount) MarshalWithEncoder(encoder *ag_binary.Encoder) error {
	// Serialize `Authorized` param:
	{
		err := encoder.Encode(*inst.Authorized)
		if err != nil {
			return err
		}
	}
	return nil
}

func (inst *InitializeNonceAccount) UnmarshalWithDecoder(decoder *ag_binary.Decoder) error {
	// Deserialize `Authorized` param:
	{
		err := decoder.Decode(&inst.Authorized)
		if err != nil {
			return err
		}
	}
	return nil
}

// NewInitializeNonceAccountInstruction declares a new InitializeNonceAccount instruction with the provided parameters and accounts.
func NewInitializeNonceAccountInstruction(
	// Parameters:
	authorized ag_solanago.PublicKey,
	// Accounts:
	nonceAccount ag_solanago.PublicKey,
	SysVarRecentBlockHashesPubkey ag_solanago.PublicKey,
	SysVarRentPubkey ag_solanago.PublicKey) *InitializeNonceAccount {
	return NewInitializeNonceAccountInstructionBuilder().
		SetAuthorized(authorized).
		SetNonceAccount(nonceAccount).
		SetSysVarRecentBlockHashesPubkeyAccount(SysVarRecentBlockHashesPubkey).
		SetSysVarRentPubkeyAccount(SysVarRentPubkey)
}

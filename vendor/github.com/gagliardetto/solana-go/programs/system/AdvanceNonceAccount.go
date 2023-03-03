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
	"fmt"

	ag_binary "github.com/gagliardetto/binary"
	ag_solanago "github.com/gagliardetto/solana-go"
	ag_format "github.com/gagliardetto/solana-go/text/format"
	ag_treeout "github.com/gagliardetto/treeout"
)

// Consumes a stored nonce, replacing it with a successor
type AdvanceNonceAccount struct {

	// [0] = [WRITE] NonceAccount
	// ··········· Nonce account
	//
	// [1] = [] $(SysVarRecentBlockHashesPubkey)
	// ··········· RecentBlockhashes sysvar
	//
	// [2] = [SIGNER] NonceAuthorityAccount
	// ··········· Nonce authority
	ag_solanago.AccountMetaSlice `bin:"-" borsh_skip:"true"`
}

// NewAdvanceNonceAccountInstructionBuilder creates a new `AdvanceNonceAccount` instruction builder.
func NewAdvanceNonceAccountInstructionBuilder() *AdvanceNonceAccount {
	nd := &AdvanceNonceAccount{
		AccountMetaSlice: make(ag_solanago.AccountMetaSlice, 3),
	}
	nd.AccountMetaSlice[1] = ag_solanago.Meta(ag_solanago.SysVarRecentBlockHashesPubkey)
	return nd
}

// Nonce account
func (inst *AdvanceNonceAccount) SetNonceAccount(nonceAccount ag_solanago.PublicKey) *AdvanceNonceAccount {
	inst.AccountMetaSlice[0] = ag_solanago.Meta(nonceAccount).WRITE()
	return inst
}

func (inst *AdvanceNonceAccount) GetNonceAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice[0]
}

// RecentBlockhashes sysvar
func (inst *AdvanceNonceAccount) SetSysVarRecentBlockHashesPubkeyAccount(SysVarRecentBlockHashesPubkey ag_solanago.PublicKey) *AdvanceNonceAccount {
	inst.AccountMetaSlice[1] = ag_solanago.Meta(SysVarRecentBlockHashesPubkey)
	return inst
}

func (inst *AdvanceNonceAccount) GetSysVarRecentBlockHashesPubkeyAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice[1]
}

// Nonce authority
func (inst *AdvanceNonceAccount) SetNonceAuthorityAccount(nonceAuthorityAccount ag_solanago.PublicKey) *AdvanceNonceAccount {
	inst.AccountMetaSlice[2] = ag_solanago.Meta(nonceAuthorityAccount).SIGNER()
	return inst
}

func (inst *AdvanceNonceAccount) GetNonceAuthorityAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice[2]
}

func (inst AdvanceNonceAccount) Build() *Instruction {
	return &Instruction{BaseVariant: ag_binary.BaseVariant{
		Impl:   inst,
		TypeID: ag_binary.TypeIDFromUint32(Instruction_AdvanceNonceAccount, binary.LittleEndian),
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst AdvanceNonceAccount) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *AdvanceNonceAccount) Validate() error {
	// Check whether all accounts are set:
	for accIndex, acc := range inst.AccountMetaSlice {
		if acc == nil {
			return fmt.Errorf("ins.AccountMetaSlice[%v] is not set", accIndex)
		}
	}
	return nil
}

func (inst *AdvanceNonceAccount) EncodeToTree(parent ag_treeout.Branches) {
	parent.Child(ag_format.Program(ProgramName, ProgramID)).
		//
		ParentFunc(func(programBranch ag_treeout.Branches) {
			programBranch.Child(ag_format.Instruction("AdvanceNonceAccount")).
				//
				ParentFunc(func(instructionBranch ag_treeout.Branches) {

					// Parameters of the instruction:
					instructionBranch.Child("Params").ParentFunc(func(paramsBranch ag_treeout.Branches) {})

					// Accounts of the instruction:
					instructionBranch.Child("Accounts").ParentFunc(func(accountsBranch ag_treeout.Branches) {
						accountsBranch.Child(ag_format.Meta("                  Nonce", inst.AccountMetaSlice[0]))
						accountsBranch.Child(ag_format.Meta("SysVarRecentBlockHashes", inst.AccountMetaSlice[1]))
						accountsBranch.Child(ag_format.Meta("         NonceAuthority", inst.AccountMetaSlice[2]))
					})
				})
		})
}

func (inst AdvanceNonceAccount) MarshalWithEncoder(encoder *ag_binary.Encoder) error {
	return nil
}

func (inst *AdvanceNonceAccount) UnmarshalWithDecoder(decoder *ag_binary.Decoder) error {
	return nil
}

// NewAdvanceNonceAccountInstruction declares a new AdvanceNonceAccount instruction with the provided parameters and accounts.
func NewAdvanceNonceAccountInstruction(
	// Accounts:
	nonceAccount ag_solanago.PublicKey,
	SysVarRecentBlockHashesPubkey ag_solanago.PublicKey,
	nonceAuthorityAccount ag_solanago.PublicKey) *AdvanceNonceAccount {
	return NewAdvanceNonceAccountInstructionBuilder().
		SetNonceAccount(nonceAccount).
		SetSysVarRecentBlockHashesPubkeyAccount(SysVarRecentBlockHashesPubkey).
		SetNonceAuthorityAccount(nonceAuthorityAccount)
}

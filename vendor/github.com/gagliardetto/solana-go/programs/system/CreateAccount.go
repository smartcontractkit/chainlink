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

// Create a new account
type CreateAccount struct {
	// Number of lamports to transfer to the new account
	Lamports *uint64

	// Number of bytes of memory to allocate
	Space *uint64

	// Address of program that will own the new account
	Owner *ag_solanago.PublicKey

	// [0] = [WRITE, SIGNER] FundingAccount
	// ··········· Funding account
	//
	// [1] = [WRITE, SIGNER] NewAccount
	// ··········· New account
	ag_solanago.AccountMetaSlice `bin:"-" borsh_skip:"true"`
}

// NewCreateAccountInstructionBuilder creates a new `CreateAccount` instruction builder.
func NewCreateAccountInstructionBuilder() *CreateAccount {
	nd := &CreateAccount{
		AccountMetaSlice: make(ag_solanago.AccountMetaSlice, 2),
	}
	return nd
}

// Number of lamports to transfer to the new account
func (inst *CreateAccount) SetLamports(lamports uint64) *CreateAccount {
	inst.Lamports = &lamports
	return inst
}

// Number of bytes of memory to allocate
func (inst *CreateAccount) SetSpace(space uint64) *CreateAccount {
	inst.Space = &space
	return inst
}

// Address of program that will own the new account
func (inst *CreateAccount) SetOwner(owner ag_solanago.PublicKey) *CreateAccount {
	inst.Owner = &owner
	return inst
}

// Funding account
func (inst *CreateAccount) SetFundingAccount(fundingAccount ag_solanago.PublicKey) *CreateAccount {
	inst.AccountMetaSlice[0] = ag_solanago.Meta(fundingAccount).WRITE().SIGNER()
	return inst
}

func (inst *CreateAccount) GetFundingAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice[0]
}

// New account
func (inst *CreateAccount) SetNewAccount(newAccount ag_solanago.PublicKey) *CreateAccount {
	inst.AccountMetaSlice[1] = ag_solanago.Meta(newAccount).WRITE().SIGNER()
	return inst
}

func (inst *CreateAccount) GetNewAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice[1]
}

func (inst CreateAccount) Build() *Instruction {
	return &Instruction{BaseVariant: ag_binary.BaseVariant{
		Impl:   inst,
		TypeID: ag_binary.TypeIDFromUint32(Instruction_CreateAccount, binary.LittleEndian),
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst CreateAccount) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *CreateAccount) Validate() error {
	// Check whether all (required) parameters are set:
	{
		if inst.Lamports == nil {
			return errors.New("Lamports parameter is not set")
		}
		if inst.Space == nil {
			return errors.New("Space parameter is not set")
		}
		if inst.Owner == nil {
			return errors.New("Owner parameter is not set")
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

func (inst *CreateAccount) EncodeToTree(parent ag_treeout.Branches) {
	parent.Child(ag_format.Program(ProgramName, ProgramID)).
		//
		ParentFunc(func(programBranch ag_treeout.Branches) {
			programBranch.Child(ag_format.Instruction("CreateAccount")).
				//
				ParentFunc(func(instructionBranch ag_treeout.Branches) {

					// Parameters of the instruction:
					instructionBranch.Child("Params").ParentFunc(func(paramsBranch ag_treeout.Branches) {
						paramsBranch.Child(ag_format.Param("Lamports", *inst.Lamports))
						paramsBranch.Child(ag_format.Param("   Space", *inst.Space))
						paramsBranch.Child(ag_format.Param("   Owner", *inst.Owner))
					})

					// Accounts of the instruction:
					instructionBranch.Child("Accounts").ParentFunc(func(accountsBranch ag_treeout.Branches) {
						accountsBranch.Child(ag_format.Meta("Funding", inst.AccountMetaSlice[0]))
						accountsBranch.Child(ag_format.Meta("    New", inst.AccountMetaSlice[1]))
					})
				})
		})
}

func (inst CreateAccount) MarshalWithEncoder(encoder *ag_binary.Encoder) error {
	// Serialize `Lamports` param:
	{
		err := encoder.Encode(*inst.Lamports)
		if err != nil {
			return err
		}
	}
	// Serialize `Space` param:
	{
		err := encoder.Encode(*inst.Space)
		if err != nil {
			return err
		}
	}
	// Serialize `Owner` param:
	{
		err := encoder.Encode(*inst.Owner)
		if err != nil {
			return err
		}
	}
	return nil
}

func (inst *CreateAccount) UnmarshalWithDecoder(decoder *ag_binary.Decoder) error {
	// Deserialize `Lamports` param:
	{
		err := decoder.Decode(&inst.Lamports)
		if err != nil {
			return err
		}
	}
	// Deserialize `Space` param:
	{
		err := decoder.Decode(&inst.Space)
		if err != nil {
			return err
		}
	}
	// Deserialize `Owner` param:
	{
		err := decoder.Decode(&inst.Owner)
		if err != nil {
			return err
		}
	}
	return nil
}

// NewCreateAccountInstruction declares a new CreateAccount instruction with the provided parameters and accounts.
func NewCreateAccountInstruction(
	// Parameters:
	lamports uint64,
	space uint64,
	owner ag_solanago.PublicKey,
	// Accounts:
	fundingAccount ag_solanago.PublicKey,
	newAccount ag_solanago.PublicKey) *CreateAccount {
	return NewCreateAccountInstructionBuilder().
		SetLamports(lamports).
		SetSpace(space).
		SetOwner(owner).
		SetFundingAccount(fundingAccount).
		SetNewAccount(newAccount)
}

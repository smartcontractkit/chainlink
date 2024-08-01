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

// Create a new account at an address derived from a base pubkey and a seed
type CreateAccountWithSeed struct {
	// Base public key
	Base *ag_solanago.PublicKey

	// String of ASCII chars, no longer than Pubkey::MAX_SEED_LEN
	Seed *string

	// Number of lamports to transfer to the new account
	Lamports *uint64

	// Number of bytes of memory to allocate
	Space *uint64

	// Owner program account address
	Owner *ag_solanago.PublicKey

	// [0] = [WRITE, SIGNER] FundingAccount
	// ··········· Funding account
	//
	// [1] = [WRITE] CreatedAccount
	// ··········· Created account
	//
	// [2] = [SIGNER] BaseAccount
	// ··········· Base account
	ag_solanago.AccountMetaSlice `bin:"-" borsh_skip:"true"`
}

// NewCreateAccountWithSeedInstructionBuilder creates a new `CreateAccountWithSeed` instruction builder.
func NewCreateAccountWithSeedInstructionBuilder() *CreateAccountWithSeed {
	nd := &CreateAccountWithSeed{
		AccountMetaSlice: make(ag_solanago.AccountMetaSlice, 3),
	}
	return nd
}

// Base public key
func (inst *CreateAccountWithSeed) SetBase(base ag_solanago.PublicKey) *CreateAccountWithSeed {
	inst.Base = &base
	return inst
}

// String of ASCII chars, no longer than Pubkey::MAX_SEED_LEN
func (inst *CreateAccountWithSeed) SetSeed(seed string) *CreateAccountWithSeed {
	inst.Seed = &seed
	return inst
}

// Number of lamports to transfer to the new account
func (inst *CreateAccountWithSeed) SetLamports(lamports uint64) *CreateAccountWithSeed {
	inst.Lamports = &lamports
	return inst
}

// Number of bytes of memory to allocate
func (inst *CreateAccountWithSeed) SetSpace(space uint64) *CreateAccountWithSeed {
	inst.Space = &space
	return inst
}

// Owner program account address
func (inst *CreateAccountWithSeed) SetOwner(owner ag_solanago.PublicKey) *CreateAccountWithSeed {
	inst.Owner = &owner
	return inst
}

// Funding account
func (inst *CreateAccountWithSeed) SetFundingAccount(fundingAccount ag_solanago.PublicKey) *CreateAccountWithSeed {
	inst.AccountMetaSlice[0] = ag_solanago.Meta(fundingAccount).WRITE().SIGNER()
	return inst
}

func (inst *CreateAccountWithSeed) GetFundingAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice[0]
}

// Created account
func (inst *CreateAccountWithSeed) SetCreatedAccount(createdAccount ag_solanago.PublicKey) *CreateAccountWithSeed {
	inst.AccountMetaSlice[1] = ag_solanago.Meta(createdAccount).WRITE()
	return inst
}

func (inst *CreateAccountWithSeed) GetCreatedAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice[1]
}

// Base account
func (inst *CreateAccountWithSeed) SetBaseAccount(baseAccount ag_solanago.PublicKey) *CreateAccountWithSeed {
	inst.AccountMetaSlice[2] = ag_solanago.Meta(baseAccount).SIGNER()
	return inst
}

func (inst *CreateAccountWithSeed) GetBaseAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice[2]
}

func (inst CreateAccountWithSeed) Build() *Instruction {
	{
		if !inst.Base.Equals(inst.GetFundingAccount().PublicKey) {
			inst.AccountMetaSlice[2] = ag_solanago.Meta(*inst.Base).SIGNER()
		}
	}
	return &Instruction{BaseVariant: ag_binary.BaseVariant{
		Impl:   inst,
		TypeID: ag_binary.TypeIDFromUint32(Instruction_CreateAccountWithSeed, binary.LittleEndian),
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst CreateAccountWithSeed) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *CreateAccountWithSeed) Validate() error {
	// Check whether all (required) parameters are set:
	{
		if inst.Base == nil {
			return errors.New("Base parameter is not set")
		}
		if inst.Seed == nil {
			return errors.New("Seed parameter is not set")
		}
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
	{
		if inst.AccountMetaSlice[0] == nil {
			return fmt.Errorf("FundingAccount is not set")
		}
		if inst.AccountMetaSlice[1] == nil {
			return fmt.Errorf("CreatedAccount is not set")
		}
	}
	return nil
}

func (inst *CreateAccountWithSeed) EncodeToTree(parent ag_treeout.Branches) {
	parent.Child(ag_format.Program(ProgramName, ProgramID)).
		//
		ParentFunc(func(programBranch ag_treeout.Branches) {
			programBranch.Child(ag_format.Instruction("CreateAccountWithSeed")).
				//
				ParentFunc(func(instructionBranch ag_treeout.Branches) {

					// Parameters of the instruction:
					instructionBranch.Child("Params").ParentFunc(func(paramsBranch ag_treeout.Branches) {
						paramsBranch.Child(ag_format.Param("    Base", *inst.Base))
						paramsBranch.Child(ag_format.Param("    Seed", *inst.Seed))
						paramsBranch.Child(ag_format.Param("Lamports", *inst.Lamports))
						paramsBranch.Child(ag_format.Param("   Space", *inst.Space))
						paramsBranch.Child(ag_format.Param("   Owner", *inst.Owner))
					})

					// Accounts of the instruction:
					instructionBranch.Child("Accounts").ParentFunc(func(accountsBranch ag_treeout.Branches) {
						accountsBranch.Child(ag_format.Meta("Funding", inst.AccountMetaSlice[0]))
						accountsBranch.Child(ag_format.Meta("Created", inst.AccountMetaSlice[1]))
						accountsBranch.Child(ag_format.Meta("   Base", inst.AccountMetaSlice[2]))
					})
				})
		})
}

func (inst CreateAccountWithSeed) MarshalWithEncoder(encoder *ag_binary.Encoder) error {
	// Serialize `Base` param:
	{
		err := encoder.Encode(*inst.Base)
		if err != nil {
			return err
		}
	}
	// Serialize `Seed` param:
	{
		err := encoder.WriteRustString(*inst.Seed)
		if err != nil {
			return err
		}
	}
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

func (inst *CreateAccountWithSeed) UnmarshalWithDecoder(decoder *ag_binary.Decoder) error {
	// Deserialize `Base` param:
	{
		err := decoder.Decode(&inst.Base)
		if err != nil {
			return err
		}
	}
	// Deserialize `Seed` param:
	{
		value, err := decoder.ReadRustString()
		if err != nil {
			return err
		}
		inst.Seed = &value
	}
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

// NewCreateAccountWithSeedInstruction declares a new CreateAccountWithSeed instruction with the provided parameters and accounts.
func NewCreateAccountWithSeedInstruction(
	// Parameters:
	base ag_solanago.PublicKey,
	seed string,
	lamports uint64,
	space uint64,
	owner ag_solanago.PublicKey,
	// Accounts:
	fundingAccount ag_solanago.PublicKey,
	createdAccount ag_solanago.PublicKey,
	baseAccount ag_solanago.PublicKey) *CreateAccountWithSeed {
	return NewCreateAccountWithSeedInstructionBuilder().
		SetBase(base).
		SetSeed(seed).
		SetLamports(lamports).
		SetSpace(space).
		SetOwner(owner).
		SetFundingAccount(fundingAccount).
		SetCreatedAccount(createdAccount).
		SetBaseAccount(baseAccount)
}

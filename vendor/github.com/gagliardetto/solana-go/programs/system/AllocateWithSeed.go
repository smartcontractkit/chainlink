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

// Allocate space for and assign an account at an address derived from a base public key and a seed
type AllocateWithSeed struct {
	// Base public key
	Base *ag_solanago.PublicKey

	// String of ASCII chars, no longer than pubkey::MAX_SEED_LEN
	Seed *string

	// Number of bytes of memory to allocate
	Space *uint64

	// Owner program account address
	Owner *ag_solanago.PublicKey

	// [0] = [WRITE] AllocatedAccount
	// ··········· Allocated account
	//
	// [1] = [SIGNER] BaseAccount
	// ··········· Base account
	ag_solanago.AccountMetaSlice `bin:"-" borsh_skip:"true"`
}

// NewAllocateWithSeedInstructionBuilder creates a new `AllocateWithSeed` instruction builder.
func NewAllocateWithSeedInstructionBuilder() *AllocateWithSeed {
	nd := &AllocateWithSeed{
		AccountMetaSlice: make(ag_solanago.AccountMetaSlice, 2),
	}
	return nd
}

// Base public key
func (inst *AllocateWithSeed) SetBase(base ag_solanago.PublicKey) *AllocateWithSeed {
	inst.Base = &base
	return inst
}

// String of ASCII chars, no longer than pubkey::MAX_SEED_LEN
func (inst *AllocateWithSeed) SetSeed(seed string) *AllocateWithSeed {
	inst.Seed = &seed
	return inst
}

// Number of bytes of memory to allocate
func (inst *AllocateWithSeed) SetSpace(space uint64) *AllocateWithSeed {
	inst.Space = &space
	return inst
}

// Owner program account address
func (inst *AllocateWithSeed) SetOwner(owner ag_solanago.PublicKey) *AllocateWithSeed {
	inst.Owner = &owner
	return inst
}

// Allocated account
func (inst *AllocateWithSeed) SetAllocatedAccount(allocatedAccount ag_solanago.PublicKey) *AllocateWithSeed {
	inst.AccountMetaSlice[0] = ag_solanago.Meta(allocatedAccount).WRITE()
	return inst
}

func (inst *AllocateWithSeed) GetAllocatedAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice[0]
}

// Base account
func (inst *AllocateWithSeed) SetBaseAccount(baseAccount ag_solanago.PublicKey) *AllocateWithSeed {
	inst.AccountMetaSlice[1] = ag_solanago.Meta(baseAccount).SIGNER()
	return inst
}

func (inst *AllocateWithSeed) GetBaseAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice[1]
}

func (inst AllocateWithSeed) Build() *Instruction {
	return &Instruction{BaseVariant: ag_binary.BaseVariant{
		Impl:   inst,
		TypeID: ag_binary.TypeIDFromUint32(Instruction_AllocateWithSeed, binary.LittleEndian),
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst AllocateWithSeed) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *AllocateWithSeed) Validate() error {
	// Check whether all (required) parameters are set:
	{
		if inst.Base == nil {
			return errors.New("Base parameter is not set")
		}
		if inst.Seed == nil {
			return errors.New("Seed parameter is not set")
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

func (inst *AllocateWithSeed) EncodeToTree(parent ag_treeout.Branches) {
	parent.Child(ag_format.Program(ProgramName, ProgramID)).
		//
		ParentFunc(func(programBranch ag_treeout.Branches) {
			programBranch.Child(ag_format.Instruction("AllocateWithSeed")).
				//
				ParentFunc(func(instructionBranch ag_treeout.Branches) {

					// Parameters of the instruction:
					instructionBranch.Child("Params").ParentFunc(func(paramsBranch ag_treeout.Branches) {
						paramsBranch.Child(ag_format.Param(" Base", *inst.Base))
						paramsBranch.Child(ag_format.Param(" Seed", *inst.Seed))
						paramsBranch.Child(ag_format.Param("Space", *inst.Space))
						paramsBranch.Child(ag_format.Param("Owner", *inst.Owner))
					})

					// Accounts of the instruction:
					instructionBranch.Child("Accounts").ParentFunc(func(accountsBranch ag_treeout.Branches) {
						accountsBranch.Child(ag_format.Meta("Allocated", inst.AccountMetaSlice[0]))
						accountsBranch.Child(ag_format.Meta("     Base", inst.AccountMetaSlice[1]))
					})
				})
		})
}

func (inst AllocateWithSeed) MarshalWithEncoder(encoder *ag_binary.Encoder) error {
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

func (inst *AllocateWithSeed) UnmarshalWithDecoder(decoder *ag_binary.Decoder) error {
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

// NewAllocateWithSeedInstruction declares a new AllocateWithSeed instruction with the provided parameters and accounts.
func NewAllocateWithSeedInstruction(
	// Parameters:
	base ag_solanago.PublicKey,
	seed string,
	space uint64,
	owner ag_solanago.PublicKey,
	// Accounts:
	allocatedAccount ag_solanago.PublicKey,
	baseAccount ag_solanago.PublicKey) *AllocateWithSeed {
	return NewAllocateWithSeedInstructionBuilder().
		SetBase(base).
		SetSeed(seed).
		SetSpace(space).
		SetOwner(owner).
		SetAllocatedAccount(allocatedAccount).
		SetBaseAccount(baseAccount)
}

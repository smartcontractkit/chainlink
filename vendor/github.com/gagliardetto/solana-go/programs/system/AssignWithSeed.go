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

// Assign account to a program based on a seed
type AssignWithSeed struct {
	// Base public key
	Base *ag_solanago.PublicKey

	// String of ASCII chars, no longer than pubkey::MAX_SEED_LEN
	Seed *string

	// Owner program account
	Owner *ag_solanago.PublicKey

	// [0] = [WRITE] AssignedAccount
	// ··········· Assigned account
	//
	// [1] = [SIGNER] BaseAccount
	// ··········· Base account
	ag_solanago.AccountMetaSlice `bin:"-" borsh_skip:"true"`
}

// NewAssignWithSeedInstructionBuilder creates a new `AssignWithSeed` instruction builder.
func NewAssignWithSeedInstructionBuilder() *AssignWithSeed {
	nd := &AssignWithSeed{
		AccountMetaSlice: make(ag_solanago.AccountMetaSlice, 2),
	}
	return nd
}

// Base public key
func (inst *AssignWithSeed) SetBase(base ag_solanago.PublicKey) *AssignWithSeed {
	inst.Base = &base
	return inst
}

// String of ASCII chars, no longer than pubkey::MAX_SEED_LEN
func (inst *AssignWithSeed) SetSeed(seed string) *AssignWithSeed {
	inst.Seed = &seed
	return inst
}

// Owner program account
func (inst *AssignWithSeed) SetOwner(owner ag_solanago.PublicKey) *AssignWithSeed {
	inst.Owner = &owner
	return inst
}

// Assigned account
func (inst *AssignWithSeed) SetAssignedAccount(assignedAccount ag_solanago.PublicKey) *AssignWithSeed {
	inst.AccountMetaSlice[0] = ag_solanago.Meta(assignedAccount).WRITE()
	return inst
}

func (inst *AssignWithSeed) GetAssignedAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice[0]
}

// Base account
func (inst *AssignWithSeed) SetBaseAccount(baseAccount ag_solanago.PublicKey) *AssignWithSeed {
	inst.AccountMetaSlice[1] = ag_solanago.Meta(baseAccount).SIGNER()
	return inst
}

func (inst *AssignWithSeed) GetBaseAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice[1]
}

func (inst AssignWithSeed) Build() *Instruction {
	return &Instruction{BaseVariant: ag_binary.BaseVariant{
		Impl:   inst,
		TypeID: ag_binary.TypeIDFromUint32(Instruction_AssignWithSeed, binary.LittleEndian),
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst AssignWithSeed) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *AssignWithSeed) Validate() error {
	// Check whether all (required) parameters are set:
	{
		if inst.Base == nil {
			return errors.New("Base parameter is not set")
		}
		if inst.Seed == nil {
			return errors.New("Seed parameter is not set")
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

func (inst *AssignWithSeed) EncodeToTree(parent ag_treeout.Branches) {
	parent.Child(ag_format.Program(ProgramName, ProgramID)).
		//
		ParentFunc(func(programBranch ag_treeout.Branches) {
			programBranch.Child(ag_format.Instruction("AssignWithSeed")).
				//
				ParentFunc(func(instructionBranch ag_treeout.Branches) {

					// Parameters of the instruction:
					instructionBranch.Child("Params").ParentFunc(func(paramsBranch ag_treeout.Branches) {
						paramsBranch.Child(ag_format.Param(" Base", *inst.Base))
						paramsBranch.Child(ag_format.Param(" Seed", *inst.Seed))
						paramsBranch.Child(ag_format.Param("Owner", *inst.Owner))
					})

					// Accounts of the instruction:
					instructionBranch.Child("Accounts").ParentFunc(func(accountsBranch ag_treeout.Branches) {
						accountsBranch.Child(ag_format.Meta("Assigned", inst.AccountMetaSlice[0]))
						accountsBranch.Child(ag_format.Meta("    Base", inst.AccountMetaSlice[1]))
					})
				})
		})
}

func (inst AssignWithSeed) MarshalWithEncoder(encoder *ag_binary.Encoder) error {
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
	// Serialize `Owner` param:
	{
		err := encoder.Encode(*inst.Owner)
		if err != nil {
			return err
		}
	}
	return nil
}

func (inst *AssignWithSeed) UnmarshalWithDecoder(decoder *ag_binary.Decoder) error {
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
	// Deserialize `Owner` param:
	{
		err := decoder.Decode(&inst.Owner)
		if err != nil {
			return err
		}
	}
	return nil
}

// NewAssignWithSeedInstruction declares a new AssignWithSeed instruction with the provided parameters and accounts.
func NewAssignWithSeedInstruction(
	// Parameters:
	base ag_solanago.PublicKey,
	seed string,
	owner ag_solanago.PublicKey,
	// Accounts:
	assignedAccount ag_solanago.PublicKey,
	baseAccount ag_solanago.PublicKey) *AssignWithSeed {
	return NewAssignWithSeedInstructionBuilder().
		SetBase(base).
		SetSeed(seed).
		SetOwner(owner).
		SetAssignedAccount(assignedAccount).
		SetBaseAccount(baseAccount)
}

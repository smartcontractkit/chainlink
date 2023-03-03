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

// Assign account to a program
type Assign struct {
	// Owner program account
	Owner *ag_solanago.PublicKey

	// [0] = [WRITE, SIGNER] AssignedAccount
	// ··········· Assigned account public key
	ag_solanago.AccountMetaSlice `bin:"-" borsh_skip:"true"`
}

// NewAssignInstructionBuilder creates a new `Assign` instruction builder.
func NewAssignInstructionBuilder() *Assign {
	nd := &Assign{
		AccountMetaSlice: make(ag_solanago.AccountMetaSlice, 1),
	}
	return nd
}

// Owner program account
func (inst *Assign) SetOwner(owner ag_solanago.PublicKey) *Assign {
	inst.Owner = &owner
	return inst
}

// Assigned account public key
func (inst *Assign) SetAssignedAccount(assignedAccount ag_solanago.PublicKey) *Assign {
	inst.AccountMetaSlice[0] = ag_solanago.Meta(assignedAccount).WRITE().SIGNER()
	return inst
}

func (inst *Assign) GetAssignedAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice[0]
}

func (inst Assign) Build() *Instruction {
	return &Instruction{BaseVariant: ag_binary.BaseVariant{
		Impl:   inst,
		TypeID: ag_binary.TypeIDFromUint32(Instruction_Assign, binary.LittleEndian),
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst Assign) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *Assign) Validate() error {
	// Check whether all (required) parameters are set:
	{
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

func (inst *Assign) EncodeToTree(parent ag_treeout.Branches) {
	parent.Child(ag_format.Program(ProgramName, ProgramID)).
		//
		ParentFunc(func(programBranch ag_treeout.Branches) {
			programBranch.Child(ag_format.Instruction("Assign")).
				//
				ParentFunc(func(instructionBranch ag_treeout.Branches) {

					// Parameters of the instruction:
					instructionBranch.Child("Params").ParentFunc(func(paramsBranch ag_treeout.Branches) {
						paramsBranch.Child(ag_format.Param("Owner", *inst.Owner))
					})

					// Accounts of the instruction:
					instructionBranch.Child("Accounts").ParentFunc(func(accountsBranch ag_treeout.Branches) {
						accountsBranch.Child(ag_format.Meta("Assigned", inst.AccountMetaSlice[0]))
					})
				})
		})
}

func (inst Assign) MarshalWithEncoder(encoder *ag_binary.Encoder) error {
	// Serialize `Owner` param:
	{
		err := encoder.Encode(*inst.Owner)
		if err != nil {
			return err
		}
	}
	return nil
}

func (inst *Assign) UnmarshalWithDecoder(decoder *ag_binary.Decoder) error {
	// Deserialize `Owner` param:
	{
		err := decoder.Decode(&inst.Owner)
		if err != nil {
			return err
		}
	}
	return nil
}

// NewAssignInstruction declares a new Assign instruction with the provided parameters and accounts.
func NewAssignInstruction(
	// Parameters:
	owner ag_solanago.PublicKey,
	// Accounts:
	assignedAccount ag_solanago.PublicKey) *Assign {
	return NewAssignInstructionBuilder().
		SetOwner(owner).
		SetAssignedAccount(assignedAccount)
}

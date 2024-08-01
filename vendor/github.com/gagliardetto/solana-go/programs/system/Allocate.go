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

// Allocate space in a (possibly new) account without funding
type Allocate struct {
	// Number of bytes of memory to allocate
	Space *uint64

	// [0] = [WRITE, SIGNER] NewAccount
	// ··········· New account
	ag_solanago.AccountMetaSlice `bin:"-" borsh_skip:"true"`
}

// NewAllocateInstructionBuilder creates a new `Allocate` instruction builder.
func NewAllocateInstructionBuilder() *Allocate {
	nd := &Allocate{
		AccountMetaSlice: make(ag_solanago.AccountMetaSlice, 1),
	}
	return nd
}

// Number of bytes of memory to allocate
func (inst *Allocate) SetSpace(space uint64) *Allocate {
	inst.Space = &space
	return inst
}

// New account
func (inst *Allocate) SetNewAccount(newAccount ag_solanago.PublicKey) *Allocate {
	inst.AccountMetaSlice[0] = ag_solanago.Meta(newAccount).WRITE().SIGNER()
	return inst
}

func (inst *Allocate) GetNewAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice[0]
}

func (inst Allocate) Build() *Instruction {
	return &Instruction{BaseVariant: ag_binary.BaseVariant{
		Impl:   inst,
		TypeID: ag_binary.TypeIDFromUint32(Instruction_Allocate, binary.LittleEndian),
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst Allocate) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *Allocate) Validate() error {
	// Check whether all (required) parameters are set:
	{
		if inst.Space == nil {
			return errors.New("Space parameter is not set")
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

func (inst *Allocate) EncodeToTree(parent ag_treeout.Branches) {
	parent.Child(ag_format.Program(ProgramName, ProgramID)).
		//
		ParentFunc(func(programBranch ag_treeout.Branches) {
			programBranch.Child(ag_format.Instruction("Allocate")).
				//
				ParentFunc(func(instructionBranch ag_treeout.Branches) {

					// Parameters of the instruction:
					instructionBranch.Child("Params").ParentFunc(func(paramsBranch ag_treeout.Branches) {
						paramsBranch.Child(ag_format.Param("Space", *inst.Space))
					})

					// Accounts of the instruction:
					instructionBranch.Child("Accounts").ParentFunc(func(accountsBranch ag_treeout.Branches) {
						accountsBranch.Child(ag_format.Meta("New", inst.AccountMetaSlice[0]))
					})
				})
		})
}

func (inst Allocate) MarshalWithEncoder(encoder *ag_binary.Encoder) error {
	// Serialize `Space` param:
	{
		err := encoder.Encode(*inst.Space)
		if err != nil {
			return err
		}
	}
	return nil
}

func (inst *Allocate) UnmarshalWithDecoder(decoder *ag_binary.Decoder) error {
	// Deserialize `Space` param:
	{
		err := decoder.Decode(&inst.Space)
		if err != nil {
			return err
		}
	}
	return nil
}

// NewAllocateInstruction declares a new Allocate instruction with the provided parameters and accounts.
func NewAllocateInstruction(
	// Parameters:
	space uint64,
	// Accounts:
	newAccount ag_solanago.PublicKey) *Allocate {
	return NewAllocateInstructionBuilder().
		SetSpace(space).
		SetNewAccount(newAccount)
}

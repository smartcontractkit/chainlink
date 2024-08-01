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

// Transfer lamports
type Transfer struct {
	// Number of lamports to transfer to the new account
	Lamports *uint64

	// [0] = [WRITE, SIGNER] FundingAccount
	// ··········· Funding account
	//
	// [1] = [WRITE] RecipientAccount
	// ··········· Recipient account
	ag_solanago.AccountMetaSlice `bin:"-" borsh_skip:"true"`
}

// NewTransferInstructionBuilder creates a new `Transfer` instruction builder.
func NewTransferInstructionBuilder() *Transfer {
	nd := &Transfer{
		AccountMetaSlice: make(ag_solanago.AccountMetaSlice, 2),
	}
	return nd
}

// Number of lamports to transfer to the new account
func (inst *Transfer) SetLamports(lamports uint64) *Transfer {
	inst.Lamports = &lamports
	return inst
}

// Funding account
func (inst *Transfer) SetFundingAccount(fundingAccount ag_solanago.PublicKey) *Transfer {
	inst.AccountMetaSlice[0] = ag_solanago.Meta(fundingAccount).WRITE().SIGNER()
	return inst
}

func (inst *Transfer) GetFundingAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice[0]
}

// Recipient account
func (inst *Transfer) SetRecipientAccount(recipientAccount ag_solanago.PublicKey) *Transfer {
	inst.AccountMetaSlice[1] = ag_solanago.Meta(recipientAccount).WRITE()
	return inst
}

func (inst *Transfer) GetRecipientAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice[1]
}

func (inst Transfer) Build() *Instruction {
	return &Instruction{BaseVariant: ag_binary.BaseVariant{
		Impl:   inst,
		TypeID: ag_binary.TypeIDFromUint32(Instruction_Transfer, binary.LittleEndian),
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst Transfer) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *Transfer) Validate() error {
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

func (inst *Transfer) EncodeToTree(parent ag_treeout.Branches) {
	parent.Child(ag_format.Program(ProgramName, ProgramID)).
		//
		ParentFunc(func(programBranch ag_treeout.Branches) {
			programBranch.Child(ag_format.Instruction("Transfer")).
				//
				ParentFunc(func(instructionBranch ag_treeout.Branches) {

					// Parameters of the instruction:
					instructionBranch.Child("Params").ParentFunc(func(paramsBranch ag_treeout.Branches) {
						paramsBranch.Child(ag_format.Param("Lamports", *inst.Lamports))
					})

					// Accounts of the instruction:
					instructionBranch.Child("Accounts").ParentFunc(func(accountsBranch ag_treeout.Branches) {
						accountsBranch.Child(ag_format.Meta("  Funding", inst.AccountMetaSlice[0]))
						accountsBranch.Child(ag_format.Meta("Recipient", inst.AccountMetaSlice[1]))
					})
				})
		})
}

func (inst Transfer) MarshalWithEncoder(encoder *ag_binary.Encoder) error {
	// Serialize `Lamports` param:
	{
		err := encoder.Encode(*inst.Lamports)
		if err != nil {
			return err
		}
	}
	return nil
}

func (inst *Transfer) UnmarshalWithDecoder(decoder *ag_binary.Decoder) error {
	// Deserialize `Lamports` param:
	{
		err := decoder.Decode(&inst.Lamports)
		if err != nil {
			return err
		}
	}
	return nil
}

// NewTransferInstruction declares a new Transfer instruction with the provided parameters and accounts.
func NewTransferInstruction(
	// Parameters:
	lamports uint64,
	// Accounts:
	fundingAccount ag_solanago.PublicKey,
	recipientAccount ag_solanago.PublicKey) *Transfer {
	return NewTransferInstructionBuilder().
		SetLamports(lamports).
		SetFundingAccount(fundingAccount).
		SetRecipientAccount(recipientAccount)
}

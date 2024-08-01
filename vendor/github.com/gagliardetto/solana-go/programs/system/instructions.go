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

// Create new accounts, allocate account data, assign accounts to owning programs,
// transfer lamports from System Program owned accounts and pay transacation fees.

package system

import (
	"bytes"
	"encoding/binary"
	"fmt"

	ag_spew "github.com/davecgh/go-spew/spew"
	ag_binary "github.com/gagliardetto/binary"
	ag_solanago "github.com/gagliardetto/solana-go"
	ag_text "github.com/gagliardetto/solana-go/text"
	ag_treeout "github.com/gagliardetto/treeout"
)

var ProgramID ag_solanago.PublicKey = ag_solanago.SystemProgramID

func SetProgramID(pubkey ag_solanago.PublicKey) {
	ProgramID = pubkey
	ag_solanago.RegisterInstructionDecoder(ProgramID, registryDecodeInstruction)
}

const ProgramName = "System"

func init() {
	ag_solanago.RegisterInstructionDecoder(ProgramID, registryDecodeInstruction)
}

const (
	// Create a new account
	Instruction_CreateAccount uint32 = iota

	// Assign account to a program
	Instruction_Assign

	// Transfer lamports
	Instruction_Transfer

	// Create a new account at an address derived from a base pubkey and a seed
	Instruction_CreateAccountWithSeed

	// Consumes a stored nonce, replacing it with a successor
	Instruction_AdvanceNonceAccount

	// Withdraw funds from a nonce account
	Instruction_WithdrawNonceAccount

	// Drive state of Uninitalized nonce account to Initialized, setting the nonce value
	Instruction_InitializeNonceAccount

	// Change the entity authorized to execute nonce instructions on the account
	Instruction_AuthorizeNonceAccount

	// Allocate space in a (possibly new) account without funding
	Instruction_Allocate

	// Allocate space for and assign an account at an address derived from a base public key and a seed
	Instruction_AllocateWithSeed

	// Assign account to a program based on a seed
	Instruction_AssignWithSeed

	// Transfer lamports from a derived address
	Instruction_TransferWithSeed
)

// InstructionIDToName returns the name of the instruction given its ID.
func InstructionIDToName(id uint32) string {
	switch id {
	case Instruction_CreateAccount:
		return "CreateAccount"
	case Instruction_Assign:
		return "Assign"
	case Instruction_Transfer:
		return "Transfer"
	case Instruction_CreateAccountWithSeed:
		return "CreateAccountWithSeed"
	case Instruction_AdvanceNonceAccount:
		return "AdvanceNonceAccount"
	case Instruction_WithdrawNonceAccount:
		return "WithdrawNonceAccount"
	case Instruction_InitializeNonceAccount:
		return "InitializeNonceAccount"
	case Instruction_AuthorizeNonceAccount:
		return "AuthorizeNonceAccount"
	case Instruction_Allocate:
		return "Allocate"
	case Instruction_AllocateWithSeed:
		return "AllocateWithSeed"
	case Instruction_AssignWithSeed:
		return "AssignWithSeed"
	case Instruction_TransferWithSeed:
		return "TransferWithSeed"
	default:
		return ""
	}
}

type Instruction struct {
	ag_binary.BaseVariant
}

func (inst *Instruction) EncodeToTree(parent ag_treeout.Branches) {
	if enToTree, ok := inst.Impl.(ag_text.EncodableToTree); ok {
		enToTree.EncodeToTree(parent)
	} else {
		parent.Child(ag_spew.Sdump(inst))
	}
}

var InstructionImplDef = ag_binary.NewVariantDefinition(
	ag_binary.Uint32TypeIDEncoding,
	[]ag_binary.VariantType{
		{
			"CreateAccount", (*CreateAccount)(nil),
		},
		{
			"Assign", (*Assign)(nil),
		},
		{
			"Transfer", (*Transfer)(nil),
		},
		{
			"CreateAccountWithSeed", (*CreateAccountWithSeed)(nil),
		},
		{
			"AdvanceNonceAccount", (*AdvanceNonceAccount)(nil),
		},
		{
			"WithdrawNonceAccount", (*WithdrawNonceAccount)(nil),
		},
		{
			"InitializeNonceAccount", (*InitializeNonceAccount)(nil),
		},
		{
			"AuthorizeNonceAccount", (*AuthorizeNonceAccount)(nil),
		},
		{
			"Allocate", (*Allocate)(nil),
		},
		{
			"AllocateWithSeed", (*AllocateWithSeed)(nil),
		},
		{
			"AssignWithSeed", (*AssignWithSeed)(nil),
		},
		{
			"TransferWithSeed", (*TransferWithSeed)(nil),
		},
	},
)

func (inst *Instruction) ProgramID() ag_solanago.PublicKey {
	return ProgramID
}

func (inst *Instruction) Accounts() (out []*ag_solanago.AccountMeta) {
	return inst.Impl.(ag_solanago.AccountsGettable).GetAccounts()
}

func (inst *Instruction) Data() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := ag_binary.NewBinEncoder(buf).Encode(inst); err != nil {
		return nil, fmt.Errorf("unable to encode instruction: %w", err)
	}
	return buf.Bytes(), nil
}

func (inst *Instruction) TextEncode(encoder *ag_text.Encoder, option *ag_text.Option) error {
	return encoder.Encode(inst.Impl, option)
}

func (inst *Instruction) UnmarshalWithDecoder(decoder *ag_binary.Decoder) error {
	return inst.BaseVariant.UnmarshalBinaryVariant(decoder, InstructionImplDef)
}

func (inst Instruction) MarshalWithEncoder(encoder *ag_binary.Encoder) error {
	err := encoder.WriteUint32(inst.TypeID.Uint32(), binary.LittleEndian)
	if err != nil {
		return fmt.Errorf("unable to write variant type: %w", err)
	}
	return encoder.Encode(inst.Impl)
}

func registryDecodeInstruction(accounts []*ag_solanago.AccountMeta, data []byte) (interface{}, error) {
	inst, err := DecodeInstruction(accounts, data)
	if err != nil {
		return nil, err
	}
	return inst, nil
}

func DecodeInstruction(accounts []*ag_solanago.AccountMeta, data []byte) (*Instruction, error) {
	inst := new(Instruction)
	if err := ag_binary.NewBinDecoder(data).Decode(inst); err != nil {
		return nil, fmt.Errorf("unable to decode instruction: %w", err)
	}
	if v, ok := inst.Impl.(ag_solanago.AccountsSettable); ok {
		err := v.SetAccounts(accounts)
		if err != nil {
			return nil, fmt.Errorf("unable to set accounts for instruction: %w", err)
		}
	}
	return inst, nil
}

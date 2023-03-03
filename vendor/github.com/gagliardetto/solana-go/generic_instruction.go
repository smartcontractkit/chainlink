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

package solana

// NewInstruction creates a generic instruction with the provided
// programID, accounts, and data bytes.
func NewInstruction(
	programID PublicKey,
	accounts AccountMetaSlice,
	data []byte,
) *GenericInstruction {
	return &GenericInstruction{
		AccountValues: accounts,
		ProgID:        programID,
		DataBytes:     data,
	}
}

var _ Instruction = &GenericInstruction{}

type GenericInstruction struct {
	AccountValues AccountMetaSlice
	ProgID        PublicKey
	DataBytes     []byte
}

func (in *GenericInstruction) ProgramID() PublicKey {
	return in.ProgID
}

func (in *GenericInstruction) Accounts() []*AccountMeta {
	return in.AccountValues
}

func (in *GenericInstruction) Data() ([]byte, error) {
	return in.DataBytes, nil
}

// Copyright 2021 github.com/gagliardetto
// This file has been modified by github.com/gagliardetto
//
// Copyright 2020 dfuse Platform Inc.
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

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

var ErrInstructionDecoderNotFound = errors.New("instruction decoder not found")

// InstructionDecoder receives the AccountMeta FOR THAT INSTRUCTION,
// and not the accounts of the *Message object. Resolve with
// CompiledInstruction.ResolveInstructionAccounts(message) beforehand.
type InstructionDecoder func(instructionAccounts []*AccountMeta, data []byte) (interface{}, error)

var instructionDecoderRegistry = newInstructionDecoderRegistry()

type decoderRegistry struct {
	mu       *sync.RWMutex
	decoders map[PublicKey]InstructionDecoder
}

func newInstructionDecoderRegistry() *decoderRegistry {
	return &decoderRegistry{
		mu:       &sync.RWMutex{},
		decoders: make(map[PublicKey]InstructionDecoder),
	}
}

func (reg *decoderRegistry) Has(programID PublicKey) bool {
	reg.mu.RLock()
	defer reg.mu.RUnlock()

	_, ok := reg.decoders[programID]
	return ok
}

func (reg *decoderRegistry) Get(programID PublicKey) (InstructionDecoder, bool) {
	reg.mu.RLock()
	defer reg.mu.RUnlock()

	decoder, ok := reg.decoders[programID]
	return decoder, ok
}

// RegisterIfNew registers the provided decoder for the provided programID ONLY if there isn't
// already a registered decoder for the programID.
// Returns true if was successfully registered right now (non-previously registered);
// returns false if there already was a decoder registered.
func (reg *decoderRegistry) RegisterIfNew(programID PublicKey, decoder InstructionDecoder) bool {
	reg.mu.Lock()
	defer reg.mu.Unlock()

	_, ok := reg.decoders[programID]
	if ok {
		return false
	}
	reg.decoders[programID] = decoder
	return true
}

func RegisterInstructionDecoder(programID PublicKey, decoder InstructionDecoder) {
	prev, has := instructionDecoderRegistry.Get(programID)
	if has {
		// If it's the same function, then OK (tollerate multiple calls with same params).
		if isSameFunction(prev, decoder) {
			return
		}
		// If it's another decoder for the same pubkey, then panic.
		panic(fmt.Sprintf("unable to re-register instruction decoder for program %s", programID))
	}
	instructionDecoderRegistry.RegisterIfNew(programID, decoder)
}

func isSameFunction(f1 interface{}, f2 interface{}) bool {
	return reflect.ValueOf(f1).Pointer() == reflect.ValueOf(f2).Pointer()
}

func DecodeInstruction(programID PublicKey, accounts []*AccountMeta, data []byte) (interface{}, error) {
	decoder, found := instructionDecoderRegistry.Get(programID)
	if !found {
		return nil, ErrInstructionDecoderNotFound
	}
	return decoder(accounts, data)
}

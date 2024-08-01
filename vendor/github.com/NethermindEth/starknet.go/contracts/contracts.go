package contracts

import (
	"encoding/json"
	"os"

	"github.com/NethermindEth/juno/core/felt"
)

type CasmClass struct {
	Prime            string                     `json:"prime"`
	Version          string                     `json:"compiler_version"`
	ByteCode         []*felt.Felt               `json:"bytecode"`
	EntryPointByType CasmClassEntryPointsByType `json:"entry_points_by_type"`
	// Hints            any                        `json:"hints"`
}

type CasmClassEntryPointsByType struct {
	Constructor []CasmClassEntryPoint `json:"CONSTRUCTOR"`
	External    []CasmClassEntryPoint `json:"EXTERNAL"`
	L1Handler   []CasmClassEntryPoint `json:"L1_HANDLER"`
}

type CasmClassEntryPoint struct {
	Selector *felt.Felt `json:"selector"`
	Offset   int        `json:"offset"`
	Builtins []string   `json:"builtins"`
}

// UnmarshalCasmClass is a function that unmarshals a CasmClass object from a file.
// CASM = Cairo instructions
//
// It takes a file path as a parameter and returns a pointer to the unmarshaled CasmClass object and an error.
func UnmarshalCasmClass(filePath string) (*CasmClass, error) {

	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var casmClass CasmClass
	err = json.Unmarshal(content, &casmClass)
	if err != nil {
		return nil, err
	}

	return &casmClass, nil
}

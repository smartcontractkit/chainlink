package rpc

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/NethermindEth/juno/core/felt"
)

// An integer number in hex format (0x...)
type NumAsHex string

// 64 bit integers, represented by hex string of length at most 16
type U64 string

// ToUint64 converts the U64 type to a uint64.
func (u U64) ToUint64() (uint64, error) {
	hexStr := strings.TrimPrefix(string(u), "0x")

	val, err := strconv.ParseUint(hexStr, 16, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse hex string: %v", err)
	}

	return val, nil
}

// 64 bit integers, represented by hex string of length at most 32
type U128 string

type DeprecatedCairoEntryPoint struct {
	// The offset of the entry point in the program
	Offset NumAsHex `json:"offset"`
	// A unique  identifier of the entry point (function) in the program
	Selector *felt.Felt `json:"selector"`
}

type ClassOutput interface{}

var _ ClassOutput = &DeprecatedContractClass{}
var _ ClassOutput = &ContractClass{}

type ABI []ABIEntry

type DeprecatedEntryPointsByType struct {
	Constructor []DeprecatedCairoEntryPoint `json:"CONSTRUCTOR"`
	External    []DeprecatedCairoEntryPoint `json:"EXTERNAL"`
	L1Handler   []DeprecatedCairoEntryPoint `json:"L1_HANDLER"`
}

type DeprecatedContractClass struct {
	// Program A base64 representation of the compressed program code
	Program string `json:"program"`

	DeprecatedEntryPointsByType DeprecatedEntryPointsByType `json:"entry_points_by_type"`

	ABI *ABI `json:"abi,omitempty"`
}

type ContractClass struct {
	// The list of Sierra instructions of which the program consists
	SierraProgram []*felt.Felt `json:"sierra_program"`

	// The version of the contract class object. Currently, the Starknet OS supports version 0.1.0
	ContractClassVersion string `json:"contract_class_version"`

	EntryPointsByType EntryPointsByType `json:"entry_points_by_type"`

	ABI string `json:"abi,omitempty"`
}

// UnmarshalJSON unmarshals the JSON content into the DeprecatedContractClass struct.
//
// It takes a byte array `content` as a parameter and returns an error if there is any.
// The function processes the `program` field in the JSON object.
// If `program` is a string, it is assigned to the `Program` field in the struct.
// Otherwise, it is encoded and assigned to the `Program` field.
// The function then processes the `entry_points_by_type` field in the JSON object.
// The value is unmarshaled into the `DeprecatedEntryPointsByType` field in the struct.
// Finally, the function processes the `abi` field in the JSON object.
// If it is missing, the function returns nil.
// Otherwise, it unmarshals the value into a slice of interfaces.
// For each element in the slice, it checks the type and assigns it to the appropriate field in the `ABI` field in the struct.
//
// Parameters:
// - content: byte array
// Returns:
// - error: error if there is any
func (c *DeprecatedContractClass) UnmarshalJSON(content []byte) error {
	v := map[string]json.RawMessage{}
	if err := json.Unmarshal(content, &v); err != nil {
		return err
	}

	// process 'program'. If it is a string, keep it, otherwise encode it.
	data, ok := v["program"]
	if !ok {
		return fmt.Errorf("missing program in json object")
	}
	program := ""
	if err := json.Unmarshal(data, &program); err != nil {
		if program, err = encodeProgram(data); err != nil {
			return err
		}
	}
	c.Program = program

	// process 'entry_points_by_type'
	data, ok = v["entry_points_by_type"]
	if !ok {
		return fmt.Errorf("missing entry_points_by_type in json object")
	}

	depEntryPointsByType := DeprecatedEntryPointsByType{}
	if err := json.Unmarshal(data, &depEntryPointsByType); err != nil {
		return err
	}
	c.DeprecatedEntryPointsByType = depEntryPointsByType

	// process 'abi'
	data, ok = v["abi"]
	if !ok {
		// contractClass can have an empty ABI for instance with ClassAt
		return nil
	}

	abis := []interface{}{}
	if err := json.Unmarshal(data, &abis); err != nil {
		return err
	}

	abiPointer := ABI{}
	for _, abi := range abis {
		if checkABI, ok := abi.(map[string]interface{}); ok {
			var ab ABIEntry
			abiType, ok := checkABI["type"].(string)
			if !ok {
				return fmt.Errorf("unknown abi type %v", checkABI["type"])
			}
			switch abiType {
			case string(ABITypeConstructor), string(ABITypeFunction), string(ABITypeL1Handler):
				ab = &FunctionABIEntry{}
			case string(ABITypeStruct):
				ab = &StructABIEntry{}
			case string(ABITypeEvent):
				ab = &EventABIEntry{}
			default:
				return fmt.Errorf("unknown ABI type %v", checkABI["type"])
			}
			data, err := json.Marshal(checkABI)
			if err != nil {
				return err
			}
			err = json.Unmarshal(data, ab)
			if err != nil {
				return err
			}
			abiPointer = append(abiPointer, ab)
		}
	}

	c.ABI = &abiPointer
	return nil
}

type SierraEntryPoint struct {
	// The index of the function in the program
	FunctionIdx int `json:"function_idx"`
	// A unique  identifier of the entry point (function) in the program
	Selector *felt.Felt `json:"selector"`
}

type EntryPointsByType struct {
	Constructor []SierraEntryPoint `json:"CONSTRUCTOR"`
	External    []SierraEntryPoint `json:"EXTERNAL"`
	L1Handler   []SierraEntryPoint `json:"L1_HANDLER"`
}

type ABIEntry interface {
	IsType() ABIType
}

type ABIType string

const (
	ABITypeConstructor ABIType = "constructor"
	ABITypeFunction    ABIType = "function"
	ABITypeL1Handler   ABIType = "l1_handler"
	ABITypeEvent       ABIType = "event"
	ABITypeStruct      ABIType = "struct"
)

type StructABIEntry struct {
	// The event type
	Type ABIType `json:"type"`

	// The event name
	Name string `json:"name"`

	// todo(minumum size should be 1)
	Size uint64 `json:"size"`

	Members []Member `json:"members"`
}

type Member struct {
	TypedParameter
	Offset int64 `json:"offset"`
}

type EventABIEntry struct {
	// The event type
	Type ABIType `json:"type"`

	// The event name
	Name string `json:"name"`

	Keys []TypedParameter `json:"keys"`

	Data []TypedParameter `json:"data"`
}

type FunctionStateMutability string

const (
	FuncStateMutVIEW FunctionStateMutability = "view"
)

type FunctionABIEntry struct {
	// The function type
	Type ABIType `json:"type"`

	// The function name
	Name string `json:"name"`

	StateMutability FunctionStateMutability `json:"stateMutability,omitempty"`

	Inputs []TypedParameter `json:"inputs"`

	Outputs []TypedParameter `json:"outputs"`
}

// IsType returns the ABIType of the StructABIEntry.
//
// Parameters:
//
//	none
//
// Returns:
// - ABIType: the ABIType
func (s *StructABIEntry) IsType() ABIType {
	return s.Type
}

// IsType returns the ABIType of the EventABIEntry.
//
// Parameters:
//
//	none
//
// Returns:
// - ABIType: the ABIType
func (e *EventABIEntry) IsType() ABIType {
	return e.Type
}

// IsType returns the ABIType of the FunctionABIEntry.
//
// Parameters:
//
//	none
//
// Returns:
// - ABIType: the ABIType
func (f *FunctionABIEntry) IsType() ABIType {
	return f.Type
}

type TypedParameter struct {
	// The parameter's name
	Name string `json:"name"`

	// The parameter's type
	Type string `json:"type"`
}

// encodeProgram encodes the content byte array using gzip compression and base64 encoding.
//
// It takes a content byte array as a parameter and returns the encoded program string and an error.
//
// Parameters:
// - content: byte array to be encoded
// Returns:
// - string: the encoded program
// - error: the error if any
func encodeProgram(content []byte) (string, error) {
	buf := bytes.NewBuffer(nil)
	gzipContent := gzip.NewWriter(buf)
	_, err := gzipContent.Write(content)
	if err != nil {
		return "", err
	}
	gzipContent.Close()
	program := base64.StdEncoding.EncodeToString(buf.Bytes())
	return program, nil
}

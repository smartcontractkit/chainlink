// Sourced from https://github.com/ethereum/go-ethereum/blob/fe91d476ba3e29316b6dc99b6efd4a571481d888/accounts/abi/selector_parser_test.go

// Copyright 2022 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package abi

import (
	"fmt"
	"log"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// https://docs.soliditylang.org/en/latest/grammar.html#a4.SolidityParser.parameterList
type parameter struct {
	typeName   string
	identifier string
}

func mkType(parameterList ...interface{}) []abi.ArgumentMarshaling {
	var result []abi.ArgumentMarshaling
	for i, p := range parameterList {
		name := fmt.Sprintf("name%d", i)

		if safeParameter, ok := p.(parameter); ok {
			name = safeParameter.identifier
			p = safeParameter.typeName
		}

		if typeName, ok := p.(string); ok {
			result = append(result, abi.ArgumentMarshaling{Name: name, Type: typeName, InternalType: typeName, Components: nil, Indexed: false})
		} else if components, ok := p.([]abi.ArgumentMarshaling); ok {
			result = append(result, abi.ArgumentMarshaling{Name: name, Type: "tuple", InternalType: "tuple", Components: components, Indexed: false})
		} else if components, ok := p.([][]abi.ArgumentMarshaling); ok {
			result = append(result, abi.ArgumentMarshaling{Name: name, Type: "tuple[]", InternalType: "tuple[]", Components: components[0], Indexed: false})
		} else {
			log.Fatalf("unexpected type %T", p)
		}
	}
	return result
}

func TestParseSelector(t *testing.T) {
	t.Parallel()

	type testCases struct {
		description    string
		input          string
		expectedOutput abi.SelectorMarshaling
	}

	for _, tc := range []testCases{
		{
			description: "No function args",
			input:       "noargs()",
			expectedOutput: abi.SelectorMarshaling{
				Name:   "noargs",
				Type:   "function",
				Inputs: []abi.ArgumentMarshaling{},
			},
		},
		{
			description: "Simple unnamed args",
			input:       "simple(uint256,uint256,uint256)",
			expectedOutput: abi.SelectorMarshaling{
				Name:   "simple",
				Type:   "function",
				Inputs: mkType("uint256", "uint256", "uint256"),
			},
		},
		{
			description: "Simple named args",
			input:       "simple(uint256 a, address b, byte c)",
			expectedOutput: abi.SelectorMarshaling{
				Name:   "simple",
				Type:   "function",
				Inputs: mkType(parameter{"uint256", "a"}, parameter{"address", "b"}, parameter{"byte", "c"}),
			},
		},
		// FAILING
		// {
		// 	description: "Extra whitespace",
		// 	input:       "simple    (uint256     a,     address b, byte c)",
		// 	expectedOutput: abi.SelectorMarshaling{
		// 		Name:   "simple",
		// 		Type:   "function",
		// 		Inputs: mkType(parameter{"uint256", "a"}, parameter{"address", "b"}, parameter{"byte", "c"}),
		// 	},
		// },
		{
			description: "Unnamed arrays",
			input:       "withArray(uint256[], address[2], uint8[4][][5])",
			expectedOutput: abi.SelectorMarshaling{
				Name:   "withArray",
				Type:   "function",
				Inputs: mkType("uint256[]", "address[2]", "uint8[4][][5]"),
			},
		},
		{
			description: "Named arrays",
			input:       "withArray(uint256[] a, address[2] b, uint8[4][][5] c)",
			expectedOutput: abi.SelectorMarshaling{
				Name: "withArray",
				Type: "function",
				Inputs: mkType(
					parameter{"uint256[]", "a"},
					parameter{"address[2]", "b"},
					parameter{"uint8[4][][5]", "c"},
				),
			},
		},
		{
			description: "Named tuple",
			input:       "addPerson((string name, uint16 age) person)",
			expectedOutput: abi.SelectorMarshaling{
				Name: "addPerson",
				Type: "function",
				Inputs: []abi.ArgumentMarshaling{
					{
						Name:         "person",
						Type:         "tuple",
						InternalType: "tuple",
						Components: []abi.ArgumentMarshaling{
							{
								Name:         "name",
								Type:         "string",
								InternalType: "string",
							},
							{
								Name:         "age",
								Type:         "uint16",
								InternalType: "uint16",
							},
						},
					},
				},
			},
		},
		// FAILING CASE
		// {
		// 	description: "Name tuple with explicit 'tuple' keyword",
		// 	input:       "addPerson(tuple(string name, uint16 age) person)",
		// 	expectedOutput: abi.SelectorMarshaling{
		// 		Name: "addPerson",
		// 		Type: "function",
		// 		Inputs: []abi.ArgumentMarshaling{
		// 			{
		// 				Name:         "person",
		// 				Type:         "tuple",
		// 				InternalType: "tuple",
		// 				Components: []abi.ArgumentMarshaling{
		// 					{
		// 						Name: "name", Type: "string", InternalType: "string", Components: nil, Indexed: false,
		// 					},
		// 					{
		// 						Name: "age", Type: "uint16", InternalType: "uint16", Components: nil, Indexed: false,
		// 					},
		// 				},
		// 			},
		// 		},
		// 	},
		// },
		{
			description: "Arrays of tuples",
			input:       "arrayNest((uint256,uint256)[],bytes32)",
			expectedOutput: abi.SelectorMarshaling{
				Name: "arrayNest",
				Type: "function",
				Inputs: mkType(
					[][]abi.ArgumentMarshaling{mkType("uint256", "uint256")},
					"bytes32",
				),
			},
		},
		{
			description: "Nested arrays of tuples",
			input:       "multiArrayNest((uint256,uint256)[],(uint256,uint256)[])",
			expectedOutput: abi.SelectorMarshaling{
				Name: "multiArrayNest",
				Type: "function",
				Inputs: mkType(
					[][]abi.ArgumentMarshaling{mkType("uint256", "uint256")},
					[][]abi.ArgumentMarshaling{mkType("uint256", "uint256")},
				),
			},
		},
		{
			description: "Nested tuples and arrays",
			input:       "multiNest(address,(uint256[],uint256),((address,bytes32),uint256))",
			expectedOutput: abi.SelectorMarshaling{
				Name: "multiNest",
				Type: "function",
				Inputs: mkType(
					"address",
					mkType("uint256[]", "uint256"),
					mkType(
						mkType("address", "bytes32"),
						"uint256",
					),
				),
			},
		},
		{
			description: "Combination of nested array of tuples and array type",
			input:       "singleArrayNestAndArray((uint256,uint256)[],bytes32[])",
			expectedOutput: abi.SelectorMarshaling{
				Name: "singleArrayNestAndArray",
				Type: "function",
				Inputs: mkType(
					[][]abi.ArgumentMarshaling{mkType("uint256", "uint256")},
					"bytes32[]",
				),
			},
		},
		{
			description: "Complex nesting with arrays and tuples",
			input:       "singleArrayNestWithArrayAndArray((uint256[],address[2],uint8[4][][5])[],bytes32[])",
			expectedOutput: abi.SelectorMarshaling{
				Name: "singleArrayNestWithArrayAndArray",
				Type: "function",
				Inputs: mkType(
					[][]abi.ArgumentMarshaling{mkType("uint256[]", "address[2]", "uint8[4][][5]")},
					"bytes32[]",
				),
			},
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			selector, err := ParseSelector(tc.input)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedOutput, selector)
		})
	}
}

func TestParseSelectorErrors(t *testing.T) {
	type errorTestCases struct {
		description   string
		input         string
		expectedError string
	}

	for _, scenario := range []errorTestCases{
		{
			description:   "invalid name",
			input:         "123()",
			expectedError: "failed to parse selector identifier '123()': invalid token start. Expected: [a-zA-Z], received: 1",
		},
		{
			description:   "missing closing parenthesis",
			input:         "noargs(",
			expectedError: "failed to parse selector args 'noargs(': expected ')', got ''",
		},
		{
			description:   "missing opening parenthesis",
			input:         "noargs)",
			expectedError: "failed to parse selector args 'noargs)': expected '(', got )",
		},
	} {
		t.Run(scenario.description, func(t *testing.T) {
			_, err := ParseSelector(scenario.input)
			require.Error(t, err)
			assert.Equal(t, scenario.expectedError, err.Error())
		})
	}
}

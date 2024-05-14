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
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

func TestParseSelector(t *testing.T) {
	t.Parallel()
	mkType := func(types ...interface{}) []abi.ArgumentMarshaling {
		var result []abi.ArgumentMarshaling
		for i, typeOrComponents := range types {
			name := fmt.Sprintf("name%d", i)
			if typeName, ok := typeOrComponents.(string); ok {
				result = append(result, abi.ArgumentMarshaling{Name: name, Type: typeName, InternalType: typeName, Components: nil, Indexed: false})
			} else if components, ok := typeOrComponents.([]abi.ArgumentMarshaling); ok {
				result = append(result, abi.ArgumentMarshaling{Name: name, Type: "tuple", InternalType: "tuple", Components: components, Indexed: false})
			} else if components, ok := typeOrComponents.([][]abi.ArgumentMarshaling); ok {
				result = append(result, abi.ArgumentMarshaling{Name: name, Type: "tuple[]", InternalType: "tuple[]", Components: components[0], Indexed: false})
			} else {
				log.Fatalf("unexpected type %T", typeOrComponents)
			}
		}
		return result
	}
	tests := []struct {
		input string
		name  string
		args  []abi.ArgumentMarshaling
	}{
		{"noargs()", "noargs", []abi.ArgumentMarshaling{}},
		{"simple(uint256,uint256,uint256)", "simple", mkType("uint256", "uint256", "uint256")},
		{"other(uint256,address)", "other", mkType("uint256", "address")},
		{"withArray(uint256[],address[2],uint8[4][][5])", "withArray", mkType("uint256[]", "address[2]", "uint8[4][][5]")},
		{"singleNest(bytes32,uint8,(uint256,uint256),address)", "singleNest", mkType("bytes32", "uint8", mkType("uint256", "uint256"), "address")},
		{"multiNest(address,(uint256[],uint256),((address,bytes32),uint256))", "multiNest",
			mkType("address", mkType("uint256[]", "uint256"), mkType(mkType("address", "bytes32"), "uint256"))},
		{"arrayNest((uint256,uint256)[],bytes32)", "arrayNest", mkType([][]abi.ArgumentMarshaling{mkType("uint256", "uint256")}, "bytes32")},
		{"multiArrayNest((uint256,uint256)[],(uint256,uint256)[])", "multiArrayNest",
			mkType([][]abi.ArgumentMarshaling{mkType("uint256", "uint256")}, [][]abi.ArgumentMarshaling{mkType("uint256", "uint256")})},
		{"singleArrayNestAndArray((uint256,uint256)[],bytes32[])", "singleArrayNestAndArray",
			mkType([][]abi.ArgumentMarshaling{mkType("uint256", "uint256")}, "bytes32[]")},
		{"singleArrayNestWithArrayAndArray((uint256[],address[2],uint8[4][][5])[],bytes32[])", "singleArrayNestWithArrayAndArray",
			mkType([][]abi.ArgumentMarshaling{mkType("uint256[]", "address[2]", "uint8[4][][5]")}, "bytes32[]")},
	}
	for i, tt := range tests {
		selector, err := ParseSelector(tt.input)
		if err != nil {
			t.Errorf("test %d: failed to parse selector '%v': %v", i, tt.input, err)
		}
		if selector.Name != tt.name {
			t.Errorf("test %d: unexpected function name: '%s' != '%s'", i, selector.Name, tt.name)
		}

		if selector.Type != "function" {
			t.Errorf("test %d: unexpected type: '%s' != '%s'", i, selector.Type, "function")
		}
		if !reflect.DeepEqual(selector.Inputs, tt.args) {
			t.Errorf("test %d: unexpected args: '%v' != '%v'", i, selector.Inputs, tt.args)
		}
	}
}

func TestParseSignature(t *testing.T) {
	t.Parallel()
	mkType := func(name string, typeOrComponents interface{}) abi.ArgumentMarshaling {
		if typeName, ok := typeOrComponents.(string); ok {
			return abi.ArgumentMarshaling{Name: name, Type: typeName, InternalType: typeName, Components: nil, Indexed: false}
		} else if components, ok := typeOrComponents.([]abi.ArgumentMarshaling); ok {
			return abi.ArgumentMarshaling{Name: name, Type: "tuple", InternalType: "tuple", Components: components, Indexed: false}
		} else if components, ok := typeOrComponents.([][]abi.ArgumentMarshaling); ok {
			return abi.ArgumentMarshaling{Name: name, Type: "tuple[]", InternalType: "tuple[]", Components: components[0], Indexed: false}
		}
		log.Fatalf("unexpected type %T", typeOrComponents)
		return abi.ArgumentMarshaling{}
	}
	tests := []struct {
		input string
		name  string
		args  []abi.ArgumentMarshaling
	}{
		{"noargs()", "noargs", []abi.ArgumentMarshaling{}},
		{"simple(a uint256, b uint256, c uint256)", "simple", []abi.ArgumentMarshaling{mkType("a", "uint256"), mkType("b", "uint256"), mkType("c", "uint256")}},
		{"other(foo uint256,    bar address)", "other", []abi.ArgumentMarshaling{mkType("foo", "uint256"), mkType("bar", "address")}},
		{"withArray(a uint256[], b address[2], c uint8[4][][5])", "withArray", []abi.ArgumentMarshaling{mkType("a", "uint256[]"), mkType("b", "address[2]"), mkType("c", "uint8[4][][5]")}},
		{"singleNest(d bytes32, e uint8, f (uint256,uint256), g address)", "singleNest", []abi.ArgumentMarshaling{mkType("d", "bytes32"), mkType("e", "uint8"), mkType("f", []abi.ArgumentMarshaling{mkType("name0", "uint256"), mkType("name1", "uint256")}), mkType("g", "address")}},
	}
	for i, tt := range tests {
		selector, err := ParseSignature(tt.input)
		if err != nil {
			t.Errorf("test %d: failed to parse selector '%v': %v", i, tt.input, err)
		}
		if selector.Name != tt.name {
			t.Errorf("test %d: unexpected function name: '%s' != '%s'", i, selector.Name, tt.name)
		}

		if selector.Type != "function" {
			t.Errorf("test %d: unexpected type: '%s' != '%s'", i, selector.Type, "function")
		}
		if !reflect.DeepEqual(selector.Inputs, tt.args) {
			t.Errorf("test %d: unexpected args: '%v' != '%v'", i, selector.Inputs, tt.args)
		}
	}
}

// Originally sourced from https://github.com/ethereum/go-ethereum/blob/fe91d476ba3e29316b6dc99b6efd4a571481d888/accounts/abi/selector_parser.go#L126
// Modified to suppor parsing selectors with argument names

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
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func isIdentifierSymbol(c byte) bool {
	return c == '$' || c == '_'
}

func parseToken(unescapedSelector string, isIdent bool) (string, string, error) {
	if len(unescapedSelector) == 0 {
		return "", "", errors.New("empty token")
	}
	firstChar := unescapedSelector[0]
	position := 1
	if !(isAlpha(firstChar) || (isIdent && isIdentifierSymbol(firstChar))) {
		return "", "", fmt.Errorf("invalid token start: %c", firstChar)
	}
	for position < len(unescapedSelector) {
		char := unescapedSelector[position]
		if !(isAlpha(char) || isDigit(char) || (isIdent && isIdentifierSymbol(char))) {
			break
		}
		position++
	}
	return unescapedSelector[:position], unescapedSelector[position:], nil
}

func parseIdentifier(unescapedSelector string) (string, string, error) {
	return parseToken(unescapedSelector, true)
}

func parseElementaryType(unescapedSelector string) (string, string, error) {
	parsedType, rest, err := parseToken(unescapedSelector, false)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse elementary type: %v", err)
	}
	// handle arrays
	for len(rest) > 0 && rest[0] == '[' {
		parsedType = parsedType + string(rest[0])
		rest = rest[1:]
		for len(rest) > 0 && isDigit(rest[0]) {
			parsedType = parsedType + string(rest[0])
			rest = rest[1:]
		}
		if len(rest) == 0 || rest[0] != ']' {
			return "", "", fmt.Errorf("failed to parse array: expected ']', got %c", unescapedSelector[0])
		}
		parsedType = parsedType + string(rest[0])
		rest = rest[1:]
	}
	return parsedType, rest, nil
}

func parseCompositeType(unescapedSelector string) ([]abi.ArgumentMarshaling, string, error) {
	if len(unescapedSelector) == 0 || unescapedSelector[0] != '(' {
		return nil, "", fmt.Errorf("expected '(', got %c", unescapedSelector[0])
	}
	rest := unescapedSelector[1:] // skip over the opening `(`
	result := []abi.ArgumentMarshaling{}
	var parsedType any
	var err error
	i := 0
	for len(rest) > 0 && rest[0] != ')' {
		// skip any leading whitespace
		for rest[0] == ' ' {
			rest = rest[1:]
		}

		parsedType, rest, err = parseType(rest[0:])
		if err != nil {
			return nil, "", fmt.Errorf("failed to parse type: %v", err)
		}

		// skip whitespace between name and identifier
		for rest[0] == ' ' {
			rest = rest[1:]
		}

		name := fmt.Sprintf("name%d", i)
		// if we're at a delimiter the parameter is unnamed
		if !(rest[0] == ',' || rest[0] == ')') {
			// attempt to parse name
			name, rest, err = parseIdentifier(rest[:])
			if err != nil {
				return nil, "", fmt.Errorf("failed to parse name: %v", err)
			}
		}

		arg, err := assembleArg(name, parsedType)
		if err != nil {
			return nil, "", fmt.Errorf("failed to parse type: %v", err)
		}

		result = append(result, arg)
		i++

		// skip trailing whitespace, consume comma
		for rest[0] == ' ' || rest[0] == ',' {
			rest = rest[1:]
		}
	}
	if len(rest) == 0 || rest[0] != ')' {
		return nil, "", fmt.Errorf("expected ')', got '%s'", rest)
	}
	if len(rest) >= 3 && rest[1] == '[' && rest[2] == ']' {
		// emits a sentinel value that later gets removed when assembling
		array, err := assembleArg("", "[]")
		if err != nil {
			panic("unreachable")
		}
		return append(result, array), rest[3:], nil
	}
	return result, rest[1:], nil
}

// type-name rule
func parseType(unescapedSelector string) (interface{}, string, error) {
	if len(unescapedSelector) == 0 {
		return nil, "", errors.New("empty type")
	}
	if unescapedSelector[0] == '(' {
		return parseCompositeType(unescapedSelector)
	}
	return parseElementaryType(unescapedSelector)
}

func assembleArg(name string, arg any) (abi.ArgumentMarshaling, error) {
	if s, ok := arg.(string); ok {
		return abi.ArgumentMarshaling{Name: name, Type: s, InternalType: s, Components: nil, Indexed: false}, nil
	} else if subArgs, ok := arg.([]abi.ArgumentMarshaling); ok {
		tupleType := "tuple"
		if len(subArgs) != 0 && subArgs[len(subArgs)-1].Type == "[]" {
			subArgs = subArgs[:len(subArgs)-1]
			tupleType = "tuple[]"
		}
		return abi.ArgumentMarshaling{Name: name, Type: tupleType, InternalType: tupleType, Components: subArgs, Indexed: false}, nil
	}
	return abi.ArgumentMarshaling{}, fmt.Errorf("failed to assemble args: unexpected type %T", arg)
}

// ParseSelector converts a method selector into a struct that can be JSON encoded
// and consumed by other functions in this package.
// Note, although uppercase letters are not part of the ABI spec, this function
// still accepts it as the general format is valid.
func ParseSelector(unescapedSelector string) (abi.SelectorMarshaling, error) {
	name, rest, err := parseIdentifier(unescapedSelector)
	if err != nil {
		return abi.SelectorMarshaling{}, fmt.Errorf("failed to parse selector identifier '%s': %v", unescapedSelector, err)
	}
	args := []abi.ArgumentMarshaling{}
	if len(rest) >= 2 && rest[0] == '(' && rest[1] == ')' {
		rest = rest[2:]
	} else {
		args, rest, err = parseCompositeType(rest)
		if err != nil {
			return abi.SelectorMarshaling{}, fmt.Errorf("failed to parse selector args '%s': %v", unescapedSelector, err)
		}
	}
	if len(rest) > 0 {
		return abi.SelectorMarshaling{}, fmt.Errorf("failed to parse selector '%s': unexpected string '%s'", unescapedSelector, rest)
	}
	return abi.SelectorMarshaling{Name: name, Type: "function", Inputs: args}, nil
}

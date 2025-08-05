package evm

import (
	"fmt"
	"math/big"
	"reflect"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

// buildMethodSignature creates a method signature string like "mint(address,uint256)" from ABI method
func (w *chainWriter) buildMethodSignature(abiMethod abi.Method) string {
	var inputTypes []string
	for _, input := range abiMethod.Inputs {
		inputTypes = append(inputTypes, input.Type.String())
	}
	return fmt.Sprintf("%s(%s)", abiMethod.Name, strings.Join(inputTypes, ","))
}

// convertArgsToTronParams converts Go struct args to Tron's flattened type-value format
// Example: struct{Address: "0x123", Amount: 100} -> []any{"address", "0x123", "uint256", "100"}
func (w *chainWriter) convertArgsToTronParams(abiMethod abi.Method, args any) ([]any, error) {
	var params []any

	// Handle the case where args is a slice/array of arguments
	argsValue := reflect.ValueOf(args)
	if argsValue.Kind() == reflect.Slice || argsValue.Kind() == reflect.Array {
		if argsValue.Len() != len(abiMethod.Inputs) {
			return nil, fmt.Errorf("argument count mismatch: got %d, expected %d", argsValue.Len(), len(abiMethod.Inputs))
		}

		for i, input := range abiMethod.Inputs {
			argValue := argsValue.Index(i)
			params = append(params, input.Type.String())

			// Convert the value to string format for Tron
			valueStr, err := w.convertValueToString(argValue.Interface(), input.Type)
			if err != nil {
				return nil, fmt.Errorf("failed to convert argument %d: %w", i, err)
			}
			params = append(params, valueStr)
		}
		return params, nil
	}

	// Handle the case where args is a struct
	if argsValue.Kind() == reflect.Struct {
		argsType := argsValue.Type()

		// We need to match struct fields to ABI inputs
		// This assumes the struct fields are in the same order as ABI inputs
		if argsValue.NumField() != len(abiMethod.Inputs) {
			return nil, fmt.Errorf("struct field count mismatch: got %d, expected %d", argsValue.NumField(), len(abiMethod.Inputs))
		}

		for i, input := range abiMethod.Inputs {
			field := argsValue.Field(i)
			params = append(params, input.Type.String())

			// Convert the value to string format for Tron
			valueStr, err := w.convertValueToString(field.Interface(), input.Type)
			if err != nil {
				return nil, fmt.Errorf("failed to convert field %s: %w", argsType.Field(i).Name, err)
			}
			params = append(params, valueStr)
		}
		return params, nil
	}

	return nil, fmt.Errorf("unsupported args type: %v", argsValue.Kind())
}

// convertValueToString converts a Go value to string format for Tron
func (w *chainWriter) convertValueToString(value any, abiType abi.Type) (string, error) {
	switch abiType.T {
	case abi.AddressTy:
		if addr, ok := value.(common.Address); ok {
			return addr.Hex(), nil
		}
		if str, ok := value.(string); ok {
			return str, nil
		}
		return "", fmt.Errorf("invalid address type: %T", value)

	case abi.UintTy, abi.IntTy:
		if bigInt, ok := value.(*big.Int); ok {
			return bigInt.String(), nil
		}
		// Handle various integer types
		val := reflect.ValueOf(value)
		switch val.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return fmt.Sprintf("%d", val.Int()), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return fmt.Sprintf("%d", val.Uint()), nil
		}
		return fmt.Sprintf("%v", value), nil

	case abi.StringTy:
		if str, ok := value.(string); ok {
			return str, nil
		}
		return fmt.Sprintf("%v", value), nil

	case abi.BoolTy:
		if b, ok := value.(bool); ok {
			return fmt.Sprintf("%t", b), nil
		}
		return fmt.Sprintf("%v", value), nil

	case abi.BytesTy, abi.FixedBytesTy:
		if bytes, ok := value.([]byte); ok {
			return fmt.Sprintf("0x%x", bytes), nil
		}
		if str, ok := value.(string); ok {
			return str, nil
		}
		return fmt.Sprintf("%v", value), nil

	default:
		// For other types, convert to string
		return fmt.Sprintf("%v", value), nil
	}
}

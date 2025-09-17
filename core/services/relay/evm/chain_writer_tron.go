package evm

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

// buildMethodSignature creates a method signature string like "mint(address,uint256)" from ABI method
func (w *chainWriter) buildMethodSignature(abiMethod abi.Method) string {
	inputTypes := make([]string, 0, len(abiMethod.Inputs))
	for _, input := range abiMethod.Inputs {
		inputTypes = append(inputTypes, input.Type.String())
	}
	return fmt.Sprintf("%s(%s)", abiMethod.Name, strings.Join(inputTypes, ","))
}

// convertArgsToTronParams converts Go args to Tron's parameter format
// Simply extracts the raw values and lets the Tron SDK handle type conversion
func (w *chainWriter) convertArgsToTronParams(abiMethod abi.Method, args any) ([]any, error) {
	argsValue := reflect.ValueOf(args)

	// Handle slice/array of arguments
	if argsValue.Kind() == reflect.Slice || argsValue.Kind() == reflect.Array {
		if argsValue.Len() != len(abiMethod.Inputs) {
			return nil, fmt.Errorf("argument count mismatch: got %d, expected %d", argsValue.Len(), len(abiMethod.Inputs))
		}

		params := make([]any, 0, argsValue.Len()*2)
		for i, input := range abiMethod.Inputs {
			argValue := argsValue.Index(i).Interface()
			w.logger.Debugf("Arg %d: type=%s, value type=%T, value=%+v", i, input.Type.String(), argValue, argValue)
			params = append(params, input.Type.String(), argValue)
		}

		w.logger.Debugf("Final Tron params: %+v", params)
		return params, nil
	}

	// Handle struct arguments
	if argsValue.Kind() == reflect.Struct {
		if argsValue.NumField() != len(abiMethod.Inputs) {
			return nil, fmt.Errorf("struct field count mismatch: got %d, expected %d", argsValue.NumField(), len(abiMethod.Inputs))
		}

		params := make([]any, 0, argsValue.NumField()*2)
		for i, input := range abiMethod.Inputs {
			fieldValue := argsValue.Field(i).Interface()
			w.logger.Debugf("Field %d: type=%s, value type=%T, value=%+v", i, input.Type.String(), fieldValue, fieldValue)
			params = append(params, input.Type.String(), fieldValue)
		}

		w.logger.Debugf("Final Tron params: %+v", params)
		return params, nil
	}

	return nil, fmt.Errorf("unsupported args type: %v", argsValue.Kind())
}

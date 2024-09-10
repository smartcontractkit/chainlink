package evm

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"reflect"
)

// Returns empty string if type is not primitive
func SolidityPrimitiveTypeToGoType(abiType abi.Type) string {
	if abiType.TupleType != nil {
		return ""
	}
	switch abiType.String() {
	case "uint256":
		return reflect.TypeOf(uint64(0)).String() // Go's uint64 can hold values up to 2^64-1, which is a subset of uint256
	case "uint8":
		return reflect.TypeOf(uint8(0)).String()
	case "uint16":
		return reflect.TypeOf(uint16(0)).String()
	case "uint32":
		return reflect.TypeOf(uint32(0)).String()
	case "uint64":
		return reflect.TypeOf(uint64(0)).String()
	case "int256":
		return reflect.TypeOf(int64(0)).String() // Go's int64 can hold values up to 2^63-1, a subset of int256
	case "int8":
		return reflect.TypeOf(int8(0)).String()
	case "int16":
		return reflect.TypeOf(int16(0)).String()
	case "int32":
		return reflect.TypeOf(int32(0)).String()
	case "int64":
		return reflect.TypeOf(int64(0)).String()
	case "bool":
		return reflect.TypeOf(false).String()
	case "string":
		return reflect.TypeOf("").String()
	case "address":
		return reflect.TypeOf(common.Address{}).String()
	case "address[]":
		return reflect.TypeOf([]common.Address{}).String()
	case "bytes2":
		return reflect.TypeOf([2]byte{}).String()
		return reflect.TypeOf([]string{}).String()
	case "bytes":
		return reflect.TypeOf([]byte{}).String()
	case "bytes32":
		return reflect.TypeOf([32]byte{}).String() // Fixed-size byte arrays
		//TODO fix, just for POC
	case "uint[32]":
		return reflect.TypeOf([]uint8{}).String()
	case "uint8[32]":
		return reflect.TypeOf([32]byte{}).String()
	case "int192":
		return reflect.TypeOf(big.Int{}).String()
	default:
		//TODO add missing cases
		return "string" // Only other option is struct types default to ""
	}
}

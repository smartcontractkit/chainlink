package types

import (
	"reflect"

	"github.com/ethereum/go-ethereum/common"
)

//go:generate go run ./gen/main.go

var typeMap = map[string]*ABIEncodingType{
	"bool": {
		native:  reflect.TypeOf(true),
		checked: reflect.TypeOf(true),
	},
	"int8": {
		native:  reflect.TypeOf(int8(0)),
		checked: reflect.TypeOf(int8(0)),
	},
	"int16": {
		native:  reflect.TypeOf(int16(0)),
		checked: reflect.TypeOf(int16(0)),
	},
	"int32": {
		native:  reflect.TypeOf(int32(0)),
		checked: reflect.TypeOf(int32(0)),
	},
	"int64": {
		native:  reflect.TypeOf(int64(0)),
		checked: reflect.TypeOf(int64(0)),
	},
	"uint8": {
		native:  reflect.TypeOf(uint8(0)),
		checked: reflect.TypeOf(uint8(0)),
	},
	"uint16": {
		native:  reflect.TypeOf(uint16(0)),
		checked: reflect.TypeOf(uint16(0)),
	},
	"uint32": {
		native:  reflect.TypeOf(uint32(0)),
		checked: reflect.TypeOf(uint32(0)),
	},
	"uint64": {
		native:  reflect.TypeOf(uint64(0)),
		checked: reflect.TypeOf(uint64(0)),
	},
	"string": {
		native:  reflect.TypeOf(""),
		checked: reflect.TypeOf(""),
	},
	"address": {
		native:  reflect.TypeOf(common.Address{}),
		checked: reflect.TypeOf(common.Address{}),
	},
	"bytes": {
		native:  reflect.TypeOf([]byte{}),
		checked: reflect.TypeOf([]byte{}),
	},
}

type ABIEncodingType struct {
	native  reflect.Type
	checked reflect.Type
}

func GetAbiEncodingType(name string) (*ABIEncodingType, bool) {
	abiType, ok := typeMap[name]
	return abiType, ok
}

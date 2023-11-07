package types

import (
	"reflect"
)

//go:generate go run gen.go

var typeMap = map[string]*abiEncodingType{
	"bool": {
		Native:  reflect.TypeOf(true),
		Checked: reflect.TypeOf(true),
	},
	"int8": {
		Native:  reflect.TypeOf(int8(0)),
		Checked: reflect.TypeOf(int8(0)),
	},
	"int16": {
		Native:  reflect.TypeOf(int16(0)),
		Checked: reflect.TypeOf(int16(0)),
	},
	"int32": {
		Native:  reflect.TypeOf(int32(0)),
		Checked: reflect.TypeOf(int32(0)),
	},
	"int64": {
		Native:  reflect.TypeOf(int64(0)),
		Checked: reflect.TypeOf(int64(0)),
	},
	"uint8": {
		Native:  reflect.TypeOf(uint8(0)),
		Checked: reflect.TypeOf(uint8(0)),
	},
	"uint16": {
		Native:  reflect.TypeOf(uint16(0)),
		Checked: reflect.TypeOf(uint16(0)),
	},
	"uint32": {
		Native:  reflect.TypeOf(uint32(0)),
		Checked: reflect.TypeOf(uint32(0)),
	},
	"uint64": {
		Native:  reflect.TypeOf(uint64(0)),
		Checked: reflect.TypeOf(uint64(0)),
	},
	"string": {
		Native:  reflect.TypeOf(""),
		Checked: reflect.TypeOf(""),
	},
}

type abiEncodingType struct {
	Native  reflect.Type
	Checked reflect.Type
}

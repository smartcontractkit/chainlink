// Copyright 2021 github.com/gagliardetto
// This file has been modified by github.com/gagliardetto
//
// Copyright 2020 dfuse Platform Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package bin

import (
	"errors"
	"fmt"
	"reflect"
	"sort"

	"go.uber.org/zap"
)

func (e *Encoder) encodePrimitive(rv reflect.Value, opt *option) (isPrimitive bool, err error) {
	isPrimitive = true
	switch rv.Kind() {
	// case reflect.Int:
	// 	err = e.WriteInt64(rv.Int(), LE)
	// case reflect.Uint:
	// 	err = e.WriteUint64(rv.Uint(), LE)
	case reflect.String:
		err = e.WriteString(rv.String())
	case reflect.Uint8:
		err = e.WriteByte(byte(rv.Uint()))
	case reflect.Int8:
		err = e.WriteByte(byte(rv.Int()))
	case reflect.Int16:
		err = e.WriteInt16(int16(rv.Int()), LE)
	case reflect.Uint16:
		err = e.WriteUint16(uint16(rv.Uint()), LE)
	case reflect.Int32:
		err = e.WriteInt32(int32(rv.Int()), LE)
	case reflect.Uint32:
		err = e.WriteUint32(uint32(rv.Uint()), LE)
	case reflect.Uint64:
		err = e.WriteUint64(rv.Uint(), LE)
	case reflect.Int64:
		err = e.WriteInt64(rv.Int(), LE)
	case reflect.Float32:
		err = e.WriteFloat32(float32(rv.Float()), LE)
	case reflect.Float64:
		err = e.WriteFloat64(rv.Float(), LE)
	case reflect.Bool:
		err = e.WriteBool(rv.Bool())
	default:
		isPrimitive = false
	}
	return
}

func (e *Encoder) encodeBorsh(rv reflect.Value, opt *option) (err error) {
	if opt == nil {
		opt = newDefaultOption()
	}
	e.currentFieldOpt = opt

	if traceEnabled {
		zlog.Debug("encode: type",
			zap.Stringer("value_kind", rv.Kind()),
			zap.Reflect("options", opt),
		)
	}

	if opt.isOptional() {
		if rv.IsZero() {
			if traceEnabled {
				zlog.Debug("encode: skipping optional value with", zap.Stringer("type", rv.Kind()))
			}
			return e.WriteBool(false)
		}
		err := e.WriteBool(true)
		if err != nil {
			return err
		}
		// The optionality has been used; stop its propagation:
		opt.setIsOptional(false)
	}
	// Reset optionality so it won't propagate to child types:
	opt = opt.clone().setIsOptional(false)

	if isZero(rv) {
		return nil
	}

	if marshaler, ok := rv.Interface().(BinaryMarshaler); ok {
		if rv.Kind() == reflect.Ptr && rv.IsZero() {
			return nil
		}
		if traceEnabled {
			zlog.Debug("encode: using MarshalerBinary method to encode type")
		}
		return marshaler.MarshalWithEncoder(e)
	}

	// Encode the value if it's a primitive type
	isPrimitive, err := e.encodePrimitive(rv, nil)
	if isPrimitive {
		return err
	}

	switch rv.Kind() {
	case reflect.Ptr:
		if rv.IsNil() {
			el := reflect.New(rv.Type().Elem()).Elem()
			return e.encodeBorsh(el, nil)
		} else {
			return e.encodeBorsh(rv.Elem(), nil)
		}
	case reflect.Interface:
		// skip
		return nil
	}

	if !rv.IsZero() && !reflect.Indirect(rv).IsZero() {
		rv = reflect.Indirect(rv)
	}
	rt := rv.Type()
	switch rt.Kind() {
	case reflect.Array:
		l := rt.Len()
		if traceEnabled {
			defer func(prev *zap.Logger) { zlog = prev }(zlog)
			zlog = zlog.Named("array")
			zlog.Debug("encode: array", zap.Int("length", l), zap.Stringer("type", rv.Kind()))
		}

		switch k := rv.Type().Elem().Kind(); k {
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			// if it's a [n]byte, accumulate and write in one command:
			if err := reflect_writeArrayOfUint_(e, l, k, rv, LE); err != nil {
				return err
			}
		default:
			for i := 0; i < l; i++ {
				if err = e.encodeBorsh(rv.Index(i), nil); err != nil {
					return
				}
			}
		}
	case reflect.Slice:
		var l int
		if opt.hasSizeOfSlice() {
			l = opt.getSizeOfSlice()
			if traceEnabled {
				zlog.Debug("encode: slice with sizeof set", zap.Int("size_of", l))
			}
		} else {
			l = rv.Len()
			if err = e.WriteUint32(uint32(l), LE); err != nil {
				return
			}
		}
		if traceEnabled {
			defer func(prev *zap.Logger) { zlog = prev }(zlog)
			zlog = zlog.Named("slice")
			zlog.Debug("encode: slice", zap.Int("length", l), zap.Stringer("type", rv.Kind()))
		}

		// we would want to skip to the correct head_offset

		switch k := rv.Type().Elem().Kind(); k {
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			// if it's a [n]byte, accumulate and write in one command:
			if err := reflect_writeArrayOfUint_(e, l, k, rv, LE); err != nil {
				return err
			}
		default:
			for i := 0; i < l; i++ {
				if err = e.encodeBorsh(rv.Index(i), nil); err != nil {
					return
				}
			}
		}

	case reflect.Struct:
		if err = e.encodeStructBorsh(rt, rv); err != nil {
			return
		}

	case reflect.Map:
		keys := rv.MapKeys()
		sort.Slice(keys, vComp(keys))

		keyCount := rv.Len()
		if traceEnabled {
			zlog.Debug("encode: map",
				zap.Int("key_count", keyCount),
				zap.String("key_type", rt.String()),
				typeField("value_type", rv),
			)
			defer func(prev *zap.Logger) { zlog = prev }(zlog)
			zlog = zlog.Named("struct")
		}

		if err = e.WriteUint32(uint32(keyCount), LE); err != nil {
			return
		}

		for _, mapKey := range keys {
			if err = e.Encode(mapKey.Interface()); err != nil {
				return
			}

			if err = e.Encode(rv.MapIndex(mapKey).Interface()); err != nil {
				return
			}
		}
	// TODO:
	// case reflect.Ptr:
	// 	if rv.IsNil() {
	// 	} else {
	// 		return e.encodeBorsh(rv.Elem(), opt)
	// 	}
	default:
		return fmt.Errorf("encode: unsupported type %q", rt)
	}
	return
}

func (enc *Encoder) encodeComplexEnumBorsh(rv reflect.Value) error {
	t := rv.Type()
	enum := BorshEnum(rv.Field(0).Uint())
	// write enum identifier
	if err := enc.WriteByte(byte(enum)); err != nil {
		return err
	}
	// write enum field, if necessary
	if int(enum)+1 >= t.NumField() {
		return errors.New("complex enum too large")
	}
	// Enum is empty
	field := rv.Field(int(enum) + 1)
	if field.Kind() == reflect.Ptr {
		field = field.Elem()
	}
	if field.Kind() == reflect.Struct {
		return enc.encodeStructBorsh(field.Type(), field)
	}
	// Encode the value if it's a primitive type
	isPrimitive, err := enc.encodePrimitive(field, nil)
	if isPrimitive {
		return err
	}
	return nil
}

type BorshEnum uint8

// EmptyVariant is an empty borsh enum variant.
type EmptyVariant struct{}

func (_ *EmptyVariant) MarshalWithEncoder(_ *Encoder) error {
	return nil
}

func (_ *EmptyVariant) UnmarshalWithDecoder(_ *Decoder) error {
	return nil
}

func (e *Encoder) encodeStructBorsh(rt reflect.Type, rv reflect.Value) (err error) {
	l := rv.NumField()

	if traceEnabled {
		zlog.Debug("encode: struct", zap.Int("fields", l), zap.Stringer("type", rv.Kind()))
	}

	// Handle complex enum:
	if rt.NumField() > 0 {
		// If the first field has type BorshEnum and is flagged with "borsh_enum"
		// we have a complex enum:
		firstField := rt.Field(0)
		if isTypeBorshEnum(firstField.Type) &&
			parseFieldTag(firstField.Tag).IsBorshEnum {
			return e.encodeComplexEnumBorsh(rv)
		}
	}

	sizeOfMap := map[string]int{}
	for i := 0; i < l; i++ {
		structField := rt.Field(i)
		fieldTag := parseFieldTag(structField.Tag)

		if fieldTag.Skip {
			if traceEnabled {
				zlog.Debug("encode: skipping struct field with skip flag",
					zap.String("struct_field_name", structField.Name),
				)
			}
			continue
		}

		rv := rv.Field(i)

		if fieldTag.SizeOf != "" {
			if traceEnabled {
				zlog.Debug("encode: struct field has sizeof tag",
					zap.String("sizeof_field_name", fieldTag.SizeOf),
					zap.String("struct_field_name", structField.Name),
				)
			}
			sizeOfMap[fieldTag.SizeOf] = sizeof(structField.Type, rv)
		}

		if !rv.CanInterface() {
			if traceEnabled {
				zlog.Debug("encode:  skipping field: unable to interface field, probably since field is not exported",
					zap.String("sizeof_field_name", fieldTag.SizeOf),
					zap.String("struct_field_name", structField.Name),
				)
			}
			continue
		}

		option := &option{
			OptionalField: fieldTag.Optional,
			Order:         fieldTag.Order,
		}

		if s, ok := sizeOfMap[structField.Name]; ok {
			if traceEnabled {
				zlog.Debug("setting sizeof option", zap.String("of", structField.Name), zap.Int("size", s))
			}
			option.setSizeOfSlice(s)
		}

		if traceEnabled {
			zlog.Debug("encode: struct field",
				zap.Stringer("struct_field_value_type", rv.Kind()),
				zap.String("struct_field_name", structField.Name),
				zap.Reflect("struct_field_tags", fieldTag),
				zap.Reflect("struct_field_option", option),
			)
		}

		if err := e.encodeBorsh(rv, option); err != nil {
			return fmt.Errorf("error while encoding %q field: %w", structField.Name, err)
		}
	}
	return nil
}

func vComp(keys []reflect.Value) func(int, int) bool {
	return func(i int, j int) bool {
		a, b := keys[i], keys[j]
		if a.Kind() == reflect.Interface {
			a = a.Elem()
			b = b.Elem()
		}
		switch a.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
			return a.Int() < b.Int()
		case reflect.Int64:
			return a.Interface().(int64) < b.Interface().(int64)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32:
			return a.Uint() < b.Uint()
		case reflect.Uint64:
			return a.Interface().(uint64) < b.Interface().(uint64)
		case reflect.Float32, reflect.Float64:
			return a.Float() < b.Float()
		case reflect.String:
			return a.String() < b.String()
		}
		panic("unsupported key compare")
	}
}

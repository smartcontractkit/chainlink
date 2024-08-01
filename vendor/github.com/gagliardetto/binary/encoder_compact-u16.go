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
	"fmt"
	"reflect"

	"go.uber.org/zap"
)

func (e *Encoder) encodeCompactU16(rv reflect.Value, opt *option) (err error) {
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

	if opt.is_Optional() {
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
		opt.set_Optional(false)
	}

	if isZero(rv) {
		return nil
	}

	if marshaler, ok := rv.Interface().(BinaryMarshaler); ok {
		if traceEnabled {
			zlog.Debug("encode: using MarshalerBinary method to encode type")
		}
		return marshaler.MarshalWithEncoder(e)
	}

	switch rv.Kind() {
	case reflect.String:
		return e.WriteString(rv.String())
	case reflect.Uint8:
		return e.WriteByte(byte(rv.Uint()))
	case reflect.Int8:
		return e.WriteByte(byte(rv.Int()))
	case reflect.Int16:
		return e.WriteInt16(int16(rv.Int()), opt.Order)
	case reflect.Uint16:
		return e.WriteUint16(uint16(rv.Uint()), opt.Order)
	case reflect.Int32:
		return e.WriteInt32(int32(rv.Int()), opt.Order)
	case reflect.Uint32:
		return e.WriteUint32(uint32(rv.Uint()), opt.Order)
	case reflect.Uint64:
		return e.WriteUint64(rv.Uint(), opt.Order)
	case reflect.Int64:
		return e.WriteInt64(rv.Int(), opt.Order)
	case reflect.Float32:
		return e.WriteFloat32(float32(rv.Float()), opt.Order)
	case reflect.Float64:
		return e.WriteFloat64(rv.Float(), opt.Order)
	case reflect.Bool:
		return e.WriteBool(rv.Bool())
	case reflect.Ptr:
		return e.encodeCompactU16(rv.Elem(), opt)
	case reflect.Interface:
		// skip
		return nil
	}

	rv = reflect.Indirect(rv)
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
				if err = e.encodeCompactU16(rv.Index(i), nil); err != nil {
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
			if err = e.WriteCompactU16Length(l); err != nil {
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
				if err = e.encodeCompactU16(rv.Index(i), nil); err != nil {
					return
				}
			}
		}
	case reflect.Struct:
		if err = e.encodeStructCompactU16(rt, rv); err != nil {
			return
		}

	case reflect.Map:
		keyCount := len(rv.MapKeys())

		if traceEnabled {
			zlog.Debug("encode: map",
				zap.Int("key_count", keyCount),
				zap.String("key_type", rt.String()),
				typeField("value_type", rv.Elem()),
			)
			defer func(prev *zap.Logger) { zlog = prev }(zlog)
			zlog = zlog.Named("struct")
		}

		if err = e.WriteCompactU16Length(keyCount); err != nil {
			return
		}

		for _, mapKey := range rv.MapKeys() {
			if err = e.Encode(mapKey.Interface()); err != nil {
				return
			}

			if err = e.Encode(rv.MapIndex(mapKey).Interface()); err != nil {
				return
			}
		}

	default:
		return fmt.Errorf("encode: unsupported type %q", rt)
	}
	return
}

func (e *Encoder) encodeStructCompactU16(rt reflect.Type, rv reflect.Value) (err error) {
	l := rv.NumField()

	if traceEnabled {
		zlog.Debug("encode: struct", zap.Int("fields", l), zap.Stringer("type", rv.Kind()))
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
			is_OptionalField: fieldTag.Option,
			Order:            fieldTag.Order,
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

		if err := e.encodeCompactU16(rv, option); err != nil {
			return fmt.Errorf("error while encoding %q field: %w", structField.Name, err)
		}
	}
	return nil
}

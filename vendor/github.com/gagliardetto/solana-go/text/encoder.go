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

package text

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/fatih/color"
)

var (
	EncoderColorCyan   = color.New(color.FgCyan)
	EncoderColorYellow = color.New(color.FgYellow)
	EncoderColorGreen  = color.New(color.FgGreen)
	EncoderColorWhite  = color.New(color.FgWhite)
)

type Encoder struct {
	output      io.Writer
	indentLevel int
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		output: w,
	}
}

func isZero(rv reflect.Value) (b bool) {
	return rv.Kind() == 0
}

func isNil(rv reflect.Value) (b bool) {
	defer func(b bool) {
		if err := recover(); err != nil {
			b = true
		}
	}(b)
	return rv.IsNil()
}

func (e *Encoder) Encode(v interface{}, option *Option) (err error) {
	if option == nil {
		option = &Option{}
	}
	return e.encode(reflect.ValueOf(v), option)
}

type Option struct {
	indent     bool
	linear     bool
	fgColor    *color.Color
	NoTypeName bool
}

func (e *Encoder) encode(rv reflect.Value, option *Option) (err error) {
	if option == nil {
		option = &Option{}
	}

	if isZero(rv) {
		return e.ToWriter("NIL_VALUE", option.indent, option.fgColor)
	}

	if enc, ok := rv.Interface().(TextEncodable); ok {
		return enc.TextEncode(e, option)
	}

	//if sr, ok := rv.Interface().(fmt.Stringer); ok {
	//	return e.ToWriter(sr.String(), option.indent, option.fgColor)
	//}

	switch rv.Kind() {
	case reflect.String:
		return e.ToWriter(rv.String(), option.indent, option.fgColor)
	case reflect.Uint8, reflect.Int8:
		return e.ToWriter(fmt.Sprintf("%d", byte(rv.Uint())), option.indent, option.fgColor)
	case reflect.Int16:
		return e.ToWriter(fmt.Sprintf("%d", int16(rv.Int())), option.indent, option.fgColor)
	case reflect.Uint16:
		return e.ToWriter(fmt.Sprintf("%d", uint16(rv.Uint())), option.indent, option.fgColor)
	case reflect.Int32:
		return e.ToWriter(fmt.Sprintf("%d", int32(rv.Int())), option.indent, option.fgColor)
	case reflect.Uint32:
		return e.ToWriter(fmt.Sprintf("%d", uint32(rv.Uint())), option.indent, option.fgColor)
	case reflect.Uint64:
		return e.ToWriter(fmt.Sprintf("%d", rv.Uint()), option.indent, option.fgColor)
	case reflect.Int64:
		return e.ToWriter(fmt.Sprintf("%d", rv.Int()), option.indent, option.fgColor)
	case reflect.Float32:
		return e.ToWriter(fmt.Sprintf("%f", float32(rv.Float())), option.indent, option.fgColor)
	case reflect.Float64:
		return e.ToWriter(fmt.Sprintf("%f", rv.Float()), option.indent, option.fgColor)
	case reflect.Bool:
		return e.ToWriter(fmt.Sprintf("%t", rv.Bool()), option.indent, option.fgColor)
	case reflect.Ptr:
		return e.encode(rv.Elem(), option)
	}

	rv = reflect.Indirect(rv)
	rt := rv.Type()
	switch rt.Kind() {
	case reflect.Array:
		l := rt.Len()
		e.indentLevel++
		for i := 0; i < l; i++ {
			if err := e.ToWriter("\n", false, nil); err != nil {
				return err
			}
			if err := e.ToWriter(fmt.Sprintf("[%d] ", i), true, nil); err != nil {
				return err
			}
			if err = e.encode(rv.Index(i), option); err != nil {
				return
			}
		}
		e.indentLevel--
	case reflect.Slice:
		l := rv.Len()
		e.indentLevel++
		for i := 0; i < l; i++ {

			if err := e.ToWriter("\n", false, nil); err != nil {
				return err
			}

			if err := e.ToWriter(fmt.Sprintf("[%d] ", i), true, nil); err != nil {
				return err
			}

			if err = e.encode(rv.Index(i), option); err != nil {
				return
			}
			//if err := e.ToWriter("\n"); err != nil {
			//	return err
			//}
		}
		e.indentLevel--
	case reflect.Struct:
		if err = e.encodeStruct(rt, rv, option); err != nil {
			return
		}

	case reflect.Map:
		for _, mapKey := range rv.MapKeys() {
			if err = e.Encode(mapKey.Interface(), option); err != nil {
				return
			}

			if err = e.Encode(rv.MapIndex(mapKey).Interface(), option); err != nil {
				return
			}
		}

	default:
		return e.ToWriter("NOT TEXTABLE", false, option.fgColor)
	}
	return
}

func (e *Encoder) ToWriter(s string, indent bool, c *color.Color) (err error) {
	if indent {
		indent := strings.Repeat(" ", e.indentLevel)
		if _, err = e.output.Write([]byte(indent)); err != nil {
			return err
		}
	}
	if c != nil {
		s = c.Sprintf("%s", s)
	}
	_, err = e.output.Write([]byte(s))

	return nil
}

func (e *Encoder) encodeStruct(rt reflect.Type, rv reflect.Value, option *Option) (err error) {
	e.indentLevel++
	defer func() {
		e.indentLevel--
	}()

	if err := e.ToWriter("\n", false, nil); err != nil {
		return err
	}

	if !option.NoTypeName {
		if err := e.ToWriter(rt.Name(), true, EncoderColorCyan); err != nil {
			return err
		}
	}

	if !option.linear {
		if err := e.ToWriter("\n", false, nil); err != nil {
			return err
		}
	} else {
		if err := e.ToWriter(" ", false, nil); err != nil {
			return err
		}
	}

	l := rv.NumField()
	for i := 0; i < l; i++ {
		structField := rt.Field(i)
		fieldTag := parseFieldTag(structField.Tag)
		fieldOption := &Option{
			fgColor:    EncoderColorWhite,
			linear:     fieldTag.Linear,
			NoTypeName: false,
		}

		fieldOption.linear = fieldTag.Linear

		if fieldTag.Skip {
			continue
		}

		rv := rv.Field(i)

		if !rv.CanInterface() {
			continue
		}

		if err := e.ToWriter(structField.Name+": ", !option.linear, EncoderColorGreen); err != nil {
			return err
		}

		if err := e.encode(rv, fieldOption); err != nil {
			return err
		}

		if !option.linear {
			if err := e.ToWriter("\n", false, nil); err != nil {
				return err
			}
		} else {
			if err := e.ToWriter(" ", false, nil); err != nil {
				return err
			}
		}
	}
	return nil
}

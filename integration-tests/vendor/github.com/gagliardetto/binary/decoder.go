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
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math"
	"reflect"
	"strings"
	"unicode/utf8"

	"go.uber.org/zap"
)

var TypeSize = struct {
	Bool int
	Byte int

	Int8  int
	Int16 int

	Uint8   int
	Uint16  int
	Uint32  int
	Uint64  int
	Uint128 int

	Float32 int
	Float64 int

	PublicKey int
	Signature int

	Tstamp         int
	BlockTimestamp int

	CurrencyName int
}{
	Byte: 1,
	Bool: 1,

	Int8:  1,
	Int16: 2,

	Uint8:   1,
	Uint16:  2,
	Uint32:  4,
	Uint64:  8,
	Uint128: 16,

	Float32: 4,
	Float64: 8,
}

// Decoder implements the EOS unpacking, similar to FC_BUFFER
type Decoder struct {
	data []byte
	pos  int

	currentFieldOpt *option

	encoding Encoding
}

func (dec *Decoder) IsBorsh() bool {
	return dec.encoding.IsBorsh()
}

func (dec *Decoder) IsBin() bool {
	return dec.encoding.IsBin()
}

func (dec *Decoder) IsCompactU16() bool {
	return dec.encoding.IsCompactU16()
}

func NewDecoderWithEncoding(data []byte, enc Encoding) *Decoder {
	if !isValidEncoding(enc) {
		panic(fmt.Sprintf("provided encoding is not valid: %s", enc))
	}
	return &Decoder{
		data:     data,
		encoding: enc,
	}
}

func NewBinDecoder(data []byte) *Decoder {
	return NewDecoderWithEncoding(data, EncodingBin)
}

func NewBorshDecoder(data []byte) *Decoder {
	return NewDecoderWithEncoding(data, EncodingBorsh)
}

func NewCompactU16Decoder(data []byte) *Decoder {
	return NewDecoderWithEncoding(data, EncodingCompactU16)
}

func (dec *Decoder) Decode(v interface{}) (err error) {
	switch dec.encoding {
	case EncodingBin:
		return dec.decodeWithOptionBin(v, nil)
	case EncodingBorsh:
		return dec.decodeWithOptionBorsh(v, nil)
	case EncodingCompactU16:
		return dec.decodeWithOptionCompactU16(v, nil)
	default:
		panic(fmt.Errorf("encoding not implemented: %s", dec.encoding))
	}
}

func sizeof(t reflect.Type, v reflect.Value) int {
	switch t.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return int(v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		n := int(v.Uint())
		// all the builtin array length types are native int
		// so this guards against weird truncation
		if n < 0 {
			return 0
		}
		return n
	default:
		panic(fmt.Sprintf("sizeof field not implemented for kind %s", t.Kind()))
	}
}

var ErrVarIntBufferSize = errors.New("varint: invalid buffer size")

func (dec *Decoder) ReadUvarint64() (uint64, error) {
	l, read := binary.Uvarint(dec.data[dec.pos:])
	if read <= 0 {
		return l, ErrVarIntBufferSize
	}
	if traceEnabled {
		zlog.Debug("decode: read uvarint64", zap.Uint64("val", l))
	}
	dec.pos += read
	return l, nil
}

func (d *Decoder) ReadVarint64() (out int64, err error) {
	l, read := binary.Varint(d.data[d.pos:])
	if read <= 0 {
		return l, ErrVarIntBufferSize
	}
	if traceEnabled {
		zlog.Debug("decode: read varint", zap.Int64("val", l))
	}
	d.pos += read
	return l, nil
}

func (dec *Decoder) ReadVarint32() (out int32, err error) {
	n, err := dec.ReadVarint64()
	if err != nil {
		return out, err
	}
	out = int32(n)
	if traceEnabled {
		zlog.Debug("decode: read varint32", zap.Int32("val", out))
	}
	return
}

func (dec *Decoder) ReadUvarint32() (out uint32, err error) {
	n, err := dec.ReadUvarint64()
	if err != nil {
		return out, err
	}
	out = uint32(n)
	if traceEnabled {
		zlog.Debug("decode: read uvarint32", zap.Uint32("val", out))
	}
	return
}

func (dec *Decoder) ReadVarint16() (out int16, err error) {
	n, err := dec.ReadVarint64()
	if err != nil {
		return out, err
	}
	out = int16(n)
	if traceEnabled {
		zlog.Debug("decode: read varint16", zap.Int16("val", out))
	}
	return
}

func (dec *Decoder) ReadUvarint16() (out uint16, err error) {
	n, err := dec.ReadUvarint64()
	if err != nil {
		return out, err
	}
	out = uint16(n)
	if traceEnabled {
		zlog.Debug("decode: read uvarint16", zap.Uint16("val", out))
	}
	return
}

func (dec *Decoder) ReadByteSlice() (out []byte, err error) {
	length, err := dec.ReadLength()
	if err != nil {
		return nil, err
	}

	if len(dec.data) < dec.pos+length {
		return nil, fmt.Errorf("byte array: varlen=%d, missing %d bytes", length, dec.pos+length-len(dec.data))
	}

	out = dec.data[dec.pos : dec.pos+length]
	dec.pos += length
	if traceEnabled {
		zlog.Debug("decode: read byte array", zap.Stringer("hex", HexBytes(out)))
	}
	return
}

func (dec *Decoder) ReadLength() (length int, err error) {
	switch dec.encoding {
	case EncodingBin:
		val, err := dec.ReadUvarint64()
		if err != nil {
			return 0, err
		}
		if val > 0x7FFF_FFFF {
			return 0, io.ErrUnexpectedEOF
		}
		length = int(val)
	case EncodingBorsh:
		val, err := dec.ReadUint32(LE)
		if err != nil {
			return 0, err
		}
		if val > 0x7FFF_FFFF {
			return 0, io.ErrUnexpectedEOF
		}
		length = int(val)
	case EncodingCompactU16:
		val, err := DecodeCompactU16LengthFromByteReader(dec)
		if err != nil {
			return 0, err
		}
		length = val
	default:
		panic(fmt.Errorf("encoding not implemented: %s", dec.encoding))
	}
	return
}

type peekAbleByteReader interface {
	io.ByteReader
	Peek(n int) ([]byte, error)
}

func readNBytes(n int, reader *Decoder) ([]byte, error) {
	if n == 0 {
		return make([]byte, 0), nil
	}
	if n < 0 || n > 0x7FFF_FFFF {
		return nil, fmt.Errorf("invalid length n: %v", n)
	}
	if reader.pos+n > len(reader.data) {
		return nil, fmt.Errorf("not enough data: %d bytes missing", reader.pos+n-len(reader.data))
	}
	out := reader.data[reader.pos : reader.pos+n]
	reader.pos += n
	return out, nil
}

func discardNBytes(n int, reader *Decoder) error {
	if n == 0 {
		return nil
	}
	if n < 0 || n > 0x7FFF_FFFF {
		return fmt.Errorf("invalid length n: %v", n)
	}
	return reader.SkipBytes(uint(n))
}

func (dec *Decoder) ReadNBytes(n int) (out []byte, err error) {
	return readNBytes(n, dec)
}

func (dec *Decoder) Discard(n int) (err error) {
	return discardNBytes(n, dec)
}

func (dec *Decoder) ReadTypeID() (out TypeID, err error) {
	discriminator, err := dec.ReadNBytes(8)
	if err != nil {
		return TypeID{}, err
	}
	return TypeIDFromBytes(discriminator), nil
}

func (dec *Decoder) Peek(n int) (out []byte, err error) {
	if n < 0 {
		err = fmt.Errorf("n not valid: %d", n)
		return
	}

	requiredSize := TypeSize.Byte * n
	if dec.Remaining() < requiredSize {
		err = fmt.Errorf("required [%d] bytes, remaining [%d]", requiredSize, dec.Remaining())
		return
	}

	out = dec.data[dec.pos : dec.pos+n]
	if traceEnabled {
		zlog.Debug("decode: peek", zap.Int("n", n), zap.Binary("out", out))
	}
	return
}

func (dec *Decoder) ReadByte() (out byte, err error) {
	if dec.Remaining() < TypeSize.Byte {
		err = fmt.Errorf("required [1] byte, remaining [%d]", dec.Remaining())
		return
	}

	out = dec.data[dec.pos]
	dec.pos++
	if traceEnabled {
		zlog.Debug("decode: read byte", zap.Uint8("byte", out), zap.String("hex", hex.EncodeToString([]byte{out})))
	}
	return
}

func (dec *Decoder) ReadBool() (out bool, err error) {
	if dec.Remaining() < TypeSize.Bool {
		err = fmt.Errorf("bool required [%d] byte, remaining [%d]", TypeSize.Bool, dec.Remaining())
		return
	}

	b, err := dec.ReadByte()

	if err != nil {
		err = fmt.Errorf("readBool, %s", err)
	}
	out = b != 0
	if traceEnabled {
		zlog.Debug("decode: read bool", zap.Bool("val", out))
	}
	return
}

func (dec *Decoder) ReadUint8() (out uint8, err error) {
	out, err = dec.ReadByte()
	return
}

func (dec *Decoder) ReadInt8() (out int8, err error) {
	b, err := dec.ReadByte()
	out = int8(b)
	if traceEnabled {
		zlog.Debug("decode: read int8", zap.Int8("val", out))
	}
	return
}

func (dec *Decoder) ReadUint16(order binary.ByteOrder) (out uint16, err error) {
	if dec.Remaining() < TypeSize.Uint16 {
		err = fmt.Errorf("uint16 required [%d] bytes, remaining [%d]", TypeSize.Uint16, dec.Remaining())
		return
	}

	out = order.Uint16(dec.data[dec.pos:])
	dec.pos += TypeSize.Uint16
	if traceEnabled {
		zlog.Debug("decode: read uint16", zap.Uint16("val", out))
	}
	return
}

func (dec *Decoder) ReadInt16(order binary.ByteOrder) (out int16, err error) {
	n, err := dec.ReadUint16(order)
	out = int16(n)
	if traceEnabled {
		zlog.Debug("decode: read int16", zap.Int16("val", out))
	}
	return
}

func (dec *Decoder) ReadInt64(order binary.ByteOrder) (out int64, err error) {
	n, err := dec.ReadUint64(order)
	out = int64(n)
	if traceEnabled {
		zlog.Debug("decode: read int64", zap.Int64("val", out))
	}
	return
}

func (dec *Decoder) ReadUint32(order binary.ByteOrder) (out uint32, err error) {
	if dec.Remaining() < TypeSize.Uint32 {
		err = fmt.Errorf("uint32 required [%d] bytes, remaining [%d]", TypeSize.Uint32, dec.Remaining())
		return
	}

	out = order.Uint32(dec.data[dec.pos:])
	dec.pos += TypeSize.Uint32
	if traceEnabled {
		zlog.Debug("decode: read uint32", zap.Uint32("val", out))
	}
	return
}

func (dec *Decoder) ReadInt32(order binary.ByteOrder) (out int32, err error) {
	n, err := dec.ReadUint32(order)
	out = int32(n)
	if traceEnabled {
		zlog.Debug("decode: read int32", zap.Int32("val", out))
	}
	return
}

func (dec *Decoder) ReadUint64(order binary.ByteOrder) (out uint64, err error) {
	if dec.Remaining() < TypeSize.Uint64 {
		err = fmt.Errorf("decode: uint64 required [%d] bytes, remaining [%d]", TypeSize.Uint64, dec.Remaining())
		return
	}

	data, err := dec.ReadNBytes(TypeSize.Uint64)
	if err != nil {
		return 0, err
	}
	out = order.Uint64(data)
	if traceEnabled {
		zlog.Debug("decode: read uint64", zap.Uint64("val", out), zap.Stringer("hex", HexBytes(data)))
	}
	return
}

func (dec *Decoder) ReadInt128(order binary.ByteOrder) (out Int128, err error) {
	v, err := dec.ReadUint128(order)
	if err != nil {
		return
	}
	return Int128(v), nil
}

func (dec *Decoder) ReadUint128(order binary.ByteOrder) (out Uint128, err error) {
	if dec.Remaining() < TypeSize.Uint128 {
		err = fmt.Errorf("uint128 required [%d] bytes, remaining [%d]", TypeSize.Uint128, dec.Remaining())
		return
	}

	data := dec.data[dec.pos : dec.pos+TypeSize.Uint128]

	if order == binary.LittleEndian {
		out.Lo = order.Uint64(data[:8])
		out.Hi = order.Uint64(data[8:])
	} else {
		// TODO: is this correct?
		out.Hi = order.Uint64(data[:8])
		out.Lo = order.Uint64(data[8:])
	}

	dec.pos += TypeSize.Uint128
	if traceEnabled {
		zlog.Debug("decode: read uint128", zap.Stringer("hex", out), zap.Uint64("hi", out.Hi), zap.Uint64("lo", out.Lo))
	}
	return
}

func (dec *Decoder) ReadFloat32(order binary.ByteOrder) (out float32, err error) {
	if dec.Remaining() < TypeSize.Float32 {
		err = fmt.Errorf("float32 required [%d] bytes, remaining [%d]", TypeSize.Float32, dec.Remaining())
		return
	}

	n := order.Uint32(dec.data[dec.pos:])
	out = math.Float32frombits(n)
	dec.pos += TypeSize.Float32
	if traceEnabled {
		zlog.Debug("decode: read float32", zap.Float32("val", out))
	}

	if dec.IsBorsh() {
		if math.IsNaN(float64(out)) {
			return 0, errors.New("NaN for float not allowed")
		}
	}
	return
}

func (dec *Decoder) ReadFloat64(order binary.ByteOrder) (out float64, err error) {
	if dec.Remaining() < TypeSize.Float64 {
		err = fmt.Errorf("float64 required [%d] bytes, remaining [%d]", TypeSize.Float64, dec.Remaining())
		return
	}

	n := order.Uint64(dec.data[dec.pos:])
	out = math.Float64frombits(n)
	dec.pos += TypeSize.Float64
	if traceEnabled {
		zlog.Debug("decode: read Float64", zap.Float64("val", out))
	}
	if dec.IsBorsh() {
		if math.IsNaN(out) {
			return 0, errors.New("NaN for float not allowed")
		}
	}
	return
}

func (dec *Decoder) ReadFloat128(order binary.ByteOrder) (out Float128, err error) {
	value, err := dec.ReadUint128(order)
	if err != nil {
		return out, fmt.Errorf("float128: %s", err)
	}
	return Float128(value), nil
}

func (dec *Decoder) SafeReadUTF8String() (out string, err error) {
	data, err := dec.ReadByteSlice()
	out = strings.Map(fixUtf, string(data))
	if traceEnabled {
		zlog.Debug("read safe UTF8 string", zap.String("val", out))
	}
	return
}

func fixUtf(r rune) rune {
	if r == utf8.RuneError {
		return '�'
	}
	return r
}

func (dec *Decoder) ReadString() (out string, err error) {
	data, err := dec.ReadByteSlice()
	out = string(data)
	if traceEnabled {
		zlog.Debug("read string", zap.String("val", out))
	}
	return
}

func (dec *Decoder) ReadRustString() (out string, err error) {
	length, err := dec.ReadUint64(binary.LittleEndian)
	if err != nil {
		return "", err
	}
	if length > 0x7FFF_FFFF {
		return "", io.ErrUnexpectedEOF
	}
	bytes, err := dec.ReadNBytes(int(length))
	if err != nil {
		return "", err
	}
	out = string(bytes)
	if traceEnabled {
		zlog.Debug("read Rust string", zap.String("val", out))
	}
	return
}

func (dec *Decoder) ReadCompactU16Length() (int, error) {
	val, err := DecodeCompactU16LengthFromByteReader(dec)
	if traceEnabled {
		zlog.Debug("read compact-u16 length", zap.Int("val", val))
	}
	return val, err
}

func (dec *Decoder) SkipBytes(count uint) error {
	if uint(dec.Remaining()) < count {
		return fmt.Errorf("request to skip %d but only %d bytes remain", count, dec.Remaining())
	}
	dec.pos += int(count)
	return nil
}

func (dec *Decoder) SetPosition(idx uint) error {
	if int(idx) < len(dec.data) {
		dec.pos = int(idx)
		return nil
	}
	return fmt.Errorf("request to set position to %d outsize of buffer (buffer size %d)", idx, len(dec.data))
}

func (dec *Decoder) Position() uint {
	return uint(dec.pos)
}

func (dec *Decoder) Remaining() int {
	return len(dec.data) - dec.pos
}

func (dec *Decoder) HasRemaining() bool {
	return dec.Remaining() > 0
}

// indirect walks down v allocating pointers as needed,
// until it gets to a non-pointer.
// if it encounters an Unmarshaler, indirect stops and returns that.
// if decodingNull is true, indirect stops at the last pointer so it can be set to nil.
//
// *Note* This is a copy of `encoding/json/decoder.go#indirect` of Golang 1.14.
//
// See here: https://github.com/golang/go/blob/go1.14.2/src/encoding/json/decode.go#L439
func indirect(v reflect.Value, decodingNull bool) (BinaryUnmarshaler, reflect.Value) {
	// Issue #24153 indicates that it is generally not a guaranteed property
	// that you may round-trip a reflect.Value by calling Value.Addr().Elem()
	// and expect the value to still be settable for values derived from
	// unexported embedded struct fields.
	//
	// The logic below effectively does this when it first addresses the value
	// (to satisfy possible pointer methods) and continues to dereference
	// subsequent pointers as necessary.
	//
	// After the first round-trip, we set v back to the original value to
	// preserve the original RW flags contained in reflect.Value.
	v0 := v
	haveAddr := false

	// If v is a named type and is addressable,
	// start with its address, so that if the type has pointer methods,
	// we find them.
	if v.Kind() != reflect.Ptr && v.Type().Name() != "" && v.CanAddr() {
		haveAddr = true
		v = v.Addr()
	}
	for {
		// Load value from interface, but only if the result will be
		// usefully addressable.
		if v.Kind() == reflect.Interface && !v.IsNil() {
			e := v.Elem()
			if e.Kind() == reflect.Ptr && !e.IsNil() && (!decodingNull || e.Elem().Kind() == reflect.Ptr) {
				haveAddr = false
				v = e
				continue
			}
		}

		if v.Kind() != reflect.Ptr {
			break
		}

		if v.Elem().Kind() != reflect.Ptr && decodingNull && v.CanSet() {
			break
		}

		// Prevent infinite loop if v is an interface pointing to its own address:
		//     var v interface{}
		//     v = &v
		if v.Elem().Kind() == reflect.Interface && v.Elem().Elem() == v {
			v = v.Elem()
			break
		}
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		if v.Type().NumMethod() > 0 && v.CanInterface() {
			if u, ok := v.Interface().(BinaryUnmarshaler); ok {
				return u, reflect.Value{}
			}
		}

		if haveAddr {
			v = v0 // restore original value after round-trip Value.Addr().Elem()
			haveAddr = false
		} else {
			v = v.Elem()
		}
	}
	return nil, v
}

func reflect_readArrayOfBytes(d *Decoder, l int, rv reflect.Value) error {
	buf, err := d.ReadNBytes(l)
	if err != nil {
		return err
	}
	switch rv.Kind() {
	case reflect.Array:
		// if the type of the array is not [n]uint8, but a custom type like [n]CustomUint8:
		if rv.Type().Elem() != typeOfUint8 {
			// if the type of the array is not [n]uint8, but a custom type like [n]CustomUint8:
			// then we need to convert each uint8 to the custom type
			for i := 0; i < l; i++ {
				rv.Index(i).Set(reflect.ValueOf(buf[i]).Convert(rv.Index(i).Type()))
			}
		} else {
			reflect.Copy(rv, reflect.ValueOf(buf))
		}
	case reflect.Slice:
		// if the type of the slice is not []uint8, but a custom type like []CustomUint8:
		if rv.Type().Elem() != typeOfUint8 {
			// convert the []uint8 to the custom type
			customSlice := reflect.MakeSlice(rv.Type(), len(buf), len(buf))
			for i := 0; i < len(buf); i++ {
				customSlice.Index(i).SetUint(uint64(buf[i]))
			}
			rv.Set(customSlice)
		} else {
			rv.Set(reflect.ValueOf(buf))
		}
	default:
		return fmt.Errorf("unsupported kind: %s", rv.Kind())
	}
	return nil
}

func reflect_readArrayOfUint16(d *Decoder, l int, rv reflect.Value, order binary.ByteOrder) error {
	buf := make([]uint16, l)
	for i := 0; i < l; i++ {
		n, err := d.ReadUint16(order)
		if err != nil {
			return err
		}
		buf[i] = n
	}
	switch rv.Kind() {
	case reflect.Array:
		// if the type of the array is not [n]uint16, but a custom type like [n]CustomUint16:
		if rv.Type().Elem() != typeOfUint16 {
			// if the type of the array is not [n]uint16, but a custom type like [n]CustomUint16:
			// then we need to convert each uint16 to the custom type
			for i := 0; i < l; i++ {
				rv.Index(i).Set(reflect.ValueOf(buf[i]).Convert(rv.Index(i).Type()))
			}
		} else {
			reflect.Copy(rv, reflect.ValueOf(buf))
		}
	case reflect.Slice:
		// if the type of the slice is not []uint16, but a custom type like []CustomUint16:
		if rv.Type().Elem() != typeOfUint16 {
			// convert the []uint16 to the custom type
			customSlice := reflect.MakeSlice(rv.Type(), len(buf), len(buf))
			for i := 0; i < len(buf); i++ {
				customSlice.Index(i).SetUint(uint64(buf[i]))
			}
			rv.Set(customSlice)
		} else {
			rv.Set(reflect.ValueOf(buf))
		}
	default:
		return fmt.Errorf("unsupported kind: %s", rv.Kind())
	}
	return nil
}

func reflect_readArrayOfUint32(d *Decoder, l int, rv reflect.Value, order binary.ByteOrder) error {
	buf := make([]uint32, l)
	for i := 0; i < l; i++ {
		n, err := d.ReadUint32(order)
		if err != nil {
			return err
		}
		buf[i] = n
	}
	switch rv.Kind() {
	case reflect.Array:
		// if the type of the array is not [n]uint32, but a custom type like [n]CustomUint32:
		if rv.Type().Elem() != typeOfUint32 {
			// if the type of the array is not [n]uint32, but a custom type like [n]CustomUint32:
			// then we need to convert each uint32 to the custom type
			for i := 0; i < l; i++ {
				rv.Index(i).Set(reflect.ValueOf(buf[i]).Convert(rv.Index(i).Type()))
			}
		} else {
			reflect.Copy(rv, reflect.ValueOf(buf))
		}
	case reflect.Slice:
		// if the type of the slice is not []uint32, but a custom type like []CustomUint32:
		if rv.Type().Elem() != typeOfUint32 {
			// convert the []uint32 to the custom type
			customSlice := reflect.MakeSlice(rv.Type(), len(buf), len(buf))
			for i := 0; i < len(buf); i++ {
				customSlice.Index(i).SetUint(uint64(buf[i]))
			}
			rv.Set(customSlice)
		} else {
			rv.Set(reflect.ValueOf(buf))
		}
	default:
		return fmt.Errorf("unsupported kind: %s", rv.Kind())
	}
	return nil
}

func init() {
	if typeOfByte != typeOfUint8 {
		panic("typeOfByte != typeOfUint8")
	}
}

var (
	typeOfByte   = reflect.TypeOf(byte(0))
	typeOfUint8  = reflect.TypeOf(uint8(0))
	typeOfUint16 = reflect.TypeOf(uint16(0))
	typeOfUint32 = reflect.TypeOf(uint32(0))
	typeOfUint64 = reflect.TypeOf(uint64(0))
)

func reflect_readArrayOfUint64(d *Decoder, l int, rv reflect.Value, order binary.ByteOrder) error {
	buf := make([]uint64, l)
	for i := 0; i < l; i++ {
		n, err := d.ReadUint64(order)
		if err != nil {
			return err
		}
		buf[i] = n
	}
	switch rv.Kind() {
	case reflect.Array:
		// if the type of the array is not [n]uint64, but a custom type like [n]CustomUint64:
		if rv.Type().Elem() != typeOfUint64 {
			// if the type of the array is not [n]uint64, but a custom type like [n]CustomUint64:
			// then we need to convert each uint64 to the custom type
			for i := 0; i < l; i++ {
				rv.Index(i).Set(reflect.ValueOf(buf[i]).Convert(rv.Index(i).Type()))
			}
		} else {
			reflect.Copy(rv, reflect.ValueOf(buf))
		}
	case reflect.Slice:
		// if the type of the slice is not []uint64, but a custom type like []CustomUint64:
		if rv.Type().Elem() != typeOfUint64 {
			// convert the []uint64 to the custom type
			customSlice := reflect.MakeSlice(rv.Type(), len(buf), len(buf))
			for i := 0; i < len(buf); i++ {
				customSlice.Index(i).SetUint(uint64(buf[i]))
			}
			rv.Set(customSlice)
		} else {
			rv.Set(reflect.ValueOf(buf))
		}
	default:
		return fmt.Errorf("unsupported kind: %s", rv.Kind())
	}
	return nil
}

// reflect_readArrayOfUint_ is used for reading arrays/slices of uints of any size.
func reflect_readArrayOfUint_(d *Decoder, l int, k reflect.Kind, rv reflect.Value, order binary.ByteOrder) error {
	switch k {
	// case reflect.Uint:
	// 	// switch on system architecture (32 or 64 bit)
	// 	if unsafe.Sizeof(uintptr(0)) == 4 {
	// 		return reflect_readArrayOfUint32(  d, l, rv, order)
	// 	}
	// 	return reflect_readArrayOfUint64(  d, l, rv, order)
	case reflect.Uint8:
		if l > d.Remaining() {
			return io.ErrUnexpectedEOF
		}
		return reflect_readArrayOfBytes(d, l, rv)
	case reflect.Uint16:
		if l*2 > d.Remaining() {
			return io.ErrUnexpectedEOF
		}
		return reflect_readArrayOfUint16(d, l, rv, order)
	case reflect.Uint32:
		if l*4 > d.Remaining() {
			return io.ErrUnexpectedEOF
		}
		return reflect_readArrayOfUint32(d, l, rv, order)
	case reflect.Uint64:
		if l*8 > d.Remaining() {
			return io.ErrUnexpectedEOF
		}
		return reflect_readArrayOfUint64(d, l, rv, order)
	default:
		return fmt.Errorf("unsupported kind: %v", k)
	}
}

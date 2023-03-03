package abi

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"reflect"
	"strconv"
	"strings"

	"github.com/umbracle/ethgo"
)

var (
	zero = big.NewInt(0)
	one  = big.NewInt(1)
)

// Encode encodes a value
func Encode(v interface{}, t *Type) ([]byte, error) {
	return encode(reflect.ValueOf(v), t)
}

func encode(v reflect.Value, t *Type) ([]byte, error) {
	if v.Kind() == reflect.Interface {
		v = v.Elem()
	}

	switch t.kind {
	case KindSlice, KindArray:
		return encodeSliceAndArray(v, t)

	case KindTuple:
		return encodeTuple(v, t)

	case KindString:
		return encodeString(v)

	case KindBool:
		return encodeBool(v)

	case KindAddress:
		return encodeAddress(v)

	case KindInt, KindUInt:
		return encodeNum(v)

	case KindBytes:
		return encodeBytes(v)

	case KindFixedBytes, KindFunction:
		return encodeFixedBytes(v)

	default:
		return nil, fmt.Errorf("encoding not available for type '%s'", t.kind)
	}
}

func encodeSliceAndArray(v reflect.Value, t *Type) ([]byte, error) {
	if v.Kind() != reflect.Array && v.Kind() != reflect.Slice {
		return nil, encodeErr(v, t.kind.String())
	}

	if v.Kind() == reflect.Array && t.kind != KindArray {
		return nil, fmt.Errorf("expected array")
	} else if v.Kind() == reflect.Slice && t.kind != KindSlice {
		return nil, fmt.Errorf("expected slice")
	}

	if t.kind == KindArray && t.size != v.Len() {
		return nil, fmt.Errorf("array len incompatible")
	}

	var ret, tail []byte
	if t.isVariableInput() {
		ret = append(ret, packNum(v.Len())...)
	}

	offset := 0
	isDynamic := t.elem.isDynamicType()
	if isDynamic {
		offset = getTypeSize(t.elem) * v.Len()
	}

	for i := 0; i < v.Len(); i++ {
		val, err := encode(v.Index(i), t.elem)
		if err != nil {
			return nil, err
		}
		if !isDynamic {
			ret = append(ret, val...)
		} else {
			ret = append(ret, packNum(offset)...)
			offset += len(val)
			tail = append(tail, val...)
		}
	}
	return append(ret, tail...), nil
}

func encodeTuple(v reflect.Value, t *Type) ([]byte, error) {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	var err error
	isList := true

	switch v.Kind() {
	case reflect.Slice, reflect.Array:
	case reflect.Map:
		isList = false

	case reflect.Struct:
		isList = false
		v, err = mapFromStruct(v)
		if err != nil {
			return nil, err
		}

	default:
		return nil, encodeErr(v, "tuple")
	}

	if v.Len() < len(t.tuple) {
		return nil, fmt.Errorf("expected at least the same length")
	}

	offset := 0
	for _, elem := range t.tuple {
		offset += getTypeSize(elem.Elem)
	}

	var ret, tail []byte
	var aux reflect.Value

	for i, elem := range t.tuple {
		if isList {
			aux = v.Index(i)
		} else {
			name := elem.Name
			if name == "" {
				name = strconv.Itoa(i)
			}
			aux = v.MapIndex(reflect.ValueOf(name))
		}
		if aux.Kind() == reflect.Invalid {
			return nil, fmt.Errorf("cannot get key %s", elem.Name)
		}

		val, err := encode(aux, elem.Elem)
		if err != nil {
			return nil, err
		}
		if elem.Elem.isDynamicType() {
			ret = append(ret, packNum(offset)...)
			tail = append(tail, val...)
			offset += len(val)
		} else {
			ret = append(ret, val...)
		}
	}

	return append(ret, tail...), nil
}

func convertArrayToBytes(value reflect.Value) reflect.Value {
	slice := reflect.MakeSlice(reflect.TypeOf([]byte{}), value.Len(), value.Len())
	reflect.Copy(slice, value)
	return slice
}

func encodeFixedBytes(v reflect.Value) ([]byte, error) {
	if v.Kind() == reflect.Array {
		v = convertArrayToBytes(v)
	}
	if v.Kind() == reflect.String {
		value, err := decodeHex(v.String())
		if err != nil {
			return nil, err
		}

		v = reflect.ValueOf(value)
	}
	return rightPad(v.Bytes(), 32), nil
}

func encodeAddress(v reflect.Value) ([]byte, error) {
	if v.Kind() == reflect.Array {
		v = convertArrayToBytes(v)
	}
	if v.Kind() == reflect.String {
		var addr ethgo.Address
		if err := addr.UnmarshalText([]byte(v.String())); err != nil {
			return nil, err
		}
		v = reflect.ValueOf(addr.Bytes())
	}
	return leftPad(v.Bytes(), 32), nil
}

func encodeBytes(v reflect.Value) ([]byte, error) {
	if v.Kind() == reflect.Array {
		v = convertArrayToBytes(v)
	}
	if v.Kind() == reflect.String {
		value, err := decodeHex(v.String())
		if err != nil {
			return nil, err
		}

		v = reflect.ValueOf(value)
	}
	return packBytesSlice(v.Bytes(), v.Len())
}

func encodeString(v reflect.Value) ([]byte, error) {
	if v.Kind() != reflect.String {
		return nil, encodeErr(v, "string")
	}
	return packBytesSlice([]byte(v.String()), v.Len())
}

func packBytesSlice(buf []byte, l int) ([]byte, error) {
	len, err := encodeNum(reflect.ValueOf(l))
	if err != nil {
		return nil, err
	}
	return append(len, rightPad(buf, (l+31)/32*32)...), nil
}

func packNum(offset int) []byte {
	n, _ := encodeNum(reflect.ValueOf(offset))
	return n
}

func encodeNum(v reflect.Value) ([]byte, error) {
	switch v.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return toU256(new(big.Int).SetUint64(v.Uint())), nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return toU256(big.NewInt(v.Int())), nil

	case reflect.Ptr:
		if v.Type() != bigIntT {
			return nil, encodeErr(v.Elem(), "number")
		}
		return toU256(v.Interface().(*big.Int)), nil

	case reflect.Float64:
		return encodeNum(reflect.ValueOf(int64(v.Float())))

	case reflect.String:
		n, ok := new(big.Int).SetString(v.String(), 10)
		if !ok {
			n, ok = new(big.Int).SetString(v.String()[2:], 16)
			if !ok {
				return nil, encodeErr(v, "number")
			}
		}
		return encodeNum(reflect.ValueOf(n))

	default:
		return nil, encodeErr(v, "number")
	}
}

func encodeBool(v reflect.Value) ([]byte, error) {
	if v.Kind() != reflect.Bool {
		return nil, encodeErr(v, "bool")
	}
	if v.Bool() {
		return leftPad(one.Bytes(), 32), nil
	}
	return leftPad(zero.Bytes(), 32), nil
}

func encodeErr(v reflect.Value, t string) error {
	return fmt.Errorf("failed to encode %s as %s", v.Kind().String(), t)
}

func mapFromStruct(v reflect.Value) (reflect.Value, error) {
	res := map[string]interface{}{}
	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		f := typ.Field(i)
		if f.PkgPath != "" {
			continue
		}

		tagValue := f.Tag.Get("abi")
		if tagValue == "-" {
			continue
		}

		name := f.Name
		if tagValue != "" {
			name = tagValue
		}

		name = strings.ToLower(name)
		if _, ok := res[name]; !ok {
			res[name] = v.Field(i).Interface()
		}
	}
	return reflect.ValueOf(res), nil
}

var (
	tt256   = new(big.Int).Lsh(big.NewInt(1), 256)   // 2 ** 256
	tt256m1 = new(big.Int).Sub(tt256, big.NewInt(1)) // 2 ** 256 - 1
)

// U256 converts a big Int into a 256bit EVM number.
func toU256(n *big.Int) []byte {
	b := new(big.Int)
	b = b.Set(n)

	if b.Sign() < 0 || b.BitLen() > 256 {
		b.And(b, tt256m1)
	}

	return leftPad(b.Bytes(), 32)
}

func padBytes(b []byte, size int, left bool) []byte {
	l := len(b)
	if l == size {
		return b
	}
	if l > size {
		return b[l-size:]
	}
	tmp := make([]byte, size)
	if left {
		copy(tmp[size-l:], b)
	} else {
		copy(tmp, b)
	}
	return tmp
}

func leftPad(b []byte, size int) []byte {
	return padBytes(b, size, true)
}

func rightPad(b []byte, size int) []byte {
	return padBytes(b, size, false)
}

func encodeHex(b []byte) string {
	return "0x" + hex.EncodeToString(b)
}

func decodeHex(str string) ([]byte, error) {
	if strings.HasPrefix(str, "0x") {
		str = str[2:]
	}
	buf, err := hex.DecodeString(str)
	if err != nil {
		return nil, fmt.Errorf("could not decode hex: %v", err)
	}
	return buf, nil
}

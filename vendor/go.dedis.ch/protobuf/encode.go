package protobuf

import (
	"bytes"
	"encoding"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"reflect"
	"time"
)

// Message fields declared to have exactly this type
// will be transmitted as fixed-size 32-bit unsigned integers.
type Ufixed32 uint32

// Message fields declared to have exactly this type
// will be transmitted as fixed-size 64-bit unsigned integers.
type Ufixed64 uint64

// Message fields declared to have exactly this type
// will be transmitted as fixed-size 32-bit signed integers.
type Sfixed32 int32

// Message fields declared to have exactly this type
// will be transmitted as fixed-size 64-bit signed integers.
type Sfixed64 int64

// Protobufs enums are transmitted as unsigned varints;
// using this type alias is optional but recommended
// to ensure they get the correct type.
type Enum uint32

type encoder struct {
	bytes.Buffer
}

// Encode a Go struct into protocol buffer format.
// The caller must pass a pointer to the struct to encode.
func Encode(structPtr interface{}) (bytes []byte, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
			bytes = nil
		}
	}()
	if structPtr == nil {
		return nil, nil
	}

	if bu, ok := structPtr.(encoding.BinaryMarshaler); ok {
		return bu.MarshalBinary()
	}

	en := encoder{}
	val := reflect.ValueOf(structPtr)
	if val.Kind() != reflect.Ptr {
		return nil, errors.New("encode takes a pointer to struct")
	}
	en.message(val.Elem())
	return en.Bytes(), nil
}

func (en *encoder) message(sval reflect.Value) {
	var index *ProtoField
	defer func() {
		if r := recover(); r != nil {
			if index != nil {
				panic(fmt.Sprintf("%s (field %s)", r, index.Field.Name))
			} else {
				panic(r)
			}
		}
	}()
	// Encode all fields in-order
	protoFields := ProtoFields(sval.Type())
	if len(protoFields) == 0 {
		return
	}
	noPublicFields := true
	for _, index = range protoFields {
		field := sval.FieldByIndex(index.Index)
		key := uint64(index.ID) << 3
		if field.CanSet() { // Skip blank/padding fields
			en.value(key, field, index.Prefix)
			noPublicFields = false
		}
	}
	if noPublicFields {
		panic("struct has no serializable fields")
	}
}

var timeType = reflect.TypeOf(time.Time{})
var durationType = reflect.TypeOf(time.Duration(0))

func (en *encoder) value(key uint64, val reflect.Value, prefix TagPrefix) {

	// Non-reflectively handle some of the fixed types
	switch v := val.Interface().(type) {
	case bool:
		en.uvarint(key | 0)
		vi := uint64(0)
		if v {
			vi = 1
		}
		en.uvarint(vi)
		return

	case int:
		en.uvarint(key | 0)
		en.svarint(int64(v))
		return

	case int32:
		en.uvarint(key | 0)
		en.svarint(int64(v))
		return

	case time.Time: // Encode time.Time as sfixed64
		t := v.UnixNano()
		en.uvarint(key | 1)
		en.u64(uint64(t))
		return

	case int64:
		en.uvarint(key | 0)
		en.svarint(v)
		return

	case uint32:
		en.uvarint(key | 0)
		en.uvarint(uint64(v))
		return

	case uint64:
		en.uvarint(key | 0)
		en.uvarint(v)
		return

	case Sfixed32:
		en.uvarint(key | 5)
		en.u32(uint32(v))
		return

	case Sfixed64:
		en.uvarint(key | 1)
		en.u64(uint64(v))
		return

	case Ufixed32:
		en.uvarint(key | 5)
		en.u32(uint32(v))
		return

	case Ufixed64:
		en.uvarint(key | 1)
		en.u64(uint64(v))
		return

	case float32:
		en.uvarint(key | 5)
		en.u32(math.Float32bits(v))
		return

	case float64:
		en.uvarint(key | 1)
		en.u64(math.Float64bits(v))
		return

	case string:
		en.uvarint(key | 2)
		b := []byte(v)
		en.uvarint(uint64(len(b)))
		en.Write(b)
		return
	}

	// Handle pointer or interface values (possibly within slices).
	// Note that this switch has to handle all the cases,
	// because custom type aliases will fail the above typeswitch.
	switch val.Kind() {
	case reflect.Bool:
		en.uvarint(key | 0)
		v := uint64(0)
		if val.Bool() {
			v = 1
		}
		en.uvarint(v)

	case reflect.Int, reflect.Int32, reflect.Int64:
		// Varint-encoded 32-bit and 64-bit signed integers.
		// Note that protobufs don't support 8- or 16-bit ints.
		en.uvarint(key | 0)
		en.svarint(val.Int())

	case reflect.Uint32, reflect.Uint64:
		// Varint-encoded 32-bit and 64-bit unsigned integers.
		en.uvarint(key | 0)
		en.uvarint(val.Uint())

	case reflect.Float32:
		// Fixed-length 32-bit floats.
		en.uvarint(key | 5)
		en.u32(math.Float32bits(float32(val.Float())))

	case reflect.Float64:
		// Fixed-length 64-bit floats.
		en.uvarint(key | 1)
		en.u64(math.Float64bits(val.Float()))

	case reflect.String:
		// Length-delimited string.
		en.uvarint(key | 2)
		b := []byte(val.String())
		en.uvarint(uint64(len(b)))
		en.Write(b)

	case reflect.Struct:
		var b []byte
		if enc, ok := val.Interface().(encoding.BinaryMarshaler); ok {
			en.uvarint(key | 2)
			var err error
			b, err = enc.MarshalBinary()
			if err != nil {
				panic(err.Error())
			}
		} else {
			// Embedded messages.
			en.uvarint(key | 2)
			emb := encoder{}
			emb.message(val)
			b = emb.Bytes()
		}
		en.uvarint(uint64(len(b)))
		en.Write(b)
	case reflect.Slice, reflect.Array:
		// Length-delimited slices or byte-vectors.
		en.slice(key, val)
		return

	case reflect.Ptr:
		// Optional field: encode only if pointer is non-nil.
		if val.IsNil() {
			if prefix == TagRequired {
				panic("required field is nil")
			}
			return
		}
		en.value(key, val.Elem(), prefix)

	case reflect.Interface:
		// Abstract interface field.
		if val.IsNil() {
			return
		}

		// If the object support self-encoding, use that.
		if enc, ok := val.Interface().(encoding.BinaryMarshaler); ok {
			en.uvarint(key | 2)
			bytes, err := enc.MarshalBinary()
			if err != nil {
				panic(err.Error())
			}

			size := len(bytes)
			var id GeneratorID
			im, ok := val.Interface().(InterfaceMarshaler)
			if ok {
				id = im.MarshalID()

				g := generators.get(id)
				ok = g != nil
				if ok {
					// add the length of the type tag
					size += len(id)
				}
			}

			en.uvarint(uint64(size))
			if ok {
				// Only write the tag if a generator exists
				en.Write(id[:])
			}
			en.Write(bytes)
			return
		}

		// Encode from the object the interface points to.
		en.value(key, val.Elem(), prefix)

	case reflect.Map:
		en.handleMap(key, val, prefix)
		return

	default:
		panic(fmt.Sprintf("unsupported field Kind %d", val.Kind()))
	}
}

func (en *encoder) slice(key uint64, slval reflect.Value) {

	// First handle common cases with a direct typeswitch
	sllen := slval.Len()
	packed := encoder{}
	switch slt := slval.Interface().(type) {
	case []bool:
		for i := 0; i < sllen; i++ {
			v := uint64(0)
			if slt[i] {
				v = 1
			}
			packed.uvarint(v)
		}

	case []int32:
		for i := 0; i < sllen; i++ {
			packed.svarint(int64(slt[i]))
		}

	case []int64:
		for i := 0; i < sllen; i++ {
			packed.svarint(slt[i])
		}

	case []uint32:
		for i := 0; i < sllen; i++ {
			packed.uvarint(uint64(slt[i]))
		}

	case []uint64:
		for i := 0; i < sllen; i++ {
			packed.uvarint(slt[i])
		}

	case []Sfixed32:
		for i := 0; i < sllen; i++ {
			packed.u32(uint32(slt[i]))
		}

	case []Sfixed64:
		for i := 0; i < sllen; i++ {
			packed.u64(uint64(slt[i]))
		}

	case []Ufixed32:
		for i := 0; i < sllen; i++ {
			packed.u32(uint32(slt[i]))
		}

	case []Ufixed64:
		for i := 0; i < sllen; i++ {
			packed.u64(uint64(slt[i]))
		}

	case []float32:
		for i := 0; i < sllen; i++ {
			packed.u32(math.Float32bits(slt[i]))
		}

	case []float64:
		for i := 0; i < sllen; i++ {
			packed.u64(math.Float64bits(slt[i]))
		}

	case []byte: // Write the whole byte-slice as one key,value pair
		en.uvarint(key | 2)
		en.uvarint(uint64(sllen))
		en.Write(slt)
		return

	case []string:
		for i := 0; i < sllen; i++ {
			subVal := slval.Index(i)
			subStr := subVal.Interface().(string)
			subSlice := []byte(subStr)
			en.uvarint(key | 2)
			en.uvarint(uint64(len(subSlice)))
			en.Write(subSlice)
		}
		return
	default: // We'll need to use the reflective path
		en.sliceReflect(key, slval)
		return
	}

	// Encode packed representation key/value pair
	en.uvarint(key | 2)
	b := packed.Bytes()
	en.uvarint(uint64(len(b)))
	en.Write(b)
}

// Handle the encoding of an arbritary map[K]V
func (en *encoder) handleMap(key uint64, mpval reflect.Value, prefix TagPrefix) {
	/*
		A map defined as
			map<key_type, value_type> map_field = N;
		is encoded in the same way as
			message MapFieldEntry {
				key_type key = 1;
				value_type value = 2;
			}
			repeated MapFieldEntry map_field = N;
	*/

	for _, mkey := range mpval.MapKeys() {
		mval := mpval.MapIndex(mkey)

		// illegal map entry values
		// - nil message pointers.
		switch kind := mval.Kind(); kind {
		case reflect.Ptr:
			if mval.IsNil() {
				panic("proto: map has nil element")
			}
		case reflect.Slice, reflect.Array:
			if mval.Type().Elem().Kind() != reflect.Uint8 {
				panic("protobuf: map only support []byte or string as repeated value")
			}
		}

		packed := encoder{}
		packed.value(1<<3, mkey, prefix)
		packed.value(2<<3, mval, prefix)

		en.uvarint(key | 2)
		b := packed.Bytes()
		en.uvarint((uint64(len(b))))
		en.Write(b)
	}
}

var bytesType = reflect.TypeOf([]byte{})

func (en *encoder) sliceReflect(key uint64, slval reflect.Value) {
	kind := slval.Kind()
	if kind != reflect.Slice && kind != reflect.Array {
		panic("no slice passed")
	}
	sllen := slval.Len()
	slelt := slval.Type().Elem()
	packed := encoder{}
	switch slelt.Kind() {
	case reflect.Bool:
		for i := 0; i < sllen; i++ {
			v := uint64(0)
			if slval.Index(i).Bool() {
				v = 1
			}
			packed.uvarint(v)
		}

	case reflect.Int, reflect.Int32, reflect.Int64:
		for i := 0; i < sllen; i++ {
			packed.svarint(slval.Index(i).Int())
		}

	case reflect.Uint32, reflect.Uint64:
		for i := 0; i < sllen; i++ {
			packed.uvarint(slval.Index(i).Uint())
		}

	case reflect.Float32:
		for i := 0; i < sllen; i++ {
			packed.u32(math.Float32bits(
				float32(slval.Index(i).Float())))
		}

	case reflect.Float64:
		for i := 0; i < sllen; i++ {
			packed.u64(math.Float64bits(slval.Index(i).Float()))
		}

	case reflect.Uint8: // Write the byte-slice as one key,value pair
		en.uvarint(key | 2)
		en.uvarint(uint64(sllen))
		var b []byte
		if slval.Kind() == reflect.Array {
			if slval.CanAddr() {
				sliceVal := slval.Slice(0, sllen)
				b = sliceVal.Convert(bytesType).Interface().([]byte)
			} else {
				sliceVal := reflect.MakeSlice(bytesType, sllen, sllen)
				reflect.Copy(sliceVal, slval)
				b = sliceVal.Interface().([]byte)
			}
		} else {
			b = slval.Convert(bytesType).Interface().([]byte)
		}
		en.Write(b)
		return

	default: // Write each element as a separate key,value pair
		t := slval.Type().Elem()
		if t.Kind() == reflect.Slice || t.Kind() == reflect.Array {
			subSlice := t.Elem()
			if subSlice.Kind() != reflect.Uint8 {
				panic("protobuf: no support for 2-dimensional array except for [][]byte")
			}
		}
		for i := 0; i < sllen; i++ {
			en.value(key, slval.Index(i), TagNone)
		}
		return
	}

	// Encode packed representation key/value pair
	en.uvarint(key | 2)
	b := packed.Bytes()
	en.uvarint(uint64(len(b)))
	en.Write(b)
}

func (en *encoder) uvarint(v uint64) {
	var b [binary.MaxVarintLen64]byte
	n := binary.PutUvarint(b[:], v)
	en.Write(b[:n])
}

func (en *encoder) svarint(v int64) {
	if v >= 0 {
		en.uvarint(uint64(v) << 1)
	} else {
		en.uvarint(^uint64(v << 1))
	}
}

func (en *encoder) u32(v uint32) {
	var b [4]byte
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)
	en.Write(b[:])
}

func (en *encoder) u64(v uint64) {
	var b [8]byte
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)
	b[4] = byte(v >> 32)
	b[5] = byte(v >> 40)
	b[6] = byte(v >> 48)
	b[7] = byte(v >> 56)
	en.Write(b[:])
}

package protobuf

import (
	"encoding"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"reflect"
	"time"
)

// Constructors represents a map defining how to instantiate any interface
// types that Decode() might encounter while reading and decoding structured
// data. The keys are reflect.Type values denoting interface types. The
// corresponding values are functions expected to instantiate, and initialize
// as necessary, an appropriate concrete object type supporting that
// interface. A caller could use this capability to support
// dynamic instantiation of objects of the concrete type
// appropriate for a given abstract type.
type Constructors map[reflect.Type]func() interface{}

// String returns an easy way to visualize what you have in your constructors.
func (c *Constructors) String() string {
	var s string
	for k := range *c {
		s += k.String() + "=>" + "(func() interface {})" + "\t"
	}
	return s
}

// Decoder is the main struct used to decode a protobuf blob.
type decoder struct {
	nm Constructors
}

// Decode a protocol buffer into a Go struct.
// The caller must pass a pointer to the struct to decode into.
//
// Decode() currently does not explicitly check that all 'required' fields
// are actually present in the input buffer being decoded.
// If required fields are missing, then the corresponding fields
// will be left unmodified, meaning they will take on
// their default Go zero values if Decode() is passed a fresh struct.
func Decode(buf []byte, structPtr interface{}) error {
	return DecodeWithConstructors(buf, structPtr, nil)
}

// DecodeWithConstructors is like Decode, but you can pass a map of
// constructors with which to instantiate interface types.
func DecodeWithConstructors(buf []byte, structPtr interface{}, cons Constructors) (err error) {
	defer func() {
		if r := recover(); r != nil {
			switch e := r.(type) {
			case string:
				err = errors.New(e)
			case error:
				err = e
			default:
				err = errors.New("Failed to decode the field")
			}
		}
	}()
	if structPtr == nil {
		return nil
	}

	if bu, ok := structPtr.(encoding.BinaryUnmarshaler); ok {
		return bu.UnmarshalBinary(buf)
	}

	de := decoder{cons}
	val := reflect.ValueOf(structPtr)
	// if its NOT a pointer, it is bad return an error
	if val.Kind() != reflect.Ptr {
		return errors.New("Decode has been given a non pointer type")
	}
	return de.message(buf, val.Elem())
}

// Decode a Protocol Buffers message into a Go struct.
// The Kind of the passed value v must be Struct.
func (de *decoder) message(buf []byte, sval reflect.Value) error {
	if sval.Kind() != reflect.Struct {
		return errors.New("not a struct")
	}

	for i := 0; i < sval.NumField(); i++ {
		switch field := sval.Field(i); field.Kind() {
		case reflect.Interface:
			// Interface are not reset because the decoder won't
			// be able to instantiate it again in some scenarios.
		default:
			if field.CanSet() {
				field.Set(reflect.Zero(field.Type()))
			}
		}
	}

	// Decode all the fields
	fields := ProtoFields(sval.Type())
	fieldi := 0
	for len(buf) > 0 {
		// Parse the key
		key, n := binary.Uvarint(buf)
		if n <= 0 {
			return errors.New("bad protobuf field key")
		}
		buf = buf[n:]
		wiretype := int(key & 7)
		fieldnum := key >> 3

		// Lookup the corresponding struct field.
		// Leave field with a zero Value if fieldnum is out-of-range.
		// In this case, as well as for blank fields,
		// value() will just skip over and discard the field content.
		var field reflect.Value
		for fieldi < len(fields) && fields[fieldi].ID < int64(fieldnum) {
			fieldi++
		}

		if fieldi < len(fields) && fields[fieldi].ID == int64(fieldnum) {
			// For fields within embedded structs, ensure the embedded values aren't nil.
			index := fields[fieldi].Index
			path := make([]int, 0, len(index))
			for _, id := range index {
				path = append(path, id)
				field = sval.FieldByIndex(path)
				if field.Kind() == reflect.Ptr && field.IsNil() {
					field.Set(reflect.New(field.Type().Elem()))
				}
			}
		}

		// For more debugging output, uncomment the following three lines.
		// if fieldi < len(fields){
		//   fmt.Printf("Decoding FieldName %+v\n", fields[fieldi].Field)
		// }
		// Decode the field's value
		rem, err := de.value(wiretype, buf, field)
		if err != nil {
			if fieldi < len(fields) && fields[fieldi] != nil {
				return fmt.Errorf("Error while decoding field %+v: %v", fields[fieldi].Field, err)
			}

			return err
		}
		buf = rem
	}
	return nil
}

// Pull a value from the buffer and put it into a reflective Value.
func (de *decoder) value(wiretype int, buf []byte,
	val reflect.Value) ([]byte, error) {

	// Break out the value from the buffer based on the wire type
	var v uint64
	var n int
	var vb []byte
	switch wiretype {
	case 0: // varint
		v, n = binary.Uvarint(buf)
		if n <= 0 {
			return nil, errors.New("bad protobuf varint value")
		}
		buf = buf[n:]

	case 5: // 32-bit
		if len(buf) < 4 {
			return nil, errors.New("bad protobuf 32-bit value")
		}
		v = uint64(buf[0]) |
			uint64(buf[1])<<8 |
			uint64(buf[2])<<16 |
			uint64(buf[3])<<24
		buf = buf[4:]

	case 1: // 64-bit
		if len(buf) < 8 {
			return nil, errors.New("bad protobuf 64-bit value")
		}
		v = uint64(buf[0]) |
			uint64(buf[1])<<8 |
			uint64(buf[2])<<16 |
			uint64(buf[3])<<24 |
			uint64(buf[4])<<32 |
			uint64(buf[5])<<40 |
			uint64(buf[6])<<48 |
			uint64(buf[7])<<56
		buf = buf[8:]

	case 2: // length-delimited
		v, n = binary.Uvarint(buf)
		if n <= 0 || v > uint64(len(buf)-n) {
			return nil, errors.New(
				"bad protobuf length-delimited value")
		}
		vb = buf[n : n+int(v) : n+int(v)]
		buf = buf[n+int(v):]

	default:
		return nil, errors.New("unknown protobuf wire-type")
	}

	// We've gotten the value out of the buffer,
	// now put it into the appropriate reflective Value.
	if err := de.putvalue(wiretype, val, v, vb); err != nil {
		return nil, err
	}
	return buf, nil
}

func (de *decoder) decodeSignedInt(wiretype int, v uint64) (int64, error) {
	if wiretype == 0 { // encoded as varint
		sv := int64(v) >> 1
		if v&1 != 0 {
			sv = ^sv
		}
		return sv, nil
	} else if wiretype == 5 { // sfixed32
		return int64(int32(v)), nil
	} else if wiretype == 1 { // sfixed64
		return int64(v), nil
	} else {
		return -1, errors.New("bad wiretype for sint")
	}
}

func (de *decoder) putvalue(wiretype int, val reflect.Value,
	v uint64, vb []byte) error {
	// If val is not settable, it either represents an out-of-range field
	// or an in-range but blank (padding) field in the struct.
	// In this case, simply ignore and discard the field's content.
	if !val.CanSet() {
		return nil
	}
	switch val.Kind() {
	case reflect.Bool:
		if wiretype != 0 {
			return errors.New("bad wiretype for bool")
		}
		if v > 1 {
			return errors.New("invalid bool value")
		}
		val.SetBool(v != 0)

	case reflect.Int, reflect.Int32, reflect.Int64:
		// Signed integers may be encoded either zigzag-varint or fixed
		// Note that protobufs don't support 8- or 16-bit ints.
		if val.Kind() == reflect.Int && val.Type().Size() < 8 {
			return errors.New("detected a 32bit machine, please use either int64 or int32")
		}
		sv, err := de.decodeSignedInt(wiretype, v)
		if err != nil {
			fmt.Println("Error Reflect.Int for v=", v, "wiretype=", wiretype, "for Value=", val.Type().Name())
			return err
		}
		val.SetInt(sv)

	case reflect.Uint, reflect.Uint32, reflect.Uint64:
		// Varint-encoded 32-bit and 64-bit unsigned integers.
		if val.Kind() == reflect.Uint && val.Type().Size() < 8 {
			return errors.New("detected a 32bit machine, please use either uint64 or uint32")
		}
		if wiretype == 0 {
			val.SetUint(v)
		} else if wiretype == 5 { // ufixed32
			val.SetUint(uint64(uint32(v)))
		} else if wiretype == 1 { // ufixed64
			val.SetUint(uint64(v))
		} else {
			return errors.New("bad wiretype for uint")
		}

	case reflect.Float32:
		// Fixed-length 32-bit floats.
		if wiretype != 5 {
			return errors.New("bad wiretype for float32")
		}
		val.SetFloat(float64(math.Float32frombits(uint32(v))))

	case reflect.Float64:
		// Fixed-length 64-bit floats.
		if wiretype != 1 {
			return errors.New("bad wiretype for float64")
		}
		val.SetFloat(math.Float64frombits(v))

	case reflect.String:
		// Length-delimited string.
		if wiretype != 2 {
			return errors.New("bad wiretype for string")
		}
		val.SetString(string(vb))

	case reflect.Struct:
		// Embedded message
		if val.Type() == timeType {
			sv, err := de.decodeSignedInt(wiretype, v)
			if err != nil {
				return err
			}
			t := time.Unix(sv/int64(time.Second), sv%int64(time.Second))
			val.Set(reflect.ValueOf(t))
			return nil
		} else if enc, ok := val.Addr().Interface().(encoding.BinaryUnmarshaler); ok {
			return enc.UnmarshalBinary(vb[:])
		}
		if wiretype != 2 {
			return errors.New("bad wiretype for embedded message")
		}
		return de.message(vb, val)

	case reflect.Ptr:
		// Optional field
		// Instantiate pointer's element type.
		if val.IsNil() {
			val.Set(de.instantiate(val.Type().Elem()))
		}
		return de.putvalue(wiretype, val.Elem(), v, vb)

	case reflect.Slice, reflect.Array:
		// Repeated field or byte-slice
		if wiretype != 2 {
			return errors.New("bad wiretype for repeated field")
		}
		return de.slice(val, vb)
	case reflect.Map:
		if wiretype != 2 {
			return errors.New("bad wiretype for repeated field")
		}
		if val.IsNil() {
			// make(map[k]v):
			val.Set(reflect.MakeMap(val.Type()))
		}
		return de.mapEntry(val, vb)
	case reflect.Interface:
		data := vb[:]

		// Abstract field: instantiate via dynamic constructor.
		if val.IsNil() {
			id := GeneratorID{}
			var g InterfaceGeneratorFunc
			if len(id) < len(vb) {
				copy(id[:], vb[:len(id)])
				g = generators.get(id)
			}

			if g == nil {
				// Backwards compatible usage of the default constructors
				val.Set(de.instantiate(val.Type()))
			} else {
				// As pointers to interface are discouraged in Go, we use
				// the generator only for interface types
				data = vb[len(id):]
				val.Set(reflect.ValueOf(g()))
			}
		}

		// If the object support self-decoding, use that.
		if enc, ok := val.Interface().(encoding.BinaryUnmarshaler); ok {
			if wiretype != 2 {
				return errors.New("bad wiretype for bytes")
			}

			return enc.UnmarshalBinary(data)
		}

		// Decode into the object the interface points to.
		// XXX perhaps better ONLY to support self-decoding
		// for interface fields?
		return Decode(vb, val.Interface())

	default:
		panic("unsupported value kind " + val.Kind().String())
	}
	return nil
}

// Instantiate an arbitrary type, handling dynamic interface types.
// Returns a Ptr value.
func (de *decoder) instantiate(t reflect.Type) reflect.Value {

	// If it's an interface type, lookup a dynamic constructor for it.
	if t.Kind() == reflect.Interface {
		newfunc, ok := de.nm[t]
		if !ok {
			panic("no constructor for interface " + t.String())
		}
		return reflect.ValueOf(newfunc())
	}

	// Otherwise, for all concrete types, just instantiate directly.
	return reflect.New(t)
}

var sfixed32type = reflect.TypeOf(Sfixed32(0))
var sfixed64type = reflect.TypeOf(Sfixed64(0))
var ufixed32type = reflect.TypeOf(Ufixed32(0))
var ufixed64type = reflect.TypeOf(Ufixed64(0))

// Handle decoding of slices
func (de *decoder) slice(slval reflect.Value, vb []byte) error {
	// Find the element type, and create a temporary instance of it.
	eltype := slval.Type().Elem()
	val := reflect.New(eltype).Elem()

	// Decide on the wiretype to use for decoding.
	var wiretype int
	switch eltype.Kind() {
	case reflect.Bool, reflect.Int32, reflect.Int64, reflect.Int,
		reflect.Uint32, reflect.Uint64, reflect.Uint:
		if (eltype.Kind() == reflect.Int || eltype.Kind() == reflect.Uint) && eltype.Size() < 8 {
			return errors.New("detected a 32bit machine, please either use (u)int64 or (u)int32")
		}
		switch eltype {

		case sfixed32type:
			wiretype = 5 // Packed 32-bit representation
		case sfixed64type:
			wiretype = 1 // Packed 64-bit representation
		case ufixed32type:
			wiretype = 5 // Packed 32-bit representation
		case ufixed64type:
			wiretype = 1 // Packed 64-bit representation
		default:
			wiretype = 0 // Packed varint representation
		}

	case reflect.Float32:
		wiretype = 5 // Packed 32-bit representation

	case reflect.Float64:
		wiretype = 1 // Packed 64-bit representation

	case reflect.Uint8: // Unpacked byte-slice
		if slval.Kind() == reflect.Array {
			if slval.Len() != len(vb) {
				return errors.New("array length and buffer length differ")
			}
			for i := 0; i < slval.Len(); i++ {
				// no SetByte method in reflect so has to pass down by uint64
				slval.Index(i).SetUint(uint64(vb[i]))
			}
		} else {
			slval.SetBytes(vb)
		}
		return nil

	default: // Other unpacked repeated types
		// Just unpack and append one value from vb.
		if err := de.putvalue(2, val, 0, vb); err != nil {
			return err
		}
		if slval.Kind() != reflect.Slice {
			return errors.New("append to non-slice")
		}
		slval.Set(reflect.Append(slval, val))
		return nil
	}

	// Decode packed values from the buffer and append them to the slice.
	for len(vb) > 0 {
		rem, err := de.value(wiretype, vb, val)
		if err != nil {
			return err
		}
		slval.Set(reflect.Append(slval, val))
		vb = rem
	}
	return nil
}

// Handles the entry k,v of a map[K]V
func (de *decoder) mapEntry(slval reflect.Value, vb []byte) error {
	mKey := reflect.New(slval.Type().Key())
	mVal := reflect.New(slval.Type().Elem())
	k := mKey.Elem()
	v := mVal.Elem()
	key, n := binary.Uvarint(vb)
	if n <= 0 {
		return errors.New("bad protobuf field key")
	}
	buf := vb[n:]
	wiretype := int(key & 7)

	var err error
	buf, err = de.value(wiretype, buf, k)
	if err != nil {
		return err
	}
	for len(buf) > 0 { // for repeated values (slices etc)
		key, n = binary.Uvarint(buf)
		if n <= 0 {
			return errors.New("bad protobuf field key")
		}
		buf = buf[n:]
		wiretype = int(key & 7)
		buf, err = de.value(wiretype, buf, v)
		if err != nil {
			return err
		}
	}

	if !k.IsValid() || !v.IsValid() {
		// We did not decode the key or the value in the map entry.
		// Either way, it's an invalid map entry.
		return errors.New("proto: bad map data: missing key/val")
	}
	slval.SetMapIndex(k, v)

	return nil
}

package fixbuf // import "go.dedis.ch/fixbuf"

import (
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
	"strings"
)

// Marshaling is a basic interface representing fixed-length (or known-length)
// cryptographic objects or structures having a built-in binary encoding.
type Marshalling interface {
	// Encode the contents of this object and write it to an io.Writer.
	MarshalTo(w io.Writer) (int, error)

	// Decode the content of this object by reading from an io.Reader.
	UnmarshalFrom(r io.Reader) (int, error)
}

// Constructor represents a generic constructor
// that takes a reflect.Type, typically for an interface type,
// and constructs some suitable concrete instance of that type.
// The crypto library uses this capability to support
// dynamic instantiation of cryptographic objects of the concrete type
// appropriate for a given abstract.Suite.
type Constructor interface {
	New(t reflect.Type) interface{}
}

// BinaryEncoding represents a simple binary encoding
// suitable for reading and writing fixed-length cryptographic objects.
// The interface allows reading and writing composite types
// such as structs, arrays, and slices,
// but the encoded size of any object must be completely defined
// by the type and size of the object itself and the ciphersuite in use.
//
// Slices must be instantiated to the correct length
// before either reading or writing:
// hence the reader must determine the correct length "out of band"
// (the encoding supports no transmission of length metadata).
type BinaryEncoding struct {
	Constructor // Constructor for instantiating abstract types

	// prevent clients from depending on the exact set of fields,
	// to reserve the right to extend in backward-compatible ways.
	hidden struct{}
}

// NewBinaryEncoding returns a BinaryEncoding. The given c Constructor
// can be nil. In that case, the BinaryEncoding MUST not be used to Read() any
// structure can contains any interface fields implementing the Marshalling
// interface, because BinaryEncoding won't be able to create the concrete type
// of the field.
func NewBinaryEncoding(c Constructor) *BinaryEncoding {
	return &BinaryEncoding{Constructor: c}
}

func prindent(depth int, format string, a ...interface{}) {
	fmt.Print(strings.Repeat("  ", depth))
	fmt.Printf(format, a...)
}

type decoder struct {
	c Constructor
	r io.Reader
}

var int32Type = reflect.TypeOf(int32(0))

// Read a series of binary objects from an io.Reader.
// The objs must be a list of pointers.
func (e BinaryEncoding) Read(r io.Reader, objs ...interface{}) error {
	de := decoder{e.Constructor, r}
	for i := 0; i < len(objs); i++ {
		// XXX check that it's a by-reference type
		// (pointer, slice, etc.) and complain if not,
		// to head of accidental misuse?
		if err := de.value(reflect.ValueOf(objs[i]), 0); err != nil {
			return err
		}
	}
	return nil
}

func (de *decoder) value(v reflect.Value, depth int) error {

	// Does the object support our self-decoding interface?
	obj := v.Interface()
	if e, ok := obj.(Marshalling); ok {
		_, err := e.UnmarshalFrom(de.r)
		//prindent(depth, "decode: %s\n", e.String())
		return err
	}
	var err error
	// Otherwise, reflectively handle composite types.
	//prindent(depth, "%s: %s\n", v.Kind().String(), v.Type().String())
	switch v.Kind() {

	case reflect.Interface:
		if v.IsNil() {
			// See if we can auto-fill certain interface variables
			t := v.Type()
			o := de.c.New(t)
			if o == nil {
				panic("unsupported null pointer type: " +
					t.String())
			}
			v.Set(reflect.ValueOf(o))
		}
		fallthrough
	case reflect.Ptr:
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		return de.value(v.Elem(), depth+1)

	case reflect.Struct:
		l := v.NumField()
		for i := 0; i < l; i++ {
			if err = de.value(v.Field(i), depth+1); err != nil {
				return err
			}
		}

	case reflect.Slice:
		if v.IsNil() {
			panic("slices must be initialized to correct length before decoding")
		}
		fallthrough
	case reflect.Array:
		l := v.Len()
		for i := 0; i < l; i++ {
			if err = de.value(v.Index(i), depth+1); err != nil {
				return err
			}
		}

	case reflect.Int:
		var i int32
		err := binary.Read(de.r, binary.BigEndian, &i)
		if err != nil {
			return fmt.Errorf("Error converting int to int32 ( %v )", err)
		}
		v.SetInt(int64(i))
		return err

	case reflect.Bool:
		var b uint8
		err := binary.Read(de.r, binary.BigEndian, &b)
		v.SetBool(b != 0)
		return err

	default:

		return binary.Read(de.r, binary.BigEndian, v.Addr().Interface())
	}
	return err
}

type encoder struct {
	w io.Writer
}

// Write a data structure containing cryptographic objects,
// using their built-in binary serialization, to an io.Writer.
// Supports writing of Points, Scalars,
// basic fixed-length data types supported by encoding/binary/Write(),
// and structs, arrays, and slices containing all of these types.
func (e BinaryEncoding) Write(w io.Writer, objs ...interface{}) error {
	en := encoder{w}
	for i := 0; i < len(objs); i++ {
		if err := en.value(objs[i], 0); err != nil {
			return err
		}
	}
	return nil
}

func (en *encoder) value(obj interface{}, depth int) error {

	// Does the object support our self-decoding interface?
	if e, ok := obj.(Marshalling); ok {
		//prindent(depth, "encode: %s\n", e.String())
		_, err := e.MarshalTo(en.w)
		return err
	}

	// Otherwise, reflectively handle composite types.
	v := reflect.ValueOf(obj)
	//prindent(depth, "%s: %s\n", v.Kind().String(), v.Type().String())
	switch v.Kind() {

	case reflect.Interface:
	case reflect.Ptr:
		return en.value(v.Elem().Interface(), depth+1)

	case reflect.Struct:
		l := v.NumField()
		for i := 0; i < l; i++ {
			if err := en.value(v.Field(i).Interface(), depth+1); err != nil {
				return err
			}
		}

	case reflect.Slice, reflect.Array:
		l := v.Len()
		for i := 0; i < l; i++ {
			if err := en.value(v.Index(i).Interface(), depth+1); err != nil {
				return err
			}
		}

	case reflect.Int:
		i := int32(obj.(int))
		if int(i) != obj.(int) {
			panic("Int does not fit into int32")
		}
		return binary.Write(en.w, binary.BigEndian, i)

	case reflect.Bool:
		b := uint8(0)
		if v.Bool() {
			b = 1
		}
		return binary.Write(en.w, binary.BigEndian, b)

	default:
		// Fall back to big-endian binary encoding
		return binary.Write(en.w, binary.BigEndian, obj)
	}
	return nil
}

// Read is the default implementation of Encoding interface Read
func Read(r io.Reader, c Constructor, objs ...interface{}) error {
	return NewBinaryEncoding(c).Read(r, objs)
}

// Write is the default implementation of Encoding interface Write
func Write(w io.Writer, objs ...interface{}) error {
	return NewBinaryEncoding(nil).Write(w, objs)
}

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
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/tidwall/gjson"
)

//
/// Variant (emulates `fc::static_variant` type)
//

type Variant interface {
	Assign(typeID TypeID, impl interface{})
	Obtain(*VariantDefinition) (typeID TypeID, typeName string, impl interface{})
}

type VariantType struct {
	Name string
	Type interface{}
}

type VariantDefinition struct {
	typeIDToType   map[TypeID]reflect.Type
	typeIDToName   map[TypeID]string
	typeNameToID   map[string]TypeID
	typeIDEncoding TypeIDEncoding
}

// TypeID defines the internal representation of an instruction type ID
// (or account type, etc. in anchor programs)
// and it's used to associate instructions to decoders in the variant tracker.
type TypeID [8]byte

func (vid TypeID) Bytes() []byte {
	return vid[:]
}

// Uvarint32 parses the TypeID to a uint32.
func (vid TypeID) Uvarint32() uint32 {
	return Uvarint32FromTypeID(vid)
}

// Uint32 parses the TypeID to a uint32.
func (vid TypeID) Uint32() uint32 {
	return Uint32FromTypeID(vid, binary.LittleEndian)
}

// Uint8 parses the TypeID to a Uint8.
func (vid TypeID) Uint8() uint8 {
	return Uint8FromTypeID(vid)
}

// Equal returns true if the provided bytes are equal to
// the bytes of the TypeID.
func (vid TypeID) Equal(b []byte) bool {
	return bytes.Equal(vid.Bytes(), b)
}

// TypeIDFromBytes converts a []byte to a TypeID.
// The provided slice must be 8 bytes long or less.
func TypeIDFromBytes(slice []byte) (id TypeID) {
	// TODO: panic if len(slice) > 8 ???
	copy(id[:], slice)
	return id
}

// TypeIDFromSighash converts a sighash bytes to a TypeID.
func TypeIDFromSighash(sh []byte) TypeID {
	return TypeIDFromBytes(sh)
}

// TypeIDFromUvarint32 converts a Uvarint to a TypeID.
func TypeIDFromUvarint32(v uint32) TypeID {
	buf := make([]byte, 8)
	l := binary.PutUvarint(buf, uint64(v))
	return TypeIDFromBytes(buf[:l])
}

// TypeIDFromUint32 converts a uint32 to a TypeID.
func TypeIDFromUint32(v uint32, bo binary.ByteOrder) TypeID {
	out := make([]byte, TypeSize.Uint32)
	bo.PutUint32(out, v)
	return TypeIDFromBytes(out)
}

// TypeIDFromUint32 converts a uint8 to a TypeID.
func TypeIDFromUint8(v uint8) TypeID {
	return TypeIDFromBytes([]byte{v})
}

// Uvarint32FromTypeID parses a TypeID bytes to a uvarint 32.
func Uvarint32FromTypeID(vid TypeID) (out uint32) {
	l, _ := binary.Uvarint(vid[:])
	out = uint32(l)
	return out
}

// Uint32FromTypeID parses a TypeID bytes to a uint32.
func Uint32FromTypeID(vid TypeID, order binary.ByteOrder) (out uint32) {
	out = order.Uint32(vid[:])
	return out
}

// Uint32FromTypeID parses a TypeID bytes to a uint8.
func Uint8FromTypeID(vid TypeID) (out uint8) {
	return vid[0]
}

type TypeIDEncoding uint32

const (
	Uvarint32TypeIDEncoding TypeIDEncoding = iota
	Uint32TypeIDEncoding
	Uint8TypeIDEncoding
	// AnchorTypeIDEncoding is the instruction ID encoding used by programs
	// written using the anchor SDK.
	// The typeID is the sighash of the instruction.
	AnchorTypeIDEncoding
	// No type ID; ONLY ONE VARIANT PER PROGRAM.
	NoTypeIDEncoding
)

var NoTypeIDDefaultID = TypeIDFromUint8(0)

// NewVariantDefinition creates a variant definition based on the *ordered* provided types.
//
//   - For anchor instructions, it's the name that defines the binary variant value.
//   - For all other types, it's the ordering that defines the binary variant value just like in native `nodeos` C++
//     and in Smart Contract via the `std::variant` type. It's important to pass the entries
//     in the right order!
//
// This variant definition can now be passed to functions of `BaseVariant` to implement
// marshal/unmarshaling functionalities for binary & JSON.
func NewVariantDefinition(typeIDEncoding TypeIDEncoding, types []VariantType) (out *VariantDefinition) {
	if len(types) < 0 {
		panic("it's not valid to create a variant definition without any types")
	}

	typeCount := len(types)
	out = &VariantDefinition{
		typeIDEncoding: typeIDEncoding,
		typeIDToType:   make(map[TypeID]reflect.Type, typeCount),
		typeIDToName:   make(map[TypeID]string, typeCount),
		typeNameToID:   make(map[string]TypeID, typeCount),
	}

	switch typeIDEncoding {
	case Uvarint32TypeIDEncoding:
		for i, typeDef := range types {
			typeID := TypeIDFromUvarint32(uint32(i))

			// FIXME: Check how the reflect.Type is used and cache all its usage in the definition.
			//        Right now, on each Unmarshal, we re-compute some expensive stuff that can be
			//        re-used like the `typeGo.Elem()` which is always the same. It would be preferable
			//        to have those already pre-defined here so we can actually speed up the
			//        Unmarshal code.
			out.typeIDToType[typeID] = reflect.TypeOf(typeDef.Type)
			out.typeIDToName[typeID] = typeDef.Name
			out.typeNameToID[typeDef.Name] = typeID
		}
	case Uint32TypeIDEncoding:
		for i, typeDef := range types {
			typeID := TypeIDFromUint32(uint32(i), binary.LittleEndian)

			// FIXME: Check how the reflect.Type is used and cache all its usage in the definition.
			//        Right now, on each Unmarshal, we re-compute some expensive stuff that can be
			//        re-used like the `typeGo.Elem()` which is always the same. It would be preferable
			//        to have those already pre-defined here so we can actually speed up the
			//        Unmarshal code.
			out.typeIDToType[typeID] = reflect.TypeOf(typeDef.Type)
			out.typeIDToName[typeID] = typeDef.Name
			out.typeNameToID[typeDef.Name] = typeID
		}
	case Uint8TypeIDEncoding:
		for i, typeDef := range types {
			typeID := TypeIDFromUint8(uint8(i))

			// FIXME: Check how the reflect.Type is used and cache all its usage in the definition.
			//        Right now, on each Unmarshal, we re-compute some expensive stuff that can be
			//        re-used like the `typeGo.Elem()` which is always the same. It would be preferable
			//        to have those already pre-defined here so we can actually speed up the
			//        Unmarshal code.
			out.typeIDToType[typeID] = reflect.TypeOf(typeDef.Type)
			out.typeIDToName[typeID] = typeDef.Name
			out.typeNameToID[typeDef.Name] = typeID
		}
	case AnchorTypeIDEncoding:
		for _, typeDef := range types {
			typeID := TypeIDFromSighash(Sighash(SIGHASH_GLOBAL_NAMESPACE, typeDef.Name))

			// FIXME: Check how the reflect.Type is used and cache all its usage in the definition.
			//        Right now, on each Unmarshal, we re-compute some expensive stuff that can be
			//        re-used like the `typeGo.Elem()` which is always the same. It would be preferable
			//        to have those already pre-defined here so we can actually speed up the
			//        Unmarshal code.
			out.typeIDToType[typeID] = reflect.TypeOf(typeDef.Type)
			out.typeIDToName[typeID] = typeDef.Name
			out.typeNameToID[typeDef.Name] = typeID
		}
	case NoTypeIDEncoding:
		if len(types) != 1 {
			panic(fmt.Sprintf("NoTypeIDEncoding can only have one variant type definition, got %v", len(types)))
		}
		typeDef := types[0]

		typeID := NoTypeIDDefaultID

		// FIXME: Check how the reflect.Type is used and cache all its usage in the definition.
		//        Right now, on each Unmarshal, we re-compute some expensive stuff that can be
		//        re-used like the `typeGo.Elem()` which is always the same. It would be preferable
		//        to have those already pre-defined here so we can actually speed up the
		//        Unmarshal code.
		out.typeIDToType[typeID] = reflect.TypeOf(typeDef.Type)
		out.typeIDToName[typeID] = typeDef.Name
		out.typeNameToID[typeDef.Name] = typeID

	default:
		panic(fmt.Errorf("unsupported TypeIDEncoding: %v", typeIDEncoding))
	}

	return out
}

func (d *VariantDefinition) TypeID(name string) TypeID {
	id, found := d.typeNameToID[name]
	if !found {
		knownNames := make([]string, len(d.typeNameToID))
		i := 0
		for name := range d.typeNameToID {
			knownNames[i] = name
			i++
		}

		panic(fmt.Errorf("trying to use an unknown type name %q, known names are %q", name, strings.Join(knownNames, ", ")))
	}

	return id
}

type VariantImplFactory = func() interface{}
type OnVariant = func(impl interface{}) error

type BaseVariant struct {
	TypeID TypeID
	Impl   interface{}
}

var _ Variant = &BaseVariant{}

func (a *BaseVariant) Assign(typeID TypeID, impl interface{}) {
	a.TypeID = typeID
	a.Impl = impl
}

func (a *BaseVariant) Obtain(def *VariantDefinition) (typeID TypeID, typeName string, impl interface{}) {
	return a.TypeID, def.typeIDToName[a.TypeID], a.Impl
}

func (a *BaseVariant) MarshalJSON(def *VariantDefinition) ([]byte, error) {
	typeName, found := def.typeIDToName[a.TypeID]
	if !found {
		return nil, fmt.Errorf("type %d is not know by variant definition", a.TypeID)
	}

	return json.Marshal([]interface{}{typeName, a.Impl})
}

func (a *BaseVariant) UnmarshalJSON(data []byte, def *VariantDefinition) error {
	typeResult := gjson.GetBytes(data, "0")
	implResult := gjson.GetBytes(data, "1")

	if !typeResult.Exists() || !implResult.Exists() {
		return fmt.Errorf("invalid format, expected '[<type>, <impl>]' pair, got %q", string(data))
	}

	typeName := typeResult.String()
	typeID, found := def.typeNameToID[typeName]
	if !found {
		return fmt.Errorf("type %q is not know by variant definition", typeName)
	}

	typeGo := def.typeIDToType[typeID]
	if typeGo == nil {
		return fmt.Errorf("no known type for %q", typeName)
	}

	a.TypeID = typeID

	if typeGo.Kind() == reflect.Ptr {
		a.Impl = reflect.New(typeGo.Elem()).Interface()
		if err := json.Unmarshal([]byte(implResult.Raw), a.Impl); err != nil {
			return err
		}
	} else {
		// This is not the most optimal way of doing things for "value"
		// types (over "pointer" types) as we always allocate a new pointer
		// element, unmarshal it and then either keep the pointer type or turn
		// it into a value type.
		//
		// However, in non-reflection based code, one would do like this and
		// avoid an `new` memory allocation:
		//
		// ```
		// name := eos.Name("")
		// json.Unmarshal(data, &name)
		// ```
		//
		// This would work without a problem. In reflection code however, I
		// did not find how one can go from `reflect.Zero(typeGo)` (which is
		// the equivalence of doing `name := eos.Name("")`) and take the
		// pointer to it so it can be unmarshalled correctly.
		//
		// A played with various iteration, and nothing got it working. Maybe
		// the next step would be to explore the `unsafe` package and obtain
		// an unsafe pointer and play with it.
		value := reflect.New(typeGo)
		if err := json.Unmarshal([]byte(implResult.Raw), value.Interface()); err != nil {
			return err
		}

		a.Impl = value.Elem().Interface()
	}

	return nil
}

func (a *BaseVariant) UnmarshalBinaryVariant(decoder *Decoder, def *VariantDefinition) (err error) {
	var typeID TypeID
	switch def.typeIDEncoding {
	case Uvarint32TypeIDEncoding:
		val, err := decoder.ReadUvarint32()
		if err != nil {
			return fmt.Errorf("uvarint32: unable to read variant type id: %s", err)
		}
		typeID = TypeIDFromUvarint32(val)
	case Uint32TypeIDEncoding:
		val, err := decoder.ReadUint32(binary.LittleEndian)
		if err != nil {
			return fmt.Errorf("uint32: unable to read variant type id: %s", err)
		}
		typeID = TypeIDFromUint32(val, binary.LittleEndian)
	case Uint8TypeIDEncoding:
		id, err := decoder.ReadUint8()
		if err != nil {
			return fmt.Errorf("uint8: unable to read variant type id: %s", err)
		}
		typeID = TypeIDFromBytes([]byte{id})
	case AnchorTypeIDEncoding:
		typeID, err = decoder.ReadTypeID()
		if err != nil {
			return fmt.Errorf("anchor: unable to read variant type id: %s", err)
		}
	case NoTypeIDEncoding:
		typeID = NoTypeIDDefaultID
	}

	a.TypeID = typeID

	typeGo := def.typeIDToType[typeID]
	if typeGo == nil {
		return fmt.Errorf("no known type for type %d", typeID)
	}

	if typeGo.Kind() == reflect.Ptr {
		a.Impl = reflect.New(typeGo.Elem()).Interface()
		if err = decoder.Decode(a.Impl); err != nil {
			return fmt.Errorf("unable to decode variant type %d: %s", typeID, err)
		}
	} else {
		// This is not the most optimal way of doing things for "value"
		// types (over "pointer" types) as we always allocate a new pointer
		// element, unmarshal it and then either keep the pointer type or turn
		// it into a value type.
		//
		// However, in non-reflection based code, one would do like this and
		// avoid an `new` memory allocation:
		//
		// ```
		// name := eos.Name("")
		// json.Unmarshal(data, &name)
		// ```
		//
		// This would work without a problem. In reflection code however, I
		// did not find how one can go from `reflect.Zero(typeGo)` (which is
		// the equivalence of doing `name := eos.Name("")`) and take the
		// pointer to it so it can be unmarshalled correctly.
		//
		// A played with various iteration, and nothing got it working. Maybe
		// the next step would be to explore the `unsafe` package and obtain
		// an unsafe pointer and play with it.
		value := reflect.New(typeGo)
		if err = decoder.Decode(value.Interface()); err != nil {
			return fmt.Errorf("unable to decode variant type %d: %s", typeID, err)
		}

		a.Impl = value.Elem().Interface()
	}
	return nil
}

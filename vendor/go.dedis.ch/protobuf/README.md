# Reflection-based Protocol Buffers

Package protobuf implements Protocol Buffers reflectively using Go types to
define message formats. This approach provides convenience similar to Gob
encoding, but with a widely-used and language-neutral wire format.
For detailed API documentation see https://godoc.org/go.dedis.ch/protobuf.
For general information on Protocol buffers see
https://developers.google.com/protocol-buffers.

## Features

- Reflection-based encoding and decoding to/from protocol buffer wire format.
- Use Go struct field tags to control protobuf fields (ID, optional/required, names).
- Generate `.proto` files from Go structures.
- Encode `time.Time` as an `sfixed64` UnixNano.
- Support for enums.

## Details

In contrast with goprotobuf, this package does not require users to write or
compile .proto files; you just define the message formats you want as Go
struct types. Consider for example this example message format definition from
the Protocol Buffers overview:

```protobuf
message Person {
  required string name = 1;
  required int32  id = 2;
  optional string email = 3;

  enum PhoneType {
    MOBILE = 0;
    HOME = 1;
    WORK = 2;
  }

  message PhoneNumber {
    required string    number = 1;
    optional PhoneType type = 2;
  }

  repeated PhoneNumber phone = 4;
}
```

The following Go type and const definitions express exactly the same format,
for purposes of encoding and decoding with this protobuf package:

```go
type Person struct {
  Name  string
  Id    int32
  Email *string
  Phone []PhoneNumber
}

type PhoneType uint32
const (
  MobilePhone PhoneType = iota
  HomePhone
  WorkPhone
)

type PhoneNumber struct {
  Number string
  Type *PhoneType
}
```

To encode a message, you simply call the Encode() function
with a pointer to the struct you wish to encode, and
Encode() returns a []byte slice containing the protobuf-encoded struct:

```go
person := Person{...}
buf := Encode(&person)
output.Write(buf)
```

To decode an encoded message, simply call Decode() on the byte-slice:

```go
err := Decode(buf,&person)
if err != nil {
  panic("Decode failed: "+err.Error())
}
```

If you want to interoperate with code in other languages
using the same message formats, you may of course still end up writing
.proto files for the code in those other languages.

However, defining message formats with native Go types enables these types
to be tailored to the code using them without affecting wire compatibility,
such as by attaching useful methods to these struct types.

## Translation Rules

The translation between a Go struct definition and a basic Protocol Buffers
message format definition is straightforward; the rules are as follows.

Go [field tags](https://golang.org/ref/spec#Struct_types) with the key
"protobuf" may be used to control naming, IDs, and optional/required state.
The options are comma-separated like so:

```go
type Tags struct {
  Field1 string `protobuf:"10,req,field_1"`
  Field2 int32  `protobuf:"20,opt,field_2"`
}
```

A message definition in a .proto file translates to a Go struct, whose fields
are implicitly assigned consecutive numbers starting from 1.

```go
type Padded struct {
  Field1 string                  // = 1
  Field2 int32  `protobuf:"3"`   // = 3
}
```

A 'required' protobuf field translates to a plain field of a corresponding
type in the Go struct. The following table summarizes the correspondence
between .proto definition types and Go field types:

Protobuf  |  Go
----------|---------
bool      | bool
enum      | Enum
int32     | uint32
int64     | uint64
uint32    | uint32
uint64    | uint64
sint32    | int32
sint64    | int64
fixed32   | Ufixed32
fixed64   | Ufixed64
sfixed32  | Sfixed32
sfixed64  | Sfixed64
float     | float32
double    | float64
string    | string
bytes     | []byte
message   | struct

An 'optional' protobuf field is expressed as a pointer field in Go.
Encode() will transmit the field only if the pointer is non-nil.
Decode() will instantiate the pointed-to type and fill in the pointer
if the field is present in the message being decoded,
leaving the pointer unmodified (usually nil) if the field is not present.

A 'repeated' protobuf field translates to a slice field in Go.
Slices of primitive bool, integer, and float types are encoded
and decoded in packed format, as if the [packed=true] option
was declared for the field in the .proto file.

For flexibility and convenience, struct fields may have interface types,
which this package interprets as having dynamic types to be bound at runtime.
Encode() follows the interface's implicit pointer and uses reflection
to determine the referred-to object's actual type for encoding
Decode() takes an optional map of interface types to constructor functions,
which it uses to instantiate concrete types for interfaces while decoding.
Furthermore, if the instantiated types support the Encoding interface,
Encode() and Decode() will invoke the methods of that interface,
allowing objects to implement their own custom encoding/decoding methods.

This package does not try to support all possible protobuf formats. It
currently does not support nonzero default value declarations for enums, the
legacy unpacked formats for repeated numeric fields, messages with extremely
sparse field numbering, or other more exotic features like extensions or
oneof. If you need to interoperate with existing protobuf code using these
features, then you should probably use goprotobuf, at least for those
particular message formats.
Many of these limitations could be fixed by creative use of
struct tag metadata (see https://golang.org/ref/spec#Struct_types).

Another downside of this reflective approach to protobuf implementation is
that reflective code is generally less efficient than statically generated
code, as gogoprotobuf produces for example. If we decide we want the
convenience of format definitions in Go with the runtime performance of static
code generation, we could in principle achieve that by adding a "Go-format"
message format compiler frontend to goprotobuf or gogoprotobuf - but we leave
this as an exercise for the reader.


## Generating .proto files

`.proto` files can be generated from Go structs using the
`GenerateProtobufDefinition()` function. The following:

```go
types := []interface{}{
  Person{},
  PhoneNumber{},
}
enums := EnumMap{
  "MobilePhone": MobilePhone,
  "HomePhone": HomePhone,
  "WorkPhone": WorkPhone,
}
GenerateProtobufDefinition(w, types, enums, nil)
```

Will generate:

```protobuf
message Person {
  required string name = 1;
  required sint32 id = 2;
  optional string email = 3;
  repeated PhoneNumber phone = 4;
}

message PhoneNumber {
  required string number = 1;
  optional uint32 type = 2;
}
```

Note: It can be quite tedious to manually synchronise the type and enum maps with the types in your package. I've found [pkgreflect](https://github.com/ungerik/pkgreflect) very useful for automating this.

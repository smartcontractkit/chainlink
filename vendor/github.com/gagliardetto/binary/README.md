# binary


### Borsh

#### Decoding borsh

```golang
 dec := bin.NewBorshDecoder(data)
 var meta token_metadata.Metadata
 err = dec.Decode(&meta)
 if err != nil {
   panic(err)
 }
```

#### Encoding borsh

```golang
buf := new(bytes.Buffer)
enc := bin.NewBorshEncoder(buf)
err := enc.Encode(meta)
if err != nil {
  panic(err)
}
// fmt.Print(buf.Bytes())
```

### Optional Types

```golang
type Person struct {
	Name string
	Age  uint8 `bin:"optional"`
}
```

Rust equivalent:
```rust
struct Person {
    name: String,
    age: Option<u8>
}
```

### Enum Types

```golang
type MyEnum struct {
	Enum  bin.BorshEnum `borsh_enum:"true"`
	One   bin.EmptyVariant
	Two   uint32
	Three int16
}
```

Rust equivalent:
```rust
enum MyEnum {
    One,
    Two(u32),
    Three(i16),
}
```

### Exported vs Unexported Fields

In this example, the `two` field will be skipped by the encoder/decoder because the
field is not exported.
```golang
type MyStruct struct {
	One   string
	two   uint32
	Three int16
}
```

### Skip Decoding/Encoding Attributes

Encoding/Decoding of exported fields can be skipped using the `borsh_skip` tag.
```golang
type MyStruct struct {
	One   string
	Two   uint32 `borsh_skip:"true"`
	Three int16
}
```

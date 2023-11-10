//go:build ignore

package main

import (
	"bytes"
	"os"
	"text/template"
)

func main() {
	genInts()
	genBytes()
}

func genBytes() {
	byteTypes := [33]ByteType{}
	for i := 1; i < 33; i++ {
		byteTypes[i-1].Size = i
	}
	runTemplate("bytes", bytesTemplate, "byte_types_gen.go", byteTypes)
}

func genInts() {
	var intTypes []*IntType

	// 8, 16, 32, and 64 bits have their own type in go that is used by abi.
	// The test use *big.Int
	for i := 24; i <= 256; i += 8 {
		if i == 32 || i == 64 {
			continue
		}

		signed := &IntType{Size: i, Signed: false}
		unsigned := &IntType{Prefix: "u", Size: i}
		intTypes = append(intTypes, signed, unsigned)
	}
	runTemplate("ints", intsTemplate, "int_types_gen.go", intTypes)
}

func runTemplate(name, rawTemplate, outputFile string, input any) {
	t, err := template.New(name).Parse(rawTemplate)
	if err != nil {
		panic(err)
	}

	br := bytes.Buffer{}
	if err = t.Execute(&br, input); err != nil {
		panic(err)
	}

	if err = os.WriteFile(outputFile, br.Bytes(), 0777); err != nil {
		panic(err)
	}
}

type IntType struct {
	Prefix string
	Size   int
	Signed bool
}

type ByteType struct {
	Size int
}

const bytesTemplate = `
package types

import "reflect"

{{ range . }}
type bytes{{.Size}} [{{.Size}}]byte
func init() {
	typeMap["bytes{{.Size}}"] = &AbiEncodingType {
		Native: reflect.TypeOf([{{.Size}}]byte{}),
		Checked: reflect.TypeOf(bytes{{.Size}}{}),
	}
}

{{ end }}
`

const intsTemplate = `
package types

import (
	"math/big"
	"reflect"

	"github.com/fxamacker/cbor/v2"
	"github.com/smartcontractkit/chainlink-relay/pkg/types"
)

type SizedBigInt interface {
	Verify() error
	private()
}

var sizedBigIntType = reflect.TypeOf((*SizedBigInt)(nil)).Elem()
func SizedBigIntType() reflect.Type {
	return sizedBigIntType
}

{{ range . }}
type {{.Prefix}}int{{.Size}} big.Int
func (i *{{.Prefix}}int{{.Size}}) UnmarshalCBOR(input []byte) error {
	bi := (*big.Int)(i)
	if err := cbor.Unmarshal(input, bi); err != nil {
		return err
	}

	return i.Verify()
}

func (i *{{.Prefix}}int{{.Size}}) UnmarshalText(input []byte) error {
	bi := (*big.Int)(i)
	if _, ok := bi.SetString(string(input), 10); !ok {
		return types.InvalidTypeError{}
	}

	return i.Verify()
}

func (i *{{.Prefix}}int{{.Size}}) Verify() error {
	bi := (*big.Int)(i)
	{{ if .Signed }}
	if !codec.FitsInNBitsSigned({{.Size}}, bi) {
		return types.InvalidTypeError{}
	}
	{{ else }}
	if bi.BitLen() > {{.Size}} {
		return types.InvalidTypeError{}
	}
	{{ end }}
	return nil
}

func (i *{{.Prefix}}int{{.Size}}) private() {}

func init() {
	typeMap["{{.Prefix}}int{{.Size}}"] = &AbiEncodingType {
		Native: reflect.TypeOf((*big.Int)(nil)),
		Checked: reflect.TypeOf((*{{.Prefix}}int{{.Size}})(nil)),
	}
}
{{ end }}
`

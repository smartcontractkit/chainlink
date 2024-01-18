package main

import (
	"bytes"
	_ "embed"
	"go/format"
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
	mustRunTemplate("bytes", bytesTemplate, "byte_types_gen.go", byteTypes)
}

func genInts() {
	var intTypes []*IntType

	// 8, 16, 32, and 64 bits have their own type in go that is used by abi.
	// The test use *big.Int
	for i := 24; i <= 256; i += 8 {
		if i == 32 || i == 64 {
			continue
		}

		signed := &IntType{Size: i, Signed: true}
		unsigned := &IntType{Prefix: "u", Size: i}
		intTypes = append(intTypes, signed, unsigned)
	}
	mustRunTemplate("ints", intsTemplate, "int_types_gen.go", intTypes)
}

func mustRunTemplate(name, rawTemplate, outputFile string, input any) {
	t := template.Must(template.New(name).Parse(rawTemplate))

	br := bytes.Buffer{}
	if err := t.Execute(&br, input); err != nil {
		panic(err)
	}

	res, err := format.Source(br.Bytes())
	if err != nil {
		panic(err)
	}

	if err = os.WriteFile(outputFile, res, 0600); err != nil {
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

//go:embed bytes.go.tmpl
var bytesTemplate string

//go:embed ints.go.tmpl
var intsTemplate string

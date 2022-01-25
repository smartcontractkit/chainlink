// Package monitoring contains a small DSL to help write more robust Avro schemas
// by taking advantage of go's type system.
package avro

import "encoding/json"

type Schema interface {
	IsSchema()
}

// Primitive types

type primitive struct {
	Typ string
}

func (p primitive) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.Typ)
}

var (
	Null    = primitive{"null"}
	Boolean = primitive{"boolean"}
	Int     = primitive{"int"}
	Long    = primitive{"long"}
	Double  = primitive{"double"}
	Bytes   = primitive{"bytes"}
	String  = primitive{"string"}
)

// Complex types

// Opts represents the optional fields of a complex type.
type Opts struct {
	Doc       string
	Namespace string
	Default   interface{}
}

type record struct {
	Name      string `json:"name"`
	Typ       string `json:"type"`
	Namespace string `json:"namespace,omitempty"`
	Doc       string `json:"doc,omitempty"`
	Fields    Fields `json:"fields"`
}

func Record(name string, opts Opts, fields Fields) Schema {
	return record{
		name,
		"record",
		opts.Namespace,
		opts.Doc,
		fields,
	}
}

type field struct {
	Name    string      `json:"name"`
	Doc     string      `json:"doc,omitempty"`
	Typ     Schema      `json:"type"`
	Default interface{} `json:"default,omitempty"`
}

type IField interface {
	IsField()
}

type Fields []IField

func Field(name string, opts Opts, typ Schema) IField {
	return field{
		name,
		opts.Doc,
		typ,
		opts.Default,
	}
}

type array struct {
	Typ   string `json:"type"`
	Items Schema `json:"items"`
}

func Array(items Schema) Schema {
	return array{
		"array",
		items,
	}
}

type Union []Schema

// Type checking

func (p primitive) IsSchema() {}
func (r record) IsSchema()    {}
func (a array) IsSchema()     {}
func (u Union) IsSchema()     {}

func (f field) IsField() {}

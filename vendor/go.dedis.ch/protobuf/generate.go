package protobuf

import (
	"fmt"
	"io"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"text/template"
	"errors"
)

const protoTemplate = `[[range $name, $values := .Enums]]
enum [[$name|$.Renamer.TypeName]] {[[range $values]]
  [[.Name|$.Renamer.ConstName]] = [[.Value]];[[end]]
}

[[end]][[range .Types]]
message [[.Name|$.Renamer.TypeName]] {[[range .|Fields]]
  [[.|TypeName]] [[.|$.Renamer.FieldName]] = [[.ID]][[.|Options]];[[end]]
}
[[end]]
`

var splitName = regexp.MustCompile(`((?:ID)|(?:[A-Z][a-z_0-9]+)|([\w\d]+))`)

func typeIndirect(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

func typeName(f ProtoField, enums enumTypeMap, renamer GeneratorNamer) (s string) {
	defer func() {
		if e := recover(); e != nil {
			s = ""
			panic(e.(string))
		}
	}()
	t := f.Field.Type
	if t.Kind() == reflect.Slice {
		if t.Elem().Kind() == reflect.Uint8 {
			return fieldPrefix(f, TagNone) + "bytes"
		}
		return "repeated " + innerTypeName(typeIndirect(t.Elem()), enums, renamer)
	}
	if t.Kind() == reflect.Ptr {
		return fieldPrefix(f, TagOptional) + innerTypeName(t.Elem(), enums, renamer)
	}
	return fieldPrefix(f, TagNone) + innerTypeName(t, enums, renamer)
}

func fieldPrefix(f ProtoField, def TagPrefix) string {
	opt := def
	if def == TagNone {
		opt = f.Prefix
	}
	switch opt {
	case TagOptional:
		return "optional "
	case TagRequired:
		return "required "
	default:
		if f.Field.Type.Kind() == reflect.Ptr {
			return "optional "
		}
		return "required "
	}
}

func innerTypeName(t reflect.Type, enums enumTypeMap, renamer GeneratorNamer) string {
	if (t.Kind() == reflect.Slice || t.Kind() == reflect.Array) && t.Elem().Kind() == reflect.Uint8 {
		return "bytes"
	}
	if t.PkgPath() == "time" {
		if t.Name() == "Time" {
			return "sfixed64"
		}
		if t.Name() == "Duration" {
			return "sint64"
		}
	}
	switch t.Name() {
	case "Ufixed32":
		return "fixed32"
	case "Ufixed64":
		return "ufixed64"
	case "Sfixed32":
		return "sfixed32"
	case "Sfixed64":
		return "sfixed64"
	}

	if _, ok := enums[t.Name()]; ok {
		return renamer.TypeName(t.Name())
	}

	switch t.Kind() {
	case reflect.Float64:
		return "double"
	case reflect.Float32:
		return "float"
	case reflect.Int32:
		return "sint32"
	case reflect.Int, reflect.Int64:
		return "sint64"
	case reflect.Bool:
		return "bool"
	case reflect.Uint32:
		return "uint32"
	case reflect.Uint, reflect.Uint64:
		return "uint64"
	case reflect.String:
		return "string"
	case reflect.Struct:
		return t.Name()
	case reflect.Map:
		// we have to do this again (otherwise we'll end up with an empty name for the value):
		var valTypeName string
		valType := t.Elem()
		if valType.Kind() == reflect.Slice {
			if valType.Elem().Kind() == reflect.Uint8 {
				valTypeName = "bytes"
			} else {
				valTypeName = innerTypeName(typeIndirect(valType.Elem()), enums, renamer)
			}
		} else if valType.Kind() == reflect.Ptr {
			valTypeName = innerTypeName(valType.Elem(), enums, renamer)
		} else {
			// here we can just use the value's type:
			valTypeName = innerTypeName(valType, enums, renamer)
		}
		return fmt.Sprintf("map<%s, %s>", innerTypeName(t.Key(), enums, renamer), valTypeName)
	default:
		panic("unsupported type " + t.Name())
	}
}

func options(f ProtoField) string {
	if f.Field.Type.Kind() == reflect.Slice {
		switch f.Field.Type.Elem().Kind() {
		case reflect.Bool,
			reflect.Int32, reflect.Int64,
			reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64:
			return " [packed=true]"
		}
	}
	return ""
}

type GeneratorNamer interface {
	FieldName(ProtoField) string
	TypeName(name string) string
	ConstName(name string) string
}

// DefaultGeneratorNamer renames symbols when mapping from Go to .proto files.
//
// The rules are:
// - Field names are mapped from SomeFieldName to some_field_name.
// - Type names are not modified.
// - Constants are mapped form SomeConstantName to SOME_CONSTANT_NAME.
type DefaultGeneratorNamer struct{}

func (d *DefaultGeneratorNamer) FieldName(f ProtoField) string {
	if f.Name != "" {
		return f.Name
	}
	parts := splitName.FindAllString(f.Field.Name, -1)
	for i := range parts {
		parts[i] = strings.ToLower(parts[i])
	}
	return strings.Join(parts, "_")
}

func (d *DefaultGeneratorNamer) TypeName(name string) string {
	return name
}

func (d *DefaultGeneratorNamer) ConstName(name string) string {
	parts := splitName.FindAllString(name, -1)
	for i := range parts {
		parts[i] = strings.ToUpper(parts[i])
	}
	return strings.Join(parts, "_")
}

type reflectedTypes []reflect.Type

func (r reflectedTypes) Len() int           { return len(r) }
func (r reflectedTypes) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r reflectedTypes) Less(i, j int) bool { return r[i].Name() < r[j].Name() }

type EnumMap map[string]interface{}

type enumValue struct {
	Name  string
	Value Enum
}

type enumValues []enumValue

func (e enumValues) Len() int           { return len(e) }
func (e enumValues) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }
func (e enumValues) Less(i, j int) bool { return e[i].Value < e[j].Value }

type enumTypeMap map[string]enumValues

// GenerateProtobufDefinition generates a .proto file from a list of structs via reflection.
// fieldNamer is a function that maps ProtoField types to generated protobuf field names.
func GenerateProtobufDefinition(w io.Writer, types []interface{}, enumMap EnumMap, renamer GeneratorNamer) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(e.(string))
		}
	}()
	enums := enumTypeMap{}
	for name, value := range enumMap {
		v := reflect.ValueOf(value)
		t := v.Type()
		if t.Kind() != reflect.Uint32 {
			return fmt.Errorf("enum type aliases must be uint32")
		}
		if t.Name() == "uint32" {
			return fmt.Errorf("enum value must be a type alias, but got uint32")
		}
		enums[t.Name()] = append(enums[t.Name()], enumValue{name, Enum(v.Uint())})
	}
	for _, values := range enums {
		sort.Sort(values)
	}
	rt := reflectedTypes{}
	for _, t := range types {
		typ := reflect.Indirect(reflect.ValueOf(t)).Type()
		if typ.Kind() != reflect.Struct {
			continue
		}
		rt = append(rt, typ)
	}
	sort.Sort(rt)
	if renamer == nil {
		renamer = &DefaultGeneratorNamer{}
	}
	t := template.Must(template.New("protobuf").Funcs(template.FuncMap{
		"Fields":   ProtoFields,
		"TypeName": func(f ProtoField) string { return typeName(f, enums, renamer) },
		"Options":  options,
	}).Delims("[[", "]]").Parse(protoTemplate))
	return t.Execute(w, map[string]interface{}{
		"Renamer": renamer,
		"Enums":   enums,
		"Types":   rt,
		"Ptr":     reflect.Ptr,
		"Slice":   reflect.Slice,
		"Map":     reflect.Map,
	})
}

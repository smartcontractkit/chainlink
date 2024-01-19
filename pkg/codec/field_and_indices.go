package codec

import (
	"fmt"
	"reflect"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

type fieldsAndIndices struct {
	pkgPath   string
	fields    []reflect.StructField
	Indices   map[string]int
	subFields map[string]*fieldsAndIndices
	transform func(reflect.Type) reflect.Type
}

func getFieldIndices(inputType reflect.Type) (*fieldsAndIndices, error) {
	typeTransform := func(t reflect.Type) reflect.Type { return t }
	for ; inputType.Kind() != reflect.Struct; inputType = inputType.Elem() {
		tmp := typeTransform
		switch inputType.Kind() {
		case reflect.Ptr:
			typeTransform = func(t reflect.Type) reflect.Type { return reflect.PtrTo(tmp(t)) }
		case reflect.Slice:
			typeTransform = func(t reflect.Type) reflect.Type { return reflect.SliceOf(tmp(t)) }
		case reflect.Array:
			typeTransform = func(t reflect.Type) reflect.Type { return reflect.ArrayOf(inputType.Len(), tmp(t)) }
		default:
			return nil, fmt.Errorf("%w: cannot get field index from kind %v", types.ErrInvalidType, inputType.Kind())
		}
	}
	length := inputType.NumField()
	fields := make([]reflect.StructField, length)
	Indices := map[string]int{}

	for i := 0; i < length; i++ {
		field := inputType.Field(i)
		Indices[field.Name] = i
		fields[i] = field
	}

	pkgPath := inputType.PkgPath()
	// types created by reflection may not have a pkgPath
	if pkgPath == "" {
		pkgPath = "github.com/smartcontractkit/chainlink-common/pkg/codec"
	}
	return &fieldsAndIndices{
		pkgPath:   pkgPath,
		fields:    fields,
		Indices:   Indices,
		subFields: map[string]*fieldsAndIndices{},
		transform: typeTransform,
	}, nil
}

func (f *fieldsAndIndices) fieldByName(name string) (*reflect.StructField, bool) {
	if index, ok := f.Indices[name]; ok {
		return &f.fields[index], true
	}
	return nil, false
}

func (f *fieldsAndIndices) populateSubFields(field string) (*fieldsAndIndices, error) {
	if subField, ok := f.subFields[field]; ok {
		return subField, nil
	} else if index, ok := f.Indices[field]; ok {
		fi, err := getFieldIndices(f.fields[index].Type)
		if err != nil {
			return nil, err
		}
		f.subFields[field] = fi
		return fi, nil
	}

	return nil, fmt.Errorf("%w: cannot find field %s", types.ErrInvalidType, field)
}

func (f *fieldsAndIndices) makeNewType() reflect.Type {
	for key, subField := range f.subFields {
		f.fields[f.Indices[key]].Type = subField.makeNewType()
	}
	return f.transform(reflect.StructOf(f.fields))
}

func (f *fieldsAndIndices) updateTypeFromSubkeyMods(key string) {
	if subField, ok := f.subFields[key]; ok {
		f.fields[f.Indices[key]].Type = subField.makeNewType()
		delete(f.subFields, key)
	}
}

func (f *fieldsAndIndices) addNewField(field reflect.StructField) {
	f.fields = append(f.fields, field)
	f.Indices[field.Name] = len(f.fields) - 1
}

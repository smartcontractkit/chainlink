package gen

import (
	"fmt"
	"reflect"

	"github.com/leanovate/gopter"
)

// Struct generates a given struct type.
// rt has to be the reflect type of the struct, gens contains a map of field generators.
// Note that the result types of the generators in gen have to match the type of the correspoinding
// field in the struct. Also note that only public fields of a struct can be generated
func Struct(rt reflect.Type, gens map[string]gopter.Gen) gopter.Gen {
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}
	if rt.Kind() != reflect.Struct {
		return Fail(rt)
	}
	fieldGens := []gopter.Gen{}
	fieldTypes := []reflect.Type{}
	assignable := reflect.New(rt).Elem()
	for i := 0; i < rt.NumField(); i++ {
		fieldName := rt.Field(i).Name
		if !assignable.Field(i).CanSet() {
			continue
		}

		gen := gens[fieldName]
		if gen != nil {
			fieldGens = append(fieldGens, gen)
			fieldTypes = append(fieldTypes, rt.Field(i).Type)
		}
	}

	buildStructType := reflect.FuncOf(fieldTypes, []reflect.Type{rt}, false)
	unbuildStructType := reflect.FuncOf([]reflect.Type{rt}, fieldTypes, false)

	buildStructFunc := reflect.MakeFunc(buildStructType, func(args []reflect.Value) []reflect.Value {
		result := reflect.New(rt)
		for i := 0; i < rt.NumField(); i++ {
			if _, ok := gens[rt.Field(i).Name]; !ok {
				continue
			}
			if !assignable.Field(i).CanSet() {
				continue
			}
			result.Elem().Field(i).Set(args[0])
			args = args[1:]
		}
		return []reflect.Value{result.Elem()}
	})
	unbuildStructFunc := reflect.MakeFunc(unbuildStructType, func(args []reflect.Value) []reflect.Value {
		s := args[0]
		results := []reflect.Value{}
		for i := 0; i < s.NumField(); i++ {
			if _, ok := gens[rt.Field(i).Name]; !ok {
				continue
			}
			if !assignable.Field(i).CanSet() {
				continue
			}
			results = append(results, s.Field(i))
		}
		return results
	})

	return gopter.DeriveGen(
		buildStructFunc.Interface(),
		unbuildStructFunc.Interface(),
		fieldGens...,
	)
}

// StructPtr generates pointers to a given struct type.
// Note that StructPtr does not generate nil, if you want to include nil in your
// testing you should combine gen.PtrOf with gen.Struct.
// rt has to be the reflect type of the struct, gens contains a map of field generators.
// Note that the result types of the generators in gen have to match the type of the correspoinding
// field in the struct. Also note that only public fields of a struct can be generated
func StructPtr(rt reflect.Type, gens map[string]gopter.Gen) gopter.Gen {
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}

	buildPtrType := reflect.FuncOf([]reflect.Type{rt}, []reflect.Type{reflect.PtrTo(rt)}, false)
	unbuildPtrType := reflect.FuncOf([]reflect.Type{reflect.PtrTo(rt)}, []reflect.Type{rt}, false)

	buildPtrFunc := reflect.MakeFunc(buildPtrType, func(args []reflect.Value) []reflect.Value {
		sp := reflect.New(rt)
		sp.Elem().Set(args[0])
		return []reflect.Value{sp}
	})
	unbuildPtrFunc := reflect.MakeFunc(unbuildPtrType, func(args []reflect.Value) []reflect.Value {
		return []reflect.Value{args[0].Elem()}
	})

	return gopter.DeriveGen(
		buildPtrFunc.Interface(),
		unbuildPtrFunc.Interface(),
		Struct(rt, gens),
	)
}

// checkFieldsMatch panics unless the keys in gens exactly match the public
// fields on rt. With an extra bool argument of value "true", it only panics if
// there's a key in gens which is not a field on rt.
func checkFieldsMatch(
	rt reflect.Type,
	gens map[string]gopter.Gen,
	allowFieldsWithNoGenerator ...bool,
) {
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}
	fields := make(map[string]bool, rt.NumField())
	for i := 0; i < rt.NumField(); i++ {
		fields[rt.Field(i).Name] = true
	}
	for field := range gens {
		if _, ok := fields[field]; !ok {
			panic(fmt.Errorf("generator for non-existent field %s on struct %s",
				field, rt.Name()))
		}
		delete(fields, field)
	}
	if len(allowFieldsWithNoGenerator) > 0 && allowFieldsWithNoGenerator[0] {
		return // Don't check that every field is present in gens
	}
	if len(allowFieldsWithNoGenerator) > 1 {
		panic("expect at most one boolean argument in StrictStruct/StrictStructPtr")
	}
	if len(fields) != 0 { // Check that every field is present in gens
		var missingFields []string
		for field := range fields {
			missingFields = append(missingFields, field)
		}
		panic(fmt.Errorf("generator missing for fields %v on struct %s",
			missingFields, rt.Name()))
	}
}

// StrictStruct behaves the same as Struct, except it requires the keys in gens
// to exactly match the public fields of rt. It panics if gens contains extra
// keys, or has missing keys.
//
// If given a third true argument, it only requires the keys of gens to be
// fields of rt. In that case, unspecified fields will remain unset.
func StrictStruct(
	rt reflect.Type,
	gens map[string]gopter.Gen,
	allowFieldsWithNoGenerator ...bool,
) gopter.Gen {
	checkFieldsMatch(rt, gens, allowFieldsWithNoGenerator...)
	return Struct(rt, gens)
}

// StrictStructPtr behaves the same as StructPtr, except it requires the keys in
// gens to exactly match the public fields of rt. It panics if gens contains
// extra keys, or has missing keys.
//
// If given a third true argument, it only requires the keys of gens to be
// fields of rt. In that case, unspecified fields will remain unset.
func StrictStructPtr(
	rt reflect.Type,
	gens map[string]gopter.Gen,
	allowFieldsWithNoGenerator ...bool,
) gopter.Gen {
	checkFieldsMatch(rt, gens, allowFieldsWithNoGenerator...)
	return StructPtr(rt, gens)
}

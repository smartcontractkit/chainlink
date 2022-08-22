package v2

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/core/utils"
)

// Validated configurations impose constraints that must be checked.
type Validated interface {
	// ValidateConfig returns nil if the config is valid, otherwise an error describing why it is invalid.
	//
	// For implementations:
	//  - A nil receiver should return nil, freeing the caller to decide whether each case is required.
	//  - Use package multierr to accumulate all errors, rather than returning the first encountered.
	ValidateConfig() error
}

// Validate returns any errors from calling Validated.ValidateConfig on cfg and any nested types that implement Validated.
func Validate(cfg interface{}) (err error) {
	return utils.MultiErrorList(validate(cfg))
}

func validate(s interface{}) (err error) {
	if vc, ok := s.(Validated); ok {
		err = multierr.Append(err, vc.ValidateConfig())
	}

	t := reflect.TypeOf(s)
	v := reflect.ValueOf(s)
	if t.Kind() == reflect.Ptr {
		if v.IsNil() {
			//TODO error if required? https://app.shortcut.com/chainlinklabs/story/33618/add-config-validate-command
			return
		}
		t = t.Elem()
		v = v.Elem()
	}
	switch t.Kind() {
	case reflect.Bool, reflect.Chan, reflect.Complex128, reflect.Complex64, reflect.Float32, reflect.Float64,
		reflect.Func, reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8, reflect.Interface,
		reflect.Invalid, reflect.Ptr, reflect.String, reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Uint8, reflect.Uintptr, reflect.UnsafePointer:
		//TODO additional field validation? e.g. struct tags? https://app.shortcut.com/chainlinklabs/story/33618/add-config-validate-command
		return
	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			ft := t.Field(i)
			if !ft.IsExported() {
				continue
			}
			fv := v.Field(i)
			if !fv.CanInterface() {
				continue
			}
			if fe := Validate(fv.Interface()); fe != nil {
				err = multierr.Append(err, namedMultiErrorList(fe, ft.Name))
			}
		}
		return
	case reflect.Map:
		iter := v.MapRange()
		for iter.Next() {
			mk := iter.Key()
			mv := iter.Value()
			if !v.CanInterface() {
				continue
			}
			if me := Validate(mv.Interface()); me != nil {
				err = multierr.Append(err, namedMultiErrorList(me, fmt.Sprintf("%s", mk.Interface())))
			}
		}
		return
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			iv := v.Index(i)
			if !v.CanInterface() {
				continue
			}
			if me := Validate(iv.Interface()); me != nil {
				err = multierr.Append(err, namedMultiErrorList(me, strconv.Itoa(i)))
			}
		}
		return
	}

	return fmt.Errorf("should be unreachable: switch missing case for kind: %s", t.Kind())
}

func namedMultiErrorList(err error, name string) error {
	err = utils.MultiErrorList(err)
	msg := strings.ReplaceAll(err.Error(), "\n", "\n\t")
	return fmt.Errorf("%s: %s", name, msg)
}

type ErrInvalid struct {
	Name  string
	Value any
	Msg   string
}

func (e ErrInvalid) Error() string {
	return fmt.Sprintf("%s: invalid value %v: %s", e.Name, e.Value, e.Msg)
}

type ErrMissing struct {
	Name string
	Msg  string
}

func (e ErrMissing) Error() string {
	return fmt.Sprintf("%s: missing: %s", e.Name, e.Msg)
}

type ErrEmpty struct {
	Name string
	Msg  string
}

func (e ErrEmpty) Error() string {
	return fmt.Sprintf("%s: empty: %s", e.Name, e.Msg)
}

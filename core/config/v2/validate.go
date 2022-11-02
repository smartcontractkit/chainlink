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
	//  - Use package multierr to accumulate all errors, rather than returning the first encountered.
	//  - If an anonymous field also implements ValidateConfig(), it must be called explicitly!
	ValidateConfig() error
}

// Validate returns any errors from calling Validated.ValidateConfig on cfg and any nested types that implement Validated.
func Validate(cfg interface{}) (err error) {
	_, err = utils.MultiErrorList(validate(reflect.ValueOf(cfg), true))
	return
}

func validate(v reflect.Value, checkInterface bool) (err error) {
	if checkInterface {
		i := v.Interface()
		if vc, ok := i.(Validated); ok {
			err = multierr.Append(err, vc.ValidateConfig())
		} else if v.CanAddr() {
			i = v.Addr().Interface()
			if vc, ok := i.(Validated); ok {
				err = multierr.Append(err, vc.ValidateConfig())
			}
		}
	}

	t := v.Type()
	if t.Kind() == reflect.Ptr {
		if v.IsNil() {
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
			if fv.Kind() == reflect.Ptr && fv.IsNil() {
				continue
			}
			// skip the interface if Anonymous, since the parent struct inherits the methods
			if fe := validate(fv, !ft.Anonymous); fe != nil {
				if ft.Anonymous {
					err = multierr.Append(err, fe)
				} else {
					err = multierr.Append(err, namedMultiErrorList(fe, ft.Name))
				}
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
			if mv.Kind() == reflect.Ptr && mv.IsNil() {
				continue
			}
			if me := validate(mv, true); me != nil {
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
			if iv.Kind() == reflect.Ptr && iv.IsNil() {
				continue
			}
			if me := validate(iv, true); me != nil {
				err = multierr.Append(err, namedMultiErrorList(me, strconv.Itoa(i)))
			}
		}
		return
	}

	return fmt.Errorf("should be unreachable: switch missing case for kind: %s", t.Kind())
}

func namedMultiErrorList(err error, name string) error {
	l, merr := utils.MultiErrorList(err)
	if l == 0 {
		return nil
	}
	msg := strings.ReplaceAll(merr.Error(), "\n", "\n\t")
	if l == 1 {
		return fmt.Errorf("%s.%s", name, msg)
	}
	return fmt.Errorf("%s: %s", name, msg)
}

type ErrInvalid struct {
	Name  string
	Value any
	Msg   string
}

// NewErrDuplicate returns an ErrInvalid with a standard duplicate message.
func NewErrDuplicate(name string, value any) ErrInvalid {
	return ErrInvalid{Name: name, Value: value, Msg: "duplicate - must be unique"}
}

func (e ErrInvalid) Error() string {
	return fmt.Sprintf("%s: invalid value (%v): %s", e.Name, e.Value, e.Msg)
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

// UniqueStrings is a helper for tracking unique values in string form.
type UniqueStrings map[string]struct{}

// IsDupeFmt is like IsDupe, but calls String().
func (u UniqueStrings) IsDupeFmt(t fmt.Stringer) bool {
	if t == nil {
		return false
	}
	if reflect.ValueOf(t).IsNil() {
		// interface holds a typed-nil value
		return false
	}
	return u.isDupe(t.String())
}

// IsDupe returns true if the set already contains the string, otherwise false.
// Non-nil/empty strings are added to the set.
func (u UniqueStrings) IsDupe(s *string) bool {
	if s == nil {
		return false
	}
	return u.isDupe(*s)
}

func (u UniqueStrings) isDupe(s string) bool {
	if s == "" {
		return false
	}
	_, ok := u[s]
	if !ok {
		u[s] = struct{}{}
	}
	return ok
}

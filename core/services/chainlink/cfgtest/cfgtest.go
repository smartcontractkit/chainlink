package cfgtest

import (
	"encoding"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/multierr"
)

func AssertFieldsNotNil(t *testing.T, s interface{}) {
	err := assertValNotNil(t, "", reflect.ValueOf(s))
	assert.NoError(t, err, MultiErrorList(err))
}

// assertFieldsNotNil recursively checks the struct s for nil fields.
func assertFieldsNotNil(t *testing.T, prefix string, s reflect.Value) (err error) {
	t.Helper()
	require.Equal(t, reflect.Struct, s.Kind())

	typ := s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		key := prefix
		if tf := typ.Field(i); !tf.Anonymous {
			if key != "" {
				key += "."
			}
			key += tf.Name
		}
		err = multierr.Combine(err, assertValNotNil(t, key, f))
	}
	return
}

// assertValuesNotNil recursively checks the map m for nil values.
func assertValuesNotNil(t *testing.T, prefix string, m reflect.Value) (err error) {
	t.Helper()
	require.Equal(t, reflect.Map, m.Kind())
	if prefix != "" {
		prefix += "."
	}

	mi := m.MapRange()
	for mi.Next() {
		key := prefix + mi.Key().String()
		err = multierr.Combine(err, assertValNotNil(t, key, mi.Value()))
	}
	return
}

// assertElementsNotNil recursively checks the slice s for nil values.
func assertElementsNotNil(t *testing.T, prefix string, s reflect.Value) (err error) {
	t.Helper()
	require.Equal(t, reflect.Slice, s.Kind())

	for i := 0; i < s.Len(); i++ {
		err = multierr.Combine(err, assertValNotNil(t, prefix, s.Index(i)))
	}
	return
}

var (
	textUnmarshaler     encoding.TextUnmarshaler
	textUnmarshalerType = reflect.TypeOf(&textUnmarshaler).Elem()
)

// assertValNotNil recursively checks that val is not nil. val must be a struct, map, slice, or point to one.
func assertValNotNil(t *testing.T, key string, val reflect.Value) error {
	t.Helper()
	k := val.Kind()
	switch k { //nolint:exhaustive
	case reflect.Ptr, reflect.Map, reflect.Slice:
		if val.IsNil() {
			return fmt.Errorf("%s: nil", key)
		}
	}
	if k == reflect.Ptr {
		if val.Type().Implements(textUnmarshalerType) {
			return nil // skip values unmarshaled from strings
		}
		val = val.Elem()
	}
	switch val.Kind() {
	case reflect.Struct:
		if val.Type().Implements(textUnmarshalerType) {
			return nil // skip values unmarshaled from strings
		}
		return assertFieldsNotNil(t, key, val)
	case reflect.Map:
		return assertValuesNotNil(t, key, val)
	case reflect.Slice:
		return assertElementsNotNil(t, key, val)
	default:
		return nil
	}
}

type multiErrorList []error

// MultiErrorList returns an error which formats underlying errors as a list, or nil if err is nil.
func MultiErrorList(err error) error {
	if err == nil {
		return nil
	}

	return multiErrorList(multierr.Errors(err))
}

func (m multiErrorList) Error() string {
	l := len(m)
	if l == 1 {
		return m[0].Error()
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d errors:", l)
	for _, e := range m {
		fmt.Fprintf(&sb, "\n\t- %v", e)
	}
	return sb.String()
}

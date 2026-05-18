package cfgtest

import (
	"encoding"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink-common/pkg/config"
)

func AssertFieldsNotNil(t *testing.T, s interface{}) {
	err := assertValNotNil(t, "", reflect.ValueOf(s))
	_, err = config.MultiErrorList(err)
	assert.NoError(t, err)
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
	case reflect.Ptr:
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
		if val.IsNil() {
			return nil // not actually a problem
		}
		return assertValuesNotNil(t, key, val)
	case reflect.Slice:
		if val.IsNil() {
			return nil // not actually a problem
		}
		return assertElementsNotNil(t, key, val)
	default:
		return nil
	}
}

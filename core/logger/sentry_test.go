package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_toMap(t *testing.T) {
	t.Run("with even number of keys/values", func(t *testing.T) {
		keysAndValues := []interface{}{
			"foo", 1, "bar", 42.43, "boggly", "str",
		}

		m := toMap(keysAndValues)

		assert.Equal(t, map[string]interface{}{"bar": 42.43, "boggly": "str", "foo": 1}, m)
	})

	t.Run("with odd number of keys/values, drops the last one", func(t *testing.T) {
		keysAndValues := []interface{}{
			"foo", 1, "bar", 42.43, "boggly", "str", "odd",
		}

		m := toMap(keysAndValues)

		assert.Equal(t, map[string]interface{}{"bar": 42.43, "boggly": "str", "foo": 1}, m)
	})
}

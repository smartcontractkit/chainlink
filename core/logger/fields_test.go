package logger_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/stretchr/testify/assert"
)

func TestFields_Merge(t *testing.T) {
	t.Parallel()

	f1 := make(logger.Fields)
	f1["key1"] = "value1"
	f2 := make(logger.Fields)
	f2["key2"] = "value2"

	merged := f1.Merge(f2)
	assert.Len(t, merged, 2)

	v1, ok1 := merged["key1"]
	assert.True(t, ok1)
	assert.Equal(t, "value1", v1)

	v2, ok2 := merged["key2"]
	assert.True(t, ok2)
	assert.Equal(t, "value2", v2)

	t.Run("self merge", func(t *testing.T) {
		t.Parallel()

		merged := f1.Merge(f1)
		assert.Len(t, merged, 1)
		assert.Equal(t, f1, merged)
	})
}

func TestFields_Slice(t *testing.T) {
	t.Parallel()

	f := make(logger.Fields)
	f["str"] = "foo"
	f["int"] = 123

	s := f.Slice()
	assert.Len(t, s, 4)
	for i := 0; i < len(s); i += 2 {
		switch s[i] {
		case "int":
			assert.Equal(t, 123, s[i+1])
		case "str":
			assert.Equal(t, "foo", s[i+1])
		}
	}

	t.Run("empty", func(t *testing.T) {
		t.Parallel()

		empty := make(logger.Fields)
		assert.Empty(t, empty.Slice())
	})
}

func TestFields_With(t *testing.T) {
	t.Parallel()

	f := make(logger.Fields)
	f["str"] = "foo"
	f["int"] = 123

	w := f.With("bool", true, "float", 3.14)
	assert.Len(t, w, 4)

	t.Run("single", func(t *testing.T) {
		t.Parallel()

		assert.Panics(t, func() {
			//lint:ignore SA5012 we expect panic here
			_ = f.With("xyz")
		}, "expected even number of arguments")
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()

		empty := make(logger.Fields).With()
		assert.Empty(t, empty)
	})
}

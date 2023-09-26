package loop

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_removeArg(t *testing.T) {
	for _, tt := range []struct {
		name string
		args []interface{}
		key  string

		wantArgs []interface{}
		wantVal  string
	}{
		{"empty", nil, "logger",
			nil, ""},
		{"simple", []any{"logger", "foo"}, "logger",
			[]any{}, "foo"},
		{"multi", []any{"logger", "foo", "bar", "baz"}, "logger",
			[]any{"bar", "baz"}, "foo"},
		{"reorder", []any{"bar", "baz", "logger", "foo"}, "logger",
			[]any{"bar", "baz"}, "foo"},

		{"invalid", []any{"logger"}, "logger",
			[]any{"logger"}, ""},
		{"invalid-multi", []any{"foo", "bar", "logger"}, "logger",
			[]any{"foo", "bar", "logger"}, ""},
		{"value", []any{"foo", "logger", "bar", "baz"}, "logger",
			[]any{"foo", "logger", "bar", "baz"}, ""},
	} {
		t.Run(tt.name, func(t *testing.T) {
			args, val := removeArg(tt.args, tt.key)
			assert.ElementsMatch(t, tt.wantArgs, args)
			assert.Equal(t, tt.wantVal, val)
		})
	}
}

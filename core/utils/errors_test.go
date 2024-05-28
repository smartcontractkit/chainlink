package utils

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/multierr"
)

func TestFlatten(t *testing.T) {
	e := []error{
		errors.New("0"),
		errors.New("1"),
		errors.New("2"),
		errors.New("3"),
	}

	// nested errors
	// [[[0, 1], 2], 3]
	err0 := errors.Join(nil, e[0])
	err0 = errors.Join(err0, e[1])
	err0 = errors.Join(err0, e[2])
	err0 = errors.Join(err0, e[3])

	// flat error
	err1 := errors.Join(e...)

	// multierr provides a flat error
	err2 := multierr.Append(nil, e[0])
	err2 = multierr.Append(err2, e[1])
	err2 = multierr.Append(err2, e[2])
	err2 = multierr.Append(err2, e[3])

	params := []struct {
		name string
		err  error
		out  []error
	}{
		{"errors.Join nested", err0, e},
		{"errors.Join flat", err1, e},
		{"multierr.Append", err2, e},
		{"nil", nil, []error{nil}},
		{"single", e[0], []error{e[0]}},
	}

	for _, p := range params {
		t.Run(p.name, func(t *testing.T) {
			assert.Equal(t, p.out, Flatten(p.err))
		})
	}
}

package oraclecreator

import (
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_wrappedOracle_Close(t *testing.T) {
	tests := []struct {
		name         string
		oracleErr    error
		closerErrors []error
		expectedErr  error
	}{
		{
			name:        "no errors",
			expectedErr: nil,
		},
		{
			name:        "oracle error",
			oracleErr:   err1,
			expectedErr: errors.New("close base oracle: err1"),
		},
		{
			name:         "oracle and closers errors",
			oracleErr:    err1,
			closerErrors: []error{nil, nil, err3},
			expectedErr:  errors.New("close base oracle: err1\nerr3"),
		},
		{
			name:         "closers only errors",
			oracleErr:    nil,
			closerErrors: []error{nil, err2, nil},
			expectedErr:  err2,
		},
		{
			name:         "no errors with closers",
			closerErrors: []error{nil, nil, nil, nil},
			expectedErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			closers := make([]io.Closer, 0, len(tt.closerErrors))
			for _, err := range tt.closerErrors {
				closers = append(closers, mockCloser{err: err})
			}

			o := newWrappedOracle(mockOracle{err: tt.oracleErr}, closers)

			err := o.Close()
			if err == nil && tt.expectedErr == nil {
				assert.NoError(t, err)
				return
			}

			assert.Error(t, err)
			assert.Equal(t, tt.expectedErr.Error(), err.Error())
		})
	}
}

type mockCloser struct{ err error }

func (m mockCloser) Close() error { return m.err }

type mockOracle struct{ err error }

func (m mockOracle) Close() error { return m.err }

func (m mockOracle) Start() error { return m.err }

var (
	err1 = errors.New("err1")
	err2 = errors.New("err2")
	err3 = errors.New("err3")
)

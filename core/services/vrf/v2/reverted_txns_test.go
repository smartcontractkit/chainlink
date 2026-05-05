package v2

import (
	"context"
	stderrors "errors"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestSkipRevertedTxnFetchFatal(t *testing.T) {
	t.Parallel()

	ctxAlive := context.Background()
	ctxCanceled, cancel := context.WithCancel(context.Background())
	cancel()

	dbErr := stderrors.New("connection reset")

	tests := []struct {
		name  string
		ctxFn func() context.Context
		err   error
		want  bool
	}{
		{
			name: "nil error",
			ctxFn: func() context.Context {
				return ctxAlive
			},
			err:  nil,
			want: false,
		},
		{
			name: "context.Canceled on alive ctx",
			ctxFn: func() context.Context {
				return ctxAlive
			},
			err:  context.Canceled,
			want: true,
		},
		{
			name: "wrapped context.Canceled",
			ctxFn: func() context.Context {
				return ctxAlive
			},
			err:  errors.Wrap(context.Canceled, "pq"),
			want: true,
		},
		{
			name: "context.DeadlineExceeded on alive ctx inner query timeout",
			ctxFn: func() context.Context {
				return ctxAlive
			},
			err:  context.DeadlineExceeded,
			want: false,
		},
		{
			name: "wrapped DeadlineExceeded on alive ctx",
			ctxFn: func() context.Context {
				return ctxAlive
			},
			err:  errors.Wrap(context.DeadlineExceeded, "timeout"),
			want: false,
		},
		{
			name: "DeadlineExceeded while outer ctx canceled",
			ctxFn: func() context.Context {
				return ctxCanceled
			},
			err:  context.DeadlineExceeded,
			want: true,
		},
		{
			name: "generic DB error on alive ctx",
			ctxFn: func() context.Context {
				return ctxAlive
			},
			err:  dbErr,
			want: false,
		},
		{
			name: "generic DB error while outer ctx canceled",
			ctxFn: func() context.Context {
				return ctxCanceled
			},
			err:  dbErr,
			want: true,
		},
		{
			name: "context.Canceled while outer ctx also canceled",
			ctxFn: func() context.Context {
				return ctxCanceled
			},
			err:  context.Canceled,
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := skipRevertedTxnFetchFatal(tt.ctxFn(), tt.err)
			assert.Equal(t, tt.want, got)
		})
	}
}

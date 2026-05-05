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
		name string
		ctx  context.Context
		err  error
		want bool
	}{
		{
			name: "nil error",
			ctx:  ctxAlive,
			err:  nil,
			want: false,
		},
		{
			name: "context.Canceled on alive ctx",
			ctx:  ctxAlive,
			err:  context.Canceled,
			want: true,
		},
		{
			name: "wrapped context.Canceled",
			ctx:  ctxAlive,
			err:  errors.Wrap(context.Canceled, "pq"),
			want: true,
		},
		{
			name: "context.DeadlineExceeded on alive ctx inner query timeout",
			ctx:  ctxAlive,
			err:  context.DeadlineExceeded,
			want: false,
		},
		{
			name: "wrapped DeadlineExceeded on alive ctx",
			ctx:  ctxAlive,
			err:  errors.Wrap(context.DeadlineExceeded, "timeout"),
			want: false,
		},
		{
			name: "DeadlineExceeded while outer ctx canceled",
			ctx:  ctxCanceled,
			err:  context.DeadlineExceeded,
			want: true,
		},
		{
			name: "generic DB error on alive ctx",
			ctx:  ctxAlive,
			err:  dbErr,
			want: false,
		},
		{
			name: "generic DB error while outer ctx canceled",
			ctx:  ctxCanceled,
			err:  dbErr,
			want: true,
		},
		{
			name: "context.Canceled while outer ctx also canceled",
			ctx:  ctxCanceled,
			err:  context.Canceled,
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := skipRevertedTxnFetchFatal(tt.ctx, tt.err)
			assert.Equal(t, tt.want, got)
		})
	}
}

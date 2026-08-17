package grpcsourcemock

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestNewJWTAuthInterceptor(t *testing.T) {
	t.Parallel()
	provider := &AcceptAllAuthProvider{}
	interceptor := NewJWTAuthInterceptor(provider)
	require.NotNil(t, interceptor)

	// Call interceptor without bearer token in context; expect unauthenticated error
	ctx := context.Background()
	handler := func(ctx context.Context, req any) (any, error) {
		return "ok", nil
	}

	res, err := interceptor(ctx, "req", &grpc.UnaryServerInfo{}, handler)
	assert.Nil(t, res)
	require.Error(t, err)
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.Unauthenticated, st.Code())
}

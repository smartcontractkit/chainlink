// Copyright 2021 Michal Witkowski. All Rights Reserved.
// See LICENSE for licensing terms.

package proxy

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// NewProxy sets up a simple proxy that forwards all requests to dst.
func NewProxy(dst grpc.ClientConnInterface, opts ...grpc.ServerOption) *grpc.Server {
	opts = append(opts, DefaultProxyOpt(dst))
	// Set up the proxy server and then serve from it like in step one.
	return grpc.NewServer(opts...)
}

// DefaultProxyOpt returns an grpc.UnknownServiceHandler with a DefaultDirector.
func DefaultProxyOpt(cc grpc.ClientConnInterface) grpc.ServerOption {
	return grpc.UnknownServiceHandler(TransparentHandler(DefaultDirector(cc)))
}

// DefaultDirector returns a very simple forwarding StreamDirector that forwards all
// calls.
func DefaultDirector(cc grpc.ClientConnInterface) StreamDirector {
	return func(ctx context.Context, fullMethodName string) (context.Context, grpc.ClientConnInterface, error) {
		md, _ := metadata.FromIncomingContext(ctx)
		ctx = metadata.NewOutgoingContext(ctx, md.Copy())
		return ctx, cc, nil
	}
}

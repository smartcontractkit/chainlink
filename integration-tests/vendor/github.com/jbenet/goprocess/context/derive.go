package goprocessctx

import (
	"context"

	goprocess "github.com/jbenet/goprocess"
)

// OnClosingContext derives a context from a given goprocess that will
// be 'Done' when the process is closing
func OnClosingContext(p goprocess.Process) context.Context {
	return WithProcessClosing(context.Background(), p)
}

// OnClosedContext derives a context from a given goprocess that will
// be 'Done' when the process is closed
func OnClosedContext(p goprocess.Process) context.Context {
	return WithProcessClosed(context.Background(), p)
}

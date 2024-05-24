package services

import (
	"io"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
)

// StartClose is a subset of the ServiceCtx interface.
type StartClose = services.StartClose

// MultiStart is a utility for starting multiple services together.
// The set of started services is tracked internally, so that they can be closed if any single service fails to start.
type MultiStart = services.MultiStart

// CloseAll closes all elements concurrently.
// Use this when you have various different types of io.Closer.
func CloseAll(cs ...io.Closer) error {
	return services.CloseAll(cs...)
}

// MultiCloser returns an io.Closer which closes all elements concurrently.
// Use this when you have a slice of a type which implements io.Closer.
// []io.Closer can be cast directly to MultiCloser.
func MultiCloser[C io.Closer](cs []C) io.Closer {
	return services.MultiCloser(cs)
}

package operations

import (
	"github.com/Masterminds/semver/v3"
)

// SequenceHandler is the function signature of a sequence handler.
// A sequence handler is a function that calls 1 or more operations or child sequence.
type SequenceHandler[IN, OUT, DEP any] func(b Bundle, deps DEP, input IN) (output OUT, err error)

type Sequence[IN, OUT, DEP any] struct {
	def     Definition
	handler SequenceHandler[IN, OUT, DEP]
}

// NewSequence creates a new sequence.
// Useful for logically grouping a list of operations and also child sequence together.
func NewSequence[IN, OUT, DEP any](
	id string, version *semver.Version, description string, handler SequenceHandler[IN, OUT, DEP],
) *Sequence[IN, OUT, DEP] {
	return &Sequence[IN, OUT, DEP]{
		def: Definition{
			ID:          id,
			Version:     version,
			Description: description,
		},
		handler: handler,
	}
}

// ID returns the sequence ID.
func (o *Sequence[IN, OUT, DEP]) ID() string {
	return o.def.ID
}

// Version returns the sequence semver version in string.
func (o *Sequence[IN, OUT, DEP]) Version() string {
	return o.def.Version.String()
}

// Description returns the sequence description.
func (o *Sequence[IN, OUT, DEP]) Description() string {
	return o.def.Description
}

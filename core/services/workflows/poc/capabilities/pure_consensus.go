package capabilities

func NewPureConsensus[I, O any](ref string, fn func([]I) (O, error)) Consensus[I, O] {
	return &pureConsensus[I, O]{ref: ref, fn: fn}
}

func NewIdenticalConsensus[I any]() Consensus[I, I] {
	return NewPureConsensus[I, I]("identical", func(observations []I) (I, error) {
		// check for the one with the most results, verify it's at least F+1
		// We might want to expose that from the core node to make life faster/easier
		return observations[0], nil
	})
}

type pureConsensus[I, O any] struct {
	ref string
	fn  func([]I) (O, error)
}

func (p pureConsensus[I, O]) Type() string {
	return LocalCodeConsensusCapability
}

func (p pureConsensus[I, O]) Ref() string {
	return p.ref
}

func (p pureConsensus[I, O]) Invoke(observations []I) (ConsensusResult[O], error) {
	results, err := p.fn(observations)
	return &consensusResult[O]{results: results}, err
}

package capabilities

func NewPureConsensus[I, O any](ref string, fn func([]I) (O, error)) Consensus[I, O] {
	return &pureConsensus[I, O]{ref: ref, fn: fn}
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

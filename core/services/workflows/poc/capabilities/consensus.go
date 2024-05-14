package capabilities

type Consensus[I, O any] interface {
	Base
	Invoke(observations []I) (ConsensusResult[O], error)
}

type ConsensusResult[O any] interface {
	Results() O
	private()
}

type consensusResult[O any] struct {
	results O
}

func (c *consensusResult[O]) Results() O {
	return c.results
}

func (c *consensusResult[O]) private() {}

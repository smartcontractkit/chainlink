package capabilities

import "errors"

// BuiltInConsensus should be created by generated code, or referenced from a capabilities' public library to assure correct usage with types
type BuiltInConsensus[I, O any] struct {
	RefName  string
	TypeName string
}

func (b *BuiltInConsensus[I, O]) Type() string {
	return b.TypeName
}

func (b *BuiltInConsensus[I, O]) Ref() string {
	return b.RefName
}

func (b *BuiltInConsensus[I, O]) Invoke(_ []I) (ConsensusResult[O], error) {
	return nil, errors.New("host should be invoking built-in consensus, this should be a placeholder")
}

var _ Consensus[any, any] = &BuiltInConsensus[any, any]{}

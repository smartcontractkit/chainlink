package capabilities

/*

Something similar to the below to allow trigger normalization
we don't have it in the engine right now, but it can/should be built out when needed.
It would work similarly to action and consensus

type PureAdaptedTrigger[I, O any] struct {
	Trigger[I]
	Adaptor func(I) (O, error)
}

func (p PureAdaptedTrigger[I, O]) Transform(a values.Value) (O, error) {
	input, err := p.Trigger.Transform(a)
	if err != nil {
		var o O
		return o, err
	}
	return p.Adaptor(input)
}

var _ Trigger[any] = (*PureAdaptedTrigger[any, any])(nil)
*/

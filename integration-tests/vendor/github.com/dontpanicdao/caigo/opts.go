package caigo

type curveOptions struct {
	initConstants bool
	paramsPath    string
}

// funcCurveOptions wraps a function that modifies curveOptions into an
// implementation of the CurveOption interface.
type funcCurveOption struct {
	f func(*curveOptions)
}

func (fso *funcCurveOption) apply(do *curveOptions) {
	fso.f(do)
}

func newFuncCurveOption(f func(*curveOptions)) *funcCurveOption {
	return &funcCurveOption{
		f: f,
	}
}

type CurveOption interface {
	apply(*curveOptions)
}

// functions that require pedersen hashes must be run on
// a curve initialized with constant points
func WithConstants(paramsPath ...string) CurveOption {
	return newFuncCurveOption(func(o *curveOptions) {
		o.initConstants = true

		if len(paramsPath) == 1 && paramsPath[0] != "" {
			o.paramsPath = paramsPath[0]
		}
	})
}

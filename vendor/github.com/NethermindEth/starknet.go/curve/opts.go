package curve

type curveOptions struct {
	initConstants bool
	paramsPath    string
}

// funcCurveOptions wraps a function that modifies curveOptions into an
// implementation of the CurveOption interface.
type funcCurveOption struct {
	f func(*curveOptions)
}

// apply applies the given curve options to the funcCurveOption.
//
// Parameters:
// - fso: a pointer to funcCurveOption
// Returns:
//  none
func (fso *funcCurveOption) apply(do *curveOptions) {
	fso.f(do)
}

// newFuncCurveOption returns a new instance of funcCurveOption.
//
// Parameters:
// - f: a function of type func(*curveOptions)
// Returns:
// - a pointer to funcCurveOption
func newFuncCurveOption(f func(*curveOptions)) *funcCurveOption {
	return &funcCurveOption{
		f: f,
	}
}

type CurveOption interface {
	apply(*curveOptions)
}

// WithConstants creates a CurveOption (a curve initialized with constant points) that initializes the constants of the curve.
//
// Parameters:
// - paramsPath: a variadic parameter of type string, representing the path(s) to the parameters
// Returns:
// - a new instance of CurveOption
func WithConstants(paramsPath ...string) CurveOption {
	return newFuncCurveOption(func(o *curveOptions) {
		o.initConstants = true

		if len(paramsPath) == 1 && paramsPath[0] != "" {
			o.paramsPath = paramsPath[0]
		}
	})
}

package capabilities

type Action[I, O any] interface {
	Base
	Invoke(input I) (O, bool, error)
}

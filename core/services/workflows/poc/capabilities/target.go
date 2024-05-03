package capabilities

type Target[I any] interface {
	Base
	Invoke(input I) error
}

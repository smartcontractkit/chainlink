package workflow

type Runner interface {
	Run(*Spec) error
}

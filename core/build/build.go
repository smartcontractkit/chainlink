// Package build utilizes build tags and package testing API to determine the environment that this binary was built to target.
//   - Prod is the default
//   - Test is automatically set in test binaries, e.g. when using `go test`
//   - Dev can be set with the 'dev' build tag, for standard builds or test binaries
package build

const (
	Prod = "prod"
	Dev  = "dev"
	Test = "test"
)

var mode string

func Mode() string { return mode }

func IsDev() bool {
	return mode == Dev
}

func IsTest() bool {
	return mode == Test
}

func IsProd() bool {
	return mode == Prod
}

package build

// The build module utilizes build tags to determine the environment that this binary was built to target
// the currently supported build modes are dev, test. Setting both tags is not allowed and will result to compilation errors.

const (
	Prod = "prod"
	Dev  = "dev"
	Test = "test"
)

func IsDev() bool {
	return mode == Dev
}

func IsTest() bool {
	return mode == Test
}

func IsProd() bool {
	return mode == Prod
}

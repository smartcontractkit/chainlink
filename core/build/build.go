package build

// The build module utilizes build tags to determine the environment that this binary was built to target
// the currently supported build modes are dev, test. Setting both tags is not allowed and will result to compilation errors.

const (
	Prod = "prod"
	Dev  = "dev"
	Test = "test"
)

func DevelopmentBuild() bool {
	return Mode == Dev
}

func TestBuild() bool {
	return Mode == Test
}

func ProdBuild() bool {
	return Mode == Prod
}

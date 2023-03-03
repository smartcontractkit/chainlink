//go:build no_runtime_type_checking

// A programming model for software-defined state
package constructs

// Building without runtime type checking enabled, so all the below just return nil

func validateConstruct_IsConstructParameters(x interface{}) error {
	return nil
}

func validateNewConstructParameters(scope Construct, id *string) error {
	return nil
}


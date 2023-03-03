//go:build no_runtime_type_checking

// This is the core library of Cloud Development Kit (CDK) for Kubernetes (cdk8s). cdk8s apps synthesize into standard Kubernetes manifests which can be applied to any Kubernetes cluster.
package cdk8s

// Building without runtime type checking enabled, so all the below just return nil

func validateJsonPatch_AddParameters(path *string, value interface{}) error {
	return nil
}

func validateJsonPatch_ApplyParameters(document interface{}) error {
	return nil
}

func validateJsonPatch_CopyParameters(from *string, path *string) error {
	return nil
}

func validateJsonPatch_MoveParameters(from *string, path *string) error {
	return nil
}

func validateJsonPatch_RemoveParameters(path *string) error {
	return nil
}

func validateJsonPatch_ReplaceParameters(path *string, value interface{}) error {
	return nil
}

func validateJsonPatch_TestParameters(path *string, value interface{}) error {
	return nil
}


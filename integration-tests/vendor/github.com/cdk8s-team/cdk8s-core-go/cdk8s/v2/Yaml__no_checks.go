//go:build no_runtime_type_checking

// This is the core library of Cloud Development Kit (CDK) for Kubernetes (cdk8s). cdk8s apps synthesize into standard Kubernetes manifests which can be applied to any Kubernetes cluster.
package cdk8s

// Building without runtime type checking enabled, so all the below just return nil

func validateYaml_FormatObjectsParameters(docs *[]interface{}) error {
	return nil
}

func validateYaml_LoadParameters(urlOrFile *string) error {
	return nil
}

func validateYaml_SaveParameters(filePath *string, docs *[]interface{}) error {
	return nil
}

func validateYaml_TmpParameters(docs *[]interface{}) error {
	return nil
}


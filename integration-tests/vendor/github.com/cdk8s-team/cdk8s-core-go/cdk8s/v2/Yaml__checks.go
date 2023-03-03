//go:build !no_runtime_type_checking

// This is the core library of Cloud Development Kit (CDK) for Kubernetes (cdk8s). cdk8s apps synthesize into standard Kubernetes manifests which can be applied to any Kubernetes cluster.
package cdk8s

import (
	"fmt"
)

func validateYaml_FormatObjectsParameters(docs *[]interface{}) error {
	if docs == nil {
		return fmt.Errorf("parameter docs is required, but nil was provided")
	}

	return nil
}

func validateYaml_LoadParameters(urlOrFile *string) error {
	if urlOrFile == nil {
		return fmt.Errorf("parameter urlOrFile is required, but nil was provided")
	}

	return nil
}

func validateYaml_SaveParameters(filePath *string, docs *[]interface{}) error {
	if filePath == nil {
		return fmt.Errorf("parameter filePath is required, but nil was provided")
	}

	if docs == nil {
		return fmt.Errorf("parameter docs is required, but nil was provided")
	}

	return nil
}

func validateYaml_TmpParameters(docs *[]interface{}) error {
	if docs == nil {
		return fmt.Errorf("parameter docs is required, but nil was provided")
	}

	return nil
}


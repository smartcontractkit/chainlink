//go:build !no_runtime_type_checking

// This is the core library of Cloud Development Kit (CDK) for Kubernetes (cdk8s). cdk8s apps synthesize into standard Kubernetes manifests which can be applied to any Kubernetes cluster.
package cdk8s

import (
	"fmt"

	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
)

func (a *jsiiProxy_ApiObjectMetadataDefinition) validateAddParameters(key *string, value interface{}) error {
	if key == nil {
		return fmt.Errorf("parameter key is required, but nil was provided")
	}

	if value == nil {
		return fmt.Errorf("parameter value is required, but nil was provided")
	}

	return nil
}

func (a *jsiiProxy_ApiObjectMetadataDefinition) validateAddAnnotationParameters(key *string, value *string) error {
	if key == nil {
		return fmt.Errorf("parameter key is required, but nil was provided")
	}

	if value == nil {
		return fmt.Errorf("parameter value is required, but nil was provided")
	}

	return nil
}

func (a *jsiiProxy_ApiObjectMetadataDefinition) validateAddLabelParameters(key *string, value *string) error {
	if key == nil {
		return fmt.Errorf("parameter key is required, but nil was provided")
	}

	if value == nil {
		return fmt.Errorf("parameter value is required, but nil was provided")
	}

	return nil
}

func (a *jsiiProxy_ApiObjectMetadataDefinition) validateAddOwnerReferenceParameters(owner *OwnerReference) error {
	if owner == nil {
		return fmt.Errorf("parameter owner is required, but nil was provided")
	}
	if err := _jsii_.ValidateStruct(owner, func() string { return "parameter owner" }); err != nil {
		return err
	}

	return nil
}

func (a *jsiiProxy_ApiObjectMetadataDefinition) validateGetLabelParameters(key *string) error {
	if key == nil {
		return fmt.Errorf("parameter key is required, but nil was provided")
	}

	return nil
}

func validateNewApiObjectMetadataDefinitionParameters(options *ApiObjectMetadata) error {
	if err := _jsii_.ValidateStruct(options, func() string { return "parameter options" }); err != nil {
		return err
	}

	return nil
}


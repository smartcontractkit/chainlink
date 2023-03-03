// This is the core library of Cloud Development Kit (CDK) for Kubernetes (cdk8s). cdk8s apps synthesize into standard Kubernetes manifests which can be applied to any Kubernetes cluster.
package cdk8s

import (
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
	_init_ "github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2/jsii"

	"github.com/aws/constructs-go/constructs/v10"
)

// Utilities for generating unique and stable names.
type Names interface {
}

// The jsii proxy struct for Names
type jsiiProxy_Names struct {
	_ byte // padding
}

// Generates a unique and stable name compatible DNS_LABEL from RFC-1123 from a path.
//
// The generated name will:
//   - contain at most 63 characters
//   - contain only lowercase alphanumeric characters or ‘-’
//   - start with an alphanumeric character
//   - end with an alphanumeric character
//
// The generated name will have the form:
//   <comp0>-<comp1>-..-<compN>-<short-hash>
//
// Where <comp> are the path components (assuming they are is separated by
// "/").
//
// Note that if the total length is longer than 63 characters, we will trim
// the first components since the last components usually encode more meaning.
func Names_ToDnsLabel(scope constructs.Construct, options *NameOptions) *string {
	_init_.Initialize()

	if err := validateNames_ToDnsLabelParameters(scope, options); err != nil {
		panic(err)
	}
	var returns *string

	_jsii_.StaticInvoke(
		"cdk8s.Names",
		"toDnsLabel",
		[]interface{}{scope, options},
		&returns,
	)

	return returns
}

// Generates a unique and stable name compatible label key name segment and label value from a path.
//
// The name segment is required and must be 63 characters or less, beginning
// and ending with an alphanumeric character ([a-z0-9A-Z]) with dashes (-),
// underscores (_), dots (.), and alphanumerics between.
//
// Valid label values must be 63 characters or less and must be empty or
// begin and end with an alphanumeric character ([a-z0-9A-Z]) with dashes
// (-), underscores (_), dots (.), and alphanumerics between.
//
// The generated name will have the form:
//   <comp0><delim><comp1><delim>..<delim><compN><delim><short-hash>
//
// Where <comp> are the path components (assuming they are is separated by
// "/").
//
// Note that if the total length is longer than 63 characters, we will trim
// the first components since the last components usually encode more meaning.
func Names_ToLabelValue(scope constructs.Construct, options *NameOptions) *string {
	_init_.Initialize()

	if err := validateNames_ToLabelValueParameters(scope, options); err != nil {
		panic(err)
	}
	var returns *string

	_jsii_.StaticInvoke(
		"cdk8s.Names",
		"toLabelValue",
		[]interface{}{scope, options},
		&returns,
	)

	return returns
}


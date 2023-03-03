// This is the core library of Cloud Development Kit (CDK) for Kubernetes (cdk8s). cdk8s apps synthesize into standard Kubernetes manifests which can be applied to any Kubernetes cluster.
package cdk8s

import (
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
	_init_ "github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2/jsii"
)

// YAML utilities.
type Yaml interface {
}

// The jsii proxy struct for Yaml
type jsiiProxy_Yaml struct {
	_ byte // padding
}

// Deprecated: use `stringify(doc[, doc, ...])`
func Yaml_FormatObjects(docs *[]interface{}) *string {
	_init_.Initialize()

	if err := validateYaml_FormatObjectsParameters(docs); err != nil {
		panic(err)
	}
	var returns *string

	_jsii_.StaticInvoke(
		"cdk8s.Yaml",
		"formatObjects",
		[]interface{}{docs},
		&returns,
	)

	return returns
}

// Downloads a set of YAML documents (k8s manifest for example) from a URL or a file and returns them as javascript objects.
//
// Empty documents are filtered out.
//
// Returns: an array of objects, each represents a document inside the YAML.
func Yaml_Load(urlOrFile *string) *[]interface{} {
	_init_.Initialize()

	if err := validateYaml_LoadParameters(urlOrFile); err != nil {
		panic(err)
	}
	var returns *[]interface{}

	_jsii_.StaticInvoke(
		"cdk8s.Yaml",
		"load",
		[]interface{}{urlOrFile},
		&returns,
	)

	return returns
}

// Saves a set of objects as a multi-document YAML file.
func Yaml_Save(filePath *string, docs *[]interface{}) {
	_init_.Initialize()

	if err := validateYaml_SaveParameters(filePath, docs); err != nil {
		panic(err)
	}
	_jsii_.StaticInvokeVoid(
		"cdk8s.Yaml",
		"save",
		[]interface{}{filePath, docs},
	)
}

// Stringify a document (or multiple documents) into YAML.
//
// We convert undefined values to null, but ignore any documents that are
// undefined.
//
// Returns: a YAML string. Multiple docs are separated by `---`.
func Yaml_Stringify(docs ...interface{}) *string {
	_init_.Initialize()

	args := []interface{}{}
	for _, a := range docs {
		args = append(args, a)
	}

	var returns *string

	_jsii_.StaticInvoke(
		"cdk8s.Yaml",
		"stringify",
		args,
		&returns,
	)

	return returns
}

// Saves a set of YAML documents into a temp file (in /tmp).
//
// Returns: the path to the temporary file.
func Yaml_Tmp(docs *[]interface{}) *string {
	_init_.Initialize()

	if err := validateYaml_TmpParameters(docs); err != nil {
		panic(err)
	}
	var returns *string

	_jsii_.StaticInvoke(
		"cdk8s.Yaml",
		"tmp",
		[]interface{}{docs},
		&returns,
	)

	return returns
}


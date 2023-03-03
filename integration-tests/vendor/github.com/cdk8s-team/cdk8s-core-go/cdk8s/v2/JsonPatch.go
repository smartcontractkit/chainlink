// This is the core library of Cloud Development Kit (CDK) for Kubernetes (cdk8s). cdk8s apps synthesize into standard Kubernetes manifests which can be applied to any Kubernetes cluster.
package cdk8s

import (
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
	_init_ "github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2/jsii"
)

// Utility for applying RFC-6902 JSON-Patch to a document.
//
// Use the the `JsonPatch.apply(doc, ...ops)` function to apply a set of
// operations to a JSON document and return the result.
//
// Operations can be created using the factory methods `JsonPatch.add()`,
// `JsonPatch.remove()`, etc.
//
// Example:
//   const output = JsonPatch.apply(input,
//    JsonPatch.replace('/world/hi/there', 'goodbye'),
//    JsonPatch.add('/world/foo/', 'boom'),
//    JsonPatch.remove('/hello'));
//
type JsonPatch interface {
}

// The jsii proxy struct for JsonPatch
type jsiiProxy_JsonPatch struct {
	_ byte // padding
}

// Adds a value to an object or inserts it into an array.
//
// In the case of an
// array, the value is inserted before the given index. The - character can be
// used instead of an index to insert at the end of an array.
//
// Example:
//   JsonPatch.add('/biscuits/1', { "name": "Ginger Nut" })
//
func JsonPatch_Add(path *string, value interface{}) JsonPatch {
	_init_.Initialize()

	if err := validateJsonPatch_AddParameters(path, value); err != nil {
		panic(err)
	}
	var returns JsonPatch

	_jsii_.StaticInvoke(
		"cdk8s.JsonPatch",
		"add",
		[]interface{}{path, value},
		&returns,
	)

	return returns
}

// Applies a set of JSON-Patch (RFC-6902) operations to `document` and returns the result.
//
// Returns: The result document.
func JsonPatch_Apply(document interface{}, ops ...JsonPatch) interface{} {
	_init_.Initialize()

	if err := validateJsonPatch_ApplyParameters(document); err != nil {
		panic(err)
	}
	args := []interface{}{document}
	for _, a := range ops {
		args = append(args, a)
	}

	var returns interface{}

	_jsii_.StaticInvoke(
		"cdk8s.JsonPatch",
		"apply",
		args,
		&returns,
	)

	return returns
}

// Copies a value from one location to another within the JSON document.
//
// Both
// from and path are JSON Pointers.
//
// Example:
//   JsonPatch.copy('/biscuits/0', '/best_biscuit')
//
func JsonPatch_Copy(from *string, path *string) JsonPatch {
	_init_.Initialize()

	if err := validateJsonPatch_CopyParameters(from, path); err != nil {
		panic(err)
	}
	var returns JsonPatch

	_jsii_.StaticInvoke(
		"cdk8s.JsonPatch",
		"copy",
		[]interface{}{from, path},
		&returns,
	)

	return returns
}

// Moves a value from one location to the other.
//
// Both from and path are JSON Pointers.
//
// Example:
//   JsonPatch.move('/biscuits', '/cookies')
//
func JsonPatch_Move(from *string, path *string) JsonPatch {
	_init_.Initialize()

	if err := validateJsonPatch_MoveParameters(from, path); err != nil {
		panic(err)
	}
	var returns JsonPatch

	_jsii_.StaticInvoke(
		"cdk8s.JsonPatch",
		"move",
		[]interface{}{from, path},
		&returns,
	)

	return returns
}

// Removes a value from an object or array.
//
// Example:
//   JsonPatch.remove('/biscuits/0')
//
func JsonPatch_Remove(path *string) JsonPatch {
	_init_.Initialize()

	if err := validateJsonPatch_RemoveParameters(path); err != nil {
		panic(err)
	}
	var returns JsonPatch

	_jsii_.StaticInvoke(
		"cdk8s.JsonPatch",
		"remove",
		[]interface{}{path},
		&returns,
	)

	return returns
}

// Replaces a value.
//
// Equivalent to a “remove” followed by an “add”.
//
// Example:
//   JsonPatch.replace('/biscuits/0/name', 'Chocolate Digestive')
//
func JsonPatch_Replace(path *string, value interface{}) JsonPatch {
	_init_.Initialize()

	if err := validateJsonPatch_ReplaceParameters(path, value); err != nil {
		panic(err)
	}
	var returns JsonPatch

	_jsii_.StaticInvoke(
		"cdk8s.JsonPatch",
		"replace",
		[]interface{}{path, value},
		&returns,
	)

	return returns
}

// Tests that the specified value is set in the document.
//
// If the test fails,
// then the patch as a whole should not apply.
//
// Example:
//   JsonPatch.test('/best_biscuit/name', 'Choco Leibniz')
//
func JsonPatch_Test(path *string, value interface{}) JsonPatch {
	_init_.Initialize()

	if err := validateJsonPatch_TestParameters(path, value); err != nil {
		panic(err)
	}
	var returns JsonPatch

	_jsii_.StaticInvoke(
		"cdk8s.JsonPatch",
		"test",
		[]interface{}{path, value},
		&returns,
	)

	return returns
}


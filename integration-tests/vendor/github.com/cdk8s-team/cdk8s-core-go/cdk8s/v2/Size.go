// This is the core library of Cloud Development Kit (CDK) for Kubernetes (cdk8s). cdk8s apps synthesize into standard Kubernetes manifests which can be applied to any Kubernetes cluster.
package cdk8s

import (
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
	_init_ "github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2/jsii"
)

// Represents the amount of digital storage.
//
// The amount can be specified either as a literal value (e.g: `10`) which
// cannot be negative.
//
// When the amount is passed as a token, unit conversion is not possible.
type Size interface {
	// Return this storage as a total number of gibibytes.
	ToGibibytes(opts *SizeConversionOptions) *float64
	// Return this storage as a total number of kibibytes.
	ToKibibytes(opts *SizeConversionOptions) *float64
	// Return this storage as a total number of mebibytes.
	ToMebibytes(opts *SizeConversionOptions) *float64
	// Return this storage as a total number of pebibytes.
	ToPebibytes(opts *SizeConversionOptions) *float64
	// Return this storage as a total number of tebibytes.
	ToTebibytes(opts *SizeConversionOptions) *float64
}

// The jsii proxy struct for Size
type jsiiProxy_Size struct {
	_ byte // padding
}

// Create a Storage representing an amount gibibytes.
//
// 1 GiB = 1024 MiB.
func Size_Gibibytes(amount *float64) Size {
	_init_.Initialize()

	if err := validateSize_GibibytesParameters(amount); err != nil {
		panic(err)
	}
	var returns Size

	_jsii_.StaticInvoke(
		"cdk8s.Size",
		"gibibytes",
		[]interface{}{amount},
		&returns,
	)

	return returns
}

// Create a Storage representing an amount kibibytes.
//
// 1 KiB = 1024 bytes.
func Size_Kibibytes(amount *float64) Size {
	_init_.Initialize()

	if err := validateSize_KibibytesParameters(amount); err != nil {
		panic(err)
	}
	var returns Size

	_jsii_.StaticInvoke(
		"cdk8s.Size",
		"kibibytes",
		[]interface{}{amount},
		&returns,
	)

	return returns
}

// Create a Storage representing an amount mebibytes.
//
// 1 MiB = 1024 KiB.
func Size_Mebibytes(amount *float64) Size {
	_init_.Initialize()

	if err := validateSize_MebibytesParameters(amount); err != nil {
		panic(err)
	}
	var returns Size

	_jsii_.StaticInvoke(
		"cdk8s.Size",
		"mebibytes",
		[]interface{}{amount},
		&returns,
	)

	return returns
}

// Create a Storage representing an amount pebibytes.
//
// 1 PiB = 1024 TiB.
func Size_Pebibyte(amount *float64) Size {
	_init_.Initialize()

	if err := validateSize_PebibyteParameters(amount); err != nil {
		panic(err)
	}
	var returns Size

	_jsii_.StaticInvoke(
		"cdk8s.Size",
		"pebibyte",
		[]interface{}{amount},
		&returns,
	)

	return returns
}

// Create a Storage representing an amount tebibytes.
//
// 1 TiB = 1024 GiB.
func Size_Tebibytes(amount *float64) Size {
	_init_.Initialize()

	if err := validateSize_TebibytesParameters(amount); err != nil {
		panic(err)
	}
	var returns Size

	_jsii_.StaticInvoke(
		"cdk8s.Size",
		"tebibytes",
		[]interface{}{amount},
		&returns,
	)

	return returns
}

func (s *jsiiProxy_Size) ToGibibytes(opts *SizeConversionOptions) *float64 {
	if err := s.validateToGibibytesParameters(opts); err != nil {
		panic(err)
	}
	var returns *float64

	_jsii_.Invoke(
		s,
		"toGibibytes",
		[]interface{}{opts},
		&returns,
	)

	return returns
}

func (s *jsiiProxy_Size) ToKibibytes(opts *SizeConversionOptions) *float64 {
	if err := s.validateToKibibytesParameters(opts); err != nil {
		panic(err)
	}
	var returns *float64

	_jsii_.Invoke(
		s,
		"toKibibytes",
		[]interface{}{opts},
		&returns,
	)

	return returns
}

func (s *jsiiProxy_Size) ToMebibytes(opts *SizeConversionOptions) *float64 {
	if err := s.validateToMebibytesParameters(opts); err != nil {
		panic(err)
	}
	var returns *float64

	_jsii_.Invoke(
		s,
		"toMebibytes",
		[]interface{}{opts},
		&returns,
	)

	return returns
}

func (s *jsiiProxy_Size) ToPebibytes(opts *SizeConversionOptions) *float64 {
	if err := s.validateToPebibytesParameters(opts); err != nil {
		panic(err)
	}
	var returns *float64

	_jsii_.Invoke(
		s,
		"toPebibytes",
		[]interface{}{opts},
		&returns,
	)

	return returns
}

func (s *jsiiProxy_Size) ToTebibytes(opts *SizeConversionOptions) *float64 {
	if err := s.validateToTebibytesParameters(opts); err != nil {
		panic(err)
	}
	var returns *float64

	_jsii_.Invoke(
		s,
		"toTebibytes",
		[]interface{}{opts},
		&returns,
	)

	return returns
}


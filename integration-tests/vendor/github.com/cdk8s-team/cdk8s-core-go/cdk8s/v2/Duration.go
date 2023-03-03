// This is the core library of Cloud Development Kit (CDK) for Kubernetes (cdk8s). cdk8s apps synthesize into standard Kubernetes manifests which can be applied to any Kubernetes cluster.
package cdk8s

import (
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
	_init_ "github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2/jsii"
)

// Represents a length of time.
//
// The amount can be specified either as a literal value (e.g: `10`) which
// cannot be negative.
type Duration interface {
	// Return the total number of days in this Duration.
	//
	// Returns: the value of this `Duration` expressed in Days.
	ToDays(opts *TimeConversionOptions) *float64
	// Return the total number of hours in this Duration.
	//
	// Returns: the value of this `Duration` expressed in Hours.
	ToHours(opts *TimeConversionOptions) *float64
	// Turn this duration into a human-readable string.
	ToHumanString() *string
	// Return an ISO 8601 representation of this period.
	//
	// Returns: a string starting with 'PT' describing the period.
	// See: https://www.iso.org/fr/standard/70907.html
	//
	ToIsoString() *string
	// Return the total number of milliseconds in this Duration.
	//
	// Returns: the value of this `Duration` expressed in Milliseconds.
	ToMilliseconds(opts *TimeConversionOptions) *float64
	// Return the total number of minutes in this Duration.
	//
	// Returns: the value of this `Duration` expressed in Minutes.
	ToMinutes(opts *TimeConversionOptions) *float64
	// Return the total number of seconds in this Duration.
	//
	// Returns: the value of this `Duration` expressed in Seconds.
	ToSeconds(opts *TimeConversionOptions) *float64
	// Return unit of Duration.
	UnitLabel() *string
}

// The jsii proxy struct for Duration
type jsiiProxy_Duration struct {
	_ byte // padding
}

// Create a Duration representing an amount of days.
//
// Returns: a new `Duration` representing `amount` Days.
func Duration_Days(amount *float64) Duration {
	_init_.Initialize()

	if err := validateDuration_DaysParameters(amount); err != nil {
		panic(err)
	}
	var returns Duration

	_jsii_.StaticInvoke(
		"cdk8s.Duration",
		"days",
		[]interface{}{amount},
		&returns,
	)

	return returns
}

// Create a Duration representing an amount of hours.
//
// Returns: a new `Duration` representing `amount` Hours.
func Duration_Hours(amount *float64) Duration {
	_init_.Initialize()

	if err := validateDuration_HoursParameters(amount); err != nil {
		panic(err)
	}
	var returns Duration

	_jsii_.StaticInvoke(
		"cdk8s.Duration",
		"hours",
		[]interface{}{amount},
		&returns,
	)

	return returns
}

// Create a Duration representing an amount of milliseconds.
//
// Returns: a new `Duration` representing `amount` ms.
func Duration_Millis(amount *float64) Duration {
	_init_.Initialize()

	if err := validateDuration_MillisParameters(amount); err != nil {
		panic(err)
	}
	var returns Duration

	_jsii_.StaticInvoke(
		"cdk8s.Duration",
		"millis",
		[]interface{}{amount},
		&returns,
	)

	return returns
}

// Create a Duration representing an amount of minutes.
//
// Returns: a new `Duration` representing `amount` Minutes.
func Duration_Minutes(amount *float64) Duration {
	_init_.Initialize()

	if err := validateDuration_MinutesParameters(amount); err != nil {
		panic(err)
	}
	var returns Duration

	_jsii_.StaticInvoke(
		"cdk8s.Duration",
		"minutes",
		[]interface{}{amount},
		&returns,
	)

	return returns
}

// Parse a period formatted according to the ISO 8601 standard.
//
// Returns: the parsed `Duration`.
// See: https://www.iso.org/fr/standard/70907.html
//
func Duration_Parse(duration *string) Duration {
	_init_.Initialize()

	if err := validateDuration_ParseParameters(duration); err != nil {
		panic(err)
	}
	var returns Duration

	_jsii_.StaticInvoke(
		"cdk8s.Duration",
		"parse",
		[]interface{}{duration},
		&returns,
	)

	return returns
}

// Create a Duration representing an amount of seconds.
//
// Returns: a new `Duration` representing `amount` Seconds.
func Duration_Seconds(amount *float64) Duration {
	_init_.Initialize()

	if err := validateDuration_SecondsParameters(amount); err != nil {
		panic(err)
	}
	var returns Duration

	_jsii_.StaticInvoke(
		"cdk8s.Duration",
		"seconds",
		[]interface{}{amount},
		&returns,
	)

	return returns
}

func (d *jsiiProxy_Duration) ToDays(opts *TimeConversionOptions) *float64 {
	if err := d.validateToDaysParameters(opts); err != nil {
		panic(err)
	}
	var returns *float64

	_jsii_.Invoke(
		d,
		"toDays",
		[]interface{}{opts},
		&returns,
	)

	return returns
}

func (d *jsiiProxy_Duration) ToHours(opts *TimeConversionOptions) *float64 {
	if err := d.validateToHoursParameters(opts); err != nil {
		panic(err)
	}
	var returns *float64

	_jsii_.Invoke(
		d,
		"toHours",
		[]interface{}{opts},
		&returns,
	)

	return returns
}

func (d *jsiiProxy_Duration) ToHumanString() *string {
	var returns *string

	_jsii_.Invoke(
		d,
		"toHumanString",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (d *jsiiProxy_Duration) ToIsoString() *string {
	var returns *string

	_jsii_.Invoke(
		d,
		"toIsoString",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (d *jsiiProxy_Duration) ToMilliseconds(opts *TimeConversionOptions) *float64 {
	if err := d.validateToMillisecondsParameters(opts); err != nil {
		panic(err)
	}
	var returns *float64

	_jsii_.Invoke(
		d,
		"toMilliseconds",
		[]interface{}{opts},
		&returns,
	)

	return returns
}

func (d *jsiiProxy_Duration) ToMinutes(opts *TimeConversionOptions) *float64 {
	if err := d.validateToMinutesParameters(opts); err != nil {
		panic(err)
	}
	var returns *float64

	_jsii_.Invoke(
		d,
		"toMinutes",
		[]interface{}{opts},
		&returns,
	)

	return returns
}

func (d *jsiiProxy_Duration) ToSeconds(opts *TimeConversionOptions) *float64 {
	if err := d.validateToSecondsParameters(opts); err != nil {
		panic(err)
	}
	var returns *float64

	_jsii_.Invoke(
		d,
		"toSeconds",
		[]interface{}{opts},
		&returns,
	)

	return returns
}

func (d *jsiiProxy_Duration) UnitLabel() *string {
	var returns *string

	_jsii_.Invoke(
		d,
		"unitLabel",
		nil, // no parameters
		&returns,
	)

	return returns
}


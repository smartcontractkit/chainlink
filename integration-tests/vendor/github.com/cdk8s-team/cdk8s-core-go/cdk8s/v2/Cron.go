// This is the core library of Cloud Development Kit (CDK) for Kubernetes (cdk8s). cdk8s apps synthesize into standard Kubernetes manifests which can be applied to any Kubernetes cluster.
package cdk8s

import (
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
	_init_ "github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2/jsii"
)

// Represents a cron schedule.
type Cron interface {
	// Retrieve the expression for this schedule.
	ExpressionString() *string
}

// The jsii proxy struct for Cron
type jsiiProxy_Cron struct {
	_ byte // padding
}

func (j *jsiiProxy_Cron) ExpressionString() *string {
	var returns *string
	_jsii_.Get(
		j,
		"expressionString",
		&returns,
	)
	return returns
}


func NewCron(cronOptions *CronOptions) Cron {
	_init_.Initialize()

	if err := validateNewCronParameters(cronOptions); err != nil {
		panic(err)
	}
	j := jsiiProxy_Cron{}

	_jsii_.Create(
		"cdk8s.Cron",
		[]interface{}{cronOptions},
		&j,
	)

	return &j
}

func NewCron_Override(c Cron, cronOptions *CronOptions) {
	_init_.Initialize()

	_jsii_.Create(
		"cdk8s.Cron",
		[]interface{}{cronOptions},
		c,
	)
}

// Create a cron schedule which runs first day of January every year.
func Cron_Annually() Cron {
	_init_.Initialize()

	var returns Cron

	_jsii_.StaticInvoke(
		"cdk8s.Cron",
		"annually",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Create a cron schedule which runs every day at midnight.
func Cron_Daily() Cron {
	_init_.Initialize()

	var returns Cron

	_jsii_.StaticInvoke(
		"cdk8s.Cron",
		"daily",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Create a cron schedule which runs every minute.
func Cron_EveryMinute() Cron {
	_init_.Initialize()

	var returns Cron

	_jsii_.StaticInvoke(
		"cdk8s.Cron",
		"everyMinute",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Create a cron schedule which runs every hour.
func Cron_Hourly() Cron {
	_init_.Initialize()

	var returns Cron

	_jsii_.StaticInvoke(
		"cdk8s.Cron",
		"hourly",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Create a cron schedule which runs first day of every month.
func Cron_Monthly() Cron {
	_init_.Initialize()

	var returns Cron

	_jsii_.StaticInvoke(
		"cdk8s.Cron",
		"monthly",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Create a custom cron schedule from a set of cron fields.
func Cron_Schedule(options *CronOptions) Cron {
	_init_.Initialize()

	if err := validateCron_ScheduleParameters(options); err != nil {
		panic(err)
	}
	var returns Cron

	_jsii_.StaticInvoke(
		"cdk8s.Cron",
		"schedule",
		[]interface{}{options},
		&returns,
	)

	return returns
}

// Create a cron schedule which runs every week on Sunday.
func Cron_Weekly() Cron {
	_init_.Initialize()

	var returns Cron

	_jsii_.StaticInvoke(
		"cdk8s.Cron",
		"weekly",
		nil, // no parameters
		&returns,
	)

	return returns
}


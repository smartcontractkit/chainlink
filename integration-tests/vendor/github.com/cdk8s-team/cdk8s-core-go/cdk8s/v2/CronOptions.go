// This is the core library of Cloud Development Kit (CDK) for Kubernetes (cdk8s). cdk8s apps synthesize into standard Kubernetes manifests which can be applied to any Kubernetes cluster.
package cdk8s


// Options to configure a cron expression.
//
// All fields are strings so you can use complex expressions. Absence of
// a field implies '*'.
type CronOptions struct {
	// The day of the month to run this rule at.
	Day *string `field:"optional" json:"day" yaml:"day"`
	// The hour to run this rule at.
	Hour *string `field:"optional" json:"hour" yaml:"hour"`
	// The minute to run this rule at.
	Minute *string `field:"optional" json:"minute" yaml:"minute"`
	// The month to run this rule at.
	Month *string `field:"optional" json:"month" yaml:"month"`
	// The day of the week to run this rule at.
	WeekDay *string `field:"optional" json:"weekDay" yaml:"weekDay"`
}


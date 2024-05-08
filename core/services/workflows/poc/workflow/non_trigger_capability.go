package workflow

type nonTriggerCapability struct {
	inputs           map[string]string
	stepDependencies []string
	ref              string
}

func (c *nonTriggerCapability) Inputs() map[string]string {
	return c.inputs
}

func (c *nonTriggerCapability) Output() string {
	return c.ref
}

func (c *nonTriggerCapability) StepDependencies() []string {
	return c.stepDependencies
}

func (c *nonTriggerCapability) private() {}

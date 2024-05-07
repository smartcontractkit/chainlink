package workflow

import (
	"fmt"
)

type nonTriggerCapability struct {
	inputs map[string]any
	ref    string
}

func (c *nonTriggerCapability) Inputs() map[string]any {
	return c.inputs
}

func (c *nonTriggerCapability) Output() string {
	return fmt.Sprintf("$(%s.outputs)", c.ref)
}

func (c *nonTriggerCapability) private() {}

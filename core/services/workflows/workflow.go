package workflows

import "gopkg.in/yaml.v3"

type Capability struct {
	Type   string         `yaml:"type"`
	Ref    string         `yaml:"ref"`
	Inputs map[string]any `yaml:"inputs"`
	Config map[string]any `yaml:"config"`
}

type Workflow struct {
	Triggers  []Capability `yaml:"triggers"`
	Actions   []Capability `yaml:"actions"`
	Consensus []Capability `yaml:"consensus"`
	Targets   []Capability `yaml:"targets"`
}

func Parse(yamlWorkflow string) (*Workflow, error) {
	wf := &Workflow{}
	err := yaml.Unmarshal([]byte(yamlWorkflow), wf)
	return wf, err
}

package cmd

type Test struct {
	Name string
	Path string
}

type TestConf struct {
	ID                    string   `yaml:"id" json:"id"`
	Path                  string   `yaml:"path" json:"path"`
	TestType              string   `yaml:"test-type" json:"testType"`
	RunsOn                string   `yaml:"runs-on" json:"runsOn"`
	Cmd                   string   `yaml:"cmd" json:"cmd"`
	RemoteRunnerTestSuite string   `yaml:"remote-runner-test-suite" json:"remoteRunnerTestSuite"`
	PyroscopeEnv          string   `yaml:"pyroscope-env" json:"pyroscopeEnv"`
	Trigger               []string `yaml:"trigger" json:"trigger"`
}

type Config struct {
	Tests []TestConf `yaml:"test-runner-matrix"`
}

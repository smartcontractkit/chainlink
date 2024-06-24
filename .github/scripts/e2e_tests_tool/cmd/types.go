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
	TestCmd               string   `yaml:"test-cmd" json:"testCmd"`
	TestConfigOverride    string   `yaml:"test-config-override" json:"testConfigOverride"`
	RemoteRunnerTestSuite string   `yaml:"remote-runner-test-suite" json:"remoteRunnerTestSuite"`
	PyroscopeEnv          string   `yaml:"pyroscope-env" json:"pyroscopeEnv"`
	Trigger               []string `yaml:"trigger" json:"trigger"`
}

type Config struct {
	Tests []TestConf `yaml:"test-runner-matrix"`
}

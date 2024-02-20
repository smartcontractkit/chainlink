package waspconfig

import (
	"errors"
)

type WaspAutoBuildConfig struct {
	Namespace           *string `toml:"namespace"`
	RepoImageVersionURI *string `toml:"repo_image_version_uri"`
	TestBinaryName      *string `toml:"test_binary_name"`
	TestName            *string `toml:"test_name"`
	TestTimeout         *string `toml:"test_timeout"`
	KeepJobs            bool    `toml:"keep_jobs"`
	WaspLogLevel        *string `toml:"wasp_log_level"`
	WaspJobs            *string `toml:"wasp_jobs"`
	UpdateImage         bool    `toml:"update_image"`
}

func (c *WaspAutoBuildConfig) Validate() error {
	if c.Namespace == nil || *c.Namespace == "" {
		return errors.New("WASP namespace name should not be empty, see WASP docs to setup it")
	}
	if c.RepoImageVersionURI == nil || *c.RepoImageVersionURI == "" {
		return errors.New("WASP image URI is empty, must be ${registry}/${repo}:${tag}")
	}
	if c.TestBinaryName == nil || *c.TestBinaryName == "" {
		return errors.New("WASP test binary is empty, should be 'ocr.test', run 'go test -c ./...' in load test dir to figure out the name")
	}
	if c.TestName == nil || *c.TestName == "" {
		return errors.New("WASP test name is empty, should be a name of go test you want to run")
	}
	if c.TestTimeout == nil || *c.TestTimeout == "" {
		return errors.New("WASP test timeout should be in Go time format: '1w2d3h4m5s'")
	}
	if c.WaspJobs == nil || *c.WaspJobs == "" {
		return errors.New("WASP jobs are empty, amount of pods to spin up in k8s")
	}
	return nil
}

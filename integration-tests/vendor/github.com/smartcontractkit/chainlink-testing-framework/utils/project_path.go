// Package utils contains some common paths used in configuration and tests
package utils

import (
	"path/filepath"
	"runtime"
)

var (
	_, b, _, _ = runtime.Caller(0)
	// ProjectRoot Root folder of this project
	ProjectRoot = filepath.Join(filepath.Dir(b), "/..")
	// SuiteRoot test suite root
	SuiteRoot = filepath.Join(ProjectRoot, "suite")
	// SoakRoot test suite soak root
	SoakRoot = filepath.Join(SuiteRoot, "soak")
	// PresetRoot root folder for environments preset
	PresetRoot = filepath.Join(ProjectRoot, "preset")
	// ContractsDir path to our contracts
	ContractsDir = filepath.Join(ProjectRoot, "contracts")
	// EthereumContractsDir path to our ethereum contracts
	EthereumContractsDir = filepath.Join(ContractsDir, "ethereum")
	// RemoteRunnerConfigLocation is the path to the remote runner config
	RemoteRunnerConfigLocation = filepath.Join(SuiteRoot, "remote_runner_config.yaml")
)

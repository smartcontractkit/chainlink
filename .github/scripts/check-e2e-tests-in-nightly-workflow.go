package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// JobConfig represents the structure of the jobs in the YAML file.
type JobConfig struct {
	Jobs map[string]struct {
		Strategy struct {
			Matrix struct {
				Test []struct {
					Path     string `yaml:"path"`
					TestOpts string `yaml:"testOpts"`
				} `yaml:"test"`
			} `yaml:"matrix"`
		} `yaml:"strategy"`
	} `yaml:"jobs"`
}

func main() {
	// Check command-line arguments
	if len(os.Args) != 4 {
		fmt.Println("Usage: go run check_test_in_pipeline.go <yaml file> <test name> <test file path>")
		os.Exit(1)
	}
	filename := os.Args[1]
	testName := os.Args[2]
	testFilePath := os.Args[3]

	// Read the YAML file
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading YAML file: %s\n", err)
		os.Exit(1)
	}

	// Parse the YAML
	var config JobConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		fmt.Printf("Error parsing YAML: %s\n", err)
		os.Exit(1)
	}

	// Look for the test in the configuration
	found := false
	for _, job := range config.Jobs {
		for _, test := range job.Strategy.Matrix.Test {
			if test.Path == testFilePath {
				// Check if -test.run is specified and if it contains the testName
				if strings.Contains(test.TestOpts, "-test.run") {
					if strings.Contains(test.TestOpts, testName) {
						fmt.Printf("Test '%s' is specifically included in the nightly CI pipeline for file '%s'.\n", testName, testFilePath)
						found = true
						break
					}
				} else {
					// If -test.run is not specified, assume all tests in the file are included
					fmt.Printf("All tests in '%s' are included in the nightly CI pipeline, including '%s'.\n", testFilePath, testName)
					found = true
					break
				}
			}
		}
		if found {
			break
		}
	}

	if !found {
		fmt.Printf("Test '%s' is NOT included in the nightly CI pipeline for file '%s'.\n", testName, testFilePath)
	}
}

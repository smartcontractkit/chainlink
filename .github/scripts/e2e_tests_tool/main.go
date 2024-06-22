package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"gopkg.in/yaml.v2"
)

// Test defines the structure of a test entry in the YAML file.
type Test struct {
	ID                    string   `yaml:"id" json:"id"`
	Name                  string   `yaml:"name" json:"name"`
	Path                  string   `yaml:"path" json:"path"`
	TestType              string   `yaml:"test-type" json:"testType"`
	RunsOn                string   `yaml:"runs-on" json:"runsOn"`
	Cmd                   string   `yaml:"cmd" json:"cmd"`
	RemoteRunnerTestSuite string   `yaml:"remote-runner-test-suite" json:"remoteRunnerTestSuite"`
	PyroscopeEnv          string   `yaml:"pyroscope-env" json:"pyroscopeEnv"`
	Trigger               []string `yaml:"trigger" json:"trigger"`
}

// Config represents the tests configuration.
type Config struct {
	Tests []Test `yaml:"test-runner-matrix"`
}

// Filter tests based on name, trigger, test type, and test IDs.
func FilterTests(tests []Test, names, trigger, testType, ids string) []Test {
	nameFilter := strings.Split(names, ",")
	triggerFilter := trigger
	typeFilter := testType
	idFilter := strings.Split(ids, ",")

	var filteredTests []Test

	for _, test := range tests {
		nameMatch := names == "" || contains(nameFilter, test.Name)
		triggerMatch := trigger == "" || contains(test.Trigger, triggerFilter)
		typeMatch := testType == "" || test.TestType == typeFilter
		idMatch := ids == "*" || ids == "" || contains(idFilter, test.ID)

		if nameMatch && triggerMatch && typeMatch && idMatch {
			filteredTests = append(filteredTests, test)
		}
	}

	return filteredTests
}

// Main function including the new test-ids flag.
func main() {
	yamlFile := flag.String("file", ".github/e2e-tests.yml", "Path to the YAML file")
	names := flag.String("name", "", "Comma-separated list of test names to filter by")
	trigger := flag.String("trigger", "", "Trigger filter")
	testType := flag.String("test-type", "", "Type of test to filter by")
	testIDs := flag.String("test-ids", "*", "Comma-separated list of test IDs to filter by")

	flag.Parse()

	data, err := ioutil.ReadFile(*yamlFile)
	if err != nil {
		log.Fatalf("Error reading YAML file: %v", err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Error parsing YAML data: %v", err)
	}

	filteredTests := FilterTests(config.Tests, *names, *trigger, *testType, *testIDs)
	matrix := map[string][]Test{"tests": filteredTests}
	matrixJSON, err := json.Marshal(matrix)
	if err != nil {
		log.Fatalf("Error marshaling matrix to JSON: %v", err)
	}

	fmt.Printf("%s", matrixJSON)
}

// Utility function to check if a slice contains a string.
func contains(slice []string, element string) bool {
	for _, s := range slice {
		if s == element {
			return true
		}
	}
	return false
}

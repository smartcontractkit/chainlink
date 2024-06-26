package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// Filter tests based on workflow, test type, and test IDs.
func filterTests(tests []CITestConf, workflow, testType, ids string) []CITestConf {
	workflowFilter := workflow
	typeFilter := testType
	idFilter := strings.Split(ids, ",")

	var filteredTests []CITestConf

	for _, test := range tests {
		workflowMatch := workflow == "" || contains(test.Workflows, workflowFilter)
		typeMatch := testType == "" || test.TestType == typeFilter
		idMatch := ids == "*" || ids == "" || contains(idFilter, test.ID)

		if workflowMatch && typeMatch && idMatch {
			filteredTests = append(filteredTests, test)
		}
	}

	return filteredTests
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

// filterCmd represents the filter command
var filterCmd = &cobra.Command{
	Use:   "filter",
	Short: "Filter test configurations based on specified criteria",
	Long: `Filters tests from a YAML configuration based on name, workflow, test type, and test IDs.
Example usage:
./e2e_tests_tool filter --file .github/e2e-tests.yml --workflow "Run Nightly E2E Tests" --test-type "docker" --test-ids "test1,test2"`,
	Run: func(cmd *cobra.Command, args []string) {
		yamlFile, _ := cmd.Flags().GetString("file")
		workflow, _ := cmd.Flags().GetString("workflow")
		testType, _ := cmd.Flags().GetString("test-type")
		testIDs, _ := cmd.Flags().GetString("test-ids")

		data, err := ioutil.ReadFile(yamlFile)
		if err != nil {
			log.Fatalf("Error reading YAML file: %v", err)
		}

		var config Config
		err = yaml.Unmarshal(data, &config)
		if err != nil {
			log.Fatalf("Error parsing YAML data: %v", err)
		}

		filteredTests := filterTests(config.Tests, workflow, testType, testIDs)
		matrix := map[string][]CITestConf{"tests": filteredTests}
		matrixJSON, err := json.Marshal(matrix)
		if err != nil {
			log.Fatalf("Error marshaling matrix to JSON: %v", err)
		}

		fmt.Printf("%s", matrixJSON)
	},
}

func init() {
	filterCmd.Flags().StringP("file", "f", "", "Path to the YAML file")
	filterCmd.Flags().StringP("workflow", "t", "", "Workflow filter")
	filterCmd.Flags().StringP("test-type", "y", "", "Type of test to filter by")
	filterCmd.Flags().StringP("test-ids", "i", "*", "Comma-separated list of test IDs to filter by")
	filterCmd.MarkFlagRequired("file")
}

package cmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// Filter tests based on workflow, test type, and test IDs.
func filterTests(allTests []CITestConf, workflow, testType, ids string) []CITestConf {
	workflowFilter := workflow
	typeFilter := testType
	idFilter := strings.Split(ids, ",")

	var filteredTests []CITestConf

	for _, test := range allTests {
		workflowMatch := workflow == "" || contains(test.Workflows, workflowFilter)
		typeMatch := testType == "" || test.TestType == typeFilter
		idMatch := ids == "*" || ids == "" || contains(idFilter, test.ID)

		if workflowMatch && typeMatch && idMatch {
			test.IDSanitized = sanitizeTestID(test.ID)
			filteredTests = append(filteredTests, test)
		}
	}

	return filteredTests
}

func filterAndMergeTests(allTests []CITestConf, base64Tests string) ([]CITestConf, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(base64Tests)
	if err != nil {
		return nil, err
	}
	var decodedTests []CITestConf
	err = yaml.Unmarshal(decodedBytes, &decodedTests)
	if err != nil {
		log.Fatalf("Error parsing YAML data: %v", err)
	}

	idFilter := make(map[string]CITestConf)
	for _, dt := range decodedTests {
		idFilter[dt.ID] = dt
	}

	var filteredTests []CITestConf
	for _, test := range allTests {
		if decodedTest, exists := idFilter[test.ID]; exists {
			// Apply config overrides
			for k, v := range decodedTest.TestConfigOverrides {
				if test.TestConfigOverrides == nil {
					test.TestConfigOverrides = make(map[string]string)
				}
				test.TestConfigOverrides[k] = v
			}
			test.IDSanitized = sanitizeTestID(test.ID)
			filteredTests = append(filteredTests, test)
		}
	}

	return filteredTests, nil
}

func sanitizeTestID(id string) string {
	// Define a regular expression that matches any character not a letter, digit, hyphen
	re := regexp.MustCompile(`[^a-zA-Z0-9-_]+`)
	// Replace all occurrences of disallowed characters with "_"
	return re.ReplaceAllString(id, "_")
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
	Run: func(cmd *cobra.Command, _ []string) {
		yamlFile, _ := cmd.Flags().GetString("file")
		workflow, _ := cmd.Flags().GetString("workflow")
		testType, _ := cmd.Flags().GetString("test-type")
		testIDs, _ := cmd.Flags().GetString("test-ids")
		testMap, _ := cmd.Flags().GetString("test-map")

		data, err := os.ReadFile(yamlFile)
		if err != nil {
			log.Fatalf("Error reading YAML file: %v", err)
		}

		var config Config
		err = yaml.Unmarshal(data, &config)
		if err != nil {
			log.Fatalf("Error parsing YAML data: %v", err)
		}

		var filteredTests []CITestConf
		if testMap == "" {
			filteredTests = filterTests(config.Tests, workflow, testType, testIDs)
		} else {
			filteredTests, err = filterAndMergeTests(config.Tests, testMap)
			if err != nil {
				log.Fatalf("Error filtering and merging tests: %v", err)
			}
		}
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
	filterCmd.Flags().String("test-map", "", "Base64 encoded list of tests (YML objects) to filter by. Can include test-config-overrides for each test.")
	filterCmd.Flags().StringP("test-ids", "i", "*", "Comma-separated list of test IDs to filter by")
	filterCmd.Flags().StringP("test-type", "y", "", "Type of test to filter by")
	filterCmd.Flags().StringP("workflow", "t", "", "Workflow filter")
	err := filterCmd.MarkFlagRequired("file")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marking flag as required: %v\n", err)
		os.Exit(1)
	}
}

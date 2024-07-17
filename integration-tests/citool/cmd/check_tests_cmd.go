package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

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

var checkTestsCmd = &cobra.Command{
	Use:   "check-tests [directory] [yaml file]",
	Short: "Check if all tests in a directory are included in the test configurations YAML file",
	Args:  cobra.ExactArgs(2),
	Run: func(_ *cobra.Command, args []string) {
		directory := args[0]
		yamlFile := args[1]
		tests, err := extractTests(directory)
		if err != nil {
			fmt.Println("Error extracting tests:", err)
			os.Exit(1)
		}

		checkTestsInPipeline(yamlFile, tests)
	},
}

func extractTests(dir string) ([]Test, error) {
	var tests []Test
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), "_test.go") {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			re := regexp.MustCompile(`func (Test\w+)`)
			matches := re.FindAllSubmatch(content, -1)
			for _, match := range matches {
				tests = append(tests, Test{
					Name: string(match[1]),
					Path: path,
				})
			}
		}
		return nil
	})
	return tests, err
}

func checkTestsInPipeline(yamlFile string, tests []Test) {
	data, err := os.ReadFile(yamlFile)
	if err != nil {
		fmt.Printf("Error reading YAML file: %s\n", err)
		return
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		fmt.Printf("Error parsing YAML: %s\n", err)
		return
	}

	missingTests := []string{} // Track missing tests

	for _, test := range tests {
		found := false
		for _, item := range config.Tests {
			if item.Path == test.Path {
				if strings.Contains(item.TestCmd, "-test.run") {
					if matchTestNameInCmd(item.TestCmd, test.Name) {
						found = true
						break
					}
				} else {
					found = true
					break
				}
			}
		}
		if !found {
			missingTests = append(missingTests, fmt.Sprintf("ERROR: Test '%s' in file '%s' does not have CI configuration in '%s'", test.Name, test.Path, yamlFile))
		}
	}

	if len(missingTests) > 0 {
		for _, missing := range missingTests {
			fmt.Println(missing)
		}
		os.Exit(1) // Exit with a failure status
	}
}

// matchTestNameInCmd checks if the given test name matches the -test.run pattern in the command string.
func matchTestNameInCmd(cmd string, testName string) bool {
	testRunRegex := regexp.MustCompile(`-test\.run ([^\s]+)`)
	matches := testRunRegex.FindStringSubmatch(cmd)
	if len(matches) > 1 {
		// Extract the regex pattern used in the -test.run command
		pattern := matches[1]

		// Escape regex metacharacters in the testName before matching
		escapedTestName := regexp.QuoteMeta(testName)

		// Check if the escaped test name matches the extracted pattern
		return regexp.MustCompile(pattern).MatchString(escapedTestName)
	}
	return false
}

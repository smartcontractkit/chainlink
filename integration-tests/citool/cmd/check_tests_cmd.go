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
		excludedDirs := []string{"../../citool"}
		tests, err := extractTests(directory, excludedDirs)
		if err != nil {
			fmt.Println("Error extracting tests:", err)
			os.Exit(1)
		}

		checkTestsInPipeline(yamlFile, tests)
	},
}

// extractTests scans the given directory and subdirectories (except the excluded ones)
// for Go test files, extracts test function names, and returns a slice of Test.
func extractTests(dir string, excludeDirs []string) ([]Test, error) {
	var tests []Test

	// Resolve to absolute path
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}

	// filepath.WalkDir provides more control and is more efficient for skipping directories
	err = filepath.WalkDir(absDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Check if the current path is one of the excluded directories
		for _, exclude := range excludeDirs {
			absExclude, _ := filepath.Abs(exclude)
			if strings.HasPrefix(path, absExclude) {
				if d.IsDir() {
					return filepath.SkipDir // Skip this directory
				}
				return nil // Skip this file
			}
		}

		if !d.IsDir() && strings.HasSuffix(d.Name(), "_test.go") {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			re := regexp.MustCompile(`func (Test\w+)`)
			matches := re.FindAllSubmatch(content, -1)
			for _, match := range matches {
				funcName := string(match[1])
				if funcName == "TestMain" { // Skip "TestMain"
					continue
				}
				tests = append(tests, Test{
					Name: funcName,
					Path: mustExtractSubpath(path, "integration-tests"),
				})
			}
		}
		return nil
	})

	return tests, err
}

// ExtractSubpath extracts a specific subpath from a given full path.
// If the subpath is not found, it returns an error.
func mustExtractSubpath(fullPath, subPath string) string {
	index := strings.Index(fullPath, subPath)
	if index == -1 {
		panic("subpath not found in the provided full path")
	}
	return fullPath[index:]
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

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// main manages the process of combining coverage data for Chainlink nodes.
// It identifies "go-coverage" directories within a given root directory,
// merges their data into a "merged" directory for each test, and then
// calculates the overall coverage percentage.
func main() {
	// Check if an argument is provided
	if len(os.Args) < 2 {
		fmt.Println("Usage: show_coverage [root-directory]")
		os.Exit(1)
	}
	root := os.Args[1] // Use the first command-line argument as the root directory

	// Validate the root directory before proceeding
	if _, err := os.Stat(root); os.IsNotExist(err) {
		fmt.Printf("No coverage dir found: %s\n", root)
		os.Exit(0)
	}

	testDirs := make(map[string][]string)

	// Walk the file system from the root
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && filepath.Base(path) == "go-coverage" && path != root {
			// Assuming path structure /var/tmp/go-coverage/TestName/node_X/go-coverage
			testName := filepath.Dir(filepath.Dir(path)) // This should get the test name directory
			testDirs[testName] = append(testDirs[testName], path)
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error walking the path:", err)
		os.Exit(1)
	}

	// Iterate over the map and run the merge command for each test
	for test, dirs := range testDirs {
		testName := filepath.Base(test)

		// Ensure the merged directory exists
		mergedDir := filepath.Join(test, "merged")
		if err := os.MkdirAll(mergedDir, 0755); err != nil {
			fmt.Printf("Failed to create merged directory %s: %v\n", mergedDir, err)
			continue
		}

		// Merge the coverage data from all chainlink nodes
		dirInput := strings.Join(dirs, ",")
		mergeCmd := exec.Command("go", "tool", "covdata", "merge", "-o", mergedDir, "-i="+dirInput)
		mergeCmd.Dir = test
		fmt.Printf("Merging coverage for %s:\n%s\n", testName, mergeCmd.String())
		output, err := mergeCmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Error executing merge command for %s: %v, output: %s\n", test, err, output)
			os.Exit(1)
		}

		// Calculate coverage percentage in the merged directory
		coverageCmd := exec.Command("go", "tool", "covdata", "percent", "-i=merged")
		coverageCmd.Dir = test
		coverageOutput, err := coverageCmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Error calculating coverage for %s: %v, %s\n", test, err, string(coverageOutput))
			continue
		}
		fmt.Printf("Total coverage for %s:\n%s\n%s\n\n", testName, coverageCmd.String(), string(coverageOutput))
	}
}

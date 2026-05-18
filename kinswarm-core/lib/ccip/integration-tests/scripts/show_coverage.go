package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// main manages the process of combining coverage data for all tests
func main() {
	// Check if the user has provided an argument
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run script.go <searchPattern>")
		os.Exit(1)
	}

	// First argument after the program name is the search pattern
	searchPattern := os.Args[1]

	// Glob pattern to find all 'merged' directories in artifact folders
	dirs, err := filepath.Glob(searchPattern)
	if err != nil {
		fmt.Printf("Failed to find directories: %v\n", err)
		os.Exit(1)
	}

	if len(dirs) == 0 {
		fmt.Println("No directories found.")
		return
	}

	fmt.Printf("Found directories with test coverage data: %v\n", dirs)

	// Join the directory paths for input
	dirInput := strings.Join(dirs, ",")

	// Ensure the merged directory exists
	mergedDir := filepath.Join(".covdata", "merged")
	if err := os.MkdirAll(mergedDir, 0755); err != nil {
		fmt.Printf("Failed to create merged directory %s: %v\n", mergedDir, err)
		os.Exit(1)
	}

	// Merge the coverage data from all chainlink nodes
	mergeCmd := exec.Command("go", "tool", "covdata", "merge", "-o", mergedDir, "-i="+dirInput)
	fmt.Printf("Merging coverage data for all tests:\n%s\n", mergeCmd.String())
	output, err := mergeCmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error executing merge command: %v, output: %s\n", err, output)
		os.Exit(1)
	}

	// Calculate coverage percentage in the merged directory
	coverageCmd := exec.Command("go", "tool", "covdata", "percent", "-i=.")
	coverageCmd.Dir = mergedDir
	fmt.Printf("Calculate total coverage for on all tests: %s\n", coverageCmd.String())
	coverageOutput, err := coverageCmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error calculating coverage percentage: %v\n", err)
		os.Exit(1)
	}

	// Save the coverage percentage to a file
	filePath, err := filepath.Abs(filepath.Join(mergedDir, "percentage.txt"))
	if err != nil {
		fmt.Printf("Error obtaining absolute path: %s\n", err)
		os.Exit(1)
	}
	if err := os.WriteFile(filePath, coverageOutput, 0600); err != nil {
		fmt.Printf("Failed to write coverage percentage to file: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Total coverage for all tests saved to %s\n", filePath)

	fmt.Printf("Total coverage for all tests:\n%s\n", string(coverageOutput))
}

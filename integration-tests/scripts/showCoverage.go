package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	// Check if an argument is provided
	if len(os.Args) < 2 {
		fmt.Println("Usage: showCoverage [root-directory]")
		os.Exit(1)
	}
	root := os.Args[1] // Use the first command-line argument as the root directory

	// Validate the root directory before proceeding
	if _, err := os.Stat(root); os.IsNotExist(err) {
		fmt.Printf("The specified root directory does not exist: %s\n", root)
		os.Exit(1)
	}

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Skip the root directory itself
		if path == root {
			return nil
		}
		if info.IsDir() && info.Name() == "go-coverage" {
			fmt.Println("Found coverage profile for:", path)
			// Run the go tool covdata percent command
			cmd := exec.Command("go", "tool", "covdata", "percent", "-i=.")
			cmd.Dir = path // set working directory
			output, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Println("Error executing command:", err)
				return err
			}
			fmt.Println("Command:", cmd.String())
			fmt.Printf("Output:\n%s\n", string(output))
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error walking the path:", err)
	}
}

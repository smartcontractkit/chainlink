package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func LoadFile(contractABILocation string) (string, string, error) {
	data, err := ioutil.ReadFile(contractABILocation)
	//TODO: For now we take the contract name from the ABI file name given that the ABI doesn't contain the contract name. We need a more elegant/customizable solution

	contractName := getContractName(contractABILocation)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	contractABI := string(data)
	return contractName, contractABI, err
}

// CreateDirectories creates directories based on the provided paths.
// If a directory already exists and removeContents is true, it removes all contents inside the directory.
func CreateDirectories(dirs []string, removeContents bool) error {
	for _, dir := range dirs {
		// Check if the directory already exists
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			// Create the directory if it doesn't exist
			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %v", dir, err)
			}
			fmt.Printf("Created directory: %s\n", dir)
		} else if removeContents {
			// Remove the contents of the directory if it exists and removeContents is true
			if err := removeDirContents(dir); err != nil {
				return fmt.Errorf("failed to remove contents of directory %s: %v", dir, err)
			}
			fmt.Printf("Removed contents of directory: %s\n", dir)
		}
	}
	return nil
}

func ListFiles(dir string, extension string) ([]string, error) {
	var abiFiles []string

	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Check if the file has a .abi extension
		if !d.IsDir() && filepath.Ext(d.Name()) == "."+extension {
			abiFiles = append(abiFiles, path)
		}

		return nil
	})

	return abiFiles, err
}

// removeDirContents removes all contents of the specified directory without removing the directory itself.
func removeDirContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()

	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}

	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

func getContractName(filePath string) string {
	// Get the base name of the file (e.g., "pepe.abi")
	baseName := filepath.Base(filePath)
	// Remove the file extension to get the contract name (e.g., "pepe")
	contractName := strings.TrimSuffix(baseName, filepath.Ext(baseName))
	return contractName
}

func GenerateFile(outputDir string, fileName string, fileExt string, content []byte) error {
	// Ensure the output directory exists
	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	// Construct the full file path
	var filePath *string
	//TODO this is weird, for some reason it complains when I used & in filepath.Join(..)
	aux := filepath.Join(outputDir, fileName)
	filePath = &aux
	if fileExt != "" {
		aux := filepath.Join(outputDir, fileName+"."+fileExt)
		filePath = &aux
	}

	// Write the content to the file
	err = ioutil.WriteFile(*filePath, content, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}
	return nil
}

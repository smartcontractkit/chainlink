package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

type Input struct {
	ProductName           string   `json:"product_name"`
	TestNameRegex         string   `json:"test_name_regex"`
	TestFile              string   `json:"test_file"`
	EthImplementationName string   `json:"eth_implementation_name"`
	DockerImages          []string `json:"docker_images"`
}

type OutputEntry struct {
	ProductName           string `json:"product_name"`
	TestNameRegex         string `json:"test_name_regex"`
	TestFile              string `json:"test_file"`
	EthImplementationName string `json:"eth_implementation_name"`
	DockerImage           string `json:"docker_image"`
}

type Output struct {
	Entries []OutputEntry `json:"tests"`
}

const OutputFile = "compatibility_test_list.json"

func main() {
	if len(os.Args) < 5 {
		fmt.Println("Usage: go run main.go <product_name> <test_name_regex> <test_file> '<eth_implementation_name> <docker_images>")
		fmt.Println("Example: go run main.go 'OCR' 'TestOCR.*' './smoke/ocr_test.go' 'besu' 'hyperledger/besu:21.0.0,hyperledger/besu:22.0.0'")
		os.Exit(1)
	}

	dockerImagesArg := os.Args[5]
	dockerImages := strings.Split(dockerImagesArg, ",")

	input := Input{
		ProductName:           os.Args[1],
		TestNameRegex:         os.Args[2],
		TestFile:              os.Args[3],
		EthImplementationName: os.Args[4],
		DockerImages:          dockerImages,
	}

	var output Output
	var file *os.File
	if _, err := os.Stat(OutputFile); err == nil {
		file, err = os.OpenFile(OutputFile, os.O_RDWR, 0644)
		if err != nil {
			fmt.Printf("Error opening file: %v\n", err)
			return
		}
		defer func() { _ = file.Close() }()

		bytes, err := io.ReadAll(file)
		if err != nil {
			fmt.Printf("Error reading file: %v\n", err)
			return
		}

		if len(bytes) > 0 {
			if err := json.Unmarshal(bytes, &output); err != nil {
				fmt.Printf("Error unmarshalling JSON: %v\n", err)
				return
			}
		}
	} else {
		file, err = os.Create(OutputFile)
		if err != nil {
			fmt.Printf("Error creating file: %v\n", err)
			return
		}
	}
	defer func() { _ = file.Close() }()

	for _, image := range dockerImages {
		if !strings.Contains(image, ":") {
			fmt.Printf("Docker image format is invalid: %s", image)
			os.Exit(1)
		}
		output.Entries = append(output.Entries, OutputEntry{
			ProductName:           input.ProductName,
			TestNameRegex:         input.TestNameRegex,
			TestFile:              input.TestFile,
			EthImplementationName: input.EthImplementationName,
			DockerImage:           image,
		})
	}

	newOutput, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		fmt.Printf("Error marshalling JSON: %v\n", err)
		return
	}

	if _, err := file.WriteAt(newOutput, 0); err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
		return
	}

	fmt.Printf("%d compatibility test(s) for %s and %s added successfully!\n", len(dockerImages), input.ProductName, input.EthImplementationName)
}

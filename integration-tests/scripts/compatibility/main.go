package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

type Input struct {
	Product           string   `json:"product"`
	TestRegex         string   `json:"test_regex"`
	File              string   `json:"file"`
	EthImplementation string   `json:"eth_implementation"`
	DockerImages      []string `json:"docker_images"`
}

type OutputEntry struct {
	Product               string `json:"product"`
	TestRegex             string `json:"test_regex"`
	File                  string `json:"file"`
	EthImplementationName string `json:"eth_implementation"`
	DockerImage           string `json:"docker_image"`
}

type Output struct {
	Entries []OutputEntry `json:"tests"`
}

const OutputFile = "compatibility_test_list.json"

func main() {
	if len(os.Args) < 5 {
		fmt.Println("Usage: go run main.go <product> <test_regex> <file> '<eth_implementation> <docker_images>")
		fmt.Println("Example: go run main.go 'ocr' 'TestOCR.*' './smoke/ocr_test.go' 'besu' 'hyperledger/besu:21.0.0,hyperledger/besu:22.0.0'")
		os.Exit(1)
	}

	dockerImagesArg := os.Args[5]
	dockerImages := strings.Split(dockerImagesArg, ",")

	input := Input{
		Product:           os.Args[1],
		TestRegex:         os.Args[2],
		File:              os.Args[3],
		EthImplementation: os.Args[4],
		DockerImages:      dockerImages,
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
			Product:               input.Product,
			TestRegex:             input.TestRegex,
			File:                  input.File,
			EthImplementationName: input.EthImplementation,
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

	fmt.Printf("%d compatibility test(s) for %s and %s added successfully!\n", len(dockerImages), input.Product, input.EthImplementation)
}

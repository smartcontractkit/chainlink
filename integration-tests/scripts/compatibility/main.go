package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
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

const (
	OutputFile          = "compatibility_test_list.json"
	InsufficientArgsErr = `Usage: go run main.go <product> <test_regex> <file> '<eth_implementation> <docker_images>'
Example: go run main.go 'ocr' 'TestOCR.*' './smoke/ocr_test.go' 'besu' 'hyperledger/besu:21.0.0,hyperledger/besu:22.0.0'`
	EmptyParameterErr = "parameter '%s' cannot be empty"
)

// this script builds a JSON file with the compatibility tests to be run for a given product and Ethereum implementation
func main() {
	if len(os.Args) < 6 {
		panic(errors.New(InsufficientArgsErr))
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

	validateInput(input)

	var output Output
	var file *os.File
	if _, err := os.Stat(OutputFile); err == nil {
		file, err = os.OpenFile(OutputFile, os.O_RDWR, 0644)
		if err != nil {
			panic(fmt.Errorf("error opening file: %v\n", err))
		}
		defer func() { _ = file.Close() }()

		bytes, err := io.ReadAll(file)
		if err != nil {
			panic(fmt.Errorf("error reading file: %v\n", err))
		}

		if len(bytes) > 0 {
			if err := json.Unmarshal(bytes, &output); err != nil {
				panic(fmt.Errorf("error unmarshalling JSON: %v\n", err))
			}
		}
	} else {
		file, err = os.Create(OutputFile)
		if err != nil {
			panic(fmt.Errorf("error creating file: %v\n", err))
		}
	}
	defer func() { _ = file.Close() }()

	for _, image := range dockerImages {
		if !strings.Contains(image, ":") {
			panic(fmt.Errorf("docker image format is invalid: %s", image))
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
		panic(fmt.Errorf("Error marshalling JSON: %v\n", err))
	}

	if _, err := file.WriteAt(newOutput, 0); err != nil {
		panic(fmt.Errorf("Error writing to file: %v\n", err))
	}

	fmt.Printf("%d compatibility test(s) for %s and %s added successfully!\n", len(dockerImages), input.Product, input.EthImplementation)
}

func validateInput(input Input) {
	if input.Product == "" {
		panic(fmt.Errorf(EmptyParameterErr, "product"))
	}
	if input.TestRegex == "" {
		panic(fmt.Errorf(EmptyParameterErr, "test_regex"))
	} else {
		if _, err := regexp.Compile(input.TestRegex); err != nil {
			panic(fmt.Errorf("failed to compile regex: %v", err))
		}
	}
	if input.File == "" {
		panic(fmt.Errorf(EmptyParameterErr, "file"))
	}
	if input.EthImplementation == "" {
		panic(fmt.Errorf(EmptyParameterErr, "eth_implementation"))
	}
	if len(input.DockerImages) == 0 || (len(input.DockerImages) == 1 && input.DockerImages[0] == "") {
		panic(fmt.Errorf(EmptyParameterErr, "docker_images"))
	}
}

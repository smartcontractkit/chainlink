package src

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

const (
	defaultArtefactsDir        = "artefacts"
	defaultNodeSetsPath        = ".cache/node_sets.json"
	defaultKeylessNodeSetsPath = ".cache/keyless_node_sets.json"
	deployedContractsJSON      = "deployed_contracts.json"
	bootstrapSpecTemplate      = "bootstrap.toml"
	streamsTriggerSpecTemplate = "streams_trigger.toml"
	oracleSpecTemplate         = "oracle.toml"
)

func writeLines(lines []string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	wc := utils.NewDeferableWriteCloser(file)
	defer wc.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	if err := w.Flush(); err != nil {
		return err
	}
	return wc.Close()
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func mustReadJSON[T any](fileName string) (output T) {
	jsonFile, err := os.Open(fileName)
	if err != nil {
		panic(fmt.Sprintf("failed to open file at %s: %v", fileName, err))
	}
	defer jsonFile.Close()
	bytes, err := io.ReadAll(jsonFile)
	if err != nil {
		panic(fmt.Sprintf("failed to read file at %s: %v", fileName, err))
	}
	err = json.Unmarshal(bytes, &output)
	if err != nil {
		panic(fmt.Sprintf("failed to unmarshal data: %v", err))
	}
	return
}

func mustWriteJSON[T any](fileName string, data T) {
	jsonFile, err := os.Create(fileName)
	if err != nil {
		panic(fmt.Sprintf("failed to create file at %s: %v", fileName, err))
	}
	defer jsonFile.Close()
	encoder := json.NewEncoder(jsonFile)
	encoder.SetIndent("", " ")
	err = encoder.Encode(data)
	if err != nil {
		panic(fmt.Sprintf("failed to encode data: %v", err))
	}
}

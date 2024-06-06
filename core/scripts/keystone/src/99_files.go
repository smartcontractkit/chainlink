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
	artefactsDir          = "artefacts"
	deployedContractsJSON = "deployed_contracts.json"
	bootstrapSpecTemplate = "bootstrap.toml"
	cribOverrideTemplate  = "crib-overrides.yaml"
	oracleSpecTemplate    = "oracle.toml"
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

func mustParseJSON[T any](fileName string) (output T) {
	jsonFile, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()
	bytes, err := io.ReadAll(jsonFile)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(bytes, &output)
	if err != nil {
		panic(err)
	}
	return
}

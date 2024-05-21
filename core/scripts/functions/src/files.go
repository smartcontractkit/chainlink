package src

import (
	"bufio"
	"fmt"
	"os"

	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

const (
	configFile            = "config.yaml"
	templatesDir          = "templates"
	artefactsDir          = "artefacts"
	ocr2ConfigJson        = "FunctionsOracleConfig.json"
	ocr2PublicKeysJSON    = "OCR2PublicKeys.json"
	bootstrapSpecTemplate = "bootstrap.toml"
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

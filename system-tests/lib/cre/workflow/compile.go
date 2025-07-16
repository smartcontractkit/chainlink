package workflow

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/andybalholm/brotli"
	"github.com/pkg/errors"
)

func CompileWorkflow(workflowPath, workflowName string) (string, error) {
	workflowWasmPath := workflowName + ".wasm"

	goModTidyCmd := exec.Command("go", "mod", "tidy")
	goModTidyCmd.Dir = filepath.Dir(workflowPath)
	if err := goModTidyCmd.Run(); err != nil {
		return "", errors.Wrap(err, "failed to run go mod tidy")
	}

	buffer := bytes.Buffer{}
	compileCmd := exec.Command("go", "build", "-o", workflowWasmPath, filepath.Base(workflowPath)) // #nosec G204 -- we control the value of the cmd so the lint/sec error is a false positive
	compileCmd.Dir = filepath.Dir(workflowPath)
	compileCmd.Env = append(os.Environ(), "GOOS=wasip1", "GOARCH=wasm")
	compileCmd.Stdout = &buffer
	compileCmd.Stderr = &buffer
	if err := compileCmd.Run(); err != nil {
		fmt.Fprint(os.Stderr, buffer.String())
		return "", errors.Wrap(err, "failed to compile workflow")
	}

	workflowWasmAbsPath, workflowWasmAbsPathErr := filepath.Abs(filepath.Join(filepath.Dir(workflowPath), workflowWasmPath))
	if workflowWasmAbsPathErr != nil {
		return "", errors.Wrap(workflowWasmAbsPathErr, "failed to get absolute path of the workflow WASM file")
	}

	compressedWorkflowWasmPath, compressedWorkflowWasmPathErr := compressWorkflow(workflowWasmAbsPath)
	if compressedWorkflowWasmPathErr != nil {
		return "", errors.Wrap(compressedWorkflowWasmPathErr, "failed to compress workflow")
	}

	defer func() {
		_ = os.Remove(workflowWasmAbsPath)
	}()

	return compressedWorkflowWasmPath, nil
}

func compressWorkflow(workflowWasmPath string) (string, error) {
	baseName := strings.TrimSuffix(workflowWasmPath, filepath.Ext(workflowWasmPath))
	outputFile := baseName + ".br.b64"

	input, inputErr := os.ReadFile(workflowWasmPath)
	if inputErr != nil {
		return "", errors.Wrap(inputErr, "failed to read workflow WASM file")
	}

	var compressed bytes.Buffer
	brotliWriter := brotli.NewWriter(&compressed)

	if _, writeErr := brotliWriter.Write(input); writeErr != nil {
		return "", errors.Wrap(writeErr, "failed to compress workflow WASM file")
	}
	brotliWriter.Close()

	outputData := []byte(base64.StdEncoding.EncodeToString(compressed.Bytes()))

	if err := os.WriteFile(outputFile, outputData, 0600); err != nil {
		return "", errors.Wrap(err, "failed to write output file")
	}

	outputFileAbsPath, outputFileAbsPathErr := filepath.Abs(outputFile)
	if outputFileAbsPathErr != nil {
		return "", errors.Wrap(outputFileAbsPathErr, "failed to get absolute path of the output file")
	}

	return outputFileAbsPath, nil
}

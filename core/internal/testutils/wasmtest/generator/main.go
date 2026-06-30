package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/andybalholm/brotli"
)

func main() {
	pkgArg := flag.String("pkg", "", "The package to compile, e.g. test/simple/cmd")
	flag.Parse()

	if *pkgArg == "" {
		log.Fatal("-pkg is required")
	}

	pkgPath := "github.com/smartcontractkit/chainlink/v2/" + *pkgArg

	listCtx, listCancel := context.WithTimeout(context.Background(), time.Minute)
	listCmd := exec.CommandContext(listCtx, "go", "list", "-f", "{{.Dir}}", pkgPath)
	listCmd.Env = append(os.Environ(), "GOOS=wasip1", "GOARCH=wasm")
	out, err := listCmd.Output()
	listCancel()
	if err != nil {
		log.Fatalf("failed to find package dir for %s: %v\n%s", pkgPath, err, string(out))
	}
	pkgDir := strings.TrimSpace(string(out))

	fixtureName := "output.wasm.br"
	testdataDir := filepath.Join(pkgDir, "testdata")
	filePath := filepath.Join(testdataDir, fixtureName)

	// Always rebuild. Correctness is enforced by CI's `make generate` + dirty-tree check:
	// any drift in source, transitive dependencies, or toolchain produces different bytes
	// here, which CI detects. No staleness heuristic to get wrong.
	log.Printf("Building WASM fixture: %s", fixtureName)

	binary, err := buildBinary(pkgPath)
	if err != nil {
		log.Fatalf("build failed: %v", err)
	}

	if err := os.MkdirAll(testdataDir, 0755); err != nil {
		log.Fatalf("failed to create testdata dir: %v", err)
	}

	tmpFile := filePath + ".tmp"
	if err := os.WriteFile(tmpFile, binary, 0600); err != nil {
		log.Fatalf("failed to write temp file: %v", err)
	}
	if err := os.Rename(tmpFile, filePath); err != nil {
		log.Fatalf("failed to rename file: %v", err)
	}

	log.Printf("Successfully generated %s", filePath)
}

func buildBinary(pkgPath string) ([]byte, error) {
	tmpDir, err := os.MkdirTemp("", "wasmtest-*")
	if err != nil {
		return nil, fmt.Errorf("create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	cmdCtx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	buildPath := filepath.Join(tmpDir, "output.wasm")
	// -trimpath strips machine-specific absolute paths from the binary so the committed
	// fixture is reproducible across machines (required for the CI dirty-tree check).
	cmd := exec.CommandContext(cmdCtx, "go", "build", "-trimpath", "-o", buildPath, pkgPath) // #nosec
	cmd.Env = append(os.Environ(), "GOOS=wasip1", "GOARCH=wasm")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("build failed: %s %w", string(output), err)
	}

	binary, err := os.ReadFile(buildPath)
	if err != nil {
		return nil, fmt.Errorf("read file failed: %w", err)
	}

	var b bytes.Buffer
	bwr := brotli.NewWriter(&b)
	if _, err = bwr.Write(binary); err != nil {
		return nil, err
	}
	if err = bwr.Close(); err != nil {
		return nil, err
	}
	binary = b.Bytes()

	return binary, nil
}

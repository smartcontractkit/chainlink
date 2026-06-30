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

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/wasmtest"
)

func main() {
	pkgArg := flag.String("pkg", "", "The package to compile, e.g. test/simple/cmd")
	compressArg := flag.Bool("compress", false, "Whether to brotli compress the output")
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

	hash, err := wasmtest.HashPackage(pkgDir)
	if err != nil {
		log.Fatalf("failed to hash package: %v", err)
	}

	cacheFile := fmt.Sprintf("output-%s.wasm", hash)
	if *compressArg {
		cacheFile += ".br"
	}
	testdataDir := filepath.Join(pkgDir, "testdata")
	filePath := filepath.Join(testdataDir, cacheFile)

	if _, statErr := os.Stat(filePath); statErr == nil {
		log.Printf("Fixture already up to date: %s", filePath)
		return
	}

	log.Printf("Building WASM fixture: %s", cacheFile)
	cleanupOldCaches(pkgDir, *compressArg)

	binary, err := buildBinary(pkgPath, *compressArg)
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

func cleanupOldCaches(pkgDir string, compress bool) {
	testdataDir := filepath.Join(pkgDir, "testdata")
	entries, err := os.ReadDir(testdataDir)
	if err != nil {
		return
	}
	suffix := ".wasm"
	if compress {
		suffix = ".wasm.br"
	}
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), "output-") && strings.HasSuffix(e.Name(), suffix) {
			_ = os.Remove(filepath.Join(testdataDir, e.Name()))
		}
	}
}

func buildBinary(pkgPath string, compress bool) ([]byte, error) {
	tmpDir, err := os.MkdirTemp("", "wasmtest-*")
	if err != nil {
		return nil, fmt.Errorf("create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	cmdCtx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	buildPath := filepath.Join(tmpDir, "output.wasm")
	cmd := exec.CommandContext(cmdCtx, "go", "build", "-o", buildPath, pkgPath) // #nosec
	cmd.Env = append(os.Environ(), "GOOS=wasip1", "GOARCH=wasm")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("build failed: %s %w", string(output), err)
	}

	binary, err := os.ReadFile(buildPath)
	if err != nil {
		return nil, fmt.Errorf("read file failed: %w", err)
	}

	if compress {
		var b bytes.Buffer
		bwr := brotli.NewWriter(&b)
		if _, err = bwr.Write(binary); err != nil {
			return nil, err
		}
		if err = bwr.Close(); err != nil {
			return nil, err
		}
		binary = b.Bytes()
	}
	return binary, nil
}

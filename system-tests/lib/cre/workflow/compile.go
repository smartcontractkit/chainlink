package workflow

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/andybalholm/brotli"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

type Language = string

const (
	LanguageGo Language = "go"
	LanguageTS Language = "typescript"
)

// Set CRE_TEST_COMPILE_CACHE_DISABLED=true to bypass the in-process compile cache.
const compileCacheDisabledEnv = "CRE_TEST_COMPILE_CACHE_DISABLED"

// CompileWorkflow compiles a workflow from a file path (absolute or relative) and returns the path to the compiled workflow.
// workflowFilePath is the path to the workflow file.
// workflowName is the name of the workflow.
// It will return the path to the compiled workflow.
// It will return an error if the workflow name is less than 10 characters long.
// It will return an error if the workflow file path is not a valid file path.
func CompileWorkflow(ctx context.Context, workflowFilePath, workflowName string) (string, error) {
	return CompileWorkflowToDir(ctx, workflowFilePath, workflowName, "")
}

// CompileWorkflowToDir compiles a workflow and stores build artifacts in outputDir.
// If outputDir is empty, a temporary directory is created automatically.
// Artifacts are cached in-process by a hash of the workflow directory, so
// identical sources compile once per test run.
func CompileWorkflowToDir(ctx context.Context, workflowFilePath, workflowName, outputDir string) (string, error) {
	if len(workflowName) < 10 {
		return "", errors.New("workflow name must be at least 10 characters long")
	}
	if outputDir == "" {
		var err error
		outputDir, err = os.MkdirTemp("", "cre-workflow-build-*")
		if err != nil {
			return "", errors.Wrap(err, "failed to create temporary workflow build dir")
		}
	}
	if mkErr := os.MkdirAll(outputDir, 0o755); mkErr != nil {
		return "", errors.Wrap(mkErr, "failed to prepare workflow build dir")
	}

	if compileCacheEnabled() {
		key, keyErr := workflowSourceCacheKey(workflowFilePath)
		if keyErr == nil {
			return compileWorkflowToDirCached(ctx, key, workflowFilePath, workflowName, outputDir)
		}
		framework.L.Warn().Err(keyErr).Str("workflow_file", workflowFilePath).Msg("failed to hash workflow sources; compiling without the compile cache")
	}

	return compileWorkflowToDirUncached(ctx, workflowFilePath, workflowName, outputDir)
}

type compileCacheEntry struct {
	done     chan struct{}
	artifact []byte
	err      error
}

var (
	compileCacheMu sync.Mutex
	compileCache   = map[string]*compileCacheEntry{}
)

func compileCacheEnabled() bool {
	return strings.TrimSpace(strings.ToLower(os.Getenv(compileCacheDisabledEnv))) != "true"
}

// getOrCreateCompileCacheEntry returns the entry for key and whether the
// caller is the leader that must compile it.
func getOrCreateCompileCacheEntry(key string) (entry *compileCacheEntry, leader bool) {
	compileCacheMu.Lock()
	defer compileCacheMu.Unlock()
	if existing, ok := compileCache[key]; ok {
		return existing, false
	}
	entry = &compileCacheEntry{done: make(chan struct{})}
	compileCache[key] = entry
	return entry, true
}

// finalizeCompileCacheEntry publishes the result to waiters. Failed entries
// are evicted so later callers retry instead of reusing the failure.
func finalizeCompileCacheEntry(key string, entry *compileCacheEntry, artifact []byte, err error) {
	entry.artifact = artifact
	entry.err = err
	compileCacheMu.Lock()
	defer compileCacheMu.Unlock()
	if err != nil {
		delete(compileCache, key)
	}
	close(entry.done)
}

func compileWorkflowToDirCached(ctx context.Context, cacheKey, workflowFilePath, workflowName, outputDir string) (string, error) {
	for {
		entry, leader := getOrCreateCompileCacheEntry(cacheKey)
		if leader {
			return compileAsCompileCacheLeader(ctx, cacheKey, entry, workflowFilePath, workflowName, outputDir)
		}

		if err := waitForCompileCacheEntry(ctx, entry); err != nil {
			// The failed entry was evicted, so retrying makes this caller the
			// leader (covers a previous leader whose context was cancelled).
			framework.L.Warn().Str("workflow_name", workflowName).Err(err).Msg("in-flight workflow compile failed; retrying as the compile leader")
			continue
		}

		framework.L.Info().Str("workflow_name", workflowName).Str("workflow_file", workflowFilePath).Msg("workflow compile cache hit; reusing compiled workflow artifact")
		return writeCompressedArtifact(outputDir, workflowName, entry.artifact)
	}
}

func compileAsCompileCacheLeader(ctx context.Context, cacheKey string, entry *compileCacheEntry, workflowFilePath, workflowName, outputDir string) (string, error) {
	framework.L.Info().Str("workflow_name", workflowName).Str("workflow_file", workflowFilePath).Msg("workflow compile cache miss; compiling workflow")

	artifactPath, compileErr := compileWorkflowToDirUncached(ctx, workflowFilePath, workflowName, outputDir)

	var artifact []byte
	cacheErr := compileErr
	if cacheErr == nil {
		b, readErr := os.ReadFile(artifactPath)
		if readErr != nil {
			// Artifact is still valid for this caller; only caching is skipped.
			cacheErr = errors.Wrap(readErr, "failed to read compiled workflow artifact for the compile cache")
		} else {
			artifact = b
		}
	}

	finalizeCompileCacheEntry(cacheKey, entry, artifact, cacheErr)
	return artifactPath, compileErr
}

// waitForCompileCacheEntry blocks until the leader finishes or ctx ends.
func waitForCompileCacheEntry(ctx context.Context, entry *compileCacheEntry) error {
	select {
	case <-entry.done:
		if entry.err != nil {
			return errors.Wrap(entry.err, "in-flight workflow compile failed")
		}
		if len(entry.artifact) == 0 {
			return errors.New("in-flight workflow compile produced no artifact")
		}
		return nil
	case <-ctx.Done():
		return errors.Wrap(ctx.Err(), "cancelled while waiting for an in-flight workflow compile")
	}
}

// workflowSourceCacheKey hashes all build inputs in the workflow's directory
// (sources, go.mod/go.sum, package.json, ...), excluding build outputs.
func workflowSourceCacheKey(workflowFilePath string) (string, error) {
	dir := filepath.Dir(workflowFilePath)

	var files []string
	walkErr := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if d.Name() == ".git" || d.Name() == "node_modules" {
				return filepath.SkipDir
			}
			return nil
		}
		if strings.HasSuffix(d.Name(), ".wasm") || strings.HasSuffix(d.Name(), ".br.b64") {
			return nil
		}
		files = append(files, path)
		return nil
	})
	if walkErr != nil {
		return "", errors.Wrapf(walkErr, "failed to walk workflow source dir %s", dir)
	}

	sort.Strings(files)
	hasher := sha256.New()
	for _, path := range files {
		rel, relErr := filepath.Rel(dir, path)
		if relErr != nil {
			return "", errors.Wrap(relErr, "failed to relativize workflow source path")
		}
		content, readErr := os.ReadFile(path)
		if readErr != nil {
			return "", errors.Wrapf(readErr, "failed to read workflow source file %s", path)
		}
		hasher.Write([]byte(rel))
		hasher.Write([]byte{0})
		hasher.Write(content)
		hasher.Write([]byte{0})
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

// writeCompressedArtifact writes a cached artifact to
// outputDir/<workflowName>.br.b64, matching compressWorkflow's naming.
func writeCompressedArtifact(outputDir, workflowName string, artifact []byte) (string, error) {
	if mkErr := os.MkdirAll(outputDir, 0o755); mkErr != nil {
		return "", errors.Wrap(mkErr, "failed to prepare workflow build dir")
	}
	outputFile := filepath.Join(outputDir, workflowName+".br.b64")
	if writeErr := os.WriteFile(outputFile, artifact, 0o644); writeErr != nil { //nolint:gosec // G306: we want it to be readable by everyone, same as compressWorkflow
		return "", errors.Wrap(writeErr, "failed to write cached workflow artifact")
	}
	return outputFile, nil
}

// compileWorkflowToDirUncached compiles the workflow and returns the path to
// the compressed artifact in outputDir.
func compileWorkflowToDirUncached(ctx context.Context, workflowFilePath, workflowName, outputDir string) (string, error) {
	language, lErr := delectLanguage(workflowFilePath)
	if lErr != nil {
		return "", errors.Wrap(lErr, "failed to detect workflow language")
	}

	var workflowWasmAbsPath string
	var err error
	switch language {
	case LanguageGo:
		workflowWasmAbsPath, err = compileGoWorkflow(ctx, workflowFilePath, workflowName, outputDir)
	case LanguageTS:
		workflowWasmAbsPath, err = compileTSWorkflow(ctx, workflowFilePath, workflowName, outputDir)
	default:
		return "", fmt.Errorf("unsupported workflow language: %s", language)
	}

	if err != nil {
		return "", fmt.Errorf("failed to compile %s workflow: %w", language, err)
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

func delectLanguage(workflowFilePath string) (Language, error) {
	ext := strings.ToLower(filepath.Ext(workflowFilePath))
	switch ext {
	case ".ts", ".tsx":
		return LanguageTS, nil
	case ".go":
		return LanguageGo, nil
	default:
		return "", fmt.Errorf("unsupported workflow file extension: %s", ext)
	}
}

func compileTSWorkflow(ctx context.Context, workflowFilePath, workflowName, outputDir string) (string, error) {
	workflowWasmPath := filepath.Join(outputDir, workflowName+".wasm")

	compileCmd := exec.CommandContext(ctx, "bun", "cre-compile", workflowFilePath, workflowWasmPath) // #nosec G204 -- we control the value of the cmd so the lint/sec error is a false positive
	if output, err := compileCmd.CombinedOutput(); err != nil {
		fmt.Fprint(os.Stderr, string(output))
		return "", errors.Wrap(err, "failed to compile workflow")
	}

	workflowWasmAbsPath, workflowWasmAbsPathErr := filepath.Abs(workflowWasmPath)
	if workflowWasmAbsPathErr != nil {
		return "", errors.Wrap(workflowWasmAbsPathErr, "failed to get absolute path of the workflow WASM file")
	}

	return workflowWasmAbsPath, nil
}

func compileGoWorkflow(ctx context.Context, workflowFilePath, workflowName, outputDir string) (string, error) {
	workflowWasmPath := filepath.Join(outputDir, workflowName+".wasm")

	goModTidyCmd := exec.CommandContext(ctx, "go", "mod", "tidy")
	goModTidyCmd.Dir = filepath.Dir(workflowFilePath)
	if output, err := goModTidyCmd.CombinedOutput(); err != nil {
		return "", errors.Wrapf(err, "failed to run go mod tidy: %s", string(output))
	}

	compileCmd := exec.CommandContext(ctx, "go", "build", "-o", workflowWasmPath, filepath.Base(workflowFilePath)) // #nosec G204 -- we control the value of the cmd so the lint/sec error is a false positive
	compileCmd.Dir = filepath.Dir(workflowFilePath)
	compileCmd.Env = append(os.Environ(), "CGO_ENABLED=0", "GOOS=wasip1", "GOARCH=wasm")
	if output, err := compileCmd.CombinedOutput(); err != nil {
		fmt.Fprint(os.Stderr, string(output))
		return "", errors.Wrap(err, "failed to compile workflow")
	}

	workflowWasmAbsPath, workflowWasmAbsPathErr := filepath.Abs(workflowWasmPath)
	if workflowWasmAbsPathErr != nil {
		return "", errors.Wrap(workflowWasmAbsPathErr, "failed to get absolute path of the workflow WASM file")
	}

	return workflowWasmAbsPath, nil
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

	// remove the file if it already exists
	_, statErr := os.Stat(outputFile)
	if statErr == nil {
		if err := os.Remove(outputFile); err != nil {
			return "", errors.Wrap(err, "failed to remove existing output file")
		}
	}

	if err := os.WriteFile(outputFile, outputData, 0644); err != nil { //nolint:gosec // G306: we want it to be readable by everyone
		return "", errors.Wrap(err, "failed to write compressed output file")
	}

	outputFileAbsPath, outputFileAbsPathErr := filepath.Abs(outputFile)
	if outputFileAbsPathErr != nil {
		return "", errors.Wrap(outputFileAbsPathErr, "failed to get absolute path of the output file")
	}

	return outputFileAbsPath, nil
}

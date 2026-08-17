package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"regexp"
	"sort"
)

type MatrixEntry struct {
	TestName string `json:"test_name"`
	TestID   string `json:"test_id"`
	RunsOn   string `json:"runs_on"`
}

func main() {
	if err := run(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("ci-matrix", flag.ContinueOnError)
	fs.SetOutput(stderr)

	dir := fs.String("dir", "system-tests/tests/smoke/cre", "target directory to scan")
	pattern := fs.String("pattern", `^TestCRE_.*_E2E$`, "regex pattern matching E2E test functions")
	runID := fs.String("run-id", "0", "github.run_id")
	attempt := fs.String("attempt", "1", "github.run_attempt")
	runner := fs.String("runner", "cpu=16/ram=64/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs", "runs-on specification")
	githubOutput := fs.Bool("github-output", false, "write output in GITHUB_OUTPUT format")

	if err := fs.Parse(args); err != nil {
		return err
	}

	testNames, err := ScanDir(*dir, *pattern)
	if err != nil {
		return err
	}

	entries := BuildMatrix(testNames, *runID, *attempt, *runner)
	jsonBytes, err := json.Marshal(entries)
	if err != nil {
		return fmt.Errorf("failed to marshal matrix: %w", err)
	}

	if *githubOutput {
		outputPath := os.Getenv("GITHUB_OUTPUT")
		if outputPath != "" {
			f, openErr := os.OpenFile(outputPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
			if openErr != nil {
				return fmt.Errorf("failed to open GITHUB_OUTPUT: %w", openErr)
			}
			defer f.Close()
			if _, writeErr := fmt.Fprintf(f, "matrix=%s\n", string(jsonBytes)); writeErr != nil {
				return fmt.Errorf("failed to write to GITHUB_OUTPUT: %w", writeErr)
			}
		}
		if _, writeErr := fmt.Fprintf(stdout, "matrix=%s\n", string(jsonBytes)); writeErr != nil {
			return writeErr
		}
	} else {
		if _, writeErr := fmt.Fprintln(stdout, string(jsonBytes)); writeErr != nil {
			return writeErr
		}
	}

	return nil
}

// ScanDir scans all Go files in dir using AST parser and returns sorted test function names matching pattern.
func ScanDir(dir, pattern string) ([]string, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid pattern %q: %w", pattern, err)
	}

	stat, statErr := os.Stat(dir)
	if statErr != nil {
		return nil, fmt.Errorf("directory error %q: %w", dir, statErr)
	}
	if !stat.IsDir() {
		return nil, fmt.Errorf("%q is not a directory", dir)
	}

	fset := token.NewFileSet()
	pkgs, parseErr := parser.ParseDir(fset, dir, nil, 0)
	if parseErr != nil {
		return nil, fmt.Errorf("failed to parse directory %q: %w", dir, parseErr)
	}

	seen := make(map[string]struct{})
	var testNames []string
	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			for _, decl := range file.Decls {
				fn, ok := decl.(*ast.FuncDecl)
				if !ok {
					continue
				}
				name := fn.Name.Name
				if re.MatchString(name) {
					if _, exists := seen[name]; !exists {
						seen[name] = struct{}{}
						testNames = append(testNames, name)
					}
				}
			}
		}
	}

	if len(testNames) == 0 {
		return nil, errors.New("no matching test functions found")
	}

	sort.Strings(testNames)
	return testNames, nil
}

// BuildMatrix constructs MatrixEntry slice with unique runs-on labels.
func BuildMatrix(testNames []string, runID, attempt, runner string) []MatrixEntry {
	entries := make([]MatrixEntry, 0, len(testNames))
	for i, name := range testNames {
		runsOnLabel := fmt.Sprintf("runs-on=%s-%d-%s/%s", runID, i, attempt, runner)
		entries = append(entries, MatrixEntry{
			TestName: name,
			TestID:   name,
			RunsOn:   runsOnLabel,
		})
	}
	return entries
}

package matrix

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// DiscoverOptions specifies filtering options when discovering test names.
type DiscoverOptions struct {
	// IncludedPatterns, when non-empty, restricts results to names matching at
	// least one pattern. IgnoredPatterns still applies afterwards.
	IncludedPatterns []string
	IgnoredPatterns  []string
}

// DiscoverGoTestNames scans a directory for Go test and example function names.
func DiscoverGoTestNames(dir string, opts DiscoverOptions) ([]string, error) {
	compilePatterns := func(patterns []string, kind string) ([]*regexp.Regexp, error) {
		regexes := make([]*regexp.Regexp, 0, len(patterns))
		for _, pattern := range patterns {
			re, err := regexp.Compile(pattern)
			if err != nil {
				return nil, fmt.Errorf("invalid %s regex %q: %w", kind, pattern, err)
			}
			regexes = append(regexes, re)
		}
		return regexes, nil
	}

	includedRegexes, err := compilePatterns(opts.IncludedPatterns, "include")
	if err != nil {
		return nil, err
	}

	ignoredRegexes, err := compilePatterns(opts.IgnoredPatterns, "ignore")
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read test directory %s: %w", dir, err)
	}

	fset := token.NewFileSet()
	testNames := make(map[string]struct{})

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), "_test.go") {
			continue
		}

		filePath := filepath.Join(dir, entry.Name())
		node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
		if err != nil {
			return nil, fmt.Errorf("failed to parse file %s: %w", filePath, err)
		}

		for _, decl := range node.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Recv != nil {
				continue
			}

			name := fn.Name.Name
			if !strings.HasPrefix(name, "Test") && !strings.HasPrefix(name, "Example") {
				continue
			}

			if !matchesAny(includedRegexes, name) {
				continue
			}

			ignored := false
			for _, re := range ignoredRegexes {
				if re.MatchString(name) {
					ignored = true
					break
				}
			}
			if !ignored {
				testNames[name] = struct{}{}
			}
		}
	}

	var result []string
	for name := range testNames {
		result = append(result, name)
	}
	sort.Strings(result)

	return result, nil
}

// matchesAny reports whether name matches at least one pattern. An empty
// pattern list matches everything.
func matchesAny(regexes []*regexp.Regexp, name string) bool {
	if len(regexes) == 0 {
		return true
	}
	for _, re := range regexes {
		if re.MatchString(name) {
			return true
		}
	}
	return false
}

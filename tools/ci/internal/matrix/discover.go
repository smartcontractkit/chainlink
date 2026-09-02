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
	IgnoredPatterns []string
}

// DiscoverGoTestNames scans a directory for Go test and example function names.
func DiscoverGoTestNames(dir string, opts DiscoverOptions) ([]string, error) {
	var ignoredRegexes []*regexp.Regexp
	for _, pattern := range opts.IgnoredPatterns {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("invalid ignore regex %q: %w", pattern, err)
		}
		ignoredRegexes = append(ignoredRegexes, re)
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

// Package manifest loads dev-tool version entries from .tool-versions and go-tools.txt.
package manifest

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Entry is a single name/version pair from a manifest file.
type Entry struct {
	Name    string
	Version string
}

// Store holds parsed entries from both manifest files.
type Store struct {
	ToolVersionsPath string
	GoToolsPath      string
	runtimes         map[string]string
	goTools          map[string]string
}

// New loads both manifest files from the given paths.
func New(toolVersionsPath, goToolsPath string) (*Store, error) {
	runtimes, err := parseFile(toolVersionsPath)
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", toolVersionsPath, err)
	}
	goTools, err := parseFile(goToolsPath)
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", goToolsPath, err)
	}
	return &Store{
		ToolVersionsPath: toolVersionsPath,
		GoToolsPath:      goToolsPath,
		runtimes:         runtimes,
		goTools:          goTools,
	}, nil
}

func parseFile(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	entries := make(map[string]string)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		name, version, ok := parseLine(scanner.Text())
		if !ok {
			continue
		}
		entries[name] = version
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return entries, nil
}

func parseLine(line string) (name, version string, ok bool) {
	line = strings.TrimSpace(line)
	if line == "" || strings.HasPrefix(line, "#") {
		return "", "", false
	}
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return "", "", false
	}
	return fields[0], fields[1], true
}

// Lookup returns the version for key. Import paths (containing "/") resolve from
// go-tools.txt; short names resolve from .tool-versions.
func (s *Store) Lookup(key string) (string, error) {
	if strings.Contains(key, "/") {
		v, ok := s.goTools[key]
		if !ok {
			return "", fmt.Errorf("unknown go tool: %s", key)
		}
		return v, nil
	}
	v, ok := s.runtimes[key]
	if !ok {
		return "", fmt.Errorf("unknown tool: %s", key)
	}
	return v, nil
}

// List returns all entries from both manifests in file order (.tool-versions first).
func (s *Store) List() ([]Entry, error) {
	var out []Entry
	for _, path := range []string{s.ToolVersionsPath, s.GoToolsPath} {
		entries, err := listFile(path)
		if err != nil {
			return nil, err
		}
		out = append(out, entries...)
	}
	return out, nil
}

func listFile(path string) ([]Entry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var out []Entry
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		name, version, ok := parseLine(scanner.Text())
		if !ok {
			continue
		}
		out = append(out, Entry{Name: name, Version: version})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

// GoToolModules returns import paths from go-tools.txt.
func (s *Store) GoToolModules() []string {
	mods := make([]string, 0, len(s.goTools))
	for name := range s.goTools {
		mods = append(mods, name)
	}
	return mods
}

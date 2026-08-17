package testshard

import (
	"bufio"
	"errors"
	"fmt"
	"hash/fnv"
	"io"
	"strings"
)

// ShardForPackage computes the shard index for a given package path using FNV-1a.
func ShardForPackage(pkg string, shardCount int) int {
	hasher := fnv.New32a()
	_, _ = hasher.Write([]byte(pkg))
	return int(int64(hasher.Sum32()) % int64(shardCount))
}

// ReadPackages reads newline-delimited package names from r.
func ReadPackages(r io.Reader) ([]string, error) {
	scanner := bufio.NewScanner(r)
	packages := make([]string, 0)
	seen := make(map[string]struct{})

	for scanner.Scan() {
		pkg := strings.TrimSpace(scanner.Text())
		if pkg == "" {
			continue
		}

		if _, exists := seen[pkg]; exists {
			return nil, fmt.Errorf("duplicate package path %q", pkg)
		}
		seen[pkg] = struct{}{}
		packages = append(packages, pkg)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if len(packages) == 0 {
		return nil, errors.New("no package paths provided on stdin")
	}

	return packages, nil
}

// ValidateShardArgs validates shardCount and shardIndex bounds.
func ValidateShardArgs(shardCount, shardIndex int) error {
	if shardCount < 1 {
		return fmt.Errorf("invalid --shard-count %d: must be >= 1", shardCount)
	}
	if shardIndex < 0 || shardIndex >= shardCount {
		return fmt.Errorf("invalid --shard-index %d: must be in [0,%d)", shardIndex, shardCount)
	}
	return nil
}

// List filters packages from r that belong to shardIndex.
func List(r io.Reader, w io.Writer, shardCount, shardIndex int) error {
	packages, err := ReadPackages(r)
	if err != nil {
		return err
	}
	if err := ValidateShardArgs(shardCount, shardIndex); err != nil {
		return err
	}

	for _, pkg := range packages {
		if ShardForPackage(pkg, shardCount) == shardIndex {
			if _, err := fmt.Fprintln(w, pkg); err != nil {
				return err
			}
		}
	}
	return nil
}

// Verify checks that all packages from r are assigned to shards without overlap.
func Verify(r io.Reader, w io.Writer, shardCount int) error {
	packages, err := ReadPackages(r)
	if err != nil {
		return err
	}
	if shardCount < 1 {
		return fmt.Errorf("invalid --shard-count %d: must be >= 1", shardCount)
	}

	shardSizes := make([]int, shardCount)
	seen := make(map[string]int, len(packages))
	for _, pkg := range packages {
		idx := ShardForPackage(pkg, shardCount)
		shardSizes[idx]++
		seen[pkg]++
	}

	for _, pkg := range packages {
		if seen[pkg] != 1 {
			return fmt.Errorf("package %q assigned %d times", pkg, seen[pkg])
		}
	}

	if _, err := fmt.Fprintf(w, "verified %d packages across %d shards\n", len(packages), shardCount); err != nil {
		return err
	}
	for i, size := range shardSizes {
		if _, err := fmt.Fprintf(w, "shard %d: %d packages\n", i, size); err != nil {
			return err
		}
	}
	return nil
}

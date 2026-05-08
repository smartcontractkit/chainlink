package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"hash/fnv"
	"io"
	"os"
	"strings"
)

func main() {
	if err := run(os.Args[1:], os.Stdin, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string, stdin io.Reader, stdout io.Writer) error {
	if len(args) == 0 {
		return usageError("expected subcommand: list or verify")
	}

	switch args[0] {
	case "list":
		return runList(args[1:], stdin, stdout)
	case "verify":
		return runVerify(args[1:], stdin, stdout)
	case "-h", "--help", "help":
		printUsage(stdout)
		return nil
	default:
		return usageError("unknown subcommand %q", args[0])
	}
}

func runList(args []string, stdin io.Reader, stdout io.Writer) error {
	fs := flag.NewFlagSet("list", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	shardCount := fs.Int("shard-count", 1, "total number of shards")
	shardIndex := fs.Int("shard-index", 0, "zero-based shard index")

	if err := fs.Parse(args); err != nil {
		return usageError("%v", err)
	}
	if fs.NArg() != 0 {
		return usageError("list takes no positional arguments")
	}

	packages, err := readPackages(stdin)
	if err != nil {
		return err
	}
	if err := validateShardArgs(*shardCount, *shardIndex); err != nil {
		return err
	}

	for _, pkg := range packages {
		if shardForPackage(pkg, *shardCount) == *shardIndex {
			if _, err := fmt.Fprintln(stdout, pkg); err != nil {
				return err
			}
		}
	}

	return nil
}

func runVerify(args []string, stdin io.Reader, stdout io.Writer) error {
	fs := flag.NewFlagSet("verify", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	shardCount := fs.Int("shard-count", 1, "total number of shards")

	if err := fs.Parse(args); err != nil {
		return usageError("%v", err)
	}
	if fs.NArg() != 0 {
		return usageError("verify takes no positional arguments")
	}

	packages, err := readPackages(stdin)
	if err != nil {
		return err
	}
	if *shardCount < 1 {
		return fmt.Errorf("invalid --shard-count %d: must be >= 1", *shardCount)
	}

	shardSizes := make([]int, *shardCount)
	seen := make(map[string]int, len(packages))
	for _, pkg := range packages {
		shardIndex := shardForPackage(pkg, *shardCount)
		shardSizes[shardIndex]++
		seen[pkg]++
	}

	for _, pkg := range packages {
		if seen[pkg] != 1 {
			return fmt.Errorf("package %q assigned %d times", pkg, seen[pkg])
		}
	}

	if _, err := fmt.Fprintf(stdout, "verified %d packages across %d shards\n", len(packages), *shardCount); err != nil {
		return err
	}
	for shardIndex, size := range shardSizes {
		if _, err := fmt.Fprintf(stdout, "shard %d: %d packages\n", shardIndex, size); err != nil {
			return err
		}
	}

	return nil
}

func readPackages(r io.Reader) ([]string, error) {
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

func validateShardArgs(shardCount, shardIndex int) error {
	if shardCount < 1 {
		return fmt.Errorf("invalid --shard-count %d: must be >= 1", shardCount)
	}
	if shardIndex < 0 || shardIndex >= shardCount {
		return fmt.Errorf("invalid --shard-index %d: must be in [0,%d)", shardIndex, shardCount)
	}
	return nil
}

func shardForPackage(pkg string, shardCount int) int {
	hasher := fnv.New32a()
	_, _ = hasher.Write([]byte(pkg)) // hash.Hash.Write on fnv (Fowler-Noll-Vo) never returns an error
	return int(int64(hasher.Sum32()) % int64(shardCount))
}

func printUsage(w io.Writer) {
	fmt.Fprintln(w, "usage: ci-testshard <list|verify> [flags]")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "Commands:")
	fmt.Fprintln(w, "  list    read newline-delimited package paths from stdin and emit one shard")
	fmt.Fprintln(w, "  verify  read newline-delimited package paths from stdin and verify shard coverage")
}

func usageError(format string, args ...any) error {
	return fmt.Errorf(format, args...)
}

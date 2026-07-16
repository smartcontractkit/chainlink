package main

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"testing"
)

func TestReadPackagesRejectsDuplicatePaths(t *testing.T) {
	_, err := readPackages(strings.NewReader("pkg/a\npkg/a\n"))
	if err == nil || !strings.Contains(err.Error(), `duplicate package path "pkg/a"`) {
		t.Fatalf("expected duplicate package error, got %v", err)
	}
}

func TestReadPackagesRejectsEmptyInput(t *testing.T) {
	_, err := readPackages(strings.NewReader("\n\n"))
	if err == nil || !strings.Contains(err.Error(), "no package paths provided on stdin") {
		t.Fatalf("expected empty input error, got %v", err)
	}
}

func TestReadPackagesTrimsWhitespace(t *testing.T) {
	packages, err := readPackages(strings.NewReader("  pkg/a  \n\tpkg/b\t\n"))
	if err != nil {
		t.Fatalf("readPackages failed: %v", err)
	}
	if len(packages) != 2 || packages[0] != "pkg/a" || packages[1] != "pkg/b" {
		t.Fatalf("unexpected packages: %#v", packages)
	}
}

func TestReadPackagesIgnoresBlankLinesBetweenPackages(t *testing.T) {
	packages, err := readPackages(strings.NewReader("pkg/a\n\n   \n\t\npkg/b\n"))
	if err != nil {
		t.Fatalf("readPackages failed: %v", err)
	}
	if len(packages) != 2 || packages[0] != "pkg/a" || packages[1] != "pkg/b" {
		t.Fatalf("unexpected packages: %#v", packages)
	}
}

func TestReadPackagesRejectsDuplicatePathsAfterTrimming(t *testing.T) {
	_, err := readPackages(strings.NewReader("pkg/a\n  pkg/a  \n"))
	if err == nil || !strings.Contains(err.Error(), `duplicate package path "pkg/a"`) {
		t.Fatalf("expected duplicate package error after trimming, got %v", err)
	}
}

func TestListReturnsPartitionWithoutOverlap(t *testing.T) {
	input := "pkg/a\npkg/b\npkg/c\npkg/d\n"

	seen := make(map[string]struct{})
	for shardIndex := 0; shardIndex < 4; shardIndex++ {
		packages := runListForTest(t, input, 4, shardIndex)
		for _, pkg := range packages {
			if _, exists := seen[pkg]; exists {
				t.Fatalf("package %s appeared in multiple shards", pkg)
			}
			seen[pkg] = struct{}{}
		}
	}

	for _, pkg := range []string{"pkg/a", "pkg/b", "pkg/c", "pkg/d"} {
		if _, exists := seen[pkg]; !exists {
			t.Fatalf("package %s missing from shard union", pkg)
		}
	}
}

func TestListWithSingleShardReturnsEntireInputInOrder(t *testing.T) {
	input := "pkg/a\npkg/b\npkg/c\n"
	packages := runListForTest(t, input, 1, 0)
	want := []string{"pkg/a", "pkg/b", "pkg/c"}
	if len(packages) != len(want) {
		t.Fatalf("unexpected package count: got %d want %d (%v)", len(packages), len(want), packages)
	}
	for i := range want {
		if packages[i] != want[i] {
			t.Fatalf("unexpected package at %d: got %q want %q", i, packages[i], want[i])
		}
	}
}

func TestListProducesDeterministicOutput(t *testing.T) {
	input := "pkg/a\npkg/b\npkg/c\npkg/d\npkg/e\n"
	first := runListOutputForTest(t, input, 4, 2)
	second := runListOutputForTest(t, input, 4, 2)
	if first != second {
		t.Fatalf("list output changed between runs:\nfirst:\n%s\nsecond:\n%s", first, second)
	}
}

func TestListCanProduceEmptyShard(t *testing.T) {
	input := "pkg/a\npkg/b\n"
	foundEmpty := false
	for shardIndex := 0; shardIndex < 10; shardIndex++ {
		if output := runListOutputForTest(t, input, 10, shardIndex); output == "" {
			foundEmpty = true
			break
		}
	}
	if !foundEmpty {
		t.Fatal("expected at least one empty shard for 2 packages across 10 shards")
	}
}

func TestListAndVerifyAgreeOnPartition(t *testing.T) {
	inputPackages := []string{
		"pkg/a",
		"pkg/b",
		"pkg/c",
		"pkg/d",
		"pkg/e",
		"pkg/f",
	}
	input := strings.Join(inputPackages, "\n") + "\n"
	seen := make(map[string]struct{}, len(inputPackages))

	for shardIndex := 0; shardIndex < 4; shardIndex++ {
		for _, pkg := range runListForTest(t, input, 4, shardIndex) {
			if _, exists := seen[pkg]; exists {
				t.Fatalf("package %s appeared in multiple shards", pkg)
			}
			seen[pkg] = struct{}{}
		}
	}

	for _, pkg := range inputPackages {
		if _, exists := seen[pkg]; !exists {
			t.Fatalf("package %s missing from shard union", pkg)
		}
	}

	var stdout bytes.Buffer
	if err := run([]string{"verify", "--shard-count", "4"}, strings.NewReader(input), &stdout); err != nil {
		t.Fatalf("verify failed: %v", err)
	}
}

func TestVerifyAllowsEmptyShard(t *testing.T) {
	var stdout bytes.Buffer
	err := run([]string{"verify", "--shard-count", "10"}, strings.NewReader("pkg/a\npkg/b\n"), &stdout)
	if err != nil {
		t.Fatalf("verify failed: %v", err)
	}
	if !strings.Contains(stdout.String(), "verified 2 packages across 10 shards") {
		t.Fatalf("unexpected verify output: %q", stdout.String())
	}
}

func TestVerifyWithSingleShardCoversEntireInput(t *testing.T) {
	var stdout bytes.Buffer
	err := run([]string{"verify", "--shard-count", "1"}, strings.NewReader("pkg/a\npkg/b\npkg/c\n"), &stdout)
	if err != nil {
		t.Fatalf("verify failed: %v", err)
	}

	output := stdout.String()
	if !strings.Contains(output, "verified 3 packages across 1 shards") {
		t.Fatalf("unexpected verify summary: %q", output)
	}
	if !strings.Contains(output, "shard 0: 3 packages") {
		t.Fatalf("unexpected shard coverage: %q", output)
	}
}

func TestVerifyRejectsDuplicatePaths(t *testing.T) {
	var stdout bytes.Buffer
	err := run([]string{"verify", "--shard-count", "2"}, strings.NewReader("pkg/a\npkg/a\n"), &stdout)
	if err == nil || !strings.Contains(err.Error(), `duplicate package path "pkg/a"`) {
		t.Fatalf("expected duplicate package failure, got %v", err)
	}
}

func TestVerifyRejectsDuplicatePathsAmongOthers(t *testing.T) {
	var stdout bytes.Buffer
	err := run([]string{"verify", "--shard-count", "2"}, strings.NewReader("pkg/a\npkg/b\npkg/c\npkg/d\npkg/e\npkg/a\n"), &stdout)
	if err == nil || !strings.Contains(err.Error(), `duplicate package path "pkg/a"`) {
		t.Fatalf("expected duplicate package failure, got %v", err)
	}
}

func TestInvalidShardParamsFail(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{name: "zero-count", args: []string{"list", "--shard-count", "0", "--shard-index", "0"}},
		{name: "negative-index", args: []string{"list", "--shard-count", "2", "--shard-index", "-1"}},
		{name: "index-out-of-range", args: []string{"list", "--shard-count", "2", "--shard-index", "2"}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := run(tc.args, strings.NewReader("pkg/a\n"), &bytes.Buffer{})
			if err == nil {
				t.Fatal("expected error")
			}
		})
	}
}

func TestUnknownSubcommandFails(t *testing.T) {
	err := run([]string{"wat"}, strings.NewReader("pkg/a\n"), &bytes.Buffer{})
	if err == nil || !strings.Contains(err.Error(), `unknown subcommand "wat"`) {
		t.Fatalf("expected unknown subcommand error, got %v", err)
	}
}

func TestExtraPositionalArgsFail(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{name: "list", args: []string{"list", "--shard-count", "2", "--shard-index", "0", "extra"}, want: "list takes no positional arguments"},
		{name: "verify", args: []string{"verify", "--shard-count", "2", "extra"}, want: "verify takes no positional arguments"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := run(tc.args, strings.NewReader("pkg/a\n"), &bytes.Buffer{})
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("expected %q, got %v", tc.want, err)
			}
		})
	}
}

func TestLargePackageListParses(t *testing.T) {
	var builder strings.Builder
	for i := 0; i < 500; i++ {
		fmt.Fprintf(&builder, "pkg/%03d\n", i)
	}

	packages, err := readPackages(strings.NewReader(builder.String()))
	if err != nil {
		t.Fatalf("readPackages failed: %v", err)
	}
	if len(packages) != 500 {
		t.Fatalf("unexpected package count: got %d want 500", len(packages))
	}
	if packages[0] != "pkg/000" || packages[499] != "pkg/499" {
		t.Fatalf("unexpected package boundaries: first=%q last=%q", packages[0], packages[499])
	}
}

func runListForTest(t *testing.T, input string, shardCount, shardIndex int) []string {
	t.Helper()
	output := runListOutputForTest(t, input, shardCount, shardIndex)
	if output == "" {
		return nil
	}
	return strings.Fields(output)
}

func runListOutputForTest(t *testing.T, input string, shardCount, shardIndex int) string {
	t.Helper()
	var stdout bytes.Buffer
	if err := run(
		[]string{"list", "--shard-count", strconv.Itoa(shardCount), "--shard-index", strconv.Itoa(shardIndex)},
		strings.NewReader(input),
		&stdout,
	); err != nil {
		t.Fatalf("list failed for shard %d/%d: %v", shardIndex, shardCount, err)
	}
	return stdout.String()
}

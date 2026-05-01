// Package main launches the tools/test CLI.
//
//go:generate go run . --sync-skills
package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/smartcontractkit/chainlink/v2/tools/test/internal/cmd"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--sync-skills" {
		if err := syncSkills(); err != nil {
			fmt.Fprintf(os.Stderr, "error syncing skills: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Synced /tools/test/.agents/skills -> /tools/test/.claude/skills.")
		os.Exit(0)
	}

	cmd.Execute()
}

const (
	skillsSrc = ".agents/skills"
	skillsDst = ".claude/skills"
)

func syncSkills() error {
	base, err := moduleDir()
	if err != nil {
		return err
	}
	return syncSkillsFromBase(base)
}

func syncSkillsFromBase(base string) error {
	base, err := filepath.Abs(base)
	if err != nil {
		return fmt.Errorf("failed to resolve base directory: %w", err)
	}
	src := filepath.Join(base, skillsSrc)
	dst := filepath.Join(base, skillsDst)
	err = ensureWithinBase(base, dst)
	if err != nil {
		return err
	}
	info, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("failed to stat source directory %q: %w", src, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("source path %q is not a directory", src)
	}

	err = os.RemoveAll(dst)
	if err != nil {
		return fmt.Errorf("failed to remove destination directory %q: %w", dst, err)
	}

	err = filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(src, path)
		if err != nil {
			return fmt.Errorf("failed to calculate relative path for %q: %w", path, err)
		}
		targetPath := filepath.Join(dst, rel)

		if info.IsDir() {
			return os.MkdirAll(targetPath, info.Mode())
		}

		return copyFile(path, targetPath)
	})

	if err != nil {
		return fmt.Errorf("error syncing skills: %w", err)
	}
	return nil
}

func moduleDir() (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("failed to resolve module directory")
	}
	return filepath.Dir(file), nil
}

func ensureWithinBase(base, target string) error {
	rel, err := filepath.Rel(base, target)
	if err != nil {
		return fmt.Errorf("failed to verify destination path %q: %w", target, err)
	}
	if rel == "." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) || rel == ".." || filepath.IsAbs(rel) {
		return fmt.Errorf("destination path %q is outside base directory %q", target, base)
	}
	return nil
}

func copyFile(src, dst string) error {
	info, err := os.Stat(src)
	if err != nil {
		return err
	}
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, info.Mode())
	if err != nil {
		return err
	}
	_, err = io.Copy(out, in)
	if err != nil {
		out.Close()
		return err
	}
	if err := out.Close(); err != nil {
		return err
	}
	return os.Chmod(dst, info.Mode())
}

// Package main launches the tools/test CLI.
//
//go:generate go run . --sync-skills
package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

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
	// 1. Wipe destination to prevent stale skills
	os.RemoveAll(skillsDst)

	// 2. Walk source and copy
	err := filepath.Walk(skillsSrc, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Calculate relative path to mirror structure
		rel, _ := filepath.Rel(skillsSrc, path)
		targetPath := filepath.Join(skillsDst, rel)

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

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return err
}

package helpers

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const aptosCLIPathEnvVar = "APTOS_CLI_PATH"

func PreferHostAptosCLI() (string, error) {
	cliPath, err := firstWorkingAptosCLI(aptosCLICandidates())
	if err != nil {
		return "", err
	}

	cliDir := filepath.Dir(cliPath)
	pathDirs := filepath.SplitList(os.Getenv("PATH"))
	if len(pathDirs) > 0 && filepath.Clean(pathDirs[0]) == cliDir {
		return cliPath, nil
	}

	currentPath := os.Getenv("PATH")
	if currentPath != "" {
		return cliPath, os.Setenv("PATH", cliDir+string(os.PathListSeparator)+currentPath)
	}
	return cliPath, os.Setenv("PATH", cliDir)
}

func aptosCLICandidates() []string {
	var candidates []string
	seen := make(map[string]struct{})
	add := func(path string) {
		path = strings.TrimSpace(path)
		if path == "" {
			return
		}
		path = filepath.Clean(path)
		if _, ok := seen[path]; ok {
			return
		}
		seen[path] = struct{}{}
		candidates = append(candidates, path)
	}

	add(os.Getenv(aptosCLIPathEnvVar))
	for _, dir := range filepath.SplitList(os.Getenv("PATH")) {
		add(filepath.Join(dir, "aptos"))
	}

	add("/tmp/aptos-cli-7.8.0/aptos")
	add("/opt/homebrew/bin/aptos")
	add("/usr/local/bin/aptos")
	return candidates
}

func firstWorkingAptosCLI(candidates []string) (string, error) {
	problems := make([]string, 0, len(candidates))
	for _, candidate := range candidates {
		err := validateAptosCLI(candidate)
		if err == nil {
			return candidate, nil
		}

		problems = append(problems, fmt.Sprintf("%s (%v)", candidate, err))
	}
	return "", fmt.Errorf("failed to find a working Aptos CLI; set %s or install a valid aptos binary. Checked: %s", aptosCLIPathEnvVar, strings.Join(problems, ", "))
}

func validateAptosCLI(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if info.IsDir() {
		return errors.New("is a directory")
	}
	if info.Mode()&0o111 == 0 {
		return errors.New("is not executable")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	out, err := exec.CommandContext(ctx, path, "--version").CombinedOutput()
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if err != nil {
		return fmt.Errorf("version check failed: %s", strings.TrimSpace(string(out)))
	}

	versionOutput := strings.ToLower(strings.TrimSpace(string(out)))
	if !strings.Contains(versionOutput, "aptos") {
		return fmt.Errorf("unexpected version output: %q", strings.TrimSpace(string(out)))
	}
	return nil
}

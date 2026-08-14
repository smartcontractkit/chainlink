package target

import (
	"context"
	"os"
	"os/exec"
)

func PositiveCases(ctx context.Context) {
	// Custom WASM build calls that MUST be flagged:
	exec.Command("go", "build", "-o", "test.wasm", "main.go")            // want `Do not invoke 'go build' directly`
	exec.CommandContext(ctx, "go", "build", "-o", "test.wasm", "main.go") // want `Do not invoke 'go build' directly`

	// Manual WASM target envs that MUST be flagged:
	var env []string
	env = append(env, "GOOS=wasip1", "GOARCH=wasm") // want `Do not set GOARCH=wasm` `Do not set GOARCH=wasm`
	_ = env

	os.Setenv("GOARCH=wasm", "1") // want `Do not set GOARCH=wasm`
	os.Setenv("GOOS=wasip1", "1") // want `Do not set GOARCH=wasm`
}

func NegativeCases(ctx context.Context) {
	// Legitimate command invocations that must NOT be flagged (no false positives):
	_ = exec.Command("git", "status")
	_ = exec.CommandContext(ctx, "docker", "ps")
	_ = exec.Command("bun", "cre-compile", "app.ts")
	_ = exec.Command("go", "version")
	_ = exec.Command("go", "test", "./...")

	var env []string
	env = append(env, "CGO_ENABLED=0", "FOO=bar")
	_ = env

	_ = os.Setenv("FOO", "bar")
}

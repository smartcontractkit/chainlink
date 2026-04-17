package environment

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuildConfigDockerBuildArgs(t *testing.T) {
	t.Setenv("GITHUB_TOKEN", "test-token")

	cfg := BuildConfig{
		LocalImage:         "test-image:latest",
		Dockerfile:         "Dockerfile.test",
		DockerCtx:          ".",
		RequireGithubToken: true,
	}

	require.Equal(t, []string{
		"build",
		"-t", "test-image:latest",
		"-f", "Dockerfile.test",
		"--build-arg", "GITHUB_TOKEN=test-token",
		".",
	}, cfg.dockerBuildArgs())
}

func TestBuildConfigDockerBuildArgs_WithoutGithubToken(t *testing.T) {
	require.NoError(t, os.Unsetenv("GITHUB_TOKEN"))

	cfg := BuildConfig{
		LocalImage: "test-image:latest",
		Dockerfile: "Dockerfile.test",
		DockerCtx:  ".",
	}

	require.Equal(t, []string{
		"build",
		"-t", "test-image:latest",
		"-f", "Dockerfile.test",
		".",
	}, cfg.dockerBuildArgs())
}

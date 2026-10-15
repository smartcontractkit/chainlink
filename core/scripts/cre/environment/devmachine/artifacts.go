package devmachine

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

// Artifacts implements phase P3: build the remote agent and the two amd64
// images, caching the resulting tarballs under ~/.cache/cre/devmachine.
func (o *Orchestrator) Artifacts(ctx context.Context) error {
	o.rep.Phase(phaseTitle(phaseArtifacts))

	cache, err := CacheDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(cache, 0o755); err != nil {
		return errors.Wrap(err, "failed to create cache directory")
	}

	if err := o.buildAgent(ctx, cache); err != nil {
		return err
	}
	if err := o.buildImages(ctx, cache); err != nil {
		return err
	}
	o.cacheDir = cache
	return nil
}

func (o *Orchestrator) buildAgent(ctx context.Context, cache string) error {
	out := filepath.Join(cache, "remote-agent")
	sha, _ := o.git(ctx, o.cfg.Repos.Chainlink, "rev-parse", "HEAD")
	marker := filepath.Join(cache, "remote-agent.sha")
	if sha != "" {
		if prev, err := os.ReadFile(marker); err == nil && strings.TrimSpace(string(prev)) == sha {
			if _, err := os.Stat(out); err == nil {
				o.rep.Info("remote-agent up to date (sha %s), skipping build", short(sha))
				return nil
			}
		}
	}

	o.rep.Info("building remote-agent for linux/amd64")
	libDir := filepath.Join(o.cfg.Repos.Chainlink, "system-tests", "lib")
	if o.run.DryRun {
		_, _ = o.run.RunEnv(ctx, []string{"GOOS=linux", "GOARCH=amd64", "CGO_ENABLED=0"},
			"go", "build", "-o", out, "./cre/environment/remoteexec/agent/cmd/local-agent")
		return nil
	}
	runner := o.run.In(libDir)
	if err := runner.StreamEnv(ctx,
		[]string{"GOOS=linux", "GOARCH=amd64", "CGO_ENABLED=0"},
		"go", "build", "-o", out, "./cre/environment/remoteexec/agent/cmd/local-agent"); err != nil {
		return errors.Wrap(err, "failed to build remote-agent")
	}
	if sha != "" {
		_ = os.WriteFile(marker, []byte(sha), 0o600)
	}
	return nil
}

func (o *Orchestrator) buildImages(ctx context.Context, cache string) error {
	clTar := filepath.Join(cache, tarName(o.cfg.Images.Chainlink))
	if err := o.ensureImage(ctx, o.cfg.Images.Chainlink, clTar, func() error {
		return o.buildChainlinkImage(ctx)
	}); err != nil {
		return err
	}

	jdTar := filepath.Join(cache, tarName(o.cfg.Images.JobDistributor))
	return o.ensureImage(ctx, o.cfg.Images.JobDistributor, jdTar, func() error {
		return o.buildJobDistributorImage(ctx)
	})
}

// ensureImage makes sure a tarball exists for image, building it first if the
// local image is missing.
func (o *Orchestrator) ensureImage(ctx context.Context, image, tarPath string, build func() error) error {
	if !o.run.DryRun {
		if _, err := os.Stat(tarPath); err == nil {
			o.rep.Info("cached tarball present for %s, skipping build", image)
			return nil
		}
		if !o.imageExists(ctx, image) {
			if err := build(); err != nil {
				return err
			}
		}
	} else {
		if err := build(); err != nil {
			return err
		}
	}
	o.rep.Info("saving %s to %s", image, tarPath)
	_, err := o.run.Run(ctx, "docker", "save", image, "-o", tarPath)
	return errors.Wrapf(err, "failed to save image %s", image)
}

func (o *Orchestrator) imageExists(ctx context.Context, image string) bool {
	_, err := o.run.Run(ctx, "docker", "image", "inspect", image)
	return err == nil
}

func (o *Orchestrator) buildChainlinkImage(ctx context.Context) error {
	repo := o.cfg.Repos.Chainlink
	sha, _ := o.git(ctx, repo, "rev-parse", "HEAD")
	tag, _ := o.git(ctx, repo, "describe", "--always")
	o.rep.Info("building %s for linux/amd64 (this can take a while)", o.cfg.Images.Chainlink)
	runner := o.run.In(repo)
	return runner.Stream(ctx, "docker", "buildx", "build",
		"--platform", "linux/amd64",
		"--build-arg", "COMMIT_SHA="+sha,
		"--build-arg", "VERSION_TAG="+tag,
		"--build-arg", "CL_AUTO_DOCKER_TAG=develop",
		"-f", "core/chainlink.Dockerfile",
		"-t", o.cfg.Images.Chainlink,
		"--load",
		".")
}

func (o *Orchestrator) buildJobDistributorImage(ctx context.Context) error {
	repo := o.cfg.Repos.JobDistributor
	o.rep.Info("building %s for linux/amd64", o.cfg.Images.JobDistributor)
	runner := o.run.In(repo)

	if !o.run.DryRun {
		if _, err := runner.Run(ctx, "git", "checkout", o.cfg.Images.JDTag); err != nil {
			return errors.Wrapf(err, "failed to checkout %s in job-distributor", o.cfg.Images.JDTag)
		}
	}
	version := strings.TrimPrefix(o.cfg.Images.JDTag, "v")
	if err := runner.StreamEnv(ctx,
		[]string{"CGO_ENABLED=0", "GOOS=linux", "GOARCH=amd64"},
		"go", "build",
		"-ldflags", "-s -w -X github.com/smartcontractkit/job-distributor/build.Version="+version,
		"-o", "job-distributor",
		"./cmd/server"); err != nil {
		return errors.Wrap(err, "failed to build job-distributor server binary")
	}
	return runner.Stream(ctx, "docker", "buildx", "build",
		"--platform", "linux/amd64",
		"-t", o.cfg.Images.JobDistributor,
		"--load",
		".")
}

func (o *Orchestrator) git(ctx context.Context, dir string, args ...string) (string, error) {
	if o.run.DryRun {
		return "", nil
	}
	return o.run.In(dir).Run(ctx, "git", args...)
}

// tarName converts an image reference into a filesystem-friendly tar name.
func tarName(image string) string {
	replacer := strings.NewReplacer("/", "_", ":", "-")
	return replacer.Replace(image) + ".tar"
}

func short(sha string) string {
	if len(sha) > 12 {
		return sha[:12]
	}
	return sha
}

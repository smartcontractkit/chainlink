package environment

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"

	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
)

// LocalBuildConfig controls building the Chainlink node image and/or CRE
// capability plugins from the developer's LOCAL working tree instead of
// pulling/using the pinned upstream versions.
//
// It is driven by the `env start` flags:
//
//	--local-node                       build the node image from local source
//	--local-capabilities <names|all>   build these capabilities from local source
//	--capabilities-path <dir>          location of the local capabilities repo
//
// Node image: a minimal self-contained image (ubuntu + runtime libs + the
// locally-built chainlink binary) is built and every node spec is pointed at
// it. The node is cross-compiled on the host, which resolves any local
// `replace` directives (e.g. a local chainlink-common) that an in-Docker build
// cannot see.
//
// Capabilities: each named plugin is cross-compiled from the local capabilities
// repo and injected into every node via CTF's CapabilitiesBinaryPaths mechanism
// (copied into /usr/local/bin at container creation), overriding whatever the
// node image ships. Unlisted capabilities keep using the image's baked version.
type LocalBuildConfig struct {
	Node             bool
	Capabilities     []string // capability names, or a single "all"
	CapabilitiesPath string
	Platform         string // GOOS/GOARCH target, e.g. "linux/arm64"

	// DevImage is the tag used for the locally built node image.
	DevImage string
	// RepoRoot is the chainlink repo root (build context for the node binary).
	RepoRoot string
	// Toolchain env for cross-compiling the CGO node binary. Defaults are the
	// zig-based wrappers under $CLCROSS (see docs). Overridable via env.
	CC        string
	CXX       string
	LibDir    string // dir with libstdc++.so + static .a backfill
	OutputDir string // where built artifacts are staged (default: os.TempDir based)
}

// capBuildSpec maps a capability name to its module directory (relative to the
// capabilities repo), the output binary name (must match binary_name in
// capability_defaults.toml), and any extra build flags.
type capBuildSpec struct {
	dir    string
	binary string
	flags  []string
}

// knownCapabilities is the set of CRE capability plugins that can be built from
// the local capabilities repo. Keys accept both dash and underscore forms.
var knownCapabilities = map[string]capBuildSpec{
	"cron":         {dir: "cron", binary: "cron", flags: []string{"-tags", "timetzdata"}},
	"consensus":    {dir: "consensus", binary: "consensus"},
	"http-action":  {dir: "http_action", binary: "http_action"},
	"http-trigger": {dir: "http_trigger", binary: "http_trigger"},
	"evm":          {dir: "chain_capabilities/evm", binary: "evm"},
}

// allCapabilities is the default set built when "all" is requested.
var allCapabilities = []string{"cron", "consensus", "http-action", "http-trigger", "evm"}

const localBuildContainerCapDir = "/usr/local/bin"

func normalizeCapName(name string) string {
	return strings.ReplaceAll(strings.TrimSpace(name), "_", "-")
}

// resolveCapabilities expands "all" and normalizes/validates the requested names.
func (c *LocalBuildConfig) resolveCapabilities() ([]capBuildSpec, error) {
	names := c.Capabilities
	if len(names) == 1 && normalizeCapName(names[0]) == "all" {
		names = allCapabilities
	}
	seen := map[string]bool{}
	var specs []capBuildSpec
	for _, n := range names {
		key := normalizeCapName(n)
		if key == "" || seen[key] {
			continue
		}
		seen[key] = true
		spec, ok := knownCapabilities[key]
		if !ok {
			valid := make([]string, 0, len(knownCapabilities))
			for k := range knownCapabilities {
				valid = append(valid, k)
			}
			sort.Strings(valid)
			return nil, fmt.Errorf("unknown capability %q; supported: %s (or 'all')", n, strings.Join(valid, ", "))
		}
		specs = append(specs, spec)
	}
	return specs, nil
}

func (c *LocalBuildConfig) goEnv() []string {
	goos, goarch := "linux", "arm64"
	if c.Platform != "" {
		parts := strings.SplitN(c.Platform, "/", 2)
		if len(parts) == 2 {
			goos, goarch = parts[0], parts[1]
		}
	}
	return []string{"GOOS=" + goos, "GOARCH=" + goarch}
}

func (c *LocalBuildConfig) applyDefaults() {
	if c.DevImage == "" {
		c.DevImage = "cre-node:local"
	}
	if c.Platform == "" {
		// Match the host arch for the linux target (Docker Desktop on arm64
		// runs arm64 containers; amd64 hosts run amd64).
		c.Platform = "linux/" + runtime.GOARCH
	}
	clcross := os.Getenv("CLCROSS")
	if clcross == "" {
		clcross = filepath.Join(os.Getenv("HOME"), "clcross")
	}
	if c.CC == "" {
		c.CC = firstNonEmpty(os.Getenv("CRE_LOCAL_CC"), filepath.Join(clcross, "bin", "cc"))
	}
	if c.CXX == "" {
		c.CXX = firstNonEmpty(os.Getenv("CRE_LOCAL_CXX"), filepath.Join(clcross, "bin", "cxx"))
	}
	if c.LibDir == "" {
		c.LibDir = firstNonEmpty(os.Getenv("CRE_LOCAL_LIBDIR"), filepath.Join(clcross, "dyn"))
	}
	if c.CapabilitiesPath == "" {
		c.CapabilitiesPath = firstNonEmpty(
			os.Getenv("CRE_CAPABILITIES_PATH"),
			filepath.Join(os.Getenv("HOME"), "go", "src", "github.com", "smartcontractkit", "capabilities"),
		)
	}
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}

// ApplyLocalBuilds builds the requested local artifacts and rewrites the loaded
// config so the environment uses them. It must run after in.Load(...) and
// before Docker image checks / environment startup.
func ApplyLocalBuilds(ctx context.Context, in *envconfig.Config, cfg LocalBuildConfig) error {
	cfg.applyDefaults()

	// A local node image built minimally ships no capability plugins, so it is
	// only useful together with locally built capabilities. Default to all.
	if cfg.Node && len(cfg.Capabilities) == 0 {
		framework.L.Info().Msg("--local-node set without --local-capabilities; defaulting to all capabilities so the DON has its plugins")
		cfg.Capabilities = []string{"all"}
	}

	var capBinaries []string // host paths of built capability binaries
	if len(cfg.Capabilities) > 0 {
		specs, err := cfg.resolveCapabilities()
		if err != nil {
			return err
		}
		capBinaries, err = cfg.buildCapabilities(ctx, specs)
		if err != nil {
			return errors.Wrap(err, "failed to build local capabilities")
		}
	}

	if cfg.Node {
		if err := cfg.buildNodeImage(ctx); err != nil {
			return errors.Wrap(err, "failed to build local node image")
		}
	}

	// Rewrite every node spec in every nodeset.
	for _, nodeSet := range in.NodeSets {
		for _, spec := range nodeSet.NodeSpecs {
			if spec == nil || spec.Node == nil {
				continue
			}
			if cfg.Node {
				spec.Node.Image = cfg.DevImage
				// Mutually exclusive with a prebuilt image (clnode guard).
				spec.Node.DockerContext = ""
				spec.Node.DockerFilePath = ""
				spec.Node.PullImage = false
			}
			if len(capBinaries) > 0 {
				spec.Node.CapabilityContainerDir = localBuildContainerCapDir
				spec.Node.CapabilitiesBinaryPaths = appendUnique(spec.Node.CapabilitiesBinaryPaths, capBinaries...)
			}
		}
	}
	return nil
}

func appendUnique(dst []string, vals ...string) []string {
	seen := map[string]bool{}
	for _, d := range dst {
		seen[d] = true
	}
	for _, v := range vals {
		if !seen[v] {
			dst = append(dst, v)
			seen[v] = true
		}
	}
	return dst
}

// buildCapabilities cross-compiles each capability from the local repo and
// returns the host paths of the built binaries (named by binary_name).
func (c *LocalBuildConfig) buildCapabilities(ctx context.Context, specs []capBuildSpec) ([]string, error) {
	if _, err := os.Stat(c.CapabilitiesPath); err != nil {
		return nil, fmt.Errorf("capabilities repo not found at %q (set --capabilities-path): %w", c.CapabilitiesPath, err)
	}
	outDir := c.stageDir("capabilities")
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return nil, err
	}
	var out []string
	for _, s := range specs {
		modDir := filepath.Join(c.CapabilitiesPath, s.dir)
		binPath := filepath.Join(outDir, s.binary)
		framework.L.Info().Msgf("Building capability %q (from %s)", s.binary, s.dir)
		args := append([]string{"build"}, s.flags...)
		args = append(args, "-o", binPath, ".")
		cmd := exec.CommandContext(ctx, "go", args...)
		cmd.Dir = modDir
		cmd.Env = append(os.Environ(), append([]string{"CGO_ENABLED=0"}, c.goEnv()...)...)
		cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
		if err := cmd.Run(); err != nil {
			return nil, errors.Wrapf(err, "failed to build capability %q in %s", s.binary, modDir)
		}
		out = append(out, binPath)
	}
	return out, nil
}

// buildNodeImage cross-compiles the chainlink node on the host and bakes it
// into a minimal self-contained image tagged cfg.DevImage.
func (c *LocalBuildConfig) buildNodeImage(ctx context.Context) error {
	// Preflight the cross toolchain so we fail with a clear message.
	for _, f := range []string{c.CC, c.CXX} {
		if _, err := os.Stat(f); err != nil {
			return fmt.Errorf("cross-compiler %q not found; set up the local toolchain or override CRE_LOCAL_CC/CXX (see docs): %w", f, err)
		}
	}
	staticStd := filepath.Join(filepath.Dir(c.LibDir), "lib", "libstdc++.a")
	staticSup := filepath.Join(filepath.Dir(c.LibDir), "lib", "libsupc++.a")

	ctxDir := c.stageDir("node-image")
	if err := os.MkdirAll(filepath.Join(ctxDir, "caps"), 0o755); err != nil {
		return err
	}
	binPath := filepath.Join(ctxDir, "chainlink")

	framework.L.Info().Msgf("Cross-compiling chainlink node from %s (%s)", c.RepoRoot, c.Platform)
	build := exec.CommandContext(ctx, "go", "build", "-o", binPath, ".")
	build.Dir = c.RepoRoot
	build.Env = append(os.Environ(),
		"CGO_ENABLED=1",
		"CC="+c.CC,
		"CXX="+c.CXX,
		"LIBRARY_PATH="+c.LibDir,
		"CGO_LDFLAGS=-L"+c.LibDir+" "+staticStd+" "+staticSup,
	)
	build.Env = append(build.Env, c.goEnv()...)
	build.Stdout, build.Stderr = os.Stdout, os.Stderr
	if err := build.Run(); err != nil {
		return errors.Wrap(err, "node cross-compile failed")
	}

	dockerfile := `FROM ubuntu:24.04
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update \
 && apt-get install -y --no-install-recommends ca-certificates curl libstdc++6 \
 && rm -rf /var/lib/apt/lists/*
RUN useradd --uid 14933 --create-home chainlink
COPY chainlink /usr/local/bin/chainlink
RUN chmod 0755 /usr/local/bin/chainlink
USER chainlink
WORKDIR /home/chainlink
EXPOSE 6688
ENTRYPOINT ["chainlink"]
CMD ["local", "node"]
`
	if err := os.WriteFile(filepath.Join(ctxDir, "Dockerfile"), []byte(dockerfile), 0o600); err != nil {
		return err
	}

	framework.L.Info().Msgf("Building node image %s", c.DevImage)
	db := exec.CommandContext(ctx, "docker", "build", "--platform", c.Platform, "-t", c.DevImage, ctxDir) //nolint:gosec // G204: local docker build with developer-controlled args
	db.Stdout, db.Stderr = os.Stdout, os.Stderr
	if err := db.Run(); err != nil {
		return errors.Wrap(err, "docker build of node image failed")
	}
	return nil
}

func (c *LocalBuildConfig) stageDir(sub string) string {
	base := c.OutputDir
	if base == "" {
		base = filepath.Join(os.TempDir(), "cre-local-build")
	}
	return filepath.Join(base, sub)
}

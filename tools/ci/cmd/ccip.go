package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink/v2/tools/ci/internal/ccip"
	"github.com/smartcontractkit/chainlink/v2/tools/ci/internal/ghaction"
)

const (
	defaultCCIPImageRepo = "public.ecr.aws/chainlink/ccip"
	defaultMaxFallback   = 12
	defaultProbeAttempts = 3
	defaultProbeDelay    = 2 * time.Second
)

func newCcipCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ccip",
		Short: "Commands for CCIP release operations",
	}
	cmd.AddCommand(newCcipResolveBaselineCmd())
	return cmd
}

func newCcipResolveBaselineCmd() *cobra.Command {
	var (
		chainlinkVersion string
		repo             string
		registryURL      string
		maxFallback      int
		probeAttempts    int
		probeRetryDelay  time.Duration
		jsonOutput       bool
	)

	cmd := &cobra.Command{
		Use:   "resolve-baseline",
		Short: "Resolve the CCIP release baseline image tag for mixed-version/rollout tests",
		Long: `Resolve the CCIP release baseline image tag for mixed-version/rollout tests.

Derives the CCIP image tag under test from CHAINLINK_VERSION and probes the
container registry for the first published baseline rc.0 image. Outputs
ccip_pr_tag, baseline_image_tag, and skip (to $GITHUB_OUTPUT when set, else
stdout). Replaces .github/scripts/resolve-ccip-release-baseline.sh.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if chainlinkVersion == "" {
				chainlinkVersion = os.Getenv("CHAINLINK_VERSION")
			}
			if repo == "" {
				repo = os.Getenv("CCIP_IMAGE_REPO")
			}
			if repo == "" {
				repo = defaultCCIPImageRepo
			}

			// Split the full image ref "<host>/<path>" into a registry base
			// URL and a repository path. registryURL overrides the derived
			// base URL (used to point at a mirror or an httptest fake).
			repoPath := repo
			baseURL := ""
			if idx := strings.Index(repo, "/"); idx >= 0 {
				baseURL = "https://" + repo[:idx]
				repoPath = repo[idx+1:]
			}
			if registryURL != "" {
				baseURL = registryURL
			}

			prober := ccip.NewRegistryProbe(baseURL, repoPath, &http.Client{Timeout: 30 * time.Second})
			res, err := ccip.Resolve(cmd.Context(), ccip.ResolveOptions{
				ChainlinkVersion: chainlinkVersion,
				MaxFallback:      maxFallback,
				ProbeAttempts:    probeAttempts,
				ProbeRetryDelay:  probeRetryDelay,
			}, prober)
			if err != nil {
				return err
			}

			out := cmd.OutOrStdout()

			if jsonOutput {
				payload := map[string]any{
					"ccip_pr_tag":        res.PRTag,
					"baseline_image_tag": res.BaselineTag,
					"skip":               res.Skip,
					"candidates":         res.Candidates,
				}
				enc := json.NewEncoder(out)
				enc.SetIndent("", "  ")
				return enc.Encode(payload)
			}

			// SetOutput writes to $GITHUB_OUTPUT when set, else falls back to
			// "k=v" lines on stdout — mirroring the bash resolver's behavior.
			act := ghaction.New(out, "", "")
			if err := act.SetOutput("ccip_pr_tag", res.PRTag); err != nil {
				return err
			}
			if err := act.SetOutput("baseline_image_tag", res.BaselineTag); err != nil {
				return err
			}
			if err := act.SetOutput("skip", strconv.FormatBool(res.Skip)); err != nil {
				return err
			}

			if res.Skip {
				act.Warningf("No published baseline CCIP rc.0 image found for %s (tried: %s). Mixed-version test will be skipped.",
					chainlinkVersion, strings.Join(res.Candidates, " "))
			} else {
				fmt.Fprintf(out, "Baseline for %s: %s:%s\n", chainlinkVersion, repo, res.BaselineTag)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&chainlinkVersion, "chainlink-version", "", "Chainlink release version, e.g. v2.63.1-rc.4 (env: CHAINLINK_VERSION)")
	cmd.Flags().StringVar(&repo, "repo", "", "Full CCIP image ref without scheme, e.g. public.ecr.aws/chainlink/ccip (env: CCIP_IMAGE_REPO)")
	cmd.Flags().StringVar(&registryURL, "registry-url", "", "Override the registry base URL (for mirrors/testing)")
	cmd.Flags().IntVar(&maxFallback, "max-fallback", defaultMaxFallback, "Maximum prior minor .0 rc.0 candidates to probe")
	cmd.Flags().IntVar(&probeAttempts, "probe-attempts", defaultProbeAttempts, "Per-candidate probe attempts on transient errors")
	cmd.Flags().DurationVar(&probeRetryDelay, "probe-retry-delay", defaultProbeDelay, "Delay between probe retries")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output result in JSON format")

	return cmd
}

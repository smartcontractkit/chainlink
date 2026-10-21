// Package ccip contains domain logic for CCIP release CI operations.
//
// The baseline resolver replicates the semantics of the former
// .github/scripts/resolve-ccip-release-baseline.sh: it derives the CCIP image
// tag under test from a Chainlink release version, builds an ordered list of
// candidate baseline rc.0 tags, and probes a container registry for the first
// published one. The mixed-version / rollout test upgrades nodes FROM that
// baseline image TO the image under test.
package ccip

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// prTagRegex captures the core semver and the trailing pre-release identifier
// of a stripped version, e.g. "2.63.1-rc.4" -> core "2.63.1", pre "rc.4".
var prTagRegex = regexp.MustCompile(`^([0-9]+\.[0-9]+\.[0-9]+)-(.*)$`)

// Prober reports whether an image tag is published in a container registry.
//
// IsPublished returns (true, nil) if the tag exists, (false, nil) if it is
// definitively absent (e.g. HTTP 404), and (_, err) for transient failures
// (network errors, 5xx, auth errors) which the caller may retry.
type Prober interface {
	IsPublished(ctx context.Context, tag string) (bool, error)
}

// ResolveOptions configures baseline resolution.
type ResolveOptions struct {
	ChainlinkVersion string
	MaxFallback      int           // prior minor .0 rc.0s to walk (default 12)
	ProbeAttempts    int           // per-candidate retries on transient errors
	ProbeRetryDelay  time.Duration // delay between retries
}

// ResolveResult is the outcome of baseline resolution.
type ResolveResult struct {
	PRTag       string   // derived CCIP image tag under test
	BaselineTag string   // resolved baseline tag (empty if Skip)
	Skip        bool     // true if no published baseline was found
	Candidates  []string // ordered candidate tags that were probed
}

// DerivePRTag replicates build-publish's "Compute CCIP image tag" step: strip
// the leading "v", then insert "-ccip-" before the pre-release identifier.
// Stable versions (no pre-release) are returned unchanged.
//
//	v2.63.1-rc.4  -> 2.63.1-ccip-rc.4
//	v2.63.0       -> 2.63.0
//	v2.63.0-beta.0 -> 2.63.0-ccip-beta.0
func DerivePRTag(chainlinkVersion string) (string, error) {
	if err := validateVersion(chainlinkVersion); err != nil {
		return "", err
	}
	ver := strings.TrimPrefix(chainlinkVersion, "v")
	return prTagRegex.ReplaceAllString(ver, "${1}-ccip-${2}"), nil
}

// Candidates builds the ordered list of baseline ccip tags to probe.
//
// The baseline is always an rc.0, chosen product-agnostically:
//
//	rcN (N>0)          -> first try rc0 of the exact same version, then fall
//	                     back to the base minor .0 rc0s.
//	rc0 / stable / beta -> base minor .0 rc0s, where the base minor is this
//	                       minor's .0 (patch>0) or the previous minor's .0
//	                       (patch==0).
//
// The base minor .0 rc0 is followed by prior minor .0 rc0s (e.g.
// 2.62.0-ccip-rc.0 -> 2.61.0-ccip-rc.0 -> ...) until MAX_FALLBACK minors are
// tried or minor drops below 0.
func Candidates(chainlinkVersion string, maxFallback int) ([]string, error) {
	if err := validateVersion(chainlinkVersion); err != nil {
		return nil, err
	}
	core, pre := splitCore(strings.TrimPrefix(chainlinkVersion, "v"))
	major, minor, patch, err := parseCore(core)
	if err != nil {
		return nil, err
	}

	// rc number N (0 for stable / beta / rc.0).
	n := 0
	if rest, ok := strings.CutPrefix(pre, "rc."); ok {
		nn, err := strconv.Atoi(rest)
		if err != nil {
			return nil, fmt.Errorf("could not parse rc number: %s", pre)
		}
		n = nn
	}

	var candidates []string
	if n > 0 {
		// rcN: first try rc0 of the exact same version.
		candidates = append(candidates, core+"-ccip-rc.0")
	}

	// Anchor minor for the rc0/stable/beta rule and the fallback walk.
	baseMinor := minor
	if patch == 0 {
		baseMinor = minor - 1
	}

	for i := 0; i <= maxFallback; i++ {
		m := baseMinor - i
		if m < 0 {
			break
		}
		candidates = append(candidates, fmt.Sprintf("%d.%d.0-ccip-rc.0", major, m))
	}
	return candidates, nil
}

// Resolve derives the PR tag, builds candidates, and probes them in order; the
// first published candidate wins. A transient probe error is retried up to
// ProbeAttempts times; a definitively-absent tag (404) moves on immediately.
// If no candidate is published, Skip is true and BaselineTag is empty.
func Resolve(ctx context.Context, opts ResolveOptions, prober Prober) (ResolveResult, error) {
	if opts.ProbeAttempts < 1 {
		opts.ProbeAttempts = 1
	}
	if opts.MaxFallback < 0 {
		opts.MaxFallback = 0
	}

	prTag, err := DerivePRTag(opts.ChainlinkVersion)
	if err != nil {
		return ResolveResult{}, err
	}
	candidates, err := Candidates(opts.ChainlinkVersion, opts.MaxFallback)
	if err != nil {
		return ResolveResult{}, err
	}

	res := ResolveResult{PRTag: prTag, Candidates: candidates, Skip: true}
	for _, c := range candidates {
		published, perr := probeWithRetry(ctx, prober, c, opts.ProbeAttempts, opts.ProbeRetryDelay)
		if perr != nil && ctx.Err() != nil {
			return res, ctx.Err()
		}
		if published {
			res.BaselineTag = c
			res.Skip = false
			break
		}
	}
	return res, nil
}

// probeWithRetry probes a tag, retrying transient errors up to attempts times.
// A definitive absent (false, nil) is returned immediately without retry.
func probeWithRetry(ctx context.Context, prober Prober, tag string, attempts int, delay time.Duration) (bool, error) {
	var lastErr error
	for attempt := range attempts {
		published, err := prober.IsPublished(ctx, tag)
		if err == nil {
			return published, nil
		}
		lastErr = err
		if attempt < attempts-1 {
			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return false, ctx.Err()
			}
		}
	}
	return false, lastErr
}

// RegistryProbe probes a Docker Registry v2-compatible registry (e.g. AWS
// public ECR) over HTTP, removing the need for a local docker daemon. It
// authenticates with an anonymous bearer token fetched from {baseURL}/token.
//
// baseURL is the registry API root (e.g. "https://public.ecr.aws"); repo is the
// repository path (e.g. "chainlink/ccip").
type RegistryProbe struct {
	baseURL string
	repo    string
	client  *http.Client
}

// NewRegistryProbe constructs a prober for the given registry root and repo.
func NewRegistryProbe(baseURL, repo string, client *http.Client) *RegistryProbe {
	if client == nil {
		client = &http.Client{Timeout: 30 * time.Second}
	}
	return &RegistryProbe{baseURL: strings.TrimRight(baseURL, "/"), repo: repo, client: client}
}

// IsPublished fetches an anonymous token then requests the manifest. A 200 means
// published, a 404 means definitively absent, and anything else is a transient
// error the caller may retry.
func (p *RegistryProbe) IsPublished(ctx context.Context, tag string) (bool, error) {
	token, err := p.fetchToken(ctx)
	if err != nil {
		return false, err
	}
	ref := fmt.Sprintf("%s/v2/%s/manifests/%s", p.baseURL, p.repo, tag)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ref, nil)
	if err != nil {
		return false, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.docker.distribution.manifest.v2+json, application/vnd.oci.image.manifest.v1+json")

	resp, err := p.client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)

	switch resp.StatusCode {
	case http.StatusOK:
		return true, nil
	case http.StatusNotFound:
		return false, nil
	default:
		return false, fmt.Errorf("registry returned status %d for %s:%s", resp.StatusCode, p.repo, tag)
	}
}

func (p *RegistryProbe) fetchToken(ctx context.Context) (string, error) {
	u := fmt.Sprintf("%s/token?service=public.ecr.aws&scope=repository:%s:pull", p.baseURL, url.QueryEscape(p.repo))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return "", err
	}
	resp, err := p.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("token endpoint returned status %d", resp.StatusCode)
	}
	var body struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return "", err
	}
	if body.Token == "" {
		return "", errors.New("token endpoint returned an empty token")
	}
	return body.Token, nil
}

// validateVersion checks the input is a non-empty "v"-prefixed version.
func validateVersion(v string) error {
	if v == "" {
		return errors.New("CHAINLINK_VERSION must be set")
	}
	if !strings.HasPrefix(v, "v") {
		return fmt.Errorf("CHAINLINK_VERSION must begin with 'v', got: %s", v)
	}
	return nil
}

// splitCore splits a stripped version into its core and pre-release parts.
// "2.63.1-rc.4" -> ("2.63.1", "rc.4"); "2.63.0" -> ("2.63.0", "").
func splitCore(ver string) (string, string) {
	core, pre, _ := strings.Cut(ver, "-")
	return core, pre
}

// parseCore parses a "major.minor.patch" core into its integer components.
func parseCore(core string) (major, minor, patch int, err error) {
	parts := strings.Split(core, ".")
	if len(parts) != 3 {
		return 0, 0, 0, fmt.Errorf("could not parse version core: %s", core)
	}
	major, err = strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, 0, fmt.Errorf("could not parse version core: %s", core)
	}
	minor, err = strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, 0, fmt.Errorf("could not parse version core: %s", core)
	}
	patch, err = strconv.Atoi(parts[2])
	if err != nil {
		return 0, 0, 0, fmt.Errorf("could not parse version core: %s", core)
	}
	return major, minor, patch, nil
}

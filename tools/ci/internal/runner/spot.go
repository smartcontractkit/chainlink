package runner

import (
	"fmt"
	"strings"
)

// SpotStrategy represents a RunsOn spot allocation strategy or on-demand setting.
type SpotStrategy string

const (
	// SpotDisabled represents on-demand EC2 instances (spot=false).
	SpotDisabled SpotStrategy = "false"
	// SpotCapacityOptimized prioritizes pools with most available capacity (spot=co).
	SpotCapacityOptimized SpotStrategy = "co"
	// SpotPriceCapacityOptimized balances price and capacity (spot=pco).
	SpotPriceCapacityOptimized SpotStrategy = "pco"
	// SpotLowestPrice selects the cheapest available pool (spot=lowest-price / spot=lp).
	SpotLowestPrice SpotStrategy = "lowest-price"
)

// SpotInput contains the environment and context factors for determining the spot setting.
type SpotInput struct {
	EventName        string       // GITHUB_EVENT_NAME (e.g., "pull_request", "push", "merge_group", "release", "schedule", "workflow_dispatch")
	Ref              string       // GITHUB_REF (e.g., "refs/heads/develop", "refs/heads/release/2.57.1", "refs/tags/v2.57.1")
	RefType          string       // GITHUB_REF_TYPE ("branch" or "tag")
	RefName          string       // GITHUB_REF_NAME (e.g., "develop", "release/2.57.1", "v2.57.1")
	BaseRef          string       // GITHUB_BASE_REF (target branch for PR, e.g., "develop", "release/2.57.1")
	HeadRef          string       // GITHUB_HEAD_REF (source branch for PR)
	ForceOnDemand    bool         // Explicit flag to disable spot
	StrategyOverride SpotStrategy // Explicit strategy override
	DefaultStrategy  SpotStrategy // Default strategy if spot is enabled (defaults to pco)
}

// SpotResult is the evaluation result for runs-on spot configuration.
type SpotResult struct {
	Spot         string       `json:"spot"`           // "false", "co", "pco", "lowest-price"
	SpotFlag     string       `json:"spot_flag"`      // "spot=false", "spot=co", "spot=pco", "spot=lowest-price"
	Enabled      bool         `json:"enabled"`        // true if spot is enabled, false if on-demand
	Strategy     SpotStrategy `json:"strategy"`       // SpotDisabled, SpotCapacityOptimized, SpotPriceCapacityOptimized, SpotLowestPrice
	Reason       string       `json:"reason"`         // Human-readable explanation of why this setting was chosen
	IsRelease    bool         `json:"is_release"`     // true if evaluated as release branch/tag
	IsMergeQueue bool         `json:"is_merge_queue"` // true if evaluated as merge queue event/branch
}

func isMergeQueue(eventName, ref, refName string) bool {
	if eventName == "merge_group" {
		return true
	}
	cleanRef := strings.TrimPrefix(ref, "refs/heads/")
	if strings.HasPrefix(cleanRef, "gh-readonly-queue/") || strings.HasPrefix(refName, "gh-readonly-queue/") {
		return true
	}
	return false
}

func isReleaseBranch(name string) bool {
	clean := strings.ToLower(strings.TrimPrefix(name, "refs/heads/"))
	if clean == "" {
		return false
	}
	if clean == "release" || clean == "releases" {
		return true
	}
	if strings.HasPrefix(clean, "release/") ||
		strings.HasPrefix(clean, "releases/") ||
		strings.HasPrefix(clean, "release-") ||
		strings.HasPrefix(clean, "hotfix/") ||
		strings.HasPrefix(clean, "hotfix-") {
		return true
	}
	return false
}

func isReleaseOrTag(eventName, ref, refType, refName, baseRef string) bool {
	if eventName == "release" {
		return true
	}
	if refType == "tag" || strings.HasPrefix(ref, "refs/tags/") {
		return true
	}
	if isReleaseBranch(ref) || isReleaseBranch(refName) || isReleaseBranch(baseRef) {
		return true
	}
	return false
}

func isDefaultBranch(ref, refName string) bool {
	clean := strings.ToLower(strings.TrimPrefix(ref, "refs/heads/"))
	cleanName := strings.ToLower(refName)
	return clean == "develop" || clean == "main" || clean == "master" ||
		cleanName == "develop" || cleanName == "main" || cleanName == "master"
}

func normalizeStrategy(s SpotStrategy) (SpotStrategy, string, error) {
	lower := strings.ToLower(strings.TrimSpace(string(s)))
	switch lower {
	case "false", "off", "0", "no", "on-demand", "ondemand", "none":
		return SpotDisabled, "false", nil
	case "co", "capacity-optimized", "capacity_optimized":
		return SpotCapacityOptimized, "co", nil
	case "pco", "price-capacity-optimized", "price_capacity_optimized", "true", "on", "1", "yes":
		return SpotPriceCapacityOptimized, "pco", nil
	case "lp", "lowest-price", "lowest_price":
		return SpotLowestPrice, "lowest-price", nil
	default:
		return SpotStrategy(lower), lower, fmt.Errorf("invalid spot strategy %q (expected 'false', 'co', 'pco', or 'lowest-price')", s)
	}
}

// ResolveSpot evaluates the input context and returns the spot allocation result.
func ResolveSpot(input SpotInput) (SpotResult, error) {
	// 1. Check explicit override: Force on-demand
	if input.ForceOnDemand {
		return SpotResult{
			Spot:         "false",
			SpotFlag:     "spot=false",
			Enabled:      false,
			Strategy:     SpotDisabled,
			Reason:       "forced on-demand via override",
			IsRelease:    isReleaseOrTag(input.EventName, input.Ref, input.RefType, input.RefName, input.BaseRef),
			IsMergeQueue: isMergeQueue(input.EventName, input.Ref, input.RefName),
		}, nil
	}

	// 2. Check explicit strategy override
	if input.StrategyOverride != "" {
		strategyOverride, spotVal, err := normalizeStrategy(input.StrategyOverride)
		if err != nil {
			return SpotResult{}, fmt.Errorf("invalid spot strategy override %q: %w", input.StrategyOverride, err)
		}
		enabled := strategyOverride != SpotDisabled
		flag := "spot=" + spotVal
		return SpotResult{
			Spot:         spotVal,
			SpotFlag:     flag,
			Enabled:      enabled,
			Strategy:     strategyOverride,
			Reason:       "explicit strategy override",
			IsRelease:    isReleaseOrTag(input.EventName, input.Ref, input.RefType, input.RefName, input.BaseRef),
			IsMergeQueue: isMergeQueue(input.EventName, input.Ref, input.RefName),
		}, nil
	}

	// 3. Merge Queue check
	if isMergeQueue(input.EventName, input.Ref, input.RefName) {
		return SpotResult{
			Spot:         "false",
			SpotFlag:     "spot=false",
			Enabled:      false,
			Strategy:     SpotDisabled,
			Reason:       "merge queue runs require on-demand to prevent queue eviction",
			IsMergeQueue: true,
		}, nil
	}

	// 4. Release branch or tag check
	if isReleaseOrTag(input.EventName, input.Ref, input.RefType, input.RefName, input.BaseRef) {
		return SpotResult{
			Spot:      "false",
			SpotFlag:  "spot=false",
			Enabled:   false,
			Strategy:  SpotDisabled,
			Reason:    "release branch/tag requires on-demand for reliable builds",
			IsRelease: true,
		}, nil
	}

	// 5. Default branch push (develop / main / master)
	if input.EventName == "push" && isDefaultBranch(input.Ref, input.RefName) {
		return SpotResult{
			Spot:     "co",
			SpotFlag: "spot=co",
			Enabled:  true,
			Strategy: SpotCapacityOptimized,
			Reason:   "default branch push uses capacity-optimized spot for mainline stability",
		}, nil
	}

	// 6. Resolve strategy for PRs, scheduled, and other events
	defaultStrategy := input.DefaultStrategy
	if defaultStrategy == "" {
		defaultStrategy = SpotPriceCapacityOptimized
	}
	resolvedStrategy, spotVal, err := normalizeStrategy(defaultStrategy)
	if err != nil {
		return SpotResult{}, fmt.Errorf("invalid default spot strategy %q: %w", input.DefaultStrategy, err)
	}
	enabled := resolvedStrategy != SpotDisabled
	flag := "spot=" + spotVal

	reason := "standard CI run uses price-capacity-optimized spot for cost efficiency"
	switch {
	case input.EventName == "schedule":
		reason = "scheduled workflow uses price-capacity-optimized spot"
	case input.EventName == "pull_request" || input.EventName == "pull_request_target":
		reason = "pull request uses price-capacity-optimized spot"
	case defaultStrategy != SpotPriceCapacityOptimized:
		reason = fmt.Sprintf("uses custom default strategy (%s)", spotVal)
	}

	return SpotResult{
		Spot:     spotVal,
		SpotFlag: flag,
		Enabled:  enabled,
		Strategy: resolvedStrategy,
		Reason:   reason,
	}, nil
}

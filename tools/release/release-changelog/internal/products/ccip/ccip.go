// Package ccip defines the CCIP product for the release changelog engine:
// which repos and pins matter for a CCIP release audit, and how to
// normalize CCIP release ref shapes (image tags, image URIs) into git refs.
//
// EDIT the Product definition below to add/remove repositories or tune which
// paths are tracked. Everything in this file is configuration: no CLI flags
// or Slack inputs are needed to change tracking behavior.
package ccip

import (
	"regexp"
	"strings"

	"github.com/smartcontractkit/chainlink/v2/tools/release/release-changelog/internal/engine"
)

// Product is the CCIP product definition consumed by the engine. To track an
// additional repo, add an entry to Repos. To narrow a repo's changelog to
// specific paths (like the core repo entry), set IncludePaths/ExcludePaths.
var Product = engine.Product{
	DisplayName:  "CCIP",
	NormalizeRef: normalizeRef,
	Repos: []engine.RepoConfig{
		{
			Name:  "chainlink-ccip",
			Owner: "smartcontractkit",
			GoModules: []string{
				"github.com/smartcontractkit/chainlink-ccip", // primary
				"github.com/smartcontractkit/chainlink-ccip/chains/evm",
				"github.com/smartcontractkit/chainlink-ccip/chains/solana",
				"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings",
			},
		},
		{
			Name:       "chainlink-aptos",
			Owner:      "smartcontractkit",
			GoModules:  []string{"github.com/smartcontractkit/chainlink-aptos/codec"},
			PluginKeys: []string{"aptos"}, // primary pin
		},
		{
			Name:       "chainlink-sui",
			Owner:      "smartcontractkit",
			GoModules:  []string{"github.com/smartcontractkit/chainlink-sui/codec"},
			PluginKeys: []string{"sui"}, // primary pin
		},
		{
			Name:       "chainlink-solana",
			Owner:      "smartcontractkit",
			PluginKeys: []string{"solana"}, // primary pin (not in root go.mod)
		},
		{
			Name:       "chainlink-ton",
			Owner:      "smartcontractkit",
			GoModules:  []string{"github.com/smartcontractkit/chainlink-ton"},
			PluginKeys: []string{"ton"}, // dual-source: checked against go.mod
		},
		{
			Name:  "chainlink-evm",
			Owner: "smartcontractkit",
			GoModules: []string{
				"github.com/smartcontractkit/chainlink-evm",
				"github.com/smartcontractkit/chainlink-evm/gethwrappers",
				// contracts/cre/gobindings intentionally not tracked (CRE-owned)
			},
			PluginKeys: []string{"evm"}, // dual-source: checked against go.mod
		},
		{
			// The core repo itself, restricted to the CCIP capability tree.
			Name:         "chainlink",
			Owner:        "smartcontractkit",
			IncludePaths: []string{"core/capabilities/ccip/"},
			ExcludePaths: []string{},
			Local:        true,
		},
	},
}

// ---------------------------------------------------------------------------
// CCIP ref normalization (assigned to Product.NormalizeRef above)
//
// Anything a release engineer might paste in — image tags, image URIs — is
// normalized to a git ref here, before resolution. The git machinery
// (engine/refs.go) never sees product-specific shapes.
// ---------------------------------------------------------------------------

// ccipImageTag matches CCIP release image tags. build-publish.yml derives
// the image tag from the pushed git tag by inserting "-ccip" before the
// prerelease identifier (v2.34.0-rc.1 -> 2.34.0-ccip-rc.1), so the mapping is
// invertible without contacting the registry. An optional arch suffix
// (-amd64/-arm64, as shown in the ECR console) is ignored.
var ccipImageTag = regexp.MustCompile(`^(\d+\.\d+\.\d+)-ccip-(.+?)(?:-(?:amd64|arm64))?$`)

// ccipGitTagWithCCIP matches the common mistake of adding a "v" prefix to an
// image tag (or inserting "-ccip-" into a git tag). Git tags are vX.Y.Z-rc.N
// (no "-ccip-"); image tags are X.Y.Z-ccip-rc.N (no "v"). Strip "-ccip-" and
// keep the "v" to recover the real git tag.
var ccipGitTagWithCCIP = regexp.MustCompile(`^v(\d+\.\d+\.\d+)-ccip-(.+?)(?:-(?:amd64|arm64))?$`)

// normalizeRef maps a CCIP image tag or image URI to the git tag that built
// the image, so release engineers can pass the version they actually work
// with (e.g. "2.56.1-ccip-rc.2" or
// "public.ecr.aws/chainlink/ccip:2.56.1-ccip-rc.2"). It also handles the
// hybrid "v2.56.1-ccip-rc.2" (image tag with a "v" prefix). Anything else is
// returned unchanged.
func normalizeRef(ref string) string {
	// Accept a full image URI: git refnames never contain ":".
	if i := strings.LastIndex(ref, ":"); i >= 0 {
		ref = ref[i+1:]
	}
	if m := ccipGitTagWithCCIP.FindStringSubmatch(ref); m != nil {
		return "v" + m[1] + "-" + m[2]
	}
	if m := ccipImageTag.FindStringSubmatch(ref); m != nil {
		return "v" + m[1] + "-" + m[2]
	}
	return ref
}

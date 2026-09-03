#!/usr/bin/env bash
# Resolve the CCIP release baseline image tag for mixed-version / rollout tests.
#
# The "Post Build, Sign and Publish Chainlink" workflow is triggered (by
# build-publish.yml) with the release git tag in CHAINLINK_VERSION (e.g.
# "v2.63.1-rc.4") but it never receives the CCIP image tag directly. The CCIP
# image lives at public.ecr.aws/chainlink/ccip and is tagged by stripping the
# leading "v" and inserting "-ccip-" before the pre-release identifier, e.g.
#   v2.63.1-rc.4 -> 2.63.1-ccip-rc.4     (rc)
#   v2.63.0       -> 2.63.0              (stable, no -ccip suffix)
# This replicates build-publish.yml's "Compute CCIP image tag" step.
#
# The rollout test upgrades nodes FROM a baseline rc.0 image TO the image under
# test. The baseline is always an rc.0, chosen product-agnostically (we don't
# know which patch belongs to CCIP vs other products), per the rule:
#   rcN (N>0)            -> <same version>-ccip-rc.0
#   rc0 / stable / beta  -> <base version>-ccip-rc.0, where the base version is
#                           this minor's .0 (patch>0) or the previous minor's .0
#                           (patch==0).
# If the preferred rc.0 image is not published, fall back to prior minor .0 rc.0s
# (e.g. 2.62.0-ccip-rc.0 -> 2.61.0-ccip-rc.0 -> ...) until one is found. If
# none is found within MAX_FALLBACK minors, skip=true.
#
# Each candidate is probed against CCIP_IMAGE_REPO (default
# public.ecr.aws/chainlink/ccip) via `docker manifest inspect`; the first
# published image wins.
#
# Outputs (to $GITHUB_OUTPUT if set, else stdout):
#   ccip_pr_tag        - derived CCIP image tag under test
#   baseline_image_tag - resolved baseline tag (empty if skip)
#   skip               - "true" if no published baseline was found
set -euo pipefail

CHAINLINK_VERSION="${CHAINLINK_VERSION:-}"
CCIP_IMAGE_REPO="${CCIP_IMAGE_REPO:-public.ecr.aws/chainlink/ccip}"
MAX_FALLBACK="${MAX_FALLBACK:-12}"
PROBE_ATTEMPTS="${PROBE_ATTEMPTS:-3}"
PROBE_RETRY_DELAY="${PROBE_RETRY_DELAY:-2}"

error() {
  echo "::error::$1" >&2
  exit 1
}

# Probe whether an image tag is published. `docker manifest inspect` against
# public ECR can transiently fail for existing tags, so retry a few times
# before treating a tag as absent. PROBE_ATTEMPTS/PROBE_RETRY_DELAY are
# overridable (tests use PROBE_ATTEMPTS=1 to skip retries).
probe() { # <tag> -> 0 if published, 1 otherwise
  local ref="${CCIP_IMAGE_REPO}:$1"
  local i
  for ((i = 0; i < PROBE_ATTEMPTS; i++)); do
    if docker manifest inspect "$ref" >/dev/null 2>&1; then
      return 0
    fi
    if ((i < PROBE_ATTEMPTS - 1)); then
      sleep "$PROBE_RETRY_DELAY"
    fi
  done
  return 1
}

out() {
  if [[ -n "${GITHUB_OUTPUT:-}" ]]; then
    echo "$1" >>"${GITHUB_OUTPUT}"
  else
    echo "$1"
  fi
}

[[ -n "$CHAINLINK_VERSION" ]] || error "CHAINLINK_VERSION must be set"
[[ "$CHAINLINK_VERSION" == v* ]] || error "CHAINLINK_VERSION must begin with 'v', got: ${CHAINLINK_VERSION}"

ver="${CHAINLINK_VERSION#v}" # 2.63.1-rc.4

# Derive the CCIP "PR" image tag (replicates build-publish compute-ccip-tag).
ccip_pr_tag=$(printf '%s' "$ver" | sed -E 's/^([0-9]+\.[0-9]+\.[0-9]+)-(.*)$/\1-ccip-\2/')
out "ccip_pr_tag=${ccip_pr_tag}"

# Split core version from pre-release identifier.
core="${ver%%-*}" # 2.63.1
pre="${ver#"${core}"}"
pre="${pre#-}" # rc.4 / beta.0 / "" (stable)

IFS=. read -r major minor patch <<<"$core"
[[ "$major" =~ ^[0-9]+$ && "$minor" =~ ^[0-9]+$ && "$patch" =~ ^[0-9]+$ ]] \
  || error "Could not parse version core: ${core}"

# rc number N (0 for stable / beta / rc.0).
n=0
if [[ "$pre" == rc.* ]]; then
  n="${pre#rc.}"
  [[ "$n" =~ ^[0-9]+$ ]] || error "Could not parse rc number: ${pre}"
fi

# Build the ordered list of candidate baseline ccip tags.
candidates=()
if (( n > 0 )); then
  # rcN: first try rc0 of the exact same version.
  candidates+=("${core}-ccip-rc.0")
fi

# Anchor minor for the rc0/stable/beta rule and the fallback walk.
if (( patch > 0 )); then
  base_minor=$((minor)) # this minor's .0
else
  base_minor=$((minor - 1)) # previous minor's .0
fi

# Add the base-version rc0, then walk prior minor .0 rc0s as fallback.
for ((i = 0; i <= MAX_FALLBACK; i++)); do
  m=$((base_minor - i))
  ((m < 0)) && break
  candidates+=("${major}.${m}.0-ccip-rc.0")
done

# Probe candidates in order; first published image wins.
skip=true
baseline_image_tag=""
for c in "${candidates[@]}"; do
  if probe "$c"; then
    baseline_image_tag="$c"
    skip=false
    break
  fi
done

out "baseline_image_tag=${baseline_image_tag}"
out "skip=${skip}"

if [[ "$skip" == "true" ]]; then
  echo "::warning::No published baseline CCIP rc.0 image found for ${CHAINLINK_VERSION} (tried: ${candidates[*]}). Mixed-version test will be skipped."
else
  echo "Baseline for ${CHAINLINK_VERSION}: ${CCIP_IMAGE_REPO}:${baseline_image_tag}"
fi

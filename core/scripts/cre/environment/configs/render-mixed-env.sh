#!/usr/bin/env bash
#
# Renders the mixed-env topology config from its template, substituting the two
# node images (PR-built and develop/baseline). The rendered file is gitignored.
#
# Usage (from core/scripts/cre/environment):
#
#   CRE_PR_IMAGE=<pr chainlink image ref> \
#   CRE_BASELINE_IMAGE=<develop chainlink image ref> \
#     ./configs/render-mixed-env.sh
#
# Then start the environment against the rendered config, WITHOUT setting
# CTF_CHAINLINK_IMAGE (a non-empty value would force every node onto one image):
#
#   CTF_CONFIGS=configs/mixed-env-don.toml go run . env start
#
set -euo pipefail

: "${CRE_PR_IMAGE:?set CRE_PR_IMAGE to the PR-built chainlink node image ref}"
: "${CRE_BASELINE_IMAGE:?set CRE_BASELINE_IMAGE to the develop/baseline chainlink node image ref}"

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
template="${script_dir}/mixed-env-don.toml.tmpl"
output="${script_dir}/mixed-env-don.toml"

if ! command -v envsubst >/dev/null 2>&1; then
  echo "error: envsubst not found (install gettext)" >&2
  exit 1
fi

# Restrict substitution to only these two vars so no other '$' in the file is touched.
envsubst '${CRE_PR_IMAGE} ${CRE_BASELINE_IMAGE}' < "${template}" > "${output}"

echo "Rendered ${output}"
echo "  CRE_PR_IMAGE=${CRE_PR_IMAGE}"
echo "  CRE_BASELINE_IMAGE=${CRE_BASELINE_IMAGE}"

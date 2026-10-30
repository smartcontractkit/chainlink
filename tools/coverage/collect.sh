#!/usr/bin/env bash
# Collect coverage artifacts for a commit and publish a PR report.
#
# Required environment:
#   GITHUB_REPOSITORY  owner/repo
#   GITHUB_TOKEN       token with actions:read, pull-requests:write
#   TARGET_SHA         commit whose workflow runs should be collected
#
# Optional:
#   PR_NUMBER          PR number; inferred from TARGET_SHA when empty
#   OUTPUT_DIR         report output directory (default: coverage-report)
set -euo pipefail

REPO="${GITHUB_REPOSITORY:?GITHUB_REPOSITORY is required}"
TARGET_SHA="${TARGET_SHA:?TARGET_SHA is required}"
OUT="${OUTPUT_DIR:-coverage-report}"
ARTIFACTS="${OUT}/artifacts"
mkdir -p "$ARTIFACTS"

if [[ -z "${PR_NUMBER:-}" ]]; then
  PR_NUMBER="$(gh api "repos/${REPO}/commits/${TARGET_SHA}/pulls" --jq '.[0].number // empty' 2>/dev/null || true)"
fi
BASE_SHA=""
if [[ -n "${PR_NUMBER:-}" ]]; then
  BASE_SHA="$(gh api "repos/${REPO}/pulls/${PR_NUMBER}" --jq '.base.sha' 2>/dev/null || true)"
fi
PATCH_BASE="${PATCH_BASE:-}"
PATCH_HEAD="${PATCH_HEAD:-$TARGET_SHA}"

if [[ -n "$BASE_SHA" ]]; then
  git fetch --no-tags origin "$BASE_SHA" "$TARGET_SHA" >/dev/null 2>&1 || true
fi
if [[ -n "$PATCH_BASE" && -n "$PATCH_HEAD" ]]; then
  git fetch --no-tags origin "$PATCH_BASE" "$PATCH_HEAD" >/dev/null 2>&1 || true
fi
if [[ -z "$PATCH_BASE" && -n "$BASE_SHA" ]]; then
  PATCH_BASE="$(git merge-base "$BASE_SHA" "$TARGET_SHA" 2>/dev/null || echo "$BASE_SHA")"
fi

echo "Collecting coverage for ${REPO}@${TARGET_SHA} (PR ${PR_NUMBER:-none}, base ${BASE_SHA:-none})"

runs="$(gh run list --repo "$REPO" --commit "$TARGET_SHA" -L 100 --json databaseId --jq '.[].databaseId')"
for run in $runs; do
  names="$(gh api "repos/${REPO}/actions/runs/${run}/artifacts?per_page=100" --jq '.artifacts[].name' 2>/dev/null | grep -E '^(coverage-(core|deployment|inmemory|tools|e2e)-|test-results-)' || true)"
  while IFS= read -r name; do
    [[ -z "$name" ]] && continue
    dest="${ARTIFACTS}/${run}/${name}"
    mkdir -p "$dest"
    echo "Downloading ${name} from run ${run}"
    gh run download "$run" --repo "$REPO" -n "$name" -D "$dest" || true
  done <<< "$names"
done

all_dirs=""
unit_dirs=""
integration_dirs=""
e2e_dirs=""

while IFS= read -r meta; do
  [[ -z "$meta" ]] && continue
  dir="$(dirname "$meta")"
  all_dirs="${all_dirs}${all_dirs:+,}${dir}"
  case "$meta" in
    */coverage-core-go_core_tests-*/*|*/coverage-core-go_core_race_tests-*/*|*/coverage-core-scripts-*/*|*/coverage-tools-*/*)
      unit_dirs="${unit_dirs}${unit_dirs:+,}${dir}" ;;
    */coverage-core-go_core_tests_integration-*/*|*/coverage-deployment-*/*|*/coverage-inmemory-*/*)
      integration_dirs="${integration_dirs}${integration_dirs:+,}${dir}" ;;
    */coverage-e2e-*/*)
      e2e_dirs="${e2e_dirs}${e2e_dirs:+,}${dir}" ;;
  esac
done < <(find "$ARTIFACTS" -type f -name 'covmeta.*' 2>/dev/null | sort -u)

merge_and_report() {
  local category="$1" dirs="$2"
  [[ -z "$dirs" ]] && return 0
  echo "Merging ${category}: ${dirs}"
  go run ./tools/coverage merge --covdata "$dirs" --out "${OUT}/${category}.txt"
  local args=(--profile "${OUT}/${category}.txt" --category "$category"
    --json "${OUT}/${category}.json" --markdown "${OUT}/${category}.md")
  if [[ -n "$PATCH_BASE" && -n "$PATCH_HEAD" ]]; then
    args+=(--base "$PATCH_BASE" --head "$PATCH_HEAD")
  fi
  go run ./tools/coverage report "${args[@]}"
}

merge_and_report unit "$unit_dirs"
merge_and_report integration "$integration_dirs"
merge_and_report e2e "$e2e_dirs"

profile_list=""
for category in unit integration e2e; do
  if [[ -f "${OUT}/${category}.txt" ]]; then
    profile_list="${profile_list}${profile_list:+,}${OUT}/${category}.txt"
  fi
done
if [[ -n "$profile_list" ]]; then
  go run ./tools/coverage merge-profiles --profile "$profile_list" --out "${OUT}/any.txt"
  args=(--profile "${OUT}/any.txt" --category any
    --json "${OUT}/any.json" --markdown "${OUT}/any.md")
  if [[ -n "$PATCH_BASE" && -n "$PATCH_HEAD" ]]; then
    args+=(--base "$PATCH_BASE" --head "$PATCH_HEAD")
  fi
  go run ./tools/coverage report "${args[@]}"
fi

test_xmls="$(find "$ARTIFACTS" -type f -name '*.xml' -path '*test-results-*' 2>/dev/null | sort -u)"
if [[ -n "$test_xmls" ]]; then
  python3 - "$OUT/test-results.json" $test_xmls <<'PY'
import json, sys, xml.etree.ElementTree as ET
files = sys.argv[2:]
totals = {"tests": 0, "failures": 0, "errors": 0, "skipped": 0, "files": len(files)}
for path in files:
    try:
        root = ET.parse(path).getroot()
    except Exception:
        continue
    suites = root.findall("testsuite") if root.tag == "testsuites" else [root]
    for suite in suites:
        for key in ("tests", "failures", "errors", "skipped"):
            try:
                totals[key] += int(suite.get(key, "0"))
            except ValueError:
                pass
with open(sys.argv[1], "w", encoding="utf-8") as f:
    json.dump(totals, f)
PY
fi

{
  echo "<!-- coverage-report -->"
  echo "# Coverage report for \`${TARGET_SHA:0:9}\`"
  echo
  echo "Commit: \`${TARGET_SHA}\`"
  if [[ -n "${PR_NUMBER:-}" ]]; then echo "PR: #${PR_NUMBER}"; fi
  if [[ -n "$BASE_SHA" ]]; then echo "Base: \`${BASE_SHA}\`"; fi
  echo
  echo "Targets: project coverage should not drop more than 1%; patch target is 80%. This report is informational."
  echo
  if [[ -f "${OUT}/test-results.json" ]]; then
    echo "### Test results"
    echo
    echo '```json'
    cat "${OUT}/test-results.json"
    echo
    echo '```'
    echo
  fi
  for category in any unit integration e2e; do
    if [[ -f "${OUT}/${category}.md" ]]; then
      cat "${OUT}/${category}.md"
      echo
    fi
  done
} > "${OUT}/summary.md"

cat "${OUT}/summary.md"

if [[ -n "${PR_NUMBER:-}" ]]; then
  marker="<!-- coverage-report -->"
  comment_id="$(gh api --paginate "repos/${REPO}/issues/${PR_NUMBER}/comments" 2>/dev/null \
    | jq -r --arg marker "$marker" '.[] | select(.body | contains($marker)) | .id' 2>/dev/null \
    | head -1 || true)"
  if [[ -n "$comment_id" ]]; then
    gh api -X PATCH "repos/${REPO}/issues/comments/${comment_id}" -F body="@${OUT}/summary.md"
  else
    gh pr comment "$PR_NUMBER" --repo "$REPO" --body-file "${OUT}/summary.md"
  fi
fi

if [[ -n "${GITHUB_STEP_SUMMARY:-}" ]]; then
  cat "${OUT}/summary.md" >> "$GITHUB_STEP_SUMMARY"
fi

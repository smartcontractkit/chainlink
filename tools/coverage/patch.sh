#!/usr/bin/env bash
# Fast patch-only coverage for the packages touched by a PR.
#
# It intentionally does NOT run the full suite or E2E. It runs the affected
# packages with -short, uses the normal Go build/module caches, computes patch
# coverage with tools/coverage, and updates a sticky PR comment.
#
# Required:
#   BASE_SHA     PR base SHA (the merge base is derived with git merge-base)
#   HEAD_SHA     PR head SHA
# Optional:
#   OUT          output directory (default: coverage-patch)
#   PR_NUMBER    PR number; when set, the report is commented on the PR
#   GH_TOKEN     token with pull-requests:write (defaults to gh auth)
#   SHORT        set to false to run the affected packages without -short
#   TEST_TIMEOUT per-package test timeout (default: 10m)
set -euo pipefail

BASE_SHA="${BASE_SHA:?BASE_SHA is required}"
HEAD_SHA="${HEAD_SHA:?HEAD_SHA is required}"
REPO_ROOT="$(git rev-parse --show-toplevel)"
cd "$REPO_ROOT"

OUT="${OUT:-coverage-patch}"
case "$OUT" in
  /*) ;;
  *) OUT="$REPO_ROOT/$OUT" ;;
esac
SHORT="${SHORT:-true}"
TEST_TIMEOUT="${TEST_TIMEOUT:-10m}"
COVDATA_ROOT="$OUT/covdata"
rm -rf "$OUT"
mkdir -p "$OUT" "$COVDATA_ROOT"

CHANGED="$OUT/changed-files.txt"
git diff --name-only --diff-filter=ACMR "$BASE_SHA...$HEAD_SHA" -- '*.go' 2>/dev/null \
  | grep -vE '(^|/)(mocks|testdata|testutils|testhelpers|cltest|fixtures)/|_test\.go$|\.pb\.go$|_gen\.go$|_codecgen\.go$|^vendor/' \
  > "$CHANGED" || true

if [[ ! -s "$CHANGED" ]]; then
  echo "No changed production Go files."
  exit 0
fi

# Map every changed file to the nearest Go module and package directory.
ROWS="$OUT/rows.tsv"
: > "$ROWS"
while IFS= read -r file; do
  dir="$(dirname "$file")"
  while [[ "$dir" != "." && ! -f "$dir/go.mod" ]]; do
    dir="$(dirname "$dir")"
  done
  if [[ -f "$dir/go.mod" ]]; then
    module_dir="$REPO_ROOT/$dir"
  else
    module_dir="$REPO_ROOT"
  fi
  printf '%s\t%s\n' "$module_dir" "$REPO_ROOT/$(dirname "$file")" >> "$ROWS"
done < "$CHANGED"

modules="$(cut -f1 "$ROWS" | sort -u)"
tests_failed=0
ran_any=0

while IFS= read -r module_dir; do
  [[ -z "$module_dir" ]] && continue
  map_file="$OUT/packages-$(echo "$module_dir" | tr '/' '_').tsv"
  if ! (cd "$module_dir" && go list -e -f '{{.Dir}}	{{.ImportPath}}' ./...) > "$map_file" 2> "$OUT/go-list.log"; then
    echo "warning: could not list packages in $module_dir; skipping" >&2
    continue
  fi

  packages=()
  while IFS= read -r package_dir; do
    [[ -z "$package_dir" ]] && continue
    import_path="$(awk -F '\t' -v d="$package_dir" '$1 == d { print $2; exit }' "$map_file")"
    [[ -z "$import_path" ]] && continue
    packages+=("$import_path")
  done < <(awk -F '\t' -v m="$module_dir" '$1 == m { print $2 }' "$ROWS" | sort -u)

  # Deduplicate package import paths while preserving a stable order.
  package_list="$(printf '%s\n' "${packages[@]}" | sort -u | paste -sd, -)"
  [[ -z "$package_list" ]] && continue

  covdir="$COVDATA_ROOT/$(echo "$module_dir" | tr '/' '_')"
  mkdir -p "$covdir"
  short_flag=()
  [[ "$SHORT" == "true" ]] && short_flag=(-short)

  echo "Testing affected packages in ${module_dir#"$REPO_ROOT"/}:"
  echo "  $package_list"
  if ! (cd "$module_dir" && go test \
      -count=1 \
      -timeout "$TEST_TIMEOUT" \
      "${short_flag[@]}" \
      -cover \
      -covermode=atomic \
      -coverpkg="$package_list" \
      "${packages[@]}" \
      -args -test.gocoverdir="$covdir") >> "$OUT/test.log" 2>&1; then
    tests_failed=1
    echo "warning: some affected-package tests failed in ${module_dir#"$REPO_ROOT"/}; see $OUT/test.log" >&2
  fi
  ran_any=1
done <<< "$modules"

if [[ "$ran_any" == "0" ]]; then
  echo "No runnable affected packages."
  exit 0
fi

covdirs="$(find "$COVDATA_ROOT" -type f -name 'covmeta.*' -exec dirname {} \; 2>/dev/null | sort -u | paste -sd, -)"
if [[ -z "$covdirs" ]]; then
  echo "No coverage data was produced for the affected packages."
  exit 0
fi

go run ./tools/coverage merge --covdata "$covdirs" --out "$OUT/patch.txt"
merge_base="$(git merge-base "$BASE_SHA" "$HEAD_SHA")"
go run ./tools/coverage report \
  --profile "$OUT/patch.txt" \
  --base "$merge_base" \
  --head "$HEAD_SHA" \
  --category "patch (affected packages)" \
  --json "$OUT/patch.json" \
  --markdown "$OUT/patch.md"

not_measured="$(jq -r '.Patch.Files[] | select(.TestFile == false and .Statements == 0) | .Path' "$OUT/patch.json" 2>/dev/null || true)"
if [[ -n "$not_measured" ]]; then
  {
    echo
    echo "### Not measured"
    echo
    echo "These changed files had no tests in their own package, so the quick run could not measure them:"
    echo
    while IFS= read -r path; do
      [[ -n "$path" ]] && echo "- \`$path\`"
    done <<< "$not_measured"
    echo
    echo "The nightly full report covers these."
  } >> "$OUT/patch.md"
fi

{
  echo "<!-- coverage-patch -->"
  echo
  echo "### Quick patch coverage"
  echo
  if [[ "$SHORT" == "true" ]]; then
    echo "Affected-package coverage for \`$HEAD_SHA\`, using \`-short\`. "
  else
    echo "Affected-package coverage for \`$HEAD_SHA\`. "
  fi
  echo "This is a fast PR signal, not the global coverage number; the nightly full report is authoritative."
  echo
  if [[ "$tests_failed" == "1" ]]; then
    echo "> Some affected-package tests failed. The coverage below may be incomplete; check the job log."
    echo
  fi
  cat "$OUT/patch.md"
} > "$OUT/comment.md"

cat "$OUT/comment.md"

if [[ -n "${PR_NUMBER:-}" && -n "${GITHUB_REPOSITORY:-}" ]]; then
  marker="<!-- coverage-patch -->"
  comment_id="$(gh api --paginate "repos/${GITHUB_REPOSITORY}/issues/${PR_NUMBER}/comments" 2>/dev/null \
    | jq -r --arg marker "$marker" '.[] | select(.user.login == "github-actions[bot]") | select(.body | contains($marker)) | .id' 2>/dev/null \
    | head -1 || true)"
  if [[ -n "$comment_id" ]]; then
    gh api -X PATCH "repos/${GITHUB_REPOSITORY}/issues/comments/${comment_id}" -F body="@$OUT/comment.md"
  else
    gh pr comment "$PR_NUMBER" --repo "$GITHUB_REPOSITORY" --body-file "$OUT/comment.md"
  fi
fi

if [[ "$tests_failed" == "1" ]]; then
  exit 1
fi

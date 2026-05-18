#!/bin/bash
set -e

# Dependencies:
# gh cli ^2.15.0 https://github.com/cli/cli/releases/tag/v2.15.0
# jq ^1.6 https://stedolan.github.io/jq/

repo=smartcontractkit/operator-ui
gitRoot="$(dirname -- "$0")/../"
cd "$gitRoot/operator_ui"

tag_file=TAG
current_tag=$(cat $tag_file)
echo "Currently pinned tag for $repo is $current_tag"

echo "Getting latest release for tag for $repo"
release=$(gh release view -R $repo --json 'tagName,body')
latest_tag=$(echo "$release" | jq -r '.tagName')
body=$(echo "$release" | jq -r '.body')

if [ "$current_tag" = "$latest_tag" ]; then
  echo "Tag $current_tag is up to date."
  exit 0
else
  echo "Tag $current_tag is out of date, updating $tag_file file to latest version..."
  echo "$latest_tag" >"$tag_file"
  echo "Tag updated $current_tag -> $latest_tag"
  if [ "$CI" ]; then
    echo "current_tag=$current_tag" >>$GITHUB_OUTPUT
    echo "latest_tag=$latest_tag" >>$GITHUB_OUTPUT

    # See https://github.com/orgs/community/discussions/26288#discussioncomment-3876281
    delimiter="$(openssl rand -hex 8)"
    echo "body<<${delimiter}" >>"${GITHUB_OUTPUT}"
    echo "$body" >>"${GITHUB_OUTPUT}"
    echo "${delimiter}" >>"${GITHUB_OUTPUT}"
  fi
fi

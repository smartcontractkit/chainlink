#!/bin/bash
set -e

# Dependencies:
# gh cli ^2.15.0 https://github.com/cli/cli/releases/tag/v2.15.0
# jq ^1.6 https://stedolan.github.io/jq/

repo=smartcontractkit/operator-ui
gitRoot=$(git rev-parse --show-toplevel)
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
    echo "::set-output name=current_tag::$current_tag"
    echo "::set-output name=latest_tag::$latest_tag"
    # See https://github.com/peter-evans/create-pull-request/blob/main/docs/examples.md#setting-the-pull-request-body-from-a-file
    body="${body//'%'/'%25'}"
    body="${body//$'\n'/'%0A'}"
    body="${body//$'\r'/'%0D'}"
    echo "::set-output name=body::$body"
  fi
fi

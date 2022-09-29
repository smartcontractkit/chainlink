#!/bin/bash
set -e

# Dependencies:
# gh cli ^2.15.0 https://github.com/cli/cli/releases/tag/v2.15.0
# jq ^1.6 https://stedolan.github.io/jq/

repo=smartcontractkit/operator-ui
gitRoot=$(git rev-parse --show-toplevel)
cd "$gitRoot/operator_ui"
unpack_dir="$gitRoot/core/web/assets"
tag=$(cat TAG)

echo "Getting release $tag for $repo"
release=$(gh release view "$tag" -R $repo --json 'assets')
asset_name=$(echo "$release" | jq -r '.assets | map(select(.contentType == "application/x-gtar"))[0].name')

echo "Downloading ${repo}:${tag} asset: $asset_name..."
echo ""
gh release download "$tag" -R "$repo" -p "$asset_name"

echo "Unpacking asset $asset_name"
tar -xvzf "$asset_name"

echo ""
echo "Removing old contents of $unpack_dir"
rm -rf "$unpack_dir"
echo "Copying contents of package/artifacts to $unpack_dir"
cp -rf package/artifacts/. "$unpack_dir" || true

echo "Cleaning up"
rm -r package
rm "$asset_name"

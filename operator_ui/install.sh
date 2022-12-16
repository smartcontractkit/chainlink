#!/usr/bin/env bash
set -e

owner=smartcontractkit
repo=operator-ui
fullRepo=${owner}/${repo}
gitRoot=$(git rev-parse --show-toplevel || pwd)
cd "$gitRoot/operator_ui"
unpack_dir="$gitRoot/core/web/assets"
tag=$(cat TAG)
# Remove the version prefix "v"
strippedTag="${tag:1}"
# Taken from https://github.com/kennyp/asdf-golang/blob/master/lib/helpers.sh
msg() {
  echo -e "\033[32m$1\033[39m" >&2
}

err() {
  echo -e "\033[31m$1\033[39m" >&2
}

fail() {
  err "$1"
  exit 1
}

msg "Getting release $tag for $fullRepo"
# https://docs.github.com/en/rest/releases/releases#get-a-release-by-tag-name
asset_name=${owner}-${repo}-${strippedTag}.tgz
download_url=https://github.com/${fullRepo}/releases/download/${tag}/${asset_name}

# Inspired from https://github.com/kennyp/asdf-golang/blob/master/bin/download#L29
msg "Download URL: ${download_url}"
# Check if we're able to download first
http_code=$(curl -LIs -w '%{http_code}' -o /dev/null "$download_url")
if [ "$http_code" -eq 404 ] || [ "$http_code" -eq 403 ]; then
  fail "URL: ${download_url} returned status ${http_code}"
fi
# Then go ahead if we get a success code
msg "Downloading ${fullRepo}:${tag} asset: $asset_name..."
msg ""
curl -L -o "$asset_name" "$download_url"

msg "Unpacking asset $asset_name"
tar -xvzf "$asset_name"

msg ""
msg "Removing old contents of $unpack_dir"
rm -rf "$unpack_dir"
msg "Copying contents of package/artifacts to $unpack_dir"
cp -rf package/artifacts/. "$unpack_dir" || true

msg "Cleaning up"
rm -r package
rm "$asset_name"
